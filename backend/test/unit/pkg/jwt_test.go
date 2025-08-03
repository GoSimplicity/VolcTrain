//go:build unit
// +build unit

package auth

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"api/pkg/auth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// JWTTestSuite JWT测试套件
type JWTTestSuite struct {
	suite.Suite
	jwtService *auth.JWTService
	secretKey  string
}

// SetupTest 每个测试前的初始化
func (suite *JWTTestSuite) SetupTest() {
	suite.secretKey = "test-secret-key-for-jwt-testing"
	suite.jwtService = auth.NewJWTService(suite.secretKey, suite.secretKey, 3600, 7*24*3600)
}

// TestGenerateTokenSuccess 测试成功生成Token
func (suite *JWTTestSuite) TestGenerateTokenSuccess() {
	// 准备测试数据
	userID := int64(1001)
	username := "testuser"

	// 执行测试
	tokenPair, err := suite.jwtService.GenerateTokenPair(userID, username, "test@example.com", []string{"user"}, []string{"read", "write"})

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), tokenPair.AccessToken)

	// 验证Token格式（JWT包含3个部分，用.分隔）
	parts := strings.Split(tokenPair.AccessToken, ".")
	assert.Len(suite.T(), parts, 3, "JWT应该包含3个部分")
}

// TestValidateTokenSuccess 测试成功验证Token
func (suite *JWTTestSuite) TestValidateTokenSuccess() {
	// 准备测试数据
	userID := int64(1001)
	username := "testuser"

	// 生成Token
	tokenPair, err := suite.jwtService.GenerateTokenPair(userID, username, "test@example.com", []string{"admin"}, []string{"read", "write", "admin"})
	require.NoError(suite.T(), err)

	// 执行测试
	claims, err := suite.jwtService.ParseAccessToken(tokenPair.AccessToken)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), claims)
	assert.Equal(suite.T(), userID, claims.UserID)
	assert.Equal(suite.T(), username, claims.Username)
	assert.Contains(suite.T(), claims.Roles, "admin")
	assert.True(suite.T(), claims.ExpiresAt.After(time.Now()))
}

// TestValidateTokenInvalid 测试验证无效Token
func (suite *JWTTestSuite) TestValidateTokenInvalid() {
	invalidTokens := []struct {
		name  string
		token string
	}{
		{"空Token", ""},
		{"格式错误", "invalid.token.format"},
		{"伪造Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.fake.signature"},
		{"错误签名", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.wrong_signature"},
	}

	for _, tc := range invalidTokens {
		suite.Run(tc.name, func() {
			// 执行测试
			claims, err := suite.jwtService.ParseAccessToken(tc.token)

			// 验证结果
			assert.Error(suite.T(), err)
			assert.Nil(suite.T(), claims)
		})
	}
}

// TestValidateTokenExpired 测试验证过期Token
func (suite *JWTTestSuite) TestValidateTokenExpired() {
	// 创建一个极短有效期的JWT服务
	shortLivedService := auth.NewJWTService(suite.secretKey, suite.secretKey, 1, 7*24*3600) // 1秒过期

	// 生成Token
	tokenPair, err := shortLivedService.GenerateTokenPair(1001, "testuser", "test@example.com", []string{"user"}, []string{"read"})
	require.NoError(suite.T(), err)

	// 等待Token过期
	time.Sleep(2 * time.Second)

	// 执行测试
	claims, err := shortLivedService.ParseAccessToken(tokenPair.AccessToken)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
	assert.Contains(suite.T(), err.Error(), "expired")
}

// TestRefreshTokenSuccess 测试成功刷新Token
func (suite *JWTTestSuite) TestRefreshTokenSuccess() {
	// 生成原始Token对
	originalTokenPair, err := suite.jwtService.GenerateTokenPair(1001, "testuser", "test@example.com", []string{"user"}, []string{"read"})
	require.NoError(suite.T(), err)

	// 稍等片刻确保时间戳不同
	time.Sleep(10 * time.Millisecond)

	// 执行测试
	newTokenPair, err := suite.jwtService.RefreshToken(originalTokenPair.RefreshToken)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), newTokenPair.AccessToken)
	assert.NotEqual(suite.T(), originalTokenPair.AccessToken, newTokenPair.AccessToken)

	// 验证新Token内容
	claims, err := suite.jwtService.ParseAccessToken(newTokenPair.AccessToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1001), claims.UserID)
	assert.Equal(suite.T(), "testuser", claims.Username)
}

// TestGetUserIDFromToken 测试从Token提取用户ID
func (suite *JWTTestSuite) TestGetUserIDFromToken() {
	// 准备测试数据
	userID := int64(1001)
	username := "testuser"

	// 生成Token
	tokenPair, err := suite.jwtService.GenerateTokenPair(userID, username, "test@example.com", []string{"admin"}, []string{"read", "write"})
	require.NoError(suite.T(), err)

	// 执行测试
	extractedUserID, err := suite.jwtService.GetUserIDFromToken(tokenPair.AccessToken)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, extractedUserID)
}

