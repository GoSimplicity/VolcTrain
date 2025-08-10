package volcano

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

// ClusterSchedulingInfo 集群调度信息
type ClusterSchedulingInfo struct {
	ClusterName        string                 `json:"clusterName"`
	SchedulerName      string                 `json:"schedulerName"`
	Queues             []QueueSchedulingInfo  `json:"queues"`
	Nodes              []NodeSchedulingInfo   `json:"nodes"`
	TotalResources     map[string]string      `json:"totalResources"`
	AllocatedResources map[string]string      `json:"allocatedResources"`
	AvailableResources map[string]string      `json:"availableResources"`
	PendingJobs        int32                  `json:"pendingJobs"`
	RunningJobs        int32                  `json:"runningJobs"`
	SchedulingActions  []SchedulingActionInfo `json:"schedulingActions"`
}

// QueueSchedulingInfo 队列调度信息
type QueueSchedulingInfo struct {
	Name               string            `json:"name"`
	Weight             int32             `json:"weight"`
	Share              float64           `json:"share"`
	State              string            `json:"state"`
	GuaranteedResource map[string]string `json:"guaranteedResource"`
	AllocatedResource  map[string]string `json:"allocatedResource"`
	PendingJobs        int32             `json:"pendingJobs"`
	RunningJobs        int32             `json:"runningJobs"`
	Capability         map[string]string `json:"capability"`
}

// NodeSchedulingInfo 节点调度信息
type NodeSchedulingInfo struct {
	Name               string            `json:"name"`
	Ready              bool              `json:"ready"`
	Schedulable        bool              `json:"schedulable"`
	TotalResources     map[string]string `json:"totalResources"`
	AllocatedResources map[string]string `json:"allocatedResources"`
	AvailableResources map[string]string `json:"availableResources"`
	Taints             []string          `json:"taints"`
	Labels             map[string]string `json:"labels"`
	Conditions         []string          `json:"conditions"`
}

// SchedulingActionInfo 调度动作信息
type SchedulingActionInfo struct {
	Action        string            `json:"action"`
	JobName       string            `json:"jobName"`
	Queue         string            `json:"queue"`
	TaskName      string            `json:"taskName"`
	NodeName      string            `json:"nodeName"`
	Reason        string            `json:"reason"`
	Timestamp     string            `json:"timestamp"`
	ResourceUsage map[string]string `json:"resourceUsage"`
}

// GetClusterSchedulingInfo 获取集群调度信息
func (c *Client) GetClusterSchedulingInfo(clusterName string) (*ClusterSchedulingInfo, error) {
	// 获取所有队列
	queues, err := c.ListQueues("")
	if err != nil {
		return nil, fmt.Errorf("获取队列列表失败: %v", err)
	}

	// 获取所有节点
	nodes, err := c.kubeClient.CoreV1().Nodes().List(
		context.TODO(),
		metav1.ListOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("获取节点列表失败: %v", err)
	}

	// 构建集群调度信息
	info := &ClusterSchedulingInfo{
		ClusterName:   clusterName,
		SchedulerName: "volcano",
		Queues:        make([]QueueSchedulingInfo, 0),
		Nodes:         make([]NodeSchedulingInfo, 0),
	}

	// 处理队列信息
	totalWeight := int32(0)
	for _, queue := range queues.Items {
		totalWeight += queue.Spec.Weight

		queueInfo := QueueSchedulingInfo{
			Name:               queue.Name,
			Weight:             queue.Spec.Weight,
			State:              getQueueState(&queue),
			GuaranteedResource: c.resourceListToMap(getGuaranteedResource(&queue)),
			AllocatedResource:  c.resourceListToMap(queue.Status.Allocated),
			PendingJobs:        queue.Status.Pending,
			RunningJobs:        queue.Status.Running,
			Capability:         c.resourceListToStringMap(queue.Spec.Capability),
		}
		info.Queues = append(info.Queues, queueInfo)
	}

	// 计算队列共享比例
	for i := range info.Queues {
		if totalWeight > 0 {
			info.Queues[i].Share = float64(info.Queues[i].Weight) / float64(totalWeight)
		}
	}

	// 处理节点信息
	totalResources := make(map[string]int64)
	allocatedResources := make(map[string]int64)

	for _, node := range nodes.Items {
		nodeInfo := c.buildNodeSchedulingInfo(&node)
		info.Nodes = append(info.Nodes, nodeInfo)

		// 累计资源
		for resource, quantity := range node.Status.Capacity {
			if total, ok := totalResources[string(resource)]; ok {
				totalResources[string(resource)] = total + quantity.Value()
			} else {
				totalResources[string(resource)] = quantity.Value()
			}
		}

		for resource, quantity := range node.Status.Allocatable {
			// 计算已分配资源 = 容量 - 可分配
			capacity := node.Status.Capacity[resource]
			allocated := capacity.DeepCopy()
			allocated.Sub(quantity)

			if total, ok := allocatedResources[string(resource)]; ok {
				allocatedResources[string(resource)] = total + allocated.Value()
			} else {
				allocatedResources[string(resource)] = allocated.Value()
			}
		}
	}

	// 设置集群资源信息
	info.TotalResources = c.int64MapToStringMap(totalResources)
	info.AllocatedResources = c.int64MapToStringMap(allocatedResources)
	info.AvailableResources = make(map[string]string)

	// 计算可用资源
	for resource, total := range totalResources {
		allocated := allocatedResources[resource]
		available := total - allocated
		if available < 0 {
			available = 0
		}
		info.AvailableResources[resource] = fmt.Sprintf("%d", available)
	}

	// 获取作业统计信息
	info.PendingJobs, info.RunningJobs = c.getJobStatistics()

	// TODO: 获取调度动作信息（需要实现事件监听器）
	info.SchedulingActions = make([]SchedulingActionInfo, 0)

	return info, nil
}

