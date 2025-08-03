package test

import (
	"fmt"
	"testing"

	"api/test/config"
	"api/test/utils"
	"github.com/stretchr/testify/assert"
)

// TestFrameworkSetup 测试框架基础设置
func TestFrameworkSetup(t *testing.T) {
	// 测试配置加载
	testConfig := config.GetTestConfig()
	assert.NotNil(t, testConfig)
	assert.Equal(t, "mysql", testConfig.Database.Driver)
	assert.Equal(t, "volctraindb", testConfig.Database.Database)

	fmt.Printf("✅ 测试配置加载成功\n")
	fmt.Printf("   数据库: %s\n", testConfig.Database.Database)
	fmt.Printf("   Redis端口: %d\n", testConfig.Redis.Port)
}

// TestUtilsSetup 测试工具函数
func TestUtilsSetup(t *testing.T) {
	// 创建测试套件
	suite := utils.SetupTestSuite(t)
	defer suite.TeardownTestSuite()

	assert.NotNil(t, suite)
	assert.NotNil(t, suite.Config)

	fmt.Printf("✅ 测试工具包初始化成功\n")
}

// TestMockDBCreation 测试Mock数据库创建
func TestMockDBCreation(t *testing.T) {
	db, mock, err := config.CreateMockDB()
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, mock)

	defer db.Close()

	// 简单的Mock测试
	mock.ExpectQuery("SELECT 1").WillReturnRows(mock.NewRows([]string{"result"}).AddRow(1))

	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	assert.NoError(t, err)
	assert.Equal(t, 1, result)

	// 验证所有期望都被满足
	assert.NoError(t, mock.ExpectationsWereMet())

	fmt.Printf("✅ Mock数据库创建和测试成功\n")
}
