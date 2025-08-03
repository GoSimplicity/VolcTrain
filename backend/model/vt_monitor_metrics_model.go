package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// VtMonitorMetrics 监控指标表模型
type VtMonitorMetrics struct {
	Id                        int64     `db:"id" json:"id"`
	Name                      string    `db:"name" json:"name"`
	DisplayName               string    `db:"display_name" json:"displayName"`
	Description               string    `db:"description" json:"description"`
	MetricType                string    `db:"metric_type" json:"metricType"`
	DataType                  string    `db:"data_type" json:"dataType"`
	Category                  string    `db:"category" json:"category"`
	Module                    string    `db:"module" json:"module"`
	SourceType                string    `db:"source_type" json:"sourceType"`
	Unit                      string    `db:"unit" json:"unit"`
	AggregationType           string    `db:"aggregation_type" json:"aggregationType"`
	CollectionIntervalSeconds int       `db:"collection_interval_seconds" json:"collectionIntervalSeconds"`
	RetentionDays             int       `db:"retention_days" json:"retentionDays"`
	NormalRangeMin            *float64  `db:"normal_range_min" json:"normalRangeMin"`
	NormalRangeMax            *float64  `db:"normal_range_max" json:"normalRangeMax"`
	WarningThreshold          *float64  `db:"warning_threshold" json:"warningThreshold"`
	CriticalThreshold         *float64  `db:"critical_threshold" json:"criticalThreshold"`
	ThresholdCondition        string    `db:"threshold_condition" json:"thresholdCondition"`
	Status                    string    `db:"status" json:"status"`
	IsBuiltin                 bool      `db:"is_builtin" json:"isBuiltin"`
	IsCore                    bool      `db:"is_core" json:"isCore"`
	DefaultLabels             string    `db:"default_labels" json:"defaultLabels"`
	Dimensions                string    `db:"dimensions" json:"dimensions"`
	Metadata                  string    `db:"metadata" json:"metadata"`
	CreatedAt                 time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt                 time.Time `db:"updated_at" json:"updatedAt"`
}

// VtMonitorMetricsModel 监控指标模型操作接口
type VtMonitorMetricsModel interface {
	Insert(data *VtMonitorMetrics) (sql.Result, error)
	FindOne(id int64) (*VtMonitorMetrics, error)
	FindOneByName(name string) (*VtMonitorMetrics, error)
	Update(data *VtMonitorMetrics) error
	Delete(id int64) error
	List(page, pageSize int, filters map[string]interface{}) ([]*VtMonitorMetrics, int64, error)
	FindByCategory(category string) ([]*VtMonitorMetrics, error)
	FindByModule(module string) ([]*VtMonitorMetrics, error)
	FindActiveMetrics() ([]*VtMonitorMetrics, error)
	FindBuiltinMetrics() ([]*VtMonitorMetrics, error)
	FindCoreMetrics() ([]*VtMonitorMetrics, error)
	UpdateStatus(id int64, status string) error
	GetMetricsByType(metricType string) ([]*VtMonitorMetrics, error)
	SearchMetrics(keyword string, offset, limit int) ([]*VtMonitorMetrics, int64, error)
}

// vtMonitorMetricsModel 监控指标模型实现
type vtMonitorMetricsModel struct {
	conn *sql.DB
}

// NewVtMonitorMetricsModel 创建监控指标模型实例
func NewVtMonitorMetricsModel(conn *sql.DB) VtMonitorMetricsModel {
	return &vtMonitorMetricsModel{conn: conn}
}

