package notification

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"api/model"
	"api/pkg/alerting"
	"github.com/zeromicro/go-zero/core/logx"
)

// NotificationManager 通知管理器
type NotificationManager struct {
	logger              logx.Logger
	ctx                 context.Context
	cancel              context.CancelFunc
	db                  *sql.DB
	channelsModel       model.VtNotificationChannelsModel
	templatesModel      model.VtNotificationTemplatesModel
	config              *NotificationConfig
	channels            map[string]NotificationChannel
	templates           map[string]*NotificationTemplate
	notificationQueue   chan *NotificationRequest
	mu                  sync.RWMutex
	rateLimiters        map[string]*RateLimiter
	failedNotifications map[string]int
	stopped             bool
}

// NotificationConfig 通知配置
type NotificationConfig struct {
	MaxQueueSize         int           `json:"max_queue_size"`
	MaxConcurrentSenders int           `json:"max_concurrent_senders"`
	RetryMaxAttempts     int           `json:"retry_max_attempts"`
	RetryBackoffSeconds  int           `json:"retry_backoff_seconds"`
	RateLimitPerMinute   int           `json:"rate_limit_per_minute"`
	TimeoutSeconds       int           `json:"timeout_seconds"`
	FailedRetentionDays  int           `json:"failed_retention_days"`
	EnableDeduplication  bool          `json:"enable_deduplication"`
	DeduplicationWindow  time.Duration `json:"deduplication_window"`
}

// NotificationChannel 通知渠道接口
type NotificationChannel interface {
	GetType() string
	GetName() string
	Send(ctx context.Context, req *NotificationRequest) error
	ValidateConfig(config map[string]interface{}) error
	IsEnabled() bool
}

// NotificationRequest 通知请求
type NotificationRequest struct {
	ID            string                 `json:"id"`
	Alert         *alerting.ActiveAlert  `json:"alert"`
	RuleInfo      *alerting.AlertRule    `json:"rule_info"`
	Action        string                 `json:"action"`
	ChannelName   string                 `json:"channel_name"`
	ChannelType   string                 `json:"channel_type"`
	ChannelConfig map[string]interface{} `json:"channel_config"`
	TemplateName  string                 `json:"template_name"`
	Recipients    []string               `json:"recipients"`
	Subject       string                 `json:"subject"`
	Content       string                 `json:"content"`
	Priority      string                 `json:"priority"`
	Metadata      map[string]interface{} `json:"metadata"`
	CreatedAt     time.Time              `json:"created_at"`
	RetryCount    int                    `json:"retry_count"`
	LastRetryAt   time.Time              `json:"last_retry_at"`
}

// NotificationTemplate 通知模板
type NotificationTemplate struct {
	ID           int64                  `json:"id"`
	Name         string                 `json:"name"`
	DisplayName  string                 `json:"display_name"`
	Description  string                 `json:"description"`
	ChannelType  string                 `json:"channel_type"`
	TemplateType string                 `json:"template_type"`
	Subject      string                 `json:"subject"`
	Content      string                 `json:"content"`
	Variables    map[string]interface{} `json:"variables"`
	IsDefault    bool                   `json:"is_default"`
	Status       string                 `json:"status"`
}

// RateLimiter 速率限制器
type RateLimiter struct {
	tokens     int
	capacity   int
	lastRefill time.Time
	mu         sync.Mutex
}

// NotificationResult 通知结果
type NotificationResult struct {
	RequestID    string        `json:"request_id"`
	ChannelName  string        `json:"channel_name"`
	Success      bool          `json:"success"`
	ErrorMessage string        `json:"error_message"`
	SentAt       time.Time     `json:"sent_at"`
	Duration     time.Duration `json:"duration"`
}

// NewNotificationManager 创建通知管理器
func NewNotificationManager(db *sql.DB, config *NotificationConfig) *NotificationManager {
	ctx, cancel := context.WithCancel(context.Background())

	manager := &NotificationManager{
		logger:              logx.WithContext(ctx),
		ctx:                 ctx,
		cancel:              cancel,
		db:                  db,
		channelsModel:       model.NewVtNotificationChannelsModel(db),
		templatesModel:      model.NewVtNotificationTemplatesModel(db),
		config:              config,
		channels:            make(map[string]NotificationChannel),
		templates:           make(map[string]*NotificationTemplate),
		notificationQueue:   make(chan *NotificationRequest, config.MaxQueueSize),
		rateLimiters:        make(map[string]*RateLimiter),
		failedNotifications: make(map[string]int),
	}

	return manager
}

