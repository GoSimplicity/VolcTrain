package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"api/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

// MySQLManager MySQL数据库管理器
type MySQLManager struct {
	db     *sql.DB
	config config.MySQLConfig
}

// NewMySQLConnection 创建MySQL数据库连接
func NewMySQLConnection(c config.MySQLConfig) (*sql.DB, error) {
	// 构建DSN连接字符串，添加更多配置参数优化连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s&timeout=30s&readTimeout=30s&writeTimeout=30s&interpolateParams=true",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.Charset, c.ParseTime, c.Loc)

	// 创建数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("无法打开数据库连接: %w", err)
	}

	// 配置连接池参数
	db.SetMaxOpenConns(c.MaxOpenConns)                                // 最大打开连接数
	db.SetMaxIdleConns(c.MaxIdleConns)                                // 最大空闲连接数
	db.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second) // 连接最大生存时间
	db.SetConnMaxIdleTime(30 * time.Minute)                           // 连接最大空闲时间

	// 测试数据库连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("数据库连接测试失败: %w", err)
	}

	log.Printf("MySQL连接池配置: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%ds",
		c.MaxOpenConns, c.MaxIdleConns, c.MaxLifetime)

	return db, nil
}

// NewMySQLManager 创建MySQL管理器
func NewMySQLManager(c config.MySQLConfig) (*MySQLManager, error) {
	db, err := NewMySQLConnection(c)
	if err != nil {
		return nil, err
	}

	return &MySQLManager{
		db:     db,
		config: c,
	}, nil
}

// GetDB 获取数据库连接
func (m *MySQLManager) GetDB() *sql.DB {
	return m.db
}

// GetStats 获取连接池统计信息
func (m *MySQLManager) GetStats() sql.DBStats {
	return m.db.Stats()
}

// Close 关闭数据库连接
func (m *MySQLManager) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

// HealthCheck 数据库健康检查
func (m *MySQLManager) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var result int
	err := m.db.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("数据库健康检查失败: %w", err)
	}

	return nil
}

// GetDSN 获取MySQL连接DSN字符串
func GetDSN(c config.MySQLConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.Charset, c.ParseTime, c.Loc)
}

// DBTransaction 数据库事务封装
type DBTransaction struct {
	tx *sql.Tx
}

// NewTransaction 创建数据库事务
func (m *MySQLManager) NewTransaction(ctx context.Context) (*DBTransaction, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("创建事务失败: %w", err)
	}

	return &DBTransaction{tx: tx}, nil
}

// Commit 提交事务
func (dt *DBTransaction) Commit() error {
	return dt.tx.Commit()
}

// Rollback 回滚事务
func (dt *DBTransaction) Rollback() error {
	return dt.tx.Rollback()
}

// Exec 执行SQL语句
func (dt *DBTransaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return dt.tx.Exec(query, args...)
}

// Query 查询数据
func (dt *DBTransaction) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return dt.tx.Query(query, args...)
}

// QueryRow 查询单行数据
func (dt *DBTransaction) QueryRow(query string, args ...interface{}) *sql.Row {
	return dt.tx.QueryRow(query, args...)
}
