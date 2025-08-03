#!/bin/bash

# VolcTrain 部署脚本
# 用途：一键部署VolcTrain AI训练调度平台

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

log_warn() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查环境变量
check_env() {
    log_info "检查环境变量..."
    
    required_vars=(
        "MYSQL_HOST"
        "MYSQL_DATABASE" 
        "MYSQL_USER"
        "MYSQL_PASSWORD"
        "DOMAIN_NAME"
    )
    
    missing_vars=()
    for var in "${required_vars[@]}"; do
        if [[ -z "${!var}" ]]; then
            missing_vars+=("$var")
        fi
    done
    
    if [[ ${#missing_vars[@]} -gt 0 ]]; then
        log_error "以下环境变量未设置："
        printf '%s\n' "${missing_vars[@]}"
        log_error "请设置这些环境变量后重新运行"
        exit 1
    fi
    
    log_success "环境变量检查通过"
}

# 检查依赖工具
check_dependencies() {
    log_info "检查依赖工具..."
    
    tools=("docker" "kubectl" "helm")
    missing_tools=()
    
    for tool in "${tools[@]}"; do
        if ! command -v "$tool" &> /dev/null; then
            missing_tools+=("$tool")
        fi
    done
    
    if [[ ${#missing_tools[@]} -gt 0 ]]; then
        log_error "以下工具未安装："
        printf '%s\n' "${missing_tools[@]}"
        log_error "请安装这些工具后重新运行"
        exit 1
    fi
    
    log_success "依赖工具检查通过"
}

# 构建Docker镜像
build_images() {
    log_info "构建Docker镜像..."
    
    # 设置镜像标签
    IMAGE_TAG="${IMAGE_TAG:-latest}"
    REGISTRY="${REGISTRY:-volctrain}"
    
    # 构建后端镜像
    log_info "构建后端镜像..."
    docker build -t "${REGISTRY}/backend:${IMAGE_TAG}" ./backend/
    
    # 构建前端镜像
    log_info "构建前端镜像..."
    docker build -t "${REGISTRY}/frontend:${IMAGE_TAG}" ./web/
    
    log_success "Docker镜像构建完成"
}

# 推送镜像到仓库
push_images() {
    if [[ "${PUSH_IMAGES:-false}" == "true" ]]; then
        log_info "推送镜像到仓库..."
        
        docker push "${REGISTRY}/backend:${IMAGE_TAG}"
        docker push "${REGISTRY}/frontend:${IMAGE_TAG}"
        
        log_success "镜像推送完成"
    fi
}

# 创建Kubernetes命名空间
create_namespace() {
    log_info "创建Kubernetes命名空间..."
    
    if kubectl get namespace volctrain &> /dev/null; then
        log_warn "命名空间 volctrain 已存在"
    else
        kubectl create namespace volctrain
        log_success "命名空间 volctrain 创建成功"
    fi
}

# 创建配置和密钥
create_config() {
    log_info "创建配置和密钥..."
    
    # 创建配置文件
    envsubst < deploy/k8s/01-namespace-config.yaml | kubectl apply -f -
    
    # 创建密钥
    kubectl create secret generic volctrain-secrets \
        --from-literal=db-password="$(echo -n "$MYSQL_PASSWORD" | base64)" \
        --from-literal=redis-password="$(echo -n "${REDIS_PASSWORD:-}" | base64)" \
        --from-literal=jwt-secret="$(echo -n "${JWT_SECRET:-volctrain-jwt-secret}" | base64)" \
        --namespace=volctrain \
        --dry-run=client -o yaml | kubectl apply -f -
    
    log_success "配置和密钥创建完成"
}

# 部署使用Docker Compose
deploy_docker_compose() {
    log_info "使用Docker Compose部署..."
    
    # 检查环境文件
    if [[ ! -f .env.prod ]]; then
        log_info "创建生产环境配置文件..."
        cat > .env.prod << EOF
# 数据库配置
MYSQL_HOST=${MYSQL_HOST}
MYSQL_PORT=${MYSQL_PORT:-3306}
MYSQL_DATABASE=${MYSQL_DATABASE}
MYSQL_USER=${MYSQL_USER}
MYSQL_PASSWORD=${MYSQL_PASSWORD}

# Redis配置
REDIS_HOST=${REDIS_HOST:-redis}
REDIS_PORT=${REDIS_PORT:-6379}
REDIS_PASSWORD=${REDIS_PASSWORD:-}

# 镜像配置
REGISTRY=${REGISTRY:-volctrain}
TAG=${IMAGE_TAG:-latest}

# 日志级别
LOG_LEVEL=${LOG_LEVEL:-info}

# Grafana配置
GRAFANA_PASSWORD=${GRAFANA_PASSWORD:-admin123}
EOF
    fi
    
    # 部署服务
    docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d
    
    log_success "Docker Compose部署完成"
}

# 部署到Kubernetes
deploy_kubernetes() {
    log_info "部署到Kubernetes..."
    
    # 应用所有配置文件
    for config_file in deploy/k8s/*.yaml; do
        log_info "应用配置文件: $config_file"
        envsubst < "$config_file" | kubectl apply -f -
    done
    
    # 等待部署完成
    log_info "等待部署完成..."
    kubectl wait --for=condition=ready pod -l app=volctrain-api --namespace=volctrain --timeout=300s
    kubectl wait --for=condition=ready pod -l app=volctrain-frontend --namespace=volctrain --timeout=300s
    
    log_success "Kubernetes部署完成"
}

# 使用Helm部署
deploy_helm() {
    log_info "使用Helm部署..."
    
    # 添加必要的Helm仓库
    helm repo add nginx-stable https://helm.nginx.com/stable
    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    helm repo update
    
    # 安装或升级
    helm upgrade --install volctrain ./deploy/helm/volctrain \
        --namespace volctrain \
        --create-namespace \
        --set database.host="$MYSQL_HOST" \
        --set database.name="$MYSQL_DATABASE" \
        --set database.user="$MYSQL_USER" \
        --set database.password="$MYSQL_PASSWORD" \
        --set redis.host="${REDIS_HOST:-redis}" \
        --set ingress.hosts[0].host="$DOMAIN_NAME" \
        --set ingress.tls[0].hosts[0]="$DOMAIN_NAME" \
        --set image.backend.tag="$IMAGE_TAG" \
        --set image.frontend.tag="$IMAGE_TAG" \
        --wait
    
    log_success "Helm部署完成"
}

# 验证部署
verify_deployment() {
    log_info "验证部署..."
    
    if [[ "$DEPLOY_METHOD" == "docker-compose" ]]; then
        # 验证Docker Compose部署
        if docker-compose -f docker-compose.prod.yml ps | grep -q "Up"; then
            log_success "服务运行正常"
        else
            log_error "部分服务未运行"
            exit 1
        fi
    else
        # 验证Kubernetes部署
        kubectl get pods -n volctrain
        
        # 检查服务状态
        if kubectl get pods -n volctrain -o jsonpath='{.items[*].status.phase}' | grep -v Running; then
            log_warn "部分Pod未就绪，请检查状态"
        else
            log_success "所有Pod运行正常"
        fi
        
        # 显示访问地址
        if [[ "$DEPLOY_METHOD" == "k8s" ]]; then
            EXTERNAL_IP=$(kubectl get service volctrain-nodeport -n volctrain -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "")
            if [[ -n "$EXTERNAL_IP" ]]; then
                log_success "应用访问地址: http://$EXTERNAL_IP:30080"
            else
                log_info "使用NodePort访问: http://<节点IP>:30080"
            fi
        else
            log_success "应用访问地址: https://$DOMAIN_NAME"
        fi
    fi
}

# 清理函数
cleanup() {
    if [[ "$1" == "all" ]]; then
        log_warn "清理所有资源..."
        
        if [[ "$DEPLOY_METHOD" == "docker-compose" ]]; then
            docker-compose -f docker-compose.prod.yml down -v
        elif [[ "$DEPLOY_METHOD" == "helm" ]]; then
            helm uninstall volctrain -n volctrain
            kubectl delete namespace volctrain
        else
            kubectl delete namespace volctrain
        fi
        
        log_success "清理完成"
    fi
}

# 显示帮助信息
show_help() {
    cat << EOF
VolcTrain 部署脚本

用法: $0 [选项] <命令>

命令:
  build       构建Docker镜像
  deploy      部署应用
  verify      验证部署
  cleanup     清理资源

选项:
  -m, --method METHOD     部署方法 (docker-compose|k8s|helm，默认: helm)
  -t, --tag TAG          镜像标签 (默认: latest)
  -r, --registry REG     镜像仓库 (默认: volctrain)
  -p, --push             构建后推送镜像
  -h, --help             显示此帮助信息

环境变量:
  MYSQL_HOST            MySQL主机地址 (必需)
  MYSQL_DATABASE        数据库名称 (必需)
  MYSQL_USER            数据库用户 (必需)
  MYSQL_PASSWORD        数据库密码 (必需)
  DOMAIN_NAME           域名 (Kubernetes部署时必需)
  REDIS_HOST            Redis主机地址 (可选)
  
示例:
  # 使用Docker Compose部署
  $0 -m docker-compose deploy
  
  # 使用Kubernetes部署
  $0 -m k8s deploy
  
  # 使用Helm部署（推荐）
  $0 -m helm deploy
  
  # 构建并推送镜像
  $0 -p build
EOF
}

# 主函数
main() {
    # 解析命令行参数
    DEPLOY_METHOD="helm"
    IMAGE_TAG="latest"
    REGISTRY="volctrain"
    PUSH_IMAGES="false"
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -m|--method)
                DEPLOY_METHOD="$2"
                shift 2
                ;;
            -t|--tag)
                IMAGE_TAG="$2"
                shift 2
                ;;
            -r|--registry)
                REGISTRY="$2"
                shift 2
                ;;
            -p|--push)
                PUSH_IMAGES="true"
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            build|deploy|verify|cleanup)
                COMMAND="$1"
                shift
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    if [[ -z "$COMMAND" ]]; then
        log_error "请指定命令"
        show_help
        exit 1
    fi
    
    # 执行命令
    case "$COMMAND" in
        build)
            check_dependencies
            build_images
            push_images
            ;;
        deploy)
            check_env
            check_dependencies
            build_images
            push_images
            
            case "$DEPLOY_METHOD" in
                docker-compose)
                    deploy_docker_compose
                    ;;
                k8s)
                    create_namespace
                    create_config
                    deploy_kubernetes
                    ;;
                helm)
                    deploy_helm
                    ;;
                *)
                    log_error "不支持的部署方法: $DEPLOY_METHOD"
                    exit 1
                    ;;
            esac
            
            verify_deployment
            ;;
        verify)
            verify_deployment
            ;;
        cleanup)
            cleanup all
            ;;
        *)
            log_error "未知命令: $COMMAND"
            show_help
            exit 1
            ;;
    esac
}

# 捕获中断信号
trap 'log_warn "部署被中断"; exit 1' INT TERM

# 运行主函数
main "$@"