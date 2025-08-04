#!/bin/bash

# VolcTrain 一键启动脚本
# 使用方法: ./start.sh [start|stop|restart|status|logs]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目名称
PROJECT_NAME="volctrain"

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# 检查Docker是否安装并运行
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi

    if ! docker info &> /dev/null; then
        log_error "Docker 未运行，请启动 Docker"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi

    log_info "Docker 环境检查通过"
}

# 检查环境文件
check_env_file() {
    if [ ! -f ".env" ]; then
        log_warn ".env 文件不存在，从 .env.example 复制"
        if [ -f ".env.example" ]; then
            cp .env.example .env
            log_info "已创建 .env 文件，请根据实际环境修改配置"
        else
            log_error ".env.example 文件不存在"
            exit 1
        fi
    fi
    log_info "环境配置文件检查通过"
}

# 检查必要的目录和文件
check_directories() {
    local files=(
        "backend/sql"
        "backend/Dockerfile"
        "backend/etc/config.yaml"
        "docker-compose.yml"
        "docker-compose.dev.yml"
        "deploy/prometheus.yml"
        "deploy/redis.conf"
    )
    
    for file in "${files[@]}"; do
        if [ ! -e "$file" ]; then
            log_error "必要的目录或文件不存在: $file"
            exit 1
        fi
    done
    
    log_info "项目目录结构检查通过"
}

# 构建镜像
build_images() {
    local mode=${1:-prod}
    log_info "开始构建 Docker 镜像 ($mode 模式)..."
    
    if [ "$mode" = "dev" ]; then
        log_info "构建开发环境镜像..."
        docker-compose -f docker-compose.dev.yml build --no-cache backend-api
    else
        log_info "构建生产环境镜像..."
        docker-compose build --no-cache backend-api
    fi
    
    log_info "Docker 镜像构建完成"
}

# 启动服务
start_services() {
    local mode=${1:-prod}
    log_info "启动 VolcTrain 服务 ($mode 模式)..."
    
    if [ "$mode" = "dev" ]; then
        # 启动开发环境
        log_info "启动开发环境..."
        docker-compose -f docker-compose.dev.yml down >/dev/null 2>&1 || true
        docker-compose -f docker-compose.dev.yml up -d
        
        # 等待服务启动
        log_info "等待服务启动..."
        sleep 15
        
        # 健康检查
        check_service_health "dev"
        
        log_info "开发环境启动完成"
        show_urls "dev"
    else
        # 启动生产环境
        log_info "启动生产环境..."
        docker-compose down >/dev/null 2>&1 || true
        docker-compose up -d
        
        # 等待服务启动
        log_info "等待服务启动..."
        sleep 20
        
        # 健康检查
        check_service_health "prod"
        
        log_info "生产环境启动完成"
        show_urls "prod"
    fi
}

# 健康检查
check_service_health() {
    local mode=${1:-prod}
    log_info "执行健康检查..."
    
    # 检查数据库连接
    check_database_connection "$mode"
    
    # 检查Redis连接
    check_redis_connection "$mode"
    
    # 检查API服务
    check_api_service
    
    log_info "所有服务健康检查完成"
}

# 检查数据库连接
check_database_connection() {
    local mode=${1:-prod}
    local container_name="volctrain-mysql"
    if [ "$mode" = "dev" ]; then
        container_name="volctrain-mysql-dev"
    fi
    
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if docker exec "$container_name" mysqladmin ping -h localhost -uroot -proot &> /dev/null; then
            log_info "✓ MySQL数据库连接正常"
            return 0
        fi
        
        log_debug "尝试连接数据库 ($attempt/$max_attempts)..."
        attempt=$((attempt + 1))
        sleep 2
    done
    
    log_error "✗ MySQL数据库连接失败，请检查配置"
    return 1
}

