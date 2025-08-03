//go:build integration
// +build integration

package integration

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"api/model"
	"api/test/config"
	"api/test/seeds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// DatabaseIntegrationTestSuite 数据库集成测试套件
type DatabaseIntegrationTestSuite struct {
	suite.Suite
	db         *sql.DB
	testConfig *config.TestConfig
	seeder     *seeds.DataSeeder
	ctx        context.Context
}

// SetupSuite 测试套件初始化
func (suite *DatabaseIntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.testConfig = config.GetTestConfig()

	// 创建数据库连接
	var err error
	suite.db, err = suite.testConfig.CreateTestDB()
	require.NoError(suite.T(), err, "创建测试数据库连接失败")

	// 初始化数据种子
	suite.seeder = seeds.NewDataSeeder(suite.db)
}

// TearDownSuite 测试套件清理
func (suite *DatabaseIntegrationTestSuite) TearDownSuite() {
	if suite.seeder != nil {
		suite.seeder.CleanupAllData()
	}
	if suite.db != nil {
		suite.db.Close()
	}
}

// SetupTest 每个测试前的初始化
func (suite *DatabaseIntegrationTestSuite) SetupTest() {
	// 清理并重新生成测试数据
	err := suite.seeder.CleanupAllData()
	require.NoError(suite.T(), err)

	err = suite.seeder.SeedAllTestData()
	require.NoError(suite.T(), err)
}

// TestUserModelCRUD 测试用户模型CRUD操作
func (suite *DatabaseIntegrationTestSuite) TestUserModelCRUD() {
	userModel := model.NewVtUsersModel(suite.db)

	suite.Run("创建用户", func() {
		user := &model.VtUsers{
			Username:     "dbtest_user",
			Email:        "dbtest@example.com",
			PasswordHash: "hashed_password",
			Salt:         "random_salt",
			RealName:     "数据库测试用户",
			UserType:     "user",
			Status:       "active",
		}

		result, err := userModel.Insert(user)
		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), result)

		// 获取插入的ID
		userID, err := result.LastInsertId()
		assert.NoError(suite.T(), err)
		assert.Greater(suite.T(), userID, int64(0))

		user.Id = userID

		// 验证插入的数据
		foundUser, err := userModel.FindOne(userID)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), user.Username, foundUser.Username)
		assert.Equal(suite.T(), user.Email, foundUser.Email)
		assert.Equal(suite.T(), user.UserType, foundUser.UserType)
	})

	suite.Run("根据用户名查找用户", func() {
		user, err := userModel.FindOneByUsername("admin_test")
		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), user)
		assert.Equal(suite.T(), "admin_test", user.Username)
		assert.Equal(suite.T(), "admin", user.UserType)
	})

	suite.Run("根据邮箱查找用户", func() {
		user, err := userModel.FindOneByEmail("admin@test.com")
		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), user)
		assert.Equal(suite.T(), "admin@test.com", user.Email)
	})

	suite.Run("更新用户信息", func() {
		user, err := userModel.FindOneByUsername("user_test")
		require.NoError(suite.T(), err)

		originalRealName := user.RealName
		user.RealName = "更新后的真实姓名"
		user.Email = "updated_email@test.com"

		err = userModel.Update(user)
		assert.NoError(suite.T(), err)

		// 验证更新
		updatedUser, err := userModel.FindOne(user.Id)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "更新后的真实姓名", updatedUser.RealName)
		assert.Equal(suite.T(), "updated_email@test.com", updatedUser.Email)
		assert.NotEqual(suite.T(), originalRealName, updatedUser.RealName)
	})

	suite.Run("获取用户列表", func() {
		users, total, err := userModel.List(1, 10, map[string]interface{}{})
		assert.NoError(suite.T(), err)
		assert.Greater(suite.T(), len(users), 0)
		assert.Greater(suite.T(), total, int64(0))

		// 测试过滤条件
		adminUsers, adminTotal, err := userModel.List(1, 10, map[string]interface{}{
			"user_type": "admin",
		})
		assert.NoError(suite.T(), err)
		assert.Greater(suite.T(), len(adminUsers), 0)
		assert.Less(suite.T(), adminTotal, total) // 管理员用户应该少于总用户数
	})

	suite.Run("删除用户（软删除）", func() {
		// 创建一个测试用户
		user := &model.VtUsers{
			Username:     "to_be_deleted",
			Email:        "delete@test.com",
			PasswordHash: "password",
			Salt:         "salt",
			UserType:     "user",
			Status:       "active",
		}

		result, err := userModel.Insert(user)
		require.NoError(suite.T(), err)
		userID, _ := result.LastInsertId()

		// 删除用户
		err = userModel.Delete(userID)
		assert.NoError(suite.T(), err)

		// 验证用户已被软删除
		deletedUser, err := userModel.FindOne(userID)
		assert.Error(suite.T(), err) // 应该找不到用户
		assert.Nil(suite.T(), deletedUser)
	})
}

