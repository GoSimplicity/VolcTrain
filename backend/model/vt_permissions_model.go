package model

import (
	"database/sql"
	"time"
)

// VtPermissions 权限表模型
type VtPermissions struct {
	Id             int64     `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	DisplayName    string    `db:"display_name" json:"displayName"`
	Description    string    `db:"description" json:"description"`
	Module         string    `db:"module" json:"module"`
	Action         string    `db:"action" json:"action"`
	Resource       string    `db:"resource" json:"resource"`
	PermissionCode string    `db:"permission_code" json:"permissionCode"`
	ParentId       int64     `db:"parent_id" json:"parentId"`
	Level          int       `db:"level" json:"level"`
	SortOrder      int       `db:"sort_order" json:"sortOrder"`
	Status         string    `db:"status" json:"status"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt"`
}

// VtPermissionsModel 权限模型操作接口
type VtPermissionsModel interface {
	Insert(data *VtPermissions) (sql.Result, error)
	FindOne(id int64) (*VtPermissions, error)
	FindOneByName(name string) (*VtPermissions, error)
	Update(data *VtPermissions) error
	Delete(id int64) error
	List(filters map[string]interface{}) ([]*VtPermissions, error)
	GetPermissionTree() ([]*VtPermissions, error)
}

type vtPermissionsModel struct {
	conn *sql.DB
}

func NewVtPermissionsModel(conn *sql.DB) VtPermissionsModel {
	return &vtPermissionsModel{conn: conn}
}

func (m *vtPermissionsModel) Insert(data *VtPermissions) (sql.Result, error) {
	query := `INSERT INTO vt_permissions (name, display_name, description, module, action, resource, permission_code, parent_id, level, sort_order, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	return m.conn.Exec(query, data.Name, data.DisplayName, data.Description, data.Module, data.Action, data.Resource, data.PermissionCode, data.ParentId, data.Level, data.SortOrder, data.Status)
}

func (m *vtPermissionsModel) FindOne(id int64) (*VtPermissions, error) {
	var perm VtPermissions
	query := `SELECT id, name, display_name, description, module, action, resource, permission_code, parent_id, level, sort_order, status, created_at, updated_at FROM vt_permissions WHERE id = ?`
	err := m.conn.QueryRow(query, id).Scan(&perm.Id, &perm.Name, &perm.DisplayName, &perm.Description, &perm.Module, &perm.Action, &perm.Resource, &perm.PermissionCode, &perm.ParentId, &perm.Level, &perm.SortOrder, &perm.Status, &perm.CreatedAt, &perm.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

func (m *vtPermissionsModel) FindOneByName(name string) (*VtPermissions, error) {
	var perm VtPermissions
	query := `SELECT id, name, display_name, description, module, action, resource, permission_code, parent_id, level, sort_order, status, created_at, updated_at FROM vt_permissions WHERE name = ?`
	err := m.conn.QueryRow(query, name).Scan(&perm.Id, &perm.Name, &perm.DisplayName, &perm.Description, &perm.Module, &perm.Action, &perm.Resource, &perm.PermissionCode, &perm.ParentId, &perm.Level, &perm.SortOrder, &perm.Status, &perm.CreatedAt, &perm.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

func (m *vtPermissionsModel) Update(data *VtPermissions) error {
	query := `UPDATE vt_permissions SET display_name = ?, description = ?, module = ?, action = ?, resource = ?, permission_code = ?, parent_id = ?, sort_order = ?, status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := m.conn.Exec(query, data.DisplayName, data.Description, data.Module, data.Action, data.Resource, data.PermissionCode, data.ParentId, data.SortOrder, data.Status, data.Id)
	return err
}

func (m *vtPermissionsModel) Delete(id int64) error {
	query := `DELETE FROM vt_permissions WHERE id = ?`
	_, err := m.conn.Exec(query, id)
	return err
}

func (m *vtPermissionsModel) List(filters map[string]interface{}) ([]*VtPermissions, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if status, ok := filters["status"]; ok {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	if module, ok := filters["module"]; ok {
		whereClause += " AND module = ?"
		args = append(args, module)
	}

	if parentId, ok := filters["parent_id"]; ok {
		whereClause += " AND parent_id = ?"
		args = append(args, parentId)
	}

	query := `SELECT id, name, display_name, description, module, action, resource, permission_code, parent_id, level, sort_order, status, created_at, updated_at FROM vt_permissions ` + whereClause + ` ORDER BY level ASC, sort_order ASC`

	rows, err := m.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*VtPermissions
	for rows.Next() {
		var perm VtPermissions
		err := rows.Scan(&perm.Id, &perm.Name, &perm.DisplayName, &perm.Description, &perm.Module, &perm.Action, &perm.Resource, &perm.PermissionCode, &perm.ParentId, &perm.Level, &perm.SortOrder, &perm.Status, &perm.CreatedAt, &perm.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &perm)
	}

	return permissions, nil
}

func (m *vtPermissionsModel) GetPermissionTree() ([]*VtPermissions, error) {
	query := `SELECT id, name, display_name, description, module, action, resource, permission_code, parent_id, level, sort_order, status, created_at, updated_at FROM vt_permissions WHERE status = 'active' ORDER BY level ASC, sort_order ASC`

	rows, err := m.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*VtPermissions
	for rows.Next() {
		var perm VtPermissions
		err := rows.Scan(&perm.Id, &perm.Name, &perm.DisplayName, &perm.Description, &perm.Module, &perm.Action, &perm.Resource, &perm.PermissionCode, &perm.ParentId, &perm.Level, &perm.SortOrder, &perm.Status, &perm.CreatedAt, &perm.UpdatedAt)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &perm)
	}

	return permissions, nil
}