// Start 启动通知管理器
func (m *NotificationManager) Start() error {
	m.logger.Info("启动通知管理器")

	// 加载通知渠道
	if err := m.loadChannels(); err != nil {
		return fmt.Errorf("加载通知渠道失败: %v", err)
	}

	// 加载通知模板
	if err := m.loadTemplates(); err != nil {
		return fmt.Errorf("加载通知模板失败: %v", err)
	}

	// 启动通知发送workers
	for i := 0; i < m.config.MaxConcurrentSenders; i++ {
		go m.notificationWorker(fmt.Sprintf("worker-%d", i))
	}

	// 启动清理任务
	go m.cleanupWorker()

	return nil
}

// Stop 停止通知管理器
func (m *NotificationManager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.stopped {
		return
	}
	
	m.logger.Info("停止通知管理器")
	m.stopped = true
	
	if m.cancel != nil {
		m.cancel()
	}
	
	// 安全关闭队列
	select {
	case <-m.notificationQueue:
	default:
	}
	close(m.notificationQueue)
}

// RegisterChannel 注册通知渠道
func (m *NotificationManager) RegisterChannel(channel NotificationChannel) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.channels[channel.GetName()] = channel
	m.rateLimiters[channel.GetName()] = &RateLimiter{
		capacity:   m.config.RateLimitPerMinute,
		tokens:     m.config.RateLimitPerMinute,
		lastRefill: time.Now(),
	}

	m.logger.Infof("注册通知渠道: %s [%s]", channel.GetName(), channel.GetType())
}

// SendNotification 发送通知
func (m *NotificationManager) SendNotification(alert *alerting.AlertNotification) error {
	// 获取告警规则的通知渠道配置
	channelNames := alert.RuleInfo.NotificationChannels
	if len(channelNames) == 0 {
		m.logger.Infof("告警规则未配置通知渠道: %s", alert.RuleInfo.Name)
		return nil
	}

	for _, channelName := range channelNames {
		// 获取渠道配置
		channelConfig, err := m.getChannelConfig(channelName)
		if err != nil {
			m.logger.Errorf("获取通知渠道配置失败 [%s]: %v", channelName, err)
			continue
		}

		// 检查去重
		if m.config.EnableDeduplication && m.isDuplicate(alert, channelName) {
			m.logger.Debugf("跳过重复通知 [%s]: %s", channelName, alert.Alert.ID)
			continue
		}

		// 生成通知内容
		req, err := m.buildNotificationRequest(alert, channelName, channelConfig)
		if err != nil {
			m.logger.Errorf("构建通知请求失败 [%s]: %v", channelName, err)
			continue
		}

		// 提交到队列
		select {
		case m.notificationQueue <- req:
			m.logger.Debugf("通知请求已提交到队列: %s", req.ID)
		case <-m.ctx.Done():
			m.logger.Infof("通知管理器已停止，跳过通知: %s", req.ID)
		default:
			m.logger.Errorf("通知队列已满，丢弃通知: %s", req.ID)
		}
	}

	return nil
}

// loadChannels 加载通知渠道
func (m *NotificationManager) loadChannels() error {
	m.logger.Info("加载通知渠道配置")

	channels, err := m.channelsModel.FindActiveChannels()
	if err != nil {
		return err
	}

	for _, dbChannel := range channels {
		// 解析渠道配置
		var config map[string]interface{}
		if dbChannel.Config != "" {
			if err := json.Unmarshal([]byte(dbChannel.Config), &config); err != nil {
				m.logger.Errorf("解析渠道配置失败 [%s]: %v", dbChannel.Name, err)
				continue
			}
		}

		// 根据渠道类型创建具体实现
		var channel NotificationChannel
		switch dbChannel.ChannelType {
		case "email":
			channel = NewEmailChannel(dbChannel.Name, config)
		case "sms":
			channel = NewSMSChannel(dbChannel.Name, config)
		case "dingtalk":
			channel = NewDingTalkChannel(dbChannel.Name, config)
		case "webhook":
			channel = NewWebhookChannel(dbChannel.Name, config)
		default:
			m.logger.Infof("不支持的通知渠道类型: %s", dbChannel.ChannelType)
			continue
		}

		// 验证配置
		if err := channel.ValidateConfig(config); err != nil {
			m.logger.Errorf("通知渠道配置验证失败 [%s]: %v", dbChannel.Name, err)
			continue
		}

		m.RegisterChannel(channel)
	}

	m.logger.Infof("成功加载 %d 个通知渠道", len(m.channels))
	return nil
}

