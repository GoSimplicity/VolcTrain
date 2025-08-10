package volcano

import (
	"context"
	"fmt"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// MonitoringService Volcano监控服务
type MonitoringService struct {
	client       *Client
	gpuManager   *GPUManager
	jobManager   *JobManager
	alertManager *AlertManager

	// 监控配置
	config    *MonitoringConfig
	stopCh    chan struct{}
	mutex     sync.RWMutex
	isRunning bool

	// 缓存的监控数据
	clusterMetrics *ClusterMetrics
	lastUpdate     time.Time
	eventQueue     chan MonitoringEvent
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	UpdateInterval       time.Duration         `json:"updateInterval"`       // 更新间隔
	AlertThresholds      AlertThresholds       `json:"alertThresholds"`      // 告警阈值
	MetricsRetention     time.Duration         `json:"metricsRetention"`     // 指标保留时间
	EnableGPUMonitor     bool                  `json:"enableGPUMonitor"`     // 启用GPU监控
	EnableJobMonitor     bool                  `json:"enableJobMonitor"`     // 启用作业监控
	EnableQueueMonitor   bool                  `json:"enableQueueMonitor"`   // 启用队列监控
	NotificationChannels []NotificationChannel `json:"notificationChannels"` // 通知渠道
}

// AlertThresholds 告警阈值配置
type AlertThresholds struct {
	GPUUtilization        float64 `json:"gpuUtilization"`        // GPU使用率阈值
	MemoryUtilization     float64 `json:"memoryUtilization"`     // 内存使用率阈值
	QueueWaitTime         int64   `json:"queueWaitTime"`         // 队列等待时间阈值(秒)
	JobFailureRate        float64 `json:"jobFailureRate"`        // 作业失败率阈值
	NodeUnhealthyRate     float64 `json:"nodeUnhealthyRate"`     // 节点异常率阈值
	ResourceFragmentation float64 `json:"resourceFragmentation"` // 资源碎片化阈值
}

// NotificationChannel 通知渠道
type NotificationChannel struct {
	Type    string                 `json:"type"` // email, slack, webhook, dingtalk
	Enabled bool                   `json:"enabled"`
	Config  map[string]interface{} `json:"config"`
	Filters []string               `json:"filters"` // 过滤规则
}

// ClusterMetrics 集群指标
type ClusterMetrics struct {
	Timestamp    time.Time         `json:"timestamp"`
	GPUMetrics   GPUClusterMetrics `json:"gpuMetrics"`
	JobMetrics   JobClusterMetrics `json:"jobMetrics"`
	QueueMetrics []QueueMetrics    `json:"queueMetrics"`
	NodeMetrics  []NodeMetrics     `json:"nodeMetrics"`
	AlertSummary AlertSummary      `json:"alertSummary"`
}

// GPUClusterMetrics GPU集群指标
type GPUClusterMetrics struct {
	TotalGPUs       int32            `json:"totalGPUs"`
	AllocatedGPUs   int32            `json:"allocatedGPUs"`
	AvailableGPUs   int32            `json:"availableGPUs"`
	UtilizationRate float64          `json:"utilizationRate"`
	AllocationRate  float64          `json:"allocationRate"`
	AvgUtilization  float64          `json:"avgUtilization"`
	GPUTypes        map[string]int32 `json:"gpuTypes"`
	UnhealthyGPUs   int32            `json:"unhealthyGPUs"`
}

// JobClusterMetrics 作业集群指标
type JobClusterMetrics struct {
	TotalJobs       int32            `json:"totalJobs"`
	RunningJobs     int32            `json:"runningJobs"`
	PendingJobs     int32            `json:"pendingJobs"`
	CompletedJobs   int32            `json:"completedJobs"`
	FailedJobs      int32            `json:"failedJobs"`
	SuccessRate     float64          `json:"successRate"`
	FailureRate     float64          `json:"failureRate"`
	AvgWaitTime     float64          `json:"avgWaitTime"`
	AvgRunTime      float64          `json:"avgRunTime"`
	JobsByFramework map[string]int32 `json:"jobsByFramework"`
}

// QueueMetrics 队列指标
type QueueMetrics struct {
	QueueName     string            `json:"queueName"`
	JobsInQueue   int32             `json:"jobsInQueue"`
	RunningJobs   int32             `json:"runningJobs"`
	PendingJobs   int32             `json:"pendingJobs"`
	ResourceUsage map[string]string `json:"resourceUsage"`
	AvgWaitTime   float64           `json:"avgWaitTime"`
	Throughput    float64           `json:"throughput"`
	Priority      int32             `json:"priority"`
	Weight        int32             `json:"weight"`
}

// NodeMetrics 节点指标
type NodeMetrics struct {
	NodeName       string  `json:"nodeName"`
	Status         string  `json:"status"`
	CPUUsage       float64 `json:"cpuUsage"`
	MemoryUsage    float64 `json:"memoryUsage"`
	GPUUsage       float64 `json:"gpuUsage"`
	PodCount       int32   `json:"podCount"`
	DiskUsage      float64 `json:"diskUsage"`
	NetworkTraffic string  `json:"networkTraffic"`
}

// AlertSummary 告警摘要
type AlertSummary struct {
	CriticalAlerts int32 `json:"criticalAlerts"`
	WarningAlerts  int32 `json:"warningAlerts"`
	InfoAlerts     int32 `json:"infoAlerts"`
	ResolvedAlerts int32 `json:"resolvedAlerts"`
	NewAlerts      int32 `json:"newAlerts"`
}

// MonitoringEvent 监控事件
type MonitoringEvent struct {
	Type      string                 `json:"type"` // metric, alert, status_change
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"` // job, queue, node, gpu
	Data      map[string]interface{} `json:"data"`
	Level     string                 `json:"level"` // info, warning, error, critical
}

// NewMonitoringService 创建监控服务
func NewMonitoringService(client *Client, config *MonitoringConfig) *MonitoringService {
	return &MonitoringService{
		client:       client,
		gpuManager:   NewGPUManager(client),
		jobManager:   NewJobManager(client),
		alertManager: NewAlertManager(),
		config:       config,
		stopCh:       make(chan struct{}),
		eventQueue:   make(chan MonitoringEvent, 1000),
	}
}

// Start 启动监控服务
func (ms *MonitoringService) Start() error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	if ms.isRunning {
		return fmt.Errorf("监控服务已经在运行")
	}

	ms.isRunning = true

	// 启动各个监控协程
	if ms.config.EnableGPUMonitor {
		go ms.runGPUMonitoring()
	}

	if ms.config.EnableJobMonitor {
		go ms.runJobMonitoring()
	}

	if ms.config.EnableQueueMonitor {
		go ms.runQueueMonitoring()
	}

	// 启动事件处理协程
	go ms.processEvents()

	// 启动主监控循环
	go ms.runMainLoop()

	return nil
}

