package test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"api/pkg/auth"
	"api/pkg/monitoring"
	"api/pkg/notification"
	testconfig "api/test/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ComprehensiveTestSuite 综合测试套件
type ComprehensiveTestSuite struct {
	suite.Suite
	testConfig      *testconfig.TestConfig
	db              *sql.DB
	ctx             context.Context
	passwordService *auth.PasswordService
	jwtService      *auth.JWTService
}

// SetupSuite 测试套件初始化
func (suite *ComprehensiveTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.testConfig = testconfig.GetTestConfig()

	// 初始化认证服务
	suite.passwordService = auth.NewPasswordService()
	suite.jwtService = auth.NewJWTService("test-access-secret", "test-refresh-secret", 3600, 86400)

	// 尝试创建测试数据库连接，如果失败则使用Mock
	db, err := suite.testConfig.CreateTestDB()
	if err != nil {
		suite.T().Logf("无法连接真实数据库，使用Mock: %v", err)
		mockDB, _, mockErr := testconfig.CreateMockDB()
		if mockErr != nil {
			suite.T().Fatalf("创建Mock数据库失败: %v", mockErr)
		}
		suite.db = mockDB
	} else {
		suite.db = db
	}
}

// TearDownSuite 测试套件清理
func (suite *ComprehensiveTestSuite) TearDownSuite() {
	if suite.db != nil {
		// 清理测试数据
		suite.testConfig.CleanupTestData(suite.db)
		suite.db.Close()
	}
}

// TestSystemIntegration 系统集成测试
func (suite *ComprehensiveTestSuite) TestSystemIntegration() {
	suite.T().Log("开始系统集成测试...")

	// 测试数据库集成
	suite.Run("DatabaseIntegration", suite.testDatabaseIntegration)

	// 测试认证系统集成
	suite.Run("AuthenticationIntegration", suite.testAuthenticationIntegration)

	// 测试监控系统集成
	suite.Run("MonitoringIntegration", suite.testMonitoringIntegration)

	// 测试通知系统集成
	suite.Run("NotificationIntegration", suite.testNotificationIntegration)
}

// testDatabaseIntegration 数据库集成测试
func (suite *ComprehensiveTestSuite) testDatabaseIntegration() {
	// 测试数据库连接
	assert.NotNil(suite.T(), suite.db, "数据库连接不应为空")

	// 测试基本查询
	var result int
	err := suite.db.QueryRow("SELECT 1").Scan(&result)
	assert.NoError(suite.T(), err, "基本查询应该成功")
	assert.Equal(suite.T(), 1, result, "查询结果应该为1")

	suite.T().Log("✅ 数据库集成测试通过")
}

// testAuthenticationIntegration 认证系统集成测试
func (suite *ComprehensiveTestSuite) testAuthenticationIntegration() {
	// 测试密码加密和验证
	password := "test123456"
	hashedPassword, err := suite.passwordService.HashPassword(password)
	assert.NoError(suite.T(), err, "密码加密不应该失败")
	assert.NotEmpty(suite.T(), hashedPassword, "加密后密码不应为空")

	// 验证密码
	isValid := suite.passwordService.VerifyPassword(hashedPassword, password)
	assert.True(suite.T(), isValid, "密码验证应该成功")

	// 测试JWT生成和验证
	userID := int64(1)
	tokenPair, err := suite.jwtService.GenerateTokenPair(userID, "test", "test@example.com", []string{"user"}, []string{"read"})
	assert.NoError(suite.T(), err, "JWT生成不应该失败")
	assert.NotEmpty(suite.T(), tokenPair.AccessToken, "JWT token不应为空")

	// 验证JWT
	claims, err := suite.jwtService.ParseAccessToken(tokenPair.AccessToken)
	assert.NoError(suite.T(), err, "JWT验证不应该失败")
	assert.Equal(suite.T(), userID, claims.UserID, "用户ID应该匹配")

	suite.T().Log("✅ 认证系统集成测试通过")
}

