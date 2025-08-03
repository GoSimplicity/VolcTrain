package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	MySQL MySQLConfig `json:",optional"`
	Redis RedisConfig `json:",optional"`
	Auth  AuthConfig  `json:",optional"`
}

// MySQL数据库配置
type MySQLConfig struct {
	Host         string `json:",default=localhost"`
	Port         int    `json:",default=3306"`
	User         string `json:",default=root"`
	Password     string `json:",optional"`
	DBName       string `json:",default=vt_volctrain"`
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
