//go:build integration
// +build integration

package integration

import (
	"context"
	"database/sql"
	"sync"
	"testing"
	"time"

	"api/internal/config"
	"api/model"
	"api/pkg/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zeromicro/go-zero/core/conf"
)

// DatabaseIntegrationTestSuite 数据库集成测试套件
type DatabaseIntegrationTestSuite struct {
	suite.Suite
	db               *sql.DB
	config           config.Config
	userModel        model.VtUsersModel
	trainingJobModel model.VtTrainingJobsModel
	gpuDeviceModel   model.VtGpuDevicesModel
}

// SetupSuite 整个测试套件开始前的初始化
func (suite *DatabaseIntegrationTestSuite) SetupSuite() {
	// 加载测试配置
	var c config.Config
	conf.MustLoad("../../../etc/config-dev.yaml", &c)

	// 使用测试数据库
	c.MySQL.DataSource = "test_user:test_password@tcp(localhost:3306)/volctrain_test?charset=utf8mb4&parseTime=true&loc=Local"
	suite.config = c

	// 初始化数据库连接
	var err error
	suite.db, err = database.NewMySQLConnection(c.MySQL)
	suite.Require().NoError(err, "无法连接到测试数据库")

	// 初始化模型
	suite.userModel = model.NewVtUsersModel(suite.db)
	suite.trainingJobModel = model.NewVtTrainingJobsModel(suite.db)
	suite.gpuDeviceModel = model.NewVtGpuDevicesModel(suite.db)

	// 创建测试表（如果不存在）
	suite.createTestTables()
}

// TearDownSuite 整个测试套件结束后的清理
func (suite *DatabaseIntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		// 清理测试数据
		suite.cleanupTestData()
		suite.db.Close()
	}
}

// SetupTest 每个测试前的初始化
func (suite *DatabaseIntegrationTestSuite) SetupTest() {
	// 清理可能存在的测试数据
	suite.cleanupTestData()
}

// TearDownTest 每个测试后的清理
func (suite *DatabaseIntegrationTestSuite) TearDownTest() {
	// 清理测试数据
	suite.cleanupTestData()
}

// createTestTables 创建测试表
func (suite *DatabaseIntegrationTestSuite) createTestTables() {
	// 这里应该包含创建测试表的SQL语句
	// 在真实环境中，可能会使用数据库迁移工具
	suite.T().Log("创建测试表...")

	// 示例：创建用户表（简化版本）
	createUserTableSQL := `
		CREATE TABLE IF NOT EXISTS vt_users_test (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			salt VARCHAR(32) NOT NULL,
			status VARCHAR(20) DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)
	`

	_, err := suite.db.Exec(createUserTableSQL)
	suite.Require().NoError(err, "创建测试用户表失败")
}

// cleanupTestData 清理测试数据
func (suite *DatabaseIntegrationTestSuite) cleanupTestData() {
	// 清理测试产生的数据
	tables := []string{
		"vt_training_jobs",
		"vt_gpu_devices",
		"vt_users",
		"vt_users_test",
	}

	for _, table := range tables {
		_, err := suite.db.Exec("DELETE FROM " + table + " WHERE 1=1")
		if err != nil {
			// 表可能不存在，继续清理其他表
			suite.T().Logf("清理表 %s 失败: %v", table, err)
		}
	}
}

// TestDatabaseConnection 测试数据库连接
func (suite *DatabaseIntegrationTestSuite) TestDatabaseConnection() {
	// 测试基本连接
	err := suite.db.Ping()
	assert.NoError(suite.T(), err, "数据库连接失败")

	// 测试查询
	var version string
	err = suite.db.QueryRow("SELECT VERSION()").Scan(&version)
	assert.NoError(suite.T(), err, "查询数据库版本失败")
	assert.NotEmpty(suite.T(), version, "数据库版本为空")

	suite.T().Logf("数据库版本: %s", version)
}

