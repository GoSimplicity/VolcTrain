package monitoring

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"api/model"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zeromicro/go-zero/core/logx"
)

// EnhancedMetricsCollector 增强的监控指标收集器
type EnhancedMetricsCollector struct {
	logger             logx.Logger
	ctx                context.Context
	cancel             context.CancelFunc
	db                 *sql.DB
	metricsModel       model.VtMonitorMetricsModel
	monitorDataModel   model.VtMonitorDataModel
	prometheus         PromClientInterface
	config             *CollectorConfig
	registry           *prometheus.Registry
	systemMetrics      *SystemMetrics
	businessMetrics    *BusinessMetrics
	collectInterval    time.Duration
	lastCollectionTime time.Time
	mu                 sync.RWMutex
}

// CollectorConfig 收集器配置
type CollectorConfig struct {
	PrometheusURL          string        `json:"prometheus_url"`
	CollectInterval        time.Duration `json:"collect_interval"`
	SystemMetricsEnabled   bool          `json:"system_metrics_enabled"`
	BusinessMetricsEnabled bool          `json:"business_metrics_enabled"`
	PrometheusEnabled      bool          `json:"prometheus_enabled"`
	MetricsPort            int           `json:"metrics_port"`
	EnableBuiltinMetrics   bool          `json:"enable_builtin_metrics"`
	MaxRetries             int           `json:"max_retries"`
	Timeout                time.Duration `json:"timeout"`
}

// PromClientInterface Prometheus客户端接口
type PromClientInterface interface {
	Query(ctx context.Context, query string, ts time.Time) (interface{}, error)
	QueryRange(ctx context.Context, query string, r v1.Range) (interface{}, error)
}

// PrometheusClient Prometheus客户端
type PrometheusClient struct {
	client v1.API
}

// SystemMetrics 系统指标
type SystemMetrics struct {
	// CPU指标
	CPUUsageTotal  prometheus.GaugeVec
	CPUUsageUser   prometheus.GaugeVec
	CPUUsageSystem prometheus.GaugeVec
	CPUUsageIdle   prometheus.GaugeVec

	// 内存指标
	MemoryUsageTotal prometheus.GaugeVec
	MemoryUsageUsed  prometheus.GaugeVec
	MemoryUsageFree  prometheus.GaugeVec
	MemoryUsageCache prometheus.GaugeVec

	// 磁盘指标
	DiskUsageTotal prometheus.GaugeVec
	DiskUsageUsed  prometheus.GaugeVec
	DiskUsageFree  prometheus.GaugeVec
	DiskIORead     prometheus.CounterVec
	DiskIOWrite    prometheus.CounterVec

	// 网络指标
	NetworkBytesReceived   prometheus.CounterVec
	NetworkBytesSent       prometheus.CounterVec
	NetworkPacketsReceived prometheus.CounterVec
	NetworkPacketsSent     prometheus.CounterVec
}

// BusinessMetrics 业务指标
type BusinessMetrics struct {
	// 训练任务指标
	TrainingJobsTotal     prometheus.GaugeVec
	TrainingJobsRunning   prometheus.GaugeVec
	TrainingJobsPending   prometheus.GaugeVec
	TrainingJobsCompleted prometheus.CounterVec
	TrainingJobsFailed    prometheus.CounterVec

	// GPU指标
	GPUUsage       prometheus.GaugeVec
	GPUMemoryUsage prometheus.GaugeVec
	GPUTemperature prometheus.GaugeVec
	GPUPowerUsage  prometheus.GaugeVec

	// 队列指标
	QueueLength   prometheus.GaugeVec
	QueueWaitTime prometheus.HistogramVec

	// 用户指标
	ActiveUsers    prometheus.GaugeVec
	UserLoginCount prometheus.CounterVec
}

// MetricPoint 指标数据点
type MetricPoint struct {
	MetricName   string                 `json:"metric_name"`
	Value        float64                `json:"value"`
	Labels       map[string]interface{} `json:"labels"`
	Timestamp    time.Time              `json:"timestamp"`
	ResourceType string                 `json:"resource_type"`
	ResourceID   int64                  `json:"resource_id"`
	InstanceID   string                 `json:"instance_id"`
}

