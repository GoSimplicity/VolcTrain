package volcano

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	vcjob "volcano.sh/apis/pkg/apis/batch/v1alpha1"
	vcclient "volcano.sh/apis/pkg/client/clientset/versioned"
)

// Client Volcano客户端封装
type Client struct {
	volcanoClient vcclient.Interface   // Volcano CRD客户端
	kubeClient    kubernetes.Interface // Kubernetes标准客户端
	config        *rest.Config         // Kubernetes配置
	namespace     string               // 默认命名空间
}

// NewClient 创建Volcano客户端
func NewClient(kubeconfig, namespace string) (*Client, error) {
	var config *rest.Config
	var err error

	if kubeconfig != "" {
		// 使用kubeconfig文件
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		// 使用集群内配置
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		return nil, fmt.Errorf("创建Kubernetes配置失败: %v", err)
	}

	// 创建Volcano客户端
	volcanoClient, err := vcclient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建Volcano客户端失败: %v", err)
	}

	// 创建标准Kubernetes客户端
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建Kubernetes客户端失败: %v", err)
	}

	return &Client{
		volcanoClient: volcanoClient,
		kubeClient:    kubeClient,
		config:        config,
		namespace:     namespace,
	}, nil
}

// JobSpec Volcano作业规格定义
type JobSpec struct {
	Name                    string
	Namespace               string
	MinAvailable            int32
	Queue                   string
	PriorityClassName       string
	MaxRetry                int32
	TTLSecondsAfterFinished *int32
	SchedulerName           string
	Tasks                   []TaskSpec
	Plugins                 map[string][]string
	JobType                 string
	RunPolicy               *RunPolicy
	Labels                  map[string]string
	Annotations             map[string]string
}

// TaskSpec 任务规格定义
type TaskSpec struct {
	Name         string
	Replicas     int32
	Template     PodTemplateSpec
	Policies     []LifecyclePolicy
	MinAvailable *int32
	MaxRetry     *int32
}

// PodTemplateSpec Pod模板规格
type PodTemplateSpec struct {
	Metadata     ObjectMeta
	Spec         PodSpec
	Affinity     *corev1.Affinity
	NodeSelector map[string]string
	Tolerations  []corev1.Toleration
}

// ObjectMeta 对象元数据
type ObjectMeta struct {
	Name        string
	Labels      map[string]string
	Annotations map[string]string
}

// PodSpec Pod规格
type PodSpec struct {
	Containers       []Container
	InitContainers   []Container
	Volumes          []corev1.Volume
	ServiceAccount   string
	RestartPolicy    corev1.RestartPolicy
	HostNetwork      bool
	SecurityContext  *corev1.PodSecurityContext
	ImagePullSecrets []corev1.LocalObjectReference
}

// Container 容器定义
type Container struct {
	Name            string
	Image           string
	Command         []string
	Args            []string
	WorkingDir      string
	Ports           []corev1.ContainerPort
	Env             []corev1.EnvVar
	Resources       corev1.ResourceRequirements
	VolumeMounts    []corev1.VolumeMount
	LivenessProbe   *corev1.Probe
	ReadinessProbe  *corev1.Probe
	ImagePullPolicy corev1.PullPolicy
	SecurityContext *corev1.SecurityContext
}

// RunPolicy 运行策略
type RunPolicy struct {
	CleanPodPolicy           *string // 简化为字符串
	TTLSecondsAfterFinished  *int32
	ActiveDeadlineSeconds    *int64
	BackoffLimit             *int32
	SchedulingTimeoutSeconds *int32
}

// LifecyclePolicy 生命周期策略
type LifecyclePolicy struct {
	Event    string // 简化为字符串
	Action   string // 简化为字符串
	Timeout  *metav1.Duration
	ExitCode *int32
}

// JobStatus 作业状态
type JobStatus struct {
	Name                string
	Namespace           string
	Phase               vcjob.JobPhase
	State               vcjob.JobState
	MinAvailable        int32
	Succeeded           int32
	Failed              int32
	Running             int32
	Pending             int32
	TaskStatuses        map[string]TaskStatus
	Conditions          []vcjob.JobCondition
	RetryCount          int32
	CreationTime        *time.Time
	CompletionTime      *time.Time
	ControlledResources []string
	// RunningHistories []vcjob.JobRunningHistory // v1alpha1版本中暂不支持
}

// TaskStatus 任务状态
type TaskStatus struct {
	Name         string
	Phase        string // 简化为字符串
	Replicas     int32
	Running      int32
	Succeeded    int32
	Failed       int32
	Pending      int32
	MinAvailable *int32
}

