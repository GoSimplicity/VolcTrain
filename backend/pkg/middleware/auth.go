package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"api/pkg/auth"
	"api/pkg/response"
)

// JWTAuthMiddleware JWT认证中间件
type JWTAuthMiddleware struct {
	jwtService     *auth.JWTService
	redisBlacklist *auth.RedisTokenBlacklist
	skipPaths      []string // 跳过认证的路径
}

// NewJWTAuthMiddleware 创建JWT认证中间件
func NewJWTAuthMiddleware(jwtService *auth.JWTService, redisBlacklist *auth.RedisTokenBlacklist) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		jwtService:     jwtService,
		redisBlacklist: redisBlacklist,
		skipPaths: []string{
			"/api/v1/auth/login",
			"/api/v1/auth/refresh",
			"/api/v1/health",
			"/health",
			"/ready",
			"/metrics",
			"/swagger",
			"/docs",
		},
	}
}

// AddSkipPath 添加跳过认证的路径
func (j *JWTAuthMiddleware) AddSkipPath(path string) {
	j.skipPaths = append(j.skipPaths, path)
}

// Handler 中间件处理函数
func (j *JWTAuthMiddleware) Handler() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 检查是否需要跳过认证
			if j.shouldSkip(r.URL.Path) {
				next(w, r)
				return
			}

			// 从请求头获取Token
			tokenString := j.extractToken(r)
			if tokenString == "" {
				response.Unauthorized(w, "缺少认证Token")
				return
			}

			// 检查Token是否在Redis黑名单中
			if j.redisBlacklist != nil {
				isBlacklisted, err := j.redisBlacklist.IsBlacklisted(r.Context(), tokenString)
				if err != nil {
					// Redis错误不影响认证流程，只记录日志
					log.Printf("检查Token黑名单失败: %v", err)
				} else if isBlacklisted {
					response.Unauthorized(w, "Token已失效")
					return
				}
			}

			// 解析和验证Token
			claims, err := j.jwtService.ParseAccessToken(tokenString)
			if err != nil {
				response.Unauthorized(w, "Token无效或已过期")
				return
			}

			// 将用户信息添加到请求上下文
			// 类型化上下文键注入
			ctx := WithValue(r.Context(), CtxKeyUser, claims)
			ctx = WithValue(ctx, CtxKeyUserID, claims.UserID)
			ctx = WithValue(ctx, CtxKeyUsername, claims.Username)
			ctx = WithValue(ctx, CtxKeyRoles, claims.Roles)
			ctx = WithValue(ctx, CtxKeyPerms, claims.Permissions)
			ctx = WithValue(ctx, CtxKeyToken, tokenString)

			// 继续处理请求
			next(w, r.WithContext(ctx))
		}
	}
}

// shouldSkip 检查是否应该跳过认证
func (j *JWTAuthMiddleware) shouldSkip(path string) bool {
	for _, skipPath := range j.skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// extractToken 从请求中提取Token
func (j *JWTAuthMiddleware) extractToken(r *http.Request) string {
	// 从Authorization头提取
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Bearer Token格式
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimPrefix(authHeader, "Bearer ")
		}
		// 直接Token格式
		return authHeader
	}

	// 从查询参数提取（不推荐，但支持）
	token := r.URL.Query().Get("token")
	if token != "" {
		return token
	}

	// 从Cookie提取
	cookie, err := r.Cookie("access_token")
	if err == nil {
		return cookie.Value
	}

	return ""
}

// GetUserFromContext 从上下文获取用户信息
func GetUserFromContext(ctx context.Context) *auth.JWTClaims {
	user, ok := Value[*auth.JWTClaims](ctx, CtxKeyUser)
	if !ok {
		return nil
	}
	return user
}

// GetUserIDFromContext 从上下文获取用户ID
func GetUserIDFromContext(ctx context.Context) int64 {
	userID, ok := Value[int64](ctx, CtxKeyUserID)
	if !ok {
		return 0
	}
	return userID
}

// GetUsernameFromContext 从上下文获取用户名
func GetUsernameFromContext(ctx context.Context) string {
	username, ok := Value[string](ctx, CtxKeyUsername)
	if !ok {
		return ""
	}
	return username
}

