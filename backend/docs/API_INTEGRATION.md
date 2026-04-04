# Bluebell API 前后端联调文档

> 文档版本：v2.0.0  
> 更新日期：2026-04-04  
> 基础URL：`http://localhost:8084/api/v1`

---

## 目录

- [通用说明](#通用说明)
- [错误码表](#错误码表)
- [认证模块](#1-认证模块)
- [文章模块](#2-文章模块)
- [标签模块](#3-标签模块)
- [栏目模块](#4-栏目模块)
- [评论模块](#5-评论模块)
- [点赞模块](#6-点赞模块)
- [作者模块](#7-作者模块)
- [统计模块](#8-统计模块)
- [管理员模块](#9-管理员模块)
- [文件上传](#10-文件上传)
- [导出模块](#11-导出模块)
- [RSS订阅](#12-rss订阅)

---

## 通用说明

### 请求格式

- Content-Type: `application/json`
- 需要认证的接口在 Header 中携带：`Authorization: Bearer <access_token>`

### 响应格式

所有接口返回统一格式：

```json
{
  "code": 1000,
  "msg": "success",
  "data": { ... }
}
```

### 分页参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码 |
| size | int | 20 | 每页数量（最大50） |

### 分页响应

```json
{
  "list": [...],
  "total": 100,
  "page": 1,
  "size": 20,
  "pages": 5
}
```

---

## 错误码表

| 错误码 | 说明 |
|--------|------|
| 1000 | 成功 |
| 1001 | 请求参数错误 |
| 1002 | 用户名已存在 |
| 1003 | 用户名不存在 |
| 1004 | 用户名或密码错误 |
| 1005 | 服务繁忙 |
| 1006 | 需要登录 |
| 1007 | 无效的token |
| 1008 | 栏目不存在 |
| 1009 | 栏目已存在 |
| 1010 | 栏目名称无效 |
| 1011 | 无权限操作该栏目 |
| 1012 | 文章已在该栏目中 |
| 1013 | 没有权限 |
| 1014 | 文章不存在 |

---

## 1. 认证模块

### 1.1 用户注册

**POST** `/auth/signup`

**请求参数：**

```json
{
  "username": "testuser_001",
  "password": "test123456",
  "re_password": "test123456"
}
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "id": "1234567890123456789",
    "username": "testuser_001"
  }
}
```

**错误响应示例：**

```json
// 用户名已存在
{
  "code": 1002,
  "msg": "用户名已存在",
  "data": null
}

// 密码不匹配
{
  "code": 1001,
  "msg": "请求参数错误",
  "data": null
}
```

---

### 1.2 用户登录

**POST** `/auth/login`

**请求参数：**

```json
{
  "username": "testuser_001",
  "password": "test123456"
}
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "user": {
      "id": "1234567890123456789",
      "username": "testuser_001",
      "email": "",
      "role": "reader",
      "status": "active",
      "nickname": "",
      "avatar": "",
      "bio": "",
      "total_words": 0,
      "total_likes": 0,
      "created_at": "2026-04-04T10:30:00Z"
    },
    "token": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_in": 7200
    }
  }
}
```

**错误响应示例：**

```json
// 用户名或密码错误
{
  "code": 1004,
  "msg": "用户名或密码错误",
  "data": null
}
```

---

### 1.3 刷新Token

**POST** `/auth/refresh`

**请求参数：**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 7200
  }
}
```

---

### 1.4 获取用户信息

**GET** `/auth/profile`

**请求头：** `Authorization: Bearer <access_token>`

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "id": "1234567890123456789",
    "username": "testuser_001",
    "email": "test@example.com",
    "role": "author",
    "status": "active",
    "nickname": "测试用户",
    "avatar": "https://example.com/avatar.jpg",
    "bio": "这是我的个人简介",
    "total_words": 12500,
    "total_likes": 328,
    "created_at": "2026-04-04T10:30:00Z"
  }
}
```

---

### 1.5 更新用户信息

**PUT** `/auth/profile`

**请求头：** `Authorization: Bearer <access_token>`

**请求参数：**

```json
{
  "nickname": "新昵称",
  "email": "newemail@example.com",
  "bio": "更新后的个人简介",
  "avatar": "https://example.com/new-avatar.jpg"
}
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "id": "1234567890123456789",
    "username": "testuser_001",
    "nickname": "新昵称",
    "email": "newemail@example.com",
    "bio": "更新后的个人简介",
    "avatar": "https://example.com/new-avatar.jpg"
  }
}
```

---

## 2. 文章模块

### 2.1 获取文章列表

**GET** `/articles`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| size | int | 否 | 每页数量，默认20 |
| sort | string | 否 | 排序方式：time/popular |
| status | string | 否 | 状态筛选：draft/published/offline/all |
| tag | string | 否 | 标签名筛选 |
| keyword | string | 否 | 关键词搜索 |
| author_name | string | 否 | 作者用户名 |

**请求示例：**

```
GET /articles?page=1&size=10&sort=time&status=published
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": "1234567890123456789",
        "title": "Go语言入门教程",
        "summary": "本文介绍Go语言的基础知识...",
        "author": {
          "id": "9876543210987654321",
          "username": "alice",
          "nickname": "Alice",
          "avatar": "https://example.com/alice.jpg"
        },
        "tags": [
          {"id": "111", "name": "Go", "article_count": 15},
          {"id": "222", "name": "教程", "article_count": 30}
        ],
        "view_count": 1250,
        "like_count": 86,
        "comment_count": 12,
        "is_featured": true,
        "is_recent": false,
        "created_at": "2026-04-01T08:00:00Z",
        "updated_at": "2026-04-03T15:30:00Z"
      }
    ],
    "total": 156,
    "page": 1,
    "size": 10,
    "pages": 16
  }
}
```

---

### 2.2 获取文章详情

**GET** `/articles/:id`

**请求示例：**

```
GET /articles/1234567890123456789
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "id": "1234567890123456789",
    "title": "Go语言入门教程",
    "content": "# Go语言入门\n\n## 1. 安装Go\n\n首先需要安装Go环境...",
    "summary": "本文介绍Go语言的基础知识...",
    "word_count": 3500,
    "author": {
      "id": "9876543210987654321",
      "username": "alice",
      "nickname": "Alice",
      "avatar": "https://example.com/alice.jpg",
      "bio": "全栈开发工程师"
    },
    "tags": [
      {"id": "111", "name": "Go"},
      {"id": "222", "name": "教程"}
    ],
    "status": "published",
    "is_featured": true,
    "featured_at": "2026-04-02T10:00:00Z",
    "allow_comment": true,
    "like_count": 86,
    "comment_count": 12,
    "view_count": 1250,
    "slug": "go-tutorial",
    "meta_keywords": "Go,Golang,教程,入门",
    "meta_description": "Go语言入门教程，适合初学者",
    "created_at": "2026-04-01T08:00:00Z",
    "updated_at": "2026-04-03T15:30:00Z"
  }
}
```

---

### 2.3 创建文章

**POST** `/articles`

**请求头：** `Authorization: Bearer <access_token>`

**请求参数：**

```json
{
  "title": "我的第一篇文章",
  "content": "# 标题\n\n这是文章内容...",
  "tags": ["Go", "教程"],
  "status": "draft",
  "allow_comment": true,
  "slug": "my-first-article",
  "meta_keywords": "关键词1,关键词2",
  "meta_description": "文章描述"
}
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "id": "1234567890123456790",
    "title": "我的第一篇文章",
    "status": "draft",
    "created_at": "2026-04-04T12:00:00Z"
  }
}
```

---

### 2.4 更新文章

**PUT** `/author/articles/:id`

**请求头：** `Authorization: Bearer <access_token>`

**请求参数：**

```json
{
  "title": "更新后的标题",
  "content": "更新后的内容...",
  "tags": ["Go", "进阶"],
  "status": "published"
}
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "id": "1234567890123456790",
    "title": "更新后的标题",
    "updated_at": "2026-04-04T14:00:00Z"
  }
}
```

---

### 2.5 删除文章

**DELETE** `/author/articles/:id`

**请求头：** `Authorization: Bearer <access_token>`

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": null
}
```

---

### 2.6 更新文章状态

**PATCH** `/author/articles/:id/status`

**请求头：** `Authorization: Bearer <access_token>`

**请求参数：**

```json
{
  "status": "published"
}
```

| status | 说明 |
|--------|------|
| draft | 草稿 |
| published | 已发布 |
| offline | 已下线 |

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "id": "1234567890123456790",
    "status": "published"
  }
}
```

---

### 2.7 设置文章精选

**PATCH** `/author/articles/:id/featured`

**请求头：** `Authorization: Bearer <access_token>`

**请求参数：**

```json
{
  "is_featured": true
}
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "id": "1234567890123456790",
    "is_featured": true,
    "featured_at": "2026-04-04T14:30:00Z"
  }
}
```

---

### 2.8 获取精选文章

**GET** `/articles/featured`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| limit | int | 否 | 获取数量，默认10 |

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": "1234567890123456789",
        "title": "Go语言入门教程",
        "summary": "本文介绍Go语言的基础知识...",
        "author": {
          "id": "9876543210987654321",
          "username": "alice",
          "nickname": "Alice"
        },
        "view_count": 1250,
        "like_count": 86,
        "is_featured": true,
        "featured_at": "2026-04-02T10:00:00Z",
        "created_at": "2026-04-01T08:00:00Z"
      }
    ]
  }
}
```

---

### 2.9 搜索文章

**GET** `/articles/search`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 是 | 搜索关键词 |
| page | int | 否 | 页码 |
| size | int | 否 | 每页数量 |

**请求示例：**

```
GET /articles/search?keyword=Go语言&page=1&size=10
```

**成功响应：** 同文章列表格式

---

### 2.10 获取作者的文章列表

**GET** `/author/articles`

**请求头：** `Authorization: Bearer <access_token>`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| status | string | 否 | 状态筛选 |
| page | int | 否 | 页码 |
| size | int | 否 | 每页数量 |

**成功响应：** 同文章列表格式

---

## 3. 标签模块

### 3.1 获取所有标签

**GET** `/tags`

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": "111",
        "name": "Go",
        "description": "Go语言相关",
        "slug": "go",
        "article_count": 15,
        "created_at": "2026-01-01T00:00:00Z"
      },
      {
        "id": "222",
        "name": "JavaScript",
        "description": "JavaScript相关",
        "slug": "javascript",
        "article_count": 28,
        "created_at": "2026-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

### 3.2 获取标签下的文章

**GET** `/tags/:id/articles`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| size | int | 否 | 每页数量 |

**成功响应：** 同文章列表格式

---

### 3.3 创建标签

**POST** `/tags`

**请求头：** `Authorization: Bearer <access_token>`

**请求参数：**

```json
{
  "name": "React",
  "description": "React框架相关",
  "slug": "react"
}
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "id": "333",
    "name": "React",
    "description": "React框架相关",
    "slug": "react",
    "created_at": "2026-04-04T12:00:00Z"
  }
}
```

---

### 3.4 更新标签

**PUT** `/tags/:id`

**请求头：** `Authorization: Bearer <access_token>`

**请求参数：**

```json
{
  "name": "React.js",
  "description": "React.js框架相关"
}
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "id": "333",
    "name": "React.js",
    "description": "React.js框架相关",
    "updated_at": "2026-04-04T14:00:00Z"
  }
}
```

---

### 3.5 删除标签

**DELETE** `/tags/:id`

**请求头：** `Authorization: Bearer <access_token>`

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": null
}
```

**注意：** 被文章使用中的标签无法删除

---

### 3.6 获取作者使用的标签

**GET** `/author/tags`

**请求头：** `Authorization: Bearer <access_token>`

**成功响应：** 同标签列表格式

---

## 4. 栏目模块

### 4.1 获取栏目列表

**GET** `/categories`

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "category_name": "技术分享",
        "introduction": "技术文章分享栏目",
        "created_by": 1,
        "article_count": 45,
        "created_at": "2026-01-01T00:00:00Z"
      },
      {
        "id": 2,
        "category_name": "生活随笔",
        "introduction": "生活感悟与随笔",
        "created_by": 1,
        "article_count": 23,
        "created_at": "2026-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

### 4.2 获取栏目详情

**GET** `/categories/:id`

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "id": 1,
    "category_name": "技术分享",
    "introduction": "技术文章分享栏目",
    "created_by": 1,
    "article_count": 45,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-03-15T10:00:00Z"
  }
}
```

---

### 4.3 获取栏目下的文章

**GET** `/categories/:id/articles`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| size | int | 否 | 每页数量 |

**成功响应：** 同文章列表格式

---

### 4.4 创建栏目

**POST** `/categories`

**请求头：** `Authorization: Bearer <access_token>`

**请求参数：**

```json
{
  "category_name": "新栏目",
  "introduction": "这是一个新栏目的简介"
}
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "id": 3,
    "category_name": "新栏目",
    "introduction": "这是一个新栏目的简介",
    "created_at": "2026-04-04T12:00:00Z"
  }
}
```

---

### 4.5 更新栏目

**PUT** `/categories/:id`

**请求头：** `Authorization: Bearer <access_token>`

**请求参数：**

```json
{
  "category_name": "更新后的栏目名",
  "introduction": "更新后的简介"
}
```

---

### 4.6 删除栏目

**DELETE** `/categories/:id`

**请求头：** `Authorization: Bearer <access_token>`

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": null
}
```

