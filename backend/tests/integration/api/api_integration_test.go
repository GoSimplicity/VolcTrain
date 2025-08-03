//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"api/internal/handler"
	"api/internal/svc"
	"api/internal/types"
	"api/test/config"
	"api/test/seeds"
	"api/test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/zeromicro/go-zero/rest"
)

// APIIntegrationTestSuite API集成测试套件
type APIIntegrationTestSuite struct {
	suite.Suite
	server     *rest.Server
	testConfig *config.TestConfig
	db         *sql.DB
	seeder     *seeds.DataSeeder
	adminToken string
	userToken  string
	baseURL    string
}

// SetupSuite 测试套件初始化
func (suite *APIIntegrationTestSuite) SetupSuite() {
	// 加载测试配置
	suite.testConfig = config.GetTestConfig()

	// 创建数据库连接
	var err error
	suite.db, err = suite.testConfig.CreateTestDB()
	require.NoError(suite.T(), err, "创建测试数据库连接失败")

	// 初始化数据种子
	suite.seeder = seeds.NewDataSeeder(suite.db)
	err = suite.seeder.SeedAllTestData()
	require.NoError(suite.T(), err, "生成测试数据失败")

	// 创建服务上下文
	serviceContext := &svc.ServiceContext{
		Config: config.Config{
			DataSource: fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				suite.testConfig.Database.Username,
				suite.testConfig.Database.Password,
				suite.testConfig.Database.Host,
				suite.testConfig.Database.Port,
				suite.testConfig.Database.Database),
		},
	}

	// 设置基础URL
	suite.baseURL = fmt.Sprintf("http://%s:%d", suite.testConfig.API.Host, suite.testConfig.API.Port)

	// 启动测试服务器
	suite.setupTestServer(serviceContext)

	// 获取测试用的Token
	suite.setupTestTokens()
}

// TearDownSuite 测试套件清理
func (suite *APIIntegrationTestSuite) TearDownSuite() {
	if suite.seeder != nil {
		suite.seeder.CleanupAllData()
	}
	if suite.db != nil {
		suite.db.Close()
	}
	if suite.server != nil {
		suite.server.Stop()
	}
}

// setupTestServer 设置测试服务器
func (suite *APIIntegrationTestSuite) setupTestServer(svcCtx *svc.ServiceContext) {
	// 这里应该设置实际的路由
	// 由于当前代码结构的限制，我们创建一个简化的测试服务器
	suite.server = rest.MustNewServer(rest.RestConf{
		Host: suite.testConfig.API.Host,
		Port: suite.testConfig.API.Port,
	})

	// 注册路由（这里需要根据实际的路由设置进行调整）
	handler.RegisterHandlers(suite.server, svcCtx)

	// 启动服务器
	go func() {
		suite.server.Start()
	}()

	// 等待服务器启动
	time.Sleep(2 * time.Second)
}

// setupTestTokens 设置测试Token
func (suite *APIIntegrationTestSuite) setupTestTokens() {
	// 管理员登录
	adminLoginReq := &types.LoginReq{
		Username: "admin_test",
		Password: "admin123",
	}

	adminResp := suite.postJSON("/api/v1/users/login", adminLoginReq)
	require.Equal(suite.T(), http.StatusOK, adminResp.Code, "管理员登录失败")

	var adminLoginResp types.LoginResp
	err := json.NewDecoder(adminResp.Body).Decode(&adminLoginResp)
	require.NoError(suite.T(), err)
	suite.adminToken = adminLoginResp.Token

	// 普通用户登录
	userLoginReq := &types.LoginReq{
		Username: "user_test",
		Password: "user123",
	}

	userResp := suite.postJSON("/api/v1/users/login", userLoginReq)
	require.Equal(suite.T(), http.StatusOK, userResp.Code, "用户登录失败")

	var userLoginResp types.LoginResp
	err = json.NewDecoder(userResp.Body).Decode(&userLoginResp)
	require.NoError(suite.T(), err)
	suite.userToken = userLoginResp.Token
}

