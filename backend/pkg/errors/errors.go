package errors

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// BizError 业务错误
type BizError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
	Stack   string `json:"stack,omitempty"`
}

// Error 实现error接口
func (e *BizError) Error() string {
	return fmt.Sprintf("[%s] %s (code: %d)", e.Type, e.Message, e.Code)
}

// GetHTTPStatus 获取对应的HTTP状态码
func (e *BizError) GetHTTPStatus() int {
	switch e.Code {
	case ErrCodeBadRequest, ErrCodeValidation, ErrCodeInvalidParam:
		return http.StatusBadRequest
	case ErrCodeUnauthorized, ErrCodeTokenInvalid, ErrCodeTokenExpired:
		return http.StatusUnauthorized
	case ErrCodeForbidden, ErrCodePermissionDenied:
		return http.StatusForbidden
	case ErrCodeNotFound, ErrCodeUserNotFound:
		return http.StatusNotFound
	case ErrCodeConflict, ErrCodeDuplicateData:
		return http.StatusConflict
	case ErrCodeTooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// 错误码定义
const (
	// 通用错误码
	ErrCodeOK                 = 200
	ErrCodeInternalError      = 500
	ErrCodeBadRequest         = 400
	ErrCodeUnauthorized       = 401
	ErrCodeForbidden          = 403
	ErrCodeNotFound           = 404
	ErrCodeConflict           = 409
	ErrCodeTooManyRequests    = 429
	ErrCodeServiceUnavailable = 503

	// 参数验证错误码 (1000-1099)
	ErrCodeValidation    = 1001
	ErrCodeInvalidParam  = 1002
	ErrCodeMissingParam  = 1003
	ErrCodeInvalidFormat = 1004

	// 认证授权错误码 (2000-2099)
	ErrCodeTokenInvalid     = 2001
	ErrCodeTokenExpired     = 2002
	ErrCodePermissionDenied = 2003
	ErrCodeLoginFailed      = 2004
	ErrCodeAccountDisabled  = 2005

	// 用户相关错误码 (3000-3099)
	ErrCodeUserNotFound  = 3001
	ErrCodeUserExist     = 3002
	ErrCodePasswordWrong = 3003
	ErrCodeEmailExist    = 3004
	ErrCodeUsernameExist = 3005

	// 数据相关错误码 (4000-4099)
	ErrCodeDuplicateData = 4001
	ErrCodeDataNotFound  = 4002
	ErrCodeDataInvalid   = 4003
	ErrCodeDatabaseError = 4004

	// 业务逻辑错误码 (5000-5099)
	ErrCodeBusinessLogic   = 5001
	ErrCodeOperationFailed = 5002
	ErrCodeResourceBusy    = 5003
	ErrCodeQuotaExceeded   = 5004

	// 外部服务错误码 (6000-6099)
	ErrCodeExternalService = 6001
	ErrCodeNetworkError    = 6002
	ErrCodeTimeout         = 6003
)

// 错误类型定义
const (
	ErrorTypeValidation = "VALIDATION_ERROR"
	ErrorTypeAuth       = "AUTH_ERROR"
	ErrorTypeBusiness   = "BUSINESS_ERROR"
	ErrorTypeSystem     = "SYSTEM_ERROR"
	ErrorTypeExternal   = "EXTERNAL_ERROR"
)

// NewBizError 创建业务错误
func NewBizError(code int, message, errorType string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
		Type:    errorType,
	}
}

// NewBizErrorWithStack 创建带堆栈信息的业务错误
func NewBizErrorWithStack(code int, message, errorType string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
		Type:    errorType,
		Stack:   getStack(),
	}
}