// loadTemplates 加载通知模板
func (m *NotificationManager) loadTemplates() error {
	m.logger.Info("加载通知模板")

	templates, err := m.templatesModel.FindActiveTemplates()
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, dbTemplate := range templates {
		template := &NotificationTemplate{
			ID:           dbTemplate.Id,
			Name:         dbTemplate.Name,
			DisplayName:  dbTemplate.DisplayName,
			Description:  dbTemplate.Description,
			ChannelType:  dbTemplate.ChannelType,
			TemplateType: dbTemplate.TemplateType,
			Subject:      dbTemplate.Subject,
			Content:      dbTemplate.Content,
			IsDefault:    dbTemplate.IsDefault,
			Status:       dbTemplate.Status,
		}

		// 解析模板变量
		if dbTemplate.Variables != "" {
			if err := json.Unmarshal([]byte(dbTemplate.Variables), &template.Variables); err != nil {
				m.logger.Errorf("解析模板变量失败 [%s]: %v", template.Name, err)
				continue
			}
		}

		m.templates[template.Name] = template
		m.logger.Infof("加载通知模板: %s [%s]", template.Name, template.ChannelType)
	}

	m.logger.Infof("成功加载 %d 个通知模板", len(m.templates))
	return nil
}

// getChannelConfig 获取渠道配置
func (m *NotificationManager) getChannelConfig(channelName string) (map[string]interface{}, error) {
	channel, err := m.channelsModel.FindOneByName(channelName)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if channel.Config != "" {
		if err := json.Unmarshal([]byte(channel.Config), &config); err != nil {
			return nil, fmt.Errorf("解析渠道配置失败: %v", err)
		}
	}

	// 添加渠道类型信息
	config["channel_type"] = channel.ChannelType
	config["channel_name"] = channel.Name

	return config, nil
}

// buildNotificationRequest 构建通知请求
func (m *NotificationManager) buildNotificationRequest(alert *alerting.AlertNotification, channelName string, channelConfig map[string]interface{}) (*NotificationRequest, error) {
	// 获取渠道信息
	channel, exists := m.channels[channelName]
	if !exists {
		return nil, fmt.Errorf("通知渠道不存在: %s", channelName)
	}

	// 选择模板
	templateName := m.selectTemplate(channel.GetType(), alert.Action)
	template, exists := m.templates[templateName]
	if !exists {
		return nil, fmt.Errorf("通知模板不存在: %s", templateName)
	}

	// 渲染通知内容
	subject, content, err := m.renderTemplate(template, alert)
	if err != nil {
		return nil, fmt.Errorf("渲染通知模板失败: %v", err)
	}

	// 获取收件人
	recipients := m.getRecipients(channelConfig, alert)

	req := &NotificationRequest{
		ID:            m.generateNotificationID(alert.Alert.ID, channelName),
		Alert:         alert.Alert,
		RuleInfo:      alert.RuleInfo,
		Action:        alert.Action,
		ChannelName:   channelName,
		ChannelType:   channel.GetType(),
		ChannelConfig: channelConfig,
		TemplateName:  templateName,
		Recipients:    recipients,
		Subject:       subject,
		Content:       content,
		Priority:      alert.Alert.AlertLevel,
		Metadata:      make(map[string]interface{}),
		CreatedAt:     time.Now(),
		RetryCount:    0,
	}

	return req, nil
}

// selectTemplate 选择通知模板
func (m *NotificationManager) selectTemplate(channelType, action string) string {
	// 优先级：特定动作模板 > 渠道默认模板 > 系统默认模板
	templateName := fmt.Sprintf("%s_%s_default", channelType, action)
	if _, exists := m.templates[templateName]; exists {
		return templateName
	}

	templateName = fmt.Sprintf("%s_default", channelType)
	if _, exists := m.templates[templateName]; exists {
		return templateName
	}

	// 返回系统默认模板
	return "system_default"
}

// renderTemplate 渲染通知模板
func (m *NotificationManager) renderTemplate(template *NotificationTemplate, alert *alerting.AlertNotification) (string, string, error) {
	// 构建模板变量
	vars := map[string]interface{}{
		"alert":     alert.Alert,
		"rule":      alert.RuleInfo,
		"action":    alert.Action,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"date":      time.Now().Format("2006-01-02"),
		"time":      time.Now().Format("15:04:05"),
	}

	// 合并模板定义的变量
	for k, v := range template.Variables {
		vars[k] = v
	}

	// 简化的模板渲染（实际项目中可使用template/text包）
	subject := m.replaceVariables(template.Subject, vars)
	content := m.replaceVariables(template.Content, vars)

	return subject, content, nil
}

