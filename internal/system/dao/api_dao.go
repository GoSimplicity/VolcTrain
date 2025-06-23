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

package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoSimplicity/VolcTrain/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ApiDAO interface {
	CreateApi(ctx context.Context, api *model.Api) error
	GetApiById(ctx context.Context, id int) (*model.Api, error)
	UpdateApi(ctx context.Context, api *model.Api) error
	DeleteApi(ctx context.Context, id int) error
	ListApis(ctx context.Context, page, size int, search string) ([]*model.Api, int64, error)
}

type apiDAO struct {
	db *gorm.DB
	l  *zap.Logger
}

func NewApiDAO(db *gorm.DB, l *zap.Logger) ApiDAO {
	return &apiDAO{
		db: db,
		l:  l,
	}
}

// CreateApi 创建新的API记录
func (a *apiDAO) CreateApi(ctx context.Context, api *model.Api) error {
	if api == nil {
		return errors.New("API对象不能为空")
	}

	if api.Name == "" {
		return errors.New("API名称不能为空")
	}

	if api.Path == "" {
		return errors.New("API路径不能为空")
	}

	if api.Method <= 0 {
		return errors.New("无效的HTTP方法")
	}

	// 检查API名称是否已存在
	var count int64
	if err := a.db.WithContext(ctx).Model(&model.Api{}).Where("name = ? AND deleted_at = ?", api.Name, 0).Count(&count).Error; err != nil {
		return fmt.Errorf("检查API名称失败: %v", err)
	}
	if count > 0 {
		return errors.New("API名称已存在")
	}

	return a.db.WithContext(ctx).Create(api).Error
}

// GetApiById 根据ID获取API记录
func (a *apiDAO) GetApiById(ctx context.Context, id int) (*model.Api, error) {
	if id <= 0 {
		return nil, errors.New("无效的API ID")
	}

	var api model.Api
	if err := a.db.WithContext(ctx).Where("id = ?", id).First(&api).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询API失败: %v", err)
	}

	return &api, nil
}

// UpdateApi 更新API记录
func (a *apiDAO) UpdateApi(ctx context.Context, api *model.Api) error {
	if api == nil {
		return errors.New("API对象不能为空")
	}

	if api.ID <= 0 {
		return errors.New("无效的API ID")
	}

	if api.Name == "" {
		return errors.New("API名称不能为空")
	}

	if api.Path == "" {
		return errors.New("API路径不能为空")
	}

	if api.Method <= 0 {
		return errors.New("无效的HTTP方法")
	}

	// 获取旧的API记录
	oldApi, err := a.GetApiById(ctx, api.ID)
	if err != nil {
		return err
	}

	if oldApi == nil {
		return errors.New("API不存在")
	}

	// 检查API名称是否已被其他记录使用
	if oldApi.Name != api.Name {
		var count int64
		if err := a.db.WithContext(ctx).Model(&model.Api{}).Where("name = ? AND id != ?", api.Name, api.ID).Count(&count).Error; err != nil {
			return fmt.Errorf("检查API名称失败: %v", err)
		}
		if count > 0 {
			return errors.New("API名称已被其他记录使用")
		}
	}

	updates := map[string]interface{}{
		"name":        api.Name,
		"path":        api.Path,
		"method":      api.Method,
		"description": api.Description,
		"version":     api.Version,
		"category":    api.Category,
		"is_public":   api.IsPublic,
	}

	// 开启事务
	tx := a.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 更新API记录
	if err := tx.Model(&model.Api{}).
		Where("id = ?", api.ID).
		Updates(updates).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新API失败: %v", err)
	}

	return tx.Commit().Error
}

// DeleteApi 软删除API记录
func (a *apiDAO) DeleteApi(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("无效的API ID")
	}

	// 检查API是否存在
	api, err := a.GetApiById(ctx, id)
	if err != nil {
		return err
	}
	if api == nil {
		return errors.New("API不存在")
	}

	// 检查API是否被角色使用
	var count int64
	if err := a.db.WithContext(ctx).Table("role_apis").Where("api_id = ?", id).Count(&count).Error; err != nil {
		return fmt.Errorf("检查API使用情况失败: %v", err)
	}
	if count > 0 {
		return errors.New("该API已被角色使用，无法删除")
	}

	result := a.db.WithContext(ctx).Model(&model.Api{}).Where("id = ?", id).Delete(&model.Api{})
	if result.Error != nil {
		return fmt.Errorf("删除API失败: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("API不存在或已被删除")
	}

	return nil
}

// ListApis 分页获取API列表
func (a *apiDAO) ListApis(ctx context.Context, page, size int, search string) ([]*model.Api, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	query := a.db.WithContext(ctx).Model(&model.Api{})
	if search != "" {
		query = query.Where("name LIKE ? OR path LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var apis []*model.Api
	var total int64

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取API总数失败: %v", err)
	}

	// 获取分页数据
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("id DESC").Find(&apis).Error; err != nil {
		return nil, 0, fmt.Errorf("获取API列表失败: %v", err)
	}

	return apis, total, nil
}
