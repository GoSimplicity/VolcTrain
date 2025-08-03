//go:build unit
// +build unit

package model

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"api/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// GPUDevicesModelTestSuite GPU设备模型测试套件
type GPUDevicesModelTestSuite struct {
	suite.Suite
	db       *sql.DB
	mock     sqlmock.Sqlmock
	gpuModel model.VtGpuDevicesModel
}

// SetupTest 每个测试前的初始化
func (suite *GPUDevicesModelTestSuite) SetupTest() {
	var err error
	suite.db, suite.mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	suite.Require().NoError(err)

	suite.gpuModel = model.NewVtGpuDevicesModel(suite.db)
}

// TearDownTest 每个测试后的清理
func (suite *GPUDevicesModelTestSuite) TearDownTest() {
	suite.db.Close()
}

// TestInsertGPUDevice 测试插入GPU设备
func (suite *GPUDevicesModelTestSuite) TestInsertGPUDevice() {
	// 准备测试数据
	now := time.Now()
	gpuDevice := &model.VtGpuDevices{
		ClusterId:      1,
		NodeId:         1,
		DeviceIndex:    0,
		DeviceUuid:     "GPU-12345678-1234-1234-1234-123456789abc",
		DeviceName:     "NVIDIA GeForce RTX 3080",
		Brand:          "NVIDIA",
		Model:          "RTX 3080",
		Architecture:   "Ampere",
		MemoryTotalMb:  10240,
		MemoryFreeMb:   10240,
		MemoryUsedMb:   0,
		PowerDrawW:     50,
		PowerLimitW:    320,
		TemperatureC:   45,
		UtilizationGpu: 0,
		UtilizationMem: 0,
		Status:         "available",
		HealthStatus:   "healthy",
		PcieBusId:      "0000:01:00.0",
		CudaVersion:    "11.8",
		DriverVersion:  "520.61.05",
		LastHeartbeat:  now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// 设置Mock期望
	suite.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO vt_gpu_devices")).
		WithArgs(
			gpuDevice.ClusterId,
			gpuDevice.NodeId,
			gpuDevice.DeviceIndex,
			gpuDevice.DeviceUuid,
			gpuDevice.DeviceName,
			gpuDevice.Brand,
			gpuDevice.Model,
			gpuDevice.Architecture,
			gpuDevice.MemoryTotalMb,
			gpuDevice.MemoryFreeMb,
			gpuDevice.MemoryUsedMb,
			gpuDevice.PowerDrawW,
			gpuDevice.PowerLimitW,
			gpuDevice.TemperatureC,
			gpuDevice.UtilizationGpu,
			gpuDevice.UtilizationMem,
			gpuDevice.Status,
			gpuDevice.HealthStatus,
			gpuDevice.PcieBusId,
			gpuDevice.CudaVersion,
			gpuDevice.DriverVersion,
			sqlmock.AnyArg(), // LastHeartbeat
		).
		WillReturnResult(sqlmock.NewResult(1001, 1))

	// 执行测试
	result, err := suite.gpuModel.Insert(gpuDevice)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)

	lastInsertId, err := result.LastInsertId()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1001), lastInsertId)

	// 验证所有期望都被满足
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// TestFindOneGPUDevice 测试根据ID查找GPU设备
func (suite *GPUDevicesModelTestSuite) TestFindOneGPUDevice() {
	// 准备测试数据
	deviceID := int64(1001)
	expectedTime := time.Now()

	// 设置Mock期望
	rows := sqlmock.NewRows([]string{
		"id", "cluster_id", "node_id", "device_index", "device_uuid", "device_name",
		"brand", "model", "architecture", "memory_total_mb", "memory_free_mb", "memory_used_mb",
		"power_draw_w", "power_limit_w", "temperature_c", "utilization_gpu", "utilization_mem",
		"status", "health_status", "pcie_bus_id", "cuda_version", "driver_version",
		"allocation_id", "allocated_job_id", "allocated_user_id", "allocated_at",
		"last_heartbeat", "created_at", "updated_at",
	}).AddRow(
		deviceID, int64(1), int64(1), 0, "GPU-12345678-1234-1234-1234-123456789abc", "NVIDIA GeForce RTX 3080",
		"NVIDIA", "RTX 3080", "Ampere", 10240, 8192, 2048,
		150, 320, 65, 85, 70,
		"allocated", "healthy", "0000:01:00.0", "11.8", "520.61.05",
		int64(5001), int64(2001), int64(1001), expectedTime,
		expectedTime, expectedTime, expectedTime,
	)

	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs(deviceID).
		WillReturnRows(rows)

	// 执行测试
	device, err := suite.gpuModel.FindOne(deviceID)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), device)
	assert.Equal(suite.T(), deviceID, device.Id)
	assert.Equal(suite.T(), "NVIDIA GeForce RTX 3080", device.DeviceName)
	assert.Equal(suite.T(), "NVIDIA", device.Brand)
	assert.Equal(suite.T(), "RTX 3080", device.Model)
	assert.Equal(suite.T(), "allocated", device.Status)
	assert.Equal(suite.T(), int(10240), device.MemoryTotalMb)
	assert.Equal(suite.T(), int64(2001), device.AllocatedJobId)

	// 验证所有期望都被满足
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// TestUpdateGPUDeviceStatus 测试更新GPU设备状态
func (suite *GPUDevicesModelTestSuite) TestUpdateGPUDeviceStatus() {
	deviceID := int64(1001)
	newStatus := "maintenance"

	// 设置Mock期望
	suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE vt_gpu_devices SET status = ?, updated_at = ? WHERE id = ?")).
		WithArgs(newStatus, sqlmock.AnyArg(), deviceID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 执行测试
	err := suite.gpuModel.UpdateStatus(deviceID, newStatus)

	// 验证结果
	assert.NoError(suite.T(), err)

	// 验证所有期望都被满足
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// TestFindAvailableGPUDevices 测试查找可用的GPU设备
func (suite *GPUDevicesModelTestSuite) TestFindAvailableGPUDevices() {
	clusterID := int64(1)
	gpuCount := 2
	expectedTime := time.Now()

	// 设置Mock期望
	rows := sqlmock.NewRows([]string{
		"id", "cluster_id", "node_id", "device_index", "device_uuid", "device_name",
		"brand", "model", "status", "health_status", "memory_total_mb", "memory_free_mb",
		"created_at", "updated_at",
	}).
		AddRow(
			int64(1001), clusterID, int64(1), 0, "GPU-device-1", "NVIDIA T4",
			"NVIDIA", "T4", "available", "healthy", 15360, 15360,
			expectedTime, expectedTime,
		).
		AddRow(
			int64(1002), clusterID, int64(1), 1, "GPU-device-2", "NVIDIA T4",
			"NVIDIA", "T4", "available", "healthy", 15360, 15360,
			expectedTime, expectedTime,
		)

	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs(clusterID, gpuCount).
		WillReturnRows(rows)

	// 执行测试
	devices, err := suite.gpuModel.FindAvailableDevices(clusterID, gpuCount)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), devices, 2)
	assert.Equal(suite.T(), "available", devices[0].Status)
	assert.Equal(suite.T(), "available", devices[1].Status)
	assert.Equal(suite.T(), "NVIDIA", devices[0].Brand)
	assert.Equal(suite.T(), "T4", devices[0].Model)

	// 验证所有期望都被满足
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// TestFindGPUDevicesByNodeId 测试根据节点ID查找GPU设备
func (suite *GPUDevicesModelTestSuite) TestFindGPUDevicesByNodeId() {
	nodeID := int64(1)
	page, pageSize := 1, 10
	status := "available"
	expectedTime := time.Now()

	// 设置Mock期望 - 数据查询
	rows := sqlmock.NewRows([]string{
		"id", "cluster_id", "node_id", "device_index", "device_name", "brand", "model",
		"status", "health_status", "memory_total_mb", "utilization_gpu", "created_at",
	}).
		AddRow(
			int64(1001), int64(1), nodeID, 0, "GPU Device 1", "NVIDIA", "RTX 3080",
			status, "healthy", 10240, 0, expectedTime,
		).
		AddRow(
			int64(1002), int64(1), nodeID, 1, "GPU Device 2", "NVIDIA", "RTX 3080",
			status, "healthy", 10240, 0, expectedTime,
		)

	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WithArgs(nodeID, status, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	// 设置Mock期望 - 计数查询
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(int64(2))
	suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*)")).
		WithArgs(nodeID, status).
		WillReturnRows(countRows)

	// 执行测试
	devices, total, err := suite.gpuModel.FindByNodeId(nodeID, page, pageSize, status)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), devices, 2)
	assert.Equal(suite.T(), int64(2), total)
	assert.Equal(suite.T(), nodeID, devices[0].NodeId)
	assert.Equal(suite.T(), nodeID, devices[1].NodeId)
	assert.Equal(suite.T(), status, devices[0].Status)

	// 验证所有期望都被满足
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// TestDeleteGPUDevice 测试删除GPU设备
func (suite *GPUDevicesModelTestSuite) TestDeleteGPUDevice() {
	deviceID := int64(1001)

	// 设置Mock期望
	suite.mock.ExpectExec(regexp.QuoteMeta("DELETE FROM vt_gpu_devices WHERE id = ?")).
		WithArgs(deviceID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 执行测试
	err := suite.gpuModel.Delete(deviceID)

	// 验证结果
	assert.NoError(suite.T(), err)

	// 验证所有期望都被满足
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// TestGPUDeviceModelErrors 测试各种错误情况
func (suite *GPUDevicesModelTestSuite) TestGPUDeviceModelErrors() {
	// 测试插入失败
	suite.Run("InsertFailure", func() {
		gpuDevice := &model.VtGpuDevices{
			ClusterId:   1,
			NodeId:      1,
			DeviceIndex: 0,
			DeviceUuid:  "test-uuid",
			DeviceName:  "Test GPU",
			Status:      "available",
		}

		suite.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO vt_gpu_devices")).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(sql.ErrConnDone)

		result, err := suite.gpuModel.Insert(gpuDevice)
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), result)
		assert.Equal(suite.T(), sql.ErrConnDone, err)
	})

	// 测试查找不存在的设备
	suite.Run("FindOneNotFound", func() {
		deviceID := int64(99999)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
			WithArgs(deviceID).
			WillReturnError(sql.ErrNoRows)

		device, err := suite.gpuModel.FindOne(deviceID)
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), device)
		assert.Equal(suite.T(), sql.ErrNoRows, err)
	})

	// 测试更新状态失败
	suite.Run("UpdateStatusFailure", func() {
		deviceID := int64(1001)
		newStatus := "error"

		suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE vt_gpu_devices SET status = ?, updated_at = ? WHERE id = ?")).
			WithArgs(newStatus, sqlmock.AnyArg(), deviceID).
			WillReturnError(sql.ErrTxDone)

		err := suite.gpuModel.UpdateStatus(deviceID, newStatus)
		assert.Error(suite.T(), err)
		assert.Equal(suite.T(), sql.ErrTxDone, err)
	})
}

