package alerting

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"sync"
	"time"

	"api/model"
	"github.com/zeromicro/go-zero/core/logx"
)

// AlertEngine 告警规则引擎
type AlertEngine struct {
	logger              logx.Logger
	ctx                 context.Context
	cancel              context.CancelFunc
	db                  *sql.DB
	alertRulesModel     model.VtAlertRulesModel
	alertRecordsModel   model.VtAlertRecordsModel
	monitorDataModel    model.VtMonitorDataModel
	metricsModel        model.VtMonitorMetricsModel
	config              *AlertEngineConfig
	rules               map[int64]*AlertRule
	activeAlerts        map[string]*ActiveAlert
	evaluationInterval  time.Duration
	lastEvaluationTime  time.Time
	mu                  sync.RWMutex
	notificationChannel chan *AlertNotification
	stopped             bool
}

// AlertEngineConfig 告警引擎配置
type AlertEngineConfig struct {
	EvaluationInterval      time.Duration `json:"evaluation_interval"`
	MaxConcurrentRules      int           `json:"max_concurrent_rules"`
	AlertRetentionDays      int           `json:"alert_retention_days"`
	EnableGrouping          bool          `json:"enable_grouping"`
	EnableSuppression       bool          `json:"enable_suppression"`
	DefaultThrottleMinutes  int           `json:"default_throttle_minutes"`
	AnomalyDetectionEnabled bool          `json:"anomaly_detection_enabled"`
}

// AlertRule 告警规则
type AlertRule struct {
	ID                          int64                  `json:"id"`
	Name                        string                 `json:"name"`
	DisplayName                 string                 `json:"display_name"`
	Description                 string                 `json:"description"`
	RuleType                    string                 `json:"rule_type"`
	ConditionExpression         string                 `json:"condition_expression"`
	QueryExpression             string                 `json:"query_expression"`
	WarningThreshold            float64                `json:"warning_threshold"`
	CriticalThreshold           float64                `json:"critical_threshold"`
	ThresholdCondition          string                 `json:"threshold_condition"`
	EvaluationWindowSeconds     int                    `json:"evaluation_window_seconds"`
	EvaluationIntervalSeconds   int                    `json:"evaluation_interval_seconds"`
	TriggerDurationSeconds      int                    `json:"trigger_duration_seconds"`
	RecoveryDurationSeconds     int                    `json:"recovery_duration_seconds"`
	FilterLabels                map[string]interface{} `json:"filter_labels"`
	FilterResources             map[string]interface{} `json:"filter_resources"`
	AlertLevel                  string                 `json:"alert_level"`
	SeverityScore               int                    `json:"severity_score"`
	NotificationChannels        []string               `json:"notification_channels"`
	NotificationThrottleMinutes int                    `json:"notification_throttle_minutes"`
	SilenceDurationSeconds      int                    `json:"silence_duration_seconds"`
	SuppressionRules            map[string]interface{} `json:"suppression_rules"`
	DependencyRules             map[string]interface{} `json:"dependency_rules"`
	Status                      string                 `json:"status"`
	LastEvaluationAt            time.Time              `json:"last_evaluation_at"`
	LastTriggerAt               time.Time              `json:"last_trigger_at"`
	TriggerCount                int                    `json:"trigger_count"`
}

// ActiveAlert 活跃告警
type ActiveAlert struct {
	ID                  string                 `json:"id"`
	RuleID              int64                  `json:"rule_id"`
	RuleName            string                 `json:"rule_name"`
	AlertLevel          string                 `json:"alert_level"`
	SeverityScore       int                    `json:"severity_score"`
	Message             string                 `json:"message"`
	Summary             string                 `json:"summary"`
	ResourceType        string                 `json:"resource_type"`
	ResourceID          int64                  `json:"resource_id"`
	ResourceName        string                 `json:"resource_name"`
	InstanceID          string                 `json:"instance_id"`
	TriggerValue        float64                `json:"trigger_value"`
	ThresholdValue      float64                `json:"threshold_value"`
	ConditionExpression string                 `json:"condition_expression"`
	EvaluationData      map[string]interface{} `json:"evaluation_data"`
	Labels              map[string]interface{} `json:"labels"`
	Annotations         map[string]interface{} `json:"annotations"`
	Context             map[string]interface{} `json:"context"`
	Status              string                 `json:"status"`
	TriggeredAt         time.Time              `json:"triggered_at"`
	FirstOccurrenceAt   time.Time              `json:"first_occurrence_at"`
	LastOccurrenceAt    time.Time              `json:"last_occurrence_at"`
	OccurrenceCount     int                    `json:"occurrence_count"`
	NotificationSent    bool                   `json:"notification_sent"`
	NotificationCount   int                    `json:"notification_count"`
	LastNotificationAt  time.Time              `json:"last_notification_at"`
	EscalationLevel     int                    `json:"escalation_level"`
	AlertGroupID        string                 `json:"alert_group_id"`
	CorrelationID       string                 `json:"correlation_id"`
}

