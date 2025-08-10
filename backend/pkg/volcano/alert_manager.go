package volcano

import (
	"fmt"
	"sync"
	"time"
)

// AlertManager 告警管理器
type AlertManager struct {
	alerts   map[string]*Alert
	rules    []AlertRule
	channels []AlertChannel
	mutex    sync.RWMutex
}

// Alert 告警定义
type Alert struct {
	ID            string            `json:"id"`
	Type          string            `json:"type"`     // resource, job, infrastructure, security
	Severity      string            `json:"severity"` // info, warning, critical
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	Source        string            `json:"source"` // 告警源
	Labels        map[string]string `json:"labels"`
	Annotations   map[string]string `json:"annotations"`
	Timestamp     time.Time         `json:"timestamp"`
	Status        string            `json:"status"` // active, resolved, silenced
	ResolvedAt    *time.Time        `json:"resolvedAt,omitempty"`
	AckedAt       *time.Time        `json:"ackedAt,omitempty"`
	AckedBy       string            `json:"ackedBy,omitempty"`
	Count         int32             `json:"count"` // 重复次数
	FirstSeen     time.Time         `json:"firstSeen"`
	LastSeen      time.Time         `json:"lastSeen"`
	SilencedUntil *time.Time        `json:"silencedUntil,omitempty"`
}

// AlertRule 告警规则
type AlertRule struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Expression  string            `json:"expression"` // 告警表达式
	Severity    string            `json:"severity"`
	Duration    time.Duration     `json:"duration"` // 持续时间
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Enabled     bool              `json:"enabled"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

// AlertChannel 告警通道
type AlertChannel struct {
	ID      string                 `json:"id"`
	Type    string                 `json:"type"` // email, slack, webhook, dingtalk, sms
	Name    string                 `json:"name"`
	Config  map[string]interface{} `json:"config"`
	Filters []AlertFilter          `json:"filters"` // 过滤条件
	Enabled bool                   `json:"enabled"`
}

// AlertFilter 告警过滤器
type AlertFilter struct {
	Field    string `json:"field"`    // severity, type, source
	Operator string `json:"operator"` // eq, ne, in, not_in, regex
	Value    string `json:"value"`
}

// NewAlertManager 创建告警管理器
func NewAlertManager() *AlertManager {
	return &AlertManager{
		alerts:   make(map[string]*Alert),
		rules:    make([]AlertRule, 0),
		channels: make([]AlertChannel, 0),
	}
}

// TriggerAlert 触发告警
func (am *AlertManager) TriggerAlert(alert Alert) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	// 检查是否是重复告警
	if existingAlert, exists := am.alerts[alert.ID]; exists {
		existingAlert.Count++
		existingAlert.LastSeen = alert.Timestamp
		existingAlert.Description = alert.Description // 更新描述

		// 如果之前已解决，重新激活
		if existingAlert.Status == "resolved" {
			existingAlert.Status = "active"
			existingAlert.ResolvedAt = nil
		}
	} else {
		// 新告警
		alert.FirstSeen = alert.Timestamp
		alert.LastSeen = alert.Timestamp
		alert.Count = 1
		am.alerts[alert.ID] = &alert
	}

	// 发送告警通知
	go am.sendAlertNotifications(&alert)

	return nil
}

// ResolveAlert 解决告警
func (am *AlertManager) ResolveAlert(alertID string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if alert, exists := am.alerts[alertID]; exists {
		alert.Status = "resolved"
		now := time.Now()
		alert.ResolvedAt = &now
		return nil
	}

	return fmt.Errorf("告警 %s 不存在", alertID)
}

// AcknowledgeAlert 确认告警
func (am *AlertManager) AcknowledgeAlert(alertID, userID string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if alert, exists := am.alerts[alertID]; exists {
		now := time.Now()
		alert.AckedAt = &now
		alert.AckedBy = userID
		return nil
	}

	return fmt.Errorf("告警 %s 不存在", alertID)
}

// SilenceAlert 静默告警
func (am *AlertManager) SilenceAlert(alertID string, duration time.Duration) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if alert, exists := am.alerts[alertID]; exists {
		alert.Status = "silenced"
		silencedUntil := time.Now().Add(duration)
		alert.SilencedUntil = &silencedUntil
		return nil
	}

	return fmt.Errorf("告警 %s 不存在", alertID)
}

// ListAlerts 列出告警
func (am *AlertManager) ListAlerts(filter AlertListFilter) ([]*Alert, error) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	var alerts []*Alert
	for _, alert := range am.alerts {
		if am.matchesFilter(alert, filter) {
			alerts = append(alerts, alert)
		}
	}

	return alerts, nil
}

// AlertListFilter 告警列表过滤器
type AlertListFilter struct {
	Status    string            `json:"status,omitempty"`   // active, resolved, silenced
	Severity  string            `json:"severity,omitempty"` // info, warning, critical
	Type      string            `json:"type,omitempty"`     // resource, job, infrastructure
	Source    string            `json:"source,omitempty"`
	TimeRange TimeRange         `json:"timeRange,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// TimeRange 时间范围