// TestTrainingJobModelCRUD 测试训练作业模型CRUD操作
func (suite *DatabaseIntegrationTestSuite) TestTrainingJobModelCRUD() {
	jobModel := model.NewVtTrainingJobsModel(suite.db)

	suite.Run("创建训练作业", func() {
		job := &model.VtTrainingJobs{
			UserId:      1002, // 使用测试用户
			Name:        "db-test-pytorch-job",
			Framework:   "pytorch",
			Image:       "pytorch/pytorch:1.12.0",
			EntryPoint:  "train.py",
			Status:      "pending",
			CPUCores:    "4",
			MemoryGb:    "8",
			GPUCount:    1,
			GPUType:     "T4",
			WorkspaceId: 3001,
			Description: "数据库集成测试作业",
		}

		result, err := jobModel.Insert(job)
		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), result)

		jobID, err := result.LastInsertId()
		assert.NoError(suite.T(), err)
		assert.Greater(suite.T(), jobID, int64(0))

		// 验证插入的数据
		foundJob, err := jobModel.FindOne(jobID)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), job.Name, foundJob.Name)
		assert.Equal(suite.T(), job.Framework, foundJob.Framework)
		assert.Equal(suite.T(), job.Status, foundJob.Status)
	})

	suite.Run("根据用户ID查找训练作业", func() {
		jobs, err := jobModel.FindByUserId(1002)
		assert.NoError(suite.T(), err)
		assert.Greater(suite.T(), len(jobs), 0)

		for _, job := range jobs {
			assert.Equal(suite.T(), int64(1002), job.UserId)
		}
	})

	suite.Run("更新训练作业状态", func() {
		job, err := jobModel.FindOne(3001)
		require.NoError(suite.T(), err)

		originalStatus := job.Status
		newStatus := "running"

		err = jobModel.UpdateStatus(job.Id, newStatus)
		assert.NoError(suite.T(), err)

		// 验证状态更新
		updatedJob, err := jobModel.FindOne(3001)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), newStatus, updatedJob.Status)
		assert.NotEqual(suite.T(), originalStatus, updatedJob.Status)
	})

	suite.Run("获取训练作业列表", func() {
		jobs, total, err := jobModel.List(1, 10, map[string]interface{}{})
		assert.NoError(suite.T(), err)
		assert.Greater(suite.T(), len(jobs), 0)
		assert.Greater(suite.T(), total, int64(0))

		// 测试按状态过滤
		runningJobs, runningTotal, err := jobModel.List(1, 10, map[string]interface{}{
			"status": "running",
		})
		assert.NoError(suite.T(), err)
		assert.LessOrEqual(suite.T(), runningTotal, total)
	})
}