# 检查Redis连接
check_redis_connection() {
    local mode=${1:-prod}
    local container_name="volctrain-redis"
    if [ "$mode" = "dev" ]; then
        container_name="volctrain-redis-dev"
    fi
    
    if docker exec "$container_name" redis-cli ping &> /dev/null; then
        log_info "✓ Redis缓存连接正常"
        return 0
    else
        log_error "✗ Redis缓存连接失败"
        return 1
    fi
}

# 检查API服务
check_api_service() {
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s http://localhost:8888/health >/dev/null 2>&1; then
            log_info "✓ API服务健康检查通过"
            return 0
        fi
        
        if [ $attempt -eq $max_attempts ]; then
            log_error "✗ API服务健康检查失败"
            show_service_logs
            return 1
        fi
        
        log_debug "等待API服务启动... ($attempt/$max_attempts)"
        sleep 2
        ((attempt++))
    done
}

# 停止服务
stop_services() {
    local mode=${1:-prod}
    log_info "停止 VolcTrain 服务 ($mode 模式)..."
    
    if [ "$mode" = "dev" ]; then
        docker-compose -f docker-compose.dev.yml down
    else
        docker-compose down
    fi
    
    log_info "所有服务已停止"
}

# 重启服务
restart_services() {
    local mode=${1:-prod}
    log_info "重启 VolcTrain 服务 ($mode 模式)..."
    stop_services "$mode"
    sleep 5
    start_services "$mode"
}

# 显示服务状态
show_status() {
    local mode=${1:-prod}
    log_info "服务状态 ($mode 模式):"
    
    if [ "$mode" = "dev" ]; then
        docker-compose -f docker-compose.dev.yml ps
    else
        docker-compose ps
    fi
}

# 显示服务访问地址  
show_urls() {
    local mode=${1:-prod}
    echo ""
    log_info "服务访问地址 ($mode 模式):"
    echo -e "${GREEN}  API服务:        ${NC}http://localhost:8888"
    echo -e "${GREEN}  Swagger文档:    ${NC}http://localhost:8888/swagger"
    echo -e "${GREEN}  健康检查:       ${NC}http://localhost:8888/health"
    echo -e "${GREEN}  MySQL数据库:    ${NC}localhost:3306 (root/root)"
    echo -e "${GREEN}  Redis缓存:      ${NC}localhost:6379"
    
    if [ "$mode" = "prod" ]; then
        echo -e "${GREEN}  Prometheus:     ${NC}http://localhost:9091"
        echo -e "${GREEN}  Grafana:        ${NC}http://localhost:3000 (admin/admin123)"
    fi
    echo ""
}

# 显示服务日志
show_service_logs() {
    local mode=${1:-prod}
    local service=${2:-backend-api}
    
    local container_name="volctrain-api"
    if [ "$mode" = "dev" ]; then
        container_name="volctrain-api-dev"
    fi
    
    log_info "显示 $service 服务日志 ($mode 模式):"
    if docker ps --format "table {{.Names}}" | grep -q "$container_name"; then
        echo "=== $service 服务日志 ==="
        docker logs --tail=20 "$container_name"
    else
        log_warn "$service 服务未运行"
    fi
}

# 显示日志
show_logs() {
    local mode=${1:-prod}
    local service=${2:-}
    
    if [ -n "$service" ]; then
        if [ "$mode" = "dev" ]; then
            log_info "显示 $service 服务日志 (开发模式):"
            docker-compose -f docker-compose.dev.yml logs -f "$service"
        else
            log_info "显示 $service 服务日志 (生产模式):"
            docker-compose logs -f "$service"
        fi
    else
        if [ "$mode" = "dev" ]; then
            log_info "显示所有服务日志 (开发模式):"
            docker-compose -f docker-compose.dev.yml logs -f
        else
            log_info "显示所有服务日志 (生产模式):"
            docker-compose logs -f
        fi
    fi
}

