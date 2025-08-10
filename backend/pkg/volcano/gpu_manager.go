package volcano

import (
	"context"
	"fmt"
	"sort"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GPUManager GPU资源管理器
type GPUManager struct {
	client *Client
}

// NewGPUManager 创建GPU管理器
func NewGPUManager(client *Client) *GPUManager {
	return &GPUManager{
		client: client,
	}
}

// GPUResourceInfo GPU资源信息
type GPUResourceInfo struct {
	NodeName       string            `json:"nodeName"`
	GPUType        string            `json:"gpuType"`
	TotalGPUs      int32             `json:"totalGPUs"`
	AvailableGPUs  int32             `json:"availableGPUs"`
	AllocatedGPUs  int32             `json:"allocatedGPUs"`
	GPUUtilization float64           `json:"gpuUtilization"`
	MemoryTotal    string            `json:"memoryTotal"`
	MemoryUsed     string            `json:"memoryUsed"`
	MemoryFree     string            `json:"memoryFree"`
	PowerDraw      string            `json:"powerDraw"`
	Temperature    string            `json:"temperature"`
	Labels         map[string]string `json:"labels"`
	Taints         []string          `json:"taints"`
	Status         string            `json:"status"` // Ready, NotReady, Unknown
}

// GPUAllocationStrategy GPU分配策略
type GPUAllocationStrategy struct {
	Strategy        string   `json:"strategy"` // binpack, spread, affinity, topology
	PreferredNodes  []string `json:"preferredNodes,omitempty"`
	ExcludedNodes   []string `json:"excludedNodes,omitempty"`
	GPUTypes        []string `json:"gpuTypes,omitempty"`        // 优先GPU类型
	MinGPUMemory    string   `json:"minGpuMemory,omitempty"`    // 最小GPU内存
	TopologyKey     string   `json:"topologyKey,omitempty"`     // 拓扑键，如"kubernetes.io/hostname"
	MaxSkew         int32    `json:"maxSkew,omitempty"`         // 最大偏差
	AntiAffinity    bool     `json:"antiAffinity"`              // 反亲和性
	ColocateWithJob string   `json:"colocateWithJob,omitempty"` // 与指定作业在同一节点
}

// GPUAllocationRequest GPU分配请求
type GPUAllocationRequest struct {
	JobName      string                `json:"jobName"`
	TaskName     string                `json:"taskName"`
	GPUCount     int32                 `json:"gpuCount"`
	GPUMemoryMin string                `json:"gpuMemoryMin,omitempty"`
	Strategy     GPUAllocationStrategy `json:"strategy"`
	NodeSelector map[string]string     `json:"nodeSelector,omitempty"`
	Tolerations  []corev1.Toleration   `json:"tolerations,omitempty"`
	Priority     int32                 `json:"priority"`
	Queue        string                `json:"queue"`
}

// GPUAllocationResult GPU分配结果
type GPUAllocationResult struct {
	Success       bool           `json:"success"`
	AllocatedGPUs []AllocatedGPU `json:"allocatedGPUs"`
	Message       string         `json:"message"`
	Reason        string         `json:"reason,omitempty"`
	Suggestions   []string       `json:"suggestions,omitempty"`
}

// AllocatedGPU 已分配的GPU
type AllocatedGPU struct {
	NodeName     string `json:"nodeName"`
	GPUIndex     int32  `json:"gpuIndex"`
	GPUUUID      string `json:"gpuUUID"`
	GPUType      string `json:"gpuType"`
	MemoryTotal  string `json:"memoryTotal"`
	AllocationID string `json:"allocationId"`
}

// GetClusterGPUResources 获取集群GPU资源信息
func (gm *GPUManager) GetClusterGPUResources() ([]GPUResourceInfo, error) {
	// 获取所有节点
	nodes, err := gm.client.kubeClient.CoreV1().Nodes().List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: "node.kubernetes.io/instance-type!=virtual-node",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("获取节点列表失败: %v", err)
	}

	var gpuInfos []GPUResourceInfo
	for _, node := range nodes.Items {
		if gpuInfo := gm.extractGPUInfo(&node); gpuInfo != nil {
			gpuInfos = append(gpuInfos, *gpuInfo)
		}
	}

	// 按节点名排序
	sort.Slice(gpuInfos, func(i, j int) bool {
		return gpuInfos[i].NodeName < gpuInfos[j].NodeName
	})

	return gpuInfos, nil
}

