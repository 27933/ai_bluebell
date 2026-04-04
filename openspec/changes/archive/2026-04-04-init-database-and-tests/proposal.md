## 为什么

本地服务已启动，但数据库是空的。需要初始化 MySQL 数据库表结构，并通过 HTTP 接口测试验证 API 功能是否正常。SQL 文件需要持久化保存，方便后续部署使用。

## 变更内容

1. **创建数据库初始化脚本** - 在 `./sql/` 目录创建 `init.sql`，包含所有表的建表语句
2. **执行数据库初始化** - 运行 SQL 脚本创建表结构
3. **运行接口测试** - 执行 `tests/api/` 下的测试用例，验证 API 功能

### 需要创建的数据库表（10个）

| 表名 | 用途 |
|------|------|
| `users` | 用户信息（角色：visitor/reader/author/admin）|
| `articles` | 文章（状态：draft/published/offline）|
| `tags` | 标签 |
| `article_tags` | 文章-标签关联 |
| `categories` | 栏目 |
| `article_categories` | 文章-栏目关联 |
| `comments` | 评论（支持嵌套回复）|
| `likes` | 点赞（支持文章/评论）|
| `article_visits` | 访问记录（防刷机制）|
| `trending_articles` | 热门排行 |

## 功能 (Capabilities)

### 新增功能

- `database-init`: 数据库初始化脚本，包含所有表的 DDL 语句
- `api-testing`: HTTP 接口测试执行流程

### 修改功能

（无）

## 影响

- **文件系统**: 新增 `./sql/init.sql` 文件
- **数据库**: 创建 10 个表（users, articles, tags, article_tags, categories, article_categories, comments, likes, article_visits, trending_articles）
- **测试**: 运行 `tests/api/` 下的所有测试用例
