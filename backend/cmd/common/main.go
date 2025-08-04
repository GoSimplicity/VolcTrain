package main

import (
	"flag"
	"fmt"
	"os"

	"api/internal/config"
	"api/internal/handler"
	"api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "", "the config file")

// getConfigFile 根据环境变量获取配置文件路径
func getConfigFile() string {
	// 如果通过命令行参数指定了配置文件，优先使用
	if *configFile != "" {
		return *configFile
	}
	
	// 根据 DEPLOY_ENV 环境变量决定配置文件
	deployEnv := os.Getenv("DEPLOY_ENV")
	switch deployEnv {
	case "production":
		return "etc/config-production.yaml"
	default:
		return "etc/config-dev.yaml"
	}
}

func main() {
	flag.Parse()

	configPath := getConfigFile()
	var c config.Config
	conf.MustLoad(configPath, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	fmt.Printf("Config file: %s\n", configPath)
	fmt.Printf("Environment: %s\n", os.Getenv("DEPLOY_ENV"))
	server.Start()
}