// Insert 插入监控指标
func (m *vtMonitorMetricsModel) Insert(data *VtMonitorMetrics) (sql.Result, error) {
	query := `INSERT INTO vt_monitor_metrics (
		name, display_name, description, metric_type, data_type, category, module, 
		source_type, unit, aggregation_type, collection_interval_seconds, retention_days,
		normal_range_min, normal_range_max, warning_threshold, critical_threshold,
		threshold_condition, status, is_builtin, is_core, default_labels, dimensions, metadata
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return m.conn.Exec(query,
		data.Name, data.DisplayName, data.Description, data.MetricType, data.DataType,
		data.Category, data.Module, data.SourceType, data.Unit, data.AggregationType,
		data.CollectionIntervalSeconds, data.RetentionDays, data.NormalRangeMin,
		data.NormalRangeMax, data.WarningThreshold, data.CriticalThreshold,
		data.ThresholdCondition, data.Status, data.IsBuiltin, data.IsCore,
		data.DefaultLabels, data.Dimensions, data.Metadata,
	)
}

// FindOne 根据ID查找监控指标
func (m *vtMonitorMetricsModel) FindOne(id int64) (*VtMonitorMetrics, error) {
	query := `SELECT id, name, display_name, description, metric_type, data_type, category, 
		module, source_type, unit, aggregation_type, collection_interval_seconds, 
		retention_days, normal_range_min, normal_range_max, warning_threshold, 
		critical_threshold, threshold_condition, status, is_builtin, is_core, 
		default_labels, dimensions, metadata, created_at, updated_at
		FROM vt_monitor_metrics WHERE id = ?`

	var metric VtMonitorMetrics
	err := m.conn.QueryRow(query, id).Scan(
		&metric.Id, &metric.Name, &metric.DisplayName, &metric.Description,
		&metric.MetricType, &metric.DataType, &metric.Category, &metric.Module,
		&metric.SourceType, &metric.Unit, &metric.AggregationType,
		&metric.CollectionIntervalSeconds, &metric.RetentionDays,
		&metric.NormalRangeMin, &metric.NormalRangeMax, &metric.WarningThreshold,
		&metric.CriticalThreshold, &metric.ThresholdCondition, &metric.Status,
		&metric.IsBuiltin, &metric.IsCore, &metric.DefaultLabels,
		&metric.Dimensions, &metric.Metadata, &metric.CreatedAt, &metric.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &metric, nil
}

// FindOneByName 根据名称查找监控指标
func (m *vtMonitorMetricsModel) FindOneByName(name string) (*VtMonitorMetrics, error) {
	query := `SELECT id, name, display_name, description, metric_type, data_type, category, 
		module, source_type, unit, aggregation_type, collection_interval_seconds, 
		retention_days, normal_range_min, normal_range_max, warning_threshold, 
		critical_threshold, threshold_condition, status, is_builtin, is_core, 
		default_labels, dimensions, metadata, created_at, updated_at
		FROM vt_monitor_metrics WHERE name = ?`

	var metric VtMonitorMetrics
	err := m.conn.QueryRow(query, name).Scan(
		&metric.Id, &metric.Name, &metric.DisplayName, &metric.Description,
		&metric.MetricType, &metric.DataType, &metric.Category, &metric.Module,
		&metric.SourceType, &metric.Unit, &metric.AggregationType,
		&metric.CollectionIntervalSeconds, &metric.RetentionDays,
		&metric.NormalRangeMin, &metric.NormalRangeMax, &metric.WarningThreshold,
		&metric.CriticalThreshold, &metric.ThresholdCondition, &metric.Status,
		&metric.IsBuiltin, &metric.IsCore, &metric.DefaultLabels,
		&metric.Dimensions, &metric.Metadata, &metric.CreatedAt, &metric.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &metric, nil
}

// Update 更新监控指标
func (m *vtMonitorMetricsModel) Update(data *VtMonitorMetrics) error {
	query := `UPDATE vt_monitor_metrics SET 
		name = ?, display_name = ?, description = ?, metric_type = ?, data_type = ?,
		category = ?, module = ?, source_type = ?, unit = ?, aggregation_type = ?,
		collection_interval_seconds = ?, retention_days = ?, normal_range_min = ?,
		normal_range_max = ?, warning_threshold = ?, critical_threshold = ?,
		threshold_condition = ?, status = ?, is_builtin = ?, is_core = ?,
		default_labels = ?, dimensions = ?, metadata = ?, updated_at = NOW()
		WHERE id = ?`

	_, err := m.conn.Exec(query,
		data.Name, data.DisplayName, data.Description, data.MetricType, data.DataType,
		data.Category, data.Module, data.SourceType, data.Unit, data.AggregationType,
		data.CollectionIntervalSeconds, data.RetentionDays, data.NormalRangeMin,
		data.NormalRangeMax, data.WarningThreshold, data.CriticalThreshold,
		data.ThresholdCondition, data.Status, data.IsBuiltin, data.IsCore,
		data.DefaultLabels, data.Dimensions, data.Metadata, data.Id,
	)
	return err
}

// Delete 删除监控指标
func (m *vtMonitorMetricsModel) Delete(id int64) error {
	query := "DELETE FROM vt_monitor_metrics WHERE id = ?"
	_, err := m.conn.Exec(query, id)
	return err
}

// List 分页查询监控指标
func (m *vtMonitorMetricsModel) List(page, pageSize int, filters map[string]interface{}) ([]*VtMonitorMetrics, int64, error) {
	offset := (page - 1) * pageSize

	// 构建WHERE条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if name, ok := filters["name"]; ok && name != "" {
		whereClause += " AND name LIKE ?"
		args = append(args, "%"+name.(string)+"%")
	}

	if metricType, ok := filters["metric_type"]; ok && metricType != "" {
		whereClause += " AND metric_type = ?"
		args = append(args, metricType)
	}

	if category, ok := filters["category"]; ok && category != "" {
		whereClause += " AND category = ?"
		args = append(args, category)
	}

	if status, ok := filters["status"]; ok && status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}

	// 查询总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM vt_monitor_metrics " + whereClause
	err := m.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	query := `SELECT id, name, display_name, description, metric_type, data_type, category, 
		module, source_type, unit, aggregation_type, collection_interval_seconds, 
		retention_days, normal_range_min, normal_range_max, warning_threshold, 
		critical_threshold, threshold_condition, status, is_builtin, is_core, 
		default_labels, dimensions, metadata, created_at, updated_at
		FROM vt_monitor_metrics ` + whereClause + ` 
		ORDER BY is_core DESC, name LIMIT ? OFFSET ?`

	args = append(args, pageSize, offset)
	rows, err := m.conn.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var metrics []*VtMonitorMetrics
	for rows.Next() {
		var metric VtMonitorMetrics
		err := rows.Scan(
			&metric.Id, &metric.Name, &metric.DisplayName, &metric.Description,
			&metric.MetricType, &metric.DataType, &metric.Category, &metric.Module,
			&metric.SourceType, &metric.Unit, &metric.AggregationType,
			&metric.CollectionIntervalSeconds, &metric.RetentionDays,
			&metric.NormalRangeMin, &metric.NormalRangeMax, &metric.WarningThreshold,
			&metric.CriticalThreshold, &metric.ThresholdCondition, &metric.Status,
			&metric.IsBuiltin, &metric.IsCore, &metric.DefaultLabels,
			&metric.Dimensions, &metric.Metadata, &metric.CreatedAt, &metric.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		metrics = append(metrics, &metric)
	}

	return metrics, total, nil
}

// FindByCategory 根据分类查找指标
func (m *vtMonitorMetricsModel) FindByCategory(category string) ([]*VtMonitorMetrics, error) {
	query := `SELECT id, name, display_name, description, metric_type, data_type, category, 
		module, source_type, unit, aggregation_type, collection_interval_seconds, 
		retention_days, normal_range_min, normal_range_max, warning_threshold, 
		critical_threshold, threshold_condition, status, is_builtin, is_core, 
		default_labels, dimensions, metadata, created_at, updated_at
		FROM vt_monitor_metrics WHERE category = ? AND status = 'active' 
		ORDER BY is_core DESC, name`

	rows, err := m.conn.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*VtMonitorMetrics
	for rows.Next() {
		var metric VtMonitorMetrics
		err := rows.Scan(
			&metric.Id, &metric.Name, &metric.DisplayName, &metric.Description,
			&metric.MetricType, &metric.DataType, &metric.Category, &metric.Module,
			&metric.SourceType, &metric.Unit, &metric.AggregationType,
			&metric.CollectionIntervalSeconds, &metric.RetentionDays,
			&metric.NormalRangeMin, &metric.NormalRangeMax, &metric.WarningThreshold,
			&metric.CriticalThreshold, &metric.ThresholdCondition, &metric.Status,
			&metric.IsBuiltin, &metric.IsCore, &metric.DefaultLabels,
			&metric.Dimensions, &metric.Metadata, &metric.CreatedAt, &metric.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, &metric)
	}

	return metrics, nil
}

// FindByModule 根据模块查找指标
func (m *vtMonitorMetricsModel) FindByModule(module string) ([]*VtMonitorMetrics, error) {
	query := `SELECT id, name, display_name, description, metric_type, data_type, category, 
		module, source_type, unit, aggregation_type, collection_interval_seconds, 
		retention_days, normal_range_min, normal_range_max, warning_threshold, 
		critical_threshold, threshold_condition, status, is_builtin, is_core, 
		default_labels, dimensions, metadata, created_at, updated_at
		FROM vt_monitor_metrics WHERE module = ? AND status = 'active' 
		ORDER BY is_core DESC, name`

	rows, err := m.conn.Query(query, module)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*VtMonitorMetrics
	for rows.Next() {
		var metric VtMonitorMetrics
		err := rows.Scan(
			&metric.Id, &metric.Name, &metric.DisplayName, &metric.Description,
			&metric.MetricType, &metric.DataType, &metric.Category, &metric.Module,
			&metric.SourceType, &metric.Unit, &metric.AggregationType,
			&metric.CollectionIntervalSeconds, &metric.RetentionDays,
			&metric.NormalRangeMin, &metric.NormalRangeMax, &metric.WarningThreshold,
			&metric.CriticalThreshold, &metric.ThresholdCondition, &metric.Status,
			&metric.IsBuiltin, &metric.IsCore, &metric.DefaultLabels,
			&metric.Dimensions, &metric.Metadata, &metric.CreatedAt, &metric.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, &metric)
	}

	return metrics, nil
}

// FindActiveMetrics 查找所有活跃的指标
func (m *vtMonitorMetricsModel) FindActiveMetrics() ([]*VtMonitorMetrics, error) {
	query := `SELECT id, name, display_name, description, metric_type, data_type, category, 
		module, source_type, unit, aggregation_type, collection_interval_seconds, 
		retention_days, normal_range_min, normal_range_max, warning_threshold, 
		critical_threshold, threshold_condition, status, is_builtin, is_core, 
		default_labels, dimensions, metadata, created_at, updated_at
		FROM vt_monitor_metrics WHERE status = 'active' 
		ORDER BY is_core DESC, category, name`

	rows, err := m.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*VtMonitorMetrics
	for rows.Next() {
		var metric VtMonitorMetrics
		err := rows.Scan(
			&metric.Id, &metric.Name, &metric.DisplayName, &metric.Description,
			&metric.MetricType, &metric.DataType, &metric.Category, &metric.Module,
			&metric.SourceType, &metric.Unit, &metric.AggregationType,
			&metric.CollectionIntervalSeconds, &metric.RetentionDays,
			&metric.NormalRangeMin, &metric.NormalRangeMax, &metric.WarningThreshold,
			&metric.CriticalThreshold, &metric.ThresholdCondition, &metric.Status,
			&metric.IsBuiltin, &metric.IsCore, &metric.DefaultLabels,
			&metric.Dimensions, &metric.Metadata, &metric.CreatedAt, &metric.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, &metric)
	}

	return metrics, nil
}

// FindBuiltinMetrics 查找所有内置指标
func (m *vtMonitorMetricsModel) FindBuiltinMetrics() ([]*VtMonitorMetrics, error) {
	query := `SELECT id, name, display_name, description, metric_type, data_type, category, 
		module, source_type, unit, aggregation_type, collection_interval_seconds, 
		retention_days, normal_range_min, normal_range_max, warning_threshold, 
		critical_threshold, threshold_condition, status, is_builtin, is_core, 
		default_labels, dimensions, metadata, created_at, updated_at
		FROM vt_monitor_metrics WHERE is_builtin = 1 AND status = 'active' 
		ORDER BY category, name`

	rows, err := m.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*VtMonitorMetrics
	for rows.Next() {
		var metric VtMonitorMetrics
		err := rows.Scan(
			&metric.Id, &metric.Name, &metric.DisplayName, &metric.Description,
			&metric.MetricType, &metric.DataType, &metric.Category, &metric.Module,
			&metric.SourceType, &metric.Unit, &metric.AggregationType,
			&metric.CollectionIntervalSeconds, &metric.RetentionDays,
			&metric.NormalRangeMin, &metric.NormalRangeMax, &metric.WarningThreshold,
			&metric.CriticalThreshold, &metric.ThresholdCondition, &metric.Status,
			&metric.IsBuiltin, &metric.IsCore, &metric.DefaultLabels,
			&metric.Dimensions, &metric.Metadata, &metric.CreatedAt, &metric.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, &metric)
	}

	return metrics, nil
}

// FindCoreMetrics 查找所有核心指标
func (m *vtMonitorMetricsModel) FindCoreMetrics() ([]*VtMonitorMetrics, error) {
	query := `SELECT id, name, display_name, description, metric_type, data_type, category, 
		module, source_type, unit, aggregation_type, collection_interval_seconds, 
		retention_days, normal_range_min, normal_range_max, warning_threshold, 
		critical_threshold, threshold_condition, status, is_builtin, is_core, 
		default_labels, dimensions, metadata, created_at, updated_at
		FROM vt_monitor_metrics WHERE is_core = 1 AND status = 'active' 
		ORDER BY category, name`

	rows, err := m.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*VtMonitorMetrics
	for rows.Next() {
		var metric VtMonitorMetrics
		err := rows.Scan(
			&metric.Id, &metric.Name, &metric.DisplayName, &metric.Description,
			&metric.MetricType, &metric.DataType, &metric.Category, &metric.Module,
			&metric.SourceType, &metric.Unit, &metric.AggregationType,
			&metric.CollectionIntervalSeconds, &metric.RetentionDays,
			&metric.NormalRangeMin, &metric.NormalRangeMax, &metric.WarningThreshold,
			&metric.CriticalThreshold, &metric.ThresholdCondition, &metric.Status,
			&metric.IsBuiltin, &metric.IsCore, &metric.DefaultLabels,
			&metric.Dimensions, &metric.Metadata, &metric.CreatedAt, &metric.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, &metric)
	}

	return metrics, nil
}

// UpdateStatus 更新指标状态
func (m *vtMonitorMetricsModel) UpdateStatus(id int64, status string) error {
	query := "UPDATE vt_monitor_metrics SET status = ?, updated_at = NOW() WHERE id = ?"
	_, err := m.conn.Exec(query, status, id)
	return err
}

// GetMetricsByType 根据指标类型获取指标
func (m *vtMonitorMetricsModel) GetMetricsByType(metricType string) ([]*VtMonitorMetrics, error) {
	query := `SELECT id, name, display_name, description, metric_type, data_type, category, 
		module, source_type, unit, aggregation_type, collection_interval_seconds, 
		retention_days, normal_range_min, normal_range_max, warning_threshold, 
		critical_threshold, threshold_condition, status, is_builtin, is_core, 
		default_labels, dimensions, metadata, created_at, updated_at
		FROM vt_monitor_metrics WHERE metric_type = ? AND status = 'active' 
		ORDER BY name`

	rows, err := m.conn.Query(query, metricType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*VtMonitorMetrics
	for rows.Next() {
		var metric VtMonitorMetrics
		err := rows.Scan(
			&metric.Id, &metric.Name, &metric.DisplayName, &metric.Description,
			&metric.MetricType, &metric.DataType, &metric.Category, &metric.Module,
			&metric.SourceType, &metric.Unit, &metric.AggregationType,
			&metric.CollectionIntervalSeconds, &metric.RetentionDays,
			&metric.NormalRangeMin, &metric.NormalRangeMax, &metric.WarningThreshold,
			&metric.CriticalThreshold, &metric.ThresholdCondition, &metric.Status,
			&metric.IsBuiltin, &metric.IsCore, &metric.DefaultLabels,
			&metric.Dimensions, &metric.Metadata, &metric.CreatedAt, &metric.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, &metric)
	}

	return metrics, nil
}

// SearchMetrics 搜索指标
func (m *vtMonitorMetricsModel) SearchMetrics(keyword string, offset, limit int) ([]*VtMonitorMetrics, int64, error) {
	whereClause := "WHERE status = 'active'"
	args := []interface{}{}

	if keyword != "" {
		whereClause += " AND (name LIKE ? OR display_name LIKE ? OR description LIKE ?)"
		likeKeyword := "%" + keyword + "%"
		args = append(args, likeKeyword, likeKeyword, likeKeyword)
	}

	// 查询总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM vt_monitor_metrics " + whereClause
	err := m.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	query := `SELECT id, name, display_name, description, metric_type, data_type, category, 
		module, source_type, unit, aggregation_type, collection_interval_seconds, 
		retention_days, normal_range_min, normal_range_max, warning_threshold, 
		critical_threshold, threshold_condition, status, is_builtin, is_core, 
		default_labels, dimensions, metadata, created_at, updated_at
		FROM vt_monitor_metrics ` + whereClause + ` 
		ORDER BY is_core DESC, name LIMIT ? OFFSET ?`

	args = append(args, limit, offset)
	rows, err := m.conn.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var metrics []*VtMonitorMetrics
	for rows.Next() {
		var metric VtMonitorMetrics
		err := rows.Scan(
			&metric.Id, &metric.Name, &metric.DisplayName, &metric.Description,
			&metric.MetricType, &metric.DataType, &metric.Category, &metric.Module,
			&metric.SourceType, &metric.Unit, &metric.AggregationType,
			&metric.CollectionIntervalSeconds, &metric.RetentionDays,
			&metric.NormalRangeMin, &metric.NormalRangeMax, &metric.WarningThreshold,
			&metric.CriticalThreshold, &metric.ThresholdCondition, &metric.Status,
			&metric.IsBuiltin, &metric.IsCore, &metric.DefaultLabels,
			&metric.Dimensions, &metric.Metadata, &metric.CreatedAt, &metric.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		metrics = append(metrics, &metric)
	}

	return metrics, total, nil
}

// GetDefaultLabelsMap 获取默认标签映射
func (m *VtMonitorMetrics) GetDefaultLabelsMap() (map[string]interface{}, error) {
	if m.DefaultLabels == "" {
		return make(map[string]interface{}), nil
	}

	var labels map[string]interface{}
	err := json.Unmarshal([]byte(m.DefaultLabels), &labels)
	if err != nil {
		return nil, fmt.Errorf("解析默认标签失败: %v", err)
	}

	return labels, nil
}

// GetDimensionsMap 获取维度映射
func (m *VtMonitorMetrics) GetDimensionsMap() (map[string]interface{}, error) {
	if m.Dimensions == "" {
		return make(map[string]interface{}), nil
	}

	var dimensions map[string]interface{}
	err := json.Unmarshal([]byte(m.Dimensions), &dimensions)
	if err != nil {
		return nil, fmt.Errorf("解析维度信息失败: %v", err)
	}

	return dimensions, nil
}

// GetMetadataMap 获取元数据映射
func (m *VtMonitorMetrics) GetMetadataMap() (map[string]interface{}, error) {
	if m.Metadata == "" {
		return make(map[string]interface{}), nil
	}

	var metadata map[string]interface{}
	err := json.Unmarshal([]byte(m.Metadata), &metadata)
	if err != nil {
		return nil, fmt.Errorf("解析元数据失败: %v", err)
	}

	return metadata, nil
}

// SetDefaultLabels 设置默认标签
func (m *VtMonitorMetrics) SetDefaultLabels(labels map[string]interface{}) error {
	if labels == nil || len(labels) == 0 {
		m.DefaultLabels = ""
		return nil
	}

	data, err := json.Marshal(labels)
	if err != nil {
		return fmt.Errorf("序列化默认标签失败: %v", err)
	}

	m.DefaultLabels = string(data)
	return nil
}

// SetDimensions 设置维度信息
func (m *VtMonitorMetrics) SetDimensions(dimensions map[string]interface{}) error {
	if dimensions == nil || len(dimensions) == 0 {
		m.Dimensions = ""
		return nil
	}

	data, err := json.Marshal(dimensions)
	if err != nil {
		return fmt.Errorf("序列化维度信息失败: %v", err)
	}

	m.Dimensions = string(data)
	return nil
}

// SetMetadata 设置元数据
func (m *VtMonitorMetrics) SetMetadata(metadata map[string]interface{}) error {
	if metadata == nil || len(metadata) == 0 {
		m.Metadata = ""
		return nil
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("序列化元数据失败: %v", err)
	}

	m.Metadata = string(data)
	return nil
}
