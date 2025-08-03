package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// VtMonitorData 监控数据表模型
type VtMonitorData struct {
	Id             int64     `db:"id" json:"id"`
	MetricId       int64     `db:"metric_id" json:"metricId"`
	ResourceType   string    `db:"resource_type" json:"resourceType"`
	ResourceId     *int64    `db:"resource_id" json:"resourceId"`
	ResourceName   string    `db:"resource_name" json:"resourceName"`
	InstanceId     string    `db:"instance_id" json:"instanceId"`
	Labels         string    `db:"labels" json:"labels"`
	Value          float64   `db:"value" json:"value"`
	ValueInt       *int64    `db:"value_int" json:"valueInt"`
	ValueStr       string    `db:"value_str" json:"valueStr"`
	ValueBool      *bool     `db:"value_bool" json:"valueBool"`
	CountValue     *int64    `db:"count_value" json:"countValue"`
	SumValue       *float64  `db:"sum_value" json:"sumValue"`
	MinValue       *float64  `db:"min_value" json:"minValue"`
	MaxValue       *float64  `db:"max_value" json:"maxValue"`
	AvgValue       *float64  `db:"avg_value" json:"avgValue"`
	Timestamp      time.Time `db:"timestamp" json:"timestamp"`
	CollectionTime time.Time `db:"collection_time" json:"collectionTime"`
	QualityScore   float64   `db:"quality_score" json:"qualityScore"`
	IsAnomaly      bool      `db:"is_anomaly" json:"isAnomaly"`
	Metadata       string    `db:"metadata" json:"metadata"`
}

// QueryMetricsRequest 查询指标数据请求
type QueryMetricsRequest struct {
	MetricId     int64
	ResourceType string
	ResourceId   int64
	InstanceId   string
	StartTime    time.Time
	EndTime      time.Time
	IsAnomaly    *bool
	Offset       int
	Limit        int
}

