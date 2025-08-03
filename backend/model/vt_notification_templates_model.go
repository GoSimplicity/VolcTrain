package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// VtNotificationTemplates 通知模板表模型
type VtNotificationTemplates struct {
	Id           int64     `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	DisplayName  string    `db:"display_name" json:"displayName"`
	Description  string    `db:"description" json:"description"`
	ChannelType  string    `db:"channel_type" json:"channelType"`
	TemplateType string    `db:"template_type" json:"templateType"`
	Subject      string    `db:"subject" json:"subject"`
	Content      string    `db:"content" json:"content"`
	Variables    string    `db:"variables" json:"variables"`
	IsDefault    bool      `db:"is_default" json:"isDefault"`
	IsBuiltin    bool      `db:"is_builtin" json:"isBuiltin"`
	Status       string    `db:"status" json:"status"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}

// VtNotificationTemplatesModel 通知模板模型操作接口
type VtNotificationTemplatesModel interface {
	Insert(data *VtNotificationTemplates) (sql.Result, error)
	FindOne(id int64) (*VtNotificationTemplates, error)
	FindOneByName(name string) (*VtNotificationTemplates, error)
	Update(data *VtNotificationTemplates) error
	Delete(id int64) error
	List(page, pageSize int, filters map[string]interface{}) ([]*VtNotificationTemplates, int64, error)
	FindActiveTemplates() ([]*VtNotificationTemplates, error)
	FindByChannelType(channelType string) ([]*VtNotificationTemplates, error)
	FindDefaultTemplate(channelType, templateType string) (*VtNotificationTemplates, error)
	UpdateStatus(id int64, status string) error
	GetBuiltinTemplates() ([]*VtNotificationTemplates, error)
	SearchTemplates(keyword string, offset, limit int) ([]*VtNotificationTemplates, int64, error)
}

// vtNotificationTemplatesModel 通知模板模型实现
type vtNotificationTemplatesModel struct {
	conn *sql.DB
}

// NewVtNotificationTemplatesModel 创建通知模板模型实例
func NewVtNotificationTemplatesModel(conn *sql.DB) VtNotificationTemplatesModel {
	return &vtNotificationTemplatesModel{conn: conn}
}

