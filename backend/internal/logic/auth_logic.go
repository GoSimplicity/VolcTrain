package logic

import (
	"context"
	"encoding/json"
	"time"

	"api/internal/svc"
	"api/internal/types"
	"api/pkg/auth"
	"api/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 查找用户
	user, err := l.svcCtx.VtUsersModel.FindByUsername(l.ctx, req.Username)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	// 验证密码
	if !auth.CheckPassword(req.Password, user.Password) {
		return nil, errors.ErrInvalidPassword
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, errors.ErrUserDisabled
	}

	// 生成JWT token
	accessToken, err := auth.GenerateToken(user.Id, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		return nil, errors.NewInternalError("生成访问令牌失败")
	}

	refreshToken, err := auth.GenerateToken(user.Id, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.RefreshExpire)
	if err != nil {
		return nil, errors.NewInternalError("生成刷新令牌失败")
	}

	// 更新用户最后登录时间
	user.LastLoginAt = time.Now()
	if err := l.svcCtx.VtUsersModel.Update(l.ctx, user); err != nil {
		l.Errorf("更新用户最后登录时间失败: %v", err)
	}

	// 构造响应
	resp = &types.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    l.svcCtx.Config.Auth.AccessExpire,
		TokenType:    "Bearer",
		UserInfo: types.UserInfo{
			ID:         user.Id,
			Username:   user.Username,
			Email:      user.Email,
			RealName:   user.RealName,
			Status:     user.Status,
			UserType:   user.UserType,
			Department: user.Department,
			CreatedAt:  user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  user.UpdatedAt.Format(time.RFC3339),
		},
	}

	l.Infof("用户登录成功: %s", req.Username)
	return resp, nil
}

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (resp *types.RefreshTokenResp, err error) {
	// 验证刷新令牌
	userID, err := auth.ValidateToken(req.RefreshToken, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	// 查找用户
	user, err := l.svcCtx.VtUsersModel.FindOne(l.ctx, userID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, errors.ErrUserDisabled
	}

	// 生成新的token
	accessToken, err := auth.GenerateToken(userID, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		return nil, errors.NewInternalError("生成访问令牌失败")
	}

	refreshToken, err := auth.GenerateToken(userID, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.RefreshExpire)
	if err != nil {
		return nil, errors.NewInternalError("生成刷新令牌失败")
	}

	resp = &types.RefreshTokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    l.svcCtx.Config.Auth.AccessExpire,
		TokenType:    "Bearer",
	}

	return resp, nil
}

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutReq) (resp *types.LogoutResp, err error) {
	// TODO: 实现token黑名单机制
	// 这里可以将token添加到Redis黑名单中
	
	l.Info("用户登出成功")
	return &types.LogoutResp{}, nil
}

type GetAccessCodesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAccessCodesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccessCodesLogic {
	return &GetAccessCodesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAccessCodesLogic) GetAccessCodes(req *types.EmptyReq) (resp *types.GetAccessCodesResp, err error) {
	// 从JWT中间件获取用户ID
	userIDValue := l.ctx.Value("userId")
	if userIDValue == nil {
		return nil, errors.NewAuthError("未找到用户信息")
	}
	
	userID, ok := userIDValue.(json.Number)
	if !ok {
		return nil, errors.NewAuthError("用户ID格式错误")
	}
	
	uid, err := userID.Int64()
	if err != nil {
		return nil, errors.NewAuthError("用户ID转换失败")
	}
	
	// 查找用户
	user, err := l.svcCtx.VtUsersModel.FindOne(l.ctx, uid)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	// TODO: 根据用户角色获取权限码
	// 这里简化实现，实际应该根据用户角色查询权限表
	codes := make([]string, 0)
	
	switch user.UserType {
	case "admin":
		codes = []string{"*"} // 管理员拥有所有权限
	case "user":
		codes = []string{
			"training:job:read",
			"training:job:create",
			"training:queue:read",
			"gpu:device:read",
		}
	default:
		codes = []string{"training:job:read"}
	}

	resp = &types.GetAccessCodesResp{
		Codes: codes,
	}

	return resp, nil
}

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.EmptyReq) (resp *types.GetUserInfoResp, err error) {
	// 从JWT中间件获取用户ID
	userIDValue := l.ctx.Value("userId")
	if userIDValue == nil {
		return nil, errors.NewAuthError("未找到用户信息")
	}
	
	userID, ok := userIDValue.(json.Number)
	if !ok {
		return nil, errors.NewAuthError("用户ID格式错误")
	}
	
	uid, err := userID.Int64()
	if err != nil {
		return nil, errors.NewAuthError("用户ID转换失败")
	}
	
	// 查找用户
	user, err := l.svcCtx.VtUsersModel.FindOne(l.ctx, uid)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	resp = &types.GetUserInfoResp{
		UserInfo: types.UserInfo{
			ID:         user.Id,
			Username:   user.Username,
			Email:      user.Email,
			RealName:   user.RealName,
			Status:     user.Status,
			UserType:   user.UserType,
			Department: user.Department,
			CreatedAt:  user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  user.UpdatedAt.Format(time.RFC3339),
		},
	}

	return resp, nil
}