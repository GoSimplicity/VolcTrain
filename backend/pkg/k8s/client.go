package k8s

import (
	"context"
	"fmt"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client Kubernetes客户端封装
type Client struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
	namespace string
}

// NewClient 创建Kubernetes客户端
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

	// 创建标准客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建Kubernetes客户端失败: %v", err)
	}

	return &Client{
		clientset: clientset,
		config:    config,
		namespace: namespace,
	}, nil
}

// JobSpec 训练作业规格
type JobSpec struct {
	Name          string
	Image         string
	Command       []string
	Args          []string
	Env           map[string]string
	Resources     ResourceRequirements
	WorkingDir    string
	QueueName     string
	Priority      int32
	Replicas      int32
	NodeSelector  map[string]string
	Tolerations   []corev1.Toleration
	Volumes       []VolumeSpec
	ConfigMaps    []ConfigMapSpec
	Secrets       []SecretSpec
	RestartPolicy string
	MaxRetryCount int32
}

// ResourceRequirements 资源需求
type ResourceRequirements struct {
	CPURequests     string
	CPULimits       string
	MemoryRequests  string
	MemoryLimits    string
	GPURequests     int
	GPULimits       int
	StorageRequests string
}

// VolumeSpec 存储卷规格
type VolumeSpec struct {
	Name      string
	Type      string // pvc, hostPath, emptyDir, configMap, secret
	MountPath string
	ReadOnly  bool
	Source    map[string]string
}

// ConfigMapSpec ConfigMap规格
type ConfigMapSpec struct {
	Name      string
	MountPath string
	Items     map[string]string
}

// SecretSpec Secret规格
type SecretSpec struct {
	Name      string
	MountPath string
	Items     map[string]string
}

// JobStatus 作业状态
type JobStatus struct {
	Phase          string
	State          string
	Reason         string
	Message        string
	StartTime      *time.Time
	CompletionTime *time.Time
	Active         int32
	Succeeded      int32
	Failed         int32
	Conditions     []JobCondition
}

// JobCondition 作业条件
type JobCondition struct {
	Type               string
	Status             string
	LastTransitionTime time.Time
	Reason             string
	Message            string
}

// CreateJob 创建Kubernetes Job
func (c *Client) CreateJob(spec *JobSpec) (*batchv1.Job, error) {
	// 构建Job对象
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.Name,
			Namespace: c.namespace,
			Labels: map[string]string{
				"app":      "training-job",
				"job-name": spec.Name,
				"queue":    spec.QueueName,
			},
		},
		Spec: batchv1.JobSpec{
			Parallelism:             &spec.Replicas,
			Completions:             &spec.Replicas,
			BackoffLimit:            &spec.MaxRetryCount,
			TTLSecondsAfterFinished: int32Ptr(3600), // 1小时后清理
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":      "training-job",
						"job-name": spec.Name,
					},
				},
				Spec: c.buildPodSpec(spec),
			},
		},
	}

	// 创建作业
	createdJob, err := c.clientset.BatchV1().Jobs(c.namespace).Create(
		context.TODO(),
		job,
		metav1.CreateOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("创建Kubernetes作业失败: %v", err)
	}

	return createdJob, nil
}

// buildPodSpec 构建Pod规格
func (c *Client) buildPodSpec(spec *JobSpec) corev1.PodSpec {
	podSpec := corev1.PodSpec{
		RestartPolicy: corev1.RestartPolicy(spec.RestartPolicy),
		Containers: []corev1.Container{
			{
				Name:            "training-container",
				Image:           spec.Image,
				Command:         spec.Command,
				Args:            spec.Args,
				WorkingDir:      spec.WorkingDir,
				ImagePullPolicy: corev1.PullIfNotPresent,
				Resources:       c.buildResourceRequirements(spec.Resources),
				Env:             c.buildEnvVars(spec.Env),
				VolumeMounts:    c.buildVolumeMounts(spec.Volumes, spec.ConfigMaps, spec.Secrets),
			},
		},
		Volumes:      c.buildVolumes(spec.Volumes, spec.ConfigMaps, spec.Secrets),
		NodeSelector: spec.NodeSelector,
		Tolerations:  spec.Tolerations,
	}

	return podSpec
}

// buildResourceRequirements 构建资源需求
func (c *Client) buildResourceRequirements(res ResourceRequirements) corev1.ResourceRequirements {
	requirements := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{},
		Limits:   corev1.ResourceList{},
	}

	if res.CPURequests != "" {
		requirements.Requests[corev1.ResourceCPU] = resource.MustParse(res.CPURequests)
	}
	if res.MemoryRequests != "" {
		requirements.Requests[corev1.ResourceMemory] = resource.MustParse(res.MemoryRequests)
	}
	if res.CPULimits != "" {
		requirements.Limits[corev1.ResourceCPU] = resource.MustParse(res.CPULimits)
	}
	if res.MemoryLimits != "" {
		requirements.Limits[corev1.ResourceMemory] = resource.MustParse(res.MemoryLimits)
	}
	if res.GPURequests > 0 {
		requirements.Requests["nvidia.com/gpu"] = resource.MustParse(fmt.Sprintf("%d", res.GPURequests))
	}
	if res.GPULimits > 0 {
		requirements.Limits["nvidia.com/gpu"] = resource.MustParse(fmt.Sprintf("%d", res.GPULimits))
	}

	return requirements
}

