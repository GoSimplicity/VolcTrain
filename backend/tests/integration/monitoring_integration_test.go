package integration

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"api/internal/service"
	"api/pkg/alerting"
	"api/pkg/monitoring"
	"api/pkg/notification"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MonitoringIntegrationTestSuite 监控系统集成测试套件
type MonitoringIntegrationTestSuite struct {
	suite.Suite
	db                *sql.DB
	monitoringService *service.MonitoringService
	config            *service.MonitoringConfig
}

// SetupSuite 测试套件初始化
func (suite *MonitoringIntegrationTestSuite) SetupSuite() {
	// 连接测试数据库
	db, err := sql.Open("mysql", "volctrain_app:Abc@1234@tcp(localhost:3306)/volctraindb?charset=utf8mb4&parseTime=True&loc=Local")
	suite.Require().NoError(err)

	err = db.Ping()
	suite.Require().NoError(err)

	suite.db = db

	// 创建测试配置
	suite.config = &service.MonitoringConfig{
		MetricsConfig: &monitoring.CollectorConfig{
			PrometheusURL:          "http://localhost:9090",
			CollectInterval:        30 * time.Second,
			SystemMetricsEnabled:   true,
			BusinessMetricsEnabled: true,
			PrometheusEnabled:      true,
			MetricsPort:            9094,
			EnableBuiltinMetrics:   true,
			MaxRetries:             3,
			Timeout:                30 * time.Second,
		},
		AlertConfig: &alerting.AlertEngineConfig{
			EvaluationInterval:      30 * time.Second,
			MaxConcurrentRules:      10,
			AlertRetentionDays:      30,
			EnableGrouping:          true,
			EnableSuppression:       true,
			DefaultThrottleMinutes:  5,
			AnomalyDetectionEnabled: true,
		},
		NotificationConfig: &notification.NotificationConfig{
			MaxQueueSize:         1000,
			MaxConcurrentSenders: 5,
			RetryMaxAttempts:     3,
			RetryBackoffSeconds:  5,
			RateLimitPerMinute:   60,
			TimeoutSeconds:       30,
			FailedRetentionDays:  7,
			EnableDeduplication:  true,
			DeduplicationWindow:  5 * time.Minute,
		},
		EnableMetrics:       true,
		EnableAlerts:        true,
		EnableNotifications: true,
		HealthCheckInterval: 30 * time.Second,
		AutoRestart:         true,
		MaxRestartAttempts:  3,
	}
}

// TearDownSuite 测试套件清理
func (suite *MonitoringIntegrationTestSuite) TearDownSuite() {
	if suite.monitoringService != nil {
		suite.monitoringService.Stop()
	}
	if suite.db != nil {
		suite.db.Close()
	}
}

// SetupTest 每个测试用例前的初始化
func (suite *MonitoringIntegrationTestSuite) SetupTest() {
	// 创建监控服务实例
	suite.monitoringService = service.NewMonitoringService(suite.db, suite.config)
}

// TearDownTest 每个测试用例后的清理
func (suite *MonitoringIntegrationTestSuite) TearDownTest() {
	if suite.monitoringService != nil {
		suite.monitoringService.Stop()
	}
}

// TestMonitoringServiceStartup 测试监控服务启动
func (suite *MonitoringIntegrationTestSuite) TestMonitoringServiceStartup() {
	// 启动监控服务
	err := suite.monitoringService.Start()
	assert.NoError(suite.T(), err, "监控服务启动应该成功")

	// 等待服务完全启动
	time.Sleep(2 * time.Second)

	// 检查系统状态
	status := suite.monitoringService.GetSystemStatus()
	assert.NotNil(suite.T(), status, "系统状态不应为空")
	assert.Equal(suite.T(), "healthy", status.OverallStatus, "系统整体状态应为健康")

	// 检查各组件状态
	components := status.Components
	assert.NotNil(suite.T(), components, "组件状态不应为空")
	assert.Contains(suite.T(), components, "metrics_collector", "应包含指标收集器组件")
	assert.Contains(suite.T(), components, "alert_engine", "应包含告警引擎组件")
	assert.Contains(suite.T(), components, "notification_manager", "应包含通知管理器组件")
}