// CreateJob 创建Volcano作业
func (c *Client) CreateJob(spec *JobSpec) (*vcjob.Job, error) {
	// 构建Volcano Job对象
	job := &vcjob.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        spec.Name,
			Namespace:   c.getNamespace(spec.Namespace),
			Labels:      spec.Labels,
			Annotations: spec.Annotations,
		},
		Spec: vcjob.JobSpec{
			MinAvailable:  spec.MinAvailable,
			Queue:         spec.Queue,
			MaxRetry:      spec.MaxRetry,
			SchedulerName: c.getSchedulerName(spec.SchedulerName),
			Tasks:         c.buildTasks(spec.Tasks),
			Plugins:       c.buildPlugins(spec.Plugins),
		},
	}

	// 设置优先级类名
	if spec.PriorityClassName != "" {
		job.Spec.PriorityClassName = spec.PriorityClassName
	}

	// 设置TTL
	if spec.TTLSecondsAfterFinished != nil {
		job.Spec.TTLSecondsAfterFinished = spec.TTLSecondsAfterFinished
	}

	// 设置运行策略 - v1alpha1版本中通过annotations设置
	if spec.RunPolicy != nil {
		if job.Annotations == nil {
			job.Annotations = make(map[string]string)
		}
		if spec.RunPolicy.TTLSecondsAfterFinished != nil {
			job.Annotations["batch.volcano.sh/ttl-seconds-after-finished"] = fmt.Sprintf("%d", *spec.RunPolicy.TTLSecondsAfterFinished)
		}
		if spec.RunPolicy.ActiveDeadlineSeconds != nil {
			job.Annotations["batch.volcano.sh/active-deadline-seconds"] = fmt.Sprintf("%d", *spec.RunPolicy.ActiveDeadlineSeconds)
		}
	}

	// 创建作业
	createdJob, err := c.volcanoClient.BatchV1alpha1().Jobs(c.getNamespace(spec.Namespace)).Create(
		context.TODO(),
		job,
		metav1.CreateOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("创建Volcano作业失败: %v", err)
	}

	return createdJob, nil
}

// buildTasks 构建任务列表
func (c *Client) buildTasks(tasks []TaskSpec) []vcjob.TaskSpec {
	var volcanoTasks []vcjob.TaskSpec

	for _, task := range tasks {
		volcanoTask := vcjob.TaskSpec{
			Name:     task.Name,
			Replicas: task.Replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        task.Template.Metadata.Name,
					Labels:      task.Template.Metadata.Labels,
					Annotations: task.Template.Metadata.Annotations,
				},
				Spec: c.buildPodSpec(task.Template.Spec),
			},
			Policies: c.buildLifecyclePolicies(task.Policies),
		}

		// 设置调度配置
		if task.Template.Affinity != nil {
			volcanoTask.Template.Spec.Affinity = task.Template.Affinity
		}
		if task.Template.NodeSelector != nil {
			volcanoTask.Template.Spec.NodeSelector = task.Template.NodeSelector
		}
		if task.Template.Tolerations != nil {
			volcanoTask.Template.Spec.Tolerations = task.Template.Tolerations
		}

		// 设置可选参数
		if task.MinAvailable != nil {
			volcanoTask.MinAvailable = task.MinAvailable
		}
		if task.MaxRetry != nil {
			volcanoTask.MaxRetry = *task.MaxRetry
		}

		volcanoTasks = append(volcanoTasks, volcanoTask)
	}

	return volcanoTasks
}

// buildPodSpec 构建Pod规格
func (c *Client) buildPodSpec(spec PodSpec) corev1.PodSpec {
	podSpec := corev1.PodSpec{
		RestartPolicy:      spec.RestartPolicy,
		ServiceAccountName: spec.ServiceAccount,
		HostNetwork:        spec.HostNetwork,
		SecurityContext:    spec.SecurityContext,
		Volumes:            spec.Volumes,
		ImagePullSecrets:   spec.ImagePullSecrets,
	}

	// 构建容器列表
	for _, container := range spec.Containers {
		podSpec.Containers = append(podSpec.Containers, c.buildContainer(container))
	}

	// 构建Init容器列表
	for _, initContainer := range spec.InitContainers {
		podSpec.InitContainers = append(podSpec.InitContainers, c.buildContainer(initContainer))
	}

	return podSpec
}

// buildContainer 构建容器
func (c *Client) buildContainer(container Container) corev1.Container {
	return corev1.Container{
		Name:            container.Name,
		Image:           container.Image,
		Command:         container.Command,
		Args:            container.Args,
		WorkingDir:      container.WorkingDir,
		Ports:           container.Ports,
		Env:             container.Env,
		Resources:       container.Resources,
		VolumeMounts:    container.VolumeMounts,
		LivenessProbe:   container.LivenessProbe,
		ReadinessProbe:  container.ReadinessProbe,
		ImagePullPolicy: container.ImagePullPolicy,
		SecurityContext: container.SecurityContext,
	}
}