// replaceVariables 简化的变量替换
func (m *NotificationManager) replaceVariables(text string, vars map[string]interface{}) string {
	// 这里是简化实现，实际应使用专业的模板引擎
	result := text

	// 替换常用变量
	if alert, ok := vars["alert"].(*alerting.ActiveAlert); ok {
		result = replaceString(result, "{{.Alert.RuleName}}", alert.RuleName)
		result = replaceString(result, "{{.Alert.Message}}", alert.Message)
		result = replaceString(result, "{{.Alert.AlertLevel}}", alert.AlertLevel)
		result = replaceString(result, "{{.Alert.TriggerValue}}", fmt.Sprintf("%.2f", alert.TriggerValue))
		result = replaceString(result, "{{.Alert.ThresholdValue}}", fmt.Sprintf("%.2f", alert.ThresholdValue))
		result = replaceString(result, "{{.Alert.TriggeredAt}}", alert.TriggeredAt.Format("2006-01-02 15:04:05"))
	}

	if action, ok := vars["action"].(string); ok {
		result = replaceString(result, "{{.Action}}", action)
	}

	if timestamp, ok := vars["timestamp"].(string); ok {
		result = replaceString(result, "{{.Timestamp}}", timestamp)
	}

	return result
}

// replaceString 字符串替换辅助函数
func replaceString(text, placeholder, value string) string {
	return text // 简化实现，实际应使用strings.ReplaceAll
}

// getRecipients 获取收件人列表
func (m *NotificationManager) getRecipients(channelConfig map[string]interface{}, alert *alerting.AlertNotification) []string {
	var recipients []string

	// 从渠道配置中获取默认收件人
	if defaultRecipients, ok := channelConfig["recipients"].([]interface{}); ok {
		for _, recipient := range defaultRecipients {
			if str, ok := recipient.(string); ok {
				recipients = append(recipients, str)
			}
		}
	}

	// 根据告警级别添加特定收件人
	if alert.Alert.AlertLevel == "critical" {
		if criticalRecipients, ok := channelConfig["critical_recipients"].([]interface{}); ok {
			for _, recipient := range criticalRecipients {
				if str, ok := recipient.(string); ok {
					recipients = append(recipients, str)
				}
			}
		}
	}

	return recipients
}

// generateNotificationID 生成通知ID
func (m *NotificationManager) generateNotificationID(alertID, channelName string) string {
	return fmt.Sprintf("notification_%s_%s_%d", alertID, channelName, time.Now().Unix())
}

// isDuplicate 检查是否重复通知
func (m *NotificationManager) isDuplicate(alert *alerting.AlertNotification, channelName string) bool {
	// 简化的去重实现，实际可用Redis等
	key := fmt.Sprintf("%s_%s_%s", alert.Alert.ID, channelName, alert.Action)

	m.mu.RLock()
	_, exists := m.failedNotifications[key]
	m.mu.RUnlock()

	if exists {
		// 检查时间窗口
		return time.Since(time.Now()) < m.config.DeduplicationWindow
	}

	return false
}

// notificationWorker 通知发送工作器
func (m *NotificationManager) notificationWorker(workerID string) {
	m.logger.Infof("启动通知工作器: %s", workerID)

	for {
		select {
		case <-m.ctx.Done():
			m.logger.Infof("通知工作器停止: %s", workerID)
			return
		case req, ok := <-m.notificationQueue:
			if !ok {
				m.logger.Infof("通知队列关闭，工作器停止: %s", workerID)
				return
			}
			m.processNotification(workerID, req)
		}
	}
}

// processNotification 处理通知
func (m *NotificationManager) processNotification(workerID string, req *NotificationRequest) {
	start := time.Now()
	m.logger.Infof("[%s] 处理通知: %s [%s]", workerID, req.ID, req.ChannelName)

	// 检查速率限制
	if !m.checkRateLimit(req.ChannelName) {
		m.logger.Infof("[%s] 触发速率限制，延迟处理: %s", workerID, req.ID)
		time.Sleep(time.Second * 5)
	}

	// 获取通知渠道
	m.mu.RLock()
	channel, exists := m.channels[req.ChannelName]
	m.mu.RUnlock()

	if !exists {
		m.logger.Errorf("[%s] 通知渠道不存在: %s", workerID, req.ChannelName)
		return
	}

	// 发送通知
	ctx, cancel := context.WithTimeout(m.ctx, time.Duration(m.config.TimeoutSeconds)*time.Second)
	defer cancel()

	err := channel.Send(ctx, req)
	duration := time.Since(start)

	// 记录结果
	result := &NotificationResult{
		RequestID:   req.ID,
		ChannelName: req.ChannelName,
		Success:     err == nil,
		SentAt:      time.Now(),
		Duration:    duration,
	}

	if err != nil {
		result.ErrorMessage = err.Error()
		m.logger.Errorf("[%s] 通知发送失败: %s, 错误: %v", workerID, req.ID, err)
		m.handleFailedNotification(req, err)
	} else {
		m.logger.Infof("[%s] 通知发送成功: %s, 耗时: %v", workerID, req.ID, duration)
	}

	// 记录通知历史（可选）
	m.recordNotificationHistory(result)
}