---

### 4.7 添加文章到栏目

**POST** `/articles/:id/categories`

**请求头：** `Authorization: Bearer <access_token>`

**请求参数：**

```json
{
  "category_ids": [1, 2, 3]
}
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": null
}
```

---

### 4.8 获取文章所属栏目

**GET** `/articles/:id/categories`

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "category_name": "技术分享",
        "introduction": "技术文章分享栏目"
      }
    ]
  }
}
```

---

## 5. 评论模块

### 5.1 获取评论列表

**GET** `/comments`

**请求头：** `Authorization: Bearer <access_token>`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| article_id | int | 是 | 文章ID |
| page | int | 否 | 页码 |
| size | int | 否 | 每页数量 |

**请求示例：**

```
GET /comments?article_id=1234567890123456789&page=1&size=20
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": "111222333",
        "content": "写得很好，学到了很多！",
        "like_count": 5,
        "status": "active",
        "author": {
          "id": "9876543210987654321",
          "username": "bob",
          "nickname": "Bob",
          "avatar": "https://example.com/bob.jpg"
        },
        "parent_id": null,
        "created_at": "2026-04-04T10:00:00Z"
      },
      {
        "id": "111222334",
        "content": "谢谢分享！",
        "like_count": 2,
        "status": "active",
        "author": {
          "id": "1111111111111111111",
          "username": "charlie",
          "nickname": "Charlie"
        },
        "parent_id": "111222333",
        "created_at": "2026-04-04T11:00:00Z"
      }
    ],
    "total": 12,
    "page": 1,
    "size": 20,
    "pages": 1
  }
}
```

---

### 5.2 创建评论

**POST** `/comments`

**请求头：** `Authorization: Bearer <access_token>`

**请求参数：**

```json
{
  "article_id": 1234567890123456789,
  "content": "这是一条评论",
  "parent_id": null
}
```

**回复评论：**

```json
{
  "article_id": 1234567890123456789,
  "content": "这是一条回复",
  "parent_id": 111222333
}
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "id": "111222335",
    "content": "这是一条评论",
    "created_at": "2026-04-04T12:00:00Z"
  }
}
```

---

### 5.3 更新评论

**PUT** `/comments/:id`

**请求头：** `Authorization: Bearer <access_token>`

**请求参数：**

```json
{
  "content": "更新后的评论内容"
}
```

---

### 5.4 删除评论

**DELETE** `/comments/:id`

**请求头：** `Authorization: Bearer <access_token>`

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": null
}
```