// TestUserManagementAPI 测试用户管理API
func (suite *APIIntegrationTestSuite) TestUserManagementAPI() {
	suite.Run("创建用户", func() {
		req := &types.CreateUserReq{
			Username: "newuser",
			Email:    "newuser@test.com",
			Password: "NewUser123!",
			RealName: "新用户",
			UserType: "user",
		}

		resp := suite.postJSONWithAuth("/api/v1/users", req, suite.adminToken)
		assert.Equal(suite.T(), http.StatusCreated, resp.Code)

		var result types.CreateUserResp
		suite.decodeJSON(resp, &result)
		assert.Greater(suite.T(), result.Id, int64(0))
		assert.Equal(suite.T(), req.Username, result.Username)
	})

	suite.Run("获取用户列表", func() {
		resp := suite.getWithAuth("/api/v1/users?page=1&pageSize=10", suite.adminToken)
		assert.Equal(suite.T(), http.StatusOK, resp.Code)

		var result types.ListUsersResp
		suite.decodeJSON(resp, &result)
		assert.Greater(suite.T(), len(result.Data.List), 0)
		assert.Greater(suite.T(), result.Data.Total, int64(0))
	})

	suite.Run("获取用户详情", func() {
		resp := suite.getWithAuth("/api/v1/users/1001", suite.adminToken)
		assert.Equal(suite.T(), http.StatusOK, resp.Code)

		var result types.GetUserResp
		suite.decodeJSON(resp, &result)
		assert.Equal(suite.T(), int64(1001), result.User.Id)
		assert.Equal(suite.T(), "admin_test", result.User.Username)
	})

	suite.Run("更新用户信息", func() {
		req := &types.UpdateUserReq{
			Id:       1002,
			RealName: "更新后的真实姓名",
			Email:    "updated@test.com",
		}

		resp := suite.putJSONWithAuth("/api/v1/users/1002", req, suite.adminToken)
		assert.Equal(suite.T(), http.StatusOK, resp.Code)

		// 验证更新是否成功
		getResp := suite.getWithAuth("/api/v1/users/1002", suite.adminToken)
		var user types.GetUserResp
		suite.decodeJSON(getResp, &user)
		assert.Equal(suite.T(), req.RealName, user.User.RealName)
		assert.Equal(suite.T(), req.Email, user.User.Email)
	})

	suite.Run("删除用户", func() {
		// 先创建一个测试用户
		createReq := &types.CreateUserReq{
			Username: "tobedeleted",
			Email:    "delete@test.com",
			Password: "Delete123!",
			UserType: "user",
		}

		createResp := suite.postJSONWithAuth("/api/v1/users", createReq, suite.adminToken)
		var createResult types.CreateUserResp
		suite.decodeJSON(createResp, &createResult)

		// 删除用户
		deleteResp := suite.deleteWithAuth(fmt.Sprintf("/api/v1/users/%d", createResult.Id), suite.adminToken)
		assert.Equal(suite.T(), http.StatusOK, deleteResp.Code)

		// 验证用户已被删除（软删除）
		getResp := suite.getWithAuth(fmt.Sprintf("/api/v1/users/%d", createResult.Id), suite.adminToken)
		assert.Equal(suite.T(), http.StatusNotFound, getResp.Code)
	})
}

// TestTrainingJobAPI 测试训练作业API
func (suite *APIIntegrationTestSuite) TestTrainingJobAPI() {
	suite.Run("创建训练作业", func() {
		req := &types.CreateTrainingJobReq{
			Name:        "api-test-pytorch-job",
			Framework:   "pytorch",
			Image:       "pytorch/pytorch:1.12.0",
			EntryPoint:  "train.py",
			CPUCores:    "4",
			MemoryGb:    "8",
			GPUCount:    1,
			GPUType:     "T4",
			WorkspaceId: 3001,
			Description: "API集成测试训练作业",
		}

		resp := suite.postJSONWithAuth("/api/v1/training/jobs", req, suite.userToken)
		assert.Equal(suite.T(), http.StatusCreated, resp.Code)

		var result types.CreateTrainingJobResp
		suite.decodeJSON(resp, &result)
		assert.Greater(suite.T(), result.JobId, int64(0))
		assert.Equal(suite.T(), req.Name, result.Name)
		assert.Equal(suite.T(), "pending", result.Status)
	})

	suite.Run("获取训练作业列表", func() {
		resp := suite.getWithAuth("/api/v1/training/jobs?page=1&pageSize=10", suite.userToken)
		assert.Equal(suite.T(), http.StatusOK, resp.Code)

		var result types.ListTrainingJobsResp
		suite.decodeJSON(resp, &result)
		assert.Greater(suite.T(), len(result.Data.List), 0)
	})

	suite.Run("获取训练作业详情", func() {
		resp := suite.getWithAuth("/api/v1/training/jobs/3001", suite.userToken)
		assert.Equal(suite.T(), http.StatusOK, resp.Code)

		var result types.GetTrainingJobResp
		suite.decodeJSON(resp, &result)
		assert.Equal(suite.T(), int64(3001), result.Job.Id)
		assert.Equal(suite.T(), "pytorch-training-test", result.Job.Name)
	})

	suite.Run("更新训练作业状态", func() {
		req := &types.UpdateTrainingJobStatusReq{
			JobId:  3001,
			Status: "running",
		}

		resp := suite.putJSONWithAuth("/api/v1/training/jobs/3001/status", req, suite.adminToken)
		assert.Equal(suite.T(), http.StatusOK, resp.Code)

		// 验证状态更新
		getResp := suite.getWithAuth("/api/v1/training/jobs/3001", suite.userToken)
		var job types.GetTrainingJobResp
		suite.decodeJSON(getResp, &job)
		assert.Equal(suite.T(), "running", job.Job.Status)
	})

	suite.Run("取消训练作业", func() {
		resp := suite.postWithAuth("/api/v1/training/jobs/3001/cancel", nil, suite.userToken)
		assert.Equal(suite.T(), http.StatusOK, resp.Code)

		// 验证作业状态变为cancelled
		getResp := suite.getWithAuth("/api/v1/training/jobs/3001", suite.userToken)
		var job types.GetTrainingJobResp
		suite.decodeJSON(getResp, &job)
		assert.Equal(suite.T(), "cancelled", job.Job.Status)
	})
}