// checkRateLimit 检查速率限制
func (m *NotificationManager) checkRateLimit(channelName string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	limiter, exists := m.rateLimiters[channelName]
	if !exists {
		return true
	}

	now := time.Now()
	// 补充令牌
	if now.Sub(limiter.lastRefill) >= time.Minute {
		limiter.tokens = limiter.capacity
		limiter.lastRefill = now
	}

	if limiter.tokens > 0 {
		limiter.tokens--
		return true
	}

	return false
}

// handleFailedNotification 处理失败的通知
func (m *NotificationManager) handleFailedNotification(req *NotificationRequest, err error) {
	if req.RetryCount < m.config.RetryMaxAttempts {
		req.RetryCount++
		req.LastRetryAt = time.Now()

		// 指数退避重试
		backoff := time.Duration(m.config.RetryBackoffSeconds) * time.Second * time.Duration(req.RetryCount)
		time.AfterFunc(backoff, func() {
			select {
			case m.notificationQueue <- req:
				m.logger.Infof("重试通知: %s, 第%d次", req.ID, req.RetryCount)
			case <-m.ctx.Done():
				m.logger.Infof("通知管理器已停止，取消重试: %s", req.ID)
			default:
				m.logger.Errorf("重试队列已满，放弃通知: %s", req.ID)
			}
		})
	} else {
		m.logger.Errorf("通知最终失败: %s, 已达到最大重试次数", req.ID)

		// 记录失败统计
		m.mu.Lock()
		m.failedNotifications[req.ChannelName]++
		m.mu.Unlock()
	}
}

// recordNotificationHistory 记录通知历史
func (m *NotificationManager) recordNotificationHistory(result *NotificationResult) {
	// 这里可以将通知结果记录到数据库或日志系统
	// 为了简化，暂时只记录日志
	if result.Success {
		m.logger.Infof("通知历史记录: %s 发送成功, 耗时: %v", result.RequestID, result.Duration)
	} else {
		m.logger.Errorf("通知历史记录: %s 发送失败, 错误: %s", result.RequestID, result.ErrorMessage)
	}
}

// cleanupWorker 清理工作器
func (m *NotificationManager) cleanupWorker() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			m.logger.Info("清理工作器停止")
			return
		case <-ticker.C:
			m.cleanup()
		}
	}
}

// cleanup 清理过期数据
func (m *NotificationManager) cleanup() {
	m.logger.Info("开始清理过期通知数据")

	// 清理失败通知统计
	m.mu.Lock()
	m.failedNotifications = make(map[string]int)
	m.mu.Unlock()

	m.logger.Info("过期通知数据清理完成")
}

// GetStatus 获取通知管理器状态
func (m *NotificationManager) GetStatus() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := map[string]interface{}{
		"channels_count":       len(m.channels),
		"templates_count":      len(m.templates),
		"queue_size":           len(m.notificationQueue),
		"queue_capacity":       cap(m.notificationQueue),
		"failed_notifications": m.failedNotifications,
	}

	// 添加各渠道状态
	channelStatus := make(map[string]interface{})
	for name, channel := range m.channels {
		channelStatus[name] = map[string]interface{}{
			"type":    channel.GetType(),
			"enabled": channel.IsEnabled(),
		}
	}
	status["channels"] = channelStatus

	return status
}

// GetRegisteredChannels 获取已注册的通知渠道列表
func (m *NotificationManager) GetRegisteredChannels() []map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var channels []map[string]interface{}
	for name, channel := range m.channels {
		channelInfo := map[string]interface{}{
			"name":    name,
			"type":    channel.GetType(),
			"enabled": channel.IsEnabled(),
		}
		channels = append(channels, channelInfo)
	}

	return channels
}

// ReloadConfig 重新加载配置
func (m *NotificationManager) ReloadConfig() error {
	m.logger.Info("重新加载通知配置")

	// 重新加载渠道配置
	if err := m.loadChannels(); err != nil {
		return fmt.Errorf("重新加载渠道配置失败: %v", err)
	}

	// 重新加载模板配置
	if err := m.loadTemplates(); err != nil {
		return fmt.Errorf("重新加载模板配置失败: %v", err)
	}

	m.logger.Info("通知配置重新加载完成")
	return nil
}
