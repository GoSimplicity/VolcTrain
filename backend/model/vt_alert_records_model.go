package model

import (
	"database/sql"
	"time"
)

// VtAlertRecords 告警记录表模型
type VtAlertRecords struct {
	Id                  int64     `db:"id" json:"id"`
	RuleId              int64     `db:"rule_id" json:"ruleId"`
	AlertId             string    `db:"alert_id" json:"alertId"`
	AlertName           string    `db:"alert_name" json:"alertName"`
	AlertLevel          string    `db:"alert_level" json:"alertLevel"`
	SeverityScore       int       `db:"severity_score" json:"severityScore"`
	Message             string    `db:"message" json:"message"`
	Summary             string    `db:"summary" json:"summary"`
	ResourceType        string    `db:"resource_type" json:"resourceType"`
	ResourceId          int64     `db:"resource_id" json:"resourceId"`
	ResourceName        string    `db:"resource_name" json:"resourceName"`
	InstanceId          string    `db:"instance_id" json:"instanceId"`
	TriggerValue        float64   `db:"trigger_value" json:"triggerValue"`
	ThresholdValue      float64   `db:"threshold_value" json:"thresholdValue"`
	ConditionExpression string    `db:"condition_expression" json:"conditionExpression"`
	EvaluationData      string    `db:"evaluation_data" json:"evaluationData"`
	Labels              string    `db:"labels" json:"labels"`
	Annotations         string    `db:"annotations" json:"annotations"`
	Context             string    `db:"context" json:"context"`
	Status              string    `db:"status" json:"status"`
	TriggeredAt         time.Time `db:"triggered_at" json:"triggeredAt"`
	FirstOccurrenceAt   time.Time `db:"first_occurrence_at" json:"firstOccurrenceAt"`
	LastOccurrenceAt    time.Time `db:"last_occurrence_at" json:"lastOccurrenceAt"`
	OccurrenceCount     int       `db:"occurrence_count" json:"occurrenceCount"`
	NotificationSent    bool      `db:"notification_sent" json:"notificationSent"`
	NotificationCount   int       `db:"notification_count" json:"notificationCount"`
	LastNotificationAt  time.Time `db:"last_notification_at" json:"lastNotificationAt"`
	EscalationLevel     int       `db:"escalation_level" json:"escalationLevel"`
	AlertGroupId        string    `db:"alert_group_id" json:"alertGroupId"`
	CorrelationId       string    `db:"correlation_id" json:"correlationId"`
	CreatedAt           time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt           time.Time `db:"updated_at" json:"updatedAt"`
}

// VtAlertRecordsModel 告警记录模型操作接口
type VtAlertRecordsModel interface {
	Insert(data *VtAlertRecords) (sql.Result, error)
	FindOne(id int64) (*VtAlertRecords, error)
	FindOneByAlertId(alertId string) (*VtAlertRecords, error)
	Update(data *VtAlertRecords) error
	Delete(id int64) error
	List(page, pageSize int, filters map[string]interface{}) ([]*VtAlertRecords, int64, error)
	FindByRuleId(ruleId int64, limit int) ([]*VtAlertRecords, error)
	FindByResourceType(resourceType string, limit int) ([]*VtAlertRecords, error)
	FindByTimeRange(startTime, endTime time.Time) ([]*VtAlertRecords, error)
	FindActiveAlerts() ([]*VtAlertRecords, error)
	UpdateByAlertId(alertId, status string, occurrenceCount int, lastOccurrenceAt time.Time) error
	DeleteOldRecords(beforeTime time.Time) (int64, error)
	GetAlertStatistics(startTime, endTime time.Time) (map[string]interface{}, error)
	FindRecentAlerts(limit int) ([]*VtAlertRecords, error)
}

// vtAlertRecordsModel 告警记录模型实现
type vtAlertRecordsModel struct {
	conn *sql.DB
}

