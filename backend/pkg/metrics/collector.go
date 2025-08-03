package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"api/model"
	"github.com/zeromicro/go-zero/core/logx"
)

// MetricsCollector 训练指标收集器
type MetricsCollector struct {
	logger   logx.Logger
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
}

// JobMetrics 作业指标数据
type JobMetrics struct {
	JobID     int64     `json:"job_id"`
	JobName   string    `json:"job_name"`
	Timestamp time.Time `json:"timestamp"`
	Phase     string    `json:"phase"`

	// 资源指标
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	GPUUsage    float64 `json:"gpu_usage"`
	GPUMemory   float64 `json:"gpu_memory"`
	DiskUsage   float64 `json:"disk_usage"`
	NetworkRX   int64   `json:"network_rx"`
	NetworkTX   int64   `json:"network_tx"`

	// 训练指标
	Epoch        int     `json:"epoch"`
	Step         int     `json:"step"`
	Loss         float64 `json:"loss"`
	Accuracy     float64 `json:"accuracy"`
	LearningRate float64 `json:"learning_rate"`
	BatchSize    int     `json:"batch_size"`
	Throughput   float64 `json:"throughput"`

	// 其他指标
	Duration   int64   `json:"duration"`
	Progress   float64 `json:"progress"`
	ErrorCount int     `json:"error_count"`
	Warnings   int     `json:"warnings"`
}

// MetricsStorage 指标存储接口
type MetricsStorage interface {
	Store(metrics *JobMetrics) error
	Query(jobID int64, startTime, endTime time.Time) ([]*JobMetrics, error)
	GetLatest(jobID int64) (*JobMetrics, error)
	Delete(jobID int64) error
}