// extractGPUInfo 提取节点的GPU信息
func (gm *GPUManager) extractGPUInfo(node *corev1.Node) *GPUResourceInfo {
	// 检查节点是否有GPU资源
	gpuCapacity, hasGPU := node.Status.Capacity["nvidia.com/gpu"]
	if !hasGPU || gpuCapacity.IsZero() {
		return nil
	}

	gpuAllocatable := node.Status.Allocatable["nvidia.com/gpu"]
	totalGPUs := int32(gpuCapacity.Value())
	availableGPUs := int32(gpuAllocatable.Value())
	allocatedGPUs := totalGPUs - availableGPUs

	gpuInfo := &GPUResourceInfo{
		NodeName:      node.Name,
		TotalGPUs:     totalGPUs,
		AvailableGPUs: availableGPUs,
		AllocatedGPUs: allocatedGPUs,
		Labels:        node.Labels,
		Status:        gm.getNodeGPUStatus(node),
	}

	// 提取GPU类型
	if gpuType, exists := node.Labels["nvidia.com/gpu.product"]; exists {
		gpuInfo.GPUType = gpuType
	} else if gpuType, exists := node.Labels["gpu.nvidia.com/class"]; exists {
		gpuInfo.GPUType = gpuType
	} else {
		gpuInfo.GPUType = "Unknown"
	}

	// 提取GPU内存信息
	if gpuMemory, exists := node.Labels["nvidia.com/gpu.memory"]; exists {
		gpuInfo.MemoryTotal = gpuMemory
	}

	// 提取污点信息
	for _, taint := range node.Spec.Taints {
		taintStr := fmt.Sprintf("%s=%s:%s", taint.Key, taint.Value, taint.Effect)
		gpuInfo.Taints = append(gpuInfo.Taints, taintStr)
	}

	// TODO: 从监控系统获取实时GPU使用率、功耗、温度等信息
	gpuInfo.GPUUtilization = gm.getNodeGPUUtilization(node.Name)
	gpuInfo.PowerDraw = "N/A"
	gpuInfo.Temperature = "N/A"

	return gpuInfo
}

// getNodeGPUStatus 获取节点GPU状态
func (gm *GPUManager) getNodeGPUStatus(node *corev1.Node) string {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			if condition.Status == corev1.ConditionTrue {
				return "Ready"
			} else {
				return "NotReady"
			}
		}
	}
	return "Unknown"
}

// getNodeGPUUtilization 获取节点GPU使用率
func (gm *GPUManager) getNodeGPUUtilization(nodeName string) float64 {
	// TODO: 集成Prometheus或其他监控系统获取实时GPU使用率
	// 这里返回模拟数据
	return 0.0
}

// AllocateGPUs 分配GPU资源
func (gm *GPUManager) AllocateGPUs(req *GPUAllocationRequest) (*GPUAllocationResult, error) {
	// 获取可用的GPU资源
	gpuResources, err := gm.GetClusterGPUResources()
	if err != nil {
		return nil, fmt.Errorf("获取GPU资源失败: %v", err)
	}

	// 过滤可用节点
	availableNodes := gm.filterAvailableNodes(gpuResources, req)
	if len(availableNodes) == 0 {
		return &GPUAllocationResult{
			Success: false,
			Message: "没有满足条件的可用GPU节点",
			Suggestions: []string{
				"检查节点选择器和容忍度配置",
				"确认集群中有足够的可用GPU资源",
				"考虑调整资源请求要求",
			},
		}, nil
	}

	// 根据分配策略选择节点和GPU
	allocation, err := gm.selectGPUsForAllocation(availableNodes, req)
	if err != nil {
		return &GPUAllocationResult{
			Success: false,
			Message: fmt.Sprintf("GPU分配失败: %v", err),
			Reason:  "策略执行错误",
		}, nil
	}

	if len(allocation) < int(req.GPUCount) {
		return &GPUAllocationResult{
			Success: false,
			Message: fmt.Sprintf("只能分配 %d 个GPU，需要 %d 个", len(allocation), req.GPUCount),
			Suggestions: []string{
				"降低GPU资源请求数量",
				"等待更多资源释放",
				"使用不同的分配策略",
			},
		}, nil
	}

	return &GPUAllocationResult{
		Success:       true,
		AllocatedGPUs: allocation,
		Message:       fmt.Sprintf("成功分配 %d 个GPU", len(allocation)),
	}, nil
}

