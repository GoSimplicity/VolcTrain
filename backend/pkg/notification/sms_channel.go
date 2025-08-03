package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// SMSChannel 短信通知渠道
type SMSChannel struct {
	name    string
	config  *SMSConfig
	enabled bool
	logger  logx.Logger
	client  *http.Client
}

// SMSConfig 短信配置
type SMSConfig struct {
	Provider           string   `json:"provider"` // 短信服务商: aliyun, tencent, huawei
	AccessKeyID        string   `json:"access_key_id"`
	AccessKeySecret    string   `json:"access_key_secret"`
	SignName           string   `json:"sign_name"`           // 短信签名
	TemplateCode       string   `json:"template_code"`       // 短信模板代码
	Region             string   `json:"region"`              // 区域
	Recipients         []string `json:"recipients"`          // 默认收件人
	CriticalRecipients []string `json:"critical_recipients"` // 紧急告警收件人
	APIEndpoint        string   `json:"api_endpoint"`        // API端点
}

// AliyunSMSRequest 阿里云短信请求
type AliyunSMSRequest struct {
	Action           string `json:"Action"`
	Format           string `json:"Format"`
	Version          string `json:"Version"`
	AccessKeyId      string `json:"AccessKeyId"`
	SignatureMethod  string `json:"SignatureMethod"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	SignatureNonce   string `json:"SignatureNonce"`
	PhoneNumbers     string `json:"PhoneNumbers"`
	SignName         string `json:"SignName"`
	TemplateCode     string `json:"TemplateCode"`
	TemplateParam    string `json:"TemplateParam"`
	Signature        string `json:"Signature"`
}

// TencentSMSRequest 腾讯云短信请求
type TencentSMSRequest struct {
	PhoneNumberSet   []string `json:"PhoneNumberSet"`
	TemplateID       string   `json:"TemplateID"`
	Sign             string   `json:"Sign"`
	TemplateParamSet []string `json:"TemplateParamSet"`
	SessionContext   string   `json:"SessionContext"`
}

// NewSMSChannel 创建短信通知渠道
func NewSMSChannel(name string, configMap map[string]interface{}) *SMSChannel {
	// 解析配置
	configData, _ := json.Marshal(configMap)
	var config SMSConfig
	json.Unmarshal(configData, &config)

	// 设置默认值
	if config.Provider == "" {
		config.Provider = "aliyun"
	}

	return &SMSChannel{
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
func (s *SMSChannel) GetType() string {
	return "sms"
}

// GetName 获取渠道名称
func (s *SMSChannel) GetName() string {
	return s.name
}

// Send 发送短信通知
func (s *SMSChannel) Send(ctx context.Context, req *NotificationRequest) error {
	if !s.enabled {
		return fmt.Errorf("短信渠道已禁用")
	}

	if len(req.Recipients) == 0 {
		return fmt.Errorf("收件人列表为空")
	}

	// 验证手机号格式
	validRecipients := s.validatePhoneNumbers(req.Recipients)
	if len(validRecipients) == 0 {
		return fmt.Errorf("没有有效的手机号码")
	}

	// 构建短信内容
	content := s.buildSMSContent(req)

	// 根据不同服务商发送短信
	switch s.config.Provider {
	case "aliyun":
		return s.sendAliyunSMS(ctx, validRecipients, content, req)
	case "tencent":
		return s.sendTencentSMS(ctx, validRecipients, content, req)
	default:
		return fmt.Errorf("不支持的短信服务商: %s", s.config.Provider)
	}
}

// ValidateConfig 验证配置
func (s *SMSChannel) ValidateConfig(config map[string]interface{}) error {
	if s.config.Provider == "" {
		return fmt.Errorf("短信服务商不能为空")
	}
	if s.config.AccessKeyID == "" {
		return fmt.Errorf("AccessKeyID不能为空")
	}
	if s.config.AccessKeySecret == "" {
		return fmt.Errorf("AccessKeySecret不能为空")
	}
	if s.config.SignName == "" {
		return fmt.Errorf("短信签名不能为空")
	}
	if s.config.TemplateCode == "" {
		return fmt.Errorf("短信模板代码不能为空")
	}
	return nil
}

// IsEnabled 检查是否启用
func (s *SMSChannel) IsEnabled() bool {
	return s.enabled
}

// validatePhoneNumbers 验证手机号格式
func (s *SMSChannel) validatePhoneNumbers(phones []string) []string {
	var validPhones []string
	for _, phone := range phones {
		// 简化的手机号验证（实际应使用正则表达式）
		if len(phone) >= 11 && strings.HasPrefix(phone, "1") {
			validPhones = append(validPhones, phone)
		} else if strings.HasPrefix(phone, "+86") && len(phone) == 14 {
			validPhones = append(validPhones, phone[3:]) // 去掉+86前缀
		}
	}
	return validPhones
}

// buildSMSContent 构建短信内容
func (s *SMSChannel) buildSMSContent(req *NotificationRequest) string {
	// 短信内容需要简洁明了
	var action string
	switch req.Action {
	case "firing":
		action = "触发"
	case "resolved":
		action = "恢复"
	default:
		action = req.Action
	}

	content := fmt.Sprintf("【%s】告警%s：%s，级别：%s，当前值：%.2f，阈值：%.2f，时间：%s",
		s.config.SignName,
		action,
		req.Alert.RuleName,
		req.Priority,
		req.Alert.TriggerValue,
		req.Alert.ThresholdValue,
		req.Alert.TriggeredAt.Format("15:04"),
	)

	// 限制短信长度（一般不超过70字符）
	if len(content) > 67 {
		content = content[:67] + "..."
	}

	return content
}

// sendAliyunSMS 发送阿里云短信
func (s *SMSChannel) sendAliyunSMS(ctx context.Context, phones []string, content string, req *NotificationRequest) error {
	// 构建模板参数
	templateParams := map[string]string{
		"rule_name":       req.Alert.RuleName,
		"alert_level":     req.Priority,
		"trigger_value":   fmt.Sprintf("%.2f", req.Alert.TriggerValue),
		"threshold_value": fmt.Sprintf("%.2f", req.Alert.ThresholdValue),
		"action":          req.Action,
		"time":            req.Alert.TriggeredAt.Format("15:04"),
	}

	templateParamJSON, _ := json.Marshal(templateParams)

	// 构建请求参数
	params := map[string]string{
		"Action":           "SendSms",
		"Format":           "JSON",
		"Version":          "2017-05-25",
		"AccessKeyId":      s.config.AccessKeyID,
		"SignatureMethod":  "HMAC-SHA1",
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"SignatureVersion": "1.0",
		"SignatureNonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
		"PhoneNumbers":     strings.Join(phones, ","),
		"SignName":         s.config.SignName,
		"TemplateCode":     s.config.TemplateCode,
		"TemplateParam":    string(templateParamJSON),
	}

	// 生成签名（简化实现，实际需要按阿里云规范生成签名）
	signature := s.generateAliyunSignature(params)
	params["Signature"] = signature

	// 构建请求URL
	endpoint := "https://dysmsapi.aliyuncs.com/"
	if s.config.APIEndpoint != "" {
		endpoint = s.config.APIEndpoint
	}

	// 构建POST请求体
	formData := url.Values{}
	for k, v := range params {
		formData.Set(k, v)
	}

	req_, err := http.NewRequestWithContext(ctx, "POST", endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}

	req_.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	resp, err := s.client.Do(req_)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 解析响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	// 检查发送结果
	if code, ok := result["Code"].(string); ok && code != "OK" {
		message, _ := result["Message"].(string)
		return fmt.Errorf("阿里云短信发送失败: %s - %s", code, message)
	}

	s.logger.Infof("阿里云短信发送成功，收件人: %v", phones)
	return nil
}

// sendTencentSMS 发送腾讯云短信
func (s *SMSChannel) sendTencentSMS(ctx context.Context, phones []string, content string, req *NotificationRequest) error {
	// 构建模板参数
	templateParams := []string{
		req.Alert.RuleName,
		req.Priority,
		fmt.Sprintf("%.2f", req.Alert.TriggerValue),
		fmt.Sprintf("%.2f", req.Alert.ThresholdValue),
		req.Action,
		req.Alert.TriggeredAt.Format("15:04"),
	}

	// 构建请求体
	requestBody := TencentSMSRequest{
		PhoneNumberSet:   phones,
		TemplateID:       s.config.TemplateCode,
		Sign:             s.config.SignName,
		TemplateParamSet: templateParams,
		SessionContext:   req.ID,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	// 构建请求
	endpoint := "https://sms.tencentcloudapi.com/"
	if s.config.APIEndpoint != "" {
		endpoint = s.config.APIEndpoint
	}

	req_, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// 设置请求头（简化实现，实际需要按腾讯云规范设置认证头）
	req_.Header.Set("Content-Type", "application/json")
	req_.Header.Set("X-TC-Action", "SendSms")
	req_.Header.Set("X-TC-Version", "2019-07-11")
	req_.Header.Set("X-TC-Region", s.config.Region)

	// 发送请求
	resp, err := s.client.Do(req_)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 解析响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	// 检查发送结果
	if errorResp, ok := result["Error"].(map[string]interface{}); ok {
		code, _ := errorResp["Code"].(string)
		message, _ := errorResp["Message"].(string)
		return fmt.Errorf("腾讯云短信发送失败: %s - %s", code, message)
	}

	s.logger.Infof("腾讯云短信发送成功，收件人: %v", phones)
	return nil
}

// generateAliyunSignature 生成阿里云API签名（简化实现）
func (s *SMSChannel) generateAliyunSignature(params map[string]string) string {
	// 实际实现需要按照阿里云API签名规范
	// 这里为了简化，返回一个假的签名
	// 实际项目中需要实现完整的HMAC-SHA1签名算法
	return "fake_signature_for_demo"
}

// GetSMSStatus 获取短信发送状态
func (s *SMSChannel) GetSMSStatus() map[string]interface{} {
	return map[string]interface{}{
		"channel_name":  s.name,
		"provider":      s.config.Provider,
		"enabled":       s.enabled,
		"sign_name":     s.config.SignName,
		"template_code": s.config.TemplateCode,
	}
}
