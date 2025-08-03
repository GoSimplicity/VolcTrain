# 构建阶段
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git ca-certificates tzdata

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download && go mod verify

# 复制源代码
COPY . .

# 构建应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/common ./cmd/common/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/monitoring ./cmd/monitoring/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/training ./cmd/training/main.go

# 生产阶段
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata mysql-client curl

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 创建非root用户
RUN addgroup -g 1001 volctrain && \
    adduser -D -u 1001 -G volctrain volctrain

# 设置工作目录
WORKDIR /home/volctrain

# 从构建阶段复制二进制文件
COPY --from=builder /app/bin/ ./bin/
COPY --from=builder /app/etc/ ./etc/
COPY --from=builder /app/sql/ ./sql/

# 复制启动脚本
COPY scripts/docker-entrypoint.sh ./
RUN chmod +x docker-entrypoint.sh

# 更改文件所有权
RUN chown -R volctrain:volctrain /home/volctrain

# 切换到非root用户
USER volctrain

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8888/health || exit 1

# 暴露端口
EXPOSE 8888 8889 8890 8891 9090

# 设置启动命令
ENTRYPOINT ["./docker-entrypoint.sh"]
CMD ["api"]