---

## 6. 点赞模块

### 6.1 点赞

**POST** `/likes`

**请求头：** `Authorization: Bearer <access_token>`

**请求参数：**

```json
{
  "target_type": "article",
  "target_id": 1234567890123456789
}
```

| target_type | 说明 |
|-------------|------|
| article | 文章 |
| comment | 评论 |

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "is_liked": true,
    "like_count": 87
  }
}
```

---

### 6.2 取消点赞

**DELETE** `/likes`

**请求头：** `Authorization: Bearer <access_token>`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| target_type | string | 是 | article/comment |
| target_id | int | 是 | 目标ID |

**请求示例：**

```
DELETE /likes?target_type=article&target_id=1234567890123456789
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "is_liked": false,
    "like_count": 86
  }
}
```

---

### 6.3 获取点赞状态

**GET** `/likes/status`

**请求头：** `Authorization: Bearer <access_token>`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| target_type | string | 是 | article/comment |
| target_id | int | 是 | 目标ID |

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "is_liked": true,
    "like_count": 86,
    "created_at": "2026-04-04T10:00:00Z"
  }
}
```

---

### 6.4 获取用户点赞列表

**GET** `/user/likes`

**请求头：** `Authorization: Bearer <access_token>`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| target_type | string | 否 | article/comment，默认article |
| page | int | 否 | 页码 |
| size | int | 否 | 每页数量 |

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": "1234567890123456789",
        "title": "Go语言入门教程",
        "summary": "本文介绍Go语言的基础知识...",
        "created_at": "2026-04-01T08:00:00Z"
      }
    ],
    "total": 15,
    "page": 1,
    "size": 20,
    "pages": 1
  }
}
```

---

## 7. 作者模块

### 7.1 获取作者信息

**GET** `/authors/:username`

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "id": "9876543210987654321",
    "username": "alice",
    "nickname": "Alice",
    "avatar": "https://example.com/alice.jpg",
    "bio": "全栈开发工程师，专注于Go和React",
    "total_words": 125000,
    "total_likes": 3280
  }
}
```

