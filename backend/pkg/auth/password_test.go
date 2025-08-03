//go:build unit
// +build unit

package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestNewPasswordService(t *testing.T) {
	service := NewPasswordService()
	assert.NotNil(t, service)
	assert.Equal(t, DefaultCost, service.cost)
}

func TestNewPasswordServiceWithCost(t *testing.T) {
	customCost := 12
	service := NewPasswordServiceWithCost(customCost)
	assert.NotNil(t, service)
	assert.Equal(t, customCost, service.cost)
}

func TestPasswordService_HashPassword(t *testing.T) {
	service := NewPasswordService()

	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{
			name:        "有效密码",
			password:    "TestPassword123",
			expectError: false,
		},
		{
			name:        "空密码",
			password:    "",
			expectError: true,
		},
		{
			name:        "短密码",
			password:    "123",
			expectError: false,
		},
		{
			name:        "长密码",
			password:    strings.Repeat("a", 100),
			expectError: true, // bcrypt有72字节限制
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := service.HashPassword(tt.password)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, hashedPassword)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hashedPassword)
				// 验证哈希密码可以被bcrypt解析
				assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(tt.password)))
			}
		})
	}
}

func TestPasswordService_VerifyPassword(t *testing.T) {
	service := NewPasswordService()
	password := "TestPassword123"
	hashedPassword, err := service.HashPassword(password)
	require.NoError(t, err)

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		expected       bool
	}{
		{
			name:           "正确密码",
			hashedPassword: hashedPassword,
			password:       password,
			expected:       true,
		},
		{
			name:           "错误密码",
			hashedPassword: hashedPassword,
			password:       "WrongPassword",
			expected:       false,
		},
		{
			name:           "空哈希密码",
			hashedPassword: "",
			password:       password,
			expected:       false,
		},
		{
			name:           "空密码",
			hashedPassword: hashedPassword,
			password:       "",
			expected:       false,
		},
		{
			name:           "无效哈希",
			hashedPassword: "invalid_hash",
			password:       password,
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.VerifyPassword(tt.hashedPassword, tt.password)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateSalt(t *testing.T) {
	// 生成多个盐值确保不重复
	salts := make(map[string]bool)
	for i := 0; i < 100; i++ {
		salt, err := GenerateSalt()
		assert.NoError(t, err)
		assert.NotEmpty(t, salt)
		assert.False(t, salts[salt], "盐值不应该重复")
		salts[salt] = true
	}
}

func TestHashPasswordWithSalt(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		salt        string
		expectError bool
	}{
		{
			name:        "有效密码和盐值",
			password:    "TestPassword123",
			salt:        "randomsalt",
			expectError: false,
		},
		{
			name:        "空密码",
			password:    "",
			salt:        "randomsalt",
			expectError: true,
		},
		{
			name:        "空盐值",
			password:    "TestPassword123",
			salt:        "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := HashPasswordWithSalt(tt.password, tt.salt)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, hashedPassword)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hashedPassword)
			}
		})
	}
}

func TestVerifyPasswordWithSalt(t *testing.T) {
	password := "TestPassword123"
	salt := "randomsalt"
	hashedPassword, err := HashPasswordWithSalt(password, salt)
	require.NoError(t, err)

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		salt           string
		expected       bool
	}{
		{
			name:           "正确密码和盐值",
			hashedPassword: hashedPassword,
			password:       password,
			salt:           salt,
			expected:       true,
		},
		{
			name:           "错误密码",
			hashedPassword: hashedPassword,
			password:       "WrongPassword",
			salt:           salt,
			expected:       false,
		},
		{
			name:           "错误盐值",
			hashedPassword: hashedPassword,
			password:       password,
			salt:           "wrongsalt",
			expected:       false,
		},
		{
			name:           "空参数",
			hashedPassword: "",
			password:       "",
			salt:           "",
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VerifyPasswordWithSalt(tt.hashedPassword, tt.password, tt.salt)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		rule        PasswordValidationRule
		expectedOK  bool
		expectedErr int
	}{
		{
			name:        "符合默认规则的密码",
			password:    "TestPass123",
			rule:        DefaultPasswordRule,
			expectedOK:  true,
			expectedErr: 0,
		},
		{
			name:        "太短的密码",
			password:    "Test1",
			rule:        DefaultPasswordRule,
			expectedOK:  false,
			expectedErr: 1, // 长度不足
		},
		{
			name:        "缺少大写字母",
			password:    "testpass123",
			rule:        DefaultPasswordRule,
			expectedOK:  false,
			expectedErr: 1, // 缺少大写字母
		},
		{
			name:        "缺少小写字母",
			password:    "TESTPASS123",
			rule:        DefaultPasswordRule,
			expectedOK:  false,
			expectedErr: 1, // 缺少小写字母
		},
		{
			name:        "缺少数字",
			password:    "TestPassword",
			rule:        DefaultPasswordRule,
			expectedOK:  false,
			expectedErr: 1, // 缺少数字
		},
		{
			name:        "符合强密码规则",
			password:    "TestPass123!",
			rule:        StrongPasswordRule,
			expectedOK:  true,
			expectedErr: 0,
		},
		{
			name:        "不符合强密码规则",
			password:    "TestPass123",
			rule:        StrongPasswordRule,
			expectedOK:  false,
			expectedErr: 2, // 长度不足 + 缺少特殊字符
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, errs := ValidatePassword(tt.password, tt.rule)
			assert.Equal(t, tt.expectedOK, ok)
			assert.Equal(t, tt.expectedErr, len(errs))
		})
	}
}

func TestGetPasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected PasswordStrength
	}{
		{
			name:     "弱密码-短且简单",
			password: "123",
			expected: PasswordWeak,
		},
		{
			name:     "弱密码-只有小写",
			password: "password",
			expected: PasswordWeak,
		},
		{
			name:     "中等密码",
			password: "Password123",
			expected: PasswordMedium,
		},
		{
			name:     "强密码",
			password: "StrongPass123!",
			expected: PasswordStrong,
		},
		{
			name:     "超强密码",
			password: "VeryStrongPassword123!@#",
			expected: PasswordStrong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strength := GetPasswordStrength(tt.password)
			assert.Equal(t, tt.expected, strength)
		})
	}
}

func TestNewLoginSecurity(t *testing.T) {
	maxAttempts := 5
	lockDuration := 10 * time.Minute

	ls := NewLoginSecurity(maxAttempts, lockDuration)
	assert.NotNil(t, ls)
	assert.Equal(t, maxAttempts, ls.maxAttempts)
	assert.Equal(t, lockDuration, ls.lockDuration)
	assert.NotNil(t, ls.attempts)
}

func TestLoginSecurity_IsBlocked(t *testing.T) {
	ls := NewLoginSecurity(3, 5*time.Minute)
	ip := "192.168.1.1"

	// 初始状态不应该被阻止
	assert.False(t, ls.IsBlocked(ip))

	// 记录失败尝试
	for i := 0; i < 3; i++ {
		ls.RecordFailedAttempt(ip)
	}

	// 达到最大尝试次数后应该被阻止
	assert.True(t, ls.IsBlocked(ip))
}

func TestLoginSecurity_RecordFailedAttempt(t *testing.T) {
	ls := NewLoginSecurity(3, 5*time.Minute)
	ip := "192.168.1.1"

	// 记录几次失败尝试
	ls.RecordFailedAttempt(ip)
	ls.RecordFailedAttempt(ip)

	assert.Equal(t, 1, ls.GetRemainingAttempts(ip))
}

func TestLoginSecurity_ClearAttempts(t *testing.T) {
	ls := NewLoginSecurity(3, 5*time.Minute)
	ip := "192.168.1.1"

	// 记录失败尝试
	ls.RecordFailedAttempt(ip)
	ls.RecordFailedAttempt(ip)

	// 清除尝试记录
	ls.ClearAttempts(ip)

	// 应该恢复到初始状态
	assert.False(t, ls.IsBlocked(ip))
	assert.Equal(t, 3, ls.GetRemainingAttempts(ip))
}

func TestLoginSecurity_GetRemainingAttempts(t *testing.T) {
	ls := NewLoginSecurity(5, 5*time.Minute)
	ip := "192.168.1.1"

	// 初始剩余尝试次数
	assert.Equal(t, 5, ls.GetRemainingAttempts(ip))

	// 记录失败尝试
	ls.RecordFailedAttempt(ip)
	assert.Equal(t, 4, ls.GetRemainingAttempts(ip))

	ls.RecordFailedAttempt(ip)
	assert.Equal(t, 3, ls.GetRemainingAttempts(ip))
}

// 基准测试
func BenchmarkPasswordService_HashPassword(b *testing.B) {
	service := NewPasswordService()
	password := "BenchmarkPassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.HashPassword(password)
	}
}

func BenchmarkPasswordService_VerifyPassword(b *testing.B) {
	service := NewPasswordService()
	password := "BenchmarkPassword123"
	hashedPassword, _ := service.HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.VerifyPassword(hashedPassword, password)
	}
}

func BenchmarkGenerateSalt(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateSalt()
	}
}