// Stop 停止监控服务
func (ms *MonitoringService) Stop() error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	if !ms.isRunning {
		return nil
	}

	close(ms.stopCh)
	ms.isRunning = false

	return nil
}

// runMainLoop 运行主监控循环
func (ms *MonitoringService) runMainLoop() {
	ticker := time.NewTicker(ms.config.UpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ms.stopCh:
			return
		case <-ticker.C:
			if err := ms.updateClusterMetrics(); err != nil {
				ms.logError(fmt.Sprintf("更新集群指标失败: %v", err))
			}
		}
	}
}

// updateClusterMetrics 更新集群指标
func (ms *MonitoringService) updateClusterMetrics() error {
	metrics := &ClusterMetrics{
		Timestamp: time.Now(),
	}

	// 收集GPU指标
	if ms.config.EnableGPUMonitor {
		gpuMetrics, err := ms.collectGPUMetrics()
		if err != nil {
			return fmt.Errorf("收集GPU指标失败: %v", err)
		}
		metrics.GPUMetrics = *gpuMetrics
	}

	// 收集作业指标
	if ms.config.EnableJobMonitor {
		jobMetrics, err := ms.collectJobMetrics()
		if err != nil {
			return fmt.Errorf("收集作业指标失败: %v", err)
		}
		metrics.JobMetrics = *jobMetrics
	}

	// 收集队列指标
	if ms.config.EnableQueueMonitor {
		queueMetrics, err := ms.collectQueueMetrics()
		if err != nil {
			return fmt.Errorf("收集队列指标失败: %v", err)
		}
		metrics.QueueMetrics = queueMetrics
	}

	// 收集节点指标
	nodeMetrics, err := ms.collectNodeMetrics()
	if err != nil {
		return fmt.Errorf("收集节点指标失败: %v", err)
	}
	metrics.NodeMetrics = nodeMetrics

	// 收集告警摘要
	alertSummary, err := ms.collectAlertSummary()
	if err != nil {
		return fmt.Errorf("收集告警摘要失败: %v", err)
	}
	metrics.AlertSummary = *alertSummary

	// 更新缓存
	ms.mutex.Lock()
	ms.clusterMetrics = metrics
	ms.lastUpdate = time.Now()
	ms.mutex.Unlock()

	// 检查告警条件
	ms.checkAlertConditions(metrics)

	// 发送监控事件
	event := MonitoringEvent{
		Type:      "metric",
		Timestamp: time.Now(),
		Source:    "cluster",
		Data:      map[string]interface{}{"metrics": metrics},
		Level:     "info",
	}
	ms.sendEvent(event)

	return nil
}

