package training

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"api/internal/svc"
	"api/internal/types"
	"api/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTrainingJobLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建训练作业
func NewCreateTrainingJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTrainingJobLogic {
	return &CreateTrainingJobLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTrainingJobLogic) CreateTrainingJob(req *types.CreateTrainingJobReq) (resp *types.CreateTrainingJobResp, err error) {
	// 验证请求参数
	if err = l.validateRequest(req); err != nil {
		l.Logger.Errorf("训练作业请求参数验证失败: %v", err)
		return nil, fmt.Errorf("请求参数验证失败: %w", err)
	}

	// 检查训练作业名称是否已存在
	exists, err := l.checkJobNameExists(req.Name)
	if err != nil {
		l.Logger.Errorf("检查训练作业名称失败: %v", err)
		return nil, fmt.Errorf("检查训练作业名称失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("训练作业名称 '%s' 已存在", req.Name)
	}

	// 创建训练作业模型
	trainingJob := &model.VtTrainingJobs{
		Name:                req.Name,
		DisplayName:         req.DisplayName,
		Description:         req.Description,
		JobType:             req.JobType,
		Framework:           req.Framework,
		FrameworkVersion:    req.FrameworkVersion,
		PythonVersion:       req.PythonVersion,
		CodeSourceType:      req.CodeSourceType,
		CodeSourceConfig:    req.CodeSourceConfig,
		EntryPoint:          req.EntryPoint,
		WorkingDir:          req.WorkingDir,
		Image:               req.Image,
		ImagePullPolicy:     req.ImagePullPolicy,
		ImagePullSecrets:    req.ImagePullSecrets,
		DatasetMountConfigs: req.DatasetMountConfigs,
		DataSourceConfig:    req.DataSourceConfig,
		ModelConfig:         req.ModelConfig,
		OutputModelName:     req.OutputModelName,
		ModelSaveStrategy:   req.ModelSaveStrategy,
		CpuCores:            req.CpuCores,
		MemoryGb:            req.MemoryGb,
		GpuCount:            int(req.GpuCount),
		GpuType:             req.GpuType,
		GpuMemoryGb:         req.GpuMemoryGb,
		StorageGb:           req.StorageGb,
		SharedMemoryGb:      req.SharedMemoryGb,
		WorkerCount:         int(req.WorkerCount),
		PsCount:             int(req.PsCount),
		MasterCount:         int(req.MasterCount),
		EnvVars:             req.EnvVars,
		CommandArgs:         req.CommandArgs,
		Secrets:             req.Secrets,
		ConfigMaps:          req.ConfigMaps,
		VolumeMounts:        req.VolumeMounts,
		QueueName:           req.QueueName,
		Priority:            int(req.Priority),
		NodeSelector:        req.NodeSelector,
		Tolerations:         req.Tolerations,
		Affinity:            req.Affinity,
		MaxRuntimeSeconds:   int(req.MaxRuntimeSeconds),
		Status:              "pending",
		SubmittedAt:         time.Now(),
	}

	// 保存到数据库
	result, err := l.svcCtx.VtTrainingJobsModel.Insert(trainingJob)
	if err != nil {
		l.Logger.Errorf("保存训练作业失败: %v", err)
		return nil, fmt.Errorf("保存训练作业失败: %w", err)
	}

	// 获取插入的ID
	jobID, err := result.LastInsertId()
	if err != nil {
		l.Logger.Errorf("获取训练作业ID失败: %v", err)
		return nil, fmt.Errorf("获取训练作业ID失败: %w", err)
	}

	// 如果需要GPU资源，尝试预分配
	if req.GpuCount > 0 {
		if err = l.preAllocateGPUResources(jobID, req); err != nil {
			l.Logger.Infof("GPU资源预分配失败，作业将等待资源: %v", err)
			// 不返回错误，让作业进入等待状态
		}
	}

	// 提交作业到调度队列
	if err = l.submitToScheduleQueue(jobID, req); err != nil {
		l.Logger.Errorf("提交作业到调度队列失败: %v", err)
		// 不返回错误，作业已创建，可以后续重新调度
	}

	l.Logger.Infof("训练作业创建成功: ID=%d, Name=%s", jobID, req.Name)

	return &types.CreateTrainingJobResp{
		Id: jobID,
	}, nil
}

// validateRequest 验证请求参数
func (l *CreateTrainingJobLogic) validateRequest(req *types.CreateTrainingJobReq) error {
	if req.Name == "" {
		return fmt.Errorf("训练作业名称不能为空")
	}

	if req.Framework == "" {
		return fmt.Errorf("训练框架不能为空")
	}

	if req.Image == "" {
		return fmt.Errorf("训练镜像不能为空")
	}

	if req.EntryPoint == "" {
		return fmt.Errorf("入口点不能为空")
	}

	// 验证GPU配置
	if req.GpuCount > 0 && req.GpuType == "" {
		return fmt.Errorf("指定GPU数量时必须指定GPU类型")
	}

	// 验证分布式训练配置
	if req.JobType == "distributed" {
		if req.WorkerCount <= 0 {
			return fmt.Errorf("分布式训练必须指定工作节点数量")
		}
	}

	// 验证资源限制
	if req.MaxRuntimeSeconds < 0 {
		return fmt.Errorf("最大运行时间不能为负数")
	}

	return nil
}

// checkJobNameExists 检查训练作业名称是否已存在
func (l *CreateTrainingJobLogic) checkJobNameExists(name string) (bool, error) {
	// 使用模型检查名称是否存在
	job, err := l.svcCtx.VtTrainingJobsModel.FindOneByName(name)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return job != nil, nil
}

// getUserIDFromContext 从上下文获取用户ID
func (l *CreateTrainingJobLogic) getUserIDFromContext() int64 {
	// 简化处理，实际应该从JWT token或上下文中获取
	// 这里返回一个默认值用于测试
	return 1001
}

// preAllocateGPUResources 预分配GPU资源
func (l *CreateTrainingJobLogic) preAllocateGPUResources(jobID int64, req *types.CreateTrainingJobReq) error {
	// 简化实现，实际应该调用GPU管理服务
	l.Logger.Infof("尝试为作业 %d 预分配 %d 个 %s GPU", jobID, req.GpuCount, req.GpuType)

	// 这里可以调用GPU分配逻辑
	// allocateReq := &types.AllocateGpuDeviceReq{
	//     JobId: jobID,
	//     GpuCount: req.GpuCount,
	//     GpuType: req.GpuType,
	// }
	// return l.svcCtx.GpuAllocator.Allocate(allocateReq)

	return nil
}

// submitToScheduleQueue 提交作业到调度队列
func (l *CreateTrainingJobLogic) submitToScheduleQueue(jobID int64, req *types.CreateTrainingJobReq) error {
	// 简化实现，实际应该调用调度服务
	l.Logger.Infof("提交作业 %d 到调度队列 %s", jobID, req.QueueName)

	// 创建调度任务
	scheduleTask := map[string]interface{}{
		"job_id":     jobID,
		"queue_name": req.QueueName,
		"priority":   req.Priority,
		"resources": map[string]interface{}{
			"cpu_cores":  req.CpuCores,
			"memory_gb":  req.MemoryGb,
			"gpu_count":  req.GpuCount,
			"gpu_type":   req.GpuType,
			"storage_gb": req.StorageGb,
		},
		"created_at": time.Now().Unix(),
	}

	// 序列化任务信息
	taskData, err := json.Marshal(scheduleTask)
	if err != nil {
		return fmt.Errorf("序列化调度任务失败: %w", err)
	}

	// 实际环境中应该发送到消息队列或调度服务
	l.Logger.Infof("调度任务数据: %s", string(taskData))

	return nil
}
