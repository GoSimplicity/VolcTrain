package test

import (
	"testing"

	"api/pkg/auth"
	"api/pkg/validation"
)

// BenchmarkInputValidation 输入验证基准测试
func BenchmarkInputValidation(b *testing.B) {
	vr := validation.NewValidationResult()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vr.ValidateUsername("testuser123", "用户名")
		vr.ValidatePassword("StrongPass123!", "密码")
		vr.ValidateNoSQLInjection("normal input", "输入")
		vr.ValidateNoXSS("normal text", "输入")
	}
}

// BenchmarkPasswordHashing 密码加密基准测试
func BenchmarkPasswordHashing(b *testing.B) {
	password := "benchmarkpassword123"
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := auth.HashPassword(password)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPasswordVerification 密码验证基准测试
func BenchmarkPasswordVerification(b *testing.B) {
	password := "benchmarkpassword123"
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		auth.CheckPassword(password, hashedPassword)
	}
}

// BenchmarkAccessTokenGeneration 访问Token生成基准测试
func BenchmarkAccessTokenGeneration(b *testing.B) {
	userID := int64(123)
	secret := "test_access_secret_key_minimum_64_characters_required_for_security_benchmark"
	expire := int64(3600)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := auth.GenerateToken(userID, secret, expire)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkAccessTokenValidation 访问Token验证基准测试
func BenchmarkAccessTokenValidation(b *testing.B) {
	userID := int64(123)
	secret := "test_access_secret_key_minimum_64_characters_required_for_security_benchmark"
	expire := int64(3600)
	
	token, err := auth.GenerateToken(userID, secret, expire)
	if err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := auth.ValidateToken(token, secret)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRefreshTokenGeneration 刷新Token生成基准测试
func BenchmarkRefreshTokenGeneration(b *testing.B) {
	userID := int64(123)
	secret := "test_refresh_secret_key_minimum_64_characters_required_for_security_benchmark"
	expire := int64(604800)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := auth.GenerateToken(userID, secret, expire)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRefreshTokenValidation 刷新Token验证基准测试
func BenchmarkRefreshTokenValidation(b *testing.B) {
	userID := int64(123)
	secret := "test_refresh_secret_key_minimum_64_characters_required_for_security_benchmark"
	expire := int64(604800)
	
	token, err := auth.GenerateToken(userID, secret, expire)
	if err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := auth.ValidateToken(token, secret)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSQLInjectionDetection SQL注入检测基准测试
func BenchmarkSQLInjectionDetection(b *testing.B) {
	vr := validation.NewValidationResult()
	
	inputs := []string{
		"normal input",
		"SELECT * FROM users",
		"1' OR '1'='1",
		"; DROP TABLE users; --",
		"<script>alert('xss')</script>",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, input := range inputs {
			vr.ValidateNoSQLInjection(input, "输入")
			vr.ValidateNoXSS(input, "输入")
		}
	}
}

// BenchmarkUsernameValidation 用户名验证基准测试
func BenchmarkUsernameValidation(b *testing.B) {
	vr := validation.NewValidationResult()
	
	username := "valid_user123"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vr.ValidateUsername(username, "用户名")
	}
}

// BenchmarkComplexValidation 复杂验证基准测试
func BenchmarkComplexValidation(b *testing.B) {
	vr := validation.NewValidationResult()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		vr.ValidateUsername("testuser_" + string(rune(i%26+65)), "用户名")
		vr.ValidatePassword("ComplexPass123!@#", "密码")
		vr.ValidateEmail("test@example.com", "邮箱", false)
		vr.ValidateNoSQLInjection("safe input", "输入")
		vr.ValidateNoXSS("safe text", "输入")
	}
}