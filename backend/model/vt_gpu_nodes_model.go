package model

import (
	"database/sql"
	"time"
)

// VtGpuNodes GPU节点表模型
type VtGpuNodes struct {
	Id            int64     `db:"id" json:"id"`
	ClusterId     int64     `db:"cluster_id" json:"clusterId"`
	Name          string    `db:"name" json:"name"`
	Hostname      string    `db:"hostname" json:"hostname"`
	InternalIp    string    `db:"internal_ip" json:"internalIp"`
	ExternalIp    string    `db:"external_ip" json:"externalIp"`
	Status        string    `db:"status" json:"status"`
	NodeType      string    `db:"node_type" json:"nodeType"`
	CpuCores      int       `db:"cpu_cores" json:"cpuCores"`
	MemoryGb      int       `db:"memory_gb" json:"memoryGb"`
	StorageGb     int       `db:"storage_gb" json:"storageGb"`
	GpuCount      int       `db:"gpu_count" json:"gpuCount"`
	AvailableGpus int       `db:"available_gpus" json:"availableGpus"`
	AllocatedGpus int       `db:"allocated_gpus" json:"allocatedGpus"`
	OsImage       string    `db:"os_image" json:"osImage"`
	KernelVersion string    `db:"kernel_version" json:"kernelVersion"`
	NodeLabels    string    `db:"node_labels" json:"nodeLabels"`
	NodeTaints    string    `db:"node_taints" json:"nodeTaints"`
	LastHeartbeat time.Time `db:"last_heartbeat" json:"lastHeartbeat"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time `db:"updated_at" json:"updatedAt"`
}

// VtGpuNodesModel GPU节点模型操作接口
type VtGpuNodesModel interface {
	Insert(data *VtGpuNodes) (sql.Result, error)
	FindOne(id int64) (*VtGpuNodes, error)
	Update(data *VtGpuNodes) error
	Delete(id int64) error
	FindAll(page, pageSize int, clusterId int64, status, nodeType, search string) ([]*VtGpuNodes, int64, error)
	FindByClusterId(clusterId int64, page, pageSize int, status, nodeType string) ([]*VtGpuNodes, int64, error)
}

// vtGpuNodesModelImpl GPU节点模型实现
type vtGpuNodesModelImpl struct {
	conn *sql.DB
}

// NewVtGpuNodesModel 创建GPU节点模型实例
func NewVtGpuNodesModel(conn *sql.DB) VtGpuNodesModel {
	return &vtGpuNodesModelImpl{
		conn: conn,
	}
}

func (m *vtGpuNodesModelImpl) Insert(data *VtGpuNodes) (sql.Result, error) {
	query := `INSERT INTO vt_gpu_nodes (
		cluster_id, name, hostname, internal_ip, external_ip, status, node_type,
		cpu_cores, memory_gb, storage_gb, gpu_count, available_gpus, allocated_gpus,
		os_image, kernel_version, node_labels, node_taints, last_heartbeat
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return m.conn.Exec(query,
		data.ClusterId, data.Name, data.Hostname, data.InternalIp, data.ExternalIp, data.Status, data.NodeType,
		data.CpuCores, data.MemoryGb, data.StorageGb, data.GpuCount, data.AvailableGpus, data.AllocatedGpus,
		data.OsImage, data.KernelVersion, data.NodeLabels, data.NodeTaints, data.LastHeartbeat,
	)
}

func (m *vtGpuNodesModelImpl) FindOne(id int64) (*VtGpuNodes, error) {
	var node VtGpuNodes
	query := `SELECT id, cluster_id, name, hostname, internal_ip, external_ip, status, node_type,
		cpu_cores, memory_gb, storage_gb, gpu_count, available_gpus, allocated_gpus,
		os_image, kernel_version, node_labels, node_taints, last_heartbeat, created_at, updated_at
		FROM vt_gpu_nodes WHERE id = ?`

	err := m.conn.QueryRow(query, id).Scan(
		&node.Id, &node.ClusterId, &node.Name, &node.Hostname, &node.InternalIp, &node.ExternalIp, &node.Status, &node.NodeType,
		&node.CpuCores, &node.MemoryGb, &node.StorageGb, &node.GpuCount, &node.AvailableGpus, &node.AllocatedGpus,
		&node.OsImage, &node.KernelVersion, &node.NodeLabels, &node.NodeTaints, &node.LastHeartbeat, &node.CreatedAt, &node.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &node, nil
}

func (m *vtGpuNodesModelImpl) Update(data *VtGpuNodes) error {
	query := `UPDATE vt_gpu_nodes SET 
		cluster_id = ?, name = ?, hostname = ?, internal_ip = ?, external_ip = ?, status = ?, node_type = ?,
		cpu_cores = ?, memory_gb = ?, storage_gb = ?, gpu_count = ?, available_gpus = ?, allocated_gpus = ?,
		os_image = ?, kernel_version = ?, node_labels = ?, node_taints = ?, last_heartbeat = ?, updated_at = NOW()
		WHERE id = ?`

	_, err := m.conn.Exec(query,
		data.ClusterId, data.Name, data.Hostname, data.InternalIp, data.ExternalIp, data.Status, data.NodeType,
		data.CpuCores, data.MemoryGb, data.StorageGb, data.GpuCount, data.AvailableGpus, data.AllocatedGpus,
		data.OsImage, data.KernelVersion, data.NodeLabels, data.NodeTaints, data.LastHeartbeat, data.Id,
	)

	return err
}

func (m *vtGpuNodesModelImpl) Delete(id int64) error {
	query := `DELETE FROM vt_gpu_nodes WHERE id = ?`
	_, err := m.conn.Exec(query, id)
	return err
}

func (m *vtGpuNodesModelImpl) FindAll(page, pageSize int, clusterId int64, status, nodeType, search string) ([]*VtGpuNodes, int64, error) {
	offset := (page - 1) * pageSize

	// 构建WHERE条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if clusterId > 0 {
		whereClause += " AND cluster_id = ?"
		args = append(args, clusterId)
	}
	if status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}
	if nodeType != "" {
		whereClause += " AND node_type = ?"
		args = append(args, nodeType)
	}
	if search != "" {
		whereClause += " AND (name LIKE ? OR hostname LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	// 查询总数
	countQuery := "SELECT COUNT(*) FROM vt_gpu_nodes " + whereClause
	var total int64
	err := m.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	query := `SELECT id, cluster_id, name, hostname, internal_ip, external_ip, status, node_type,
		cpu_cores, memory_gb, storage_gb, gpu_count, available_gpus, allocated_gpus,
		os_image, kernel_version, node_labels, node_taints, last_heartbeat, created_at, updated_at
		FROM vt_gpu_nodes ` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`

	args = append(args, pageSize, offset)
	rows, err := m.conn.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var nodes []*VtGpuNodes
	for rows.Next() {
		var node VtGpuNodes
		err := rows.Scan(
			&node.Id, &node.ClusterId, &node.Name, &node.Hostname, &node.InternalIp, &node.ExternalIp, &node.Status, &node.NodeType,
			&node.CpuCores, &node.MemoryGb, &node.StorageGb, &node.GpuCount, &node.AvailableGpus, &node.AllocatedGpus,
			&node.OsImage, &node.KernelVersion, &node.NodeLabels, &node.NodeTaints, &node.LastHeartbeat, &node.CreatedAt, &node.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		nodes = append(nodes, &node)
	}

	return nodes, total, nil
}

func (m *vtGpuNodesModelImpl) FindByClusterId(clusterId int64, page, pageSize int, status, nodeType string) ([]*VtGpuNodes, int64, error) {
	return m.FindAll(page, pageSize, clusterId, status, nodeType, "")
}
