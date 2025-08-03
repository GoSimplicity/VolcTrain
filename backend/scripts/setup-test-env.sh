#!/bin/bash

# VolcTrain 测试环境设置脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

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

# 检查Docker是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi
    
    log_info "Docker 环境检查通过"
}

# 启动测试环境
start_test_env() {
    log_info "启动测试环境..."
    
    # 停止并清理现有容器
    docker-compose -f docker-compose.test.yml down -v 2>/dev/null || true
    
    # 启动测试服务
    docker-compose -f docker-compose.test.yml up -d
    
    log_info "等待服务启动..."
    sleep 10
    
    # 检查服务状态
    check_services_health
}

# 检查服务健康状态
check_services_health() {
    log_info "检查服务健康状态..."
    
    # 检查MySQL
    local mysql_ready=false
    for i in {1..30}; do
        if docker exec volctrain-mysql-test mysqladmin ping -h localhost -u root -proot123 &>/dev/null; then
            mysql_ready=true
            break
        fi
        log_info "等待MySQL启动... ($i/30)"
        sleep 2
    done
    
    if [ "$mysql_ready" = false ]; then
        log_error "MySQL 启动失败"
        exit 1
    fi
    log_info "✅ MySQL 已就绪"
    
    # 检查Redis
    if docker exec volctrain-redis-test redis-cli ping &>/dev/null; then
        log_info "✅ Redis 已就绪"
    else
        log_error "Redis 启动失败"
        exit 1
    fi
    
    # 检查Prometheus
    if curl -s http://localhost:9091/-/healthy &>/dev/null; then
        log_info "✅ Prometheus 已就绪"
    else
        log_warn "Prometheus 可能未完全启动"
    fi
    
    log_info "🎉 测试环境启动完成！"
}