// TestGPUManagementAPI 测试GPU管理API
func (suite *APIIntegrationTestSuite) TestGPUManagementAPI() {
	suite.Run("创建GPU集群", func() {
		req := &types.CreateGpuClusterReq{
			Name:        "api-test-cluster",
			Description: "API测试GPU集群",
			Kubeconfig:  "test-kubeconfig-content",
			Region:      "us-west-1",
			Provider:    "aws",
		}

		resp := suite.postJSONWithAuth("/api/v1/gpu/clusters", req, suite.adminToken)
		assert.Equal(suite.T(), http.StatusCreated, resp.Code)

		var result types.CreateGpuClusterResp
		suite.decodeJSON(resp, &result)
		assert.Greater(suite.T(), result.ClusterId, int64(0))
		assert.Equal(suite.T(), req.Name, result.Name)
	})

	suite.Run("获取GPU集群列表", func() {
		resp := suite.getWithAuth("/api/v1/gpu/clusters", suite.adminToken)
		assert.Equal(suite.T(), http.StatusOK, resp.Code)

		var result types.ListGpuClustersResp
		suite.decodeJSON(resp, &result)
		assert.Greater(suite.T(), len(result.Data.List), 0)
	})

	suite.Run("获取GPU集群详情", func() {
		resp := suite.getWithAuth("/api/v1/gpu/clusters/2001", suite.adminToken)
		assert.Equal(suite.T(), http.StatusOK, resp.Code)

		var result types.GetGpuClusterResp
		suite.decodeJSON(resp, &result)
		assert.Equal(suite.T(), int64(2001), result.Cluster.Id)
		assert.Equal(suite.T(), "test-cluster-1", result.Cluster.Name)
	})

	suite.Run("分配GPU设备", func() {
		req := &types.AllocateGpuDeviceReq{
			JobId:     3001,
			ClusterId: 2001,
			GPUCount:  1,
			GPUType:   "T4",
		}

		resp := suite.postJSONWithAuth("/api/v1/gpu/devices/allocate", req, suite.adminToken)
		assert.Equal(suite.T(), http.StatusOK, resp.Code)

		var result types.AllocateGpuDeviceResp
		suite.decodeJSON(resp, &result)
		assert.Equal(suite.T(), "allocated", result.Status)
		assert.Greater(suite.T(), len(result.AllocatedDevices), 0)
	})

	suite.Run("释放GPU设备", func() {
		req := &types.ReleaseGpuDeviceReq{
			JobId: 3001,
		}

		resp := suite.postJSONWithAuth("/api/v1/gpu/devices/release", req, suite.adminToken)
		assert.Equal(suite.T(), http.StatusOK, resp.Code)
	})
}

// TestAuthenticationAPI 测试认证API
func (suite *APIIntegrationTestSuite) TestAuthenticationAPI() {
	suite.Run("用户登录成功", func() {
		req := &types.LoginReq{
			Username: "user_test",
			Password: "user123",
		}

		resp := suite.postJSON("/api/v1/users/login", req)
		assert.Equal(suite.T(), http.StatusOK, resp.Code)

		var result types.LoginResp
		suite.decodeJSON(resp, &result)
		assert.NotEmpty(suite.T(), result.Token)
		assert.Equal(suite.T(), req.Username, result.Username)
	})

	suite.Run("用户登录失败-错误密码", func() {
		req := &types.LoginReq{
			Username: "user_test",
			Password: "wrongpassword",
		}

		resp := suite.postJSON("/api/v1/users/login", req)
		assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
	})

	suite.Run("获取用户信息", func() {
		resp := suite.getWithAuth("/api/v1/users/profile", suite.userToken)
		assert.Equal(suite.T(), http.StatusOK, resp.Code)

		var result types.GetUserProfileResp
		suite.decodeJSON(resp, &result)
		assert.Equal(suite.T(), "user_test", result.User.Username)
	})

	suite.Run("刷新Token", func() {
		req := &types.RefreshTokenReq{
			Token: suite.userToken,
		}

		resp := suite.postJSON("/api/v1/users/refresh", req)
		assert.Equal(suite.T(), http.StatusOK, resp.Code)

		var result types.RefreshTokenResp
		suite.decodeJSON(resp, &result)
		assert.NotEmpty(suite.T(), result.Token)
		assert.NotEqual(suite.T(), suite.userToken, result.Token)
	})
}

