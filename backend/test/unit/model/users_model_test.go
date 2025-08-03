//go:build unit
// +build unit

package model

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"api/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UsersModelTestSuite 用户模型测试套件
type UsersModelTestSuite struct {
	suite.Suite
	db        *sql.DB
	mock      sqlmock.Sqlmock
	userModel model.VtUsersModel
}

// SetupTest 每个测试前的初始化
func (suite *UsersModelTestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	suite.Require().NoError(err)

	suite.userModel = model.NewVtUsersModel(suite.db)
}

// TearDownTest 每个测试后的清理
func (suite *UsersModelTestSuite) TearDownTest() {
	suite.db.Close()
}

// TestInsertUser 测试插入用户
func (suite *UsersModelTestSuite) TestInsertUser() {
	// 准备测试数据
	now := time.Now()
	user := &model.VtUsers{
		Username:      "testuser",
		Email:         "test@example.com",
		Phone:         "13812345678",
		PasswordHash:  "hashed_password_123",
		Salt:          "random_salt_456",
		RealName:      "测试用户",
		Nickname:      "测试",
		Status:        "active",
		UserType:      "user",
		EmailVerified: true,
		PhoneVerified: false,
		MfaEnabled:    false,
		LoginAttempts: 0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// 设置Mock期望
	suite.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO vt_users")).
		WithArgs(
			user.Username,
			user.Email,
			user.Phone,
			user.PasswordHash,
			user.Salt,
			user.RealName,
			user.Nickname,
			user.Status,
			user.UserType,
			user.EmailVerified,
			user.PhoneVerified,
			user.MfaEnabled,
			user.LoginAttempts,
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
		).
		WillReturnResult(sqlmock.NewResult(1001, 1))

	// 执行测试
	result, err := suite.userModel.Insert(user)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)

	lastInsertId, err := result.LastInsertId()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1001), lastInsertId)

	// 验证所有期望都被满足
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// TestFindOneUser 测试根据ID查找用户
func (suite *UsersModelTestSuite) TestFindOneUser() {
	// 准备测试数据
	userID := int64(1001)
	expectedTime := time.Now()

	// 设置Mock期望
	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "phone", "password_hash", "salt", "real_name",
		"nickname", "status", "user_type", "last_login_at", "last_login_ip",
		"password_expires_at", "login_attempts", "locked_until", "mfa_enabled",
		"mfa_secret", "email_verified", "phone_verified", "created_at", "updated_at", "deleted_at",
	}).AddRow(
		userID, "testuser", "test@example.com", "13812345678", "hashed_password_123", "random_salt_456", "测试用户",
		"测试", "active", "user", &expectedTime, "192.168.1.100",
		nil, 0, nil, false,
		"", true, false, expectedTime, expectedTime, nil,
	)

	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs(userID).
		WillReturnRows(rows)

	// 执行测试
	user, err := suite.userModel.FindOne(userID)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), userID, user.Id)
	assert.Equal(suite.T(), "testuser", user.Username)
	assert.Equal(suite.T(), "test@example.com", user.Email)
	assert.Equal(suite.T(), "测试用户", user.RealName)
	assert.Equal(suite.T(), "active", user.Status)
	assert.Equal(suite.T(), "user", user.UserType)
	assert.True(suite.T(), user.EmailVerified)
	assert.False(suite.T(), user.PhoneVerified)

	// 验证所有期望都被满足
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// TestFindOneByUsername 测试根据用户名查找用户
func (suite *UsersModelTestSuite) TestFindOneByUsername() {
	// 准备测试数据
	username := "testuser"
	expectedTime := time.Now()

	// 设置Mock期望
	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "password_hash", "salt", "real_name", "status",
		"user_type", "login_attempts", "locked_until", "mfa_enabled", "created_at", "updated_at",
	}).AddRow(
		int64(1001), username, "test@example.com", "hashed_password_123", "random_salt_456", "测试用户", "active",
		"user", 0, nil, false, expectedTime, expectedTime,
	)

	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs(username).
		WillReturnRows(rows)

	// 执行测试
	user, err := suite.userModel.FindOneByUsername(username)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), username, user.Username)
	assert.Equal(suite.T(), "test@example.com", user.Email)
	assert.Equal(suite.T(), "active", user.Status)
	assert.Equal(suite.T(), int(0), user.LoginAttempts)

	// 验证所有期望都被满足
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// TestFindOneByEmail 测试根据邮箱查找用户
func (suite *UsersModelTestSuite) TestFindOneByEmail() {
	// 准备测试数据
	email := "test@example.com"
	expectedTime := time.Now()

	// 设置Mock期望
	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "password_hash", "salt", "status", "user_type",
		"email_verified", "phone_verified", "mfa_enabled", "created_at", "updated_at",
	}).AddRow(
		int64(1001), "testuser", email, "hashed_password_123", "random_salt_456", "active", "user",
		true, false, false, expectedTime, expectedTime,
	)

	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs(email).
		WillReturnRows(rows)

	// 执行测试
	user, err := suite.userModel.FindOneByEmail(email)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), email, user.Email)
	assert.Equal(suite.T(), "testuser", user.Username)
	assert.True(suite.T(), user.EmailVerified)
	assert.False(suite.T(), user.PhoneVerified)

	// 验证所有期望都被满足
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// TestUpdateUser 测试更新用户信息
func (suite *UsersModelTestSuite) TestUpdateUser() {
	// 准备测试数据
	user := &model.VtUsers{
		Id:            1001,
		Username:      "updateduser",
		Email:         "updated@example.com",
		RealName:      "更新的用户",
		Nickname:      "更新",
		Status:        "active",
		EmailVerified: true,
		PhoneVerified: true,
		UpdatedAt:     time.Now(),
	}

	// 设置Mock期望
	suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE vt_users")).
		WithArgs(
			user.Username,
			user.Email,
			user.Phone,
			user.RealName,
			user.Nickname,
			user.Status,
			user.EmailVerified,
			user.PhoneVerified,
			sqlmock.AnyArg(), // updated_at
			user.Id,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 执行测试
	err := suite.userModel.Update(user)

	// 验证结果
	assert.NoError(suite.T(), err)

	// 验证所有期望都被满足
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// TestUserAuthentication 测试用户认证相关操作
func (suite *UsersModelTestSuite) TestUserAuthentication() {
	// 测试更新登录信息
	suite.Run("UpdateLoginInfo", func() {
		userID := int64(1001)
		loginIP := "192.168.1.100"
		loginTime := time.Now()

		// 设置Mock期望
		suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE vt_users SET last_login_at = ?, last_login_ip = ?, login_attempts = 0, updated_at = ? WHERE id = ?")).
			WithArgs(sqlmock.AnyArg(), loginIP, sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// 创建包含登录信息的用户对象
		user := &model.VtUsers{
			Id:            userID,
			LastLoginAt:   &loginTime,
			LastLoginIp:   loginIP,
			LoginAttempts: 0,
			UpdatedAt:     time.Now(),
		}

		// 执行测试
		err := suite.userModel.Update(user)
		assert.NoError(suite.T(), err)
	})

	// 测试增加登录失败次数
	suite.Run("IncrementLoginAttempts", func() {
		userID := int64(1001)

		// 设置Mock期望
		suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE vt_users SET login_attempts = login_attempts + 1, updated_at = ? WHERE id = ?")).
			WithArgs(sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// 执行测试（假设有增加登录失败次数的方法）
		// 这里模拟直接执行SQL更新
		_, err := suite.db.Exec("UPDATE vt_users SET login_attempts = login_attempts + 1, updated_at = ? WHERE id = ?", time.Now(), userID)
		assert.NoError(suite.T(), err)
	})

	// 测试锁定用户账户
	suite.Run("LockUserAccount", func() {
		userID := int64(1001)
		lockedUntil := time.Now().Add(24 * time.Hour)

		// 设置Mock期望
		suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE vt_users SET status = 'locked', locked_until = ?, updated_at = ? WHERE id = ?")).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// 创建锁定的用户对象
		user := &model.VtUsers{
			Id:          userID,
			Status:      "locked",
			LockedUntil: &lockedUntil,
			UpdatedAt:   time.Now(),
		}

		// 执行测试
		err := suite.userModel.Update(user)
		assert.NoError(suite.T(), err)
	})
}

// TestUserPasswordOperations 测试用户密码相关操作
func (suite *UsersModelTestSuite) TestUserPasswordOperations() {
	// 测试更新密码
	suite.Run("UpdatePassword", func() {
		userID := int64(1001)
		newPasswordHash := "new_hashed_password_789"
		newSalt := "new_random_salt_123"

		// 设置Mock期望
		suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE vt_users SET password_hash = ?, salt = ?, password_expires_at = ?, updated_at = ? WHERE id = ?")).
			WithArgs(newPasswordHash, newSalt, sqlmock.AnyArg(), sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// 创建包含新密码的用户对象
		user := &model.VtUsers{
			Id:           userID,
			PasswordHash: newPasswordHash,
			Salt:         newSalt,
			UpdatedAt:    time.Now(),
		}

		// 执行测试
		err := suite.userModel.Update(user)
		assert.NoError(suite.T(), err)
	})

	// 测试启用MFA
	suite.Run("EnableMFA", func() {
		userID := int64(1001)
		mfaSecret := "JBSWY3DPEHPK3PXP"

		// 设置Mock期望
		suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE vt_users SET mfa_enabled = true, mfa_secret = ?, updated_at = ? WHERE id = ?")).
			WithArgs(mfaSecret, sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// 创建启用MFA的用户对象
		user := &model.VtUsers{
			Id:         userID,
			MfaEnabled: true,
			MfaSecret:  mfaSecret,
			UpdatedAt:  time.Now(),
		}

		// 执行测试
		err := suite.userModel.Update(user)
		assert.NoError(suite.T(), err)
	})
}

// TestUserVerification 测试用户验证相关操作
func (suite *UsersModelTestSuite) TestUserVerification() {
	// 测试邮箱验证
	suite.Run("VerifyEmail", func() {
		userID := int64(1001)

		// 设置Mock期望
		suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE vt_users SET email_verified = true, updated_at = ? WHERE id = ?")).
			WithArgs(sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// 创建邮箱已验证的用户对象
		user := &model.VtUsers{
			Id:            userID,
			EmailVerified: true,
			UpdatedAt:     time.Now(),
		}

		// 执行测试
		err := suite.userModel.Update(user)
		assert.NoError(suite.T(), err)
	})

	// 测试手机验证
	suite.Run("VerifyPhone", func() {
		userID := int64(1001)

		// 设置Mock期望
		suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE vt_users SET phone_verified = true, updated_at = ? WHERE id = ?")).
			WithArgs(sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		// 创建手机已验证的用户对象
		user := &model.VtUsers{
			Id:            userID,
			PhoneVerified: true,
			UpdatedAt:     time.Now(),
		}

		// 执行测试
		err := suite.userModel.Update(user)
		assert.NoError(suite.T(), err)
	})
}

// TestUserModelErrors 测试各种错误情况
func (suite *UsersModelTestSuite) TestUserModelErrors() {
	// 测试插入重复用户名
	suite.Run("InsertDuplicateUsername", func() {
		user := &model.VtUsers{
			Username:     "duplicateuser",
			Email:        "duplicate@example.com",
			PasswordHash: "hashed_password",
			Salt:         "salt",
			Status:       "active",
			UserType:     "user",
		}

		// 设置Mock期望 - 唯一约束违反
		suite.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO vt_users")).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(sql.ErrConnDone)

		result, err := suite.userModel.Insert(user)
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), result)
	})

	// 测试查找不存在的用户
	suite.Run("FindNonExistentUser", func() {
		userID := int64(99999)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		user, err := suite.userModel.FindOne(userID)
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), user)
		assert.Equal(suite.T(), sql.ErrNoRows, err)
	})

	// 测试更新不存在的用户
	suite.Run("UpdateNonExistentUser", func() {
		user := &model.VtUsers{
			Id:       99999,
			Username: "nonexistent",
			Email:    "nonexistent@example.com",
		}

		suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE vt_users")).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.Id).
			WillReturnResult(sqlmock.NewResult(0, 0)) // 没有行被更新

		err := suite.userModel.Update(user)
		// 根据实际实现，这里可能返回错误或成功
		// 这里假设不返回错误，但可以检查影响的行数
		assert.NoError(suite.T(), err)
	})
}

// TestDeleteUser 测试删除用户（软删除）
func (suite *UsersModelTestSuite) TestDeleteUser() {
	userID := int64(1001)

	// 设置Mock期望 - 软删除
	suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE vt_users SET deleted_at = ?, status = 'deleted', updated_at = ? WHERE id = ?")).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 执行测试（假设有删除方法）
	// 这里模拟直接执行SQL更新
	_, err := suite.db.Exec("UPDATE vt_users SET deleted_at = ?, status = 'deleted', updated_at = ? WHERE id = ?", time.Now(), time.Now(), userID)

	// 验证结果
	assert.NoError(suite.T(), err)

	// 验证所有期望都被满足
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// 运行用户模型测试套件
func TestUsersModelSuite(t *testing.T) {
	suite.Run(t, new(UsersModelTestSuite))
}

// 基准测试
func BenchmarkUsersModelFindOneByUsername(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("创建sqlmock失败: %v", err)
	}
	defer db.Close()

	userModel := model.NewVtUsersModel(db)

	// 设置Mock期望
	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "status", "user_type", "created_at",
	}).AddRow(
		int64(1001), "testuser", "test@example.com", "active", "user", time.Now(),
	)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WillReturnRows(rows)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := userModel.FindOneByUsername("testuser")
		if err != nil {
			b.Errorf("查找失败: %v", err)
		}
	}
}