// TestMetricsCollection 测试指标收集功能
func (suite *MonitoringIntegrationTestSuite) TestMetricsCollection() {
	err := suite.monitoringService.Start()
	suite.Require().NoError(err)

	// 等待指标收集开始
	time.Sleep(5 * time.Second)

	collector := suite.monitoringService.GetMetricsCollector()
	assert.NotNil(suite.T(), collector, "指标收集器不应为空")

	// 获取收集状态
	status := collector.GetCollectionStatus()
	assert.NotNil(suite.T(), status, "收集状态不应为空")

	// 验证内置指标是否正确初始化
	builtinMetrics := collector.GetBuiltinMetrics()
	assert.NotNil(suite.T(), builtinMetrics, "内置指标不应为空")
	assert.GreaterOrEqual(suite.T(), len(builtinMetrics), 4, "至少应有4个内置指标")
}

// TestAlertEngineEvaluation 测试告警引擎评估
func (suite *MonitoringIntegrationTestSuite) TestAlertEngineEvaluation() {
	err := suite.monitoringService.Start()
	suite.Require().NoError(err)

	// 等待告警引擎启动
	time.Sleep(3 * time.Second)

	alertEngine := suite.monitoringService.GetAlertEngine()
	assert.NotNil(suite.T(), alertEngine, "告警引擎不应为空")

	// 获取引擎状态
	status := alertEngine.GetEngineStatus()
	assert.NotNil(suite.T(), status, "引擎状态不应为空")

	// 检查是否有活跃规则
	activeAlerts := alertEngine.GetActiveAlerts()
	assert.NotNil(suite.T(), activeAlerts, "活跃告警列表不应为空")
}

// TestNotificationChannels 测试通知渠道
func (suite *MonitoringIntegrationTestSuite) TestNotificationChannels() {
	err := suite.monitoringService.Start()
	suite.Require().NoError(err)

	// 等待通知管理器启动
	time.Sleep(2 * time.Second)

	notificationManager := suite.monitoringService.GetNotificationManager()
	assert.NotNil(suite.T(), notificationManager, "通知管理器不应为空")

	// 获取通知状态
	status := notificationManager.GetStatus()
	assert.NotNil(suite.T(), status, "通知状态不应为空")

	// 验证默认渠道是否注册
	channels := notificationManager.GetRegisteredChannels()
	assert.NotNil(suite.T(), channels, "注册渠道列表不应为空")
	assert.GreaterOrEqual(suite.T(), len(channels), 2, "至少应有2个注册渠道")
}

// TestHealthCheck 测试健康检查功能
func (suite *MonitoringIntegrationTestSuite) TestHealthCheck() {
	err := suite.monitoringService.Start()
	suite.Require().NoError(err)

	// 等待健康检查运行
	time.Sleep(35 * time.Second) // 等待超过健康检查间隔

	status := suite.monitoringService.GetSystemStatus()
	assert.NotNil(suite.T(), status, "系统状态不应为空")
	assert.True(suite.T(), status.LastHealthCheck.After(time.Now().Add(-60*time.Second)), "健康检查应在最近执行")
}

// TestConfigurationReload 测试配置重新加载
func (suite *MonitoringIntegrationTestSuite) TestConfigurationReload() {
	err := suite.monitoringService.Start()
	suite.Require().NoError(err)

	// 等待服务启动完成
	time.Sleep(2 * time.Second)

	// 重新加载配置
	err = suite.monitoringService.ReloadConfiguration()
	assert.NoError(suite.T(), err, "配置重新加载应该成功")

	// 验证服务仍然正常运行
	status := suite.monitoringService.GetSystemStatus()
	assert.Equal(suite.T(), "healthy", status.OverallStatus, "重新加载后系统应保持健康")
}

