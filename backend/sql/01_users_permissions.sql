-- 用户表
CREATE TABLE vt_users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(64) NOT NULL COMMENT '用户名',
    email VARCHAR(128) NOT NULL COMMENT '邮箱',
    phone VARCHAR(32) COMMENT '手机号',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    salt VARCHAR(64) NOT NULL COMMENT '密码盐值',
    real_name VARCHAR(128) COMMENT '真实姓名',
    nickname VARCHAR(128) COMMENT '昵称',
    status ENUM('active', 'inactive', 'locked', 'pending') DEFAULT 'pending' COMMENT '用户状态',
    user_type ENUM('admin', 'user', 'service') DEFAULT 'user' COMMENT '用户类型',
    last_login_at TIMESTAMP NULL COMMENT '最后登录时间',
    last_login_ip VARCHAR(64) COMMENT '最后登录IP',
    password_expires_at TIMESTAMP NULL COMMENT '密码过期时间',
    login_attempts INT DEFAULT 0 COMMENT '登录尝试次数',
    locked_until TIMESTAMP NULL COMMENT '锁定到期时间',
    mfa_enabled TINYINT(1) DEFAULT 0 COMMENT '是否启用双因素认证',
    mfa_secret VARCHAR(128) COMMENT '双因素认证密钥',
    email_verified TINYINT(1) DEFAULT 0 COMMENT '邮箱是否验证',
    phone_verified TINYINT(1) DEFAULT 0 COMMENT '手机是否验证',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    UNIQUE KEY uk_username (username),
    UNIQUE KEY uk_email (email),
    INDEX idx_status (status),
    INDEX idx_user_type (user_type),
    INDEX idx_email_verified (email_verified),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_last_login_at (last_login_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户表';
-- 部门表
CREATE TABLE vt_departments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL COMMENT '部门名称',
    parent_id BIGINT DEFAULT NULL COMMENT '父部门ID',
    level INT DEFAULT 1 COMMENT '层级',
    path VARCHAR(512) COMMENT '层级路径',
    description TEXT COMMENT '部门描述',
    department_code VARCHAR(64) COMMENT '部门编码',
    sort_order INT DEFAULT 0 COMMENT '排序',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    UNIQUE KEY uk_name (name),
    UNIQUE KEY uk_department_code (department_code),
    INDEX idx_parent_id (parent_id),
    INDEX idx_status (status),
    INDEX idx_level (level),
    INDEX idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '部门表';
-- 权限表
CREATE TABLE vt_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL COMMENT '权限标识',
    display_name VARCHAR(256) COMMENT '权限名称',
    description TEXT COMMENT '权限描述',
    module VARCHAR(64) COMMENT '所属模块',
    action VARCHAR(64) COMMENT '操作类型',
    resource VARCHAR(64) COMMENT '资源类型',
    permission_code VARCHAR(128) COMMENT '权限编码',
    parent_id BIGINT DEFAULT 0 COMMENT '父权限ID',
    level INT DEFAULT 1 COMMENT '权限层级',
    sort_order INT DEFAULT 0 COMMENT '排序',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_name (name),
    UNIQUE KEY uk_permission_code (permission_code),
    INDEX idx_module (module),
    INDEX idx_action (action),
    INDEX idx_resource (resource),
    INDEX idx_parent_id (parent_id),
    INDEX idx_status (status)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '权限表';
-- 角色表
CREATE TABLE vt_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL COMMENT '角色标识',
    display_name VARCHAR(256) COMMENT '角色名称',
    description TEXT COMMENT '角色描述',
    role_code VARCHAR(128) COMMENT '角色编码',
    role_type ENUM('system', 'custom') DEFAULT 'custom' COMMENT '角色类型',
    sort_order INT DEFAULT 0 COMMENT '排序',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    UNIQUE KEY uk_name (name),
    UNIQUE KEY uk_role_code (role_code),
    INDEX idx_role_type (role_type),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色表';
-- 角色权限关联表
CREATE TABLE vt_role_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_id BIGINT NOT NULL COMMENT '角色ID',
    permission_id BIGINT NOT NULL COMMENT '权限ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_role_permission (role_id, permission_id),
    INDEX idx_role_id (role_id),
    INDEX idx_permission_id (permission_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色权限关联表';
-- 用户角色关联表
CREATE TABLE vt_user_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    workspace_id BIGINT COMMENT '工作空间ID(空表示全局角色)',
    expires_at TIMESTAMP NULL COMMENT '过期时间',
    assigned_by BIGINT COMMENT '分配人ID',
    status ENUM('active', 'inactive', 'expired') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_user_role_workspace (user_id, role_id, workspace_id),
    INDEX idx_user_id (user_id),
    INDEX idx_role_id (role_id),
    INDEX idx_workspace_id (workspace_id),
    INDEX idx_expires_at (expires_at),
    INDEX idx_assigned_by (assigned_by),
    INDEX idx_status (status)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户角色关联表';
-- 用户部门关联表
CREATE TABLE vt_user_departments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    department_id BIGINT NOT NULL COMMENT '部门ID',
    position VARCHAR(128) COMMENT '职位',
    is_primary TINYINT(1) DEFAULT 1 COMMENT '是否主部门',
    join_date DATE COMMENT '入职日期',
    leave_date DATE COMMENT '离职日期',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_user_department (user_id, department_id),
    INDEX idx_user_id (user_id),
    INDEX idx_department_id (department_id),
    INDEX idx_is_primary (is_primary),
    INDEX idx_status (status),
    INDEX idx_join_date (join_date),
    INDEX idx_leave_date (leave_date)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户部门关联表';
-- 部门管理者关联表
CREATE TABLE vt_department_managers (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    department_id BIGINT NOT NULL COMMENT '部门ID',
    user_id BIGINT NOT NULL COMMENT '管理者用户ID',
    manager_type ENUM('primary', 'deputy', 'assistant') DEFAULT 'primary' COMMENT '管理者类型',
    start_date DATE COMMENT '开始日期',
    end_date DATE COMMENT '结束日期',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_dept_user_type (department_id, user_id, manager_type),
    INDEX idx_department_id (department_id),
    INDEX idx_user_id (user_id),
    INDEX idx_manager_type (manager_type),
    INDEX idx_status (status)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '部门管理者关联表';