// filterAvailableNodes 过滤可用节点
func (gm *GPUManager) filterAvailableNodes(gpuResources []GPUResourceInfo, req *GPUAllocationRequest) []GPUResourceInfo {
	var availableNodes []GPUResourceInfo

	for _, gpuInfo := range gpuResources {
		// 检查节点状态
		if gpuInfo.Status != "Ready" {
			continue
		}

		// 检查可用GPU数量
		if gpuInfo.AvailableGPUs <= 0 {
			continue
		}

		// 检查节点选择器
		if !gm.matchesNodeSelector(gpuInfo.Labels, req.NodeSelector) {
			continue
		}

		// 检查污点和容忍度
		if !gm.canTolerateTaints(gpuInfo.Taints, req.Tolerations) {
			continue
		}

		// 检查GPU类型限制
		if len(req.Strategy.GPUTypes) > 0 && !gm.containsString(req.Strategy.GPUTypes, gpuInfo.GPUType) {
			continue
		}

		// 检查最小GPU内存要求
		if req.GPUMemoryMin != "" && !gm.hasEnoughGPUMemory(gpuInfo.MemoryTotal, req.GPUMemoryMin) {
			continue
		}

		// 检查排除节点
		if gm.containsString(req.Strategy.ExcludedNodes, gpuInfo.NodeName) {
			continue
		}

		availableNodes = append(availableNodes, gpuInfo)
	}

	return availableNodes
}

// matchesNodeSelector 检查节点是否匹配选择器
func (gm *GPUManager) matchesNodeSelector(nodeLabels, nodeSelector map[string]string) bool {
	for key, value := range nodeSelector {
		if nodeValue, exists := nodeLabels[key]; !exists || nodeValue != value {
			return false
		}
	}
	return true
}

// canTolerateTaints 检查是否能够容忍污点
func (gm *GPUManager) canTolerateTaints(nodeTaints []string, tolerations []corev1.Toleration) bool {
	// 简化实现：假设没有不可容忍的污点
	// 实际实现需要详细匹配污点和容忍度
	return true
}

// hasEnoughGPUMemory 检查GPU内存是否足够
func (gm *GPUManager) hasEnoughGPUMemory(totalMemory, requiredMemory string) bool {
	if totalMemory == "" || requiredMemory == "" {
		return true // 如果无法判断，默认允许
	}

	totalQ, err1 := resource.ParseQuantity(totalMemory)
	requiredQ, err2 := resource.ParseQuantity(requiredMemory)

	if err1 != nil || err2 != nil {
		return true
	}

	return totalQ.Cmp(requiredQ) >= 0
}

// containsString 检查字符串切片是否包含指定字符串
func (gm *GPUManager) containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// selectGPUsForAllocation 根据策略选择GPU进行分配
func (gm *GPUManager) selectGPUsForAllocation(availableNodes []GPUResourceInfo, req *GPUAllocationRequest) ([]AllocatedGPU, error) {
	switch req.Strategy.Strategy {
	case "binpack":
		return gm.binpackAllocation(availableNodes, req)
	case "spread":
		return gm.spreadAllocation(availableNodes, req)
	case "affinity":
		return gm.affinityAllocation(availableNodes, req)
	case "topology":
		return gm.topologyAllocation(availableNodes, req)
	default:
		return gm.binpackAllocation(availableNodes, req) // 默认使用binpack策略
	}
}

