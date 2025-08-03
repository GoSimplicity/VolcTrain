package notification

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// WebhookChannel Webhook通知渠道
type WebhookChannel struct {
	name    string
	config  *WebhookConfig
	enabled bool
	logger  logx.Logger
	client  *http.Client
}

// WebhookConfig Webhook配置
type WebhookConfig struct {
	URL             string            `json:"url"`               // Webhook URL
	Method          string            `json:"method"`            // HTTP方法，默认POST
	Headers         map[string]string `json:"headers"`           // 自定义请求头
	Secret          string            `json:"secret"`            // 签名密钥
	SignatureHeader string            `json:"signature_header"`  // 签名头名称
	SignaturePrefix string            `json:"signature_prefix"`  // 签名前缀
	TimeoutSeconds  int               `json:"timeout_seconds"`   // 超时时间
	RetryCount      int               `json:"retry_count"`       // 重试次数
	ContentType     string            `json:"content_type"`      // 内容类型
	TemplateFormat  string            `json:"template_format"`   // 模板格式: json, form, custom
	CustomTemplate  string            `json:"custom_template"`   // 自定义模板
	BasicAuth       *BasicAuth        `json:"basic_auth"`        // Basic认证
	BearerToken     string            `json:"bearer_token"`      // Bearer Token
	EnableVerifySSL bool              `json:"enable_verify_ssl"` // 是否验证SSL
}

// BasicAuth Basic认证配置
type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// WebhookPayload Webhook负载
type WebhookPayload struct {
	Timestamp int64                  `json:"timestamp"`
	Alert     *WebhookAlert          `json:"alert"`
	Rule      *WebhookRule           `json:"rule"`
	Action    string                 `json:"action"`
	Status    string                 `json:"status"`
	Metadata  map[string]interface{} `json:"metadata"`
	Signature string                 `json:"signature,omitempty"`
}

// WebhookAlert Webhook告警信息
type WebhookAlert struct {
	ID              string                 `json:"id"`
	RuleName        string                 `json:"rule_name"`
	AlertLevel      string                 `json:"alert_level"`
	SeverityScore   int                    `json:"severity_score"`
	Message         string                 `json:"message"`
	Summary         string                 `json:"summary"`
	ResourceType    string                 `json:"resource_type"`
	ResourceID      int64                  `json:"resource_id"`
	ResourceName    string                 `json:"resource_name"`
	InstanceID      string                 `json:"instance_id"`
	TriggerValue    float64                `json:"trigger_value"`
	ThresholdValue  float64                `json:"threshold_value"`
	Condition       string                 `json:"condition"`
	Labels          map[string]interface{} `json:"labels"`
	Annotations     map[string]interface{} `json:"annotations"`
	TriggeredAt     time.Time              `json:"triggered_at"`
	FirstOccurrence time.Time              `json:"first_occurrence"`
	LastOccurrence  time.Time              `json:"last_occurrence"`
	OccurrenceCount int                    `json:"occurrence_count"`
	GroupID         string                 `json:"group_id"`
	CorrelationID   string                 `json:"correlation_id"`
}

// WebhookRule Webhook规则信息
type WebhookRule struct {
	ID                  int64                  `json:"id"`
	Name                string                 `json:"name"`
	DisplayName         string                 `json:"display_name"`
	Description         string                 `json:"description"`
	RuleType            string                 `json:"rule_type"`
	ConditionExpression string                 `json:"condition_expression"`
	QueryExpression     string                 `json:"query_expression"`
	WarningThreshold    float64                `json:"warning_threshold"`
	CriticalThreshold   float64                `json:"critical_threshold"`
	EvaluationWindow    int                    `json:"evaluation_window_seconds"`
	FilterLabels        map[string]interface{} `json:"filter_labels"`
	FilterResources     map[string]interface{} `json:"filter_resources"`
}

// NewWebhookChannel 创建Webhook通知渠道
func NewWebhookChannel(name string, configMap map[string]interface{}) *WebhookChannel {
	// 解析配置
	configData, _ := json.Marshal(configMap)
	var config WebhookConfig
	json.Unmarshal(configData, &config)

	// 设置默认值
	if config.Method == "" {
		config.Method = "POST"
	}
	if config.ContentType == "" {
		config.ContentType = "application/json"
	}
	if config.TemplateFormat == "" {
		config.TemplateFormat = "json"
	}
	if config.TimeoutSeconds == 0 {
		config.TimeoutSeconds = 30
	}
	if config.SignatureHeader == "" {
		config.SignatureHeader = "X-Webhook-Signature"
	}
	if config.SignaturePrefix == "" {
		config.SignaturePrefix = "sha256="
	}

	return &WebhookChannel{
		name:    name,
		config:  &config,
		enabled: true,
		logger:  logx.WithContext(context.Background()),
		client: &http.Client{
			Timeout: time.Duration(config.TimeoutSeconds) * time.Second,
		},
	}
}

// GetType 获取渠道类型
func (w *WebhookChannel) GetType() string {
	return "webhook"
}

// GetName 获取渠道名称
func (w *WebhookChannel) GetName() string {
	return w.name
}

