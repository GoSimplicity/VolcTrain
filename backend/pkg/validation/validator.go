package validation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// ValidationResult 验证结果
type ValidationResult struct {
	IsValid bool
	Errors  []string
}

// NewValidationResult 创建验证结果
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		IsValid: true,
		Errors:  make([]string, 0),
	}
}

// AddError 添加错误信息
func (vr *ValidationResult) AddError(message string) {
	vr.IsValid = false
	vr.Errors = append(vr.Errors, message)
}

// AddErrorIf 添加条件错误
func (vr *ValidationResult) AddErrorIf(condition bool, message string) {
	if condition {
		vr.AddError(message)
	}
}

// ValidateString 验证字符串
func (vr *ValidationResult) ValidateString(value, fieldName string, required bool, minLength, maxLength int) {
	if required && (value == "" || strings.TrimSpace(value) == "") {
		vr.AddError(fmt.Sprintf("%s不能为空", fieldName))
		return
	}

	if !required && value == "" {
		return
	}

	length := len(strings.TrimSpace(value))
	if minLength > 0 && length < minLength {
		vr.AddError(fmt.Sprintf("%s长度不能少于%d个字符", fieldName, minLength))
	}

	if maxLength > 0 && length > maxLength {
		vr.AddError(fmt.Sprintf("%s长度不能超过%d个字符", fieldName, maxLength))
	}
}

// ValidateEmail 验证邮箱格式
func (vr *ValidationResult) ValidateEmail(value, fieldName string, required bool) {
	if required && (value == "" || strings.TrimSpace(value) == "") {
		vr.AddError(fmt.Sprintf("%s不能为空", fieldName))
		return
	}

	if !required && value == "" {
		return
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		vr.AddError(fmt.Sprintf("%s格式不正确", fieldName))
	}
}

// ValidatePassword 验证密码强度
func (vr *ValidationResult) ValidatePassword(value, fieldName string) {
	if value == "" {
		vr.AddError(fmt.Sprintf("%s不能为空", fieldName))
		return
	}

	if len(value) < 8 {
		vr.AddError(fmt.Sprintf("%s长度不能少于8个字符", fieldName))
	}

	if len(value) > 128 {
		vr.AddError(fmt.Sprintf("%s长度不能超过128个字符", fieldName))
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range value {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
		}
	}

	requiredTypes := 0
	if hasUpper {
		requiredTypes++
	}
	if hasLower {
		requiredTypes++
	}
	if hasDigit {
		requiredTypes++
	}
	if hasSpecial {
		requiredTypes++
	}

	if requiredTypes < 3 {
		vr.AddError(fmt.Sprintf("%s必须包含大写字母、小写字母、数字和特殊字符中的至少3种", fieldName))
	}

	// 检查常见弱密码
	weakPasswords := []string{
		"12345678", "password", "qwerty123", "abc123", "11111111",
		"87654321", "password123", "admin123", "letmein", "welcome",
		"qwerty123!", "qwerty", "123456", "123123", "admin", "root",
	}
	
	for _, weakPass := range weakPasswords {
		if strings.ToLower(value) == weakPass {
			vr.AddError(fmt.Sprintf("%s过于简单，请使用更复杂的密码", fieldName))
			break
		}
	}
}

// ValidateUsername 验证用户名
func (vr *ValidationResult) ValidateUsername(value, fieldName string) {
	if value == "" {
		vr.AddError(fmt.Sprintf("%s不能为空", fieldName))
		return
	}

	if len(value) < 3 {
		vr.AddError(fmt.Sprintf("%s长度不能少于3个字符", fieldName))
	}

	if len(value) > 50 {
		vr.AddError(fmt.Sprintf("%s长度不能超过50个字符", fieldName))
	}

	// 用户名只能包含字母、数字、下划线和点
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
	if !usernameRegex.MatchString(value) {
		vr.AddError(fmt.Sprintf("%s只能包含字母、数字、下划线和点", fieldName))
	}

	// 不能以点或下划线开头/结尾
	if strings.HasPrefix(value, ".") || strings.HasPrefix(value, "_") ||
		strings.HasSuffix(value, ".") || strings.HasSuffix(value, "_") {
		vr.AddError(fmt.Sprintf("%s不能以点或下划线开头或结尾", fieldName))
	}

	// 不能连续使用点或下划线
	if strings.Contains(value, "..") || strings.Contains(value, "__") {
		vr.AddError(fmt.Sprintf("%s不能连续使用点或下划线", fieldName))
	}
}

// ValidateInteger 验证整数
func (vr *ValidationResult) ValidateInteger(value string, fieldName string, required bool, min, max int64) {
	if required && value == "" {
		vr.AddError(fmt.Sprintf("%s不能为空", fieldName))
		return
	}

	if !required && value == "" {
		return
	}

	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		vr.AddError(fmt.Sprintf("%s必须是有效的整数", fieldName))
		return
	}

	if min != 0 && intValue < min {
		vr.AddError(fmt.Sprintf("%s不能小于%d", fieldName, min))
	}

	if max != 0 && intValue > max {
		vr.AddError(fmt.Sprintf("%s不能大于%d", fieldName, max))
	}
}

// ValidateFloat 验证浮点数
func (vr *ValidationResult) ValidateFloat(value string, fieldName string, required bool, min, max float64) {
	if required && value == "" {
		vr.AddError(fmt.Sprintf("%s不能为空", fieldName))
		return
	}

	if !required && value == "" {
		return
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		vr.AddError(fmt.Sprintf("%s必须是有效的数字", fieldName))
		return
	}

	if min != 0 && floatValue < min {
		vr.AddError(fmt.Sprintf("%s不能小于%f", fieldName, min))
	}

	if max != 0 && floatValue > max {
		vr.AddError(fmt.Sprintf("%s不能大于%f", fieldName, max))
	}
}

