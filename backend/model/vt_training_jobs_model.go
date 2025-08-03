package model

import (
	"database/sql"
	"time"
)

// VtTrainingJobs 训练作业表模型
type VtTrainingJobs struct {
	Id                        int64      `db:"id" json:"id"`
	Name                      string     `db:"name" json:"name"`
	DisplayName               string     `db:"display_name" json:"displayName"`
	Description               string     `db:"description" json:"description"`
	JobType                   string     `db:"job_type" json:"jobType"`
	Framework                 string     `db:"framework" json:"framework"`
	FrameworkVersion          string     `db:"framework_version" json:"frameworkVersion"`
	PythonVersion             string     `db:"python_version" json:"pythonVersion"`
	CodeSourceType            string     `db:"code_source_type" json:"codeSourceType"`
	CodeSourceConfig          string     `db:"code_source_config" json:"codeSourceConfig"`
	EntryPoint                string     `db:"entry_point" json:"entryPoint"`
	WorkingDir                string     `db:"working_dir" json:"workingDir"`
	Image                     string     `db:"image" json:"image"`
	ImagePullPolicy           string     `db:"image_pull_policy" json:"imagePullPolicy"`
	ImagePullSecrets          string     `db:"image_pull_secrets" json:"imagePullSecrets"`
	DatasetMountConfigs       string     `db:"dataset_mount_configs" json:"datasetMountConfigs"`
	DataSourceConfig          string     `db:"data_source_config" json:"dataSourceConfig"`
	ModelConfig               string     `db:"model_config" json:"modelConfig"`
	OutputModelName           string     `db:"output_model_name" json:"outputModelName"`
	ModelSaveStrategy         string     `db:"model_save_strategy" json:"modelSaveStrategy"`
	CpuCores                  string     `db:"cpu_cores" json:"cpuCores"`
	MemoryGb                  string     `db:"memory_gb" json:"memoryGb"`
	GpuCount                  int        `db:"gpu_count" json:"gpuCount"`
	GpuType                   string     `db:"gpu_type" json:"gpuType"`
	GpuMemoryGb               string     `db:"gpu_memory_gb" json:"gpuMemoryGb"`
	StorageGb                 string     `db:"storage_gb" json:"storageGb"`
	SharedMemoryGb            string     `db:"shared_memory_gb" json:"sharedMemoryGb"`
	WorkerCount               int        `db:"worker_count" json:"workerCount"`
	PsCount                   int        `db:"ps_count" json:"psCount"`
	MasterCount               int        `db:"master_count" json:"masterCount"`
	EnvVars                   string     `db:"env_vars" json:"envVars"`
	CommandArgs               string     `db:"command_args" json:"commandArgs"`
	Secrets                   string     `db:"secrets" json:"secrets"`
	ConfigMaps                string     `db:"config_maps" json:"configMaps"`
	VolumeMounts              string     `db:"volume_mounts" json:"volumeMounts"`
	QueueName                 string     `db:"queue_name" json:"queueName"`
	Priority                  int        `db:"priority" json:"priority"`
	NodeSelector              string     `db:"node_selector" json:"nodeSelector"`
	Tolerations               string     `db:"tolerations" json:"tolerations"`
	Affinity                  string     `db:"affinity" json:"affinity"`
	MaxRuntimeSeconds         int        `db:"max_runtime_seconds" json:"maxRuntimeSeconds"`
	MaxIdleSeconds            int        `db:"max_idle_seconds" json:"maxIdleSeconds"`
	AutoRestart               bool       `db:"auto_restart" json:"autoRestart"`
	MaxRetryCount             int        `db:"max_retry_count" json:"maxRetryCount"`
	VolcanoJobName            string     `db:"volcano_job_name" json:"volcanoJobName"`
	VolcanoQueue              string     `db:"volcano_queue" json:"volcanoQueue"`
	MinAvailable              int        `db:"min_available" json:"minAvailable"`
	Status                    string     `db:"status" json:"status"`
	Phase                     string     `db:"phase" json:"phase"`
	Namespace                 string     `db:"namespace" json:"namespace"`
	ClusterName               string     `db:"cluster_name" json:"clusterName"`
	ErrorMessage              string     `db:"error_message" json:"errorMessage"`
	ErrorCode                 string     `db:"error_code" json:"errorCode"`
	ExitCode                  int        `db:"exit_code" json:"exitCode"`
	FailureReason             string     `db:"failure_reason" json:"failureReason"`
	SubmittedAt               time.Time  `db:"submitted_at" json:"submittedAt"`
	QueuedAt                  *time.Time `db:"queued_at" json:"queuedAt"`
	ScheduledAt               *time.Time `db:"scheduled_at" json:"scheduledAt"`
	StartTime                 *time.Time `db:"start_time" json:"startTime"`
	EndTime                   *time.Time `db:"end_time" json:"endTime"`
	DurationSeconds           int        `db:"duration_seconds" json:"durationSeconds"`
	ActualCpuUsage            string     `db:"actual_cpu_usage" json:"actualCpuUsage"`
	ActualMemoryUsageGb       string     `db:"actual_memory_usage_gb" json:"actualMemoryUsageGb"`
	ActualGpuUsage            string     `db:"actual_gpu_usage" json:"actualGpuUsage"`
	PeakMemoryUsageGb         string     `db:"peak_memory_usage_gb" json:"peakMemoryUsageGb"`
	TotalGpuHours             string     `db:"total_gpu_hours" json:"totalGpuHours"`
	WorkspacePath             string     `db:"workspace_path" json:"workspacePath"`
	LogsPath                  string     `db:"logs_path" json:"logsPath"`
	OutputPath                string     `db:"output_path" json:"outputPath"`
	CheckpointPath            string     `db:"checkpoint_path" json:"checkpointPath"`
	TensorboardPath           string     `db:"tensorboard_path" json:"tensorboardPath"`
	Hyperparameters           string     `db:"hyperparameters" json:"hyperparameters"`
	TrainingConfig            string     `db:"training_config" json:"trainingConfig"`
	OptimizerConfig           string     `db:"optimizer_config" json:"optimizerConfig"`
	SchedulerConfig           string     `db:"scheduler_config" json:"schedulerConfig"`
	EnableTensorboard         bool       `db:"enable_tensorboard" json:"enableTensorboard"`
	EnableProfiling           bool       `db:"enable_profiling" json:"enableProfiling"`
	MetricsCollectionInterval int        `db:"metrics_collection_interval" json:"metricsCollectionInterval"`
	NotificationConfig        string     `db:"notification_config" json:"notificationConfig"`
	Tags                      string     `db:"tags" json:"tags"`
	Annotations               string     `db:"annotations" json:"annotations"`
	Metadata                  string     `db:"metadata" json:"metadata"`
	CreatedAt                 time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt                 time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt                 *time.Time `db:"deleted_at" json:"deletedAt"`
}

