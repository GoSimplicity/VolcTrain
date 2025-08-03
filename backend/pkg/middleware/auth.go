package middleware

import (
	"context"
	"net/http"
	"strings"

	"api/pkg/auth"
	"api/pkg/response"
)

// JWTAuthMiddleware JWT认证中间件
type JWTAuthMiddleware struct {
	jwtService *auth.JWTService
	blacklist  *auth.TokenBlacklist
	skipPaths  []string // 跳过认证的路径
}

// NewJWTAuthMiddleware 创建JWT认证中间件
func NewJWTAuthMiddleware(jwtService *auth.JWTService, blacklist *auth.TokenBlacklist) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		jwtService: jwtService,
		blacklist:  blacklist,
		skipPaths: []string{
			"/api/v1/users/login",
			"/api/v1/users/register",
			"/api/v1/health",
			"/health",
			"/ready",
			"/metrics",
		},
	}
}

// AddSkipPath 添加跳过认证的路径
func (j *JWTAuthMiddleware) AddSkipPath(path string) {
	j.skipPaths = append(j.skipPaths, path)
}

// Handler 中间件处理函数
func (j *JWTAuthMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 检查是否需要跳过认证
			if j.shouldSkip(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// 从请求头获取Token
			tokenString := j.extractToken(r)
			if tokenString == "" {
				response.Unauthorized(w, "缺少认证Token")
				return
			}

			// 检查Token是否在黑名单中
			if j.blacklist != nil && j.blacklist.IsBlacklisted(tokenString) {
				response.Unauthorized(w, "Token已失效")
				return
			}

			// 解析和验证Token
			claims, err := j.jwtService.ParseAccessToken(tokenString)
			if err != nil {
				response.Unauthorized(w, "Token无效或已过期")
				return
			}

			// 将用户信息添加到请求上下文
			ctx := context.WithValue(r.Context(), "user", claims)
			ctx = context.WithValue(ctx, "userID", claims.UserID)
			ctx = context.WithValue(ctx, "username", claims.Username)
			ctx = context.WithValue(ctx, "roles", claims.Roles)
			ctx = context.WithValue(ctx, "permissions", claims.Permissions)

			// 继续处理请求
			next.ServeHTTP(w, r.WithContext(ctx))
		})
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
	user, ok := ctx.Value("user").(*auth.JWTClaims)
	if !ok {
		return nil
	}
	return user
}

// GetUserIDFromContext 从上下文获取用户ID
func GetUserIDFromContext(ctx context.Context) int64 {
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		return 0
	}
	return userID
}

// GetUsernameFromContext 从上下文获取用户名
func GetUsernameFromContext(ctx context.Context) string {
	username, ok := ctx.Value("username").(string)
	if !ok {
		return ""
	}
	return username
}

// GetRolesFromContext 从上下文获取角色列表
func GetRolesFromContext(ctx context.Context) []string {
	roles, ok := ctx.Value("roles").([]string)
	if !ok {
		return []string{}
	}
	return roles
}

// GetPermissionsFromContext 从上下文获取权限列表
func GetPermissionsFromContext(ctx context.Context) []string {
	permissions, ok := ctx.Value("permissions").([]string)
	if !ok {
		return []string{}
	}
	return permissions
}

// RequirePermission 权限检查中间件
func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			permissions := GetPermissionsFromContext(r.Context())

			// 检查是否有所需权限
			if !containsPermission(permissions, permission) {
				response.Forbidden(w, "权限不足")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole 角色检查中间件
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles := GetRolesFromContext(r.Context())

			// 检查是否有所需角色
			if !containsRole(roles, role) {
				response.Forbidden(w, "角色权限不足")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyPermission 任一权限检查中间件
func RequireAnyPermission(permissions ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyRole 任一角色检查中间件
func RequireAnyRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			next.ServeHTTP(w, r)
		})
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
func (j *JWTAuthMiddleware) OptionalAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 尝试获取Token
			tokenString := j.extractToken(r)
			if tokenString != "" {
				// 如果有Token，尝试解析
				if claims, err := j.jwtService.ParseAccessToken(tokenString); err == nil {
					// 检查Token是否在黑名单中
					if j.blacklist == nil || !j.blacklist.IsBlacklisted(tokenString) {
						// 将用户信息添加到上下文
						ctx := context.WithValue(r.Context(), "user", claims)
						ctx = context.WithValue(ctx, "userID", claims.UserID)
						ctx = context.WithValue(ctx, "username", claims.Username)
						ctx = context.WithValue(ctx, "roles", claims.Roles)
						ctx = context.WithValue(ctx, "permissions", claims.Permissions)
						r = r.WithContext(ctx)
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
