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

package di

import (
	systemApi "github.com/GoSimplicity/VolcTrain/internal/system/api"
	userApi "github.com/GoSimplicity/VolcTrain/internal/user/api"

	"github.com/gin-gonic/gin"
)

// InitGinServer 初始化web服务
func InitGinServer(
	m []gin.HandlerFunc,
	userHdl *userApi.UserHandler,
	authApiHdl *systemApi.ApiHandler,
	authRoleHdl *systemApi.RoleHandler,
) *gin.Engine {
	server := gin.Default()
	server.Use(m...)
	userHdl.RegisterRoutes(server)
	authApiHdl.RegisterRouters(server)
	authRoleHdl.RegisterRouters(server)
	return server
}
