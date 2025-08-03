package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// VtAlertRules 告警规则表模型
type VtAlertRules struct {
	Id                          int64     `db:"id" json:"id"`
	Name                        string    `db:"name" json:"name"`
	DisplayName                 string    `db:"display_name" json:"displayName"`
	Description                 string    `db:"description" json:"description"`
	RuleType                    string    `db:"rule_type" json:"ruleType"`
	ConditionExpression         string    `db:"condition_expression" json:"conditionExpression"`
	QueryExpression             string    `db:"query_expression" json:"queryExpression"`
	WarningThreshold            float64   `db:"warning_threshold" json:"warningThreshold"`
	CriticalThreshold           float64   `db:"critical_threshold" json:"criticalThreshold"`
	ThresholdCondition          string    `db:"threshold_condition" json:"thresholdCondition"`
	EvaluationWindowSeconds     int       `db:"evaluation_window_seconds" json:"evaluationWindowSeconds"`
	EvaluationIntervalSeconds   int       `db:"evaluation_interval_seconds" json:"evaluationIntervalSeconds"`
	TriggerDurationSeconds      int       `db:"trigger_duration_seconds" json:"triggerDurationSeconds"`
	RecoveryDurationSeconds     int       `db:"recovery_duration_seconds" json:"recoveryDurationSeconds"`
	FilterLabels                string    `db:"filter_labels" json:"filterLabels"`
	FilterResources             string    `db:"filter_resources" json:"filterResources"`
	AlertLevel                  string    `db:"alert_level" json:"alertLevel"`
	SeverityScore               int       `db:"severity_score" json:"severityScore"`
	NotificationChannels        string    `db:"notification_channels" json:"notificationChannels"`
	NotificationThrottleMinutes int       `db:"notification_throttle_minutes" json:"notificationThrottleMinutes"`
	SilenceDurationSeconds      int       `db:"silence_duration_seconds" json:"silenceDurationSeconds"`
	SuppressionRules            string    `db:"suppression_rules" json:"suppressionRules"`
	DependencyRules             string    `db:"dependency_rules" json:"dependencyRules"`
	Status                      string    `db:"status" json:"status"`
	IsBuiltin                   bool      `db:"is_builtin" json:"isBuiltin"`
	LastEvaluationAt            time.Time `db:"last_evaluation_at" json:"lastEvaluationAt"`
	LastTriggerAt               time.Time `db:"last_trigger_at" json:"lastTriggerAt"`
	TriggerCount                int       `db:"trigger_count" json:"triggerCount"`
	CreatedAt                   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt                   time.Time `db:"updated_at" json:"updatedAt"`
}

// VtAlertRulesModel 告警规则模型操作接口
type VtAlertRulesModel interface {
	Insert(data *VtAlertRules) (sql.Result, error)
	FindOne(id int64) (*VtAlertRules, error)
	FindOneByName(name string) (*VtAlertRules, error)
	Update(data *VtAlertRules) error
	Delete(id int64) error
	List(page, pageSize int, filters map[string]interface{}) ([]*VtAlertRules, int64, error)
	FindActiveRules() ([]*VtAlertRules, error)
	FindByRuleType(ruleType string) ([]*VtAlertRules, error)
	UpdateStatus(id int64, status string) error
	UpdateEvaluationTime(id int64, evaluationTime time.Time) error
	IncrementTriggerCount(id int64) error
	GetBuiltinRules() ([]*VtAlertRules, error)
	SearchRules(keyword string, offset, limit int) ([]*VtAlertRules, int64, error)
}

// vtAlertRulesModel 告警规则模型实现
type vtAlertRulesModel struct {
	conn *sql.DB
}

// NewVtAlertRulesModel 创建告警规则模型实例
func NewVtAlertRulesModel(conn *sql.DB) VtAlertRulesModel {
	return &vtAlertRulesModel{conn: conn}
}

