# 构建阶段
FROM golang:1.22-alpine AS builder

# 安装基础工具
RUN apk add --no-cache git make build-base

# 设置工作目录
WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 go build -v -o bin/go-api-mono ./cmd/app

# 运行阶段
FROM alpine:latest

# 安装基础工具和CA证书
RUN apk --no-cache add ca-certificates tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# 创建非root用户
RUN adduser -D -g '' appuser

# 创建必要的目录
RUN mkdir -p /app/configs /app/logs \
    && chown -R appuser:appuser /app

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/bin/go-api-mono .
COPY --from=builder /app/configs ./configs

# 使用非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 设置健康检查
HEALTHCHECK --interval=30s --timeout=3s \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
ENTRYPOINT ["./go-api-mono"] 