// TestTransactionRollback 测试事务回滚
func (suite *DatabaseIntegrationTestSuite) TestTransactionRollback() {
	// 开始事务
	tx, err := suite.db.Begin()
	assert.NoError(suite.T(), err, "开始事务失败")

	// 在事务中插入数据
	_, err = tx.Exec(`
		INSERT INTO vt_users_test (username, email, password_hash, salt) 
		VALUES (?, ?, ?, ?)
	`, "test_rollback", "rollback@test.com", "hashed_password", "salt")
	assert.NoError(suite.T(), err, "事务中插入数据失败")

	// 验证数据在事务中存在
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM vt_users_test WHERE username = ?", "test_rollback").Scan(&count)
	assert.NoError(suite.T(), err, "事务中查询数据失败")
	assert.Equal(suite.T(), 1, count, "事务中数据不存在")

	// 回滚事务
	err = tx.Rollback()
	assert.NoError(suite.T(), err, "回滚事务失败")

	// 验证数据已回滚
	err = suite.db.QueryRow("SELECT COUNT(*) FROM vt_users_test WHERE username = ?", "test_rollback").Scan(&count)
	assert.NoError(suite.T(), err, "回滚后查询数据失败")
	assert.Equal(suite.T(), 0, count, "数据未正确回滚")
}

// TestTransactionCommit 测试事务提交
func (suite *DatabaseIntegrationTestSuite) TestTransactionCommit() {
	// 开始事务
	tx, err := suite.db.Begin()
	assert.NoError(suite.T(), err, "开始事务失败")

	// 在事务中插入数据
	_, err = tx.Exec(`
		INSERT INTO vt_users_test (username, email, password_hash, salt) 
		VALUES (?, ?, ?, ?)
	`, "test_commit", "commit@test.com", "hashed_password", "salt")
	assert.NoError(suite.T(), err, "事务中插入数据失败")

	// 提交事务
	err = tx.Commit()
	assert.NoError(suite.T(), err, "提交事务失败")

	// 验证数据已提交
	var count int
	err = suite.db.QueryRow("SELECT COUNT(*) FROM vt_users_test WHERE username = ?", "test_commit").Scan(&count)
	assert.NoError(suite.T(), err, "提交后查询数据失败")
	assert.Equal(suite.T(), 1, count, "数据未正确提交")
}

// TestNestedTransaction 测试嵌套事务
func (suite *DatabaseIntegrationTestSuite) TestNestedTransaction() {
	// MySQL不支持真正的嵌套事务，但我们可以测试保存点

	// 开始主事务
	tx, err := suite.db.Begin()
	assert.NoError(suite.T(), err, "开始主事务失败")
	defer tx.Rollback()

	// 插入第一条数据
	_, err = tx.Exec(`
		INSERT INTO vt_users_test (username, email, password_hash, salt) 
		VALUES (?, ?, ?, ?)
	`, "user1", "user1@test.com", "hashed_password", "salt")
	assert.NoError(suite.T(), err, "插入第一条数据失败")

	// 创建保存点
	_, err = tx.Exec("SAVEPOINT sp1")
	assert.NoError(suite.T(), err, "创建保存点失败")

	// 插入第二条数据
	_, err = tx.Exec(`
		INSERT INTO vt_users_test (username, email, password_hash, salt) 
		VALUES (?, ?, ?, ?)
	`, "user2", "user2@test.com", "hashed_password", "salt")
	assert.NoError(suite.T(), err, "插入第二条数据失败")

	// 回滚到保存点
	_, err = tx.Exec("ROLLBACK TO SAVEPOINT sp1")
	assert.NoError(suite.T(), err, "回滚到保存点失败")

	// 验证只有第一条数据存在
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM vt_users_test").Scan(&count)
	assert.NoError(suite.T(), err, "查询数据数量失败")
	assert.Equal(suite.T(), 1, count, "保存点回滚后数据数量不正确")

	// 验证具体数据
	var username string
	err = tx.QueryRow("SELECT username FROM vt_users_test").Scan(&username)
	assert.NoError(suite.T(), err, "查询用户名失败")
	assert.Equal(suite.T(), "user1", username, "保存点回滚后数据不正确")
}

// TestConnectionPool 测试连接池
func (suite *DatabaseIntegrationTestSuite) TestConnectionPool() {
	// 设置连接池参数
	suite.db.SetMaxOpenConns(10)
	suite.db.SetMaxIdleConns(5)
	suite.db.SetConnMaxLifetime(time.Hour)

	// 并发测试连接池
	concurrency := 20
	var wg sync.WaitGroup
	errors := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			// 执行数据库操作
			var result int
			err := suite.db.QueryRow("SELECT ? + ?", index, index).Scan(&result)
			if err != nil {
				errors <- err
				return
			}

			expected := index + index
			if result != expected {
				errors <- assert.AnError
				return
			}

			errors <- nil
		}(i)
	}

	wg.Wait()
	close(errors)

	// 检查结果
	for err := range errors {
		assert.NoError(suite.T(), err, "连接池并发测试失败")
	}

	// 验证连接池统计信息
	stats := suite.db.Stats()
	suite.T().Logf("连接池统计: OpenConnections=%d, InUse=%d, Idle=%d",
		stats.OpenConnections, stats.InUse, stats.Idle)

	assert.LessOrEqual(suite.T(), stats.OpenConnections, 10, "连接数超过最大限制")
}

