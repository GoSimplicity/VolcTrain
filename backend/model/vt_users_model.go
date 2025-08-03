package model

import (
	"database/sql"
	"time"
)

// VtUsers 用户表模型
type VtUsers struct {
	Id                int64      `db:"id" json:"id"`
	Username          string     `db:"username" json:"username"`
	Email             string     `db:"email" json:"email"`
	Phone             string     `db:"phone" json:"phone"`
	PasswordHash      string     `db:"password_hash" json:"-"` // 密码不返回给前端
	Salt              string     `db:"salt" json:"-"`          // 盐值不返回给前端
	RealName          string     `db:"real_name" json:"realName"`
	Nickname          string     `db:"nickname" json:"nickname"`
	Status            string     `db:"status" json:"status"`      // active, inactive, locked, pending
	UserType          string     `db:"user_type" json:"userType"` // admin, user, service
	LastLoginAt       *time.Time `db:"last_login_at" json:"lastLoginAt"`
	LastLoginIp       string     `db:"last_login_ip" json:"lastLoginIp"`
	PasswordExpiresAt *time.Time `db:"password_expires_at" json:"passwordExpiresAt"`
	LoginAttempts     int        `db:"login_attempts" json:"loginAttempts"`
	LockedUntil       *time.Time `db:"locked_until" json:"lockedUntil"`
	MfaEnabled        bool       `db:"mfa_enabled" json:"mfaEnabled"`
	MfaSecret         string     `db:"mfa_secret" json:"-"` // MFA密钥不返回给前端
	EmailVerified     bool       `db:"email_verified" json:"emailVerified"`
	PhoneVerified     bool       `db:"phone_verified" json:"phoneVerified"`
	CreatedAt         time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt         *time.Time `db:"deleted_at" json:"deletedAt"`
}

// VtUsersModel 用户模型操作接口
type VtUsersModel interface {
	Insert(data *VtUsers) (sql.Result, error)
	FindOne(id int64) (*VtUsers, error)
	FindOneByUsername(username string) (*VtUsers, error)
	FindOneByEmail(email string) (*VtUsers, error)
	Update(data *VtUsers) error
	Delete(id int64) error
	List(page, pageSize int, filters map[string]interface{}) ([]*VtUsers, int64, error)
}

// vtUsersModel 用户模型实现
type vtUsersModel struct {
	conn *sql.DB
}

// NewVtUsersModel 创建用户模型实例
func NewVtUsersModel(conn *sql.DB) VtUsersModel {
	return &vtUsersModel{
		conn: conn,
	}
}

