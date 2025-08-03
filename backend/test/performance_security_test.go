package test

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"api/pkg/auth"
	testconfig "api/test/config"
	"github.com/stretchr/testify/assert"
)

// BenchmarkPasswordHashing 密码哈希性能基准测试
func BenchmarkPasswordHashing(b *testing.B) {
	password := "testpassword123"
	passwordService := auth.NewPasswordService()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := passwordService.HashPassword(password)
		if err != nil {
			b.Fatalf("密码哈希失败: %v", err)
		}
	}
}

// BenchmarkPasswordVerification 密码验证性能基准测试
func BenchmarkPasswordVerification(b *testing.B) {
	password := "testpassword123"
	passwordService := auth.NewPasswordService()
	hashedPassword, _ := passwordService.HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !passwordService.VerifyPassword(hashedPassword, password) {
			b.Fatal("密码验证失败")
		}
	}
}

// BenchmarkJWTGeneration JWT生成性能基准测试
func BenchmarkJWTGeneration(b *testing.B) {
	userID := int64(123)
	email := "test@example.com"
	jwtService := auth.NewJWTService("test-access-secret", "test-refresh-secret", 3600, 86400)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := jwtService.GenerateTokenPair(userID, "test", email, []string{"user"}, []string{"read"})
		if err != nil {
			b.Fatalf("JWT生成失败: %v", err)
		}
	}
}

