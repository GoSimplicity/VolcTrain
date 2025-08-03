//go:build unit
// +build unit

package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommonResponse_Structure(t *testing.T) {
	resp := CommonResponse{
		Code:      StatusOK,
		Message:   "test message",
		Data:      map[string]string{"key": "value"},
		Timestamp: time.Now().Unix(),
		RequestID: "test-request-id",
	}

	assert.Equal(t, StatusOK, resp.Code)
	assert.Equal(t, "test message", resp.Message)
	assert.NotNil(t, resp.Data)
	assert.NotZero(t, resp.Timestamp)
	assert.Equal(t, "test-request-id", resp.RequestID)
}

func TestListResponse_Structure(t *testing.T) {
	items := []string{"item1", "item2", "item3"}
	resp := ListResponse{
		Items:    items,
		Total:    int64(len(items)),
		Page:     1,
		PageSize: 10,
	}

	assert.Equal(t, items, resp.Items)
	assert.Equal(t, int64(3), resp.Total)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 10, resp.PageSize)
}

func TestSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	testData := map[string]string{"test": "data"}

	Success(w, testData)

	assert.Equal(t, http.StatusOK, w.Code)

	var response CommonResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, StatusOK, response.Code)
	assert.Equal(t, "success", response.Message)
	assert.NotNil(t, response.Data)
	assert.NotZero(t, response.Timestamp)
}

func TestCreated(t *testing.T) {
	w := httptest.NewRecorder()
	testData := map[string]string{"created": "resource"}

	Created(w, testData)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response CommonResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, StatusCreated, response.Code)
	assert.Equal(t, "创建成功", response.Message)
	assert.NotNil(t, response.Data)
}

func TestBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	errorMessage := "Invalid parameters"

	BadRequest(w, errorMessage)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response CommonResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, StatusBadRequest, response.Code)
	assert.Equal(t, errorMessage, response.Message)
	assert.Nil(t, response.Data)
}

func TestUnauthorized(t *testing.T) {
	tests := []struct {
		name            string
		inputMessage    string
		expectedMessage string
	}{
		{
			name:            "自定义消息",
			inputMessage:    "Custom unauthorized message",
			expectedMessage: "Custom unauthorized message",
		},
		{
			name:            "空消息使用默认值",
			inputMessage:    "",
			expectedMessage: "未认证或token已过期",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			Unauthorized(w, tt.inputMessage)

			assert.Equal(t, http.StatusUnauthorized, w.Code)

			var response CommonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, StatusUnauthorized, response.Code)
			assert.Equal(t, tt.expectedMessage, response.Message)
		})
	}
}

func TestForbidden(t *testing.T) {
	tests := []struct {
		name            string
		inputMessage    string
		expectedMessage string
	}{
		{
			name:            "自定义消息",
			inputMessage:    "Access denied",
			expectedMessage: "Access denied",
		},
		{
			name:            "空消息使用默认值",
			inputMessage:    "",
			expectedMessage: "权限不足",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			Forbidden(w, tt.inputMessage)

			assert.Equal(t, http.StatusForbidden, w.Code)

			var response CommonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, StatusForbidden, response.Code)
			assert.Equal(t, tt.expectedMessage, response.Message)
		})
	}
}

func TestNotFound(t *testing.T) {
	tests := []struct {
		name            string
		inputMessage    string
		expectedMessage string
	}{
		{
			name:            "自定义消息",
			inputMessage:    "User not found",
			expectedMessage: "User not found",
		},
		{
			name:            "空消息使用默认值",
			inputMessage:    "",
			expectedMessage: "资源不存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			NotFound(w, tt.inputMessage)

			assert.Equal(t, http.StatusNotFound, w.Code)

			var response CommonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, StatusNotFound, response.Code)
			assert.Equal(t, tt.expectedMessage, response.Message)
		})
	}
}

func TestConflict(t *testing.T) {
	w := httptest.NewRecorder()
	errorMessage := "Resource already exists"

	Conflict(w, errorMessage)

	assert.Equal(t, http.StatusConflict, w.Code)

	var response CommonResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, StatusConflict, response.Code)
	assert.Equal(t, errorMessage, response.Message)
}

func TestTooManyRequests(t *testing.T) {
	tests := []struct {
		name            string
		inputMessage    string
		expectedMessage string
	}{
		{
			name:            "自定义消息",
			inputMessage:    "Rate limit exceeded",
			expectedMessage: "Rate limit exceeded",
		},
		{
			name:            "空消息使用默认值",
			inputMessage:    "",
			expectedMessage: "请求过于频繁，请稍后再试",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			TooManyRequests(w, tt.inputMessage)

			assert.Equal(t, http.StatusTooManyRequests, w.Code)

			var response CommonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, StatusTooManyRequests, response.Code)
			assert.Equal(t, tt.expectedMessage, response.Message)
		})
	}
}

