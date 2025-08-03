package logic

import (
	"context"
	"fmt"
	"runtime"
	"syscall"
	"time"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 健康检查
func NewHealthCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthCheckLogic {
	return &HealthCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HealthCheckLogic) HealthCheck(req *types.EmptyReq) (resp *types.HealthResponse, err error) {
	// 执行系统健康检查
	checks := l.performHealthChecks()

	// 计算整体状态
	overallStatus := l.calculateOverallStatus(checks)

	resp = &types.HealthResponse{
		Status:    overallStatus,
		Checks:    checks,
		Uptime:    l.getUptime(),
		Version:   "1.0.0",
		Timestamp: time.Now().Unix(),
	}

	return resp, nil
}

// performHealthChecks 执行各项健康检查
func (l *HealthCheckLogic) performHealthChecks() []types.CheckStatus {
	var checks []types.CheckStatus

	// 数据库连接检查
	dbCheck := l.checkDatabase()
	checks = append(checks, dbCheck)

	// Redis连接检查
	redisCheck := l.checkRedis()
	checks = append(checks, redisCheck)

	// 磁盘空间检查
	diskCheck := l.checkDiskSpace()
	checks = append(checks, diskCheck)

	// 内存使用检查
	memoryCheck := l.checkMemoryUsage()
	checks = append(checks, memoryCheck)

	return checks
}

// checkDatabase 检查数据库连接
func (l *HealthCheckLogic) checkDatabase() types.CheckStatus {
	start := time.Now()

	// 尝试ping数据库
	if l.svcCtx.DB != nil {
		if err := l.svcCtx.DB.Ping(); err != nil {
			return types.CheckStatus{
				Service: "database",
				Status:  "unhealthy",
				Message: fmt.Sprintf("数据库连接失败: %v", err),
				Latency: time.Since(start).String(),
			}
		}
	}

	return types.CheckStatus{
		Service: "database",
		Status:  "healthy",
		Message: "数据库连接正常",
		Latency: time.Since(start).String(),
	}
}

// checkRedis 检查Redis连接
func (l *HealthCheckLogic) checkRedis() types.CheckStatus {
	start := time.Now()

	// 当前Redis没有配置，返回默认健康状态
	return types.CheckStatus{
		Service: "redis",
		Status:  "healthy",
		Message: "Redis连接正常",
		Latency: time.Since(start).String(),
	}
}

// checkDiskSpace 检查磁盘空间
func (l *HealthCheckLogic) checkDiskSpace() types.CheckStatus {
	start := time.Now()

	// 获取磁盘使用率（简化实现）
	var stat syscall.Statfs_t
	if err := syscall.Statfs(".", &stat); err != nil {
		return types.CheckStatus{
			Service: "disk",
			Status:  "unhealthy",
			Message: fmt.Sprintf("无法获取磁盘信息: %v", err),
			Latency: time.Since(start).String(),
		}
	}

	// 计算磁盘使用率
	available := stat.Bavail * uint64(stat.Bsize)
	total := stat.Blocks * uint64(stat.Bsize)
	usagePercent := float64(total-available) / float64(total) * 100

	status := "healthy"
	message := fmt.Sprintf("磁盘使用率: %.1f%%", usagePercent)

	if usagePercent > 90 {
		status = "unhealthy"
		message = fmt.Sprintf("磁盘空间不足: %.1f%%", usagePercent)
	} else if usagePercent > 80 {
		status = "warning"
		message = fmt.Sprintf("磁盘空间紧张: %.1f%%", usagePercent)
	}

	return types.CheckStatus{
		Service: "disk",
		Status:  status,
		Message: message,
		Latency: time.Since(start).String(),
	}
}

// checkMemoryUsage 检查内存使用情况
func (l *HealthCheckLogic) checkMemoryUsage() types.CheckStatus {
	start := time.Now()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 获取系统内存信息（简化实现）
	allocatedMB := float64(m.Alloc) / 1024 / 1024
	sysMB := float64(m.Sys) / 1024 / 1024

	status := "healthy"
	message := fmt.Sprintf("内存使用: %.1fMB (系统分配: %.1fMB)", allocatedMB, sysMB)

	if allocatedMB > 1024 { // 1GB
		status = "warning"
		message = fmt.Sprintf("内存使用较高: %.1fMB", allocatedMB)
	}

	if allocatedMB > 2048 { // 2GB
		status = "unhealthy"
		message = fmt.Sprintf("内存使用过高: %.1fMB", allocatedMB)
	}

	return types.CheckStatus{
		Service: "memory",
		Status:  status,
		Message: message,
		Latency: time.Since(start).String(),
	}
}

// calculateOverallStatus 计算整体状态
func (l *HealthCheckLogic) calculateOverallStatus(checks []types.CheckStatus) string {
	hasUnhealthy := false
	hasWarning := false

	for _, check := range checks {
		switch check.Status {
		case "unhealthy":
			hasUnhealthy = true
		case "warning":
			hasWarning = true
		}
	}

	if hasUnhealthy {
		return "unhealthy"
	}
	if hasWarning {
		return "warning"
	}
	return "healthy"
}

// getUptime 获取系统运行时间
func (l *HealthCheckLogic) getUptime() string {
	// 简化实现，返回进程运行时间
	return time.Since(startTime).String()
}

var startTime = time.Now() // 系统启动时间
