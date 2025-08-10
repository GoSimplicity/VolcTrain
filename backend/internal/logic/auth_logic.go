package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"api/internal/svc"
	"api/internal/types"
	"api/model"
	"api/pkg/auth"
	"api/pkg/errors"
	"api/pkg/validation"

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

// validateJWTSecret 验证JWT密钥长度和复杂度
func (l *LoginLogic) validateJWTSecret(secret, secretType string) error {
	if len(secret) < 32 {
		l.Errorf("%s密钥长度不足，至少需要32个字符，当前长度: %d", secretType, len(secret))
		return errors.NewInternalError(fmt.Sprintf("%s密钥长度不足", secretType))
	}
	
	// 检查密钥复杂度
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	
	for _, char := range secret {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char >= '!' && char <= '/' || char >= ':' && char <= '@' || char >= '[' && char <= '`' || char >= '{' && char <= '~':
			hasSpecial = true
		}
	}
	
	// 至少包含三种字符类型
	complexityCount := 0
	if hasUpper {
		complexityCount++
	}
	if hasLower {
		complexityCount++
	}
	if hasDigit {
		complexityCount++
	}
	if hasSpecial {
		complexityCount++
	}
	
	if complexityCount < 3 {
		l.Errorf("%s密钥复杂度不足，至少需要包含三种字符类型（大写字母、小写字母、数字、特殊字符）", secretType)
		return errors.NewInternalError(fmt.Sprintf("%s密钥复杂度不足", secretType))
	}
	
	return nil
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 输入验证
	vr := validation.NewValidationResult()
	
	// 验证用户名
	vr.ValidateUsername(req.Username, "用户名")
	
	// 验证密码
	vr.ValidatePassword(req.Password, "密码")
	
	// 验证防止SQL注入
	vr.ValidateNoSQLInjection(req.Username, "用户名")
	vr.ValidateNoSQLInjection(req.Password, "密码")
	
	if !vr.IsValid {
		return nil, errors.NewValidationError(strings.Join(vr.Errors, "; "))
	}

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

	// 验证JWT密钥长度和复杂度
	if err := l.validateJWTSecret(l.svcCtx.Config.Auth.AccessSecret, "访问令牌"); err != nil {
		return nil, err
	}
	if err := l.validateJWTSecret(l.svcCtx.Config.Auth.RefreshSecret, "刷新令牌"); err != nil {
		return nil, err
	}

	// 生成JWT token
	accessToken, err := auth.GenerateToken(user.Id, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		return nil, errors.NewInternalError("生成访问令牌失败")
	}

	// 生成刷新令牌
	refreshToken, err := auth.GenerateToken(user.Id, l.svcCtx.Config.Auth.RefreshSecret, l.svcCtx.Config.Auth.RefreshExpire)
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

