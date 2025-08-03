package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// VtNotificationChannels 通知渠道表模型
type VtNotificationChannels struct {
	Id          int64      `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	DisplayName string     `db:"display_name" json:"displayName"`
	Description string     `db:"description" json:"description"`
	ChannelType string     `db:"channel_type" json:"channelType"`
	Config      string     `db:"config" json:"config"`
	Recipients  string     `db:"recipients" json:"recipients"`
	IsDefault   bool       `db:"is_default" json:"isDefault"`
	Status      string     `db:"status" json:"status"`
	TestStatus  string     `db:"test_status" json:"testStatus"`
	LastTestAt  *time.Time `db:"last_test_at" json:"lastTestAt"`
	ErrorCount  int        `db:"error_count" json:"errorCount"`
	LastError   string     `db:"last_error" json:"lastError"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
}

// VtNotificationChannelsModel 通知渠道模型操作接口
type VtNotificationChannelsModel interface {
	Insert(data *VtNotificationChannels) (sql.Result, error)
	FindOne(id int64) (*VtNotificationChannels, error)
	FindOneByName(name string) (*VtNotificationChannels, error)
	Update(data *VtNotificationChannels) error
	Delete(id int64) error
	List(page, pageSize int, filters map[string]interface{}) ([]*VtNotificationChannels, int64, error)
	FindActiveChannels() ([]*VtNotificationChannels, error)
	FindByChannelType(channelType string) ([]*VtNotificationChannels, error)
	UpdateStatus(id int64, status string) error
	UpdateTestResult(id int64, testStatus, errorMsg string) error
	IncrementErrorCount(id int64) error
	ResetErrorCount(id int64) error
}

// vtNotificationChannelsModel 通知渠道模型实现
type vtNotificationChannelsModel struct {
	conn *sql.DB
}

// NewVtNotificationChannelsModel 创建通知渠道模型实例
func NewVtNotificationChannelsModel(conn *sql.DB) VtNotificationChannelsModel {
	return &vtNotificationChannelsModel{conn: conn}
}

