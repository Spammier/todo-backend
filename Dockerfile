# ---- 构建阶段 ----
# 使用与 go.mod 中指定的 Go 版本匹配的官方 Go 镜像
FROM golang:1.23-alpine as builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件，并下载依赖
# 这样可以利用 Docker 的层缓存机制，在依赖不变的情况下加快构建速度
COPY go.mod go.sum ./
RUN go mod download

# 复制所有源代码
COPY . .

# 编译 Go 应用
# -o /app/server 指定输出文件路径和名称
# CGO_ENABLED=0 禁用 CGO，生成静态链接的可执行文件，更适合容器环境
# ./cmd/api/main.go 是你的主程序入口
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/server ./cmd/api/main.go

# ---- 运行阶段 ----
# 使用一个轻量级的 Alpine 镜像作为最终运行环境
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制编译好的可执行文件
COPY --from=builder /app/server .

# (可选) 如果你的应用需要读取 .env 文件或其他配置文件，也需要复制
# COPY .env .env
# 注意：在 Docker 中更推荐使用环境变量来传递配置，而不是打包 .env 文件

# 暴露应用程序监听的端口 (根据你的 main.go，默认为 8080，或由 PORT 环境变量指定)
# 这里我们先暴露 8080，你可以在运行时通过 -e PORT=xxx 修改
EXPOSE 8080

# 运行 Go 应用
CMD ["./server"] 