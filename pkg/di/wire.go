//go:build wireinject

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
	authHandler "github.com/GoSimplicity/VolcTrain/internal/system/api"
	authDao "github.com/GoSimplicity/VolcTrain/internal/system/dao"
	authService "github.com/GoSimplicity/VolcTrain/internal/system/service"
	userHandler "github.com/GoSimplicity/VolcTrain/internal/user/api"
	userDao "github.com/GoSimplicity/VolcTrain/internal/user/dao"
	userService "github.com/GoSimplicity/VolcTrain/internal/user/service"
	ijwt "github.com/GoSimplicity/VolcTrain/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	_ "github.com/google/wire"
	"github.com/hibiken/asynq"
)

type Cmd struct {
	Server *gin.Engine
	Asynq  *asynq.Server
}

var HandlerSet = wire.NewSet(
	authHandler.NewRoleHandler,
	authHandler.NewApiHandler,
	userHandler.NewUserHandler,
)

var ServiceSet = wire.NewSet(
	userService.NewUserService,
	authService.NewApiService,
	authService.NewRoleService,
)

var DaoSet = wire.NewSet(
	userDao.NewUserDAO,
	authDao.NewRoleDAO,
	authDao.NewApiDAO,
)

var UtilSet = wire.NewSet(
	ijwt.NewJWTHandler,
)

var Injector = wire.NewSet(
	InitMiddlewares,
	InitGinServer,
	InitLogger,
	InitRedis,
	InitDB,
	InitViper,
	InitAsynqClient,
	InitAsynqServer,
	InitScheduler,
	wire.Struct(new(Cmd), "*"),
)

func ProvideCmd() *Cmd {
	wire.Build(
		Injector,
		HandlerSet,
		ServiceSet,
		DaoSet,
		UtilSet,
	)
	return &Cmd{}
}