// GetRolesFromContext 从上下文获取角色列表
func GetRolesFromContext(ctx context.Context) []string {
	roles, ok := Value[[]string](ctx, CtxKeyRoles)
	if !ok {
		return []string{}
	}
	return roles
}

// GetPermissionsFromContext 从上下文获取权限列表
func GetPermissionsFromContext(ctx context.Context) []string {
	permissions, ok := Value[[]string](ctx, CtxKeyPerms)
	if !ok {
		return []string{}
	}
	return permissions
}

// RequirePermission 权限检查中间件
func RequirePermission(permission string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			permissions := GetPermissionsFromContext(r.Context())

			// 检查是否有所需权限
			if !containsPermission(permissions, permission) {
				response.Forbidden(w, "权限不足")
				return
			}

			next(w, r)
		}
	}
}

// RequireRole 角色检查中间件
func RequireRole(role string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			roles := GetRolesFromContext(r.Context())

			// 检查是否有所需角色
			if !containsRole(roles, role) {
				response.Forbidden(w, "角色权限不足")
				return
			}

			next(w, r)
		}
	}
}

// RequireAnyPermission 任一权限检查中间件
func RequireAnyPermission(permissions ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userPermissions := GetPermissionsFromContext(r.Context())

			// 检查是否有任一所需权限
			hasPermission := false
			for _, permission := range permissions {
				if containsPermission(userPermissions, permission) {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				response.Forbidden(w, "权限不足")
				return
			}

			next(w, r)
		}
	}
}

// RequireAnyRole 任一角色检查中间件
func RequireAnyRole(roles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userRoles := GetRolesFromContext(r.Context())

			// 检查是否有任一所需角色
			hasRole := false
			for _, role := range roles {
				if containsRole(userRoles, role) {
					hasRole = true
					break
				}
			}

			if !hasRole {
				response.Forbidden(w, "角色权限不足")
				return
			}

			next(w, r)
		}
	}
}

// containsPermission 检查权限列表是否包含指定权限
func containsPermission(permissions []string, permission string) bool {
	for _, p := range permissions {
		if p == permission || p == "*" { // "*" 表示所有权限
			return true
		}
	}
	return false
}

// containsRole 检查角色列表是否包含指定角色
func containsRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role || r == "admin" { // admin角色拥有所有权限
			return true
		}
	}
	return false
}

// OptionalAuth 可选认证中间件（不强制要求认证，但如果有Token会解析）
func (j *JWTAuthMiddleware) OptionalAuth() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 尝试获取Token
			tokenString := j.extractToken(r)
			if tokenString != "" {
				// 如果有Token，尝试解析
				if claims, err := j.jwtService.ParseAccessToken(tokenString); err == nil {
					// 检查Token是否在Redis黑名单中
					if j.redisBlacklist != nil {
						isBlacklisted, err := j.redisBlacklist.IsBlacklisted(r.Context(), tokenString)
						if err != nil {
							// Redis错误不影响认证流程，只记录日志
							log.Printf("检查Token黑名单失败: %v", err)
						} else if !isBlacklisted {
							// 将用户信息添加到上下文
							ctx := WithValue(r.Context(), CtxKeyUser, claims)
							ctx = WithValue(ctx, CtxKeyUserID, claims.UserID)
							ctx = WithValue(ctx, CtxKeyUsername, claims.Username)
							ctx = WithValue(ctx, CtxKeyRoles, claims.Roles)
							ctx = WithValue(ctx, CtxKeyPerms, claims.Permissions)
							r = r.WithContext(ctx)
						}
					} else {
						// 没有Redis黑名单服务，直接将用户信息添加到上下文
						ctx := WithValue(r.Context(), CtxKeyUser, claims)
						ctx = WithValue(ctx, CtxKeyUserID, claims.UserID)
						ctx = WithValue(ctx, CtxKeyUsername, claims.Username)
						ctx = WithValue(ctx, CtxKeyRoles, claims.Roles)
						ctx = WithValue(ctx, CtxKeyPerms, claims.Permissions)
						r = r.WithContext(ctx)
					}
				}
			}

			next(w, r)
		}
	}
}
