package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

// TestConfig 测试配置
type TestConfig struct {
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	API      APIConfig      `json:"api"`
}

type DatabaseConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Database int    `json:"database"`
}

type APIConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// GetTestConfig 获取测试配置
func GetTestConfig() *TestConfig {
	return &TestConfig{
		Database: DatabaseConfig{
			Driver:   "mysql",
			Host:     getEnv("TEST_DB_HOST", "localhost"),
			Port:     3306,
			Database: getEnv("TEST_DB_NAME", "volctraindb"),
			Username: getEnv("TEST_DB_USER", "volctrain_app"),
			Password: getEnv("TEST_DB_PASS", "Abc@1234"),
		},
		Redis: RedisConfig{
			Host:     getEnv("TEST_REDIS_HOST", "localhost"),
			Port:     6379,
			Password: getEnv("TEST_REDIS_PASS", ""),
			Database: 1, // 使用数据库1进行测试
		},
		API: APIConfig{
			Host: getEnv("TEST_API_HOST", "localhost"),
			Port: 8888,
		},
	}
}

// CreateTestDB 创建测试数据库连接
func (c *TestConfig) CreateTestDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database)

	db, err := sql.Open(c.Database.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("创建数据库连接失败: %w", err)
	}

	// 配置连接池
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// 测试连接
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("数据库连接测试失败: %w", err)
	}

	return db, nil
}

// CreateMockDB 创建Mock数据库
func CreateMockDB() (*sql.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, fmt.Errorf("创建Mock数据库失败: %w", err)
	}
	return db, mock, nil
}

// CreateTestRedis 创建测试Redis连接
func (c *TestConfig) CreateTestRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password,
		DB:       c.Redis.Database,
	})
	return rdb
}

// CleanupTestData 清理测试数据
func (c *TestConfig) CleanupTestData(db *sql.DB) error {
	// 清理测试数据的SQL语句（按依赖关系排序）
	cleanupQueries := map[string]string{
		"vt_training_job_instances": "DELETE FROM vt_training_job_instances WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)",
		"vt_training_jobs":          "DELETE FROM vt_training_jobs WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)",
		"vt_training_queues":        "DELETE FROM vt_training_queues WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)",
		"vt_gpu_devices":            "DELETE FROM vt_gpu_devices WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)",
		"vt_gpu_nodes":              "DELETE FROM vt_gpu_nodes WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)",
		"vt_gpu_clusters":           "DELETE FROM vt_gpu_clusters WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)",
		"vt_monitor_data":           "DELETE FROM vt_monitor_data WHERE collection_time >= DATE_SUB(NOW(), INTERVAL 1 HOUR)",
		"vt_alert_records":          "DELETE FROM vt_alert_records WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)",
		"vt_notification_channels":  "DELETE FROM vt_notification_channels WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)",
		"vt_users":                  "DELETE FROM vt_users WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)",
	}

	for table, query := range cleanupQueries {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("清理表 %s 失败: %v", table, err)
		}
	}

	return nil
}

// getEnv 获取环境变量，提供默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
