package notification

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// DingTalkChannel 钉钉通知渠道
type DingTalkChannel struct {
	name    string
	config  *DingTalkConfig
	enabled bool
	logger  logx.Logger
	client  *http.Client
}

// DingTalkConfig 钉钉配置
type DingTalkConfig struct {
	WebhookURL   string   `json:"webhook_url"`   // 钉钉机器人Webhook地址
	Secret       string   `json:"secret"`        // 钉钉机器人密钥
	AtMobiles    []string `json:"at_mobiles"`    // @的手机号列表
	AtUserIds    []string `json:"at_user_ids"`   // @的用户ID列表
	IsAtAll      bool     `json:"is_at_all"`     // 是否@所有人
	EnableSecret bool     `json:"enable_secret"` // 是否启用加密
	MessageType  string   `json:"message_type"`  // 消息类型: text, markdown, actionCard
}

// DingTalkMessage 钉钉消息
type DingTalkMessage struct {
	MsgType    string              `json:"msgtype"`
	Text       *DingTalkText       `json:"text,omitempty"`
	Markdown   *DingTalkMarkdown   `json:"markdown,omitempty"`
	ActionCard *DingTalkActionCard `json:"actionCard,omitempty"`
	At         *DingTalkAt         `json:"at,omitempty"`
}

// DingTalkText 文本消息
type DingTalkText struct {
	Content string `json:"content"`
}

// DingTalkMarkdown Markdown消息
type DingTalkMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// DingTalkActionCard ActionCard消息
type DingTalkActionCard struct {
	Title          string `json:"title"`
	Text           string `json:"text"`
	HideAvatar     string `json:"hideAvatar"`
	BtnOrientation string `json:"btnOrientation"`
	SingleTitle    string `json:"singleTitle"`
	SingleURL      string `json:"singleURL"`
}

// DingTalkAt @信息
type DingTalkAt struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	AtUserIds []string `json:"atUserIds,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}

// DingTalkResponse 钉钉响应
type DingTalkResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// NewDingTalkChannel 创建钉钉通知渠道
func NewDingTalkChannel(name string, configMap map[string]interface{}) *DingTalkChannel {
	// 解析配置
	configData, _ := json.Marshal(configMap)
	var config DingTalkConfig
	json.Unmarshal(configData, &config)

	// 设置默认值
	if config.MessageType == "" {
		config.MessageType = "markdown"
	}

	return &DingTalkChannel{
		name:    name,
		config:  &config,
		enabled: true,
		logger:  logx.WithContext(context.Background()),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetType 获取渠道类型
func (d *DingTalkChannel) GetType() string {
	return "dingtalk"
}

// GetName 获取渠道名称
func (d *DingTalkChannel) GetName() string {
	return d.name
}

// Send 发送钉钉通知
func (d *DingTalkChannel) Send(ctx context.Context, req *NotificationRequest) error {
	if !d.enabled {
		return fmt.Errorf("钉钉渠道已禁用")
	}

	if d.config.WebhookURL == "" {
		return fmt.Errorf("钉钉Webhook地址未配置")
	}

	// 构建钉钉消息
	message := d.buildDingTalkMessage(req)

	// 构建请求URL
	requestURL := d.buildRequestURL()

	// 发送消息
	return d.sendMessage(ctx, requestURL, message)
}

// ValidateConfig 验证配置
func (d *DingTalkChannel) ValidateConfig(config map[string]interface{}) error {
	if d.config.WebhookURL == "" {
		return fmt.Errorf("钉钉Webhook地址不能为空")
	}

	// 验证Webhook URL格式
	if _, err := url.Parse(d.config.WebhookURL); err != nil {
		return fmt.Errorf("钉钉Webhook地址格式无效: %v", err)
	}

	return nil
}

// IsEnabled 检查是否启用
func (d *DingTalkChannel) IsEnabled() bool {
	return d.enabled
}

// buildDingTalkMessage 构建钉钉消息
func (d *DingTalkChannel) buildDingTalkMessage(req *NotificationRequest) *DingTalkMessage {
	message := &DingTalkMessage{
		At: &DingTalkAt{
			AtMobiles: d.config.AtMobiles,
			AtUserIds: d.config.AtUserIds,
			IsAtAll:   d.config.IsAtAll,
		},
	}

	switch d.config.MessageType {
	case "text":
		message.MsgType = "text"
		message.Text = &DingTalkText{
			Content: d.buildTextContent(req),
		}
	case "actionCard":
		message.MsgType = "actionCard"
		message.ActionCard = d.buildActionCard(req)
	default: // markdown
		message.MsgType = "markdown"
		message.Markdown = d.buildMarkdown(req)
	}

	return message
}

// buildTextContent 构建文本内容
func (d *DingTalkChannel) buildTextContent(req *NotificationRequest) string {
	var action string
	switch req.Action {
	case "firing":
		action = "🚨 告警触发"
	case "resolved":
		action = "✅ 告警恢复"
	default:
		action = "📢 告警通知"
	}

	content := fmt.Sprintf(`%s

规则名称：%s
告警级别：%s
资源类型：%s
资源名称：%s
当前值：%.2f
阈值：%.2f
触发时间：%s
持续时间：%d分钟

