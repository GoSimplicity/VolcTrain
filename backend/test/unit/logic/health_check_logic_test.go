//go:build unit
// +build unit

package logic

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"api/internal/logic"
	"api/internal/svc"
	"api/internal/types"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockDatabase 模拟数据库连接
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDatabase) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockRedisClient 模拟Redis客户端
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	args := m.Called(ctx)
	cmd := redis.NewStatusCmd(ctx)
	if err := args.Error(0); err != nil {
		cmd.SetErr(err)
	} else {
		cmd.SetVal("PONG")
	}
	return cmd
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

// HealthCheckLogicTestSuite 健康检查逻辑测试套件
type HealthCheckLogicTestSuite struct {
	suite.Suite
	mockDB      *MockDatabase
	mockRedis   *MockRedisClient
	svcCtx      *svc.ServiceContext
	healthLogic *logic.HealthCheckLogic
	ctx         context.Context
}

// SetupTest 每个测试前的初始化
func (suite *HealthCheckLogicTestSuite) SetupTest() {
	suite.mockDB = new(MockDatabase)
	suite.mockRedis = new(MockRedisClient)
	suite.ctx = context.Background()

	// 创建模拟的服务上下文
	suite.svcCtx = &svc.ServiceContext{
		DatabaseConn: suite.mockDB,
		RedisClient:  suite.mockRedis,
	}

	suite.healthLogic = logic.NewHealthCheckLogic(suite.ctx, suite.svcCtx)
}

// TearDownTest 每个测试后的清理
func (suite *HealthCheckLogicTestSuite) TearDownTest() {
	suite.mockDB.AssertExpectations(suite.T())
	suite.mockRedis.AssertExpectations(suite.T())
}

// TestHealthCheckAllHealthy 测试所有组件健康的情况
func (suite *HealthCheckLogicTestSuite) TestHealthCheckAllHealthy() {
	// 设置Mock期望
	suite.mockDB.On("Ping").Return(nil)
	suite.mockRedis.On("Ping", suite.ctx).Return(nil)

	// 执行测试
	req := &types.EmptyReq{}
	resp, err := suite.healthLogic.HealthCheck(req)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "healthy", resp.Status)
	assert.NotEmpty(suite.T(), resp.Version)
	assert.Greater(suite.T(), resp.Timestamp, int64(0))
	assert.NotEmpty(suite.T(), resp.Uptime)

	// 验证检查项目
	assert.Len(suite.T(), resp.Checks, 4) // database, redis, disk, memory

	// 验证数据库检查
	dbCheck := findCheck(resp.Checks, "database")
	assert.NotNil(suite.T(), dbCheck)
	assert.Equal(suite.T(), "healthy", dbCheck.Status)
	assert.Equal(suite.T(), "数据库连接正常", dbCheck.Message)
	assert.NotEmpty(suite.T(), dbCheck.Latency)

	// 验证Redis检查
	redisCheck := findCheck(resp.Checks, "redis")
	assert.NotNil(suite.T(), redisCheck)
	assert.Equal(suite.T(), "healthy", redisCheck.Status)
	assert.Equal(suite.T(), "Redis连接正常", redisCheck.Message)

	// 验证磁盘检查
	diskCheck := findCheck(resp.Checks, "disk")
	assert.NotNil(suite.T(), diskCheck)
	assert.Contains(suite.T(), []string{"healthy", "warning"}, diskCheck.Status)
	assert.Contains(suite.T(), diskCheck.Message, "磁盘使用率")

	// 验证内存检查
	memoryCheck := findCheck(resp.Checks, "memory")
	assert.NotNil(suite.T(), memoryCheck)
	assert.Contains(suite.T(), []string{"healthy", "warning"}, memoryCheck.Status)
	assert.Contains(suite.T(), memoryCheck.Message, "内存使用")
}

// TestHealthCheckDatabaseUnhealthy 测试数据库不健康的情况
func (suite *HealthCheckLogicTestSuite) TestHealthCheckDatabaseUnhealthy() {
	// 设置Mock期望 - 数据库连接失败
	suite.mockDB.On("Ping").Return(errors.New("connection timeout"))
	suite.mockRedis.On("Ping", suite.ctx).Return(nil)

	// 执行测试
	req := &types.EmptyReq{}
	resp, err := suite.healthLogic.HealthCheck(req)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "unhealthy", resp.Status) // 整体状态应该是不健康

	// 验证数据库检查
	dbCheck := findCheck(resp.Checks, "database")
	assert.NotNil(suite.T(), dbCheck)
	assert.Equal(suite.T(), "unhealthy", dbCheck.Status)
	assert.Contains(suite.T(), dbCheck.Message, "数据库连接失败")
	assert.Contains(suite.T(), dbCheck.Message, "connection timeout")
}

// TestHealthCheckRedisUnhealthy 测试Redis不健康的情况
func (suite *HealthCheckLogicTestSuite) TestHealthCheckRedisUnhealthy() {
	// 设置Mock期望 - Redis连接失败
	suite.mockDB.On("Ping").Return(nil)
	suite.mockRedis.On("Ping", suite.ctx).Return(errors.New("redis connection failed"))

	// 执行测试
	req := &types.EmptyReq{}
	resp, err := suite.healthLogic.HealthCheck(req)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "unhealthy", resp.Status)

	// 验证Redis检查
	redisCheck := findCheck(resp.Checks, "redis")
	assert.NotNil(suite.T(), redisCheck)
	assert.Equal(suite.T(), "unhealthy", redisCheck.Status)
	assert.Contains(suite.T(), redisCheck.Message, "Redis连接失败")
	assert.Contains(suite.T(), redisCheck.Message, "redis connection failed")
}

