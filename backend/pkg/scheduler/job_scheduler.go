package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"api/model"
	"api/pkg/k8s"
	"github.com/zeromicro/go-zero/core/logx"
	corev1 "k8s.io/api/core/v1"
)

// JobScheduler 训练作业调度器
type JobScheduler struct {
	k8sClient *k8s.Client
	namespace string
	logger    logx.Logger
}

// NewJobScheduler 创建作业调度器
func NewJobScheduler(k8sClient *k8s.Client, namespace string) *JobScheduler {
	return &JobScheduler{
		k8sClient: k8sClient,
		namespace: namespace,
		logger:    logx.WithContext(context.Background()),
	}
}

// ScheduleJob 调度训练作业到Kubernetes
func (s *JobScheduler) ScheduleJob(job *model.VtTrainingJobs) error {
	// 构建Kubernetes作业规格
	jobSpec, err := s.buildJobSpec(job)
	if err != nil {
		return fmt.Errorf("构建作业规格失败: %v", err)
	}

	// 创建Kubernetes作业
	k8sJob, err := s.k8sClient.CreateJob(jobSpec)
	if err != nil {
		return fmt.Errorf("创建Kubernetes作业失败: %v", err)
	}

	s.logger.Infof("成功创建Kubernetes作业: %s, UID: %s", k8sJob.Name, k8sJob.UID)
	return nil
}

// buildJobSpec 构建Kubernetes作业规格
func (s *JobScheduler) buildJobSpec(job *model.VtTrainingJobs) (*k8s.JobSpec, error) {
	spec := &k8s.JobSpec{
		Name:          s.generateK8sJobName(job.Name),
		Image:         job.Image,
		QueueName:     job.QueueName,
		Priority:      int32(job.Priority),
		Replicas:      int32(job.WorkerCount),
		RestartPolicy: "Never",
		MaxRetryCount: int32(job.MaxRetryCount),
		WorkingDir:    job.WorkingDir,
	}

	// 设置命令和参数
	if err := s.setCommandAndArgs(spec, job); err != nil {
		return nil, err
	}

	// 设置环境变量
	if err := s.setEnvironmentVariables(spec, job); err != nil {
		return nil, err
	}

	// 设置资源需求
	if err := s.setResourceRequirements(spec, job); err != nil {
		return nil, err
	}

	// 设置节点选择和容忍度
	if err := s.setNodeSelectorAndTolerations(spec, job); err != nil {
		return nil, err
	}

	// 设置存储卷
	if err := s.setVolumes(spec, job); err != nil {
		return nil, err
	}

	return spec, nil
}

// generateK8sJobName 生成Kubernetes作业名称
func (s *JobScheduler) generateK8sJobName(jobName string) string {
	// K8s资源名称需要符合DNS-1123规范
	k8sName := strings.ToLower(jobName)
	k8sName = strings.ReplaceAll(k8sName, "_", "-")

	// 添加时间戳确保唯一性
	timestamp := time.Now().Format("20060102-150405")
	return fmt.Sprintf("%s-%s", k8sName, timestamp)
}

// setCommandAndArgs 设置命令和参数
func (s *JobScheduler) setCommandAndArgs(spec *k8s.JobSpec, job *model.VtTrainingJobs) error {
	// 基础命令
	spec.Command = []string{"python", "-u"}

	// 设置程序入口
	spec.Args = []string{job.EntryPoint}

	// 解析命令参数
	if job.CommandArgs != "" {
		var commandArgs []string
		if err := json.Unmarshal([]byte(job.CommandArgs), &commandArgs); err != nil {
			s.logger.Errorf("解析命令参数失败: %v", err)
			return err
		}
		spec.Args = append(spec.Args, commandArgs...)
	}

	return nil
}

// setEnvironmentVariables 设置环境变量
func (s *JobScheduler) setEnvironmentVariables(spec *k8s.JobSpec, job *model.VtTrainingJobs) error {
	spec.Env = make(map[string]string)

	// 设置基础环境变量
	spec.Env["JOB_NAME"] = job.Name
	spec.Env["JOB_TYPE"] = job.JobType
	spec.Env["FRAMEWORK"] = job.Framework
	spec.Env["FRAMEWORK_VERSION"] = job.FrameworkVersion
	spec.Env["PYTHON_VERSION"] = job.PythonVersion
	spec.Env["QUEUE_NAME"] = job.QueueName

	// 设置分布式训练环境变量
	if job.WorkerCount > 1 {
		spec.Env["WORLD_SIZE"] = strconv.Itoa(job.WorkerCount)
		spec.Env["NPROC_PER_NODE"] = strconv.Itoa(job.GpuCount)
	}

	// 解析用户自定义环境变量
	if job.EnvVars != "" {
		var envVars map[string]string
		if err := json.Unmarshal([]byte(job.EnvVars), &envVars); err != nil {
			s.logger.Errorf("解析环境变量失败: %v", err)
			return err
		}

		for key, value := range envVars {
			spec.Env[key] = value
		}
	}

	return nil
}

// setResourceRequirements 设置资源需求
func (s *JobScheduler) setResourceRequirements(spec *k8s.JobSpec, job *model.VtTrainingJobs) error {
	spec.Resources = k8s.ResourceRequirements{
		CPURequests:    job.CpuCores,
		CPULimits:      job.CpuCores,
		MemoryRequests: job.MemoryGb + "Gi",
		MemoryLimits:   job.MemoryGb + "Gi",
		GPURequests:    job.GpuCount,
		GPULimits:      job.GpuCount,
	}

	// 设置存储需求
	if job.StorageGb != "" {
		spec.Resources.StorageRequests = job.StorageGb + "Gi"
	}

	return nil
}

