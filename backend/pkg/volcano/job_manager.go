package volcano

import (
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	vcjob "volcano.sh/apis/pkg/apis/batch/v1alpha1"
)

// JobManager Volcano作业生命周期管理器
type JobManager struct {
	client *Client
}

// NewJobManager 创建作业管理器
func NewJobManager(client *Client) *JobManager {
	return &JobManager{
		client: client,
	}
}

// CreateTrainingJob 创建训练作业
func (jm *JobManager) CreateTrainingJob(spec *TrainingJobSpec) (*vcjob.Job, error) {
	// 验证作业规格
	if err := jm.validateJobSpec(spec); err != nil {
		return nil, fmt.Errorf("作业规格验证失败: %v", err)
	}

	// 构建Volcano作业规格
	volcanoJobSpec := jm.buildVolcanoJobSpec(spec)

	// 创建作业
	job, err := jm.client.CreateJob(volcanoJobSpec)
	if err != nil {
		return nil, fmt.Errorf("创建Volcano训练作业失败: %v", err)
	}

	return job, nil
}

// TrainingJobSpec 训练作业规格定义
type TrainingJobSpec struct {
	Name             string
	Namespace        string
	JobType          string // pytorch, tensorflow, mpi, paddle
	Framework        string
	FrameworkVersion string
	Image            string
	Command          []string
	Args             []string
	WorkingDir       string

	// 资源配置
	CPURequest     string
	MemoryRequest  string
	GPUCount       int64
	GPUType        string
	StorageRequest string

	// 分布式训练配置
	MasterReplicas int32
	WorkerReplicas int32
	PSReplicas     int32
	MinAvailable   int32

	// 调度配置
	QueueName         string
	Priority          int64
	PriorityClassName string
	NodeSelector      map[string]string
	Tolerations       []corev1.Toleration
	Affinity          *corev1.Affinity

	// 生命周期配置
	MaxRetry         int32
	TTLAfterFinished *int32
	ActiveDeadline   *int64
	BackoffLimit     *int32

	// 环境配置
	EnvVars      map[string]string
	ConfigMaps   map[string]string
	Secrets      map[string]string
	VolumeMounts []VolumeMountSpec

	// 监控配置
	EnableTensorboard bool
	EnableProfiling   bool
	MetricsPort       int32

	// 元数据
	Labels      map[string]string
	Annotations map[string]string
}

// VolumeMountSpec 存储卷挂载规格
type VolumeMountSpec struct {
	Name      string
	MountPath string
	ReadOnly  bool
	VolumeSource
}

// VolumeSource 存储卷源
type VolumeSource struct {
	Type      string // pvc, hostPath, emptyDir, configMap, secret, nfs
	PVCName   string
	HostPath  string
	ConfigMap string
	Secret    string
	NFSServer string
	NFSPath   string
}

// validateJobSpec 验证作业规格
func (jm *JobManager) validateJobSpec(spec *TrainingJobSpec) error {
	if spec.Name == "" {
		return fmt.Errorf("作业名称不能为空")
	}
	if spec.Image == "" {
		return fmt.Errorf("镜像不能为空")
	}
	if spec.QueueName == "" {
		return fmt.Errorf("队列名称不能为空")
	}
	if spec.MinAvailable <= 0 {
		return fmt.Errorf("最小可用实例数必须大于0")
	}
	return nil
}

// buildVolcanoJobSpec 构建Volcano作业规格
func (jm *JobManager) buildVolcanoJobSpec(spec *TrainingJobSpec) *JobSpec {
	volcanoSpec := &JobSpec{
		Name:                    spec.Name,
		Namespace:               spec.Namespace,
		MinAvailable:            spec.MinAvailable,
		Queue:                   spec.QueueName,
		PriorityClassName:       spec.PriorityClassName,
		MaxRetry:                spec.MaxRetry,
		TTLSecondsAfterFinished: spec.TTLAfterFinished,
		SchedulerName:           "volcano",
		JobType:                 spec.JobType,
		Labels:                  spec.Labels,
		Annotations:             spec.Annotations,
	}

	// 构建运行策略
	volcanoSpec.RunPolicy = &RunPolicy{
		ActiveDeadlineSeconds:   spec.ActiveDeadline,
		BackoffLimit:            spec.BackoffLimit,
		TTLSecondsAfterFinished: spec.TTLAfterFinished,
	}

	// 构建任务列表
	volcanoSpec.Tasks = jm.buildTasks(spec)

	// 构建插件配置
	volcanoSpec.Plugins = jm.buildPluginsForJobType(spec.JobType)

	return volcanoSpec
}