# 初始化测试数据库
init_test_database() {
    log_info "初始化测试数据库..."
    
    # 等待MySQL完全启动
    sleep 5
    
    # 检查数据库是否存在
    if docker exec volctrain-mysql-test mysql -u volctrain -pvolctrain -e "USE volctrain_test;" &>/dev/null; then
        log_info "数据库 volctrain_test 已存在"
    else
        log_error "数据库初始化失败"
        exit 1
    fi
    
    # 执行SQL脚本初始化表结构
    if [ -d "./sql" ]; then
        log_info "执行数据库迁移脚本..."
        for sql_file in ./sql/*.sql; do
            if [ -f "$sql_file" ]; then
                log_info "执行: $(basename $sql_file)"
                docker exec -i volctrain-mysql-test mysql -u volctrain -pvolctrain volctrain_test < "$sql_file"
            fi
        done
        log_info "✅ 数据库初始化完成"
    else
        log_warn "未找到SQL脚本目录"
    fi
}

# 设置测试环境变量
setup_test_env_vars() {
    log_info "设置测试环境变量..."
    
    cat > .env.test << EOF
# 测试环境配置
TEST_DB_HOST=localhost
TEST_DB_PORT=3307
TEST_DB_NAME=volctrain_test
TEST_DB_USER=volctrain
TEST_DB_PASS=volctrain

TEST_REDIS_HOST=localhost
TEST_REDIS_PORT=6380
TEST_REDIS_PASS=

TEST_API_HOST=localhost
TEST_API_PORT=8888

# 监控配置
PROMETHEUS_URL=http://localhost:9091
GRAFANA_URL=http://localhost:3001
JAEGER_URL=http://localhost:16687

# 测试标志
TESTING_MODE=true
LOG_LEVEL=debug
EOF
    
    log_info "✅ 环境变量配置完成"
}

# 运行测试数据种子
seed_test_data() {
    log_info "生成测试数据..."
    
    # 检查是否有种子脚本
    if [ -f "./test/seeds/run_seeds.go" ]; then
        go run ./test/seeds/run_seeds.go
        log_info "✅ 测试数据生成完成"
    else
        log_warn "未找到测试数据种子脚本，将在测试中动态生成"
    fi
}

# 验证测试环境
verify_test_env() {
    log_info "验证测试环境..."
    
    # 运行基础测试
    export TEST_DB_HOST=localhost
    export TEST_DB_PORT=3307
    export TEST_DB_NAME=volctrain_test
    export TEST_DB_USER=volctrain
    export TEST_DB_PASS=volctrain
    export TEST_REDIS_HOST=localhost
    export TEST_REDIS_PORT=6380
    
    log_info "运行测试框架验证..."
    if go test -v ./test/framework_test.go; then
        log_info "✅ 测试环境验证通过"
    else
        log_error "测试环境验证失败"
        exit 1
    fi
}

# 停止测试环境
stop_test_env() {
    log_info "停止测试环境..."
    docker-compose -f docker-compose.test.yml down -v
    log_info "✅ 测试环境已停止"
}

# 清理测试环境
cleanup_test_env() {
    log_info "清理测试环境..."
    docker-compose -f docker-compose.test.yml down -v --rmi local
    docker volume prune -f
    rm -f .env.test
    log_info "✅ 测试环境清理完成"
}

# 显示帮助信息
show_help() {
    cat << EOF
VolcTrain 测试环境管理脚本

用法: $0 [命令]

命令:
  start       启动测试环境
  stop        停止测试环境
  restart     重启测试环境
  status      查看环境状态
  init        初始化数据库
  seed        生成测试数据
  verify      验证测试环境
  cleanup     清理测试环境
  logs        查看服务日志
  help        显示帮助信息

示例:
  $0 start      # 启动完整测试环境
  $0 verify     # 验证环境是否正常
  $0 cleanup    # 清理所有测试资源
EOF
}

# 查看环境状态
show_status() {
    log_info "测试环境状态:"
    docker-compose -f docker-compose.test.yml ps
    
    echo ""
    log_info "服务健康检查:"
    
    # MySQL状态
    if docker exec volctrain-mysql-test mysqladmin ping -h localhost -u root -proot123 &>/dev/null; then
        echo "MySQL: ✅ 运行中"
    else
        echo "MySQL: ❌ 未运行"
    fi
    
    # Redis状态
    if docker exec volctrain-redis-test redis-cli ping &>/dev/null; then
        echo "Redis: ✅ 运行中"
    else
        echo "Redis: ❌ 未运行"
    fi
    
    # Prometheus状态
    if curl -s http://localhost:9091/-/healthy &>/dev/null; then
        echo "Prometheus: ✅ 运行中"
    else
        echo "Prometheus: ❌ 未运行"
    fi
}

# 查看日志
show_logs() {
    if [ -n "$2" ]; then
        docker-compose -f docker-compose.test.yml logs -f "$2"
    else
        docker-compose -f docker-compose.test.yml logs -f
    fi
}

# 主函数
main() {
    case "${1:-help}" in
        "start")
            check_docker
            start_test_env
            init_test_database
            setup_test_env_vars
            seed_test_data
            verify_test_env
            log_info "🚀 测试环境已完全就绪！"
            echo ""
            echo "测试服务地址:"
            echo "  MySQL: localhost:3307"
            echo "  Redis: localhost:6380"
            echo "  Prometheus: http://localhost:9091"
            echo "  Grafana: http://localhost:3001 (admin/admin123)"
            echo "  Jaeger: http://localhost:16687"
            ;;
        "stop")
            stop_test_env
            ;;
        "restart")
            stop_test_env
            sleep 2
            start_test_env
            init_test_database
            ;;
        "status")
            show_status
            ;;
        "init")
            init_test_database
            ;;
        "seed")
            seed_test_data
            ;;
        "verify")
            verify_test_env
            ;;
        "cleanup")
            cleanup_test_env
            ;;
        "logs")
            show_logs "$@"
            ;;
        "help"|*)
            show_help
            ;;
    esac
}

# 执行主函数
main "$@"