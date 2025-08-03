#!/bin/bash

# VolcTrain监控系统集成测试脚本

set -e

# 配置变量
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="$PROJECT_DIR/backend"
TEST_DB_NAME="volctrain_test"
TEST_CONFIG="$BACKEND_DIR/etc/monitoring-test.yaml"

echo "🚀 开始VolcTrain监控系统集成测试"
echo "项目目录: $PROJECT_DIR"
echo "后端目录: $BACKEND_DIR"

# 颜色输出函数
print_success() {
    echo -e "\033[32m✅ $1\033[0m"
}

print_error() {
    echo -e "\033[31m❌ $1\033[0m"
}

print_info() {
    echo -e "\033[34mℹ️  $1\033[0m"
}

print_warning() {
    echo -e "\033[33m⚠️  $1\033[0m"
}

# 检查依赖
check_dependencies() {
    print_info "检查系统依赖..."
    
    # 检查Go
    if ! command -v go &> /dev/null; then
        print_error "Go未安装，请先安装Go 1.19+"
        exit 1
    fi
    
    # 检查MySQL
    if ! command -v mysql &> /dev/null; then
        print_error "MySQL客户端未安装"
        exit 1
    fi
    
    # 检查Docker（可选）
    if command -v docker &> /dev/null; then
        print_success "Docker已安装"
    else
        print_warning "Docker未安装，跳过容器化测试"
    fi
    
    print_success "依赖检查完成"
}

# 设置测试数据库
setup_test_database() {
    print_info "设置测试数据库..."
    
    # 创建测试数据库
    mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS $TEST_DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" || {
        print_error "创建测试数据库失败"
        exit 1
    }
    
    # 导入表结构
    if [ -f "$BACKEND_DIR/sql/create_tables.sql" ]; then
        mysql -u root -p $TEST_DB_NAME < "$BACKEND_DIR/sql/create_tables.sql" || {
            print_error "导入表结构失败"
            exit 1
        }
    fi
    
    print_success "测试数据库设置完成"
}

# 创建测试配置文件
create_test_config() {
    print_info "创建测试配置文件..."
    
    cat > "$TEST_CONFIG" << EOF
# 测试环境配置
database:
  host: localhost
  port: 3306
  username: root
  password: ""
  database: $TEST_DB_NAME
  max_open_conns: 10
  max_idle_conns: 5

monitoring:
  metrics:
    prometheus_url: "http://localhost:9090"
    collect_interval: "10s"
    system_metrics_enabled: true
    business_metrics_enabled: true
    prometheus_enabled: false  # 测试环境禁用Prometheus
    metrics_port: 9092
    enable_builtin_metrics: true
    max_retries: 2
    timeout: "10s"

  alerts:
    evaluation_interval: "10s"
    max_concurrent_rules: 5
    alert_retention_days: 7
    enable_grouping: true
    enable_suppression: true
    default_throttle_minutes: 1
    anomaly_detection_enabled: true

  notifications:
    max_queue_size: 100
    max_concurrent_senders: 2
    retry_max_attempts: 2
    retry_backoff_seconds: 2
    rate_limit_per_minute: 30
    timeout_seconds: 10
    failed_retention_days: 1
    enable_deduplication: true
    deduplication_window: "2m"

  system:
    enable_metrics: true
    enable_alerts: true
    enable_notifications: true
    health_check_interval: "10s"
    auto_restart: true
    max_restart_attempts: 2

port: 8081
EOF
    
    print_success "测试配置文件创建完成"
}

# 编译项目
build_project() {
    print_info "编译项目..."
    
    cd "$BACKEND_DIR"
    
    # 下载依赖
    go mod tidy || {
        print_error "下载Go依赖失败"
        exit 1
    }
    
    # 编译监控服务
    go build -o bin/volctrain-monitoring cmd/monitoring/main.go || {
        print_error "编译监控服务失败"
        exit 1
    }
    
    print_success "项目编译完成"
}

# 运行单元测试
run_unit_tests() {
    print_info "运行单元测试..."
    
    cd "$BACKEND_DIR"
    
    # 运行pkg包的单元测试
    go test -v ./pkg/monitoring/... || {
        print_error "监控模块单元测试失败"
        return 1
    }
    
    go test -v ./pkg/alerting/... || {
        print_error "告警模块单元测试失败"
        return 1
    }
    
    go test -v ./pkg/notification/... || {
        print_error "通知模块单元测试失败"
        return 1
    }
    
    print_success "单元测试通过"
}

# 运行集成测试
run_integration_tests() {
    print_info "运行集成测试..."
    
    cd "$BACKEND_DIR"
    
    # 设置测试环境变量
    export TEST_DB_DSN="root:@tcp(localhost:3306)/$TEST_DB_NAME?charset=utf8mb4&parseTime=True&loc=Local"
    
    # 运行集成测试
    go test -v -tags=integration ./tests/integration/... || {
        print_error "集成测试失败"
        return 1
    }
    
    print_success "集成测试通过"
}

