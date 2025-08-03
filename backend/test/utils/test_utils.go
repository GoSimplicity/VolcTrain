package utils

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"api/test/config"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSuite 测试套件基础结构
type TestSuite struct {
	T      *testing.T
	DB     *sql.DB
	Redis  *redis.Client
	Config *config.TestConfig
	Server *httptest.Server
	Client *http.Client
	Token  string // JWT Token for API tests
}

// SetupTestSuite 创建测试套件
func SetupTestSuite(t *testing.T) *TestSuite {
	suite := &TestSuite{
		T:      t,
		Config: config.GetTestConfig(),
		Client: &http.Client{Timeout: 30 * time.Second},
	}

	// 设置测试数据库
	var err error
	suite.DB, err = suite.Config.CreateTestDB()
	require.NoError(t, err, "创建测试数据库失败")

	// 设置测试Redis
	suite.Redis = suite.Config.CreateTestRedis()

	return suite
}

// TeardownTestSuite 清理测试套件
func (s *TestSuite) TeardownTestSuite() {
	if s.DB != nil {
		s.Config.CleanupTestData(s.DB)
		s.DB.Close()
	}
	if s.Redis != nil {
		s.Redis.Close()
	}
	if s.Server != nil {
		s.Server.Close()
	}
}

// CreateTestUser 创建测试用户
func (s *TestSuite) CreateTestUser(username, email, userType string) int64 {
	query := `INSERT INTO vt_users (username, email, password_hash, salt, user_type, status, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, 'active', NOW(), NOW())`

	result, err := s.DB.Exec(query, username, email, "hashed_password", "salt", userType)
	require.NoError(s.T, err, "创建测试用户失败")

	userID, err := result.LastInsertId()
	require.NoError(s.T, err, "获取用户ID失败")

	return userID
}

// CreateTestGPUCluster 创建测试GPU集群
func (s *TestSuite) CreateTestGPUCluster(name, kubeconfig string) int64 {
	query := `INSERT INTO vt_gpu_clusters (name, description, kubeconfig, status, created_at, updated_at) 
			  VALUES (?, ?, ?, 'active', NOW(), NOW())`

	result, err := s.DB.Exec(query, name, "测试集群", kubeconfig)
	require.NoError(s.T, err, "创建测试GPU集群失败")

	clusterID, err := result.LastInsertId()
	require.NoError(s.T, err, "获取集群ID失败")

	return clusterID
}

// CreateTestTrainingJob 创建测试训练作业
func (s *TestSuite) CreateTestTrainingJob(userID int64, name, framework string) int64 {
	query := `INSERT INTO vt_training_jobs (user_id, name, framework, image, status, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, 'pending', NOW(), NOW())`

	result, err := s.DB.Exec(query, userID, name, framework, "pytorch/pytorch:latest")
	require.NoError(s.T, err, "创建测试训练作业失败")

	jobID, err := result.LastInsertId()
	require.NoError(s.T, err, "获取作业ID失败")

	return jobID
}

// AssertDBCount 断言数据库记录数量
func (s *TestSuite) AssertDBCount(table string, expectedCount int) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	err := s.DB.QueryRow(query).Scan(&count)
	require.NoError(s.T, err)
	assert.Equal(s.T, expectedCount, count, "表 %s 记录数量不符合预期", table)
}

// AssertDBExists 断言数据库记录存在
func (s *TestSuite) AssertDBExists(table, whereClause string, args ...interface{}) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, whereClause)
	err := s.DB.QueryRow(query, args...).Scan(&count)
	require.NoError(s.T, err)
	assert.Greater(s.T, count, 0, "记录不存在: %s WHERE %s", table, whereClause)
}

// PostJSON 发送JSON POST请求
func (s *TestSuite) PostJSON(url string, data interface{}) *httptest.ResponseRecorder {
	jsonData, err := json.Marshal(data)
	require.NoError(s.T, err)

	req := httptest.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	if s.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.Token))
	}

	w := httptest.NewRecorder()
	return w
}

// GetWithAuth 发送带认证的GET请求
func (s *TestSuite) GetWithAuth(url string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", url, nil)

	if s.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.Token))
	}

	w := httptest.NewRecorder()
	return w
}

// DecodeJSON 解析JSON响应
func (s *TestSuite) DecodeJSON(w *httptest.ResponseRecorder, v interface{}) {
	err := json.NewDecoder(w.Body).Decode(v)
	require.NoError(s.T, err, "解析JSON响应失败")
}

// AssertJSONResponse 断言JSON响应
func (s *TestSuite) AssertJSONResponse(w *httptest.ResponseRecorder, expectedStatus int, expectedData interface{}) {
	assert.Equal(s.T, expectedStatus, w.Code)
	assert.True(s.T, strings.Contains(w.Header().Get("Content-Type"), "application/json"))

	if expectedData != nil {
		var actualData interface{}
		s.DecodeJSON(w, &actualData)
		assert.Equal(s.T, expectedData, actualData)
	}
}

// WaitForCondition 等待条件满足
func (s *TestSuite) WaitForCondition(condition func() bool, timeout time.Duration, msg string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			require.Fail(s.T, fmt.Sprintf("等待条件超时: %s", msg))
		case <-ticker.C:
			if condition() {
				return
			}
		}
	}
}

// SetupRedisTestData 设置Redis测试数据
func (s *TestSuite) SetupRedisTestData() {
	ctx := context.Background()

	// 清理Redis测试数据
	s.Redis.FlushDB(ctx)

	// 设置一些测试数据
	s.Redis.Set(ctx, "test:key1", "value1", time.Hour)
	s.Redis.Set(ctx, "test:key2", "value2", time.Hour)
}

// CleanupRedisTestData 清理Redis测试数据
func (s *TestSuite) CleanupRedisTestData() {
	ctx := context.Background()
	s.Redis.FlushDB(ctx)
}

// GenerateTestJWT 生成测试JWT Token
func (s *TestSuite) GenerateTestJWT(userID int64, username, userType string) string {
	// 这里应该使用实际的JWT生成逻辑
	// 为了测试目的，可以使用一个固定的token或者mock JWT
	return "test_jwt_token_" + username
}

// AssertAPIError 断言API错误响应
func (s *TestSuite) AssertAPIError(w *httptest.ResponseRecorder, expectedStatus int, expectedMessage string) {
	assert.Equal(s.T, expectedStatus, w.Code)

	var response map[string]interface{}
	s.DecodeJSON(w, &response)

	if message, ok := response["message"]; ok {
		assert.Contains(s.T, message.(string), expectedMessage)
	}
}

// CreateTestTables 创建测试所需的数据表
func (s *TestSuite) CreateTestTables() {
	// 这里可以添加创建测试表的SQL语句
	// 或者运行数据库迁移脚本

	// 示例：创建用户表
	createUserTableSQL := `
	CREATE TABLE IF NOT EXISTS vt_users_test (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		username VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		salt VARCHAR(255) NOT NULL,
		user_type VARCHAR(50) DEFAULT 'user',
		status VARCHAR(50) DEFAULT 'active',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	_, err := s.DB.Exec(createUserTableSQL)
	require.NoError(s.T, err, "创建测试表失败")
}
