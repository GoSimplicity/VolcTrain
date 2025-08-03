package seeds

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// TestData 测试数据结构
type TestData struct {
	Users         []UserSeed         `json:"users"`
	GPUClusters   []GPUClusterSeed   `json:"gpu_clusters"`
	TrainingJobs  []TrainingJobSeed  `json:"training_jobs"`
	Notifications []NotificationSeed `json:"notifications"`
}

type UserSeed struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Salt         string `json:"salt"`
	RealName     string `json:"real_name"`
	UserType     string `json:"user_type"`
	Status       string `json:"status"`
}

type GPUClusterSeed struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Kubeconfig  string `json:"kubeconfig"`
	Status      string `json:"status"`
	NodeCount   int    `json:"node_count"`
}

type TrainingJobSeed struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Name      string `json:"name"`
	Framework string `json:"framework"`
	Image     string `json:"image"`
	Status    string `json:"status"`
	CPUCores  string `json:"cpu_cores"`
	MemoryGb  string `json:"memory_gb"`
	GPUCount  int    `json:"gpu_count"`
}

type NotificationSeed struct {
	ID          int64  `json:"id"`
	ChannelType string `json:"channel_type"`
	Name        string `json:"name"`
	Config      string `json:"config"`
	Enabled     bool   `json:"enabled"`
}

// DataSeeder 数据种子生成器
type DataSeeder struct {
	db *sql.DB
}

// NewDataSeeder 创建数据种子生成器
func NewDataSeeder(db *sql.DB) *DataSeeder {
	return &DataSeeder{db: db}
}

// SeedAllTestData 生成所有测试数据
func (s *DataSeeder) SeedAllTestData() error {
	// 清理现有数据
	if err := s.CleanupAllData(); err != nil {
		return fmt.Errorf("清理数据失败: %w", err)
	}

	// 生成基础测试数据
	testData := s.generateTestData()

	// 插入用户数据
	if err := s.seedUsers(testData.Users); err != nil {
		return fmt.Errorf("插入用户数据失败: %w", err)
	}

	// 插入GPU集群数据
	if err := s.seedGPUClusters(testData.GPUClusters); err != nil {
		return fmt.Errorf("插入GPU集群数据失败: %w", err)
	}

	// 插入训练作业数据
	if err := s.seedTrainingJobs(testData.TrainingJobs); err != nil {
		return fmt.Errorf("插入训练作业数据失败: %w", err)
	}

	// 插入通知渠道数据
	if err := s.seedNotifications(testData.Notifications); err != nil {
		return fmt.Errorf("插入通知渠道数据失败: %w", err)
	}

	return nil
}

// generateTestData 生成测试数据
func (s *DataSeeder) generateTestData() *TestData {
	return &TestData{
		Users: []UserSeed{
			{
				ID:           1001,
				Username:     "admin_test",
				Email:        "admin@test.com",
				PasswordHash: "hashed_admin_password",
				Salt:         "admin_salt",
				RealName:     "测试管理员",
				UserType:     "admin",
				Status:       "active",
			},
			{
				ID:           1002,
				Username:     "user_test",
				Email:        "user@test.com",
				PasswordHash: "hashed_user_password",
				Salt:         "user_salt",
				RealName:     "测试用户",
				UserType:     "user",
				Status:       "active",
			},
			{
				ID:           1003,
				Username:     "data_scientist",
				Email:        "ds@test.com",
				PasswordHash: "hashed_ds_password",
				Salt:         "ds_salt",
				RealName:     "数据科学家",
				UserType:     "user",
				Status:       "active",
			},
		},
		GPUClusters: []GPUClusterSeed{
			{
				ID:          2001,
				Name:        "test-cluster-1",
				Description: "测试GPU集群1",
				Kubeconfig:  "test_kubeconfig_1",
				Status:      "active",
				NodeCount:   3,
			},
			{
				ID:          2002,
				Name:        "test-cluster-2",
				Description: "测试GPU集群2",
				Kubeconfig:  "test_kubeconfig_2",
				Status:      "active",
				NodeCount:   5,
			},
		},
		TrainingJobs: []TrainingJobSeed{
			{
				ID:        3001,
				UserID:    1002,
				Name:      "pytorch-training-test",
				Framework: "pytorch",
				Image:     "pytorch/pytorch:1.12.0",
				Status:    "running",
				CPUCores:  "4",
				MemoryGb:  "8",
				GPUCount:  1,
			},
			{
				ID:        3002,
				UserID:    1003,
				Name:      "tensorflow-training-test",
				Framework: "tensorflow",
				Image:     "tensorflow/tensorflow:2.10.0",
				Status:    "pending",
				CPUCores:  "8",
				MemoryGb:  "16",
				GPUCount:  2,
			},
		},
		Notifications: []NotificationSeed{
			{
				ID:          4001,
				ChannelType: "email",
				Name:        "测试邮件通知",
				Config:      `{"smtp_host":"localhost","smtp_port":587,"username":"test","password":"test","from":"test@example.com"}`,
				Enabled:     true,
			},
			{
				ID:          4002,
				ChannelType: "webhook",
				Name:        "测试Webhook通知",
				Config:      `{"url":"http://localhost:3000/webhook","secret":"test_secret"}`,
				Enabled:     true,
			},
		},
	}
}

