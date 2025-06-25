-- K8s资源映射表
CREATE TABLE vt_k8s_resources (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    resource_type VARCHAR(50) NOT NULL COMMENT '资源类型：Job, Pod, Service, ConfigMap, Secret等',
    resource_name VARCHAR(255) NOT NULL COMMENT 'K8s资源名称',
    namespace VARCHAR(255) NOT NULL DEFAULT 'vt-platform' COMMENT 'K8s命名空间',
    uid VARCHAR(255) NOT NULL COMMENT 'K8s资源UID',
    api_version VARCHAR(50) NOT NULL COMMENT 'API版本',
    kind VARCHAR(50) NOT NULL COMMENT '资源种类',
    cluster_name VARCHAR(128) COMMENT '所属集群名称',
    related_type VARCHAR(50) COMMENT '关联资源类型：training_job, model_deployment等',
    related_id BIGINT COMMENT '关联资源ID',
    labels JSON COMMENT 'K8s标签',
    annotations JSON COMMENT 'K8s注解',
    spec JSON COMMENT '资源规格',
    status JSON COMMENT '资源状态',
    k8s_created_at TIMESTAMP NULL COMMENT 'K8s创建时间',
    k8s_deleted_at TIMESTAMP NULL COMMENT 'K8s删除时间',
    last_sync_at TIMESTAMP NULL COMMENT '最后同步时间',
    sync_status ENUM('synced', 'pending', 'failed') DEFAULT 'pending' COMMENT '同步状态',
    sync_error TEXT COMMENT '同步错误信息',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_namespace_name (namespace, resource_name, resource_type),
    UNIQUE KEY uk_uid (uid),
    INDEX idx_related (related_type, related_id),
    INDEX idx_sync_status (sync_status, last_sync_at),
    INDEX idx_k8s_created (k8s_created_at),
    INDEX idx_cluster_name (cluster_name)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'K8s资源映射表';
-- Volcano作业详情表
CREATE TABLE vt_volcano_jobs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_name VARCHAR(255) NOT NULL COMMENT 'Volcano作业名称',
    namespace VARCHAR(255) NOT NULL DEFAULT 'default' COMMENT 'K8s命名空间',
    uid VARCHAR(255) NOT NULL COMMENT 'K8s UID',
    training_job_id BIGINT NOT NULL COMMENT '训练任务ID',
    queue_name VARCHAR(100) DEFAULT 'default' COMMENT 'Volcano队列名称',
    priority INT DEFAULT 0 COMMENT '优先级',
    min_available INT DEFAULT 1 COMMENT '最小可用实例数',
    scheduling_policy JSON COMMENT '调度策略',
    plugins JSON COMMENT '插件配置',
    scheduler_name VARCHAR(100) DEFAULT 'volcano' COMMENT '调度器名称',
    task_specs JSON COMMENT '任务规格列表',
    volumes JSON COMMENT '卷配置',
    ttl_seconds_after_finished INT COMMENT '完成后保留时间(秒)',
    active_deadline_seconds INT COMMENT '活跃截止时间(秒)',
    backoff_limit INT DEFAULT 3 COMMENT '重试次数限制',
    phase VARCHAR(50) COMMENT '作业阶段：Pending, Running, Completed, Failed等',
    conditions JSON COMMENT '作业条件',
    status_message TEXT COMMENT '状态消息',
    start_time TIMESTAMP NULL COMMENT '开始时间',
    completion_time TIMESTAMP NULL COMMENT '完成时间',
    running_duration INT COMMENT '运行时长(秒)',
    total_cpu_request DECIMAL(10, 2) COMMENT '总CPU请求',
    total_memory_request BIGINT COMMENT '总内存请求(bytes)',
    total_gpu_request INT COMMENT '总GPU请求',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_namespace_name (namespace, job_name),
    UNIQUE KEY uk_uid (uid),
    INDEX idx_training_job_id (training_job_id),
    INDEX idx_queue_name (queue_name),
    INDEX idx_phase (phase),
    INDEX idx_start_time (start_time),
    INDEX idx_scheduler_name (scheduler_name)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'Volcano作业详情表';
-- Volcano队列配置表
CREATE TABLE vt_volcano_queues (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    queue_name VARCHAR(100) NOT NULL COMMENT '队列名称',
    namespace VARCHAR(255) NOT NULL DEFAULT 'volcano-system' COMMENT '命名空间',
    weight INT DEFAULT 1 COMMENT '权重',
    capability JSON COMMENT '队列容量配置',
    cpu_quota DECIMAL(10, 2) COMMENT 'CPU配额',
    memory_quota BIGINT COMMENT '内存配额(bytes)',
    gpu_quota INT COMMENT 'GPU配额',
    state ENUM('Open', 'Closed') DEFAULT 'Open' COMMENT '队列状态',
    running_jobs INT DEFAULT 0 COMMENT '运行中的作业数',
    pending_jobs INT DEFAULT 0 COMMENT '等待中的作业数',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_queue_name (queue_name, namespace),
    INDEX idx_state (state)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'Volcano队列配置表';
-- Pod事件记录表
CREATE TABLE vt_pod_events (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    pod_name VARCHAR(255) NOT NULL COMMENT 'Pod名称',
    namespace VARCHAR(255) NOT NULL COMMENT '命名空间',
    uid VARCHAR(255) COMMENT 'Pod UID',
    event_type VARCHAR(50) NOT NULL COMMENT '事件类型：Normal, Warning',
    reason VARCHAR(100) NOT NULL COMMENT '事件原因',
    message TEXT COMMENT '事件消息',
    source_component VARCHAR(100) COMMENT '来源组件',
    source_host VARCHAR(255) COMMENT '来源主机',
    first_timestamp TIMESTAMP NOT NULL COMMENT '首次发生时间',
    last_timestamp TIMESTAMP NOT NULL COMMENT '最后发生时间',
    count INT DEFAULT 1 COMMENT '发生次数',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_pod (namespace, pod_name),
    INDEX idx_event_type (event_type),
    INDEX idx_timestamp (last_timestamp)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'Pod事件记录表';
-- 资源关联表
CREATE TABLE vt_resource_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    source_type VARCHAR(64) NOT NULL COMMENT '源资源类型',
    source_id BIGINT NOT NULL COMMENT '源资源ID',
    target_type VARCHAR(64) NOT NULL COMMENT '目标资源类型',
    target_id BIGINT NOT NULL COMMENT '目标资源ID',
    relation_type ENUM(
        'owns',
        'uses',
        'depends_on',
        'produces',
        'contains',
        'references'
    ) NOT NULL COMMENT '关联类型',
    metadata JSON COMMENT '关联元数据',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_resource_relation (
        source_type,
        source_id,
        target_type,
        target_id,
        relation_type
    ),
    INDEX idx_source (source_type, source_id),
    INDEX idx_target (target_type, target_id),
    INDEX idx_relation_type (relation_type)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '资源关联表';
-- 用户收藏表
CREATE TABLE vt_user_favorites (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    resource_type ENUM(
        'dataset',
        'model',
        'template',
        'dashboard',
        'workspace'
    ) NOT NULL COMMENT '资源类型',
    resource_id BIGINT NOT NULL COMMENT '资源ID',
    folder_name VARCHAR(128) COMMENT '收藏夹名称',
    notes TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_resource (resource_type, resource_id),
    INDEX idx_folder_name (folder_name)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户收藏表';
-- 资源标签关联表
CREATE TABLE vt_resource_tags (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    resource_type VARCHAR(64) NOT NULL COMMENT '资源类型',
    resource_id BIGINT NOT NULL COMMENT '资源ID',
    tag_name VARCHAR(128) NOT NULL COMMENT '标签名称',
    tag_value VARCHAR(256) COMMENT '标签值',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_resource_tag (resource_type, resource_id, tag_name),
    INDEX idx_resource (resource_type, resource_id),
    INDEX idx_tag_name (tag_name)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '资源标签关联表';
-- K8s资源关联表
CREATE TABLE vt_k8s_resource_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    k8s_resource_id BIGINT NOT NULL COMMENT 'K8s资源ID',
    entity_type VARCHAR(64) NOT NULL COMMENT '实体类型(workspace, training_job等)',
    entity_id BIGINT NOT NULL COMMENT '实体ID',
    relation_type VARCHAR(64) NOT NULL COMMENT '关联类型(workspace, training_job等)',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_k8s_entity_relation (
        k8s_resource_id,
        entity_type,
        entity_id,
        relation_type
    ),
    INDEX idx_k8s_resource_id (k8s_resource_id),
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_relation_type (relation_type),
    INDEX idx_status (status)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'K8s资源关联表';
-- Volcano作业关联表
CREATE TABLE vt_volcano_job_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    volcano_job_id BIGINT NOT NULL COMMENT 'Volcano作业ID',
    entity_type VARCHAR(64) NOT NULL COMMENT '实体类型(k8s_resource, training_job等)',
    entity_id BIGINT NOT NULL COMMENT '实体ID',
    relation_type VARCHAR(64) NOT NULL COMMENT '关联类型(k8s_resource, training_job等)',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_volcano_entity_relation (
        volcano_job_id,
        entity_type,
        entity_id,
        relation_type
    ),
    INDEX idx_volcano_job_id (volcano_job_id),
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_relation_type (relation_type),
    INDEX idx_status (status)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'Volcano作业关联表';
-- Pod事件关联表
CREATE TABLE vt_pod_event_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    pod_event_id BIGINT NOT NULL COMMENT 'Pod事件ID',
    entity_type VARCHAR(64) NOT NULL COMMENT '实体类型(training_job, k8s_resource等)',
    entity_id BIGINT NOT NULL COMMENT '实体ID',
    relation_type VARCHAR(64) NOT NULL COMMENT '关联类型(training_job, k8s_resource等)',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_pod_event_entity_relation (
        pod_event_id,
        entity_type,
        entity_id,
        relation_type
    ),
    INDEX idx_pod_event_id (pod_event_id),
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_relation_type (relation_type),
    INDEX idx_status (status)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'Pod事件关联表';
-- 用户收藏关联表
CREATE TABLE vt_user_favorite_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    favorite_id BIGINT NOT NULL COMMENT '收藏ID',
    entity_type VARCHAR(64) NOT NULL COMMENT '实体类型(user等)',
    entity_id BIGINT NOT NULL COMMENT '实体ID',
    relation_type VARCHAR(64) NOT NULL COMMENT '关联类型(user等)',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_favorite_entity_relation (
        favorite_id,
        entity_type,
        entity_id,
        relation_type
    ),
    INDEX idx_favorite_id (favorite_id),
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_relation_type (relation_type),
    INDEX idx_status (status)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户收藏关联表';
-- 资源标签关联表
CREATE TABLE vt_resource_tag_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    resource_tag_id BIGINT NOT NULL COMMENT '资源标签ID',
    entity_type VARCHAR(64) NOT NULL COMMENT '实体类型(user等)',
    entity_id BIGINT NOT NULL COMMENT '实体ID',
    relation_type VARCHAR(64) NOT NULL COMMENT '关联类型(tagged_by等)',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_tag_entity_relation (
        resource_tag_id,
        entity_type,
        entity_id,
        relation_type
    ),
    INDEX idx_resource_tag_id (resource_tag_id),
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_relation_type (relation_type),
    INDEX idx_status (status)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '资源标签关联表';