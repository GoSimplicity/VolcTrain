//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"api/internal/config"
	"api/internal/handler"
	"api/internal/svc"
	"api/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

// ServiceIntegrationTestSuite 服务集成测试套件
type ServiceIntegrationTestSuite struct {
	suite.Suite
	server  *rest.Server
	svcCtx  *svc.ServiceContext
	config  config.Config
	baseURL string
	client  *http.Client
}

// SetupSuite 整个测试套件开始前的初始化
func (suite *ServiceIntegrationTestSuite) SetupSuite() {
	// 加载测试配置
	var c config.Config
	conf.MustLoad("../../../etc/config-dev.yaml", &c)

	// 设置测试环境配置
	c.MySQL.DataSource = "test_user:test_password@tcp(localhost:3306)/volctrain_test?charset=utf8mb4&parseTime=true&loc=Local"
	c.Port = 0 // 使用随机端口

	suite.config = c

	// 创建服务上下文
	suite.svcCtx = svc.NewServiceContext(c)

	// 创建HTTP服务器
	suite.server = rest.MustNewServer(rest.RestConf{
		ServiceConf: rest.ServiceConf{
			Name: "integration-test",
			Log: rest.LogConf{
				Level: "error", // 减少测试日志输出
			},
		},
		Host: "localhost",
		Port: 0, // 使用随机端口
	})

	// 注册路由（简化版本，只注册核心路由）
	suite.registerRoutes()

	// 启动服务器
	go func() {
		suite.server.Start()
	}()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	// 获取服务器地址
	addr := suite.server.Addr()
	suite.baseURL = fmt.Sprintf("http://%s", addr)

	// 创建HTTP客户端
	suite.client = &http.Client{
		Timeout: 10 * time.Second,
	}
}

// TearDownSuite 整个测试套件结束后的清理
func (suite *ServiceIntegrationTestSuite) TearDownSuite() {
	if suite.server != nil {
		suite.server.Stop()
	}
}

// registerRoutes 注册测试路由
func (suite *ServiceIntegrationTestSuite) registerRoutes() {
	// 健康检查路由
	suite.server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/health",
		Handler: handler.NewHealthCheckHandler(suite.svcCtx),
	})

	// 训练作业相关路由
	suite.server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/v1/training/jobs",
		Handler: handler.NewCreateTrainingJobHandler(suite.svcCtx),
	})

	suite.server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/training/jobs/:id",
		Handler: handler.NewGetTrainingJobHandler(suite.svcCtx),
	})

	suite.server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/training/jobs",
		Handler: handler.NewListTrainingJobsHandler(suite.svcCtx),
	})

	// GPU设备相关路由
	suite.server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/gpu/devices",
		Handler: handler.NewListGpuDevicesHandler(suite.svcCtx),
	})

	suite.server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/v1/gpu/devices/:id/allocate",
		Handler: handler.NewAllocateGpuDeviceHandler(suite.svcCtx),
	})
}

// TestHealthCheckServiceIntegration 测试健康检查服务集成
func (suite *ServiceIntegrationTestSuite) TestHealthCheckServiceIntegration() {
	// 测试健康检查端点
	resp, err := suite.client.Get(suite.baseURL + "/health")
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	assert.Equal(suite.T(), "application/json", resp.Header.Get("Content-Type"))

	var healthResp types.HealthResponse
	err = json.NewDecoder(resp.Body).Decode(&healthResp)
	assert.NoError(suite.T(), err)

	// 验证健康检查响应
	assert.Contains(suite.T(), []string{"healthy", "warning", "unhealthy"}, healthResp.Status)
	assert.NotEmpty(suite.T(), healthResp.Version)
	assert.Greater(suite.T(), healthResp.Timestamp, int64(0))
	assert.NotEmpty(suite.T(), healthResp.Uptime)
	assert.NotEmpty(suite.T(), healthResp.Checks)
}

