#!/bin/sh

# VolcTrain服务启动脚本
set -e

# 检查环境变量
check_env() {
    if [ -z "$DB_HOST" ]; then
        echo "错误: 环境变量 DB_HOST 未设置"
        exit 1
    fi
    
    if [ -z "$DB_NAME" ]; then
        echo "错误: 环境变量 DB_NAME 未设置"
        exit 1
    fi
    
    if [ -z "$DB_USER" ]; then
        echo "错误: 环境变量 DB_USER 未设置"
        exit 1
    fi
    
    if [ -z "$DB_PASSWORD" ]; then
        echo "错误: 环境变量 DB_PASSWORD 未设置"
        exit 1
    fi
}

# 等待数据库连接
wait_for_db() {
    echo "等待数据库连接..."
    max_attempts=30
    attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if mysql -h"$DB_HOST" -P"${DB_PORT:-3306}" -u"$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1" > /dev/null 2>&1; then
            echo "数据库连接成功"
            return 0
        fi
        
        echo "尝试连接数据库 ($attempt/$max_attempts)..."
        attempt=$((attempt + 1))
        sleep 2
    done
    
    echo "数据库连接失败，请检查配置"
    exit 1
}

# 初始化数据库
init_database() {
    echo "检查数据库初始化状态..."
    
    # 检查数据库是否存在表
    table_count=$(mysql -h"$DB_HOST" -P"${DB_PORT:-3306}" -u"$DB_USER" -p"$DB_PASSWORD" -D"$DB_NAME" -e "SHOW TABLES;" 2>/dev/null | wc -l)
    
    if [ "$table_count" -lt 10 ]; then
        echo "初始化数据库结构..."
        
        # 执行SQL脚本
        for sql_file in sql/*.sql; do
            if [ -f "$sql_file" ]; then
                echo "执行 $sql_file"
                mysql -h"$DB_HOST" -P"${DB_PORT:-3306}" -u"$DB_USER" -p"$DB_PASSWORD" -D"$DB_NAME" < "$sql_file"
            fi
        done
        
        echo "数据库初始化完成"
    else
        echo "数据库已初始化，跳过"
    fi
}

# 主要启动逻辑
main() {
    echo "启动 VolcTrain 服务..."
    
    # 设置默认环境变量
    export SERVICE_NAME="${SERVICE_NAME:-volctrain}"
    export DB_PORT="${DB_PORT:-3306}"
    export REDIS_HOST="${REDIS_HOST:-redis}"
    export REDIS_PORT="${REDIS_PORT:-6379}"
    export LOG_LEVEL="${LOG_LEVEL:-info}"
    
    # 检查环境变量
    check_env
    
    # 等待数据库
    wait_for_db
    
    # 初始化数据库
    init_database
    
    # 根据传入的参数决定启动哪个服务
    service_type="${1:-api}"
    
    case "$service_type" in
        api)
            echo "启动 API 服务..."
            # 使用环境变量来决定配置文件
            if [ "$DEPLOY_ENV" = "production" ]; then
                exec ./bin/api -f etc/config-production.yaml
            else
                exec ./bin/api -f etc/config-dev.yaml
            fi
            ;;
        common)
            echo "启动 Common 服务..."
            # 使用环境变量来决定配置文件
            if [ "$DEPLOY_ENV" = "production" ]; then
                exec ./bin/common -f etc/config-production.yaml
            else
                exec ./bin/common -f etc/config-dev.yaml
            fi
            ;;
        monitoring)
            echo "启动 Monitoring 服务..."
            exec ./bin/monitoring -f etc/monitoring.yaml
            ;;
        training)
            echo "启动 Training 服务..."
            exec ./bin/training -f etc/training_service.yaml
            ;;
        *)
            echo "未知的服务类型: $service_type"
            echo "支持的服务类型: api, common, monitoring, training"
            exit 1
            ;;
    esac
}

# 捕获信号以优雅关闭
trap 'echo "接收到终止信号，正在关闭服务..."; exit 0' TERM INT

# 执行主函数
main "$@"