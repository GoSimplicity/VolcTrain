//go:build unit
// +build unit

package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTService(t *testing.T) {
	accessSecret := "test-access-secret"
	refreshSecret := "test-refresh-secret"
	accessExpire := int64(3600)  // 1小时
	refreshExpire := int64(7200) // 2小时

	service := NewJWTService(accessSecret, refreshSecret, accessExpire, refreshExpire)

	assert.NotNil(t, service)
	assert.Equal(t, accessSecret, service.accessSecret)
	assert.Equal(t, refreshSecret, service.refreshSecret)
	assert.Equal(t, time.Duration(accessExpire)*time.Second, service.accessExpire)
	assert.Equal(t, time.Duration(refreshExpire)*time.Second, service.refreshExpire)
	assert.Equal(t, "volctrain", service.issuer)
}

func TestJWTService_GenerateTokenPair(t *testing.T) {
	service := NewJWTService("access-secret", "refresh-secret", 3600, 7200)

	userID := int64(123)
	username := "testuser"
	email := "test@example.com"
	roles := []string{"user", "admin"}
	permissions := []string{"read", "write"}

	tokenPair, err := service.GenerateTokenPair(userID, username, email, roles, permissions)

	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	assert.Equal(t, int64(3600), tokenPair.ExpiresIn)
	assert.Equal(t, "Bearer", tokenPair.TokenType)
}

func TestJWTService_ParseAccessToken(t *testing.T) {
	service := NewJWTService("access-secret", "refresh-secret", 3600, 7200)

	userID := int64(123)
	username := "testuser"
	email := "test@example.com"
	roles := []string{"user", "admin"}
	permissions := []string{"read", "write"}

	// 生成Token
	tokenPair, err := service.GenerateTokenPair(userID, username, email, roles, permissions)
	require.NoError(t, err)

	// 解析访问Token
	claims, err := service.ParseAccessToken(tokenPair.AccessToken)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, roles, claims.Roles)
	assert.Equal(t, permissions, claims.Permissions)
	assert.Equal(t, "access", claims.TokenType)
	assert.Equal(t, "volctrain", claims.Issuer)
}

func TestJWTService_ParseRefreshToken(t *testing.T) {
	service := NewJWTService("access-secret", "refresh-secret", 3600, 7200)

	userID := int64(123)
	username := "testuser"
	email := "test@example.com"
	roles := []string{"user"}
	permissions := []string{"read"}

	// 生成Token
	tokenPair, err := service.GenerateTokenPair(userID, username, email, roles, permissions)
	require.NoError(t, err)

	// 解析刷新Token
	claims, err := service.ParseRefreshToken(tokenPair.RefreshToken)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, "refresh", claims.TokenType)
}

