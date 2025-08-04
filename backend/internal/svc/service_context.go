package svc

import (
	"database/sql"
	"log"

	"api/internal/config"
	"api/model"
	"api/pkg/database"
	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config config.Config
	DB     *sql.DB
	Redis  *redis.Client
	Cache  *database.RedisCache

	// 用户相关模型
	VtUsersModel       model.VtUsersSimpleModel
	VtRolesModel       model.VtRolesModel
	VtPermissionsModel model.VtPermissionsModel

	// 训练相关模型
	VtTrainingQueuesModel model.VtTrainingQueuesModel
	VtTrainingJobsModel   model.VtTrainingJobsModel

	// GPU相关模型
	VtGpuClustersModel model.VtGpuClustersModel
	VtGpuNodesModel    model.VtGpuNodesModel
	VtGpuDevicesModel  model.VtGpuDevicesModel

	// 监控相关模型
	VtMonitorDataModel         model.VtMonitorDataModel
	VtMonitorMetricsModel      model.VtMonitorMetricsModel
	VtAlertRecordsModel        model.VtAlertRecordsModel
	VtAlertRulesModel          model.VtAlertRulesModel
	VtNotificationChannelsModel model.VtNotificationChannelsModel
	VtNotificationTemplatesModel model.VtNotificationTemplatesModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化MySQL连接
	db, err := database.NewMySQLConnection(c.MySQL)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// 初始化Redis连接
	rdb, err := database.NewRedisClient(c.Redis)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		rdb = nil // Redis连接失败时设为nil，应用仍可运行
	}

	// 创建Redis缓存（如果Redis可用的话）
	var cache *database.RedisCache
	if rdb != nil {
		cache = database.NewRedisCache(rdb)
	}

	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  rdb,
		Cache:  cache,

		// 初始化所有模型
		VtUsersModel:       model.NewVtUsersSimpleModel(db),
		VtRolesModel:       model.NewVtRolesModel(db),
		VtPermissionsModel: model.NewVtPermissionsModel(db),

		VtTrainingQueuesModel: model.NewVtTrainingQueuesModel(db),
		VtTrainingJobsModel:   model.NewVtTrainingJobsModel(db),

		VtGpuClustersModel: model.NewVtGpuClustersModel(db),
		VtGpuNodesModel:    model.NewVtGpuNodesModel(db),
		VtGpuDevicesModel:  model.NewVtGpuDevicesModel(db),

		VtMonitorDataModel:           model.NewVtMonitorDataModel(db),
		VtMonitorMetricsModel:        model.NewVtMonitorMetricsModel(db),
		VtAlertRecordsModel:          model.NewVtAlertRecordsModel(db),
		VtAlertRulesModel:            model.NewVtAlertRulesModel(db),
		VtNotificationChannelsModel:  model.NewVtNotificationChannelsModel(db),
		VtNotificationTemplatesModel: model.NewVtNotificationTemplatesModel(db),
	}
}
