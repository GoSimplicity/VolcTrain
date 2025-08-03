# 构建阶段
FROM node:18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制package文件
COPY pnpm-lock.yaml package.json ./
COPY apps/web-antd/package.json ./apps/web-antd/

# 安装pnpm
RUN npm install -g pnpm

# 安装依赖
RUN pnpm install --frozen-lockfile

# 复制源代码
COPY . .

# 构建应用
WORKDIR /app/apps/web-antd
RUN pnpm build

# 生产阶段
FROM nginx:alpine

# 安装必要工具
RUN apk add --no-cache curl

# 复制构建结果
COPY --from=builder /app/apps/web-antd/dist /usr/share/nginx/html

# 复制nginx配置
COPY deploy/nginx.conf /etc/nginx/nginx.conf
COPY deploy/default.conf /etc/nginx/conf.d/default.conf

# 创建日志目录
RUN mkdir -p /var/log/nginx && \
    chown -R nginx:nginx /var/log/nginx && \
    chown -R nginx:nginx /usr/share/nginx/html

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:80/health || exit 1

# 暴露端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]