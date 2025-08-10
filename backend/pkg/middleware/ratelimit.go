package middleware

import (
	"net"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

// 简单基于 Redis 的固定窗口限流中间件
func NewRateLimitMiddleware(rdb *redis.Client, requests int, window time.Duration) func(http.HandlerFunc) http.HandlerFunc {
	if rdb == nil || requests <= 0 || window <= 0 {
		// 不启用
		return func(next http.HandlerFunc) http.HandlerFunc { return next }
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			if ip == "" {
				ip = "unknown"
			}
			key := "ratelimit:" + ip + ":" + r.URL.Path
			// 增加计数并设置窗口
			pipe := rdb.TxPipeline()
			incr := pipe.Incr(r.Context(), key)
			pipe.Expire(r.Context(), key, window)
			_, _ = pipe.Exec(r.Context())
			if incr.Val() > int64(requests) {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("rate limit exceeded"))
				return
			}
			next(w, r)
		}
	}
}