// Send 发送Webhook通知
func (w *WebhookChannel) Send(ctx context.Context, req *NotificationRequest) error {
	if !w.enabled {
		return fmt.Errorf("Webhook渠道已禁用")
	}

	if w.config.URL == "" {
		return fmt.Errorf("Webhook URL未配置")
	}

	// 构建Webhook负载
	payload := w.buildWebhookPayload(req)

	// 生成请求体
	body, contentType, err := w.buildRequestBody(payload)
	if err != nil {
		return fmt.Errorf("构建请求体失败: %v", err)
	}

	// 发送Webhook请求
	return w.sendWebhookWithRetry(ctx, body, contentType)
}

// ValidateConfig 验证配置
func (w *WebhookChannel) ValidateConfig(config map[string]interface{}) error {
	if w.config.URL == "" {
		return fmt.Errorf("Webhook URL不能为空")
	}

	// 验证URL格式
	if _, err := url.Parse(w.config.URL); err != nil {
		return fmt.Errorf("Webhook URL格式无效: %v", err)
	}

	// 验证HTTP方法
	validMethods := []string{"GET", "POST", "PUT", "PATCH"}
	methodValid := false
	for _, method := range validMethods {
		if strings.ToUpper(w.config.Method) == method {
			methodValid = true
			break
		}
	}
	if !methodValid {
		return fmt.Errorf("不支持的HTTP方法: %s", w.config.Method)
	}

	return nil
}

// IsEnabled 检查是否启用
func (w *WebhookChannel) IsEnabled() bool {
	return w.enabled
}

// buildWebhookPayload 构建Webhook负载
func (w *WebhookChannel) buildWebhookPayload(req *NotificationRequest) *WebhookPayload {
	// 构建告警信息
	alert := &WebhookAlert{
		ID:              req.Alert.ID,
		RuleName:        req.Alert.RuleName,
		AlertLevel:      req.Alert.AlertLevel,
		SeverityScore:   req.Alert.SeverityScore,
		Message:         req.Alert.Message,
		Summary:         req.Alert.Summary,
		ResourceType:    req.Alert.ResourceType,
		ResourceID:      req.Alert.ResourceID,
		ResourceName:    req.Alert.ResourceName,
		InstanceID:      req.Alert.InstanceID,
		TriggerValue:    req.Alert.TriggerValue,
		ThresholdValue:  req.Alert.ThresholdValue,
		Condition:       req.Alert.ConditionExpression,
		Labels:          req.Alert.Labels,
		Annotations:     req.Alert.Annotations,
		TriggeredAt:     req.Alert.TriggeredAt,
		FirstOccurrence: req.Alert.FirstOccurrenceAt,
		LastOccurrence:  req.Alert.LastOccurrenceAt,
		OccurrenceCount: req.Alert.OccurrenceCount,
		GroupID:         req.Alert.AlertGroupID,
		CorrelationID:   req.Alert.CorrelationID,
	}

	// 构建规则信息
	rule := &WebhookRule{
		ID:                  req.RuleInfo.ID,
		Name:                req.RuleInfo.Name,
		DisplayName:         req.RuleInfo.DisplayName,
		Description:         req.RuleInfo.Description,
		RuleType:            req.RuleInfo.RuleType,
		ConditionExpression: req.RuleInfo.ConditionExpression,
		QueryExpression:     req.RuleInfo.QueryExpression,
		WarningThreshold:    req.RuleInfo.WarningThreshold,
		CriticalThreshold:   req.RuleInfo.CriticalThreshold,
		EvaluationWindow:    req.RuleInfo.EvaluationWindowSeconds,
		FilterLabels:        req.RuleInfo.FilterLabels,
		FilterResources:     req.RuleInfo.FilterResources,
	}

	// 构建负载
	payload := &WebhookPayload{
		Timestamp: time.Now().Unix(),
		Alert:     alert,
		Rule:      rule,
		Action:    req.Action,
		Status:    req.Alert.Status,
		Metadata:  req.Metadata,
	}

	return payload
}

// buildRequestBody 构建请求体
func (w *WebhookChannel) buildRequestBody(payload *WebhookPayload) ([]byte, string, error) {
	switch w.config.TemplateFormat {
	case "json":
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, "", err
		}

		// 生成签名
		if w.config.Secret != "" {
			signature := w.generateSignature(data, w.config.Secret)
			payload.Signature = signature

			// 重新序列化包含签名的负载
			data, err = json.Marshal(payload)
			if err != nil {
				return nil, "", err
			}
		}

		return data, "application/json", nil

	case "form":
		return w.buildFormData(payload)

	case "custom":
		return w.buildCustomTemplate(payload)

	default:
		return nil, "", fmt.Errorf("不支持的模板格式: %s", w.config.TemplateFormat)
	}
}

