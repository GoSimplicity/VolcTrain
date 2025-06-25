-- API密钥表
CREATE TABLE vt_api_keys (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    key_name VARCHAR(128) NOT NULL COMMENT '密钥名称',
    api_key VARCHAR(256) NOT NULL COMMENT 'API密钥',
    api_secret VARCHAR(256) NOT NULL COMMENT 'API密钥密文',
    key_type ENUM('personal', 'service', 'temporary') DEFAULT 'personal' COMMENT '密钥类型',
    description TEXT COMMENT '密钥描述',
    permissions JSON COMMENT '权限配置',
    rate_limit_per_minute INT DEFAULT 60 COMMENT '每分钟限流',
    rate_limit_per_hour INT DEFAULT 3600 COMMENT '每小时限流',
    allowed_ips JSON COMMENT '允许的IP列表',
    allowed_domains JSON COMMENT '允许的域名列表',
    expires_at TIMESTAMP NULL COMMENT '过期时间',
    last_used_at TIMESTAMP NULL COMMENT '最后使用时间',
    last_used_ip VARCHAR(64) COMMENT '最后使用IP',
    usage_count INT DEFAULT 0 COMMENT '使用次数',
    status ENUM('active', 'inactive', 'revoked', 'expired') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    revoked_at TIMESTAMP NULL COMMENT '撤销时间',
    UNIQUE KEY uk_api_key (api_key),
    INDEX idx_key_type (key_type),
    INDEX idx_status (status),
    INDEX idx_expires_at (expires_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'API密钥表';
-- 用户会话表
CREATE TABLE vt_user_sessions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    session_id VARCHAR(128) NOT NULL COMMENT '会话ID',
    device_type ENUM('web', 'mobile', 'desktop', 'api') DEFAULT 'web' COMMENT '设备类型',
    device_info VARCHAR(256) COMMENT '设备信息',
    ip_address VARCHAR(64) COMMENT 'IP地址',
    user_agent TEXT COMMENT '用户代理',
    login_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    last_activity_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '最后活动时间',
    expires_at TIMESTAMP NOT NULL COMMENT '过期时间',
    status ENUM('active', 'expired', 'revoked') DEFAULT 'active' COMMENT '状态',
    UNIQUE KEY uk_session_id (session_id),
    INDEX idx_device_type (device_type),
    INDEX idx_status (status),
    INDEX idx_expires_at (expires_at),
    INDEX idx_last_activity_at (last_activity_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户会话表';
-- 用户偏好设置表
CREATE TABLE vt_user_preferences (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    preference_key VARCHAR(128) NOT NULL COMMENT '设置键',
    preference_value TEXT COMMENT '设置值',
    value_type ENUM('string', 'number', 'boolean', 'json') DEFAULT 'string' COMMENT '值类型',
    category VARCHAR(64) COMMENT '设置分类',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_category (category)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户偏好设置表';
-- 系统配置表
CREATE TABLE vt_system_configs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    config_key VARCHAR(128) NOT NULL COMMENT '配置键',
    config_value TEXT COMMENT '配置值',
    config_type ENUM('string', 'number', 'boolean', 'json') DEFAULT 'string' COMMENT '配置类型',
    category VARCHAR(64) COMMENT '配置分类',
    description TEXT COMMENT '配置描述',
    is_encrypted TINYINT(1) DEFAULT 0 COMMENT '是否加密',
    is_readonly TINYINT(1) DEFAULT 0 COMMENT '是否只读',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_config_key (config_key),
    INDEX idx_category (category)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统配置表';
-- 通知消息表
CREATE TABLE vt_notifications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(256) NOT NULL COMMENT '通知标题',
    content TEXT COMMENT '通知内容',
    notification_type ENUM(
        'system',
        'training',
        'deployment',
        'alert',
        'workspace'
    ) NOT NULL COMMENT '通知类型',
    priority ENUM('low', 'normal', 'high', 'urgent') DEFAULT 'normal' COMMENT '优先级',
    status ENUM('unread', 'read', 'archived') DEFAULT 'unread' COMMENT '状态',
    resource_type VARCHAR(64) COMMENT '关联资源类型',
    resource_id BIGINT COMMENT '关联资源ID',
    metadata JSON COMMENT '扩展数据',
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
    read_at TIMESTAMP NULL COMMENT '阅读时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_notification_type (notification_type),
    INDEX idx_priority (priority),
    INDEX idx_status (status),
    INDEX idx_resource (resource_type, resource_id),
    INDEX idx_sent_at (sent_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '通知消息表';
-- 系统资源配额表
CREATE TABLE vt_resource_quotas (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    quota_type ENUM('user', 'workspace', 'global') NOT NULL COMMENT '配额类型',
    target_id BIGINT COMMENT '目标ID(用户ID或工作空间ID)',
    resource_type ENUM(
        'storage',
        'compute',
        'gpu',
        'memory',
        'dataset',
        'model'
    ) NOT NULL COMMENT '资源类型',
    quota_limit DECIMAL(20, 6) NOT NULL DEFAULT 0 COMMENT '配额限制',
    quota_used DECIMAL(20, 6) NOT NULL DEFAULT 0 COMMENT '已使用配额',
    unit VARCHAR(32) NOT NULL COMMENT '单位',
    reset_cycle ENUM('daily', 'weekly', 'monthly', 'never') DEFAULT 'never' COMMENT '重置周期',
    last_reset_at TIMESTAMP NULL COMMENT '最后重置时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_quota_target_resource (quota_type, target_id, resource_type),
    INDEX idx_quota_type (quota_type),
    INDEX idx_target_id (target_id),
    INDEX idx_resource_type (resource_type)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统资源配额表';
-- 统一审计日志表
CREATE TABLE vt_audit_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    resource_type VARCHAR(64) NOT NULL COMMENT '资源类型',
    resource_id BIGINT COMMENT '资源ID',
    action VARCHAR(64) NOT NULL COMMENT '操作动作',
    details JSON COMMENT '操作详情',
    ip_address VARCHAR(64) COMMENT 'IP地址',
    user_agent TEXT COMMENT '用户代理',
    status ENUM('success', 'failed') DEFAULT 'success' COMMENT '操作状态',
    error_message TEXT COMMENT '错误信息',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_resource (resource_type, resource_id),
    INDEX idx_action (action),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '统一审计日志表';
-- 操作日志表
CREATE TABLE vt_operation_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    operation_type ENUM(
        'create',
        'update',
        'delete',
        'login',
        'logout',
        'export',
        'import',
        'execute',
        'deploy',
        'start',
        'stop',
        'restart'
    ) NOT NULL COMMENT '操作类型',
    resource_type VARCHAR(64) COMMENT '资源类型',
    resource_id BIGINT COMMENT '资源ID',
    resource_name VARCHAR(256) COMMENT '资源名称',
    request_method VARCHAR(16) COMMENT '请求方法',
    request_path VARCHAR(512) COMMENT '请求路径',
    request_params JSON COMMENT '请求参数',
    request_body JSON COMMENT '请求体',
    response_code INT COMMENT '响应代码',
    response_time_ms INT COMMENT '响应时间(毫秒)',
    error_message TEXT COMMENT '错误信息',
    ip_address VARCHAR(64) COMMENT 'IP地址',
    user_agent TEXT COMMENT '用户代理',
    session_id VARCHAR(128) COMMENT '会话ID',
    old_value JSON COMMENT '旧值',
    new_value JSON COMMENT '新值',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_operation_type (operation_type),
    INDEX idx_resource (resource_type, resource_id),
    INDEX idx_created_at (created_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '操作日志表';
-- API密钥关联表
CREATE TABLE vt_api_key_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    api_key_id BIGINT NOT NULL COMMENT 'API密钥ID',
    entity_type VARCHAR(64) NOT NULL COMMENT '实体类型(user等)',
    entity_id BIGINT NOT NULL COMMENT '实体ID',
    relation_type VARCHAR(64) NOT NULL COMMENT '关联类型(owner, revoked_by等)',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_api_key_entity_relation (
        api_key_id,
        entity_type,
        entity_id,
        relation_type
    ),
    INDEX idx_api_key_id (api_key_id),
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_relation_type (relation_type),
    INDEX idx_status (status)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'API密钥关联表';
-- 用户会话关联表
CREATE TABLE vt_user_session_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    session_id BIGINT NOT NULL COMMENT '会话ID',
    entity_type VARCHAR(64) NOT NULL COMMENT '实体类型(user等)',
    entity_id BIGINT NOT NULL COMMENT '实体ID',
    relation_type VARCHAR(64) NOT NULL COMMENT '关联类型(user等)',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_session_entity_relation (
        session_id,
        entity_type,
        entity_id,
        relation_type
    ),
    INDEX idx_session_id (session_id),
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_relation_type (relation_type),
    INDEX idx_status (status)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户会话关联表';
-- 用户偏好设置关联表
CREATE TABLE vt_user_preference_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    preference_id BIGINT NOT NULL COMMENT '偏好设置ID',
    entity_type VARCHAR(64) NOT NULL COMMENT '实体类型(user等)',
    entity_id BIGINT NOT NULL COMMENT '实体ID',
    relation_type VARCHAR(64) NOT NULL COMMENT '关联类型(user等)',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_preference_entity_relation (
        preference_id,
        entity_type,
        entity_id,
        relation_type
    ),
    INDEX idx_preference_id (preference_id),
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_relation_type (relation_type),
    INDEX idx_status (status)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户偏好设置关联表';
-- 系统配置关联表
CREATE TABLE vt_system_config_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    config_id BIGINT NOT NULL COMMENT '系统配置ID',
    entity_type VARCHAR(64) NOT NULL COMMENT '实体类型(user等)',
    entity_id BIGINT NOT NULL COMMENT '实体ID',
    relation_type VARCHAR(64) NOT NULL COMMENT '关联类型(updated_by等)',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_config_entity_relation (config_id, entity_type, entity_id, relation_type),
    INDEX idx_config_id (config_id),
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_relation_type (relation_type),
    INDEX idx_status (status)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统配置关联表';
-- 通知消息关联表
CREATE TABLE vt_notification_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    notification_id BIGINT NOT NULL COMMENT '通知消息ID',
    entity_type VARCHAR(64) NOT NULL COMMENT '实体类型(user等)',
    entity_id BIGINT NOT NULL COMMENT '实体ID',
    relation_type VARCHAR(64) NOT NULL COMMENT '关联类型(recipient等)',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_notification_entity_relation (
        notification_id,
        entity_type,
        entity_id,
        relation_type
    ),
    INDEX idx_notification_id (notification_id),
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_relation_type (relation_type),
    INDEX idx_status (status)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '通知消息关联表';
-- 审计日志关联表
CREATE TABLE vt_audit_log_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    audit_log_id BIGINT NOT NULL COMMENT '审计日志ID',
    entity_type VARCHAR(64) NOT NULL COMMENT '实体类型(user等)',
    entity_id BIGINT NOT NULL COMMENT '实体ID',
    relation_type VARCHAR(64) NOT NULL COMMENT '关联类型(user等)',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_audit_entity_relation (
        audit_log_id,
        entity_type,
        entity_id,
        relation_type
    ),
    INDEX idx_audit_log_id (audit_log_id),
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_relation_type (relation_type),
    INDEX idx_status (status)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '审计日志关联表';
-- 操作日志关联表
CREATE TABLE vt_operation_log_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    operation_log_id BIGINT NOT NULL COMMENT '操作日志ID',
    entity_type VARCHAR(64) NOT NULL COMMENT '实体类型(user, api_key等)',
    entity_id BIGINT NOT NULL COMMENT '实体ID',
    relation_type VARCHAR(64) NOT NULL COMMENT '关联类型(user, api_key等)',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_operation_entity_relation (
        operation_log_id,
        entity_type,
        entity_id,
        relation_type
    ),
    INDEX idx_operation_log_id (operation_log_id),
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_relation_type (relation_type),
    INDEX idx_status (status)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '操作日志关联表';