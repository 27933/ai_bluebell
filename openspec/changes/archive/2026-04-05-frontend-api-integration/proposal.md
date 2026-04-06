## 为什么

前端页面中的 API 调用存在多处问题：部分功能未接入后端接口（使用硬编码数据或 TODO 注释），部分接口调用参数与后端 API 不匹配，导致功能无法正常工作。需要修复这些问题以实现完整的前后端联调。

## 变更内容

### Home.vue
- 接入 `GET /articles/search` 实现搜索功能（当前只有 TODO 注释）
- 接入 `GET /tags` 获取真实标签列表（当前是硬编码数据）
- 修正排序参数映射（前端 `hot/latest/comments` → 后端 `popular/time`）

### ArticleDetail.vue
- 修复 author 字段显示（从 `article.author` 改为 `article.author.username/nickname`）
- 修复评论用户名显示（从 `comment.user_name` 改为 `comment.author.username`）

### Dashboard.vue
- 修复趋势图 API 调用（添加必填的 `article_id` 参数，移除不存在的 `time_range` 参数）

### Profile.vue
- 接入 `GET /auth/profile` 获取最新用户数据（当前只读取本地缓存）

## 功能 (Capabilities)

### 新增功能

（无新增功能，本次变更为修复现有功能的 API 集成问题）

### 修改功能

（无规范级别的需求变更，仅修复实现层面的 API 调用问题）

## 影响

- **前端文件**：
  - `frontend/src/views/Home.vue`
  - `frontend/src/views/ArticleDetail.vue`
  - `frontend/src/views/Dashboard.vue`
  - `frontend/src/views/Profile.vue`

- **API 依赖**：
  - `GET /articles/search`
  - `GET /tags`
  - `GET /article-stats/trend`
  - `GET /auth/profile`

- **无破坏性变更**：所有修改都是修复现有功能，不影响其他模块
