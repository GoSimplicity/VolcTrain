//go:build unit
// +build unit

package auth

import (
	"testing"

	"api/pkg/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// PasswordTestSuite 密码测试套件
type PasswordTestSuite struct {
	suite.Suite
	passwordService *auth.PasswordService
}

// SetupTest 每个测试前的初始化
func (suite *PasswordTestSuite) SetupTest() {
	suite.passwordService = auth.NewPasswordService()
}

// TestHashPasswordSuccess 测试成功哈希密码
func (suite *PasswordTestSuite) TestHashPasswordSuccess() {
	password := "mySecurePassword123!"

	// 执行测试
	hashedPassword, err := suite.passwordService.HashPassword(password)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), hashedPassword)
	assert.NotEqual(suite.T(), password, hashedPassword)
}

// TestHashPasswordEmptyPassword 测试空密码
func (suite *PasswordTestSuite) TestHashPasswordEmptyPassword() {
	password := ""

	// 执行测试
	hashedPassword, err := suite.passwordService.HashPassword(password)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), hashedPassword)
	assert.Contains(suite.T(), err.Error(), "密码不能为空")
}

// TestHashPasswordDifferentResults 测试相同密码产生不同哈希
func (suite *PasswordTestSuite) TestHashPasswordDifferentResults() {
	password := "samePassword123"

	// 执行两次哈希
	hash1, err1 := suite.passwordService.HashPassword(password)
	hash2, err2 := suite.passwordService.HashPassword(password)

	// 验证结果
	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	assert.NotEqual(suite.T(), hash1, hash2, "相同密码应产生不同哈希值(因为bcrypt包含随机盐)")
}

// TestVerifyPasswordSuccess 测试成功验证密码
func (suite *PasswordTestSuite) TestVerifyPasswordSuccess() {
	password := "correctPassword123"

	// 先哈希密码
	hashedPassword, err := suite.passwordService.HashPassword(password)
	require.NoError(suite.T(), err)

	// 执行测试
	isValid := suite.passwordService.VerifyPassword(hashedPassword, password)

	// 验证结果
	assert.True(suite.T(), isValid)
}

// TestVerifyPasswordIncorrect 测试错误密码验证
func (suite *PasswordTestSuite) TestVerifyPasswordIncorrect() {
	correctPassword := "correctPassword123"
	incorrectPassword := "wrongPassword456"

	// 先哈希正确密码
	hashedPassword, err := suite.passwordService.HashPassword(correctPassword)
	require.NoError(suite.T(), err)

	// 执行测试 - 用错误密码验证
	isValid := suite.passwordService.VerifyPassword(hashedPassword, incorrectPassword)

	// 验证结果
	assert.False(suite.T(), isValid)
}

// TestVerifyPasswordEmptyInputs 测试空输入验证
func (suite *PasswordTestSuite) TestVerifyPasswordEmptyInputs() {
	testCases := []struct {
		name           string
		password       string
		hashedPassword string
		expected       bool
	}{
		{"空密码", "", "hash", false},
		{"空哈希", "password", "", false},
		{"全空", "", "", false},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// 执行测试
			isValid := suite.passwordService.VerifyPassword(tc.hashedPassword, tc.password)

			// 验证结果
			assert.Equal(suite.T(), tc.expected, isValid)
		})
	}
}

// TestPasswordValidation 测试密码强度验证
func (suite *PasswordTestSuite) TestPasswordValidation() {
	// 测试默认规则
	valid, errors := auth.ValidatePassword("StrongPass123!", auth.DefaultPasswordRule)
	assert.True(suite.T(), valid)
	assert.Empty(suite.T(), errors)

	// 测试弱密码
	valid, errors = auth.ValidatePassword("weak", auth.DefaultPasswordRule)
	assert.False(suite.T(), valid)
	assert.NotEmpty(suite.T(), errors)
}

// TestGenerateRandomSalt 测试生成随机盐值
func (suite *PasswordTestSuite) TestGenerateRandomSalt() {
	// 生成多个盐值
	salts := make([]string, 100)
	for i := 0; i < 100; i++ {
		salt, err := auth.GenerateSalt()
		assert.NoError(suite.T(), err)
		assert.NotEmpty(suite.T(), salt)
		salts[i] = salt
	}

	// 验证所有盐值都不相同
	saltSet := make(map[string]bool)
	for _, salt := range salts {
		assert.False(suite.T(), saltSet[salt], "盐值应该是唯一的")
		saltSet[salt] = true
	}
}

// 运行密码测试套件
func TestPasswordSuite(t *testing.T) {
	suite.Run(t, new(PasswordTestSuite))
}
