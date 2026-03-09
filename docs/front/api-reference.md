# Bluebell API 接口文档

> 供前端开发参考，基于集成测试验证通过的接口（84/84 PASS）

## 基础信息

- **Base URL**: `http://<host>:8084/api/v1`
- **内容类型**: `application/json`
- **认证方式**: JWT Bearer Token，通过 `Authorization: Bearer <access_token>` 请求头传递

## 统一响应格式

所有接口返回统一的 JSON 结构：

```json
{
  "code": 1000,
  "msg": "success",
  "data": { ... }
}
```

## 错误码一览

| 错误码 | 含义 | 前端处理建议 |
|--------|------|-------------|
| 1000 | 成功 | 正常处理 data |
| 1001 | 请求参数错误 | 提示用户检查输入 |
| 1002 | 用户名已存在 | 注册时提示换用户名 |
| 1003 | 用户名不存在 | 登录时提示 |
| 1004 | 用户名或密码错误 | 登录时提示 |
| 1005 | 服务繁忙 | 提示稍后重试 |
| 1006 | 需要登录 | 跳转登录页/弹出登录框 |
| 1007 | 无效的 Token | 清除本地 Token，跳转登录 |
| 1008 | 栏目不存在 | 提示用户 |
| 1009 | 栏目已存在 | 提示用户换名称 |
| 1013 | 没有权限 | 提示权限不足 |
| 1014 | 文章不存在 | 提示用户 |

## ID 说明

所有实体 ID（文章、用户、评论、标签等）均使用 Snowflake 算法生成的 `int64`，**在 JSON 中以字符串形式返回**以避免 JavaScript 精度丢失。

```json
{ "id": "1893668041289216000" }
```

前端存储和传递 ID 时应使用**字符串类型**。

---

## 1. 认证接口（无需登录）

### 1.1 用户注册

```
POST /auth/signup
```

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |
| re_password | string | 是 | 确认密码，必须与 password 一致 |

```json
{
  "username": "zhangsan",
  "password": "test123456",
  "re_password": "test123456"
}
```

**成功响应：** `data: null`

**可能的错误码：** 1001, 1002

---

### 1.2 用户登录

