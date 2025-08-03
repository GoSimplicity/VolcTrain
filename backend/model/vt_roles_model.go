package model

import (
	"database/sql"
	"time"
)

// VtRoles 角色表模型
type VtRoles struct {
	Id          int64      `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	DisplayName string     `db:"display_name" json:"displayName"`
	Description string     `db:"description" json:"description"`
	RoleCode    string     `db:"role_code" json:"roleCode"`
	RoleType    string     `db:"role_type" json:"roleType"` // system, custom
	SortOrder   int        `db:"sort_order" json:"sortOrder"`
	Status      string     `db:"status" json:"status"` // active, inactive
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt"`
}

// VtRolesModel 角色模型操作接口
type VtRolesModel interface {
	Insert(data *VtRoles) (sql.Result, error)
	FindOne(id int64) (*VtRoles, error)
	FindOneByName(name string) (*VtRoles, error)
	Update(data *VtRoles) error
	Delete(id int64) error
	List(page, pageSize int, filters map[string]interface{}) ([]*VtRoles, int64, error)
}

type vtRolesModel struct {
	conn *sql.DB
}

func NewVtRolesModel(conn *sql.DB) VtRolesModel {
	return &vtRolesModel{conn: conn}
}

func (m *vtRolesModel) Insert(data *VtRoles) (sql.Result, error) {
	query := `INSERT INTO vt_roles (name, display_name, description, role_code, role_type, sort_order, status) VALUES (?, ?, ?, ?, ?, ?, ?)`
	return m.conn.Exec(query, data.Name, data.DisplayName, data.Description, data.RoleCode, data.RoleType, data.SortOrder, data.Status)
}

func (m *vtRolesModel) FindOne(id int64) (*VtRoles, error) {
	var role VtRoles
	query := `SELECT id, name, display_name, description, role_code, role_type, sort_order, status, created_at, updated_at, deleted_at FROM vt_roles WHERE id = ? AND deleted_at IS NULL`
	err := m.conn.QueryRow(query, id).Scan(&role.Id, &role.Name, &role.DisplayName, &role.Description, &role.RoleCode, &role.RoleType, &role.SortOrder, &role.Status, &role.CreatedAt, &role.UpdatedAt, &role.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (m *vtRolesModel) FindOneByName(name string) (*VtRoles, error) {
	var role VtRoles
	query := `SELECT id, name, display_name, description, role_code, role_type, sort_order, status, created_at, updated_at, deleted_at FROM vt_roles WHERE name = ? AND deleted_at IS NULL`
	err := m.conn.QueryRow(query, name).Scan(&role.Id, &role.Name, &role.DisplayName, &role.Description, &role.RoleCode, &role.RoleType, &role.SortOrder, &role.Status, &role.CreatedAt, &role.UpdatedAt, &role.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (m *vtRolesModel) Update(data *VtRoles) error {
	query := `UPDATE vt_roles SET display_name = ?, description = ?, role_code = ?, role_type = ?, sort_order = ?, status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := m.conn.Exec(query, data.DisplayName, data.Description, data.RoleCode, data.RoleType, data.SortOrder, data.Status, data.Id)
	return err
}

func (m *vtRolesModel) Delete(id int64) error {
	query := `UPDATE vt_roles SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := m.conn.Exec(query, id)
	return err
}

func (m *vtRolesModel) List(page, pageSize int, filters map[string]interface{}) ([]*VtRoles, int64, error) {
	offset := (page - 1) * pageSize

	whereClause := "WHERE deleted_at IS NULL"
	args := []interface{}{}

	if status, ok := filters["status"]; ok {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	if roleType, ok := filters["role_type"]; ok {
		whereClause += " AND role_type = ?"
		args = append(args, roleType)
	}

	if keyword, ok := filters["keyword"]; ok {
		whereClause += " AND (name LIKE ? OR display_name LIKE ?)"
		searchPattern := "%" + keyword.(string) + "%"
		args = append(args, searchPattern, searchPattern)
	}

	// 查询总数
	countQuery := "SELECT COUNT(*) FROM vt_roles " + whereClause
	var total int64
	err := m.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据列表
	listQuery := `SELECT id, name, display_name, description, role_code, role_type, sort_order, status, created_at, updated_at FROM vt_roles ` + whereClause + ` ORDER BY sort_order ASC, created_at DESC LIMIT ? OFFSET ?`
	listArgs := append(args, pageSize, offset)

	rows, err := m.conn.Query(listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var roles []*VtRoles
	for rows.Next() {
		var role VtRoles
		err := rows.Scan(&role.Id, &role.Name, &role.DisplayName, &role.Description, &role.RoleCode, &role.RoleType, &role.SortOrder, &role.Status, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		roles = append(roles, &role)
	}

	return roles, total, nil
}
