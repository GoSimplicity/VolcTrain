-- 权限表
CREATE TABLE vt_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL COMMENT '权限标识',
    display_name VARCHAR(256) COMMENT '权限名称',
    description TEXT COMMENT '权限描述',
    module VARCHAR(64) COMMENT '所属模块',
    action VARCHAR(64) COMMENT '操作类型',
    resource VARCHAR(64) COMMENT '资源类型',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at BIGINT DEFAULT 0 COMMENT '创建时间戳',
    updated_at BIGINT DEFAULT 0 COMMENT '更新时间戳',
    -- 索引
    INDEX idx_name (name),
    INDEX idx_module (module),
    INDEX idx_status (status),
    UNIQUE KEY uk_name (name)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '权限表';
-- 角色表
CREATE TABLE vt_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL COMMENT '角色标识',
    display_name VARCHAR(256) COMMENT '角色名称',
    description TEXT COMMENT '角色描述',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at BIGINT DEFAULT 0 COMMENT '创建时间戳',
    updated_at BIGINT DEFAULT 0 COMMENT '更新时间戳',
    deleted_at BIGINT DEFAULT 0 COMMENT '软删除时间戳',
    -- 索引
    INDEX idx_name (name),
    INDEX idx_status (status),
    UNIQUE KEY uk_name (name)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色表';
-- 角色权限关联表
CREATE TABLE vt_role_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_id BIGINT NOT NULL COMMENT '角色ID',
    permission_id BIGINT NOT NULL COMMENT '权限ID',
    created_at BIGINT DEFAULT 0 COMMENT '创建时间戳',
    -- 索引
    INDEX idx_role_id (role_id),
    INDEX idx_permission_id (permission_id),
    UNIQUE KEY uk_role_permission (role_id, permission_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色权限关联表';
-- 用户权限表（直接授权给用户的权限）
CREATE TABLE vt_user_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    permission_id BIGINT NOT NULL COMMENT '权限ID',
    resource_type VARCHAR(64) COMMENT '资源类型',
    resource_id BIGINT COMMENT '资源ID',
    granted TINYINT(1) DEFAULT 1 COMMENT '是否授权(1:授权 0:拒绝)',
    expires_at BIGINT DEFAULT 0 COMMENT '过期时间戳(0:永不过期)',
    created_by BIGINT COMMENT '授权人ID',
    created_at BIGINT DEFAULT 0 COMMENT '创建时间戳',
    updated_at BIGINT DEFAULT 0 COMMENT '更新时间戳',
    -- 索引
    INDEX idx_user_id (user_id),
    INDEX idx_permission_id (permission_id),
    INDEX idx_resource (resource_type, resource_id),
    INDEX idx_expires_at (expires_at),
    UNIQUE KEY uk_user_permission_resource (
        user_id,
        permission_id,
        resource_type,
        resource_id
    )
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户权限表';
-- 用户角色关联表
CREATE TABLE vt_user_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    resource_type VARCHAR(64) COMMENT '资源类型',
    resource_id BIGINT COMMENT '资源ID',
    expires_at BIGINT DEFAULT 0 COMMENT '过期时间戳(0:永不过期)',
    assigned_by BIGINT COMMENT '分配人ID',
    created_at BIGINT DEFAULT 0 COMMENT '创建时间戳',
    updated_at BIGINT DEFAULT 0 COMMENT '更新时间戳',
    -- 索引
    INDEX idx_user_id (user_id),
    INDEX idx_role_id (role_id),
    INDEX idx_resource (resource_type, resource_id),
    INDEX idx_expires_at (expires_at),
    UNIQUE KEY uk_user_role_resource (user_id, role_id, resource_type, resource_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户角色关联表';
-- 部门权限表（授权给部门的权限）
CREATE TABLE vt_department_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    department_id BIGINT NOT NULL COMMENT '部门ID',
    permission_id BIGINT NOT NULL COMMENT '权限ID',
    resource_type VARCHAR(64) COMMENT '资源类型',
    resource_id BIGINT COMMENT '资源ID',
    granted TINYINT(1) DEFAULT 1 COMMENT '是否授权(1:授权 0:拒绝)',
    expires_at BIGINT DEFAULT 0 COMMENT '过期时间戳(0:永不过期)',
    created_by BIGINT COMMENT '授权人ID',
    created_at BIGINT DEFAULT 0 COMMENT '创建时间戳',
    updated_at BIGINT DEFAULT 0 COMMENT '更新时间戳',
    -- 索引
    INDEX idx_department_id (department_id),
    INDEX idx_permission_id (permission_id),
    INDEX idx_resource (resource_type, resource_id),
    INDEX idx_expires_at (expires_at),
    UNIQUE KEY uk_dept_permission_resource (
        department_id,
        permission_id,
        resource_type,
        resource_id
    )
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '部门权限表';
-- 部门角色关联表
CREATE TABLE vt_department_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    department_id BIGINT NOT NULL COMMENT '部门ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    resource_type VARCHAR(64) COMMENT '资源类型',
    resource_id BIGINT COMMENT '资源ID',
    expires_at BIGINT DEFAULT 0 COMMENT '过期时间戳(0:永不过期)',
    assigned_by BIGINT COMMENT '分配人ID',
    created_at BIGINT DEFAULT 0 COMMENT '创建时间戳',
    updated_at BIGINT DEFAULT 0 COMMENT '更新时间戳',
    -- 索引
    INDEX idx_department_id (department_id),
    INDEX idx_role_id (role_id),
    INDEX idx_resource (resource_type, resource_id),
    INDEX idx_expires_at (expires_at),
    UNIQUE KEY uk_dept_role_resource (
        department_id,
        role_id,
        resource_type,
        resource_id
    )
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '部门角色关联表';