// validateJWTSecret 验证JWT密钥长度和复杂度
func (l *RefreshTokenLogic) validateJWTSecret(secret, secretType string) error {
	if len(secret) < 32 {
		l.Errorf("%s密钥长度不足，至少需要32个字符，当前长度: %d", secretType, len(secret))
		return errors.NewInternalError(fmt.Sprintf("%s密钥长度不足", secretType))
	}
	
	// 检查密钥复杂度
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	
	for _, char := range secret {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char >= '!' && char <= '/' || char >= ':' && char <= '@' || char >= '[' && char <= '`' || char >= '{' && char <= '~':
			hasSpecial = true
		}
	}
	
	// 至少包含三种字符类型
	complexityCount := 0
	if hasUpper {
		complexityCount++
	}
	if hasLower {
		complexityCount++
	}
	if hasDigit {
		complexityCount++
	}
	if hasSpecial {
		complexityCount++
	}
	
	if complexityCount < 3 {
		l.Errorf("%s密钥复杂度不足，至少需要包含三种字符类型（大写字母、小写字母、数字、特殊字符）", secretType)
		return errors.NewInternalError(fmt.Sprintf("%s密钥复杂度不足", secretType))
	}
	
	return nil
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (resp *types.RefreshTokenResp, err error) {
	// 验证JWT密钥长度和复杂度
	if err := l.validateJWTSecret(l.svcCtx.Config.Auth.AccessSecret, "访问令牌"); err != nil {
		return nil, err
	}
	if err := l.validateJWTSecret(l.svcCtx.Config.Auth.RefreshSecret, "刷新令牌"); err != nil {
		return nil, err
	}
	
	// 使用刷新令牌密钥验证刷新令牌
	userID, err := auth.ValidateToken(req.RefreshToken, l.svcCtx.Config.Auth.RefreshSecret)
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

	// 生成新的token对（使用各自的密钥）
	accessToken, err := auth.GenerateToken(userID, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		return nil, errors.NewInternalError("生成访问令牌失败")
	}

	refreshToken, err := auth.GenerateToken(userID, l.svcCtx.Config.Auth.RefreshSecret, l.svcCtx.Config.Auth.RefreshExpire)
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
	// 从Authorization头获取token
	tokenString := l.extractTokenFromRequest()
	if tokenString == "" {
		l.Errorf("登出请求中未找到Token")
		return &types.LogoutResp{Message: "登出成功"}, nil
	}

	// 获取用户信息（用于日志记录）
	userIDValue := l.ctx.Value("userId")
	if userIDValue != nil {
		l.Infof("用户登出: userID=%v", userIDValue)
	} else {
		l.Info("匿名用户登出")
	}

	// 将token添加到Redis黑名单中
	if l.svcCtx.Redis != nil {
		// 创建token黑名单服务
		blacklist := auth.NewRedisTokenBlacklist(l.svcCtx.Redis)
		
		// 解析token获取过期时间
		_, err := auth.ValidateToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
		if err == nil {
			// 如果是有效的访问token，将其加入黑名单直到过期
			expireAt := time.Now().Add(time.Duration(l.svcCtx.Config.Auth.AccessExpire) * time.Second)
			err = blacklist.AddToken(l.ctx, tokenString, expireAt)
			if err != nil {
				l.Errorf("将Token加入黑名单失败: %v", err)
			} else {
				l.Info("Token已加入黑名单")
			}
		} else {
			// 尝试解析为刷新token
			if l.svcCtx.Config.Auth.RefreshSecret != "" {
				_, err := auth.ValidateToken(tokenString, l.svcCtx.Config.Auth.RefreshSecret)
				if err == nil {
					// 如果是有效的刷新token，将其加入黑名单直到过期
					expireAt := time.Now().Add(time.Duration(l.svcCtx.Config.Auth.RefreshExpire) * time.Second)
					err = blacklist.AddToken(l.ctx, tokenString, expireAt)
					if err != nil {
						l.Errorf("将刷新Token加入黑名单失败: %v", err)
					} else {
						l.Info("刷新Token已加入黑名单")
					}
				}
			}
		}
	} else {
		l.Errorf("Redis未连接，无法将Token加入黑名单")
	}
	
	return &types.LogoutResp{
		Message: "登出成功",
	}, nil
}

// getUserPermissions 获取用户权限列表
func (l *GetAccessCodesLogic) getUserPermissions(user *model.VtUsersSimple) []string {
	// 基础权限
	basePermissions := []string{
		"training:job:read",
		"user:profile:read",
	}
	
	// 根据用户类型添加权限
	switch user.UserType {
	case "admin":
		// 管理员拥有所有权限
		return []string{"*"}
	case "user":
		// 普通用户权限
		userPermissions := []string{
			"training:job:create",
			"training:job:update",
			"training:job:cancel",
			"training:queue:read",
			"gpu:device:read",
			"model:read",
			"dataset:read",
			"workspace:read",
		}
		return append(basePermissions, userPermissions...)
	case "service":
		// 服务账户权限
		servicePermissions := []string{
			"training:job:read",
			"training:job:update:status",
			"gpu:device:read",
			"monitoring:read",
		}
		return append(basePermissions, servicePermissions...)
	default:
		// 未知用户类型，只给基础权限
		return basePermissions
	}
}

// extractTokenFromRequest 从请求中提取token
func (l *LogoutLogic) extractTokenFromRequest() string {
	// 从中间件注入的上下文中获取token
	tokenValue := l.ctx.Value("token")
	if tokenValue != nil {
		if token, ok := tokenValue.(string); ok {
			return token
		}
	}
	
	return ""
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

	// 根据用户角色获取权限码
	codes := l.getUserPermissions(user)

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