// AggregatedMetrics 聚合指标数据
type AggregatedMetrics struct {
	MetricId  int64     `json:"metric_id"`
	Count     int64     `json:"count"`
	Sum       float64   `json:"sum"`
	Avg       float64   `json:"avg"`
	Min       float64   `json:"min"`
	Max       float64   `json:"max"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// VtMonitorDataModel 监控数据模型操作接口
type VtMonitorDataModel interface {
	Insert(data *VtMonitorData) (sql.Result, error)
	FindOne(id int64) (*VtMonitorData, error)
	Update(data *VtMonitorData) error
	Delete(id int64) error
	FindByMetricId(metricId int64, limit int) ([]*VtMonitorData, error)
	FindByResourceType(resourceType string, limit int) ([]*VtMonitorData, error)
	FindByTimeRange(metricId int64, startTime, endTime time.Time) ([]*VtMonitorData, error)
	QueryMetricsData(req *QueryMetricsRequest) ([]*VtMonitorData, int64, error)
	GetLatestByMetric(metricId int64) (*VtMonitorData, error)
	GetLatestByResource(resourceType string, resourceId int64) ([]*VtMonitorData, error)
	BatchInsert(data []*VtMonitorData) error
	DeleteOldData(beforeTime time.Time) (int64, error)
	GetAggregatedData(metricId int64, startTime, endTime time.Time, aggregationType string) (*AggregatedMetrics, error)
	FindAnomalyData(metricId int64, startTime, endTime time.Time) ([]*VtMonitorData, error)
}

// vtMonitorDataModel 监控数据模型实现
type vtMonitorDataModel struct {
	conn *sql.DB
}

// NewVtMonitorDataModel 创建监控数据模型实例
func NewVtMonitorDataModel(conn *sql.DB) VtMonitorDataModel {
	return &vtMonitorDataModel{conn: conn}
}

// Insert 插入监控数据
func (m *vtMonitorDataModel) Insert(data *VtMonitorData) (sql.Result, error) {
	query := `INSERT INTO vt_monitor_data (
		metric_id, resource_type, resource_id, resource_name, instance_id, labels,
		value, value_int, value_str, value_bool, count_value, sum_value, min_value,
		max_value, avg_value, timestamp, collection_time, quality_score, is_anomaly, metadata
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return m.conn.Exec(query,
		data.MetricId, data.ResourceType, data.ResourceId, data.ResourceName,
		data.InstanceId, data.Labels, data.Value, data.ValueInt, data.ValueStr,
		data.ValueBool, data.CountValue, data.SumValue, data.MinValue,
		data.MaxValue, data.AvgValue, data.Timestamp, data.CollectionTime,
		data.QualityScore, data.IsAnomaly, data.Metadata,
	)
}

// FindOne 根据ID查找监控数据
func (m *vtMonitorDataModel) FindOne(id int64) (*VtMonitorData, error) {
	query := `SELECT id, metric_id, resource_type, resource_id, resource_name, instance_id,
		labels, value, value_int, value_str, value_bool, count_value, sum_value, min_value,
		max_value, avg_value, timestamp, collection_time, quality_score, is_anomaly, metadata
		FROM vt_monitor_data WHERE id = ?`

	var data VtMonitorData
	err := m.conn.QueryRow(query, id).Scan(
		&data.Id, &data.MetricId, &data.ResourceType, &data.ResourceId,
		&data.ResourceName, &data.InstanceId, &data.Labels, &data.Value,
		&data.ValueInt, &data.ValueStr, &data.ValueBool, &data.CountValue,
		&data.SumValue, &data.MinValue, &data.MaxValue, &data.AvgValue,
		&data.Timestamp, &data.CollectionTime, &data.QualityScore,
		&data.IsAnomaly, &data.Metadata,
	)

	if err != nil {
		return nil, err
	}
	return &data, nil
}

// Update 更新监控数据
func (m *vtMonitorDataModel) Update(data *VtMonitorData) error {
	query := `UPDATE vt_monitor_data SET 
		metric_id = ?, resource_type = ?, resource_id = ?, resource_name = ?,
		instance_id = ?, labels = ?, value = ?, value_int = ?, value_str = ?,
		value_bool = ?, count_value = ?, sum_value = ?, min_value = ?,
		max_value = ?, avg_value = ?, timestamp = ?, collection_time = ?,
		quality_score = ?, is_anomaly = ?, metadata = ?
		WHERE id = ?`

	_, err := m.conn.Exec(query,
		data.MetricId, data.ResourceType, data.ResourceId, data.ResourceName,
		data.InstanceId, data.Labels, data.Value, data.ValueInt, data.ValueStr,
		data.ValueBool, data.CountValue, data.SumValue, data.MinValue,
		data.MaxValue, data.AvgValue, data.Timestamp, data.CollectionTime,
		data.QualityScore, data.IsAnomaly, data.Metadata, data.Id,
	)
	return err
}

// Delete 删除监控数据
func (m *vtMonitorDataModel) Delete(id int64) error {
	query := "DELETE FROM vt_monitor_data WHERE id = ?"
	_, err := m.conn.Exec(query, id)
	return err
}

// FindByMetricId 根据指标ID查找数据
func (m *vtMonitorDataModel) FindByMetricId(metricId int64, limit int) ([]*VtMonitorData, error) {
	query := `SELECT id, metric_id, resource_type, resource_id, resource_name, instance_id,
		labels, value, value_int, value_str, value_bool, count_value, sum_value, min_value,
		max_value, avg_value, timestamp, collection_time, quality_score, is_anomaly, metadata
		FROM vt_monitor_data WHERE metric_id = ? ORDER BY timestamp DESC LIMIT ?`

	rows, err := m.conn.Query(query, metricId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []*VtMonitorData
	for rows.Next() {
		var data VtMonitorData
		err := rows.Scan(
			&data.Id, &data.MetricId, &data.ResourceType, &data.ResourceId,
			&data.ResourceName, &data.InstanceId, &data.Labels, &data.Value,
			&data.ValueInt, &data.ValueStr, &data.ValueBool, &data.CountValue,
			&data.SumValue, &data.MinValue, &data.MaxValue, &data.AvgValue,
			&data.Timestamp, &data.CollectionTime, &data.QualityScore,
			&data.IsAnomaly, &data.Metadata,
		)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, &data)
	}

	return dataList, nil
}

// FindByResourceType 根据资源类型查找数据
func (m *vtMonitorDataModel) FindByResourceType(resourceType string, limit int) ([]*VtMonitorData, error) {
	query := `SELECT id, metric_id, resource_type, resource_id, resource_name, instance_id,
		labels, value, value_int, value_str, value_bool, count_value, sum_value, min_value,
		max_value, avg_value, timestamp, collection_time, quality_score, is_anomaly, metadata
		FROM vt_monitor_data WHERE resource_type = ? ORDER BY timestamp DESC LIMIT ?`

	rows, err := m.conn.Query(query, resourceType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []*VtMonitorData
	for rows.Next() {
		var data VtMonitorData
		err := rows.Scan(
			&data.Id, &data.MetricId, &data.ResourceType, &data.ResourceId,
			&data.ResourceName, &data.InstanceId, &data.Labels, &data.Value,
			&data.ValueInt, &data.ValueStr, &data.ValueBool, &data.CountValue,
			&data.SumValue, &data.MinValue, &data.MaxValue, &data.AvgValue,
			&data.Timestamp, &data.CollectionTime, &data.QualityScore,
			&data.IsAnomaly, &data.Metadata,
		)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, &data)
	}

	return dataList, nil
}

// FindByTimeRange 根据时间范围查找数据
func (m *vtMonitorDataModel) FindByTimeRange(metricId int64, startTime, endTime time.Time) ([]*VtMonitorData, error) {
	query := `SELECT id, metric_id, resource_type, resource_id, resource_name, instance_id,
		labels, value, value_int, value_str, value_bool, count_value, sum_value, min_value,
		max_value, avg_value, timestamp, collection_time, quality_score, is_anomaly, metadata
		FROM vt_monitor_data WHERE metric_id = ? AND timestamp BETWEEN ? AND ? ORDER BY timestamp`

	rows, err := m.conn.Query(query, metricId, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []*VtMonitorData
	for rows.Next() {
		var data VtMonitorData
		err := rows.Scan(
			&data.Id, &data.MetricId, &data.ResourceType, &data.ResourceId,
			&data.ResourceName, &data.InstanceId, &data.Labels, &data.Value,
			&data.ValueInt, &data.ValueStr, &data.ValueBool, &data.CountValue,
			&data.SumValue, &data.MinValue, &data.MaxValue, &data.AvgValue,
			&data.Timestamp, &data.CollectionTime, &data.QualityScore,
			&data.IsAnomaly, &data.Metadata,
		)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, &data)
	}

	return dataList, nil
}

// QueryMetricsData 查询指标数据
func (m *vtMonitorDataModel) QueryMetricsData(req *QueryMetricsRequest) ([]*VtMonitorData, int64, error) {
	conditions := []string{}
	args := []interface{}{}

	if req.MetricId > 0 {
		conditions = append(conditions, "metric_id = ?")
		args = append(args, req.MetricId)
	}

	if req.ResourceType != "" {
		conditions = append(conditions, "resource_type = ?")
		args = append(args, req.ResourceType)
	}

	if req.ResourceId > 0 {
		conditions = append(conditions, "resource_id = ?")
		args = append(args, req.ResourceId)
	}

	if req.InstanceId != "" {
		conditions = append(conditions, "instance_id = ?")
		args = append(args, req.InstanceId)
	}

	if !req.StartTime.IsZero() {
		conditions = append(conditions, "timestamp >= ?")
		args = append(args, req.StartTime)
	}

	if !req.EndTime.IsZero() {
		conditions = append(conditions, "timestamp <= ?")
		args = append(args, req.EndTime)
	}

	if req.IsAnomaly != nil {
		conditions = append(conditions, "is_anomaly = ?")
		args = append(args, *req.IsAnomaly)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 查询总数
	var total int64
	countQuery := "SELECT COUNT(*) FROM vt_monitor_data " + whereClause
	err := m.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	query := `SELECT id, metric_id, resource_type, resource_id, resource_name, instance_id,
		labels, value, value_int, value_str, value_bool, count_value, sum_value, min_value,
		max_value, avg_value, timestamp, collection_time, quality_score, is_anomaly, metadata
		FROM vt_monitor_data ` + whereClause + ` ORDER BY timestamp DESC LIMIT ? OFFSET ?`

	args = append(args, req.Limit, req.Offset)
	rows, err := m.conn.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var dataList []*VtMonitorData
	for rows.Next() {
		var data VtMonitorData
		err := rows.Scan(
			&data.Id, &data.MetricId, &data.ResourceType, &data.ResourceId,
			&data.ResourceName, &data.InstanceId, &data.Labels, &data.Value,
			&data.ValueInt, &data.ValueStr, &data.ValueBool, &data.CountValue,
			&data.SumValue, &data.MinValue, &data.MaxValue, &data.AvgValue,
			&data.Timestamp, &data.CollectionTime, &data.QualityScore,
			&data.IsAnomaly, &data.Metadata,
		)
		if err != nil {
			return nil, 0, err
		}
		dataList = append(dataList, &data)
	}

	return dataList, total, nil
}

// GetLatestByMetric 获取指标的最新数据
func (m *vtMonitorDataModel) GetLatestByMetric(metricId int64) (*VtMonitorData, error) {
	query := `SELECT id, metric_id, resource_type, resource_id, resource_name, instance_id,
		labels, value, value_int, value_str, value_bool, count_value, sum_value, min_value,
		max_value, avg_value, timestamp, collection_time, quality_score, is_anomaly, metadata
		FROM vt_monitor_data WHERE metric_id = ? ORDER BY timestamp DESC LIMIT 1`

	var data VtMonitorData
	err := m.conn.QueryRow(query, metricId).Scan(
		&data.Id, &data.MetricId, &data.ResourceType, &data.ResourceId,
		&data.ResourceName, &data.InstanceId, &data.Labels, &data.Value,
		&data.ValueInt, &data.ValueStr, &data.ValueBool, &data.CountValue,
		&data.SumValue, &data.MinValue, &data.MaxValue, &data.AvgValue,
		&data.Timestamp, &data.CollectionTime, &data.QualityScore,
		&data.IsAnomaly, &data.Metadata,
	)

	if err != nil {
		return nil, err
	}
	return &data, nil
}

// GetLatestByResource 获取资源的最新指标数据
func (m *vtMonitorDataModel) GetLatestByResource(resourceType string, resourceId int64) ([]*VtMonitorData, error) {
	query := `SELECT id, metric_id, resource_type, resource_id, resource_name, instance_id,
		labels, value, value_int, value_str, value_bool, count_value, sum_value, min_value,
		max_value, avg_value, timestamp, collection_time, quality_score, is_anomaly, metadata
		FROM vt_monitor_data 
		WHERE resource_type = ? AND resource_id = ? 
		  AND timestamp = (
		    SELECT MAX(timestamp) FROM vt_monitor_data 
		    WHERE resource_type = ? AND resource_id = ?
		  )
		ORDER BY metric_id`

	rows, err := m.conn.Query(query, resourceType, resourceId, resourceType, resourceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []*VtMonitorData
	for rows.Next() {
		var data VtMonitorData
		err := rows.Scan(
			&data.Id, &data.MetricId, &data.ResourceType, &data.ResourceId,
			&data.ResourceName, &data.InstanceId, &data.Labels, &data.Value,
			&data.ValueInt, &data.ValueStr, &data.ValueBool, &data.CountValue,
			&data.SumValue, &data.MinValue, &data.MaxValue, &data.AvgValue,
			&data.Timestamp, &data.CollectionTime, &data.QualityScore,
			&data.IsAnomaly, &data.Metadata,
		)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, &data)
	}

	return dataList, nil
}

// BatchInsert 批量插入监控数据
func (m *vtMonitorDataModel) BatchInsert(data []*VtMonitorData) error {
	if len(data) == 0 {
		return nil
	}

	// 构建批量插入SQL
	valueStrings := make([]string, 0, len(data))
	valueArgs := make([]interface{}, 0, len(data)*20)

	for _, item := range data {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs,
			item.MetricId, item.ResourceType, item.ResourceId, item.ResourceName,
			item.InstanceId, item.Labels, item.Value, item.ValueInt, item.ValueStr,
			item.ValueBool, item.CountValue, item.SumValue, item.MinValue,
			item.MaxValue, item.AvgValue, item.Timestamp, item.CollectionTime,
			item.QualityScore, item.IsAnomaly, item.Metadata,
		)
	}

	query := `INSERT INTO vt_monitor_data (
		metric_id, resource_type, resource_id, resource_name, instance_id, labels,
		value, value_int, value_str, value_bool, count_value, sum_value, min_value,
		max_value, avg_value, timestamp, collection_time, quality_score, is_anomaly, metadata
	) VALUES ` + strings.Join(valueStrings, ",")

	_, err := m.conn.Exec(query, valueArgs...)
	return err
}

// DeleteOldData 删除旧数据
func (m *vtMonitorDataModel) DeleteOldData(beforeTime time.Time) (int64, error) {
	query := "DELETE FROM vt_monitor_data WHERE timestamp < ?"
	result, err := m.conn.Exec(query, beforeTime)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// GetAggregatedData 获取聚合数据
func (m *vtMonitorDataModel) GetAggregatedData(metricId int64, startTime, endTime time.Time, aggregationType string) (*AggregatedMetrics, error) {
	var query string
	switch aggregationType {
	case "sum":
		query = "SELECT COUNT(*), SUM(value), SUM(value), MIN(value), MAX(value) FROM vt_monitor_data WHERE metric_id = ? AND timestamp BETWEEN ? AND ?"
	case "avg":
		query = "SELECT COUNT(*), SUM(value), AVG(value), MIN(value), MAX(value) FROM vt_monitor_data WHERE metric_id = ? AND timestamp BETWEEN ? AND ?"
	case "min":
		query = "SELECT COUNT(*), SUM(value), AVG(value), MIN(value), MAX(value) FROM vt_monitor_data WHERE metric_id = ? AND timestamp BETWEEN ? AND ?"
	case "max":
		query = "SELECT COUNT(*), SUM(value), AVG(value), MIN(value), MAX(value) FROM vt_monitor_data WHERE metric_id = ? AND timestamp BETWEEN ? AND ?"
	default:
		query = "SELECT COUNT(*), SUM(value), AVG(value), MIN(value), MAX(value) FROM vt_monitor_data WHERE metric_id = ? AND timestamp BETWEEN ? AND ?"
	}

	var count int64
	var sum, avg, min, max float64
	err := m.conn.QueryRow(query, metricId, startTime, endTime).Scan(&count, &sum, &avg, &min, &max)
	if err != nil {
		return nil, err
	}

	result := &AggregatedMetrics{
		MetricId:  metricId,
		Count:     count,
		Sum:       sum,
		Avg:       avg,
		Min:       min,
		Max:       max,
		StartTime: startTime,
		EndTime:   endTime,
	}

	return result, nil
}

// FindAnomalyData 查找异常数据
func (m *vtMonitorDataModel) FindAnomalyData(metricId int64, startTime, endTime time.Time) ([]*VtMonitorData, error) {
	query := `SELECT id, metric_id, resource_type, resource_id, resource_name, instance_id,
		labels, value, value_int, value_str, value_bool, count_value, sum_value, min_value,
		max_value, avg_value, timestamp, collection_time, quality_score, is_anomaly, metadata
		FROM vt_monitor_data WHERE metric_id = ? AND is_anomaly = 1 AND timestamp BETWEEN ? AND ? 
		ORDER BY timestamp DESC`

	rows, err := m.conn.Query(query, metricId, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []*VtMonitorData
	for rows.Next() {
		var data VtMonitorData
		err := rows.Scan(
			&data.Id, &data.MetricId, &data.ResourceType, &data.ResourceId,
			&data.ResourceName, &data.InstanceId, &data.Labels, &data.Value,
			&data.ValueInt, &data.ValueStr, &data.ValueBool, &data.CountValue,
			&data.SumValue, &data.MinValue, &data.MaxValue, &data.AvgValue,
			&data.Timestamp, &data.CollectionTime, &data.QualityScore,
			&data.IsAnomaly, &data.Metadata,
		)
		if err != nil {
			return nil, err
		}
		dataList = append(dataList, &data)
	}

	return dataList, nil
}

// GetLabelsMap 获取标签映射
func (m *VtMonitorData) GetLabelsMap() (map[string]interface{}, error) {
	if m.Labels == "" {
		return make(map[string]interface{}), nil
	}

	var labels map[string]interface{}
	err := json.Unmarshal([]byte(m.Labels), &labels)
	if err != nil {
		return nil, fmt.Errorf("解析标签失败: %v", err)
	}

	return labels, nil
}

// SetLabels 设置标签
func (m *VtMonitorData) SetLabels(labels map[string]interface{}) error {
	if labels == nil || len(labels) == 0 {
		m.Labels = ""
		return nil
	}

	data, err := json.Marshal(labels)
	if err != nil {
		return fmt.Errorf("序列化标签失败: %v", err)
	}

	m.Labels = string(data)
	return nil
}

// GetMetadataMap 获取元数据映射
func (m *VtMonitorData) GetMetadataMap() (map[string]interface{}, error) {
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

// SetMetadata 设置元数据
func (m *VtMonitorData) SetMetadata(metadata map[string]interface{}) error {
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