// buildFormData 构建表单数据
func (w *WebhookChannel) buildFormData(payload *WebhookPayload) ([]byte, string, error) {
	values := url.Values{}

	// 扁平化负载数据
	values.Set("timestamp", strconv.FormatInt(payload.Timestamp, 10))
	values.Set("action", payload.Action)
	values.Set("status", payload.Status)
	values.Set("alert_id", payload.Alert.ID)
	values.Set("rule_name", payload.Alert.RuleName)
	values.Set("alert_level", payload.Alert.AlertLevel)
	values.Set("message", payload.Alert.Message)
	values.Set("resource_type", payload.Alert.ResourceType)
	values.Set("resource_name", payload.Alert.ResourceName)
	values.Set("trigger_value", fmt.Sprintf("%.2f", payload.Alert.TriggerValue))
	values.Set("threshold_value", fmt.Sprintf("%.2f", payload.Alert.ThresholdValue))
	values.Set("triggered_at", payload.Alert.TriggeredAt.Format(time.RFC3339))

	data := values.Encode()

	// 生成签名
	if w.config.Secret != "" {
		signature := w.generateSignature([]byte(data), w.config.Secret)
		values.Set("signature", signature)
		data = values.Encode()
	}

	return []byte(data), "application/x-www-form-urlencoded", nil
}

// buildCustomTemplate 构建自定义模板
func (w *WebhookChannel) buildCustomTemplate(payload *WebhookPayload) ([]byte, string, error) {
	if w.config.CustomTemplate == "" {
		return nil, "", fmt.Errorf("自定义模板未配置")
	}

	// 简化的模板替换（实际应使用模板引擎）
	template := w.config.CustomTemplate
	template = strings.ReplaceAll(template, "{{.Alert.ID}}", payload.Alert.ID)
	template = strings.ReplaceAll(template, "{{.Alert.RuleName}}", payload.Alert.RuleName)
	template = strings.ReplaceAll(template, "{{.Alert.AlertLevel}}", payload.Alert.AlertLevel)
	template = strings.ReplaceAll(template, "{{.Alert.Message}}", payload.Alert.Message)
	template = strings.ReplaceAll(template, "{{.Action}}", payload.Action)
	template = strings.ReplaceAll(template, "{{.Timestamp}}", strconv.FormatInt(payload.Timestamp, 10))

	data := []byte(template)

	// 生成签名
	if w.config.Secret != "" {
		signature := w.generateSignature(data, w.config.Secret)
		// 在模板中替换签名占位符
		template = strings.ReplaceAll(template, "{{.Signature}}", signature)
		data = []byte(template)
	}

	return data, w.config.ContentType, nil
}

// generateSignature 生成签名
func (w *WebhookChannel) generateSignature(data []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// sendWebhookWithRetry 发送Webhook请求（带重试）
func (w *WebhookChannel) sendWebhookWithRetry(ctx context.Context, body []byte, contentType string) error {
	maxRetries := w.config.RetryCount
	if maxRetries <= 0 {
		maxRetries = 1
	}

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// 指数退避
			backoff := time.Duration(attempt*attempt) * time.Second
			time.Sleep(backoff)
			w.logger.Infof("重试Webhook请求，第%d次尝试", attempt+1)
		}

		err := w.sendWebhook(ctx, body, contentType)
		if err == nil {
			w.logger.Infof("Webhook发送成功")
			return nil
		}

		lastErr = err
		w.logger.Errorf("Webhook发送失败，第%d次尝试: %v", attempt+1, err)
	}

	return fmt.Errorf("Webhook发送最终失败，已重试%d次: %v", maxRetries, lastErr)
}

// sendWebhook 发送Webhook请求
func (w *WebhookChannel) sendWebhook(ctx context.Context, body []byte, contentType string) error {
	req, err := http.NewRequestWithContext(ctx, w.config.Method, w.config.URL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置Content-Type
	req.Header.Set("Content-Type", contentType)

	// 设置自定义头
	for key, value := range w.config.Headers {
		req.Header.Set(key, value)
	}

	// 设置签名头
	if w.config.Secret != "" {
		signature := w.generateSignature(body, w.config.Secret)
		req.Header.Set(w.config.SignatureHeader, w.config.SignaturePrefix+signature)
	}

	// 设置认证
	if w.config.BasicAuth != nil {
		req.SetBasicAuth(w.config.BasicAuth.Username, w.config.BasicAuth.Password)
	}
	if w.config.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+w.config.BearerToken)
	}

	// 设置User-Agent
	req.Header.Set("User-Agent", "VolcTrain-Monitor/1.0")

	// 发送请求
	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	w.logger.Infof("Webhook请求成功，状态码: %d", resp.StatusCode)
	return nil
}

// GetWebhookStatus 获取Webhook渠道状态
func (w *WebhookChannel) GetWebhookStatus() map[string]interface{} {
	return map[string]interface{}{
		"channel_name":      w.name,
		"enabled":           w.enabled,
		"url":               w.config.URL,
		"method":            w.config.Method,
		"content_type":      w.config.ContentType,
		"template_format":   w.config.TemplateFormat,
		"timeout_seconds":   w.config.TimeoutSeconds,
		"retry_count":       w.config.RetryCount,
		"enable_verify_ssl": w.config.EnableVerifySSL,
		"has_secret":        w.config.Secret != "",
		"has_basic_auth":    w.config.BasicAuth != nil,
		"has_bearer_token":  w.config.BearerToken != "",
	}
}
