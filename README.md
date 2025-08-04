# VolcTrain AI训练平台 - 后端服务

## 项目概述

VolcTrain 是一个企业级AI训练平台，提供GPU资源管理、训练任务调度、模型管理等功能。本项目是平台的后端API服务。

## 技术栈

- **框架**: Go-Zero 微服务框架
- **数据库**: MySQL 8.0
- **缓存**: Redis 7.0
- **认证**: JWT Token
- **文档**: Swagger/OpenAPI
- **容器**: Docker & Docker Compose

## 功能特性

### 核心功能
- 🔐 **用户认证**: JWT认证、角色权限管理
- ⚙️ **训练管理**: 训练任务创建、调度、监控
- 🖥️ **GPU管理**: GPU集群、节点、设备管理
- 📊 **监控告警**: 系统监控、指标采集、告警通知
- 📁 **文件管理**: 数据集、模型、日志文件管理

### API文档
- Swagger UI: `http://localhost:8888/swagger`
- API定义: `http://localhost:8888/swagger.json`
- 健康检查: `http://localhost:8888/health`

## 快速开始

### 环境要求
- Go 1.21+
- Docker & Docker Compose
- MySQL 8.0
- Redis 7.0

### 1. 克隆项目
```bash
git clone <项目地址>
cd VolcTrain/backend
```

### 2. 配置环境
```bash
# 复制环境配置
cp .env.example .env

# 修改数据库和Redis连接信息
vim .env
```

### 3. 启动服务

#### 方式一：使用启动脚本（推荐）
```bash
# 启动开发环境
./start.sh start dev

# 启动生产环境  
./start.sh start prod

# 查看服务状态
./start.sh status

# 查看日志
./start.sh logs
```

#### 方式二：Docker Compose
```bash
# 开发环境
docker-compose -f docker-compose.dev.yml up -d

# 生产环境
docker-compose up -d
```

#### 方式三：本地编译运行
```bash
# 编译
cd backend && go build -o bin/api cmd/api/main.go

# 运行
JWT_ACCESS_SECRET="your-secret-key" ./bin/api -f etc/config.yaml
```

### 4. 验证部署
```bash
# 健康检查
curl http://localhost:8888/health

# 访问API文档
open http://localhost:8888/swagger
```

## 项目结构

```
VolcTrain/
├── backend/
│   ├── cmd/
│   │   ├── api/           # API服务入口
│   │   └── test/          # 测试服务器
│   ├── internal/
│   │   ├── config/        # 配置定义
│   │   ├── handler/       # HTTP处理器
│   │   ├── logic/         # 业务逻辑
│   │   ├── svc/          # 服务上下文
│   │   └── types/        # 类型定义
│   ├── model/            # 数据模型
│   ├── pkg/              # 公共包
│   │   ├── auth/         # 认证相关
│   │   ├── database/     # 数据库连接
│   │   ├── docs/         # API文档
│   │   └── errors/       # 错误处理
│   ├── etc/              # 配置文件
│   ├── sql/              # 数据库脚本
│   └── bin/              # 编译输出
├── deploy/               # 部署配置
├── docs/                 # 项目文档
└── web/                  # 前端项目
```

## API接口

### 认证接口
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/refresh` - 刷新令牌
- `POST /api/v1/auth/logout` - 用户登出
- `GET /api/v1/auth/codes` - 获取权限码

### 用户管理
- `GET /api/v1/user/info` - 获取用户信息

### 训练任务
- `GET /api/v1/training/jobs` - 获取训练任务列表
- `POST /api/v1/training/jobs` - 创建训练任务
- `GET /api/v1/training/jobs/{id}` - 获取任务详情
- `PUT /api/v1/training/jobs/{id}` - 更新任务
- `DELETE /api/v1/training/jobs/{id}` - 删除任务

### GPU管理
- `GET /api/v1/gpuclusters` - 获取GPU集群
- `POST /api/v1/gpuclusters` - 创建GPU集群
- `GET /api/v1/gpudevices` - 获取GPU设备
- `POST /api/v1/gpudevices` - 创建GPU设备

## 配置说明

### 主要配置项
```yaml
# 服务配置
Name: volctrain-api
Host: 0.0.0.0
Port: 8888

# 数据库配置
MySQL:
  Host: localhost
  Port: 3306
  User: root
  Password: root
  DBName: volctraindb

# JWT配置
Auth:
  AccessSecret: your-secret-key
  AccessExpire: 86400
  RefreshExpire: 604800
```

### 环境变量
- `JWT_ACCESS_SECRET`: JWT密钥（必须）
- `MYSQL_HOST`: MySQL主机地址
- `MYSQL_PASSWORD`: MySQL密码
- `REDIS_HOST`: Redis主机地址

## 开发指南

### 本地开发
1. 安装依赖: `go mod tidy`
2. 配置数据库和Redis
3. 运行数据库初始化脚本
4. 启动服务: `go run cmd/api/main.go -f etc/config.yaml`

### 测试
```bash
# 运行单元测试
go test ./...

# 运行API测试
./backend/test_api.sh

# 启动测试服务器
go run backend/cmd/test/main.go
```

### 数据库管理
```bash
# 初始化数据库
mysql -u root -p < backend/sql/01_users_permissions.sql

# 查看数据库状态
./start.sh status
```

## 部署

### 生产环境部署
1. 配置环境变量
2. 构建Docker镜像
3. 使用Docker Compose部署
4. 配置反向代理（Nginx）
5. 设置监控和日志

### 监控
- Prometheus指标: `/metrics`
- 健康检查: `/health`
- 日志级别: info/debug/error

## 故障排除

### 常见问题
1. **数据库连接失败**: 检查MySQL服务状态和配置
2. **JWT认证失败**: 确认密钥配置正确
3. **端口占用**: 修改配置文件中的端口号
4. **Docker构建失败**: 检查网络连接和镜像源

### 日志查看
```bash
# 查看服务日志
./start.sh logs backend-api

# 查看所有服务日志  
./start.sh logs
```

## 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。