// TestGPUDeviceResourceMonitoring 测试GPU资源监控相关操作
func (suite *GPUDevicesModelTestSuite) TestGPUDeviceResourceMonitoring() {
	// 测试更新GPU设备利用率和内存使用情况
	deviceID := int64(1001)
	updatedDevice := &model.VtGpuDevices{
		Id:             deviceID,
		MemoryFreeMb:   8192,
		MemoryUsedMb:   2048,
		UtilizationGpu: 85,
		UtilizationMem: 70,
		TemperatureC:   75,
		PowerDrawW:     200,
		UpdatedAt:      time.Now(),
	}

	// 设置Mock期望
	suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE vt_gpu_devices")).
		WithArgs(
			updatedDevice.MemoryFreeMb,
			updatedDevice.MemoryUsedMb,
			updatedDevice.UtilizationGpu,
			updatedDevice.UtilizationMem,
			updatedDevice.TemperatureC,
			updatedDevice.PowerDrawW,
			sqlmock.AnyArg(), // updated_at
			deviceID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 执行测试
	err := suite.gpuModel.Update(updatedDevice)

	// 验证结果
	assert.NoError(suite.T(), err)

	// 验证所有期望都被满足
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// TestGPUDeviceAllocation 测试GPU设备分配相关操作
func (suite *GPUDevicesModelTestSuite) TestGPUDeviceAllocation() {
	// 模拟分配GPU设备给特定作业
	deviceID := int64(1001)
	jobID := int64(2001)
	userID := int64(1001)
	allocationID := int64(5001)

	// 设置Mock期望 - 更新分配信息
	suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE vt_gpu_devices SET allocation_id = ?, allocated_job_id = ?, allocated_user_id = ?, allocated_at = ?, status = ?, updated_at = ? WHERE id = ?")).
		WithArgs(allocationID, jobID, userID, sqlmock.AnyArg(), "allocated", sqlmock.AnyArg(), deviceID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 创建包含分配信息的设备对象
	allocatedDevice := &model.VtGpuDevices{
		Id:              deviceID,
		AllocationId:    allocationID,
		AllocatedJobId:  jobID,
		AllocatedUserId: userID,
		AllocatedAt:     time.Now(),
		Status:          "allocated",
		UpdatedAt:       time.Now(),
	}

	// 执行测试
	err := suite.gpuModel.Update(allocatedDevice)

	// 验证结果
	assert.NoError(suite.T(), err)

	// 验证所有期望都被满足
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// TestGPUDeviceHealthCheck 测试GPU设备健康检查
func (suite *GPUDevicesModelTestSuite) TestGPUDeviceHealthCheck() {
	// 测试更新设备健康状态和心跳时间
	deviceID := int64(1001)

	// 设置Mock期望
	suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE vt_gpu_devices SET health_status = ?, last_heartbeat = ?, updated_at = ? WHERE id = ?")).
		WithArgs("healthy", sqlmock.AnyArg(), sqlmock.AnyArg(), deviceID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 创建包含健康检查信息的设备对象
	healthyDevice := &model.VtGpuDevices{
		Id:            deviceID,
		HealthStatus:  "healthy",
		LastHeartbeat: time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 执行测试
	err := suite.gpuModel.Update(healthyDevice)

	// 验证结果
	assert.NoError(suite.T(), err)

	// 验证所有期望都被满足
	assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// 运行GPU设备模型测试套件
func TestGPUDevicesModelSuite(t *testing.T) {
	suite.Run(t, new(GPUDevicesModelTestSuite))
}

// 基准测试
func BenchmarkGPUDevicesModelFindOne(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("创建sqlmock失败: %v", err)
	}
	defer db.Close()

	gpuModel := model.NewVtGpuDevicesModel(db)

	// 设置Mock期望
	rows := sqlmock.NewRows([]string{
		"id", "cluster_id", "node_id", "device_index", "device_name", "brand",
		"status", "health_status", "created_at", "updated_at",
	}).AddRow(
		int64(1001), int64(1), int64(1), 0, "Test GPU", "NVIDIA",
		"available", "healthy", time.Now(), time.Now(),
	)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT")).
		WillReturnRows(rows)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := gpuModel.FindOne(1001)
		if err != nil {
			b.Errorf("查找失败: %v", err)
		}
	}
}
