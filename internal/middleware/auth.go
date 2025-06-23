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

package middleware

import (
	"strings"

	"github.com/GoSimplicity/VolcTrain/internal/system/service"
	"github.com/GoSimplicity/VolcTrain/pkg/utils"
	"github.com/gin-gonic/gin"
)

// 预定义跳过权限校验的路径
var skipAuthPaths = map[string]bool{
	"/api/user/login":                                   true,
	"/api/user/logout":                                  true,
	"/api/user/refresh_token":                           true,
	"/api/user/signup":                                  true,
	"/api/not_auth/getTreeNodeBindIps":                  true,
	"/api/monitor/prometheus_configs/prometheus":        true,
	"/api/monitor/prometheus_configs/prometheus_alert":  true,
	"/api/monitor/prometheus_configs/prometheus_record": true,
	"/api/monitor/prometheus_configs/alertManager":      true,
}

type AuthMiddleware struct {
	roleService service.RoleService
}

func NewAuthMiddleware(roleService service.RoleService) *AuthMiddleware {
	return &AuthMiddleware{
		roleService: roleService,
	}
}

func (am *AuthMiddleware) CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 快速检查是否需要跳过权限校验
		if skipAuthPaths[path] {
			c.Next()
			return
		}

		// 跳过静态资源和WebSocket路径
		if path == "/" ||
			strings.HasPrefix(path, "/api/ai/chat/ws") ||
			strings.HasPrefix(path, "/assets") ||
			strings.HasPrefix(path, "/_app.config.js") ||
			strings.HasPrefix(path, "/jse/") ||
			strings.HasPrefix(path, "/favicon.ico") ||
			strings.HasPrefix(path, "/js/") ||
			strings.HasPrefix(path, "/css/") {
			c.Next()
			return
		}

		user := c.MustGet("user").(utils.UserClaims)
		if user.Username == "admin" {
			c.Next()
			return
		}
		// TODO: 实现权限校验
		c.Next()
	}
}