// 预定义错误
var (
	// 参数验证错误
	ErrInvalidParam     = NewBizError(ErrCodeInvalidParam, "参数无效", ErrorTypeValidation)
	ErrMissingParam     = NewBizError(ErrCodeMissingParam, "缺少必要参数", ErrorTypeValidation)
	ErrValidationFailed = NewBizError(ErrCodeValidation, "参数验证失败", ErrorTypeValidation)
	ErrInvalidFormat    = NewBizError(ErrCodeInvalidFormat, "格式错误", ErrorTypeValidation)

	// 认证授权错误
	ErrTokenInvalid     = NewBizError(ErrCodeTokenInvalid, "token无效", ErrorTypeAuth)
	ErrTokenExpired     = NewBizError(ErrCodeTokenExpired, "token已过期", ErrorTypeAuth)
	ErrPermissionDenied = NewBizError(ErrCodePermissionDenied, "权限不足", ErrorTypeAuth)
	ErrLoginFailed      = NewBizError(ErrCodeLoginFailed, "登录失败", ErrorTypeAuth)
	ErrAccountDisabled  = NewBizError(ErrCodeAccountDisabled, "账户已禁用", ErrorTypeAuth)

	// 用户相关错误
	ErrUserNotFound  = NewBizError(ErrCodeUserNotFound, "用户不存在", ErrorTypeBusiness)
	ErrUserExist     = NewBizError(ErrCodeUserExist, "用户已存在", ErrorTypeBusiness)
	ErrPasswordWrong = NewBizError(ErrCodePasswordWrong, "密码错误", ErrorTypeBusiness)
	ErrEmailExist    = NewBizError(ErrCodeEmailExist, "邮箱已存在", ErrorTypeBusiness)
	ErrUsernameExist = NewBizError(ErrCodeUsernameExist, "用户名已存在", ErrorTypeBusiness)

	// 数据相关错误
	ErrDuplicateData = NewBizError(ErrCodeDuplicateData, "数据重复", ErrorTypeBusiness)
	ErrDataNotFound  = NewBizError(ErrCodeDataNotFound, "数据不存在", ErrorTypeBusiness)
	ErrDataInvalid   = NewBizError(ErrCodeDataInvalid, "数据无效", ErrorTypeBusiness)
	ErrDatabaseError = NewBizError(ErrCodeDatabaseError, "数据库操作失败", ErrorTypeSystem)

	// 业务逻辑错误
	ErrBusinessLogic   = NewBizError(ErrCodeBusinessLogic, "业务逻辑错误", ErrorTypeBusiness)
	ErrOperationFailed = NewBizError(ErrCodeOperationFailed, "操作失败", ErrorTypeBusiness)
	ErrResourceBusy    = NewBizError(ErrCodeResourceBusy, "资源忙碌", ErrorTypeBusiness)
	ErrQuotaExceeded   = NewBizError(ErrCodeQuotaExceeded, "配额已超限", ErrorTypeBusiness)

	// 外部服务错误
	ErrExternalService = NewBizError(ErrCodeExternalService, "外部服务错误", ErrorTypeExternal)
	ErrNetworkError    = NewBizError(ErrCodeNetworkError, "网络错误", ErrorTypeExternal)
	ErrTimeout         = NewBizError(ErrCodeTimeout, "请求超时", ErrorTypeExternal)
)

// WrapError 包装标准错误为业务错误
func WrapError(err error, code int, message string) *BizError {
	if err == nil {
		return nil
	}

	if bizErr, ok := err.(*BizError); ok {
		return bizErr
	}

	return &BizError{
		Code:    code,
		Message: fmt.Sprintf("%s: %s", message, err.Error()),
		Type:    ErrorTypeSystem,
		Stack:   getStack(),
	}
}

// IsBizError 判断是否为业务错误
func IsBizError(err error) bool {
	_, ok := err.(*BizError)
	return ok
}

// GetBizError 获取业务错误
func GetBizError(err error) *BizError {
	if bizErr, ok := err.(*BizError); ok {
		return bizErr
	}
	return nil
}

// NewValidationError 创建参数验证错误
func NewValidationError(message string) *BizError {
	return NewBizError(ErrCodeValidation, message, ErrorTypeValidation)
}

// NewAuthError 创建认证错误
func NewAuthError(message string) *BizError {
	return NewBizError(ErrCodeUnauthorized, message, ErrorTypeAuth)
}

// NewBusinessError 创建业务错误
func NewBusinessError(code int, message string) *BizError {
	return NewBizError(code, message, ErrorTypeBusiness)
}

// NewInternalError 创建内部错误
func NewInternalError(message string) *BizError {
	return NewBizError(ErrCodeInternalError, message, ErrorTypeSystem)
}

// 快速错误定义
var (
	ErrInvalidPassword = NewBizError(ErrCodePasswordWrong, "密码错误", ErrorTypeAuth)
	ErrUserDisabled    = NewBizError(ErrCodeAccountDisabled, "用户已被禁用", ErrorTypeAuth)
	ErrInvalidToken    = NewBizError(ErrCodeTokenInvalid, "令牌无效", ErrorTypeAuth)
)

// getStack 获取调用堆栈
func getStack() string {
	var buf strings.Builder

	// 跳过前面几层调用栈
	for i := 2; i < 10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}

		// 过滤掉runtime相关的调用
		if strings.Contains(fn.Name(), "runtime.") {
			continue
		}

		// 格式化输出
		buf.WriteString(fmt.Sprintf("  at %s (%s:%d)\n", fn.Name(), file, line))
	}

	return buf.String()
}