// buildTasks 构建任务列表
func (jm *JobManager) buildTasks(spec *TrainingJobSpec) []TaskSpec {
	var tasks []TaskSpec

	// 构建Master任务（如果需要）
	if spec.MasterReplicas > 0 {
		masterTask := jm.buildMasterTask(spec)
		tasks = append(tasks, masterTask)
	}

	// 构建Worker任务
	if spec.WorkerReplicas > 0 {
		workerTask := jm.buildWorkerTask(spec)
		tasks = append(tasks, workerTask)
	}

	// 构建PS任务（如果是TensorFlow分布式训练）
	if spec.PSReplicas > 0 && strings.ToLower(spec.Framework) == "tensorflow" {
		psTask := jm.buildPSTask(spec)
		tasks = append(tasks, psTask)
	}

	return tasks
}

// buildMasterTask 构建Master任务
func (jm *JobManager) buildMasterTask(spec *TrainingJobSpec) TaskSpec {
	return TaskSpec{
		Name:     "master",
		Replicas: spec.MasterReplicas,
		Template: PodTemplateSpec{
			Metadata: ObjectMeta{
				Labels:      jm.buildTaskLabels(spec, "master"),
				Annotations: spec.Annotations,
			},
			Spec:         jm.buildPodSpec(spec, "master"),
			Affinity:     spec.Affinity,
			NodeSelector: spec.NodeSelector,
			Tolerations:  spec.Tolerations,
		},
		Policies:     jm.buildTaskPolicies("master", spec.JobType),
		MinAvailable: &spec.MinAvailable,
		MaxRetry:     &spec.MaxRetry,
	}
}

// buildWorkerTask 构建Worker任务
func (jm *JobManager) buildWorkerTask(spec *TrainingJobSpec) TaskSpec {
	return TaskSpec{
		Name:     "worker",
		Replicas: spec.WorkerReplicas,
		Template: PodTemplateSpec{
			Metadata: ObjectMeta{
				Labels:      jm.buildTaskLabels(spec, "worker"),
				Annotations: spec.Annotations,
			},
			Spec:         jm.buildPodSpec(spec, "worker"),
			Affinity:     spec.Affinity,
			NodeSelector: spec.NodeSelector,
			Tolerations:  spec.Tolerations,
		},
		Policies: jm.buildTaskPolicies("worker", spec.JobType),
	}
}

// buildPSTask 构建PS任务
func (jm *JobManager) buildPSTask(spec *TrainingJobSpec) TaskSpec {
	return TaskSpec{
		Name:     "ps",
		Replicas: spec.PSReplicas,
		Template: PodTemplateSpec{
			Metadata: ObjectMeta{
				Labels:      jm.buildTaskLabels(spec, "ps"),
				Annotations: spec.Annotations,
			},
			Spec:         jm.buildPodSpec(spec, "ps"),
			Affinity:     spec.Affinity,
			NodeSelector: spec.NodeSelector,
			Tolerations:  spec.Tolerations,
		},
		Policies: jm.buildTaskPolicies("ps", spec.JobType),
	}
}

// buildTaskLabels 构建任务标签
func (jm *JobManager) buildTaskLabels(spec *TrainingJobSpec, taskType string) map[string]string {
	labels := make(map[string]string)

	// 复制用户标签
	for k, v := range spec.Labels {
		labels[k] = v
	}

	// 添加系统标签
	labels["app"] = "volcano-training-job"
	labels["job-name"] = spec.Name
	labels["task-type"] = taskType
	labels["framework"] = spec.Framework
	labels["job-type"] = spec.JobType
	labels["volcano.sh/queue-name"] = spec.QueueName

	return labels
}