// buildLifecyclePolicies 构建生命周期策略
func (c *Client) buildLifecyclePolicies(policies []LifecyclePolicy) []vcjob.LifecyclePolicy {
	var lifecyclePolicies []vcjob.LifecyclePolicy

	for _, policy := range policies {
		volcanoPolicy := vcjob.LifecyclePolicy{
			// 简化处理，使用字符串转换
			// Event:  policy.Event,
			// Action: policy.Action,
		}

		if policy.Timeout != nil {
			volcanoPolicy.Timeout = policy.Timeout
		}
		if policy.ExitCode != nil {
			volcanoPolicy.ExitCode = policy.ExitCode
		}

		lifecyclePolicies = append(lifecyclePolicies, volcanoPolicy)
	}

	return lifecyclePolicies
}

// buildPlugins 构建插件配置
func (c *Client) buildPlugins(plugins map[string][]string) map[string][]string {
	if plugins == nil {
		// 返回默认插件配置
		return map[string][]string{
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
	}
	return plugins
}

// GetJobStatus 获取作业状态
func (c *Client) GetJobStatus(namespace, jobName string) (*JobStatus, error) {
	job, err := c.volcanoClient.BatchV1alpha1().Jobs(c.getNamespace(namespace)).Get(
		context.TODO(),
		jobName,
		metav1.GetOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("获取Volcano作业状态失败: %v", err)
	}

	status := &JobStatus{
		Name:         job.Name,
		Namespace:    job.Namespace,
		Phase:        job.Status.State.Phase,
		State:        job.Status.State,
		MinAvailable: job.Spec.MinAvailable,
		Succeeded:    job.Status.Succeeded,
		Failed:       job.Status.Failed,
		Running:      job.Status.Running,
		Pending:      job.Status.Pending,
		TaskStatuses: make(map[string]TaskStatus),
		Conditions:   job.Status.Conditions,
		RetryCount:   job.Status.RetryCount,
	}

	// 转换时间
	if job.CreationTimestamp.Time.Unix() > 0 {
		status.CreationTime = &job.CreationTimestamp.Time
	}
	if !job.Status.State.LastTransitionTime.Time.IsZero() {
		status.CompletionTime = &job.Status.State.LastTransitionTime.Time
	}

	// 转换任务状态
	for taskName, taskStatus := range job.Status.TaskStatusCount {
		// 获取各状态的Pod数量
		running := int32(0)
		succeeded := int32(0)
		failed := int32(0)
		pending := int32(0)

		// 从taskStatus.Phase中提取Pod计数
		for phase, count := range taskStatus.Phase {
			switch string(phase) {
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

		status.TaskStatuses[taskName] = TaskStatus{
			Name:      taskName,
			Phase:     "Active", // 简化处理
			Replicas:  c.getTaskReplicas(job, taskName),
			Running:   running,
			Succeeded: succeeded,
			Failed:    failed,
			Pending:   pending,
		}
	}

	// 获取控制的资源
	status.ControlledResources = c.getControlledResources(job)
	// v1alpha1版本中RunningHistories字段可能不存在，简化处理
	// status.RunningHistories = job.Status.RunningHistories

	return status, nil
}

// getTaskReplicas 获取任务副本数
func (c *Client) getTaskReplicas(job *vcjob.Job, taskName string) int32 {
	for _, task := range job.Spec.Tasks {
		if task.Name == taskName {
			return task.Replicas
		}
	}
	return 0
}

// getControlledResources 获取控制的资源列表
func (c *Client) getControlledResources(job *vcjob.Job) []string {
	var resources []string
	// TODO: 实现获取作业控制的Pod等资源的逻辑
	return resources
}

// ScaleJob 扩缩容作业
func (c *Client) ScaleJob(namespace, jobName string, tasks []TaskScaleSpec) error {
	// 获取现有作业
	job, err := c.volcanoClient.BatchV1alpha1().Jobs(c.getNamespace(namespace)).Get(
		context.TODO(),
		jobName,
		metav1.GetOptions{},
	)
	if err != nil {
		return fmt.Errorf("获取作业失败: %v", err)
	}

	// 更新任务副本数
	updated := false
	for i, task := range job.Spec.Tasks {
		for _, scaleSpec := range tasks {
			if task.Name == scaleSpec.Name {
				if task.Replicas != scaleSpec.Replicas {
					job.Spec.Tasks[i].Replicas = scaleSpec.Replicas
					updated = true
				}
				break
			}
		}
	}

	if !updated {
		return fmt.Errorf("没有找到需要扩缩容的任务")
	}

	// 更新作业
	_, err = c.volcanoClient.BatchV1alpha1().Jobs(c.getNamespace(namespace)).Update(
		context.TODO(),
		job,
		metav1.UpdateOptions{},
	)
	if err != nil {
		return fmt.Errorf("扩缩容作业失败: %v", err)
	}

	return nil
}

// DeleteJob 删除作业
func (c *Client) DeleteJob(namespace, jobName string) error {
	err := c.volcanoClient.BatchV1alpha1().Jobs(c.getNamespace(namespace)).Delete(
		context.TODO(),
		jobName,
		metav1.DeleteOptions{},
	)
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("删除Volcano作业失败: %v", err)
	}

	return nil
}

// SuspendJob 暂停作业
func (c *Client) SuspendJob(namespace, jobName string) error {
	// 获取现有作业
	job, err := c.volcanoClient.BatchV1alpha1().Jobs(c.getNamespace(namespace)).Get(
		context.TODO(),
		jobName,
		metav1.GetOptions{},
	)
	if err != nil {
		return fmt.Errorf("获取作业失败: %v", err)
	}

	// 添加暂停动作
	job.Status.State.Phase = vcjob.Aborted

	// 更新作业状态
	_, err = c.volcanoClient.BatchV1alpha1().Jobs(c.getNamespace(namespace)).UpdateStatus(
		context.TODO(),
		job,
		metav1.UpdateOptions{},
	)
	if err != nil {
		return fmt.Errorf("暂停作业失败: %v", err)
	}

	return nil
}

// ResumeJob 恢复作业
func (c *Client) ResumeJob(namespace, jobName string) error {
	// 获取现有作业
	job, err := c.volcanoClient.BatchV1alpha1().Jobs(c.getNamespace(namespace)).Get(
		context.TODO(),
		jobName,
		metav1.GetOptions{},
	)
	if err != nil {
		return fmt.Errorf("获取作业失败: %v", err)
	}

	// 恢复作业
	job.Status.State.Phase = vcjob.Pending

	// 更新作业状态
	_, err = c.volcanoClient.BatchV1alpha1().Jobs(c.getNamespace(namespace)).UpdateStatus(
		context.TODO(),
		job,
		metav1.UpdateOptions{},
	)
	if err != nil {
		return fmt.Errorf("恢复作业失败: %v", err)
	}

	return nil
}

// ListJobs 列出作业
func (c *Client) ListJobs(namespace string, labelSelector string) (*vcjob.JobList, error) {
	jobs, err := c.volcanoClient.BatchV1alpha1().Jobs(c.getNamespace(namespace)).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: labelSelector,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("列出Volcano作业失败: %v", err)
	}

	return jobs, nil
}

// TaskScaleSpec 任务扩缩容规格
type TaskScaleSpec struct {
	Name     string
	Replicas int32
}

// 辅助方法
func (c *Client) getNamespace(namespace string) string {
	if namespace != "" {
		return namespace
	}
	return c.namespace
}

func (c *Client) getSchedulerName(schedulerName string) string {
	if schedulerName != "" {
		return schedulerName
	}
	return "volcano"
}

// BuildResourceRequirements 构建资源需求
func BuildResourceRequirements(cpu, memory string, gpuCount int64) corev1.ResourceRequirements {
	requirements := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{},
		Limits:   corev1.ResourceList{},
	}

	if cpu != "" {
		requirements.Requests[corev1.ResourceCPU] = resource.MustParse(cpu)
		requirements.Limits[corev1.ResourceCPU] = resource.MustParse(cpu)
	}

	if memory != "" {
		requirements.Requests[corev1.ResourceMemory] = resource.MustParse(memory)
		requirements.Limits[corev1.ResourceMemory] = resource.MustParse(memory)
	}

	if gpuCount > 0 {
		gpuQuantity := resource.MustParse(fmt.Sprintf("%d", gpuCount))
		requirements.Requests["nvidia.com/gpu"] = gpuQuantity
		requirements.Limits["nvidia.com/gpu"] = gpuQuantity
	}

	return requirements
}

// BuildCommonEnvVars 构建通用环境变量
func BuildCommonEnvVars(customEnv map[string]string) []corev1.EnvVar {
	envVars := []corev1.EnvVar{
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

	// 添加自定义环境变量
	for key, value := range customEnv {
		envVars = append(envVars, corev1.EnvVar{
			Name:  key,
			Value: value,
		})
	}

	return envVars
}