// Insert 插入用户记录
func (m *vtUsersModel) Insert(data *VtUsers) (sql.Result, error) {
	query := `INSERT INTO vt_users (username, email, phone, password_hash, salt, real_name, nickname, status, user_type, mfa_enabled, email_verified, phone_verified) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	return m.conn.Exec(query, data.Username, data.Email, data.Phone, data.PasswordHash, data.Salt, data.RealName, data.Nickname, data.Status, data.UserType, data.MfaEnabled, data.EmailVerified, data.PhoneVerified)
}

// FindOne 根据ID查找用户
func (m *vtUsersModel) FindOne(id int64) (*VtUsers, error) {
	var user VtUsers
	query := `SELECT id, username, email, phone, password_hash, salt, real_name, nickname, status, user_type, last_login_at, last_login_ip, password_expires_at, login_attempts, locked_until, mfa_enabled, mfa_secret, email_verified, phone_verified, created_at, updated_at, deleted_at FROM vt_users WHERE id = ? AND deleted_at IS NULL`
	err := m.conn.QueryRow(query, id).Scan(
		&user.Id, &user.Username, &user.Email, &user.Phone, &user.PasswordHash, &user.Salt,
		&user.RealName, &user.Nickname, &user.Status, &user.UserType, &user.LastLoginAt,
		&user.LastLoginIp, &user.PasswordExpiresAt, &user.LoginAttempts, &user.LockedUntil,
		&user.MfaEnabled, &user.MfaSecret, &user.EmailVerified, &user.PhoneVerified,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindOneByUsername 根据用户名查找用户
func (m *vtUsersModel) FindOneByUsername(username string) (*VtUsers, error) {
	var user VtUsers
	query := `SELECT id, username, email, phone, password_hash, salt, real_name, nickname, status, user_type, last_login_at, last_login_ip, password_expires_at, login_attempts, locked_until, mfa_enabled, mfa_secret, email_verified, phone_verified, created_at, updated_at, deleted_at FROM vt_users WHERE username = ? AND deleted_at IS NULL`
	err := m.conn.QueryRow(query, username).Scan(
		&user.Id, &user.Username, &user.Email, &user.Phone, &user.PasswordHash, &user.Salt,
		&user.RealName, &user.Nickname, &user.Status, &user.UserType, &user.LastLoginAt,
		&user.LastLoginIp, &user.PasswordExpiresAt, &user.LoginAttempts, &user.LockedUntil,
		&user.MfaEnabled, &user.MfaSecret, &user.EmailVerified, &user.PhoneVerified,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindOneByEmail 根据邮箱查找用户
func (m *vtUsersModel) FindOneByEmail(email string) (*VtUsers, error) {
	var user VtUsers
	query := `SELECT id, username, email, phone, password_hash, salt, real_name, nickname, status, user_type, last_login_at, last_login_ip, password_expires_at, login_attempts, locked_until, mfa_enabled, mfa_secret, email_verified, phone_verified, created_at, updated_at, deleted_at FROM vt_users WHERE email = ? AND deleted_at IS NULL`
	err := m.conn.QueryRow(query, email).Scan(
		&user.Id, &user.Username, &user.Email, &user.Phone, &user.PasswordHash, &user.Salt,
		&user.RealName, &user.Nickname, &user.Status, &user.UserType, &user.LastLoginAt,
		&user.LastLoginIp, &user.PasswordExpiresAt, &user.LoginAttempts, &user.LockedUntil,
		&user.MfaEnabled, &user.MfaSecret, &user.EmailVerified, &user.PhoneVerified,
		&user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (m *vtUsersModel) Update(data *VtUsers) error {
	query := `UPDATE vt_users SET email = ?, phone = ?, real_name = ?, nickname = ?, status = ?, user_type = ?, last_login_at = ?, last_login_ip = ?, login_attempts = ?, locked_until = ?, mfa_enabled = ?, email_verified = ?, phone_verified = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := m.conn.Exec(query, data.Email, data.Phone, data.RealName, data.Nickname, data.Status, data.UserType, data.LastLoginAt, data.LastLoginIp, data.LoginAttempts, data.LockedUntil, data.MfaEnabled, data.EmailVerified, data.PhoneVerified, data.Id)
	return err
}

// Delete 软删除用户（设置deleted_at时间）
func (m *vtUsersModel) Delete(id int64) error {
	query := `UPDATE vt_users SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := m.conn.Exec(query, id)
	return err
}

// List 获取用户列表（支持分页和过滤）
func (m *vtUsersModel) List(page, pageSize int, filters map[string]interface{}) ([]*VtUsers, int64, error) {
	offset := (page - 1) * pageSize

	// 构建查询条件
	whereClause := "WHERE deleted_at IS NULL"
	args := []interface{}{}

	if status, ok := filters["status"]; ok {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	if userType, ok := filters["user_type"]; ok {
		whereClause += " AND user_type = ?"
		args = append(args, userType)
	}

	if keyword, ok := filters["keyword"]; ok {
		whereClause += " AND (username LIKE ? OR email LIKE ? OR real_name LIKE ?)"
		searchPattern := "%" + keyword.(string) + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}

	// 查询总数
	countQuery := "SELECT COUNT(*) FROM vt_users " + whereClause
	var total int64
	err := m.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据列表
	listQuery := `SELECT id, username, email, phone, real_name, nickname, status, user_type, last_login_at, last_login_ip, login_attempts, locked_until, mfa_enabled, email_verified, phone_verified, created_at, updated_at FROM vt_users ` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	listArgs := append(args, pageSize, offset)

	rows, err := m.conn.Query(listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*VtUsers
	for rows.Next() {
		var user VtUsers
		err := rows.Scan(
			&user.Id, &user.Username, &user.Email, &user.Phone, &user.RealName,
			&user.Nickname, &user.Status, &user.UserType, &user.LastLoginAt,
			&user.LastLoginIp, &user.LoginAttempts, &user.LockedUntil,
			&user.MfaEnabled, &user.EmailVerified, &user.PhoneVerified,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
