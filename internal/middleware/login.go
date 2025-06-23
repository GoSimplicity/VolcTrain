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

	ijwt "github.com/GoSimplicity/VolcTrain/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type JWTMiddleware struct {
	ijwt.Handler
}

func NewJWTMiddleware(hdl ijwt.Handler) *JWTMiddleware {
	return &JWTMiddleware{
		Handler: hdl,
	}
}

// CheckLogin 校验JWT
func (m *JWTMiddleware) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		// 如果请求的路径是下述路径，则不进行token验证
		if path == "/api/user/login" ||
			//path == "/api/user/signup" ||   // 不允许用户自己注册账号
			path == "/api/user/logout" ||
			path == "/api/user/refresh_token" ||
			path == "/api/user/signup" ||
			path == "/api/not_auth/getTreeNodeBindIps" ||
			path == "/api/monitor/prometheus_configs/prometheus" ||
			path == "/api/monitor/prometheus_configs/prometheus_alert" ||
			path == "/api/monitor/prometheus_configs/prometheus_record" ||
			path == "/api/monitor/prometheus_configs/alertManager" ||
			path == "/" ||
			strings.HasPrefix(path, "/api/ai/chat/ws") ||
			strings.HasPrefix(path, "/assets") ||
			strings.HasPrefix(path, "/_app.config.js") ||
			strings.HasPrefix(path, "/jse/") ||
			strings.HasPrefix(path, "/favicon.ico") ||
			strings.HasPrefix(path, "/js/") ||
			strings.HasPrefix(path, "/css/") {
			ctx.Next()
			return
		}

		var uc ijwt.UserClaims
		var tokenStr string

		// 如果是/api/tree/ecs/console开头的路径，从查询参数获取token
		if strings.HasPrefix(path, "/api/tree/ecs/console") {
			tokenStr = ctx.Query("token")
		} else {
			// 从请求中提取token
			tokenStr = m.ExtractToken(ctx)
		}

		key := []byte(viper.GetString("jwt.key1"))
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})

		if err != nil {
			// token 错误
			ctx.AbortWithStatus(401)
			return
		}

		if token == nil || !token.Valid {
			// token 非法或过期
			ctx.AbortWithStatus(401)
			return
		}

		// 检查是否携带ua头
		if uc.UserAgent == "" {
			ctx.AbortWithStatus(401)
			return
		}

		// 检查会话是否有效
		err = m.CheckSession(ctx, uc.Ssid)

		if err != nil {
			ctx.AbortWithStatus(401)
			return
		}

		ctx.Set("user", uc)
	}
}