// AlertNotification 告警通知
type AlertNotification struct {
	Alert    *ActiveAlert `json:"alert"`
	RuleInfo *AlertRule   `json:"rule_info"`
	Action   string       `json:"action"` // firing, resolved, acknowledged
}

// EvaluationResult 评估结果
type EvaluationResult struct {
	RuleID      int64                  `json:"rule_id"`
	Triggered   bool                   `json:"triggered"`
	Value       float64                `json:"value"`
	Threshold   float64                `json:"threshold"`
	Message     string                 `json:"message"`
	EvaluatedAt time.Time              `json:"evaluated_at"`
	Labels      map[string]interface{} `json:"labels"`
	Context     map[string]interface{} `json:"context"`
}

// ThresholdCondition 阈值条件类型
type ThresholdCondition string

const (
	ConditionGT      ThresholdCondition = "gt"      // 大于
	ConditionGTE     ThresholdCondition = "gte"     // 大于等于
	ConditionLT      ThresholdCondition = "lt"      // 小于
	ConditionLTE     ThresholdCondition = "lte"     // 小于等于
	ConditionEQ      ThresholdCondition = "eq"      // 等于
	ConditionNEQ     ThresholdCondition = "neq"     // 不等于
	ConditionBetween ThresholdCondition = "between" // 在范围内
	ConditionOutside ThresholdCondition = "outside" // 在范围外
)

// AlertStatus 告警状态
type AlertStatus string

const (
	StatusFiring       AlertStatus = "firing"       // 触发中
	StatusResolved     AlertStatus = "resolved"     // 已解决
	StatusAcknowledged AlertStatus = "acknowledged" // 已确认
	StatusSuppressed   AlertStatus = "suppressed"   // 已抑制
	StatusSilenced     AlertStatus = "silenced"     // 已静默
)

// AlertLevel 告警级别
type AlertLevel string

const (
	LevelInfo     AlertLevel = "info"     // 信息
	LevelWarning  AlertLevel = "warning"  // 警告
	LevelCritical AlertLevel = "critical" // 严重
	LevelFatal    AlertLevel = "fatal"    // 致命
)

// NewAlertEngine 创建告警引擎
func NewAlertEngine(db *sql.DB, config *AlertEngineConfig) *AlertEngine {
	ctx, cancel := context.WithCancel(context.Background())

	engine := &AlertEngine{
		logger:              logx.WithContext(ctx),
		ctx:                 ctx,
		cancel:              cancel,
		db:                  db,
		alertRulesModel:     model.NewVtAlertRulesModel(db),
		alertRecordsModel:   model.NewVtAlertRecordsModel(db),
		monitorDataModel:    model.NewVtMonitorDataModel(db),
		metricsModel:        model.NewVtMonitorMetricsModel(db),
		config:              config,
		rules:               make(map[int64]*AlertRule),
		activeAlerts:        make(map[string]*ActiveAlert),
		evaluationInterval:  config.EvaluationInterval,
		notificationChannel: make(chan *AlertNotification, 1000),
	}

	return engine
}

// Start 启动告警引擎
func (e *AlertEngine) Start() error {
	e.logger.Infof("启动告警规则引擎，评估间隔: %v", e.evaluationInterval)

	// 加载告警规则
	if err := e.loadAlertRules(); err != nil {
		return fmt.Errorf("加载告警规则失败: %v", err)
	}

	// 启动评估循环
	go e.evaluationLoop()

	// 启动通知处理
	go e.notificationLoop()

	// 启动清理任务
	go e.cleanupLoop()

	return nil
}

