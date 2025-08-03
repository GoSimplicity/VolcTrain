package notification

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// EmailChannel 邮件通知渠道
type EmailChannel struct {
	name    string
	config  *EmailConfig
	enabled bool
	logger  logx.Logger
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTPHost           string   `json:"smtp_host"`
	SMTPPort           int      `json:"smtp_port"`
	Username           string   `json:"username"`
	Password           string   `json:"password"`
	From               string   `json:"from"`
	FromName           string   `json:"from_name"`
	UseTLS             bool     `json:"use_tls"`
	UseSTARTTLS        bool     `json:"use_starttls"`
	Recipients         []string `json:"recipients"`
	CriticalRecipients []string `json:"critical_recipients"`
}

// NewEmailChannel 创建邮件通知渠道
func NewEmailChannel(name string, configMap map[string]interface{}) *EmailChannel {
	// 解析配置
	configData, _ := json.Marshal(configMap)
	var config EmailConfig
	json.Unmarshal(configData, &config)

	return &EmailChannel{
		name:    name,
		config:  &config,
		enabled: true,
		logger:  logx.WithContext(context.Background()),
	}
}

// GetType 获取渠道类型
func (e *EmailChannel) GetType() string {
	return "email"
}

// GetName 获取渠道名称
func (e *EmailChannel) GetName() string {
	return e.name
}

// Send 发送邮件通知
func (e *EmailChannel) Send(ctx context.Context, req *NotificationRequest) error {
	if !e.enabled {
		return fmt.Errorf("邮件渠道已禁用")
	}

	if len(req.Recipients) == 0 {
		return fmt.Errorf("收件人列表为空")
	}

	// 构建邮件内容
	message := e.buildEmailMessage(req)

	// 发送邮件
	return e.sendSMTP(req.Recipients, message)
}

// ValidateConfig 验证配置
func (e *EmailChannel) ValidateConfig(config map[string]interface{}) error {
	if e.config.SMTPHost == "" {
		return fmt.Errorf("SMTP主机不能为空")
	}
	if e.config.SMTPPort == 0 {
		return fmt.Errorf("SMTP端口不能为空")
	}
	if e.config.Username == "" {
		return fmt.Errorf("SMTP用户名不能为空")
	}
	if e.config.Password == "" {
		return fmt.Errorf("SMTP密码不能为空")
	}
	if e.config.From == "" {
		return fmt.Errorf("发件人地址不能为空")
	}
	return nil
}

// IsEnabled 检查是否启用
func (e *EmailChannel) IsEnabled() bool {
	return e.enabled
}

// buildEmailMessage 构建邮件消息
func (e *EmailChannel) buildEmailMessage(req *NotificationRequest) string {
	fromName := e.config.FromName
	if fromName == "" {
		fromName = "VolcTrain监控系统"
	}

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", fromName, e.config.From)
	headers["To"] = strings.Join(req.Recipients, ", ")
	headers["Subject"] = req.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"
	headers["Date"] = time.Now().Format(time.RFC1123Z)

	// 构建邮件头
	var message strings.Builder
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")

	// 构建HTML邮件体
	htmlContent := e.buildHTMLContent(req)
	message.WriteString(htmlContent)

	return message.String()
}

// buildHTMLContent 构建HTML邮件内容
func (e *EmailChannel) buildHTMLContent(req *NotificationRequest) string {
	// 根据告警级别选择颜色
	var levelColor, levelBgColor string
	switch req.Priority {
	case "critical":
		levelColor = "#721c24"
		levelBgColor = "#f8d7da"
	case "warning":
		levelColor = "#856404"
		levelBgColor = "#fff3cd"
	default:
		levelColor = "#155724"
		levelBgColor = "#d4edda"
	}

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>%s</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 600px; margin: 0 auto; background-color: white; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { background-color: #007bff; color: white; padding: 20px; border-radius: 8px 8px 0 0; }
        .content { padding: 20px; }
        .alert-level { display: inline-block; padding: 4px 8px; border-radius: 4px; font-weight: bold; color: %s; background-color: %s; }
        .details { margin: 20px 0; }
        .details table { width: 100%%; border-collapse: collapse; }
        .details th, .details td { padding: 8px; text-align: left; border-bottom: 1px solid #ddd; }
        .details th { background-color: #f8f9fa; font-weight: bold; }
        .footer { background-color: #f8f9fa; padding: 15px 20px; border-radius: 0 0 8px 8px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>🚨 VolcTrain 监控告警</h2>
            <p>告警通知 - %s</p>
        </div>
        <div class="content">
            <div class="alert-level">%s级别告警</div>
            <h3>%s</h3>
            <p>%s</p>
            
            <div class="details">
                <h4>告警详情</h4>
                <table>
                    <tr><th>告警规则</th><td>%s</td></tr>
                    <tr><th>资源类型</th><td>%s</td></tr>
                    <tr><th>资源名称</th><td>%s</td></tr>
                    <tr><th>触发值</th><td>%.2f</td></tr>
                    <tr><th>阈值</th><td>%.2f</td></tr>
                    <tr><th>触发时间</th><td>%s</td></tr>
                    <tr><th>持续时间</th><td>%d 分钟</td></tr>
                </table>
            </div>
        </div>
        <div class="footer">
            <p>此邮件由 VolcTrain 监控系统自动发送，请勿直接回复。</p>
            <p>发送时间: %s</p>
        </div>
    </div>
</body>
</html>`,
		req.Subject,
		levelColor, levelBgColor,
		req.Action,
		req.Priority,
		req.Alert.RuleName,
		req.Content,
		req.Alert.RuleName,
		req.Alert.ResourceType,
		req.Alert.ResourceName,
		req.Alert.TriggerValue,
		req.Alert.ThresholdValue,
		req.Alert.TriggeredAt.Format("2006-01-02 15:04:05"),
		int(time.Since(req.Alert.TriggeredAt).Minutes()),
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return html
}

// sendSMTP 发送SMTP邮件
func (e *EmailChannel) sendSMTP(to []string, message string) error {
	auth := smtp.PlainAuth("", e.config.Username, e.config.Password, e.config.SMTPHost)
	addr := fmt.Sprintf("%s:%d", e.config.SMTPHost, e.config.SMTPPort)

	if e.config.UseTLS {
		return e.sendSMTPTLS(to, message, auth, addr)
	}

	return smtp.SendMail(addr, auth, e.config.From, to, []byte(message))
}

// sendSMTPTLS 发送TLS加密邮件
func (e *EmailChannel) sendSMTPTLS(to []string, message string, auth smtp.Auth, addr string) error {
	// 建立TLS连接
	tlsConfig := &tls.Config{
		ServerName: e.config.SMTPHost,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, e.config.SMTPHost)
	if err != nil {
		return err
	}
	defer client.Quit()

	// 认证
	if err := client.Auth(auth); err != nil {
		return err
	}

	// 设置发件人
	if err := client.Mail(e.config.From); err != nil {
		return err
	}

	// 设置收件人
	for _, addr := range to {
		if err := client.Rcpt(addr); err != nil {
			return err
		}
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.Write([]byte(message))
	return err
}