func TestInternalServerError(t *testing.T) {
	tests := []struct {
		name            string
		inputMessage    string
		expectedMessage string
	}{
		{
			name:            "自定义消息",
			inputMessage:    "Database connection failed",
			expectedMessage: "Database connection failed",
		},
		{
			name:            "空消息使用默认值",
			inputMessage:    "",
			expectedMessage: "服务器内部错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			InternalServerError(w, tt.inputMessage)

			assert.Equal(t, http.StatusInternalServerError, w.Code)

			var response CommonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, StatusInternalServerError, response.Code)
			assert.Equal(t, tt.expectedMessage, response.Message)
		})
	}
}

func TestServiceUnavailable(t *testing.T) {
	tests := []struct {
		name            string
		inputMessage    string
		expectedMessage string
	}{
		{
			name:            "自定义消息",
			inputMessage:    "Maintenance mode",
			expectedMessage: "Maintenance mode",
		},
		{
			name:            "空消息使用默认值",
			inputMessage:    "",
			expectedMessage: "服务暂时不可用",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			ServiceUnavailable(w, tt.inputMessage)

			assert.Equal(t, http.StatusServiceUnavailable, w.Code)

			var response CommonResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, StatusServiceUnavailable, response.Code)
			assert.Equal(t, tt.expectedMessage, response.Message)
		})
	}
}

func TestPagedSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	items := []string{"item1", "item2", "item3"}
	total := int64(100)
	page := 2
	pageSize := 10

	PagedSuccess(w, items, total, page, pageSize)

	assert.Equal(t, http.StatusOK, w.Code)

	var response PagedResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, StatusOK, response.Code)
	assert.Equal(t, "success", response.Message)
	// 验证数据内容
	itemsFromResponse := response.Data.Items.([]interface{})
	for i, item := range items {
		assert.Equal(t, item, itemsFromResponse[i])
	}
	assert.Equal(t, total, response.Data.Total)
	assert.Equal(t, page, response.Data.Page)
	assert.Equal(t, pageSize, response.Data.PageSize)
	assert.NotZero(t, response.Timestamp)
}

func TestCustomResponse(t *testing.T) {
	w := httptest.NewRecorder()
	httpStatus := http.StatusAccepted
	code := 202
	message := "Accepted"
	data := map[string]string{"status": "processing"}

	CustomResponse(w, httpStatus, code, message, data)

	assert.Equal(t, httpStatus, w.Code)

	var response CommonResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, code, response.Code)
	assert.Equal(t, message, response.Message)
	// 验证数据内容
	dataFromResponse := response.Data.(map[string]interface{})
	assert.Equal(t, "processing", dataFromResponse["status"])
	assert.NotZero(t, response.Timestamp)
}

func TestWithRequestID(t *testing.T) {
	response := &CommonResponse{
		Code:      StatusOK,
		Message:   "test",
		Timestamp: time.Now().Unix(),
	}

	requestID := "test-request-123"
	WithRequestID(response, requestID)

	assert.Equal(t, requestID, response.RequestID)
}

func TestStatusConstants(t *testing.T) {
	assert.Equal(t, 200, StatusOK)
	assert.Equal(t, 201, StatusCreated)
	assert.Equal(t, 400, StatusBadRequest)
	assert.Equal(t, 401, StatusUnauthorized)
	assert.Equal(t, 403, StatusForbidden)
	assert.Equal(t, 404, StatusNotFound)
	assert.Equal(t, 409, StatusConflict)
	assert.Equal(t, 429, StatusTooManyRequests)
	assert.Equal(t, 500, StatusInternalServerError)
	assert.Equal(t, 503, StatusServiceUnavailable)
}

func TestResponseJSONFormat(t *testing.T) {
	w := httptest.NewRecorder()
	testData := map[string]interface{}{
		"id":     123,
		"name":   "test",
		"active": true,
	}

	Success(w, testData)

	// 验证Content-Type是application/json
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	// 验证JSON格式正确
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// 验证必要字段存在
	assert.Contains(t, response, "code")
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "data")
	assert.Contains(t, response, "timestamp")
}

// 基准测试
func BenchmarkSuccess(b *testing.B) {
	testData := map[string]string{"test": "data"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		Success(w, testData)
	}
}

func BenchmarkPagedSuccess(b *testing.B) {
	items := make([]string, 100)
	for i := range items {
		items[i] = "item"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		PagedSuccess(w, items, 1000, 1, 100)
	}
}