// TestLongRunningTransaction 测试长时间运行的事务
func (suite *DatabaseIntegrationTestSuite) TestLongRunningTransaction() {
	// 开始事务
	tx, err := suite.db.Begin()
	assert.NoError(suite.T(), err, "开始长事务失败")
	defer tx.Rollback()

	// 插入数据
	_, err = tx.Exec(`
		INSERT INTO vt_users_test (username, email, password_hash, salt) 
		VALUES (?, ?, ?, ?)
	`, "long_tx_user", "longtx@test.com", "hashed_password", "salt")
	assert.NoError(suite.T(), err, "长事务中插入数据失败")

	// 模拟长时间操作
	time.Sleep(100 * time.Millisecond)

	// 验证事务仍然有效
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM vt_users_test WHERE username = ?", "long_tx_user").Scan(&count)
	assert.NoError(suite.T(), err, "长事务查询失败")
	assert.Equal(suite.T(), 1, count, "长事务数据不一致")

	// 在事务外查询，应该看不到未提交的数据
	err = suite.db.QueryRow("SELECT COUNT(*) FROM vt_users_test WHERE username = ?", "long_tx_user").Scan(&count)
	assert.NoError(suite.T(), err, "事务外查询失败")
	assert.Equal(suite.T(), 0, count, "事务隔离失败")
}

// TestDeadlockDetection 测试死锁检测
func (suite *DatabaseIntegrationTestSuite) TestDeadlockDetection() {
	// 插入测试数据
	_, err := suite.db.Exec(`
		INSERT INTO vt_users_test (username, email, password_hash, salt) 
		VALUES (?, ?, ?, ?), (?, ?, ?, ?)
	`, "deadlock1", "deadlock1@test.com", "hash1", "salt1",
		"deadlock2", "deadlock2@test.com", "hash2", "salt2")
	assert.NoError(suite.T(), err, "插入死锁测试数据失败")

	// 创建两个并发事务来模拟死锁
	var wg sync.WaitGroup
	deadlockDetected := false

	wg.Add(2)

	// 事务1
	go func() {
		defer wg.Done()

		tx1, err := suite.db.Begin()
		if err != nil {
			suite.T().Errorf("事务1开始失败: %v", err)
			return
		}
		defer tx1.Rollback()

		// 锁定第一条记录
		_, err = tx1.Exec("SELECT * FROM vt_users_test WHERE username = ? FOR UPDATE", "deadlock1")
		if err != nil {
			suite.T().Errorf("事务1锁定记录失败: %v", err)
			return
		}

		time.Sleep(50 * time.Millisecond)

		// 尝试锁定第二条记录（可能导致死锁）
		_, err = tx1.Exec("SELECT * FROM vt_users_test WHERE username = ? FOR UPDATE", "deadlock2")
		if err != nil {
			suite.T().Logf("事务1检测到死锁或锁等待: %v", err)
			deadlockDetected = true
		}
	}()

	// 事务2
	go func() {
		defer wg.Done()

		tx2, err := suite.db.Begin()
		if err != nil {
			suite.T().Errorf("事务2开始失败: %v", err)
			return
		}
		defer tx2.Rollback()

		// 锁定第二条记录
		_, err = tx2.Exec("SELECT * FROM vt_users_test WHERE username = ? FOR UPDATE", "deadlock2")
		if err != nil {
			suite.T().Errorf("事务2锁定记录失败: %v", err)
			return
		}

		time.Sleep(50 * time.Millisecond)

		// 尝试锁定第一条记录（可能导致死锁）
		_, err = tx2.Exec("SELECT * FROM vt_users_test WHERE username = ? FOR UPDATE", "deadlock1")
		if err != nil {
			suite.T().Logf("事务2检测到死锁或锁等待: %v", err)
			deadlockDetected = true
		}
	}()

	wg.Wait()

	// 在真实的死锁场景中，MySQL会自动检测并回滚其中一个事务
	suite.T().Logf("死锁检测结果: %v", deadlockDetected)
}