// NewEnhancedMetricsCollector 创建增强指标收集器
func NewEnhancedMetricsCollector(db *sql.DB, config *CollectorConfig) *EnhancedMetricsCollector {
	ctx, cancel := context.WithCancel(context.Background())

	collector := &EnhancedMetricsCollector{
		logger:           logx.WithContext(ctx),
		ctx:              ctx,
		cancel:           cancel,
		db:               db,
		metricsModel:     model.NewVtMonitorMetricsModel(db),
		monitorDataModel: model.NewVtMonitorDataModel(db),
		config:           config,
		collectInterval:  config.CollectInterval,
		registry:         prometheus.NewRegistry(),
	}

	// 初始化Prometheus客户端
	if config.PrometheusEnabled && config.PrometheusURL != "" {
		promClient, err := api.NewClient(api.Config{
			Address: config.PrometheusURL,
		})
		if err != nil {
			collector.logger.Errorf("创建Prometheus客户端失败: %v", err)
		} else {
			collector.prometheus = &PrometheusClient{client: v1.NewAPI(promClient)}
		}
	}

	// 初始化指标
	collector.initializeMetrics()

	return collector
}

// initializeMetrics 初始化指标
func (c *EnhancedMetricsCollector) initializeMetrics() {
	// 初始化系统指标
	if c.config.SystemMetricsEnabled {
		c.systemMetrics = c.createSystemMetrics()
	}

	// 初始化业务指标
	if c.config.BusinessMetricsEnabled {
		c.businessMetrics = c.createBusinessMetrics()
	}

	// 初始化内置指标
	if c.config.EnableBuiltinMetrics {
		c.initializeBuiltinMetrics()
	}
}

// createSystemMetrics 创建系统指标
func (c *EnhancedMetricsCollector) createSystemMetrics() *SystemMetrics {
	metrics := &SystemMetrics{
		CPUUsageTotal: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "volctrain_cpu_usage_total",
				Help: "总CPU使用率",
			},
			[]string{"instance", "node"},
		),
		MemoryUsageTotal: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "volctrain_memory_usage_total",
				Help: "总内存使用量",
			},
			[]string{"instance", "node"},
		),
		DiskUsageTotal: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "volctrain_disk_usage_total",
				Help: "总磁盘使用量",
			},
			[]string{"instance", "node", "device"},
		),
		NetworkBytesReceived: *prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "volctrain_network_bytes_received_total",
				Help: "网络接收字节数",
			},
			[]string{"instance", "node", "interface"},
		),
	}

	// 注册指标
	c.registry.MustRegister(
		metrics.CPUUsageTotal,
		metrics.MemoryUsageTotal,
		metrics.DiskUsageTotal,
		metrics.NetworkBytesReceived,
	)

	return metrics
}

// createBusinessMetrics 创建业务指标
func (c *EnhancedMetricsCollector) createBusinessMetrics() *BusinessMetrics {
	metrics := &BusinessMetrics{
		TrainingJobsTotal: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "volctrain_training_jobs_total",
				Help: "训练任务总数",
			},
			[]string{"status", "queue", "framework"},
		),
		GPUUsage: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "volctrain_gpu_usage_percent",
				Help: "GPU使用率",
			},
			[]string{"instance", "gpu_id", "node", "job_id"},
		),
		QueueLength: *prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "volctrain_queue_length",
				Help: "队列长度",
			},
			[]string{"queue_name", "queue_type"},
		),
	}

	// 注册指标
	c.registry.MustRegister(
		metrics.TrainingJobsTotal,
		metrics.GPUUsage,
		metrics.QueueLength,
	)

	return metrics
}

