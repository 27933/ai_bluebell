## 上下文

后端 API (Go/Gin) 已完全实现并通过集成测试。前端是一个全新项目，需要从零开始构建。设计规范文档（`FRONTEND_MVP_SPEC.md`）已编写，包含完整的功能需求、权限模型、交互逻辑、API映射等。

**现有资产**：
- 完整的 API 文档（`backend/docs/API_INTEGRATION.md`）
- HTML 原型演示（`frontend_prototype.html`）
- 设计规范（`FRONTEND_MVP_SPEC.md`）
- 后端已通过 API 测试

**约束**：
- 仅支持桌面版本（1200px+），移动端后续迭代
- MVP 阶段需在1-2周内完成核心功能
- 后端有2处修改需要同步进行

## 目标 / 非目标

**目标：**
1. 实现完整的 Vue 3 + TypeScript 前端应用
2. 支持所有MVP功能：浏览、搜索、登录、发布文章、评论、点赞等
3. 完整的 Markdown 编辑器（使用CodeMirror库）
4. 实时交互反馈（点赞、评论立即更新）
5. 统一的错误处理（Toast提示）
6. 跟后端一起部署

**非目标：**
- 管理员后台（后续迭代）
- 粉丝/关注系统（后续迭代）
- 移动端响应式（后续迭代）
- 深色主题（后续迭代）
- 导出/RSS功能（后续迭代）

## 决策

### 1. 技术栈选择

**决策**：Vue 3 + TypeScript + 现成UI组件库

**替代方案考虑**：
- React + TypeScript：更广泛的生态，但学习曲线更陡
- 纯 JavaScript：快速但类型不安全
- 其他框架（Angular、Svelte）：不在考虑范围

**理由**：
- Vue 3 Composition API 提供清晰的代码组织
- TypeScript 提供类型安全，减少运行时错误
- 现成UI组件库（如 Element Plus、Ant Design Vue）加快开发速度
- 与后端 Go 的简洁哲学保持一致

### 2. 状态管理

**决策**：使用 Pinia（Vue 3推荐的状态管理库）

**主要模块**：
- `auth`：用户登录状态、Token管理
- `article`：文章列表、详情、编辑状态
- `ui`：全局 UI 状态（加载、错误提示等）

**Token存储**：
- access_token 和 refresh_token 存储在 localStorage
- 自动刷新机制：收到 1007 错误时触发 refresh_token 流程

### 3. Markdown 编辑器

**决策**：使用 CodeMirror + marked.js

**具体方案**：
- CodeMirror 6：提供语法高亮、编辑辅助
- marked.js：将 Markdown 转换为 HTML
- 实时分屏预览（左编辑，右预览）
- 工具栏快捷键（加粗、标题、链接等）

**草稿保存**：
- 保存到浏览器 LocalStorage（不上传后端）
- 用户手动点击"保存草稿"或"发布"时才发送 API 请求

### 4. 导航与路由

**决策**：使用 Vue Router 4 实现 SPA 路由

**路由设计**：
```
/                          首页
/articles/:id              文章详情
/search                    搜索结果
/author/:username          作者主页
/tags                      标签页
/login                     登录
/signup                    注册
/dashboard                 作者仪表板（需认证）
/dashboard/articles        我的文章（需认证）
/dashboard/edit/:id        编辑文章（需认证）
/profile                   个人资料（需认证）
```

**权限检查**：
- 路由guard检查认证状态
- 需要登录的页面检查 token 有效性

### 5. 组件架构

**分层设计**：

```
views/
  ├─ Home.vue              首页
  ├─ ArticleDetail.vue     文章详情
  ├─ Login.vue             登录
  ├─ Signup.vue            注册
  ├─ Dashboard.vue         仪表板
  ├─ ArticleEditor.vue     文章编辑
  └─ ...

components/
  ├─ Navbar.vue            导航栏
  ├─ ArticleCard.vue       文章卡片
  ├─ CommentList.vue       评论列表
  ├─ MarkdownEditor.vue    Markdown编辑器
  └─ ...

services/
  ├─ api.ts               API 调用封装
  ├─ auth.ts              认证相关
  └─ ...

stores/
  ├─ auth.ts              用户状态
  ├─ article.ts           文章状态
  └─ ...
```