// collectGPUMetrics 收集GPU指标
func (ms *MonitoringService) collectGPUMetrics() (*GPUClusterMetrics, error) {
	gpuResources, err := ms.gpuManager.GetClusterGPUResources()
	if err != nil {
		return nil, err
	}

	metrics := &GPUClusterMetrics{
		GPUTypes: make(map[string]int32),
	}

	totalUtilization := float64(0)
	nodeCount := 0

	for _, gpuInfo := range gpuResources {
		metrics.TotalGPUs += gpuInfo.TotalGPUs
		metrics.AllocatedGPUs += gpuInfo.AllocatedGPUs
		metrics.AvailableGPUs += gpuInfo.AvailableGPUs

		if gpuInfo.Status != "Ready" {
			metrics.UnhealthyGPUs += gpuInfo.TotalGPUs
		}

		metrics.GPUTypes[gpuInfo.GPUType] += gpuInfo.TotalGPUs
		totalUtilization += gpuInfo.GPUUtilization
		nodeCount++
	}

	if metrics.TotalGPUs > 0 {
		metrics.AllocationRate = float64(metrics.AllocatedGPUs) / float64(metrics.TotalGPUs) * 100
		metrics.UtilizationRate = float64(metrics.TotalGPUs-metrics.AvailableGPUs) / float64(metrics.TotalGPUs) * 100
	}

	if nodeCount > 0 {
		metrics.AvgUtilization = totalUtilization / float64(nodeCount)
	}

	return metrics, nil
}

// collectJobMetrics 收集作业指标
func (ms *MonitoringService) collectJobMetrics() (*JobClusterMetrics, error) {
	// 获取所有作业
	jobs, err := ms.client.ListJobs("", "")
	if err != nil {
		return nil, err
	}

	metrics := &JobClusterMetrics{
		JobsByFramework: make(map[string]int32),
	}

	totalWaitTime := float64(0)
	totalRunTime := float64(0)
	jobsWithTiming := 0

	for _, job := range jobs.Items {
		metrics.TotalJobs++

		switch job.Status.State.Phase {
		case "Running":
			metrics.RunningJobs++
		case "Pending":
			metrics.PendingJobs++
		case "Completed":
			metrics.CompletedJobs++
		case "Failed":
			metrics.FailedJobs++
		}

		// 统计框架分布
		if framework, exists := job.Labels["framework"]; exists {
			metrics.JobsByFramework[framework]++
		}

		// 计算时间统计
		if !job.Status.State.LastTransitionTime.Time.IsZero() {
			if job.CreationTimestamp.Time.Before(job.Status.State.LastTransitionTime.Time) {
				waitTime := job.Status.State.LastTransitionTime.Time.Sub(job.CreationTimestamp.Time).Seconds()
				totalWaitTime += waitTime
				jobsWithTiming++
			}
		}
	}

	// 计算成功率和失败率
	if metrics.TotalJobs > 0 {
		metrics.SuccessRate = float64(metrics.CompletedJobs) / float64(metrics.TotalJobs) * 100
		metrics.FailureRate = float64(metrics.FailedJobs) / float64(metrics.TotalJobs) * 100
	}

	// 计算平均等待时间
	if jobsWithTiming > 0 {
		metrics.AvgWaitTime = totalWaitTime / float64(jobsWithTiming)
		metrics.AvgRunTime = totalRunTime / float64(jobsWithTiming)
	}

	return metrics, nil
}

// collectQueueMetrics 收集队列指标
func (ms *MonitoringService) collectQueueMetrics() ([]QueueMetrics, error) {
	queues, err := ms.client.ListQueues("")
	if err != nil {
		return nil, err
	}

	var queueMetrics []QueueMetrics
	for _, queue := range queues.Items {
		metrics := QueueMetrics{
			QueueName:     queue.Name,
			RunningJobs:   queue.Status.Running,
			PendingJobs:   queue.Status.Pending,
			JobsInQueue:   queue.Status.Inqueue,
			Priority:      int32(queue.Spec.Weight), // 简化处理
			Weight:        queue.Spec.Weight,
			ResourceUsage: map[string]string{}, // 暂时为空，Volcano v1beta1版本中没有Used字段
		}

		// TODO: 计算平均等待时间和吞吐量
		metrics.AvgWaitTime = 0
		metrics.Throughput = 0

		queueMetrics = append(queueMetrics, metrics)
	}

	return queueMetrics, nil
}

