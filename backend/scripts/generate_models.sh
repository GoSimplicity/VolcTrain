#!/bin/bash

# =====================================
# VolcTrain数据库模型层生成脚本  
# =====================================
# 使用goctl工具从SQL脚本生成Go数据模型
# =====================================

set -e

# 颜色输出
GREEN='\033[0;32m'
BLUE='\033[0;34m' 
RED='\033[0;31m'
NC='\033[0m'

print_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 检查goctl是否安装
if ! command -v goctl &> /dev/null; then
    print_error "goctl未安装，请先安装: go install github.com/zeromicro/go-zero/tools/goctl@latest"
    exit 1
fi

print_info "开始生成数据库模型层..."

# 数据库连接配置
DB_USER="root"
DB_PASSWORD=""
DB_HOST="127.0.0.1" 
DB_PORT="3306"
DB_NAME="vt_volctrain"
DSN="${DB_USER}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"

# 定义要生成模型的表列表
TABLES=(
    # 用户权限模块
    "vt_users"
    "vt_departments" 
    "vt_permissions"
    "vt_roles"
    "vt_user_roles"
    "vt_user_departments"
    "vt_department_managers"
    
    # 文件管理模块
    "vt_files"
    "vt_file_versions"
    "vt_file_downloads"
    
    # 工作空间模块
    "vt_workspaces"
    "vt_workspace_members"
    "vt_workspace_projects"
    "vt_workspace_resources"
    "vt_workspace_quotas"
    
    # 数据集模块
    "vt_datasets"
    "vt_dataset_versions"
    "vt_dataset_files"
    "vt_dataset_annotations"
    "vt_dataset_labels"
    "vt_dataset_splits"
    
    # GPU集群模块
    "vt_gpu_clusters"
    "vt_gpu_nodes"
    "vt_gpu_devices"
    "vt_gpu_allocations"
    "vt_gpu_usage_records"
    
    # 训练调度模块
    "vt_training_queues"
    "vt_training_jobs"
    "vt_training_job_instances"
    "vt_training_metrics"
    "vt_training_logs"
    "vt_training_checkpoints"
    "vt_training_job_relations"
    "vt_training_job_events"
    
    # 模型管理模块
    "vt_models"
    "vt_model_versions"
    "vt_model_deployments"
    "vt_model_inference_logs"
    "vt_model_files"
    "vt_model_metrics"
    
    # 系统支持模块
    "vt_system_configs"
    "vt_audit_logs"
    "vt_operation_logs"
    "vt_notification_channels"
    
    # 监控告警模块
    "vt_monitor_metrics"
    "vt_monitor_data"
    "vt_alert_rules"
    "vt_alert_records"
    "vt_alert_subscriptions"
    "vt_notification_logs"
    "vt_dashboard_configs"
    
    # K8s关联模块
    "vt_k8s_clusters"
    "vt_k8s_namespaces"
    "vt_k8s_resources"
)

# 清理现有model目录
rm -rf model/*

print_info "生成模型文件到 model/ 目录..."

# 批量生成模型
for table in "${TABLES[@]}"; do
    print_info "正在生成表 ${table} 的模型..."
    
    # 生成模型文件
    goctl model mysql datasource -url="${DSN}" -table="${table}" -dir="./model" -cache=false --style=goZero
    
    if [ $? -eq 0 ]; then
        print_success "表 ${table} 模型生成成功"
    else
        print_error "表 ${table} 模型生成失败"
    fi
done

print_success "所有数据库模型生成完成！"
print_info "生成的模型文件位于: ./model/"
print_info "可以通过 ls -la model/ 查看生成的文件"

echo ""
print_info "生成的模型文件统计:"
find model/ -name "*.go" | wc -l | xargs echo "Go文件数量:"
find model/ -name "*model.go" | wc -l | xargs echo "模型文件数量:"