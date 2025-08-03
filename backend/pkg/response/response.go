package response

import (
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// CommonResponse 统一响应结构
type CommonResponse struct {
	Code      int         `json:"code"`                // 响应状态码
	Message   string      `json:"message"`             // 响应消息
	Data      interface{} `json:"data,omitempty"`      // 响应数据
	Timestamp int64       `json:"timestamp"`           // 响应时间戳
	RequestID string      `json:"requestId,omitempty"` // 请求唯一标识
}

// ListResponse 列表响应结构
type ListResponse struct {
	Items    interface{} `json:"items"`    // 数据列表
	Total    int64       `json:"total"`    // 总记录数
	Page     int         `json:"page"`     // 当前页码
	PageSize int         `json:"pageSize"` // 每页数量
}

// PagedResponse 分页响应结构
type PagedResponse struct {
	Code      int          `json:"code"`
	Message   string       `json:"message"`
	Data      ListResponse `json:"data"`
	Timestamp int64        `json:"timestamp"`
	RequestID string       `json:"requestId,omitempty"`
}

// 状态码常量
const (
	StatusOK                  = 200 // 成功
	StatusCreated             = 201 // 创建成功
	StatusBadRequest          = 400 // 请求错误
	StatusUnauthorized        = 401 // 未认证
	StatusForbidden           = 403 // 权限不足
	StatusNotFound            = 404 // 资源不存在
	StatusConflict            = 409 // 冲突
	StatusTooManyRequests     = 429 // 请求过多
	StatusInternalServerError = 500 // 服务器错误
	StatusServiceUnavailable  = 503 // 服务不可用
)

// Success 成功响应
func Success(w http.ResponseWriter, data interface{}) {
	response := CommonResponse{
		Code:      StatusOK,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	}

	httpx.OkJson(w, response)
}

// Created 创建成功响应
func Created(w http.ResponseWriter, data interface{}) {
	response := CommonResponse{
		Code:      StatusCreated,
		Message:   "创建成功",
		Data:      data,
		Timestamp: time.Now().Unix(),
	}

	w.WriteHeader(http.StatusCreated)
	httpx.OkJson(w, response)
}

// BadRequest 请求错误响应
func BadRequest(w http.ResponseWriter, message string) {
	response := CommonResponse{
		Code:      StatusBadRequest,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}

	w.WriteHeader(http.StatusBadRequest)
	httpx.OkJson(w, response)
}

// Unauthorized 未认证响应
func Unauthorized(w http.ResponseWriter, message string) {
	if message == "" {
		message = "未认证或token已过期"
	}

	response := CommonResponse{
		Code:      StatusUnauthorized,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}

	w.WriteHeader(http.StatusUnauthorized)
	httpx.OkJson(w, response)
}

// Forbidden 权限不足响应
func Forbidden(w http.ResponseWriter, message string) {
	if message == "" {
		message = "权限不足"
	}

	response := CommonResponse{
		Code:      StatusForbidden,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}

	w.WriteHeader(http.StatusForbidden)
	httpx.OkJson(w, response)
}

// NotFound 资源不存在响应
func NotFound(w http.ResponseWriter, message string) {
	if message == "" {
		message = "资源不存在"
	}

	response := CommonResponse{
		Code:      StatusNotFound,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}

	w.WriteHeader(http.StatusNotFound)
	httpx.OkJson(w, response)
}

// Conflict 冲突响应
func Conflict(w http.ResponseWriter, message string) {
	response := CommonResponse{
		Code:      StatusConflict,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}

	w.WriteHeader(http.StatusConflict)
	httpx.OkJson(w, response)
}

// TooManyRequests 请求过多响应
func TooManyRequests(w http.ResponseWriter, message string) {
	if message == "" {
		message = "请求过于频繁，请稍后再试"
	}

	response := CommonResponse{
		Code:      StatusTooManyRequests,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}

	w.WriteHeader(http.StatusTooManyRequests)
	httpx.OkJson(w, response)
}

// InternalServerError 服务器错误响应
func InternalServerError(w http.ResponseWriter, message string) {
	if message == "" {
		message = "服务器内部错误"
	}

	response := CommonResponse{
		Code:      StatusInternalServerError,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}

	w.WriteHeader(http.StatusInternalServerError)
	httpx.OkJson(w, response)
}

// ServiceUnavailable 服务不可用响应
func ServiceUnavailable(w http.ResponseWriter, message string) {
	if message == "" {
		message = "服务暂时不可用"
	}

	response := CommonResponse{
		Code:      StatusServiceUnavailable,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}

	w.WriteHeader(http.StatusServiceUnavailable)
	httpx.OkJson(w, response)
}

// PagedSuccess 分页成功响应
func PagedSuccess(w http.ResponseWriter, items interface{}, total int64, page, pageSize int) {
	listData := ListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	response := PagedResponse{
		Code:      StatusOK,
		Message:   "success",
		Data:      listData,
		Timestamp: time.Now().Unix(),
	}

	httpx.OkJson(w, response)
}

// CustomResponse 自定义响应
func CustomResponse(w http.ResponseWriter, httpStatus, code int, message string, data interface{}) {
	response := CommonResponse{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}

	w.WriteHeader(httpStatus)
	httpx.OkJson(w, response)
}

// WithRequestID 添加请求ID到响应中
func WithRequestID(response *CommonResponse, requestID string) {
	response.RequestID = requestID
}
