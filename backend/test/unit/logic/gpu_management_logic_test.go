//go:build unit
// +build unit

package logic

import (
	"context"
	"fmt"
	"testing"
	"time"

	"api/internal/logic/gpu_cluster"
	"api/internal/logic/gpu_device"
	"api/internal/svc"
	"api/internal/types"
	"api/model"
	"api/test/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// GPUManagementLogicTestSuite GPU管理逻辑测试套件
type GPUManagementLogicTestSuite struct {
	suite.Suite
	ctrl                *gomock.Controller
	mockGPUClusterModel *mocks.MockVtGpuClustersModel
	mockGPUDeviceModel  *mocks.MockVtGpuDevicesModel
	mockGPUNodeModel    *mocks.MockVtGpuNodesModel
	mockSvcCtx          *svc.ServiceContext
	clusterLogic        *gpu_cluster.CreateGpuClusterLogic
	deviceLogic         *gpu_device.AllocateGpuDeviceLogic
	ctx                 context.Context
}

// SetupTest 每个测试前的初始化
func (suite *GPUManagementLogicTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockGPUClusterModel = mocks.NewMockVtGpuClustersModel(suite.ctrl)
	suite.mockGPUDeviceModel = mocks.NewMockVtGpuDevicesModel(suite.ctrl)
	suite.mockGPUNodeModel = mocks.NewMockVtGpuNodesModel(suite.ctrl)

	suite.ctx = context.Background()

	// 创建模拟的服务上下文
	suite.mockSvcCtx = &svc.ServiceContext{
		VtGpuClustersModel: suite.mockGPUClusterModel,
		VtGpuDevicesModel:  suite.mockGPUDeviceModel,
		VtGpuNodesModel:    suite.mockGPUNodeModel,
	}

	suite.clusterLogic = gpu_cluster.NewCreateGpuClusterLogic(suite.ctx, suite.mockSvcCtx)
	suite.deviceLogic = gpu_device.NewAllocateGpuDeviceLogic(suite.ctx, suite.mockSvcCtx)
}

// TearDownTest 每个测试后的清理
func (suite *GPUManagementLogicTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

// TestCreateGPUClusterSuccess 测试成功创建GPU集群
func (suite *GPUManagementLogicTestSuite) TestCreateGPUClusterSuccess() {
	// 准备测试数据
	req := &types.CreateGpuClusterReq{
		Name:        "test-gpu-cluster",
		Description: "测试GPU集群",
		Kubeconfig:  "test-kubeconfig-content",
		Region:      "us-west-1",
		Zone:        "us-west-1a",
		Provider:    "aws",
		ClusterType: "kubernetes",
	}

	expectedCluster := &model.VtGpuClusters{
		Id:          2001,
		Name:        req.Name,
		Description: req.Description,
		Kubeconfig:  req.Kubeconfig,
		Region:      req.Region,
		Zone:        req.Zone,
		Provider:    req.Provider,
		ClusterType: req.ClusterType,
		Status:      "pending",
		NodeCount:   0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 设置Mock期望
	suite.mockGPUClusterModel.EXPECT().
		CheckNameExists(req.Name).
		Return(false, nil).
		Times(1)

	suite.mockGPUClusterModel.EXPECT().
		Insert(gomock.Any()).
		DoAndReturn(func(cluster *model.VtGpuClusters) (*model.VtGpuClusters, error) {
			cluster.Id = 2001
			cluster.Status = "pending"
			cluster.CreatedAt = time.Now()
			cluster.UpdatedAt = time.Now()
			return cluster, nil
		}).
		Times(1)

	// 执行测试
	resp, err := suite.clusterLogic.CreateGpuCluster(req)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), int64(2001), resp.ClusterId)
	assert.Equal(suite.T(), req.Name, resp.Name)
	assert.Equal(suite.T(), "pending", resp.Status)
}

// TestCreateGPUClusterNameExists 测试集群名称已存在
func (suite *GPUManagementLogicTestSuite) TestCreateGPUClusterNameExists() {
	req := &types.CreateGpuClusterReq{
		Name:        "existing-cluster",
		Description: "已存在的集群",
		Kubeconfig:  "test-kubeconfig",
	}

	// 设置Mock期望 - 名称已存在
	suite.mockGPUClusterModel.EXPECT().
		CheckNameExists(req.Name).
		Return(true, nil).
		Times(1)

	// 执行测试
	resp, err := suite.clusterLogic.CreateGpuCluster(req)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "集群名称已存在")
}