// collectNodeMetrics 收集节点指标
func (ms *MonitoringService) collectNodeMetrics() ([]NodeMetrics, error) {
	nodes, err := ms.client.kubeClient.CoreV1().Nodes().List(
		context.TODO(),
		metav1.ListOptions{},
	)
	if err != nil {
		return nil, err
	}

	var nodeMetrics []NodeMetrics
	for _, node := range nodes.Items {
		metrics := NodeMetrics{
			NodeName: node.Name,
			Status:   string(getNodeConditionStatus(&node, "Ready")),
		}

		// TODO: 从监控系统获取实际的资源使用率
		metrics.CPUUsage = 0
		metrics.MemoryUsage = 0
		metrics.GPUUsage = 0
		metrics.DiskUsage = 0
		metrics.NetworkTraffic = "N/A"

		// 获取Pod数量
		pods, err := ms.client.kubeClient.CoreV1().Pods("").List(
			context.TODO(),
			metav1.ListOptions{
				FieldSelector: fmt.Sprintf("spec.nodeName=%s", node.Name),
			},
		)
		if err == nil {
			metrics.PodCount = int32(len(pods.Items))
		}

		nodeMetrics = append(nodeMetrics, metrics)
	}

	return nodeMetrics, nil
}

// collectAlertSummary 收集告警摘要
func (ms *MonitoringService) collectAlertSummary() (*AlertSummary, error) {
	// TODO: 从告警管理器获取告警统计
	return &AlertSummary{
		CriticalAlerts: 0,
		WarningAlerts:  0,
		InfoAlerts:     0,
		ResolvedAlerts: 0,
		NewAlerts:      0,
	}, nil
}

// checkAlertConditions 检查告警条件
func (ms *MonitoringService) checkAlertConditions(metrics *ClusterMetrics) {
	// 检查GPU使用率告警
	if metrics.GPUMetrics.AvgUtilization > ms.config.AlertThresholds.GPUUtilization {
		alert := Alert{
			ID:       fmt.Sprintf("gpu-utilization-%d", time.Now().Unix()),
			Type:     "resource",
			Severity: "warning",
			Title:    "GPU使用率过高",
			Description: fmt.Sprintf("集群平均GPU使用率 %.2f%% 超过阈值 %.2f%%",
				metrics.GPUMetrics.AvgUtilization, ms.config.AlertThresholds.GPUUtilization),
			Source:    "monitoring",
			Timestamp: time.Now(),
			Status:    "active",
		}
		ms.alertManager.TriggerAlert(alert)
	}

	// 检查作业失败率告警
	if metrics.JobMetrics.FailureRate > ms.config.AlertThresholds.JobFailureRate {
		alert := Alert{
			ID:       fmt.Sprintf("job-failure-rate-%d", time.Now().Unix()),
			Type:     "job",
			Severity: "critical",
			Title:    "作业失败率过高",
			Description: fmt.Sprintf("作业失败率 %.2f%% 超过阈值 %.2f%%",
				metrics.JobMetrics.FailureRate, ms.config.AlertThresholds.JobFailureRate),
			Source:    "monitoring",
			Timestamp: time.Now(),
			Status:    "active",
		}
		ms.alertManager.TriggerAlert(alert)
	}

	// 检查节点健康状况
	unhealthyNodes := 0
	for _, nodeMetric := range metrics.NodeMetrics {
		if nodeMetric.Status != "True" {
			unhealthyNodes++
		}
	}

	if len(metrics.NodeMetrics) > 0 {
		unhealthyRate := float64(unhealthyNodes) / float64(len(metrics.NodeMetrics))
		if unhealthyRate > ms.config.AlertThresholds.NodeUnhealthyRate {
			alert := Alert{
				ID:       fmt.Sprintf("node-unhealthy-%d", time.Now().Unix()),
				Type:     "infrastructure",
				Severity: "critical",
				Title:    "节点异常率过高",
				Description: fmt.Sprintf("异常节点比例 %.2f%% 超过阈值 %.2f%%",
					unhealthyRate*100, ms.config.AlertThresholds.NodeUnhealthyRate*100),
				Source:    "monitoring",
				Timestamp: time.Now(),
				Status:    "active",
			}
			ms.alertManager.TriggerAlert(alert)
		}
	}
}