// buildNodeSchedulingInfo 构建节点调度信息
func (c *Client) buildNodeSchedulingInfo(node *corev1.Node) NodeSchedulingInfo {
	nodeInfo := NodeSchedulingInfo{
		Name:               node.Name,
		Ready:              c.isNodeReady(node),
		Schedulable:        !node.Spec.Unschedulable,
		TotalResources:     c.resourceListToMap(node.Status.Capacity),
		AllocatedResources: make(map[string]string),
		AvailableResources: c.resourceListToMap(node.Status.Allocatable),
		Taints:             make([]string, 0),
		Labels:             node.Labels,
		Conditions:         make([]string, 0),
	}

	// 计算已分配资源
	for resource, capacity := range node.Status.Capacity {
		allocatable := node.Status.Allocatable[resource]
		allocated := capacity.DeepCopy()
		allocated.Sub(allocatable)
		nodeInfo.AllocatedResources[string(resource)] = allocated.String()
	}

	// 处理污点
	for _, taint := range node.Spec.Taints {
		taintStr := fmt.Sprintf("%s=%s:%s", taint.Key, taint.Value, taint.Effect)
		nodeInfo.Taints = append(nodeInfo.Taints, taintStr)
	}

	// 处理条件
	for _, condition := range node.Status.Conditions {
		if condition.Status == corev1.ConditionTrue {
			nodeInfo.Conditions = append(nodeInfo.Conditions, string(condition.Type))
		}
	}

	return nodeInfo
}

// isNodeReady 检查节点是否就绪
func (c *Client) isNodeReady(node *corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}

// getJobStatistics 获取作业统计信息
func (c *Client) getJobStatistics() (pendingJobs, runningJobs int32) {
	// 获取所有作业
	jobs, err := c.ListJobs("", "")
	if err != nil {
		return 0, 0
	}

	for _, job := range jobs.Items {
		switch job.Status.State.Phase {
		case "Pending":
			pendingJobs++
		case "Running":
			runningJobs++
		}
	}

	return pendingJobs, runningJobs
}

// capabilityToMap 将能力配置转换为字符串map
func (c *Client) capabilityToMap(capability map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for key, value := range capability {
		result[key] = fmt.Sprintf("%v", value)
	}
	return result
}

// int64MapToStringMap 将int64 map转换为string map
func (c *Client) int64MapToStringMap(m map[string]int64) map[string]string {
	result := make(map[string]string)
	for key, value := range m {
		result[key] = fmt.Sprintf("%d", value)
	}
	return result
}