// NewMetricsCollector 创建指标收集器
func NewMetricsCollector(interval time.Duration) *MetricsCollector {
	ctx, cancel := context.WithCancel(context.Background())

	return &MetricsCollector{
		logger:   logx.WithContext(ctx),
		interval: interval,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start 启动指标收集
func (mc *MetricsCollector) Start(jobModel model.VtTrainingJobsModel, storage MetricsStorage) {
	ticker := time.NewTicker(mc.interval)
	defer ticker.Stop()

	mc.logger.Infof("训练指标收集器启动，收集间隔: %v", mc.interval)

	for {
		select {
		case <-mc.ctx.Done():
			mc.logger.Info("训练指标收集器停止")
			return
		case <-ticker.C:
			mc.collectMetrics(jobModel, storage)
		}
	}
}

// Stop 停止指标收集
func (mc *MetricsCollector) Stop() {
	if mc.cancel != nil {
		mc.cancel()
	}
}

// collectMetrics 收集所有运行中作业的指标
func (mc *MetricsCollector) collectMetrics(jobModel model.VtTrainingJobsModel, storage MetricsStorage) {
	// 获取所有运行中的作业
	runningJobs, err := jobModel.GetByStatus("running")
	if err != nil {
		mc.logger.Errorf("获取运行中作业失败: %v", err)
		return
	}

	mc.logger.Infof("收集 %d 个运行中作业的指标", len(runningJobs))

	// 并发收集指标
	for _, job := range runningJobs {
		go mc.collectJobMetrics(job, storage)
	}
}

// collectJobMetrics 收集单个作业的指标
func (mc *MetricsCollector) collectJobMetrics(job *model.VtTrainingJobs, storage MetricsStorage) {
	metrics := &JobMetrics{
		JobID:     job.Id,
		JobName:   job.Name,
		Timestamp: time.Now(),
		Phase:     job.Phase,
	}

	// 收集资源指标
	if err := mc.collectResourceMetrics(job, metrics); err != nil {
		mc.logger.Errorf("收集作业 %s 资源指标失败: %v", job.Name, err)
	}

	// 收集训练指标
	if err := mc.collectTrainingMetrics(job, metrics); err != nil {
		mc.logger.Errorf("收集作业 %s 训练指标失败: %v", job.Name, err)
	}

	// 计算派生指标
	mc.calculateDerivedMetrics(job, metrics)

	// 存储指标
	if err := storage.Store(metrics); err != nil {
		mc.logger.Errorf("存储作业 %s 指标失败: %v", job.Name, err)
	} else {
		mc.logger.Debugf("成功收集作业 %s 指标", job.Name)
	}
}

// collectResourceMetrics 收集资源指标
func (mc *MetricsCollector) collectResourceMetrics(job *model.VtTrainingJobs, metrics *JobMetrics) error {
	// TODO: 实际环境中需要从Kubernetes API或监控系统获取指标
	// 这里提供模拟数据

	// 模拟CPU使用率 (50-90%)
	metrics.CPUUsage = 50.0 + float64(time.Now().Unix()%40)

	// 模拟内存使用率 (60-85%)
	metrics.MemoryUsage = 60.0 + float64(time.Now().Unix()%25)

	// 如果有GPU，模拟GPU指标
	if job.GpuCount > 0 {
		metrics.GPUUsage = 70.0 + float64(time.Now().Unix()%30)
		metrics.GPUMemory = 65.0 + float64(time.Now().Unix()%35)
	}

	// 模拟磁盘和网络使用
	metrics.DiskUsage = 30.0 + float64(time.Now().Unix()%20)
	metrics.NetworkRX = time.Now().Unix() * 1024 * 1024 // MB
	metrics.NetworkTX = time.Now().Unix() * 512 * 1024  // MB

	return nil
}

// collectTrainingMetrics 收集训练指标
func (mc *MetricsCollector) collectTrainingMetrics(job *model.VtTrainingJobs, metrics *JobMetrics) error {
	// TODO: 实际环境中需要解析训练日志或从训练框架获取指标
	// 这里提供模拟数据

	// 模拟训练进度
	elapsed := time.Since(job.SubmittedAt).Minutes()

	// 模拟epoch和step
	metrics.Epoch = int(elapsed / 10) // 假设每10分钟一个epoch
	metrics.Step = int(elapsed * 100) // 假设每分钟100步

	// 模拟训练指标
	if metrics.Epoch > 0 {
		// 模拟loss下降趋势
		metrics.Loss = 2.0 / (1.0 + float64(metrics.Epoch)*0.1)

		// 模拟accuracy上升趋势
		metrics.Accuracy = 1.0 - metrics.Loss/2.0

		// 模拟学习率衰减
		metrics.LearningRate = 0.01 * (0.95 / float64(1+metrics.Epoch))

		// 模拟其他指标
		metrics.BatchSize = 32
		metrics.Throughput = 50.0 + float64(time.Now().Unix()%20)
	}

	return nil
}

// calculateDerivedMetrics 计算派生指标
func (mc *MetricsCollector) calculateDerivedMetrics(job *model.VtTrainingJobs, metrics *JobMetrics) {
	// 计算运行时长
	if job.StartTime != nil {
		metrics.Duration = int64(time.Since(*job.StartTime).Seconds())
	}

	// 计算训练进度百分比
	if job.MaxRuntimeSeconds > 0 && metrics.Duration > 0 {
		metrics.Progress = float64(metrics.Duration) / float64(job.MaxRuntimeSeconds) * 100
		if metrics.Progress > 100 {
			metrics.Progress = 100
		}
	} else if metrics.Epoch > 0 {
		// 假设总共需要100个epoch
		estimatedTotalEpochs := 100
		metrics.Progress = float64(metrics.Epoch) / float64(estimatedTotalEpochs) * 100
		if metrics.Progress > 100 {
			metrics.Progress = 100
		}
	}
}

// GetJobMetrics 获取作业指标
func (mc *MetricsCollector) GetJobMetrics(jobID int64, timeRange string, storage MetricsStorage) ([]*JobMetrics, error) {
	// 解析时间范围
	endTime := time.Now()
	var startTime time.Time

	switch timeRange {
	case "1h":
		startTime = endTime.Add(-1 * time.Hour)
	case "6h":
		startTime = endTime.Add(-6 * time.Hour)
	case "1d":
		startTime = endTime.Add(-24 * time.Hour)
	case "7d":
		startTime = endTime.Add(-7 * 24 * time.Hour)
	default:
		startTime = endTime.Add(-1 * time.Hour)
	}

	return storage.Query(jobID, startTime, endTime)
}

// InMemoryMetricsStorage 内存中的指标存储实现
type InMemoryMetricsStorage struct {
	metrics map[int64][]*JobMetrics
	logger  logx.Logger
}

// NewInMemoryMetricsStorage 创建内存指标存储
func NewInMemoryMetricsStorage() *InMemoryMetricsStorage {
	return &InMemoryMetricsStorage{
		metrics: make(map[int64][]*JobMetrics),
		logger:  logx.WithContext(context.Background()),
	}
}

// Store 存储指标
func (s *InMemoryMetricsStorage) Store(metrics *JobMetrics) error {
	if s.metrics[metrics.JobID] == nil {
		s.metrics[metrics.JobID] = make([]*JobMetrics, 0)
	}

	s.metrics[metrics.JobID] = append(s.metrics[metrics.JobID], metrics)

	// 保留最近1000条记录
	if len(s.metrics[metrics.JobID]) > 1000 {
		s.metrics[metrics.JobID] = s.metrics[metrics.JobID][len(s.metrics[metrics.JobID])-1000:]
	}

	s.logger.Debugf("存储作业 %d 指标，当前共有 %d 条记录", metrics.JobID, len(s.metrics[metrics.JobID]))
	return nil
}

// Query 查询指标
func (s *InMemoryMetricsStorage) Query(jobID int64, startTime, endTime time.Time) ([]*JobMetrics, error) {
	jobMetrics, exists := s.metrics[jobID]
	if !exists {
		return []*JobMetrics{}, nil
	}

	var result []*JobMetrics
	for _, metric := range jobMetrics {
		if metric.Timestamp.After(startTime) && metric.Timestamp.Before(endTime) {
			result = append(result, metric)
		}
	}

	s.logger.Debugf("查询作业 %d 指标，时间范围 %v-%v，返回 %d 条记录",
		jobID, startTime.Format("15:04:05"), endTime.Format("15:04:05"), len(result))
	return result, nil
}

// GetLatest 获取最新指标
func (s *InMemoryMetricsStorage) GetLatest(jobID int64) (*JobMetrics, error) {
	jobMetrics, exists := s.metrics[jobID]
	if !exists || len(jobMetrics) == 0 {
		return nil, fmt.Errorf("作业 %d 没有指标数据", jobID)
	}

	return jobMetrics[len(jobMetrics)-1], nil
}

// Delete 删除指标
func (s *InMemoryMetricsStorage) Delete(jobID int64) error {
	delete(s.metrics, jobID)
	s.logger.Infof("删除作业 %d 的所有指标数据", jobID)
	return nil
}

// ExportMetrics 导出指标为JSON
func (s *InMemoryMetricsStorage) ExportMetrics(jobID int64) (string, error) {
	jobMetrics, exists := s.metrics[jobID]
	if !exists {
		return "", fmt.Errorf("作业 %d 没有指标数据", jobID)
	}

	data, err := json.MarshalIndent(jobMetrics, "", "  ")
	if err != nil {
		return "", fmt.Errorf("序列化指标数据失败: %v", err)
	}

	return string(data), nil
}
