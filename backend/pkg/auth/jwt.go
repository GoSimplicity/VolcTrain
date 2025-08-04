package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 简化的JWT工具函数，用于快速集成

// GenerateToken 生成JWT token
func GenerateToken(userID int64, secret string, expire int64) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"userId": userID,
		"iat":    now.Unix(),
		"exp":    now.Add(time.Duration(expire) * time.Second).Unix(),
		"iss":    "volctrain",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken 验证token并返回用户ID
func ValidateToken(tokenString, secret string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("无效的签名方法: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("token解析失败: %w", err)
	}

	if !token.Valid {
		return 0, errors.New("token无效")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("无效的token声明")
	}

	userID, ok := claims["userId"].(float64)
	if !ok {
		return 0, errors.New("无效的用户ID")
	}

	return int64(userID), nil
}

// JWTClaims JWT声明
type JWTClaims struct {
	UserID      int64    `json:"userId"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	WorkspaceID int64    `json:"workspaceId,omitempty"`
	TokenType   string   `json:"tokenType"` // access or refresh
	jwt.RegisteredClaims
}

// JWTService JWT服务
type JWTService struct {
	accessSecret  string
	refreshSecret string
	accessExpire  time.Duration
	refreshExpire time.Duration
	issuer        string
}

// NewJWTService 创建JWT服务
func NewJWTService(accessSecret, refreshSecret string, accessExpire, refreshExpire int64) *JWTService {
	return &JWTService{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessExpire:  time.Duration(accessExpire) * time.Second,
		refreshExpire: time.Duration(refreshExpire) * time.Second,
		issuer:        "volctrain",
	}
}

// TokenPair Token对
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
}

// GenerateTokenPair 生成Token对
func (j *JWTService) GenerateTokenPair(userID int64, username, email string, roles, permissions []string) (*TokenPair, error) {
	now := time.Now()

	// 生成访问Token
	accessClaims := JWTClaims{
		UserID:      userID,
		Username:    username,
		Email:       email,
		Roles:       roles,
		Permissions: permissions,
		TokenType:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   fmt.Sprintf("user:%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessExpire)),
		},
	}

	accessToken, err := j.generateToken(accessClaims, j.accessSecret)
	if err != nil {
		return nil, fmt.Errorf("生成访问Token失败: %w", err)
	}

	// 生成刷新Token
	refreshClaims := JWTClaims{
		UserID:    userID,
		Username:  username,
		Email:     email,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   fmt.Sprintf("user:%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshExpire)),
		},
	}

	refreshToken, err := j.generateToken(refreshClaims, j.refreshSecret)
	if err != nil {
		return nil, fmt.Errorf("生成刷新Token失败: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(j.accessExpire.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// generateToken 生成Token
func (j *JWTService) generateToken(claims JWTClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseAccessToken 解析访问Token
func (j *JWTService) ParseAccessToken(tokenString string) (*JWTClaims, error) {
	return j.parseToken(tokenString, j.accessSecret, "access")
}

// ParseRefreshToken 解析刷新Token
func (j *JWTService) ParseRefreshToken(tokenString string) (*JWTClaims, error) {
	return j.parseToken(tokenString, j.refreshSecret, "refresh")
}

// parseToken 解析Token
func (j *JWTService) parseToken(tokenString, secret, expectedType string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("无效的签名方法: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("Token解析失败: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("Token无效")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("无效的Token声明")
	}

	// 验证Token类型
	if claims.TokenType != expectedType {
		return nil, fmt.Errorf("Token类型不匹配，期望: %s, 实际: %s", expectedType, claims.TokenType)
	}

	// 验证发行者
	if claims.Issuer != j.issuer {
		return nil, fmt.Errorf("无效的发行者: %s", claims.Issuer)
	}

	return claims, nil
}

// RefreshToken 刷新Token
func (j *JWTService) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	// 解析刷新Token
	claims, err := j.ParseRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("刷新Token无效: %w", err)
	}

	// 检查Token是否过期
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("刷新Token已过期")
	}

	// 生成新的Token对（注意：这里需要重新获取用户的最新权限信息）
	return j.GenerateTokenPair(claims.UserID, claims.Username, claims.Email, claims.Roles, claims.Permissions)
}

// ValidateToken 验证Token有效性
func (j *JWTService) ValidateToken(tokenString string) bool {
	_, err := j.ParseAccessToken(tokenString)
	return err == nil
}

// GetUserIDFromToken 从Token中获取用户ID
func (j *JWTService) GetUserIDFromToken(tokenString string) (int64, error) {
	claims, err := j.ParseAccessToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// GetClaimsFromToken 从Token中获取所有声明
func (j *JWTService) GetClaimsFromToken(tokenString string) (*JWTClaims, error) {
	return j.ParseAccessToken(tokenString)
}

// BlacklistToken Token黑名单管理（简单内存实现，生产环境建议使用Redis）
type TokenBlacklist struct {
	tokens map[string]time.Time // token -> 过期时间
}

// NewTokenBlacklist 创建Token黑名单
func NewTokenBlacklist() *TokenBlacklist {
	return &TokenBlacklist{
		tokens: make(map[string]time.Time),
	}
}

// AddToken 添加Token到黑名单
func (tb *TokenBlacklist) AddToken(tokenString string, expireAt time.Time) {
	tb.tokens[tokenString] = expireAt

	// 清理过期的Token
	tb.cleanup()
}

// IsBlacklisted 检查Token是否在黑名单中
func (tb *TokenBlacklist) IsBlacklisted(tokenString string) bool {
	expireAt, exists := tb.tokens[tokenString]
	if !exists {
		return false
	}

	// 如果Token已过期，从黑名单中移除
	if time.Now().After(expireAt) {
		delete(tb.tokens, tokenString)
		return false
	}

	return true
}

// cleanup 清理过期的Token
func (tb *TokenBlacklist) cleanup() {
	now := time.Now()
	for token, expireAt := range tb.tokens {
		if now.After(expireAt) {
			delete(tb.tokens, token)
		}
	}
}

// RemoveToken 从黑名单中移除Token
func (tb *TokenBlacklist) RemoveToken(tokenString string) {
	delete(tb.tokens, tokenString)
}

// Clear 清空黑名单
func (tb *TokenBlacklist) Clear() {
	tb.tokens = make(map[string]time.Time)
}

// Size 获取黑名单大小
func (tb *TokenBlacklist) Size() int {
	tb.cleanup()
	return len(tb.tokens)
}
