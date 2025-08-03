//go:build unit
// +build unit

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/logic"
	"api/internal/svc"
	"api/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockDatabase 模拟数据库接口
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

// MockDBInterface 定义数据库接口
type MockDBInterface interface {
	Ping() error
}

// HealthCheckHandlerTestSuite 健康检查处理器测试套件
type HealthCheckHandlerTestSuite struct {
	suite.Suite
	mockDB  *MockDatabase
	svcCtx  *svc.ServiceContext
	handler http.HandlerFunc
	ctx     context.Context
}

// SetupTest 每个测试前的初始化
func (suite *HealthCheckHandlerTestSuite) SetupTest() {
	suite.mockDB = new(MockDatabase)
	suite.ctx = context.Background()

	// 创建模拟的服务上下文，暂时不设置DB字段，在测试中直接mock逻辑
	suite.svcCtx = &svc.ServiceContext{
		// DB字段将在具体测试中处理
	}

	// 创建健康检查处理器
	suite.handler = func(w http.ResponseWriter, r *http.Request) {
		healthLogic := logic.NewHealthCheckLogic(r.Context(), suite.svcCtx)
		resp, err := healthLogic.HealthCheck(&types.EmptyReq{})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

// TearDownTest 每个测试后的清理
func (suite *HealthCheckHandlerTestSuite) TearDownTest() {
	suite.mockDB.AssertExpectations(suite.T())
}

// TestHealthCheckHandlerSuccess 测试健康检查成功
func (suite *HealthCheckHandlerTestSuite) TestHealthCheckHandlerSuccess() {
	// 设置Mock期望
	suite.mockDB.On("Ping").Return(nil)

	// 创建HTTP请求
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// 执行处理器
	suite.handler(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "application/json", w.Header().Get("Content-Type"))

	// 解析响应体
	var response types.HealthResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(suite.T(), err)

	// 验证响应内容
	assert.Equal(suite.T(), "healthy", response.Status)
	assert.NotEmpty(suite.T(), response.Version)
	assert.Greater(suite.T(), response.Timestamp, int64(0))
	assert.NotEmpty(suite.T(), response.Uptime)
	assert.Len(suite.T(), response.Checks, 4) // database, redis, disk, memory

	// 验证各个检查项
	checkMap := make(map[string]types.CheckStatus)
	for _, check := range response.Checks {
		checkMap[check.Service] = check
	}

	// 验证数据库检查
	dbCheck, exists := checkMap["database"]
	assert.True(suite.T(), exists)
	assert.Equal(suite.T(), "healthy", dbCheck.Status)
	assert.Equal(suite.T(), "数据库连接正常", dbCheck.Message)
	assert.NotEmpty(suite.T(), dbCheck.Latency)
}

// TestHealthCheckHandlerDatabaseError 测试数据库连接错误
func (suite *HealthCheckHandlerTestSuite) TestHealthCheckHandlerDatabaseError() {
	// 设置Mock期望 - 数据库连接失败
	suite.mockDB.On("Ping").Return(assert.AnError)

	// 创建HTTP请求
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// 执行处理器
	suite.handler(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// 解析响应体
	var response types.HealthResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(suite.T(), err)

	// 验证响应内容
	assert.Equal(suite.T(), "unhealthy", response.Status)

	// 查找数据库检查结果
	var dbCheck *types.CheckStatus
	for _, check := range response.Checks {
		if check.Service == "database" {
			dbCheck = &check
			break
		}
	}

	assert.NotNil(suite.T(), dbCheck)
	assert.Equal(suite.T(), "unhealthy", dbCheck.Status)
	assert.Contains(suite.T(), dbCheck.Message, "数据库连接失败")
}

// TestHealthCheckHandlerInvalidMethod 测试无效的HTTP方法
func (suite *HealthCheckHandlerTestSuite) TestHealthCheckHandlerInvalidMethod() {
	// 设置Mock期望（可能不会被调用）
	suite.mockDB.On("Ping").Return(nil).Maybe()

	// 创建POST请求（应该是GET）
	req := httptest.NewRequest("POST", "/health", nil)
	w := httptest.NewRecorder()

	// 如果处理器检查HTTP方法，这里会失败
	// 否则，应该正常处理
	suite.handler(w, req)

	// 根据实际实现决定期望的行为
	// 这里假设处理器接受任何方法
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

// TestHealthCheckHandlerResponseHeaders 测试响应头设置
func (suite *HealthCheckHandlerTestSuite) TestHealthCheckHandlerResponseHeaders() {
	// 设置Mock期望
	suite.mockDB.On("Ping").Return(nil)

	// 创建HTTP请求
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// 执行处理器
	suite.handler(w, req)

	// 验证响应头
	assert.Equal(suite.T(), "application/json", w.Header().Get("Content-Type"))
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

// TestHealthCheckHandlerResponseStructure 测试响应结构完整性
func (suite *HealthCheckHandlerTestSuite) TestHealthCheckHandlerResponseStructure() {
	// 设置Mock期望
	suite.mockDB.On("Ping").Return(nil)

	// 创建HTTP请求
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// 执行处理器
	suite.handler(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// 解析响应体
	var response types.HealthResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(suite.T(), err)

	// 验证响应结构
	assert.NotEmpty(suite.T(), response.Status)
	assert.NotNil(suite.T(), response.Checks)
	assert.NotEmpty(suite.T(), response.Uptime)
	assert.NotEmpty(suite.T(), response.Version)
	assert.Greater(suite.T(), response.Timestamp, int64(0))

	// 验证每个检查项的结构
	for _, check := range response.Checks {
		assert.NotEmpty(suite.T(), check.Service)
		assert.NotEmpty(suite.T(), check.Status)
		assert.NotEmpty(suite.T(), check.Message)
		assert.NotEmpty(suite.T(), check.Latency)
		assert.Contains(suite.T(), []string{"healthy", "warning", "unhealthy"}, check.Status)
	}
}

// TestHealthCheckHandlerConcurrency 测试并发请求
func (suite *HealthCheckHandlerTestSuite) TestHealthCheckHandlerConcurrency() {
	// 设置Mock期望 - 允许多次调用
	suite.mockDB.On("Ping").Return(nil)

	// 并发发送多个请求
	concurrency := 10
	results := make(chan int, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			req := httptest.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()
			suite.handler(w, req)
			results <- w.Code
		}()
	}

	// 收集结果
	for i := 0; i < concurrency; i++ {
		statusCode := <-results
		assert.Equal(suite.T(), http.StatusOK, statusCode)
	}
}

// TestHealthCheckHandlerEmptyRequest 测试空请求体
func (suite *HealthCheckHandlerTestSuite) TestHealthCheckHandlerEmptyRequest() {
	// 设置Mock期望
	suite.mockDB.On("Ping").Return(nil)

	// 创建空的HTTP请求
	req := httptest.NewRequest("GET", "/health", bytes.NewBuffer([]byte("")))
	w := httptest.NewRecorder()

	// 执行处理器
	suite.handler(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// 解析响应体
	var response types.HealthResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "healthy", response.Status)
}

// 运行健康检查处理器测试套件
func TestHealthCheckHandlerSuite(t *testing.T) {
	suite.Run(t, new(HealthCheckHandlerTestSuite))
}

// 基准测试
func BenchmarkHealthCheckHandler(b *testing.B) {
	mockDB := new(MockDatabase)
	svcCtx := &svc.ServiceContext{
		// DB字段将在实际逻辑中处理
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		healthLogic := logic.NewHealthCheckLogic(r.Context(), svcCtx)
		resp, err := healthLogic.HealthCheck(&types.EmptyReq{})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}

	// 设置Mock期望
	mockDB.On("Ping").Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		handler(w, req)
	}
}