type TimeRange struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// matchesFilter 检查告警是否匹配过滤条件
func (am *AlertManager) matchesFilter(alert *Alert, filter AlertListFilter) bool {
	if filter.Status != "" && alert.Status != filter.Status {
		return false
	}

	if filter.Severity != "" && alert.Severity != filter.Severity {
		return false
	}

	if filter.Type != "" && alert.Type != filter.Type {
		return false
	}

	if filter.Source != "" && alert.Source != filter.Source {
		return false
	}

	// 检查时间范围
	if !filter.TimeRange.StartTime.IsZero() && alert.Timestamp.Before(filter.TimeRange.StartTime) {
		return false
	}
	if !filter.TimeRange.EndTime.IsZero() && alert.Timestamp.After(filter.TimeRange.EndTime) {
		return false
	}

	// 检查标签匹配
	for key, value := range filter.Labels {
		if alertValue, exists := alert.Labels[key]; !exists || alertValue != value {
			return false
		}
	}

	return true
}

// GetAlertStatistics 获取告警统计
func (am *AlertManager) GetAlertStatistics() *AlertStatistics {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	stats := &AlertStatistics{
		Total:      int32(len(am.alerts)),
		BySeverity: make(map[string]int32),
		ByType:     make(map[string]int32),
		ByStatus:   make(map[string]int32),
		BySource:   make(map[string]int32),
	}

	for _, alert := range am.alerts {
		stats.BySeverity[alert.Severity]++
		stats.ByType[alert.Type]++
		stats.ByStatus[alert.Status]++
		stats.BySource[alert.Source]++

		switch alert.Status {
		case "active":
			stats.Active++
		case "resolved":
			stats.Resolved++
		case "silenced":
			stats.Silenced++
		}

		switch alert.Severity {
		case "critical":
			stats.Critical++
		case "warning":
			stats.Warning++
		case "info":
			stats.Info++
		}
	}

	return stats
}

// AlertStatistics 告警统计
type AlertStatistics struct {
	Total      int32            `json:"total"`
	Active     int32            `json:"active"`
	Resolved   int32            `json:"resolved"`
	Silenced   int32            `json:"silenced"`
	Critical   int32            `json:"critical"`
	Warning    int32            `json:"warning"`
	Info       int32            `json:"info"`
	BySeverity map[string]int32 `json:"bySeverity"`
	ByType     map[string]int32 `json:"byType"`
	ByStatus   map[string]int32 `json:"byStatus"`
	BySource   map[string]int32 `json:"bySource"`
}

// AddAlertRule 添加告警规则
func (am *AlertManager) AddAlertRule(rule AlertRule) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()
	am.rules = append(am.rules, rule)

	return nil
}

// UpdateAlertRule 更新告警规则
func (am *AlertManager) UpdateAlertRule(ruleID string, rule AlertRule) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	for i, r := range am.rules {
		if r.ID == ruleID {
			rule.ID = ruleID
			rule.CreatedAt = r.CreatedAt
			rule.UpdatedAt = time.Now()
			am.rules[i] = rule
			return nil
		}
	}

	return fmt.Errorf("告警规则 %s 不存在", ruleID)
}

