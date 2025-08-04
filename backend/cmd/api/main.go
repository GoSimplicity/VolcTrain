package main

import (
	"flag"
	"fmt"

	"api/internal/config"
	"api/internal/handler"
	"api/internal/svc"
	"api/pkg/docs"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 注册Swagger文档
	docs.RegisterSwaggerHandler(server)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	fmt.Printf("Swagger UI: http://%s:%d/swagger\n", c.Host, c.Port)
	fmt.Printf("API Documentation: http://%s:%d/docs\n", c.Host, c.Port)
	server.Start()
}
