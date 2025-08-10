package svc

import (
	"database/sql"
	"log"

	"api/internal/config"
	"api/model"
	"api/pkg/auth"
	"api/pkg/database"

	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config    config.Config
	DB        *sql.DB
	DBManager *database.MySQLManager
	Redis     *redis.Client
	Cache     *database.RedisCache

	// 认证相关服务
	JWTService     *auth.JWTService
	TokenBlacklist *auth.RedisTokenBlacklist

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
	VtMonitorDataModel           model.VtMonitorDataModel
	VtMonitorMetricsModel        model.VtMonitorMetricsModel
	VtAlertRecordsModel          model.VtAlertRecordsModel
	VtAlertRulesModel            model.VtAlertRulesModel
	VtNotificationChannelsModel  model.VtNotificationChannelsModel
	VtNotificationTemplatesModel model.VtNotificationTemplatesModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化MySQL连接
	db, err := database.NewMySQLConnection(c.MySQL)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// 初始化MySQL管理器
	dbManager, err := database.NewMySQLManager(c.MySQL)
	if err != nil {
		log.Fatalf("Failed to create MySQL manager: %v", err)
	}

	// 初始化Redis连接
	rdb, err := database.NewRedisClient(c.Redis)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		rdb = nil // Redis连接失败时设为nil，应用仍可运行
	}

	// 创建Redis缓存（如果Redis可用的话）
	var cache *database.RedisCache
	var tokenBlacklist *auth.RedisTokenBlacklist
	if rdb != nil {
		cache = database.NewRedisCache(rdb)
		tokenBlacklist = auth.NewRedisTokenBlacklist(rdb)
	}

	// 初始化JWT服务（go-zero 项目下本仓库提供的实现需要 4 个参数）
	jwtService := auth.NewJWTService(c.Auth.AccessSecret, c.Auth.RefreshSecret, c.Auth.AccessExpire, c.Auth.RefreshExpire)

	return &ServiceContext{
		Config:         c,
		DB:             db,
		DBManager:      dbManager,
		Redis:          rdb,
		Cache:          cache,
		JWTService:     jwtService,
		TokenBlacklist: tokenBlacklist,

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
