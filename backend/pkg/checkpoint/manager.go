package checkpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// CheckpointManager 检查点管理器
type CheckpointManager struct {
	logger  logx.Logger
	baseDir string
}

// CheckpointInfo 检查点信息
type CheckpointInfo struct {
	ID          string            `json:"id"`
	JobID       int64             `json:"job_id"`
	JobName     string            `json:"job_name"`
	Epoch       int               `json:"epoch"`
	Step        int               `json:"step"`
	Loss        float64           `json:"loss"`
	Accuracy    float64           `json:"accuracy"`
	ModelPath   string            `json:"model_path"`
	Size        int64             `json:"size"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
	Description string            `json:"description"`
	IsAutoSaved bool              `json:"is_auto_saved"`
}

// CheckpointConfig 检查点配置
type CheckpointConfig struct {
	SaveInterval    time.Duration `json:"save_interval"`     // 自动保存间隔
	MaxCheckpoints  int           `json:"max_checkpoints"`   // 最大保留数量
	SaveOnEpochEnd  bool          `json:"save_on_epoch_end"` // 是否在epoch结束时保存
	SaveBestOnly    bool          `json:"save_best_only"`    // 是否只保存最优模型
	MonitorMetric   string        `json:"monitor_metric"`    // 监控的指标
	MinDelta        float64       `json:"min_delta"`         // 最小改进阈值
	Patience        int           `json:"patience"`          // 早停耐心值
	RestoreFromBest bool          `json:"restore_from_best"` // 是否从最优检查点恢复
}

// CheckpointStorage 检查点存储接口
type CheckpointStorage interface {
	Save(info *CheckpointInfo, modelData []byte) error
	Load(checkpointID string) (*CheckpointInfo, []byte, error)
	List(jobID int64) ([]*CheckpointInfo, error)
	Delete(checkpointID string) error
	GetBest(jobID int64, metric string) (*CheckpointInfo, error)
	Cleanup(jobID int64, keepCount int) error
}

// NewCheckpointManager 创建检查点管理器
func NewCheckpointManager(baseDir string) *CheckpointManager {
	return &CheckpointManager{
		logger:  logx.WithContext(context.Background()),
		baseDir: baseDir,
	}
}

// SaveCheckpoint 保存检查点
func (cm *CheckpointManager) SaveCheckpoint(info *CheckpointInfo, config *CheckpointConfig, storage CheckpointStorage) error {
	// 生成检查点ID
	info.ID = cm.generateCheckpointID(info.JobID, info.Epoch, info.Step)
	info.CreatedAt = time.Now()

	// 构建模型路径
	info.ModelPath = cm.buildModelPath(info.JobID, info.ID)

	// 模拟模型数据（实际应该是训练框架保存的模型文件）
	modelData := cm.generateMockModelData(info)

	// 检查是否需要保存
	if config.SaveBestOnly {
		shouldSave, err := cm.shouldSaveCheckpoint(info, config, storage)
		if err != nil {
			return fmt.Errorf("检查是否应该保存检查点失败: %v", err)
		}
		if !shouldSave {
			cm.logger.Infof("跳过保存检查点 %s，指标未改进", info.ID)
			return nil
		}
	}

	// 保存检查点
	if err := storage.Save(info, modelData); err != nil {
		return fmt.Errorf("保存检查点失败: %v", err)
	}

	cm.logger.Infof("成功保存检查点: %s (epoch=%d, step=%d, loss=%.4f)",
		info.ID, info.Epoch, info.Step, info.Loss)

	// 清理旧检查点
	if config.MaxCheckpoints > 0 {
		if err := storage.Cleanup(info.JobID, config.MaxCheckpoints); err != nil {
			cm.logger.Errorf("清理旧检查点失败: %v", err)
		}
	}

	return nil
}

// RestoreCheckpoint 恢复检查点
func (cm *CheckpointManager) RestoreCheckpoint(checkpointID string, storage CheckpointStorage) (*CheckpointInfo, error) {
	info, modelData, err := storage.Load(checkpointID)
	if err != nil {
		return nil, fmt.Errorf("加载检查点失败: %v", err)
	}

	// 这里应该将模型数据加载到训练框架中
	cm.logger.Infof("模拟恢复检查点 %s，模型数据大小: %d bytes", checkpointID, len(modelData))

	cm.logger.Infof("成功恢复检查点: %s (epoch=%d, step=%d)",
		info.ID, info.Epoch, info.Step)

	return info, nil
}

// GetCheckpointHistory 获取检查点历史
func (cm *CheckpointManager) GetCheckpointHistory(jobID int64, storage CheckpointStorage) ([]*CheckpointInfo, error) {
	checkpoints, err := storage.List(jobID)
	if err != nil {
		return nil, fmt.Errorf("获取检查点列表失败: %v", err)
	}

	// 按创建时间排序
	sort.Slice(checkpoints, func(i, j int) bool {
		return checkpoints[i].CreatedAt.After(checkpoints[j].CreatedAt)
	})

	return checkpoints, nil
}

// GetBestCheckpoint 获取最优检查点
func (cm *CheckpointManager) GetBestCheckpoint(jobID int64, metric string, storage CheckpointStorage) (*CheckpointInfo, error) {
	return storage.GetBest(jobID, metric)
}

// DeleteCheckpoint 删除检查点
func (cm *CheckpointManager) DeleteCheckpoint(checkpointID string, storage CheckpointStorage) error {
	if err := storage.Delete(checkpointID); err != nil {
		return fmt.Errorf("删除检查点失败: %v", err)
	}

	cm.logger.Infof("成功删除检查点: %s", checkpointID)
	return nil
}

// shouldSaveCheckpoint 检查是否应该保存检查点
func (cm *CheckpointManager) shouldSaveCheckpoint(info *CheckpointInfo, config *CheckpointConfig, storage CheckpointStorage) (bool, error) {
	if config.MonitorMetric == "" {
		return true, nil
	}

	// 获取当前最优检查点
	bestCheckpoint, err := storage.GetBest(info.JobID, config.MonitorMetric)
	if err != nil || bestCheckpoint == nil {
		// 没有历史检查点，保存当前的
		return true, nil
	}

	// 比较指标
	currentValue := cm.getMetricValue(info, config.MonitorMetric)
	bestValue := cm.getMetricValue(bestCheckpoint, config.MonitorMetric)

	// 检查是否有改进
	var improved bool
	switch config.MonitorMetric {
	case "loss":
		improved = currentValue < bestValue-config.MinDelta
	case "accuracy":
		improved = currentValue > bestValue+config.MinDelta
	default:
		improved = currentValue > bestValue+config.MinDelta
	}

	return improved, nil
}

// getMetricValue 获取指标值
func (cm *CheckpointManager) getMetricValue(info *CheckpointInfo, metric string) float64 {
	switch metric {
	case "loss":
		return info.Loss
	case "accuracy":
		return info.Accuracy
	default:
		return 0.0
	}
}

// generateCheckpointID 生成检查点ID
func (cm *CheckpointManager) generateCheckpointID(jobID int64, epoch, step int) string {
	timestamp := time.Now().Format("20060102-150405")
	return fmt.Sprintf("checkpoint-%d-epoch%d-step%d-%s", jobID, epoch, step, timestamp)
}

// buildModelPath 构建模型路径
func (cm *CheckpointManager) buildModelPath(jobID int64, checkpointID string) string {
	return filepath.Join(cm.baseDir, fmt.Sprintf("job-%d", jobID), checkpointID, "model.pkl")
}

// generateMockModelData 生成模拟模型数据
func (cm *CheckpointManager) generateMockModelData(info *CheckpointInfo) []byte {
	// 模拟模型数据结构
	mockModel := map[string]interface{}{
		"architecture": "ResNet50",
		"parameters": map[string]interface{}{
			"layers":  50,
			"filters": []int{64, 128, 256, 512},
		},
		"training_info": map[string]interface{}{
			"epoch":    info.Epoch,
			"step":     info.Step,
			"loss":     info.Loss,
			"accuracy": info.Accuracy,
		},
		"metadata":   info.Metadata,
		"created_at": info.CreatedAt.Format(time.RFC3339),
	}

	data, _ := json.Marshal(mockModel)
	return data
}

// FileSystemCheckpointStorage 文件系统检查点存储实现
type FileSystemCheckpointStorage struct {
	baseDir string
	logger  logx.Logger
}

// NewFileSystemCheckpointStorage 创建文件系统检查点存储
func NewFileSystemCheckpointStorage(baseDir string) *FileSystemCheckpointStorage {
	return &FileSystemCheckpointStorage{
		baseDir: baseDir,
		logger:  logx.WithContext(context.Background()),
	}
}

// Save 保存检查点
func (fs *FileSystemCheckpointStorage) Save(info *CheckpointInfo, modelData []byte) error {
	// 创建目录
	checkpointDir := filepath.Join(fs.baseDir, fmt.Sprintf("job-%d", info.JobID), info.ID)
	if err := os.MkdirAll(checkpointDir, 0755); err != nil {
		return fmt.Errorf("创建检查点目录失败: %v", err)
	}

	// 保存模型数据
	modelPath := filepath.Join(checkpointDir, "model.pkl")
	if err := os.WriteFile(modelPath, modelData, 0644); err != nil {
		return fmt.Errorf("保存模型文件失败: %v", err)
	}

	// 保存检查点信息
	infoPath := filepath.Join(checkpointDir, "checkpoint.json")
	infoData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化检查点信息失败: %v", err)
	}

	if err := os.WriteFile(infoPath, infoData, 0644); err != nil {
		return fmt.Errorf("保存检查点信息失败: %v", err)
	}

	// 更新文件大小
	if stat, err := os.Stat(modelPath); err == nil {
		info.Size = stat.Size()
	}

	fs.logger.Infof("检查点已保存到: %s", checkpointDir)
	return nil
}

// Load 加载检查点
func (fs *FileSystemCheckpointStorage) Load(checkpointID string) (*CheckpointInfo, []byte, error) {
	// 查找检查点目录
	checkpointDir, err := fs.findCheckpointDir(checkpointID)
	if err != nil {
		return nil, nil, err
	}

	// 加载检查点信息
	infoPath := filepath.Join(checkpointDir, "checkpoint.json")
	infoData, err := os.ReadFile(infoPath)
	if err != nil {
		return nil, nil, fmt.Errorf("读取检查点信息失败: %v", err)
	}

	var info CheckpointInfo
	if err := json.Unmarshal(infoData, &info); err != nil {
		return nil, nil, fmt.Errorf("解析检查点信息失败: %v", err)
	}

	// 加载模型数据
	modelPath := filepath.Join(checkpointDir, "model.pkl")
	modelData, err := os.ReadFile(modelPath)
	if err != nil {
		return nil, nil, fmt.Errorf("读取模型文件失败: %v", err)
	}

	return &info, modelData, nil
}

// List 列出检查点
func (fs *FileSystemCheckpointStorage) List(jobID int64) ([]*CheckpointInfo, error) {
	jobDir := filepath.Join(fs.baseDir, fmt.Sprintf("job-%d", jobID))

	entries, err := os.ReadDir(jobDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*CheckpointInfo{}, nil
		}
		return nil, fmt.Errorf("读取作业目录失败: %v", err)
	}

	var checkpoints []*CheckpointInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		infoPath := filepath.Join(jobDir, entry.Name(), "checkpoint.json")
		infoData, err := os.ReadFile(infoPath)
		if err != nil {
			fs.logger.Errorf("读取检查点信息失败: %s, %v", infoPath, err)
			continue
		}

		var info CheckpointInfo
		if err := json.Unmarshal(infoData, &info); err != nil {
			fs.logger.Errorf("解析检查点信息失败: %s, %v", infoPath, err)
			continue
		}

		checkpoints = append(checkpoints, &info)
	}

	return checkpoints, nil
}

// Delete 删除检查点
func (fs *FileSystemCheckpointStorage) Delete(checkpointID string) error {
	checkpointDir, err := fs.findCheckpointDir(checkpointID)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(checkpointDir); err != nil {
		return fmt.Errorf("删除检查点目录失败: %v", err)
	}

	return nil
}

// GetBest 获取最优检查点
func (fs *FileSystemCheckpointStorage) GetBest(jobID int64, metric string) (*CheckpointInfo, error) {
	checkpoints, err := fs.List(jobID)
	if err != nil {
		return nil, err
	}

	if len(checkpoints) == 0 {
		return nil, nil
	}

	var best *CheckpointInfo
	for _, checkpoint := range checkpoints {
		if best == nil {
			best = checkpoint
			continue
		}

		switch metric {
		case "loss":
			if checkpoint.Loss < best.Loss {
				best = checkpoint
			}
		case "accuracy":
			if checkpoint.Accuracy > best.Accuracy {
				best = checkpoint
			}
		}
	}

	return best, nil
}

// Cleanup 清理旧检查点
func (fs *FileSystemCheckpointStorage) Cleanup(jobID int64, keepCount int) error {
	checkpoints, err := fs.List(jobID)
	if err != nil {
		return err
	}

	if len(checkpoints) <= keepCount {
		return nil
	}

	// 按创建时间排序，保留最新的
	sort.Slice(checkpoints, func(i, j int) bool {
		return checkpoints[i].CreatedAt.After(checkpoints[j].CreatedAt)
	})

	// 删除多余的检查点
	for i := keepCount; i < len(checkpoints); i++ {
		if err := fs.Delete(checkpoints[i].ID); err != nil {
			fs.logger.Errorf("清理检查点失败: %s, %v", checkpoints[i].ID, err)
		} else {
			fs.logger.Infof("已清理检查点: %s", checkpoints[i].ID)
		}
	}

	return nil
}

// findCheckpointDir 查找检查点目录
func (fs *FileSystemCheckpointStorage) findCheckpointDir(checkpointID string) (string, error) {
	// 遍历所有作业目录查找检查点
	entries, err := os.ReadDir(fs.baseDir)
	if err != nil {
		return "", fmt.Errorf("读取基础目录失败: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "job-") {
			continue
		}

		checkpointDir := filepath.Join(fs.baseDir, entry.Name(), checkpointID)
		if _, err := os.Stat(checkpointDir); err == nil {
			return checkpointDir, nil
		}
	}

	return "", fmt.Errorf("检查点不存在: %s", checkpointID)
}

// strings.HasPrefix 的简单实现（避免导入strings包）
func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}
