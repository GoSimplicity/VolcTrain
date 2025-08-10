package main

import (
	"fmt"
	"log"
	"time"

	"api/pkg/volcano"

	corev1 "k8s.io/api/core/v1"
)

// 这是一个展示如何使用基于Volcano的AI训练平台的综合示例
func main() {
	// 1. 创建Volcano客户端
	client, err := volcano.NewClient("", "volcano-training")
	if err != nil {
		log.Fatalf("创建Volcano客户端失败: %v", err)
	}

	// 2. 创建GPU管理器
	gpuManager := volcano.NewGPUManager(client)

	// 3. 创建作业管理器
	jobManager := volcano.NewJobManager(client)

	// 4. 创建监控服务
	monitoringConfig := &volcano.MonitoringConfig{
		UpdateInterval:     30 * time.Second,
		EnableGPUMonitor:   true,
		EnableJobMonitor:   true,
		EnableQueueMonitor: true,
		AlertThresholds: volcano.AlertThresholds{
			GPUUtilization:    80.0,
			MemoryUtilization: 85.0,
			QueueWaitTime:     300, // 5分钟
			JobFailureRate:    10.0,
			NodeUnhealthyRate: 0.2,
		},
	}

	monitoringService := volcano.NewMonitoringService(client, monitoringConfig)

	// 启动监控服务
	if err := monitoringService.Start(); err != nil {
		log.Printf("启动监控服务失败: %v", err)
	}

	// 示例1：创建训练队列
	fmt.Println("=== 创建训练队列 ===")
	queueSpec := &volcano.QueueSpec{
		Name:               "pytorch-training",
		Weight:             10,
		Capability:         make(map[string]interface{}),
		ShareWeight:        floatPtr(1.0),
		GuaranteedResource: volcano.BuildResourceList("10", "50Gi", "4"),
		MaxResource:        volcano.BuildResourceList("50", "200Gi", "16"),
		State:              "Open",
		Labels: map[string]string{
			"purpose": "ml-training",
			"team":    "ai-research",
		},
	}

	queue, err := client.CreateQueue(queueSpec)
	if err != nil {
		log.Printf("创建队列失败: %v", err)
	} else {
		fmt.Printf("成功创建队列: %s\n", queue.Name)
	}

	// 示例2：获取集群GPU资源信息
	fmt.Println("\n=== 获取集群GPU资源 ===")
	gpuResources, err := gpuManager.GetClusterGPUResources()
	if err != nil {
		log.Printf("获取GPU资源失败: %v", err)
	} else {
		for _, gpuInfo := range gpuResources {
			fmt.Printf("节点: %s, GPU类型: %s, 总数: %d, 可用: %d\n",
				gpuInfo.NodeName, gpuInfo.GPUType, gpuInfo.TotalGPUs, gpuInfo.AvailableGPUs)
		}
	}

	// 示例3：创建PyTorch分布式训练作业
	fmt.Println("\n=== 创建PyTorch分布式训练作业 ===")
	trainingJobSpec := &volcano.TrainingJobSpec{
		Name:             "pytorch-bert-training",
		Namespace:        "volcano-training",
		JobType:          "pytorch",
		Framework:        "pytorch",
		FrameworkVersion: "1.12.0",
		Image:            "pytorch/pytorch:1.12.0-cuda11.3-cudnn8-devel",
		Command:          []string{"python"},
		Args:             []string{"/workspace/train.py", "--distributed"},
		WorkingDir:       "/workspace",

		// 资源配置
		CPURequest:     "4",
		MemoryRequest:  "16Gi",
		GPUCount:       2,
		GPUType:        "V100",
		StorageRequest: "100Gi",

		// 分布式训练配置
		MasterReplicas: 1,
		WorkerReplicas: 4,
		MinAvailable:   3, // Gang调度：至少需要3个Pod才开始调度

		// 调度配置
		QueueName:         "pytorch-training",
		Priority:          100,
		PriorityClassName: "high-priority",
		NodeSelector: map[string]string{
			"node-type": "gpu",
		},
		Tolerations: []corev1.Toleration{
			{
				Key:      "nvidia.com/gpu",
				Operator: corev1.TolerationOpEqual,
				Value:    "present",
				Effect:   corev1.TaintEffectNoSchedule,
			},
		},

		// 生命周期配置
		MaxRetry:         3,
		TTLAfterFinished: int32Ptr(3600),  // 1小时后清理
		ActiveDeadline:   int64Ptr(86400), // 24小时超时

		// 环境配置
		EnvVars: map[string]string{
			"NCCL_DEBUG": "INFO",
			"PYTHONPATH": "/workspace",
		},

		// 存储卷配置
		VolumeMounts: []volcano.VolumeMountSpec{
			{
				Name:      "training-data",
				MountPath: "/data",
				ReadOnly:  true,
				VolumeSource: volcano.VolumeSource{
					Type:    "pvc",
					PVCName: "training-dataset-pvc",
				},
			},
			{
				Name:      "model-output",
				MountPath: "/output",
				ReadOnly:  false,
				VolumeSource: volcano.VolumeSource{
					Type:    "pvc",
					PVCName: "model-output-pvc",
				},
			},
		},

		// 监控配置
		EnableTensorboard: true,
		EnableProfiling:   false,
		MetricsPort:       9090,

		// 元数据
		Labels: map[string]string{
			"model":      "bert",
			"experiment": "exp-001",
			"team":       "nlp",
		},
		Annotations: map[string]string{
			"description": "BERT模型分布式训练实验",
			"contact":     "ml-team@example.com",
		},
	}

	job, err := jobManager.CreateTrainingJob(trainingJobSpec)
	if err != nil {
		log.Printf("创建训练作业失败: %v", err)
	} else {
		fmt.Printf("成功创建训练作业: %s\n", job.Name)
	}

	// 示例4：GPU智能分配
	fmt.Println("\n=== GPU智能分配 ===")
	gpuRequest := &volcano.GPUAllocationRequest{
		JobName:      "pytorch-bert-training",
		TaskName:     "worker",
		GPUCount:     8,
		GPUMemoryMin: "16Gi",
		Strategy: volcano.GPUAllocationStrategy{
			Strategy:     "topology", // 拓扑感知分配
			GPUTypes:     []string{"V100", "A100"},
			TopologyKey:  "kubernetes.io/hostname",
			MaxSkew:      2,
			AntiAffinity: false,
		},
		Priority: 100,
		Queue:    "pytorch-training",
	}

	allocationResult, err := gpuManager.AllocateGPUs(gpuRequest)
	if err != nil {
		log.Printf("GPU分配失败: %v", err)
	} else if allocationResult.Success {
		fmt.Printf("成功分配%d个GPU:\n", len(allocationResult.AllocatedGPUs))
		for _, gpu := range allocationResult.AllocatedGPUs {
			fmt.Printf("  - 节点: %s, GPU: %s (%s)\n", gpu.NodeName, gpu.GPUUUID, gpu.GPUType)
		}
	} else {
		fmt.Printf("GPU分配失败: %s\n", allocationResult.Message)
		for _, suggestion := range allocationResult.Suggestions {
			fmt.Printf("  建议: %s\n", suggestion)
		}
	}

	// 示例5：监控作业状态
	fmt.Println("\n=== 监控作业状态 ===")
	go func() {
		for {
			time.Sleep(10 * time.Second)

			status, err := client.GetJobStatus("volcano-training", "pytorch-bert-training")
			if err != nil {
				log.Printf("获取作业状态失败: %v", err)
				continue
			}

			fmt.Printf("作业状态: %s, 阶段: %s, 运行中: %d, 成功: %d, 失败: %d\n",
				status.State, status.Phase, status.Running, status.Succeeded, status.Failed)

			// 检查作业是否完成
			if status.Phase == "Succeeded" || status.Phase == "Failed" {
				fmt.Printf("作业已完成，最终状态: %s\n", status.Phase)
				break
			}
		}
	}()

	// 示例6：扩缩容作业
	fmt.Println("\n=== 作业扩缩容 ===")
	time.Sleep(30 * time.Second) // 等待作业启动

	scaleSpecs := []volcano.TaskScaleSpec{
		{
			Name:     "worker",
			Replicas: 6, // 从4个worker扩展到6个
		},
	}

	err = client.ScaleJob("volcano-training", "pytorch-bert-training", scaleSpecs)
	if err != nil {
		log.Printf("扩缩容失败: %v", err)
	} else {
		fmt.Println("成功扩缩容作业")
	}

	// 示例7：获取GPU使用率报告
	fmt.Println("\n=== GPU使用率报告 ===")
	report, err := gpuManager.GetGPUUtilizationReport("1h")
	if err != nil {
		log.Printf("获取GPU使用率报告失败: %v", err)
	} else {
		fmt.Printf("集群GPU摘要:\n")
		fmt.Printf("  总GPU数: %d\n", report.ClusterSummary.TotalGPUs)
		fmt.Printf("  已分配: %d\n", report.ClusterSummary.AllocatedGPUs)
		fmt.Printf("  可用: %d\n", report.ClusterSummary.AvailableGPUs)
		fmt.Printf("  平均使用率: %.2f%%\n", report.ClusterSummary.AvgUtilization)
		fmt.Printf("  分配率: %.2f%%\n", report.ClusterSummary.AllocationRate)

		fmt.Println("节点详细信息:")
		for _, nodeReport := range report.NodeReports {
			fmt.Printf("  节点 %s: GPU类型=%s, 总数=%d, 已分配=%d, 使用率=%.2f%%\n",
				nodeReport.NodeName, nodeReport.GPUType, nodeReport.TotalGPUs,
				nodeReport.AllocatedGPUs, nodeReport.Utilization)
		}
	}

	// 示例8：获取集群调度信息
	fmt.Println("\n=== 集群调度信息 ===")
	schedulingInfo, err := client.GetClusterSchedulingInfo("main-cluster")
	if err != nil {
		log.Printf("获取集群调度信息失败: %v", err)
	} else {
		fmt.Printf("集群: %s, 调度器: %s\n", schedulingInfo.ClusterName, schedulingInfo.SchedulerName)
		fmt.Printf("总资源: %v\n", schedulingInfo.TotalResources)
		fmt.Printf("已分配: %v\n", schedulingInfo.AllocatedResources)
		fmt.Printf("可用: %v\n", schedulingInfo.AvailableResources)
		fmt.Printf("待调度作业: %d, 运行中作业: %d\n", schedulingInfo.PendingJobs, schedulingInfo.RunningJobs)

		fmt.Println("队列信息:")
		for _, queueInfo := range schedulingInfo.Queues {
			fmt.Printf("  队列 %s: 权重=%d, 共享=%.2f%%, 待调度=%d, 运行中=%d\n",
				queueInfo.Name, queueInfo.Weight, queueInfo.Share*100,
				queueInfo.PendingJobs, queueInfo.RunningJobs)
		}
	}

	// 示例9：健康检查和告警
	fmt.Println("\n=== 健康检查 ===")
	healthStatus, err := monitoringService.GetHealthStatus()
	if err != nil {
		log.Printf("获取健康状态失败: %v", err)
	} else {
		fmt.Printf("整体健康状况: %s\n", healthStatus.Overall)
		for component, health := range healthStatus.Components {
			fmt.Printf("  组件 %s: 状态=%s\n", component, health.Status)
			for _, issue := range health.Issues {
				fmt.Printf("    问题: %s\n", issue)
			}
		}
	}

	// 示例10：优化建议
	fmt.Println("\n=== GPU放置优化建议 ===")
	optimizationReq := []volcano.GPUAllocationRequest{*gpuRequest}
	optimization, err := gpuManager.OptimizeGPUPlacement(optimizationReq)
	if err != nil {
		log.Printf("获取优化建议失败: %v", err)
	} else {
		fmt.Println("资源告警:")
		for _, alert := range optimization.ResourceAlerts {
			fmt.Printf("  %s: %s (使用率: %.2f%%)\n", alert.Type, alert.Message, alert.UsageRate*100)
		}

		fmt.Println("优化建议:")
		for _, rec := range optimization.Recommendations {
			fmt.Printf("  - %s\n", rec)
		}
	}

	// 等待一段时间以观察监控数据
	fmt.Println("\n=== 等待监控数据收集 ===")
	time.Sleep(2 * time.Minute)

	// 获取集群指标
	metrics, err := monitoringService.GetClusterMetrics()
	if err != nil {
		log.Printf("获取集群指标失败: %v", err)
	} else {
		fmt.Printf("集群指标更新时间: %s\n", metrics.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("GPU指标: 总数=%d, 已分配=%d, 平均使用率=%.2f%%\n",
			metrics.GPUMetrics.TotalGPUs, metrics.GPUMetrics.AllocatedGPUs, metrics.GPUMetrics.AvgUtilization)
		fmt.Printf("作业指标: 总数=%d, 运行中=%d, 成功率=%.2f%%\n",
			metrics.JobMetrics.TotalJobs, metrics.JobMetrics.RunningJobs, metrics.JobMetrics.SuccessRate)
		fmt.Printf("告警摘要: 严重=%d, 警告=%d, 信息=%d\n",
			metrics.AlertSummary.CriticalAlerts, metrics.AlertSummary.WarningAlerts, metrics.AlertSummary.InfoAlerts)
	}

	fmt.Println("\n=== 示例完成 ===")
	fmt.Println("基于Volcano的AI训练平台功能演示已完成")
	fmt.Println("包含的主要功能:")
	fmt.Println("✓ 队列管理和资源配额")
	fmt.Println("✓ 分布式训练作业创建")
	fmt.Println("✓ GPU智能分配和拓扑感知调度")
	fmt.Println("✓ 作业生命周期管理（扩缩容、暂停、恢复）")
	fmt.Println("✓ 实时监控和健康检查")
	fmt.Println("✓ 告警和优化建议")
	fmt.Println("✓ Gang调度和抢占式调度")

	// 停止监控服务
	monitoringService.Stop()
}

// 辅助函数
func int32Ptr(i int32) *int32 {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}

// 演示如何使用不同框架的作业
func createTensorFlowJob() *volcano.TrainingJobSpec {
	return &volcano.TrainingJobSpec{
		Name:             "tensorflow-resnet-training",
		Namespace:        "volcano-training",
		JobType:          "tensorflow",
		Framework:        "tensorflow",
		FrameworkVersion: "2.8.0",
		Image:            "tensorflow/tensorflow:2.8.0-gpu",
		Command:          []string{"python"},
		Args:             []string{"/workspace/train_resnet.py"},
		WorkingDir:       "/workspace",

		CPURequest:    "8",
		MemoryRequest: "32Gi",
		GPUCount:      4,

		MasterReplicas: 1,
		WorkerReplicas: 3,
		PSReplicas:     2, // TensorFlow Parameter Server
		MinAvailable:   4, // 需要所有Pod都就绪

		QueueName: "tensorflow-training",
		Priority:  50,

		EnvVars: map[string]string{
			"TF_CONFIG": "", // 将在运行时动态生成
		},

		Labels: map[string]string{
			"framework": "tensorflow",
			"model":     "resnet50",
		},
	}
}

func createMPIJob() *volcano.TrainingJobSpec {
	return &volcano.TrainingJobSpec{
		Name:             "mpi-horovod-training",
		Namespace:        "volcano-training",
		JobType:          "mpi",
		Framework:        "horovod",
		FrameworkVersion: "0.24.0",
		Image:            "horovod/horovod:0.24.0-tf2.6.0-torch1.9.0-mxnet1.8.0-py3.8-gpu",
		Command:          []string{"mpirun"},
		Args:             []string{"-np", "8", "python", "/workspace/train_horovod.py"},
		WorkingDir:       "/workspace",

		CPURequest:    "4",
		MemoryRequest: "16Gi",
		GPUCount:      1, // 每个worker使用1个GPU

		MasterReplicas: 1, // MPI master
		WorkerReplicas: 8, // 8个MPI worker
		MinAvailable:   9, // master + workers

		QueueName: "mpi-training",
		Priority:  75,

		EnvVars: map[string]string{
			"OMPI_ALLOW_RUN_AS_ROOT":         "1",
			"OMPI_ALLOW_RUN_AS_ROOT_CONFIRM": "1",
			"HOROVOD_CUDA_HOME":              "/usr/local/cuda",
		},

		Labels: map[string]string{
			"framework": "horovod",
			"protocol":  "mpi",
		},
	}
}
