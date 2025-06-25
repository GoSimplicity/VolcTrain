#!/bin/bash

# =====================================
# AI-GPU机器学习平台数据库一键初始化脚本
# =====================================
# 作者: Bamboo
# 版本: 1.0.0
# 兼容: MySQL 8.0+
# 功能: 自动初始化volctraindb数据库和所有表结构
# =====================================

set -e  # 遇到错误时退出

# 颜色输出函数
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 配置参数
DB_HOST="localhost"
DB_PORT="3306"
DB_USER="root"
DB_PASSWORD="root"
DB_NAME="volctraindb"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 显示脚本信息
print_info "=========================================="
print_info "AI-GPU机器学习平台数据库初始化脚本"
print_info "=========================================="
print_info "数据库主机: $DB_HOST:$DB_PORT"
print_info "数据库用户: $DB_USER"
print_info "目标数据库: $DB_NAME"
print_info "脚本目录: $SCRIPT_DIR"
print_info "=========================================="

# 检查MySQL是否可连接
print_info "检查MySQL连接..."
if ! mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1;" &>/dev/null; then
    print_error "无法连接到MySQL服务器"
    print_error "请检查："
    print_error "1. MySQL服务是否启动"
    print_error "2. 连接参数是否正确(host=$DB_HOST, port=$DB_PORT, user=$DB_USER)"
    print_error "3. 密码是否正确"
    exit 1
fi
print_success "MySQL连接成功"

# 检查MySQL版本
print_info "检查MySQL版本..."
MYSQL_VERSION=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -sN -e "SELECT VERSION();" 2>/dev/null)
print_info "MySQL版本: $MYSQL_VERSION"

# 检查版本是否支持(MySQL 8.0+)
if [[ "$MYSQL_VERSION" < "8.0" ]]; then
    print_warning "检测到MySQL版本低于8.0，可能存在兼容性问题"
    print_warning "建议使用MySQL 8.0或更高版本"
fi

# 检查必要的SQL文件是否存在
print_info "检查SQL文件..."
SQL_FILES=(
    "00_create_schema.sql"
    "01_users_permissions.sql"
    "02_files.sql"
    "03_workspaces.sql"
    "04_datasets.sql"
    "05_gpu_clusters.sql"
    "06_training_jobs.sql"
    "07_models.sql"
    "08_system_support.sql"
    "09_monitoring.sql"
    "10_k8s_relations.sql"
)

for file in "${SQL_FILES[@]}"; do
    if [[ ! -f "$SCRIPT_DIR/$file" ]]; then
        print_error "SQL文件不存在: $file"
        exit 1
    fi
done
print_success "所有SQL文件检查通过"

# 确认是否继续
echo ""
print_warning "======== 重要提醒 ========"
print_warning "此操作将会:"
print_warning "1. 如果数据库 $DB_NAME 已存在，将直接使用"
print_warning "2. 删除并重新创建所有表结构"
print_warning "3. 插入初始化数据"
print_warning "4. 这会清空现有数据!"
print_warning "=========================="
echo ""

read -p "确认继续初始化吗? (输入 'yes' 确认): " confirm
if [[ "$confirm" != "yes" ]]; then
    print_info "已取消初始化操作"
    exit 0
fi

print_info "开始初始化数据库..."

# 备份现有数据库(如果存在)
print_info "检查数据库是否存在..."
if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "USE $DB_NAME;" &>/dev/null; then
    print_warning "数据库 $DB_NAME 已存在"
    BACKUP_FILE="volctraindb_backup_$(date +%Y%m%d_%H%M%S).sql"
    print_info "创建备份文件: $BACKUP_FILE"
    
    if mysqldump -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" \
        --single-transaction --routines --triggers "$DB_NAME" > "$SCRIPT_DIR/$BACKUP_FILE" 2>/dev/null; then
        print_success "数据库备份成功: $BACKUP_FILE"
    else
        print_warning "数据库备份失败，但继续执行初始化"
    fi
else
    print_info "数据库 $DB_NAME 不存在，将创建新数据库"
fi

# 执行初始化脚本
print_info "执行数据库初始化..."

# 使用主初始化脚本
if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" < "$SCRIPT_DIR/00_create_schema.sql"; then
    print_success "数据库初始化完成!"
else
    print_error "数据库初始化失败"
    exit 1
fi

# 验证初始化结果
print_info "验证初始化结果..."

# 检查数据库是否存在
if ! mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "USE $DB_NAME;" &>/dev/null; then
    print_error "数据库 $DB_NAME 验证失败"
    exit 1
fi

# 统计表数量
TABLE_COUNT=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -sN \
    -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = '$DB_NAME';" 2>/dev/null)

if [[ "$TABLE_COUNT" -gt 0 ]]; then
    print_success "验证通过 - 共创建 $TABLE_COUNT 个表"
else
    print_error "验证失败 - 未发现任何表"
    exit 1
fi

# 检查关键表是否存在
KEY_TABLES=("vt_users" "vt_workspaces" "vt_training_jobs" "vt_models" "vt_datasets")
for table in "${KEY_TABLES[@]}"; do
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" \
        -e "SELECT 1 FROM $DB_NAME.$table LIMIT 1;" &>/dev/null; then
        print_success "表 $table 创建成功"
    else
        print_error "表 $table 创建失败或无法访问"
        exit 1
    fi
done

# 显示数据库信息
print_info "=========================================="
print_info "数据库初始化完成!"
print_info "=========================================="

# 获取并显示数据库统计信息
print_info "数据库统计信息:"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" \
    -e "SELECT 
        COUNT(*) as '总表数',
        (SELECT COUNT(*) FROM information_schema.tables 
         WHERE table_schema = '$DB_NAME' AND table_name LIKE '%_relations') as '关联表数',
        (SELECT COUNT(*) FROM information_schema.tables 
         WHERE table_schema = '$DB_NAME' AND table_name NOT LIKE '%_relations') as '业务表数'
    FROM information_schema.tables 
    WHERE table_schema = '$DB_NAME';" 2>/dev/null

# 显示用户权限信息
print_info "用户权限信息:"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" \
    -e "SHOW GRANTS FOR 'volctrain_app'@'%';" 2>/dev/null || print_warning "volctrain_app用户权限查询失败"

mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" \
    -e "SHOW GRANTS FOR 'volctrain_readonly'@'%';" 2>/dev/null || print_warning "volctrain_readonly用户权限查询失败"

# 连接信息
print_info "=========================================="
print_info "数据库连接信息:"
print_info "主机: $DB_HOST:$DB_PORT"
print_info "数据库: $DB_NAME"
print_info "应用用户: volctrain_app (密码: Abc@1234)"
print_info "只读用户: volctrain_readonly (密码: volctrain_readonly_password_2024)"
print_info "=========================================="

# 下一步建议
print_info "建议下一步操作:"
print_info "1. 测试应用连接数据库"
print_info "2. 根据需要调整用户密码"
print_info "3. 配置应用的数据库连接参数"
print_info "4. 如需要，可以导入测试数据"

print_success "初始化脚本执行完成!"

# 脚本执行完成
exit 0