// TestCreateGPUClusterInvalidKubeconfig 测试无效的Kubeconfig
func (suite *GPUManagementLogicTestSuite) TestCreateGPUClusterInvalidKubeconfig() {
	req := &types.CreateGpuClusterReq{
		Name:        "test-cluster",
		Description: "测试集群",
		Kubeconfig:  "invalid-kubeconfig-content",
	}

	// 设置Mock期望
	suite.mockGPUClusterModel.EXPECT().
		CheckNameExists(req.Name).
		Return(false, nil).
		Times(1)

	suite.mockGPUClusterModel.EXPECT().
		ValidateKubeconfig(req.Kubeconfig).
		Return(false, fmt.Errorf("Kubeconfig格式无效")).
		Times(1)

	// 执行测试
	resp, err := suite.clusterLogic.CreateGpuCluster(req)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "Kubeconfig格式无效")
}

// TestAllocateGPUDeviceSuccess 测试成功分配GPU设备
func (suite *GPUManagementLogicTestSuite) TestAllocateGPUDeviceSuccess() {
	// 准备测试数据
	req := &types.AllocateGpuDeviceReq{
		JobId:     3001,
		ClusterId: 2001,
		GPUCount:  2,
		GPUType:   "T4",
		Requirements: map[string]interface{}{
			"memory_gb":    16,
			"cuda_version": "11.2",
		},
	}

	availableDevices := []*model.VtGpuDevices{
		{
			Id:        5001,
			ClusterId: 2001,
			NodeId:    4001,
			DeviceId:  "gpu-0",
			GPUType:   "T4",
			Status:    "available",
			MemoryGb:  16,
		},
		{
			Id:        5002,
			ClusterId: 2001,
			NodeId:    4001,
			DeviceId:  "gpu-1",
			GPUType:   "T4",
			Status:    "available",
			MemoryGb:  16,
		},
	}

	// 设置Mock期望
	suite.mockGPUDeviceModel.EXPECT().
		FindAvailableDevices(req.ClusterId, req.GPUType, req.GPUCount).
		Return(availableDevices, nil).
		Times(1)

	suite.mockGPUDeviceModel.EXPECT().
		AllocateDevices(gomock.Any(), req.JobId).
		DoAndReturn(func(deviceIds []int64, jobId int64) error {
			assert.Len(suite.T(), deviceIds, 2)
			assert.Equal(suite.T(), req.JobId, jobId)
			return nil
		}).
		Times(1)

	// 执行测试
	resp, err := suite.deviceLogic.AllocateGpuDevice(req)

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Len(suite.T(), resp.AllocatedDevices, 2)
	assert.Equal(suite.T(), "allocated", resp.Status)
}

// TestAllocateGPUDeviceInsufficientResources 测试GPU资源不足
func (suite *GPUManagementLogicTestSuite) TestAllocateGPUDeviceInsufficientResources() {
	req := &types.AllocateGpuDeviceReq{
		JobId:     3001,
		ClusterId: 2001,
		GPUCount:  4, // 请求4个GPU
		GPUType:   "V100",
	}

	// 设置Mock期望 - 只有1个可用设备
	availableDevices := []*model.VtGpuDevices{
		{
			Id:        5001,
			ClusterId: 2001,
			GPUType:   "V100",
			Status:    "available",
		},
	}

	suite.mockGPUDeviceModel.EXPECT().
		FindAvailableDevices(req.ClusterId, req.GPUType, req.GPUCount).
		Return(availableDevices, nil).
		Times(1)

	// 执行测试
	resp, err := suite.deviceLogic.AllocateGpuDevice(req)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "GPU资源不足")
}