// setNodeSelectorAndTolerations 设置节点选择器和容忍度
func (s *JobScheduler) setNodeSelectorAndTolerations(spec *k8s.JobSpec, job *model.VtTrainingJobs) error {
	// 解析节点选择器
	if job.NodeSelector != "" {
		var nodeSelector map[string]string
		if err := json.Unmarshal([]byte(job.NodeSelector), &nodeSelector); err != nil {
			s.logger.Errorf("解析节点选择器失败: %v", err)
			return err
		}
		spec.NodeSelector = nodeSelector
	}

	// 解析容忍度
	if job.Tolerations != "" {
		var tolerations []corev1.Toleration
		if err := json.Unmarshal([]byte(job.Tolerations), &tolerations); err != nil {
			s.logger.Errorf("解析容忍度失败: %v", err)
			return err
		}
		spec.Tolerations = tolerations
	}

	// 添加GPU节点容忍度
	if job.GpuCount > 0 {
		gpuToleration := corev1.Toleration{
			Key:      "nvidia.com/gpu",
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoSchedule,
		}
		spec.Tolerations = append(spec.Tolerations, gpuToleration)

		// 添加GPU节点选择器
		if spec.NodeSelector == nil {
			spec.NodeSelector = make(map[string]string)
		}
		spec.NodeSelector["accelerator"] = "nvidia-tesla-v100" // 可配置
	}

	return nil
}

// setVolumes 设置存储卷
func (s *JobScheduler) setVolumes(spec *k8s.JobSpec, job *model.VtTrainingJobs) error {
	// 工作目录卷
	workspaceVolume := k8s.VolumeSpec{
		Name:      "workspace",
		Type:      "pvc",
		MountPath: "/workspace",
		ReadOnly:  false,
		Source: map[string]string{
			"claimName": fmt.Sprintf("workspace-%s", job.Name),
		},
	}
	spec.Volumes = append(spec.Volumes, workspaceVolume)

	// 数据集卷
	if err := s.addDatasetVolumes(spec, job); err != nil {
		return err
	}

	// 输出卷
	if job.OutputPath != "" {
		outputVolume := k8s.VolumeSpec{
			Name:      "output",
			Type:      "pvc",
			MountPath: "/output",
			ReadOnly:  false,
			Source: map[string]string{
				"claimName": fmt.Sprintf("output-%s", job.Name),
			},
		}
		spec.Volumes = append(spec.Volumes, outputVolume)
	}

	// 日志卷
	logVolume := k8s.VolumeSpec{
		Name:      "logs",
		Type:      "emptyDir",
		MountPath: "/logs",
		ReadOnly:  false,
	}
	spec.Volumes = append(spec.Volumes, logVolume)

	return nil
}

// addDatasetVolumes 添加数据集卷
func (s *JobScheduler) addDatasetVolumes(spec *k8s.JobSpec, job *model.VtTrainingJobs) error {
	if job.DatasetMountConfigs == "" {
		return nil
	}

	var mountConfigs []map[string]interface{}
	if err := json.Unmarshal([]byte(job.DatasetMountConfigs), &mountConfigs); err != nil {
		s.logger.Errorf("解析数据集挂载配置失败: %v", err)
		return err
	}

	for i, config := range mountConfigs {
		datasetName := fmt.Sprintf("dataset-%d", i)
		mountPath := "/data"

		if path, ok := config["mount_path"].(string); ok {
			mountPath = path
		}

		volume := k8s.VolumeSpec{
			Name:      datasetName,
			Type:      "pvc",
			MountPath: mountPath,
			ReadOnly:  true,
			Source: map[string]string{
				"claimName": fmt.Sprintf("dataset-%s", config["dataset_name"].(string)),
			},
		}
		spec.Volumes = append(spec.Volumes, volume)
	}

	return nil
}

// GetJobStatus 获取作业状态
func (s *JobScheduler) GetJobStatus(k8sJobName string) (*k8s.JobStatus, error) {
	return s.k8sClient.GetJobStatus(k8sJobName)
}

// StopJob 停止作业
func (s *JobScheduler) StopJob(k8sJobName string) error {
	return s.k8sClient.DeleteJob(k8sJobName)
}

// UpdateJobStatus 更新作业状态
func (s *JobScheduler) UpdateJobStatus(jobModel model.VtTrainingJobsModel, jobID int64, k8sJobName string) error {
	status, err := s.GetJobStatus(k8sJobName)
	if err != nil {
		s.logger.Errorf("获取K8s作业状态失败: %v", err)
		return err
	}

	// 映射K8s状态到业务状态
	var businessStatus, businessPhase string
	switch status.Phase {
	case "Pending":
		businessStatus = "pending"
		businessPhase = "waiting"
	case "Running":
		businessStatus = "running"
		businessPhase = "training"
	case "Succeeded":
		businessStatus = "completed"
		businessPhase = "succeeded"
	case "Failed":
		businessStatus = "failed"
		businessPhase = "failed"
	default:
		businessStatus = "unknown"
		businessPhase = "unknown"
	}

	// 更新数据库中的作业状态
	return jobModel.UpdateStatus(jobID, businessStatus, businessPhase)
}

// SyncJobStatus 同步所有作业状态
func (s *JobScheduler) SyncJobStatus(jobModel model.VtTrainingJobsModel) error {
	// 获取所有运行中的作业
	runningJobs, err := jobModel.GetByStatus("running")
	if err != nil {
		return fmt.Errorf("获取运行中作业失败: %v", err)
	}

	// 逐个同步状态
	for _, job := range runningJobs {
		if job.VolcanoJobName != "" {
			if err := s.UpdateJobStatus(jobModel, job.Id, job.VolcanoJobName); err != nil {
				s.logger.Errorf("同步作业状态失败: job=%s, error=%v", job.Name, err)
			}
		}
	}

	return nil
}