// binpackAllocation Binpack分配策略：优先填满单个节点
func (gm *GPUManager) binpackAllocation(availableNodes []GPUResourceInfo, req *GPUAllocationRequest) ([]AllocatedGPU, error) {
	// 按可用GPU数量降序排序，优先使用GPU较多的节点
	sort.Slice(availableNodes, func(i, j int) bool {
		return availableNodes[i].AvailableGPUs > availableNodes[j].AvailableGPUs
	})

	var allocatedGPUs []AllocatedGPU
	remainingGPUs := req.GPUCount

	for _, node := range availableNodes {
		if remainingGPUs <= 0 {
			break
		}

		// 在当前节点分配GPU
		nodeAllocation := int32(0)
		if node.AvailableGPUs >= remainingGPUs {
			nodeAllocation = remainingGPUs
		} else {
			nodeAllocation = node.AvailableGPUs
		}

		// 创建GPU分配记录
		for i := int32(0); i < nodeAllocation; i++ {
			allocatedGPU := AllocatedGPU{
				NodeName:     node.NodeName,
				GPUIndex:     i,
				GPUUUID:      fmt.Sprintf("GPU-%s-%d", node.NodeName, i),
				GPUType:      node.GPUType,
				MemoryTotal:  node.MemoryTotal,
				AllocationID: fmt.Sprintf("%s-%s-%d", req.JobName, req.TaskName, len(allocatedGPUs)),
			}
			allocatedGPUs = append(allocatedGPUs, allocatedGPU)
		}

		remainingGPUs -= nodeAllocation
	}

	return allocatedGPUs, nil
}

// spreadAllocation Spread分配策略：将GPU分散到多个节点
func (gm *GPUManager) spreadAllocation(availableNodes []GPUResourceInfo, req *GPUAllocationRequest) ([]AllocatedGPU, error) {
	// 按可用GPU数量升序排序，优先使用GPU较少的节点
	sort.Slice(availableNodes, func(i, j int) bool {
		return availableNodes[i].AvailableGPUs < availableNodes[j].AvailableGPUs
	})

	var allocatedGPUs []AllocatedGPU
	remainingGPUs := req.GPUCount

	// 轮询分配GPU
	for remainingGPUs > 0 {
		allocated := false
		for i, node := range availableNodes {
			if remainingGPUs <= 0 {
				break
			}
			if node.AvailableGPUs > 0 {
				// 分配一个GPU
				allocatedGPU := AllocatedGPU{
					NodeName:     node.NodeName,
					GPUIndex:     node.TotalGPUs - node.AvailableGPUs,
					GPUUUID:      fmt.Sprintf("GPU-%s-%d", node.NodeName, node.TotalGPUs-node.AvailableGPUs),
					GPUType:      node.GPUType,
					MemoryTotal:  node.MemoryTotal,
					AllocationID: fmt.Sprintf("%s-%s-%d", req.JobName, req.TaskName, len(allocatedGPUs)),
				}
				allocatedGPUs = append(allocatedGPUs, allocatedGPU)

				// 更新节点可用GPU数量
				availableNodes[i].AvailableGPUs--
				remainingGPUs--
				allocated = true
			}
		}
		if !allocated {
			break // 没有更多可用GPU
		}
	}

	return allocatedGPUs, nil
}

// affinityAllocation 亲和性分配策略：优先分配到首选节点
func (gm *GPUManager) affinityAllocation(availableNodes []GPUResourceInfo, req *GPUAllocationRequest) ([]AllocatedGPU, error) {
	// 优先处理首选节点
	var preferredNodes, otherNodes []GPUResourceInfo

	for _, node := range availableNodes {
		if gm.containsString(req.Strategy.PreferredNodes, node.NodeName) {
			preferredNodes = append(preferredNodes, node)
		} else {
			otherNodes = append(otherNodes, node)
		}
	}

	// 首先在首选节点分配
	reorderedNodes := append(preferredNodes, otherNodes...)

	return gm.binpackAllocation(reorderedNodes, req)
}

