-- 训练队列表
CREATE TABLE vt_training_queues (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL COMMENT '队列名称',
    display_name VARCHAR(256) COMMENT '显示名称',
    description TEXT COMMENT '队列描述',
    queue_type ENUM(
        'default',
        'high_priority',
        'gpu_intensive',
        'cpu_intensive',
        'experimental'
    ) DEFAULT 'default' COMMENT '队列类型',
    priority INT DEFAULT 0 COMMENT '队列优先级(数值越大优先级越高)',
    max_concurrent_jobs INT DEFAULT 10 COMMENT '最大并发任务数',
    max_queue_size INT DEFAULT 100 COMMENT '最大队列大小',
    max_job_duration_hours INT DEFAULT 168 COMMENT '最大任务时长(小时)',
    resource_quota JSON COMMENT '资源配额配置',
    gpu_quota INT COMMENT 'GPU配额',
    cpu_quota DECIMAL(8, 2) COMMENT 'CPU配额',
    memory_quota_gb INT COMMENT '内存配额(GB)',
    storage_quota_gb INT COMMENT '存储配额(GB)',
    scheduling_policy ENUM(
        'fifo',
        'priority',
        'fair_share',
        'shortest_job_first'
    ) DEFAULT 'fifo' COMMENT '调度策略',
    preemption_enabled TINYINT(1) DEFAULT 0 COMMENT '是否启用抢占',
    gang_scheduling TINYINT(1) DEFAULT 0 COMMENT '是否启用gang调度',
    workspace_ids JSON COMMENT '允许访问的工作空间ID列表',
    user_ids JSON COMMENT '允许访问的用户ID列表',
    department_ids JSON COMMENT '允许访问的部门ID列表',
    cluster_ids JSON COMMENT '可用集群ID列表',
    node_selector JSON COMMENT '节点选择器',
    tolerations JSON COMMENT '容忍度配置',
    status ENUM('active', 'disabled', 'maintenance') DEFAULT 'active' COMMENT '状态',
    current_jobs INT DEFAULT 0 COMMENT '当前任务数',
    pending_jobs INT DEFAULT 0 COMMENT '等待任务数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    UNIQUE KEY uk_name (name),
    INDEX idx_queue_type (queue_type),
    INDEX idx_priority (priority),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练队列表';
-- 训练作业表
CREATE TABLE vt_training_jobs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(128) NOT NULL COMMENT '任务名称',
    display_name VARCHAR(256) COMMENT '显示名称',
    description TEXT COMMENT '任务描述',
    job_type ENUM(
        'single',
        'distributed',
        'horovod',
        'pytorch_ddp',
        'parameter_server',
        'federated'
    ) DEFAULT 'single' COMMENT '任务类型',
    framework ENUM(
        'pytorch',
        'tensorflow',
        'paddle',
        'mindspore',
        'keras',
        'sklearn',
        'xgboost',
        'custom'
    ) NOT NULL COMMENT '训练框架',
    framework_version VARCHAR(32) COMMENT '框架版本',
    python_version VARCHAR(16) DEFAULT '3.8' COMMENT 'Python版本',
    code_source_type ENUM('git', 'upload', 'builtin', 'notebook') DEFAULT 'upload' COMMENT '代码来源类型',
    code_source_config JSON COMMENT '代码来源配置',
    entry_point VARCHAR(256) NOT NULL COMMENT '入口脚本',
    working_dir VARCHAR(256) DEFAULT '/workspace' COMMENT '工作目录',
    image VARCHAR(512) NOT NULL COMMENT '训练镜像',
    image_pull_policy ENUM('Always', 'IfNotPresent', 'Never') DEFAULT 'IfNotPresent' COMMENT '镜像拉取策略',
    image_pull_secrets JSON COMMENT '镜像拉取密钥',
    dataset_mount_configs JSON COMMENT '数据集挂载配置',
    data_source_config JSON COMMENT '数据源配置',
    model_config JSON COMMENT '模型配置',
    output_model_name VARCHAR(128) COMMENT '输出模型名称',
    model_save_strategy ENUM('all', 'best', 'last', 'custom') DEFAULT 'best' COMMENT '模型保存策略',
    cpu_cores DECIMAL(6, 3) COMMENT 'CPU核心数',
    memory_gb DECIMAL(8, 2) COMMENT '内存(GB)',
    gpu_count INT DEFAULT 0 COMMENT 'GPU数量',
    gpu_type VARCHAR(64) COMMENT 'GPU类型',
    gpu_memory_gb DECIMAL(8, 2) COMMENT 'GPU显存(GB)',
    storage_gb DECIMAL(8, 2) COMMENT '存储(GB)',
    shared_memory_gb DECIMAL(8, 2) COMMENT '共享内存(GB)',
    worker_count INT DEFAULT 1 COMMENT 'Worker数量',
    ps_count INT DEFAULT 0 COMMENT 'Parameter Server数量',
    master_count INT DEFAULT 1 COMMENT 'Master数量',
    env_vars JSON COMMENT '环境变量',
    command_args JSON COMMENT '命令行参数',
    secrets JSON COMMENT '密钥配置',
    config_maps JSON COMMENT '配置映射',
    volume_mounts JSON COMMENT '挂载卷配置',
    queue_name VARCHAR(64) DEFAULT 'default' COMMENT '队列名称',
    priority INT DEFAULT 0 COMMENT '优先级(0-100)',
    node_selector JSON COMMENT '节点选择器',
    tolerations JSON COMMENT '容忍度设置',
    affinity JSON COMMENT '亲和性配置',
    max_runtime_seconds INT DEFAULT 86400 COMMENT '最大运行时间(秒)',
    max_idle_seconds INT DEFAULT 3600 COMMENT '最大空闲时间(秒)',
    auto_restart TINYINT(1) DEFAULT 0 COMMENT '是否自动重启',
    max_retry_count INT DEFAULT 3 COMMENT '最大重试次数',
    volcano_job_name VARCHAR(128) COMMENT 'Volcano作业名',
    volcano_queue VARCHAR(64) COMMENT 'Volcano队列',
    min_available INT DEFAULT 1 COMMENT '最小可用实例数',
    status ENUM(
        'pending',
        'queued',
        'scheduling',
        'running',
        'succeeded',
        'failed',
        'cancelled',
        'suspended',
        'timeout',
        'oom_killed'
    ) DEFAULT 'pending' COMMENT '任务状态',
    phase ENUM(
        'creating',
        'scheduling',
        'initializing',
        'training',
        'saving',
        'completed'
    ) DEFAULT 'creating' COMMENT '执行阶段',
    namespace VARCHAR(128) COMMENT 'K8s命名空间',
    cluster_name VARCHAR(128) COMMENT '集群名称',
    error_message TEXT COMMENT '错误信息',
    error_code VARCHAR(64) COMMENT '错误代码',
    exit_code INT COMMENT '退出码',
    failure_reason VARCHAR(256) COMMENT '失败原因',
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '提交时间',
    queued_at TIMESTAMP NULL COMMENT '进入队列时间',
    scheduled_at TIMESTAMP NULL COMMENT '调度时间',
    start_time TIMESTAMP NULL COMMENT '开始时间',
    end_time TIMESTAMP NULL COMMENT '结束时间',
    duration_seconds INT COMMENT '执行时长(秒)',
    actual_cpu_usage DECIMAL(8, 4) COMMENT '实际CPU使用',
    actual_memory_usage_gb DECIMAL(8, 2) COMMENT '实际内存使用(GB)',
    actual_gpu_usage DECIMAL(5, 2) COMMENT '实际GPU使用率',
    peak_memory_usage_gb DECIMAL(8, 2) COMMENT '峰值内存使用(GB)',
    total_gpu_hours DECIMAL(10, 2) COMMENT '总GPU小时数',
    workspace_path VARCHAR(512) COMMENT '工作空间路径',
    logs_path VARCHAR(512) COMMENT '日志路径',
    output_path VARCHAR(512) COMMENT '输出路径',
    checkpoint_path VARCHAR(512) COMMENT '检查点路径',
    tensorboard_path VARCHAR(512) COMMENT 'TensorBoard路径',
    hyperparameters JSON COMMENT '超参数',
    training_config JSON COMMENT '训练配置',
    optimizer_config JSON COMMENT '优化器配置',
    scheduler_config JSON COMMENT '调度器配置',
    enable_tensorboard TINYINT(1) DEFAULT 1 COMMENT '是否启用TensorBoard',
    enable_profiling TINYINT(1) DEFAULT 0 COMMENT '是否启用性能分析',
    metrics_collection_interval INT DEFAULT 60 COMMENT '指标收集间隔(秒)',
    notification_config JSON COMMENT '通知配置',
    tags JSON COMMENT '标签',
    annotations JSON COMMENT '注解',
    metadata JSON COMMENT '元数据',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_job_type (job_type),
    INDEX idx_framework (framework),
    INDEX idx_status (status),
    INDEX idx_phase (phase),
    INDEX idx_queue_name (queue_name),
    INDEX idx_priority (priority),
    INDEX idx_volcano_job (volcano_job_name),
    INDEX idx_cluster_name (cluster_name),
    INDEX idx_submitted_at (submitted_at),
    INDEX idx_start_time (start_time),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_status_priority (status, priority)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练作业表';
-- 训练任务实例表
CREATE TABLE vt_training_job_instances (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_id BIGINT NOT NULL COMMENT '所属训练作业ID',
    instance_name VARCHAR(128) NOT NULL COMMENT '实例名称',
    instance_type ENUM('master', 'worker', 'ps', 'evaluator', 'chief') DEFAULT 'worker' COMMENT '实例类型',
    instance_index INT NOT NULL COMMENT '实例索引',
    replica_index INT COMMENT '副本索引',
    pod_name VARCHAR(128) COMMENT 'Pod名称',
    namespace VARCHAR(128) COMMENT '命名空间',
    node_name VARCHAR(128) COMMENT '所在节点',
    node_ip VARCHAR(64) COMMENT '节点IP',
    pod_ip VARCHAR(64) COMMENT 'Pod IP',
    container_id VARCHAR(128) COMMENT '容器ID',
    allocated_cpu_cores DECIMAL(6, 3) COMMENT '分配的CPU核心数',
    allocated_memory_gb DECIMAL(8, 2) COMMENT '分配的内存(GB)',
    allocated_gpu_devices JSON COMMENT '分配的GPU设备信息',
    allocated_storage_gb DECIMAL(8, 2) COMMENT '分配的存储(GB)',
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
    ready TINYINT(1) DEFAULT 0 COMMENT '是否就绪',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    scheduled_at TIMESTAMP NULL COMMENT '调度时间',
    start_time TIMESTAMP NULL COMMENT '开始时间',
    end_time TIMESTAMP NULL COMMENT '结束时间',
    last_transition_time TIMESTAMP NULL COMMENT '最后状态变更时间',
    restart_count INT DEFAULT 0 COMMENT '重启次数',
    last_restart_time TIMESTAMP NULL COMMENT '最后重启时间',
    exit_code INT COMMENT '退出码',
    termination_reason VARCHAR(256) COMMENT '终止原因',
    cpu_usage_percent DECIMAL(5, 2) DEFAULT 0 COMMENT 'CPU使用率',
    memory_usage_percent DECIMAL(5, 2) DEFAULT 0 COMMENT '内存使用率',
    gpu_usage_percent DECIMAL(5, 2) DEFAULT 0 COMMENT 'GPU使用率',
    logs_path VARCHAR(512) COMMENT '日志路径',
    labels JSON COMMENT '标签',
    annotations JSON COMMENT '注解',
    INDEX idx_job_id (job_id),
    INDEX idx_instance_name (instance_name),
    INDEX idx_instance_type (instance_type),
    INDEX idx_pod_name (pod_name),
    INDEX idx_node_name (node_name),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_start_time (start_time),
    UNIQUE KEY uk_job_instance (job_id, instance_name)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练任务实例表';
-- 训练指标表
CREATE TABLE vt_training_metrics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_id BIGINT NOT NULL COMMENT '训练作业ID',
    instance_id BIGINT COMMENT '实例ID',
    metric_name VARCHAR(128) NOT NULL COMMENT '指标名称',
    metric_type ENUM(
        'scalar',
        'histogram',
        'image',
        'text',
        'audio',
        'video'
    ) DEFAULT 'scalar' COMMENT '指标类型',
    metric_value DECIMAL(20, 8) COMMENT '指标值',
    metric_data JSON COMMENT '复杂指标数据',
    step BIGINT COMMENT '训练步数',
    epoch INT COMMENT '训练轮次',
    global_step BIGINT COMMENT '全局步数',
    batch_idx INT COMMENT '批次索引',
    tag VARCHAR(128) COMMENT '标签',
    category VARCHAR(64) COMMENT '类别',
    phase ENUM('train', 'val', 'test') COMMENT '阶段',
    metric_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '指标时间',
    wall_time DECIMAL(15, 3) COMMENT '墙钟时间',
    relative_time_seconds DECIMAL(15, 3) COMMENT '相对时间(秒)',
    min_value DECIMAL(20, 8) COMMENT '最小值',
    max_value DECIMAL(20, 8) COMMENT '最大值',
    avg_value DECIMAL(20, 8) COMMENT '平均值',
    std_value DECIMAL(20, 8) COMMENT '标准差',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_job_id (job_id),
    INDEX idx_instance_id (instance_id),
    INDEX idx_job_metric (job_id, metric_name),
    INDEX idx_job_step (job_id, step),
    INDEX idx_job_epoch (job_id, epoch),
    INDEX idx_metric_time (metric_time),
    INDEX idx_tag (tag),
    INDEX idx_category (category),
    INDEX idx_phase (phase),
    INDEX idx_created_at (created_at)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练指标表';
-- 训练日志表
CREATE TABLE vt_training_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_id BIGINT NOT NULL COMMENT '训练作业ID',
    instance_id BIGINT COMMENT '实例ID',
    log_level ENUM(
        'TRACE',
        'DEBUG',
        'INFO',
        'WARN',
        'ERROR',
        'FATAL'
    ) DEFAULT 'INFO' COMMENT '日志级别',
    log_source VARCHAR(64) COMMENT '日志来源',
    log_content TEXT NOT NULL COMMENT '日志内容',
    log_format ENUM('text', 'json', 'structured') DEFAULT 'text' COMMENT '日志格式',
    log_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '日志时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    file_name VARCHAR(256) COMMENT '文件名',
    line_number INT COMMENT '行号',
    function_name VARCHAR(128) COMMENT '函数名',
    thread_id VARCHAR(64) COMMENT '线程ID',
    process_id VARCHAR(64) COMMENT '进程ID',
    context JSON COMMENT '上下文信息',
    correlation_id VARCHAR(128) COMMENT '关联ID',
    category VARCHAR(64) COMMENT '类别',
    tags JSON COMMENT '标签',
    INDEX idx_job_id (job_id),
    INDEX idx_instance_id (instance_id),
    INDEX idx_log_level (log_level),
    INDEX idx_log_source (log_source),
    INDEX idx_log_time (log_time),
    INDEX idx_category (category),
    INDEX idx_correlation_id (correlation_id),
    INDEX idx_created_at (created_at),
    INDEX idx_job_log_time (job_id, log_time)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练日志表';
-- 训练检查点表
CREATE TABLE vt_training_checkpoints (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_id BIGINT NOT NULL COMMENT '训练作业ID',
    checkpoint_name VARCHAR(128) NOT NULL COMMENT '检查点名称',
    checkpoint_type ENUM('auto', 'manual', 'best', 'final', 'scheduled') DEFAULT 'auto' COMMENT '检查点类型',
    checkpoint_format ENUM(
        'pytorch',
        'tensorflow',
        'onnx',
        'pickle',
        'hdf5'
    ) COMMENT '检查点格式',
    step BIGINT COMMENT '训练步数',
    epoch INT COMMENT '训练轮次',
    global_step BIGINT COMMENT '全局步数',
    storage_path VARCHAR(512) NOT NULL COMMENT '存储路径',
    file_size BIGINT COMMENT '文件大小(字节)',
    checksum VARCHAR(128) COMMENT '文件校验和',
    compression_type ENUM('none', 'gzip', 'bzip2', 'lz4') DEFAULT 'none' COMMENT '压缩类型',
    metrics JSON COMMENT '性能指标',
    loss_value DECIMAL(20, 8) COMMENT '损失值',
    accuracy DECIMAL(8, 6) COMMENT '准确率',
    validation_score DECIMAL(8, 6) COMMENT '验证分数',
    model_config JSON COMMENT '模型配置',
    optimizer_state JSON COMMENT '优化器状态',
    scheduler_state JSON COMMENT '调度器状态',
    status ENUM(
        'saving',
        'saved',
        'failed',
        'deleted',
        'corrupted'
    ) DEFAULT 'saving' COMMENT '状态',
    is_best TINYINT(1) DEFAULT 0 COMMENT '是否最佳检查点',
    is_latest TINYINT(1) DEFAULT 0 COMMENT '是否最新检查点',
    tags JSON COMMENT '标签',
    metadata JSON COMMENT '元数据',
    description TEXT COMMENT '描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    saved_at TIMESTAMP NULL COMMENT '保存完成时间',
    INDEX idx_job_id (job_id),
    INDEX idx_checkpoint_type (checkpoint_type),
    INDEX idx_step (step),
    INDEX idx_epoch (epoch),
    INDEX idx_status (status),
    INDEX idx_is_best (is_best),
    INDEX idx_is_latest (is_latest),
    INDEX idx_created_at (created_at),
    INDEX idx_saved_at (saved_at),
    UNIQUE KEY uk_job_checkpoint (job_id, checkpoint_name)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练检查点表';
-- 训练作业关联关系表 (通用关联表，替代各种外键关系)
CREATE TABLE vt_training_job_relations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    job_id BIGINT NOT NULL COMMENT '训练作业ID',
    entity_type VARCHAR(64) NOT NULL COMMENT '实体类型(workspace, project, user, dataset, model, code_file等)',
    entity_id BIGINT NOT NULL COMMENT '实体ID',
    relation_type VARCHAR(64) NOT NULL COMMENT '关联类型(workspace, project, creator, dataset, base_model, code_file等)',
    is_primary TINYINT(1) DEFAULT 0 COMMENT '是否主要关联',
    sort_order INT DEFAULT 0 COMMENT '排序',
    status ENUM('active', 'inactive', 'pending', 'deleted') DEFAULT 'active' COMMENT '状态',
    metadata JSON COMMENT '元数据',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_job_entity_relation (job_id, entity_type, entity_id, relation_type),
    INDEX idx_job_id (job_id),
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_relation_type (relation_type),
    INDEX idx_status (status),
    INDEX idx_entity_relation (entity_type, relation_type)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练作业关联关系表';
-- 训练检查点文件关联表
CREATE TABLE vt_training_checkpoint_files (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    checkpoint_id BIGINT NOT NULL COMMENT '检查点ID',
    file_id BIGINT NOT NULL COMMENT '文件ID',
    file_type VARCHAR(64) DEFAULT 'checkpoint' COMMENT '文件类型',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY uk_checkpoint_file (checkpoint_id, file_id),
    INDEX idx_checkpoint_id (checkpoint_id),
    INDEX idx_file_id (file_id),
    INDEX idx_file_type (file_type),
    INDEX idx_status (status)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '训练检查点文件关联表';