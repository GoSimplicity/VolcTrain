package health

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// HealthService 健康检查服务
type HealthService struct {
	db    *sql.DB
	redis *redis.Client
}

// NewHealthService 创建健康检查服务实例
func NewHealthService(db *sql.DB, redis *redis.Client) *HealthService {
	return &HealthService{
		db:    db,
		redis: redis,
	}
}

// CheckStatus 健康状态
type CheckStatus struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Latency string `json:"latency"`
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string        `json:"status"`
	Checks    []CheckStatus `json:"checks"`
	Uptime    string        `json:"uptime"`
	Version   string        `json:"version"`
	Timestamp int64         `json:"timestamp"`
}

// CheckHealth 执行健康检查
func (h *HealthService) CheckHealth(ctx context.Context) *HealthResponse {
	var checks []CheckStatus
	overallStatus := "healthy"

	// 检查MySQL连接
	mysqlStatus := h.checkMySQL(ctx)
	checks = append(checks, mysqlStatus)
	if mysqlStatus.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	// 检查Redis连接
	redisStatus := h.checkRedis(ctx)
	checks = append(checks, redisStatus)
	if redisStatus.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	return &HealthResponse{
		Status:    overallStatus,
		Checks:    checks,
		Uptime:    "0s", // TODO: 实际计算运行时间
		Version:   "1.0.0",
		Timestamp: time.Now().Unix(),
	}
}

// checkMySQL 检查MySQL连接
func (h *HealthService) checkMySQL(ctx context.Context) CheckStatus {
	start := time.Now()

	// 创建带超时的context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 执行简单查询测试连接
	var result int
	err := h.db.QueryRowContext(ctx, "SELECT 1").Scan(&result)

	latency := time.Since(start)

	if err != nil {
		return CheckStatus{
			Service: "mysql",
			Status:  "unhealthy",
			Message: fmt.Sprintf("MySQL连接失败: %v", err),
			Latency: latency.String(),
		}
	}

	return CheckStatus{
		Service: "mysql",
		Status:  "healthy",
		Message: "MySQL连接正常",
		Latency: latency.String(),
	}
}

// checkRedis 检查Redis连接
func (h *HealthService) checkRedis(ctx context.Context) CheckStatus {
	start := time.Now()

	// 创建带超时的context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 执行ping命令测试连接
	_, err := h.redis.Ping(ctx).Result()

	latency := time.Since(start)

	if err != nil {
		return CheckStatus{
			Service: "redis",
			Status:  "unhealthy",
			Message: fmt.Sprintf("Redis连接失败: %v", err),
			Latency: latency.String(),
		}
	}

	return CheckStatus{
		Service: "redis",
		Status:  "healthy",
		Message: "Redis连接正常",
		Latency: latency.String(),
	}
}