// buildPodSpec 构建Pod规格
func (jm *JobManager) buildPodSpec(spec *TrainingJobSpec, taskType string) PodSpec {
	podSpec := PodSpec{
		RestartPolicy: corev1.RestartPolicyNever,
		Containers: []Container{
			jm.buildMainContainer(spec, taskType),
		},
		ServiceAccount: fmt.Sprintf("%s-service-account", spec.Name),
	}

	// 添加初始化容器（如果需要）
	if jm.needInitContainer(spec) {
		podSpec.InitContainers = []Container{
			jm.buildInitContainer(spec),
		}
	}

	// 构建存储卷
	podSpec.Volumes = jm.buildVolumes(spec)

	// 添加监控容器（如果启用）
	if spec.EnableProfiling {
		podSpec.Containers = append(podSpec.Containers, jm.buildProfilingContainer(spec))
	}

	return podSpec
}

// buildMainContainer 构建主要容器
func (jm *JobManager) buildMainContainer(spec *TrainingJobSpec, taskType string) Container {
	container := Container{
		Name:            "training-container",
		Image:           spec.Image,
		Command:         spec.Command,
		Args:            jm.buildContainerArgs(spec, taskType),
		WorkingDir:      spec.WorkingDir,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Resources:       jm.buildResourceRequirements(spec),
		Env:             jm.buildEnvironmentVariables(spec, taskType),
		VolumeMounts:    jm.buildVolumeMounts(spec),
	}

	// 添加健康检查（如果需要）
	if jm.needHealthCheck(spec, taskType) {
		container.ReadinessProbe = jm.buildReadinessProbe(spec)
		container.LivenessProbe = jm.buildLivenessProbe(spec)
	}

	// 添加端口配置
	container.Ports = jm.buildContainerPorts(spec, taskType)

	return container
}

// buildResourceRequirements 构建资源需求
func (jm *JobManager) buildResourceRequirements(spec *TrainingJobSpec) corev1.ResourceRequirements {
	resources := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{},
		Limits:   corev1.ResourceList{},
	}

	// CPU资源
	if spec.CPURequest != "" {
		if cpu, err := resource.ParseQuantity(spec.CPURequest); err == nil {
			resources.Requests[corev1.ResourceCPU] = cpu
			resources.Limits[corev1.ResourceCPU] = cpu
		}
	}

	// 内存资源
	if spec.MemoryRequest != "" {
		if memory, err := resource.ParseQuantity(spec.MemoryRequest); err == nil {
			resources.Requests[corev1.ResourceMemory] = memory
			resources.Limits[corev1.ResourceMemory] = memory
		}
	}

	// GPU资源
	if spec.GPUCount > 0 {
		gpuQuantity := resource.MustParse(fmt.Sprintf("%d", spec.GPUCount))
		resources.Requests["nvidia.com/gpu"] = gpuQuantity
		resources.Limits["nvidia.com/gpu"] = gpuQuantity
	}

	return resources
}

// buildEnvironmentVariables 构建环境变量
func (jm *JobManager) buildEnvironmentVariables(spec *TrainingJobSpec, taskType string) []corev1.EnvVar {
	envVars := []corev1.EnvVar{
		// 基础环境变量
		{Name: "TASK_TYPE", Value: taskType},
		{Name: "JOB_NAME", Value: spec.Name},
		{Name: "FRAMEWORK", Value: spec.Framework},
		{Name: "FRAMEWORK_VERSION", Value: spec.FrameworkVersion},
		{Name: "QUEUE_NAME", Value: spec.QueueName},
		// Pod相关环境变量
		{
			Name: "POD_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
		{
			Name: "POD_NAMESPACE",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.namespace",
				},
			},
		},
		{
			Name: "POD_IP",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		},
	}

	// 添加分布式训练相关环境变量
	envVars = append(envVars, jm.buildDistributedTrainingEnvVars(spec, taskType)...)

	// 添加用户自定义环境变量
	for key, value := range spec.EnvVars {
		envVars = append(envVars, corev1.EnvVar{
			Name:  key,
			Value: value,
		})
	}

	return envVars
}