// ValidateEnum 验证枚举值
func (vr *ValidationResult) ValidateEnum(value, fieldName string, allowedValues []string, required bool) {
	if required && value == "" {
		vr.AddError(fmt.Sprintf("%s不能为空", fieldName))
		return
	}

	if !required && value == "" {
		return
	}

	found := false
	for _, allowed := range allowedValues {
		if value == allowed {
			found = true
			break
		}
	}

	if !found {
		vr.AddError(fmt.Sprintf("%s必须是以下值之一: %s", fieldName, strings.Join(allowedValues, ", ")))
	}
}

// ValidateRegex 验证正则表达式
func (vr *ValidationResult) ValidateRegex(value, fieldName, pattern string, required bool) {
	if required && value == "" {
		vr.AddError(fmt.Sprintf("%s不能为空", fieldName))
		return
	}

	if !required && value == "" {
		return
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		vr.AddError(fmt.Sprintf("%s验证规则配置错误", fieldName))
		return
	}

	if !regex.MatchString(value) {
		vr.AddError(fmt.Sprintf("%s格式不正确", fieldName))
	}
}

// ValidateNoSQLInjection 验证防止SQL注入
func (vr *ValidationResult) ValidateNoSQLInjection(value, fieldName string) {
	sqlInjectionPatterns := []string{
		`(?i)\b(SELECT|INSERT|UPDATE|DELETE|DROP|CREATE|ALTER|TRUNCATE|UNION|EXEC|EXECUTE|DECLARE|CAST)\b`,
		`(?i)\b(OR|AND)\s+\d+\s*=\s*\d+`,
		`(?i)\b(OR|AND)\s+['"][^'"]*['"]\s*=\s*['"][^'"]*['"]`,
		`(?i)\b(WAITFOR\s+DELAY|PG_SLEEP|SLEEP)\b`,
		`(?i)\b(XP_|SP_)\w+\b`,
		`(?i)[;'"\\]`,
		`(?i)\b(IF|THEN|ELSE|END|CASE|WHEN)\b`,
		`(?i)\b(INFORMATION_SCHEMA|SYS\.|MYSQL\.)\b`,
	}

	for _, pattern := range sqlInjectionPatterns {
		regex := regexp.MustCompile(pattern)
		if regex.MatchString(value) {
			vr.AddError(fmt.Sprintf("%s包含不安全字符或SQL注入尝试", fieldName))
			return
		}
	}
}

// ValidateNoXSS 验证防止XSS攻击
func (vr *ValidationResult) ValidateNoXSS(value, fieldName string) {
	xssPatterns := []string{
		`(?i)<script[^>]*>.*?</script>`,
		`(?i)<iframe[^>]*>.*?</iframe>`,
		`(?i)<object[^>]*>.*?</object>`,
		`(?i)<embed[^>]*>.*?</embed>`,
		`(?i)<applet[^>]*>.*?</applet>`,
		`(?i)<meta[^>]*>`,
		`(?i)<link[^>]*>`,
		`(?i)javascript:`,
		`(?i)vbscript:`,
		`(?i)on(load|click|mouseover|error|abort|focus|blur|submit|reset|change|select|keydown|keyup|keypress)\s*=`,
		`(?i)eval\s*\(`,
		`(?i)expression\s*\(`,
		`(?i)url\s*\(`,
	}

	for _, pattern := range xssPatterns {
		regex := regexp.MustCompile(pattern)
		if regex.MatchString(value) {
			vr.AddError(fmt.Sprintf("%s包含潜在的XSS攻击代码", fieldName))
			return
		}
	}
}

// ValidateFilePath 验证文件路径安全性
func (vr *ValidationResult) ValidateFilePath(value, fieldName string) {
	if value == "" {
		return
	}

	// 检查路径遍历攻击
	pathTraversalPatterns := []string{
		`\.\./`,  // ../
		`\.\.\\`, // ..\
		`~`,      // 家目录
		`/etc/`,  // 系统目录
		`C:\\`,   // Windows系统目录
		`\$\{`,   // 环境变量注入
	}

	for _, pattern := range pathTraversalPatterns {
		regex := regexp.MustCompile(pattern)
		if regex.MatchString(value) {
			vr.AddError(fmt.Sprintf("%s包含不安全的路径字符", fieldName))
			return
		}
	}
}

// SanitizeString 清理字符串输入
func SanitizeString(input string) string {
	// 移除危险字符
	result := strings.ReplaceAll(input, "'", "''")
	result = strings.ReplaceAll(result, "\"", "\"\"")
	result = strings.ReplaceAll(result, "\\", "\\\\")
	result = strings.ReplaceAll(result, "\x00", "") // 移除空字符
	
	// 移除控制字符
	var cleaned strings.Builder
	for _, char := range result {
		if char >= 32 || char == '\n' || char == '\r' || char == '\t' {
			cleaned.WriteRune(char)
		}
	}
	
	return cleaned.String()
}

// SanitizeForLike 为LIKE查询清理字符串
func SanitizeForLike(input string) string {
	// 转义LIKE通配符
	result := strings.ReplaceAll(input, "%", "\\%")
	result = strings.ReplaceAll(result, "_", "\\_")
	result = strings.ReplaceAll(result, "[", "\\[")
	result = strings.ReplaceAll(result, "]", "\\]")
	result = strings.ReplaceAll(result, "'", "''")
	
	return result
}