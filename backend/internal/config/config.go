package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	MySQL        MySQLConfig        `json:",optional"`
	Redis        RedisConfig        `json:",optional"`
	Auth         AuthConfig         `json:",optional"`
	Security     SecurityConfig     `json:",optional"`
	Storage      StorageConfig      `json:",optional"`
	K8s          K8sConfig          `json:",optional"`
	Notification NotificationConfig `json:",optional"`
}

// MySQL数据库配置
type MySQLConfig struct {
	Host         string `json:",default=localhost"`
	Port         int    `json:",default=3306"`
	User         string `json:",default=root"`
	Password     string `json:",optional"`
	DBName       string `json:",default=volctraindb"`
	Charset      string `json:",default=utf8mb4"`
	ParseTime    bool   `json:",default=true"`
	Loc          string `json:",default=Asia/Shanghai"`
	MaxOpenConns int    `json:",default=100"`
	MaxIdleConns int    `json:",default=10"`
	MaxLifetime  int    `json:",default=3600"`
}

// Redis配置
type RedisConfig struct {
	Host         string `json:",default=localhost"`
	Port         int    `json:",default=6379"`
	Password     string `json:",optional"`
	DB           int    `json:",default=0"`
	PoolSize     int    `json:",default=100"`
	MinIdleConns int    `json:",default=10"`
}

// JWT认证配置
type AuthConfig struct {
	AccessSecret  string `json:",optional"`
	AccessExpire  int64  `json:",default=86400"`
	RefreshExpire int64  `json:",default=604800"`
}

// Prometheus监控配置
type PrometheusConfig struct {
	Host string `json:",default=0.0.0.0"`
	Port int    `json:",default=9090"`
	Path string `json:",default=/metrics"`
}

// 安全配置
type SecurityConfig struct {
	RateLimit RateLimitConfig `json:",optional"`
	CORS      CORSConfig      `json:",optional"`
}

// 限流配置
type RateLimitConfig struct {
	Requests int `json:",default=1000"`
	Window   int `json:",default=3600"`
}

// CORS配置
type CORSConfig struct {
	AllowedOrigins   string `json:",default=*"`
	AllowedMethods   string `json:",default=GET,POST,PUT,DELETE,OPTIONS"`
	AllowedHeaders   string `json:",default=Content-Type,Authorization,X-Requested-With"`
	AllowCredentials bool   `json:",default=true"`
}

// 存储配置
type StorageConfig struct {
	WorkspacePath  string `json:",default=/data/workspace"`
	DatasetPath    string `json:",default=/data/datasets"`
	ModelPath      string `json:",default=/data/models"`
	LogsPath       string `json:",default=/data/logs"`
	CheckpointPath string `json:",default=/data/checkpoints"`
}

// K8s配置
type K8sConfig struct {
	ConfigPath          string `json:",optional"`
	Namespace           string `json:",default=default"`
	EnablePodMonitoring bool   `json:",default=true"`
}

// 通知配置
type NotificationConfig struct {
	Enabled  bool              `json:",default=false"`
	Email    EmailConfig       `json:",optional"`
	DingTalk DingTalkConfig    `json:",optional"`
}

// 邮件配置
type EmailConfig struct {
	SMTPHost string `json:",optional"`
	SMTPPort int    `json:",default=587"`
	Username string `json:",optional"`
	Password string `json:",optional"`
}

// 钉钉配置
type DingTalkConfig struct {
	WebhookURL string `json:",optional"`
	Secret     string `json:",optional"`
}
