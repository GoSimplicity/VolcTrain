-- MySQL dump 10.13  Distrib 8.2.0, for macos13 (arm64)
--
-- Host: localhost    Database: volctraindb
-- ------------------------------------------------------
-- Server version	8.2.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `vt_alert_records`
--

DROP TABLE IF EXISTS `vt_alert_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_alert_records` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `rule_id` bigint NOT NULL COMMENT '规则ID',
  `alert_id` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '告警ID',
  `alert_name` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '告警名称',
  `alert_level` enum('info','warning','critical','fatal') COLLATE utf8mb4_unicode_ci DEFAULT 'warning' COMMENT '告警级别',
  `severity_score` int DEFAULT '1' COMMENT '严重程度评分',
  `message` text COLLATE utf8mb4_unicode_ci COMMENT '告警信息',
  `summary` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '告警摘要',
  `resource_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '资源类型',
  `resource_id` bigint DEFAULT NULL COMMENT '资源ID',
  `resource_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '资源名称',
  `instance_id` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '实例ID',
  `trigger_value` decimal(20,6) DEFAULT NULL COMMENT '触发值',
  `threshold_value` decimal(20,6) DEFAULT NULL COMMENT '阈值',
  `condition_expression` text COLLATE utf8mb4_unicode_ci COMMENT '触发条件',
  `evaluation_data` json DEFAULT NULL COMMENT '评估数据',
  `triggered_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '触发时间',
  `first_occurrence_at` timestamp NULL DEFAULT NULL COMMENT '首次发生时间',
  `last_occurrence_at` timestamp NULL DEFAULT NULL COMMENT '最后发生时间',
  `resolved_at` timestamp NULL DEFAULT NULL COMMENT '解决时间',
  `acknowledged_at` timestamp NULL DEFAULT NULL COMMENT '确认时间',
  `status` enum('firing','resolved','acknowledged','suppressed','silenced') COLLATE utf8mb4_unicode_ci DEFAULT 'firing' COMMENT '告警状态',
  `occurrence_count` int DEFAULT '1' COMMENT '发生次数',
  `resolution_notes` text COLLATE utf8mb4_unicode_ci COMMENT '解决备注',
  `root_cause` text COLLATE utf8mb4_unicode_ci COMMENT '根本原因',
  `notification_sent` tinyint(1) DEFAULT '0' COMMENT '是否已发送通知',
  `notification_channels` json DEFAULT NULL COMMENT '通知渠道',
  `notification_count` int DEFAULT '0' COMMENT '通知次数',
  `last_notification_at` timestamp NULL DEFAULT NULL COMMENT '最后通知时间',
  `escalation_level` int DEFAULT '0' COMMENT '升级级别',
  `alert_group_id` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '告警组ID',
  `correlation_id` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '关联ID',
  `labels` json DEFAULT NULL COMMENT '标签',
  `annotations` json DEFAULT NULL COMMENT '注解',
  `context` json DEFAULT NULL COMMENT '上下文信息',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_alert_id` (`alert_id`),
  KEY `idx_rule_id` (`rule_id`),
  KEY `idx_alert_level` (`alert_level`),
  KEY `idx_status` (`status`),
  KEY `idx_resource` (`resource_type`,`resource_id`),
  KEY `idx_triggered_at` (`triggered_at`),
  KEY `idx_resolved_at` (`resolved_at`),
  KEY `idx_alert_group_id` (`alert_group_id`),
  KEY `idx_correlation_id` (`correlation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_alert_records`
--

LOCK TABLES `vt_alert_records` WRITE;
/*!40000 ALTER TABLE `vt_alert_records` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_alert_records` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_alert_rules`
--

DROP TABLE IF EXISTS `vt_alert_rules`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_alert_rules` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '规则名称',
  `display_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '显示名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '规则描述',
  `rule_type` enum('threshold','trend','anomaly','composite','custom') COLLATE utf8mb4_unicode_ci DEFAULT 'threshold' COMMENT '规则类型',
  `condition_expression` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '条件表达式',
  `query_expression` text COLLATE utf8mb4_unicode_ci COMMENT '查询表达式',
  `warning_threshold` decimal(20,6) DEFAULT NULL COMMENT '警告阈值',
  `critical_threshold` decimal(20,6) DEFAULT NULL COMMENT '严重阈值',
  `threshold_condition` enum('gt','gte','lt','lte','eq','neq','between','outside') COLLATE utf8mb4_unicode_ci DEFAULT 'gt' COMMENT '阈值条件',
  `evaluation_window_seconds` int DEFAULT '300' COMMENT '评估窗口(秒)',
  `evaluation_interval_seconds` int DEFAULT '60' COMMENT '评估间隔(秒)',
  `trigger_duration_seconds` int DEFAULT '300' COMMENT '持续时间(秒)',
  `recovery_duration_seconds` int DEFAULT '60' COMMENT '恢复时间(秒)',
  `filter_labels` json DEFAULT NULL COMMENT '过滤标签',
  `filter_resources` json DEFAULT NULL COMMENT '过滤资源',
  `time_range_filter` json DEFAULT NULL COMMENT '时间范围过滤',
  `alert_level` enum('info','warning','critical','fatal') COLLATE utf8mb4_unicode_ci DEFAULT 'warning' COMMENT '告警级别',
  `severity_score` int DEFAULT '1' COMMENT '严重程度评分(1-10)',
  `notification_channels` json DEFAULT NULL COMMENT '通知渠道',
  `notification_throttle_minutes` int DEFAULT '60' COMMENT '通知限流(分钟)',
  `escalation_policy` json DEFAULT NULL COMMENT '升级策略',
  `silence_duration_seconds` int DEFAULT '3600' COMMENT '静默时间(秒)',
  `suppression_rules` json DEFAULT NULL COMMENT '抑制规则',
  `dependency_rules` json DEFAULT NULL COMMENT '依赖规则',
  `workspace_ids` json DEFAULT NULL COMMENT '适用工作空间ID列表',
  `resource_scope` json DEFAULT NULL COMMENT '资源范围',
  `status` enum('active','inactive','paused') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '规则状态',
  `trigger_count` int DEFAULT '0' COMMENT '触发次数',
  `last_evaluation_at` timestamp NULL DEFAULT NULL COMMENT '最后评估时间',
  `last_trigger_at` timestamp NULL DEFAULT NULL COMMENT '最后触发时间',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_rule_type` (`rule_type`),
  KEY `idx_alert_level` (`alert_level`),
  KEY `idx_status` (`status`),
  KEY `idx_last_evaluation_at` (`last_evaluation_at`),
  KEY `idx_last_trigger_at` (`last_trigger_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='告警规则表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_alert_rules`
--

LOCK TABLES `vt_alert_rules` WRITE;
/*!40000 ALTER TABLE `vt_alert_rules` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_alert_rules` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_api_key_relations`
--

DROP TABLE IF EXISTS `vt_api_key_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_api_key_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `api_key_id` bigint NOT NULL COMMENT 'API密钥ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(user等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(owner, revoked_by等)',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_api_key_entity_relation` (`api_key_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_api_key_id` (`api_key_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='API密钥关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_api_key_relations`
--

LOCK TABLES `vt_api_key_relations` WRITE;
/*!40000 ALTER TABLE `vt_api_key_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_api_key_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_api_keys`
--

DROP TABLE IF EXISTS `vt_api_keys`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_api_keys` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `key_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密钥名称',
  `api_key` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'API密钥',
  `api_secret` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'API密钥密文',
  `key_type` enum('personal','service','temporary') COLLATE utf8mb4_unicode_ci DEFAULT 'personal' COMMENT '密钥类型',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '密钥描述',
  `permissions` json DEFAULT NULL COMMENT '权限配置',
  `rate_limit_per_minute` int DEFAULT '60' COMMENT '每分钟限流',
  `rate_limit_per_hour` int DEFAULT '3600' COMMENT '每小时限流',
  `allowed_ips` json DEFAULT NULL COMMENT '允许的IP列表',
  `allowed_domains` json DEFAULT NULL COMMENT '允许的域名列表',
  `expires_at` timestamp NULL DEFAULT NULL COMMENT '过期时间',
  `last_used_at` timestamp NULL DEFAULT NULL COMMENT '最后使用时间',
  `last_used_ip` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '最后使用IP',
  `usage_count` int DEFAULT '0' COMMENT '使用次数',
  `status` enum('active','inactive','revoked','expired') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `revoked_at` timestamp NULL DEFAULT NULL COMMENT '撤销时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_api_key` (`api_key`),
  KEY `idx_key_type` (`key_type`),
  KEY `idx_status` (`status`),
  KEY `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='API密钥表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_api_keys`
--

LOCK TABLES `vt_api_keys` WRITE;
/*!40000 ALTER TABLE `vt_api_keys` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_api_keys` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_audit_log_relations`
--

DROP TABLE IF EXISTS `vt_audit_log_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_audit_log_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `audit_log_id` bigint NOT NULL COMMENT '审计日志ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(user等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(user等)',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_audit_entity_relation` (`audit_log_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_audit_log_id` (`audit_log_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审计日志关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_audit_log_relations`
--

LOCK TABLES `vt_audit_log_relations` WRITE;
/*!40000 ALTER TABLE `vt_audit_log_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_audit_log_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_audit_logs`
--

DROP TABLE IF EXISTS `vt_audit_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_audit_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `resource_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资源类型',
  `resource_id` bigint DEFAULT NULL COMMENT '资源ID',
  `action` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作动作',
  `details` json DEFAULT NULL COMMENT '操作详情',
  `ip_address` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'IP地址',
  `user_agent` text COLLATE utf8mb4_unicode_ci COMMENT '用户代理',
  `status` enum('success','failed') COLLATE utf8mb4_unicode_ci DEFAULT 'success' COMMENT '操作状态',
  `error_message` text COLLATE utf8mb4_unicode_ci COMMENT '错误信息',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_resource` (`resource_type`,`resource_id`),
  KEY `idx_action` (`action`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='统一审计日志表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_audit_logs`
--

LOCK TABLES `vt_audit_logs` WRITE;
/*!40000 ALTER TABLE `vt_audit_logs` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_audit_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_dataset_file_annotations`
--

DROP TABLE IF EXISTS `vt_dataset_file_annotations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_dataset_file_annotations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `dataset_file_id` bigint NOT NULL COMMENT '数据集文件ID',
  `user_id` bigint NOT NULL COMMENT '标注用户ID',
  `annotation_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT 'primary' COMMENT '标注类型',
  `annotation_status` enum('in_progress','completed','verified','rejected') COLLATE utf8mb4_unicode_ci DEFAULT 'in_progress' COMMENT '标注状态',
  `annotation_data` json DEFAULT NULL COMMENT '标注数据',
  `annotation_time_seconds` int DEFAULT NULL COMMENT '标注耗时(秒)',
  `quality_score` decimal(3,2) DEFAULT NULL COMMENT '标注质量评分',
  `review_comments` text COLLATE utf8mb4_unicode_ci COMMENT '审核意见',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_file_user_type` (`dataset_file_id`,`user_id`,`annotation_type`),
  KEY `idx_dataset_file_id` (`dataset_file_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_annotation_status` (`annotation_status`),
  KEY `idx_annotation_type` (`annotation_type`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据集文件标注关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_dataset_file_annotations`
--

LOCK TABLES `vt_dataset_file_annotations` WRITE;
/*!40000 ALTER TABLE `vt_dataset_file_annotations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_dataset_file_annotations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_dataset_files`
--

DROP TABLE IF EXISTS `vt_dataset_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_dataset_files` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `dataset_id` bigint NOT NULL COMMENT '数据集ID',
  `version_id` bigint DEFAULT NULL COMMENT '版本ID',
  `file_id` bigint NOT NULL COMMENT '文件ID',
  `relative_path` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '相对路径',
  `file_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件类型',
  `split_type` enum('train','val','test','all','unlabeled') COLLATE utf8mb4_unicode_ci DEFAULT 'all' COMMENT '数据分割类型',
  `category` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '类别',
  `annotation_status` enum('unlabeled','labeled','verified','rejected') COLLATE utf8mb4_unicode_ci DEFAULT 'unlabeled' COMMENT '标注状态',
  `annotation_data` json DEFAULT NULL COMMENT '标注数据',
  `annotation_at` timestamp NULL DEFAULT NULL COMMENT '标注时间',
  `process_status` enum('pending','processing','completed','failed','skipped') COLLATE utf8mb4_unicode_ci DEFAULT 'pending' COMMENT '处理状态',
  `process_result` json DEFAULT NULL COMMENT '处理结果',
  `error_message` text COLLATE utf8mb4_unicode_ci COMMENT '错误信息',
  `quality_score` decimal(3,2) DEFAULT NULL COMMENT '质量评分',
  `quality_issues` json DEFAULT NULL COMMENT '质量问题',
  `metadata` json DEFAULT NULL COMMENT '文件元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_dataset_file` (`dataset_id`,`file_id`),
  KEY `idx_dataset_id` (`dataset_id`),
  KEY `idx_version_id` (`version_id`),
  KEY `idx_file_id` (`file_id`),
  KEY `idx_split_type` (`split_type`),
  KEY `idx_annotation_status` (`annotation_status`),
  KEY `idx_process_status` (`process_status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据集文件表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_dataset_files`
--

LOCK TABLES `vt_dataset_files` WRITE;
/*!40000 ALTER TABLE `vt_dataset_files` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_dataset_files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_dataset_relations`
--

DROP TABLE IF EXISTS `vt_dataset_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_dataset_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `dataset_id` bigint NOT NULL COMMENT '数据集ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(workspace, user, model, training_job等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(owner, creator, training_dataset, validation_dataset等)',
  `workspace_id` bigint DEFAULT NULL COMMENT '所属工作空间ID',
  `is_primary` tinyint(1) DEFAULT '0' COMMENT '是否主要关联',
  `sort_order` int DEFAULT '0' COMMENT '排序',
  `status` enum('active','inactive','pending','deleted') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_dataset_entity_relation` (`dataset_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_dataset_id` (`dataset_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_workspace_id` (`workspace_id`),
  KEY `idx_status` (`status`),
  KEY `idx_entity_relation` (`entity_type`,`relation_type`),
  KEY `idx_dataset_relation_composite` (`entity_type`,`entity_id`,`relation_type`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据集关联关系表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_dataset_relations`
--

LOCK TABLES `vt_dataset_relations` WRITE;
/*!40000 ALTER TABLE `vt_dataset_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_dataset_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_dataset_version_relations`
--

DROP TABLE IF EXISTS `vt_dataset_version_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_dataset_version_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `version_id` bigint NOT NULL COMMENT '版本ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(user等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(created_by等)',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_version_entity_relation` (`version_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_version_id` (`version_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_status` (`status`),
  KEY `idx_dataset_version_composite` (`version_id`,`entity_type`,`entity_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据集版本关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_dataset_version_relations`
--

LOCK TABLES `vt_dataset_version_relations` WRITE;
/*!40000 ALTER TABLE `vt_dataset_version_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_dataset_version_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_dataset_versions`
--

DROP TABLE IF EXISTS `vt_dataset_versions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_dataset_versions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `dataset_id` bigint NOT NULL COMMENT '数据集ID',
  `version` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '版本号',
  `version_name` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '版本名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '版本描述',
  `change_log` text COLLATE utf8mb4_unicode_ci COMMENT '变更日志',
  `parent_version_id` bigint DEFAULT NULL COMMENT '父版本ID',
  `total_size` bigint DEFAULT '0' COMMENT '总大小(字节)',
  `total_count` int DEFAULT '0' COMMENT '总记录数',
  `train_count` int DEFAULT '0' COMMENT '训练集数量',
  `val_count` int DEFAULT '0' COMMENT '验证集数量',
  `test_count` int DEFAULT '0' COMMENT '测试集数量',
  `storage_path` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '存储路径',
  `storage_config` json DEFAULT NULL COMMENT '存储配置',
  `checksum` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '数据校验和',
  `split_config` json DEFAULT NULL COMMENT '数据分割配置',
  `transform_config` json DEFAULT NULL COMMENT '数据转换配置',
  `preprocessing_config` json DEFAULT NULL COMMENT '预处理配置',
  `status` enum('creating','processing','ready','error','deprecated') COLLATE utf8mb4_unicode_ci DEFAULT 'creating' COMMENT '版本状态',
  `is_default` tinyint(1) DEFAULT '0' COMMENT '是否默认版本',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_dataset_version` (`dataset_id`,`version`),
  KEY `idx_dataset_id` (`dataset_id`),
  KEY `idx_version` (`version`),
  KEY `idx_status` (`status`),
  KEY `idx_is_default` (`is_default`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_parent_version_id` (`parent_version_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据集版本表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_dataset_versions`
--

LOCK TABLES `vt_dataset_versions` WRITE;
/*!40000 ALTER TABLE `vt_dataset_versions` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_dataset_versions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_datasets`
--

DROP TABLE IF EXISTS `vt_datasets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_datasets` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '数据集名称',
  `display_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '显示名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '数据集描述',
  `dataset_type` enum('image','text','audio','video','tabular','time_series','graph','mixed') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '数据集类型',
  `format` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '数据格式(json, csv, parquet等)',
  `version` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT '1.0.0' COMMENT '版本号',
  `total_size` bigint DEFAULT '0' COMMENT '总大小(字节)',
  `total_count` int DEFAULT '0' COMMENT '总记录数',
  `train_count` int DEFAULT '0' COMMENT '训练集数量',
  `val_count` int DEFAULT '0' COMMENT '验证集数量',
  `test_count` int DEFAULT '0' COMMENT '测试集数量',
  `storage_type` enum('local','s3','oss','hdfs','nfs','minio') COLLATE utf8mb4_unicode_ci DEFAULT 'local' COMMENT '存储类型',
  `storage_path` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '存储路径',
  `storage_config` json DEFAULT NULL COMMENT '存储配置',
  `annotation_type` enum('classification','detection','segmentation','regression','nlp','custom','none') COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '标注类型',
  `label_config` json DEFAULT NULL COMMENT '标签配置',
  `classes` json DEFAULT NULL COMMENT '类别信息',
  `quality_score` decimal(3,2) DEFAULT NULL COMMENT '数据质量评分(0-1)',
  `quality_report` json DEFAULT NULL COMMENT '质量报告',
  `data_profile` json DEFAULT NULL COMMENT '数据剖析',
  `status` enum('creating','processing','ready','error','archived') COLLATE utf8mb4_unicode_ci DEFAULT 'creating' COMMENT '数据集状态',
  `visibility` enum('public','private','workspace','shared') COLLATE utf8mb4_unicode_ci DEFAULT 'private' COMMENT '可见性',
  `is_featured` tinyint(1) DEFAULT '0' COMMENT '是否推荐数据集',
  `download_count` int DEFAULT '0' COMMENT '下载次数',
  `view_count` int DEFAULT '0' COMMENT '查看次数',
  `usage_count` int DEFAULT '0' COMMENT '使用次数',
  `star_count` int DEFAULT '0' COMMENT '收藏次数',
  `tags` json DEFAULT NULL COMMENT '标签',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `schema_config` json DEFAULT NULL COMMENT '数据结构配置',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_dataset_type` (`dataset_type`),
  KEY `idx_format` (`format`),
  KEY `idx_status` (`status`),
  KEY `idx_visibility` (`visibility`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_is_featured` (`is_featured`),
  KEY `idx_quality_score` (`quality_score`),
  KEY `idx_type_visibility` (`dataset_type`,`visibility`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据集表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_datasets`
--

LOCK TABLES `vt_datasets` WRITE;
/*!40000 ALTER TABLE `vt_datasets` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_datasets` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_department_managers`
--

DROP TABLE IF EXISTS `vt_department_managers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_department_managers` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `department_id` bigint NOT NULL COMMENT '部门ID',
  `user_id` bigint NOT NULL COMMENT '管理者用户ID',
  `manager_type` enum('primary','deputy','assistant') COLLATE utf8mb4_unicode_ci DEFAULT 'primary' COMMENT '管理者类型',
  `start_date` date DEFAULT NULL COMMENT '开始日期',
  `end_date` date DEFAULT NULL COMMENT '结束日期',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_dept_user_type` (`department_id`,`user_id`,`manager_type`),
  KEY `idx_department_id` (`department_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_manager_type` (`manager_type`),
  KEY `idx_status` (`status`),
  KEY `idx_user_dept_composite` (`department_id`,`user_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门管理者关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_department_managers`
--

LOCK TABLES `vt_department_managers` WRITE;
/*!40000 ALTER TABLE `vt_department_managers` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_department_managers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_departments`
--

DROP TABLE IF EXISTS `vt_departments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_departments` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '部门名称',
  `parent_id` bigint DEFAULT NULL COMMENT '父部门ID',
  `level` int DEFAULT '1' COMMENT '层级',
  `path` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '层级路径',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '部门描述',
  `department_code` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '部门编码',
  `sort_order` int DEFAULT '0' COMMENT '排序',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  UNIQUE KEY `uk_department_code` (`department_code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_level` (`level`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_departments`
--

LOCK TABLES `vt_departments` WRITE;
/*!40000 ALTER TABLE `vt_departments` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_departments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_file_access_logs`
--

DROP TABLE IF EXISTS `vt_file_access_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_file_access_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `file_id` bigint NOT NULL COMMENT '文件ID',
  `user_id` bigint DEFAULT NULL COMMENT '用户ID',
  `workspace_id` bigint DEFAULT NULL COMMENT '工作空间ID',
  `action_type` enum('upload','download','view','delete','copy','move') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作类型',
  `access_ip` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '访问IP',
  `user_agent` text COLLATE utf8mb4_unicode_ci COMMENT '用户代理',
  `file_size` bigint DEFAULT NULL COMMENT '文件大小',
  `transfer_duration_ms` int DEFAULT NULL COMMENT '传输时长(毫秒)',
  `status` enum('success','failed','partial') COLLATE utf8mb4_unicode_ci DEFAULT 'success' COMMENT '操作状态',
  `error_message` text COLLATE utf8mb4_unicode_ci COMMENT '错误信息',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_file_id` (`file_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_workspace_id` (`workspace_id`),
  KEY `idx_action_type` (`action_type`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_file_action` (`file_id`,`action_type`),
  KEY `idx_user_action` (`user_id`,`action_type`),
  KEY `idx_file_access_composite` (`file_id`,`action_type`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件访问记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_file_access_logs`
--

LOCK TABLES `vt_file_access_logs` WRITE;
/*!40000 ALTER TABLE `vt_file_access_logs` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_file_access_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_file_relations`
--

DROP TABLE IF EXISTS `vt_file_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_file_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `file_id` bigint NOT NULL COMMENT '文件ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(user, workspace, dataset, model等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(avatar, attachment, model_file, dataset_file等)',
  `workspace_id` bigint DEFAULT NULL COMMENT '所属工作空间ID',
  `owner_id` bigint DEFAULT NULL COMMENT '拥有者ID',
  `is_primary` tinyint(1) DEFAULT '0' COMMENT '是否主要文件',
  `sort_order` int DEFAULT '0' COMMENT '排序',
  `status` enum('active','inactive','pending','deleted') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_file_entity_relation` (`file_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_file_id` (`file_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_workspace_id` (`workspace_id`),
  KEY `idx_owner_id` (`owner_id`),
  KEY `idx_status` (`status`),
  KEY `idx_entity_relation` (`entity_type`,`relation_type`),
  KEY `idx_file_entity_composite` (`entity_type`,`entity_id`,`relation_type`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件关联关系表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_file_relations`
--

LOCK TABLES `vt_file_relations` WRITE;
/*!40000 ALTER TABLE `vt_file_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_file_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_files`
--

DROP TABLE IF EXISTS `vt_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_files` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `original_name` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '原始文件名',
  `file_name` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '存储文件名',
  `file_path` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件路径',
  `file_size` bigint NOT NULL COMMENT '文件大小(字节)',
  `mime_type` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'MIME类型',
  `file_extension` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件扩展名',
  `file_hash` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件哈希值',
  `checksum` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '校验和',
  `storage_type` enum('local','s3','oss','minio','hdfs') COLLATE utf8mb4_unicode_ci DEFAULT 'local' COMMENT '存储类型',
  `storage_config` json DEFAULT NULL COMMENT '存储配置',
  `bucket_name` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '存储桶名称',
  `file_category` enum('image','document','video','audio','archive','code','model','dataset','other') COLLATE utf8mb4_unicode_ci DEFAULT 'other' COMMENT '文件类别',
  `upload_status` enum('uploading','completed','failed','deleted') COLLATE utf8mb4_unicode_ci DEFAULT 'uploading' COMMENT '上传状态',
  `is_public` tinyint(1) DEFAULT '0' COMMENT '是否公开',
  `download_count` int DEFAULT '0' COMMENT '下载次数',
  `virus_scan_status` enum('pending','clean','infected','error','skipped') COLLATE utf8mb4_unicode_ci DEFAULT 'pending' COMMENT '病毒扫描状态',
  `virus_scan_result` json DEFAULT NULL COMMENT '病毒扫描结果',
  `compression_type` enum('none','gzip','zip','tar') COLLATE utf8mb4_unicode_ci DEFAULT 'none' COMMENT '压缩类型',
  `metadata` json DEFAULT NULL COMMENT '文件元数据',
  `tags` json DEFAULT NULL COMMENT '标签',
  `expire_at` timestamp NULL DEFAULT NULL COMMENT '过期时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_file_name` (`file_name`),
  KEY `idx_file_hash` (`file_hash`),
  KEY `idx_storage_type` (`storage_type`),
  KEY `idx_file_category` (`file_category`),
  KEY `idx_upload_status` (`upload_status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_expire_at` (`expire_at`),
  KEY `idx_category_status` (`file_category`,`upload_status`),
  KEY `idx_hash_size` (`file_hash`,`file_size`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件存储表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_files`
--

LOCK TABLES `vt_files` WRITE;
/*!40000 ALTER TABLE `vt_files` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_gpu_cluster_nodes`
--

DROP TABLE IF EXISTS `vt_gpu_cluster_nodes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_gpu_cluster_nodes` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `cluster_id` bigint NOT NULL COMMENT '集群ID',
  `node_id` bigint NOT NULL COMMENT '节点ID',
  `join_date` date DEFAULT NULL COMMENT '加入日期',
  `leave_date` date DEFAULT NULL COMMENT '离开日期',
  `node_role` enum('master','worker','edge') COLLATE utf8mb4_unicode_ci DEFAULT 'worker' COMMENT '节点角色',
  `status` enum('active','inactive','draining') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_cluster_node` (`cluster_id`,`node_id`),
  KEY `idx_cluster_id` (`cluster_id`),
  KEY `idx_node_id` (`node_id`),
  KEY `idx_node_role` (`node_role`),
  KEY `idx_status` (`status`),
  KEY `idx_gpu_cluster_node_composite` (`cluster_id`,`node_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='GPU集群节点关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_gpu_cluster_nodes`
--

LOCK TABLES `vt_gpu_cluster_nodes` WRITE;
/*!40000 ALTER TABLE `vt_gpu_cluster_nodes` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_gpu_cluster_nodes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_gpu_clusters`
--

DROP TABLE IF EXISTS `vt_gpu_clusters`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_gpu_clusters` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '集群名称',
  `display_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '显示名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '集群描述',
  `cluster_type` enum('k8s','slurm','yarn','standalone') COLLATE utf8mb4_unicode_ci DEFAULT 'k8s' COMMENT '集群类型',
  `cluster_endpoint` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '集群端点',
  `cluster_config` json DEFAULT NULL COMMENT '集群配置',
  `total_nodes` int DEFAULT '0' COMMENT '总节点数',
  `available_nodes` int DEFAULT '0' COMMENT '可用节点数',
  `total_gpus` int DEFAULT '0' COMMENT '总GPU数',
  `available_gpus` int DEFAULT '0' COMMENT '可用GPU数',
  `region` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '区域',
  `zone` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '可用区',
  `datacenter` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '数据中心',
  `status` enum('active','inactive','maintenance','error') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '集群状态',
  `health_status` enum('healthy','warning','critical','unknown') COLLATE utf8mb4_unicode_ci DEFAULT 'unknown' COMMENT '健康状态',
  `last_heartbeat_at` timestamp NULL DEFAULT NULL COMMENT '最后心跳时间',
  `cpu_usage_percent` decimal(5,2) DEFAULT '0.00' COMMENT '集群CPU使用率',
  `memory_usage_percent` decimal(5,2) DEFAULT '0.00' COMMENT '集群内存使用率',
  `gpu_usage_percent` decimal(5,2) DEFAULT '0.00' COMMENT '集群GPU使用率',
  `tags` json DEFAULT NULL COMMENT '标签',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_cluster_type` (`cluster_type`),
  KEY `idx_region_zone` (`region`,`zone`),
  KEY `idx_status` (`status`),
  KEY `idx_health_status` (`health_status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='GPU集群表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_gpu_clusters`
--

LOCK TABLES `vt_gpu_clusters` WRITE;
/*!40000 ALTER TABLE `vt_gpu_clusters` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_gpu_clusters` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_gpu_device_allocations`
--

DROP TABLE IF EXISTS `vt_gpu_device_allocations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_gpu_device_allocations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `device_id` bigint NOT NULL COMMENT '设备ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(user, workspace, training_job等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `allocation_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT 'exclusive' COMMENT '分配类型(exclusive, shared)',
  `allocated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '分配时间',
  `released_at` timestamp NULL DEFAULT NULL COMMENT '释放时间',
  `expected_duration_seconds` int DEFAULT NULL COMMENT '预期使用时长(秒)',
  `priority` int DEFAULT '0' COMMENT '优先级',
  `status` enum('pending','active','releasing','released','expired') COLLATE utf8mb4_unicode_ci DEFAULT 'pending' COMMENT '分配状态',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_device_id` (`device_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_allocation_type` (`allocation_type`),
  KEY `idx_allocated_at` (`allocated_at`),
  KEY `idx_status` (`status`),
  KEY `idx_device_status` (`device_id`,`status`),
  KEY `idx_gpu_device_alloc_composite` (`device_id`,`entity_type`,`entity_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='GPU设备分配关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_gpu_device_allocations`
--

LOCK TABLES `vt_gpu_device_allocations` WRITE;
/*!40000 ALTER TABLE `vt_gpu_device_allocations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_gpu_device_allocations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_gpu_devices`
--

DROP TABLE IF EXISTS `vt_gpu_devices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_gpu_devices` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `device_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '设备名称',
  `device_index` int NOT NULL COMMENT '设备索引',
  `device_uuid` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '设备UUID',
  `pci_bus_id` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'PCI总线ID',
  `minor_number` int DEFAULT NULL COMMENT '设备号',
  `k8s_resource_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'K8s资源名称',
  `k8s_node_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'K8s节点名称',
  `gpu_model` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'GPU型号',
  `gpu_brand` enum('nvidia','amd','intel','other') COLLATE utf8mb4_unicode_ci DEFAULT 'nvidia' COMMENT 'GPU品牌',
  `architecture` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'GPU架构',
  `cuda_version` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'CUDA版本',
  `driver_version` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '驱动版本',
  `vbios_version` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'VBIOS版本',
  `mig_enabled` tinyint(1) DEFAULT '0' COMMENT '是否启用MIG',
  `mig_devices` json DEFAULT NULL COMMENT 'MIG设备列表',
  `memory_total_mb` int DEFAULT NULL COMMENT '显存总量(MB)',
  `memory_free_mb` int DEFAULT NULL COMMENT '可用显存(MB)',
  `memory_used_mb` int DEFAULT NULL COMMENT '已用显存(MB)',
  `core_count` int DEFAULT NULL COMMENT '核心数',
  `sm_count` int DEFAULT NULL COMMENT 'SM数量',
  `base_clock_mhz` int DEFAULT NULL COMMENT '基础频率(MHz)',
  `boost_clock_mhz` int DEFAULT NULL COMMENT '加速频率(MHz)',
  `memory_clock_mhz` int DEFAULT NULL COMMENT '显存频率(MHz)',
  `memory_bus_width` int DEFAULT NULL COMMENT '显存位宽',
  `compute_capability` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '计算能力',
  `max_power_limit_watts` int DEFAULT NULL COMMENT '最大功耗限制(瓦)',
  `status` enum('available','occupied','reserved','maintenance','error','offline') COLLATE utf8mb4_unicode_ci DEFAULT 'available' COMMENT '设备状态',
  `health_status` enum('healthy','warning','critical','unknown') COLLATE utf8mb4_unicode_ci DEFAULT 'unknown' COMMENT '健康状态',
  `utilization_percent` decimal(5,2) DEFAULT '0.00' COMMENT 'GPU使用率',
  `memory_utilization_percent` decimal(5,2) DEFAULT '0.00' COMMENT '显存使用率',
  `encoder_utilization_percent` decimal(5,2) DEFAULT '0.00' COMMENT '编码器使用率',
  `decoder_utilization_percent` decimal(5,2) DEFAULT '0.00' COMMENT '解码器使用率',
  `temperature_celsius` decimal(5,2) DEFAULT NULL COMMENT '温度(摄氏度)',
  `power_usage_watts` decimal(7,2) DEFAULT NULL COMMENT '功耗(瓦)',
  `power_limit_watts` decimal(7,2) DEFAULT NULL COMMENT '功耗限制(瓦)',
  `fan_speed_percent` decimal(5,2) DEFAULT NULL COMMENT '风扇转速百分比',
  `fan_speed_rpm` int DEFAULT NULL COMMENT '风扇转速(RPM)',
  `performance_state` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '性能状态(P0-P12)',
  `gpu_mode` enum('normal','compute','graphics') COLLATE utf8mb4_unicode_ci DEFAULT 'normal' COMMENT 'GPU模式',
  `multi_gpu_board` tinyint(1) DEFAULT '0' COMMENT '是否多GPU板卡',
  `board_id` int DEFAULT NULL COMMENT '板卡ID',
  `capabilities` json DEFAULT NULL COMMENT '设备能力',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_device_uuid` (`device_uuid`),
  KEY `idx_device_name` (`device_name`),
  KEY `idx_gpu_model` (`gpu_model`),
  KEY `idx_gpu_brand` (`gpu_brand`),
  KEY `idx_status` (`status`),
  KEY `idx_health_status` (`health_status`),
  KEY `idx_k8s_node_name` (`k8s_node_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='GPU设备表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_gpu_devices`
--

LOCK TABLES `vt_gpu_devices` WRITE;
/*!40000 ALTER TABLE `vt_gpu_devices` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_gpu_devices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_gpu_node_devices`
--

DROP TABLE IF EXISTS `vt_gpu_node_devices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_gpu_node_devices` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `node_id` bigint NOT NULL COMMENT '节点ID',
  `device_id` bigint NOT NULL COMMENT '设备ID',
  `device_index` int NOT NULL COMMENT '设备在节点上的索引',
  `mount_path` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '挂载路径',
  `status` enum('active','inactive','error') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_node_device` (`node_id`,`device_id`),
  UNIQUE KEY `uk_node_device_index` (`node_id`,`device_index`),
  KEY `idx_node_id` (`node_id`),
  KEY `idx_device_id` (`device_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='GPU节点设备关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_gpu_node_devices`
--

LOCK TABLES `vt_gpu_node_devices` WRITE;
/*!40000 ALTER TABLE `vt_gpu_node_devices` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_gpu_node_devices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_gpu_nodes`
--

DROP TABLE IF EXISTS `vt_gpu_nodes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_gpu_nodes` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '节点名称',
  `hostname` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主机名',
  `ip_address` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'IP地址',
  `node_type` enum('physical','virtual','cloud') COLLATE utf8mb4_unicode_ci DEFAULT 'physical' COMMENT '节点类型',
  `os_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '操作系统类型',
  `os_version` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '操作系统版本',
  `architecture` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '架构(x86_64, arm64等)',
  `kernel_version` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '内核版本',
  `cpu_model` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'CPU型号',
  `cpu_cores` int DEFAULT NULL COMMENT 'CPU核心数',
  `cpu_threads` int DEFAULT NULL COMMENT 'CPU线程数',
  `memory_gb` int DEFAULT NULL COMMENT '内存大小(GB)',
  `storage_gb` int DEFAULT NULL COMMENT '存储大小(GB)',
  `network_bandwidth_mbps` int DEFAULT NULL COMMENT '网络带宽(Mbps)',
  `gpu_count` int DEFAULT '0' COMMENT 'GPU数量',
  `gpu_total_memory_gb` decimal(8,2) DEFAULT '0.00' COMMENT 'GPU总显存(GB)',
  `status` enum('online','offline','maintenance','error','draining') COLLATE utf8mb4_unicode_ci DEFAULT 'offline' COMMENT '节点状态',
  `health_status` enum('healthy','warning','critical','unknown') COLLATE utf8mb4_unicode_ci DEFAULT 'unknown' COMMENT '健康状态',
  `schedulable` tinyint(1) DEFAULT '1' COMMENT '是否可调度',
  `last_heartbeat_at` timestamp NULL DEFAULT NULL COMMENT '最后心跳时间',
  `boot_time` timestamp NULL DEFAULT NULL COMMENT '启动时间',
  `cpu_usage_percent` decimal(5,2) DEFAULT '0.00' COMMENT 'CPU使用率',
  `memory_usage_percent` decimal(5,2) DEFAULT '0.00' COMMENT '内存使用率',
  `storage_usage_percent` decimal(5,2) DEFAULT '0.00' COMMENT '存储使用率',
  `network_usage_mbps` decimal(10,2) DEFAULT '0.00' COMMENT '网络使用率(Mbps)',
  `temperature_celsius` decimal(5,2) DEFAULT NULL COMMENT '温度(摄氏度)',
  `load_average_1m` decimal(5,2) DEFAULT NULL COMMENT '1分钟负载',
  `load_average_5m` decimal(5,2) DEFAULT NULL COMMENT '5分钟负载',
  `load_average_15m` decimal(5,2) DEFAULT NULL COMMENT '15分钟负载',
  `allocatable_cpu_cores` decimal(6,3) DEFAULT NULL COMMENT '可分配CPU核心数',
  `allocatable_memory_gb` decimal(8,2) DEFAULT NULL COMMENT '可分配内存(GB)',
  `allocatable_gpu_count` int DEFAULT '0' COMMENT '可分配GPU数量',
  `node_group` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '节点组',
  `labels` json DEFAULT NULL COMMENT 'K8s标签',
  `annotations` json DEFAULT NULL COMMENT 'K8s注解',
  `taints` json DEFAULT NULL COMMENT 'K8s污点',
  `tags` json DEFAULT NULL COMMENT '标签',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_hostname` (`hostname`),
  UNIQUE KEY `uk_ip_address` (`ip_address`),
  KEY `idx_name` (`name`),
  KEY `idx_node_type` (`node_type`),
  KEY `idx_status` (`status`),
  KEY `idx_health_status` (`health_status`),
  KEY `idx_schedulable` (`schedulable`),
  KEY `idx_node_group` (`node_group`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='GPU节点表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_gpu_nodes`
--

LOCK TABLES `vt_gpu_nodes` WRITE;
/*!40000 ALTER TABLE `vt_gpu_nodes` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_gpu_nodes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_gpu_usage_records`
--

DROP TABLE IF EXISTS `vt_gpu_usage_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_gpu_usage_records` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `task_type` enum('training','inference','development','testing') COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '任务类型',
  `started_at` timestamp NOT NULL COMMENT '开始时间',
  `ended_at` timestamp NULL DEFAULT NULL COMMENT '结束时间',
  `duration_seconds` int DEFAULT NULL COMMENT '使用时长(秒)',
  `avg_utilization_percent` decimal(5,2) DEFAULT NULL COMMENT '平均GPU使用率',
  `max_utilization_percent` decimal(5,2) DEFAULT NULL COMMENT '最大GPU使用率',
  `avg_memory_utilization_percent` decimal(5,2) DEFAULT NULL COMMENT '平均显存使用率',
  `max_memory_utilization_percent` decimal(5,2) DEFAULT NULL COMMENT '最大显存使用率',
  `avg_power_usage_watts` decimal(7,2) DEFAULT NULL COMMENT '平均功耗(瓦)',
  `max_power_usage_watts` decimal(7,2) DEFAULT NULL COMMENT '最大功耗(瓦)',
  `peak_temperature_celsius` decimal(5,2) DEFAULT NULL COMMENT '峰值温度(摄氏度)',
  `compute_units` decimal(12,6) DEFAULT NULL COMMENT '计算单元',
  `cost_amount` decimal(10,4) DEFAULT NULL COMMENT '费用金额',
  `billing_mode` enum('hourly','per_task','monthly') COLLATE utf8mb4_unicode_ci DEFAULT 'hourly' COMMENT '计费模式',
  `status` enum('running','completed','failed','cancelled') COLLATE utf8mb4_unicode_ci DEFAULT 'running' COMMENT '记录状态',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_started_at` (`started_at`),
  KEY `idx_ended_at` (`ended_at`),
  KEY `idx_status` (`status`),
  KEY `idx_task_type` (`task_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='GPU使用记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_gpu_usage_records`
--

LOCK TABLES `vt_gpu_usage_records` WRITE;
/*!40000 ALTER TABLE `vt_gpu_usage_records` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_gpu_usage_records` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_gpu_usage_relations`
--

DROP TABLE IF EXISTS `vt_gpu_usage_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_gpu_usage_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `usage_record_id` bigint NOT NULL COMMENT 'GPU使用记录ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(gpu_device, user, workspace, task等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(device, user, workspace, task等)',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_usage_entity_relation` (`usage_record_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_usage_record_id` (`usage_record_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_entity_relation` (`entity_type`,`relation_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='GPU使用记录关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_gpu_usage_relations`
--

LOCK TABLES `vt_gpu_usage_relations` WRITE;
/*!40000 ALTER TABLE `vt_gpu_usage_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_gpu_usage_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_k8s_resource_relations`
--

DROP TABLE IF EXISTS `vt_k8s_resource_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_k8s_resource_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `k8s_resource_id` bigint NOT NULL COMMENT 'K8s资源ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(workspace, training_job等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(workspace, training_job等)',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_k8s_entity_relation` (`k8s_resource_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_k8s_resource_id` (`k8s_resource_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_status` (`status`),
  KEY `idx_k8s_resource_composite` (`entity_type`,`entity_id`,`relation_type`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='K8s资源关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_k8s_resource_relations`
--

LOCK TABLES `vt_k8s_resource_relations` WRITE;
/*!40000 ALTER TABLE `vt_k8s_resource_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_k8s_resource_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_k8s_resources`
--

DROP TABLE IF EXISTS `vt_k8s_resources`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_k8s_resources` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `resource_type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资源类型：Job, Pod, Service, ConfigMap, Secret等',
  `resource_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'K8s资源名称',
  `namespace` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'vt-platform' COMMENT 'K8s命名空间',
  `uid` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'K8s资源UID',
  `api_version` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'API版本',
  `kind` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资源种类',
  `cluster_name` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '所属集群名称',
  `related_type` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '关联资源类型：training_job, model_deployment等',
  `related_id` bigint DEFAULT NULL COMMENT '关联资源ID',
  `labels` json DEFAULT NULL COMMENT 'K8s标签',
  `annotations` json DEFAULT NULL COMMENT 'K8s注解',
  `spec` json DEFAULT NULL COMMENT '资源规格',
  `status` json DEFAULT NULL COMMENT '资源状态',
  `k8s_created_at` timestamp NULL DEFAULT NULL COMMENT 'K8s创建时间',
  `k8s_deleted_at` timestamp NULL DEFAULT NULL COMMENT 'K8s删除时间',
  `last_sync_at` timestamp NULL DEFAULT NULL COMMENT '最后同步时间',
  `sync_status` enum('synced','pending','failed') COLLATE utf8mb4_unicode_ci DEFAULT 'pending' COMMENT '同步状态',
  `sync_error` text COLLATE utf8mb4_unicode_ci COMMENT '同步错误信息',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_namespace_name` (`namespace`,`resource_name`,`resource_type`),
  UNIQUE KEY `uk_uid` (`uid`),
  KEY `idx_related` (`related_type`,`related_id`),
  KEY `idx_sync_status` (`sync_status`,`last_sync_at`),
  KEY `idx_k8s_created` (`k8s_created_at`),
  KEY `idx_cluster_name` (`cluster_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='K8s资源映射表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_k8s_resources`
--

LOCK TABLES `vt_k8s_resources` WRITE;
/*!40000 ALTER TABLE `vt_k8s_resources` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_k8s_resources` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_model_deployment_relations`
--

DROP TABLE IF EXISTS `vt_model_deployment_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_model_deployment_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `deployment_id` bigint NOT NULL COMMENT '部署ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(workspace, user等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(workspace, deployed_by等)',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_deployment_entity_relation` (`deployment_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_deployment_id` (`deployment_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_status` (`status`),
  KEY `idx_model_deployment_composite` (`entity_type`,`entity_id`,`relation_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模型部署关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_model_deployment_relations`
--

LOCK TABLES `vt_model_deployment_relations` WRITE;
/*!40000 ALTER TABLE `vt_model_deployment_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_model_deployment_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_model_deployments`
--

DROP TABLE IF EXISTS `vt_model_deployments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_model_deployments` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `model_id` bigint NOT NULL COMMENT '模型ID',
  `version_id` bigint NOT NULL COMMENT '版本ID',
  `deployment_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '部署名称',
  `deployment_type` enum('api','batch','edge','streaming','a_b_test') COLLATE utf8mb4_unicode_ci DEFAULT 'api' COMMENT '部署类型',
  `environment` enum('development','testing','staging','production') COLLATE utf8mb4_unicode_ci DEFAULT 'development' COMMENT '环境',
  `service_name` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '服务名称',
  `service_port` int DEFAULT '8080' COMMENT '服务端口',
  `protocol` enum('http','https','grpc','tcp') COLLATE utf8mb4_unicode_ci DEFAULT 'http' COMMENT '协议',
  `base_path` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT '/' COMMENT '基础路径',
  `cpu_cores` decimal(6,3) DEFAULT NULL COMMENT 'CPU核心数',
  `memory_gb` decimal(8,2) DEFAULT NULL COMMENT '内存(GB)',
  `gpu_count` int DEFAULT '0' COMMENT 'GPU数量',
  `gpu_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'GPU类型',
  `storage_gb` decimal(8,2) DEFAULT NULL COMMENT '存储(GB)',
  `min_replicas` int DEFAULT '1' COMMENT '最小副本数',
  `max_replicas` int DEFAULT '1' COMMENT '最大副本数',
  `auto_scaling` tinyint(1) DEFAULT '0' COMMENT '是否自动扩容',
  `scaling_policy` json DEFAULT NULL COMMENT '扩容策略',
  `target_cpu_utilization` int DEFAULT '70' COMMENT '目标CPU使用率',
  `target_memory_utilization` int DEFAULT '70' COMMENT '目标内存使用率',
  `target_qps` int DEFAULT NULL COMMENT '目标QPS',
  `endpoint_url` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '访问端点',
  `internal_endpoint` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '内部端点',
  `load_balancer_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '负载均衡类型',
  `ingress_config` json DEFAULT NULL COMMENT 'Ingress配置',
  `docker_image` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Docker镜像',
  `command` json DEFAULT NULL COMMENT '启动命令',
  `args` json DEFAULT NULL COMMENT '命令参数',
  `environment_vars` json DEFAULT NULL COMMENT '环境变量',
  `volume_mounts` json DEFAULT NULL COMMENT '挂载卷',
  `config_maps` json DEFAULT NULL COMMENT '配置映射',
  `secrets` json DEFAULT NULL COMMENT '密钥配置',
  `health_check_path` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT '/health' COMMENT '健康检查路径',
  `readiness_probe` json DEFAULT NULL COMMENT '就绪探针',
  `liveness_probe` json DEFAULT NULL COMMENT '存活探针',
  `startup_probe` json DEFAULT NULL COMMENT '启动探针',
  `status` enum('deploying','running','updating','stopped','failed','terminating') COLLATE utf8mb4_unicode_ci DEFAULT 'deploying' COMMENT '部署状态',
  `health_status` enum('healthy','unhealthy','unknown') COLLATE utf8mb4_unicode_ci DEFAULT 'unknown' COMMENT '健康状态',
  `replica_count` int DEFAULT '0' COMMENT '当前副本数',
  `ready_replicas` int DEFAULT '0' COMMENT '就绪副本数',
  `available_replicas` int DEFAULT '0' COMMENT '可用副本数',
  `current_qps` decimal(12,2) DEFAULT '0.00' COMMENT '当前QPS',
  `avg_response_time_ms` decimal(10,2) DEFAULT '0.00' COMMENT '平均响应时间(毫秒)',
  `p99_response_time_ms` decimal(10,2) DEFAULT '0.00' COMMENT 'P99响应时间(毫秒)',
  `error_rate_percent` decimal(5,2) DEFAULT '0.00' COMMENT '错误率',
  `cpu_usage_percent` decimal(5,2) DEFAULT '0.00' COMMENT 'CPU使用率',
  `memory_usage_percent` decimal(5,2) DEFAULT '0.00' COMMENT '内存使用率',
  `gpu_usage_percent` decimal(5,2) DEFAULT '0.00' COMMENT 'GPU使用率',
  `rate_limit_rpm` int DEFAULT NULL COMMENT '每分钟请求限制',
  `quota_limit_per_day` int DEFAULT NULL COMMENT '每日配额限制',
  `concurrent_requests_limit` int DEFAULT NULL COMMENT '并发请求限制',
  `deployed_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '部署时间',
  `last_updated_at` timestamp NULL DEFAULT NULL COMMENT '最后更新时间',
  `stopped_at` timestamp NULL DEFAULT NULL COMMENT '停止时间',
  `traffic_split` json DEFAULT NULL COMMENT '流量分割配置',
  `canary_config` json DEFAULT NULL COMMENT '金丝雀部署配置',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `deployment_config` json DEFAULT NULL COMMENT '部署配置',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_deployment_name` (`deployment_name`),
  KEY `idx_model_id` (`model_id`),
  KEY `idx_version_id` (`version_id`),
  KEY `idx_deployment_type` (`deployment_type`),
  KEY `idx_environment` (`environment`),
  KEY `idx_status` (`status`),
  KEY `idx_health_status` (`health_status`),
  KEY `idx_deployed_at` (`deployed_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模型部署表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_model_deployments`
--

LOCK TABLES `vt_model_deployments` WRITE;
/*!40000 ALTER TABLE `vt_model_deployments` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_model_deployments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_model_inference_logs`
--

DROP TABLE IF EXISTS `vt_model_inference_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_model_inference_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `deployment_id` bigint NOT NULL COMMENT '部署ID',
  `model_id` bigint NOT NULL COMMENT '模型ID',
  `version_id` bigint NOT NULL COMMENT '版本ID',
  `request_id` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '请求ID',
  `client_ip` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '客户端IP',
  `user_agent` text COLLATE utf8mb4_unicode_ci COMMENT '用户代理',
  `request_method` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '请求方法',
  `request_path` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '请求路径',
  `input_data` json DEFAULT NULL COMMENT '输入数据',
  `output_data` json DEFAULT NULL COMMENT '输出数据',
  `inference_time_ms` decimal(10,2) DEFAULT NULL COMMENT '推理时间(毫秒)',
  `total_time_ms` decimal(10,2) DEFAULT NULL COMMENT '总时间(毫秒)',
  `status` enum('success','error','timeout') COLLATE utf8mb4_unicode_ci DEFAULT 'success' COMMENT '状态',
  `error_message` text COLLATE utf8mb4_unicode_ci COMMENT '错误信息',
  `http_status_code` int DEFAULT NULL COMMENT 'HTTP状态码',
  `request_time` timestamp NOT NULL COMMENT '请求时间',
  `response_time` timestamp NULL DEFAULT NULL COMMENT '响应时间',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  PRIMARY KEY (`id`),
  KEY `idx_deployment_id` (`deployment_id`),
  KEY `idx_model_id` (`model_id`),
  KEY `idx_version_id` (`version_id`),
  KEY `idx_request_id` (`request_id`),
  KEY `idx_status` (`status`),
  KEY `idx_request_time` (`request_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模型推理记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_model_inference_logs`
--

LOCK TABLES `vt_model_inference_logs` WRITE;
/*!40000 ALTER TABLE `vt_model_inference_logs` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_model_inference_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_model_inference_relations`
--

DROP TABLE IF EXISTS `vt_model_inference_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_model_inference_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `inference_log_id` bigint NOT NULL COMMENT '推理日志ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(user, workspace, api_key等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(user, workspace, api_key等)',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_inference_entity_relation` (`inference_log_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_inference_log_id` (`inference_log_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_entity_relation` (`entity_type`,`relation_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模型推理日志关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_model_inference_relations`
--

LOCK TABLES `vt_model_inference_relations` WRITE;
/*!40000 ALTER TABLE `vt_model_inference_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_model_inference_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_model_relations`
--

DROP TABLE IF EXISTS `vt_model_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_model_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `model_id` bigint NOT NULL COMMENT '模型ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(workspace, user, dataset, training_job, file等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(workspace, owner, creator, training_dataset, validation_dataset, test_dataset, training_job, model_file, config_file, weights_file, tokenizer_file等)',
  `is_primary` tinyint(1) DEFAULT '0' COMMENT '是否主要关联',
  `sort_order` int DEFAULT '0' COMMENT '排序',
  `status` enum('active','inactive','pending','deleted') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_model_entity_relation` (`model_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_model_id` (`model_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_status` (`status`),
  KEY `idx_entity_relation` (`entity_type`,`relation_type`),
  KEY `idx_model_relation_composite` (`entity_type`,`entity_id`,`relation_type`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模型关联关系表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_model_relations`
--

LOCK TABLES `vt_model_relations` WRITE;
/*!40000 ALTER TABLE `vt_model_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_model_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_model_version_relations`
--

DROP TABLE IF EXISTS `vt_model_version_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_model_version_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `version_id` bigint NOT NULL COMMENT '版本ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(user, file等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(created_by, model_file, config_file, weights_file等)',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_version_entity_relation` (`version_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_version_id` (`version_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模型版本关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_model_version_relations`
--

LOCK TABLES `vt_model_version_relations` WRITE;
/*!40000 ALTER TABLE `vt_model_version_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_model_version_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_model_versions`
--

DROP TABLE IF EXISTS `vt_model_versions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_model_versions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `model_id` bigint NOT NULL COMMENT '模型ID',
  `version` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '版本号',
  `version_name` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '版本名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '版本描述',
  `change_log` text COLLATE utf8mb4_unicode_ci COMMENT '变更日志',
  `parent_version_id` bigint DEFAULT NULL COMMENT '父版本ID',
  `version_type` enum('major','minor','patch','hotfix') COLLATE utf8mb4_unicode_ci DEFAULT 'minor' COMMENT '版本类型',
  `model_size_mb` decimal(12,2) DEFAULT NULL COMMENT '模型大小(MB)',
  `model_hash` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '模型文件哈希',
  `checksum` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '校验和',
  `accuracy` decimal(8,6) DEFAULT NULL COMMENT '准确率',
  `precision_score` decimal(8,6) DEFAULT NULL COMMENT '精确率',
  `recall` decimal(8,6) DEFAULT NULL COMMENT '召回率',
  `f1_score` decimal(8,6) DEFAULT NULL COMMENT 'F1分数',
  `auc_score` decimal(8,6) DEFAULT NULL COMMENT 'AUC分数',
  `inference_time_ms` decimal(10,2) DEFAULT NULL COMMENT '推理时间(毫秒)',
  `throughput_qps` decimal(12,2) DEFAULT NULL COMMENT '吞吐量(QPS)',
  `memory_usage_mb` decimal(10,2) DEFAULT NULL COMMENT '内存使用(MB)',
  `training_config` json DEFAULT NULL COMMENT '训练配置',
  `hyperparameters` json DEFAULT NULL COMMENT '超参数',
  `training_metrics` json DEFAULT NULL COMMENT '训练指标',
  `evaluation_results` json DEFAULT NULL COMMENT '评估结果',
  `benchmark_results` json DEFAULT NULL COMMENT '基准测试结果',
  `status` enum('training','trained','testing','ready','deprecated','error') COLLATE utf8mb4_unicode_ci DEFAULT 'training' COMMENT '版本状态',
  `is_default` tinyint(1) DEFAULT '0' COMMENT '是否默认版本',
  `is_stable` tinyint(1) DEFAULT '0' COMMENT '是否稳定版本',
  `framework_version` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '框架版本',
  `python_version` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Python版本',
  `cuda_version` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'CUDA版本',
  `dependencies` json DEFAULT NULL COMMENT '依赖库',
  `docker_image` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Docker镜像',
  `compatibility` json DEFAULT NULL COMMENT '兼容性信息',
  `migration_guide` text COLLATE utf8mb4_unicode_ci COMMENT '迁移指南',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_model_version` (`model_id`,`version`),
  KEY `idx_model_id` (`model_id`),
  KEY `idx_version` (`version`),
  KEY `idx_version_type` (`version_type`),
  KEY `idx_status` (`status`),
  KEY `idx_is_default` (`is_default`),
  KEY `idx_is_stable` (`is_stable`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_parent_version_id` (`parent_version_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模型版本表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_model_versions`
--

LOCK TABLES `vt_model_versions` WRITE;
/*!40000 ALTER TABLE `vt_model_versions` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_model_versions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_models`
--

DROP TABLE IF EXISTS `vt_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_models` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模型名称',
  `display_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '显示名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '模型描述',
  `model_type` enum('classification','detection','segmentation','nlp','generative','recommendation','time_series','custom') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模型类型',
  `framework` enum('pytorch','tensorflow','onnx','keras','sklearn','huggingface','paddlepaddle','mindspore','custom') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '框架类型',
  `framework_version` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '框架版本',
  `architecture` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '模型架构',
  `base_model` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '基础模型',
  `model_family` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '模型系列',
  `version` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT '1.0.0' COMMENT '版本号',
  `version_description` text COLLATE utf8mb4_unicode_ci COMMENT '版本说明',
  `is_latest` tinyint(1) DEFAULT '1' COMMENT '是否最新版本',
  `model_size_mb` decimal(12,2) DEFAULT NULL COMMENT '模型大小(MB)',
  `parameter_count` bigint DEFAULT NULL COMMENT '参数数量',
  `flops` bigint DEFAULT NULL COMMENT '浮点运算数',
  `model_depth` int DEFAULT NULL COMMENT '模型深度',
  `accuracy` decimal(8,6) DEFAULT NULL COMMENT '准确率',
  `precision_score` decimal(8,6) DEFAULT NULL COMMENT '精确率',
  `recall` decimal(8,6) DEFAULT NULL COMMENT '召回率',
  `f1_score` decimal(8,6) DEFAULT NULL COMMENT 'F1分数',
  `auc_score` decimal(8,6) DEFAULT NULL COMMENT 'AUC分数',
  `loss_value` decimal(15,8) DEFAULT NULL COMMENT '损失值',
  `inference_time_ms` decimal(10,2) DEFAULT NULL COMMENT '推理时间(毫秒)',
  `throughput_qps` decimal(12,2) DEFAULT NULL COMMENT '吞吐量(QPS)',
  `storage_type` enum('local','s3','oss','registry','git_lfs') COLLATE utf8mb4_unicode_ci DEFAULT 'local' COMMENT '存储类型',
  `storage_path` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '存储路径',
  `input_shape` json DEFAULT NULL COMMENT '输入形状',
  `output_shape` json DEFAULT NULL COMMENT '输出形状',
  `input_dtype` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '输入数据类型',
  `output_dtype` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '输出数据类型',
  `input_format` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '输入格式',
  `output_format` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '输出格式',
  `training_config` json DEFAULT NULL COMMENT '训练配置',
  `hyperparameters` json DEFAULT NULL COMMENT '超参数',
  `training_metrics` json DEFAULT NULL COMMENT '训练指标',
  `training_duration_hours` decimal(8,2) DEFAULT NULL COMMENT '训练时长(小时)',
  `status` enum('training','trained','testing','ready','deploying','deployed','error','archived','deprecated') COLLATE utf8mb4_unicode_ci DEFAULT 'training' COMMENT '模型状态',
  `visibility` enum('public','private','workspace','shared') COLLATE utf8mb4_unicode_ci DEFAULT 'private' COMMENT '可见性',
  `is_featured` tinyint(1) DEFAULT '0' COMMENT '是否推荐模型',
  `is_verified` tinyint(1) DEFAULT '0' COMMENT '是否已验证',
  `deployment_count` int DEFAULT '0' COMMENT '部署次数',
  `active_deployments` int DEFAULT '0' COMMENT '活跃部署数',
  `download_count` int DEFAULT '0' COMMENT '下载次数',
  `view_count` int DEFAULT '0' COMMENT '查看次数',
  `usage_count` int DEFAULT '0' COMMENT '使用次数',
  `star_count` int DEFAULT '0' COMMENT '收藏次数',
  `fork_count` int DEFAULT '0' COMMENT '分叉次数',
  `dependencies` json DEFAULT NULL COMMENT '依赖库',
  `requirements` json DEFAULT NULL COMMENT '环境要求',
  `docker_image` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Docker镜像',
  `python_version` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Python版本',
  `cuda_version` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'CUDA版本',
  `license` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '许可证',
  `license_url` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '许可证URL',
  `ethical_considerations` text COLLATE utf8mb4_unicode_ci COMMENT '伦理考虑',
  `bias_report` json DEFAULT NULL COMMENT '偏见报告',
  `tags` json DEFAULT NULL COMMENT '标签',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `readme` text COLLATE utf8mb4_unicode_ci COMMENT '说明文档',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_model_type` (`model_type`),
  KEY `idx_framework` (`framework`),
  KEY `idx_architecture` (`architecture`),
  KEY `idx_status` (`status`),
  KEY `idx_visibility` (`visibility`),
  KEY `idx_is_featured` (`is_featured`),
  KEY `idx_is_verified` (`is_verified`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_type_framework` (`model_type`,`framework`),
  KEY `idx_status_visibility` (`status`,`visibility`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='模型表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_models`
--

LOCK TABLES `vt_models` WRITE;
/*!40000 ALTER TABLE `vt_models` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_models` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_monitor_data`
--

DROP TABLE IF EXISTS `vt_monitor_data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_monitor_data` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `metric_id` bigint NOT NULL COMMENT '指标ID',
  `resource_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '资源类型',
  `resource_id` bigint DEFAULT NULL COMMENT '资源ID',
  `resource_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '资源名称',
  `instance_id` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '实例ID',
  `labels` json DEFAULT NULL COMMENT '标签',
  `value` decimal(20,6) NOT NULL COMMENT '指标值',
  `value_int` bigint DEFAULT NULL COMMENT '整数值',
  `value_str` text COLLATE utf8mb4_unicode_ci COMMENT '字符串值',
  `value_bool` tinyint(1) DEFAULT NULL COMMENT '布尔值',
  `count_value` bigint DEFAULT NULL COMMENT '计数值',
  `sum_value` decimal(20,6) DEFAULT NULL COMMENT '求和值',
  `min_value` decimal(20,6) DEFAULT NULL COMMENT '最小值',
  `max_value` decimal(20,6) DEFAULT NULL COMMENT '最大值',
  `avg_value` decimal(20,6) DEFAULT NULL COMMENT '平均值',
  `timestamp` timestamp NOT NULL COMMENT '时间戳',
  `collection_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '采集时间',
  `quality_score` decimal(3,2) DEFAULT '1.00' COMMENT '数据质量评分',
  `is_anomaly` tinyint(1) DEFAULT '0' COMMENT '是否异常数据',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  PRIMARY KEY (`id`),
  KEY `idx_metric_id` (`metric_id`),
  KEY `idx_resource` (`resource_type`,`resource_id`),
  KEY `idx_instance_id` (`instance_id`),
  KEY `idx_timestamp` (`timestamp`),
  KEY `idx_metric_timestamp` (`metric_id`,`timestamp`),
  KEY `idx_resource_timestamp` (`resource_type`,`resource_id`,`timestamp`),
  KEY `idx_collection_time` (`collection_time`),
  KEY `idx_is_anomaly` (`is_anomaly`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='监控数据表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_monitor_data`
--

LOCK TABLES `vt_monitor_data` WRITE;
/*!40000 ALTER TABLE `vt_monitor_data` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_monitor_data` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_monitor_metrics`
--

DROP TABLE IF EXISTS `vt_monitor_metrics`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_monitor_metrics` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '指标名称',
  `display_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '显示名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '指标描述',
  `metric_type` enum('counter','gauge','histogram','summary') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '指标类型',
  `data_type` enum('integer','float','boolean','string') COLLATE utf8mb4_unicode_ci DEFAULT 'float' COMMENT '数据类型',
  `category` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '指标类别',
  `module` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '所属模块',
  `source_type` enum('system','application','business','custom') COLLATE utf8mb4_unicode_ci DEFAULT 'system' COMMENT '数据源类型',
  `unit` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '单位',
  `aggregation_type` enum('sum','avg','min','max','count','last','rate') COLLATE utf8mb4_unicode_ci DEFAULT 'avg' COMMENT '聚合方式',
  `collection_interval_seconds` int DEFAULT '60' COMMENT '采集间隔(秒)',
  `retention_days` int DEFAULT '30' COMMENT '保留天数',
  `normal_range_min` decimal(20,6) DEFAULT NULL COMMENT '正常范围最小值',
  `normal_range_max` decimal(20,6) DEFAULT NULL COMMENT '正常范围最大值',
  `warning_threshold` decimal(20,6) DEFAULT NULL COMMENT '警告阈值',
  `critical_threshold` decimal(20,6) DEFAULT NULL COMMENT '严重阈值',
  `threshold_condition` enum('gt','gte','lt','lte','eq','neq','between','outside') COLLATE utf8mb4_unicode_ci DEFAULT 'gt' COMMENT '阈值条件',
  `status` enum('active','inactive','deprecated') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `is_builtin` tinyint(1) DEFAULT '0' COMMENT '是否内置指标',
  `is_core` tinyint(1) DEFAULT '0' COMMENT '是否核心指标',
  `default_labels` json DEFAULT NULL COMMENT '默认标签',
  `dimensions` json DEFAULT NULL COMMENT '维度定义',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_metric_type` (`metric_type`),
  KEY `idx_category` (`category`),
  KEY `idx_module` (`module`),
  KEY `idx_source_type` (`source_type`),
  KEY `idx_status` (`status`),
  KEY `idx_is_builtin` (`is_builtin`),
  KEY `idx_is_core` (`is_core`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='监控指标定义表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_monitor_metrics`
--

LOCK TABLES `vt_monitor_metrics` WRITE;
/*!40000 ALTER TABLE `vt_monitor_metrics` DISABLE KEYS */;
INSERT INTO `vt_monitor_metrics` VALUES (1,'cpu_usage_percent','CPU使用率','系统CPU使用率百分比','gauge','float','system','system','system','%','avg',60,30,NULL,NULL,NULL,NULL,'gt','active',0,0,NULL,NULL,NULL,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(2,'memory_usage_percent','内存使用率','系统内存使用率百分比','gauge','float','system','system','system','%','avg',60,30,NULL,NULL,NULL,NULL,'gt','active',0,0,NULL,NULL,NULL,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(3,'gpu_usage_percent','GPU使用率','GPU使用率百分比','gauge','float','gpu','gpu','system','%','avg',60,30,NULL,NULL,NULL,NULL,'gt','active',0,0,NULL,NULL,NULL,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(4,'gpu_memory_usage_percent','GPU显存使用率','GPU显存使用率百分比','gauge','float','gpu','gpu','system','%','avg',60,30,NULL,NULL,NULL,NULL,'gt','active',0,0,NULL,NULL,NULL,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(5,'disk_usage_percent','磁盘使用率','磁盘使用率百分比','gauge','float','storage','system','system','%','avg',60,30,NULL,NULL,NULL,NULL,'gt','active',0,0,NULL,NULL,NULL,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(6,'network_bytes_in','网络入流量','网络入站流量','counter','float','network','system','system','bytes','avg',60,30,NULL,NULL,NULL,NULL,'gt','active',0,0,NULL,NULL,NULL,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(7,'network_bytes_out','网络出流量','网络出站流量','counter','float','network','system','system','bytes','avg',60,30,NULL,NULL,NULL,NULL,'gt','active',0,0,NULL,NULL,NULL,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(8,'training_job_count','训练任务数','当前运行的训练任务数量','gauge','float','training','training','system','count','avg',60,30,NULL,NULL,NULL,NULL,'gt','active',0,0,NULL,NULL,NULL,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(9,'model_inference_qps','模型推理QPS','模型推理每秒请求数','gauge','float','inference','model','system','qps','avg',60,30,NULL,NULL,NULL,NULL,'gt','active',0,0,NULL,NULL,NULL,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(10,'model_inference_latency_ms','模型推理延迟','模型推理平均延迟','histogram','float','inference','model','system','ms','avg',60,30,NULL,NULL,NULL,NULL,'gt','active',0,0,NULL,NULL,NULL,'2025-06-25 04:42:50','2025-06-25 04:42:50');
/*!40000 ALTER TABLE `vt_monitor_metrics` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_monitor_relations`
--

DROP TABLE IF EXISTS `vt_monitor_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_monitor_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联实体类型',
  `relation_entity_id` bigint NOT NULL COMMENT '关联实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_monitor_relation` (`entity_type`,`entity_id`,`relation_entity_type`,`relation_entity_id`,`relation_type`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_entity` (`relation_entity_type`,`relation_entity_id`),
  KEY `idx_relation_type` (`relation_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='监控关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_monitor_relations`
--

LOCK TABLES `vt_monitor_relations` WRITE;
/*!40000 ALTER TABLE `vt_monitor_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_monitor_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_notification_channels`
--

DROP TABLE IF EXISTS `vt_notification_channels`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_notification_channels` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '渠道名称',
  `display_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '显示名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '渠道描述',
  `channel_type` enum('email','sms','webhook','slack','dingtalk','wechat','teams','pagerduty','custom') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '渠道类型',
  `config` json NOT NULL COMMENT '渠道配置',
  `auth_config` json DEFAULT NULL COMMENT '认证配置',
  `template_config` json DEFAULT NULL COMMENT '模板配置',
  `rate_limit_per_minute` int DEFAULT '60' COMMENT '每分钟限制次数',
  `rate_limit_per_hour` int DEFAULT '1000' COMMENT '每小时限制次数',
  `rate_limit_per_day` int DEFAULT '10000' COMMENT '每天限制次数',
  `retry_enabled` tinyint(1) DEFAULT '1' COMMENT '是否启用重试',
  `max_retry_count` int DEFAULT '3' COMMENT '最大重试次数',
  `retry_interval_seconds` int DEFAULT '60' COMMENT '重试间隔(秒)',
  `alert_level_filter` json DEFAULT NULL COMMENT '告警级别过滤',
  `time_filter` json DEFAULT NULL COMMENT '时间过滤',
  `content_filter` json DEFAULT NULL COMMENT '内容过滤',
  `status` enum('active','inactive','error') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '渠道状态',
  `last_used_at` timestamp NULL DEFAULT NULL COMMENT '最后使用时间',
  `success_count` int DEFAULT '0' COMMENT '成功次数',
  `error_count` int DEFAULT '0' COMMENT '失败次数',
  `total_sent` int DEFAULT '0' COMMENT '总发送次数',
  `test_enabled` tinyint(1) DEFAULT '1' COMMENT '是否启用测试',
  `last_test_at` timestamp NULL DEFAULT NULL COMMENT '最后测试时间',
  `test_result` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '测试结果',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_channel_type` (`channel_type`),
  KEY `idx_status` (`status`),
  KEY `idx_last_used_at` (`last_used_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知渠道表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_notification_channels`
--

LOCK TABLES `vt_notification_channels` WRITE;
/*!40000 ALTER TABLE `vt_notification_channels` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_notification_channels` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_notification_relations`
--

DROP TABLE IF EXISTS `vt_notification_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_notification_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `notification_id` bigint NOT NULL COMMENT '通知消息ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(user等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(recipient等)',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_notification_entity_relation` (`notification_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_notification_id` (`notification_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知消息关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_notification_relations`
--

LOCK TABLES `vt_notification_relations` WRITE;
/*!40000 ALTER TABLE `vt_notification_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_notification_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_notification_templates`
--

DROP TABLE IF EXISTS `vt_notification_templates`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_notification_templates` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模板名称',
  `display_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '显示名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '模板描述',
  `template_type` enum('alert','report','reminder','custom') COLLATE utf8mb4_unicode_ci DEFAULT 'alert' COMMENT '模板类型',
  `channel_type` enum('email','sms','webhook','slack','dingtalk','wechat','teams','all') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '适用渠道类型',
  `subject_template` text COLLATE utf8mb4_unicode_ci COMMENT '主题模板',
  `body_template` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '内容模板',
  `html_template` text COLLATE utf8mb4_unicode_ci COMMENT 'HTML模板',
  `variables` json DEFAULT NULL COMMENT '模板变量定义',
  `default_values` json DEFAULT NULL COMMENT '默认值',
  `format_type` enum('text','html','markdown','json') COLLATE utf8mb4_unicode_ci DEFAULT 'text' COMMENT '格式类型',
  `encoding` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'utf-8' COMMENT '编码',
  `locale` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT 'zh-CN' COMMENT '语言区域',
  `timezone` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT 'Asia/Shanghai' COMMENT '时区',
  `status` enum('active','inactive','deprecated') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `is_default` tinyint(1) DEFAULT '0' COMMENT '是否默认模板',
  `is_system` tinyint(1) DEFAULT '0' COMMENT '是否系统模板',
  `version` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'v1.0' COMMENT '版本',
  `parent_template_id` bigint DEFAULT NULL COMMENT '父模板ID',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_template_type` (`template_type`),
  KEY `idx_channel_type` (`channel_type`),
  KEY `idx_status` (`status`),
  KEY `idx_is_default` (`is_default`),
  KEY `idx_is_system` (`is_system`),
  KEY `idx_parent_template_id` (`parent_template_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知模板表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_notification_templates`
--

LOCK TABLES `vt_notification_templates` WRITE;
/*!40000 ALTER TABLE `vt_notification_templates` DISABLE KEYS */;
INSERT INTO `vt_notification_templates` VALUES (1,'training_job_completed','训练任务完成通知',NULL,'alert','email','训练任务 {{job_name}} 已完成','您的训练任务 \"{{job_name}}\" 已成功完成。\n\n任务详情：\n- 任务ID: {{job_id}}\n- 状态: {{status}}\n- 开始时间: {{start_time}}\n- 结束时间: {{end_time}}\n- 运行时长: {{duration}}\n\n请登录平台查看详细结果。',NULL,NULL,NULL,'text','utf-8','zh-CN','Asia/Shanghai','active',0,0,'v1.0',NULL,NULL,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(2,'training_job_failed','训练任务失败通知',NULL,'alert','email','训练任务 {{job_name}} 执行失败','您的训练任务 \"{{job_name}}\" 执行失败。\n\n任务详情：\n- 任务ID: {{job_id}}\n- 状态: {{status}}\n- 错误信息: {{error_message}}\n- 失败时间: {{failed_time}}\n\n请检查任务配置和代码，然后重新提交。',NULL,NULL,NULL,'text','utf-8','zh-CN','Asia/Shanghai','active',0,0,'v1.0',NULL,NULL,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(3,'gpu_usage_high','GPU使用率告警',NULL,'alert','email','GPU使用率过高告警','GPU使用率超过阈值\n\n告警详情：\n- 节点: {{node_name}}\n- GPU设备: {{gpu_device}}\n- 当前使用率: {{current_usage}}%\n- 告警阈值: {{threshold}}%\n- 告警时间: {{alert_time}}\n\n请及时检查GPU资源使用情况。',NULL,NULL,NULL,'text','utf-8','zh-CN','Asia/Shanghai','active',0,0,'v1.0',NULL,NULL,'2025-06-25 04:42:50','2025-06-25 04:42:50');
/*!40000 ALTER TABLE `vt_notification_templates` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_notifications`
--

DROP TABLE IF EXISTS `vt_notifications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_notifications` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `title` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '通知标题',
  `content` text COLLATE utf8mb4_unicode_ci COMMENT '通知内容',
  `notification_type` enum('system','training','deployment','alert','workspace') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '通知类型',
  `priority` enum('low','normal','high','urgent') COLLATE utf8mb4_unicode_ci DEFAULT 'normal' COMMENT '优先级',
  `status` enum('unread','read','archived') COLLATE utf8mb4_unicode_ci DEFAULT 'unread' COMMENT '状态',
  `resource_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '关联资源类型',
  `resource_id` bigint DEFAULT NULL COMMENT '关联资源ID',
  `metadata` json DEFAULT NULL COMMENT '扩展数据',
  `sent_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
  `read_at` timestamp NULL DEFAULT NULL COMMENT '阅读时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_notification_type` (`notification_type`),
  KEY `idx_priority` (`priority`),
  KEY `idx_status` (`status`),
  KEY `idx_resource` (`resource_type`,`resource_id`),
  KEY `idx_sent_at` (`sent_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知消息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_notifications`
--

LOCK TABLES `vt_notifications` WRITE;
/*!40000 ALTER TABLE `vt_notifications` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_notifications` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_operation_log_relations`
--

DROP TABLE IF EXISTS `vt_operation_log_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_operation_log_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `operation_log_id` bigint NOT NULL COMMENT '操作日志ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(user, api_key等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(user, api_key等)',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_operation_entity_relation` (`operation_log_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_operation_log_id` (`operation_log_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_operation_log_relations`
--

LOCK TABLES `vt_operation_log_relations` WRITE;
/*!40000 ALTER TABLE `vt_operation_log_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_operation_log_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_operation_logs`
--

DROP TABLE IF EXISTS `vt_operation_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_operation_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `operation_type` enum('create','update','delete','login','logout','export','import','execute','deploy','start','stop','restart') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作类型',
  `resource_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '资源类型',
  `resource_id` bigint DEFAULT NULL COMMENT '资源ID',
  `resource_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '资源名称',
  `request_method` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '请求方法',
  `request_path` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '请求路径',
  `request_params` json DEFAULT NULL COMMENT '请求参数',
  `request_body` json DEFAULT NULL COMMENT '请求体',
  `response_code` int DEFAULT NULL COMMENT '响应代码',
  `response_time_ms` int DEFAULT NULL COMMENT '响应时间(毫秒)',
  `error_message` text COLLATE utf8mb4_unicode_ci COMMENT '错误信息',
  `ip_address` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'IP地址',
  `user_agent` text COLLATE utf8mb4_unicode_ci COMMENT '用户代理',
  `session_id` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '会话ID',
  `old_value` json DEFAULT NULL COMMENT '旧值',
  `new_value` json DEFAULT NULL COMMENT '新值',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_operation_type` (`operation_type`),
  KEY `idx_resource` (`resource_type`,`resource_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_operation_logs`
--

LOCK TABLES `vt_operation_logs` WRITE;
/*!40000 ALTER TABLE `vt_operation_logs` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_operation_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_permissions`
--

DROP TABLE IF EXISTS `vt_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_permissions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限标识',
  `display_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '权限名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '权限描述',
  `module` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '所属模块',
  `action` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '操作类型',
  `resource` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '资源类型',
  `permission_code` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '权限编码',
  `parent_id` bigint DEFAULT '0' COMMENT '父权限ID',
  `level` int DEFAULT '1' COMMENT '权限层级',
  `sort_order` int DEFAULT '0' COMMENT '排序',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  UNIQUE KEY `uk_permission_code` (`permission_code`),
  KEY `idx_module` (`module`),
  KEY `idx_action` (`action`),
  KEY `idx_resource` (`resource`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_permissions`
--

LOCK TABLES `vt_permissions` WRITE;
/*!40000 ALTER TABLE `vt_permissions` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_pod_event_relations`
--

DROP TABLE IF EXISTS `vt_pod_event_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_pod_event_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `pod_event_id` bigint NOT NULL COMMENT 'Pod事件ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(training_job, k8s_resource等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(training_job, k8s_resource等)',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_pod_event_entity_relation` (`pod_event_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_pod_event_id` (`pod_event_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Pod事件关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_pod_event_relations`
--

LOCK TABLES `vt_pod_event_relations` WRITE;
/*!40000 ALTER TABLE `vt_pod_event_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_pod_event_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_pod_events`
--

DROP TABLE IF EXISTS `vt_pod_events`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_pod_events` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `pod_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Pod名称',
  `namespace` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '命名空间',
  `uid` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Pod UID',
  `event_type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '事件类型：Normal, Warning',
  `reason` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '事件原因',
  `message` text COLLATE utf8mb4_unicode_ci COMMENT '事件消息',
  `source_component` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '来源组件',
  `source_host` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '来源主机',
  `first_timestamp` timestamp NOT NULL COMMENT '首次发生时间',
  `last_timestamp` timestamp NOT NULL COMMENT '最后发生时间',
  `count` int DEFAULT '1' COMMENT '发生次数',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_pod` (`namespace`,`pod_name`),
  KEY `idx_event_type` (`event_type`),
  KEY `idx_timestamp` (`last_timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Pod事件记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_pod_events`
--

LOCK TABLES `vt_pod_events` WRITE;
/*!40000 ALTER TABLE `vt_pod_events` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_pod_events` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_project_owners`
--

DROP TABLE IF EXISTS `vt_project_owners`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_project_owners` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `project_id` bigint NOT NULL COMMENT '项目ID',
  `user_id` bigint NOT NULL COMMENT '所有者用户ID',
  `owner_type` enum('primary','co_owner') COLLATE utf8mb4_unicode_ci DEFAULT 'primary' COMMENT '所有者类型',
  `start_date` date DEFAULT NULL COMMENT '开始日期',
  `end_date` date DEFAULT NULL COMMENT '结束日期',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_project_user_type` (`project_id`,`user_id`,`owner_type`),
  KEY `idx_project_id` (`project_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_owner_type` (`owner_type`),
  KEY `idx_status` (`status`),
  KEY `idx_project_owner_composite` (`project_id`,`user_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='项目所有者关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_project_owners`
--

LOCK TABLES `vt_project_owners` WRITE;
/*!40000 ALTER TABLE `vt_project_owners` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_project_owners` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_resource_quotas`
--

DROP TABLE IF EXISTS `vt_resource_quotas`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_resource_quotas` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `quota_type` enum('user','workspace','global') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配额类型',
  `target_id` bigint DEFAULT NULL COMMENT '目标ID(用户ID或工作空间ID)',
  `resource_type` enum('storage','compute','gpu','memory','dataset','model') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资源类型',
  `quota_limit` decimal(20,6) NOT NULL DEFAULT '0.000000' COMMENT '配额限制',
  `quota_used` decimal(20,6) NOT NULL DEFAULT '0.000000' COMMENT '已使用配额',
  `unit` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '单位',
  `reset_cycle` enum('daily','weekly','monthly','never') COLLATE utf8mb4_unicode_ci DEFAULT 'never' COMMENT '重置周期',
  `last_reset_at` timestamp NULL DEFAULT NULL COMMENT '最后重置时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_quota_target_resource` (`quota_type`,`target_id`,`resource_type`),
  KEY `idx_quota_type` (`quota_type`),
  KEY `idx_target_id` (`target_id`),
  KEY `idx_resource_type` (`resource_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统资源配额表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_resource_quotas`
--

LOCK TABLES `vt_resource_quotas` WRITE;
/*!40000 ALTER TABLE `vt_resource_quotas` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_resource_quotas` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_resource_relations`
--

DROP TABLE IF EXISTS `vt_resource_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_resource_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `source_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '源资源类型',
  `source_id` bigint NOT NULL COMMENT '源资源ID',
  `target_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '目标资源类型',
  `target_id` bigint NOT NULL COMMENT '目标资源ID',
  `relation_type` enum('owns','uses','depends_on','produces','contains','references') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型',
  `metadata` json DEFAULT NULL COMMENT '关联元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_resource_relation` (`source_type`,`source_id`,`target_type`,`target_id`,`relation_type`),
  KEY `idx_source` (`source_type`,`source_id`),
  KEY `idx_target` (`target_type`,`target_id`),
  KEY `idx_relation_type` (`relation_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='资源关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_resource_relations`
--

LOCK TABLES `vt_resource_relations` WRITE;
/*!40000 ALTER TABLE `vt_resource_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_resource_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_resource_tag_relations`
--

DROP TABLE IF EXISTS `vt_resource_tag_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_resource_tag_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `resource_tag_id` bigint NOT NULL COMMENT '资源标签ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(user等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(tagged_by等)',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tag_entity_relation` (`resource_tag_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_resource_tag_id` (`resource_tag_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='资源标签关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_resource_tag_relations`
--

LOCK TABLES `vt_resource_tag_relations` WRITE;
/*!40000 ALTER TABLE `vt_resource_tag_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_resource_tag_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_resource_tags`
--

DROP TABLE IF EXISTS `vt_resource_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_resource_tags` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `resource_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资源类型',
  `resource_id` bigint NOT NULL COMMENT '资源ID',
  `tag_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '标签名称',
  `tag_value` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '标签值',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_resource_tag` (`resource_type`,`resource_id`,`tag_name`),
  KEY `idx_resource` (`resource_type`,`resource_id`),
  KEY `idx_tag_name` (`tag_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='资源标签关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_resource_tags`
--

LOCK TABLES `vt_resource_tags` WRITE;
/*!40000 ALTER TABLE `vt_resource_tags` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_resource_tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_role_permissions`
--

DROP TABLE IF EXISTS `vt_role_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_role_permissions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `permission_id` bigint NOT NULL COMMENT '权限ID',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`,`permission_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_permission_id` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_role_permissions`
--

LOCK TABLES `vt_role_permissions` WRITE;
/*!40000 ALTER TABLE `vt_role_permissions` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_role_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_roles`
--

DROP TABLE IF EXISTS `vt_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_roles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色标识',
  `display_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '角色名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '角色描述',
  `role_code` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '角色编码',
  `role_type` enum('system','custom') COLLATE utf8mb4_unicode_ci DEFAULT 'custom' COMMENT '角色类型',
  `sort_order` int DEFAULT '0' COMMENT '排序',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  UNIQUE KEY `uk_role_code` (`role_code`),
  KEY `idx_role_type` (`role_type`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_roles`
--

LOCK TABLES `vt_roles` WRITE;
/*!40000 ALTER TABLE `vt_roles` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_system_config_relations`
--

DROP TABLE IF EXISTS `vt_system_config_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_system_config_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `config_id` bigint NOT NULL COMMENT '系统配置ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(user等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(updated_by等)',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_entity_relation` (`config_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_config_id` (`config_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_system_config_relations`
--

LOCK TABLES `vt_system_config_relations` WRITE;
/*!40000 ALTER TABLE `vt_system_config_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_system_config_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_system_configs`
--

DROP TABLE IF EXISTS `vt_system_configs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_system_configs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `config_key` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置键',
  `config_value` text COLLATE utf8mb4_unicode_ci COMMENT '配置值',
  `config_type` enum('string','number','boolean','json') COLLATE utf8mb4_unicode_ci DEFAULT 'string' COMMENT '配置类型',
  `category` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '配置分类',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '配置描述',
  `is_encrypted` tinyint(1) DEFAULT '0' COMMENT '是否加密',
  `is_readonly` tinyint(1) DEFAULT '0' COMMENT '是否只读',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_key` (`config_key`),
  KEY `idx_category` (`category`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_system_configs`
--

LOCK TABLES `vt_system_configs` WRITE;
/*!40000 ALTER TABLE `vt_system_configs` DISABLE KEYS */;
INSERT INTO `vt_system_configs` VALUES (1,'system.version','1.0.0','string','system','系统版本号',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(2,'system.timezone','Asia/Shanghai','string','system','系统时区',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(3,'system.max_file_size_mb','1024','number','file','最大文件上传大小(MB)',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(4,'system.session_timeout_hours','24','number','security','会话超时时间(小时)',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(5,'system.enable_gpu_monitoring','true','boolean','gpu','是否启用GPU监控',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(6,'system.default_gpu_quota','2','number','gpu','默认GPU配额',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(7,'system.enable_auto_scaling','true','boolean','cluster','是否启用自动扩缩容',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(8,'training.default_timeout_hours','168','number','training','默认训练超时时间(小时)',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(9,'training.max_concurrent_jobs','50','number','training','最大并发训练任务数',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(10,'model.default_storage_type','local','string','model','默认模型存储类型',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(11,'notification.enable_email','true','boolean','notification','是否启用邮件通知',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(12,'monitoring.data_retention_days','90','number','monitoring','监控数据保留天数',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50');
/*!40000 ALTER TABLE `vt_system_configs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_training_checkpoint_files`
--

DROP TABLE IF EXISTS `vt_training_checkpoint_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_training_checkpoint_files` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `checkpoint_id` bigint NOT NULL COMMENT '检查点ID',
  `file_id` bigint NOT NULL COMMENT '文件ID',
  `file_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT 'checkpoint' COMMENT '文件类型',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_checkpoint_file` (`checkpoint_id`,`file_id`),
  KEY `idx_checkpoint_id` (`checkpoint_id`),
  KEY `idx_file_id` (`file_id`),
  KEY `idx_file_type` (`file_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='训练检查点文件关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_training_checkpoint_files`
--

LOCK TABLES `vt_training_checkpoint_files` WRITE;
/*!40000 ALTER TABLE `vt_training_checkpoint_files` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_training_checkpoint_files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_training_checkpoints`
--

DROP TABLE IF EXISTS `vt_training_checkpoints`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_training_checkpoints` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `job_id` bigint NOT NULL COMMENT '训练作业ID',
  `checkpoint_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '检查点名称',
  `checkpoint_type` enum('auto','manual','best','final','scheduled') COLLATE utf8mb4_unicode_ci DEFAULT 'auto' COMMENT '检查点类型',
  `checkpoint_format` enum('pytorch','tensorflow','onnx','pickle','hdf5') COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '检查点格式',
  `step` bigint DEFAULT NULL COMMENT '训练步数',
  `epoch` int DEFAULT NULL COMMENT '训练轮次',
  `global_step` bigint DEFAULT NULL COMMENT '全局步数',
  `storage_path` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '存储路径',
  `file_size` bigint DEFAULT NULL COMMENT '文件大小(字节)',
  `checksum` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件校验和',
  `compression_type` enum('none','gzip','bzip2','lz4') COLLATE utf8mb4_unicode_ci DEFAULT 'none' COMMENT '压缩类型',
  `metrics` json DEFAULT NULL COMMENT '性能指标',
  `loss_value` decimal(20,8) DEFAULT NULL COMMENT '损失值',
  `accuracy` decimal(8,6) DEFAULT NULL COMMENT '准确率',
  `validation_score` decimal(8,6) DEFAULT NULL COMMENT '验证分数',
  `model_config` json DEFAULT NULL COMMENT '模型配置',
  `optimizer_state` json DEFAULT NULL COMMENT '优化器状态',
  `scheduler_state` json DEFAULT NULL COMMENT '调度器状态',
  `status` enum('saving','saved','failed','deleted','corrupted') COLLATE utf8mb4_unicode_ci DEFAULT 'saving' COMMENT '状态',
  `is_best` tinyint(1) DEFAULT '0' COMMENT '是否最佳检查点',
  `is_latest` tinyint(1) DEFAULT '0' COMMENT '是否最新检查点',
  `tags` json DEFAULT NULL COMMENT '标签',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '描述',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `saved_at` timestamp NULL DEFAULT NULL COMMENT '保存完成时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_job_checkpoint` (`job_id`,`checkpoint_name`),
  KEY `idx_job_id` (`job_id`),
  KEY `idx_checkpoint_type` (`checkpoint_type`),
  KEY `idx_step` (`step`),
  KEY `idx_epoch` (`epoch`),
  KEY `idx_status` (`status`),
  KEY `idx_is_best` (`is_best`),
  KEY `idx_is_latest` (`is_latest`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_saved_at` (`saved_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='训练检查点表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_training_checkpoints`
--

LOCK TABLES `vt_training_checkpoints` WRITE;
/*!40000 ALTER TABLE `vt_training_checkpoints` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_training_checkpoints` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_training_job_instances`
--

DROP TABLE IF EXISTS `vt_training_job_instances`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_training_job_instances` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `job_id` bigint NOT NULL COMMENT '所属训练作业ID',
  `instance_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实例名称',
  `instance_type` enum('master','worker','ps','evaluator','chief') COLLATE utf8mb4_unicode_ci DEFAULT 'worker' COMMENT '实例类型',
  `instance_index` int NOT NULL COMMENT '实例索引',
  `replica_index` int DEFAULT NULL COMMENT '副本索引',
  `pod_name` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Pod名称',
  `namespace` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '命名空间',
  `node_name` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '所在节点',
  `node_ip` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '节点IP',
  `pod_ip` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Pod IP',
  `container_id` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '容器ID',
  `allocated_cpu_cores` decimal(6,3) DEFAULT NULL COMMENT '分配的CPU核心数',
  `allocated_memory_gb` decimal(8,2) DEFAULT NULL COMMENT '分配的内存(GB)',
  `allocated_gpu_devices` json DEFAULT NULL COMMENT '分配的GPU设备信息',
  `allocated_storage_gb` decimal(8,2) DEFAULT NULL COMMENT '分配的存储(GB)',
  `status` enum('pending','creating','running','succeeded','failed','killed','unknown') COLLATE utf8mb4_unicode_ci DEFAULT 'pending' COMMENT '实例状态',
  `phase` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Pod阶段',
  `reason` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '状态原因',
  `message` text COLLATE utf8mb4_unicode_ci COMMENT '状态消息',
  `ready` tinyint(1) DEFAULT '0' COMMENT '是否就绪',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `scheduled_at` timestamp NULL DEFAULT NULL COMMENT '调度时间',
  `start_time` timestamp NULL DEFAULT NULL COMMENT '开始时间',
  `end_time` timestamp NULL DEFAULT NULL COMMENT '结束时间',
  `last_transition_time` timestamp NULL DEFAULT NULL COMMENT '最后状态变更时间',
  `restart_count` int DEFAULT '0' COMMENT '重启次数',
  `last_restart_time` timestamp NULL DEFAULT NULL COMMENT '最后重启时间',
  `exit_code` int DEFAULT NULL COMMENT '退出码',
  `termination_reason` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '终止原因',
  `cpu_usage_percent` decimal(5,2) DEFAULT '0.00' COMMENT 'CPU使用率',
  `memory_usage_percent` decimal(5,2) DEFAULT '0.00' COMMENT '内存使用率',
  `gpu_usage_percent` decimal(5,2) DEFAULT '0.00' COMMENT 'GPU使用率',
  `logs_path` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '日志路径',
  `labels` json DEFAULT NULL COMMENT '标签',
  `annotations` json DEFAULT NULL COMMENT '注解',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_job_instance` (`job_id`,`instance_name`),
  KEY `idx_job_id` (`job_id`),
  KEY `idx_instance_name` (`instance_name`),
  KEY `idx_instance_type` (`instance_type`),
  KEY `idx_pod_name` (`pod_name`),
  KEY `idx_node_name` (`node_name`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_start_time` (`start_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='训练任务实例表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_training_job_instances`
--

LOCK TABLES `vt_training_job_instances` WRITE;
/*!40000 ALTER TABLE `vt_training_job_instances` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_training_job_instances` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_training_job_relations`
--

DROP TABLE IF EXISTS `vt_training_job_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_training_job_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `job_id` bigint NOT NULL COMMENT '训练作业ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(workspace, project, user, dataset, model, code_file等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(workspace, project, creator, dataset, base_model, code_file等)',
  `is_primary` tinyint(1) DEFAULT '0' COMMENT '是否主要关联',
  `sort_order` int DEFAULT '0' COMMENT '排序',
  `status` enum('active','inactive','pending','deleted') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_job_entity_relation` (`job_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_job_id` (`job_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_status` (`status`),
  KEY `idx_entity_relation` (`entity_type`,`relation_type`),
  KEY `idx_training_relation_composite` (`entity_type`,`entity_id`,`relation_type`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='训练作业关联关系表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_training_job_relations`
--

LOCK TABLES `vt_training_job_relations` WRITE;
/*!40000 ALTER TABLE `vt_training_job_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_training_job_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_training_jobs`
--

DROP TABLE IF EXISTS `vt_training_jobs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_training_jobs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '任务名称',
  `display_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '显示名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '任务描述',
  `job_type` enum('single','distributed','horovod','pytorch_ddp','parameter_server','federated') COLLATE utf8mb4_unicode_ci DEFAULT 'single' COMMENT '任务类型',
  `framework` enum('pytorch','tensorflow','paddle','mindspore','keras','sklearn','xgboost','custom') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '训练框架',
  `framework_version` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '框架版本',
  `python_version` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT '3.8' COMMENT 'Python版本',
  `code_source_type` enum('git','upload','builtin','notebook') COLLATE utf8mb4_unicode_ci DEFAULT 'upload' COMMENT '代码来源类型',
  `code_source_config` json DEFAULT NULL COMMENT '代码来源配置',
  `entry_point` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '入口脚本',
  `working_dir` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT '/workspace' COMMENT '工作目录',
  `image` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '训练镜像',
  `image_pull_policy` enum('Always','IfNotPresent','Never') COLLATE utf8mb4_unicode_ci DEFAULT 'IfNotPresent' COMMENT '镜像拉取策略',
  `image_pull_secrets` json DEFAULT NULL COMMENT '镜像拉取密钥',
  `dataset_mount_configs` json DEFAULT NULL COMMENT '数据集挂载配置',
  `data_source_config` json DEFAULT NULL COMMENT '数据源配置',
  `model_config` json DEFAULT NULL COMMENT '模型配置',
  `output_model_name` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '输出模型名称',
  `model_save_strategy` enum('all','best','last','custom') COLLATE utf8mb4_unicode_ci DEFAULT 'best' COMMENT '模型保存策略',
  `cpu_cores` decimal(6,3) DEFAULT NULL COMMENT 'CPU核心数',
  `memory_gb` decimal(8,2) DEFAULT NULL COMMENT '内存(GB)',
  `gpu_count` int DEFAULT '0' COMMENT 'GPU数量',
  `gpu_type` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'GPU类型',
  `gpu_memory_gb` decimal(8,2) DEFAULT NULL COMMENT 'GPU显存(GB)',
  `storage_gb` decimal(8,2) DEFAULT NULL COMMENT '存储(GB)',
  `shared_memory_gb` decimal(8,2) DEFAULT NULL COMMENT '共享内存(GB)',
  `worker_count` int DEFAULT '1' COMMENT 'Worker数量',
  `ps_count` int DEFAULT '0' COMMENT 'Parameter Server数量',
  `master_count` int DEFAULT '1' COMMENT 'Master数量',
  `env_vars` json DEFAULT NULL COMMENT '环境变量',
  `command_args` json DEFAULT NULL COMMENT '命令行参数',
  `secrets` json DEFAULT NULL COMMENT '密钥配置',
  `config_maps` json DEFAULT NULL COMMENT '配置映射',
  `volume_mounts` json DEFAULT NULL COMMENT '挂载卷配置',
  `queue_name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT 'default' COMMENT '队列名称',
  `priority` int DEFAULT '0' COMMENT '优先级(0-100)',
  `node_selector` json DEFAULT NULL COMMENT '节点选择器',
  `tolerations` json DEFAULT NULL COMMENT '容忍度设置',
  `affinity` json DEFAULT NULL COMMENT '亲和性配置',
  `max_runtime_seconds` int DEFAULT '86400' COMMENT '最大运行时间(秒)',
  `max_idle_seconds` int DEFAULT '3600' COMMENT '最大空闲时间(秒)',
  `auto_restart` tinyint(1) DEFAULT '0' COMMENT '是否自动重启',
  `max_retry_count` int DEFAULT '3' COMMENT '最大重试次数',
  `volcano_job_name` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Volcano作业名',
  `volcano_queue` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Volcano队列',
  `min_available` int DEFAULT '1' COMMENT '最小可用实例数',
  `status` enum('pending','queued','scheduling','running','succeeded','failed','cancelled','suspended','timeout','oom_killed') COLLATE utf8mb4_unicode_ci DEFAULT 'pending' COMMENT '任务状态',
  `phase` enum('creating','scheduling','initializing','training','saving','completed') COLLATE utf8mb4_unicode_ci DEFAULT 'creating' COMMENT '执行阶段',
  `namespace` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'K8s命名空间',
  `cluster_name` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '集群名称',
  `error_message` text COLLATE utf8mb4_unicode_ci COMMENT '错误信息',
  `error_code` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '错误代码',
  `exit_code` int DEFAULT NULL COMMENT '退出码',
  `failure_reason` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '失败原因',
  `submitted_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '提交时间',
  `queued_at` timestamp NULL DEFAULT NULL COMMENT '进入队列时间',
  `scheduled_at` timestamp NULL DEFAULT NULL COMMENT '调度时间',
  `start_time` timestamp NULL DEFAULT NULL COMMENT '开始时间',
  `end_time` timestamp NULL DEFAULT NULL COMMENT '结束时间',
  `duration_seconds` int DEFAULT NULL COMMENT '执行时长(秒)',
  `actual_cpu_usage` decimal(8,4) DEFAULT NULL COMMENT '实际CPU使用',
  `actual_memory_usage_gb` decimal(8,2) DEFAULT NULL COMMENT '实际内存使用(GB)',
  `actual_gpu_usage` decimal(5,2) DEFAULT NULL COMMENT '实际GPU使用率',
  `peak_memory_usage_gb` decimal(8,2) DEFAULT NULL COMMENT '峰值内存使用(GB)',
  `total_gpu_hours` decimal(10,2) DEFAULT NULL COMMENT '总GPU小时数',
  `workspace_path` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '工作空间路径',
  `logs_path` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '日志路径',
  `output_path` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '输出路径',
  `checkpoint_path` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '检查点路径',
  `tensorboard_path` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'TensorBoard路径',
  `hyperparameters` json DEFAULT NULL COMMENT '超参数',
  `training_config` json DEFAULT NULL COMMENT '训练配置',
  `optimizer_config` json DEFAULT NULL COMMENT '优化器配置',
  `scheduler_config` json DEFAULT NULL COMMENT '调度器配置',
  `enable_tensorboard` tinyint(1) DEFAULT '1' COMMENT '是否启用TensorBoard',
  `enable_profiling` tinyint(1) DEFAULT '0' COMMENT '是否启用性能分析',
  `metrics_collection_interval` int DEFAULT '60' COMMENT '指标收集间隔(秒)',
  `notification_config` json DEFAULT NULL COMMENT '通知配置',
  `tags` json DEFAULT NULL COMMENT '标签',
  `annotations` json DEFAULT NULL COMMENT '注解',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_job_type` (`job_type`),
  KEY `idx_framework` (`framework`),
  KEY `idx_status` (`status`),
  KEY `idx_phase` (`phase`),
  KEY `idx_queue_name` (`queue_name`),
  KEY `idx_priority` (`priority`),
  KEY `idx_volcano_job` (`volcano_job_name`),
  KEY `idx_cluster_name` (`cluster_name`),
  KEY `idx_submitted_at` (`submitted_at`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_status_priority` (`status`,`priority`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='训练作业表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_training_jobs`
--

LOCK TABLES `vt_training_jobs` WRITE;
/*!40000 ALTER TABLE `vt_training_jobs` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_training_jobs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_training_logs`
--

DROP TABLE IF EXISTS `vt_training_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_training_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `job_id` bigint NOT NULL COMMENT '训练作业ID',
  `instance_id` bigint DEFAULT NULL COMMENT '实例ID',
  `log_level` enum('TRACE','DEBUG','INFO','WARN','ERROR','FATAL') COLLATE utf8mb4_unicode_ci DEFAULT 'INFO' COMMENT '日志级别',
  `log_source` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '日志来源',
  `log_content` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '日志内容',
  `log_format` enum('text','json','structured') COLLATE utf8mb4_unicode_ci DEFAULT 'text' COMMENT '日志格式',
  `log_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '日志时间',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `file_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件名',
  `line_number` int DEFAULT NULL COMMENT '行号',
  `function_name` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '函数名',
  `thread_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '线程ID',
  `process_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '进程ID',
  `context` json DEFAULT NULL COMMENT '上下文信息',
  `correlation_id` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '关联ID',
  `category` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '类别',
  `tags` json DEFAULT NULL COMMENT '标签',
  PRIMARY KEY (`id`),
  KEY `idx_job_id` (`job_id`),
  KEY `idx_instance_id` (`instance_id`),
  KEY `idx_log_level` (`log_level`),
  KEY `idx_log_source` (`log_source`),
  KEY `idx_log_time` (`log_time`),
  KEY `idx_category` (`category`),
  KEY `idx_correlation_id` (`correlation_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_job_log_time` (`job_id`,`log_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='训练日志表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_training_logs`
--

LOCK TABLES `vt_training_logs` WRITE;
/*!40000 ALTER TABLE `vt_training_logs` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_training_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_training_metrics`
--

DROP TABLE IF EXISTS `vt_training_metrics`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_training_metrics` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `job_id` bigint NOT NULL COMMENT '训练作业ID',
  `instance_id` bigint DEFAULT NULL COMMENT '实例ID',
  `metric_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '指标名称',
  `metric_type` enum('scalar','histogram','image','text','audio','video') COLLATE utf8mb4_unicode_ci DEFAULT 'scalar' COMMENT '指标类型',
  `metric_value` decimal(20,8) DEFAULT NULL COMMENT '指标值',
  `metric_data` json DEFAULT NULL COMMENT '复杂指标数据',
  `step` bigint DEFAULT NULL COMMENT '训练步数',
  `epoch` int DEFAULT NULL COMMENT '训练轮次',
  `global_step` bigint DEFAULT NULL COMMENT '全局步数',
  `batch_idx` int DEFAULT NULL COMMENT '批次索引',
  `tag` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '标签',
  `category` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '类别',
  `phase` enum('train','val','test') COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '阶段',
  `metric_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '指标时间',
  `wall_time` decimal(15,3) DEFAULT NULL COMMENT '墙钟时间',
  `relative_time_seconds` decimal(15,3) DEFAULT NULL COMMENT '相对时间(秒)',
  `min_value` decimal(20,8) DEFAULT NULL COMMENT '最小值',
  `max_value` decimal(20,8) DEFAULT NULL COMMENT '最大值',
  `avg_value` decimal(20,8) DEFAULT NULL COMMENT '平均值',
  `std_value` decimal(20,8) DEFAULT NULL COMMENT '标准差',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_job_id` (`job_id`),
  KEY `idx_instance_id` (`instance_id`),
  KEY `idx_job_metric` (`job_id`,`metric_name`),
  KEY `idx_job_step` (`job_id`,`step`),
  KEY `idx_job_epoch` (`job_id`,`epoch`),
  KEY `idx_metric_time` (`metric_time`),
  KEY `idx_tag` (`tag`),
  KEY `idx_category` (`category`),
  KEY `idx_phase` (`phase`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='训练指标表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_training_metrics`
--

LOCK TABLES `vt_training_metrics` WRITE;
/*!40000 ALTER TABLE `vt_training_metrics` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_training_metrics` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_training_queues`
--

DROP TABLE IF EXISTS `vt_training_queues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_training_queues` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '队列名称',
  `display_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '显示名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '队列描述',
  `queue_type` enum('default','high_priority','gpu_intensive','cpu_intensive','experimental') COLLATE utf8mb4_unicode_ci DEFAULT 'default' COMMENT '队列类型',
  `priority` int DEFAULT '0' COMMENT '队列优先级(数值越大优先级越高)',
  `max_concurrent_jobs` int DEFAULT '10' COMMENT '最大并发任务数',
  `max_queue_size` int DEFAULT '100' COMMENT '最大队列大小',
  `max_job_duration_hours` int DEFAULT '168' COMMENT '最大任务时长(小时)',
  `resource_quota` json DEFAULT NULL COMMENT '资源配额配置',
  `gpu_quota` int DEFAULT NULL COMMENT 'GPU配额',
  `cpu_quota` decimal(8,2) DEFAULT NULL COMMENT 'CPU配额',
  `memory_quota_gb` int DEFAULT NULL COMMENT '内存配额(GB)',
  `storage_quota_gb` int DEFAULT NULL COMMENT '存储配额(GB)',
  `scheduling_policy` enum('fifo','priority','fair_share','shortest_job_first') COLLATE utf8mb4_unicode_ci DEFAULT 'fifo' COMMENT '调度策略',
  `preemption_enabled` tinyint(1) DEFAULT '0' COMMENT '是否启用抢占',
  `gang_scheduling` tinyint(1) DEFAULT '0' COMMENT '是否启用gang调度',
  `workspace_ids` json DEFAULT NULL COMMENT '允许访问的工作空间ID列表',
  `user_ids` json DEFAULT NULL COMMENT '允许访问的用户ID列表',
  `department_ids` json DEFAULT NULL COMMENT '允许访问的部门ID列表',
  `cluster_ids` json DEFAULT NULL COMMENT '可用集群ID列表',
  `node_selector` json DEFAULT NULL COMMENT '节点选择器',
  `tolerations` json DEFAULT NULL COMMENT '容忍度配置',
  `status` enum('active','disabled','maintenance') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `current_jobs` int DEFAULT '0' COMMENT '当前任务数',
  `pending_jobs` int DEFAULT '0' COMMENT '等待任务数',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_queue_type` (`queue_type`),
  KEY `idx_priority` (`priority`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='训练队列表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_training_queues`
--

LOCK TABLES `vt_training_queues` WRITE;
/*!40000 ALTER TABLE `vt_training_queues` DISABLE KEYS */;
INSERT INTO `vt_training_queues` VALUES (1,'default','默认队列','系统默认训练队列','default',0,10,100,168,NULL,NULL,NULL,NULL,NULL,'fifo',0,0,NULL,NULL,NULL,NULL,NULL,NULL,'active',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50',NULL),(2,'high_priority','高优先级队列','高优先级训练任务队列','high_priority',100,5,100,168,NULL,NULL,NULL,NULL,NULL,'fifo',0,0,NULL,NULL,NULL,NULL,NULL,NULL,'active',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50',NULL),(3,'gpu_intensive','GPU密集型队列','GPU密集型训练任务队列','gpu_intensive',50,3,100,168,NULL,NULL,NULL,NULL,NULL,'fifo',0,0,NULL,NULL,NULL,NULL,NULL,NULL,'active',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50',NULL),(4,'experimental','实验队列','实验性训练任务队列','experimental',-10,20,100,168,NULL,NULL,NULL,NULL,NULL,'fifo',0,0,NULL,NULL,NULL,NULL,NULL,NULL,'active',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50',NULL);
/*!40000 ALTER TABLE `vt_training_queues` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_user_departments`
--

DROP TABLE IF EXISTS `vt_user_departments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_user_departments` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `department_id` bigint NOT NULL COMMENT '部门ID',
  `position` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '职位',
  `is_primary` tinyint(1) DEFAULT '1' COMMENT '是否主部门',
  `join_date` date DEFAULT NULL COMMENT '入职日期',
  `leave_date` date DEFAULT NULL COMMENT '离职日期',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_department` (`user_id`,`department_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_department_id` (`department_id`),
  KEY `idx_is_primary` (`is_primary`),
  KEY `idx_status` (`status`),
  KEY `idx_join_date` (`join_date`),
  KEY `idx_leave_date` (`leave_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户部门关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_user_departments`
--

LOCK TABLES `vt_user_departments` WRITE;
/*!40000 ALTER TABLE `vt_user_departments` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_user_departments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_user_favorite_relations`
--

DROP TABLE IF EXISTS `vt_user_favorite_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_user_favorite_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `favorite_id` bigint NOT NULL COMMENT '收藏ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(user等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(user等)',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_favorite_entity_relation` (`favorite_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_favorite_id` (`favorite_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户收藏关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_user_favorite_relations`
--

LOCK TABLES `vt_user_favorite_relations` WRITE;
/*!40000 ALTER TABLE `vt_user_favorite_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_user_favorite_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_user_favorites`
--

DROP TABLE IF EXISTS `vt_user_favorites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_user_favorites` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `resource_type` enum('dataset','model','template','dashboard','workspace') COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资源类型',
  `resource_id` bigint NOT NULL COMMENT '资源ID',
  `folder_name` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '收藏夹名称',
  `notes` text COLLATE utf8mb4_unicode_ci COMMENT '备注',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_resource` (`resource_type`,`resource_id`),
  KEY `idx_folder_name` (`folder_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户收藏表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_user_favorites`
--

LOCK TABLES `vt_user_favorites` WRITE;
/*!40000 ALTER TABLE `vt_user_favorites` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_user_favorites` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_user_preference_relations`
--

DROP TABLE IF EXISTS `vt_user_preference_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_user_preference_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `preference_id` bigint NOT NULL COMMENT '偏好设置ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(user等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(user等)',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_preference_entity_relation` (`preference_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_preference_id` (`preference_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户偏好设置关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_user_preference_relations`
--

LOCK TABLES `vt_user_preference_relations` WRITE;
/*!40000 ALTER TABLE `vt_user_preference_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_user_preference_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_user_preferences`
--

DROP TABLE IF EXISTS `vt_user_preferences`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_user_preferences` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `preference_key` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '设置键',
  `preference_value` text COLLATE utf8mb4_unicode_ci COMMENT '设置值',
  `value_type` enum('string','number','boolean','json') COLLATE utf8mb4_unicode_ci DEFAULT 'string' COMMENT '值类型',
  `category` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '设置分类',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_category` (`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户偏好设置表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_user_preferences`
--

LOCK TABLES `vt_user_preferences` WRITE;
/*!40000 ALTER TABLE `vt_user_preferences` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_user_preferences` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_user_roles`
--

DROP TABLE IF EXISTS `vt_user_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_user_roles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `workspace_id` bigint DEFAULT NULL COMMENT '工作空间ID(空表示全局角色)',
  `expires_at` timestamp NULL DEFAULT NULL COMMENT '过期时间',
  `assigned_by` bigint DEFAULT NULL COMMENT '分配人ID',
  `status` enum('active','inactive','expired') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role_workspace` (`user_id`,`role_id`,`workspace_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_workspace_id` (`workspace_id`),
  KEY `idx_expires_at` (`expires_at`),
  KEY `idx_assigned_by` (`assigned_by`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_user_roles`
--

LOCK TABLES `vt_user_roles` WRITE;
/*!40000 ALTER TABLE `vt_user_roles` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_user_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_user_session_relations`
--

DROP TABLE IF EXISTS `vt_user_session_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_user_session_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `session_id` bigint NOT NULL COMMENT '会话ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(user等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(user等)',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_session_entity_relation` (`session_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_session_id` (`session_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户会话关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_user_session_relations`
--

LOCK TABLES `vt_user_session_relations` WRITE;
/*!40000 ALTER TABLE `vt_user_session_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_user_session_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_user_sessions`
--

DROP TABLE IF EXISTS `vt_user_sessions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_user_sessions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `session_id` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '会话ID',
  `device_type` enum('web','mobile','desktop','api') COLLATE utf8mb4_unicode_ci DEFAULT 'web' COMMENT '设备类型',
  `device_info` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '设备信息',
  `ip_address` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'IP地址',
  `user_agent` text COLLATE utf8mb4_unicode_ci COMMENT '用户代理',
  `login_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  `last_activity_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后活动时间',
  `expires_at` timestamp NOT NULL COMMENT '过期时间',
  `status` enum('active','expired','revoked') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_session_id` (`session_id`),
  KEY `idx_device_type` (`device_type`),
  KEY `idx_status` (`status`),
  KEY `idx_expires_at` (`expires_at`),
  KEY `idx_last_activity_at` (`last_activity_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户会话表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_user_sessions`
--

LOCK TABLES `vt_user_sessions` WRITE;
/*!40000 ALTER TABLE `vt_user_sessions` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_user_sessions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_users`
--

DROP TABLE IF EXISTS `vt_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
  `email` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮箱',
  `phone` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '手机号',
  `password_hash` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码哈希',
  `salt` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码盐值',
  `real_name` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '真实姓名',
  `nickname` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '昵称',
  `status` enum('active','inactive','locked','pending') COLLATE utf8mb4_unicode_ci DEFAULT 'pending' COMMENT '用户状态',
  `user_type` enum('admin','user','service') COLLATE utf8mb4_unicode_ci DEFAULT 'user' COMMENT '用户类型',
  `last_login_at` timestamp NULL DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '最后登录IP',
  `password_expires_at` timestamp NULL DEFAULT NULL COMMENT '密码过期时间',
  `login_attempts` int DEFAULT '0' COMMENT '登录尝试次数',
  `locked_until` timestamp NULL DEFAULT NULL COMMENT '锁定到期时间',
  `mfa_enabled` tinyint(1) DEFAULT '0' COMMENT '是否启用双因素认证',
  `mfa_secret` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '双因素认证密钥',
  `email_verified` tinyint(1) DEFAULT '0' COMMENT '邮箱是否验证',
  `phone_verified` tinyint(1) DEFAULT '0' COMMENT '手机是否验证',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`),
  KEY `idx_status` (`status`),
  KEY `idx_user_type` (`user_type`),
  KEY `idx_email_verified` (`email_verified`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_last_login_at` (`last_login_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_users`
--

LOCK TABLES `vt_users` WRITE;
/*!40000 ALTER TABLE `vt_users` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_volcano_job_relations`
--

DROP TABLE IF EXISTS `vt_volcano_job_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_volcano_job_relations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `volcano_job_id` bigint NOT NULL COMMENT 'Volcano作业ID',
  `entity_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实体类型(k8s_resource, training_job等)',
  `entity_id` bigint NOT NULL COMMENT '实体ID',
  `relation_type` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联类型(k8s_resource, training_job等)',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_volcano_entity_relation` (`volcano_job_id`,`entity_type`,`entity_id`,`relation_type`),
  KEY `idx_volcano_job_id` (`volcano_job_id`),
  KEY `idx_entity` (`entity_type`,`entity_id`),
  KEY `idx_relation_type` (`relation_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Volcano作业关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_volcano_job_relations`
--

LOCK TABLES `vt_volcano_job_relations` WRITE;
/*!40000 ALTER TABLE `vt_volcano_job_relations` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_volcano_job_relations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_volcano_jobs`
--

DROP TABLE IF EXISTS `vt_volcano_jobs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_volcano_jobs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `job_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Volcano作业名称',
  `namespace` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'default' COMMENT 'K8s命名空间',
  `uid` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'K8s UID',
  `training_job_id` bigint NOT NULL COMMENT '训练任务ID',
  `queue_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT 'default' COMMENT 'Volcano队列名称',
  `priority` int DEFAULT '0' COMMENT '优先级',
  `min_available` int DEFAULT '1' COMMENT '最小可用实例数',
  `scheduling_policy` json DEFAULT NULL COMMENT '调度策略',
  `plugins` json DEFAULT NULL COMMENT '插件配置',
  `scheduler_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT 'volcano' COMMENT '调度器名称',
  `task_specs` json DEFAULT NULL COMMENT '任务规格列表',
  `volumes` json DEFAULT NULL COMMENT '卷配置',
  `ttl_seconds_after_finished` int DEFAULT NULL COMMENT '完成后保留时间(秒)',
  `active_deadline_seconds` int DEFAULT NULL COMMENT '活跃截止时间(秒)',
  `backoff_limit` int DEFAULT '3' COMMENT '重试次数限制',
  `phase` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '作业阶段：Pending, Running, Completed, Failed等',
  `conditions` json DEFAULT NULL COMMENT '作业条件',
  `status_message` text COLLATE utf8mb4_unicode_ci COMMENT '状态消息',
  `start_time` timestamp NULL DEFAULT NULL COMMENT '开始时间',
  `completion_time` timestamp NULL DEFAULT NULL COMMENT '完成时间',
  `running_duration` int DEFAULT NULL COMMENT '运行时长(秒)',
  `total_cpu_request` decimal(10,2) DEFAULT NULL COMMENT '总CPU请求',
  `total_memory_request` bigint DEFAULT NULL COMMENT '总内存请求(bytes)',
  `total_gpu_request` int DEFAULT NULL COMMENT '总GPU请求',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_namespace_name` (`namespace`,`job_name`),
  UNIQUE KEY `uk_uid` (`uid`),
  KEY `idx_training_job_id` (`training_job_id`),
  KEY `idx_queue_name` (`queue_name`),
  KEY `idx_phase` (`phase`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_scheduler_name` (`scheduler_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Volcano作业详情表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_volcano_jobs`
--

LOCK TABLES `vt_volcano_jobs` WRITE;
/*!40000 ALTER TABLE `vt_volcano_jobs` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_volcano_jobs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_volcano_queues`
--

DROP TABLE IF EXISTS `vt_volcano_queues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_volcano_queues` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `queue_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '队列名称',
  `namespace` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'volcano-system' COMMENT '命名空间',
  `weight` int DEFAULT '1' COMMENT '权重',
  `capability` json DEFAULT NULL COMMENT '队列容量配置',
  `cpu_quota` decimal(10,2) DEFAULT NULL COMMENT 'CPU配额',
  `memory_quota` bigint DEFAULT NULL COMMENT '内存配额(bytes)',
  `gpu_quota` int DEFAULT NULL COMMENT 'GPU配额',
  `state` enum('Open','Closed') COLLATE utf8mb4_unicode_ci DEFAULT 'Open' COMMENT '队列状态',
  `running_jobs` int DEFAULT '0' COMMENT '运行中的作业数',
  `pending_jobs` int DEFAULT '0' COMMENT '等待中的作业数',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_queue_name` (`queue_name`,`namespace`),
  KEY `idx_state` (`state`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Volcano队列配置表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_volcano_queues`
--

LOCK TABLES `vt_volcano_queues` WRITE;
/*!40000 ALTER TABLE `vt_volcano_queues` DISABLE KEYS */;
INSERT INTO `vt_volcano_queues` VALUES (1,'default','volcano-system',1,NULL,NULL,NULL,NULL,'Open',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(2,'high-priority','volcano-system',5,NULL,NULL,NULL,NULL,'Open',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50'),(3,'gpu-intensive','volcano-system',3,NULL,NULL,NULL,NULL,'Open',0,0,'2025-06-25 04:42:50','2025-06-25 04:42:50');
/*!40000 ALTER TABLE `vt_volcano_queues` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_workspace_members`
--

DROP TABLE IF EXISTS `vt_workspace_members`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_workspace_members` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `workspace_id` bigint NOT NULL COMMENT '工作空间ID',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `role` enum('owner','admin','member','viewer','guest') COLLATE utf8mb4_unicode_ci DEFAULT 'member' COMMENT '成员角色',
  `invited_by` bigint DEFAULT NULL COMMENT '邀请人ID',
  `invited_at` timestamp NULL DEFAULT NULL COMMENT '邀请时间',
  `invitation_token` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '邀请令牌',
  `invitation_message` text COLLATE utf8mb4_unicode_ci COMMENT '邀请消息',
  `invitation_expires_at` timestamp NULL DEFAULT NULL COMMENT '邀请过期时间',
  `joined_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
  `last_activity_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后活动时间',
  `status` enum('active','inactive','pending','rejected','removed') COLLATE utf8mb4_unicode_ci DEFAULT 'pending' COMMENT '成员状态',
  `custom_permissions` json DEFAULT NULL COMMENT '自定义权限',
  `resource_limits` json DEFAULT NULL COMMENT '资源限制',
  `notification_settings` json DEFAULT NULL COMMENT '通知设置',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_workspace_user` (`workspace_id`,`user_id`),
  KEY `idx_workspace_id` (`workspace_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role` (`role`),
  KEY `idx_status` (`status`),
  KEY `idx_invited_by` (`invited_by`),
  KEY `idx_invitation_token` (`invitation_token`),
  KEY `idx_invitation_expires_at` (`invitation_expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工作空间成员表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_workspace_members`
--

LOCK TABLES `vt_workspace_members` WRITE;
/*!40000 ALTER TABLE `vt_workspace_members` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_workspace_members` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_workspace_owners`
--

DROP TABLE IF EXISTS `vt_workspace_owners`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_workspace_owners` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `workspace_id` bigint NOT NULL COMMENT '工作空间ID',
  `user_id` bigint NOT NULL COMMENT '所有者用户ID',
  `owner_type` enum('primary','co_owner') COLLATE utf8mb4_unicode_ci DEFAULT 'primary' COMMENT '所有者类型',
  `start_date` date DEFAULT NULL COMMENT '开始日期',
  `end_date` date DEFAULT NULL COMMENT '结束日期',
  `status` enum('active','inactive') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_workspace_user_type` (`workspace_id`,`user_id`,`owner_type`),
  KEY `idx_workspace_id` (`workspace_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_owner_type` (`owner_type`),
  KEY `idx_status` (`status`),
  KEY `idx_workspace_owner_composite` (`workspace_id`,`user_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工作空间所有者关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_workspace_owners`
--

LOCK TABLES `vt_workspace_owners` WRITE;
/*!40000 ALTER TABLE `vt_workspace_owners` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_workspace_owners` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_workspace_projects`
--

DROP TABLE IF EXISTS `vt_workspace_projects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_workspace_projects` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `workspace_id` bigint NOT NULL COMMENT '工作空间ID',
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '项目名称',
  `display_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '显示名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '项目描述',
  `project_type` enum('ml_training','data_analysis','model_serving','research','experiment','custom') COLLATE utf8mb4_unicode_ci DEFAULT 'ml_training' COMMENT '项目类型',
  `repository_url` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '代码仓库URL',
  `storage_used_gb` decimal(10,2) DEFAULT '0.00' COMMENT '已用存储(GB)',
  `compute_used_hours` decimal(10,2) DEFAULT '0.00' COMMENT '已用计算(小时)',
  `gpu_used_hours` decimal(10,2) DEFAULT '0.00' COMMENT '已用GPU(小时)',
  `status` enum('active','inactive','completed','archived') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '项目状态',
  `visibility` enum('public','private','workspace') COLLATE utf8mb4_unicode_ci DEFAULT 'workspace' COMMENT '可见性',
  `priority` enum('low','normal','high','urgent') COLLATE utf8mb4_unicode_ci DEFAULT 'normal' COMMENT '优先级',
  `started_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
  `planned_end_at` timestamp NULL DEFAULT NULL COMMENT '计划结束时间',
  `completed_at` timestamp NULL DEFAULT NULL COMMENT '完成时间',
  `last_activity_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后活动时间',
  `tags` json DEFAULT NULL COMMENT '标签',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `settings` json DEFAULT NULL COMMENT '项目设置',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_workspace_project_name` (`workspace_id`,`name`),
  KEY `idx_workspace_id` (`workspace_id`),
  KEY `idx_project_type` (`project_type`),
  KEY `idx_status` (`status`),
  KEY `idx_visibility` (`visibility`),
  KEY `idx_priority` (`priority`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工作空间项目表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_workspace_projects`
--

LOCK TABLES `vt_workspace_projects` WRITE;
/*!40000 ALTER TABLE `vt_workspace_projects` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_workspace_projects` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vt_workspaces`
--

DROP TABLE IF EXISTS `vt_workspaces`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `vt_workspaces` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '工作空间名称',
  `display_name` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '显示名称',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '工作空间描述',
  `workspace_type` enum('personal','team','organization','project') COLLATE utf8mb4_unicode_ci DEFAULT 'team' COMMENT '工作空间类型',
  `visibility` enum('public','private','internal') COLLATE utf8mb4_unicode_ci DEFAULT 'private' COMMENT '可见性',
  `storage_quota_gb` decimal(10,2) DEFAULT '100.00' COMMENT '存储配额(GB)',
  `compute_quota_hours` decimal(10,2) DEFAULT '1000.00' COMMENT '计算配额(小时)',
  `gpu_quota_hours` decimal(10,2) DEFAULT '100.00' COMMENT 'GPU配额(小时)',
  `dataset_quota_count` int DEFAULT '100' COMMENT '数据集配额数量',
  `model_quota_count` int DEFAULT '50' COMMENT '模型配额数量',
  `member_quota_count` int DEFAULT '20' COMMENT '成员配额数量',
  `storage_used_gb` decimal(10,2) DEFAULT '0.00' COMMENT '已用存储(GB)',
  `compute_used_hours` decimal(10,2) DEFAULT '0.00' COMMENT '已用计算(小时)',
  `gpu_used_hours` decimal(10,2) DEFAULT '0.00' COMMENT '已用GPU(小时)',
  `dataset_count` int DEFAULT '0' COMMENT '数据集数量',
  `model_count` int DEFAULT '0' COMMENT '模型数量',
  `member_count` int DEFAULT '0' COMMENT '成员数量',
  `status` enum('active','inactive','suspended','archived') COLLATE utf8mb4_unicode_ci DEFAULT 'active' COMMENT '状态',
  `features_enabled` json DEFAULT NULL COMMENT '启用的功能列表',
  `settings` json DEFAULT NULL COMMENT '工作空间设置',
  `default_permissions` json DEFAULT NULL COMMENT '默认权限配置',
  `billing_type` enum('free','subscription','pay_as_go') COLLATE utf8mb4_unicode_ci DEFAULT 'free' COMMENT '计费类型',
  `billing_config` json DEFAULT NULL COMMENT '计费配置',
  `tags` json DEFAULT NULL COMMENT '标签',
  `metadata` json DEFAULT NULL COMMENT '元数据',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
  `last_activity_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后活动时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_workspace_type` (`workspace_type`),
  KEY `idx_visibility` (`visibility`),
  KEY `idx_status` (`status`),
  KEY `idx_billing_type` (`billing_type`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_last_activity_at` (`last_activity_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='工作空间表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vt_workspaces`
--

LOCK TABLES `vt_workspaces` WRITE;
/*!40000 ALTER TABLE `vt_workspaces` DISABLE KEYS */;
/*!40000 ALTER TABLE `vt_workspaces` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'volctraindb'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-08-03 21:37:37
