package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api/internal/service"
	"api/pkg/alerting"
	"api/pkg/monitoring"
	"api/pkg/notification"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("c", "etc/monitoring.yaml", "配置文件路径")

// Config 应用配置
type Config struct {
	// 数据库配置
	Database struct {
		Host         string `json:"host"`
		Port         int    `json:"port"`
		Username     string `json:"username"`
		Password     string `json:"password"`
		Database     string `json:"database"`
		MaxOpenConns int    `json:"max_open_conns"`
		MaxIdleConns int    `json:"max_idle_conns"`
	} `json:"database"`

	// 监控配置
	Monitoring struct {
		// 指标收集配置
		Metrics struct {
			PrometheusURL          string        `json:"prometheus_url"`
			CollectInterval        time.Duration `json:"collect_interval"`
			SystemMetricsEnabled   bool          `json:"system_metrics_enabled"`
			BusinessMetricsEnabled bool          `json:"business_metrics_enabled"`
			PrometheusEnabled      bool          `json:"prometheus_enabled"`
			MetricsPort            int           `json:"metrics_port"`
			EnableBuiltinMetrics   bool          `json:"enable_builtin_metrics"`
			MaxRetries             int           `json:"max_retries"`
			Timeout                time.Duration `json:"timeout"`
		} `json:"metrics"`

		// 告警配置
		Alerts struct {
			EvaluationInterval      time.Duration `json:"evaluation_interval"`
			MaxConcurrentRules      int           `json:"max_concurrent_rules"`
			AlertRetentionDays      int           `json:"alert_retention_days"`
			EnableGrouping          bool          `json:"enable_grouping"`
			EnableSuppression       bool          `json:"enable_suppression"`
			DefaultThrottleMinutes  int           `json:"default_throttle_minutes"`
			AnomalyDetectionEnabled bool          `json:"anomaly_detection_enabled"`
		} `json:"alerts"`

		// 通知配置
		Notifications struct {
			MaxQueueSize         int           `json:"max_queue_size"`
			MaxConcurrentSenders int           `json:"max_concurrent_senders"`
			RetryMaxAttempts     int           `json:"retry_max_attempts"`
			RetryBackoffSeconds  int           `json:"retry_backoff_seconds"`
			RateLimitPerMinute   int           `json:"rate_limit_per_minute"`
			TimeoutSeconds       int           `json:"timeout_seconds"`
			FailedRetentionDays  int           `json:"failed_retention_days"`
			EnableDeduplication  bool          `json:"enable_deduplication"`
			DeduplicationWindow  time.Duration `json:"deduplication_window"`
		} `json:"notifications"`

		// 系统配置
		System struct {
			EnableMetrics       bool          `json:"enable_metrics"`
			EnableAlerts        bool          `json:"enable_alerts"`
			EnableNotifications bool          `json:"enable_notifications"`
			HealthCheckInterval time.Duration `json:"health_check_interval"`
			AutoRestart         bool          `json:"auto_restart"`
			MaxRestartAttempts  int           `json:"max_restart_attempts"`
		} `json:"system"`
	} `json:"monitoring"`

	// 服务端口
	Port int `json:"port"`
}