// Stop 停止告警引擎
func (e *AlertEngine) Stop() {
	e.logger.Info("停止告警规则引擎")
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.stopped {
		return
	}

	if e.cancel != nil {
		e.cancel()
		e.cancel = nil
	}

	if e.notificationChannel != nil {
		close(e.notificationChannel)
		e.notificationChannel = nil
	}

	e.stopped = true
}

// loadAlertRules 加载告警规则
func (e *AlertEngine) loadAlertRules() error {
	e.logger.Info("加载告警规则")

	// 查询活跃的告警规则
	rules, err := e.alertRulesModel.FindActiveRules()
	if err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// 清空现有规则
	e.rules = make(map[int64]*AlertRule)

	// 加载规则
	for _, rule := range rules {
		alertRule := e.convertToAlertRule(rule)
		e.rules[alertRule.ID] = alertRule
		e.logger.Infof("加载告警规则: %s [ID:%d]", alertRule.Name, alertRule.ID)
	}

	e.logger.Infof("成功加载 %d 条告警规则", len(e.rules))
	return nil
}

// convertToAlertRule 转换为告警规则
func (e *AlertEngine) convertToAlertRule(dbRule *model.VtAlertRules) *AlertRule {
	rule := &AlertRule{
		ID:                          dbRule.Id,
		Name:                        dbRule.Name,
		DisplayName:                 dbRule.DisplayName,
		Description:                 dbRule.Description,
		RuleType:                    dbRule.RuleType,
		ConditionExpression:         dbRule.ConditionExpression,
		QueryExpression:             dbRule.QueryExpression,
		WarningThreshold:            dbRule.WarningThreshold,
		CriticalThreshold:           dbRule.CriticalThreshold,
		ThresholdCondition:          dbRule.ThresholdCondition,
		EvaluationWindowSeconds:     dbRule.EvaluationWindowSeconds,
		EvaluationIntervalSeconds:   dbRule.EvaluationIntervalSeconds,
		TriggerDurationSeconds:      dbRule.TriggerDurationSeconds,
		RecoveryDurationSeconds:     dbRule.RecoveryDurationSeconds,
		AlertLevel:                  dbRule.AlertLevel,
		SeverityScore:               dbRule.SeverityScore,
		NotificationThrottleMinutes: dbRule.NotificationThrottleMinutes,
		SilenceDurationSeconds:      dbRule.SilenceDurationSeconds,
		Status:                      dbRule.Status,
		LastEvaluationAt:            dbRule.LastEvaluationAt,
		LastTriggerAt:               dbRule.LastTriggerAt,
		TriggerCount:                dbRule.TriggerCount,
	}

	// 解析JSON字段
	if dbRule.FilterLabels != "" {
		json.Unmarshal([]byte(dbRule.FilterLabels), &rule.FilterLabels)
	}
	if dbRule.FilterResources != "" {
		json.Unmarshal([]byte(dbRule.FilterResources), &rule.FilterResources)
	}
	if dbRule.NotificationChannels != "" {
		json.Unmarshal([]byte(dbRule.NotificationChannels), &rule.NotificationChannels)
	}
	if dbRule.SuppressionRules != "" {
		json.Unmarshal([]byte(dbRule.SuppressionRules), &rule.SuppressionRules)
	}
	if dbRule.DependencyRules != "" {
		json.Unmarshal([]byte(dbRule.DependencyRules), &rule.DependencyRules)
	}

	return rule
}

// evaluationLoop 评估循环
func (e *AlertEngine) evaluationLoop() {
	ticker := time.NewTicker(e.evaluationInterval)
	defer ticker.Stop()

	for {
		select {
		case <-e.ctx.Done():
			e.logger.Info("告警评估循环停止")
			return
		case <-ticker.C:
			e.evaluateAllRules()
		}
	}
}

// evaluateAllRules 评估所有规则
func (e *AlertEngine) evaluateAllRules() {
	start := time.Now()
	e.logger.Debug("开始评估告警规则")

	e.mu.RLock()
	rules := make([]*AlertRule, 0, len(e.rules))
	for _, rule := range e.rules {
		rules = append(rules, rule)
	}
	e.mu.RUnlock()

	// 并发评估规则
	semaphore := make(chan struct{}, e.config.MaxConcurrentRules)
	var wg sync.WaitGroup

	for _, rule := range rules {
		wg.Add(1)
		go func(r *AlertRule) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			e.evaluateRule(r)
		}(rule)
	}

	wg.Wait()

	e.lastEvaluationTime = time.Now()
	duration := e.lastEvaluationTime.Sub(start)
	e.logger.Debugf("告警规则评估完成，耗时: %v", duration)
}