// VtTrainingJobsModel 训练作业模型操作接口
type VtTrainingJobsModel interface {
	Insert(data *VtTrainingJobs) (sql.Result, error)
	FindOne(id int64) (*VtTrainingJobs, error)
	FindOneByName(name string) (*VtTrainingJobs, error)
	Update(data *VtTrainingJobs) error
	UpdateStatus(id int64, status, phase string) error
	Delete(id int64) error
	List(page, pageSize int, filters map[string]interface{}) ([]*VtTrainingJobs, int64, error)
	GetByStatus(status string) ([]*VtTrainingJobs, error)
}

type vtTrainingJobsModel struct {
	conn *sql.DB
}

func NewVtTrainingJobsModel(conn *sql.DB) VtTrainingJobsModel {
	return &vtTrainingJobsModel{conn: conn}
}

func (m *vtTrainingJobsModel) Insert(data *VtTrainingJobs) (sql.Result, error) {
	query := `INSERT INTO vt_training_jobs (name, display_name, description, job_type, framework, framework_version, python_version, code_source_type, code_source_config, entry_point, working_dir, image, image_pull_policy, gpu_count, gpu_type, worker_count, master_count, queue_name, priority, max_runtime_seconds, max_idle_seconds, auto_restart, max_retry_count, min_available, status, phase, enable_tensorboard, enable_profiling, metrics_collection_interval) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	return m.conn.Exec(query, data.Name, data.DisplayName, data.Description, data.JobType, data.Framework, data.FrameworkVersion, data.PythonVersion, data.CodeSourceType, data.CodeSourceConfig, data.EntryPoint, data.WorkingDir, data.Image, data.ImagePullPolicy, data.GpuCount, data.GpuType, data.WorkerCount, data.MasterCount, data.QueueName, data.Priority, data.MaxRuntimeSeconds, data.MaxIdleSeconds, data.AutoRestart, data.MaxRetryCount, data.MinAvailable, data.Status, data.Phase, data.EnableTensorboard, data.EnableProfiling, data.MetricsCollectionInterval)
}

func (m *vtTrainingJobsModel) FindOne(id int64) (*VtTrainingJobs, error) {
	var job VtTrainingJobs
	query := `SELECT id, name, display_name, description, job_type, framework, framework_version, python_version, code_source_type, entry_point, working_dir, image, gpu_count, gpu_type, worker_count, master_count, queue_name, priority, status, phase, submitted_at, start_time, end_time, created_at, updated_at FROM vt_training_jobs WHERE id = ? AND deleted_at IS NULL`
	err := m.conn.QueryRow(query, id).Scan(&job.Id, &job.Name, &job.DisplayName, &job.Description, &job.JobType, &job.Framework, &job.FrameworkVersion, &job.PythonVersion, &job.CodeSourceType, &job.EntryPoint, &job.WorkingDir, &job.Image, &job.GpuCount, &job.GpuType, &job.WorkerCount, &job.MasterCount, &job.QueueName, &job.Priority, &job.Status, &job.Phase, &job.SubmittedAt, &job.StartTime, &job.EndTime, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (m *vtTrainingJobsModel) FindOneByName(name string) (*VtTrainingJobs, error) {
	var job VtTrainingJobs
	query := `SELECT id, name, display_name, description, job_type, framework, status, phase, created_at FROM vt_training_jobs WHERE name = ? AND deleted_at IS NULL`
	err := m.conn.QueryRow(query, name).Scan(&job.Id, &job.Name, &job.DisplayName, &job.Description, &job.JobType, &job.Framework, &job.Status, &job.Phase, &job.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (m *vtTrainingJobsModel) Update(data *VtTrainingJobs) error {
	query := `UPDATE vt_training_jobs SET display_name = ?, description = ?, priority = ?, status = ?, phase = ?, error_message = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := m.conn.Exec(query, data.DisplayName, data.Description, data.Priority, data.Status, data.Phase, data.ErrorMessage, data.Id)
	return err
}