// seedUsers 插入用户数据
func (s *DataSeeder) seedUsers(users []UserSeed) error {
	query := `INSERT INTO vt_users (id, username, email, password_hash, salt, real_name, user_type, status, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	for _, user := range users {
		_, err := s.db.Exec(query,
			user.ID, user.Username, user.Email, user.PasswordHash,
			user.Salt, user.RealName, user.UserType, user.Status)
		if err != nil {
			return fmt.Errorf("插入用户 %s 失败: %w", user.Username, err)
		}
	}

	return nil
}

// seedGPUClusters 插入GPU集群数据
func (s *DataSeeder) seedGPUClusters(clusters []GPUClusterSeed) error {
	query := `INSERT INTO vt_gpu_clusters (id, name, description, kubeconfig, status, node_count, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`

	for _, cluster := range clusters {
		_, err := s.db.Exec(query,
			cluster.ID, cluster.Name, cluster.Description,
			cluster.Kubeconfig, cluster.Status, cluster.NodeCount)
		if err != nil {
			return fmt.Errorf("插入GPU集群 %s 失败: %w", cluster.Name, err)
		}
	}

	return nil
}

// seedTrainingJobs 插入训练作业数据
func (s *DataSeeder) seedTrainingJobs(jobs []TrainingJobSeed) error {
	query := `INSERT INTO vt_training_jobs (id, user_id, name, framework, image, status, cpu_cores, memory_gb, gpu_count, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	for _, job := range jobs {
		_, err := s.db.Exec(query,
			job.ID, job.UserID, job.Name, job.Framework, job.Image,
			job.Status, job.CPUCores, job.MemoryGb, job.GPUCount)
		if err != nil {
			return fmt.Errorf("插入训练作业 %s 失败: %w", job.Name, err)
		}
	}

	return nil
}

// seedNotifications 插入通知渠道数据
func (s *DataSeeder) seedNotifications(notifications []NotificationSeed) error {
	query := `INSERT INTO vt_notification_channels (id, channel_type, name, config, enabled, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, NOW(), NOW())`

	for _, notification := range notifications {
		_, err := s.db.Exec(query,
			notification.ID, notification.ChannelType, notification.Name,
			notification.Config, notification.Enabled)
		if err != nil {
			return fmt.Errorf("插入通知渠道 %s 失败: %w", notification.Name, err)
		}
	}

	return nil
}

// CleanupAllData 清理所有测试数据
func (s *DataSeeder) CleanupAllData() error {
	// 按依赖关系顺序清理表
	tables := []string{
		"vt_training_job_instances",
		"vt_training_jobs",
		"vt_training_queues",
		"vt_gpu_devices",
		"vt_gpu_nodes",
		"vt_gpu_clusters",
		"vt_alert_records",
		"vt_notification_channels",
		"vt_users",
	}

	for _, table := range tables {
		_, err := s.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id >= 1000", table))
		if err != nil {
			return fmt.Errorf("清理表 %s 失败: %w", table, err)
		}
	}

	return nil
}

// LoadFromJSON 从JSON文件加载测试数据
func (s *DataSeeder) LoadFromJSON(jsonData []byte) error {
	var testData TestData
	if err := json.Unmarshal(jsonData, &testData); err != nil {
		return fmt.Errorf("解析JSON数据失败: %w", err)
	}

	// 清理现有数据
	if err := s.CleanupAllData(); err != nil {
		return err
	}

	// 插入数据
	if err := s.seedUsers(testData.Users); err != nil {
		return err
	}

	if err := s.seedGPUClusters(testData.GPUClusters); err != nil {
		return err
	}

	if err := s.seedTrainingJobs(testData.TrainingJobs); err != nil {
		return err
	}

	if err := s.seedNotifications(testData.Notifications); err != nil {
		return err
	}

	return nil
}

// CreateTestWorkspace 创建测试工作空间
func (s *DataSeeder) CreateTestWorkspace(ownerID int64, name string) (int64, error) {
	query := `INSERT INTO vt_workspaces (owner_id, name, display_name, description, visibility, status, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, 'private', 'active', NOW(), NOW())`

	result, err := s.db.Exec(query, ownerID, name, name+"_display", "测试工作空间")
	if err != nil {
		return 0, fmt.Errorf("创建测试工作空间失败: %w", err)
	}

	workspaceID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取工作空间ID失败: %w", err)
	}

	return workspaceID, nil
}

// GetTestDataStats 获取测试数据统计信息
func (s *DataSeeder) GetTestDataStats() (map[string]int, error) {
	stats := make(map[string]int)

	tables := []string{
		"vt_users", "vt_gpu_clusters", "vt_training_jobs",
		"vt_notification_channels", "vt_workspaces",
	}

	for _, table := range tables {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id >= 1000", table)
		err := s.db.QueryRow(query).Scan(&count)
		if err != nil {
			return nil, fmt.Errorf("获取表 %s 统计信息失败: %w", table, err)
		}
		stats[table] = count
	}

	return stats, nil
}
