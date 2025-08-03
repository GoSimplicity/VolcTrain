package service

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"api/pkg/alerting"
	"api/pkg/monitoring"
	"api/pkg/notification"
	"github.com/zeromicro/go-zero/core/logx"
)

// MonitoringService 监控服务集成器
type MonitoringService struct {
	logger              logx.Logger
	ctx                 context.Context
	cancel              context.CancelFunc
	db                  *sql.DB
	config              *MonitoringConfig
	metricsCollector    *monitoring.EnhancedMetricsCollector
	alertEngine         *alerting.AlertEngine
	notificationManager *notification.NotificationManager
	healthStatus        map[string]bool
	mu                  sync.RWMutex
}

// MonitoringConfig 监控服务配置
type MonitoringConfig struct {
	// 指标收集配置
	MetricsConfig *monitoring.CollectorConfig `json:"metrics_config"`

	// 告警引擎配置
	AlertConfig *alerting.AlertEngineConfig `json:"alert_config"`

	// 通知管理配置
	NotificationConfig *notification.NotificationConfig `json:"notification_config"`

	// 系统配置
	EnableMetrics       bool          `json:"enable_metrics"`
	EnableAlerts        bool          `json:"enable_alerts"`
	EnableNotifications bool          `json:"enable_notifications"`
	HealthCheckInterval time.Duration `json:"health_check_interval"`
	AutoRestart         bool          `json:"auto_restart"`
	MaxRestartAttempts  int           `json:"max_restart_attempts"`
}