func (m *vtTrainingJobsModel) UpdateStatus(id int64, status, phase string) error {
	query := `UPDATE vt_training_jobs SET status = ?, phase = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := m.conn.Exec(query, status, phase, id)
	return err
}

func (m *vtTrainingJobsModel) Delete(id int64) error {
	query := `UPDATE vt_training_jobs SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := m.conn.Exec(query, id)
	return err
}

func (m *vtTrainingJobsModel) List(page, pageSize int, filters map[string]interface{}) ([]*VtTrainingJobs, int64, error) {
	offset := (page - 1) * pageSize

	whereClause := "WHERE deleted_at IS NULL"
	args := []interface{}{}

	if status, ok := filters["status"]; ok {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	if framework, ok := filters["framework"]; ok {
		whereClause += " AND framework = ?"
		args = append(args, framework)
	}

	if queueName, ok := filters["queue_name"]; ok {
		whereClause += " AND queue_name = ?"
		args = append(args, queueName)
	}

	if keyword, ok := filters["keyword"]; ok {
		whereClause += " AND (name LIKE ? OR display_name LIKE ?)"
		searchPattern := "%" + keyword.(string) + "%"
		args = append(args, searchPattern, searchPattern)
	}

	// 查询总数
	countQuery := "SELECT COUNT(*) FROM vt_training_jobs " + whereClause
	var total int64
	err := m.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据列表
	listQuery := `SELECT id, name, display_name, description, job_type, framework, status, phase, priority, gpu_count, submitted_at, start_time, end_time, created_at FROM vt_training_jobs ` + whereClause + ` ORDER BY submitted_at DESC LIMIT ? OFFSET ?`
	listArgs := append(args, pageSize, offset)

	rows, err := m.conn.Query(listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var jobs []*VtTrainingJobs
	for rows.Next() {
		var job VtTrainingJobs
		err := rows.Scan(&job.Id, &job.Name, &job.DisplayName, &job.Description, &job.JobType, &job.Framework, &job.Status, &job.Phase, &job.Priority, &job.GpuCount, &job.SubmittedAt, &job.StartTime, &job.EndTime, &job.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		jobs = append(jobs, &job)
	}

	return jobs, total, nil
}

func (m *vtTrainingJobsModel) GetByStatus(status string) ([]*VtTrainingJobs, error) {
	query := `SELECT id, name, status, phase, volcano_job_name, cluster_name FROM vt_training_jobs WHERE status = ? AND deleted_at IS NULL`
	rows, err := m.conn.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*VtTrainingJobs
	for rows.Next() {
		var job VtTrainingJobs
		err := rows.Scan(&job.Id, &job.Name, &job.Status, &job.Phase, &job.VolcanoJobName, &job.ClusterName)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, &job)
	}

	return jobs, nil
}