// topologyAllocation 拓扑感知分配策略：考虑节点拓扑分布
func (gm *GPUManager) topologyAllocation(availableNodes []GPUResourceInfo, req *GPUAllocationRequest) ([]AllocatedGPU, error) {
	if req.Strategy.TopologyKey == "" {
		return gm.binpackAllocation(availableNodes, req)
	}

	// 按拓扑键分组节点
	topologyGroups := make(map[string][]GPUResourceInfo)
	for _, node := range availableNodes {
		if topologyValue, exists := node.Labels[req.Strategy.TopologyKey]; exists {
			topologyGroups[topologyValue] = append(topologyGroups[topologyValue], node)
		}
	}

	// 在拓扑组内进行分配，尽量保持均匀分布
	var allocatedGPUs []AllocatedGPU
	remainingGPUs := req.GPUCount

	for _, group := range topologyGroups {
		if remainingGPUs <= 0 {
			break
		}

		// 计算当前组应该分配的GPU数量
		groupAllocation := remainingGPUs
		if int32(len(group)) < remainingGPUs {
			groupAllocation = int32(len(group))
		}

		// 在组内使用spread策略分配
		subReq := *req
		subReq.GPUCount = groupAllocation

		groupAllocated, err := gm.spreadAllocation(group, &subReq)
		if err != nil {
			return nil, err
		}

		allocatedGPUs = append(allocatedGPUs, groupAllocated...)
		remainingGPUs -= int32(len(groupAllocated))
	}

	return allocatedGPUs, nil
}

// GetGPUUtilizationReport 获取GPU使用率报告
func (gm *GPUManager) GetGPUUtilizationReport(timeRange string) (*GPUUtilizationReport, error) {
	gpuResources, err := gm.GetClusterGPUResources()
	if err != nil {
		return nil, err
	}

	report := &GPUUtilizationReport{
		Timestamp:      time.Now().Format("2006-01-02 15:04:05"),
		TimeRange:      timeRange,
		NodeReports:    make([]NodeGPUReport, 0),
		ClusterSummary: ClusterGPUSummary{},
	}

	totalGPUs := int32(0)
	allocatedGPUs := int32(0)
	avgUtilization := float64(0)

	for _, gpuInfo := range gpuResources {
		nodeReport := NodeGPUReport{
			NodeName:      gpuInfo.NodeName,
			GPUType:       gpuInfo.GPUType,
			TotalGPUs:     gpuInfo.TotalGPUs,
			AllocatedGPUs: gpuInfo.AllocatedGPUs,
			AvailableGPUs: gpuInfo.AvailableGPUs,
			Utilization:   gpuInfo.GPUUtilization,
			MemoryUsed:    gpuInfo.MemoryUsed,
			MemoryTotal:   gpuInfo.MemoryTotal,
			Status:        gpuInfo.Status,
		}
		report.NodeReports = append(report.NodeReports, nodeReport)

		totalGPUs += gpuInfo.TotalGPUs
		allocatedGPUs += gpuInfo.AllocatedGPUs
		avgUtilization += gpuInfo.GPUUtilization
	}

	if len(gpuResources) > 0 {
		avgUtilization /= float64(len(gpuResources))
	}

	report.ClusterSummary = ClusterGPUSummary{
		TotalGPUs:      totalGPUs,
		AllocatedGPUs:  allocatedGPUs,
		AvailableGPUs:  totalGPUs - allocatedGPUs,
		AvgUtilization: avgUtilization,
		AllocationRate: float64(allocatedGPUs) / float64(totalGPUs) * 100,
	}

	return report, nil
}

// GPUUtilizationReport GPU使用率报告
type GPUUtilizationReport struct {
	Timestamp      string            `json:"timestamp"`
	TimeRange      string            `json:"timeRange"`
	NodeReports    []NodeGPUReport   `json:"nodeReports"`
	ClusterSummary ClusterGPUSummary `json:"clusterSummary"`
}