// TestHealthCheckNilConnections 测试连接为nil的情况
func (suite *HealthCheckLogicTestSuite) TestHealthCheckNilConnections() {
	// 创建没有连接的服务上下文
	suite.svcCtx = &svc.ServiceContext{
		DatabaseConn: nil,
		RedisClient:  nil,
	}
	suite.healthLogic = logic.NewHealthCheckLogic(suite.ctx, suite.svcCtx)

	// 执行测试
	req := &types.EmptyReq{}
	resp, err := suite.healthLogic.HealthCheck(req)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "healthy", resp.Status) // 没有连接也算健康

	// 验证数据库检查
	dbCheck := findCheck(resp.Checks, "database")
	assert.NotNil(suite.T(), dbCheck)
	assert.Equal(suite.T(), "healthy", dbCheck.Status)

	// 验证Redis检查
	redisCheck := findCheck(resp.Checks, "redis")
	assert.NotNil(suite.T(), redisCheck)
	assert.Equal(suite.T(), "healthy", redisCheck.Status)
}

// TestHealthCheckMultipleUnhealthy 测试多个组件不健康的情况
func (suite *HealthCheckLogicTestSuite) TestHealthCheckMultipleUnhealthy() {
	// 设置Mock期望 - 都连接失败
	suite.mockDB.On("Ping").Return(errors.New("db error"))
	suite.mockRedis.On("Ping", suite.ctx).Return(errors.New("redis error"))

	// 执行测试
	req := &types.EmptyReq{}
	resp, err := suite.healthLogic.HealthCheck(req)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "unhealthy", resp.Status)

	// 验证两个检查都失败
	dbCheck := findCheck(resp.Checks, "database")
	assert.Equal(suite.T(), "unhealthy", dbCheck.Status)

	redisCheck := findCheck(resp.Checks, "redis")
	assert.Equal(suite.T(), "unhealthy", redisCheck.Status)
}

// TestHealthCheckResponseStructure 测试响应结构的完整性
func (suite *HealthCheckLogicTestSuite) TestHealthCheckResponseStructure() {
	// 设置Mock期望
	suite.mockDB.On("Ping").Return(nil)
	suite.mockRedis.On("Ping", suite.ctx).Return(nil)

	// 执行测试
	req := &types.EmptyReq{}
	resp, err := suite.healthLogic.HealthCheck(req)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)

	// 验证响应结构
	assert.NotEmpty(suite.T(), resp.Status)
	assert.NotNil(suite.T(), resp.Checks)
	assert.NotEmpty(suite.T(), resp.Uptime)
	assert.NotEmpty(suite.T(), resp.Version)
	assert.Greater(suite.T(), resp.Timestamp, int64(0))

	// 验证时间戳的合理性（应该接近当前时间）
	now := time.Now().Unix()
	assert.InDelta(suite.T(), now, resp.Timestamp, 5) // 允许5秒误差

	// 验证版本号格式
	assert.Equal(suite.T(), "1.0.0", resp.Version)

	// 验证运行时间格式
	assert.Contains(suite.T(), resp.Uptime, "s") // 应该包含时间单位
}

// TestHealthCheckLatencyMeasurement 测试延迟测量
func (suite *HealthCheckLogicTestSuite) TestHealthCheckLatencyMeasurement() {
	// 设置Mock期望 - 添加一些延迟
	suite.mockDB.On("Ping").Return(nil).After(10 * time.Millisecond)
	suite.mockRedis.On("Ping", suite.ctx).Return(nil).After(5 * time.Millisecond)

	// 执行测试
	req := &types.EmptyReq{}
	resp, err := suite.healthLogic.HealthCheck(req)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)

	// 验证延迟测量
	for _, check := range resp.Checks {
		assert.NotEmpty(suite.T(), check.Latency)
		assert.Contains(suite.T(), check.Latency, "s") // 延迟应该包含时间单位
	}
}

// findCheck 查找指定服务的检查结果
func findCheck(checks []types.CheckStatus, service string) *types.CheckStatus {
	for _, check := range checks {
		if check.Service == service {
			return &check
		}
	}
	return nil
}

// 运行健康检查逻辑测试套件
func TestHealthCheckLogicSuite(t *testing.T) {
	suite.Run(t, new(HealthCheckLogicTestSuite))
}

// 基准测试
func BenchmarkHealthCheck(b *testing.B) {
	mockDB := new(MockDatabase)
	mockRedis := new(MockRedisClient)

	svcCtx := &svc.ServiceContext{
		DatabaseConn: mockDB,
		RedisClient:  mockRedis,
	}

	healthLogic := logic.NewHealthCheckLogic(context.Background(), svcCtx)
	req := &types.EmptyReq{}

	// 设置Mock期望
	mockDB.On("Ping").Return(nil)
	mockRedis.On("Ping", mock.Anything).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := healthLogic.HealthCheck(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}
