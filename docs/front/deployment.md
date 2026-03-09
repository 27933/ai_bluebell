# Bluebell 容器化部署文档

## 项目架构

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  bluebell   │────→│   MySQL 8   │     │  Redis 7    │
│  (Go/Gin)   │────→│  Port 3306  │     │  Port 6379  │
│  Port 8084  │     └─────────────┘     └─────────────┘
└─────────────┘
```

## 快速启动

### 前置要求

- Docker >= 20.10
- Docker Compose >= 2.0

### 生产部署

```bash
# 1. 克隆项目
git clone <repo-url> bluebell && cd bluebell

# 2. 修改配置文件中的 MySQL/Redis 地址
#    生产环境需将 conf/dev.yaml 中的 host 改为 docker-compose 中的服务名
cp conf/dev.yaml conf/prod.yaml
# 编辑 prod.yaml：
#   mysql.host: "mysql"
#   redis.host: "redis"

# 3. 构建并启动
docker-compose up -d --build

# 4. 查看日志
docker-compose logs -f app

# 5. 验证服务
curl http://localhost:8084/ping
# 应返回: pong
```

### 开发/测试环境

开发模式会启动一个挂载源码的 Go 容器，适合运行测试：

```bash
# 启动含开发容器的环境
docker-compose --profile dev up -d

# 进入开发容器
docker exec -it bluebell-ai bash

# 在容器内编译运行
cd /app
go build -o bluebell . && ./bluebell conf/dev.yaml &

# 运行集成测试
go test ./tests/api/ -v -count=1
```

## Docker 配置说明

### Dockerfile（多阶段构建）

```
阶段 1: 构建（golang:1.21）
  → go mod download（缓存依赖）
  → go build（编译为静态二进制）

阶段 2: 运行（debian:bookworm-slim）
  → 仅包含二进制 + 配置文件
  → 镜像体积约 50MB
```

### docker-compose.yml 服务说明

| 服务 | 镜像 | 端口 | 说明 |
|------|------|------|------|
| mysql | mysql:8.0 | 13306:3306 | 数据库，数据持久化到 volume |
| redis | redis:7-alpine | 16379:6379 | 缓存/统计，数据持久化到 volume |
| app | 自建镜像 | 8084:8084 | 生产应用服务 |
| dev | golang:1.21 | - | 开发容器（需 `--profile dev`） |

### 配置文件（conf/dev.yaml）

```yaml
name: "bluebell"
mode: "dev"               # dev/release
port: 8084                 # 应用端口

auth:
  jwt_expire: 8            # Access Token 有效期（小时）

log:
  level: "info"            # debug/info/warn/error
  filename: "web_app.log"
  max_size: 200            # 日志文件最大 MB
  max_age: 30              # 日志保留天数
  max_backups: 7           # 最大备份数

mysql:
  host: "mysql"            # docker-compose 服务名 或 IP 地址
  port: 3306
  user: "root"
  password: "123456"
  dbname: "bluebell"
  max_open_conns: 200
  max_idle_conns: 50

redis:
  host: "redis"            # docker-compose 服务名 或 IP 地址
  port: 6379
  password: "123456"
  db: 6
  pool_size: 100
```

> **生产环境注意事项：**
> - 务必修改 MySQL 和 Redis 密码
> - 将 `mode` 设置为 `release`
> - 根据实际情况调整连接池参数

## 健康检查

```bash
# 应用健康检查
curl http://localhost:8084/ping

# MySQL 健康检查
docker exec bluebell-mysql mysqladmin ping -h localhost

# Redis 健康检查
docker exec bluebell-redis redis-cli -a 123456 ping
```

## 常用运维命令

```bash
# 查看所有服务状态
docker-compose ps

# 重启应用
docker-compose restart app

# 查看应用日志
docker-compose logs -f app

# 停止所有服务
docker-compose down

# 停止并清除数据卷（慎用）
docker-compose down -v

# 重新构建应用镜像
docker-compose build app

# 进入 MySQL
docker exec -it bluebell-mysql mysql -uroot -p123456 bluebell

# 进入 Redis
docker exec -it bluebell-redis redis-cli -a 123456
```

## API 访问地址

| 环境 | 地址 |
|------|------|
| 本地开发 | http://localhost:8084/api/v1 |
| Swagger 文档 | http://localhost:8084/swagger/index.html |
| 健康检查 | http://localhost:8084/ping |
