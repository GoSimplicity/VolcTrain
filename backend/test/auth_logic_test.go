package test

import (
	"context"
	"strings"
	"testing"

	"api/internal/config"
	"api/internal/logic"
	"api/internal/svc"
	"api/internal/types"
	"api/pkg/auth"
	"api/pkg/errors"
	"api/pkg/validation"

	"github.com/stretchr/testify/suite"
)

// TestAuthLogicSuite 认证逻辑测试套件
type TestAuthLogicSuite struct {
	suite.Suite
	svcCtx *svc.ServiceContext
}

// SetupSuite 测试套件初始化
func (s *TestAuthLogicSuite) SetupSuite() {
	// 创建测试配置（不依赖数据库）
	cfg := config.Config{
		Auth: config.AuthConfig{
			AccessSecret:   "test_access_secret_key_minimum_64_characters_required_for_security",
			RefreshSecret:  "test_refresh_secret_key_minimum_64_characters_required_for_security_different",
			AccessExpire:   3600,
			RefreshExpire:  604800,
		},
	}

	// 创建服务上下文（不包含数据库连接）
	s.svcCtx = &svc.ServiceContext{
		Config: cfg,
	}
}

// TestLoginInputValidation 测试登录输入验证
func (s *TestAuthLogicSuite) TestLoginInputValidation() {
	ctx := context.Background()

	s.Run("EmptyUsername", func() {
		loginLogic := logic.NewLoginLogic(ctx, s.svcCtx)
		req := &types.LoginReq{
			Username: "",
			Password: "testpass123",
		}

		resp, err := loginLogic.Login(req)
		s.Error(err)
		s.Nil(resp)
		
		// 检查错误类型
		bizErr := errors.GetBizError(err)
		s.NotNil(bizErr)
		s.Equal(errors.ErrCodeValidation, bizErr.Code)
	})

	s.Run("EmptyPassword", func() {
		loginLogic := logic.NewLoginLogic(ctx, s.svcCtx)
		req := &types.LoginReq{
			Username: "testuser",
			Password: "",
		}

		resp, err := loginLogic.Login(req)
		s.Error(err)
		s.Nil(resp)
		
		bizErr := errors.GetBizError(err)
		s.NotNil(bizErr)
		s.Equal(errors.ErrCodeValidation, bizErr.Code)
	})

	s.Run("WeakPassword", func() {
		loginLogic := logic.NewLoginLogic(ctx, s.svcCtx)
		req := &types.LoginReq{
			Username: "testuser",
			Password: "123",
		}

		resp, err := loginLogic.Login(req)
		s.Error(err)
		s.Nil(resp)
		
		bizErr := errors.GetBizError(err)
		s.NotNil(bizErr)
		s.Equal(errors.ErrCodeValidation, bizErr.Code)
	})

	s.Run("InvalidUsername", func() {
		loginLogic := logic.NewLoginLogic(ctx, s.svcCtx)
		req := &types.LoginReq{
			Username: "test user; DROP TABLE users;", // SQL注入尝试
			Password: "testpass123",
		}

		resp, err := loginLogic.Login(req)
		s.Error(err)
		s.Nil(resp)
		
		bizErr := errors.GetBizError(err)
		s.NotNil(bizErr)
		s.Equal(errors.ErrCodeValidation, bizErr.Code)
	})

	s.Run("SQLInjectionInPassword", func() {
		loginLogic := logic.NewLoginLogic(ctx, s.svcCtx)
		req := &types.LoginReq{
			Username: "testuser",
			Password: "testpass123'; DROP TABLE users; --",
		}

		resp, err := loginLogic.Login(req)
		s.Error(err)
		s.Nil(resp)
		
		bizErr := errors.GetBizError(err)
		s.NotNil(bizErr)
		s.Equal(errors.ErrCodeValidation, bizErr.Code)
	})
}