// evaluateRule 评估单个规则
func (e *AlertEngine) evaluateRule(rule *AlertRule) {
	// 检查规则是否应该被评估
	if !e.shouldEvaluateRule(rule) {
		return
	}

	// 获取指标数据
	metricData, err := e.getMetricDataForRule(rule)
	if err != nil {
		e.logger.Errorf("获取规则指标数据失败 [%s]: %v", rule.Name, err)
		return
	}

	// 评估规则条件
	result := e.evaluateRuleCondition(rule, metricData)

	// 更新规则评估时间
	e.updateRuleEvaluationTime(rule.ID)

	// 处理评估结果
	e.handleEvaluationResult(rule, result)
}

// shouldEvaluateRule 检查规则是否应该被评估
func (e *AlertEngine) shouldEvaluateRule(rule *AlertRule) bool {
	if rule.Status != "active" {
		return false
	}

	// 检查评估间隔
	if rule.EvaluationIntervalSeconds > 0 {
		interval := time.Duration(rule.EvaluationIntervalSeconds) * time.Second
		if time.Since(rule.LastEvaluationAt) < interval {
			return false
		}
	}

	return true
}

// getMetricDataForRule 获取规则的指标数据
func (e *AlertEngine) getMetricDataForRule(rule *AlertRule) ([]*model.VtMonitorData, error) {
	// 解析查询表达式或使用条件表达式
	query := rule.QueryExpression
	if query == "" {
		query = rule.ConditionExpression
	}

	// 简化实现：根据规则获取相关指标数据
	endTime := time.Now()
	startTime := endTime.Add(-time.Duration(rule.EvaluationWindowSeconds) * time.Second)

	// 这里简化处理，实际需要解析query获取metric_id
	// 假设从condition_expression中提取metric名称
	metricName := e.extractMetricNameFromExpression(query)
	if metricName == "" {
		return nil, fmt.Errorf("无法从表达式中提取指标名称: %s", query)
	}

	// 查找指标定义
	metric, err := e.metricsModel.FindOneByName(metricName)
	if err != nil {
		return nil, fmt.Errorf("查找指标定义失败: %v", err)
	}

	// 查询指标数据
	return e.monitorDataModel.FindByTimeRange(metric.Id, startTime, endTime)
}

// extractMetricNameFromExpression 从表达式中提取指标名称
func (e *AlertEngine) extractMetricNameFromExpression(expression string) string {
	// 简化实现：假设表达式格式为 "metric_name > threshold"
	// 实际需要更复杂的解析逻辑
	if len(expression) == 0 {
		return ""
	}

	// 简单的关键词匹配
	keywords := []string{"cpu_usage", "memory_usage", "gpu_usage", "training_job_count"}
	for _, keyword := range keywords {
		if containsIgnoreCase(expression, keyword) {
			return "system_" + keyword
		}
	}

	return ""
}

