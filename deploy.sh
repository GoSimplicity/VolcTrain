#!/bin/bash
# VolcTrain 快速部署脚本
# 使用方法: ./deploy.sh [environment]
# 环境选项: dev, staging, production

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查必要工具
check_requirements() {
    log_info "检查系统要求..."
    
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    # 检查Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi
    
    # 检查Docker是否运行
    if ! docker info &> /dev/null; then
        log_error "Docker 服务未运行，请启动 Docker"
        exit 1
    fi
    
    log_success "系统要求检查通过"
}

# 创建环境变量文件
setup_environment() {
    local env=${1:-production}
    log_info "设置 ${env} 环境..."
    
    if [ ! -f .env ]; then
        if [ -f .env.example ]; then
            cp .env.example .env
            log_info "已从 .env.example 创建 .env 文件"
            log_warning "请编辑 .env 文件配置你的环境变量"
            
            # 生成随机密码和密钥
            generate_secrets
        else
            log_error ".env.example 文件不存在"
            exit 1
        fi
    else
        log_info ".env 文件已存在，跳过创建"
    fi
}

# 生成安全密钥
generate_secrets() {
    log_info "生成安全密钥..."
    
    # 生成随机密码
    MYSQL_ROOT_PASSWORD=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-25)
    MYSQL_PASSWORD=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-25)
    REDIS_PASSWORD=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-25)
    JWT_ACCESS_SECRET=$(openssl rand -base64 64 | tr -d "=+/" | cut -c1-64)
    GRAFANA_PASSWORD=$(openssl rand -base64 16 | tr -d "=+/" | cut -c1-16)
    DATA_ENCRYPTION_KEY=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-32)
    
    # 更新 .env 文件
    sed -i.bak "s/your_secure_root_password_here/${MYSQL_ROOT_PASSWORD}/g" .env
    sed -i.bak "s/your_secure_app_password_here/${MYSQL_PASSWORD}/g" .env
    sed -i.bak "s/your_secure_redis_password_here/${REDIS_PASSWORD}/g" .env
    sed -i.bak "s/your_jwt_access_secret_here_64_chars_min/${JWT_ACCESS_SECRET}/g" .env
    sed -i.bak "s/your_grafana_admin_password_here/${GRAFANA_PASSWORD}/g" .env
    sed -i.bak "s/your_32_character_encryption_key_here/${DATA_ENCRYPTION_KEY}/g" .env
    
    rm .env.bak
    log_success "安全密钥生成完成"
}

# 清理旧容器和镜像
cleanup() {
    log_info "清理旧的容器和镜像..."
    
    # 停止并删除容器
    docker-compose down --remove-orphans || true
    
    # 删除相关镜像（可选）
    if [ "${1}" == "--clean-images" ]; then
        docker images | grep volctrain | awk '{print $3}' | xargs -r docker rmi || true
        log_info "清理旧镜像完成"
    fi
    
    log_success "清理完成"
}

# 构建服务
build_services() {
    log_info "构建服务镜像..."
    
    # 构建后端服务
    log_info "构建后端服务..."
    docker-compose build backend-api backend-common backend-monitoring backend-training
    
    # 构建前端服务
    log_info "构建前端服务..."
    docker-compose build frontend
    
    log_success "服务构建完成"
}

# 初始化数据库
init_database() {
    log_info "初始化数据库..."
    
    # 启动数据库服务
    docker-compose up -d mysql redis
    
    # 等待数据库启动
    log_info "等待数据库启动..."
    sleep 30
    
    # 检查数据库连接
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if docker-compose exec -T mysql mysql -u root -p"${MYSQL_ROOT_PASSWORD:-Abc@1234}" -e "SELECT 1" > /dev/null 2>&1; then
            log_success "数据库连接成功"
            break
        fi
        
        log_info "尝试连接数据库 ($attempt/$max_attempts)..."
        attempt=$((attempt + 1))
        sleep 2
    done
    
    if [ $attempt -gt $max_attempts ]; then
        log_error "数据库连接失败"
        exit 1
    fi
    
    log_success "数据库初始化完成"
}