---

### 7.2 获取作者文章列表

**GET** `/authors/:username/articles`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| sort | string | 否 | 排序方式：time/hot |
| page | int | 否 | 页码 |
| size | int | 否 | 每页数量 |

**请求示例：**

```
GET /authors/alice/articles?sort=time&page=1&size=10
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": "1234567890123456789",
        "title": "Go语言入门教程",
        "summary": "本文介绍Go语言的基础知识...",
        "view_count": 1250,
        "like_count": 86,
        "comment_count": 12,
        "is_featured": true,
        "is_recent": false,
        "created_at": "2026-04-01T08:00:00Z",
        "updated_at": "2026-04-03T15:30:00Z"
      }
    ],
    "total": 45,
    "page": 1,
    "size": 10,
    "pages": 5
  }
}
```

---

## 8. 统计模块

### 8.1 获取热门文章排行

**GET** `/articles/trending`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| period | string | 否 | 周期：daily/weekly/monthly，默认daily |
| page | int | 否 | 页码 |
| size | int | 否 | 每页数量 |

**请求示例：**

```
GET /articles/trending?period=weekly&page=1&size=10
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "article_id": 1234567890123456789,
        "period_type": "weekly",
        "period_date": "2026-04-01",
        "view_count": 5600,
        "unique_view_count": 3200,
        "rank_position": 1,
        "title": "Go语言入门教程",
        "summary": "本文介绍Go语言的基础知识...",
        "author_username": "alice",
        "author_nickname": "Alice",
        "like_count": 86,
        "comment_count": 12
      }
    ],
    "page": 1,
    "size": 10
  }
}
```