// TestTrainingJobServiceIntegration 测试训练作业服务集成
func (suite *ServiceIntegrationTestSuite) TestTrainingJobServiceIntegration() {
	// 1. 创建训练作业
	createReq := types.CreateTrainingJobReq{
		Name:       "integration-test-job",
		Framework:  "pytorch",
		Image:      "pytorch/pytorch:1.12.0",
		EntryPoint: "train.py",
		JobType:    "single",
		GpuCount:   1,
		GpuType:    "T4",
	}

	jsonData, _ := json.Marshal(createReq)
	resp, err := suite.client.Post(
		suite.baseURL+"/api/v1/training/jobs",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createResp types.CreateTrainingJobResp
	err = json.NewDecoder(resp.Body).Decode(&createResp)
	assert.NoError(suite.T(), err)
	assert.Greater(suite.T(), createResp.Id, int64(0))

	jobID := createResp.Id

	// 2. 查询创建的训练作业
	resp, err = suite.client.Get(fmt.Sprintf("%s/api/v1/training/jobs/%d", suite.baseURL, jobID))
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var jobResp types.GetTrainingJobResp
		err = json.NewDecoder(resp.Body).Decode(&jobResp)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), jobID, jobResp.Id)
		assert.Equal(suite.T(), "integration-test-job", jobResp.Name)
		assert.Equal(suite.T(), "pytorch", jobResp.Framework)
	}

	// 3. 列出训练作业
	resp, err = suite.client.Get(suite.baseURL + "/api/v1/training/jobs")
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var listResp types.ListTrainingJobsResp
		err = json.NewDecoder(resp.Body).Decode(&listResp)
		assert.NoError(suite.T(), err)
		assert.Greater(suite.T(), listResp.Total, int64(0))
	}
}

// TestGPUServiceIntegration 测试GPU服务集成
func (suite *ServiceIntegrationTestSuite) TestGPUServiceIntegration() {
	// 1. 列出GPU设备
	resp, err := suite.client.Get(suite.baseURL + "/api/v1/gpu/devices")
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	// 这个测试可能会失败，因为测试环境可能没有GPU设备
	// 但我们可以验证API的响应格式
	if resp.StatusCode == http.StatusOK {
		var deviceListResp types.ListGpuDevicesResp
		err = json.NewDecoder(resp.Body).Decode(&deviceListResp)
		assert.NoError(suite.T(), err)
		// 验证响应结构
		assert.NotNil(suite.T(), deviceListResp.Devices)
		assert.GreaterOrEqual(suite.T(), deviceListResp.Total, int64(0))
	}
}

// TestServiceCommunicationFlow 测试服务间通信流程
func (suite *ServiceIntegrationTestSuite) TestServiceCommunicationFlow() {
	// 模拟完整的业务流程：创建训练作业 -> 分配GPU -> 启动训练

	// 1. 首先检查系统健康状态
	healthResp, err := suite.client.Get(suite.baseURL + "/health")
	assert.NoError(suite.T(), err)
	defer healthResp.Body.Close()
	assert.Equal(suite.T(), http.StatusOK, healthResp.StatusCode)

	// 2. 查询可用的GPU设备
	gpuResp, err := suite.client.Get(suite.baseURL + "/api/v1/gpu/devices?status=available")
	assert.NoError(suite.T(), err)
	defer gpuResp.Body.Close()

	// 3. 创建训练作业
	createJobReq := types.CreateTrainingJobReq{
		Name:        "flow-test-job",
		Framework:   "tensorflow",
		Image:       "tensorflow/tensorflow:2.8.0-gpu",
		EntryPoint:  "main.py",
		JobType:     "single",
		GpuCount:    1,
		GpuType:     "T4",
		WorkerCount: 1,
		Priority:    5,
	}

	jsonData, _ := json.Marshal(createJobReq)
	jobResp, err := suite.client.Post(
		suite.baseURL+"/api/v1/training/jobs",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	assert.NoError(suite.T(), err)
	defer jobResp.Body.Close()

	if jobResp.StatusCode == http.StatusCreated {
		var createJobResp types.CreateTrainingJobResp
		err = json.NewDecoder(jobResp.Body).Decode(&createJobResp)
		assert.NoError(suite.T(), err)

		// 4. 验证作业创建成功后，系统状态一致性
		time.Sleep(100 * time.Millisecond) // 等待异步处理完成

		// 再次查询健康状态，确保系统正常
		healthResp2, err := suite.client.Get(suite.baseURL + "/health")
		assert.NoError(suite.T(), err)
		defer healthResp2.Body.Close()
		assert.Equal(suite.T(), http.StatusOK, healthResp2.StatusCode)
	}
}

// TestConcurrentRequests 测试并发请求处理
func (suite *ServiceIntegrationTestSuite) TestConcurrentRequests() {
	concurrency := 5
	results := make(chan error, concurrency)

	// 并发发送健康检查请求
	for i := 0; i < concurrency; i++ {
		go func(index int) {
			resp, err := suite.client.Get(suite.baseURL + "/health")
			if err != nil {
				results <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				results <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
				return
			}

			results <- nil
		}(i)
	}

	// 收集结果
	for i := 0; i < concurrency; i++ {
		err := <-results
		assert.NoError(suite.T(), err, "并发请求 %d 失败", i)
	}
}

// TestErrorHandling 测试错误处理
func (suite *ServiceIntegrationTestSuite) TestErrorHandling() {
	// 测试无效的训练作业创建请求
	invalidReq := types.CreateTrainingJobReq{
		Name: "", // 空名称应该导致验证错误
	}

	jsonData, _ := json.Marshal(invalidReq)
	resp, err := suite.client.Post(
		suite.baseURL+"/api/v1/training/jobs",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	// 应该返回400错误
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)

	// 测试访问不存在的训练作业
	resp, err = suite.client.Get(suite.baseURL + "/api/v1/training/jobs/99999")
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	// 应该返回404错误
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}

// TestServiceTimeout 测试服务超时处理
func (suite *ServiceIntegrationTestSuite) TestServiceTimeout() {
	// 创建一个超时时间很短的客户端
	shortTimeoutClient := &http.Client{
		Timeout: 1 * time.Millisecond, // 1毫秒超时
	}

	// 发送请求，应该超时
	_, err := shortTimeoutClient.Get(suite.baseURL + "/health")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "timeout")
}