### 6. 错误处理

**决策**：统一使用 Toast 提示所有错误

**实现**：
- Pinia store 中维护一个全局错误状态
- 为 axios/fetch 配置拦截器，自动捕获错误
- Toast 自动隐藏或用户手动关闭
- HTTP 错误码映射到用户友好的消息

**特殊处理**：
- 1006（需要登录）：弹出提示并重定向登录
- 1007（Token 无效）：自动尝试刷新
- 1013（无权限）：显示"无权限"提示

### 7. 分页和搜索

**分页**：
- 默认每页 20 条
- 用户在首页/搜索页可选择每页数量（20/30/50）
- 保存到 URL query 参数便于分享

**搜索**：
- 简单的输入框 + 搜索按钮
- 进入搜索结果页面显示结果
- 不需要自动补全或搜索历史

### 8. 图表库

**决策**：使用 ECharts 绘制阅读趋势折线图

**配置**：
- 调用 `GET /article-stats/trend` API
- 参数：group_by = "day" 或 "week" 或 "month"
- 支持周统计和月统计切换

## 风险 / 权衡

### 风险1：Token 自动刷新的边界情况

**风险**：如果用户在请求过程中 token 过期，可能出现重试循环

**缓解**：
- 实现最多重试 1 次的机制
- 刷新失败后清空 token 强制跳转登录

### 风险2：Markdown 编辑器库体积

**风险**：CodeMirror + marked.js 可能增加较多包体积

**缓解**：
- 使用 Tree-shaking 移除未使用代码
- 后续可考虑使用轻量级替代方案

### 风险3：后端权限修改的时机

**风险**：前端依赖后端的权限修改，如果后端修改未同步可能功能异常

**缓解**：
- 并行开发，协调部署时机
- 前端做防御性编程，处理权限错误

### 权衡：桌面版本 Only

**权衡**：MVP 不支持移动端

**理由**：
- 优先完成核心功能
- 减少初期开发工作量
- 移动端响应式可后续增加

## 部署步骤

### 第1阶段：本地开发（1周）

1. 创建 Vue 3 + TS 项目
2. 配置路由、状态管理、UI 组件库
3. 实现页面和组件
4. 集成后端 API
5. 本地测试

### 第2阶段：后端修改同步（并行）

1. 修改 `GET /comments` 允许访客访问
2. 验证 `POST /articles` 权限检查
3. 测试 API 权限逻辑

### 第3阶段：构建和部署

1. 构建前端生产包
2. 将静态文件放到后端服务的 `public/` 或配置反向代理
3. 一起部署到服务器

## 技术决策表

| 决策项 | 选择 | 备选 | 理由 |
|--------|------|------|------|
| 框架 | Vue 3 | React | 框架明确性、学习曲线 |
| 语言 | TypeScript | JavaScript | 类型安全 |
| 状态管理 | Pinia | Vuex | Vue 3推荐 |
| 路由 | Vue Router 4 | 其他 | 官方库 |
| UI 组件库 | Element Plus | Ant Design Vue | 用户选择 |
| Markdown 编辑 | CodeMirror | Monaco | 体积和功能平衡 |
| 图表 | ECharts | Chart.js | 功能更丰富 |
| HTTP 客户端 | Axios | Fetch | 拦截器支持 |

## 待定问题

1. **UI 组件库的最终选择**：Element Plus 还是 Ant Design Vue？
2. **前端项目的具体目录结构**：是否需要特定的规范？
3. **CSS 预处理器**：使用 SCSS/Less 还是普通 CSS？
4. **构建工具**：使用 Vite 还是 Webpack？
5. **测试框架**：是否需要单元测试？覆盖率要求？