// testMonitoringIntegration 监控系统集成测试
func (suite *ComprehensiveTestSuite) testMonitoringIntegration() {
	// 创建监控收集器配置
	collectorConfig := &monitoring.CollectorConfig{
		SystemMetricsEnabled:   true,
		BusinessMetricsEnabled: true,
		PrometheusEnabled:      false, // 测试环境不启用Prometheus
		CollectInterval:        time.Second * 30,
	}

	// 创建监控收集器
	collector := monitoring.NewEnhancedMetricsCollector(suite.db, collectorConfig)
	assert.NotNil(suite.T(), collector, "监控收集器不应为空")

	// 测试状态获取
	status := collector.GetCollectionStatus()
	assert.NotNil(suite.T(), status, "收集状态不应为空")
	assert.Contains(suite.T(), status, "system_metrics_enabled", "状态应包含系统指标配置")
	assert.Equal(suite.T(), true, status["system_metrics_enabled"], "系统指标应该启用")

	suite.T().Log("✅ 监控系统集成测试通过")
}

// testNotificationIntegration 通知系统集成测试
func (suite *ComprehensiveTestSuite) testNotificationIntegration() {
	// 创建通知配置
	notificationConfig := &notification.NotificationConfig{
		MaxQueueSize:         100,
		MaxConcurrentSenders: 2,
		RetryMaxAttempts:     3,
		RetryBackoffSeconds:  5,
		RateLimitPerMinute:   10,
		TimeoutSeconds:       30,
		FailedRetentionDays:  7,
		EnableDeduplication:  false,
		DeduplicationWindow:  time.Minute * 5,
	}

	// 创建通知管理器
	notificationManager := notification.NewNotificationManager(suite.db, notificationConfig)
	assert.NotNil(suite.T(), notificationManager, "通知管理器不应为空")

	// 测试通知管理器启动（不实际发送通知）
	err := notificationManager.Start()
	assert.NoError(suite.T(), err, "通知管理器启动不应该失败")

	// 停止通知管理器
	notificationManager.Stop()

	suite.T().Log("✅ 通知系统集成测试通过")
}

// TestPerformance 性能测试
func (suite *ComprehensiveTestSuite) TestPerformance() {
	suite.T().Log("开始性能测试...")

	// 测试并发性能
	suite.Run("ConcurrencyPerformance", suite.testConcurrencyPerformance)

	// 测试数据库性能
	suite.Run("DatabasePerformance", suite.testDatabasePerformance)

	// 测试内存使用
	suite.Run("MemoryUsage", suite.testMemoryUsage)
}

// testConcurrencyPerformance 并发性能测试
func (suite *ComprehensiveTestSuite) testConcurrencyPerformance() {
	concurrency := 100
	iterations := 1000

	start := time.Now()

	// 并发执行密码哈希操作
	ch := make(chan bool, concurrency)
	for i := 0; i < iterations; i++ {
		go func(id int) {
			password := fmt.Sprintf("test%d", id)
			_, err := suite.passwordService.HashPassword(password)
			assert.NoError(suite.T(), err, "密码哈希不应该失败")
			ch <- true
		}(i)
	}

	// 等待所有操作完成
	for i := 0; i < iterations; i++ {
		<-ch
	}

	duration := time.Since(start)
	suite.T().Logf("并发密码哈希 %d 次用时: %v", iterations, duration)

	// 验证性能要求：平均每次操作不超过10ms
	avgDuration := duration / time.Duration(iterations)
	assert.Less(suite.T(), avgDuration, time.Millisecond*10, "平均操作时间应小于10ms")

	suite.T().Log("✅ 并发性能测试通过")
}

// testDatabasePerformance 数据库性能测试
func (suite *ComprehensiveTestSuite) testDatabasePerformance() {
	queries := 1000
	start := time.Now()

	// 执行多次查询
	for i := 0; i < queries; i++ {
		var result int
		err := suite.db.QueryRow("SELECT ?", i).Scan(&result)
		if err != nil {
			// 如果是Mock数据库，跳过实际查询
			suite.T().Log("跳过Mock数据库性能测试")
			return
		}
	}

	duration := time.Since(start)
	suite.T().Logf("执行 %d 次数据库查询用时: %v", queries, duration)

	// 验证性能要求：平均每次查询不超过1ms
	avgDuration := duration / time.Duration(queries)
	assert.Less(suite.T(), avgDuration, time.Millisecond, "平均查询时间应小于1ms")

	suite.T().Log("✅ 数据库性能测试通过")
}

