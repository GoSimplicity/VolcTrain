package monitoring

import (
	"context"
	"database/sql"
	"encoding/json"

	"api/model"
	"github.com/zeromicro/go-zero/core/logx"
)

// BuiltinRulesTemplates 内置规则和模板管理器
type BuiltinRulesTemplates struct {
	logger                     logx.Logger
	db                         *sql.DB
	alertRulesModel            model.VtAlertRulesModel
	notificationTemplatesModel model.VtNotificationTemplatesModel
}

// NewBuiltinRulesTemplates 创建内置规则模板管理器
func NewBuiltinRulesTemplates(db *sql.DB) *BuiltinRulesTemplates {
	return &BuiltinRulesTemplates{
		logger:                     logx.WithContext(context.Background()),
		db:                         db,
		alertRulesModel:            model.NewVtAlertRulesModel(db),
		notificationTemplatesModel: model.NewVtNotificationTemplatesModel(db),
	}
}

// InitializeBuiltinRules 初始化内置告警规则
func (b *BuiltinRulesTemplates) InitializeBuiltinRules() error {
	b.logger.Info("初始化内置告警规则")

	builtinRules := []*model.VtAlertRules{
		// 系统资源告警规则
		{
			Name:                        "system_cpu_high_usage",
			DisplayName:                 "系统CPU使用率过高",
			Description:                 "当系统CPU使用率超过80%时触发告警",
			RuleType:                    "threshold",
			ConditionExpression:         "system_cpu_usage > 80",
			QueryExpression:             "system_cpu_usage",
			WarningThreshold:            80.0,
			CriticalThreshold:           90.0,
			ThresholdCondition:          "gt",
			EvaluationWindowSeconds:     300,
			EvaluationIntervalSeconds:   60,
			TriggerDurationSeconds:      120,
			RecoveryDurationSeconds:     60,
			AlertLevel:                  "warning",
			SeverityScore:               70,
			NotificationThrottleMinutes: 5,
			Status:                      "active",
			IsBuiltin:                   true,
		},
		{
			Name:                        "system_memory_high_usage",
			DisplayName:                 "系统内存使用率过高",
			Description:                 "当系统内存使用率超过85%时触发告警",
			RuleType:                    "threshold",
			ConditionExpression:         "system_memory_usage > 85",
			QueryExpression:             "system_memory_usage",
			WarningThreshold:            85.0,
			CriticalThreshold:           95.0,
			ThresholdCondition:          "gt",
			EvaluationWindowSeconds:     300,
			EvaluationIntervalSeconds:   60,
			TriggerDurationSeconds:      120,
			RecoveryDurationSeconds:     60,
			AlertLevel:                  "warning",
			SeverityScore:               80,
			NotificationThrottleMinutes: 5,
			Status:                      "active",
			IsBuiltin:                   true,
		},
		{
			Name:                        "gpu_high_usage",
			DisplayName:                 "GPU使用率过高",
			Description:                 "当GPU使用率超过90%时触发告警",
			RuleType:                    "threshold",
			ConditionExpression:         "gpu_usage > 90",
			QueryExpression:             "gpu_usage",
			WarningThreshold:            85.0,
			CriticalThreshold:           95.0,
			ThresholdCondition:          "gt",
			EvaluationWindowSeconds:     180,
			EvaluationIntervalSeconds:   30,
			TriggerDurationSeconds:      60,
			RecoveryDurationSeconds:     30,
			AlertLevel:                  "warning",
			SeverityScore:               85,
			NotificationThrottleMinutes: 3,
			Status:                      "active",
			IsBuiltin:                   true,
		},
		{
			Name:                        "training_job_failed",
			DisplayName:                 "训练任务失败",
			Description:                 "当训练任务失败时立即触发告警",
			RuleType:                    "threshold",
			ConditionExpression:         "training_job_failed_count > 0",
			QueryExpression:             "training_job_count{status=\"failed\"}",
			WarningThreshold:            1.0,
			CriticalThreshold:           3.0,
			ThresholdCondition:          "gte",
			EvaluationWindowSeconds:     60,
			EvaluationIntervalSeconds:   30,
			TriggerDurationSeconds:      0,
			RecoveryDurationSeconds:     30,
			AlertLevel:                  "critical",
			SeverityScore:               90,
			NotificationThrottleMinutes: 1,
			Status:                      "active",
			IsBuiltin:                   true,
		},
		{
			Name:                        "training_queue_backlog",
			DisplayName:                 "训练队列积压",
			Description:                 "当训练队列长度超过50时触发告警",
			RuleType:                    "threshold",
			ConditionExpression:         "training_queue_length > 50",
			QueryExpression:             "training_queue_length",
			WarningThreshold:            50.0,
			CriticalThreshold:           100.0,
			ThresholdCondition:          "gt",
			EvaluationWindowSeconds:     300,
			EvaluationIntervalSeconds:   60,
			TriggerDurationSeconds:      180,
			RecoveryDurationSeconds:     120,
			AlertLevel:                  "warning",
			SeverityScore:               60,
			NotificationThrottleMinutes: 10,
			Status:                      "active",
			IsBuiltin:                   true,
		},
		{
			Name:                        "service_availability",
			DisplayName:                 "服务可用性检测",
			Description:                 "当关键服务不可用时触发告警",
			RuleType:                    "threshold",
			ConditionExpression:         "service_up < 1",
			QueryExpression:             "up{job=\"volctrain-api\"}",
			WarningThreshold:            1.0,
			CriticalThreshold:           1.0,
			ThresholdCondition:          "lt",
			EvaluationWindowSeconds:     60,
			EvaluationIntervalSeconds:   15,
			TriggerDurationSeconds:      30,
			RecoveryDurationSeconds:     15,
			AlertLevel:                  "critical",
			SeverityScore:               100,
			NotificationThrottleMinutes: 1,
			Status:                      "active",
			IsBuiltin:                   true,
		},
		// 业务指标告警规则
		{
			Name:                        "api_high_latency",
			DisplayName:                 "API响应延迟过高",
			Description:                 "当API 95%分位响应时间超过2秒时触发告警",
			RuleType:                    "threshold",
			ConditionExpression:         "api_latency_p95 > 2000",
			QueryExpression:             "histogram_quantile(0.95, api_request_duration_ms)",
			WarningThreshold:            2000.0,
			CriticalThreshold:           5000.0,
			ThresholdCondition:          "gt",
			EvaluationWindowSeconds:     300,
			EvaluationIntervalSeconds:   60,
			TriggerDurationSeconds:      120,
			RecoveryDurationSeconds:     60,
			AlertLevel:                  "warning",
			SeverityScore:               70,
			NotificationThrottleMinutes: 5,
			Status:                      "active",
			IsBuiltin:                   true,
		},
		{
			Name:                        "api_error_rate_high",
			DisplayName:                 "API错误率过高",
			Description:                 "当API错误率超过5%时触发告警",
			RuleType:                    "threshold",
			ConditionExpression:         "api_error_rate > 0.05",
			QueryExpression:             "api_error_rate",
			WarningThreshold:            0.05,
			CriticalThreshold:           0.10,
			ThresholdCondition:          "gt",
			EvaluationWindowSeconds:     300,
			EvaluationIntervalSeconds:   60,
			TriggerDurationSeconds:      120,
			RecoveryDurationSeconds:     60,
			AlertLevel:                  "warning",
			SeverityScore:               75,
			NotificationThrottleMinutes: 5,
			Status:                      "active",
			IsBuiltin:                   true,
		},
	}

	for _, rule := range builtinRules {
		// 检查规则是否已存在
		existing, err := b.alertRulesModel.FindOneByName(rule.Name)
		if err != nil && err != sql.ErrNoRows {
			b.logger.Errorf("检查内置告警规则失败 [%s]: %v", rule.Name, err)
			continue
		}

		// 如果不存在则插入
		if existing == nil {
			// 设置通知渠道（默认使用所有渠道）
			channels := []string{"default_email", "default_dingtalk"}
			if err := rule.SetNotificationChannels(channels); err != nil {
				b.logger.Errorf("设置通知渠道失败 [%s]: %v", rule.Name, err)
				continue
			}

			// 设置过滤标签
			filterLabels := map[string]interface{}{
				"service": "volctrain",
				"env":     "production",
			}
			if err := rule.SetFilterLabels(filterLabels); err != nil {
				b.logger.Errorf("设置过滤标签失败 [%s]: %v", rule.Name, err)
				continue
			}

			_, err := b.alertRulesModel.Insert(rule)
			if err != nil {
				b.logger.Errorf("插入内置告警规则失败 [%s]: %v", rule.Name, err)
			} else {
				b.logger.Infof("成功插入内置告警规则: %s", rule.Name)
			}
		}
	}

	b.logger.Info("内置告警规则初始化完成")
	return nil
}