// buildDistributedTrainingEnvVars 构建分布式训练环境变量
func (jm *JobManager) buildDistributedTrainingEnvVars(spec *TrainingJobSpec, taskType string) []corev1.EnvVar {
	var envVars []corev1.EnvVar

	switch strings.ToLower(spec.Framework) {
	case "pytorch":
		envVars = append(envVars,
			corev1.EnvVar{Name: "MASTER_ADDR", Value: fmt.Sprintf("%s-master-0.%s-master", spec.Name, spec.Name)},
			corev1.EnvVar{Name: "MASTER_PORT", Value: "23456"},
			corev1.EnvVar{Name: "WORLD_SIZE", Value: fmt.Sprintf("%d", spec.WorkerReplicas)},
			corev1.EnvVar{Name: "NCCL_DEBUG", Value: "INFO"},
		)
	case "tensorflow":
		envVars = append(envVars,
			corev1.EnvVar{Name: "TF_CONFIG", Value: jm.buildTFConfig(spec, taskType)},
		)
	case "mpi":
		envVars = append(envVars,
			corev1.EnvVar{Name: "OMPI_ALLOW_RUN_AS_ROOT", Value: "1"},
			corev1.EnvVar{Name: "OMPI_ALLOW_RUN_AS_ROOT_CONFIRM", Value: "1"},
		)
	}

	return envVars
}

// buildTFConfig 构建TensorFlow配置
func (jm *JobManager) buildTFConfig(spec *TrainingJobSpec, taskType string) string {
	// 简化的TF_CONFIG构建逻辑
	// 实际实现中需要根据具体需求构建完整的TF_CONFIG JSON
	return fmt.Sprintf(`{"cluster": {"worker": ["%s-worker:2222"], "ps": ["%s-ps:2222"]}, "task": {"type": "%s", "index": 0}}`,
		spec.Name, spec.Name, taskType)
}

// buildContainerArgs 构建容器参数
func (jm *JobManager) buildContainerArgs(spec *TrainingJobSpec, taskType string) []string {
	args := make([]string, len(spec.Args))
	copy(args, spec.Args)

	// 根据框架类型添加特定参数
	switch strings.ToLower(spec.Framework) {
	case "pytorch":
		args = jm.addPyTorchArgs(args, spec, taskType)
	case "tensorflow":
		args = jm.addTensorFlowArgs(args, spec, taskType)
	}

	return args
}

// addPyTorchArgs 添加PyTorch特定参数
func (jm *JobManager) addPyTorchArgs(args []string, spec *TrainingJobSpec, taskType string) []string {
	if taskType == "master" {
		args = append(args, "--rank=0")
		args = append(args, fmt.Sprintf("--world-size=%d", spec.WorkerReplicas))
	}
	return args
}

// addTensorFlowArgs 添加TensorFlow特定参数
func (jm *JobManager) addTensorFlowArgs(args []string, spec *TrainingJobSpec, taskType string) []string {
	args = append(args, fmt.Sprintf("--task-type=%s", taskType))
	return args
}

// buildVolumeMounts 构建存储卷挂载
func (jm *JobManager) buildVolumeMounts(spec *TrainingJobSpec) []corev1.VolumeMount {
	var mounts []corev1.VolumeMount

	for _, vm := range spec.VolumeMounts {
		mounts = append(mounts, corev1.VolumeMount{
			Name:      vm.Name,
			MountPath: vm.MountPath,
			ReadOnly:  vm.ReadOnly,
		})
	}

	return mounts
}

// buildVolumes 构建存储卷
func (jm *JobManager) buildVolumes(spec *TrainingJobSpec) []corev1.Volume {
	var volumes []corev1.Volume

	for _, vm := range spec.VolumeMounts {
		volume := corev1.Volume{
			Name: vm.Name,
		}

		switch vm.Type {
		case "pvc":
			volume.VolumeSource = corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: vm.PVCName,
				},
			}
		case "hostPath":
			volume.VolumeSource = corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: vm.HostPath,
				},
			}
		case "configMap":
			volume.VolumeSource = corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: vm.ConfigMap,
					},
				},
			}
		case "secret":
			volume.VolumeSource = corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: vm.Secret,
				},
			}
		case "emptyDir":
			volume.VolumeSource = corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			}
		case "nfs":
			volume.VolumeSource = corev1.VolumeSource{
				NFS: &corev1.NFSVolumeSource{
					Server: vm.NFSServer,
					Path:   vm.NFSPath,
				},
			}
		}

		volumes = append(volumes, volume)
	}

	return volumes
}