# 清理资源
cleanup() {
    local mode=${1:-both}
    log_info "清理 Docker 资源 ($mode)..."
    
    if [ "$mode" = "dev" ] || [ "$mode" = "both" ]; then
        docker-compose -f docker-compose.dev.yml down -v --remove-orphans >/dev/null 2>&1 || true
    fi
    
    if [ "$mode" = "prod" ] || [ "$mode" = "both" ]; then
        docker-compose down -v --remove-orphans >/dev/null 2>&1 || true
    fi
    
    # 清理系统资源
    docker system prune -f >/dev/null 2>&1 || true
    
    # 询问是否删除数据卷
    read -p "是否删除所有数据卷 (包括数据库数据)? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        docker volume rm $(docker volume ls -q | grep volctrain) >/dev/null 2>&1 || true
        log_info "数据卷已清理"
    fi
    
    log_info "清理完成"
}

# 显示帮助信息
show_help() {
    echo "VolcTrain AI训练平台启动脚本"
    echo ""
    echo "使用方法:"
    echo "  ./start.sh [COMMAND] [MODE]"
    echo ""
    echo "命令:"
    echo "  start [dev|prod]      启动服务 (默认prod)"
    echo "  stop [dev|prod]       停止服务"
    echo "  restart [dev|prod]    重启服务"
    echo "  status [dev|prod]     显示服务状态"
    echo "  logs [dev|prod] [service]  显示日志"
    echo "  build [dev|prod]      构建镜像"
    echo "  cleanup [dev|prod|both]    清理资源"
    echo "  test                  执行API测试"
    echo "  help                  显示帮助"
    echo ""
    echo "模式说明:"
    echo "  dev     开发环境 (仅MySQL + Redis + API)"
    echo "  prod    生产环境 (包含监控服务)"
    echo ""
    echo "示例:"
    echo "  ./start.sh start dev           # 启动开发环境"
    echo "  ./start.sh logs prod backend-api    # 显示生产环境API服务日志"
    echo "  ./start.sh cleanup dev         # 清理开发环境"
    echo ""
}

# API测试函数
test_api() {
    log_info "执行API功能测试..."
    
    # 检查API服务是否运行
    if ! curl -s http://localhost:8888/health >/dev/null 2>&1; then
        log_error "API服务未运行，请先启动服务"
        return 1
    fi
    
    # 测试健康检查端点
    log_info "测试健康检查端点..."
    response=$(curl -s http://localhost:8888/health)
    if echo "$response" | grep -q "status"; then
        log_info "✓ 健康检查端点测试通过"
    else
        log_error "✗ 健康检查端点测试失败"
    fi
    
    # 测试Swagger文档端点
    log_info "测试Swagger文档端点..."
    if curl -s http://localhost:8888/swagger >/dev/null 2>&1; then
        log_info "✓ Swagger文档端点测试通过"
    else
        log_error "✗ Swagger文档端点测试失败"
    fi
    
    # 测试Swagger JSON端点
    log_info "测试Swagger JSON端点..."
    if curl -s http://localhost:8888/swagger.json >/dev/null 2>&1; then
        log_info "✓ Swagger JSON端点测试通过"
    else
        log_error "✗ Swagger JSON端点测试失败"
    fi
    
    log_info "API测试完成"
    log_info "访问 http://localhost:8888/swagger 查看完整API文档"
}

# 主函数
main() {
    local command=${1:-start}
    local mode=${2:-prod}
    
    case "$command" in
        start)
            check_docker
            check_env_file
            check_directories
            build_images "$mode"
            start_services "$mode"
            ;;
        stop)
            stop_services "$mode"
            ;;
        restart)
            check_docker
            check_env_file
            check_directories
            restart_services "$mode"
            ;;
        status)
            show_status "$mode"
            show_urls "$mode"
            ;;
        logs)
            show_logs "$mode" "$3"
            ;;
        build)
            check_docker
            check_env_file
            check_directories
            build_images "$mode"
            ;;
        cleanup)
            cleanup "$mode"
            ;;
        test)
            test_api
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

# 捕获信号以优雅关闭
trap 'echo ""; log_info "接收到中断信号，正在退出..."; exit 0' INT TERM

# 执行主函数
main "$@"