---

### 8.2 获取文章日统计

**GET** `/article-stats/daily`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| article_id | int | 是 | 文章ID |
| days | int | 否 | 统计天数，默认30，最大90 |

**请求示例：**

```
GET /article-stats/daily?article_id=1234567890123456789&days=7
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "article_id": 1234567890123456789,
    "total_views": 1250,
    "total_uv": 890,
    "daily_stats": {
      "2026-04-01": 180,
      "2026-04-02": 220,
      "2026-04-03": 195,
      "2026-04-04": 165
    }
  }
}
```

---

### 8.3 获取文章访问趋势

**GET** `/article-stats/trend`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| article_id | int | 是 | 文章ID |
| days | int | 否 | 统计天数，默认30 |
| group_by | string | 是 | 分组方式：hour/day/week/month |

**请求示例：**

```
GET /article-stats/trend?article_id=1234567890123456789&days=7&group_by=day
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "article_id": 1234567890123456789,
    "days": 7,
    "group_by": "day",
    "trend": [
      {"label": "2026-04-01", "value": 180},
      {"label": "2026-04-02", "value": 220},
      {"label": "2026-04-03", "value": 195},
      {"label": "2026-04-04", "value": 165}
    ]
  }
}
```

---

### 8.4 批量获取文章统计

**GET** `/article-stats/batch`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| ids | string | 是 | 文章ID列表，逗号分隔 |

**请求示例：**

```
GET /article-stats/batch?ids=1234567890123456789,1234567890123456790
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "stats": {
      "1234567890123456789": {
        "article_id": 1234567890123456789,
        "total_views": 1250,
        "total_uv": 890
      },
      "1234567890123456790": {
        "article_id": 1234567890123456790,
        "total_views": 560,
        "total_uv": 320
      }
    }
  }
}
```

---

### 8.5 记录文章浏览（防刷）

