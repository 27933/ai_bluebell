## 新增需求

### 需求:SQL文件必须存放在指定目录
系统必须将数据库初始化 SQL 文件存放在 `./sql/` 目录下，文件名为 `init.sql`。

#### 场景:SQL文件位置正确
- **当** 执行初始化脚本生成后
- **那么** 文件路径必须为 `./sql/init.sql`

### 需求:必须创建用户表
系统必须创建 `users` 表，包含用户认证和个人信息字段。

#### 场景:用户表结构完整
- **当** 执行建表语句后
- **那么** `users` 表必须包含以下字段：id, username, email, password, role, status, nickname, avatar, bio, total_words, total_likes, wechat_openid, github_id, ip_address, user_agent, last_login_at, extra, created_at, updated_at

### 需求:必须创建文章表
系统必须创建 `articles` 表，支持文章的完整生命周期管理。

#### 场景:文章表结构完整
- **当** 执行建表语句后
- **那么** `articles` 表必须包含以下字段：id, title, content, summary, word_count, author_id, status, is_featured, featured_at, allow_comment, like_count, comment_count, slug, meta_keywords, meta_description, view_count, ip_address, user_agent, extra, created_at, updated_at

### 需求:必须创建标签相关表
系统必须创建 `tags` 表和 `article_tags` 关联表，支持文章的多标签管理。

#### 场景:标签表结构完整
- **当** 执行建表语句后
- **那么** `tags` 表必须包含：id, name, description, slug, created_at, updated_at
- **且** `article_tags` 表必须包含：id, article_id, tag_id, created_at

### 需求:必须创建栏目相关表
系统必须创建 `categories` 表和 `article_categories` 关联表，支持文章的栏目分类。

#### 场景:栏目表结构完整
- **当** 执行建表语句后
- **那么** `categories` 表必须包含：id, category_name, introduction, created_by, created_at, updated_at
- **且** `article_categories` 表必须包含：id, article_id, category_id, created_at

### 需求:必须创建评论表
系统必须创建 `comments` 表，支持嵌套回复。

#### 场景:评论表结构完整
- **当** 执行建表语句后
- **那么** `comments` 表必须包含：id, article_id, user_id, parent_id, content, like_count, status, ip_address, user_agent, created_at, updated_at

### 需求:必须创建点赞表
系统必须创建 `likes` 表，支持对文章和评论的点赞。

#### 场景:点赞表结构完整
- **当** 执行建表语句后
- **那么** `likes` 表必须包含：id, user_id, target_type, target_id, created_at
- **且** `target_type` 必须支持 'article' 和 'comment' 两种类型

### 需求:必须创建访问统计相关表
系统必须创建 `article_visits` 和 `trending_articles` 表，支持访问统计和热门排行。

#### 场景:访问统计表结构完整
- **当** 执行建表语句后
- **那么** `article_visits` 表必须包含：id, article_id, user_id, ip_address, visit_date, visit_hour, created_at
- **且** `trending_articles` 表必须包含：id, article_id, period_type, period_date, view_count, unique_view_count, rank_position, created_at, updated_at

### 需求:必须使用utf8mb4编码
所有表必须使用 `utf8mb4` 字符集和 `utf8mb4_unicode_ci` 排序规则。

#### 场景:字符集正确
- **当** 检查任意表的字符集时
- **那么** 必须是 `utf8mb4` 和 `utf8mb4_unicode_ci`

### 需求:必须创建必要的索引
系统必须为常用查询字段创建索引以保证性能。

#### 场景:索引存在
- **当** 检查表索引时
- **那么** `users` 表必须有 username 唯一索引
- **且** `articles` 表必须有 author_id, status 索引
- **且** `likes` 表必须有 (user_id, target_type, target_id) 唯一索引

## 修改需求

## 移除需求