# 启动服务
start_services() {
    log_info "启动所有服务..."
    
    # 使用 docker-compose 启动所有服务
    docker-compose up -d
    
    log_success "服务启动完成"
}

# 检查服务健康状态
check_health() {
    log_info "检查服务健康状态..."
    
    local services=("backend-api:8888" "backend-common:8889" "backend-monitoring:8890" "backend-training:8891" "frontend:80")
    local max_attempts=30
    local healthy_count=0
    
    for service_endpoint in "${services[@]}"; do
        local service_name=$(echo $service_endpoint | cut -d: -f1)
        local port=$(echo $service_endpoint | cut -d: -f2)
        local attempt=1
        
        log_info "检查 ${service_name} 服务..."
        
        while [ $attempt -le $max_attempts ]; do
            if curl -f -s http://localhost:${port}/health > /dev/null 2>&1; then
                log_success "${service_name} 服务健康"
                healthy_count=$((healthy_count + 1))
                break
            fi
            
            sleep 2
            attempt=$((attempt + 1))
        done
        
        if [ $attempt -gt $max_attempts ]; then
            log_warning "${service_name} 服务健康检查超时"
        fi
    done
    
    log_info "健康服务数量: ${healthy_count}/${#services[@]}"
}

# 显示部署信息
show_deployment_info() {
    log_success "VolcTrain 部署完成！"
    echo
    echo "========================================"
    echo "服务访问地址:"
    echo "========================================"
    echo "🌐 前端应用: http://localhost:80"
    echo "🔧 API服务: http://localhost:8888"
    echo "📊 监控面板: http://localhost:3000"
    echo "📈 Prometheus: http://localhost:9091"
    echo "========================================"
    echo "默认登录信息:"
    echo "========================================"
    echo "Grafana:"
    echo "  用户名: admin"
    echo "  密码: ${GRAFANA_PASSWORD:-admin123}"
    echo "========================================"
    echo
    log_info "使用 'docker-compose logs -f' 查看日志"
    log_info "使用 'docker-compose ps' 查看服务状态"
    log_info "使用 './deploy.sh stop' 停止服务"
}

# 停止服务
stop_services() {
    log_info "停止所有服务..."
    docker-compose down
    log_success "服务已停止"
}

# 重启服务
restart_services() {
    log_info "重启服务..."
    docker-compose restart
    log_success "服务重启完成"
}

# 显示帮助信息
show_help() {
    echo "VolcTrain 部署脚本"
    echo
    echo "使用方法:"
    echo "  $0 [command] [options]"
    echo
    echo "命令:"
    echo "  start     启动服务 (默认)"
    echo "  stop      停止服务"
    echo "  restart   重启服务"
    echo "  build     重新构建镜像"
    echo "  clean     清理容器和镜像"
    echo "  logs      查看日志"
    echo "  status    查看服务状态"
    echo "  help      显示帮助信息"
    echo
    echo "选项:"
    echo "  --clean-images  清理时同时删除镜像"
    echo "  --no-cache      构建时不使用缓存"
    echo
    echo "示例:"
    echo "  $0 start              # 启动服务"
    echo "  $0 build --no-cache   # 重新构建（不使用缓存）"
    echo "  $0 clean --clean-images  # 清理容器和镜像"
}

# 主函数
main() {
    local command=${1:-start}
    local option=${2}
    
    case $command in
        start)
            check_requirements
            setup_environment
            build_services
            init_database
            start_services
            sleep 10
            check_health
            show_deployment_info
            ;;
        stop)
            stop_services
            ;;
        restart)
            restart_services
            ;;
        build)
            check_requirements
            if [ "$option" == "--no-cache" ]; then
                docker-compose build --no-cache
            else
                build_services
            fi
            ;;
        clean)
            cleanup "$option"
            ;;
        logs)
            docker-compose logs -f
            ;;
        status)
            docker-compose ps
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "未知命令: $command"
            show_help
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"