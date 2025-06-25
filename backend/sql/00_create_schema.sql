-- =====================================
-- AI-GPU机器学习平台完整数据库架构初始化脚本
-- =====================================
-- 本脚本是完整的数据库初始化脚本，包含：
-- 1. 数据库和用户创建
-- 2. 所有表结构创建
-- 3. 索引和约束设置
-- 4. 初始数据插入
-- =====================================

-- 设置字符集和排序规则
SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建数据库
CREATE DATABASE IF NOT EXISTS volctraindb 
    CHARACTER SET utf8mb4 
    COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE volctraindb;

-- 创建应用用户（如果不存在）
-- 注意：在生产环境中，请使用更强的密码
CREATE USER IF NOT EXISTS 'volctrain_app'@'%' IDENTIFIED BY 'Abc@1234';
CREATE USER IF NOT EXISTS 'volctrain_readonly'@'%' IDENTIFIED BY 'volctrain_readonly_password_2024';

-- 授予权限
GRANT ALL PRIVILEGES ON volctraindb.* TO 'volctrain_app'@'%';
GRANT SELECT ON volctraindb.* TO 'volctrain_readonly'@'%';
FLUSH PRIVILEGES;

-- 设置全局配置（根据服务器配置调整）
-- SET GLOBAL innodb_buffer_pool_size = 2147483648; -- 2GB
-- SET GLOBAL max_connections = 500;

-- 设置会话配置
SET SESSION sql_mode = 'STRICT_TRANS_TABLES,NO_ZERO_DATE,NO_ZERO_IN_DATE,ERROR_FOR_DIVISION_BY_ZERO';
SET SESSION time_zone = '+08:00';
SET SESSION foreign_key_checks = 0; -- 禁用外键检查，因为我们不使用外键

-- =====================================
-- 开始创建表结构
-- =====================================

-- 1. 用户和权限模块表
SOURCE 01_users_permissions.sql;

-- 2. 文件管理模块表  
SOURCE 02_files.sql;

-- 3. 工作空间模块表
SOURCE 03_workspaces.sql;

-- 4. 数据集模块表
SOURCE 04_datasets.sql;

-- 5. GPU集群模块表
SOURCE 05_gpu_clusters.sql;

-- 6. 训练任务模块表
SOURCE 06_training_jobs.sql;

-- 7. 模型模块表
SOURCE 07_models.sql;

-- 8. 系统支持模块表
SOURCE 08_system_support.sql;

-- 9. 监控模块表
SOURCE 09_monitoring.sql;

-- 10. K8s关联关系模块表
SOURCE 10_k8s_relations.sql;

-- =====================================
-- 创建额外的索引优化
-- =====================================

-- 为关联表添加复合索引，提高查询性能
-- 这些索引针对常见的查询模式进行优化

-- 用户相关关联表索引
CREATE INDEX idx_user_dept_composite ON vt_department_managers(department_id, user_id, status);

-- 文件相关关联表索引  
CREATE INDEX idx_file_entity_composite ON vt_file_relations(entity_type, entity_id, relation_type, status);
CREATE INDEX idx_file_access_composite ON vt_file_access_logs(file_id, action_type, created_at);

-- 工作空间相关关联表索引
CREATE INDEX idx_workspace_owner_composite ON vt_workspace_owners(workspace_id, user_id, status);
CREATE INDEX idx_project_owner_composite ON vt_project_owners(project_id, user_id, status);

-- 数据集相关关联表索引
CREATE INDEX idx_dataset_relation_composite ON vt_dataset_relations(entity_type, entity_id, relation_type, status);
CREATE INDEX idx_dataset_version_composite ON vt_dataset_version_relations(version_id, entity_type, entity_id);

-- GPU集群相关关联表索引
CREATE INDEX idx_gpu_cluster_node_composite ON vt_gpu_cluster_nodes(cluster_id, node_id, status);
CREATE INDEX idx_gpu_device_alloc_composite ON vt_gpu_device_allocations(device_id, entity_type, entity_id, status);

-- 训练任务相关关联表索引
CREATE INDEX idx_training_relation_composite ON vt_training_job_relations(entity_type, entity_id, relation_type, status);

-- 模型相关关联表索引
CREATE INDEX idx_model_relation_composite ON vt_model_relations(entity_type, entity_id, relation_type, status);
CREATE INDEX idx_model_deployment_composite ON vt_model_deployment_relations(entity_type, entity_id, relation_type);

-- K8s相关关联表索引
CREATE INDEX idx_k8s_resource_composite ON vt_k8s_resource_relations(entity_type, entity_id, relation_type, status);

-- =====================================
-- 插入初始系统数据
-- =====================================

-- 插入默认系统配置
INSERT INTO vt_system_configs (config_key, config_value, config_type, category, description) VALUES
('system.version', '1.0.0', 'string', 'system', '系统版本号'),
('system.timezone', 'Asia/Shanghai', 'string', 'system', '系统时区'),
('system.max_file_size_mb', '1024', 'number', 'file', '最大文件上传大小(MB)'),
('system.session_timeout_hours', '24', 'number', 'security', '会话超时时间(小时)'),
('system.enable_gpu_monitoring', 'true', 'boolean', 'gpu', '是否启用GPU监控'),
('system.default_gpu_quota', '2', 'number', 'gpu', '默认GPU配额'),
('system.enable_auto_scaling', 'true', 'boolean', 'cluster', '是否启用自动扩缩容'),
('training.default_timeout_hours', '168', 'number', 'training', '默认训练超时时间(小时)'),
('training.max_concurrent_jobs', '50', 'number', 'training', '最大并发训练任务数'),
('model.default_storage_type', 'local', 'string', 'model', '默认模型存储类型'),
('notification.enable_email', 'true', 'boolean', 'notification', '是否启用邮件通知'),
('monitoring.data_retention_days', '90', 'number', 'monitoring', '监控数据保留天数');

