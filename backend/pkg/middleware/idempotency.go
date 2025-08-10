package middleware

import (
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewIdempotencyMiddleware 基于 Redis 的简单幂等性中间件
// 客户端在写操作（POST/PUT/DELETE）提供 Idempotency-Key 头部值，避免短期重复提交
func NewIdempotencyMiddleware(rdb *redis.Client, ttl time.Duration) func(http.HandlerFunc) http.HandlerFunc {
	if rdb == nil || ttl <= 0 {
		return func(next http.HandlerFunc) http.HandlerFunc { return next }
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodDelete {
				next(w, r)
				return
			}

			key := r.Header.Get("Idempotency-Key")
			if key == "" {
				// 未提供幂等键则直接放行（可按需改为强制）
				next(w, r)
				return
			}

			redisKey := "idem:" + key
			ok, err := rdb.SetNX(r.Context(), redisKey, "1", ttl).Result()
			if err == nil && ok {
				next(w, r)
				return
			}
			// 重复请求或 Redis 异常都返回冲突，避免重复写入
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("duplicate request"))
		}
	}
}