// DeleteAlertRule 删除告警规则
func (am *AlertManager) DeleteAlertRule(ruleID string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	for i, rule := range am.rules {
		if rule.ID == ruleID {
			am.rules = append(am.rules[:i], am.rules[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("告警规则 %s 不存在", ruleID)
}

// ListAlertRules 列出告警规则
func (am *AlertManager) ListAlertRules() ([]AlertRule, error) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	return am.rules, nil
}

// AddAlertChannel 添加告警通道
func (am *AlertManager) AddAlertChannel(channel AlertChannel) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	am.channels = append(am.channels, channel)
	return nil
}

// UpdateAlertChannel 更新告警通道
func (am *AlertManager) UpdateAlertChannel(channelID string, channel AlertChannel) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	for i, c := range am.channels {
		if c.ID == channelID {
			channel.ID = channelID
			am.channels[i] = channel
			return nil
		}
	}

	return fmt.Errorf("告警通道 %s 不存在", channelID)
}

// DeleteAlertChannel 删除告警通道
func (am *AlertManager) DeleteAlertChannel(channelID string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	for i, channel := range am.channels {
		if channel.ID == channelID {
			am.channels = append(am.channels[:i], am.channels[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("告警通道 %s 不存在", channelID)
}

// ListAlertChannels 列出告警通道
func (am *AlertManager) ListAlertChannels() ([]AlertChannel, error) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	return am.channels, nil
}

// sendAlertNotifications 发送告警通知
func (am *AlertManager) sendAlertNotifications(alert *Alert) {
	for _, channel := range am.channels {
		if !channel.Enabled {
			continue
		}

		if !am.shouldSendToChannel(alert, channel) {
			continue
		}

		if err := am.sendToChannel(alert, channel); err != nil {
			// TODO: 记录发送失败的日志
			fmt.Printf("发送告警到通道 %s 失败: %v\n", channel.Name, err)
		}
	}
}

// shouldSendToChannel 检查是否应该发送到指定通道
func (am *AlertManager) shouldSendToChannel(alert *Alert, channel AlertChannel) bool {
	for _, filter := range channel.Filters {
		if !am.applyFilter(alert, filter) {
			return false
		}
	}
	return true
}

// applyFilter 应用过滤器
func (am *AlertManager) applyFilter(alert *Alert, filter AlertFilter) bool {
	var value string

	switch filter.Field {
	case "severity":
		value = alert.Severity
	case "type":
		value = alert.Type
	case "source":
		value = alert.Source
	case "status":
		value = alert.Status
	default:
		if labelValue, exists := alert.Labels[filter.Field]; exists {
			value = labelValue
		} else {
			return false
		}
	}

	switch filter.Operator {
	case "eq":
		return value == filter.Value
	case "ne":
		return value != filter.Value
	case "in":
		// TODO: 实现 in 操作符
		return true
	case "not_in":
		// TODO: 实现 not_in 操作符
		return true
	case "regex":
		// TODO: 实现正则表达式匹配
		return true
	default:
		return false
	}
}

// sendToChannel 发送告警到指定通道
func (am *AlertManager) sendToChannel(alert *Alert, channel AlertChannel) error {
	switch channel.Type {
	case "email":
		return am.sendEmailAlert(alert, channel)
	case "slack":
		return am.sendSlackAlert(alert, channel)
	case "webhook":
		return am.sendWebhookAlert(alert, channel)
	case "dingtalk":
		return am.sendDingTalkAlert(alert, channel)
	case "sms":
		return am.sendSMSAlert(alert, channel)
	default:
		return fmt.Errorf("不支持的通道类型: %s", channel.Type)
	}
}

// sendEmailAlert 发送邮件告警
func (am *AlertManager) sendEmailAlert(alert *Alert, channel AlertChannel) error {
	// TODO: 实现邮件发送逻辑
	fmt.Printf("发送邮件告警: %s - %s\n", alert.Title, alert.Description)
	return nil
}

// sendSlackAlert 发送Slack告警
func (am *AlertManager) sendSlackAlert(alert *Alert, channel AlertChannel) error {
	// TODO: 实现Slack发送逻辑
	fmt.Printf("发送Slack告警: %s - %s\n", alert.Title, alert.Description)
	return nil
}

// sendWebhookAlert 发送Webhook告警
func (am *AlertManager) sendWebhookAlert(alert *Alert, channel AlertChannel) error {
	// TODO: 实现Webhook发送逻辑
	fmt.Printf("发送Webhook告警: %s - %s\n", alert.Title, alert.Description)
	return nil
}

// sendDingTalkAlert 发送钉钉告警
func (am *AlertManager) sendDingTalkAlert(alert *Alert, channel AlertChannel) error {
	// TODO: 实现钉钉发送逻辑
	fmt.Printf("发送钉钉告警: %s - %s\n", alert.Title, alert.Description)
	return nil
}

// sendSMSAlert 发送短信告警
func (am *AlertManager) sendSMSAlert(alert *Alert, channel AlertChannel) error {
	// TODO: 实现短信发送逻辑
	fmt.Printf("发送短信告警: %s - %s\n", alert.Title, alert.Description)
	return nil
}

// EvaluateRules 评估告警规则
func (am *AlertManager) EvaluateRules(metrics map[string]interface{}) error {
	am.mutex.RLock()
	rules := make([]AlertRule, len(am.rules))
	copy(rules, am.rules)
	am.mutex.RUnlock()

	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		if am.evaluateExpression(rule.Expression, metrics) {
			alert := Alert{
				ID:          fmt.Sprintf("%s-%d", rule.ID, time.Now().Unix()),
				Type:        "rule",
				Severity:    rule.Severity,
				Title:       rule.Name,
				Description: fmt.Sprintf("规则 %s 被触发", rule.Name),
				Source:      "alert-rule",
				Labels:      rule.Labels,
				Annotations: rule.Annotations,
				Timestamp:   time.Now(),
				Status:      "active",
			}

			am.TriggerAlert(alert)
		}
	}

	return nil
}

// evaluateExpression 评估告警表达式
func (am *AlertManager) evaluateExpression(expression string, metrics map[string]interface{}) bool {
	// TODO: 实现告警表达式评估逻辑
	// 这里是简化实现，实际需要支持复杂的表达式解析
	return false
}

// GetAlertHistory 获取告警历史
func (am *AlertManager) GetAlertHistory(timeRange TimeRange, limit int) ([]*AlertHistoryEntry, error) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	var history []*AlertHistoryEntry

	for _, alert := range am.alerts {
		// 检查时间范围
		if !timeRange.StartTime.IsZero() && alert.Timestamp.Before(timeRange.StartTime) {
			continue
		}
		if !timeRange.EndTime.IsZero() && alert.Timestamp.After(timeRange.EndTime) {
			continue
		}

		entry := &AlertHistoryEntry{
			AlertID:     alert.ID,
			Title:       alert.Title,
			Description: alert.Description,
			Severity:    alert.Severity,
			Type:        alert.Type,
			Source:      alert.Source,
			Status:      alert.Status,
			Timestamp:   alert.Timestamp,
			Count:       alert.Count,
		}

		if alert.ResolvedAt != nil {
			entry.Duration = alert.ResolvedAt.Sub(alert.Timestamp)
		}

		history = append(history, entry)
	}

	// 按时间降序排序
	for i := 0; i < len(history)-1; i++ {
		for j := i + 1; j < len(history); j++ {
			if history[i].Timestamp.Before(history[j].Timestamp) {
				history[i], history[j] = history[j], history[i]
			}
		}
	}

	// 限制返回数量
	if limit > 0 && len(history) > limit {
		history = history[:limit]
	}

	return history, nil
}

// AlertHistoryEntry 告警历史条目
type AlertHistoryEntry struct {
	AlertID     string        `json:"alertId"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Severity    string        `json:"severity"`
	Type        string        `json:"type"`
	Source      string        `json:"source"`
	Status      string        `json:"status"`
	Timestamp   time.Time     `json:"timestamp"`
	Duration    time.Duration `json:"duration,omitempty"`
	Count       int32         `json:"count"`
}

// CleanupResolvedAlerts 清理已解决的告警
func (am *AlertManager) CleanupResolvedAlerts(olderThan time.Duration) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	cutoffTime := time.Now().Add(-olderThan)

	for id, alert := range am.alerts {
		if alert.Status == "resolved" && alert.ResolvedAt != nil && alert.ResolvedAt.Before(cutoffTime) {
			delete(am.alerts, id)
		}
	}

	return nil
}

// ExportAlerts 导出告警数据
func (am *AlertManager) ExportAlerts(format string) ([]byte, error) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	switch format {
	case "json":
		// TODO: 实现JSON导出
		return nil, nil
	case "csv":
		// TODO: 实现CSV导出
		return nil, nil
	default:
		return nil, fmt.Errorf("不支持的导出格式: %s", format)
	}
}

// ImportAlerts 导入告警数据
func (am *AlertManager) ImportAlerts(data []byte, format string) error {
	switch format {
	case "json":
		// TODO: 实现JSON导入
		return nil
	case "csv":
		// TODO: 实现CSV导入
		return nil
	default:
		return fmt.Errorf("不支持的导入格式: %s", format)
	}
}