// Insert 插入告警规则
func (m *vtAlertRulesModel) Insert(data *VtAlertRules) (sql.Result, error) {
	query := `INSERT INTO vt_alert_rules (
		name, display_name, description, rule_type, condition_expression, query_expression,
		warning_threshold, critical_threshold, threshold_condition, evaluation_window_seconds,
		evaluation_interval_seconds, trigger_duration_seconds, recovery_duration_seconds,
		filter_labels, filter_resources, alert_level, severity_score, notification_channels,
		notification_throttle_minutes, silence_duration_seconds, suppression_rules,
		dependency_rules, status, is_builtin
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return m.conn.Exec(query,
		data.Name, data.DisplayName, data.Description, data.RuleType,
		data.ConditionExpression, data.QueryExpression, data.WarningThreshold,
		data.CriticalThreshold, data.ThresholdCondition, data.EvaluationWindowSeconds,
		data.EvaluationIntervalSeconds, data.TriggerDurationSeconds,
		data.RecoveryDurationSeconds, data.FilterLabels, data.FilterResources,
		data.AlertLevel, data.SeverityScore, data.NotificationChannels,
		data.NotificationThrottleMinutes, data.SilenceDurationSeconds,
		data.SuppressionRules, data.DependencyRules, data.Status, data.IsBuiltin,
	)
}

// FindOne 根据ID查找告警规则
func (m *vtAlertRulesModel) FindOne(id int64) (*VtAlertRules, error) {
	query := `SELECT id, name, display_name, description, rule_type, condition_expression,
		query_expression, warning_threshold, critical_threshold, threshold_condition,
		evaluation_window_seconds, evaluation_interval_seconds, trigger_duration_seconds,
		recovery_duration_seconds, filter_labels, filter_resources, alert_level,
		severity_score, notification_channels, notification_throttle_minutes,
		silence_duration_seconds, suppression_rules, dependency_rules, status,
		is_builtin, last_evaluation_at, last_trigger_at, trigger_count, created_at, updated_at
		FROM vt_alert_rules WHERE id = ?`

	var rule VtAlertRules
	err := m.conn.QueryRow(query, id).Scan(
		&rule.Id, &rule.Name, &rule.DisplayName, &rule.Description, &rule.RuleType,
		&rule.ConditionExpression, &rule.QueryExpression, &rule.WarningThreshold,
		&rule.CriticalThreshold, &rule.ThresholdCondition, &rule.EvaluationWindowSeconds,
		&rule.EvaluationIntervalSeconds, &rule.TriggerDurationSeconds,
		&rule.RecoveryDurationSeconds, &rule.FilterLabels, &rule.FilterResources,
		&rule.AlertLevel, &rule.SeverityScore, &rule.NotificationChannels,
		&rule.NotificationThrottleMinutes, &rule.SilenceDurationSeconds,
		&rule.SuppressionRules, &rule.DependencyRules, &rule.Status, &rule.IsBuiltin,
		&rule.LastEvaluationAt, &rule.LastTriggerAt, &rule.TriggerCount,
		&rule.CreatedAt, &rule.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// FindOneByName 根据名称查找告警规则
func (m *vtAlertRulesModel) FindOneByName(name string) (*VtAlertRules, error) {
	query := `SELECT id, name, display_name, description, rule_type, condition_expression,
		query_expression, warning_threshold, critical_threshold, threshold_condition,
		evaluation_window_seconds, evaluation_interval_seconds, trigger_duration_seconds,
		recovery_duration_seconds, filter_labels, filter_resources, alert_level,
		severity_score, notification_channels, notification_throttle_minutes,
		silence_duration_seconds, suppression_rules, dependency_rules, status,
		is_builtin, last_evaluation_at, last_trigger_at, trigger_count, created_at, updated_at
		FROM vt_alert_rules WHERE name = ?`

	var rule VtAlertRules
	err := m.conn.QueryRow(query, name).Scan(
		&rule.Id, &rule.Name, &rule.DisplayName, &rule.Description, &rule.RuleType,
		&rule.ConditionExpression, &rule.QueryExpression, &rule.WarningThreshold,
		&rule.CriticalThreshold, &rule.ThresholdCondition, &rule.EvaluationWindowSeconds,
		&rule.EvaluationIntervalSeconds, &rule.TriggerDurationSeconds,
		&rule.RecoveryDurationSeconds, &rule.FilterLabels, &rule.FilterResources,
		&rule.AlertLevel, &rule.SeverityScore, &rule.NotificationChannels,
		&rule.NotificationThrottleMinutes, &rule.SilenceDurationSeconds,
		&rule.SuppressionRules, &rule.DependencyRules, &rule.Status, &rule.IsBuiltin,
		&rule.LastEvaluationAt, &rule.LastTriggerAt, &rule.TriggerCount,
		&rule.CreatedAt, &rule.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// Update 更新告警规则
func (m *vtAlertRulesModel) Update(data *VtAlertRules) error {
	query := `UPDATE vt_alert_rules SET 
		name = ?, display_name = ?, description = ?, rule_type = ?,
		condition_expression = ?, query_expression = ?, warning_threshold = ?,
		critical_threshold = ?, threshold_condition = ?, evaluation_window_seconds = ?,
		evaluation_interval_seconds = ?, trigger_duration_seconds = ?,
		recovery_duration_seconds = ?, filter_labels = ?, filter_resources = ?,
		alert_level = ?, severity_score = ?, notification_channels = ?,
		notification_throttle_minutes = ?, silence_duration_seconds = ?,
		suppression_rules = ?, dependency_rules = ?, status = ?, updated_at = NOW()
		WHERE id = ?`

	_, err := m.conn.Exec(query,
		data.Name, data.DisplayName, data.Description, data.RuleType,
		data.ConditionExpression, data.QueryExpression, data.WarningThreshold,
		data.CriticalThreshold, data.ThresholdCondition, data.EvaluationWindowSeconds,
		data.EvaluationIntervalSeconds, data.TriggerDurationSeconds,
		data.RecoveryDurationSeconds, data.FilterLabels, data.FilterResources,
		data.AlertLevel, data.SeverityScore, data.NotificationChannels,
		data.NotificationThrottleMinutes, data.SilenceDurationSeconds,
		data.SuppressionRules, data.DependencyRules, data.Status, data.Id,
	)
	return err
}

// Delete 删除告警规则
func (m *vtAlertRulesModel) Delete(id int64) error {
	query := "DELETE FROM vt_alert_rules WHERE id = ? AND is_builtin = 0"
	_, err := m.conn.Exec(query, id)
	return err
}

// List 分页查询告警规则
func (m *vtAlertRulesModel) List(page, pageSize int, filters map[string]interface{}) ([]*VtAlertRules, int64, error) {
	offset := (page - 1) * pageSize

	// 构建WHERE条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if name, ok := filters["name"]; ok && name != "" {
		whereClause += " AND name LIKE ?"
		args = append(args, "%"+name.(string)+"%")
	}

	if ruleType, ok := filters["rule_type"]; ok && ruleType != "" {
		whereClause += " AND rule_type = ?"
		args = append(args, ruleType)
	}

	if alertLevel, ok := filters["alert_level"]; ok && alertLevel != "" {
		whereClause += " AND alert_level = ?"
		args = append(args, alertLevel)
	}

	if status, ok := filters["status"]; ok && status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	if isBuiltin, ok := filters["is_builtin"]; ok {
		// 忽略 is_builtin 过滤条件，因为字段不存在
		_ = isBuiltin
	}

	// 查询总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM vt_alert_rules " + whereClause
	err := m.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	query := `SELECT id, name, display_name, description, rule_type, condition_expression,
		query_expression, warning_threshold, critical_threshold, threshold_condition,
		evaluation_window_seconds, evaluation_interval_seconds, alert_level,
		severity_score, notification_channels, status, last_evaluation_at,
		last_trigger_at, trigger_count, created_at, updated_at
		FROM vt_alert_rules ` + whereClause + ` 
		ORDER BY severity_score DESC, name LIMIT ? OFFSET ?`

	args = append(args, pageSize, offset)
	rows, err := m.conn.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var rules []*VtAlertRules
	for rows.Next() {
		var rule VtAlertRules
		err := rows.Scan(
			&rule.Id, &rule.Name, &rule.DisplayName, &rule.Description, &rule.RuleType,
			&rule.ConditionExpression, &rule.QueryExpression, &rule.WarningThreshold,
			&rule.CriticalThreshold, &rule.ThresholdCondition, &rule.EvaluationWindowSeconds,
			&rule.EvaluationIntervalSeconds, &rule.AlertLevel, &rule.SeverityScore,
			&rule.NotificationChannels, &rule.Status,
			&rule.LastEvaluationAt, &rule.LastTriggerAt, &rule.TriggerCount,
			&rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		// 为了兼容性，设置默认值
		rule.IsBuiltin = false
		rules = append(rules, &rule)
	}

	return rules, total, nil
}

// FindActiveRules 查找所有活跃的告警规则
func (m *vtAlertRulesModel) FindActiveRules() ([]*VtAlertRules, error) {
	query := `SELECT id, name, display_name, description, rule_type, condition_expression,
		query_expression, warning_threshold, critical_threshold, threshold_condition,
		evaluation_window_seconds, evaluation_interval_seconds, trigger_duration_seconds,
		recovery_duration_seconds, filter_labels, filter_resources, alert_level,
		severity_score, notification_channels, notification_throttle_minutes,
		silence_duration_seconds, suppression_rules, dependency_rules, status,
		last_evaluation_at, last_trigger_at, trigger_count, created_at, updated_at
		FROM vt_alert_rules WHERE status = 'active' 
		ORDER BY severity_score DESC, name`

	rows, err := m.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*VtAlertRules
	for rows.Next() {
		var rule VtAlertRules
		err := rows.Scan(
			&rule.Id, &rule.Name, &rule.DisplayName, &rule.Description, &rule.RuleType,
			&rule.ConditionExpression, &rule.QueryExpression, &rule.WarningThreshold,
			&rule.CriticalThreshold, &rule.ThresholdCondition, &rule.EvaluationWindowSeconds,
			&rule.EvaluationIntervalSeconds, &rule.TriggerDurationSeconds,
			&rule.RecoveryDurationSeconds, &rule.FilterLabels, &rule.FilterResources,
			&rule.AlertLevel, &rule.SeverityScore, &rule.NotificationChannels,
			&rule.NotificationThrottleMinutes, &rule.SilenceDurationSeconds,
			&rule.SuppressionRules, &rule.DependencyRules, &rule.Status,
			&rule.LastEvaluationAt, &rule.LastTriggerAt, &rule.TriggerCount,
			&rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		// 为了兼容性，设置默认值
		rule.IsBuiltin = false
		rules = append(rules, &rule)
	}

	return rules, nil
}

// FindByRuleType 根据规则类型查找告警规则
func (m *vtAlertRulesModel) FindByRuleType(ruleType string) ([]*VtAlertRules, error) {
	query := `SELECT id, name, display_name, description, rule_type, condition_expression,
		alert_level, severity_score, status, created_at
		FROM vt_alert_rules WHERE rule_type = ? AND status = 'active'
		ORDER BY severity_score DESC, name`

	rows, err := m.conn.Query(query, ruleType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*VtAlertRules
	for rows.Next() {
		var rule VtAlertRules
		err := rows.Scan(
			&rule.Id, &rule.Name, &rule.DisplayName, &rule.Description,
			&rule.RuleType, &rule.ConditionExpression, &rule.AlertLevel,
			&rule.SeverityScore, &rule.Status, &rule.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		rules = append(rules, &rule)
	}

	return rules, nil
}

// UpdateStatus 更新规则状态
func (m *vtAlertRulesModel) UpdateStatus(id int64, status string) error {
	query := "UPDATE vt_alert_rules SET status = ?, updated_at = NOW() WHERE id = ?"
	_, err := m.conn.Exec(query, status, id)
	return err
}

// UpdateEvaluationTime 更新评估时间
func (m *vtAlertRulesModel) UpdateEvaluationTime(id int64, evaluationTime time.Time) error {
	query := "UPDATE vt_alert_rules SET last_evaluation_at = ?, updated_at = NOW() WHERE id = ?"
	_, err := m.conn.Exec(query, evaluationTime, id)
	return err
}

// IncrementTriggerCount 增加触发次数
func (m *vtAlertRulesModel) IncrementTriggerCount(id int64) error {
	query := `UPDATE vt_alert_rules SET 
		trigger_count = trigger_count + 1, last_trigger_at = NOW(), updated_at = NOW() 
		WHERE id = ?`
	_, err := m.conn.Exec(query, id)
	return err
}

// GetBuiltinRules 获取内置规则
func (m *vtAlertRulesModel) GetBuiltinRules() ([]*VtAlertRules, error) {
	query := `SELECT id, name, display_name, description, rule_type, condition_expression,
		alert_level, severity_score, status, created_at
		FROM vt_alert_rules WHERE is_builtin = 1 AND status = 'active'
		ORDER BY severity_score DESC, name`

	rows, err := m.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*VtAlertRules
	for rows.Next() {
		var rule VtAlertRules
		err := rows.Scan(
			&rule.Id, &rule.Name, &rule.DisplayName, &rule.Description,
			&rule.RuleType, &rule.ConditionExpression, &rule.AlertLevel,
			&rule.SeverityScore, &rule.Status, &rule.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		rules = append(rules, &rule)
	}

	return rules, nil
}

// SearchRules 搜索规则
func (m *vtAlertRulesModel) SearchRules(keyword string, offset, limit int) ([]*VtAlertRules, int64, error) {
	whereClause := "WHERE status = 'active'"
	args := []interface{}{}

	if keyword != "" {
		whereClause += " AND (name LIKE ? OR display_name LIKE ? OR description LIKE ?)"
		likeKeyword := "%" + keyword + "%"
		args = append(args, likeKeyword, likeKeyword, likeKeyword)
	}

	// 查询总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM vt_alert_rules " + whereClause
	err := m.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	query := `SELECT id, name, display_name, description, rule_type, condition_expression,
		alert_level, severity_score, status, is_builtin, last_evaluation_at,
		last_trigger_at, trigger_count, created_at, updated_at
		FROM vt_alert_rules ` + whereClause + ` 
		ORDER BY is_builtin DESC, severity_score DESC, name LIMIT ? OFFSET ?`

	args = append(args, limit, offset)
	rows, err := m.conn.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var rules []*VtAlertRules
	for rows.Next() {
		var rule VtAlertRules
		err := rows.Scan(
			&rule.Id, &rule.Name, &rule.DisplayName, &rule.Description, &rule.RuleType,
			&rule.ConditionExpression, &rule.AlertLevel, &rule.SeverityScore,
			&rule.Status, &rule.IsBuiltin, &rule.LastEvaluationAt,
			&rule.LastTriggerAt, &rule.TriggerCount, &rule.CreatedAt, &rule.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		rules = append(rules, &rule)
	}

	return rules, total, nil
}

// GetFilterLabelsMap 获取过滤标签映射
func (r *VtAlertRules) GetFilterLabelsMap() (map[string]interface{}, error) {
	if r.FilterLabels == "" {
		return make(map[string]interface{}), nil
	}

	var labels map[string]interface{}
	err := json.Unmarshal([]byte(r.FilterLabels), &labels)
	if err != nil {
		return nil, fmt.Errorf("解析过滤标签失败: %v", err)
	}

	return labels, nil
}

// SetFilterLabels 设置过滤标签
func (r *VtAlertRules) SetFilterLabels(labels map[string]interface{}) error {
	if labels == nil || len(labels) == 0 {
		r.FilterLabels = ""
		return nil
	}

	data, err := json.Marshal(labels)
	if err != nil {
		return fmt.Errorf("序列化过滤标签失败: %v", err)
	}

	r.FilterLabels = string(data)
	return nil
}

// GetNotificationChannelsList 获取通知渠道列表
func (r *VtAlertRules) GetNotificationChannelsList() ([]string, error) {
	if r.NotificationChannels == "" {
		return []string{}, nil
	}

	var channels []string
	err := json.Unmarshal([]byte(r.NotificationChannels), &channels)
	if err != nil {
		return nil, fmt.Errorf("解析通知渠道失败: %v", err)
	}

	return channels, nil
}

// SetNotificationChannels 设置通知渠道
func (r *VtAlertRules) SetNotificationChannels(channels []string) error {
	if channels == nil || len(channels) == 0 {
		r.NotificationChannels = ""
		return nil
	}

	data, err := json.Marshal(channels)
	if err != nil {
		return fmt.Errorf("序列化通知渠道失败: %v", err)
	}

	r.NotificationChannels = string(data)
	return nil
}
