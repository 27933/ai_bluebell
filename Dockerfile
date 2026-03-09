FROM golang:1.21 AS builder

WORKDIR /build

# 先复制依赖文件，利用 Docker 缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bluebell .

# --- 运行阶段 ---
FROM debian:bookworm-slim

WORKDIR /app

# 安装必要的运行时依赖和等待工具
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates netcat-openbsd && \
    rm -rf /var/lib/apt/lists/*

# 从构建阶段复制二进制
COPY --from=builder /build/bluebell .
COPY --from=builder /build/wait-for.sh .
RUN chmod +x wait-for.sh

# 复制配置文件（运行时可通过挂载覆盖）
COPY conf/ ./conf/

EXPOSE 8084

# 默认使用 dev.yaml 配置启动，可通过 CMD 覆盖
CMD ["./bluebell", "conf/dev.yaml"]