// TestTokenWithDifferentSecrets 测试不同密钥生成的Token
func (suite *JWTTestSuite) TestTokenWithDifferentSecrets() {
	// 使用不同密钥创建另一个JWT服务
	differentService := auth.NewJWTService("different-secret-key", "different-secret-key", 3600, 7*24*3600)

	// 用原服务生成Token
	tokenPair, err := suite.jwtService.GenerateTokenPair(1001, "testuser", "test@example.com", []string{"user"}, []string{"read"})
	require.NoError(suite.T(), err)

	// 用不同密钥的服务验证Token
	claims, err := differentService.ParseAccessToken(tokenPair.AccessToken)

	// 验证结果 - 应该失败
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
}

// TestTokenClaims 测试Token声明字段
func (suite *JWTTestSuite) TestTokenClaims() {
	userID := int64(1001)
	username := "testuser"

	tokenPair, err := suite.jwtService.GenerateTokenPair(userID, username, "test@example.com", []string{"admin"}, []string{"read", "write"})
	require.NoError(suite.T(), err)

	// 手动解析Token进行详细验证
	parsedToken, err := jwt.ParseWithClaims(tokenPair.AccessToken, &auth.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(suite.secretKey), nil
	})

	assert.NoError(suite.T(), err)
	assert.True(suite.T(), parsedToken.Valid)

	claims := parsedToken.Claims.(*auth.JWTClaims)
	assert.Equal(suite.T(), userID, claims.UserID)
	assert.Equal(suite.T(), username, claims.Username)
	assert.Contains(suite.T(), claims.Roles, "admin")
	assert.Contains(suite.T(), claims.Permissions, "read")
	assert.True(suite.T(), claims.IssuedAt.Before(time.Now()))
	assert.True(suite.T(), claims.ExpiresAt.After(time.Now()))
}

// TestInvalidUserData 测试无效用户数据
func (suite *JWTTestSuite) TestInvalidUserData() {
	testCases := []struct {
		name      string
		userID    int64
		username  string
		email     string
		expectErr bool
	}{
		{"有效数据", 1001, "validuser", "test@example.com", false},
		{"无效用户ID", 0, "validuser", "test@example.com", false},  // JWTService不验证ID
		{"负数用户ID", -1, "validuser", "test@example.com", false}, // JWTService不验证ID
		{"空用户名", 1001, "", "test@example.com", false},          // JWTService不验证用户名
		{"空邮箱", 1001, "validuser", "", false},                  // JWTService不验证邮箱
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// 执行测试
			tokenPair, err := suite.jwtService.GenerateTokenPair(tc.userID, tc.username, tc.email, []string{"user"}, []string{"read"})

			if tc.expectErr {
				assert.Error(suite.T(), err)
				assert.Nil(suite.T(), tokenPair)
			} else {
				assert.NoError(suite.T(), err)
				assert.NotNil(suite.T(), tokenPair)
				assert.NotEmpty(suite.T(), tokenPair.AccessToken)
			}
		})
	}
}

// TestConcurrentTokenOperations 测试并发Token操作
func (suite *JWTTestSuite) TestConcurrentTokenOperations() {
	const numGoroutines = 100
	results := make(chan struct {
		token string
		err   error
	}, numGoroutines)

	// 并发生成Token
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			userID := int64(1000 + id)
			username := fmt.Sprintf("user%d", id)
			email := fmt.Sprintf("user%d@example.com", id)
			tokenPair, err := suite.jwtService.GenerateTokenPair(userID, username, email, []string{"user"}, []string{"read"})
			if err != nil {
				results <- struct {
					token string
					err   error
				}{"", err}
			} else {
				results <- struct {
					token string
					err   error
				}{tokenPair.AccessToken, nil}
			}
		}(i)
	}

	// 收集结果
	tokens := make([]string, 0, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		result := <-results
		assert.NoError(suite.T(), result.err)
		assert.NotEmpty(suite.T(), result.token)
		tokens = append(tokens, result.token)
	}

	// 验证所有Token都是唯一的
	tokenSet := make(map[string]bool)
	for _, token := range tokens {
		assert.False(suite.T(), tokenSet[token], "Token应该是唯一的")
		tokenSet[token] = true

		// 验证每个Token都有效
		claims, err := suite.jwtService.ParseAccessToken(token)
		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), claims)
	}
}

// 运行JWT测试套件
func TestJWTSuite(t *testing.T) {
	suite.Run(t, new(JWTTestSuite))
}