// InitializeBuiltinTemplates 初始化内置通知模板
func (b *BuiltinRulesTemplates) InitializeBuiltinTemplates() error {
	b.logger.Info("初始化内置通知模板")

	builtinTemplates := []*model.VtNotificationTemplates{
		// 邮件模板
		{
			Name:         "email_firing_default",
			DisplayName:  "邮件告警触发模板",
			Description:  "默认的邮件告警触发通知模板",
			ChannelType:  "email",
			TemplateType: "firing",
			Subject:      "🚨 VolcTrain告警 - {{.Alert.RuleName}}",
			Content: `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>VolcTrain告警通知</title>
</head>
<body style="font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5;">
    <div style="max-width: 600px; margin: 0 auto; background-color: white; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
        <div style="background-color: #ff4d4f; color: white; padding: 20px; border-radius: 8px 8px 0 0;">
            <h2 style="margin: 0;">🚨 VolcTrain 监控告警</h2>
            <p style="margin: 5px 0 0 0;">告警触发通知</p>
        </div>
        <div style="padding: 20px;">
            <div style="display: inline-block; padding: 4px 8px; border-radius: 4px; font-weight: bold; color: #721c24; background-color: #f8d7da;">
                {{.Alert.AlertLevel}}级别告警
            </div>
            <h3 style="margin: 15px 0 10px 0;">{{.Alert.RuleName}}</h3>
            <p style="margin: 10px 0;">{{.Alert.Message}}</p>
            
            <h4>告警详情</h4>
            <table style="width: 100%; border-collapse: collapse;">
                <tr><th style="padding: 8px; text-align: left; border-bottom: 1px solid #ddd; background-color: #f8f9fa;">告警规则</th><td style="padding: 8px; border-bottom: 1px solid #ddd;">{{.Alert.RuleName}}</td></tr>
                <tr><th style="padding: 8px; text-align: left; border-bottom: 1px solid #ddd; background-color: #f8f9fa;">资源类型</th><td style="padding: 8px; border-bottom: 1px solid #ddd;">{{.Alert.ResourceType}}</td></tr>
                <tr><th style="padding: 8px; text-align: left; border-bottom: 1px solid #ddd; background-color: #f8f9fa;">资源名称</th><td style="padding: 8px; border-bottom: 1px solid #ddd;">{{.Alert.ResourceName}}</td></tr>
                <tr><th style="padding: 8px; text-align: left; border-bottom: 1px solid #ddd; background-color: #f8f9fa;">触发值</th><td style="padding: 8px; border-bottom: 1px solid #ddd;">{{.Alert.TriggerValue}}</td></tr>
                <tr><th style="padding: 8px; text-align: left; border-bottom: 1px solid #ddd; background-color: #f8f9fa;">阈值</th><td style="padding: 8px; border-bottom: 1px solid #ddd;">{{.Alert.ThresholdValue}}</td></tr>
                <tr><th style="padding: 8px; text-align: left; border-bottom: 1px solid #ddd; background-color: #f8f9fa;">触发时间</th><td style="padding: 8px; border-bottom: 1px solid #ddd;">{{.Alert.TriggeredAt}}</td></tr>
            </table>
        </div>
        <div style="background-color: #f8f9fa; padding: 15px 20px; border-radius: 0 0 8px 8px; font-size: 12px; color: #666;">
            <p style="margin: 0;">此邮件由 VolcTrain 监控系统自动发送，请勿直接回复。</p>
            <p style="margin: 5px 0 0 0;">发送时间: {{.Timestamp}}</p>
        </div>
    </div>
</body>
</html>`,
			IsDefault: true,
			IsBuiltin: true,
			Status:    "active",
		},
		{
			Name:         "email_resolved_default",
			DisplayName:  "邮件告警恢复模板",
			Description:  "默认的邮件告警恢复通知模板",
			ChannelType:  "email",
			TemplateType: "resolved",
			Subject:      "✅ VolcTrain告警恢复 - {{.Alert.RuleName}}",
			Content: `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>VolcTrain告警恢复通知</title>
</head>
<body style="font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5;">
    <div style="max-width: 600px; margin: 0 auto; background-color: white; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
        <div style="background-color: #52c41a; color: white; padding: 20px; border-radius: 8px 8px 0 0;">
            <h2 style="margin: 0;">✅ VolcTrain 告警恢复</h2>
            <p style="margin: 5px 0 0 0;">告警已恢复正常</p>
        </div>
        <div style="padding: 20px;">
            <h3 style="margin: 15px 0 10px 0;">{{.Alert.RuleName}}</h3>
            <p style="margin: 10px 0;">告警已恢复正常状态</p>
            
            <h4>恢复详情</h4>
            <table style="width: 100%; border-collapse: collapse;">
                <tr><th style="padding: 8px; text-align: left; border-bottom: 1px solid #ddd; background-color: #f8f9fa;">告警规则</th><td style="padding: 8px; border-bottom: 1px solid #ddd;">{{.Alert.RuleName}}</td></tr>
                <tr><th style="padding: 8px; text-align: left; border-bottom: 1px solid #ddd; background-color: #f8f9fa;">资源</th><td style="padding: 8px; border-bottom: 1px solid #ddd;">{{.Alert.ResourceType}}/{{.Alert.ResourceName}}</td></tr>
                <tr><th style="padding: 8px; text-align: left; border-bottom: 1px solid #ddd; background-color: #f8f9fa;">恢复时间</th><td style="padding: 8px; border-bottom: 1px solid #ddd;">{{.Timestamp}}</td></tr>
            </table>
        </div>
    </div>
</body>
</html>`,
			IsDefault: true,
			IsBuiltin: true,
			Status:    "active",
		},
		// 钉钉模板
		{
			Name:         "dingtalk_firing_default",
			DisplayName:  "钉钉告警触发模板",
			Description:  "默认的钉钉告警触发通知模板",
			ChannelType:  "dingtalk",
			TemplateType: "firing",
			Subject:      "VolcTrain告警",
			Content: `## 🚨 VolcTrain 告警触发

---

### 📋 告警信息
- **规则名称**：{{.Alert.RuleName}}
- **告警级别**：<font color="#FF4D4F">🔴 {{.Alert.AlertLevel}}</font>
- **触发时间**：{{.Alert.TriggeredAt}}

### 🎯 资源详情
- **资源类型**：{{.Alert.ResourceType}}
- **资源名称**：{{.Alert.ResourceName}}

### 📊 指标详情
- **当前值**：**{{.Alert.TriggerValue}}**
- **阈值**：{{.Alert.ThresholdValue}}

### 📝 告警描述
{{.Alert.Message}}

---
> 📅 **发送时间**：{{.Timestamp}}  
> 🤖 **来源**：VolcTrain监控系统`,
			IsDefault: true,
			IsBuiltin: true,
			Status:    "active",
		},
		{
			Name:         "dingtalk_resolved_default",
			DisplayName:  "钉钉告警恢复模板",
			Description:  "默认的钉钉告警恢复通知模板",
			ChannelType:  "dingtalk",
			TemplateType: "resolved",
			Subject:      "VolcTrain告警恢复",
			Content: `## ✅ VolcTrain 告警恢复

---

### 📋 恢复信息
- **规则名称**：{{.Alert.RuleName}}
- **告警级别**：<font color="#52C41A">🟢 已恢复</font>
- **恢复时间**：{{.Timestamp}}

### 🎯 资源详情
- **资源类型**：{{.Alert.ResourceType}}
- **资源名称**：{{.Alert.ResourceName}}

---
> 📅 **发送时间**：{{.Timestamp}}  
> 🤖 **来源**：VolcTrain监控系统`,
			IsDefault: true,
			IsBuiltin: true,
			Status:    "active",
		},
		// 短信模板
		{
			Name:         "sms_firing_default",
			DisplayName:  "短信告警触发模板",
			Description:  "默认的短信告警触发通知模板",
			ChannelType:  "sms",
			TemplateType: "firing",
			Subject:      "VolcTrain告警",
			Content:      "【VolcTrain】告警触发：{{.Alert.RuleName}}，级别：{{.Alert.AlertLevel}}，当前值：{{.Alert.TriggerValue}}，阈值：{{.Alert.ThresholdValue}}，时间：{{.Time}}",
			IsDefault:    true,
			IsBuiltin:    true,
			Status:       "active",
		},
		{
			Name:         "sms_resolved_default",
			DisplayName:  "短信告警恢复模板",
			Description:  "默认的短信告警恢复通知模板",
			ChannelType:  "sms",
			TemplateType: "resolved",
			Subject:      "VolcTrain告警恢复",
			Content:      "【VolcTrain】告警恢复：{{.Alert.RuleName}}已恢复正常，时间：{{.Time}}",
			IsDefault:    true,
			IsBuiltin:    true,
			Status:       "active",
		},
		// Webhook模板
		{
			Name:         "webhook_default",
			DisplayName:  "Webhook默认模板",
			Description:  "默认的Webhook通知模板",
			ChannelType:  "webhook",
			TemplateType: "default",
			Subject:      "VolcTrain Alert",
			Content: `{
  "alert_id": "{{.Alert.ID}}",
  "rule_name": "{{.Alert.RuleName}}",
  "alert_level": "{{.Alert.AlertLevel}}",
  "action": "{{.Action}}",
  "message": "{{.Alert.Message}}",
  "resource_type": "{{.Alert.ResourceType}}",
  "resource_name": "{{.Alert.ResourceName}}",
  "trigger_value": {{.Alert.TriggerValue}},
  "threshold_value": {{.Alert.ThresholdValue}},
  "triggered_at": "{{.Alert.TriggeredAt}}",
  "timestamp": {{.Timestamp}}
}`,
			IsDefault: true,
			IsBuiltin: true,
			Status:    "active",
		},
	}

	for _, template := range builtinTemplates {
		// 检查模板是否已存在
		existing, err := b.notificationTemplatesModel.FindOneByName(template.Name)
		if err != nil && err != sql.ErrNoRows {
			b.logger.Errorf("检查内置通知模板失败 [%s]: %v", template.Name, err)
			continue
		}

		// 如果不存在则插入
		if existing == nil {
			// 设置模板变量
			variables := map[string]interface{}{
				"alert_description": "告警详细信息",
				"timestamp_format":  "2006-01-02 15:04:05",
				"support_contact":   "support@volctrain.com",
			}
			variablesJSON, _ := json.Marshal(variables)
			template.Variables = string(variablesJSON)

			_, err := b.notificationTemplatesModel.Insert(template)
			if err != nil {
				b.logger.Errorf("插入内置通知模板失败 [%s]: %v", template.Name, err)
			} else {
				b.logger.Infof("成功插入内置通知模板: %s", template.Name)
			}
		}
	}

	b.logger.Info("内置通知模板初始化完成")
	return nil
}

// Initialize 初始化所有内置规则和模板
func (b *BuiltinRulesTemplates) Initialize() error {
	if err := b.InitializeBuiltinRules(); err != nil {
		return err
	}

	if err := b.InitializeBuiltinTemplates(); err != nil {
		return err
	}

	return nil
}
