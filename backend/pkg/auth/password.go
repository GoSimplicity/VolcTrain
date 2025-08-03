package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	// 默认加密成本
	DefaultCost = bcrypt.DefaultCost
	// 盐值长度
	SaltLength = 32
)

// PasswordService 密码服务
type PasswordService struct {
	cost int // bcrypt加密成本
}

// NewPasswordService 创建密码服务
func NewPasswordService() *PasswordService {
	return &PasswordService{
		cost: DefaultCost,
	}
}

// NewPasswordServiceWithCost 创建指定成本的密码服务
func NewPasswordServiceWithCost(cost int) *PasswordService {
	return &PasswordService{
		cost: cost,
	}
}

// HashPassword 加密密码
func (p *PasswordService) HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("密码不能为空")
	}

	// 使用bcrypt加密密码
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return "", fmt.Errorf("密码加密失败: %w", err)
	}

	return string(hashedBytes), nil
}

// VerifyPassword 验证密码
func (p *PasswordService) VerifyPassword(hashedPassword, password string) bool {
	if hashedPassword == "" || password == "" {
		return false
	}

	// 使用bcrypt验证密码
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GenerateSalt 生成随机盐值
func GenerateSalt() (string, error) {
	salt := make([]byte, SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("生成盐值失败: %w", err)
	}

	return base64.URLEncoding.EncodeToString(salt), nil
}

// HashPasswordWithSalt 使用自定义盐值加密密码
func HashPasswordWithSalt(password, salt string) (string, error) {
	if password == "" || salt == "" {
		return "", fmt.Errorf("密码和盐值不能为空")
	}

	// 组合密码和盐值
	saltedPassword := password + salt

	// 使用SHA256哈希
	hash := sha256.Sum256([]byte(saltedPassword))

	return base64.URLEncoding.EncodeToString(hash[:]), nil
}

// VerifyPasswordWithSalt 使用自定义盐值验证密码
func VerifyPasswordWithSalt(hashedPassword, password, salt string) bool {
	if hashedPassword == "" || password == "" || salt == "" {
		return false
	}

	// 重新计算哈希值
	computedHash, err := HashPasswordWithSalt(password, salt)
	if err != nil {
		return false
	}

	// 比较哈希值
	return hashedPassword == computedHash
}

// PasswordStrength 密码强度
type PasswordStrength int

const (
	PasswordWeak   PasswordStrength = iota // 弱密码
	PasswordMedium                         // 中等密码
	PasswordStrong                         // 强密码
)

// PasswordValidationRule 密码验证规则
type PasswordValidationRule struct {
	MinLength      int  // 最小长度
	RequireUpper   bool // 需要大写字母
	RequireLower   bool // 需要小写字母
	RequireDigit   bool // 需要数字
	RequireSpecial bool // 需要特殊字符
}

// DefaultPasswordRule 默认密码规则
var DefaultPasswordRule = PasswordValidationRule{
	MinLength:      8,
	RequireUpper:   true,
	RequireLower:   true,
	RequireDigit:   true,
	RequireSpecial: false,
}

// StrongPasswordRule 强密码规则
var StrongPasswordRule = PasswordValidationRule{
	MinLength:      12,
	RequireUpper:   true,
	RequireLower:   true,
	RequireDigit:   true,
	RequireSpecial: true,
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string, rule PasswordValidationRule) (bool, []string) {
	var errors []string

	// 检查长度
	if len(password) < rule.MinLength {
		errors = append(errors, fmt.Sprintf("密码长度至少为%d位", rule.MinLength))
	}

	// 检查大写字母
	if rule.RequireUpper {
		hasUpper := false
		for _, char := range password {
			if char >= 'A' && char <= 'Z' {
				hasUpper = true
				break
			}
		}
		if !hasUpper {
			errors = append(errors, "密码必须包含大写字母")
		}
	}

	// 检查小写字母
	if rule.RequireLower {
		hasLower := false
		for _, char := range password {
			if char >= 'a' && char <= 'z' {
				hasLower = true
				break
			}
		}
		if !hasLower {
			errors = append(errors, "密码必须包含小写字母")
		}
	}

	// 检查数字
	if rule.RequireDigit {
		hasDigit := false
		for _, char := range password {
			if char >= '0' && char <= '9' {
				hasDigit = true
				break
			}
		}
		if !hasDigit {
			errors = append(errors, "密码必须包含数字")
		}
	}

	// 检查特殊字符
	if rule.RequireSpecial {
		specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"
		hasSpecial := false
		for _, char := range password {
			for _, special := range specialChars {
				if char == special {
					hasSpecial = true
					break
				}
			}
			if hasSpecial {
				break
			}
		}
		if !hasSpecial {
			errors = append(errors, "密码必须包含特殊字符")
		}
	}

	return len(errors) == 0, errors
}

// GetPasswordStrength 获取密码强度
func GetPasswordStrength(password string) PasswordStrength {
	score := 0

	// 长度评分
	if len(password) >= 8 {
		score++
	}
	if len(password) >= 12 {
		score++
	}

	// 字符类型评分
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		} else if char >= 'a' && char <= 'z' {
			hasLower = true
		} else if char >= '0' && char <= '9' {
			hasDigit = true
		} else {
			hasSpecial = true
		}
	}

	if hasUpper {
		score++
	}
	if hasLower {
		score++
	}
	if hasDigit {
		score++
	}
	if hasSpecial {
		score++
	}

	// 判断强度
	if score >= 5 {
		return PasswordStrong
	} else if score >= 3 {
		return PasswordMedium
	} else {
		return PasswordWeak
	}
}