# 启动监控服务进行功能测试
start_monitoring_service() {
    print_info "启动监控服务进行功能测试..."
    
    cd "$BACKEND_DIR"
    
    # 后台启动监控服务
    ./bin/volctrain-monitoring -c "$TEST_CONFIG" > /tmp/volctrain-monitoring.log 2>&1 &
    MONITORING_PID=$!
    
    # 等待服务启动
    sleep 5
    
    # 检查服务是否正常启动
    if ! kill -0 $MONITORING_PID 2>/dev/null; then
        print_error "监控服务启动失败"
        cat /tmp/volctrain-monitoring.log
        return 1
    fi
    
    print_success "监控服务启动成功 (PID: $MONITORING_PID)"
    
    # 功能测试
    test_monitoring_endpoints
    
    # 停止服务
    kill $MONITORING_PID 2>/dev/null || true
    wait $MONITORING_PID 2>/dev/null || true
    
    print_success "监控服务功能测试完成"
}

# 测试监控端点
test_monitoring_endpoints() {
    print_info "测试监控服务端点..."
    
    # 健康检查端点
    if curl -s http://localhost:8081/health | grep -q "healthy"; then
        print_success "健康检查端点正常"
    else
        print_error "健康检查端点异常"
        return 1
    fi
    
    # 就绪检查端点
    if curl -s http://localhost:8081/ready | grep -q "ready"; then
        print_success "就绪检查端点正常"
    else
        print_error "就绪检查端点异常"
        return 1
    fi
    
    # 系统状态端点
    if curl -s http://localhost:8081/api/v1/status | grep -q "overall_status"; then
        print_success "系统状态端点正常"
    else
        print_error "系统状态端点异常"
        return 1
    fi
    
    print_success "所有端点测试通过"
}

# 性能测试
run_performance_tests() {
    print_info "运行性能测试..."
    
    cd "$BACKEND_DIR"
    
    # 如果有性能测试文件，运行它们
    if [ -f "tests/performance/monitoring_performance_test.go" ]; then
        go test -v -bench=. ./tests/performance/... || {
            print_warning "性能测试失败，但不影响主要功能"
        }
    else
        print_warning "性能测试文件不存在，跳过性能测试"
    fi
    
    print_success "性能测试完成"
}

# 清理测试环境
cleanup() {
    print_info "清理测试环境..."
    
    # 停止可能运行的服务
    pkill -f "volctrain-monitoring" 2>/dev/null || true
    
    # 删除测试数据库
    mysql -u root -p -e "DROP DATABASE IF EXISTS $TEST_DB_NAME;" 2>/dev/null || true
    
    # 删除测试配置文件
    rm -f "$TEST_CONFIG"
    
    # 删除日志文件
    rm -f /tmp/volctrain-monitoring.log
    
    print_success "测试环境清理完成"
}

# 生成测试报告
generate_test_report() {
    print_info "生成测试报告..."
    
    REPORT_FILE="$PROJECT_DIR/test-report-$(date +%Y%m%d-%H%M%S).md"
    
    cat > "$REPORT_FILE" << EOF
# VolcTrain监控系统集成测试报告

## 测试概要
- 测试时间: $(date)
- 测试环境: $(uname -a)
- Go版本: $(go version)
- 项目目录: $PROJECT_DIR

## 测试结果

### ✅ 通过的测试
- 依赖检查
- 项目编译
- 数据库设置
- 单元测试
- 集成测试
- 监控服务启动
- API端点测试

### 📊 测试覆盖的功能模块
- 指标收集系统
- 告警引擎
- 通知管理器
- 监控服务集成
- 健康检查机制
- HTTP API接口

### 🔧 测试配置
- 测试数据库: $TEST_DB_NAME
- 监控端口: 8081
- 指标端口: 9092

## 结论
✅ VolcTrain监控系统集成测试全部通过，系统功能正常。

---
生成时间: $(date)
EOF
    
    print_success "测试报告已生成: $REPORT_FILE"
}

# 主测试流程
main() {
    echo "=================================================="
    echo "🔍 VolcTrain监控系统集成测试开始"
    echo "=================================================="
    
    # 设置错误处理
    trap cleanup EXIT
    
    # 执行测试步骤
    check_dependencies
    setup_test_database
    create_test_config
    build_project
    
    print_info "开始测试流程..."
    
    # 运行各种测试
    if run_unit_tests && run_integration_tests && start_monitoring_service; then
        print_success "🎉 所有测试通过！"
        
        # 可选的性能测试
        run_performance_tests
        
        # 生成报告
        generate_test_report
        
        echo "=================================================="
        echo "✅ VolcTrain监控系统集成测试成功完成"
        echo "=================================================="
        
        exit 0
    else
        print_error "❌ 测试失败"
        echo "=================================================="
        echo "❌ VolcTrain监控系统集成测试失败"
        echo "=================================================="
        
        exit 1
    fi
}

# 帮助信息
show_help() {
    echo "VolcTrain监控系统集成测试脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help     显示此帮助信息"
    echo "  --unit-only    仅运行单元测试"
    echo "  --integration-only  仅运行集成测试"
    echo "  --cleanup      仅执行清理操作"
    echo ""
    echo "示例:"
    echo "  $0              # 运行完整测试"
    echo "  $0 --unit-only  # 仅运行单元测试"
    echo "  $0 --cleanup    # 清理测试环境"
}

# 解析命令行参数
case "${1:-}" in
    -h|--help)
        show_help
        exit 0
        ;;
    --unit-only)
        check_dependencies
        build_project
        run_unit_tests
        exit $?
        ;;
    --integration-only)
        check_dependencies
        setup_test_database
        create_test_config
        build_project
        run_integration_tests
        cleanup
        exit $?
        ;;
    --cleanup)
        cleanup
        exit 0
        ;;
    "")
        main
        ;;
    *)
        echo "未知选项: $1"
        show_help
        exit 1
        ;;
esac