// containsIgnoreCase 不区分大小写的包含检查
func containsIgnoreCase(s, substr string) bool {
	// 简化实现
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

// evaluateRuleCondition 评估规则条件
func (e *AlertEngine) evaluateRuleCondition(rule *AlertRule, data []*model.VtMonitorData) *EvaluationResult {
	result := &EvaluationResult{
		RuleID:      rule.ID,
		Triggered:   false,
		EvaluatedAt: time.Now(),
		Labels:      make(map[string]interface{}),
		Context:     make(map[string]interface{}),
	}

	if len(data) == 0 {
		result.Message = "无指标数据"
		return result
	}

	// 获取最新数据点
	latestData := data[len(data)-1]
	result.Value = latestData.Value

	// 根据规则类型进行评估
	switch rule.RuleType {
	case "threshold":
		result = e.evaluateThresholdRule(rule, latestData, result)
	case "trend":
		result = e.evaluateTrendRule(rule, data, result)
	case "anomaly":
		result = e.evaluateAnomalyRule(rule, data, result)
	default:
		result = e.evaluateThresholdRule(rule, latestData, result)
	}

	// 添加上下文信息
	result.Context["rule_type"] = rule.RuleType
	result.Context["data_points"] = len(data)
	result.Context["evaluation_window"] = rule.EvaluationWindowSeconds

	return result
}

// evaluateThresholdRule 评估阈值规则
func (e *AlertEngine) evaluateThresholdRule(rule *AlertRule, data *model.VtMonitorData, result *EvaluationResult) *EvaluationResult {
	// 确定使用哪个阈值
	threshold := rule.CriticalThreshold
	alertLevel := "critical"

	if rule.WarningThreshold != 0 && rule.CriticalThreshold != 0 {
		// 同时配置了警告和严重阈值
		if e.checkThresholdCondition(data.Value, rule.WarningThreshold, rule.ThresholdCondition) {
			threshold = rule.WarningThreshold
			alertLevel = "warning"
		}
		if e.checkThresholdCondition(data.Value, rule.CriticalThreshold, rule.ThresholdCondition) {
			threshold = rule.CriticalThreshold
			alertLevel = "critical"
		}
	} else if rule.WarningThreshold != 0 {
		threshold = rule.WarningThreshold
		alertLevel = "warning"
	}

	result.Threshold = threshold
	result.Triggered = e.checkThresholdCondition(data.Value, threshold, rule.ThresholdCondition)

	if result.Triggered {
		result.Message = fmt.Sprintf("指标值 %.2f %s 阈值 %.2f",
			data.Value, e.getThresholdOperator(rule.ThresholdCondition), threshold)
		result.Labels["alert_level"] = alertLevel
		result.Labels["threshold"] = threshold
	} else {
		result.Message = fmt.Sprintf("指标值 %.2f 正常", data.Value)
	}

	return result
}

// evaluateTrendRule 评估趋势规则
func (e *AlertEngine) evaluateTrendRule(rule *AlertRule, data []*model.VtMonitorData, result *EvaluationResult) *EvaluationResult {
	if len(data) < 2 {
		result.Message = "数据点不足以分析趋势"
		return result
	}

	// 计算趋势（简化：比较最新值与平均值）
	var sum float64
	for _, point := range data {
		sum += point.Value
	}
	average := sum / float64(len(data))
	latest := data[len(data)-1].Value

	changePercent := ((latest - average) / average) * 100
	threshold := rule.CriticalThreshold // 趋势阈值（百分比）

	result.Value = changePercent
	result.Threshold = threshold
	result.Triggered = e.checkThresholdCondition(changePercent, threshold, rule.ThresholdCondition)

	if result.Triggered {
		result.Message = fmt.Sprintf("趋势变化 %.2f%% 超过阈值 %.2f%%", changePercent, threshold)
	} else {
		result.Message = fmt.Sprintf("趋势变化 %.2f%% 正常", changePercent)
	}

	result.Context["trend_change_percent"] = changePercent
	result.Context["average_value"] = average

	return result
}

// evaluateAnomalyRule 评估异常检测规则
func (e *AlertEngine) evaluateAnomalyRule(rule *AlertRule, data []*model.VtMonitorData, result *EvaluationResult) *EvaluationResult {
	if len(data) < 10 {
		result.Message = "数据点不足以进行异常检测"
		return result
	}

	// 简化的异常检测：使用Z-score
	values := make([]float64, len(data))
	var sum float64
	for i, point := range data {
		values[i] = point.Value
		sum += point.Value
	}

	mean := sum / float64(len(values))

	// 计算标准差
	var variance float64
	for _, value := range values {
		variance += math.Pow(value-mean, 2)
	}
	stdDev := math.Sqrt(variance / float64(len(values)))

	// 计算最新值的Z-score
	latest := data[len(data)-1].Value
	zScore := math.Abs((latest - mean) / stdDev)

	threshold := rule.CriticalThreshold // Z-score阈值
	if threshold == 0 {
		threshold = 3.0 // 默认3个标准差
	}

	result.Value = zScore
	result.Threshold = threshold
	result.Triggered = zScore > threshold

	if result.Triggered {
		result.Message = fmt.Sprintf("异常检测: Z-score %.2f 超过阈值 %.2f", zScore, threshold)
	} else {
		result.Message = fmt.Sprintf("异常检测: Z-score %.2f 正常", zScore)
	}

	result.Context["z_score"] = zScore
	result.Context["mean"] = mean
	result.Context["std_dev"] = stdDev

	return result
}

// checkThresholdCondition 检查阈值条件
func (e *AlertEngine) checkThresholdCondition(value, threshold float64, condition string) bool {
	switch ThresholdCondition(condition) {
	case ConditionGT:
		return value > threshold
	case ConditionGTE:
		return value >= threshold
	case ConditionLT:
		return value < threshold
	case ConditionLTE:
		return value <= threshold
	case ConditionEQ:
		return math.Abs(value-threshold) < 0.001 // 浮点数相等比较
	case ConditionNEQ:
		return math.Abs(value-threshold) >= 0.001
	default:
		return value > threshold // 默认大于
	}
}

// getThresholdOperator 获取阈值操作符描述
func (e *AlertEngine) getThresholdOperator(condition string) string {
	switch ThresholdCondition(condition) {
	case ConditionGT:
		return ">"
	case ConditionGTE:
		return ">="
	case ConditionLT:
		return "<"
	case ConditionLTE:
		return "<="
	case ConditionEQ:
		return "=="
	case ConditionNEQ:
		return "!="
	default:
		return ">"
	}
}

// handleEvaluationResult 处理评估结果
func (e *AlertEngine) handleEvaluationResult(rule *AlertRule, result *EvaluationResult) {
	alertID := e.generateAlertID(rule, result)

	e.mu.Lock()
	existingAlert, exists := e.activeAlerts[alertID]
	e.mu.Unlock()

	if result.Triggered {
		if exists {
			// 更新现有告警
			e.updateActiveAlert(existingAlert, result)
		} else {
			// 创建新告警
			alert := e.createActiveAlert(rule, result, alertID)
			e.mu.Lock()
			e.activeAlerts[alertID] = alert
			e.mu.Unlock()

			// 保存到数据库
			e.saveAlertRecord(alert)

			// 发送通知
			e.sendNotification(alert, rule, "firing")
		}

		// 更新规则触发统计
		e.updateRuleTriggerCount(rule.ID)
	} else {
		if exists && existingAlert.Status == string(StatusFiring) {
			// 解决告警
			e.resolveActiveAlert(existingAlert)

			// 发送恢复通知
			e.sendNotification(existingAlert, rule, "resolved")
		}
	}
}

// generateAlertID 生成告警ID
func (e *AlertEngine) generateAlertID(rule *AlertRule, result *EvaluationResult) string {
	// 简化实现：rule_id + 时间戳的小时部分
	hour := time.Now().Hour()
	return fmt.Sprintf("alert_%d_%d", rule.ID, hour)
}

// createActiveAlert 创建活跃告警
func (e *AlertEngine) createActiveAlert(rule *AlertRule, result *EvaluationResult, alertID string) *ActiveAlert {
	now := time.Now()

	alert := &ActiveAlert{
		ID:                  alertID,
		RuleID:              rule.ID,
		RuleName:            rule.Name,
		AlertLevel:          rule.AlertLevel,
		SeverityScore:       rule.SeverityScore,
		Message:             result.Message,
		Summary:             fmt.Sprintf("%s: %s", rule.DisplayName, result.Message),
		TriggerValue:        result.Value,
		ThresholdValue:      result.Threshold,
		ConditionExpression: rule.ConditionExpression,
		EvaluationData:      result.Context,
		Labels:              result.Labels,
		Status:              string(StatusFiring),
		TriggeredAt:         now,
		FirstOccurrenceAt:   now,
		LastOccurrenceAt:    now,
		OccurrenceCount:     1,
		NotificationSent:    false,
		NotificationCount:   0,
		EscalationLevel:     0,
		AlertGroupID:        "",
		CorrelationID:       alertID,
	}

	return alert
}

// updateActiveAlert 更新活跃告警
func (e *AlertEngine) updateActiveAlert(alert *ActiveAlert, result *EvaluationResult) {
	alert.LastOccurrenceAt = time.Now()
	alert.OccurrenceCount++
	alert.TriggerValue = result.Value
	alert.Message = result.Message
	alert.EvaluationData = result.Context

	// 更新数据库记录
	e.updateAlertRecord(alert)
}

// resolveActiveAlert 解决活跃告警
func (e *AlertEngine) resolveActiveAlert(alert *ActiveAlert) {
	alert.Status = string(StatusResolved)
	alert.LastOccurrenceAt = time.Now()

	// 更新数据库记录
	e.updateAlertRecord(alert)

	// 从活跃告警中移除
	e.mu.Lock()
	delete(e.activeAlerts, alert.ID)
	e.mu.Unlock()
}

// sendNotification 发送通知
func (e *AlertEngine) sendNotification(alert *ActiveAlert, rule *AlertRule, action string) {
	// 检查通知限流
	if !e.shouldSendNotification(alert, rule) {
		return
	}

	notification := &AlertNotification{
		Alert:    alert,
		RuleInfo: rule,
		Action:   action,
	}

	select {
	case e.notificationChannel <- notification:
		alert.NotificationSent = true
		alert.NotificationCount++
		alert.LastNotificationAt = time.Now()
	default:
		e.logger.Errorf("通知队列已满，丢弃告警通知: %s", alert.ID)
	}
}

// shouldSendNotification 检查是否应该发送通知
func (e *AlertEngine) shouldSendNotification(alert *ActiveAlert, rule *AlertRule) bool {
	// 检查通知限流
	if rule.NotificationThrottleMinutes > 0 {
		throttleDuration := time.Duration(rule.NotificationThrottleMinutes) * time.Minute
		if time.Since(alert.LastNotificationAt) < throttleDuration {
			return false
		}
	}

	// 检查抑制规则
	if e.config.EnableSuppression && e.isAlertSuppressed(alert, rule) {
		return false
	}

	return true
}

// isAlertSuppressed 检查告警是否被抑制
func (e *AlertEngine) isAlertSuppressed(alert *ActiveAlert, rule *AlertRule) bool {
	// 简化实现：检查是否有更高级别的告警
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, otherAlert := range e.activeAlerts {
		if otherAlert.ID != alert.ID &&
			otherAlert.SeverityScore > alert.SeverityScore &&
			otherAlert.ResourceType == alert.ResourceType &&
			otherAlert.ResourceID == alert.ResourceID {
			return true
		}
	}

	return false
}

// notificationLoop 通知处理循环
func (e *AlertEngine) notificationLoop() {
	defer func() {
		if r := recover(); r != nil {
			e.logger.Errorf("通知循环异常恢复: %v", r)
		}
	}()

	for {
		select {
		case <-e.ctx.Done():
			e.logger.Info("告警通知循环停止")
			return
		case notification, ok := <-e.notificationChannel:
			if !ok {
				e.logger.Info("通知渠道已关闭，退出通知循环")
				return
			}
			e.processNotification(notification)
		}
	}
}

// processNotification 处理通知
func (e *AlertEngine) processNotification(notification *AlertNotification) {
	e.logger.Infof("处理告警通知: %s [%s]", notification.Alert.ID, notification.Action)

	// 这里集成具体的通知渠道（邮件、短信、钉钉等）
	// 实际实现中会调用通知系统的接口

	// 记录通知日志
	e.logger.Infof("告警通知已发送: 规则=%s, 级别=%s, 动作=%s",
		notification.RuleInfo.Name,
		notification.Alert.AlertLevel,
		notification.Action)
}

// cleanupLoop 清理循环
func (e *AlertEngine) cleanupLoop() {
	ticker := time.NewTicker(24 * time.Hour) // 每天清理一次
	defer ticker.Stop()

	for {
		select {
		case <-e.ctx.Done():
			e.logger.Info("告警清理循环停止")
			return
		case <-ticker.C:
			e.cleanupOldAlerts()
		}
	}
}

// cleanupOldAlerts 清理旧告警
func (e *AlertEngine) cleanupOldAlerts() {
	if e.config.AlertRetentionDays <= 0 {
		return
	}

	cutoffTime := time.Now().AddDate(0, 0, -e.config.AlertRetentionDays)

	// 删除旧的告警记录
	_, err := e.alertRecordsModel.DeleteOldRecords(cutoffTime)
	if err != nil {
		e.logger.Errorf("清理旧告警记录失败: %v", err)
	} else {
		e.logger.Infof("清理了 %d 天前的告警记录", e.config.AlertRetentionDays)
	}
}

// updateRuleEvaluationTime 更新规则评估时间
func (e *AlertEngine) updateRuleEvaluationTime(ruleID int64) {
	err := e.alertRulesModel.UpdateEvaluationTime(ruleID, time.Now())
	if err != nil {
		e.logger.Errorf("更新规则评估时间失败 [%d]: %v", ruleID, err)
	}
}

// updateRuleTriggerCount 更新规则触发次数
func (e *AlertEngine) updateRuleTriggerCount(ruleID int64) {
	err := e.alertRulesModel.IncrementTriggerCount(ruleID)
	if err != nil {
		e.logger.Errorf("更新规则触发次数失败 [%d]: %v", ruleID, err)
	}
}

// saveAlertRecord 保存告警记录
func (e *AlertEngine) saveAlertRecord(alert *ActiveAlert) {
	record := &model.VtAlertRecords{
		RuleId:              alert.RuleID,
		AlertId:             alert.ID,
		AlertName:           alert.RuleName,
		AlertLevel:          alert.AlertLevel,
		SeverityScore:       alert.SeverityScore,
		Message:             alert.Message,
		Summary:             alert.Summary,
		ResourceType:        alert.ResourceType,
		ResourceId:          alert.ResourceID,
		ResourceName:        alert.ResourceName,
		InstanceId:          alert.InstanceID,
		TriggerValue:        alert.TriggerValue,
		ThresholdValue:      alert.ThresholdValue,
		ConditionExpression: alert.ConditionExpression,
		TriggeredAt:         alert.TriggeredAt,
		FirstOccurrenceAt:   alert.FirstOccurrenceAt,
		LastOccurrenceAt:    alert.LastOccurrenceAt,
		Status:              alert.Status,
		OccurrenceCount:     alert.OccurrenceCount,
		NotificationSent:    alert.NotificationSent,
		NotificationCount:   alert.NotificationCount,
		LastNotificationAt:  alert.LastNotificationAt,
		EscalationLevel:     alert.EscalationLevel,
		AlertGroupId:        alert.AlertGroupID,
		CorrelationId:       alert.CorrelationID,
	}

	// 设置JSON字段
	if alert.EvaluationData != nil {
		data, _ := json.Marshal(alert.EvaluationData)
		record.EvaluationData = string(data)
	}
	if alert.Labels != nil {
		data, _ := json.Marshal(alert.Labels)
		record.Labels = string(data)
	}
	if alert.Annotations != nil {
		data, _ := json.Marshal(alert.Annotations)
		record.Annotations = string(data)
	}
	if alert.Context != nil {
		data, _ := json.Marshal(alert.Context)
		record.Context = string(data)
	}

	_, err := e.alertRecordsModel.Insert(record)
	if err != nil {
		e.logger.Errorf("保存告警记录失败: %v", err)
	}
}

// updateAlertRecord 更新告警记录
func (e *AlertEngine) updateAlertRecord(alert *ActiveAlert) {
	err := e.alertRecordsModel.UpdateByAlertId(alert.ID, alert.Status, alert.OccurrenceCount, alert.LastOccurrenceAt)
	if err != nil {
		e.logger.Errorf("更新告警记录失败: %v", err)
	}
}

// ReloadRules 重新加载规则
func (e *AlertEngine) ReloadRules() error {
	e.logger.Info("重新加载告警规则")
	return e.loadAlertRules()
}

// GetActiveAlerts 获取活跃告警
func (e *AlertEngine) GetActiveAlerts() map[string]*ActiveAlert {
	e.mu.RLock()
	defer e.mu.RUnlock()

	alerts := make(map[string]*ActiveAlert)
	for id, alert := range e.activeAlerts {
		alerts[id] = alert
	}

	return alerts
}

// GetEngineStatus 获取引擎状态
func (e *AlertEngine) GetEngineStatus() map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return map[string]interface{}{
		"rules_count":             len(e.rules),
		"active_alerts_count":     len(e.activeAlerts),
		"last_evaluation_time":    e.lastEvaluationTime,
		"evaluation_interval":     e.evaluationInterval,
		"notification_queue_size": len(e.notificationChannel),
	}
}
