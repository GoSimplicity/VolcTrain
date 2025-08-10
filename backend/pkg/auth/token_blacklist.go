package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisTokenBlacklist 基于Redis的Token黑名单
type RedisTokenBlacklist struct {
	client *redis.Client
	prefix string
}

// NewRedisTokenBlacklist 创建Redis Token黑名单
func NewRedisTokenBlacklist(client *redis.Client) *RedisTokenBlacklist {
	return &RedisTokenBlacklist{
		client: client,
		prefix: "token:blacklist:",
	}
}

// AddToken 添加Token到黑名单
func (tb *RedisTokenBlacklist) AddToken(ctx context.Context, tokenString string, expireAt time.Time) error {
	if tb.client == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}

	// 生成token的唯一标识（使用SHA256哈希）
	tokenHash := tb.hashToken(tokenString)
	key := tb.prefix + tokenHash

	// 计算剩余有效期
	ttl := time.Until(expireAt)
	if ttl <= 0 {
		ttl = time.Minute // 至少保留1分钟
	}

	// 将token添加到黑名单
	err := tb.client.Set(ctx, key, "1", ttl).Err()
	if err != nil {
		return fmt.Errorf("添加Token到黑名单失败: %w", err)
	}

	return nil
}

// IsBlacklisted 检查Token是否在黑名单中
func (tb *RedisTokenBlacklist) IsBlacklisted(ctx context.Context, tokenString string) (bool, error) {
	if tb.client == nil {
		return false, nil // 如果Redis不可用，则认为不在黑名单中
	}

	tokenHash := tb.hashToken(tokenString)
	key := tb.prefix + tokenHash

	// 检查token是否在黑名单中
	exists, err := tb.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("检查Token黑名单失败: %w", err)
	}

	return exists > 0, nil
}

// RemoveToken 从黑名单中移除Token
func (tb *RedisTokenBlacklist) RemoveToken(ctx context.Context, tokenString string) error {
	if tb.client == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}

	tokenHash := tb.hashToken(tokenString)
	key := tb.prefix + tokenHash

	// 从黑名单中移除token
	err := tb.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("从黑名单移除Token失败: %w", err)
	}

	return nil
}

// Clear 清空黑名单（谨慎使用）
func (tb *RedisTokenBlacklist) Clear(ctx context.Context) error {
	if tb.client == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}

	// 获取所有黑名单key
	pattern := tb.prefix + "*"
	keys, err := tb.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("获取黑名单Keys失败: %w", err)
	}

	// 批量删除
	if len(keys) > 0 {
		err = tb.client.Del(ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("清空黑名单失败: %w", err)
		}
	}

	return nil
}

// Size 获取黑名单大小
func (tb *RedisTokenBlacklist) Size(ctx context.Context) (int64, error) {
	if tb.client == nil {
		return 0, fmt.Errorf("Redis客户端未初始化")
	}

	pattern := tb.prefix + "*"
	keys, err := tb.client.Keys(ctx, pattern).Result()
	if err != nil {
		return 0, fmt.Errorf("获取黑名单Keys失败: %w", err)
	}

	return int64(len(keys)), nil
}

// CleanupExpired 清理过期的Token（Redis会自动清理，此方法主要用于统计）
func (tb *RedisTokenBlacklist) CleanupExpired(ctx context.Context) error {
	// Redis会自动清理过期的key，这里主要用于统计和日志记录
	size, err := tb.Size(ctx)
	if err != nil {
		return err
	}

	// 记录黑名单大小（可以根据需要扩展）
	fmt.Printf("当前Token黑名单大小: %d\n", size)
	return nil
}

// hashToken 生成token的哈希值
func (tb *RedisTokenBlacklist) hashToken(tokenString string) string {
	hash := sha256.Sum256([]byte(tokenString))
	return hex.EncodeToString(hash[:])
}

// GetUserTokens 获取用户的所有黑名单Token（用于用户登出时清理所有token）
func (tb *RedisTokenBlacklist) GetUserTokens(ctx context.Context, userID int64) ([]string, error) {
	if tb.client == nil {
		return nil, fmt.Errorf("Redis客户端未初始化")
	}

	// 这里需要存储用户ID和token的映射关系
	// 可以在添加token时同时存储一个用户ID到token的映射
	// 由于实现较复杂，这里提供基本框架
	return nil, fmt.Errorf("未实现：需要存储用户ID和token的映射关系")
}

// BlacklistUserTokens 将用户的所有Token加入黑名单
func (tb *RedisTokenBlacklist) BlacklistUserTokens(ctx context.Context, userID int64, expireAt time.Time) error {
	// 这里需要先获取用户的所有有效token，然后逐一加入黑名单
	// 由于实现较复杂，这里提供基本框架
	return fmt.Errorf("未实现：需要存储用户ID和token的映射关系")
}