func TestJWTService_ParseToken_InvalidCases(t *testing.T) {
	service := NewJWTService("access-secret", "refresh-secret", 3600, 7200)

	tests := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "空Token",
			token:       "",
			expectError: true,
		},
		{
			name:        "无效Token格式",
			token:       "invalid.token.format",
			expectError: true,
		},
		{
			name:        "错误签名的Token",
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.ParseAccessToken(tt.token)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestJWTService_RefreshToken(t *testing.T) {
	service := NewJWTService("access-secret", "refresh-secret", 3600, 7200)

	userID := int64(123)
	username := "testuser"
	email := "test@example.com"
	roles := []string{"user"}
	permissions := []string{"read"}

	// 生成初始Token对
	tokenPair, err := service.GenerateTokenPair(userID, username, email, roles, permissions)
	require.NoError(t, err)

	// 刷新Token测试 - 添加小延迟确保时间戳不同
	time.Sleep(1 * time.Millisecond)

	// 使用刷新Token生成新的Token对
	newTokenPair, err := service.RefreshToken(tokenPair.RefreshToken)

	assert.NoError(t, err)
	assert.NotNil(t, newTokenPair)
	assert.NotEmpty(t, newTokenPair.AccessToken)
	assert.NotEmpty(t, newTokenPair.RefreshToken)

	// 验证新Token包含正确信息
	newClaims, err := service.ParseAccessToken(newTokenPair.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, userID, newClaims.UserID)
}

func TestJWTService_RefreshToken_ExpiredToken(t *testing.T) {
	// 创建一个短期过期的JWT服务
	service := NewJWTService("access-secret", "refresh-secret", 1, 1) // 1秒过期

	userID := int64(123)
	username := "testuser"
	email := "test@example.com"
	roles := []string{"user"}
	permissions := []string{"read"}

	// 生成Token
	tokenPair, err := service.GenerateTokenPair(userID, username, email, roles, permissions)
	require.NoError(t, err)

	// 等待Token过期
	time.Sleep(2 * time.Second)

	// 尝试刷新过期的Token
	_, err = service.RefreshToken(tokenPair.RefreshToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}

func TestJWTService_ValidateToken(t *testing.T) {
	service := NewJWTService("access-secret", "refresh-secret", 3600, 7200)

	userID := int64(123)
	username := "testuser"
	email := "test@example.com"
	roles := []string{"user"}
	permissions := []string{"read"}

	// 生成Token
	tokenPair, err := service.GenerateTokenPair(userID, username, email, roles, permissions)
	require.NoError(t, err)

	// 验证有效Token
	assert.True(t, service.ValidateToken(tokenPair.AccessToken))

	// 验证无效Token
	assert.False(t, service.ValidateToken("invalid.token"))
	assert.False(t, service.ValidateToken(""))
}

func TestJWTService_GetUserIDFromToken(t *testing.T) {
	service := NewJWTService("access-secret", "refresh-secret", 3600, 7200)

	userID := int64(12345)
	username := "testuser"
	email := "test@example.com"
	roles := []string{"user"}
	permissions := []string{"read"}

	// 生成Token
	tokenPair, err := service.GenerateTokenPair(userID, username, email, roles, permissions)
	require.NoError(t, err)

	// 从Token中获取用户ID
	extractedUserID, err := service.GetUserIDFromToken(tokenPair.AccessToken)

	assert.NoError(t, err)
	assert.Equal(t, userID, extractedUserID)
}

func TestJWTService_GetClaimsFromToken(t *testing.T) {
	service := NewJWTService("access-secret", "refresh-secret", 3600, 7200)

	userID := int64(123)
	username := "testuser"
	email := "test@example.com"
	roles := []string{"user", "admin"}
	permissions := []string{"read", "write", "delete"}

	// 生成Token
	tokenPair, err := service.GenerateTokenPair(userID, username, email, roles, permissions)
	require.NoError(t, err)

	// 从Token中获取所有声明
	claims, err := service.GetClaimsFromToken(tokenPair.AccessToken)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, roles, claims.Roles)
	assert.Equal(t, permissions, claims.Permissions)
	assert.Equal(t, "access", claims.TokenType)
}

func TestJWTService_TokenTypeMismatch(t *testing.T) {
	userID := int64(123)
	username := "testuser"
	email := "test@example.com"
	roles := []string{"user"}
	permissions := []string{"read"}

	// 尝试用刷新Token解析方法解析访问Token（使用相同secret）
	service := NewJWTService("same-secret", "same-secret", 3600, 7200)
	tokenPair, err := service.GenerateTokenPair(userID, username, email, roles, permissions)
	require.NoError(t, err)

	_, err = service.ParseRefreshToken(tokenPair.AccessToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Token类型不匹配")

	// 尝试用访问Token解析方法解析刷新Token
	_, err = service.ParseAccessToken(tokenPair.RefreshToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Token类型不匹配")
}

func TestTokenBlacklist(t *testing.T) {
	blacklist := NewTokenBlacklist()

	token1 := "token1"
	token2 := "token2"
	expireAt := time.Now().Add(1 * time.Hour)

	// 初始状态
	assert.Equal(t, 0, blacklist.Size())
	assert.False(t, blacklist.IsBlacklisted(token1))

	// 添加Token到黑名单
	blacklist.AddToken(token1, expireAt)
	assert.Equal(t, 1, blacklist.Size())
	assert.True(t, blacklist.IsBlacklisted(token1))
	assert.False(t, blacklist.IsBlacklisted(token2))

	// 添加更多Token
	blacklist.AddToken(token2, expireAt)
	assert.Equal(t, 2, blacklist.Size())
	assert.True(t, blacklist.IsBlacklisted(token2))

	// 移除Token
	blacklist.RemoveToken(token1)
	assert.Equal(t, 1, blacklist.Size())
	assert.False(t, blacklist.IsBlacklisted(token1))
	assert.True(t, blacklist.IsBlacklisted(token2))

	// 清空黑名单
	blacklist.Clear()
	assert.Equal(t, 0, blacklist.Size())
	assert.False(t, blacklist.IsBlacklisted(token2))
}

func TestTokenBlacklist_ExpiredTokens(t *testing.T) {
	blacklist := NewTokenBlacklist()

	token1 := "expired_token"
	token2 := "valid_token"

	// 添加已过期的Token
	expiredTime := time.Now().Add(-1 * time.Hour)
	blacklist.AddToken(token1, expiredTime)

	// 添加有效的Token
	validTime := time.Now().Add(1 * time.Hour)
	blacklist.AddToken(token2, validTime)

	// 过期Token不应该被认为在黑名单中
	assert.False(t, blacklist.IsBlacklisted(token1))
	assert.True(t, blacklist.IsBlacklisted(token2))

	// Size应该只计算有效的Token
	assert.Equal(t, 1, blacklist.Size())
}

func TestJWTService_EdgeCases(t *testing.T) {
	service := NewJWTService("access-secret", "refresh-secret", 3600, 7200)

	t.Run("空用户信息", func(t *testing.T) {
		tokenPair, err := service.GenerateTokenPair(0, "", "", nil, nil)
		assert.NoError(t, err)
		assert.NotNil(t, tokenPair)

		claims, err := service.ParseAccessToken(tokenPair.AccessToken)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), claims.UserID)
		assert.Equal(t, "", claims.Username)
		assert.Empty(t, claims.Roles)
		assert.Empty(t, claims.Permissions)
	})

	t.Run("特殊字符用户名", func(t *testing.T) {
		specialUsername := "用户@#$%^&*()"
		tokenPair, err := service.GenerateTokenPair(123, specialUsername, "test@example.com", nil, nil)
		assert.NoError(t, err)

		claims, err := service.ParseAccessToken(tokenPair.AccessToken)
		assert.NoError(t, err)
		assert.Equal(t, specialUsername, claims.Username)
	})
}

// 基准测试
func BenchmarkJWTService_GenerateTokenPair(b *testing.B) {
	service := NewJWTService("access-secret", "refresh-secret", 3600, 7200)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GenerateTokenPair(123, "testuser", "test@example.com", []string{"user"}, []string{"read"})
	}
}

func BenchmarkJWTService_ParseAccessToken(b *testing.B) {
	service := NewJWTService("access-secret", "refresh-secret", 3600, 7200)
	tokenPair, _ := service.GenerateTokenPair(123, "testuser", "test@example.com", []string{"user"}, []string{"read"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ParseAccessToken(tokenPair.AccessToken)
	}
}

func BenchmarkJWTService_ValidateToken(b *testing.B) {
	service := NewJWTService("access-secret", "refresh-secret", 3600, 7200)
	tokenPair, _ := service.GenerateTokenPair(123, "testuser", "test@example.com", []string{"user"}, []string{"read"})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.ValidateToken(tokenPair.AccessToken)
	}
}
