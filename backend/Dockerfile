FROM golang:1.21-alpine AS builder

WORKDIR /build

# 先复制依赖文件，利用 Docker 缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bluebell .

# --- 运行阶段 ---
FROM golang:1.21-alpine

WORKDIR /app

# 从构建阶段复制二进制
COPY --from=builder /build/bluebell .

# 复制配置文件（运行时可通过挂载覆盖）
COPY conf/ ./conf/


EXPOSE 8084

# 默认使用 local.yaml 配置启动，可通过 CMD 覆盖
CMD ["./bluebell", "conf/local.yaml"]
