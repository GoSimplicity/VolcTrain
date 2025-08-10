package volcano

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apiResource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vcqueue "volcano.sh/apis/pkg/apis/scheduling/v1beta1"
)

// QueueSpec Volcano队列规格定义
type QueueSpec struct {
	Name               string
	Weight             int32
	Capability         map[string]interface{}
	GuaranteedResource corev1.ResourceList
	MaxResource        corev1.ResourceList
	DeservedResource   corev1.ResourceList
	ShareWeight        *float64
	Reclaimable        *bool
	HierarchyEnabled   bool
	ParentQueue        string
	State              vcqueue.QueueState
	Labels             map[string]string
	Annotations        map[string]string
}

// QueueStatus Volcano队列状态
type QueueStatus struct {
	Name               string
	State              vcqueue.QueueState
	Weight             int32
	Share              float64
	GuaranteedResource corev1.ResourceList
	AllocatedResource  corev1.ResourceList
	PendingResource    corev1.ResourceList
	UsedResource       corev1.ResourceList
	Capability         map[string]interface{}
	Inqueue            int32
	Pending            int32
	Running            int32
	Unknown            int32
}

// CreateQueue 创建Volcano队列
func (c *Client) CreateQueue(spec *QueueSpec) (*vcqueue.Queue, error) {
	// 构建Volcano Queue对象
	queue := &vcqueue.Queue{
		ObjectMeta: metav1.ObjectMeta{
			Name:        spec.Name,
			Labels:      spec.Labels,
			Annotations: spec.Annotations,
		},
		Spec: vcqueue.QueueSpec{
			Weight: spec.Weight,
			// Capability:  spec.Capability, // 需要类型转换
			Reclaimable: spec.Reclaimable,
		},
	}

	// 设置资源配额 - v1beta1版本中通过Capability设置
	if spec.GuaranteedResource != nil {
		// 将GuaranteedResource合并到Capability中 - v1beta1中需要类型转换
		for k, v := range spec.GuaranteedResource {
			if queue.Spec.Capability == nil {
				queue.Spec.Capability = make(corev1.ResourceList)
			}
			queue.Spec.Capability[corev1.ResourceName(k)] = v
		}
	}

	// 设置最大资源
	if spec.MaxResource != nil {
		queue.Spec.Capability["nvidia.com/gpu"] = spec.MaxResource["nvidia.com/gpu"]
		queue.Spec.Capability["cpu"] = spec.MaxResource["cpu"]
		queue.Spec.Capability["memory"] = spec.MaxResource["memory"]
	}

	// 设置期望资源 - v1beta1版本中通过Capability和annotation设置
	if spec.DeservedResource != nil {
		if queue.Annotations == nil {
			queue.Annotations = make(map[string]string)
		}
		// 通过annotation存储期望资源信息
		for k, v := range spec.DeservedResource {
			queue.Annotations["scheduling.volcano.sh/deserved-"+k.String()] = v.String()
		}
	}

	// 设置层级队列
	if spec.HierarchyEnabled && spec.ParentQueue != "" {
		if queue.Annotations == nil {
			queue.Annotations = make(map[string]string)
		}
		queue.Annotations["scheduling.volcano.sh/parent-queue"] = spec.ParentQueue
	}

	// 设置队列状态 - v1beta1版本中通过annotation设置
	if queue.Annotations == nil {
		queue.Annotations = make(map[string]string)
	}

	queueState := string(spec.State)
	if queueState == "" {
		queueState = "Open"
	}
	queue.Annotations["scheduling.volcano.sh/queue-state"] = queueState

	// 创建队列
	createdQueue, err := c.volcanoClient.SchedulingV1beta1().Queues().Create(
		context.TODO(),
		queue,
		metav1.CreateOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("创建Volcano队列失败: %v", err)
	}

	return createdQueue, nil
}

// GetQueueStatus 获取队列状态
func (c *Client) GetQueueStatus(queueName string) (*QueueStatus, error) {
	queue, err := c.volcanoClient.SchedulingV1beta1().Queues().Get(
		context.TODO(),
		queueName,
		metav1.GetOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("获取Volcano队列失败: %v", err)
	}

	status := &QueueStatus{
		Name:               queue.Name,
		State:              "Open", // 简化处理
		Weight:             queue.Spec.Weight,
		GuaranteedResource: getGuaranteedResource(queue),
		AllocatedResource:  queue.Status.Allocated,
		PendingResource:    getResourceListFromInt(queue.Status.Pending),
		UsedResource:       corev1.ResourceList{}, // v1beta1版本中没有Used字段
		Capability:         c.resourceListToMapInterface(queue.Spec.Capability),
		Inqueue:            queue.Status.Inqueue,
		Pending:            queue.Status.Pending,
		Running:            queue.Status.Running,
		Unknown:            queue.Status.Unknown,
	}

	// 计算资源共享比例
	if totalWeight := c.getTotalQueueWeight(); totalWeight > 0 {
		status.Share = float64(queue.Spec.Weight) / float64(totalWeight)
	}

	return status, nil
}

