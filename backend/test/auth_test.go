package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"api/internal/config"
	"api/internal/logic"
	"api/internal/svc"
	"api/internal/types"
	"api/model"
	"api/pkg/auth"
	"api/pkg/database"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TestSuite 测试套件
type TestSuite struct {
	suite.Suite
	svcCtx   *svc.ServiceContext
	testDB   *sql.DB
	testUser *model.VtUsers
}

// SetupSuite 测试套件初始化
func (s *TestSuite) SetupSuite() {
	// 创建测试配置
	cfg := config.Config{
		MySQL: config.MySQLConfig{
			Host:         "localhost",
			Port:         3306,
			User:         "root",
			Password:     "root",
			DBName:       "volctraindb_test",
			Charset:      "utf8mb4",
			ParseTime:    true,
			Loc:          "Asia/Shanghai",
			MaxOpenConns: 10,
			MaxIdleConns: 5,
			MaxLifetime:  3600,
		},
		Redis: config.RedisConfig{
			Host:         "localhost",
			Port:         6379,
			Password:     "",
			DB:           1, // 使用测试数据库
			PoolSize:     10,
			MinIdleConns: 2,
		},
		Auth: config.AuthConfig{
			AccessSecret:  "test_access_secret_key_minimum_64_characters_required_for_security",
			AccessExpire:  3600,
			RefreshExpire: 604800,
		},
	}

	// 初始化数据库连接
	db, err := database.NewMySQLConnection(cfg.MySQL)
	s.Require().NoError(err)
	s.testDB = db

	// 初始化服务上下文
	s.svcCtx = &svc.ServiceContext{
		Config:        cfg,
		DB:            db,
		VtUsersModel:  model.NewVtUsersModel(db),
		VtRolesModel:  model.NewVtRolesModel(db),
		// 其他模型...
	}

	// 创建测试数据
	s.setupTestData()
}

// TearDownSuite 测试套件清理
func (s *TestSuite) TearDownSuite() {
	// 清理测试数据
	s.cleanupTestData()
	if s.testDB != nil {
		s.testDB.Close()
	}
}

// setupTestData 创建测试数据
func (s *TestSuite) setupTestData() {
	// 创建测试用户
	hashedPassword, err := auth.HashPassword("testpass123")
	s.Require().NoError(err)

	s.testUser = &model.VtUsers{
		Username:      "testuser_" + time.Now().Format("20060102150405"),
		Email:         "test@example.com",
		Password:      hashedPassword,
		Salt:          "test_salt",
		RealName:      "Test User",
		Nickname:      "TestNick",
		Department:    "Test Dept",
		Status:        "active",
		UserType:      "user",
		EmailVerified: true,
		PhoneVerified: false,
	}

	result, err := s.svcCtx.VtUsersModel.Insert(context.Background(), s.testUser)
	s.Require().NoError(err)

	id, err := result.LastInsertId()
	s.Require().NoError(err)
	s.testUser.Id = id
}

// cleanupTestData 清理测试数据
func (s *TestSuite) cleanupTestData() {
	if s.testUser != nil {
		s.svcCtx.VtUsersModel.Delete(context.Background(), s.testUser.Id)
	}
}

// TestAuthLogic 测试认证逻辑
func (s *TestSuite) TestAuthLogic() {
	ctx := context.Background()

	// 测试登录
	s.Run("Login", func() {
		loginLogic := logic.NewLoginLogic(ctx, s.svcCtx)
		
		// 测试成功登录
		req := &types.LoginReq{
			Username: s.testUser.Username,
			Password: "testpass123",
		}
		
		resp, err := loginLogic.Login(req)
		s.NoError(err)
		s.NotNil(resp)
		s.NotEmpty(resp.AccessToken)
		s.NotEmpty(resp.RefreshToken)
		s.Equal("Bearer", resp.TokenType)
		s.Equal(s.testUser.Id, resp.UserInfo.ID)
		
		// 测试错误密码
		req.Password = "wrongpassword"
		resp, err = loginLogic.Login(req)
		s.Error(err)
		s.Nil(resp)
		
		// 测试不存在的用户
		req.Username = "nonexistent"
		req.Password = "testpass123"
		resp, err = loginLogic.Login(req)
		s.Error(err)
		s.Nil(resp)
	})

	// 测试刷新Token
	s.Run("RefreshToken", func() {
		// 先登录获取token
		loginLogic := logic.NewLoginLogic(ctx, s.svcCtx)
		loginResp, err := loginLogic.Login(&types.LoginReq{
			Username: s.testUser.Username,
			Password: "testpass123",
		})
		s.Require().NoError(err)
		s.Require().NotNil(loginResp)
		
		// 测试刷新token
		refreshLogic := logic.NewRefreshTokenLogic(ctx, s.svcCtx)
		refreshReq := &types.RefreshTokenReq{
			RefreshToken: loginResp.RefreshToken,
		}
		
		refreshResp, err := refreshLogic.RefreshToken(refreshReq)
		s.NoError(err)
		s.NotNil(refreshResp)
		s.NotEmpty(refreshResp.AccessToken)
		s.NotEmpty(refreshResp.RefreshToken)
		
		// 测试无效token
		refreshReq.RefreshToken = "invalid_token"
		refreshResp, err = refreshLogic.RefreshToken(refreshReq)
		s.Error(err)
		s.Nil(refreshResp)
	})

	// 测试获取用户信息
	s.Run("GetUserInfo", func() {
		// 创建带用户ID的上下文 (使用json.Number类型以匹配JWT中间件)
		ctxWithUser := context.WithValue(ctx, "userId", json.Number(fmt.Sprintf("%d", s.testUser.Id)))
		
		getUserInfoLogic := logic.NewGetUserInfoLogic(ctxWithUser, s.svcCtx)
		
		resp, err := getUserInfoLogic.GetUserInfo(&types.EmptyReq{})
		s.NoError(err)
		s.NotNil(resp)
		s.Equal(s.testUser.Id, resp.UserInfo.ID)
		s.Equal(s.testUser.Username, resp.UserInfo.Username)
		s.Equal(s.testUser.Email, resp.UserInfo.Email)
	})

	// 测试获取权限码
	s.Run("GetAccessCodes", func() {
		// 创建带用户ID的上下文 (使用json.Number类型以匹配JWT中间件)
		ctxWithUser := context.WithValue(ctx, "userId", json.Number(fmt.Sprintf("%d", s.testUser.Id)))
		
		getAccessCodesLogic := logic.NewGetAccessCodesLogic(ctxWithUser, s.svcCtx)
		
		resp, err := getAccessCodesLogic.GetAccessCodes(&types.EmptyReq{})
		s.NoError(err)
		s.NotNil(resp)
		s.NotEmpty(resp.Codes)
		// 普通用户应该有基础权限
		s.Contains(resp.Codes, "training:job:read")
	})
}