-- 插入默认训练队列
INSERT INTO vt_training_queues (name, display_name, description, queue_type, priority, max_concurrent_jobs) VALUES
('default', '默认队列', '系统默认训练队列', 'default', 0, 10),
('high_priority', '高优先级队列', '高优先级训练任务队列', 'high_priority', 100, 5),
('gpu_intensive', 'GPU密集型队列', 'GPU密集型训练任务队列', 'gpu_intensive', 50, 3),
('experimental', '实验队列', '实验性训练任务队列', 'experimental', -10, 20);

-- 插入默认Volcano队列
INSERT INTO vt_volcano_queues (queue_name, namespace, weight, state) VALUES
('default', 'volcano-system', 1, 'Open'),
('high-priority', 'volcano-system', 5, 'Open'),
('gpu-intensive', 'volcano-system', 3, 'Open');

-- 插入默认监控指标
INSERT INTO vt_monitor_metrics (name, display_name, description, metric_type, category, module, unit) VALUES
('cpu_usage_percent', 'CPU使用率', '系统CPU使用率百分比', 'gauge', 'system', 'system', '%'),
('memory_usage_percent', '内存使用率', '系统内存使用率百分比', 'gauge', 'system', 'system', '%'),
('gpu_usage_percent', 'GPU使用率', 'GPU使用率百分比', 'gauge', 'gpu', 'gpu', '%'),
('gpu_memory_usage_percent', 'GPU显存使用率', 'GPU显存使用率百分比', 'gauge', 'gpu', 'gpu', '%'),
('disk_usage_percent', '磁盘使用率', '磁盘使用率百分比', 'gauge', 'storage', 'system', '%'),
('network_bytes_in', '网络入流量', '网络入站流量', 'counter', 'network', 'system', 'bytes'),
('network_bytes_out', '网络出流量', '网络出站流量', 'counter', 'network', 'system', 'bytes'),
('training_job_count', '训练任务数', '当前运行的训练任务数量', 'gauge', 'training', 'training', 'count'),
('model_inference_qps', '模型推理QPS', '模型推理每秒请求数', 'gauge', 'inference', 'model', 'qps'),
('model_inference_latency_ms', '模型推理延迟', '模型推理平均延迟', 'histogram', 'inference', 'model', 'ms');

-- 插入默认通知模板
INSERT INTO vt_notification_templates (name, display_name, template_type, channel_type, subject_template, body_template) VALUES
('training_job_completed', '训练任务完成通知', 'alert', 'email', '训练任务 {{job_name}} 已完成', '您的训练任务 "{{job_name}}" 已成功完成。\n\n任务详情：\n- 任务ID: {{job_id}}\n- 状态: {{status}}\n- 开始时间: {{start_time}}\n- 结束时间: {{end_time}}\n- 运行时长: {{duration}}\n\n请登录平台查看详细结果。'),
('training_job_failed', '训练任务失败通知', 'alert', 'email', '训练任务 {{job_name}} 执行失败', '您的训练任务 "{{job_name}}" 执行失败。\n\n任务详情：\n- 任务ID: {{job_id}}\n- 状态: {{status}}\n- 错误信息: {{error_message}}\n- 失败时间: {{failed_time}}\n\n请检查任务配置和代码，然后重新提交。'),
('gpu_usage_high', 'GPU使用率告警', 'alert', 'email', 'GPU使用率过高告警', 'GPU使用率超过阈值\n\n告警详情：\n- 节点: {{node_name}}\n- GPU设备: {{gpu_device}}\n- 当前使用率: {{current_usage}}%\n- 告警阈值: {{threshold}}%\n- 告警时间: {{alert_time}}\n\n请及时检查GPU资源使用情况。');

-- =====================================
-- 完成初始化
-- =====================================

-- 重新启用外键检查（虽然我们不使用外键）
SET SESSION foreign_key_checks = 1;

-- 提交事务
COMMIT;

-- 显示初始化完成信息
SELECT 
    'AI-GPU机器学习平台数据库初始化完成！' as '状态',
    NOW() as '完成时间',
    COUNT(*) as '创建表数量'
FROM information_schema.tables 
WHERE table_schema = 'volctraindb';

-- 显示数据库信息
SELECT 
    SCHEMA_NAME as '数据库名',
    DEFAULT_CHARACTER_SET_NAME as '字符集',
    DEFAULT_COLLATION_NAME as '排序规则'
FROM information_schema.SCHEMATA
WHERE SCHEMA_NAME = 'volctraindb';

-- 显示所有创建的表
SELECT 
    table_name as '表名', 
    table_comment as '表注释',
    IFNULL(table_rows, 0) as '初始行数'
FROM information_schema.tables 
WHERE table_schema = 'volctraindb' 
    AND table_type = 'BASE TABLE'
ORDER BY table_name;

-- 显示关联表统计
SELECT 
    '关联表统计' as '类型',
    COUNT(*) as '数量'
FROM information_schema.tables 
WHERE table_schema = 'volctraindb' 
    AND table_name LIKE '%_relations'
UNION ALL
SELECT 
    '业务表统计' as '类型',
    COUNT(*) as '数量'
FROM information_schema.tables 
WHERE table_schema = 'volctraindb' 
    AND table_name NOT LIKE '%_relations'
    AND table_type = 'BASE TABLE';

-- 显示用户权限
SHOW GRANTS FOR 'volctrain_app'@'%';
SHOW GRANTS FOR 'volctrain_readonly'@'%';

-- =====================================
-- 初始化脚本完成
-- =====================================