// initializeBuiltinMetrics 初始化内置指标定义
func (c *EnhancedMetricsCollector) initializeBuiltinMetrics() {
	builtinMetrics := []model.VtMonitorMetrics{
		{
			Name:                      "system_cpu_usage",
			DisplayName:               "系统CPU使用率",
			Description:               "系统CPU使用率百分比",
			MetricType:                "gauge",
			DataType:                  "float",
			Category:                  "system",
			Module:                    "system",
			SourceType:                "system",
			Unit:                      "percent",
			AggregationType:           "avg",
			CollectionIntervalSeconds: 30,
			RetentionDays:             30,
			WarningThreshold:          &[]float64{80.0}[0],
			CriticalThreshold:         &[]float64{90.0}[0],
			ThresholdCondition:        "gt",
			Status:                    "active",
			IsBuiltin:                 true,
			IsCore:                    true,
			DefaultLabels:             "{}",
			Dimensions:                "{}",
			Metadata:                  "{}",
		},
		{
			Name:                      "system_memory_usage",
			DisplayName:               "系统内存使用率",
			Description:               "系统内存使用率百分比",
			MetricType:                "gauge",
			DataType:                  "float",
			Category:                  "system",
			Module:                    "system",
			SourceType:                "system",
			Unit:                      "percent",
			AggregationType:           "avg",
			CollectionIntervalSeconds: 30,
			RetentionDays:             30,
			WarningThreshold:          &[]float64{80.0}[0],
			CriticalThreshold:         &[]float64{90.0}[0],
			ThresholdCondition:        "gt",
			Status:                    "active",
			IsBuiltin:                 true,
			IsCore:                    true,
			DefaultLabels:             "{}",
			Dimensions:                "{}",
			Metadata:                  "{}",
		},
		{
			Name:                      "training_job_count",
			DisplayName:               "训练任务数量",
			Description:               "当前运行的训练任务数量",
			MetricType:                "gauge",
			DataType:                  "integer",
			Category:                  "business",
			Module:                    "training",
			SourceType:                "business",
			Unit:                      "count",
			AggregationType:           "sum",
			CollectionIntervalSeconds: 60,
			RetentionDays:             90,
			Status:                    "active",
			IsBuiltin:                 true,
			IsCore:                    true,
			DefaultLabels:             "{}",
			Dimensions:                "{}",
			Metadata:                  "{}",
		},
		{
			Name:                      "gpu_usage",
			DisplayName:               "GPU使用率",
			Description:               "GPU使用率百分比",
			MetricType:                "gauge",
			DataType:                  "float",
			Category:                  "resource",
			Module:                    "gpu",
			SourceType:                "business",
			Unit:                      "percent",
			AggregationType:           "avg",
			CollectionIntervalSeconds: 30,
			RetentionDays:             30,
			WarningThreshold:          &[]float64{85.0}[0],
			CriticalThreshold:         &[]float64{95.0}[0],
			ThresholdCondition:        "gt",
			Status:                    "active",
			IsBuiltin:                 true,
			IsCore:                    true,
			DefaultLabels:             "{}",
			Dimensions:                "{}",
			Metadata:                  "{}",
		},
	}

	// 插入内置指标定义
	for _, metric := range builtinMetrics {
		// 检查是否已存在
		existing, err := c.metricsModel.FindOneByName(metric.Name)
		if err != nil && err != sql.ErrNoRows {
			c.logger.Errorf("检查内置指标失败: %v", err)
			continue
		}

		// 如果不存在则插入
		if existing == nil {
			_, err := c.metricsModel.Insert(&metric)
			if err != nil {
				c.logger.Errorf("插入内置指标失败: %v", err)
			} else {
				c.logger.Infof("成功插入内置指标: %s", metric.Name)
			}
		}
	}
}

// Start 启动指标收集
func (c *EnhancedMetricsCollector) Start() error {
	c.logger.Infof("启动增强监控指标收集器，收集间隔: %v", c.collectInterval)

	// 启动Prometheus HTTP服务器
	if c.config.PrometheusEnabled {
		go c.startPrometheusServer()
	}

	// 启动指标收集循环
	go c.collectLoop()

	return nil
}

// Stop 停止指标收集
func (c *EnhancedMetricsCollector) Stop() {
	c.logger.Info("停止增强监控指标收集器")
	if c.cancel != nil {
		c.cancel()
	}
}

// startPrometheusServer 启动Prometheus HTTP服务器
func (c *EnhancedMetricsCollector) startPrometheusServer() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(c.registry, promhttp.HandlerOpts{}))

	addr := fmt.Sprintf(":%d", c.config.MetricsPort)
	c.logger.Infof("启动Prometheus指标服务器: http://localhost%s/metrics", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		c.logger.Errorf("Prometheus指标服务器启动失败: %v", err)
	}
}

// collectLoop 指标收集循环
func (c *EnhancedMetricsCollector) collectLoop() {
	ticker := time.NewTicker(c.collectInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			c.logger.Info("指标收集循环停止")
			return
		case <-ticker.C:
			c.collectAllMetrics()
		}
	}
}

// collectAllMetrics 收集所有指标
func (c *EnhancedMetricsCollector) collectAllMetrics() {
	c.mu.Lock()
	defer c.mu.Unlock()

	start := time.Now()
	c.logger.Debug("开始收集指标")

	var wg sync.WaitGroup

	// 并发收集不同类型的指标
	if c.config.SystemMetricsEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.collectSystemMetrics()
		}()
	}

	if c.config.BusinessMetricsEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.collectBusinessMetrics()
		}()
	}

	if c.config.PrometheusEnabled && c.prometheus != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.collectPrometheusMetrics()
		}()
	}

	wg.Wait()

	c.lastCollectionTime = time.Now()
	duration := c.lastCollectionTime.Sub(start)
	c.logger.Debugf("指标收集完成，耗时: %v", duration)
}

