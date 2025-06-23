# ColcTrain - AI 训练与调度平台

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/colctrain)](https://goreportcard.com/report/github.com/yourusername/colctrain)
[![Vue](https://img.shields.io/badge/Vue-3.x-brightgreen)](https://vuejs.org/)

ColcTrain 是一个现代化的 AI 训练任务管理与调度平台，基于 Go 语言和 Volcano 调度器构建，提供高效、可扩展的 AI 工作负载管理解决方案。

![ColcTrain Dashboard Preview](https://via.placeholder.com/800x400?text=ColcTrain+Dashboard+Preview)

## ✨ 核心特性

- **分布式训练支持**：无缝管理分布式 AI 训练任务（TensorFlow/PyTorch）
- **智能资源调度**：基于 Volcano 的高性能批处理调度器
- **可视化监控**：实时监控 GPU 利用率、任务进度和资源分配
- **多租户支持**：完善的用户/团队资源隔离与配额管理
- **流水线编排**：可视化 AI 训练流水线构建器
- **模型版本管理**：训练模型的全生命周期追踪
- **弹性伸缩**：根据负载自动扩缩容计算资源
- **多云支持**：可在本地集群或主流云平台上部署

## 🚀 技术栈

### 后端服务

- **语言**: Go 1.20+
- **调度引擎**: [Volcano](https://volcano.sh/) v1.7+
- **编排框架**: Kubernetes Operators
- **API 服务**: Gin Web Framework
- **存储**: PostgreSQL + Redis
- **监控**: Prometheus + Grafana

### 前端界面

- **框架**: Vue 3 + Composition API
- **语言**: TypeScript
- **UI 库**: Element Plus
- **可视化**: ECharts
- **构建工具**: Vite

## 🛠 安装指南

### 前提条件

- Kubernetes 集群 1.24+
- Helm 3.8+
- PostgreSQL 12+
- Redis 6+

### 使用 Helm 部署

```bash
# 添加仓库
helm repo add colctrain https://colctrain.github.io/helm-charts

# 安装
helm install colctrain colctrain/colctrain \
  --namespace colctrain-system \
  --create-namespace \
  --set global.postgresql.url="postgres://user:pass@postgres:5432/colctrain" \
  --set global.redis.url="redis://redis:6379"
```

### 本地开发环境

```bash
# 启动后端服务
make run-backend

# 启动前端开发服务器
cd frontend
npm install
npm run dev
```

## 🧭 使用示例

### 1. 提交训练任务 (YAML 示例)

```yaml
apiVersion: colctrain.ai/v1
kind: TrainingJob
metadata:
  name: resnet50-training
spec:
  framework: pytorch
  workers: 4
  gpuPerWorker: 2
  image: nvcr.io/nvidia/pytorch:22.12-py3
  command:
    - "python"
    - "/workspace/train.py"
    - "--epochs=50"
    - "--batch-size=256"
  datasets:
    - name: imagenet-1k
      mountPath: /data
  resourceProfile: high-mem-gpu
```

### 2. 通过 REST API 管理任务

```bash
# 提交新任务
curl -X POST https://api.colctrain.io/jobs \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"mnist-training", "framework":"tensorflow", ...}'

# 获取任务状态
curl -X GET https://api.colctrain.io/jobs/mnist-training/status
```

### 3. 前端界面操作

```typescript
// Vue组件示例
import { useJobStore } from "@/stores/job";

const jobStore = useJobStore();

// 提交训练任务
const submitJob = async () => {
  const jobConfig = {
    name: "text-classification",
    framework: "huggingface",
    // ...其他配置
  };

  await jobStore.createJob(jobConfig);
};
```

## 📊 系统架构

```mermaid
graph TD
    A[Web UI] --> B[API Gateway]
    B --> C[Job Manager]
    C --> D[Volcano Scheduler]
    D --> E[Kubernetes Cluster]
    C --> F[Resource Monitor]
    F --> G[Prometheus]
    C --> H[Model Registry]
    H --> I[PostgreSQL]
    C --> J[Dataset Manager]
    J --> K[NFS/Cloud Storage]
```

## 🤝 参与贡献

我们欢迎各种形式的贡献！请阅读我们的[贡献指南](CONTRIBUTING.md)开始参与。

1. 提交 Issue 报告 bug 或建议新功能
2. Fork 仓库并创建分支 (`git checkout -b feature/your-feature`)
3. 提交 Pull Request
4. 确保所有测试通过并更新文档

## 📞 联系我们

- 邮箱: 13664854532@163.com
- 微信: GoSimplicity

---

**让 AI 训练任务管理变得简单高效** - ColcTrain 团队 🚀