// LoginAttempt 登录尝试记录
type LoginAttempt struct {
	UserID    int64     `json:"userId"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"userAgent"`
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

// LoginSecurity 登录安全管理
type LoginSecurity struct {
	maxAttempts  int                    // 最大尝试次数
	lockDuration time.Duration          // 锁定时间
	attempts     map[string][]time.Time // IP尝试记录
}

// NewLoginSecurity 创建登录安全管理器
func NewLoginSecurity(maxAttempts int, lockDuration time.Duration) *LoginSecurity {
	return &LoginSecurity{
		maxAttempts:  maxAttempts,
		lockDuration: lockDuration,
		attempts:     make(map[string][]time.Time),
	}
}

// IsBlocked 检查IP是否被锁定
func (ls *LoginSecurity) IsBlocked(ip string) bool {
	attempts, exists := ls.attempts[ip]
	if !exists {
		return false
	}

	// 清理过期的尝试记录
	now := time.Now()
	var validAttempts []time.Time
	for _, attempt := range attempts {
		if now.Sub(attempt) < ls.lockDuration {
			validAttempts = append(validAttempts, attempt)
		}
	}
	ls.attempts[ip] = validAttempts

	// 检查是否超过最大尝试次数
	return len(validAttempts) >= ls.maxAttempts
}

// RecordFailedAttempt 记录失败尝试
func (ls *LoginSecurity) RecordFailedAttempt(ip string) {
	now := time.Now()
	if _, exists := ls.attempts[ip]; !exists {
		ls.attempts[ip] = []time.Time{}
	}
	ls.attempts[ip] = append(ls.attempts[ip], now)
}

// ClearAttempts 清除IP的尝试记录
func (ls *LoginSecurity) ClearAttempts(ip string) {
	delete(ls.attempts, ip)
}

// GetRemainingAttempts 获取剩余尝试次数
func (ls *LoginSecurity) GetRemainingAttempts(ip string) int {
	attempts, exists := ls.attempts[ip]
	if !exists {
		return ls.maxAttempts
	}

	// 清理过期的尝试记录
	now := time.Now()
	var validAttempts []time.Time
	for _, attempt := range attempts {
		if now.Sub(attempt) < ls.lockDuration {
			validAttempts = append(validAttempts, attempt)
		}
	}
	ls.attempts[ip] = validAttempts

	remaining := ls.maxAttempts - len(validAttempts)
	if remaining < 0 {
		return 0
	}
	return remaining
}