// TestAllocateGPUDeviceClusterNotFound 测试集群不存在
func (suite *GPUManagementLogicTestSuite) TestAllocateGPUDeviceClusterNotFound() {
	req := &types.AllocateGpuDeviceReq{
		JobId:     3001,
		ClusterId: 9999, // 不存在的集群ID
		GPUCount:  1,
		GPUType:   "T4",
	}

	// 设置Mock期望
	suite.mockGPUDeviceModel.EXPECT().
		FindAvailableDevices(req.ClusterId, req.GPUType, req.GPUCount).
		Return(nil, fmt.Errorf("集群不存在")).
		Times(1)

	// 执行测试
	resp, err := suite.deviceLogic.AllocateGpuDevice(req)

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Contains(suite.T(), err.Error(), "集群不存在")
}

// TestReleaseGPUDevice 测试释放GPU设备
func (suite *GPUManagementLogicTestSuite) TestReleaseGPUDevice() {
	jobId := int64(3001)
	allocatedDevices := []*model.VtGpuDevices{
		{Id: 5001, Status: "allocated", JobId: &jobId},
		{Id: 5002, Status: "allocated", JobId: &jobId},
	}

	// 设置Mock期望
	suite.mockGPUDeviceModel.EXPECT().
		FindByJobId(jobId).
		Return(allocatedDevices, nil).
		Times(1)

	suite.mockGPUDeviceModel.EXPECT().
		ReleaseDevices(gomock.Any()).
		DoAndReturn(func(deviceIds []int64) error {
			assert.Len(suite.T(), deviceIds, 2)
			assert.Contains(suite.T(), deviceIds, int64(5001))
			assert.Contains(suite.T(), deviceIds, int64(5002))
			return nil
		}).
		Times(1)

	// 创建释放设备逻辑
	releaseLogic := gpu_device.NewReleaseGpuDeviceLogic(suite.ctx, suite.mockSvcCtx)

	// 执行测试
	err := releaseLogic.ReleaseGpuDevice(&types.ReleaseGpuDeviceReq{JobId: jobId})

	// 验证结果
	assert.NoError(suite.T(), err)
}

// TestGetGPUClusterStatus 测试获取GPU集群状态
func (suite *GPUManagementLogicTestSuite) TestGetGPUClusterStatus() {
	clusterId := int64(2001)

	expectedCluster := &model.VtGpuClusters{
		Id:        clusterId,
		Name:      "test-cluster",
		Status:    "active",
		NodeCount: 3,
		Region:    "us-west-1",
		Provider:  "aws",
	}

	clusterNodes := []*model.VtGpuNodes{
		{Id: 4001, ClusterId: clusterId, Status: "ready", GPUCount: 4},
		{Id: 4002, ClusterId: clusterId, Status: "ready", GPUCount: 4},
		{Id: 4003, ClusterId: clusterId, Status: "ready", GPUCount: 4},
	}

	clusterDevices := []*model.VtGpuDevices{
		{Id: 5001, ClusterId: clusterId, Status: "available"},
		{Id: 5002, ClusterId: clusterId, Status: "available"},
		{Id: 5003, ClusterId: clusterId, Status: "allocated"},
	}

	// 设置Mock期望
	suite.mockGPUClusterModel.EXPECT().
		FindOne(clusterId).
		Return(expectedCluster, nil).
		Times(1)

	suite.mockGPUNodeModel.EXPECT().
		FindByClusterId(clusterId).
		Return(clusterNodes, nil).
		Times(1)

	suite.mockGPUDeviceModel.EXPECT().
		FindByClusterId(clusterId).
		Return(clusterDevices, nil).
		Times(1)

	// 创建获取集群状态逻辑
	getLogic := gpu_cluster.NewGetGpuClusterLogic(suite.ctx, suite.mockSvcCtx)

	// 执行测试
	resp, err := getLogic.GetGpuCluster(&types.GetGpuClusterReq{ClusterId: clusterId})

	// 验证结果
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), clusterId, resp.Cluster.Id)
	assert.Equal(suite.T(), "test-cluster", resp.Cluster.Name)
	assert.Equal(suite.T(), "active", resp.Cluster.Status)
	assert.Len(suite.T(), resp.Nodes, 3)
	assert.Len(suite.T(), resp.Devices, 3)
}