// runGPUMonitoring 运行GPU监控
func (ms *MonitoringService) runGPUMonitoring() {
	// TODO: 实现GPU实时监控逻辑
}

// runJobMonitoring 运行作业监控
func (ms *MonitoringService) runJobMonitoring() {
	// 监控作业状态变化
	watcher, err := ms.client.volcanoClient.BatchV1alpha1().Jobs("").Watch(
		context.TODO(),
		metav1.ListOptions{},
	)
	if err != nil {
		ms.logError(fmt.Sprintf("创建作业监控失败: %v", err))
		return
	}
	defer watcher.Stop()

	for {
		select {
		case <-ms.stopCh:
			return
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return
			}
			ms.handleJobEvent(event)
		}
	}
}

// runQueueMonitoring 运行队列监控
func (ms *MonitoringService) runQueueMonitoring() {
	// 监控队列状态变化
	watcher, err := ms.client.volcanoClient.SchedulingV1beta1().Queues().Watch(
		context.TODO(),
		metav1.ListOptions{},
	)
	if err != nil {
		ms.logError(fmt.Sprintf("创建队列监控失败: %v", err))
		return
	}
	defer watcher.Stop()

	for {
		select {
		case <-ms.stopCh:
			return
		case event, ok := <-watcher.ResultChan():
			if !ok {
				return
			}
			ms.handleQueueEvent(event)
		}
	}
}

// handleJobEvent 处理作业事件
func (ms *MonitoringService) handleJobEvent(event watch.Event) {
	monitoringEvent := MonitoringEvent{
		Type:      "status_change",
		Timestamp: time.Now(),
		Source:    "job",
		Data: map[string]interface{}{
			"eventType": string(event.Type),
			"object":    event.Object,
		},
		Level: "info",
	}

	ms.sendEvent(monitoringEvent)
}

// handleQueueEvent 处理队列事件
func (ms *MonitoringService) handleQueueEvent(event watch.Event) {
	monitoringEvent := MonitoringEvent{
		Type:      "status_change",
		Timestamp: time.Now(),
		Source:    "queue",
		Data: map[string]interface{}{
			"eventType": string(event.Type),
			"object":    event.Object,
		},
		Level: "info",
	}

	ms.sendEvent(monitoringEvent)
}

// processEvents 处理监控事件
func (ms *MonitoringService) processEvents() {
	for {
		select {
		case <-ms.stopCh:
			return
		case event := <-ms.eventQueue:
			ms.handleMonitoringEvent(event)
		}
	}
}

// handleMonitoringEvent 处理监控事件
func (ms *MonitoringService) handleMonitoringEvent(event MonitoringEvent) {
	// TODO: 实现事件处理逻辑，如持久化、分析、通知等
}

// sendEvent 发送监控事件
func (ms *MonitoringService) sendEvent(event MonitoringEvent) {
	select {
	case ms.eventQueue <- event:
	default:
		// 事件队列满，丢弃事件
		ms.logError("监控事件队列已满，丢弃事件")
	}
}

// GetClusterMetrics 获取集群指标
func (ms *MonitoringService) GetClusterMetrics() (*ClusterMetrics, error) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	if ms.clusterMetrics == nil {
		return nil, fmt.Errorf("集群指标尚未初始化")
	}

	return ms.clusterMetrics, nil
}

// GetHealthStatus 获取健康状态
func (ms *MonitoringService) GetHealthStatus() (*HealthStatus, error) {
	metrics, err := ms.GetClusterMetrics()
	if err != nil {
		return nil, err
	}

	status := &HealthStatus{
		Overall:    "healthy",
		Timestamp:  time.Now(),
		Components: make(map[string]ComponentHealth),
	}

	// 评估各组件健康状况
	status.Components["gpu"] = ms.evaluateGPUHealth(&metrics.GPUMetrics)
	status.Components["jobs"] = ms.evaluateJobHealth(&metrics.JobMetrics)
	status.Components["nodes"] = ms.evaluateNodeHealth(metrics.NodeMetrics)

	// 计算整体健康状况
	criticalIssues := 0
	warningIssues := 0

	for _, component := range status.Components {
		switch component.Status {
		case "critical":
			criticalIssues++
		case "warning":
			warningIssues++
		}
	}

	if criticalIssues > 0 {
		status.Overall = "critical"
	} else if warningIssues > 0 {
		status.Overall = "warning"
	}

	return status, nil
}