// getTotalQueueWeight 获取所有队列的总权重（简化实现）
func (c *Client) getTotalQueueWeight() int32 {
	// TODO: 实现获取所有队列总权重的逻辑
	return 100 // 临时返回值
}

// UpdateQueue 更新队列
func (c *Client) UpdateQueue(queueName string, spec *QueueSpec) (*vcqueue.Queue, error) {
	// 获取现有队列
	queue, err := c.volcanoClient.SchedulingV1beta1().Queues().Get(
		context.TODO(),
		queueName,
		metav1.GetOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("获取队列失败: %v", err)
	}

	// 更新队列规格
	if spec.Weight > 0 {
		queue.Spec.Weight = spec.Weight
	}

	if spec.Capability != nil {
		// v1beta1中需要类型转换
		for k, v := range spec.Capability {
			if queue.Spec.Capability == nil {
				queue.Spec.Capability = make(corev1.ResourceList)
			}
			if quantity, err := apiResource.ParseQuantity(fmt.Sprintf("%v", v)); err == nil {
				queue.Spec.Capability[corev1.ResourceName(k)] = quantity
			}
		}
	}

	if spec.GuaranteedResource != nil {
		// v1beta1版本中通过annotations存储保证资源
		if queue.Annotations == nil {
			queue.Annotations = make(map[string]string)
		}
		for k, v := range spec.GuaranteedResource {
			queue.Annotations["scheduling.volcano.sh/guaranteed-"+k.String()] = v.String()
		}
	}

	if spec.Reclaimable != nil {
		queue.Spec.Reclaimable = spec.Reclaimable
	}

	if spec.State != "" {
		// v1beta1版本中通过annotations存储队列状态
		if queue.Annotations == nil {
			queue.Annotations = make(map[string]string)
		}
		queue.Annotations["scheduling.volcano.sh/queue-state"] = string(spec.State)
	}

	// 更新标签和注解
	if spec.Labels != nil {
		if queue.Labels == nil {
			queue.Labels = make(map[string]string)
		}
		for k, v := range spec.Labels {
			queue.Labels[k] = v
		}
	}

	if spec.Annotations != nil {
		if queue.Annotations == nil {
			queue.Annotations = make(map[string]string)
		}
		for k, v := range spec.Annotations {
			queue.Annotations[k] = v
		}
	}

	// 更新队列
	updatedQueue, err := c.volcanoClient.SchedulingV1beta1().Queues().Update(
		context.TODO(),
		queue,
		metav1.UpdateOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("更新Volcano队列失败: %v", err)
	}

	return updatedQueue, nil
}

// DeleteQueue 删除队列
func (c *Client) DeleteQueue(queueName string) error {
	err := c.volcanoClient.SchedulingV1beta1().Queues().Delete(
		context.TODO(),
		queueName,
		metav1.DeleteOptions{},
	)
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("删除Volcano队列失败: %v", err)
	}

	return nil
}

// ListQueues 列出队列
func (c *Client) ListQueues(labelSelector string) (*vcqueue.QueueList, error) {
	queues, err := c.volcanoClient.SchedulingV1beta1().Queues().List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: labelSelector,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("列出Volcano队列失败: %v", err)
	}

	return queues, nil
}

// OpenQueue 开放队列
func (c *Client) OpenQueue(queueName string) error {
	return c.updateQueueState(queueName, vcqueue.QueueStateOpen)
}

// CloseQueue 关闭队列
func (c *Client) CloseQueue(queueName string) error {
	return c.updateQueueState(queueName, vcqueue.QueueStateClosed)
}

// updateQueueState 更新队列状态
func (c *Client) updateQueueState(queueName string, state vcqueue.QueueState) error {
	// 获取现有队列
	queue, err := c.volcanoClient.SchedulingV1beta1().Queues().Get(
		context.TODO(),
		queueName,
		metav1.GetOptions{},
	)
	if err != nil {
		return fmt.Errorf("获取队列失败: %v", err)
	}

	// 更新状态 - v1beta1版本中通过annotations存储
	if queue.Annotations == nil {
		queue.Annotations = make(map[string]string)
	}
	queue.Annotations["scheduling.volcano.sh/queue-state"] = string(state)

	// 更新队列
	_, err = c.volcanoClient.SchedulingV1beta1().Queues().Update(
		context.TODO(),
		queue,
		metav1.UpdateOptions{},
	)
	if err != nil {
		return fmt.Errorf("更新队列状态失败: %v", err)
	}

	return nil
}

// GetQueueResourceUsage 获取队列资源使用情况
func (c *Client) GetQueueResourceUsage(queueName string) (*QueueResourceUsage, error) {
	queue, err := c.volcanoClient.SchedulingV1beta1().Queues().Get(
		context.TODO(),
		queueName,
		metav1.GetOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("获取队列失败: %v", err)
	}

	usage := &QueueResourceUsage{
		QueueName:          queue.Name,
		Allocated:          c.resourceListToMap(queue.Status.Allocated),
		Used:               map[string]string{}, // v1beta1版本中没有Used字段
		Pending:            c.intToResourceMapString(queue.Status.Pending),
		GuaranteedResource: c.resourceListToMap(getGuaranteedResource(queue)),
		Capacity:           c.getQueueCapacity(queue),
		Inqueue:            queue.Status.Inqueue,
		Running:            queue.Status.Running,
		Weight:             queue.Spec.Weight,
	}

	// 计算共享比例
	if totalWeight := c.getTotalQueueWeight(); totalWeight > 0 {
		usage.Share = float64(queue.Spec.Weight) / float64(totalWeight)
	}

	return usage, nil
}