// Insert 插入通知渠道
func (m *vtNotificationChannelsModel) Insert(data *VtNotificationChannels) (sql.Result, error) {
	query := `INSERT INTO vt_notification_channels (
		name, display_name, description, channel_type, config, recipients,
		is_default, status, test_status
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return m.conn.Exec(query,
		data.Name, data.DisplayName, data.Description, data.ChannelType,
		data.Config, data.Recipients, data.IsDefault, data.Status, data.TestStatus,
	)
}

// FindOne 根据ID查找通知渠道
func (m *vtNotificationChannelsModel) FindOne(id int64) (*VtNotificationChannels, error) {
	query := `SELECT id, name, display_name, description, channel_type, config, recipients,
		is_default, status, test_status, last_test_at, error_count, last_error, created_at, updated_at
		FROM vt_notification_channels WHERE id = ?`

	var channel VtNotificationChannels
	err := m.conn.QueryRow(query, id).Scan(
		&channel.Id, &channel.Name, &channel.DisplayName, &channel.Description,
		&channel.ChannelType, &channel.Config, &channel.Recipients,
		&channel.IsDefault, &channel.Status, &channel.TestStatus, &channel.LastTestAt,
		&channel.ErrorCount, &channel.LastError, &channel.CreatedAt, &channel.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// FindOneByName 根据名称查找通知渠道
func (m *vtNotificationChannelsModel) FindOneByName(name string) (*VtNotificationChannels, error) {
	query := `SELECT id, name, display_name, description, channel_type, config, recipients,
		is_default, status, test_status, last_test_at, error_count, last_error, created_at, updated_at
		FROM vt_notification_channels WHERE name = ?`

	var channel VtNotificationChannels
	err := m.conn.QueryRow(query, name).Scan(
		&channel.Id, &channel.Name, &channel.DisplayName, &channel.Description,
		&channel.ChannelType, &channel.Config, &channel.Recipients,
		&channel.IsDefault, &channel.Status, &channel.TestStatus, &channel.LastTestAt,
		&channel.ErrorCount, &channel.LastError, &channel.CreatedAt, &channel.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// Update 更新通知渠道
func (m *vtNotificationChannelsModel) Update(data *VtNotificationChannels) error {
	query := `UPDATE vt_notification_channels SET 
		name = ?, display_name = ?, description = ?, channel_type = ?,
		config = ?, recipients = ?, is_default = ?, status = ?,
		updated_at = NOW() WHERE id = ?`

	_, err := m.conn.Exec(query,
		data.Name, data.DisplayName, data.Description, data.ChannelType,
		data.Config, data.Recipients, data.IsDefault, data.Status, data.Id,
	)
	return err
}

// Delete 删除通知渠道
func (m *vtNotificationChannelsModel) Delete(id int64) error {
	query := "DELETE FROM vt_notification_channels WHERE id = ?"
	_, err := m.conn.Exec(query, id)
	return err
}

// List 分页查询通知渠道
func (m *vtNotificationChannelsModel) List(page, pageSize int, filters map[string]interface{}) ([]*VtNotificationChannels, int64, error) {
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

	if status, ok := filters["status"]; ok && status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	// 查询总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM vt_notification_channels " + whereClause
	err := m.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	query := `SELECT id, name, display_name, description, channel_type, config, recipients,
		is_default, status, test_status, last_test_at, error_count, last_error, created_at, updated_at
		FROM vt_notification_channels ` + whereClause + ` 
		ORDER BY is_default DESC, name LIMIT ? OFFSET ?`

	args = append(args, pageSize, offset)
	rows, err := m.conn.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var channels []*VtNotificationChannels
	for rows.Next() {
		var channel VtNotificationChannels
		err := rows.Scan(
			&channel.Id, &channel.Name, &channel.DisplayName, &channel.Description,
			&channel.ChannelType, &channel.Config, &channel.Recipients,
			&channel.IsDefault, &channel.Status, &channel.TestStatus, &channel.LastTestAt,
			&channel.ErrorCount, &channel.LastError, &channel.CreatedAt, &channel.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		channels = append(channels, &channel)
	}

	return channels, total, nil
}

// FindActiveChannels 查找所有活跃的通知渠道
func (m *vtNotificationChannelsModel) FindActiveChannels() ([]*VtNotificationChannels, error) {
	query := `SELECT id, name, display_name, description, channel_type, config,
		status, created_at, updated_at
		FROM vt_notification_channels WHERE status = 'active' 
		ORDER BY name`

	rows, err := m.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []*VtNotificationChannels
	for rows.Next() {
		var channel VtNotificationChannels
		err := rows.Scan(
			&channel.Id, &channel.Name, &channel.DisplayName, &channel.Description,
			&channel.ChannelType, &channel.Config,
			&channel.Status, &channel.CreatedAt, &channel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		// 设置默认值
		channel.Recipients = ""
		channel.IsDefault = false
		channel.TestStatus = ""
		channel.ErrorCount = 0
		channel.LastError = ""
		channels = append(channels, &channel)
	}

	return channels, nil
}

// FindByChannelType 根据渠道类型查找通知渠道
func (m *vtNotificationChannelsModel) FindByChannelType(channelType string) ([]*VtNotificationChannels, error) {
	query := `SELECT id, name, display_name, description, channel_type, config, recipients,
		is_default, status, test_status, last_test_at, error_count, last_error, created_at, updated_at
		FROM vt_notification_channels WHERE channel_type = ? AND status = 'active'
		ORDER BY name`

	rows, err := m.conn.Query(query, channelType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []*VtNotificationChannels
	for rows.Next() {
		var channel VtNotificationChannels
		err := rows.Scan(
			&channel.Id, &channel.Name, &channel.DisplayName, &channel.Description,
			&channel.ChannelType, &channel.Config, &channel.Recipients,
			&channel.IsDefault, &channel.Status, &channel.TestStatus, &channel.LastTestAt,
			&channel.ErrorCount, &channel.LastError, &channel.CreatedAt, &channel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		channels = append(channels, &channel)
	}

	return channels, nil
}

// UpdateStatus 更新渠道状态
func (m *vtNotificationChannelsModel) UpdateStatus(id int64, status string) error {
	query := "UPDATE vt_notification_channels SET status = ?, updated_at = NOW() WHERE id = ?"
	_, err := m.conn.Exec(query, status, id)
	return err
}

// UpdateTestResult 更新测试结果
func (m *vtNotificationChannelsModel) UpdateTestResult(id int64, testStatus, errorMsg string) error {
	query := `UPDATE vt_notification_channels SET 
		test_status = ?, last_test_at = NOW(), last_error = ?, updated_at = NOW() 
		WHERE id = ?`
	_, err := m.conn.Exec(query, testStatus, errorMsg, id)
	return err
}

// IncrementErrorCount 增加错误计数
func (m *vtNotificationChannelsModel) IncrementErrorCount(id int64) error {
	query := "UPDATE vt_notification_channels SET error_count = error_count + 1, updated_at = NOW() WHERE id = ?"
	_, err := m.conn.Exec(query, id)
	return err
}

// ResetErrorCount 重置错误计数
func (m *vtNotificationChannelsModel) ResetErrorCount(id int64) error {
	query := "UPDATE vt_notification_channels SET error_count = 0, updated_at = NOW() WHERE id = ?"
	_, err := m.conn.Exec(query, id)
	return err
}

// GetConfigMap 获取配置映射
func (c *VtNotificationChannels) GetConfigMap() (map[string]interface{}, error) {
	if c.Config == "" {
		return make(map[string]interface{}), nil
	}

	var config map[string]interface{}
	err := json.Unmarshal([]byte(c.Config), &config)
	if err != nil {
		return nil, fmt.Errorf("解析渠道配置失败: %v", err)
	}

	return config, nil
}

// SetConfig 设置配置
func (c *VtNotificationChannels) SetConfig(config map[string]interface{}) error {
	if config == nil || len(config) == 0 {
		c.Config = ""
		return nil
	}

	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化渠道配置失败: %v", err)
	}

	c.Config = string(data)
	return nil
}

// GetRecipientsList 获取收件人列表
func (c *VtNotificationChannels) GetRecipientsList() ([]string, error) {
	if c.Recipients == "" {
		return []string{}, nil
	}

	var recipients []string
	err := json.Unmarshal([]byte(c.Recipients), &recipients)
	if err != nil {
		// 兼容逗号分隔的格式
		recipients = strings.Split(c.Recipients, ",")
		for i, recipient := range recipients {
			recipients[i] = strings.TrimSpace(recipient)
		}
	}

	return recipients, nil
}

// SetRecipients 设置收件人
func (c *VtNotificationChannels) SetRecipients(recipients []string) error {
	if recipients == nil || len(recipients) == 0 {
		c.Recipients = ""
		return nil
	}

	data, err := json.Marshal(recipients)
	if err != nil {
		return fmt.Errorf("序列化收件人列表失败: %v", err)
	}

	c.Recipients = string(data)
	return nil
}