// TestGPUClusterModelCRUD 测试GPU集群模型CRUD操作
func (suite *DatabaseIntegrationTestSuite) TestGPUClusterModelCRUD() {
	clusterModel := model.NewVtGpuClustersModel(suite.db)

	suite.Run("创建GPU集群", func() {
		cluster := &model.VtGpuClusters{
			Name:        "db-test-cluster",
			Description: "数据库测试GPU集群",
			Kubeconfig:  "test-kubeconfig-content",
			Region:      "us-west-1",
			Zone:        "us-west-1a",
			Provider:    "aws",
			ClusterType: "kubernetes",
			Status:      "pending",
			NodeCount:   0,
		}

		result, err := clusterModel.Insert(cluster)
		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), result)

		clusterID, err := result.LastInsertId()
		assert.NoError(suite.T(), err)
		assert.Greater(suite.T(), clusterID, int64(0))

		// 验证插入的数据
		foundCluster, err := clusterModel.FindOne(clusterID)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), cluster.Name, foundCluster.Name)
		assert.Equal(suite.T(), cluster.Provider, foundCluster.Provider)
		assert.Equal(suite.T(), cluster.Status, foundCluster.Status)
	})

	suite.Run("检查集群名称是否存在", func() {
		exists, err := clusterModel.CheckNameExists("test-cluster-1")
		assert.NoError(suite.T(), err)
		assert.True(suite.T(), exists)

		exists, err = clusterModel.CheckNameExists("non-existent-cluster")
		assert.NoError(suite.T(), err)
		assert.False(suite.T(), exists)
	})

	suite.Run("更新集群状态", func() {
		cluster, err := clusterModel.FindOne(2001)
		require.NoError(suite.T(), err)

		newStatus := "active"
		err = clusterModel.UpdateStatus(cluster.Id, newStatus)
		assert.NoError(suite.T(), err)

		// 验证状态更新
		updatedCluster, err := clusterModel.FindOne(2001)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), newStatus, updatedCluster.Status)
	})

	suite.Run("获取集群列表", func() {
		clusters, total, err := clusterModel.List(1, 10, map[string]interface{}{})
		assert.NoError(suite.T(), err)
		assert.Greater(suite.T(), len(clusters), 0)
		assert.Greater(suite.T(), total, int64(0))

		// 测试按提供商过滤
		awsClusters, awsTotal, err := clusterModel.List(1, 10, map[string]interface{}{
			"provider": "aws",
		})
		assert.NoError(suite.T(), err)
		assert.LessOrEqual(suite.T(), awsTotal, total)
	})
}

// TestDatabaseTransactions 测试数据库事务
func (suite *DatabaseIntegrationTestSuite) TestDatabaseTransactions() {
	suite.Run("成功的事务", func() {
		tx, err := suite.db.BeginTx(suite.ctx, nil)
		require.NoError(suite.T(), err)

		// 在事务中创建用户
		_, err = tx.ExecContext(suite.ctx,
			"INSERT INTO vt_users (username, email, password_hash, salt, user_type, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())",
			"tx_test_user", "tx@test.com", "password", "salt", "user", "active")
		assert.NoError(suite.T(), err)

		// 提交事务
		err = tx.Commit()
		assert.NoError(suite.T(), err)

		// 验证用户已创建
		var count int
		err = suite.db.QueryRowContext(suite.ctx,
			"SELECT COUNT(*) FROM vt_users WHERE username = ?", "tx_test_user").Scan(&count)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), 1, count)
	})

	suite.Run("回滚的事务", func() {
		tx, err := suite.db.BeginTx(suite.ctx, nil)
		require.NoError(suite.T(), err)

		// 在事务中创建用户
		_, err = tx.ExecContext(suite.ctx,
			"INSERT INTO vt_users (username, email, password_hash, salt, user_type, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())",
			"rollback_test_user", "rollback@test.com", "password", "salt", "user", "active")
		assert.NoError(suite.T(), err)

		// 回滚事务
		err = tx.Rollback()
		assert.NoError(suite.T(), err)

		// 验证用户未创建
		var count int
		err = suite.db.QueryRowContext(suite.ctx,
			"SELECT COUNT(*) FROM vt_users WHERE username = ?", "rollback_test_user").Scan(&count)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), 0, count)
	})
}