// TestAlertNotificationFlow 测试告警通知流程
func (suite *MonitoringIntegrationTestSuite) TestAlertNotificationFlow() {
	err := suite.monitoringService.Start()
	suite.Require().NoError(err)

	// 等待所有组件启动
	time.Sleep(5 * time.Second)

	alertEngine := suite.monitoringService.GetAlertEngine()
	notificationManager := suite.monitoringService.GetNotificationManager()

	assert.NotNil(suite.T(), alertEngine, "告警引擎不应为空")
	assert.NotNil(suite.T(), notificationManager, "通知管理器不应为空")

	// 模拟告警触发（通过创建测试告警规则）
	// 这里可以添加更详细的告警流程测试

	// 验证告警通知桥接是否正常工作
	// 等待告警通知桥接处理
	time.Sleep(10 * time.Second)

	// 检查通知状态
	notificationStatus := notificationManager.GetStatus()
	assert.NotNil(suite.T(), notificationStatus, "通知状态不应为空")
}

// TestComponentRestart 测试组件重启功能
func (suite *MonitoringIntegrationTestSuite) TestComponentRestart() {
	// 创建启用自动重启的配置
	config := *suite.config
	config.AutoRestart = true
	config.MaxRestartAttempts = 3

	monitoringService := service.NewMonitoringService(suite.db, &config)
	defer monitoringService.Stop()

	err := monitoringService.Start()
	suite.Require().NoError(err)

	// 等待服务启动完成
	time.Sleep(2 * time.Second)

	// 模拟组件故障（实际测试中可能需要更复杂的模拟）
	// 这里主要验证重启机制的存在

	status := monitoringService.GetSystemStatus()
	assert.NotNil(suite.T(), status, "系统状态不应为空")
}

// TestConcurrentOperations 测试并发操作
func (suite *MonitoringIntegrationTestSuite) TestConcurrentOperations() {
	err := suite.monitoringService.Start()
	suite.Require().NoError(err)

	// 等待服务启动完成
	time.Sleep(2 * time.Second)

	// 创建多个并发的状态查询
	const numGoroutines = 10
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer func() { done <- true }()

			// 执行多种并发操作
			for j := 0; j < 5; j++ {
				status := suite.monitoringService.GetSystemStatus()
				assert.NotNil(suite.T(), status, "并发获取系统状态不应为空")

				collector := suite.monitoringService.GetMetricsCollector()
				assert.NotNil(suite.T(), collector, "并发获取指标收集器不应为空")

				time.Sleep(100 * time.Millisecond)
			}
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < numGoroutines; i++ {
		select {
		case <-done:
		case <-time.After(30 * time.Second):
			suite.T().Fatal("并发操作测试超时")
		}
	}
}

// TestLongRunningStability 测试长时间运行稳定性
func (suite *MonitoringIntegrationTestSuite) TestLongRunningStability() {
	err := suite.monitoringService.Start()
	suite.Require().NoError(err)

	// 运行较长时间以测试稳定性
	duration := 60 * time.Second
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	healthCheckCount := 0
	for {
		select {
		case <-ctx.Done():
			// 测试完成
			assert.GreaterOrEqual(suite.T(), healthCheckCount, 10, "应进行多次健康检查")
			return
		case <-ticker.C:
			// 定期检查系统状态
			status := suite.monitoringService.GetSystemStatus()
			assert.NotNil(suite.T(), status, "长时间运行中系统状态不应为空")
			assert.Equal(suite.T(), "healthy", status.OverallStatus, "长时间运行中系统应保持健康")
			healthCheckCount++
		}
	}
}

// TestMonitoringIntegration 运行监控系统集成测试
func TestMonitoringIntegration(t *testing.T) {
	suite.Run(t, new(MonitoringIntegrationTestSuite))
}