// NodeGPUReport 节点GPU报告
type NodeGPUReport struct {
	NodeName      string  `json:"nodeName"`
	GPUType       string  `json:"gpuType"`
	TotalGPUs     int32   `json:"totalGPUs"`
	AllocatedGPUs int32   `json:"allocatedGPUs"`
	AvailableGPUs int32   `json:"availableGPUs"`
	Utilization   float64 `json:"utilization"`
	MemoryUsed    string  `json:"memoryUsed"`
	MemoryTotal   string  `json:"memoryTotal"`
	Status        string  `json:"status"`
}

// ClusterGPUSummary 集群GPU摘要
type ClusterGPUSummary struct {
	TotalGPUs      int32   `json:"totalGPUs"`
	AllocatedGPUs  int32   `json:"allocatedGPUs"`
	AvailableGPUs  int32   `json:"availableGPUs"`
	AvgUtilization float64 `json:"avgUtilization"`
	AllocationRate float64 `json:"allocationRate"`
}

// OptimizeGPUPlacement GPU放置优化建议
func (gm *GPUManager) OptimizeGPUPlacement(jobRequirements []GPUAllocationRequest) (*GPUOptimizationSuggestion, error) {
	gpuResources, err := gm.GetClusterGPUResources()
	if err != nil {
		return nil, err
	}

	suggestion := &GPUOptimizationSuggestion{
		Timestamp:       time.Now().Format("2006-01-02 15:04:05"),
		Optimizations:   make([]PlacementOptimization, 0),
		ResourceAlerts:  make([]ResourceAlert, 0),
		Recommendations: make([]string, 0),
	}

	// 分析当前资源状况
	totalGPUs := int32(0)
	availableGPUs := int32(0)

	for _, gpuInfo := range gpuResources {
		totalGPUs += gpuInfo.TotalGPUs
		availableGPUs += gpuInfo.AvailableGPUs

		// 检查资源告警
		utilizationRate := float64(gpuInfo.AllocatedGPUs) / float64(gpuInfo.TotalGPUs)
		if utilizationRate > 0.9 {
			alert := ResourceAlert{
				Type:      "critical",
				Resource:  "gpu",
				Message:   fmt.Sprintf("节点 %s GPU使用率过高 (%.1f%%)", gpuInfo.NodeName, utilizationRate*100),
				UsageRate: utilizationRate,
				Threshold: 0.9,
				Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			}
			suggestion.ResourceAlerts = append(suggestion.ResourceAlerts, alert)
		}
	}

	// 分析作业需求和优化建议
	totalRequestedGPUs := int32(0)
	for _, req := range jobRequirements {
		totalRequestedGPUs += req.GPUCount
	}

	// 生成优化建议
	if totalRequestedGPUs > availableGPUs {
		suggestion.Recommendations = append(suggestion.Recommendations,
			fmt.Sprintf("需要 %d 个GPU，但只有 %d 个可用，建议调整作业优先级或增加集群容量", totalRequestedGPUs, availableGPUs))
	}

	// 检查GPU类型分布
	gpuTypeCount := make(map[string]int32)
	for _, gpuInfo := range gpuResources {
		gpuTypeCount[gpuInfo.GPUType] += gpuInfo.TotalGPUs
	}

	if len(gpuTypeCount) > 1 {
		suggestion.Recommendations = append(suggestion.Recommendations,
			"集群中存在多种GPU类型，建议为不同类型的作业指定合适的GPU类型以提高资源利用率")
	}

	return suggestion, nil
}

// GPUOptimizationSuggestion GPU优化建议
type GPUOptimizationSuggestion struct {
	Timestamp       string                  `json:"timestamp"`
	Optimizations   []PlacementOptimization `json:"optimizations"`
	ResourceAlerts  []ResourceAlert         `json:"resourceAlerts"`
	Recommendations []string                `json:"recommendations"`
}