// HealthStatus 健康状态
type HealthStatus struct {
	Overall    string                     `json:"overall"`
	Timestamp  time.Time                  `json:"timestamp"`
	Components map[string]ComponentHealth `json:"components"`
}

// ComponentHealth 组件健康状况
type ComponentHealth struct {
	Status  string                 `json:"status"` // healthy, warning, critical
	Issues  []string               `json:"issues"`
	Metrics map[string]interface{} `json:"metrics"`
}

// evaluateGPUHealth 评估GPU健康状况
func (ms *MonitoringService) evaluateGPUHealth(metrics *GPUClusterMetrics) ComponentHealth {
	health := ComponentHealth{
		Status:  "healthy",
		Issues:  make([]string, 0),
		Metrics: make(map[string]interface{}),
	}

	health.Metrics["totalGPUs"] = metrics.TotalGPUs
	health.Metrics["availableGPUs"] = metrics.AvailableGPUs
	health.Metrics["utilizationRate"] = metrics.UtilizationRate

	if metrics.UnhealthyGPUs > 0 {
		health.Status = "warning"
		health.Issues = append(health.Issues, fmt.Sprintf("%d个GPU状态异常", metrics.UnhealthyGPUs))
	}

	if metrics.AllocationRate > 90 {
		health.Status = "warning"
		health.Issues = append(health.Issues, fmt.Sprintf("GPU分配率过高 (%.1f%%)", metrics.AllocationRate))
	}

	return health
}

// evaluateJobHealth 评估作业健康状况
func (ms *MonitoringService) evaluateJobHealth(metrics *JobClusterMetrics) ComponentHealth {
	health := ComponentHealth{
		Status:  "healthy",
		Issues:  make([]string, 0),
		Metrics: make(map[string]interface{}),
	}

	health.Metrics["totalJobs"] = metrics.TotalJobs
	health.Metrics["successRate"] = metrics.SuccessRate
	health.Metrics["failureRate"] = metrics.FailureRate

	if metrics.FailureRate > ms.config.AlertThresholds.JobFailureRate {
		health.Status = "critical"
		health.Issues = append(health.Issues, fmt.Sprintf("作业失败率过高 (%.1f%%)", metrics.FailureRate))
	}

	if metrics.PendingJobs > 0 && metrics.AvgWaitTime > float64(ms.config.AlertThresholds.QueueWaitTime) {
		health.Status = "warning"
		health.Issues = append(health.Issues, fmt.Sprintf("作业平均等待时间过长 (%.1f秒)", metrics.AvgWaitTime))
	}

	return health
}

// evaluateNodeHealth 评估节点健康状况
func (ms *MonitoringService) evaluateNodeHealth(nodeMetrics []NodeMetrics) ComponentHealth {
	health := ComponentHealth{
		Status:  "healthy",
		Issues:  make([]string, 0),
		Metrics: make(map[string]interface{}),
	}

	totalNodes := len(nodeMetrics)
	unhealthyNodes := 0

	for _, node := range nodeMetrics {
		if node.Status != "True" {
			unhealthyNodes++
		}
	}

	health.Metrics["totalNodes"] = totalNodes
	health.Metrics["unhealthyNodes"] = unhealthyNodes

	if totalNodes > 0 {
		unhealthyRate := float64(unhealthyNodes) / float64(totalNodes)
		health.Metrics["unhealthyRate"] = unhealthyRate * 100

		if unhealthyRate > ms.config.AlertThresholds.NodeUnhealthyRate {
			health.Status = "critical"
			health.Issues = append(health.Issues, fmt.Sprintf("节点异常率过高 (%.1f%%)", unhealthyRate*100))
		} else if unhealthyNodes > 0 {
			health.Status = "warning"
			health.Issues = append(health.Issues, fmt.Sprintf("%d个节点状态异常", unhealthyNodes))
		}
	}

	return health
}

// logError 记录错误日志
func (ms *MonitoringService) logError(message string) {
	// TODO: 集成日志系统
	fmt.Printf("[ERROR] %s: %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}

// 辅助函数
func getNodeConditionStatus(node *corev1.Node, conditionType corev1.NodeConditionType) corev1.ConditionStatus {
	for _, condition := range node.Status.Conditions {
		if condition.Type == conditionType {
			return condition.Status
		}
	}
	return corev1.ConditionUnknown
}
