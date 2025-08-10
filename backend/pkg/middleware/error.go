package middleware

import (
	"log"
	"net/http"
	"time"

	"api/pkg/errors"
	"api/pkg/response"
)

// ErrorHandlerMiddleware 统一错误处理中间件
func ErrorHandlerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 使用recover捕获panic
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)

				// 记录panic详细信息
				log.Printf("Request: %s %s", r.Method, r.URL.Path)
				log.Printf("Remote Addr: %s", r.RemoteAddr)
				log.Printf("User Agent: %s", r.UserAgent())

				// 返回500错误
				response.InternalServerError(w, "服务器内部错误")
			}
		}()

		// 继续处理请求
		next(w, r)
	}
}

// ErrorResponseWriter 错误响应写入器
type ErrorResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// NewErrorResponseWriter 创建错误响应写入器
func NewErrorResponseWriter(w http.ResponseWriter) *ErrorResponseWriter {
	return &ErrorResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		written:        false,
	}
}

// WriteHeader 重写WriteHeader方法
func (w *ErrorResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

// Write 重写Write方法
func (w *ErrorResponseWriter) Write(data []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(data)
}

// RequestLogMiddleware 请求日志中间件
func RequestLogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 包装ResponseWriter以获取状态码
		wrappedWriter := NewErrorResponseWriter(w)

		// 生成请求ID
		requestID := generateRequestID()

		// 将请求ID添加到响应头
		wrappedWriter.Header().Set("X-Request-ID", requestID)

		// 将请求ID写入上下文
		ctx := WithValue(r.Context(), CtxKeyRequestID, requestID)
		r = r.WithContext(ctx)

		// 记录请求开始
		log.Printf("[%s] %s %s - Start", requestID, r.Method, r.URL.Path)

		// 处理请求
		next(wrappedWriter, r)

		// 计算处理时间
		duration := time.Since(start)

		// 记录请求结束
		log.Printf("[%s] %s %s - %d - %v",
			requestID, r.Method, r.URL.Path, wrappedWriter.statusCode, duration)

		// 如果是错误状态码，记录更多信息
		if wrappedWriter.statusCode >= 400 {
			log.Printf("[%s] Error Response - Status: %d, RemoteAddr: %s, UserAgent: %s",
				requestID, wrappedWriter.statusCode, r.RemoteAddr, r.UserAgent())
		}
	}
}

// HandleBizError 处理业务错误
func HandleBizError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	// 判断是否为业务错误
	if bizErr := errors.GetBizError(err); bizErr != nil {
		// 记录业务错误
		log.Printf("Business Error: %s", bizErr.Error())

		// 获取对应的HTTP状态码
		httpStatus := bizErr.GetHTTPStatus()

		// 根据错误类型返回不同的响应
		switch bizErr.Code {
		case errors.ErrCodeBadRequest, errors.ErrCodeValidation, errors.ErrCodeInvalidParam:
			response.BadRequest(w, bizErr.Message)
		case errors.ErrCodeUnauthorized, errors.ErrCodeTokenInvalid, errors.ErrCodeTokenExpired:
			response.Unauthorized(w, bizErr.Message)
		case errors.ErrCodeForbidden, errors.ErrCodePermissionDenied:
			response.Forbidden(w, bizErr.Message)
		case errors.ErrCodeNotFound, errors.ErrCodeUserNotFound:
			response.NotFound(w, bizErr.Message)
		case errors.ErrCodeConflict, errors.ErrCodeDuplicateData:
			response.Conflict(w, bizErr.Message)
		case errors.ErrCodeTooManyRequests:
			response.TooManyRequests(w, bizErr.Message)
		case errors.ErrCodeServiceUnavailable:
			response.ServiceUnavailable(w, bizErr.Message)
		default:
			response.CustomResponse(w, httpStatus, bizErr.Code, bizErr.Message, nil)
		}
		return
	}

	// 非业务错误，记录详细错误信息
	log.Printf("System Error: %v", err)
	response.InternalServerError(w, "系统错误")
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return generateUUID()
}

// generateUUID 生成UUID
func generateUUID() string {
	// 简单的UUID生成，实际生产环境建议使用专业的UUID库
	return time.Now().Format("20060102150405") + "-" +
		string(rune(time.Now().Nanosecond()%26+65)) +
		string(rune(time.Now().Nanosecond()%26+65))
}