// TestPasswordSecurity 测试密码安全
func (s *TestSuite) TestPasswordSecurity() {
	s.Run("PasswordHashing", func() {
		password := "testpassword123"
		
		// 测试密码加密
		hashedPassword, err := auth.HashPassword(password)
		s.NoError(err)
		s.NotEmpty(hashedPassword)
		s.NotEqual(password, hashedPassword)
		
		// 测试密码验证
		isValid := auth.CheckPassword(password, hashedPassword)
		s.True(isValid)
		
		// 测试错误密码
		isValid = auth.CheckPassword("wrongpassword", hashedPassword)
		s.False(isValid)
	})

	s.Run("PasswordValidation", func() {
		// 测试弱密码
		isValid, errors := auth.ValidatePassword("123", auth.DefaultPasswordRule)
		s.False(isValid)
		s.NotEmpty(errors)
		
		// 测试强密码
		isValid, errors = auth.ValidatePassword("StrongPass123!", auth.StrongPasswordRule)
		s.True(isValid)
		s.Empty(errors)
	})
}

// TestDatabaseOperations 测试数据库操作
func (s *TestSuite) TestDatabaseOperations() {
	ctx := context.Background()

	s.Run("UserCRUD", func() {
		// 测试创建用户
		hashedPassword, err := auth.HashPassword("newuserpass")
		s.Require().NoError(err)

		newUser := &model.VtUsers{
			Username:      "newuser_" + time.Now().Format("20060102150405"),
			Email:         "newuser@example.com",
			Password:      hashedPassword,
			Salt:          "new_salt",
			RealName:      "New User",
			Status:        "active",
			UserType:      "user",
			EmailVerified: true,
		}

		// 插入
		result, err := s.svcCtx.VtUsersModel.Insert(ctx, newUser)
		s.NoError(err)
		s.NotNil(result)

		id, err := result.LastInsertId()
		s.NoError(err)
		s.Greater(id, int64(0))
		newUser.Id = id

		// 查询
		foundUser, err := s.svcCtx.VtUsersModel.FindOne(ctx, id)
		s.NoError(err)
		s.NotNil(foundUser)
		s.Equal(newUser.Username, foundUser.Username)
		s.Equal(newUser.Email, foundUser.Email)

		// 通过用户名查询
		foundByUsername, err := s.svcCtx.VtUsersModel.FindByUsername(ctx, newUser.Username)
		s.NoError(err)
		s.NotNil(foundByUsername)
		s.Equal(newUser.Id, foundByUsername.Id)

		// 更新
		foundUser.RealName = "Updated Name"
		err = s.svcCtx.VtUsersModel.Update(ctx, foundUser)
		s.NoError(err)

		// 验证更新
		updatedUser, err := s.svcCtx.VtUsersModel.FindOne(ctx, id)
		s.NoError(err)
		s.Equal("Updated Name", updatedUser.RealName)

		// 删除
		err = s.svcCtx.VtUsersModel.Delete(ctx, id)
		s.NoError(err)

		// 验证删除（软删除，应该查询不到）
		deletedUser, err := s.svcCtx.VtUsersModel.FindOne(ctx, id)
		s.Error(err)
		s.Nil(deletedUser)
	})
}

// TestJWTSecurity 测试JWT安全性
func (s *TestSuite) TestJWTSecurity() {
	s.Run("TokenGeneration", func() {
		userID := s.testUser.Id
		secret := s.svcCtx.Config.Auth.AccessSecret
		expire := s.svcCtx.Config.Auth.AccessExpire

		// 生成token
		token, err := auth.GenerateToken(userID, secret, expire)
		s.NoError(err)
		s.NotEmpty(token)

		// 验证token
		parsedUserID, err := auth.ValidateToken(token, secret)
		s.NoError(err)
		s.Equal(userID, parsedUserID)

		// 测试错误密钥
		_, err = auth.ValidateToken(token, "wrong_secret")
		s.Error(err)

		// 测试无效token
		_, err = auth.ValidateToken("invalid_token", secret)
		s.Error(err)
	})
}

// TestSuite 运行测试套件
func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

// 基准测试
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

func BenchmarkTokenGeneration(b *testing.B) {
	userID := int64(123)
	secret := "test_secret_key_minimum_64_characters_required_for_security_benchmark"
	expire := int64(3600)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := auth.GenerateToken(userID, secret, expire)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTokenValidation(b *testing.B) {
	userID := int64(123)
	secret := "test_secret_key_minimum_64_characters_required_for_security_benchmark"
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