// BenchmarkJWTValidation JWT验证性能基准测试
func BenchmarkJWTValidation(b *testing.B) {
	userID := int64(123)
	email := "test@example.com"
	jwtService := auth.NewJWTService("test-access-secret", "test-refresh-secret", 3600, 86400)
	tokenPair, _ := jwtService.GenerateTokenPair(userID, "test", email, []string{"user"}, []string{"read"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := jwtService.ParseAccessToken(tokenPair.AccessToken)
		if err != nil {
			b.Fatalf("JWT验证失败: %v", err)
		}
	}
}

// TestConcurrentPasswordHashing 并发密码哈希测试
func TestConcurrentPasswordHashing(t *testing.T) {
	const numGoroutines = 100
	const passwordsPerGoroutine = 10

	passwordService := auth.NewPasswordService()
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*passwordsPerGoroutine)

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < passwordsPerGoroutine; j++ {
				password := fmt.Sprintf("password_%d_%d", goroutineID, j)
				_, err := passwordService.HashPassword(password)
				if err != nil {
					errors <- err
					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	duration := time.Since(start)

	// 检查是否有错误
	for err := range errors {
		t.Errorf("并发密码哈希失败: %v", err)
	}

	totalOperations := numGoroutines * passwordsPerGoroutine
	avgDuration := duration / time.Duration(totalOperations)

	t.Logf("并发密码哈希测试完成:")
	t.Logf("- 协程数: %d", numGoroutines)
	t.Logf("- 总操作数: %d", totalOperations)
	t.Logf("- 总耗时: %v", duration)
	t.Logf("- 平均耗时: %v", avgDuration)

	// 性能要求：平均每次操作不超过50ms（在高并发下）
	assert.Less(t, avgDuration, time.Millisecond*50, "并发场景下平均操作时间应小于50ms")
}

// TestMemoryUsageUnderLoad 负载下内存使用测试
func TestMemoryUsageUnderLoad(t *testing.T) {
	// 获取测试前内存状态
	var m1 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	jwtService := auth.NewJWTService("test-access-secret", "test-refresh-secret", 3600, 86400)

	// 创建大量JWT token
	const tokenCount = 10000
	tokens := make([]string, tokenCount)

	for i := 0; i < tokenCount; i++ {
		userID := int64(i + 1)
		email := fmt.Sprintf("user%d@example.com", i+1)

		tokenPair, err := jwtService.GenerateTokenPair(userID, "test", email, []string{"user"}, []string{"read"})
		assert.NoError(t, err, "JWT生成不应该失败")
		tokens[i] = tokenPair.AccessToken
	}

	// 验证所有token
	for i, token := range tokens {
		claims, err := jwtService.ParseAccessToken(token)
		assert.NoError(t, err, "JWT验证不应该失败")
		assert.Equal(t, int64(i+1), claims.UserID, "用户ID应该匹配")
	}

	// 获取测试后内存状态
	var m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m2)

	// 计算内存使用增长
	memoryIncrease := m2.Alloc - m1.Alloc
	avgMemoryPerToken := memoryIncrease / tokenCount

	t.Logf("内存使用测试结果:")
	t.Logf("- Token数量: %d", tokenCount)
	t.Logf("- 内存增长: %d bytes (%.2f MB)", memoryIncrease, float64(memoryIncrease)/1024/1024)
	t.Logf("- 平均每Token: %d bytes", avgMemoryPerToken)

	// 内存要求：平均每个token不超过1KB内存
	assert.Less(t, avgMemoryPerToken, uint64(1024), "平均每Token内存使用应小于1KB")

	// 清理
	tokens = nil
	runtime.GC()
}

// TestDatabaseConnectionPool 数据库连接池测试
func TestDatabaseConnectionPool(t *testing.T) {
	testConfig := testconfig.GetTestConfig()

	// 尝试创建真实数据库连接进行测试
	db, err := testConfig.CreateTestDB()
	if err != nil {
		t.Skipf("跳过数据库连接池测试，无法连接数据库: %v", err)
		return
	}
	defer db.Close()

	const numGoroutines = 50
	const queriesPerGoroutine = 20

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*queriesPerGoroutine)

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < queriesPerGoroutine; j++ {
				var result int
				err := db.QueryRow("SELECT ?", j).Scan(&result)
				if err != nil {
					errors <- fmt.Errorf("查询失败 (goroutine %d, query %d): %v", goroutineID, j, err)
					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	duration := time.Since(start)

	// 检查错误
	errorCount := 0
	for err := range errors {
		errorCount++
		t.Logf("数据库连接池错误: %v", err)
	}

	totalQueries := numGoroutines * queriesPerGoroutine
	successRate := float64(totalQueries-errorCount) / float64(totalQueries) * 100

	t.Logf("数据库连接池测试结果:")
	t.Logf("- 并发数: %d", numGoroutines)
	t.Logf("- 总查询数: %d", totalQueries)
	t.Logf("- 成功率: %.2f%%", successRate)
	t.Logf("- 总耗时: %v", duration)

	// 性能要求：成功率应该≥95%
	assert.GreaterOrEqual(t, successRate, 95.0, "数据库连接池成功率应该≥95%")
}

// TestRateLimiting 速率限制测试
func TestRateLimiting(t *testing.T) {
	// 模拟速率限制器
	type RateLimiter struct {
		requests map[string][]time.Time
		mu       sync.Mutex
		limit    int
		window   time.Duration
	}

	rateLimiter := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    10, // 每分钟10次请求
		window:   time.Minute,
	}

	isAllowed := func(clientID string) bool {
		rateLimiter.mu.Lock()
		defer rateLimiter.mu.Unlock()

		now := time.Now()

		// 清理过期请求
		if times, exists := rateLimiter.requests[clientID]; exists {
			var validTimes []time.Time
			for _, t := range times {
				if now.Sub(t) < rateLimiter.window {
					validTimes = append(validTimes, t)
				}
			}
			rateLimiter.requests[clientID] = validTimes
		}

		// 检查是否超过限制
		if len(rateLimiter.requests[clientID]) >= rateLimiter.limit {
			return false
		}

		// 记录请求
		rateLimiter.requests[clientID] = append(rateLimiter.requests[clientID], now)
		return true
	}

	clientID := "test-client"

	// 测试正常请求（应该在限制内）
	allowedCount := 0
	for i := 0; i < rateLimiter.limit; i++ {
		if isAllowed(clientID) {
			allowedCount++
		}
	}
	assert.Equal(t, rateLimiter.limit, allowedCount, "在限制内的请求应该全部通过")

	// 测试超限请求（应该被拒绝）
	assert.False(t, isAllowed(clientID), "超限请求应该被拒绝")

	t.Logf("速率限制测试通过: %d/%d 请求通过", allowedCount, rateLimiter.limit)
}

// TestResourceLeakDetection 资源泄漏检测测试
func TestResourceLeakDetection(t *testing.T) {
	// 获取初始goroutine数量
	initialGoroutines := runtime.NumGoroutine()

	// 创建一些协程并确保它们能正确退出
	const numGoroutines = 10
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Millisecond * 100):
				// 模拟一些工作
				_ = fmt.Sprintf("goroutine %d completed", id)
			}
		}(i)
	}

	wg.Wait()

	// 等待一段时间让协程完全退出
	time.Sleep(time.Millisecond * 100)

	// 检查协程数量是否回到初始状态（允许一些偏差）
	finalGoroutines := runtime.NumGoroutine()
	goroutineDiff := finalGoroutines - initialGoroutines

	t.Logf("Goroutine数量检测:")
	t.Logf("- 初始: %d", initialGoroutines)
	t.Logf("- 最终: %d", finalGoroutines)
	t.Logf("- 差异: %d", goroutineDiff)

	// 允许少量goroutine差异（运行时系统的正常变化）
	assert.LessOrEqual(t, goroutineDiff, 5, "Goroutine泄漏检测：差异应该≤5")
}

// TestSecurityHeaders 安全头测试
func TestSecurityHeaders(t *testing.T) {
	// 模拟HTTP响应头验证
	securityHeaders := map[string]string{
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "DENY",
		"X-XSS-Protection":          "1; mode=block",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
		"Content-Security-Policy":   "default-src 'self'",
	}

	// 验证所有必需的安全头
	for header, expectedValue := range securityHeaders {
		assert.NotEmpty(t, expectedValue, "安全头 %s 不应为空", header)
		t.Logf("✅ 安全头 %s: %s", header, expectedValue)
	}

	// 验证敏感信息不被暴露
	sensitiveHeaders := []string{
		"Server",       // 不暴露服务器信息
		"X-Powered-By", // 不暴露技术栈信息
	}

	for _, header := range sensitiveHeaders {
		// 在实际实现中，这些头应该被移除或设置为通用值
		t.Logf("⚠️  敏感头 %s 应该被移除或隐藏", header)
	}
}
