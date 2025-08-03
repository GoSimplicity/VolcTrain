package mocks

//go:generate mockgen -source=../../model/vt_users_model.go -destination=user_model_mock.go -package=mocks
//go:generate mockgen -source=../../model/vt_training_jobs_model.go -destination=training_job_model_mock.go -package=mocks
//go:generate mockgen -source=../../model/vt_gpu_clusters_model.go -destination=gpu_cluster_model_mock.go -package=mocks
//go:generate mockgen -source=../../model/vt_gpu_devices_model.go -destination=gpu_device_model_mock.go -package=mocks
//go:generate mockgen -source=../../pkg/auth/jwt.go -destination=jwt_mock.go -package=mocks
//go:generate mockgen -source=../../pkg/database/mysql.go -destination=database_mock.go -package=mocks
//go:generate mockgen -source=../../pkg/database/redis.go -destination=redis_mock.go -package=mocks

// 这个文件用于生成Mock对象
// 运行命令: go generate ./test/mocks