// QueueResourceUsage 队列资源使用情况
type QueueResourceUsage struct {
	QueueName          string            `json:"queueName"`
	Allocated          map[string]string `json:"allocated"`
	Used               map[string]string `json:"used"`
	Pending            map[string]string `json:"pending"`
	Share              float64           `json:"share"`
	Weight             int32             `json:"weight"`
	GuaranteedResource map[string]string `json:"guaranteedResource"`
	Capacity           map[string]string `json:"capacity"`
	Inqueue            int32             `json:"inqueue"`
	Running            int32             `json:"running"`
}

// resourceListToMap 将ResourceList转换为map
func (c *Client) resourceListToMap(rl corev1.ResourceList) map[string]string {
	result := make(map[string]string)
	for name, quantity := range rl {
		result[string(name)] = quantity.String()
	}
	return result
}

// getQueueCapacity 获取队列容量
func (c *Client) getQueueCapacity(queue *vcqueue.Queue) map[string]string {
	capacity := make(map[string]string)

	// 从Capability中获取容量信息
	for resourceName, quantity := range queue.Spec.Capability {
		capacity[string(resourceName)] = quantity.String()
	}

	return capacity
}

// BuildResourceList 构建资源列表
func BuildResourceList(cpu, memory, gpu string) corev1.ResourceList {
	resources := corev1.ResourceList{}

	if cpu != "" {
		if quantity, err := apiResource.ParseQuantity(cpu); err == nil {
			resources[corev1.ResourceCPU] = quantity
		}
	}

	if memory != "" {
		if quantity, err := apiResource.ParseQuantity(memory); err == nil {
			resources[corev1.ResourceMemory] = quantity
		}
	}

	if gpu != "" {
		if quantity, err := apiResource.ParseQuantity(gpu); err == nil {
			resources["nvidia.com/gpu"] = quantity
		}
	}

	return resources
}

// BuildQueueCapability 构建队列能力
func BuildQueueCapability(cpu, memory, gpu string, customCapability map[string]interface{}) map[string]interface{} {
	capability := map[string]interface{}{}

	// 添加自定义能力
	for k, v := range customCapability {
		capability[k] = v
	}

	// 设置资源容量
	if cpu != "" {
		capability["cpu"] = cpu
	}
	if memory != "" {
		capability["memory"] = memory
	}
	if gpu != "" {
		capability["nvidia.com/gpu"] = gpu
	}

	return capability
}

// 辅助函数：获取队列状态
func getQueueState(queue *vcqueue.Queue) string {
	if queue.Annotations != nil {
		if state, ok := queue.Annotations["scheduling.volcano.sh/queue-state"]; ok {
			return state
		}
	}
	return string(vcqueue.QueueStateOpen) // 默认状态
}

// 辅助函数：获取保证资源
func getGuaranteedResource(queue *vcqueue.Queue) corev1.ResourceList {
	// 从Capability中提取或从annotations中获取
	resources := make(corev1.ResourceList)
	if queue.Annotations != nil {
		for k, v := range queue.Annotations {
			if strings.HasPrefix(k, "scheduling.volcano.sh/guaranteed-") {
				resourceName := strings.TrimPrefix(k, "scheduling.volcano.sh/guaranteed-")
				if quantity, err := apiResource.ParseQuantity(v); err == nil {
					resources[corev1.ResourceName(resourceName)] = quantity
				}
			}
		}
	}
	return resources
}

// 辅助函数：将int32转换为ResourceList（用于兼容性）
func getResourceListFromInt(pending int32) corev1.ResourceList {
	// 简化处理，实际应该根据具体情况映射
	return corev1.ResourceList{
		corev1.ResourcePods: *apiResource.NewQuantity(int64(pending), apiResource.DecimalSI),
	}
}

// 辅助函数：将int32转换为资源映射（用于兼容性）
func (c *Client) intToResourceMap(pending int32) map[string]interface{} {
	return map[string]interface{}{
		"pods": pending,
	}
}

// 辅助函数：将int32转换为资源映射字符串（用于兼容性）
func (c *Client) intToResourceMapString(pending int32) map[string]string {
	return map[string]string{
		"pods": fmt.Sprintf("%d", pending),
	}
}

// resourceListToMapInterface 将ResourceList转换为map[string]interface{}
func (c *Client) resourceListToMapInterface(rl corev1.ResourceList) map[string]interface{} {
	result := make(map[string]interface{})
	for name, quantity := range rl {
		result[string(name)] = quantity.String()
	}
	return result
}