```
POST /auth/login
```

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "user": {
      "id": "1893668041289216000",
      "username": "zhangsan",
      "email": "",
      "role": "reader",
      "status": "active"
    },
    "token": {
      "access_token": "eyJhbGciOi...",
      "refresh_token": "eyJhbGciOi...",
      "expires_in": 28800
    }
  }
}
```

> **前端处理：** 登录成功后保存 `access_token` 和 `refresh_token`。`access_token` 有效期 8 小时。

**可能的错误码：** 1001, 1004

---

### 1.3 刷新 Token

```
POST /auth/refresh
```

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| refresh_token | string | 是 | 登录时获取的 refresh_token |

**成功响应：**

```json
{
  "data": {
    "access_token": "eyJhbGciOi...",
    "refresh_token": "eyJhbGciOi...",
    "expires_in": 28800
  }
}
```

> **前端处理：** 当 access_token 过期（收到 1007 错误）时，用 refresh_token 换取新 Token。如果刷新也失败，则跳转登录页。

**可能的错误码：** 1001, 1007

---

## 2. 用户资料（需要登录）

### 2.1 获取个人资料

```
GET /auth/profile
Authorization: Bearer <access_token>
```

**成功响应：**

```json
{
  "data": {
    "id": "1893668041289216000",
    "username": "zhangsan",
    "email": "test@example.com",
    "role": "reader",
    "status": "active",
    "nickname": "张三",
    "avatar": "/uploads/images/avatar.png",
    "bio": "这是个人简介",
    "total_words": 15000,
    "total_likes": 42,
    "last_login_at": "2026-03-09T10:00:00Z",
    "created_at": "2026-01-01T00:00:00Z"
  }
}
```

---

### 2.2 更新个人资料

```
PUT /auth/profile
Authorization: Bearer <access_token>
```

**请求体（所有字段可选）：**

| 字段 | 类型 | 限制 | 说明 |
|------|------|------|------|
| nickname | string | 1-100 字符 | 昵称 |
| email | string | 合法邮箱，最长 100 | 邮箱 |
| bio | string | 最长 500 | 个人简介 |
| avatar | string | 最长 500 | 头像 URL |

**成功响应：** `data: null`

---

## 3. 文章接口（公开）

### 3.1 获取文章列表

```
GET /articles?page=1&size=20&sort=time
```

**查询参数：**

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码（从 1 开始） |
| size | int | 20 | 每页条数（1-50） |
| sort | string | "time" | 排序方式：`time`（时间）/ `popular`（热度） |
| status | string | "published" | 状态筛选：`draft` / `published` / `offline` / `all` |
| tag | string | - | 按标签筛选 |
| author_id | string | - | 按作者 ID 筛选 |
| keyword | string | - | 搜索关键词 |
| days | int | - | 天数范围（1-365） |

**成功响应：**

```json
{
  "data": {
    "list": [
      {
        "id": "1893668041289216000",
        "title": "文章标题",
        "summary": "文章摘要...",
        "author": {
          "id": "1893668041289216000",
          "username": "zhangsan",
          "nickname": "张三",
          "avatar": "/uploads/images/avatar.png"
        },
        "tags": [
          { "id": "1893668041289216001", "name": "Go" }
        ],
        "view_count": 100,
        "like_count": 20,
        "comment_count": 5,
        "is_featured": false,
        "is_recent": true,
        "created_at": "2026-03-01T10:00:00Z",
        "updated_at": "2026-03-08T15:30:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "size": 20,
    "pages": 5
  }
}
```

---

### 3.2 获取文章详情

```
GET /articles/:id
```

**路径参数：** `id` — 文章 ID

**成功响应：**

```json
{
  "data": {
    "id": "1893668041289216000",
    "title": "文章标题",
    "content": "文章完整内容（Markdown）...",
    "summary": "摘要",
    "word_count": 3500,
    "author": {
      "id": "1893668041289216000",
      "username": "zhangsan",
      "nickname": "张三",
      "avatar": "/uploads/images/avatar.png",
      "bio": "个人简介"
    },
    "tags": [{ "id": "...", "name": "Go" }],
    "status": "published",
    "is_featured": false,
    "allow_comment": true,
    "like_count": 20,
    "comment_count": 5,
    "view_count": 100,
    "slug": "article-slug",
    "meta_keywords": "Go,编程",
    "meta_description": "SEO 描述",
    "created_at": "2026-03-01T10:00:00Z",
    "updated_at": "2026-03-08T15:30:00Z"
  }
}
```

> **副作用：** 访问详情会自动记录浏览量（含防刷机制）。

---

### 3.3 获取精选文章

```
GET /articles/featured?limit=3
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| limit | int | 3 | 返回数量（最大 10） |

**成功响应：** `data: [Article]`（文章数组）

---

### 3.4 搜索文章

```
GET /articles/search?keyword=Go&page=1&size=20
```

参数同文章列表，至少需要提供 `keyword`、`author_id`、`author_name` 或 `tag` 之一。

**成功响应：** 同文章列表分页格式。

---

### 3.5 热门文章排行

```
GET /articles/trending?period=daily&page=1&size=20
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| period | string | "daily" | 排行周期：`daily` / `weekly` / `monthly` |
| page | int | 1 | 页码 |
| size | int | 20 | 每页条数（最大 100） |

**成功响应：**

```json
{
  "data": {
    "list": [
      {
        "article_id": "...",
        "title": "热门文章",
        "view_count": 500,
        "unique_view_count": 320,
        "rank_position": 1,
        "author_username": "zhangsan"
      }
    ],
    "page": 1,
    "size": 20
  }
}
```

---

## 4. 文章统计（公开）

### 4.1 文章日统计

```
GET /article-stats/daily?article_id=123&days=30
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| article_id | int | 必填 | 文章 ID |
| days | int | 30 | 天数（最大 90） |

---

### 4.2 文章访问趋势

```
GET /article-stats/trend?article_id=123&days=30&group_by=day
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| article_id | int | 必填 | 文章 ID |
| days | int | 30 | 天数（1-90） |
| group_by | string | 必填 | 分组方式：`hour` / `day` / `week` / `month` |

**成功响应：**

```json
{
  "data": {
    "article_id": 123,
    "days": 30,
    "group_by": "day",
    "trend": [
      { "label": "2026-03-01", "value": 42 },
      { "label": "2026-03-02", "value": 58 }
    ]
  }
}
```

---

### 4.3 批量文章统计

```
GET /article-stats/batch?ids=1,2,3
```

| 参数 | 类型 | 说明 |
|------|------|------|
| ids | string | 逗号分隔的文章 ID（最多 100 个） |

---

### 4.4 记录文章浏览（防刷）

```
POST /articles/view?article_id=123
```

---

## 5. 作者写作接口（需要登录）

### 5.1 创建文章

```
POST /articles
Authorization: Bearer <access_token>
```

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 标题（1-200 字符） |
| content | string | 是 | 内容（Markdown） |
| tags | string[] | 否 | 标签名称数组 |
| status | string | 否 | 状态：`draft`（默认）/ `published` / `offline` |
| allow_comment | bool | 否 | 是否允许评论 |
| slug | string | 否 | URL slug（最长 200） |
| meta_keywords | string | 否 | SEO 关键词（最长 500） |
| meta_description | string | 否 | SEO 描述（最长 500） |

```json
{
  "title": "我的第一篇文章",
  "content": "# Hello\n这是文章内容...",
  "tags": ["Go", "教程"],
  "status": "draft",
  "allow_comment": true
}
```

**成功响应：** `data: Article`（返回创建的文章对象，含 `id`）

---

### 5.2 更新文章

```
PUT /author/articles/:id
Authorization: Bearer <access_token>
```

> **权限：** 仅文章作者可更新。

**请求体（所有字段可选）：**

| 字段 | 类型 | 说明 |
|------|------|------|
| title | string | 标题 |
| content | string | 内容 |
| tags | string[] | 标签 |
| summary | string | 摘要（最长 500） |
| status | string | 状态 |
| allow_comment | bool | 是否允许评论 |
| slug | string | URL slug |

**成功响应：** `data: null`

---

### 5.3 删除文章

```
DELETE /author/articles/:id
Authorization: Bearer <access_token>
```

> **权限：** 仅文章作者可删除。

**成功响应：** `data: null`

---

### 5.4 更新文章状态

```
PATCH /author/articles/:id/status
Authorization: Bearer <access_token>
```

**请求体：**

```json
{ "status": "published" }
```

可选值：`draft` / `published` / `offline`

---

### 5.5 设置文章精选

```
PATCH /author/articles/:id/featured
Authorization: Bearer <access_token>
```

**请求体：**

```json
{ "is_featured": true }
```

---

### 5.6 获取我的文章列表

```
GET /author/articles?page=1&size=20&status=all
Authorization: Bearer <access_token>
```

**成功响应：** 同文章列表分页格式（包含所有状态的文章）。

---

### 5.7 导出文章

```
GET /author/articles/:id/export
Authorization: Bearer <access_token>
```

**成功响应：**

```json
{
  "data": {
    "filename": "article-title.md",
    "content": "# Title\n...",
    "size": 1234
  }
}
```

---

### 5.8 批量导出文章

```
POST /author/articles/export
Authorization: Bearer <access_token>
```

**请求体：**

```json
{ "article_ids": [1, 2, 3] }
```

（1-50 个 ID）

**成功响应：**

```json
{
  "data": {
    "batch_id": "uuid",
    "files": [{ "filename": "...", "content": "...", "size": 0 }],
    "file_count": 3,
    "total_size": 5000
  }
}
```

---

## 6. 评论接口（需要登录）

### 6.1 创建评论

```
POST /comments
Authorization: Bearer <access_token>
```

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| article_id | int64 | 是 | 文章 ID |
| parent_id | int64 | 否 | 父评论 ID（回复某评论时填） |
| content | string | 是 | 评论内容（1-1000 字符） |

**成功响应：**

```json
{
  "data": {
    "id": "1893668041289216000",
    "content": "写得真好！",
    "like_count": 0,
    "status": "active",
    "author": {
      "id": "...",
      "username": "zhangsan",
      "nickname": "张三",
      "avatar": "/uploads/images/avatar.png"
    },
    "parent_id": null,
    "created_at": "2026-03-09T10:00:00Z"
  }
}
```

---

### 6.2 获取评论列表

```
GET /comments?article_id=123&page=1&size=20
Authorization: Bearer <access_token>
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| article_id | int64 | 是 | 文章 ID |
| page | int | 否 | 页码（默认 1） |
| size | int | 否 | 每页条数（默认 20，最大 50） |

**成功响应：** 分页格式，`list` 为 `ApiComment` 数组。

---

### 6.3 更新评论

```
PUT /comments/:id
Authorization: Bearer <access_token>
```

> **权限：** 仅评论作者可更新。

**请求体：**

```json
{ "content": "更新后的评论内容" }
```

---

### 6.4 删除评论

```
DELETE /comments/:id
Authorization: Bearer <access_token>
```

> **权限：** 仅评论作者可删除。

---

## 7. 点赞接口（需要登录）

### 7.1 点赞

```
POST /likes
Authorization: Bearer <access_token>
```

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| target_type | string | 是 | 目标类型：`article` / `comment` |
| target_id | int64 | 是 | 目标 ID |

> 重复点赞为幂等操作，不会报错。

---

### 7.2 取消点赞

```
DELETE /likes?target_type=article&target_id=123
Authorization: Bearer <access_token>
```

> 取消未点赞的内容也不会报错。

---

### 7.3 获取点赞状态

```
GET /likes/status?target_type=article&target_id=123
Authorization: Bearer <access_token>
```

**成功响应：**

```json
{
  "data": {
    "is_liked": true,
    "like_count": 42,
    "created_at": "2026-03-09T10:00:00Z"
  }
}
```

---

### 7.4 获取用户点赞列表

```
GET /user/likes?target_type=article&page=1&size=20
Authorization: Bearer <access_token>
```

---

## 8. 标签接口

### 8.1 获取所有标签（公开）

```
GET /tags
```

**成功响应：**

```json
{
  "data": [
    {
      "id": "1893668041289216000",
      "name": "Go",
      "description": "Go 语言相关",
      "slug": "go",
      "article_count": 15,
      "created_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

> **注意：** 标签列表直接返回数组，不是包裹在 `list` 中。

---

### 8.2 获取标签下文章（公开）

```
GET /tags/:id/articles?page=1&size=20
```

**成功响应：** 分页文章列表。

---

### 8.3 创建标签（需要登录）

```
POST /tags
Authorization: Bearer <access_token>
```

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 标签名称（1-50 字符） |
| description | string | 否 | 描述（最长 200） |
| slug | string | 否 | URL slug（最长 50） |

**成功响应：** `data: Tag`（返回创建的标签对象）

---

### 8.4 更新标签（需要登录）

```
PUT /tags/:id
Authorization: Bearer <access_token>
```

---

### 8.5 删除标签（需要登录）

```
DELETE /tags/:id
Authorization: Bearer <access_token>
```

---

### 8.6 获取作者标签（需要登录）

```
GET /author/tags
Authorization: Bearer <access_token>
```

返回当前登录用户使用过的标签列表。

---

## 9. 栏目接口

### 9.1 获取栏目列表（公开）

```
GET /categories
```

**成功响应：**

```json
{
  "data": {
    "list": [
      {
        "id": 1,
        "category_name": "技术",
        "introduction": "技术类文章",
        "created_by": 1,
        "created_at": "2026-01-01T00:00:00Z",
        "article_count": 30
      }
    ]
  }
}
```

---

### 9.2 获取栏目详情（公开）

```
GET /categories/:id
```

---

### 9.3 获取栏目下文章（公开）

```
GET /categories/:id/articles?page=1&size=20
```

---

### 9.4 获取文章所属栏目（公开）

```
GET /articles/:id/categories
```

**成功响应：** `data: { "list": [Category] }`

---

### 9.5 创建栏目（需要登录）

```
POST /categories
Authorization: Bearer <access_token>
```

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| category_name | string | 是 | 栏目名称（2-128 字符） |
| introduction | string | 是 | 栏目介绍（2-256 字符） |

---

### 9.6 更新栏目（需要登录）

```
PUT /categories/:id
Authorization: Bearer <access_token>
```

> **权限：** 仅栏目创建者或管理员可更新。

**请求体：**

| 字段 | 类型 | 说明 |
|------|------|------|
| category_name | string | 栏目名称（2-128 字符） |
| introduction | string | 栏目介绍（2-256 字符） |

---

### 9.7 删除栏目（需要登录）

```
DELETE /categories/:id
Authorization: Bearer <access_token>
```

> **权限：** 仅栏目创建者或管理员可删除。

---

### 9.8 添加文章到栏目（需要登录）

```
POST /articles/:id/categories
Authorization: Bearer <access_token>
```

**请求体：**

```json
{ "category_ids": [1, 2, 3] }
```

> **权限：** 仅文章作者可操作。

---

## 10. 作者主页（公开）

### 10.1 获取作者信息

```
GET /authors/:username
```

**成功响应：**

```json
{
  "data": {
    "username": "zhangsan",
    "nickname": "张三",
    "bio": "Go 开发者",
    "join_date": "2026-01-01",
    "article_count": 30,
    "total_views": 5000,
    "total_likes": 200
  }
}
```

---

### 10.2 获取作者文章列表

```
GET /authors/:username/articles?page=1&size=20&sort=time
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码 |
| size | int | 20 | 每页条数（最大 100） |
| sort | string | "time" | `time` 或 `hot` |

---

## 11. 文件上传（需要登录）

### 11.1 上传图片

```
POST /upload/image
Authorization: Bearer <access_token>
Content-Type: multipart/form-data
```

表单字段名：`file`

**成功响应：**

```json
{
  "data": {
    "url": "/uploads/images/1234567890.png",
    "filename": "1234567890.png",
    "size": 102400
  }
}
```

---

### 11.2 上传附件

```
POST /upload/attachment
Authorization: Bearer <access_token>
Content-Type: multipart/form-data
```

---

## 12. RSS 订阅（公开）

```
GET /rss
```

返回 JSON 格式的 RSS 数据。

---

## 13. 管理员接口（需要管理员权限）

> 所有管理员接口需要 JWT + admin 角色。普通用户访问会返回权限错误。

### 13.1 管理员获取文章列表

```
GET /admin/articles?page=1&size=20
```

---

### 13.2 设置文章精选

```
PATCH /admin/articles/:id/featured
```

**请求体：** `{ "is_featured": true }`

---

### 13.3 获取用户列表

```
GET /admin/users?page=1&size=20&role=all&status=all
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| role | string | "all" | `all` / `visitor` / `reader` / `author` / `admin` |
| status | string | "all" | `all` / `active` / `inactive` |
| page | int | 1 | 页码 |
| size | int | 20 | 每页条数（最大 50） |

---

### 13.4 获取用户详情

```
GET /admin/users/:id
```

---

### 13.5 更新用户状态

```
PATCH /admin/users/:id/status
```

**请求体：** `{ "status": "active" }` 或 `{ "status": "inactive" }`

---

### 13.6 更新用户角色

```
PATCH /admin/users/:id/role
```

**请求体：** `{ "role": "admin" }`

可选值：`admin` / `author` / `reader`

---

### 13.7 批量更新用户状态

```
PATCH /admin/users/batch/status
```

**请求体：**

```json
{
  "user_ids": [1, 2, 3],
  "status": "inactive"
}
```

（1-100 个用户 ID）

---

### 13.8 系统概览

```
GET /admin/stats/overview
```

**成功响应：**

```json
{
  "data": {
    "user_count": 100,
    "article_count": 500,
    "comment_count": 2000,
    "today_new_user_count": 5,
    "today_new_article_count": 10
  }
}
```

---

### 13.9 系统日统计

```
GET /admin/stats/daily?start_date=2026-01-01&end_date=2026-01-31
```

---

### 13.10 系统实时性能指标

```
GET /admin/metrics/realtime
```

**成功响应：**

```json
{
  "data": {
    "cpu_usage": 45.2,
    "memory_usage": 67.8,
    "disk_usage": 55.0,
    "active_users": 120,
    "timestamp": "2026-03-09T10:00:00Z"
  }
}
```

---

### 13.11 系统性能历史

```
GET /admin/metrics/history?start_time=1709900000&end_time=1709903600&metric_type=cpu
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| start_time | int64 | 当前时间 - 1 小时 | Unix 时间戳（秒） |
| end_time | int64 | 当前时间 | Unix 时间戳（秒） |
| metric_type | string | "cpu" | `cpu` / `memory` / `disk` |

> **注意：** `start_time` 必须小于 `end_time`，否则返回参数错误。
