#!/bin/bash

# VolcTrain 开发环境启动脚本
# 用于同时启动后端API服务器和前端开发服务器

set -e

echo "🚀 启动 VolcTrain 开发环境"
echo "================================"

# 检查必要工具
echo "📋 检查开发环境依赖..."

if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装或未在PATH中"
    exit 1
fi

if ! command -v pnpm &> /dev/null; then
    echo "❌ pnpm 未安装或未在PATH中"
    exit 1
fi

echo "✅ 开发环境检查通过"

# 获取项目根目录
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
echo "📁 项目根目录: $PROJECT_ROOT"

# 启动后端API服务器
echo ""
echo "🔧 启动后端API服务器..."
echo "--------------------------------"
cd "$PROJECT_ROOT/backend"

# 检查配置文件
if [ ! -f "./etc/config-dev.yaml" ]; then
    echo "❌ 后端配置文件不存在: ./etc/config-dev.yaml"
    exit 1
fi

# 后台启动API服务器
echo "📡 启动API服务器 (端口: 8888)..."
go run cmd/api/main.go &
API_PID=$!
echo "✅ 后端API服务器已启动 (PID: $API_PID)"

# 等待API服务器启动
echo "⏳ 等待API服务器启动..."
sleep 3

# 检查API服务器状态
if ps -p $API_PID > /dev/null; then
    echo "✅ API服务器运行正常"
else
    echo "❌ API服务器启动失败"
    exit 1
fi

# 启动前端开发服务器
echo ""
echo "🎨 启动前端开发服务器..."
echo "--------------------------------"
cd "$PROJECT_ROOT/web"

# 检查依赖
if [ ! -d "./node_modules" ]; then
    echo "📦 安装前端依赖..."
    pnpm install
fi

echo "🌐 启动前端服务器 (端口: 5173)..."
pnpm dev &
WEB_PID=$!
echo "✅ 前端开发服务器已启动 (PID: $WEB_PID)"

# 创建清理函数
cleanup() {
    echo ""
    echo "🧹 清理进程..."
    if ps -p $API_PID > /dev/null; then
        kill $API_PID
        echo "🔴 API服务器已停止"
    fi
    if ps -p $WEB_PID > /dev/null; then
        kill $WEB_PID
        echo "🔴 前端服务器已停止"
    fi
    echo "👋 开发环境已关闭"
}

# 注册清理函数
trap cleanup EXIT INT TERM

echo ""
echo "🎉 VolcTrain 开发环境启动成功!"
echo "================================"
echo "📡 后端API: http://localhost:8888"
echo "🌐 前端界面: http://localhost:5173"
echo "📊 API文档: http://localhost:8888/swagger"
echo ""
echo "按 Ctrl+C 停止所有服务"
echo ""

# 等待用户中断
wait