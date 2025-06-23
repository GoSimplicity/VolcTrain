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

package api

import (
	"github.com/GoSimplicity/VolcTrain/internal/model"
	"github.com/GoSimplicity/VolcTrain/internal/system/service"
	"github.com/GoSimplicity/VolcTrain/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	svc service.ApiService
}

func NewApiHandler(svc service.ApiService) *ApiHandler {
	return &ApiHandler{
		svc: svc,
	}
}

func (h *ApiHandler) RegisterRouters(server *gin.Engine) {
	apiGroup := server.Group("/api/apis")

	apiGroup.GET("/list", h.ListApis)
	apiGroup.POST("/create", h.CreateAPI)
	apiGroup.PUT("/update/:id", h.UpdateAPI)
	apiGroup.DELETE("/delete/:id", h.DeleteAPI)
	apiGroup.GET("/detail/:id", h.DetailAPI)
}

// ListApis 获取API列表
func (a *ApiHandler) ListApis(ctx *gin.Context) {
	var req model.ListApisRequest

	utils.HandleRequest(ctx, &req, func() (interface{}, error) {
		return a.svc.ListApis(ctx, &req)
	})
}

// CreateAPI 创建新的API
func (a *ApiHandler) CreateAPI(ctx *gin.Context) {
	var req model.CreateApiRequest

	utils.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, a.svc.CreateApi(ctx, &req)
	})
}

// UpdateAPI 更新API信息
func (a *ApiHandler) UpdateAPI(ctx *gin.Context) {
	var req model.UpdateApiRequest

	id, err := utils.GetParamID(ctx)
	if err != nil {
		utils.ErrorWithMessage(ctx, err.Error())
		return
	}

	req.ID = id

	utils.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, a.svc.UpdateApi(ctx, &req)
	})
}

// DeleteAPI 删除API
func (a *ApiHandler) DeleteAPI(ctx *gin.Context) {
	var req model.DeleteApiRequest

	id, err := utils.GetParamID(ctx)
	if err != nil {
		utils.ErrorWithMessage(ctx, err.Error())
		return
	}

	req.ID = id

	utils.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, a.svc.DeleteApi(ctx, req.ID)
	})
}

// DetailAPI 获取API详情
func (a *ApiHandler) DetailAPI(ctx *gin.Context) {
	var req model.GetApiRequest

	id, err := utils.GetParamID(ctx)
	if err != nil {
		utils.ErrorWithMessage(ctx, err.Error())
		return
	}

	req.ID = id

	utils.HandleRequest(ctx, &req, func() (interface{}, error) {
		return a.svc.GetApiById(ctx, id)
	})
}
