-- 用户表
CREATE TABLE vt_users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(64) NOT NULL COMMENT '用户名',
    email VARCHAR(128) NOT NULL COMMENT '邮箱',
    phone VARCHAR(32) COMMENT '手机号',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    real_name VARCHAR(128) COMMENT '真实姓名',
    nickname VARCHAR(128) COMMENT '昵称',
    avatar_url VARCHAR(512) COMMENT '头像URL',
    status ENUM('active', 'inactive', 'locked') DEFAULT 'active' COMMENT '用户状态',
    last_login_at BIGINT DEFAULT 0 COMMENT '最后登录时间戳',
    last_login_ip VARCHAR(64) COMMENT '最后登录IP',
    created_at BIGINT DEFAULT 0 COMMENT '创建时间戳',
    updated_at BIGINT DEFAULT 0 COMMENT '更新时间戳',
    deleted_at BIGINT DEFAULT 0 COMMENT '软删除时间戳',
    -- 索引
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    UNIQUE KEY uk_username (username),
    UNIQUE KEY uk_email (email)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户表';
-- 部门表
CREATE TABLE vt_departments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL COMMENT '部门名称',
    parent_id BIGINT DEFAULT 0 COMMENT '父部门ID',
    level INT DEFAULT 1 COMMENT '层级',
    path VARCHAR(512) COMMENT '路径',
    manager_id BIGINT COMMENT '部门负责人ID',
    description TEXT COMMENT '部门描述',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at BIGINT DEFAULT 0 COMMENT '创建时间戳',
    updated_at BIGINT DEFAULT 0 COMMENT '更新时间戳',
    deleted_at BIGINT DEFAULT 0 COMMENT '软删除时间戳',
    -- 索引
    INDEX idx_name (name),
    INDEX idx_parent_id (parent_id),
    INDEX idx_manager_id (manager_id),
    INDEX idx_status (status),
    UNIQUE KEY uk_name (name)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '部门表';
-- 用户部门关联表
CREATE TABLE vt_user_departments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    department_id BIGINT NOT NULL COMMENT '部门ID',
    position VARCHAR(128) COMMENT '职位',
    is_primary TINYINT(1) DEFAULT 1 COMMENT '是否主部门',
    created_at BIGINT DEFAULT 0 COMMENT '创建时间戳',
    -- 索引
    INDEX idx_user_id (user_id),
    INDEX idx_department_id (department_id),
    UNIQUE KEY uk_user_department (user_id, department_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户部门关联表';
-- 用户会话表
CREATE TABLE vt_user_sessions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    session_id VARCHAR(128) NOT NULL COMMENT '会话ID',
    device_type ENUM('web', 'mobile', 'desktop') DEFAULT 'web' COMMENT '设备类型',
    ip_address VARCHAR(64) COMMENT 'IP地址',
    user_agent TEXT COMMENT '用户代理',
    login_at BIGINT DEFAULT 0 COMMENT '登录时间戳',
    last_activity_at BIGINT DEFAULT 0 COMMENT '最后活动时间戳',
    expires_at BIGINT DEFAULT 0 COMMENT '过期时间戳',
    status ENUM('active', 'expired', 'revoked') DEFAULT 'active' COMMENT '状态',
    -- 索引
    INDEX idx_user_id (user_id),
    INDEX idx_session_id (session_id),
    INDEX idx_status (status),
    INDEX idx_expires_at (expires_at),
    UNIQUE KEY uk_session_id (session_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户会话表';