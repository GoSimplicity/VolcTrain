package auth

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

// BenchmarkPasswordHashing 密码哈希性能基准测试
func BenchmarkPasswordHashing(b *testing.B) {
	password := "testpassword123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			b.Fatalf("密码哈希失败: %v", err)
		}
	}
}

// BenchmarkPasswordVerification 密码验证性能基准测试
func BenchmarkPasswordVerification(b *testing.B) {
	password := "testpassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
		if err != nil {
			b.Fatal("密码验证失败")
		}
	}
}

// TestPasswordSecurity 密码安全测试
func TestPasswordSecurity(t *testing.T) {
	// 测试弱密码处理
	weakPasswords := []string{"123", "password", "admin", ""}

	for _, weak := range weakPasswords {
		// 即使是弱密码，加密功能也应该正常工作
		hashed, err := bcrypt.GenerateFromPassword([]byte(weak), bcrypt.DefaultCost)
		if err != nil {
			t.Errorf("密码加密失败: %v", err)
			continue
		}

		if weak == string(hashed) {
			t.Error("加密后密码不应与原密码相同")
		}

		if len(hashed) <= 20 {
			t.Error("加密后密码长度应该足够长")
		}
	}

	// 测试相同密码生成不同哈希
	password := "samepassword"
	hash1, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hash2, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if string(hash1) == string(hash2) {
		t.Error("相同密码应该生成不同的哈希值")
	}

	t.Log("✅ 密码安全测试通过")
}