// ComponentStatus 组件状态
type ComponentStatus struct {
	Name         string                 `json:"name"`
	Status       string                 `json:"status"`
	LastUpdate   time.Time              `json:"last_update"`
	Uptime       time.Duration          `json:"uptime"`
	RestartCount int                    `json:"restart_count"`
	ErrorCount   int                    `json:"error_count"`
	LastError    string                 `json:"last_error"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// SystemStatus 系统状态
type SystemStatus struct {
	OverallStatus      string                      `json:"overall_status"`
	Components         map[string]*ComponentStatus `json:"components"`
	StartTime          time.Time                   `json:"start_time"`
	Uptime             time.Duration               `json:"uptime"`
	ActiveAlerts       int                         `json:"active_alerts"`
	TotalNotifications int                         `json:"total_notifications"`
	MetricsCollected   int64                       `json:"metrics_collected"`
	LastHealthCheck    time.Time                   `json:"last_health_check"`
}

// NewMonitoringService 创建监控服务
func NewMonitoringService(db *sql.DB, config *MonitoringConfig) *MonitoringService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &MonitoringService{
		logger:       logx.WithContext(ctx),
		ctx:          ctx,
		cancel:       cancel,
		db:           db,
		config:       config,
		healthStatus: make(map[string]bool),
	}

	return service
}

// Start 启动监控服务
func (s *MonitoringService) Start() error {
	s.logger.Info("启动VolcTrain监控服务集成器")

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	// 启动指标收集器
	if s.config.EnableMetrics {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.startMetricsCollector(); err != nil {
				errChan <- fmt.Errorf("启动指标收集器失败: %v", err)
			}
		}()
	}

	// 启动告警引擎
	if s.config.EnableAlerts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.startAlertEngine(); err != nil {
				errChan <- fmt.Errorf("启动告警引擎失败: %v", err)
			}
		}()
	}

	// 启动通知管理器
	if s.config.EnableNotifications {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.startNotificationManager(); err != nil {
				errChan <- fmt.Errorf("启动通知管理器失败: %v", err)
			}
		}()
	}

	// 等待所有组件启动完成
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// 检查启动错误
	for err := range errChan {
		if err != nil {
			s.Stop()
			return err
		}
	}

	// 启动健康检查
	go s.healthCheckLoop()

	// 集成告警引擎和通知管理器
	if s.alertEngine != nil && s.notificationManager != nil {
		go s.alertNotificationBridge()
	}

	s.logger.Info("监控服务集成器启动完成")
	return nil
}

// Stop 停止监控服务
func (s *MonitoringService) Stop() {
	s.logger.Info("停止VolcTrain监控服务集成器")

	if s.cancel != nil {
		s.cancel()
	}

	// 停止各个组件
	if s.metricsCollector != nil {
		s.metricsCollector.Stop()
	}

	if s.alertEngine != nil {
		s.alertEngine.Stop()
	}

	if s.notificationManager != nil {
		s.notificationManager.Stop()
	}

	s.logger.Info("监控服务集成器已停止")
}

// startMetricsCollector 启动指标收集器
func (s *MonitoringService) startMetricsCollector() error {
	s.logger.Info("启动指标收集器")

	// 创建指标收集器
	s.metricsCollector = monitoring.NewEnhancedMetricsCollector(s.db, s.config.MetricsConfig)

	// 启动收集器
	if err := s.metricsCollector.Start(); err != nil {
		return fmt.Errorf("指标收集器启动失败: %v", err)
	}

	s.mu.Lock()
	s.healthStatus["metrics_collector"] = true
	s.mu.Unlock()

	s.logger.Info("指标收集器启动成功")
	return nil
}

// startAlertEngine 启动告警引擎
func (s *MonitoringService) startAlertEngine() error {
	s.logger.Info("启动告警引擎")

	// 创建告警引擎
	s.alertEngine = alerting.NewAlertEngine(s.db, s.config.AlertConfig)

	// 启动引擎
	if err := s.alertEngine.Start(); err != nil {
		return fmt.Errorf("告警引擎启动失败: %v", err)
	}

	s.mu.Lock()
	s.healthStatus["alert_engine"] = true
	s.mu.Unlock()

	s.logger.Info("告警引擎启动成功")
	return nil
}

// startNotificationManager 启动通知管理器
func (s *MonitoringService) startNotificationManager() error {
	s.logger.Info("启动通知管理器")

	// 创建通知管理器
	s.notificationManager = notification.NewNotificationManager(s.db, s.config.NotificationConfig)

	// 注册通知渠道
	s.registerNotificationChannels()

	// 启动管理器
	if err := s.notificationManager.Start(); err != nil {
		return fmt.Errorf("通知管理器启动失败: %v", err)
	}

	s.mu.Lock()
	s.healthStatus["notification_manager"] = true
	s.mu.Unlock()

	s.logger.Info("通知管理器启动成功")
	return nil
}

// registerNotificationChannels 注册通知渠道
func (s *MonitoringService) registerNotificationChannels() {
	// 这里可以根据配置动态注册不同的通知渠道
	// 实际项目中会从数据库或配置文件中读取渠道配置

	// 示例：注册默认邮件渠道
	emailConfig := map[string]interface{}{
		"smtp_host":    "smtp.example.com",
		"smtp_port":    587,
		"username":     "volctrain@example.com",
		"password":     "password",
		"from":         "volctrain@example.com",
		"from_name":    "VolcTrain监控系统",
		"use_starttls": true,
	}
	emailChannel := notification.NewEmailChannel("default_email", emailConfig)
	s.notificationManager.RegisterChannel(emailChannel)

	// 示例：注册钉钉渠道
	dingtalkConfig := map[string]interface{}{
		"webhook_url":   "https://oapi.dingtalk.com/robot/send?access_token=xxx",
		"secret":        "xxx",
		"message_type":  "markdown",
		"enable_secret": true,
	}
	dingtalkChannel := notification.NewDingTalkChannel("default_dingtalk", dingtalkConfig)
	s.notificationManager.RegisterChannel(dingtalkChannel)
}

// alertNotificationBridge 告警与通知的桥接
func (s *MonitoringService) alertNotificationBridge() {
	s.logger.Info("启动告警通知桥接服务")

	// 这里可以实现更复杂的告警路由逻辑
	// 简化实现：监听告警引擎的通知，转发给通知管理器

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("告警通知桥接服务停止")
			return
		case <-ticker.C:
			// 检查活跃告警并发送通知
			s.processActiveAlerts()
		}
	}
}

// processActiveAlerts 处理活跃告警
func (s *MonitoringService) processActiveAlerts() {
	if s.alertEngine == nil || s.notificationManager == nil {
		return
	}

	// 获取活跃告警
	activeAlerts := s.alertEngine.GetActiveAlerts()

	for _, alert := range activeAlerts {
		// 构造告警通知
		alertNotification := &alerting.AlertNotification{
			Alert:    alert,
			RuleInfo: nil, // 需要从告警引擎获取规则信息
			Action:   "firing",
		}

		// 发送通知
		if err := s.notificationManager.SendNotification(alertNotification); err != nil {
			s.logger.Errorf("发送告警通知失败: %v", err)
		}
	}
}

// healthCheckLoop 健康检查循环
func (s *MonitoringService) healthCheckLoop() {
	ticker := time.NewTicker(s.config.HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("健康检查循环停止")
			return
		case <-ticker.C:
			s.performHealthCheck()
		}
	}
}

// performHealthCheck 执行健康检查
func (s *MonitoringService) performHealthCheck() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查指标收集器
	if s.metricsCollector != nil {
		status := s.metricsCollector.GetCollectionStatus()
		s.healthStatus["metrics_collector"] = s.isHealthy(status)
	}

	// 检查告警引擎
	if s.alertEngine != nil {
		status := s.alertEngine.GetEngineStatus()
		s.healthStatus["alert_engine"] = s.isHealthy(status)
	}

	// 检查通知管理器
	if s.notificationManager != nil {
		status := s.notificationManager.GetStatus()
		s.healthStatus["notification_manager"] = s.isHealthy(status)
	}

	// 检查是否需要重启不健康的组件
	if s.config.AutoRestart {
		s.restartUnhealthyComponents()
	}
}

// isHealthy 判断组件是否健康
func (s *MonitoringService) isHealthy(status map[string]interface{}) bool {
	// 简化的健康检查逻辑
	// 实际实现中可以根据具体的状态指标进行更复杂的判断
	return true
}

// restartUnhealthyComponents 重启不健康的组件
func (s *MonitoringService) restartUnhealthyComponents() {
	for component, healthy := range s.healthStatus {
		if !healthy {
			s.logger.Infof("组件 %s 不健康，尝试重启", component)

			switch component {
			case "metrics_collector":
				s.restartMetricsCollector()
			case "alert_engine":
				s.restartAlertEngine()
			case "notification_manager":
				s.restartNotificationManager()
			}
		}
	}
}

// restartMetricsCollector 重启指标收集器
func (s *MonitoringService) restartMetricsCollector() {
	if s.metricsCollector != nil {
		s.metricsCollector.Stop()
		time.Sleep(2 * time.Second)
	}

	if err := s.startMetricsCollector(); err != nil {
		s.logger.Errorf("重启指标收集器失败: %v", err)
	} else {
		s.logger.Info("指标收集器重启成功")
	}
}

// restartAlertEngine 重启告警引擎
func (s *MonitoringService) restartAlertEngine() {
	if s.alertEngine != nil {
		s.alertEngine.Stop()
		time.Sleep(2 * time.Second)
	}

	if err := s.startAlertEngine(); err != nil {
		s.logger.Errorf("重启告警引擎失败: %v", err)
	} else {
		s.logger.Info("告警引擎重启成功")
	}
}

// restartNotificationManager 重启通知管理器
func (s *MonitoringService) restartNotificationManager() {
	if s.notificationManager != nil {
		s.notificationManager.Stop()
		time.Sleep(2 * time.Second)
	}

	if err := s.startNotificationManager(); err != nil {
		s.logger.Errorf("重启通知管理器失败: %v", err)
	} else {
		s.logger.Info("通知管理器重启成功")
	}
}

// GetSystemStatus 获取系统状态
func (s *MonitoringService) GetSystemStatus() *SystemStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	components := make(map[string]*ComponentStatus)
	overallHealthy := true

	// 收集各组件状态
	for name, healthy := range s.healthStatus {
		status := "healthy"
		if !healthy {
			status = "unhealthy"
			overallHealthy = false
		}

		components[name] = &ComponentStatus{
			Name:       name,
			Status:     status,
			LastUpdate: time.Now(),
		}
	}

	// 获取活跃告警数量
	var activeAlerts int
	if s.alertEngine != nil {
		activeAlerts = len(s.alertEngine.GetActiveAlerts())
	}

	overallStatus := "healthy"
	if !overallHealthy {
		overallStatus = "degraded"
	}

	return &SystemStatus{
		OverallStatus:      overallStatus,
		Components:         components,
		StartTime:          time.Now(), // 简化实现
		Uptime:             time.Since(time.Now()),
		ActiveAlerts:       activeAlerts,
		TotalNotifications: 0, // 可以从通知管理器获取
		MetricsCollected:   0, // 可以从指标收集器获取
		LastHealthCheck:    time.Now(),
	}
}

// ReloadConfiguration 重新加载配置
func (s *MonitoringService) ReloadConfiguration() error {
	s.logger.Info("重新加载监控服务配置")

	var errors []string

	// 重新加载告警规则
	if s.alertEngine != nil {
		if err := s.alertEngine.ReloadRules(); err != nil {
			errors = append(errors, fmt.Sprintf("重新加载告警规则失败: %v", err))
		}
	}

	// 重新加载通知配置
	if s.notificationManager != nil {
		if err := s.notificationManager.ReloadConfig(); err != nil {
			errors = append(errors, fmt.Sprintf("重新加载通知配置失败: %v", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("配置重新加载部分失败: %v", errors)
	}

	s.logger.Info("监控服务配置重新加载完成")
	return nil
}

// GetMetricsCollector 获取指标收集器
func (s *MonitoringService) GetMetricsCollector() *monitoring.EnhancedMetricsCollector {
	return s.metricsCollector
}

// GetAlertEngine 获取告警引擎
func (s *MonitoringService) GetAlertEngine() *alerting.AlertEngine {
	return s.alertEngine
}

// GetNotificationManager 获取通知管理器
func (s *MonitoringService) GetNotificationManager() *notification.NotificationManager {
	return s.notificationManager
}