告警详情：%s`,
		action,
		req.Alert.RuleName,
		req.Priority,
		req.Alert.ResourceType,
		req.Alert.ResourceName,
		req.Alert.TriggerValue,
		req.Alert.ThresholdValue,
		req.Alert.TriggeredAt.Format("2006-01-02 15:04:05"),
		int(time.Since(req.Alert.TriggeredAt).Minutes()),
		req.Alert.Message,
	)

	// 添加@信息
	if len(d.config.AtMobiles) > 0 {
		content += "\n\n@" + fmt.Sprintf("%v", d.config.AtMobiles)
	}

	return content
}

// buildMarkdown 构建Markdown消息
func (d *DingTalkChannel) buildMarkdown(req *NotificationRequest) *DingTalkMarkdown {
	var title, emoji, color string
	switch req.Priority {
	case "critical":
		title = "🚨 严重告警"
		emoji = "🔴"
		color = "#FF4D4F"
	case "warning":
		title = "⚠️ 警告告警"
		emoji = "🟡"
		color = "#FAAD14"
	default:
		title = "ℹ️ 信息告警"
		emoji = "🔵"
		color = "#1890FF"
	}

	var actionText string
	switch req.Action {
	case "firing":
		actionText = "告警触发"
	case "resolved":
		actionText = "告警恢复"
	default:
		actionText = "告警通知"
	}

	markdownText := fmt.Sprintf(`## %s %s

---

### 📋 告警信息
- **规则名称**：%s
- **告警级别**：<font color="%s">%s %s</font>
- **动作类型**：%s
- **触发时间**：%s

### 🎯 资源详情
- **资源类型**：%s
- **资源名称**：%s
- **实例ID**：%s

### 📊 指标详情
- **当前值**：**%.2f**
- **阈值**：%.2f
- **条件表达式**：%s

### 📝 告警描述
%s

---
> 📅 **发送时间**：%s  
> 🤖 **来源**：VolcTrain监控系统`,
		emoji, title,
		req.Alert.RuleName,
		color, emoji, req.Priority,
		actionText,
		req.Alert.TriggeredAt.Format("2006-01-02 15:04:05"),
		req.Alert.ResourceType,
		req.Alert.ResourceName,
		req.Alert.InstanceID,
		req.Alert.TriggerValue,
		req.Alert.ThresholdValue,
		req.Alert.ConditionExpression,
		req.Alert.Message,
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return &DingTalkMarkdown{
		Title: title,
		Text:  markdownText,
	}
}

// buildActionCard 构建ActionCard消息
func (d *DingTalkChannel) buildActionCard(req *NotificationRequest) *DingTalkActionCard {
	var title, color string
	switch req.Priority {
	case "critical":
		title = "🚨 严重告警"
		color = "#FF4D4F"
	case "warning":
		title = "⚠️ 警告告警"
		color = "#FAAD14"
	default:
		title = "ℹ️ 信息告警"
		color = "#1890FF"
	}

	text := fmt.Sprintf(`<div style="border-left: 4px solid %s; padding-left: 10px;">
<h3>%s</h3>
<p><strong>规则名称：</strong>%s</p>
<p><strong>告警级别：</strong><span style="color: %s;">%s</span></p>
<p><strong>当前值：</strong>%.2f，<strong>阈值：</strong>%.2f</p>
<p><strong>触发时间：</strong>%s</p>
<p><strong>资源：</strong>%s/%s</p>
<p><strong>描述：</strong>%s</p>
</div>`,
		color,
		title,
		req.Alert.RuleName,
		color, req.Priority,
		req.Alert.TriggerValue, req.Alert.ThresholdValue,
		req.Alert.TriggeredAt.Format("2006-01-02 15:04:05"),
		req.Alert.ResourceType, req.Alert.ResourceName,
		req.Alert.Message,
	)

	return &DingTalkActionCard{
		Title:          title,
		Text:           text,
		HideAvatar:     "0",
		BtnOrientation: "0",
		SingleTitle:    "查看详情",
		SingleURL:      "https://your-monitoring-dashboard.com", // 可配置的监控面板链接
	}
}

// buildRequestURL 构建请求URL
func (d *DingTalkChannel) buildRequestURL() string {
	requestURL := d.config.WebhookURL

	// 如果启用了加密，添加时间戳和签名
	if d.config.EnableSecret && d.config.Secret != "" {
		timestamp := time.Now().UnixNano() / 1e6
		sign := d.generateSign(timestamp, d.config.Secret)

		requestURL += fmt.Sprintf("&timestamp=%d&sign=%s", timestamp, url.QueryEscape(sign))
	}

	return requestURL
}

// generateSign 生成钉钉签名
func (d *DingTalkChannel) generateSign(timestamp int64, secret string) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// sendMessage 发送消息
func (d *DingTalkChannel) sendMessage(ctx context.Context, url string, message *DingTalkMessage) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化钉钉消息失败: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %v", err)
	}

	var response DingTalkResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if response.ErrCode != 0 {
		return fmt.Errorf("钉钉API返回错误: %d - %s", response.ErrCode, response.ErrMsg)
	}

	d.logger.Infof("钉钉消息发送成功")
	return nil
}

// GetDingTalkStatus 获取钉钉渠道状态
func (d *DingTalkChannel) GetDingTalkStatus() map[string]interface{} {
	return map[string]interface{}{
		"channel_name":  d.name,
		"enabled":       d.enabled,
		"message_type":  d.config.MessageType,
		"enable_secret": d.config.EnableSecret,
		"at_mobiles":    d.config.AtMobiles,
		"at_user_ids":   d.config.AtUserIds,
		"is_at_all":     d.config.IsAtAll,
	}
}
