/*
 * Apache License
 * Version 2.0, January 2004
 * http://www.apache.org/licenses/
 *
 * TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION
 *
 * Copyright 2025 Bamboo
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GoSimplicity/VolcTrain/pkg/di"
	"github.com/GoSimplicity/VolcTrain/pkg/utils"
	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// 开发模式下注释此行，生产模式下取消注释
// //go:embed ui/apps/web-antd/dist/*
var embeddedFiles embed.FS

func main() {
	if err := Init(); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}
}

func Init() error {
	// 初始化配置
	if err := di.InitViper(); err != nil {
		return fmt.Errorf("初始化配置失败: %v", err)
	}

	// 初始化 Web 服务器和其他组件
	cmd := di.ProvideCmd()

	// 设置中间件
	cmd.Server.Use(cors.Default())
	cmd.Server.Use(gzip.Gzip(gzip.BestCompression))

	// 设置请求头打印路由
	cmd.Server.GET("/headers", printHeaders)

	// 判断是否为生产模式（通过检查嵌入文件是否可用）
	isProductionMode := true
	_, err := embeddedFiles.ReadFile("ui/apps/web-antd/dist/index.html")
	if err != nil {
		isProductionMode = false
		log.Println("运行在开发模式，仅提供API服务")
	} else {
		log.Println("运行在生产模式，提供完整前后端服务")
	}

	// 只在生产模式下挂载静态文件
	if isProductionMode {
		// 挂载静态文件
		assetsFS, _ := fs.Sub(embeddedFiles, "ui/apps/web-antd/dist/assets")
		cmd.Server.StaticFS("/assets", http.FS(assetsFS))

		// 直接返回 index.html
		cmd.Server.GET("/", func(c *gin.Context) {
			index, err := embeddedFiles.ReadFile("ui/apps/web-antd/dist/index.html")
			if err != nil {
				c.String(http.StatusInternalServerError, "Internal Server Error")
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", index)
		})

		// 处理 favicon.ico 请求
		cmd.Server.GET("/favicon.ico", func(c *gin.Context) {
			favicon, err := embeddedFiles.ReadFile("ui/apps/web-antd/dist/favicon.ico")
			if err != nil {
				c.Status(http.StatusNoContent)
				return
			}
			c.Data(http.StatusOK, "image/x-icon", favicon)
		})

		// 处理 _app.config.js 请求
		cmd.Server.GET("/_app.config.js", func(c *gin.Context) {
			config, err := embeddedFiles.ReadFile("ui/apps/web-antd/dist/_app.config.js")
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
			c.Data(http.StatusOK, "application/javascript", config)
		})

		// 处理 jse 目录下的文件请求
		cmd.Server.GET("/jse/:filename", func(c *gin.Context) {
			filename := c.Param("filename")
			file, err := embeddedFiles.ReadFile("ui/apps/web-antd/dist/jse/" + filename)
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
			c.Data(http.StatusOK, "application/javascript", file)
		})

		// 处理 css 目录下的文件请求
		cmd.Server.GET("/css/:filename", func(c *gin.Context) {
			filename := c.Param("filename")
			file, err := embeddedFiles.ReadFile("ui/apps/web-antd/dist/css/" + filename)
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
			c.Data(http.StatusOK, "text/css", file)
		})

		// 处理 js 目录下的文件请求
		cmd.Server.GET("/js/:filename", func(c *gin.Context) {
			filename := c.Param("filename")
			file, err := embeddedFiles.ReadFile("ui/apps/web-antd/dist/js/" + filename)
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
			c.Data(http.StatusOK, "application/javascript", file)
		})
	} else {
		// 开发模式下，提供一个简单的首页
		cmd.Server.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "AI-CloudOps API 服务运行中 (开发模式)",
				"status":  "running",
			})
		})
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    ":" + viper.GetString("server.port"),
		Handler: cmd.Server,
	}

	// 创建系统信号接收器
	quit := make(chan os.Signal, 1)
	// 监听 SIGINT 和 SIGTERM 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 在goroutine中启动服务器
	go func() {
		showBootInfo(viper.GetString("server.port"), isProductionMode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号
	<-quit
	log.Println("正在关闭服务器...")

	// 设置关闭超时时间为30秒
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭HTTP服务器,等待所有连接处理完成
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("服务器关闭异常: %v", err)
		return fmt.Errorf("服务器关闭失败: %v", err)
	}

	// 等待所有goroutine完成
	time.Sleep(2 * time.Second)

	log.Println("服务器已成功关闭")
	return nil
}

// printHeaders 打印请求头信息
func printHeaders(c *gin.Context) {
	headers := make(map[string]string)
	for key, values := range c.Request.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}
	c.JSON(http.StatusOK, headers)
}

func showBootInfo(port string, isProductionMode bool) {
	// 获取本机所有 IP 地址
	ips, err := utils.GetLocalIPs()
	if err != nil {
		log.Printf("获取本机 IP 失败: %v", err)
		return
	}

	// 打印启动信息
	modeText := "生产模式"
	if !isProductionMode {
		modeText = "开发模式 (仅API)"
	}

	color.Green("VolcTrain 启动成功 (%s)", modeText)
	fmt.Printf("%s  ", color.GreenString("➜"))
	fmt.Printf("%s    ", color.New(color.Bold).Sprint("Local:"))
	fmt.Printf("%s\n", color.MagentaString("http://localhost:%s/", port))

	for _, ip := range ips {
		fmt.Printf("%s  ", color.GreenString("➜"))
		fmt.Printf("%s  ", color.New(color.Bold).Sprint("Network:"))
		fmt.Printf("%s\n", color.MagentaString("http://%s:%s/", ip, port))
	}
}