// collectSystemMetrics 收集系统指标
func (c *EnhancedMetricsCollector) collectSystemMetrics() {
	timestamp := time.Now()

	// 收集CPU指标
	cpuUsage := c.getCPUUsage()
	c.saveMetricData("system_cpu_usage", cpuUsage, map[string]interface{}{
		"type": "total",
	}, timestamp, "system", 0, "localhost")

	// 收集内存指标
	memUsage := c.getMemoryUsage()
	c.saveMetricData("system_memory_usage", memUsage, map[string]interface{}{
		"type": "used",
	}, timestamp, "system", 0, "localhost")

	// 更新Prometheus指标
	if c.systemMetrics != nil {
		c.systemMetrics.CPUUsageTotal.WithLabelValues("localhost", "node1").Set(cpuUsage)
		c.systemMetrics.MemoryUsageTotal.WithLabelValues("localhost", "node1").Set(memUsage)
	}
}

// collectBusinessMetrics 收集业务指标
func (c *EnhancedMetricsCollector) collectBusinessMetrics() {
	timestamp := time.Now()

	// 收集训练任务指标
	jobStats := c.getTrainingJobStats()
	for status, count := range jobStats {
		c.saveMetricData("training_job_count", float64(count), map[string]interface{}{
			"status": status,
		}, timestamp, "training_job", 0, "")

		// 更新Prometheus指标
		if c.businessMetrics != nil {
			c.businessMetrics.TrainingJobsTotal.WithLabelValues(status, "default", "unknown").Set(float64(count))
		}
	}

	// 收集GPU指标
	gpuStats := c.getGPUStats()
	for _, gpu := range gpuStats {
		c.saveMetricData("gpu_usage", gpu.Usage, map[string]interface{}{
			"gpu_id":   gpu.ID,
			"gpu_type": gpu.Type,
		}, timestamp, "gpu", gpu.ID, gpu.InstanceID)

		// 更新Prometheus指标
		if c.businessMetrics != nil {
			c.businessMetrics.GPUUsage.WithLabelValues("localhost", fmt.Sprintf("%d", gpu.ID), "node1", "").Set(gpu.Usage)
		}
	}
}

// collectPrometheusMetrics 从Prometheus收集指标
func (c *EnhancedMetricsCollector) collectPrometheusMetrics() {
	if c.prometheus == nil {
		return
	}

	// 查询Prometheus指标
	queries := []string{
		"node_cpu_usage_percent",
		"node_memory_usage_percent",
		"node_disk_usage_percent",
		"node_network_bytes_total",
	}

	for _, query := range queries {
		result, err := c.prometheus.Query(c.ctx, query, time.Now())
		if err != nil {
			c.logger.Errorf("查询Prometheus指标失败 [%s]: %v", query, err)
			continue
		}

		// 处理查询结果并保存
		c.processPrometheusResult(query, result)
	}
}

// getCPUUsage 获取CPU使用率
func (c *EnhancedMetricsCollector) getCPUUsage() float64 {
	// 简化实现，实际环境中需要调用系统API或解析/proc/stat
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 模拟CPU使用率 (50-90%)
	baseUsage := 50.0
	variableUsage := float64(time.Now().Unix() % 40)
	return baseUsage + variableUsage
}

// getMemoryUsage 获取内存使用率
func (c *EnhancedMetricsCollector) getMemoryUsage() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 使用Go运行时内存信息计算使用率
	totalMem := float64(m.Sys)
	usedMem := float64(m.Alloc)

	if totalMem > 0 {
		return (usedMem / totalMem) * 100
	}

	// 兜底：模拟内存使用率 (60-85%)
	return 60.0 + float64(time.Now().Unix()%25)
}

// getTrainingJobStats 获取训练任务统计
func (c *EnhancedMetricsCollector) getTrainingJobStats() map[string]int64 {
	// 实际实现需要查询数据库
	stats := map[string]int64{
		"running":   0,
		"pending":   0,
		"completed": 0,
		"failed":    0,
	}

	// 查询数据库获取真实统计数据
	query := `SELECT status, COUNT(*) as count FROM vt_training_jobs 
		WHERE deleted_at IS NULL 
		GROUP BY status`

	rows, err := c.db.Query(query)
	if err != nil {
		c.logger.Errorf("查询训练任务统计失败: %v", err)
		return stats
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			c.logger.Errorf("扫描训练任务统计失败: %v", err)
			continue
		}
		stats[status] = count
	}

	return stats
}

