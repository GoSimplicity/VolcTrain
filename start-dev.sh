#!/bin/bash
# 开发环境快速启动脚本

set -e

echo "🚀 启动 VolcTrain 开发环境..."

# 检查 .env 文件
if [ ! -f .env ]; then
    echo "📋 创建开发环境配置..."
    cp .env.example .env
    
    # 设置开发环境默认配置
    sed -i.bak 's/DEPLOY_ENV=production/DEPLOY_ENV=development/g' .env
    sed -i.bak 's/LOG_LEVEL=info/LOG_LEVEL=debug/g' .env
    rm .env.bak
fi

# 启动基础服务（数据库 + 缓存）
echo "🗄️  启动数据库和缓存服务..."
docker-compose up -d mysql redis

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 10

# 启动后端服务
echo "🔧 启动后端服务..."
docker-compose up -d backend-api backend-common backend-monitoring backend-training

# 等待后端服务启动
sleep 5

# 启动前端和监控
echo "🌐 启动前端和监控服务..."
docker-compose up -d frontend prometheus grafana

echo "✅ 开发环境启动完成！"
echo ""
echo "📱 访问地址:"
echo "   前端: http://localhost:80"
echo "   API: http://localhost:8888"
echo "   Grafana: http://localhost:3000 (admin/admin123)"
echo ""
echo "📊 查看日志: docker-compose logs -f"
echo "🛑 停止服务: docker-compose down"