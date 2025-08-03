package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"api/internal/config"
	"github.com/redis/go-redis/v9"
)

// RedisManager Redis管理器
type RedisManager struct {
	client *redis.Client
	config config.RedisConfig
}

// NewRedisClient 创建Redis客户端连接
func NewRedisClient(c config.RedisConfig) (*redis.Client, error) {
	// 创建Redis客户端配置
	opt := &redis.Options{
		Addr:            fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password:        c.Password,
		DB:              c.DB,
		PoolSize:        c.PoolSize,
		MinIdleConns:    c.MinIdleConns,
		DialTimeout:     10 * time.Second,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		PoolTimeout:     30 * time.Second,
		MaxRetries:      3,                      // 最大重试次数
		MaxRetryBackoff: 2 * time.Second,        // 最大重试间隔
		MinRetryBackoff: 100 * time.Millisecond, // 最小重试间隔
	}

	// 创建Redis客户端
	client := redis.NewClient(opt)

	// 测试Redis连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("Redis连接测试失败: %w", err)
	}

	log.Printf("Redis连接池配置: Host=%s:%d, DB=%d, PoolSize=%d, MinIdleConns=%d",
		c.Host, c.Port, c.DB, c.PoolSize, c.MinIdleConns)

	return client, nil
}

// NewRedisManager 创建Redis管理器
func NewRedisManager(c config.RedisConfig) (*RedisManager, error) {
	client, err := NewRedisClient(c)
	if err != nil {
		return nil, err
	}

	return &RedisManager{
		client: client,
		config: c,
	}, nil
}

// GetClient 获取Redis客户端
func (r *RedisManager) GetClient() *redis.Client {
	return r.client
}

// Close 关闭Redis连接
func (r *RedisManager) Close() error {
	if r.client != nil {
		return r.client.Close()
	}
	return nil
}

// HealthCheck Redis健康检查
func (r *RedisManager) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis健康检查失败: %w", err)
	}

	return nil
}

// GetStats 获取Redis统计信息
func (r *RedisManager) GetStats(ctx context.Context) (map[string]string, error) {
	stats := r.client.PoolStats()
	info := map[string]string{
		"hits":       fmt.Sprintf("%d", stats.Hits),
		"misses":     fmt.Sprintf("%d", stats.Misses),
		"timeouts":   fmt.Sprintf("%d", stats.Timeouts),
		"totalConns": fmt.Sprintf("%d", stats.TotalConns),
		"idleConns":  fmt.Sprintf("%d", stats.IdleConns),
		"staleConns": fmt.Sprintf("%d", stats.StaleConns),
	}
	return info, nil
}

// RedisCache Redis缓存操作封装
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建Redis缓存操作实例
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

// Set 设置缓存
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取缓存
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// GetObject 获取对象缓存
func (r *RedisCache) GetObject(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// SetObject 设置对象缓存
func (r *RedisCache) SetObject(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化对象失败: %w", err)
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

// Del 删除缓存
func (r *RedisCache) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Exists 检查key是否存在
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result > 0, err
}

// Expire 设置key过期时间
func (r *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

// TTL 获取key剩余过期时间
func (r *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

// Incr 自增
func (r *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// Decr 自减
func (r *RedisCache) Decr(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, key).Result()
}

// HSet 设置哈希字段
func (r *RedisCache) HSet(ctx context.Context, key, field string, value interface{}) error {
	return r.client.HSet(ctx, key, field, value).Err()
}

// HGet 获取哈希字段
func (r *RedisCache) HGet(ctx context.Context, key, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

// HGetAll 获取所有哈希字段
func (r *RedisCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func (r *RedisCache) HDel(ctx context.Context, key string, fields ...string) error {
	return r.client.HDel(ctx, key, fields...).Err()
}

// LPush 从左侧推入列表
func (r *RedisCache) LPush(ctx context.Context, key string, values ...interface{}) error {
	return r.client.LPush(ctx, key, values...).Err()
}

// RPush 从右侧推入列表
func (r *RedisCache) RPush(ctx context.Context, key string, values ...interface{}) error {
	return r.client.RPush(ctx, key, values...).Err()
}

// LPop 从左侧弹出列表元素
func (r *RedisCache) LPop(ctx context.Context, key string) (string, error) {
	return r.client.LPop(ctx, key).Result()
}

// RPop 从右侧弹出列表元素
func (r *RedisCache) RPop(ctx context.Context, key string) (string, error) {
	return r.client.RPop(ctx, key).Result()
}

// LLen 获取列表长度
func (r *RedisCache) LLen(ctx context.Context, key string) (int64, error) {
	return r.client.LLen(ctx, key).Result()
}

// SAdd 添加到集合
func (r *RedisCache) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func (r *RedisCache) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

// SIsMember 检查是否为集合成员
func (r *RedisCache) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return r.client.SIsMember(ctx, key, member).Result()
}

// SRem 从集合移除成员
func (r *RedisCache) SRem(ctx context.Context, key string, members ...interface{}) error {
	return r.client.SRem(ctx, key, members...).Err()
}
