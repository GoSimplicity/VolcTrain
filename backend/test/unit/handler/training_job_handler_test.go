//go:build unit
// +build unit

package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api/internal/logic/training"
	"api/internal/svc"
	"api/internal/types"
	"api/model"
	"api/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

// mockResult 实现sql.Result接口
type mockResult struct {
	lastInsertId int64
	rowsAffected int64
}

func (r *mockResult) LastInsertId() (int64, error) {
	return r.lastInsertId, nil
}

func (r *mockResult) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

// TrainingJobHandlerTestSuite 训练作业处理器测试套件
type TrainingJobHandlerTestSuite struct {
	suite.Suite
	ctrl              *gomock.Controller
	mockTrainingModel *mocks.MockVtTrainingJobsModel
	svcCtx            *svc.ServiceContext
	handler           http.HandlerFunc
	ctx               context.Context
}

// SetupTest 每个测试前的初始化
func (suite *TrainingJobHandlerTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockTrainingModel = mocks.NewMockVtTrainingJobsModel(suite.ctrl)
	suite.ctx = context.Background()

	// 创建模拟的服务上下文
	suite.svcCtx = &svc.ServiceContext{
		VtTrainingJobsModel: suite.mockTrainingModel,
	}

	// 创建训练作业创建处理器
	suite.handler = func(w http.ResponseWriter, r *http.Request) {
		// 解析请求体
		var req types.CreateTrainingJobReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "无效的请求格式", http.StatusBadRequest)
			return
		}

		// 创建逻辑处理器
		createLogic := training.NewCreateTrainingJobLogic(r.Context(), suite.svcCtx)
		resp, err := createLogic.CreateTrainingJob(&req)

		if err != nil {
			// 根据错误类型返回不同的状态码
			if isValidationError(err) {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else if isConflictError(err) {
				http.Error(w, err.Error(), http.StatusConflict)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// 返回成功响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}

// TearDownTest 每个测试后的清理
func (suite *TrainingJobHandlerTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

// TestCreateTrainingJobHandlerSuccess 测试成功创建训练作业
func (suite *TrainingJobHandlerTestSuite) TestCreateTrainingJobHandlerSuccess() {
	// 准备测试数据
	requestData := types.CreateTrainingJobReq{
		Name:       "test-pytorch-job",
		Framework:  "pytorch",
		Image:      "pytorch/pytorch:1.12.0",
		EntryPoint: "train.py",
		JobType:    "single",
		GpuCount:   1,
		GpuType:    "T4",
	}

	// 设置Mock期望
	suite.mockTrainingModel.EXPECT().
		FindOneByName("test-pytorch-job").
		Return(nil, sql.ErrNoRows).
		Times(1)

	suite.mockTrainingModel.EXPECT().
		Insert(gomock.Any()).
		Return(&mockResult{lastInsertId: 12345, rowsAffected: 1}, nil).
		Times(1)

	// 创建HTTP请求
	jsonData, _ := json.Marshal(requestData)
	req := httptest.NewRequest("POST", "/api/v1/training/jobs", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行处理器
	suite.handler(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	assert.Equal(suite.T(), "application/json", w.Header().Get("Content-Type"))

	// 解析响应体
	var response types.CreateTrainingJobResp
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(12345), response.Id)
}

// TestCreateTrainingJobHandlerValidationError 测试请求参数验证错误
func (suite *TrainingJobHandlerTestSuite) TestCreateTrainingJobHandlerValidationError() {
	testCases := []struct {
		name         string
		requestData  types.CreateTrainingJobReq
		expectedCode int
		expectedMsg  string
	}{
		{
			name: "空名称",
			requestData: types.CreateTrainingJobReq{
				Name:       "",
				Framework:  "pytorch",
				Image:      "pytorch/pytorch:1.12.0",
				EntryPoint: "train.py",
			},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "训练作业名称不能为空",
		},
		{
			name: "空框架",
			requestData: types.CreateTrainingJobReq{
				Name:       "test-job",
				Framework:  "",
				Image:      "pytorch/pytorch:1.12.0",
				EntryPoint: "train.py",
			},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "训练框架不能为空",
		},
		{
			name: "GPU配置错误",
			requestData: types.CreateTrainingJobReq{
				Name:       "test-job",
				Framework:  "pytorch",
				Image:      "pytorch/pytorch:1.12.0",
				EntryPoint: "train.py",
				GpuCount:   2,
				GpuType:    "",
			},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "指定GPU数量时必须指定GPU类型",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// 创建HTTP请求
			jsonData, _ := json.Marshal(tc.requestData)
			req := httptest.NewRequest("POST", "/api/v1/training/jobs", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// 执行处理器
			suite.handler(w, req)

			// 验证响应
			assert.Equal(suite.T(), tc.expectedCode, w.Code)
			assert.Contains(suite.T(), w.Body.String(), tc.expectedMsg)
		})
	}
}

// TestCreateTrainingJobHandlerNameConflict 测试作业名称冲突
func (suite *TrainingJobHandlerTestSuite) TestCreateTrainingJobHandlerNameConflict() {
	// 准备测试数据
	requestData := types.CreateTrainingJobReq{
		Name:       "existing-job",
		Framework:  "pytorch",
		Image:      "pytorch/pytorch:1.12.0",
		EntryPoint: "train.py",
	}

	// 设置Mock期望 - 返回已存在的作业
	suite.mockTrainingModel.EXPECT().
		FindOneByName("existing-job").
		Return(&model.VtTrainingJobs{}, nil).
		Times(1)

	// 创建HTTP请求
	jsonData, _ := json.Marshal(requestData)
	req := httptest.NewRequest("POST", "/api/v1/training/jobs", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行处理器
	suite.handler(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusConflict, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "训练作业名称 'existing-job' 已存在")
}

// TestCreateTrainingJobHandlerInvalidJSON 测试无效的JSON格式
func (suite *TrainingJobHandlerTestSuite) TestCreateTrainingJobHandlerInvalidJSON() {
	// 创建无效的JSON数据
	invalidJSON := `{"name": "test-job", "framework": "pytorch", "image": }`

	req := httptest.NewRequest("POST", "/api/v1/training/jobs", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行处理器
	suite.handler(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "无效的请求格式")
}

// TestCreateTrainingJobHandlerDatabaseError 测试数据库错误
func (suite *TrainingJobHandlerTestSuite) TestCreateTrainingJobHandlerDatabaseError() {
	// 准备测试数据
	requestData := types.CreateTrainingJobReq{
		Name:       "test-job",
		Framework:  "pytorch",
		Image:      "pytorch/pytorch:1.12.0",
		EntryPoint: "train.py",
	}

	// 设置Mock期望 - 数据库操作失败
	suite.mockTrainingModel.EXPECT().
		FindOneByName("test-job").
		Return(nil, sql.ErrNoRows).
		Times(1)

	suite.mockTrainingModel.EXPECT().
		Insert(gomock.Any()).
		Return(nil, assert.AnError).
		Times(1)

	// 创建HTTP请求
	jsonData, _ := json.Marshal(requestData)
	req := httptest.NewRequest("POST", "/api/v1/training/jobs", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行处理器
	suite.handler(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "保存训练作业失败")
}

// TestCreateTrainingJobHandlerMissingContentType 测试缺少Content-Type头
func (suite *TrainingJobHandlerTestSuite) TestCreateTrainingJobHandlerMissingContentType() {
	// 准备测试数据
	requestData := types.CreateTrainingJobReq{
		Name:       "test-job",
		Framework:  "pytorch",
		Image:      "pytorch/pytorch:1.12.0",
		EntryPoint: "train.py",
	}

	// 创建HTTP请求但不设置Content-Type
	jsonData, _ := json.Marshal(requestData)
	req := httptest.NewRequest("POST", "/api/v1/training/jobs", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()

	// 设置Mock期望（如果请求被处理）
	suite.mockTrainingModel.EXPECT().
		FindOneByName("test-job").
		Return(nil, sql.ErrNoRows).
		Times(1)

	suite.mockTrainingModel.EXPECT().
		Insert(gomock.Any()).
		Return(&mockResult{lastInsertId: 12345, rowsAffected: 1}, nil).
		Times(1)

	// 执行处理器
	suite.handler(w, req)

	// 验证响应（应该仍然成功，因为Go的json.Decoder可以处理）
	assert.True(suite.T(), w.Code == http.StatusCreated || w.Code == http.StatusBadRequest)
}

// TestCreateTrainingJobHandlerLargePayload 测试大负载请求
func (suite *TrainingJobHandlerTestSuite) TestCreateTrainingJobHandlerLargePayload() {
	// 准备包含大量数据的测试请求
	requestData := types.CreateTrainingJobReq{
		Name:        "test-large-job",
		Framework:   "pytorch",
		Image:       "pytorch/pytorch:1.12.0",
		EntryPoint:  "train.py",
		Description: generateLargeString(1000), // 1KB description
		EnvVars:     generateLargeString(5000), // 5KB env vars
	}

	// 设置Mock期望
	suite.mockTrainingModel.EXPECT().
		FindOneByName("test-large-job").
		Return(nil, sql.ErrNoRows).
		Times(1)

	suite.mockTrainingModel.EXPECT().
		Insert(gomock.Any()).
		Return(&mockResult{lastInsertId: 12345, rowsAffected: 1}, nil).
		Times(1)

	// 创建HTTP请求
	jsonData, _ := json.Marshal(requestData)
	req := httptest.NewRequest("POST", "/api/v1/training/jobs", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 执行处理器
	suite.handler(w, req)

	// 验证响应
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	// 解析响应体
	var response types.CreateTrainingJobResp
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(12345), response.Id)
}

// TestCreateTrainingJobHandlerConcurrency 测试并发请求
func (suite *TrainingJobHandlerTestSuite) TestCreateTrainingJobHandlerConcurrency() {
	// 设置Mock期望 - 允许多次调用
	suite.mockTrainingModel.EXPECT().
		FindOneByName(gomock.Any()).
		Return(nil, sql.ErrNoRows).
		AnyTimes()

	suite.mockTrainingModel.EXPECT().
		Insert(gomock.Any()).
		Return(&mockResult{lastInsertId: 12345, rowsAffected: 1}, nil).
		AnyTimes()

	// 并发发送多个请求
	concurrency := 5
	results := make(chan int, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(index int) {
			requestData := types.CreateTrainingJobReq{
				Name:       fmt.Sprintf("concurrent-job-%d", index),
				Framework:  "pytorch",
				Image:      "pytorch/pytorch:1.12.0",
				EntryPoint: "train.py",
			}

			jsonData, _ := json.Marshal(requestData)
			req := httptest.NewRequest("POST", "/api/v1/training/jobs", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			suite.handler(w, req)
			results <- w.Code
		}(i)
	}

	// 收集结果
	for i := 0; i < concurrency; i++ {
		statusCode := <-results
		assert.Equal(suite.T(), http.StatusCreated, statusCode)
	}
}

// 辅助函数

// isValidationError 检查是否为验证错误
func isValidationError(err error) bool {
	msg := err.Error()
	validationKeywords := []string{
		"不能为空", "无效", "格式不正确", "必须指定", "不能为负数",
	}
	for _, keyword := range validationKeywords {
		if strings.Contains(msg, keyword) {
			return true
		}
	}
	return false
}

// isConflictError 检查是否为冲突错误
func isConflictError(err error) bool {
	return strings.Contains(err.Error(), "已存在")
}

// generateLargeString 生成指定大小的字符串
func generateLargeString(size int) string {
	data := make([]byte, size)
	for i := range data {
		data[i] = 'a' + byte(i%26)
	}
	return string(data)
}

// 运行训练作业处理器测试套件
func TestTrainingJobHandlerSuite(t *testing.T) {
	suite.Run(t, new(TrainingJobHandlerTestSuite))
}

// 基准测试
func BenchmarkCreateTrainingJobHandler(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockTrainingModel := mocks.NewMockVtTrainingJobsModel(ctrl)
	svcCtx := &svc.ServiceContext{
		VtTrainingJobsModel: mockTrainingModel,
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateTrainingJobReq
		json.NewDecoder(r.Body).Decode(&req)

		createLogic := training.NewCreateTrainingJobLogic(r.Context(), svcCtx)
		resp, err := createLogic.CreateTrainingJob(&req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}

	// 设置Mock期望
	mockTrainingModel.EXPECT().
		FindOneByName(gomock.Any()).
		Return(nil, sql.ErrNoRows).
		AnyTimes()

	mockTrainingModel.EXPECT().
		Insert(gomock.Any()).
		Return(&mockResult{lastInsertId: 1, rowsAffected: 1}, nil).
		AnyTimes()

	requestData := types.CreateTrainingJobReq{
		Name:       "benchmark-job",
		Framework:  "pytorch",
		Image:      "pytorch/pytorch:1.12.0",
		EntryPoint: "train.py",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsonData, _ := json.Marshal(requestData)
		req := httptest.NewRequest("POST", "/api/v1/training/jobs", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler(w, req)
	}
}
