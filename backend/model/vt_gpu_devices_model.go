package model

import (
	"database/sql"
	"time"
)

// VtGpuDevices GPU设备表模型
type VtGpuDevices struct {
	Id              int64     `db:"id" json:"id"`
	ClusterId       int64     `db:"cluster_id" json:"clusterId"`
	NodeId          int64     `db:"node_id" json:"nodeId"`
	DeviceIndex     int       `db:"device_index" json:"deviceIndex"`
	DeviceUuid      string    `db:"device_uuid" json:"deviceUuid"`
	DeviceName      string    `db:"device_name" json:"deviceName"`
	Brand           string    `db:"brand" json:"brand"`
	Model           string    `db:"model" json:"model"`
	Architecture    string    `db:"architecture" json:"architecture"`
	MemoryTotalMb   int       `db:"memory_total_mb" json:"memoryTotalMb"`
	MemoryFreeMb    int       `db:"memory_free_mb" json:"memoryFreeMb"`
	MemoryUsedMb    int       `db:"memory_used_mb" json:"memoryUsedMb"`
	PowerDrawW      int       `db:"power_draw_w" json:"powerDrawW"`
	PowerLimitW     int       `db:"power_limit_w" json:"powerLimitW"`
	TemperatureC    int       `db:"temperature_c" json:"temperatureC"`
	UtilizationGpu  int       `db:"utilization_gpu" json:"utilizationGpu"`
	UtilizationMem  int       `db:"utilization_mem" json:"utilizationMem"`
	Status          string    `db:"status" json:"status"`
	HealthStatus    string    `db:"health_status" json:"healthStatus"`
	PcieBusId       string    `db:"pcie_bus_id" json:"pcieBusId"`
	CudaVersion     string    `db:"cuda_version" json:"cudaVersion"`
	DriverVersion   string    `db:"driver_version" json:"driverVersion"`
	AllocationId    int64     `db:"allocation_id" json:"allocationId"`
	AllocatedJobId  int64     `db:"allocated_job_id" json:"allocatedJobId"`
	AllocatedUserId int64     `db:"allocated_user_id" json:"allocatedUserId"`
	AllocatedAt     time.Time `db:"allocated_at" json:"allocatedAt"`
	LastHeartbeat   time.Time `db:"last_heartbeat" json:"lastHeartbeat"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time `db:"updated_at" json:"updatedAt"`
}

// VtGpuDevicesModel GPU设备模型操作接口
type VtGpuDevicesModel interface {
	Insert(data *VtGpuDevices) (sql.Result, error)
	FindOne(id int64) (*VtGpuDevices, error)
	Update(data *VtGpuDevices) error
	Delete(id int64) error
	FindAll(page, pageSize int, clusterId, nodeId int64, status, brand, model, healthStatus, search string) ([]*VtGpuDevices, int64, error)
	FindByNodeId(nodeId int64, page, pageSize int, status string) ([]*VtGpuDevices, int64, error)
	FindAvailableDevices(clusterId int64, gpuCount int) ([]*VtGpuDevices, error)
	UpdateStatus(id int64, status string) error
}

// vtGpuDevicesModelImpl GPU设备模型实现
type vtGpuDevicesModelImpl struct {
	conn *sql.DB
}

// NewVtGpuDevicesModel 创建GPU设备模型实例
func NewVtGpuDevicesModel(conn *sql.DB) VtGpuDevicesModel {
	return &vtGpuDevicesModelImpl{
		conn: conn,
	}
}

func (m *vtGpuDevicesModelImpl) Insert(data *VtGpuDevices) (sql.Result, error) {
	query := `INSERT INTO vt_gpu_devices (
		cluster_id, node_id, device_index, device_uuid, device_name, brand, model, architecture,
		memory_total_mb, memory_free_mb, memory_used_mb, power_draw_w, power_limit_w,
		temperature_c, utilization_gpu, utilization_mem, status, health_status,
		pcie_bus_id, cuda_version, driver_version, allocation_id, allocated_job_id,
		allocated_user_id, allocated_at, last_heartbeat
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return m.conn.Exec(query,
		data.ClusterId, data.NodeId, data.DeviceIndex, data.DeviceUuid, data.DeviceName, data.Brand, data.Model, data.Architecture,
		data.MemoryTotalMb, data.MemoryFreeMb, data.MemoryUsedMb, data.PowerDrawW, data.PowerLimitW,
		data.TemperatureC, data.UtilizationGpu, data.UtilizationMem, data.Status, data.HealthStatus,
		data.PcieBusId, data.CudaVersion, data.DriverVersion, data.AllocationId, data.AllocatedJobId,
		data.AllocatedUserId, data.AllocatedAt, data.LastHeartbeat,
	)
}

func (m *vtGpuDevicesModelImpl) FindOne(id int64) (*VtGpuDevices, error) {
	var device VtGpuDevices
	query := `SELECT id, cluster_id, node_id, device_index, device_uuid, device_name, brand, model, architecture,
		memory_total_mb, memory_free_mb, memory_used_mb, power_draw_w, power_limit_w,
		temperature_c, utilization_gpu, utilization_mem, status, health_status,
		pcie_bus_id, cuda_version, driver_version, allocation_id, allocated_job_id,
		allocated_user_id, allocated_at, last_heartbeat, created_at, updated_at
		FROM vt_gpu_devices WHERE id = ?`

	err := m.conn.QueryRow(query, id).Scan(
		&device.Id, &device.ClusterId, &device.NodeId, &device.DeviceIndex, &device.DeviceUuid, &device.DeviceName, &device.Brand, &device.Model, &device.Architecture,
		&device.MemoryTotalMb, &device.MemoryFreeMb, &device.MemoryUsedMb, &device.PowerDrawW, &device.PowerLimitW,
		&device.TemperatureC, &device.UtilizationGpu, &device.UtilizationMem, &device.Status, &device.HealthStatus,
		&device.PcieBusId, &device.CudaVersion, &device.DriverVersion, &device.AllocationId, &device.AllocatedJobId,
		&device.AllocatedUserId, &device.AllocatedAt, &device.LastHeartbeat, &device.CreatedAt, &device.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &device, nil
}

func (m *vtGpuDevicesModelImpl) Update(data *VtGpuDevices) error {
	query := `UPDATE vt_gpu_devices SET 
		cluster_id = ?, node_id = ?, device_index = ?, device_uuid = ?, device_name = ?, brand = ?, model = ?, architecture = ?,
		memory_total_mb = ?, memory_free_mb = ?, memory_used_mb = ?, power_draw_w = ?, power_limit_w = ?,
		temperature_c = ?, utilization_gpu = ?, utilization_mem = ?, status = ?, health_status = ?,
		pcie_bus_id = ?, cuda_version = ?, driver_version = ?, allocation_id = ?, allocated_job_id = ?,
		allocated_user_id = ?, allocated_at = ?, last_heartbeat = ?, updated_at = NOW()
		WHERE id = ?`

	_, err := m.conn.Exec(query,
		data.ClusterId, data.NodeId, data.DeviceIndex, data.DeviceUuid, data.DeviceName, data.Brand, data.Model, data.Architecture,
		data.MemoryTotalMb, data.MemoryFreeMb, data.MemoryUsedMb, data.PowerDrawW, data.PowerLimitW,
		data.TemperatureC, data.UtilizationGpu, data.UtilizationMem, data.Status, data.HealthStatus,
		data.PcieBusId, data.CudaVersion, data.DriverVersion, data.AllocationId, data.AllocatedJobId,
		data.AllocatedUserId, data.AllocatedAt, data.LastHeartbeat, data.Id,
	)

	return err
}

func (m *vtGpuDevicesModelImpl) Delete(id int64) error {
	query := `DELETE FROM vt_gpu_devices WHERE id = ?`
	_, err := m.conn.Exec(query, id)
	return err
}

func (m *vtGpuDevicesModelImpl) FindAll(page, pageSize int, clusterId, nodeId int64, status, brand, model, healthStatus, search string) ([]*VtGpuDevices, int64, error) {
	offset := (page - 1) * pageSize

	// 构建WHERE条件
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if clusterId > 0 {
		whereClause += " AND cluster_id = ?"
		args = append(args, clusterId)
	}
	if nodeId > 0 {
		whereClause += " AND node_id = ?"
		args = append(args, nodeId)
	}
	if status != "" {
		whereClause += " AND status = ?"
		args = append(args, status)
	}
	if brand != "" {
		whereClause += " AND brand = ?"
		args = append(args, brand)
	}
	if model != "" {
		whereClause += " AND model = ?"
		args = append(args, model)
	}
	if healthStatus != "" {
		whereClause += " AND health_status = ?"
		args = append(args, healthStatus)
	}
	if search != "" {
		whereClause += " AND (device_name LIKE ? OR device_uuid LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	// 查询总数
	countQuery := "SELECT COUNT(*) FROM vt_gpu_devices " + whereClause
	var total int64
	err := m.conn.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	query := `SELECT id, cluster_id, node_id, device_index, device_uuid, device_name, brand, model, architecture,
		memory_total_mb, memory_free_mb, memory_used_mb, power_draw_w, power_limit_w,
		temperature_c, utilization_gpu, utilization_mem, status, health_status,
		pcie_bus_id, cuda_version, driver_version, allocation_id, allocated_job_id,
		allocated_user_id, allocated_at, last_heartbeat, created_at, updated_at
		FROM vt_gpu_devices ` + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`

	args = append(args, pageSize, offset)
	rows, err := m.conn.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var devices []*VtGpuDevices
	for rows.Next() {
		var device VtGpuDevices
		err := rows.Scan(
			&device.Id, &device.ClusterId, &device.NodeId, &device.DeviceIndex, &device.DeviceUuid, &device.DeviceName, &device.Brand, &device.Model, &device.Architecture,
			&device.MemoryTotalMb, &device.MemoryFreeMb, &device.MemoryUsedMb, &device.PowerDrawW, &device.PowerLimitW,
			&device.TemperatureC, &device.UtilizationGpu, &device.UtilizationMem, &device.Status, &device.HealthStatus,
			&device.PcieBusId, &device.CudaVersion, &device.DriverVersion, &device.AllocationId, &device.AllocatedJobId,
			&device.AllocatedUserId, &device.AllocatedAt, &device.LastHeartbeat, &device.CreatedAt, &device.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		devices = append(devices, &device)
	}

	return devices, total, nil
}

func (m *vtGpuDevicesModelImpl) FindByNodeId(nodeId int64, page, pageSize int, status string) ([]*VtGpuDevices, int64, error) {
	return m.FindAll(page, pageSize, 0, nodeId, status, "", "", "", "")
}

func (m *vtGpuDevicesModelImpl) FindAvailableDevices(clusterId int64, gpuCount int) ([]*VtGpuDevices, error) {
	whereClause := "WHERE status = 'available' AND health_status = 'healthy'"
	args := []interface{}{}

	if clusterId > 0 {
		whereClause += " AND cluster_id = ?"
		args = append(args, clusterId)
	}

	query := `SELECT id, cluster_id, node_id, device_index, device_uuid, device_name, brand, model, architecture,
		memory_total_mb, memory_free_mb, memory_used_mb, power_draw_w, power_limit_w,
		temperature_c, utilization_gpu, utilization_mem, status, health_status,
		pcie_bus_id, cuda_version, driver_version, allocation_id, allocated_job_id,
		allocated_user_id, allocated_at, last_heartbeat, created_at, updated_at
		FROM vt_gpu_devices ` + whereClause + ` ORDER BY utilization_gpu ASC LIMIT ?`

	args = append(args, gpuCount)
	rows, err := m.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*VtGpuDevices
	for rows.Next() {
		var device VtGpuDevices
		err := rows.Scan(
			&device.Id, &device.ClusterId, &device.NodeId, &device.DeviceIndex, &device.DeviceUuid, &device.DeviceName, &device.Brand, &device.Model, &device.Architecture,
			&device.MemoryTotalMb, &device.MemoryFreeMb, &device.MemoryUsedMb, &device.PowerDrawW, &device.PowerLimitW,
			&device.TemperatureC, &device.UtilizationGpu, &device.UtilizationMem, &device.Status, &device.HealthStatus,
			&device.PcieBusId, &device.CudaVersion, &device.DriverVersion, &device.AllocationId, &device.AllocatedJobId,
			&device.AllocatedUserId, &device.AllocatedAt, &device.LastHeartbeat, &device.CreatedAt, &device.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		devices = append(devices, &device)
	}

	return devices, nil
}

func (m *vtGpuDevicesModelImpl) UpdateStatus(id int64, status string) error {
	query := `UPDATE vt_gpu_devices SET status = ?, updated_at = NOW() WHERE id = ?`
	_, err := m.conn.Exec(query, status, id)
	return err
}
