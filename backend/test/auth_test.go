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
	"github.com/stretchr/testify/suite"
)

// TestSuite 测试套件
type TestSuite struct {
	suite.Suite
	svcCtx   *svc.ServiceContext
	testDB   *sql.DB
	testUser *model.VtUsersSimple
}

// SetupSuite 测试套件初始化
func (s *TestSuite) SetupSuite() {
	// 创建测试配置
	cfg := config.Config{
		MySQL: config.MySQLConfig{
			Host:         "localhost",
			Port:         3306,
			User:         "root",
			Password:     "",
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

	// 初始化数据库连接（不可用时跳过集成测试）
	db, err := database.NewMySQLConnection(cfg.MySQL)
	if err != nil {
		s.T().Skipf("Skipping integration tests: MySQL not available (%v)", err)
		return
	}
	s.testDB = db

	// 初始化服务上下文
	s.svcCtx = &svc.ServiceContext{
		Config:       cfg,
		DB:           db,
		VtUsersModel: model.NewVtUsersSimpleModel(db),
		VtRolesModel: model.NewVtRolesModel(db),
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

	username := "testuser_" + time.Now().Format("20060102150405")

	// 直接使用SQL插入测试用户（简化模型可能不支持完整的CRUD）
	query := `INSERT INTO vt_users (username, email, password_hash, real_name, department, status, user_type, email_verified, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`
	result, err := s.testDB.ExecContext(context.Background(), query,
		username, "test@example.com", hashedPassword, "Test User", "Test Dept", "active", "user", true)
	s.Require().NoError(err)

	id, err := result.LastInsertId()
	s.Require().NoError(err)

	s.testUser = &model.VtUsersSimple{
		Id:          id,
		Username:    username,
		Email:       "test@example.com",
		Password:    hashedPassword,
		RealName:    "Test User",
		Department:  "Test Dept",
		Status:      "active",
		UserType:    "user",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		LastLoginAt: time.Now(),
	}
}

// cleanupTestData 清理测试数据
func (s *TestSuite) cleanupTestData() {
	if s.testUser != nil {
		// 直接使用SQL删除测试用户
		_, err := s.testDB.ExecContext(context.Background(), "DELETE FROM vt_users WHERE id = ?", s.testUser.Id)
		s.Require().NoError(err)
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
		// 测试创建用户（简化版，主要测试查询功能）
		hashedPassword, err := auth.HashPassword("newuserpass")
		s.Require().NoError(err)

		username := "newuser_" + time.Now().Format("20060102150405")

		// 直接插入测试用户
		query := `INSERT INTO vt_users (username, email, password_hash, real_name, status, user_type, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`
		result, err := s.testDB.ExecContext(ctx, query, username, "newuser@example.com", hashedPassword, "New User", "active", "user")
		s.NoError(err)

		id, err := result.LastInsertId()
		s.NoError(err)
		s.Greater(id, int64(0))

		// 测试查询功能
		foundUser, err := s.svcCtx.VtUsersModel.FindOne(ctx, id)
		s.NoError(err)
		s.NotNil(foundUser)
		s.Equal(username, foundUser.Username)
		s.Equal("newuser@example.com", foundUser.Email)

		// 测试通过用户名查询
		foundByUsername, err := s.svcCtx.VtUsersModel.FindByUsername(ctx, username)
		s.NoError(err)
		s.NotNil(foundByUsername)
		s.Equal(id, foundByUsername.Id)

		// 测试更新功能
		updatedUser := &model.VtUsersSimple{
			Id:         id,
			Username:   username,
			Email:      "newuser@example.com",
			Password:   hashedPassword,
			RealName:   "Updated Name",
			Department: "Test Dept",
			Status:     "active",
			UserType:   "user",
		}
		err = s.svcCtx.VtUsersModel.Update(ctx, updatedUser)
		s.NoError(err)

		// 验证更新
		retrievedUser, err := s.svcCtx.VtUsersModel.FindOne(ctx, id)
		s.NoError(err)
		s.Equal("Updated Name", retrievedUser.RealName)

		// 清理测试数据
		_, err = s.testDB.ExecContext(ctx, "DELETE FROM vt_users WHERE id = ?", id)
		s.NoError(err)
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
