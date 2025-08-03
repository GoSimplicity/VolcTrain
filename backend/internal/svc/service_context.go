package svc

import (
	"database/sql"
	"log"

	"api/internal/config"
	"api/model"
	"api/pkg/database"
)

type ServiceContext struct {
	Config config.Config
	DB     *sql.DB
	// Redis              *redis.Client
	// Cache              *database.RedisCache
	VtUsersModel          model.VtUsersModel
	VtTrainingQueuesModel model.VtTrainingQueuesModel
	VtTrainingJobsModel   model.VtTrainingJobsModel
	VtGpuClustersModel    model.VtGpuClustersModel
	VtGpuNodesModel       model.VtGpuNodesModel
	VtGpuDevicesModel     model.VtGpuDevicesModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化MySQL连接
	db, err := database.NewMySQLConnection(c.MySQL)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// 创建一个简单的Redis缓存（如果需要的话）
	// cache := database.NewRedisCache(rdb)

	return &ServiceContext{
		Config: c,
		DB:     db,
		// Redis:              rdb,
		// Cache:              cache,
		VtUsersModel:          model.NewVtUsersModel(db),
		VtTrainingQueuesModel: model.NewVtTrainingQueuesModel(db),
		VtTrainingJobsModel:   model.NewVtTrainingJobsModel(db),
		VtGpuClustersModel:    model.NewVtGpuClustersModel(db),
		VtGpuNodesModel:       model.NewVtGpuNodesModel(db),
		VtGpuDevicesModel:     model.NewVtGpuDevicesModel(db),
	}
}
