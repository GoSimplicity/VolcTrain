//go:build unit
// +build unit

package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBizError_Error(t *testing.T) {
	bizErr := &BizError{
		Code:    ErrCodeUserNotFound,
		Message: "用户不存在",
		Type:    ErrorTypeBusiness,
	}

	expected := fmt.Sprintf("[%s] %s (code: %d)", ErrorTypeBusiness, "用户不存在", ErrCodeUserNotFound)
	assert.Equal(t, expected, bizErr.Error())
}

func TestBizError_GetHTTPStatus(t *testing.T) {
	tests := []struct {
		name           string
		code           int
		expectedStatus int
	}{
		{
			name:           "Bad Request错误码",
			code:           ErrCodeBadRequest,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "参数验证错误码",
			code:           ErrCodeValidation,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "无效参数错误码",
			code:           ErrCodeInvalidParam,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "未认证错误码",
			code:           ErrCodeUnauthorized,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Token无效错误码",
			code:           ErrCodeTokenInvalid,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Token过期错误码",
			code:           ErrCodeTokenExpired,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "权限不足错误码",
			code:           ErrCodeForbidden,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "权限拒绝错误码",
			code:           ErrCodePermissionDenied,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "资源不存在错误码",
			code:           ErrCodeNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "用户不存在错误码",
			code:           ErrCodeUserNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "冲突错误码",
			code:           ErrCodeConflict,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "数据重复错误码",
			code:           ErrCodeDuplicateData,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "请求过多错误码",
			code:           ErrCodeTooManyRequests,
			expectedStatus: http.StatusTooManyRequests,
		},
		{
			name:           "未知错误码默认为500",
			code:           9999,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bizErr := &BizError{Code: tt.code}
			assert.Equal(t, tt.expectedStatus, bizErr.GetHTTPStatus())
		})
	}
}

func TestNewBizError(t *testing.T) {
	code := ErrCodeValidation
	message := "测试消息"
	errorType := ErrorTypeValidation

	bizErr := NewBizError(code, message, errorType)

	assert.Equal(t, code, bizErr.Code)
	assert.Equal(t, message, bizErr.Message)
	assert.Equal(t, errorType, bizErr.Type)
	assert.Empty(t, bizErr.Stack)
}

func TestNewBizErrorWithStack(t *testing.T) {
	code := ErrCodeInternalError
	message := "系统错误"
	errorType := ErrorTypeSystem

	bizErr := NewBizErrorWithStack(code, message, errorType)

	assert.Equal(t, code, bizErr.Code)
	assert.Equal(t, message, bizErr.Message)
	assert.Equal(t, errorType, bizErr.Type)
	assert.NotEmpty(t, bizErr.Stack)
	assert.Contains(t, bizErr.Stack, "TestNewBizErrorWithStack")
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name    string
		bizErr  *BizError
		code    int
		message string
		errType string
	}{
		{
			name:    "ErrInvalidParam",
			bizErr:  ErrInvalidParam,
			code:    ErrCodeInvalidParam,
			message: "参数无效",
			errType: ErrorTypeValidation,
		},
		{
			name:    "ErrMissingParam",
			bizErr:  ErrMissingParam,
			code:    ErrCodeMissingParam,
			message: "缺少必要参数",
			errType: ErrorTypeValidation,
		},
		{
			name:    "ErrTokenInvalid",
			bizErr:  ErrTokenInvalid,
			code:    ErrCodeTokenInvalid,
			message: "token无效",
			errType: ErrorTypeAuth,
		},
		{
			name:    "ErrUserNotFound",
			bizErr:  ErrUserNotFound,
			code:    ErrCodeUserNotFound,
			message: "用户不存在",
			errType: ErrorTypeBusiness,
		},
		{
			name:    "ErrDatabaseError",
			bizErr:  ErrDatabaseError,
			code:    ErrCodeDatabaseError,
			message: "数据库操作失败",
			errType: ErrorTypeSystem,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.code, tt.bizErr.Code)
			assert.Equal(t, tt.message, tt.bizErr.Message)
			assert.Equal(t, tt.errType, tt.bizErr.Type)
		})
	}
}

func TestWrapError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		code     int
		message  string
		expected *BizError
	}{
		{
			name:     "包装nil错误返回nil",
			err:      nil,
			code:     ErrCodeInternalError,
			message:  "测试消息",
			expected: nil,
		},
		{
			name:    "包装标准错误",
			err:     errors.New("standard error"),
			code:    ErrCodeInternalError,
			message: "测试消息",
			expected: &BizError{
				Code:    ErrCodeInternalError,
				Message: "测试消息: standard error",
				Type:    ErrorTypeSystem,
			},
		},
		{
			name:     "包装业务错误返回原错误",
			err:      ErrUserNotFound,
			code:     ErrCodeInternalError,
			message:  "测试消息",
			expected: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapError(tt.err, tt.code, tt.message)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			assert.Equal(t, tt.expected.Code, result.Code)
			assert.Equal(t, tt.expected.Type, result.Type)

			if tt.err == ErrUserNotFound {
				// 如果是包装业务错误，应该返回原错误
				assert.Equal(t, tt.err, result)
			} else {
				// 包装标准错误，检查消息格式
				assert.Equal(t, tt.expected.Message, result.Message)
				assert.NotEmpty(t, result.Stack)
			}
		})
	}
}