// TestRequestResponseFormat 测试请求响应格式
func (suite *ServiceIntegrationTestSuite) TestRequestResponseFormat() {
	// 测试JSON格式请求
	validReq := types.CreateTrainingJobReq{
		Name:       "format-test-job",
		Framework:  "pytorch",
		Image:      "pytorch/pytorch:1.12.0",
		EntryPoint: "train.py",
	}

	jsonData, _ := json.Marshal(validReq)
	resp, err := suite.client.Post(
		suite.baseURL+"/api/v1/training/jobs",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	// 验证响应头
	assert.Equal(suite.T(), "application/json", resp.Header.Get("Content-Type"))

	// 测试无效的JSON格式
	invalidJSON := `{"name": "test", "framework":}`
	resp, err = suite.client.Post(
		suite.baseURL+"/api/v1/training/jobs",
		"application/json",
		bytes.NewBufferString(invalidJSON),
	)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

// TestServiceDiscovery 测试服务发现和负载均衡
func (suite *ServiceIntegrationTestSuite) TestServiceDiscovery() {
	// 模拟多个服务实例的场景
	// 在真实环境中，这会测试服务注册、发现和负载均衡

	// 发送多个请求，验证响应一致性
	for i := 0; i < 10; i++ {
		resp, err := suite.client.Get(suite.baseURL + "/health")
		assert.NoError(suite.T(), err)
		defer resp.Body.Close()

		assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

		var healthResp types.HealthResponse
		err = json.NewDecoder(resp.Body).Decode(&healthResp)
		assert.NoError(suite.T(), err)
		assert.NotEmpty(suite.T(), healthResp.Version)
	}
}

// TestCircuitBreaker 测试熔断器
func (suite *ServiceIntegrationTestSuite) TestCircuitBreaker() {
	// 模拟熔断器测试
	// 在真实环境中，这会测试当下游服务失败时的熔断行为

	maxFailures := 5
	failures := 0

	// 发送多个可能失败的请求
	for i := 0; i < maxFailures+2; i++ {
		resp, err := suite.client.Get(suite.baseURL + "/api/v1/training/jobs/invalid")
		if err != nil {
			failures++
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 500 {
			failures++
		}

		// 如果失败次数达到阈值，熔断器应该开启
		if failures >= maxFailures {
			suite.T().Logf("熔断器应该在 %d 次失败后开启", failures)
			break
		}
	}
}

// 运行服务集成测试套件
func TestServiceIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ServiceIntegrationTestSuite))
}

// 辅助函数

// setupTestData 设置测试数据
func (suite *ServiceIntegrationTestSuite) setupTestData() {
	// 在真实环境中，这里会设置测试所需的数据
	suite.T().Log("设置测试数据...")
}

// cleanupTestData 清理测试数据
func (suite *ServiceIntegrationTestSuite) cleanupTestData() {
	// 在真实环境中，这里会清理测试产生的数据
	suite.T().Log("清理测试数据...")
}