// buildContainerPorts 构建容器端口
func (jm *JobManager) buildContainerPorts(spec *TrainingJobSpec, taskType string) []corev1.ContainerPort {
	var ports []corev1.ContainerPort

	// 添加框架特定端口
	switch strings.ToLower(spec.Framework) {
	case "pytorch":
		if taskType == "master" || taskType == "worker" {
			ports = append(ports, corev1.ContainerPort{
				Name:          "pytorch-port",
				ContainerPort: 23456,
				Protocol:      corev1.ProtocolTCP,
			})
		}
	case "tensorflow":
		ports = append(ports, corev1.ContainerPort{
			Name:          "tf-port",
			ContainerPort: 2222,
			Protocol:      corev1.ProtocolTCP,
		})
	}

	// 添加监控端口
	if spec.MetricsPort > 0 {
		ports = append(ports, corev1.ContainerPort{
			Name:          "metrics",
			ContainerPort: spec.MetricsPort,
			Protocol:      corev1.ProtocolTCP,
		})
	}

	return ports
}

// buildTaskPolicies 构建任务策略
func (jm *JobManager) buildTaskPolicies(taskType, jobType string) []LifecyclePolicy {
	policies := []LifecyclePolicy{
		{
			Event:  "PodFailed",
			Action: "RestartJob",
		},
		{
			Event:  "PodEvicted",
			Action: "RestartJob",
		},
	}

	// 根据任务类型添加特定策略
	if taskType == "master" {
		policies = append(policies, LifecyclePolicy{
			Event:  "TaskCompleted",
			Action: "CompleteJob",
		})
	}

	return policies
}

// buildPluginsForJobType 为不同作业类型构建插件配置
func (jm *JobManager) buildPluginsForJobType(jobType string) map[string][]string {
	basePlugins := map[string][]string{
		"env":         {},
		"svc":         {},
		"ssh":         {},
		"gang":        {},
		"proportion":  {},
		"drf":         {},
		"predicates":  {},
		"nodeorder":   {},
		"binpack":     {},
		"conformance": {},
	}

	switch strings.ToLower(jobType) {
	case "mpi":
		basePlugins["ssh"] = []string{"--enable-init-container=true"}
	case "tensorflow":
		basePlugins["svc"] = []string{"--enable-headless-service=true"}
	}

	return basePlugins
}

// 健康检查相关方法
func (jm *JobManager) needHealthCheck(spec *TrainingJobSpec, taskType string) bool {
	// 根据框架和任务类型决定是否需要健康检查
	return taskType == "master" || taskType == "worker"
}

func (jm *JobManager) buildReadinessProbe(spec *TrainingJobSpec) *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt(int(spec.MetricsPort)),
			},
		},
		InitialDelaySeconds: 30,
		PeriodSeconds:       10,
		TimeoutSeconds:      5,
		FailureThreshold:    3,
		SuccessThreshold:    1,
	}
}

func (jm *JobManager) buildLivenessProbe(spec *TrainingJobSpec) *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt(int(spec.MetricsPort)),
			},
		},
		InitialDelaySeconds: 60,
		PeriodSeconds:       30,
		TimeoutSeconds:      10,
		FailureThreshold:    3,
	}
}

// 初始化容器相关方法
func (jm *JobManager) needInitContainer(spec *TrainingJobSpec) bool {
	// 根据作业类型决定是否需要初始化容器
	return strings.ToLower(spec.JobType) == "mpi"
}

func (jm *JobManager) buildInitContainer(spec *TrainingJobSpec) Container {
	return Container{
		Name:    "init-container",
		Image:   fmt.Sprintf("%s-init:latest", spec.Framework),
		Command: []string{"sh", "-c", "echo 'Initializing...' && sleep 10"},
	}
}

// 监控容器相关方法
func (jm *JobManager) buildProfilingContainer(spec *TrainingJobSpec) Container {
	return Container{
		Name:  "profiling",
		Image: "prom/node-exporter:latest",
		Args:  []string{"--path.rootfs=/host"},
		Ports: []corev1.ContainerPort{
			{
				ContainerPort: 9100,
				Name:          "profiling",
			},
		},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("100m"),
				corev1.ResourceMemory: resource.MustParse("128Mi"),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("200m"),
				corev1.ResourceMemory: resource.MustParse("256Mi"),
			},
		},
	}
}