// Insert 插入通知模板
func (m *vtNotificationTemplatesModel) Insert(data *VtNotificationTemplates) (sql.Result, error) {
	query := `INSERT INTO vt_notification_templates (
		name, display_name, description, channel_type, template_type, subject,
		content, variables, is_default, is_builtin, status
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return m.conn.Exec(query,
		data.Name, data.DisplayName, data.Description, data.ChannelType,
		data.TemplateType, data.Subject, data.Content, data.Variables,
		data.IsDefault, data.IsBuiltin, data.Status,
	)
}

// FindOne 根据ID查找通知模板
func (m *vtNotificationTemplatesModel) FindOne(id int64) (*VtNotificationTemplates, error) {
	query := `SELECT id, name, display_name, description, channel_type, template_type,
		subject, content, variables, is_default, is_builtin, status, created_at, updated_at
		FROM vt_notification_templates WHERE id = ?`

	var template VtNotificationTemplates
	err := m.conn.QueryRow(query, id).Scan(
		&template.Id, &template.Name, &template.DisplayName, &template.Description,
		&template.ChannelType, &template.TemplateType, &template.Subject,
		&template.Content, &template.Variables, &template.IsDefault,
		&template.IsBuiltin, &template.Status, &template.CreatedAt, &template.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &template, nil
}

// FindOneByName 根据名称查找通知模板
func (m *vtNotificationTemplatesModel) FindOneByName(name string) (*VtNotificationTemplates, error) {
	query := `SELECT id, name, display_name, description, channel_type, template_type,
		subject, content, variables, is_default, is_builtin, status, created_at, updated_at
		FROM vt_notification_templates WHERE name = ?`

	var template VtNotificationTemplates
	err := m.conn.QueryRow(query, name).Scan(
		&template.Id, &template.Name, &template.DisplayName, &template.Description,
		&template.ChannelType, &template.TemplateType, &template.Subject,
		&template.Content, &template.Variables, &template.IsDefault,
		&template.IsBuiltin, &template.Status, &template.CreatedAt, &template.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &template, nil
}

// Update 更新通知模板
func (m *vtNotificationTemplatesModel) Update(data *VtNotificationTemplates) error {
	query := `UPDATE vt_notification_templates SET 
		name = ?, display_name = ?, description = ?, channel_type = ?,
		template_type = ?, subject = ?, content = ?, variables = ?,
		is_default = ?, status = ?, updated_at = NOW() WHERE id = ?`

	_, err := m.conn.Exec(query,
		data.Name, data.DisplayName, data.Description, data.ChannelType,
		data.TemplateType, data.Subject, data.Content, data.Variables,
		data.IsDefault, data.Status, data.Id,
	)
	return err
}

// Delete 删除通知模板
func (m *vtNotificationTemplatesModel) Delete(id int64) error {
	query := "DELETE FROM vt_notification_templates WHERE id = ? AND is_builtin = 0"
	_, err := m.conn.Exec(query, id)
	return err
}

// List 分页查询通知模板
func (m *vtNotificationTemplatesModel) List(page, pageSize int, filters map[string]interface{}) ([]*VtNotificationTemplates, int64, error) {
	offset := (page - 1) * pageSize

	// 构建WHERE条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if name, ok := filters["name"]; ok && name != "" {
		whereClause += " AND name LIKE ?"
		args = append(args, "%"+name.(string)+"%")
	}

	if channelType, ok := filters["channel_type"]; ok && channelType != "" {
		whereClause += " AND channel_type = ?"
		args = append(args, channelType)
	}

	if templateType, ok := filters["template_type"]; ok && templateType != "" {
		whereClause += " AND template_type = ?"
		args = append(args, templateType)
	}

	if status, ok := filters["status"]; ok && status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	if isBuiltin, ok := filters["is_builtin"]; ok {
		whereClause += " AND is_builtin = ?"
		args = append(args, isBuiltin)
	}

	// 查询总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM vt_notification_templates " + whereClause
	err := m.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	query := `SELECT id, name, display_name, description, channel_type, template_type,
		subject, content, variables, is_default, is_builtin, status, created_at, updated_at
		FROM vt_notification_templates ` + whereClause + ` 
		ORDER BY is_default DESC, is_builtin DESC, name LIMIT ? OFFSET ?`

	args = append(args, pageSize, offset)
	rows, err := m.conn.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var templates []*VtNotificationTemplates
	for rows.Next() {
		var template VtNotificationTemplates
		err := rows.Scan(
			&template.Id, &template.Name, &template.DisplayName, &template.Description,
			&template.ChannelType, &template.TemplateType, &template.Subject,
			&template.Content, &template.Variables, &template.IsDefault,
			&template.IsBuiltin, &template.Status, &template.CreatedAt, &template.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		templates = append(templates, &template)
	}

	return templates, total, nil
}

// FindActiveTemplates 查找所有活跃的通知模板
func (m *vtNotificationTemplatesModel) FindActiveTemplates() ([]*VtNotificationTemplates, error) {
	query := `SELECT id, name, display_name, description, channel_type, template_type,
		subject, content, variables, is_default, is_builtin, status, created_at, updated_at
		FROM vt_notification_templates WHERE status = 'active' 
		ORDER BY is_default DESC, is_builtin DESC, channel_type, name`

	rows, err := m.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*VtNotificationTemplates
	for rows.Next() {
		var template VtNotificationTemplates
		err := rows.Scan(
			&template.Id, &template.Name, &template.DisplayName, &template.Description,
			&template.ChannelType, &template.TemplateType, &template.Subject,
			&template.Content, &template.Variables, &template.IsDefault,
			&template.IsBuiltin, &template.Status, &template.CreatedAt, &template.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, &template)
	}

	return templates, nil
}

// FindByChannelType 根据渠道类型查找通知模板
func (m *vtNotificationTemplatesModel) FindByChannelType(channelType string) ([]*VtNotificationTemplates, error) {
	query := `SELECT id, name, display_name, description, channel_type, template_type,
		subject, content, variables, is_default, is_builtin, status, created_at, updated_at
		FROM vt_notification_templates WHERE channel_type = ? AND status = 'active'
		ORDER BY is_default DESC, name`

	rows, err := m.conn.Query(query, channelType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*VtNotificationTemplates
	for rows.Next() {
		var template VtNotificationTemplates
		err := rows.Scan(
			&template.Id, &template.Name, &template.DisplayName, &template.Description,
			&template.ChannelType, &template.TemplateType, &template.Subject,
			&template.Content, &template.Variables, &template.IsDefault,
			&template.IsBuiltin, &template.Status, &template.CreatedAt, &template.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, &template)
	}

	return templates, nil
}

// FindDefaultTemplate 查找默认模板
func (m *vtNotificationTemplatesModel) FindDefaultTemplate(channelType, templateType string) (*VtNotificationTemplates, error) {
	query := `SELECT id, name, display_name, description, channel_type, template_type,
		subject, content, variables, is_default, is_builtin, status, created_at, updated_at
		FROM vt_notification_templates 
		WHERE channel_type = ? AND template_type = ? AND is_default = 1 AND status = 'active'
		ORDER BY is_builtin DESC
		LIMIT 1`

	var template VtNotificationTemplates
	err := m.conn.QueryRow(query, channelType, templateType).Scan(
		&template.Id, &template.Name, &template.DisplayName, &template.Description,
		&template.ChannelType, &template.TemplateType, &template.Subject,
		&template.Content, &template.Variables, &template.IsDefault,
		&template.IsBuiltin, &template.Status, &template.CreatedAt, &template.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &template, nil
}

// UpdateStatus 更新模板状态
func (m *vtNotificationTemplatesModel) UpdateStatus(id int64, status string) error {
	query := "UPDATE vt_notification_templates SET status = ?, updated_at = NOW() WHERE id = ?"
	_, err := m.conn.Exec(query, status, id)
	return err
}

// GetBuiltinTemplates 获取内置模板
func (m *vtNotificationTemplatesModel) GetBuiltinTemplates() ([]*VtNotificationTemplates, error) {
	query := `SELECT id, name, display_name, description, channel_type, template_type,
		subject, content, variables, is_default, is_builtin, status, created_at, updated_at
		FROM vt_notification_templates WHERE is_builtin = 1 AND status = 'active'
		ORDER BY channel_type, name`

	rows, err := m.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*VtNotificationTemplates
	for rows.Next() {
		var template VtNotificationTemplates
		err := rows.Scan(
			&template.Id, &template.Name, &template.DisplayName, &template.Description,
			&template.ChannelType, &template.TemplateType, &template.Subject,
			&template.Content, &template.Variables, &template.IsDefault,
			&template.IsBuiltin, &template.Status, &template.CreatedAt, &template.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, &template)
	}

	return templates, nil
}

// SearchTemplates 搜索模板
func (m *vtNotificationTemplatesModel) SearchTemplates(keyword string, offset, limit int) ([]*VtNotificationTemplates, int64, error) {
	whereClause := "WHERE status = 'active'"
	args := []interface{}{}

	if keyword != "" {
		whereClause += " AND (name LIKE ? OR display_name LIKE ? OR description LIKE ?)"
		likeKeyword := "%" + keyword + "%"
		args = append(args, likeKeyword, likeKeyword, likeKeyword)
	}

	// 查询总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM vt_notification_templates " + whereClause
	err := m.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	query := `SELECT id, name, display_name, description, channel_type, template_type,
		subject, content, variables, is_default, is_builtin, status, created_at, updated_at
		FROM vt_notification_templates ` + whereClause + ` 
		ORDER BY is_default DESC, is_builtin DESC, name LIMIT ? OFFSET ?`

	args = append(args, limit, offset)
	rows, err := m.conn.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var templates []*VtNotificationTemplates
	for rows.Next() {
		var template VtNotificationTemplates
		err := rows.Scan(
			&template.Id, &template.Name, &template.DisplayName, &template.Description,
			&template.ChannelType, &template.TemplateType, &template.Subject,
			&template.Content, &template.Variables, &template.IsDefault,
			&template.IsBuiltin, &template.Status, &template.CreatedAt, &template.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		templates = append(templates, &template)
	}

	return templates, total, nil
}

// GetVariablesMap 获取变量映射
func (t *VtNotificationTemplates) GetVariablesMap() (map[string]interface{}, error) {
	if t.Variables == "" {
		return make(map[string]interface{}), nil
	}

	var variables map[string]interface{}
	err := json.Unmarshal([]byte(t.Variables), &variables)
	if err != nil {
		return nil, fmt.Errorf("解析模板变量失败: %v", err)
	}

	return variables, nil
}

// SetVariables 设置变量
func (t *VtNotificationTemplates) SetVariables(variables map[string]interface{}) error {
	if variables == nil || len(variables) == 0 {
		t.Variables = ""
		return nil
	}

	data, err := json.Marshal(variables)
	if err != nil {
		return fmt.Errorf("序列化模板变量失败: %v", err)
	}

	t.Variables = string(data)
	return nil
}

// GetAvailableVariables 获取可用变量列表
func (t *VtNotificationTemplates) GetAvailableVariables() []string {
	content := t.Subject + " " + t.Content
	variables := []string{}

	// 简单的变量提取（实际应使用正则表达式）
	standardVars := []string{
		"{{.Alert.RuleName}}",
		"{{.Alert.Message}}",
		"{{.Alert.AlertLevel}}",
		"{{.Alert.TriggerValue}}",
		"{{.Alert.ThresholdValue}}",
		"{{.Alert.TriggeredAt}}",
		"{{.Alert.ResourceType}}",
		"{{.Alert.ResourceName}}",
		"{{.Action}}",
		"{{.Timestamp}}",
		"{{.Date}}",
		"{{.Time}}",
	}

	for _, variable := range standardVars {
		if strings.Contains(content, variable) {
			variables = append(variables, variable)
		}
	}

	return variables
}

// ValidateTemplate 验证模板
func (t *VtNotificationTemplates) ValidateTemplate() error {
	if t.Name == "" {
		return fmt.Errorf("模板名称不能为空")
	}
	if t.ChannelType == "" {
		return fmt.Errorf("渠道类型不能为空")
	}
	if t.TemplateType == "" {
		return fmt.Errorf("模板类型不能为空")
	}
	if t.Content == "" {
		return fmt.Errorf("模板内容不能为空")
	}

	// 验证渠道类型
	validChannelTypes := []string{"email", "sms", "dingtalk", "webhook"}
	channelTypeValid := false
	for _, validType := range validChannelTypes {
		if t.ChannelType == validType {
			channelTypeValid = true
			break
		}
	}
	if !channelTypeValid {
		return fmt.Errorf("不支持的渠道类型: %s", t.ChannelType)
	}

	// 验证模板类型
	validTemplateTypes := []string{"firing", "resolved", "default"}
	templateTypeValid := false
	for _, validType := range validTemplateTypes {
		if t.TemplateType == validType {
			templateTypeValid = true
			break
		}
	}
	if !templateTypeValid {
		return fmt.Errorf("不支持的模板类型: %s", t.TemplateType)
	}

	return nil
}