// buildEnvVars 构建环境变量
func (c *Client) buildEnvVars(env map[string]string) []corev1.EnvVar {
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
	}

	// 添加用户自定义环境变量
	for key, value := range env {
		envVars = append(envVars, corev1.EnvVar{
			Name:  key,
			Value: value,
		})
	}

	return envVars
}

// buildVolumeMounts 构建存储卷挂载
func (c *Client) buildVolumeMounts(volumes []VolumeSpec, configMaps []ConfigMapSpec, secrets []SecretSpec) []corev1.VolumeMount {
	var mounts []corev1.VolumeMount

	// 添加存储卷挂载
	for _, vol := range volumes {
		mounts = append(mounts, corev1.VolumeMount{
			Name:      vol.Name,
			MountPath: vol.MountPath,
			ReadOnly:  vol.ReadOnly,
		})
	}

	// 添加ConfigMap挂载
	for _, cm := range configMaps {
		mounts = append(mounts, corev1.VolumeMount{
			Name:      cm.Name,
			MountPath: cm.MountPath,
			ReadOnly:  true,
		})
	}

	// 添加Secret挂载
	for _, secret := range secrets {
		mounts = append(mounts, corev1.VolumeMount{
			Name:      secret.Name,
			MountPath: secret.MountPath,
			ReadOnly:  true,
		})
	}

	return mounts
}

// buildVolumes 构建存储卷
func (c *Client) buildVolumes(volumeSpecs []VolumeSpec, configMaps []ConfigMapSpec, secrets []SecretSpec) []corev1.Volume {
	var volumes []corev1.Volume

	// 添加存储卷
	for _, spec := range volumeSpecs {
		volume := corev1.Volume{
			Name: spec.Name,
		}

		switch spec.Type {
		case "pvc":
			volume.VolumeSource = corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: spec.Source["claimName"],
				},
			}
		case "hostPath":
			volume.VolumeSource = corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: spec.Source["path"],
				},
			}
		case "emptyDir":
			volume.VolumeSource = corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			}
		}

		volumes = append(volumes, volume)
	}

	// 添加ConfigMap卷
	for _, cm := range configMaps {
		volumes = append(volumes, corev1.Volume{
			Name: cm.Name,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: cm.Name,
					},
				},
			},
		})
	}

	// 添加Secret卷
	for _, secret := range secrets {
		volumes = append(volumes, corev1.Volume{
			Name: secret.Name,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: secret.Name,
				},
			},
		})
	}

	return volumes
}

// GetJobStatus 获取作业状态
func (c *Client) GetJobStatus(jobName string) (*JobStatus, error) {
	job, err := c.clientset.BatchV1().Jobs(c.namespace).Get(
		context.TODO(),
		jobName,
		metav1.GetOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("获取作业状态失败: %v", err)
	}

	status := &JobStatus{
		Active:    job.Status.Active,
		Succeeded: job.Status.Succeeded,
		Failed:    job.Status.Failed,
	}

	// 确定作业阶段
	if job.Status.Succeeded > 0 {
		status.Phase = "Succeeded"
		status.State = "Completed"
	} else if job.Status.Failed > 0 {
		status.Phase = "Failed"
		status.State = "Failed"
	} else if job.Status.Active > 0 {
		status.Phase = "Running"
		status.State = "Running"
	} else {
		status.Phase = "Pending"
		status.State = "Pending"
	}

	// 转换时间
	if job.Status.StartTime != nil {
		t := job.Status.StartTime.Time
		status.StartTime = &t
	}
	if job.Status.CompletionTime != nil {
		t := job.Status.CompletionTime.Time
		status.CompletionTime = &t
	}

	// 转换条件
	for _, cond := range job.Status.Conditions {
		status.Conditions = append(status.Conditions, JobCondition{
			Type:               string(cond.Type),
			Status:             string(cond.Status),
			LastTransitionTime: cond.LastTransitionTime.Time,
			Reason:             cond.Reason,
			Message:            cond.Message,
		})
	}

	return status, nil
}

// DeleteJob 删除作业
func (c *Client) DeleteJob(jobName string) error {
	err := c.clientset.BatchV1().Jobs(c.namespace).Delete(
		context.TODO(),
		jobName,
		metav1.DeleteOptions{},
	)
	if err != nil {
		return fmt.Errorf("删除作业失败: %v", err)
	}

	return nil
}

// ListJobs 列出作业
func (c *Client) ListJobs(labelSelector string) (*batchv1.JobList, error) {
	jobs, err := c.clientset.BatchV1().Jobs(c.namespace).List(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: labelSelector,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("列出作业失败: %v", err)
	}

	return jobs, nil
}

// int32Ptr 返回int32指针
func int32Ptr(i int32) *int32 {
	return &i
}
