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

package service

import (
	"context"
	"errors"

	"github.com/GoSimplicity/VolcTrain/internal/model"
	"github.com/GoSimplicity/VolcTrain/internal/system/dao"
	"go.uber.org/zap"
)

type ApiService interface {
	CreateApi(ctx context.Context, req *model.CreateApiRequest) error
	GetApiById(ctx context.Context, id int) (*model.Api, error)
	UpdateApi(ctx context.Context, req *model.UpdateApiRequest) error
	DeleteApi(ctx context.Context, id int) error
	ListApis(ctx context.Context, req *model.ListApisRequest) (model.ListResp[*model.Api], error)
}

type apiService struct {
	l   *zap.Logger
	dao dao.ApiDAO
}

func NewApiService(l *zap.Logger, dao dao.ApiDAO) ApiService {
	return &apiService{
		l:   l,
		dao: dao,
	}
}

// CreateApi 创建新的API
func (a *apiService) CreateApi(ctx context.Context, req *model.CreateApiRequest) error {
	if req == nil {
		a.l.Warn("API不能为空")
		return errors.New("api不能为空")
	}

	return a.dao.CreateApi(ctx, a.buildCreateApi(req))
}

// GetApiById 根据ID获取API
func (a *apiService) GetApiById(ctx context.Context, id int) (*model.Api, error) {
	if id <= 0 {
		a.l.Warn("API ID无效", zap.Int("ID", id))
		return nil, errors.New("api id无效")
	}

	return a.dao.GetApiById(ctx, id)
}

// UpdateApi 更新API信息
func (a *apiService) UpdateApi(ctx context.Context, req *model.UpdateApiRequest) error {
	if req == nil {
		a.l.Warn("API不能为空")
		return errors.New("api不能为空")
	}

	return a.dao.UpdateApi(ctx, a.buildUpdateApi(req))
}

// DeleteApi 删除指定ID的API
func (a *apiService) DeleteApi(ctx context.Context, id int) error {
	if id <= 0 {
		a.l.Warn("API ID无效", zap.Int("ID", id))
		return errors.New("api id无效")
	}

	return a.dao.DeleteApi(ctx, id)
}

// ListApis 分页获取API列表
func (a *apiService) ListApis(ctx context.Context, req *model.ListApisRequest) (model.ListResp[*model.Api], error) {
	if req.Page < 1 || req.Size < 1 {
		a.l.Warn("分页参数无效", zap.Int("页码", req.Page), zap.Int("每页数量", req.Size))
		return model.ListResp[*model.Api]{}, errors.New("分页参数无效")
	}

	apis, total, err := a.dao.ListApis(ctx, req.Page, req.Size, req.Search)
	if err != nil {
		a.l.Error("获取API列表失败", zap.Error(err))
		return model.ListResp[*model.Api]{}, err
	}

	return model.ListResp[*model.Api]{
		Items: apis,
		Total: total,
	}, nil
}

func (a *apiService) buildCreateApi(req *model.CreateApiRequest) *model.Api {
	return &model.Api{
		Name:        req.Name,
		Path:        req.Path,
		Method:      int8(req.Method),
		Description: req.Description,
		Version:     req.Version,
		Category:    int8(req.Category),
		IsPublic:    int8(req.IsPublic),
	}
}

func (a *apiService) buildUpdateApi(req *model.UpdateApiRequest) *model.Api {
	return &model.Api{
		Model: model.Model{
			ID: req.ID,
		},
		Name:        req.Name,
		Path:        req.Path,
		Method:      int8(req.Method),
		Description: req.Description,
		Version:     req.Version,
		Category:    int8(req.Category),
		IsPublic:    int8(req.IsPublic),
	}
}