// TestDatabasePerformance 测试数据库性能
func (suite *DatabaseIntegrationTestSuite) TestDatabasePerformance() {
	// 批量插入性能测试
	batchSize := 1000
	start := time.Now()

	tx, err := suite.db.Begin()
	assert.NoError(suite.T(), err, "开始批量插入事务失败")
	defer tx.Rollback()

	// 准备批量插入语句
	stmt, err := tx.Prepare(`
		INSERT INTO vt_users_test (username, email, password_hash, salt) 
		VALUES (?, ?, ?, ?)
	`)
	assert.NoError(suite.T(), err, "准备批量插入语句失败")
	defer stmt.Close()

	// 执行批量插入
	for i := 0; i < batchSize; i++ {
		_, err = stmt.Exec(
			fmt.Sprintf("perf_user_%d", i),
			fmt.Sprintf("perf_user_%d@test.com", i),
			"hashed_password",
			"salt",
		)
		assert.NoError(suite.T(), err, "批量插入第 %d 条记录失败", i)
	}

	err = tx.Commit()
	assert.NoError(suite.T(), err, "提交批量插入事务失败")

	duration := time.Since(start)
	suite.T().Logf("批量插入 %d 条记录耗时: %v", batchSize, duration)

	// 性能指标验证
	assert.Less(suite.T(), duration, 10*time.Second, "批量插入性能过慢")

	// 查询性能测试
	start = time.Now()

	var count int
	err = suite.db.QueryRow("SELECT COUNT(*) FROM vt_users_test").Scan(&count)
	assert.NoError(suite.T(), err, "查询总数失败")
	assert.Equal(suite.T(), batchSize, count, "插入数据数量不正确")

	queryDuration := time.Since(start)
	suite.T().Logf("查询 %d 条记录耗时: %v", batchSize, queryDuration)

	assert.Less(suite.T(), queryDuration, 1*time.Second, "查询性能过慢")
}

// TestConnectionTimeout 测试连接超时
func (suite *DatabaseIntegrationTestSuite) TestConnectionTimeout() {
	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// 执行一个可能耗时的查询
	var result int
	err := suite.db.QueryRowContext(ctx, "SELECT SLEEP(2)").Scan(&result)

	// 应该超时
	assert.Error(suite.T(), err, "查询应该超时")
	assert.Contains(suite.T(), err.Error(), "context deadline exceeded", "错误类型不正确")
}

// TestReadWriteSplit 测试读写分离（如果配置了主从复制）
func (suite *DatabaseIntegrationTestSuite) TestReadWriteSplit() {
	// 这个测试需要配置主从数据库
	// 在单机测试环境中，我们可以模拟这个行为

	// 写操作（应该到主库）
	_, err := suite.db.Exec(`
		INSERT INTO vt_users_test (username, email, password_hash, salt) 
		VALUES (?, ?, ?, ?)
	`, "rw_split_user", "rwsplit@test.com", "hashed_password", "salt")
	assert.NoError(suite.T(), err, "写入主库失败")

	// 读操作（可能从从库读取）
	var count int
	err = suite.db.QueryRow("SELECT COUNT(*) FROM vt_users_test WHERE username = ?", "rw_split_user").Scan(&count)
	assert.NoError(suite.T(), err, "从从库读取失败")
	assert.Equal(suite.T(), 1, count, "读写一致性检查失败")
}

// TestDatabaseMigration 测试数据库迁移
func (suite *DatabaseIntegrationTestSuite) TestDatabaseMigration() {
	// 创建迁移表
	migrationSQL := `
		CREATE TABLE IF NOT EXISTS vt_test_migration (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(100),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`

	_, err := suite.db.Exec(migrationSQL)
	assert.NoError(suite.T(), err, "创建迁移表失败")

	// 插入测试数据
	_, err = suite.db.Exec("INSERT INTO vt_test_migration (name) VALUES (?)", "migration_test")
	assert.NoError(suite.T(), err, "插入迁移测试数据失败")

	// 执行表结构变更（添加列）
	alterSQL := "ALTER TABLE vt_test_migration ADD COLUMN description TEXT"
	_, err = suite.db.Exec(alterSQL)
	assert.NoError(suite.T(), err, "表结构迁移失败")

	// 验证新列可用
	_, err = suite.db.Exec("UPDATE vt_test_migration SET description = ? WHERE name = ?", "test description", "migration_test")
	assert.NoError(suite.T(), err, "更新新列失败")

	// 清理迁移表
	_, err = suite.db.Exec("DROP TABLE vt_test_migration")
	assert.NoError(suite.T(), err, "清理迁移表失败")
}

// 运行数据库集成测试套件
func TestDatabaseIntegrationSuite(t *testing.T) {
	suite.Run(t, new(DatabaseIntegrationTestSuite))
}