**POST** `/articles/view`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| article_id | int | 是 | 文章ID |

**请求示例：**

```
POST /articles/view?article_id=1234567890123456789
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "message": "访问记录成功"
  }
}
```

**注意：** 同一IP每天对同一文章最多记录1次有效浏览

---

## 9. 管理员模块

> 以下接口需要管理员权限（role=admin）

### 9.1 获取文章列表（管理员）

**GET** `/admin/articles`

**请求头：** `Authorization: Bearer <admin_token>`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| status | string | 否 | 状态筛选 |
| page | int | 否 | 页码 |
| size | int | 否 | 每页数量 |

**成功响应：** 同文章列表格式

---

### 9.2 管理员设置精选

**PATCH** `/admin/articles/:id/featured`

**请求头：** `Authorization: Bearer <admin_token>`

**请求参数：**

```json
{
  "is_featured": true
}
```

---

### 9.3 获取用户列表

**GET** `/admin/users`

**请求头：** `Authorization: Bearer <admin_token>`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| role | string | 否 | 角色筛选：visitor/reader/author/admin/all |
| status | string | 否 | 状态筛选：active/inactive/all |
| page | int | 否 | 页码 |
| size | int | 否 | 每页数量 |

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "list": [
      {
        "id": "9876543210987654321",
        "username": "alice",
        "email": "alice@example.com",
        "role": "author",
        "status": "active",
        "nickname": "Alice",
        "avatar": "https://example.com/alice.jpg",
        "total_words": 125000,
        "total_likes": 3280,
        "created_at": "2026-01-15T08:00:00Z",
        "last_login_at": "2026-04-04T09:30:00Z"
      }
    ],
    "total": 156,
    "page": 1,
    "size": 20,
    "pages": 8
  }
}
```

---

### 9.4 获取用户详情

**GET** `/admin/users/:id`

**请求头：** `Authorization: Bearer <admin_token>`

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "id": "9876543210987654321",
    "username": "alice",
    "email": "alice@example.com",
    "role": "author",
    "status": "active",
    "nickname": "Alice",
    "avatar": "https://example.com/alice.jpg",
    "bio": "全栈开发工程师",
    "total_words": 125000,
    "total_likes": 3280,
    "created_at": "2026-01-15T08:00:00Z",
    "updated_at": "2026-04-04T10:00:00Z",
    "last_login_at": "2026-04-04T09:30:00Z"
  }
}
```

---

### 9.5 更新用户状态

**PATCH** `/admin/users/:id/status`

**请求头：** `Authorization: Bearer <admin_token>`

**请求参数：**

```json
{
  "status": "inactive"
}
```

| status | 说明 |
|--------|------|
| active | 正常 |
| inactive | 禁用 |

---

### 9.6 更新用户角色

**PATCH** `/admin/users/:id/role`

**请求头：** `Authorization: Bearer <admin_token>`

**请求参数：**

```json
{
  "role": "author"
}
```

| role | 说明 |
|------|------|
| visitor | 访客 |
| reader | 读者 |
| author | 作者 |
| admin | 管理员 |

---

### 9.7 批量更新用户状态

**PATCH** `/admin/users/batch/status`

**请求头：** `Authorization: Bearer <admin_token>`

**请求参数：**

```json
{
  "user_ids": [1, 2, 3],
  "status": "inactive"
}
```

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "updated_count": 3
  }
}
```

---

### 9.8 获取系统概览

**GET** `/admin/stats/overview`

**请求头：** `Authorization: Bearer <admin_token>`

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "total_users": 1560,
    "total_articles": 890,
    "total_comments": 12500,
    "new_users_today": 15,
    "new_articles_today": 8,
    "new_comments_today": 125
  }
}
```

---

### 9.9 获取系统日统计

**GET** `/admin/stats/daily`

**请求头：** `Authorization: Bearer <admin_token>`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| start_date | string | 否 | 开始日期 YYYY-MM-DD |
| end_date | string | 否 | 结束日期 YYYY-MM-DD |

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "list": [
      {
        "date": "2026-04-04",
        "new_user_count": 15,
        "new_article_count": 8,
        "new_comment_count": 125
      },
      {
        "date": "2026-04-03",
        "new_user_count": 12,
        "new_article_count": 10,
        "new_comment_count": 98
      }
    ]
  }
}
```

---

### 9.10 获取实时系统指标

**GET** `/admin/metrics/realtime`

**请求头：** `Authorization: Bearer <admin_token>`

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "cpu_usage": 35.5,
    "memory_usage": 62.3,
    "goroutines": 156,
    "gc_pause_ms": 1.2,
    "uptime_seconds": 86400,
    "requests_per_second": 125.5,
    "timestamp": 1712235600
  }
}
```