// TestValidator 测试验证器功能
func (s *TestAuthLogicSuite) TestValidator() {
	s.Run("UsernameValidation", func() {
		tests := []struct {
			username string
			valid    bool
		}{
			{"validuser", true},
			{"valid_user", true},
			{"valid.user", true},
			{"valid_user123", true},
			{"", false},
			{"ab", false}, // 太短
			{strings.Repeat("a", 51), false}, // 太长
			{"user name", false}, // 包含空格
			{"user@name", false}, // 包含特殊字符
			{"user;name", false}, // SQL注入字符
			{"--username", false}, // 以特殊字符开头
			{"username--", false}, // 以特殊字符结尾
			{"user__name", false}, // 连续下划线
			{"user..name", false}, // 连续点
		}

		for _, test := range tests {
			vr := validation.NewValidationResult()
			vr.ValidateUsername(test.username, "用户名")
			
			if test.valid {
				s.True(vr.IsValid, "用户名 %s 应该有效", test.username)
			} else {
				s.False(vr.IsValid, "用户名 %s 应该无效", test.username)
			}
		}
	})

	s.Run("PasswordValidation", func() {
		tests := []struct {
			password string
			valid    bool
		}{
			{"StrongPass123!", true},
			{"MyP@ssw0rd", true},
			{"Complex123!@#", true},
			{"ABCDEFGH123!", true}, // 包含大写字母、数字、特殊字符
			{"abcdefgh123!", true}, // 包含小写字母、数字、特殊字符
			{"ABCDefgh123", true}, // 包含大写字母、小写字母、数字
			{"", false},
			{"simple", false}, // 太简单
			{"12345678", false}, // 只有数字
			{"password", false}, // 常见弱密码
			{"abcdefgh", false}, // 只有字母
			{"ABCDefgh!", true}, // 包含大写字母、小写字母、特殊字符（满足3种类型要求）
			{"password123", false}, // 常见弱密码
			{"qwerty123!", false}, // 常见弱密码（qwerty是常见键盘模式）
		}

		for _, test := range tests {
			vr := validation.NewValidationResult()
			vr.ValidatePassword(test.password, "密码")
			
			if test.valid {
				s.True(vr.IsValid, "密码 %s 应该有效", test.password)
			} else {
				s.False(vr.IsValid, "密码 %s 应该无效", test.password)
			}
		}
	})

	s.Run("SQLInjectionDetection", func() {
		tests := []struct {
			input string
			valid bool
		}{
			{"normal input", true},
			{"user123", true},
			{"SELECT * FROM users", false},
			{"1' OR '1'='1", false},
			{"; DROP TABLE users; --", false},
			{"admin'--", false},
			{"1; INSERT INTO users VALUES", false},
			{"UNION SELECT username, password FROM users", false},
			{"' OR 1=1 --", false},
			{"<script>alert('xss')</script>", false},
			{"javascript:alert('xss')", false},
		}

		for _, test := range tests {
			vr := validation.NewValidationResult()
			vr.ValidateNoSQLInjection(test.input, "输入")
			
			if test.valid {
				s.True(vr.IsValid, "输入 '%s' 应该有效", test.input)
			} else {
				s.False(vr.IsValid, "输入 '%s' 应该被拒绝（SQL注入）", test.input)
			}
		}
	})

	s.Run("XSSDetection", func() {
		tests := []struct {
			input string
			valid bool
		}{
			{"normal text", true},
			{"user description", true},
			{"<script>alert('xss')</script>", false},
			{"<img src=x onerror=alert('xss')>", false},
			{"javascript:alert('xss')", false},
			{"<iframe src='javascript:alert(\"xss\")'></iframe>", false},
			{"<div onclick='alert(\"xss\")'>click</div>", false},
			{"onload=alert('xss')", false},
			{"eval('alert(\"xss\")')", false},
			{"expression(alert('xss'))", false},
		}

		for _, test := range tests {
			vr := validation.NewValidationResult()
			vr.ValidateNoXSS(test.input, "输入")
			
			if test.valid {
				s.True(vr.IsValid, "输入 '%s' 应该有效", test.input)
			} else {
				s.False(vr.IsValid, "输入 '%s' 应该被拒绝（XSS攻击）", test.input)
			}
		}
	})
}