// PlacementOptimization 放置优化建议
type PlacementOptimization struct {
	JobName           string   `json:"jobName"`
	CurrentStrategy   string   `json:"currentStrategy"`
	SuggestedStrategy string   `json:"suggestedStrategy"`
	Reason            string   `json:"reason"`
	ExpectedBenefit   string   `json:"expectedBenefit"`
	PreferredNodes    []string `json:"preferredNodes"`
}

// MonitorGPUHealth 监控GPU健康状况
func (gm *GPUManager) MonitorGPUHealth() (*GPUHealthReport, error) {
	gpuResources, err := gm.GetClusterGPUResources()
	if err != nil {
		return nil, err
	}

	healthReport := &GPUHealthReport{
		Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
		OverallHealth: "Healthy",
		NodeHealth:    make([]NodeHealthStatus, 0),
		Issues:        make([]HealthIssue, 0),
		Statistics:    HealthStatistics{},
	}

	healthyNodes := 0
	totalNodes := len(gpuResources)
	totalIssues := 0

	for _, gpuInfo := range gpuResources {
		nodeHealth := NodeHealthStatus{
			NodeName:    gpuInfo.NodeName,
			Status:      gpuInfo.Status,
			TotalGPUs:   gpuInfo.TotalGPUs,
			HealthyGPUs: gpuInfo.TotalGPUs, // 简化：假设所有GPU都是健康的
			Issues:      make([]string, 0),
		}

		// 检查节点健康状况
		if gpuInfo.Status == "Ready" {
			healthyNodes++
		} else {
			issue := HealthIssue{
				Severity:   "high",
				NodeName:   gpuInfo.NodeName,
				Component:  "node",
				Issue:      fmt.Sprintf("节点状态异常: %s", gpuInfo.Status),
				Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
				Suggestion: "检查节点状态和kubelet日志",
			}
			healthReport.Issues = append(healthReport.Issues, issue)
			nodeHealth.Issues = append(nodeHealth.Issues, issue.Issue)
			totalIssues++
		}

		// TODO: 检查GPU温度、功耗等健康指标

		healthReport.NodeHealth = append(healthReport.NodeHealth, nodeHealth)
	}

	// 计算整体健康状况
	if totalIssues == 0 {
		healthReport.OverallHealth = "Healthy"
	} else if float64(healthyNodes)/float64(totalNodes) > 0.8 {
		healthReport.OverallHealth = "Warning"
	} else {
		healthReport.OverallHealth = "Critical"
	}

	healthReport.Statistics = HealthStatistics{
		TotalNodes:   int32(totalNodes),
		HealthyNodes: int32(healthyNodes),
		TotalIssues:  int32(totalIssues),
		HealthRate:   float64(healthyNodes) / float64(totalNodes) * 100,
	}

	return healthReport, nil
}

// GPUHealthReport GPU健康报告
type GPUHealthReport struct {
	Timestamp     string             `json:"timestamp"`
	OverallHealth string             `json:"overallHealth"`
	NodeHealth    []NodeHealthStatus `json:"nodeHealth"`
	Issues        []HealthIssue      `json:"issues"`
	Statistics    HealthStatistics   `json:"statistics"`
}

// NodeHealthStatus 节点健康状态
type NodeHealthStatus struct {
	NodeName    string   `json:"nodeName"`
	Status      string   `json:"status"`
	TotalGPUs   int32    `json:"totalGPUs"`
	HealthyGPUs int32    `json:"healthyGPUs"`
	Issues      []string `json:"issues"`
}

// HealthIssue 健康问题
type HealthIssue struct {
	Severity   string `json:"severity"` // low, medium, high, critical
	NodeName   string `json:"nodeName"`
	Component  string `json:"component"` // node, gpu, driver
	Issue      string `json:"issue"`
	Timestamp  string `json:"timestamp"`
	Suggestion string `json:"suggestion"`
}

// HealthStatistics 健康统计
type HealthStatistics struct {
	TotalNodes   int32   `json:"totalNodes"`
	HealthyNodes int32   `json:"healthyNodes"`
	TotalIssues  int32   `json:"totalIssues"`
	HealthRate   float64 `json:"healthRate"`
}
