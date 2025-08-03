package model

import (
	"database/sql"
	"time"
)

// VtGpuClusters GPU集群表模型
type VtGpuClusters struct {
	Id             int64     `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	DisplayName    string    `db:"display_name" json:"displayName"`
	Description    string    `db:"description" json:"description"`
	ClusterType    string    `db:"cluster_type" json:"clusterType"`
	Status         string    `db:"status" json:"status"`
	KubeConfig     string    `db:"kube_config" json:"kubeConfig"`
	ApiEndpoint    string    `db:"api_endpoint" json:"apiEndpoint"`
	Region         string    `db:"region" json:"region"`
	Zone           string    `db:"zone" json:"zone"`
	TotalNodes     int       `db:"total_nodes" json:"totalNodes"`
	ActiveNodes    int       `db:"active_nodes" json:"activeNodes"`
	TotalGpus      int       `db:"total_gpus" json:"totalGpus"`
	AvailableGpus  int       `db:"available_gpus" json:"availableGpus"`
	AllocatedGpus  int       `db:"allocated_gpus" json:"allocatedGpus"`
	ResourceLabels string    `db:"resource_labels" json:"resourceLabels"`
	MetricsConfig  string    `db:"metrics_config" json:"metricsConfig"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt"`
}

// VtGpuClustersModel GPU集群模型操作接口
type VtGpuClustersModel interface {
	Insert(data *VtGpuClusters) (sql.Result, error)
	FindOne(id int64) (*VtGpuClusters, error)
	Update(data *VtGpuClusters) error
	Delete(id int64) error
	FindAll(page, pageSize int, status, clusterType, region, search string) ([]*VtGpuClusters, int64, error)
}

// vtGpuClustersModelImpl GPU集群模型实现
type vtGpuClustersModelImpl struct {
	conn *sql.DB
}

// NewVtGpuClustersModel 创建GPU集群模型实例
func NewVtGpuClustersModel(conn *sql.DB) VtGpuClustersModel {
	return &vtGpuClustersModelImpl{
		conn: conn,
	}
}

func (m *vtGpuClustersModelImpl) Insert(data *VtGpuClusters) (sql.Result, error) {
	query := `INSERT INTO vt_gpu_clusters (
		name, display_name, description, cluster_type, status, 
		kube_config, api_endpoint, region, zone, 
		total_nodes, active_nodes, total_gpus, available_gpus, allocated_gpus,
		resource_labels, metrics_config
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return m.conn.Exec(query,
		data.Name, data.DisplayName, data.Description, data.ClusterType, data.Status,
		data.KubeConfig, data.ApiEndpoint, data.Region, data.Zone,
		data.TotalNodes, data.ActiveNodes, data.TotalGpus, data.AvailableGpus, data.AllocatedGpus,
		data.ResourceLabels, data.MetricsConfig,
	)
}

func (m *vtGpuClustersModelImpl) FindOne(id int64) (*VtGpuClusters, error) {
	var cluster VtGpuClusters
	query := `SELECT id, name, display_name, description, cluster_type, status,
		kube_config, api_endpoint, region, zone,
		total_nodes, active_nodes, total_gpus, available_gpus, allocated_gpus,
		resource_labels, metrics_config, created_at, updated_at
		FROM vt_gpu_clusters WHERE id = ?`

	err := m.conn.QueryRow(query, id).Scan(
		&cluster.Id, &cluster.Name, &cluster.DisplayName, &cluster.Description, &cluster.ClusterType, &cluster.Status,
		&cluster.KubeConfig, &cluster.ApiEndpoint, &cluster.Region, &cluster.Zone,
		&cluster.TotalNodes, &cluster.ActiveNodes, &cluster.TotalGpus, &cluster.AvailableGpus, &cluster.AllocatedGpus,
		&cluster.ResourceLabels, &cluster.MetricsConfig, &cluster.CreatedAt, &cluster.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &cluster, nil
}

func (m *vtGpuClustersModelImpl) Update(data *VtGpuClusters) error {
	query := `UPDATE vt_gpu_clusters SET 
		name = ?, display_name = ?, description = ?, cluster_type = ?, status = ?,
		kube_config = ?, api_endpoint = ?, region = ?, zone = ?,
		total_nodes = ?, active_nodes = ?, total_gpus = ?, available_gpus = ?, allocated_gpus = ?,
		resource_labels = ?, metrics_config = ?, updated_at = NOW()
		WHERE id = ?`

	_, err := m.conn.Exec(query,
		data.Name, data.DisplayName, data.Description, data.ClusterType, data.Status,
		data.KubeConfig, data.ApiEndpoint, data.Region, data.Zone,
		data.TotalNodes, data.ActiveNodes, data.TotalGpus, data.AvailableGpus, data.AllocatedGpus,
		data.ResourceLabels, data.MetricsConfig, data.Id,
	)

	return err
}

func (m *vtGpuClustersModelImpl) Delete(id int64) error {
	query := `DELETE FROM vt_gpu_clusters WHERE id = ?`
	_, err := m.conn.Exec(query, id)
	return err
}

func (m *vtGpuClustersModelImpl) FindAll(page, pageSize int, status, clusterType, region, search string) ([]*VtGpuClusters, int64, error) {
	offset := (page - 1) * pageSize

	// 构建WHERE条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}
	if clusterType != "" {
		whereClause += " AND cluster_type = ?"
		args = append(args, clusterType)
	}
	if region != "" {
		whereClause += " AND region = ?"
		args = append(args, region)
	}
	if search != "" {
		whereClause += " AND (name LIKE ? OR display_name LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	// 查询总数
	countQuery := "SELECT COUNT(*) FROM vt_gpu_clusters " + whereClause
	var total int64
	err := m.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	query := `SELECT id, name, display_name, description, cluster_type, status,
		kube_config, api_endpoint, region, zone,
		total_nodes, active_nodes, total_gpus, available_gpus, allocated_gpus,
		resource_labels, metrics_config, created_at, updated_at
		FROM vt_gpu_clusters ` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`

	args = append(args, pageSize, offset)
	rows, err := m.conn.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var clusters []*VtGpuClusters
	for rows.Next() {
		var cluster VtGpuClusters
		err := rows.Scan(
			&cluster.Id, &cluster.Name, &cluster.DisplayName, &cluster.Description, &cluster.ClusterType, &cluster.Status,
			&cluster.KubeConfig, &cluster.ApiEndpoint, &cluster.Region, &cluster.Zone,
			&cluster.TotalNodes, &cluster.ActiveNodes, &cluster.TotalGpus, &cluster.AvailableGpus, &cluster.AllocatedGpus,
			&cluster.ResourceLabels, &cluster.MetricsConfig, &cluster.CreatedAt, &cluster.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		clusters = append(clusters, &cluster)
	}

	return clusters, total, nil
}