// TestJWTSecurity 测试JWT安全性
func (s *TestAuthLogicSuite) TestJWTSecurity() {
	s.Run("TokenGeneration", func() {
		userID := int64(123)
		accessSecret := s.svcCtx.Config.Auth.AccessSecret
		refreshSecret := s.svcCtx.Config.Auth.RefreshSecret
		accessExpire := s.svcCtx.Config.Auth.AccessExpire
		refreshExpire := s.svcCtx.Config.Auth.RefreshExpire

		// 生成访问token
		accessToken, err := auth.GenerateToken(userID, accessSecret, accessExpire)
		s.NoError(err)
		s.NotEmpty(accessToken)

		// 生成刷新token（使用不同的密钥）
		refreshToken, err := auth.GenerateToken(userID, refreshSecret, refreshExpire)
		s.NoError(err)
		s.NotEmpty(refreshToken)

		// 验证访问token
		parsedUserID, err := auth.ValidateToken(accessToken, accessSecret)
		s.NoError(err)
		s.Equal(userID, parsedUserID)

		// 验证刷新token
		parsedUserID, err = auth.ValidateToken(refreshToken, refreshSecret)
		s.NoError(err)
		s.Equal(userID, parsedUserID)

		// 测试错误密钥（访问token不能用刷新密钥验证）
		_, err = auth.ValidateToken(accessToken, refreshSecret)
		s.Error(err)

		// 测试错误密钥（刷新token不能用访问密钥验证）
		_, err = auth.ValidateToken(refreshToken, accessSecret)
		s.Error(err)

		// 测试无效token
		_, err = auth.ValidateToken("invalid_token", accessSecret)
		s.Error(err)
	})
}

// TestErrorHandling 测试错误处理
func (s *TestAuthLogicSuite) TestErrorHandling() {
	s.Run("ErrorTypes", func() {
		// 测试验证错误
		validationErr := errors.NewValidationError("输入验证失败")
		s.NotNil(validationErr)
		s.Equal(errors.ErrCodeValidation, validationErr.Code)
		s.Equal(errors.ErrorTypeValidation, validationErr.Type)

		// 测试认证错误
		authErr := errors.NewAuthError("认证失败")
		s.NotNil(authErr)
		s.Equal(errors.ErrCodeUnauthorized, authErr.Code)
		s.Equal(errors.ErrorTypeAuth, authErr.Type)

		// 测试业务错误
		businessErr := errors.NewBusinessError(errors.ErrCodeUserNotFound, "用户不存在")
		s.NotNil(businessErr)
		s.Equal(errors.ErrCodeUserNotFound, businessErr.Code)
		s.Equal(errors.ErrorTypeBusiness, businessErr.Type)

		// 测试系统错误
		systemErr := errors.NewInternalError("系统内部错误")
		s.NotNil(systemErr)
		s.Equal(errors.ErrCodeInternalError, systemErr.Code)
		s.Equal(errors.ErrorTypeSystem, systemErr.Type)
	})

	s.Run("ErrorHTTPStatus", func() {
		tests := []struct {
			err           *errors.BizError
			expectedStatus int
		}{
			{errors.ErrInvalidParam, 400},
			{errors.ErrTokenInvalid, 401},
			{errors.ErrPermissionDenied, 403},
			{errors.ErrUserNotFound, 404},
			{errors.ErrDuplicateData, 409},
			{errors.NewBizError(429, "请求过多", "VALIDATION_ERROR"), 429},
			{errors.ErrDatabaseError, 500},
		}

		for _, test := range tests {
			status := test.err.GetHTTPStatus()
			s.Equal(test.expectedStatus, status, "错误码 %d 应该返回HTTP状态码 %d", test.err.Code, test.expectedStatus)
		}
	})
}

// RunAuthLogicTests 运行认证逻辑测试
func TestRunAuthLogicTests(t *testing.T) {
	suite.Run(t, new(TestAuthLogicSuite))
}