---

### 9.11 获取系统指标历史

**GET** `/admin/metrics/history`

**请求头：** `Authorization: Bearer <admin_token>`

**查询参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| metric_type | string | 是 | 指标类型：cpu/memory/requests |
| start_time | int | 是 | 开始时间戳 |
| end_time | int | 是 | 结束时间戳 |

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "metric_type": "cpu",
    "start_time": 1712232000,
    "end_time": 1712235600,
    "data": [
      {"timestamp": 1712232000, "value": 32.5},
      {"timestamp": 1712232300, "value": 35.2},
      {"timestamp": 1712232600, "value": 28.8}
    ]
  }
}
```

---

## 10. 文件上传

### 10.1 上传图片

**POST** `/upload/image`

**请求头：**
- `Authorization: Bearer <access_token>`
- `Content-Type: multipart/form-data`

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 图片文件 |

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "url": "/uploads/images/1712235600.jpg",
    "filename": "1712235600.jpg",
    "size": 125678
  }
}
```

---

### 10.2 上传附件

**POST** `/upload/attachment`

**请求头：**
- `Authorization: Bearer <access_token>`
- `Content-Type: multipart/form-data`

**请求参数：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 附件文件 |

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "url": "/uploads/attachments/document.pdf",
    "filename": "document.pdf",
    "size": 2456789
  }
}
```

---

## 11. 导出模块

### 11.1 导出单篇文章

**GET** `/author/articles/:id/export`

**请求头：** `Authorization: Bearer <access_token>`

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "filename": "go-tutorial.md",
    "content": "# Go语言入门教程\n\n## 1. 安装Go\n\n首先需要安装Go环境...",
    "size": 3500
  }
}
```

---

### 11.2 批量导出文章

**POST** `/author/articles/export`

**请求头：** `Authorization: Bearer <access_token>`

**请求参数：**

```json
{
  "article_ids": [1234567890123456789, 1234567890123456790]
}
```

**注意：** 最多支持50篇文章

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "batch_id": "export_20260404_120000",
    "files": [
      {
        "article_id": 1234567890123456789,
        "filename": "go-tutorial.md",
        "content": "# Go语言入门教程...",
        "size": 3500
      },
      {
        "article_id": 1234567890123456790,
        "filename": "react-basics.md",
        "content": "# React基础教程...",
        "size": 2800
      }
    ],
    "file_count": 2,
    "total_size": 6300
  }
}
```

---

## 12. RSS订阅

### 12.1 获取RSS订阅

**GET** `/rss`

**成功响应：**

```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "title": "知识博客",
    "link": "https://example.com",
    "description": "知识博客 - 分享技术，传播知识",
    "language": "zh-CN",
    "pub_date": "Fri, 04 Apr 2026 12:00:00 +0800",
    "items": [
      {
        "title": "Go语言入门教程",
        "link": "https://example.com/articles/1234567890123456789",
        "description": "本文介绍Go语言的基础知识...",
        "pub_date": "Mon, 01 Apr 2026 08:00:00 +0800",
        "author": "alice"
      },
      {
        "title": "React基础教程",
        "link": "https://example.com/articles/1234567890123456790",
        "description": "本文介绍React的基础知识...",
        "pub_date": "Tue, 02 Apr 2026 10:00:00 +0800",
        "author": "bob"
      }
    ]
  }
}
```

---

## 附录

### A. 用户角色权限

| 角色 | 权限说明 |
|------|----------|
| visitor | 仅可浏览公开内容 |
| reader | 可评论、点赞 |
| author | 可发布文章、管理自己的内容 |
| admin | 全部权限 |

### B. 文章状态流转

```
draft(草稿) -> published(已发布) -> offline(已下线)
                    ^                    |
                    |____________________|
```

### C. 测试账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| testuser_001 | test123456 | reader |
| testauthor_001 | test123456 | author |
| admin | admin123456 | admin |

---

**文档生成时间：** 2026-04-04  
**API版本：** v2.0.0
