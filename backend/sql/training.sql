-- 训练作业表
CREATE TABLE vt_training_jobs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL COMMENT '任务名称',
    display_name VARCHAR(256) COMMENT '显示名称',
    description TEXT COMMENT '任务描述',
    workspace_id BIGINT NOT NULL COMMENT '所属工作空间',
    creator_id BIGINT NOT NULL COMMENT '创建者ID',
    -- 任务配置
    job_type ENUM(
        'single',
        'distributed',
        'horovod',
        'pytorch_ddp'
    ) DEFAULT 'single' COMMENT '任务类型',
    framework ENUM(
        'pytorch',
        'tensorflow',
        'paddle',
        'mindspore',
        'custom'
    ) NOT NULL COMMENT '训练框架',
    framework_version VARCHAR(32) COMMENT '框架版本',
    python_version VARCHAR(16) DEFAULT '3.8' COMMENT 'Python版本',
    -- 镜像和代码配置
    image VARCHAR(512) NOT NULL COMMENT '训练镜像',
    code_source_type ENUM('git', 'upload', 'builtin') DEFAULT 'upload' COMMENT '代码来源类型',
    code_source_config JSON COMMENT '代码来源配置',
    entry_point VARCHAR(256) NOT NULL COMMENT '入口脚本',
    working_dir VARCHAR(256) DEFAULT '/workspace' COMMENT '工作目录',
    -- 数据配置
    dataset_ids JSON COMMENT '关联数据集ID列表',
    dataset_mount_path VARCHAR(256) DEFAULT '/data' COMMENT '数据集挂载路径',
    -- 模型配置
    base_model_id BIGINT COMMENT '基础模型ID',
    model_mount_path VARCHAR(256) DEFAULT '/models' COMMENT '模型挂载路径',
    output_model_name VARCHAR(128) COMMENT '输出模型名称',
    -- 资源配置
    resource_spec JSON NOT NULL COMMENT '资源需求规格',
    node_selector JSON COMMENT '节点选择器',
    tolerations JSON COMMENT '容忍度设置',
    -- 运行时配置
    env_vars JSON COMMENT '环境变量',
    command_args JSON COMMENT '命令行参数',
    max_runtime_seconds INT DEFAULT 86400 COMMENT '最大运行时间(秒)',
    auto_restart TINYINT(1) DEFAULT 0 COMMENT '是否自动重启',
    max_retry_count INT DEFAULT 3 COMMENT '最大重试次数',
    -- 调度配置
    priority INT DEFAULT 0 COMMENT '优先级(0-100)',
    queue_name VARCHAR(64) DEFAULT 'default' COMMENT '队列名称',
    volcano_job_name VARCHAR(128) COMMENT 'Volcano作业名',
    -- 状态信息
    status ENUM(
        'pending',
        'queued',
        'running',
        'succeeded',
        'failed',
        'cancelled',
        'suspended',
        'timeout'
    ) DEFAULT 'pending' COMMENT '任务状态',
    phase ENUM(
        'creating',
        'scheduling',
        'training',
        'completed'
    ) DEFAULT 'creating' COMMENT '执行阶段',
    error_message TEXT COMMENT '错误信息',
    exit_code INT COMMENT '退出码',
    -- 时间信息
    submitted_at BIGINT DEFAULT 0 COMMENT '提交时间戳',
    scheduled_at BIGINT DEFAULT 0 COMMENT '调度时间戳',
    start_time BIGINT DEFAULT 0 COMMENT '开始时间戳',
    end_time BIGINT DEFAULT 0 COMMENT '结束时间戳',
    duration_seconds INT COMMENT '执行时长(秒)',
    -- 存储路径
    workspace_path VARCHAR(512) COMMENT '工作空间路径',
    logs_path VARCHAR(512) COMMENT '日志路径',
    output_path VARCHAR(512) COMMENT '输出路径',
    checkpoint_path VARCHAR(512) COMMENT '检查点路径',
    -- 扩展信息
    tags JSON COMMENT '标签',
    annotations JSON COMMENT '注解',
    metadata JSON COMMENT '元数据',
    -- 审计字段
    created_at BIGINT DEFAULT 0 COMMENT '创建时间戳',
    updated_at BIGINT DEFAULT 0 COMMENT '更新时间戳',
    deleted_at BIGINT DEFAULT 0 COMMENT '软删除时间戳',
    -- 索引
    INDEX idx_workspace_id (workspace_id),
    INDEX idx_creator_id (creator_id),
    INDEX idx_status (status),
    INDEX idx_framework (framework),
    INDEX idx_volcano_job (volcano_job_name),
    INDEX idx_submitted_at (submitted_at),
    INDEX idx_priority (priority),
    INDEX idx_queue (queue_name),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_created_at (created_at),
    INDEX idx_updated_at (updated_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练作业表';
-- 训练任务实例表 (Pod级别)
CREATE TABLE vt_training_job_instances (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_id BIGINT NOT NULL COMMENT '所属训练作业ID',
    -- 实例标识
    instance_name VARCHAR(128) NOT NULL COMMENT '实例名称',
    instance_type ENUM('master', 'worker', 'ps', 'evaluator') DEFAULT 'worker' COMMENT '实例类型',
    instance_index INT NOT NULL COMMENT '实例索引',
    replica_index INT COMMENT '副本索引',
    -- K8s信息
    pod_name VARCHAR(128) COMMENT 'Pod名称',
    namespace VARCHAR(128) COMMENT '命名空间',
    node_name VARCHAR(128) COMMENT '所在节点',
    node_ip VARCHAR(64) COMMENT '节点IP',
    pod_ip VARCHAR(64) COMMENT 'Pod IP',
    -- 资源分配
    allocated_resources JSON COMMENT '分配的资源',
    gpu_cards JSON COMMENT '分配的GPU卡信息',
    cpu_cores DECIMAL(4, 2) COMMENT '分配的CPU核心数',
    memory_mb INT COMMENT '分配的内存(MB)',
    -- 状态信息
    status ENUM(
        'pending',
        'creating',
        'running',
        'succeeded',
        'failed',
        'killed',
        'unknown'
    ) DEFAULT 'pending' COMMENT '实例状态',
    phase VARCHAR(32) COMMENT 'Pod阶段',
    reason VARCHAR(128) COMMENT '状态原因',
    message TEXT COMMENT '状态消息',
    -- 时间信息
    created_at BIGINT DEFAULT 0 COMMENT '创建时间戳',
    updated_at BIGINT DEFAULT 0 COMMENT '更新时间戳',
    scheduled_at BIGINT DEFAULT 0 COMMENT '调度时间戳',
    start_time BIGINT DEFAULT 0 COMMENT '开始时间戳',
    end_time BIGINT DEFAULT 0 COMMENT '结束时间戳',
    last_transition_time BIGINT DEFAULT 0 COMMENT '最后状态变更时间戳',
    -- 重启信息
    restart_count INT DEFAULT 0 COMMENT '重启次数',
    last_restart_time BIGINT DEFAULT 0 COMMENT '最后重启时间戳',
    exit_code INT COMMENT '退出码',
    -- 存储路径
    logs_path VARCHAR(512) COMMENT '日志路径',
    -- 扩展信息
    labels JSON COMMENT '标签',
    annotations JSON COMMENT '注解',
    -- 索引
    INDEX idx_job_id (job_id),
    INDEX idx_pod_name (pod_name),
    INDEX idx_node_name (node_name),
    INDEX idx_status (status),
    INDEX idx_instance_type (instance_type),
    INDEX idx_created_at (created_at),
    INDEX idx_updated_at (updated_at),
    UNIQUE KEY uk_job_instance (job_id, instance_name)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练任务实例表';
-- 训练作业实例关联表
CREATE TABLE vt_training_job_instance_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_id BIGINT NOT NULL COMMENT '训练作业ID',
    instance_id BIGINT NOT NULL COMMENT '实例ID',
    INDEX idx_job_id (job_id),
    INDEX idx_instance_id (instance_id),
    UNIQUE KEY uk_job_instance (job_id, instance_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练作业实例关联表';
-- 训练指标表
CREATE TABLE vt_training_metrics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_id BIGINT NOT NULL COMMENT '训练作业ID',
    instance_id BIGINT COMMENT '实例ID',
    -- 指标信息
    metric_name VARCHAR(128) NOT NULL COMMENT '指标名称',
    metric_type ENUM('scalar', 'histogram', 'image', 'text', 'audio') DEFAULT 'scalar' COMMENT '指标类型',
    metric_value DECIMAL(20, 6) COMMENT '指标值',
    metric_data JSON COMMENT '复杂指标数据',
    -- 训练进度
    step BIGINT COMMENT '训练步数',
    epoch INT COMMENT '训练轮次',
    global_step BIGINT COMMENT '全局步数',
    -- 时间戳
    metric_time BIGINT DEFAULT 0 COMMENT '指标时间戳',
    wall_time DECIMAL(15, 3) COMMENT '墙钟时间',
    created_at BIGINT DEFAULT 0 COMMENT '创建时间戳',
    -- 分组信息
    tag VARCHAR(128) COMMENT '标签',
    category VARCHAR(64) COMMENT '类别',
    -- 索引
    INDEX idx_job_metric (job_id, metric_name),
    INDEX idx_job_step (job_id, step),
    INDEX idx_job_epoch (job_id, epoch),
    INDEX idx_metric_time (metric_time),
    INDEX idx_instance_id (instance_id),
    INDEX idx_created_at (created_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练指标表';
-- 训练指标关联表
CREATE TABLE vt_training_metric_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_id BIGINT NOT NULL COMMENT '训练作业ID',
    metric_id BIGINT NOT NULL COMMENT '指标ID',
    instance_id BIGINT COMMENT '实例ID',
    INDEX idx_job_id (job_id),
    INDEX idx_metric_id (metric_id),
    INDEX idx_instance_id (instance_id),
    UNIQUE KEY uk_job_metric (job_id, metric_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练指标关联表';
-- 训练日志表
CREATE TABLE vt_training_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_id BIGINT NOT NULL COMMENT '训练作业ID',
    instance_id BIGINT COMMENT '实例ID',
    -- 日志信息
    log_level ENUM('DEBUG', 'INFO', 'WARN', 'ERROR', 'FATAL') DEFAULT 'INFO' COMMENT '日志级别',
    log_source VARCHAR(64) COMMENT '日志来源',
    log_content TEXT NOT NULL COMMENT '日志内容',
    -- 时间信息
    log_time BIGINT DEFAULT 0 COMMENT '日志时间戳',
    created_at BIGINT DEFAULT 0 COMMENT '创建时间戳',
    -- 位置信息
    file_name VARCHAR(256) COMMENT '文件名',
    line_number INT COMMENT '行号',
    function_name VARCHAR(128) COMMENT '函数名',
    -- 索引
    INDEX idx_job_id (job_id),
    INDEX idx_instance_id (instance_id),
    INDEX idx_log_level (log_level),
    INDEX idx_log_time (log_time),
    INDEX idx_created_at (created_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练日志表';
-- 训练日志关联表
CREATE TABLE vt_training_log_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_id BIGINT NOT NULL COMMENT '训练作业ID',
    log_id BIGINT NOT NULL COMMENT '日志ID',
    instance_id BIGINT COMMENT '实例ID',
    INDEX idx_job_id (job_id),
    INDEX idx_log_id (log_id),
    INDEX idx_instance_id (instance_id),
    UNIQUE KEY uk_job_log (job_id, log_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练日志关联表';
-- 训练检查点表
CREATE TABLE vt_training_checkpoints (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_id BIGINT NOT NULL COMMENT '训练作业ID',
    -- 检查点信息
    checkpoint_name VARCHAR(128) NOT NULL COMMENT '检查点名称',
    checkpoint_type ENUM('auto', 'manual', 'best', 'final') DEFAULT 'auto' COMMENT '检查点类型',
    step BIGINT COMMENT '训练步数',
    epoch INT COMMENT '训练轮次',
    -- 存储信息
    storage_path VARCHAR(512) NOT NULL COMMENT '存储路径',
    file_size BIGINT COMMENT '文件大小(字节)',
    checksum VARCHAR(64) COMMENT '文件校验和',
    -- 性能指标
    metrics JSON COMMENT '性能指标',
    loss_value DECIMAL(15, 6) COMMENT '损失值',
    accuracy DECIMAL(8, 6) COMMENT '准确率',
    -- 状态信息
    status ENUM('saving', 'saved', 'failed', 'deleted') DEFAULT 'saving' COMMENT '状态',
    -- 时间信息
    created_at BIGINT DEFAULT 0 COMMENT '创建时间戳',
    updated_at BIGINT DEFAULT 0 COMMENT '更新时间戳',
    saved_at BIGINT DEFAULT 0 COMMENT '保存完成时间戳',
    -- 索引
    INDEX idx_job_id (job_id),
    INDEX idx_checkpoint_type (checkpoint_type),
    INDEX idx_step (step),
    INDEX idx_epoch (epoch),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_updated_at (updated_at),
    UNIQUE KEY uk_job_checkpoint (job_id, checkpoint_name)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练检查点表';
-- 训练检查点关联表
CREATE TABLE vt_training_checkpoint_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_id BIGINT NOT NULL COMMENT '训练作业ID',
    checkpoint_id BIGINT NOT NULL COMMENT '检查点ID',
    INDEX idx_job_id (job_id),
    INDEX idx_checkpoint_id (checkpoint_id),
    UNIQUE KEY uk_job_checkpoint (job_id, checkpoint_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练检查点关联表';
-- 训练模板表
CREATE TABLE vt_training_templates (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL COMMENT '模板名称',
    display_name VARCHAR(256) COMMENT '显示名称',
    description TEXT COMMENT '模板描述',
    -- 模板分类
    category VARCHAR(64) COMMENT '模板类别',
    framework ENUM(
        'pytorch',
        'tensorflow',
        'paddle',
        'mindspore',
        'custom'
    ) NOT NULL COMMENT '框架',
    framework_version VARCHAR(32) COMMENT '框架版本',
    -- 模板配置
    template_config JSON NOT NULL COMMENT '模板配置',
    default_image VARCHAR(512) COMMENT '默认镜像',
    default_resources JSON COMMENT '默认资源配置',
    -- 参数定义
    parameters JSON COMMENT '参数定义',
    parameter_schema JSON COMMENT '参数Schema',
    -- 访问控制
    visibility ENUM('public', 'private', 'workspace') DEFAULT 'private' COMMENT '可见性',
    workspace_id BIGINT COMMENT '所属工作空间',
    creator_id BIGINT NOT NULL COMMENT '创建者',
    -- 使用统计
    usage_count INT DEFAULT 0 COMMENT '使用次数',
    star_count INT DEFAULT 0 COMMENT '收藏次数',
    -- 状态信息
    status ENUM('active', 'deprecated', 'deleted') DEFAULT 'active' COMMENT '状态',
    version VARCHAR(32) DEFAULT 'v1.0' COMMENT '版本',
    -- 扩展信息
    tags JSON COMMENT '标签',
    readme TEXT COMMENT '说明文档',
    created_at BIGINT DEFAULT 0 COMMENT '创建时间戳',
    updated_at BIGINT DEFAULT 0 COMMENT '更新时间戳',
    deleted_at BIGINT DEFAULT 0 COMMENT '软删除时间戳',
    -- 索引
    INDEX idx_framework (framework),
    INDEX idx_category (category),
    INDEX idx_visibility (visibility),
    INDEX idx_workspace_id (workspace_id),
    INDEX idx_creator_id (creator_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_updated_at (updated_at),
    INDEX idx_deleted_at (deleted_at),
    UNIQUE KEY uk_name_version (name, version)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练模板表';
-- 训练队列表
CREATE TABLE vt_training_queues (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL COMMENT '队列名称',
    display_name VARCHAR(256) COMMENT '显示名称',
    description TEXT COMMENT '队列描述',
    -- 队列配置
    priority INT DEFAULT 0 COMMENT '队列优先级',
    max_concurrent_jobs INT DEFAULT 10 COMMENT '最大并发任务数',
    max_queue_size INT DEFAULT 100 COMMENT '最大队列大小',
    -- 资源配额
    resource_quota JSON COMMENT '资源配额',
    gpu_quota INT COMMENT 'GPU配额',
    cpu_quota DECIMAL(8, 2) COMMENT 'CPU配额',
    memory_quota_gb INT COMMENT '内存配额(GB)',
    -- 调度策略
    scheduling_policy ENUM('fifo', 'priority', 'fair_share') DEFAULT 'fifo' COMMENT '调度策略',
    preemption_enabled TINYINT(1) DEFAULT 0 COMMENT '是否启用抢占',
    -- 访问控制
    workspace_ids JSON COMMENT '允许访问的工作空间ID列表',
    user_ids JSON COMMENT '允许访问的用户ID列表',
    -- 状态信息
    status ENUM('active', 'disabled', 'maintenance') DEFAULT 'active' COMMENT '状态',
    current_jobs INT DEFAULT 0 COMMENT '当前任务数',
    pending_jobs INT DEFAULT 0 COMMENT '等待任务数',
    created_at BIGINT DEFAULT 0 COMMENT '创建时间戳',
    updated_at BIGINT DEFAULT 0 COMMENT '更新时间戳',
    deleted_at BIGINT DEFAULT 0 COMMENT '软删除时间戳',
    INDEX idx_status (status),
    INDEX idx_priority (priority),
    INDEX idx_created_at (created_at),
    INDEX idx_updated_at (updated_at),
    INDEX idx_deleted_at (deleted_at),
    UNIQUE KEY uk_name (name)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练队列表';
-- 工作空间关联表
CREATE TABLE vt_workspace_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    workspace_id BIGINT NOT NULL COMMENT '工作空间ID',
    resource_type ENUM('job', 'template', 'queue') NOT NULL COMMENT '资源类型',
    resource_id BIGINT NOT NULL COMMENT '资源ID',
    INDEX idx_workspace_id (workspace_id),
    INDEX idx_resource_type (resource_type),
    INDEX idx_resource_id (resource_id),
    UNIQUE KEY uk_workspace_resource (workspace_id, resource_type, resource_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '工作空间关联表';
-- 用户关联表
CREATE TABLE vt_user_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    resource_type ENUM('job', 'template', 'queue') NOT NULL COMMENT '资源类型',
    resource_id BIGINT NOT NULL COMMENT '资源ID',
    relation_type ENUM('owner', 'creator', 'collaborator') NOT NULL COMMENT '关联类型',
    INDEX idx_user_id (user_id),
    INDEX idx_resource_type (resource_type),
    INDEX idx_resource_id (resource_id),
    INDEX idx_relation_type (relation_type),
    UNIQUE KEY uk_user_resource (
        user_id,
        resource_type,
        resource_id,
        relation_type
    )
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户关联表';