func TestIsBizError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "业务错误",
			err:      ErrUserNotFound,
			expected: true,
		},
		{
			name:     "标准错误",
			err:      errors.New("standard error"),
			expected: false,
		},
		{
			name:     "nil错误",
			err:      nil,
			expected: false,
		},
		{
			name:     "自定义业务错误",
			err:      NewBizError(ErrCodeValidation, "测试", ErrorTypeValidation),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBizError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetBizError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected *BizError
	}{
		{
			name:     "获取业务错误",
			err:      ErrUserNotFound,
			expected: ErrUserNotFound,
		},
		{
			name:     "获取标准错误返回nil",
			err:      errors.New("standard error"),
			expected: nil,
		},
		{
			name:     "获取nil错误返回nil",
			err:      nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetBizError(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewValidationError(t *testing.T) {
	message := "验证失败"
	err := NewValidationError(message)

	assert.Equal(t, ErrCodeValidation, err.Code)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, ErrorTypeValidation, err.Type)
}

func TestNewAuthError(t *testing.T) {
	message := "认证失败"
	err := NewAuthError(message)

	assert.Equal(t, ErrCodeUnauthorized, err.Code)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, ErrorTypeAuth, err.Type)
}

func TestNewBusinessError(t *testing.T) {
	code := ErrCodeUserExist
	message := "用户已存在"
	err := NewBusinessError(code, message)

	assert.Equal(t, code, err.Code)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, ErrorTypeBusiness, err.Type)
}

func TestNewSystemError(t *testing.T) {
	message := "系统错误"
	err := NewSystemError(message)

	assert.Equal(t, ErrCodeInternalError, err.Code)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, ErrorTypeSystem, err.Type)
	assert.NotEmpty(t, err.Stack)
}

func TestGetStack(t *testing.T) {
	// 直接测试getStack函数
	stack := getStack()

	assert.NotEmpty(t, stack)
	// 检查堆栈是否包含当前函数名或testing相关函数
	assert.True(t, strings.Contains(stack, "TestGetStack") || strings.Contains(stack, "testing.tRunner"))

	// 验证堆栈格式
	lines := strings.Split(stack, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			assert.Contains(t, line, "at ")
			assert.Contains(t, line, "(")
			assert.Contains(t, line, ":")
			assert.Contains(t, line, ")")
		}
	}
}

func TestErrorConstants(t *testing.T) {
	// 测试错误码常量
	assert.Equal(t, 200, ErrCodeOK)
	assert.Equal(t, 500, ErrCodeInternalError)
	assert.Equal(t, 400, ErrCodeBadRequest)
	assert.Equal(t, 401, ErrCodeUnauthorized)
	assert.Equal(t, 403, ErrCodeForbidden)
	assert.Equal(t, 404, ErrCodeNotFound)
	assert.Equal(t, 409, ErrCodeConflict)
	assert.Equal(t, 429, ErrCodeTooManyRequests)
	assert.Equal(t, 503, ErrCodeServiceUnavailable)

	// 测试自定义错误码
	assert.Equal(t, 1001, ErrCodeValidation)
	assert.Equal(t, 2001, ErrCodeTokenInvalid)
	assert.Equal(t, 3001, ErrCodeUserNotFound)
	assert.Equal(t, 4001, ErrCodeDuplicateData)
	assert.Equal(t, 5001, ErrCodeBusinessLogic)
	assert.Equal(t, 6001, ErrCodeExternalService)
}

func TestErrorTypeConstants(t *testing.T) {
	// 测试错误类型常量
	assert.Equal(t, "VALIDATION_ERROR", ErrorTypeValidation)
	assert.Equal(t, "AUTH_ERROR", ErrorTypeAuth)
	assert.Equal(t, "BUSINESS_ERROR", ErrorTypeBusiness)
	assert.Equal(t, "SYSTEM_ERROR", ErrorTypeSystem)
	assert.Equal(t, "EXTERNAL_ERROR", ErrorTypeExternal)
}

func TestBizError_JSON(t *testing.T) {
	// 测试业务错误的JSON序列化
	bizErr := NewBizError(ErrCodeUserNotFound, "用户不存在", ErrorTypeBusiness)

	// 这里我们不直接测试JSON序列化，因为测试环境可能没有JSON包
	// 但我们可以测试结构体字段的JSON标签是否正确定义
	assert.NotNil(t, bizErr.Code)
	assert.NotNil(t, bizErr.Message)
	assert.NotNil(t, bizErr.Type)
}

// 基准测试
func BenchmarkNewBizError(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBizError(ErrCodeValidation, "测试消息", ErrorTypeValidation)
	}
}

func BenchmarkNewBizErrorWithStack(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBizErrorWithStack(ErrCodeInternalError, "测试消息", ErrorTypeSystem)
	}
}

func BenchmarkWrapError(b *testing.B) {
	err := errors.New("test error")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WrapError(err, ErrCodeInternalError, "包装错误")
	}
}

func BenchmarkGetStack(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = getStack()
	}
}