// GetJobStatusDetail 获取作业详细状态
func (jm *JobManager) GetJobStatusDetail(namespace, jobName string) (*JobStatusDetail, error) {
	// 获取基础状态
	status, err := jm.client.GetJobStatus(namespace, jobName)
	if err != nil {
		return nil, err
	}

	// 获取Pod详细信息
	pods, err := jm.client.kubeClient.CoreV1().Pods(namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: fmt.Sprintf("job-name=%s", jobName),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("获取Pod列表失败: %v", err)
	}

	// 构建详细状态
	detail := &JobStatusDetail{
		JobStatus:     *status,
		PodStatuses:   make([]PodStatusDetail, 0),
		ResourceUsage: jm.calculateJobResourceUsage(pods.Items),
		Events:        []string{}, // TODO: 获取事件信息
	}

	// 处理Pod状态
	for _, pod := range pods.Items {
		podDetail := jm.buildPodStatusDetail(&pod)
		detail.PodStatuses = append(detail.PodStatuses, podDetail)
	}

	return detail, nil
}

// JobStatusDetail 作业详细状态
type JobStatusDetail struct {
	JobStatus     JobStatus         `json:"jobStatus"`
	PodStatuses   []PodStatusDetail `json:"podStatuses"`
	ResourceUsage ResourceUsage     `json:"resourceUsage"`
	Events        []string          `json:"events"`
}

// PodStatusDetail Pod详细状态
type PodStatusDetail struct {
	Name            string            `json:"name"`
	Phase           string            `json:"phase"`
	NodeName        string            `json:"nodeName"`
	PodIP           string            `json:"podIP"`
	StartTime       *time.Time        `json:"startTime"`
	ContainerStates []ContainerState  `json:"containerStates"`
	Conditions      []string          `json:"conditions"`
	ResourceUsage   map[string]string `json:"resourceUsage"`
}

// ContainerState 容器状态
type ContainerState struct {
	Name         string `json:"name"`
	Ready        bool   `json:"ready"`
	RestartCount int32  `json:"restartCount"`
	State        string `json:"state"`
	Reason       string `json:"reason"`
	Message      string `json:"message"`
}

// ResourceUsage 资源使用情况
type ResourceUsage struct {
	CPUUsage    string `json:"cpuUsage"`
	MemoryUsage string `json:"memoryUsage"`
	GPUUsage    string `json:"gpuUsage"`
}

// buildPodStatusDetail 构建Pod详细状态
func (jm *JobManager) buildPodStatusDetail(pod *corev1.Pod) PodStatusDetail {
	detail := PodStatusDetail{
		Name:     pod.Name,
		Phase:    string(pod.Status.Phase),
		NodeName: pod.Spec.NodeName,
		PodIP:    pod.Status.PodIP,
	}

	if pod.Status.StartTime != nil {
		detail.StartTime = &pod.Status.StartTime.Time
	}

	// 处理容器状态
	for _, containerStatus := range pod.Status.ContainerStatuses {
		state := ContainerState{
			Name:         containerStatus.Name,
			Ready:        containerStatus.Ready,
			RestartCount: containerStatus.RestartCount,
		}

		// 解析容器状态
		if containerStatus.State.Running != nil {
			state.State = "Running"
		} else if containerStatus.State.Waiting != nil {
			state.State = "Waiting"
			state.Reason = containerStatus.State.Waiting.Reason
			state.Message = containerStatus.State.Waiting.Message
		} else if containerStatus.State.Terminated != nil {
			state.State = "Terminated"
			state.Reason = containerStatus.State.Terminated.Reason
			state.Message = containerStatus.State.Terminated.Message
		}

		detail.ContainerStates = append(detail.ContainerStates, state)
	}

	// 处理Pod条件
	for _, condition := range pod.Status.Conditions {
		if condition.Status == corev1.ConditionTrue {
			detail.Conditions = append(detail.Conditions, string(condition.Type))
		}
	}

	return detail
}

// calculateJobResourceUsage 计算作业资源使用
func (jm *JobManager) calculateJobResourceUsage(pods []corev1.Pod) ResourceUsage {
	// 简化的资源使用计算
	// 实际实现中需要调用Metrics API获取实时资源使用情况
	return ResourceUsage{
		CPUUsage:    "计算中...",
		MemoryUsage: "计算中...",
		GPUUsage:    "计算中...",
	}
}