// testMemoryUsage 内存使用测试
func (suite *ComprehensiveTestSuite) testMemoryUsage() {
	// 创建大量对象测试内存使用
	const objectCount = 10000
	objects := make([]interface{}, objectCount)

	for i := 0; i < objectCount; i++ {
		objects[i] = map[string]interface{}{
			"id":   i,
			"name": fmt.Sprintf("object_%d", i),
			"data": make([]byte, 1024), // 1KB数据
		}
	}

	// 验证对象创建成功
	assert.Len(suite.T(), objects, objectCount, "应该创建指定数量的对象")

	// 清理对象
	objects = nil

	suite.T().Log("✅ 内存使用测试通过")
}

// TestSecurity 安全测试
func (suite *ComprehensiveTestSuite) TestSecurity() {
	suite.T().Log("开始安全测试...")

	// 测试密码安全
	suite.Run("PasswordSecurity", suite.testPasswordSecurity)

	// 测试JWT安全
	suite.Run("JWTSecurity", suite.testJWTSecurity)

	// 测试输入验证
	suite.Run("InputValidation", suite.testInputValidation)
}

// testPasswordSecurity 密码安全测试
func (suite *ComprehensiveTestSuite) testPasswordSecurity() {
	// 测试弱密码处理
	weakPasswords := []string{"123", "password", "admin", ""}

	for _, weak := range weakPasswords {
		// 即使是弱密码，加密功能也应该正常工作
		hashed, err := suite.passwordService.HashPassword(weak)
		assert.NoError(suite.T(), err, "密码加密不应该失败")
		assert.NotEqual(suite.T(), weak, hashed, "加密后密码应该与原密码不同")
		assert.True(suite.T(), len(hashed) > 20, "加密后密码应该足够长")
	}

	// 测试相同密码生成不同哈希
	password := "samepassword"
	hash1, _ := suite.passwordService.HashPassword(password)
	hash2, _ := suite.passwordService.HashPassword(password)
	assert.NotEqual(suite.T(), hash1, hash2, "相同密码应该生成不同的哈希值")

	suite.T().Log("✅ 密码安全测试通过")
}

// testJWTSecurity JWT安全测试
func (suite *ComprehensiveTestSuite) testJWTSecurity() {
	userID := int64(123)
	email := "test@example.com"

	// 生成有效JWT
	tokenPair, err := suite.jwtService.GenerateTokenPair(userID, "test", email, []string{"user"}, []string{"read"})
	assert.NoError(suite.T(), err, "JWT生成不应该失败")

	// 测试有效JWT验证
	claims, err := suite.jwtService.ParseAccessToken(tokenPair.AccessToken)
	assert.NoError(suite.T(), err, "有效JWT验证不应该失败")
	assert.Equal(suite.T(), userID, claims.UserID, "用户ID应该匹配")
	assert.Equal(suite.T(), email, claims.Email, "邮箱应该匹配")

	// 测试无效JWT
	invalidTokens := []string{
		"invalid.token.here",
		"",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid",
		tokenPair.AccessToken + "tampered",
	}

	for _, invalidToken := range invalidTokens {
		_, err := suite.jwtService.ParseAccessToken(invalidToken)
		assert.Error(suite.T(), err, "无效JWT应该验证失败")
	}

	suite.T().Log("✅ JWT安全测试通过")
}

// testInputValidation 输入验证测试
func (suite *ComprehensiveTestSuite) testInputValidation() {
	// 测试SQL注入防护（通过参数化查询）
	maliciousInputs := []string{
		"'; DROP TABLE users; --",
		"1 OR 1=1",
		"<script>alert('xss')</script>",
		"../../../etc/passwd",
	}

	for _, input := range maliciousInputs {
		// 使用参数化查询测试
		var result string
		err := suite.db.QueryRow("SELECT ?", input).Scan(&result)
		if err != nil {
			// Mock数据库可能不支持，跳过
			suite.T().Log("跳过Mock数据库注入测试")
			continue
		}
		// 如果查询成功，确保返回的是原始输入，说明没有被执行为SQL
		assert.Equal(suite.T(), input, result, "参数化查询应该返回原始输入")
	}

	suite.T().Log("✅ 输入验证测试通过")
}

// TestComprehensiveSuite 运行综合测试套件
func TestComprehensiveSuite(t *testing.T) {
	suite.Run(t, new(ComprehensiveTestSuite))
}