func main() {
	flag.Parse()

	// 加载配置
	var c Config
	conf.MustLoad(*configFile, &c)

	// 初始化日志
	logx.MustSetup(logx.LogConf{
		ServiceName: "volctrain-monitoring",
		Mode:        "file",
		Path:        "logs",
		Level:       "info",
	})

	logger := logx.WithContext(context.Background())
	logger.Info("启动VolcTrain监控服务")

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.Username, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 设置数据库连接池
	db.SetMaxOpenConns(c.Database.MaxOpenConns)
	db.SetMaxIdleConns(c.Database.MaxIdleConns)

	// 测试数据库连接
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}

	logger.Info("数据库连接成功")

	// 创建监控服务配置
	monitoringConfig := &service.MonitoringConfig{
		MetricsConfig: &monitoring.CollectorConfig{
			PrometheusURL:          c.Monitoring.Metrics.PrometheusURL,
			CollectInterval:        c.Monitoring.Metrics.CollectInterval,
			SystemMetricsEnabled:   c.Monitoring.Metrics.SystemMetricsEnabled,
			BusinessMetricsEnabled: c.Monitoring.Metrics.BusinessMetricsEnabled,
			PrometheusEnabled:      c.Monitoring.Metrics.PrometheusEnabled,
			MetricsPort:            c.Monitoring.Metrics.MetricsPort,
			EnableBuiltinMetrics:   c.Monitoring.Metrics.EnableBuiltinMetrics,
			MaxRetries:             c.Monitoring.Metrics.MaxRetries,
			Timeout:                c.Monitoring.Metrics.Timeout,
		},
		AlertConfig: &alerting.AlertEngineConfig{
			EvaluationInterval:      c.Monitoring.Alerts.EvaluationInterval,
			MaxConcurrentRules:      c.Monitoring.Alerts.MaxConcurrentRules,
			AlertRetentionDays:      c.Monitoring.Alerts.AlertRetentionDays,
			EnableGrouping:          c.Monitoring.Alerts.EnableGrouping,
			EnableSuppression:       c.Monitoring.Alerts.EnableSuppression,
			DefaultThrottleMinutes:  c.Monitoring.Alerts.DefaultThrottleMinutes,
			AnomalyDetectionEnabled: c.Monitoring.Alerts.AnomalyDetectionEnabled,
		},
		NotificationConfig: &notification.NotificationConfig{
			MaxQueueSize:         c.Monitoring.Notifications.MaxQueueSize,
			MaxConcurrentSenders: c.Monitoring.Notifications.MaxConcurrentSenders,
			RetryMaxAttempts:     c.Monitoring.Notifications.RetryMaxAttempts,
			RetryBackoffSeconds:  c.Monitoring.Notifications.RetryBackoffSeconds,
			RateLimitPerMinute:   c.Monitoring.Notifications.RateLimitPerMinute,
			TimeoutSeconds:       c.Monitoring.Notifications.TimeoutSeconds,
			FailedRetentionDays:  c.Monitoring.Notifications.FailedRetentionDays,
			EnableDeduplication:  c.Monitoring.Notifications.EnableDeduplication,
			DeduplicationWindow:  c.Monitoring.Notifications.DeduplicationWindow,
		},
		EnableMetrics:       c.Monitoring.System.EnableMetrics,
		EnableAlerts:        c.Monitoring.System.EnableAlerts,
		EnableNotifications: c.Monitoring.System.EnableNotifications,
		HealthCheckInterval: c.Monitoring.System.HealthCheckInterval,
		AutoRestart:         c.Monitoring.System.AutoRestart,
		MaxRestartAttempts:  c.Monitoring.System.MaxRestartAttempts,
	}

	// 创建监控服务
	monitoringService := service.NewMonitoringService(db, monitoringConfig)

	// 启动监控服务
	if err := monitoringService.Start(); err != nil {
		log.Fatalf("启动监控服务失败: %v", err)
	}

	// 初始化内置规则和模板
	builtinManager := monitoring.NewBuiltinRulesTemplates(db)
	if err := builtinManager.Initialize(); err != nil {
		logger.Errorf("初始化内置规则和模板失败: %v", err)
	} else {
		logger.Info("内置规则和模板初始化完成")
	}

	// 设置HTTP API服务器
	mux := http.NewServeMux()

	// 健康检查端点
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		status := monitoringService.GetSystemStatus()
		w.Header().Set("Content-Type", "application/json")
		if status.OverallStatus == "healthy" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		fmt.Fprintf(w, `{"status": "%s", "timestamp": "%s"}`,
			status.OverallStatus, time.Now().Format(time.RFC3339))
	})

	// 就绪检查端点
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "ready", "timestamp": "%s"}`,
			time.Now().Format(time.RFC3339))
	})

	// 系统状态端点
	mux.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		status := monitoringService.GetSystemStatus()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// 简化的JSON输出
		fmt.Fprintf(w, `{
			"overall_status": "%s",
			"components": %d,
			"active_alerts": %d,
			"uptime": "%.0f",
			"last_health_check": "%s"
		}`,
			status.OverallStatus,
			len(status.Components),
			status.ActiveAlerts,
			status.Uptime.Seconds(),
			status.LastHealthCheck.Format(time.RFC3339))
	})

	// 配置重新加载端点
	mux.HandleFunc("/api/v1/reload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := monitoringService.ReloadConfiguration(); err != nil {
			http.Error(w, fmt.Sprintf("重新加载配置失败: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "success", "timestamp": "%s"}`,
			time.Now().Format(time.RFC3339))
	})

	// 启动HTTP服务器
	port := c.Port
	if port == 0 {
		port = 8080
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	// 启动服务器
	go func() {
		logger.Infof("HTTP服务器启动在端口 %d", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP服务器启动失败: %v", err)
		}
	}()

	logger.Info("VolcTrain监控服务启动完成")

	// 等待信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("收到退出信号，开始优雅关闭...")

	// 关闭HTTP服务器
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("HTTP服务器关闭失败: %v", err)
	}

	// 停止监控服务
	monitoringService.Stop()

	logger.Info("VolcTrain监控服务已停止")
}
