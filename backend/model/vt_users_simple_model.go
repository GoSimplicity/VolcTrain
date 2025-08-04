package model

import (
	"context"
	"database/sql"
	"time"
)

// VtUsersSimple 简化的用户表模型，匹配当前数据库结构
type VtUsersSimple struct {
	Id          int64      `db:"id" json:"id"`
	Username    string     `db:"username" json:"username"`
	Password    string     `db:"password_hash" json:"-"` // 映射到password_hash字段
	Email       string     `db:"email" json:"email"`
	RealName    string     `db:"real_name" json:"realName"`
	Status      string     `db:"status" json:"status"`
	UserType    string     `db:"user_type" json:"userType"`
	Department  string     `db:"department" json:"department"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
	LastLoginAt time.Time  `db:"last_login_at" json:"lastLoginAt"`
}

// VtUsersSimpleModel 简化用户模型操作接口
type VtUsersSimpleModel interface {
	FindOne(ctx context.Context, id int64) (*VtUsersSimple, error)
	FindByUsername(ctx context.Context, username string) (*VtUsersSimple, error)
	Update(ctx context.Context, data *VtUsersSimple) error
}

// vtUsersSimpleModel 简化用户模型实现
type vtUsersSimpleModel struct {
	conn *sql.DB
}

// NewVtUsersSimpleModel 创建简化用户模型实例
func NewVtUsersSimpleModel(conn *sql.DB) VtUsersSimpleModel {
	return &vtUsersSimpleModel{
		conn: conn,
	}
}

// FindOne 根据ID查找用户
func (m *vtUsersSimpleModel) FindOne(ctx context.Context, id int64) (*VtUsersSimple, error) {
	var user VtUsersSimple
	query := `SELECT id, username, password_hash, email, real_name, status, user_type, department, created_at, updated_at, COALESCE(last_login_at, '1970-01-01 00:00:01') FROM vt_users WHERE id = ?`
	err := m.conn.QueryRowContext(ctx, query, id).Scan(
		&user.Id, &user.Username, &user.Password, &user.Email, &user.RealName,
		&user.Status, &user.UserType, &user.Department, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (m *vtUsersSimpleModel) FindByUsername(ctx context.Context, username string) (*VtUsersSimple, error) {
	var user VtUsersSimple
	query := `SELECT id, username, password_hash, email, real_name, status, user_type, department, created_at, updated_at, COALESCE(last_login_at, '1970-01-01 00:00:01') FROM vt_users WHERE username = ?`
	err := m.conn.QueryRowContext(ctx, query, username).Scan(
		&user.Id, &user.Username, &user.Password, &user.Email, &user.RealName,
		&user.Status, &user.UserType, &user.Department, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (m *vtUsersSimpleModel) Update(ctx context.Context, data *VtUsersSimple) error {
	query := `UPDATE vt_users SET username=?, email=?, real_name=?, status=?, user_type=?, department=?, last_login_at=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`
	_, err := m.conn.ExecContext(ctx, query, data.Username, data.Email, data.RealName, data.Status, data.UserType, data.Department, data.LastLoginAt, data.Id)
	return err
}