// TestDatabaseConstraints 测试数据库约束
func (suite *DatabaseIntegrationTestSuite) TestDatabaseConstraints() {
	suite.Run("唯一约束测试", func() {
		userModel := model.NewVtUsersModel(suite.db)

		// 尝试创建重复用户名的用户
		user1 := &model.VtUsers{
			Username:     "unique_test",
			Email:        "unique1@test.com",
			PasswordHash: "password",
			Salt:         "salt",
			UserType:     "user",
			Status:       "active",
		}

		_, err := userModel.Insert(user1)
		assert.NoError(suite.T(), err)

		// 尝试创建相同用户名的用户
		user2 := &model.VtUsers{
			Username:     "unique_test", // 重复用户名
			Email:        "unique2@test.com",
			PasswordHash: "password",
			Salt:         "salt",
			UserType:     "user",
			Status:       "active",
		}

		_, err = userModel.Insert(user2)
		assert.Error(suite.T(), err) // 应该报错
		assert.Contains(suite.T(), err.Error(), "Duplicate entry")
	})

	suite.Run("外键约束测试", func() {
		jobModel := model.NewVtTrainingJobsModel(suite.db)

		// 尝试创建引用不存在用户的训练作业
		job := &model.VtTrainingJobs{
			UserId:    99999, // 不存在的用户ID
			Name:      "fk_test_job",
			Framework: "pytorch",
			Status:    "pending",
		}

		_, err := jobModel.Insert(job)
		assert.Error(suite.T(), err) // 应该报错
		assert.Contains(suite.T(), err.Error(), "foreign key constraint")
	})
}

// TestDatabasePerformance 测试数据库性能
func (suite *DatabaseIntegrationTestSuite) TestDatabasePerformance() {
	suite.Run("批量插入性能测试", func() {
		userModel := model.NewVtUsersModel(suite.db)

		start := time.Now()

		// 批量创建100个用户
		for i := 0; i < 100; i++ {
			user := &model.VtUsers{
				Username:     fmt.Sprintf("perf_user_%d", i),
				Email:        fmt.Sprintf("perf%d@test.com", i),
				PasswordHash: "password",
				Salt:         "salt",
				UserType:     "user",
				Status:       "active",
			}

			_, err := userModel.Insert(user)
			assert.NoError(suite.T(), err)
		}

		duration := time.Since(start)
		suite.T().Logf("批量插入100个用户耗时: %v", duration)

		// 性能要求：100个用户插入应在5秒内完成
		assert.Less(suite.T(), duration, 5*time.Second)
	})

	suite.Run("查询性能测试", func() {
		userModel := model.NewVtUsersModel(suite.db)

		start := time.Now()

		// 执行100次查询
		for i := 0; i < 100; i++ {
			_, _, err := userModel.List(1, 10, map[string]interface{}{})
			assert.NoError(suite.T(), err)
		}

		duration := time.Since(start)
		suite.T().Logf("100次列表查询耗时: %v", duration)

		// 性能要求：100次查询应在2秒内完成
		assert.Less(suite.T(), duration, 2*time.Second)
	})
}

// TestDatabaseConnectionPool 测试数据库连接池
func (suite *DatabaseIntegrationTestSuite) TestDatabaseConnectionPool() {
	suite.Run("并发查询测试", func() {
		const numGoroutines = 20
		results := make(chan error, numGoroutines)

		userModel := model.NewVtUsersModel(suite.db)

		// 启动多个并发查询
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				// 执行查询操作
				_, _, err := userModel.List(1, 10, map[string]interface{}{})
				results <- err
			}(i)
		}

		// 等待所有查询完成
		for i := 0; i < numGoroutines; i++ {
			err := <-results
			assert.NoError(suite.T(), err, "并发查询失败")
		}
	})
}

// 运行数据库集成测试套件
func TestDatabaseIntegrationSuite(t *testing.T) {
	suite.Run(t, new(DatabaseIntegrationTestSuite))
}