// GetJobSchedulingHistory 获取作业调度历史
func (c *Client) GetJobSchedulingHistory(namespace, jobName string) ([]JobSchedulingHistory, error) {
	// 获取作业相关的事件
	fieldSelector := fields.AndSelectors(
		fields.OneTermEqualSelector("involvedObject.name", jobName),
		fields.OneTermEqualSelector("involvedObject.kind", "Job"),
	)

	events, err := c.kubeClient.CoreV1().Events(c.getNamespace(namespace)).List(
		context.TODO(),
		metav1.ListOptions{
			FieldSelector: fieldSelector.String(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("获取作业事件失败: %v", err)
	}

	// 转换为调度历史
	history := make([]JobSchedulingHistory, 0)
	for _, event := range events.Items {
		if c.isSchedulingEvent(&event) {
			historyItem := JobSchedulingHistory{
				JobName:   jobName,
				Action:    c.mapEventToAction(event.Reason),
				Reason:    event.Reason,
				Message:   event.Message,
				Timestamp: event.FirstTimestamp.Format("2006-01-02 15:04:05"),
			}

			// 从事件中提取节点信息
			if event.Source.Host != "" {
				historyItem.CurrentNode = event.Source.Host
			}

			history = append(history, historyItem)
		}
	}

	return history, nil
}

// JobSchedulingHistory 作业调度历史
type JobSchedulingHistory struct {
	Id            int64             `json:"id"`
	JobId         int64             `json:"jobId"`
	JobName       string            `json:"jobName"`
	Action        string            `json:"action"` // scheduled, preempted, rescheduled, failed
	Reason        string            `json:"reason"`
	Message       string            `json:"message"`
	PreviousNode  string            `json:"previousNode,omitempty"`
	CurrentNode   string            `json:"currentNode,omitempty"`
	ResourceUsage map[string]string `json:"resourceUsage,omitempty"`
	Queue         string            `json:"queue"`
	Priority      int64             `json:"priority"`
	Timestamp     string            `json:"timestamp"`
	Duration      int64             `json:"duration,omitempty"`
}

// isSchedulingEvent 判断是否为调度相关事件
func (c *Client) isSchedulingEvent(event *corev1.Event) bool {
	schedulingReasons := map[string]bool{
		"Scheduled":             true,
		"FailedScheduling":      true,
		"Preempted":             true,
		"EvictJob":              true,
		"JobUnschedulable":      true,
		"PodGroupUnschedulable": true,
	}

	return schedulingReasons[event.Reason]
}

// mapEventToAction 将事件原因映射为调度动作
func (c *Client) mapEventToAction(reason string) string {
	actionMap := map[string]string{
		"Scheduled":             "scheduled",
		"FailedScheduling":      "failed",
		"Preempted":             "preempted",
		"EvictJob":              "preempted",
		"JobUnschedulable":      "failed",
		"PodGroupUnschedulable": "failed",
	}

	if action, ok := actionMap[reason]; ok {
		return action
	}
	return "unknown"
}

// GetQueueJobs 获取队列中的作业
func (c *Client) GetQueueJobs(queueName string) ([]JobInfo, error) {
	// 使用标签选择器获取队列中的作业
	labelSelector := labels.Set{"volcano.sh/queue-name": queueName}.AsSelector()

	jobs, err := c.ListJobs("", labelSelector.String())
	if err != nil {
		return nil, fmt.Errorf("获取队列作业失败: %v", err)
	}

	jobInfos := make([]JobInfo, 0)
	for _, job := range jobs.Items {
		jobInfo := JobInfo{
			Name:      job.Name,
			Namespace: job.Namespace,
			Queue:     queueName,
			Phase:     string(job.Status.State.Phase),
			State:     string(job.Status.State.Phase), // State和Phase使用相同值
			Tasks:     make(map[string]TaskInfo),
		}

		// 转换任务信息
		for taskName, taskStatus := range job.Status.TaskStatusCount {
			// 获取各状态的Pod数量
			running := int32(0)
			succeeded := int32(0)
			failed := int32(0)
			pending := int32(0)

			// 从taskStatus.Phase中提取Pod计数
			for phase, count := range taskStatus.Phase {
				switch phase {
				case "Running":
					running = count
				case "Succeeded":
					succeeded = count
				case "Failed":
					failed = count
				case "Pending":
					pending = count
				}
			}

			jobInfo.Tasks[taskName] = TaskInfo{
				Name:      taskName,
				Phase:     "Active", // 简化处理
				Replicas:  c.getTaskReplicas(&job, taskName),
				Running:   running,
				Succeeded: succeeded,
				Failed:    failed,
				Pending:   pending,
			}
		}

		jobInfos = append(jobInfos, jobInfo)
	}

	return jobInfos, nil
}

// JobInfo 作业信息
type JobInfo struct {
	Name      string              `json:"name"`
	Namespace string              `json:"namespace"`
	Queue     string              `json:"queue"`
	Phase     string              `json:"phase"`
	State     string              `json:"state"`
	Tasks     map[string]TaskInfo `json:"tasks"`
}

// TaskInfo 任务信息
type TaskInfo struct {
	Name      string `json:"name"`
	Phase     string `json:"phase"`
	Replicas  int32  `json:"replicas"`
	Running   int32  `json:"running"`
	Succeeded int32  `json:"succeeded"`
	Failed    int32  `json:"failed"`
	Pending   int32  `json:"pending"`
}

// MonitorQueueResources 监控队列资源使用
func (c *Client) MonitorQueueResources(queueName string) (*QueueResourceMonitoring, error) {
	// 获取队列状态
	queueStatus, err := c.GetQueueStatus(queueName)
	if err != nil {
		return nil, err
	}

	// 获取队列中的作业
	jobs, err := c.GetQueueJobs(queueName)
	if err != nil {
		return nil, err
	}

	monitoring := &QueueResourceMonitoring{
		QueueName:     queueName,
		QueueStatus:   *queueStatus,
		Jobs:          jobs,
		ResourceUsage: c.calculateResourceUsageRate(queueStatus),
		Alerts:        c.generateResourceAlerts(queueStatus),
	}

	return monitoring, nil
}

// QueueResourceMonitoring 队列资源监控
type QueueResourceMonitoring struct {
	QueueName     string             `json:"queueName"`
	QueueStatus   QueueStatus        `json:"queueStatus"`
	Jobs          []JobInfo          `json:"jobs"`
	ResourceUsage map[string]float64 `json:"resourceUsage"`
	Alerts        []ResourceAlert    `json:"alerts"`
}

// ResourceAlert 资源告警
type ResourceAlert struct {
	Type      string  `json:"type"`     // warning, critical
	Resource  string  `json:"resource"` // cpu, memory, gpu
	Message   string  `json:"message"`
	UsageRate float64 `json:"usageRate"`
	Threshold float64 `json:"threshold"`
	Timestamp string  `json:"timestamp"`
}

// calculateResourceUsageRate 计算资源使用率
func (c *Client) calculateResourceUsageRate(status *QueueStatus) map[string]float64 {
	usage := make(map[string]float64)

	for resourceName, allocated := range status.AllocatedResource {
		if guaranteed, ok := status.GuaranteedResource[resourceName]; ok {
			allocatedVal := c.parseResourceValue(allocated.String())
			guaranteedVal := c.parseResourceValue(guaranteed.String())
			if guaranteedVal > 0 {
				usage[string(resourceName)] = float64(allocatedVal) / float64(guaranteedVal)
			}
		}
	}

	return usage
}

// parseResourceValue 解析资源值
func (c *Client) parseResourceValue(value string) int64 {
	if quantity, err := resource.ParseQuantity(value); err == nil {
		return quantity.Value()
	}
	return 0
}

// generateResourceAlerts 生成资源告警
func (c *Client) generateResourceAlerts(status *QueueStatus) []ResourceAlert {
	alerts := make([]ResourceAlert, 0)

	usage := c.calculateResourceUsageRate(status)
	for resource, rate := range usage {
		if rate > 0.9 { // 使用率超过90%触发告警
			alert := ResourceAlert{
				Type:      "critical",
				Resource:  resource,
				Message:   fmt.Sprintf("队列 %s 的 %s 资源使用率过高", status.Name, resource),
				UsageRate: rate,
				Threshold: 0.9,
				Timestamp: metav1.Now().Format("2006-01-02 15:04:05"),
			}
			alerts = append(alerts, alert)
		} else if rate > 0.8 { // 使用率超过80%触发警告
			alert := ResourceAlert{
				Type:      "warning",
				Resource:  resource,
				Message:   fmt.Sprintf("队列 %s 的 %s 资源使用率较高", status.Name, resource),
				UsageRate: rate,
				Threshold: 0.8,
				Timestamp: metav1.Now().Format("2006-01-02 15:04:05"),
			}
			alerts = append(alerts, alert)
		}
	}

	return alerts
}

// resourceListToStringMap 将ResourceList转换为map[string]string
func (c *Client) resourceListToStringMap(rl corev1.ResourceList) map[string]string {
	result := make(map[string]string)
	for name, quantity := range rl {
		result[string(name)] = quantity.String()
	}
	return result
}
