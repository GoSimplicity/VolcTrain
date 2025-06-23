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

	"github.com/GoSimplicity/VolcTrain/internal/constants"
	"github.com/GoSimplicity/VolcTrain/internal/system/service"
	"go.uber.org/zap"

	"github.com/GoSimplicity/VolcTrain/internal/model"
	"github.com/GoSimplicity/VolcTrain/internal/user/dao"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	SignUp(ctx context.Context, user *model.UserSignUpReq) error
	Login(ctx context.Context, user *model.UserLoginReq) (*model.User, error)
	GetProfile(ctx context.Context, uid int) (*model.User, error)
	GetPermCode(ctx context.Context, uid int) ([]string, error)
	GetUserDetail(ctx context.Context, uid int) (*model.User, error)
	GetUserList(ctx context.Context, req *model.ListReq) (model.ListResp[*model.User], error)
	ChangePassword(ctx context.Context, req *model.ChangePasswordReq) error
	WriteOff(ctx context.Context, username, password string) error
	UpdateProfile(ctx context.Context, user *model.UpdateProfileReq) error
	DeleteUser(ctx context.Context, uid int) error
}

type userService struct {
	dao     dao.UserDAO
	roleSvc service.RoleService
	l       *zap.Logger
}

func NewUserService(dao dao.UserDAO, roleSvc service.RoleService, l *zap.Logger) UserService {
	return &userService{
		dao:     dao,
		roleSvc: roleSvc,
		l:       l,
	}
}

// SignUp 用户注册
func (us *userService) SignUp(ctx context.Context, user *model.UserSignUpReq) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)

	if err := us.dao.CreateUser(ctx, &model.User{
		Username:     user.Username,
		Password:     user.Password,
		RealName:     user.RealName,
		Desc:         user.Desc,
		Mobile:       user.Mobile,
		FeiShuUserId: user.FeiShuUserId,
		AccountType:  user.AccountType,
		Enable:       1,
	}); err != nil {
		return err
	}

	return nil
}

// Login 用户登录
func (us *userService) Login(ctx context.Context, user *model.UserLoginReq) (*model.User, error) {
	u, err := us.dao.GetUserByUsername(ctx, user.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &model.User{}, constants.ErrorUserNotExist
	} else if err != nil {
		return &model.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		return &model.User{}, constants.ErrorPasswordIncorrect
	}

	return u, nil
}

// GetProfile 获取用户信息
func (us *userService) GetProfile(ctx context.Context, uid int) (*model.User, error) {
	return us.dao.GetUserByID(ctx, uid)
}

// GetPermCode 获取用户权限
func (us *userService) GetPermCode(ctx context.Context, uid int) ([]string, error) {
	codes, err := us.dao.GetPermCode(ctx, uid)
	if err != nil {
		return nil, err
	}

	return codes, nil
}

// GetUserList 获取用户列表
func (us *userService) GetUserList(ctx context.Context, req *model.ListReq) (model.ListResp[*model.User], error) {
	users, count, err := us.dao.GetAllUsers(ctx, req.Page, req.Size, req.Search)
	if err != nil {
		return model.ListResp[*model.User]{}, err
	}

	return model.ListResp[*model.User]{
		Items: users,
		Total: count,
	}, nil
}

// ChangePassword 修改密码
func (us *userService) ChangePassword(ctx context.Context, req *model.ChangePasswordReq) error {
	if req.Password == req.NewPassword {
		return errors.New("新密码不能与旧密码相同")
	}

	// 验证旧密码是否正确
	user, err := us.dao.GetUserByID(ctx, req.UserID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return constants.ErrorPasswordIncorrect
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 修改密码
	return us.dao.ChangePassword(ctx, req.UserID, string(hash))
}

// UpdateProfile 修改用户信息
func (us *userService) UpdateProfile(ctx context.Context, req *model.UpdateProfileReq) error {
	// 验证用户是否存在
	user, err := us.dao.GetUserByID(ctx, req.ID)
	if err != nil {
		return err
	}

	// 更新用户信息
	user.RealName = req.RealName
	user.Desc = req.Desc
	user.Mobile = req.Mobile
	user.FeiShuUserId = req.FeiShuUserId
	user.AccountType = req.AccountType
	user.HomePath = req.HomePath
	user.Enable = req.Enable

	return us.dao.UpdateProfile(ctx, user)
}

// WriteOff 注销账号
func (us *userService) WriteOff(ctx context.Context, username string, password string) error {
	// 验证用户是否存在
	user, err := us.dao.GetUserByUsername(ctx, username)
	if err != nil {
		return err
	}

	// 验证密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return constants.ErrorPasswordIncorrect
	}

	return us.dao.WriteOff(ctx, username, password)
}

func (us *userService) DeleteUser(ctx context.Context, uid int) error {
	// 删除用户角色关联
	if err := us.roleSvc.DeleteRole(ctx, uid); err != nil {
		us.l.Error("删除用户角色关联失败", zap.Int("uid", uid), zap.Error(err))
		return err
	}

	return us.dao.DeleteUser(ctx, uid)
}

func (us *userService) GetUserDetail(ctx context.Context, uid int) (*model.User, error) {
	user, err := us.dao.GetUserByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return user, nil
}
