# Bluebell 知识博客平台

> 一个由 AI 全程辅助开发的全栈博客平台——从需求分析、架构设计到前后端代码，均通过与 Claude AI 对话协作完成。

## 📊 AI 代码生成数据

| 指标 | 数据 |
|------|------|
| **AI 代码生成比例** | **~100%** |
| 后端代码（Go） | 14,528 行 / 69 个文件 |
| 前端代码（Vue/TS） | 6,510 行 / 25 个文件 |
| **代码总量** | **~21,000 行** |
| Git 提交次数 | 24 次 |
| 开发周期 | 约 1 周 |
| 人工手写代码 | 接近 0 行 |

> 开发者全程通过自然语言与 Claude Code 对话，**不直接编写任何业务代码**。
> 从数据库表结构、后端 API、前端页面到 Docker 部署配置，全部由 AI 生成并迭代优化。

---

## 🤖 AI 辅助开发

本项目是 **人机协作开发** 的实践案例。开发者通过自然语言描述需求，由 [Claude Code](https://claude.ai/code) 负责生成、调试和迭代全部代码。

### Claude 完成的工作

```
✅ 数据库表结构设计（articles / users / comments / likes / tags 等 10+ 张表）
✅ 后端 RESTful API 全部接口（约 50 个端点）
✅ JWT 双 Token 认证体系
✅ 前端 Vue 3 全部页面（首页、文章详情、编辑器、仪表板、管理后台等 15 个页面）
✅ Axios 拦截器 + Token 自动刷新
✅ ECharts 数据可视化图表
✅ Docker / Docker Compose / Nginx 部署配置
✅ 安全加固（限流、MIME 检测、XSS 防护、JWT Secret 迁移）
✅ Bug 诊断和修复（gin.H int64 类型问题、路由冲突、跨域等）
✅ 生产环境部署排错
```

### 开发者的角色

```
开发者：描述功能需求、审查代码、测试验证、做产品决策
Claude：架构设计、代码生成、Bug 修复、安全加固、部署配置
```

### 典型对话示例

| 开发者输入 | Claude 输出 |
|-----------|------------|
| "新建用户需要风控，不然服务器被打挂" | 实现 per-IP 令牌桶限速中间件，注册 10s/次，登录 3s/次 |
| "图片上传有安全漏洞吗？" | 发现仅校验扩展名的问题，补充文件头魔数检测 |
| "JWT secret 在源码里，帮我移出来" | 重构 jwt.go，迁移至 viper config，生成随机密钥，更新 gitignore |
| "精选标识有点丑，优化一下，加个精选筛选" | 重设计圆形角标 + 金色边框，新增前后端联动筛选 |

### 开发工作流

```
用户描述需求
    ↓
Claude 探索代码库（/opsx:explore）
    ↓
提案 + 设计文档（OpenSpec）
    ↓
Claude 生成代码 / 用户审查
    ↓
本地测试验证
    ↓
Claude 修复问题
    ↓
提交 & 部署
```

---

## 🏗 项目架构

```
┌─────────────────────────────────────────────────────────────┐
│                         用户浏览器                           │
└──────────────────┬──────────────────────────────────────────┘
                   │ HTTP / HTTPS
┌──────────────────▼──────────────────────────────────────────┐
│               Nginx（反向代理 + 静态资源）                    │
│    /          → Vue 3 前端静态文件                           │
│    /api/v1/*  → 后端 :8084                                  │
└──────────────────┬──────────────────────────────────────────┘
                   │
┌──────────────────▼──────────────────────────────────────────┐
│              Go 后端（Gin 框架，:8084）                       │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────────────┐  │
│  │controller│→│  logic   │→│   dao    │→│  MySQL/Redis  │  │
│  └──────────┘ └──────────┘ └──────────┘ └───────────────┘  │
│                                                              │
│  中间件：JWT 认证 · per-IP 限流 · CORS · pprof               │
└─────────────────────────────────────────────────────────────┘
```

### 目录结构

```
ai_bluebell/
├── backend/                  # Go 后端
│   ├── controller/           # HTTP 处理层
│   ├── logic/                # 业务逻辑层
│   ├── dao/mysql/            # 数据访问层
│   ├── models/               # 数据模型
│   ├── middlewares/          # 中间件（JWT、限流）
│   ├── pkg/jwt/              # JWT 工具
│   ├── router/               # 路由注册
│   └── conf/                 # 配置文件（dev / cloud）
├── frontend/                 # Vue 3 前端
│   └── src/
│       ├── views/            # 页面组件
│       │   └── admin/        # 管理后台页面
│       ├── components/       # 公共组件
│       ├── stores/           # Pinia 状态管理
│       ├── services/         # Axios 封装
│       └── router/           # 路由配置
├── docker-compose.yml        # 本地开发环境
├── docker-compose.cloud.yml  # 生产部署
└── openspec/                 # AI 协作设计文档
```

---

## ✨ 功能特性

### 内容管理
- **文章系统**：Markdown 编辑器（CodeMirror）、实时预览、草稿自动保存
- **精选文章**：管理员设置精选，首页角标展示 + 一键筛选
- **标签系统**：多标签管理，标签页浏览、筛选
- **RSS 订阅**：标准 RSS 2.0 XML，支持所有 RSS 阅读器

### 用户系统
- JWT 双 Token 认证（Access + Refresh，自动续签）
- 角色权限：`visitor` / `reader` / `author` / `admin`
- 用户注册、登录、个人资料编辑

### 互动功能
- 文章点赞（幂等处理）
- 评论系统（发表、回复、编辑、删除，墓碑删除模式）
- 文章浏览统计、热度排行

### 管理后台
- 用户管理（角色修改、封禁/解封、批量操作）
- 文章管理（状态控制、精选设置）
- 数据统计（概览卡片 + ECharts 趋势图）

### 安全
- per-IP 请求限流（注册 10s/次，登录 3s/次）
- 图片上传文件头魔数验证（防扩展名伪装）
- XSS 防护（DOMPurify sanitize Markdown 输出）
- JWT Secret 外置配置（不进代码仓库）
- 密码最小长度 8 位

---

## 🛠 技术栈

| 层级 | 技术 |
|------|------|
| 前端框架 | Vue 3 + TypeScript + Vite |
| UI 组件 | Element Plus + Bootstrap Icons |
| 状态管理 | Pinia |
| HTTP 客户端 | Axios（含 Token 自动刷新拦截器） |
| Markdown | marked + DOMPurify |
| 图表 | ECharts |
| 编辑器 | CodeMirror 6 |
| 后端框架 | Go 1.21 + Gin |
| 数据库 | MySQL 8.0 |
| 缓存 | Redis |
| ORM | sqlx（参数化查询，防 SQL 注入） |
| 认证 | JWT（dgrijalva/jwt-go） |
| 配置 | Viper |
| 日志 | Zap |
| 部署 | Docker + Docker Compose + Nginx |

---

## 🚀 本地开发

### 前置条件
- Docker & Docker Compose
- Node.js 18+
- Go 1.21+

### 启动后端

```bash
# 启动 MySQL + Redis（首次运行）
docker run -d --name mysql8 --network host --restart always \
  -e MYSQL_ROOT_PASSWORD=your_password \
  -e MYSQL_DATABASE=bluebell \
  mysql:8.0 --default-authentication-plugin=mysql_native_password \
  --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

docker run -d --name redis --network host --restart always \
  redis:latest redis-server --requirepass your_password

# 启动后端（源码热编译）
docker run -d --name bluebell-backend \
  --network host \
  -v $(pwd)/backend:/app \
  golang:1.21 \
  sh -c 'cd /app && go build -o bluebell . && ./bluebell conf/dev.yaml'

# 修改代码后重启
docker restart bluebell-backend
```

### 启动前端

```bash
cd frontend
npm install
npm run dev
# 访问 http://localhost:5173
```

### 配置文件

复制模板并填入实际密码和 JWT Secret：

```bash
cp backend/conf/cloud.yaml.example backend/conf/cloud.yaml
# 编辑 cloud.yaml，填入数据库密码和随机 JWT Secret
# 生成 Secret：openssl rand -hex 32
```

---

## 📦 生产部署

```bash
# 1. 克隆代码
git clone https://github.com/your-username/ai_bluebell.git
cd ai_bluebell

# 2. 创建生产配置（参考 backend/conf/cloud.yaml.example）
cp backend/conf/cloud.yaml.example backend/conf/cloud.yaml
# 编辑 cloud.yaml 填入真实密码和 JWT Secret

# 3. 构建并启动
docker-compose -f docker-compose.cloud.yml up -d --build

# 4. 验证
curl http://127.0.0.1:8084/ping  # 返回 pong 即正常
```

---

## 📄 开源协议

MIT License
