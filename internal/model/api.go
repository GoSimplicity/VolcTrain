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

package model

type Api struct {
	Model
	Name        string  `json:"name" gorm:"type:varchar(50);uniqueIndex:idx_name_del;not null;comment:API名称"`       // API名称，唯一且非空
	Path        string  `json:"path" gorm:"type:varchar(255);not null;comment:API路径"`                               // API路径，非空
	Method      int8    `json:"method" gorm:"type:tinyint(1);not null;comment:HTTP请求方法 1GET 2POST 3PUT 4DELETE"`    // 请求方法，使用int8节省空间
	Description string  `json:"description" gorm:"type:varchar(500);comment:API描述"`                                 // API描述
	Version     string  `json:"version" gorm:"type:varchar(20);default:v1;comment:API版本"`                           // API版本，默认v1
	Category    int8    `json:"category" gorm:"type:tinyint(1);not null;comment:API分类 1系统 2业务" binding:"oneof=1 2"` // API分类，使用int8节省空间
	IsPublic    int8    `json:"is_public" gorm:"type:tinyint(1);default:0;comment:是否公开 0否 1是" binding:"oneof=0 1"`  // 是否公开，使用int8节省空间
	Users       []*User `json:"users" gorm:"many2many:user_apis;comment:关联用户"`                                      // 多对多关联用户
}

type CreateApiRequest struct {
	Name        string `json:"name" binding:"required"`       // API名称
	Path        string `json:"path" binding:"required"`       // API路径
	Method      int    `json:"method" binding:"required"`     // 请求方法
	Description string `json:"description"`                   // API描述
	Version     string `json:"version"`                       // API版本
	Category    int    `json:"category"`                      // API分类
	IsPublic    int    `json:"is_public" binding:"oneof=0 1"` // 是否公开
}

type UpdateApiRequest struct {
	ID          int    `json:"id" binding:"required,gt=0"`    // API ID
	Name        string `json:"name" binding:"required"`       // API名称
	Path        string `json:"path" binding:"required"`       // API路径
	Method      int    `json:"method" binding:"required"`     // 请求方法
	Description string `json:"description"`                   // API描述
	Version     string `json:"version"`                       // API版本
	Category    int    `json:"category"`                      // API分类
	IsPublic    int    `json:"is_public" binding:"oneof=0 1"` // 是否公开
}

type DeleteApiRequest struct {
	ID int `json:"id" binding:"required,gt=0"` // API ID
}

type GetApiRequest struct {
	ID int `json:"id" binding:"required,gt=0"` // API ID
}

type ListApisRequest struct {
	ListReq
}