// TestUpdateGPUClusterStatus 测试更新GPU集群状态
func (suite *GPUManagementLogicTestSuite) TestUpdateGPUClusterStatus() {
	clusterId := int64(2001)
	newStatus := "maintenance"

	existingCluster := &model.VtGpuClusters{
		Id:     clusterId,
		Name:   "test-cluster",
		Status: "active",
	}

	// 设置Mock期望
	suite.mockGPUClusterModel.EXPECT().
		FindOne(clusterId).
		Return(existingCluster, nil).
		Times(1)

	suite.mockGPUClusterModel.EXPECT().
		UpdateStatus(clusterId, newStatus).
		Return(nil).
		Times(1)

	// 创建更新集群逻辑
	updateLogic := gpu_cluster.NewUpdateGpuClusterLogic(suite.ctx, suite.mockSvcCtx)

	// 执行测试
	err := updateLogic.UpdateClusterStatus(clusterId, newStatus)

	// 验证结果
	assert.NoError(suite.T(), err)
}

// TestDeleteGPUCluster 测试删除GPU集群
func (suite *GPUManagementLogicTestSuite) TestDeleteGPUCluster() {
	clusterId := int64(2001)

	existingCluster := &model.VtGpuClusters{
		Id:     clusterId,
		Name:   "test-cluster",
		Status: "inactive",
	}

	// 设置Mock期望
	suite.mockGPUClusterModel.EXPECT().
		FindOne(clusterId).
		Return(existingCluster, nil).
		Times(1)

	// 检查是否有活跃的设备
	suite.mockGPUDeviceModel.EXPECT().
		CountActiveDevices(clusterId).
		Return(0, nil).
		Times(1)

	suite.mockGPUClusterModel.EXPECT().
		Delete(clusterId).
		Return(nil).
		Times(1)

	// 创建删除集群逻辑
	deleteLogic := gpu_cluster.NewDeleteGpuClusterLogic(suite.ctx, suite.mockSvcCtx)

	// 执行测试
	err := deleteLogic.DeleteGpuCluster(&types.DeleteGpuClusterReq{ClusterId: clusterId})

	// 验证结果
	assert.NoError(suite.T(), err)
}

// TestDeleteGPUClusterWithActiveDevices 测试删除有活跃设备的集群
func (suite *GPUManagementLogicTestSuite) TestDeleteGPUClusterWithActiveDevices() {
	clusterId := int64(2001)

	existingCluster := &model.VtGpuClusters{
		Id:     clusterId,
		Name:   "test-cluster",
		Status: "active",
	}

	// 设置Mock期望
	suite.mockGPUClusterModel.EXPECT().
		FindOne(clusterId).
		Return(existingCluster, nil).
		Times(1)

	// 检查有活跃的设备
	suite.mockGPUDeviceModel.EXPECT().
		CountActiveDevices(clusterId).
		Return(5, nil).
		Times(1)

	// 创建删除集群逻辑
	deleteLogic := gpu_cluster.NewDeleteGpuClusterLogic(suite.ctx, suite.mockSvcCtx)

	// 执行测试
	err := deleteLogic.DeleteGpuCluster(&types.DeleteGpuClusterReq{ClusterId: clusterId})

	// 验证结果
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "集群中还有活跃的GPU设备")
}

// 运行GPU管理逻辑测试套件
func TestGPUManagementLogicSuite(t *testing.T) {
	suite.Run(t, new(GPUManagementLogicTestSuite))
}