// NewVtAlertRecordsModel 创建告警记录模型实例
func NewVtAlertRecordsModel(conn *sql.DB) VtAlertRecordsModel {
	return &vtAlertRecordsModel{conn: conn}
}

// Insert 插入告警记录
func (m *vtAlertRecordsModel) Insert(data *VtAlertRecords) (sql.Result, error) {
	query := `INSERT INTO vt_alert_records (
		rule_id, alert_id, alert_name, alert_level, severity_score, message, summary,
		resource_type, resource_id, resource_name, instance_id, trigger_value, threshold_value,
		condition_expression, evaluation_data, labels, annotations, context, status,
		triggered_at, first_occurrence_at, last_occurrence_at, occurrence_count,
		notification_sent, notification_count, last_notification_at, escalation_level,
		alert_group_id, correlation_id
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return m.conn.Exec(query,
		data.RuleId, data.AlertId, data.AlertName, data.AlertLevel, data.SeverityScore,
		data.Message, data.Summary, data.ResourceType, data.ResourceId, data.ResourceName,
		data.InstanceId, data.TriggerValue, data.ThresholdValue, data.ConditionExpression,
		data.EvaluationData, data.Labels, data.Annotations, data.Context, data.Status,
		data.TriggeredAt, data.FirstOccurrenceAt, data.LastOccurrenceAt, data.OccurrenceCount,
		data.NotificationSent, data.NotificationCount, data.LastNotificationAt,
		data.EscalationLevel, data.AlertGroupId, data.CorrelationId,
	)
}

// FindOne 根据ID查找告警记录
func (m *vtAlertRecordsModel) FindOne(id int64) (*VtAlertRecords, error) {
	query := `SELECT id, rule_id, alert_id, alert_name, alert_level, severity_score, message, summary,
		resource_type, resource_id, resource_name, instance_id, trigger_value, threshold_value,
		condition_expression, evaluation_data, labels, annotations, context, status,
		triggered_at, first_occurrence_at, last_occurrence_at, occurrence_count,
		notification_sent, notification_count, last_notification_at, escalation_level,
		alert_group_id, correlation_id, created_at, updated_at
		FROM vt_alert_records WHERE id = ?`

	var record VtAlertRecords
	err := m.conn.QueryRow(query, id).Scan(
		&record.Id, &record.RuleId, &record.AlertId, &record.AlertName, &record.AlertLevel,
		&record.SeverityScore, &record.Message, &record.Summary, &record.ResourceType,
		&record.ResourceId, &record.ResourceName, &record.InstanceId, &record.TriggerValue,
		&record.ThresholdValue, &record.ConditionExpression, &record.EvaluationData,
		&record.Labels, &record.Annotations, &record.Context, &record.Status,
		&record.TriggeredAt, &record.FirstOccurrenceAt, &record.LastOccurrenceAt,
		&record.OccurrenceCount, &record.NotificationSent, &record.NotificationCount,
		&record.LastNotificationAt, &record.EscalationLevel, &record.AlertGroupId,
		&record.CorrelationId, &record.CreatedAt, &record.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &record, nil
}

// FindOneByAlertId 根据告警ID查找告警记录
func (m *vtAlertRecordsModel) FindOneByAlertId(alertId string) (*VtAlertRecords, error) {
	query := `SELECT id, rule_id, alert_id, alert_name, alert_level, severity_score, message, summary,
		resource_type, resource_id, resource_name, instance_id, trigger_value, threshold_value,
		condition_expression, evaluation_data, labels, annotations, context, status,
		triggered_at, first_occurrence_at, last_occurrence_at, occurrence_count,
		notification_sent, notification_count, last_notification_at, escalation_level,
		alert_group_id, correlation_id, created_at, updated_at
		FROM vt_alert_records WHERE alert_id = ?`

	var record VtAlertRecords
	err := m.conn.QueryRow(query, alertId).Scan(
		&record.Id, &record.RuleId, &record.AlertId, &record.AlertName, &record.AlertLevel,
		&record.SeverityScore, &record.Message, &record.Summary, &record.ResourceType,
		&record.ResourceId, &record.ResourceName, &record.InstanceId, &record.TriggerValue,
		&record.ThresholdValue, &record.ConditionExpression, &record.EvaluationData,
		&record.Labels, &record.Annotations, &record.Context, &record.Status,
		&record.TriggeredAt, &record.FirstOccurrenceAt, &record.LastOccurrenceAt,
		&record.OccurrenceCount, &record.NotificationSent, &record.NotificationCount,
		&record.LastNotificationAt, &record.EscalationLevel, &record.AlertGroupId,
		&record.CorrelationId, &record.CreatedAt, &record.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &record, nil
}

// Update 更新告警记录
func (m *vtAlertRecordsModel) Update(data *VtAlertRecords) error {
	query := `UPDATE vt_alert_records SET 
		alert_level = ?, message = ?, summary = ?, trigger_value = ?,
		threshold_value = ?, status = ?, last_occurrence_at = ?,
		occurrence_count = ?, notification_sent = ?, notification_count = ?,
		last_notification_at = ?, escalation_level = ?, updated_at = NOW()
		WHERE id = ?`

	_, err := m.conn.Exec(query,
		data.AlertLevel, data.Message, data.Summary, data.TriggerValue,
		data.ThresholdValue, data.Status, data.LastOccurrenceAt,
		data.OccurrenceCount, data.NotificationSent, data.NotificationCount,
		data.LastNotificationAt, data.EscalationLevel, data.Id,
	)
	return err
}

// Delete 删除告警记录
func (m *vtAlertRecordsModel) Delete(id int64) error {
	query := "DELETE FROM vt_alert_records WHERE id = ?"
	_, err := m.conn.Exec(query, id)
	return err
}

// List 分页查询告警记录
func (m *vtAlertRecordsModel) List(page, pageSize int, filters map[string]interface{}) ([]*VtAlertRecords, int64, error) {
	offset := (page - 1) * pageSize

	// 构建WHERE条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if ruleId, ok := filters["rule_id"]; ok {
		whereClause += " AND rule_id = ?"
		args = append(args, ruleId)
	}

	if alertLevel, ok := filters["alert_level"]; ok && alertLevel != "" {
		whereClause += " AND alert_level = ?"
		args = append(args, alertLevel)
	}

	if status, ok := filters["status"]; ok && status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	if resourceType, ok := filters["resource_type"]; ok && resourceType != "" {
		whereClause += " AND resource_type = ?"
		args = append(args, resourceType)
	}

	if keyword, ok := filters["keyword"]; ok && keyword != "" {
		whereClause += " AND (alert_name LIKE ? OR message LIKE ? OR resource_name LIKE ?)"
		searchPattern := "%" + keyword.(string) + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}

	if startTime, ok := filters["start_time"]; ok {
		whereClause += " AND triggered_at >= ?"
		args = append(args, startTime)
	}

	if endTime, ok := filters["end_time"]; ok {
		whereClause += " AND triggered_at <= ?"
		args = append(args, endTime)
	}

	// 查询总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM vt_alert_records " + whereClause
	err := m.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	query := `SELECT id, rule_id, alert_id, alert_name, alert_level, severity_score, message, summary,
		resource_type, resource_id, resource_name, instance_id, trigger_value, threshold_value,
		condition_expression, evaluation_data, labels, annotations, context, status,
		triggered_at, first_occurrence_at, last_occurrence_at, occurrence_count,
		notification_sent, notification_count, last_notification_at, escalation_level,
		alert_group_id, correlation_id, created_at, updated_at
		FROM vt_alert_records ` + whereClause + ` 
		ORDER BY triggered_at DESC LIMIT ? OFFSET ?`

	args = append(args, pageSize, offset)
	rows, err := m.conn.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var records []*VtAlertRecords
	for rows.Next() {
		var record VtAlertRecords
		err := rows.Scan(
			&record.Id, &record.RuleId, &record.AlertId, &record.AlertName, &record.AlertLevel,
			&record.SeverityScore, &record.Message, &record.Summary, &record.ResourceType,
			&record.ResourceId, &record.ResourceName, &record.InstanceId, &record.TriggerValue,
			&record.ThresholdValue, &record.ConditionExpression, &record.EvaluationData,
			&record.Labels, &record.Annotations, &record.Context, &record.Status,
			&record.TriggeredAt, &record.FirstOccurrenceAt, &record.LastOccurrenceAt,
			&record.OccurrenceCount, &record.NotificationSent, &record.NotificationCount,
			&record.LastNotificationAt, &record.EscalationLevel, &record.AlertGroupId,
			&record.CorrelationId, &record.CreatedAt, &record.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		records = append(records, &record)
	}

	return records, total, nil
}

// FindByRuleId 根据规则ID查找告警记录
func (m *vtAlertRecordsModel) FindByRuleId(ruleId int64, limit int) ([]*VtAlertRecords, error) {
	query := `SELECT id, rule_id, alert_id, alert_name, alert_level, severity_score, message,
		resource_type, resource_name, trigger_value, threshold_value, status,
		triggered_at, occurrence_count, created_at
		FROM vt_alert_records WHERE rule_id = ? 
		ORDER BY triggered_at DESC LIMIT ?`

	rows, err := m.conn.Query(query, ruleId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*VtAlertRecords
	for rows.Next() {
		var record VtAlertRecords
		err := rows.Scan(
			&record.Id, &record.RuleId, &record.AlertId, &record.AlertName,
			&record.AlertLevel, &record.SeverityScore, &record.Message,
			&record.ResourceType, &record.ResourceName, &record.TriggerValue,
			&record.ThresholdValue, &record.Status, &record.TriggeredAt,
			&record.OccurrenceCount, &record.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, nil
}

// FindByResourceType 根据资源类型查找告警记录
func (m *vtAlertRecordsModel) FindByResourceType(resourceType string, limit int) ([]*VtAlertRecords, error) {
	query := `SELECT id, rule_id, alert_id, alert_name, alert_level, resource_type,
		resource_name, trigger_value, status, triggered_at
		FROM vt_alert_records WHERE resource_type = ? 
		ORDER BY triggered_at DESC LIMIT ?`

	rows, err := m.conn.Query(query, resourceType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*VtAlertRecords
	for rows.Next() {
		var record VtAlertRecords
		err := rows.Scan(
			&record.Id, &record.RuleId, &record.AlertId, &record.AlertName,
			&record.AlertLevel, &record.ResourceType, &record.ResourceName,
			&record.TriggerValue, &record.Status, &record.TriggeredAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, nil
}

// FindByTimeRange 根据时间范围查找告警记录
func (m *vtAlertRecordsModel) FindByTimeRange(startTime, endTime time.Time) ([]*VtAlertRecords, error) {
	query := `SELECT id, rule_id, alert_id, alert_name, alert_level, resource_type,
		resource_name, trigger_value, status, triggered_at
		FROM vt_alert_records 
		WHERE triggered_at BETWEEN ? AND ? 
		ORDER BY triggered_at DESC`

	rows, err := m.conn.Query(query, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*VtAlertRecords
	for rows.Next() {
		var record VtAlertRecords
		err := rows.Scan(
			&record.Id, &record.RuleId, &record.AlertId, &record.AlertName,
			&record.AlertLevel, &record.ResourceType, &record.ResourceName,
			&record.TriggerValue, &record.Status, &record.TriggeredAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, nil
}

// FindActiveAlerts 查找活跃的告警
func (m *vtAlertRecordsModel) FindActiveAlerts() ([]*VtAlertRecords, error) {
	query := `SELECT id, rule_id, alert_id, alert_name, alert_level, severity_score,
		resource_type, resource_name, trigger_value, threshold_value, status,
		triggered_at, occurrence_count
		FROM vt_alert_records 
		WHERE status = 'firing' 
		ORDER BY severity_score DESC, triggered_at DESC`

	rows, err := m.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*VtAlertRecords
	for rows.Next() {
		var record VtAlertRecords
		err := rows.Scan(
			&record.Id, &record.RuleId, &record.AlertId, &record.AlertName,
			&record.AlertLevel, &record.SeverityScore, &record.ResourceType,
			&record.ResourceName, &record.TriggerValue, &record.ThresholdValue,
			&record.Status, &record.TriggeredAt, &record.OccurrenceCount,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, nil
}

// UpdateByAlertId 根据告警ID更新记录
func (m *vtAlertRecordsModel) UpdateByAlertId(alertId, status string, occurrenceCount int, lastOccurrenceAt time.Time) error {
	query := `UPDATE vt_alert_records SET 
		status = ?, occurrence_count = ?, last_occurrence_at = ?, updated_at = NOW()
		WHERE alert_id = ?`

	_, err := m.conn.Exec(query, status, occurrenceCount, lastOccurrenceAt, alertId)
	return err
}

// DeleteOldRecords 删除旧记录
func (m *vtAlertRecordsModel) DeleteOldRecords(beforeTime time.Time) (int64, error) {
	query := "DELETE FROM vt_alert_records WHERE created_at < ?"
	result, err := m.conn.Exec(query, beforeTime)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// GetAlertStatistics 获取告警统计
func (m *vtAlertRecordsModel) GetAlertStatistics(startTime, endTime time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总告警数量
	var totalCount int64
	err := m.conn.QueryRow("SELECT COUNT(*) FROM vt_alert_records WHERE triggered_at BETWEEN ? AND ?",
		startTime, endTime).Scan(&totalCount)
	if err != nil {
		return nil, err
	}
	stats["total_count"] = totalCount

	// 按级别统计
	levelQuery := `SELECT alert_level, COUNT(*) as count 
		FROM vt_alert_records 
		WHERE triggered_at BETWEEN ? AND ? 
		GROUP BY alert_level`

	rows, err := m.conn.Query(levelQuery, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	levelStats := make(map[string]int64)
	for rows.Next() {
		var level string
		var count int64
		if err := rows.Scan(&level, &count); err != nil {
			continue
		}
		levelStats[level] = count
	}
	stats["by_level"] = levelStats

	// 按状态统计
	statusQuery := `SELECT status, COUNT(*) as count 
		FROM vt_alert_records 
		WHERE triggered_at BETWEEN ? AND ? 
		GROUP BY status`

	rows, err = m.conn.Query(statusQuery, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statusStats := make(map[string]int64)
	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			continue
		}
		statusStats[status] = count
	}
	stats["by_status"] = statusStats

	return stats, nil
}

// FindRecentAlerts 查找最近的告警
func (m *vtAlertRecordsModel) FindRecentAlerts(limit int) ([]*VtAlertRecords, error) {
	query := `SELECT id, rule_id, alert_id, alert_name, alert_level, resource_type,
		resource_name, trigger_value, status, triggered_at, occurrence_count
		FROM vt_alert_records 
		ORDER BY triggered_at DESC 
		LIMIT ?`

	rows, err := m.conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*VtAlertRecords
	for rows.Next() {
		var record VtAlertRecords
		err := rows.Scan(
			&record.Id, &record.RuleId, &record.AlertId, &record.AlertName,
			&record.AlertLevel, &record.ResourceType, &record.ResourceName,
			&record.TriggerValue, &record.Status, &record.TriggeredAt,
			&record.OccurrenceCount,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, nil
}