// TestAPIPermissions 测试API权限控制
func (suite *APIIntegrationTestSuite) TestAPIPermissions() {
	suite.Run("未认证访问受保护资源", func() {
		resp := suite.get("/api/v1/users/profile")
		assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)
	})

	suite.Run("普通用户访问管理员资源", func() {
		resp := suite.getWithAuth("/api/v1/admin/users", suite.userToken)
		assert.Equal(suite.T(), http.StatusForbidden, resp.Code)
	})

	suite.Run("管理员访问管理员资源", func() {
		resp := suite.getWithAuth("/api/v1/admin/users", suite.adminToken)
		assert.Equal(suite.T(), http.StatusOK, resp.Code)
	})

	suite.Run("用户只能访问自己的资源", func() {
		// 用户尝试访问别人的训练作业
		resp := suite.getWithAuth("/api/v1/training/jobs/3002", suite.userToken)
		// 这应该失败或者只返回属于该用户的作业
		assert.Contains(suite.T(), []int{http.StatusForbidden, http.StatusNotFound}, resp.Code)
	})
}

// TestAPIValidation 测试API参数验证
func (suite *APIIntegrationTestSuite) TestAPIValidation() {
	suite.Run("创建用户-缺少必要字段", func() {
		req := &types.CreateUserReq{
			Username: "",              // 空用户名
			Email:    "invalid-email", // 无效邮箱
			Password: "123",           // 密码太短
		}

		resp := suite.postJSONWithAuth("/api/v1/users", req, suite.adminToken)
		assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
	})

	suite.Run("创建训练作业-无效参数", func() {
		req := &types.CreateTrainingJobReq{
			Name:      "",                  // 空名称
			Framework: "invalid-framework", // 无效框架
			GPUCount:  -1,                  // 负数GPU
		}

		resp := suite.postJSONWithAuth("/api/v1/training/jobs", req, suite.userToken)
		assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
	})

	suite.Run("分页参数验证", func() {
		// 无效的分页参数
		resp := suite.getWithAuth("/api/v1/users?page=-1&pageSize=0", suite.adminToken)
		assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)
	})
}

// HTTP请求辅助方法
func (suite *APIIntegrationTestSuite) postJSON(path string, data interface{}) *httptest.ResponseRecorder {
	return suite.makeRequest("POST", path, data, "")
}

func (suite *APIIntegrationTestSuite) postJSONWithAuth(path string, data interface{}, token string) *httptest.ResponseRecorder {
	return suite.makeRequest("POST", path, data, token)
}

func (suite *APIIntegrationTestSuite) putJSONWithAuth(path string, data interface{}, token string) *httptest.ResponseRecorder {
	return suite.makeRequest("PUT", path, data, token)
}

func (suite *APIIntegrationTestSuite) get(path string) *httptest.ResponseRecorder {
	return suite.makeRequest("GET", path, nil, "")
}

func (suite *APIIntegrationTestSuite) getWithAuth(path string, token string) *httptest.ResponseRecorder {
	return suite.makeRequest("GET", path, nil, token)
}

func (suite *APIIntegrationTestSuite) deleteWithAuth(path string, token string) *httptest.ResponseRecorder {
	return suite.makeRequest("DELETE", path, nil, token)
}

func (suite *APIIntegrationTestSuite) postWithAuth(path string, data interface{}, token string) *httptest.ResponseRecorder {
	return suite.makeRequest("POST", path, data, token)
}

func (suite *APIIntegrationTestSuite) makeRequest(method, path string, data interface{}, token string) *httptest.ResponseRecorder {
	var body *bytes.Buffer

	if data != nil {
		jsonData, err := json.Marshal(data)
		require.NoError(suite.T(), err)
		body = bytes.NewBuffer(jsonData)
	} else {
		body = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, body)
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	w := httptest.NewRecorder()

	// 这里需要根据实际的HTTP处理器进行调整
	// suite.server.ServeHTTP(w, req)

	return w
}

func (suite *APIIntegrationTestSuite) decodeJSON(w *httptest.ResponseRecorder, v interface{}) {
	err := json.NewDecoder(w.Body).Decode(v)
	require.NoError(suite.T(), err, "解析JSON响应失败")
}

// 运行API集成测试套件
func TestAPIIntegrationSuite(t *testing.T) {
	suite.Run(t, new(APIIntegrationTestSuite))
}