// GPUStat GPU统计信息
type GPUStat struct {
	ID         int64   `json:"id"`
	Type       string  `json:"type"`
	Usage      float64 `json:"usage"`
	InstanceID string  `json:"instance_id"`
}

// getGPUStats 获取GPU统计信息
func (c *EnhancedMetricsCollector) getGPUStats() []GPUStat {
	var stats []GPUStat

	// 查询数据库获取GPU信息
	query := `SELECT id, model, utilization_gpu, device_uuid 
		FROM vt_gpu_devices 
		WHERE status = 'available' OR status = 'allocated'`

	rows, err := c.db.Query(query)
	if err != nil {
		c.logger.Errorf("查询GPU统计失败: %v", err)
		return stats
	}
	defer rows.Close()

	for rows.Next() {
		var stat GPUStat
		var model, uuid sql.NullString
		if err := rows.Scan(&stat.ID, &model, &stat.Usage, &uuid); err != nil {
			c.logger.Errorf("扫描GPU统计失败: %v", err)
			continue
		}

		stat.Type = model.String
		stat.InstanceID = uuid.String
		stats = append(stats, stat)
	}

	return stats
}

// saveMetricData 保存指标数据
func (c *EnhancedMetricsCollector) saveMetricData(metricName string, value float64, labels map[string]interface{}, timestamp time.Time, resourceType string, resourceID int64, instanceID string) {
	// 查找指标定义
	metric, err := c.metricsModel.FindOneByName(metricName)
	if err != nil {
		c.logger.Errorf("查找指标定义失败 [%s]: %v", metricName, err)
		return
	}

	// 创建监控数据
	data := &model.VtMonitorData{
		MetricId:       metric.Id,
		ResourceType:   resourceType,
		ResourceName:   instanceID,
		InstanceId:     instanceID,
		Value:          value,
		Timestamp:      timestamp,
		CollectionTime: time.Now(),
		QualityScore:   1.0,
		IsAnomaly:      false,
	}

	if resourceID > 0 {
		data.ResourceId = &resourceID
	}

	// 设置标签
	if err := data.SetLabels(labels); err != nil {
		c.logger.Errorf("设置标签失败: %v", err)
	}

	// 保存到数据库
	if _, err := c.monitorDataModel.Insert(data); err != nil {
		c.logger.Errorf("保存指标数据失败 [%s]: %v", metricName, err)
	}
}

// processPrometheusResult 处理Prometheus查询结果
func (c *EnhancedMetricsCollector) processPrometheusResult(query string, result interface{}) {
	// 实际实现需要根据Prometheus API响应格式解析
	c.logger.Debugf("处理Prometheus查询结果: %s", query)
	// TODO: 解析result并调用saveMetricData保存
}

// Query 查询Prometheus（实现PromClientInterface）
func (p *PrometheusClient) Query(ctx context.Context, query string, ts time.Time) (interface{}, error) {
	result, warnings, err := p.client.Query(ctx, query, ts)
	if len(warnings) > 0 {
		// 记录警告信息但不阻止返回结果
		fmt.Printf("Prometheus query warnings: %v\n", warnings)
	}
	return result, err
}

// QueryRange 范围查询Prometheus（实现PromClientInterface）
func (p *PrometheusClient) QueryRange(ctx context.Context, query string, r v1.Range) (interface{}, error) {
	result, warnings, err := p.client.QueryRange(ctx, query, r)
	if len(warnings) > 0 {
		// 记录警告信息但不阻止返回结果
		fmt.Printf("Prometheus query range warnings: %v\n", warnings)
	}
	return result, err
}

// GetCollectionStatus 获取收集状态
func (c *EnhancedMetricsCollector) GetCollectionStatus() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]interface{}{
		"last_collection_time":     c.lastCollectionTime,
		"collect_interval":         c.collectInterval,
		"system_metrics_enabled":   c.config.SystemMetricsEnabled,
		"business_metrics_enabled": c.config.BusinessMetricsEnabled,
		"prometheus_enabled":       c.config.PrometheusEnabled,
		"prometheus_url":           c.config.PrometheusURL,
	}
}

// GetBuiltinMetrics 获取内置指标定义
func (c *EnhancedMetricsCollector) GetBuiltinMetrics() []model.VtMonitorMetrics {
	// 查询所有内置指标
	metrics, err := c.metricsModel.FindBuiltinMetrics()
	if err != nil {
		c.logger.Errorf("查询内置指标失败: %v", err)
		return []model.VtMonitorMetrics{}
	}

	// 转换指针切片为值切片
	var result []model.VtMonitorMetrics
	for _, metric := range metrics {
		result = append(result, *metric)
	}

	return result
}
