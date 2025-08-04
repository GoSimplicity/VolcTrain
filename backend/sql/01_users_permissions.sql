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
    department VARCHAR(128) COMMENT '部门',
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
-- 插入初始数据
-- 插入默认管理员用户（密码：admin123）
INSERT INTO vt_users (username, email, password_hash, salt, real_name, nickname, department, status, user_type, email_verified) 
VALUES ('admin', 'admin@volctrain.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IcnhVkuF2lFhgq6pf9s8cPLHvQOj2O', 'default_salt', '系统管理员', '管理员', '技术部', 'active', 'admin', 1);

-- 插入测试用户（密码：user123）
INSERT INTO vt_users (username, email, password_hash, salt, real_name, nickname, department, status, user_type, email_verified) 
VALUES ('testuser', 'user@volctrain.com', '$2a$10$XHpZc8gOCe9ZwGQYjmqAdeNKIqNhVF4gJcUFqFOhqpEE7aXGt4OIy', 'default_salt', '测试用户', '用户', '测试部', 'active', 'user', 1);

-- 插入基础权限
INSERT INTO vt_permissions (name, display_name, description, module, action, resource, permission_code) VALUES
('training:job:read', '查看训练任务', '查看训练任务列表和详情', 'training', 'read', 'job', 'training.job.read'),
('training:job:create', '创建训练任务', '创建新的训练任务', 'training', 'create', 'job', 'training.job.create'),
('training:job:update', '更新训练任务', '修改训练任务配置', 'training', 'update', 'job', 'training.job.update'),
('training:job:delete', '删除训练任务', '删除训练任务', 'training', 'delete', 'job', 'training.job.delete'),
('training:queue:read', '查看训练队列', '查看训练队列信息', 'training', 'read', 'queue', 'training.queue.read'),
('training:queue:create', '创建训练队列', '创建新的训练队列', 'training', 'create', 'queue', 'training.queue.create'),
('gpu:device:read', '查看GPU设备', '查看GPU设备信息', 'gpu', 'read', 'device', 'gpu.device.read'),
('gpu:device:manage', '管理GPU设备', '管理GPU设备分配', 'gpu', 'manage', 'device', 'gpu.device.manage'),
('system:manage', '系统管理', '系统配置和管理', 'system', 'manage', '*', 'system.manage');

-- 插入基础角色
INSERT INTO vt_roles (name, display_name, description, role_code, role_type) VALUES
('admin', '系统管理员', '拥有系统所有权限的管理员角色', 'ADMIN', 'system'),
('user', '普通用户', '基础用户角色，具有基本的训练任务权限', 'USER', 'system'),
('developer', '开发者', '开发者角色，具有更多权限', 'DEVELOPER', 'custom');

-- 分配管理员权限（管理员拥有所有权限）
INSERT INTO vt_role_permissions (role_id, permission_id) 
SELECT 1, id FROM vt_permissions;

-- 分配普通用户权限
INSERT INTO vt_role_permissions (role_id, permission_id)
SELECT 2, id FROM vt_permissions WHERE name IN ('training:job:read', 'training:job:create', 'training:queue:read', 'gpu:device:read');

-- 分配用户角色
INSERT INTO vt_user_roles (user_id, role_id, status) VALUES
(1, 1, 'active'), -- admin用户分配管理员角色
(2, 2, 'active'); -- testuser分配普通用户角色