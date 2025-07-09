# --- Build Stage ---
# *** 修改点：将 Go 版本从 1.21 升级到 1.23 ***
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum 文件并下载依赖
# 这一步可以利用 Docker 的层缓存，如果依赖没有变化，就不需要重新下载
COPY go.mod go.sum ./
RUN go mod download

# 复制所有剩余的源代码
COPY . .

# 编译应用，生成一个静态链接的二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /sub-node-cvt main.go

# --- Final Stage ---
# 使用一个非常小的基础镜像来减小最终镜像的体积
FROM alpine:latest

# 在容器中为应用创建一个非 root 用户，增加安全性
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /app

# 从构建阶段复制编译好的二进制文件和必要的静态资源
COPY --from=builder /sub-node-cvt .
COPY frontend ./frontend
COPY templates ./templates
COPY rulesets ./rulesets

# 暴露应用运行的端口
EXPOSE 8080

# 容器启动时运行的命令
CMD ["./sub-node-cvt"]