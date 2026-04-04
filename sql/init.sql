-- Bluebell Database Initialization Script
-- Generated based on DAO layer code analysis
-- Character set: utf8mb4, Collation: utf8mb4_unicode_ci

-- ============================================
-- 1. Users Table
-- ============================================
CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(255) NOT NULL COMMENT '用户名',
    `email` VARCHAR(255) DEFAULT NULL COMMENT '邮箱',
    `password` VARCHAR(255) NOT NULL COMMENT '密码（MD5加密）',
    `role` VARCHAR(20) NOT NULL DEFAULT 'reader' COMMENT '角色：visitor/reader/author/admin',
    `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态：active/inactive',
    `nickname` VARCHAR(255) DEFAULT NULL COMMENT '昵称',
    `avatar` VARCHAR(500) DEFAULT NULL COMMENT '头像URL',
    `bio` TEXT DEFAULT NULL COMMENT '个人简介',
    `total_words` BIGINT NOT NULL DEFAULT 0 COMMENT '总字数',
    `total_likes` BIGINT NOT NULL DEFAULT 0 COMMENT '总点赞数',
    `wechat_openid` VARCHAR(255) DEFAULT NULL COMMENT '微信OpenID',
    `github_id` VARCHAR(255) DEFAULT NULL COMMENT 'GitHub ID',
    `ip_address` VARCHAR(50) DEFAULT NULL COMMENT '最后登录IP',
    `user_agent` VARCHAR(500) DEFAULT NULL COMMENT '最后登录UA',
    `last_login_at` TIMESTAMP NULL DEFAULT NULL COMMENT '最后登录时间',
    `extra` TEXT DEFAULT NULL COMMENT '扩展字段（JSON）',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`),
    UNIQUE KEY `idx_email` (`email`),
    KEY `idx_role` (`role`),
    KEY `idx_status` (`status`),
    KEY `idx_wechat_openid` (`wechat_openid`),
    KEY `idx_github_id` (`github_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================
-- 2. Articles Table
-- ============================================
CREATE TABLE IF NOT EXISTS `articles` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `title` VARCHAR(255) NOT NULL COMMENT '文章标题',
    `content` LONGTEXT NOT NULL COMMENT '文章内容',
    `summary` VARCHAR(500) DEFAULT NULL COMMENT '文章摘要',
    `word_count` INT NOT NULL DEFAULT 0 COMMENT '字数统计',
    `author_id` BIGINT NOT NULL COMMENT '作者ID',
    `status` VARCHAR(20) NOT NULL DEFAULT 'draft' COMMENT '状态：draft/published/offline',
    `is_featured` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否精选',
    `featured_at` TIMESTAMP NULL DEFAULT NULL COMMENT '精选时间',
    `allow_comment` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否允许评论',
    `like_count` INT NOT NULL DEFAULT 0 COMMENT '点赞数',
    `comment_count` INT NOT NULL DEFAULT 0 COMMENT '评论数',
    `slug` VARCHAR(255) DEFAULT NULL COMMENT 'URL别名',
    `meta_keywords` VARCHAR(500) DEFAULT NULL COMMENT 'SEO关键词',
    `meta_description` VARCHAR(500) DEFAULT NULL COMMENT 'SEO描述',
    `view_count` INT NOT NULL DEFAULT 0 COMMENT '浏览量',
    `ip_address` VARCHAR(50) DEFAULT NULL COMMENT '发布IP',
    `user_agent` VARCHAR(500) DEFAULT NULL COMMENT '发布UA',
    `extra` TEXT DEFAULT NULL COMMENT '扩展字段（JSON）',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_author_id` (`author_id`),
    KEY `idx_status` (`status`),
    KEY `idx_is_featured` (`is_featured`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_view_count` (`view_count`),
    KEY `idx_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章表';

-- ============================================
-- 3. Tags Table
-- ============================================
CREATE TABLE IF NOT EXISTS `tags` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(50) NOT NULL COMMENT '标签名',
    `description` VARCHAR(200) DEFAULT NULL COMMENT '标签描述',
    `slug` VARCHAR(50) DEFAULT NULL COMMENT 'URL别名',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name` (`name`),
    KEY `idx_slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='标签表';

-- ============================================
-- 4. Article-Tags Association Table
-- ============================================
CREATE TABLE IF NOT EXISTS `article_tags` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `article_id` BIGINT NOT NULL COMMENT '文章ID',
    `tag_id` BIGINT NOT NULL COMMENT '标签ID',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_article_tag` (`article_id`, `tag_id`),
    KEY `idx_tag_id` (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章-标签关联表';

-- ============================================
-- 5. Categories Table
-- ============================================
CREATE TABLE IF NOT EXISTS `categories` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `category_name` VARCHAR(128) NOT NULL COMMENT '栏目名称',
    `introduction` VARCHAR(256) NOT NULL COMMENT '栏目简介',
    `created_by` BIGINT NOT NULL COMMENT '创建者ID',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_category_name` (`category_name`),
    KEY `idx_created_by` (`created_by`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='栏目表';

-- ============================================
-- 6. Article-Categories Association Table
-- ============================================
CREATE TABLE IF NOT EXISTS `article_categories` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `article_id` BIGINT NOT NULL COMMENT '文章ID',
    `category_id` BIGINT NOT NULL COMMENT '栏目ID',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_article_category` (`article_id`, `category_id`),
    KEY `idx_category_id` (`category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章-栏目关联表';

-- ============================================
-- 7. Comments Table
-- ============================================
CREATE TABLE IF NOT EXISTS `comments` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `article_id` BIGINT NOT NULL COMMENT '文章ID',
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `parent_id` BIGINT DEFAULT NULL COMMENT '父评论ID（用于嵌套回复）',
    `content` TEXT NOT NULL COMMENT '评论内容',
    `like_count` INT NOT NULL DEFAULT 0 COMMENT '点赞数',
    `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态：active/deleted',
    `ip_address` VARCHAR(50) DEFAULT NULL COMMENT '评论IP',
    `user_agent` VARCHAR(500) DEFAULT NULL COMMENT '评论UA',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_article_id` (`article_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

-- ============================================
-- 8. Likes Table
-- ============================================
CREATE TABLE IF NOT EXISTS `likes` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT NOT NULL COMMENT '用户ID',
    `target_type` VARCHAR(20) NOT NULL COMMENT '目标类型：article/comment',
    `target_id` BIGINT NOT NULL COMMENT '目标ID',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_target` (`user_id`, `target_type`, `target_id`),
    KEY `idx_target` (`target_type`, `target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='点赞表';

-- ============================================
-- 9. Article Visits Table (for analytics)
-- ============================================
CREATE TABLE IF NOT EXISTS `article_visits` (
    `id` BIGINT NOT NULL,
    `article_id` BIGINT NOT NULL COMMENT '文章ID',
    `user_id` BIGINT DEFAULT NULL COMMENT '用户ID（登录用户）',
    `ip_address` VARCHAR(50) NOT NULL COMMENT '访问IP',
    `visit_date` DATE NOT NULL COMMENT '访问日期',
    `visit_hour` TINYINT NOT NULL COMMENT '访问小时（0-23）',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_article_date` (`article_id`, `visit_date`),
    KEY `idx_user_article_date` (`user_id`, `article_id`, `visit_date`),
    KEY `idx_ip_article_date` (`ip_address`, `article_id`, `visit_date`),
    KEY `idx_visit_date` (`visit_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章访问记录表';

-- ============================================
-- 10. Trending Articles Table
-- ============================================
CREATE TABLE IF NOT EXISTS `trending_articles` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `article_id` BIGINT NOT NULL COMMENT '文章ID',
    `period_type` VARCHAR(20) NOT NULL COMMENT '周期类型：daily/weekly/monthly',
    `period_date` DATE NOT NULL COMMENT '周期日期',
    `view_count` INT NOT NULL DEFAULT 0 COMMENT '浏览量',
    `unique_view_count` INT NOT NULL DEFAULT 0 COMMENT '独立访客数',
    `rank_position` INT NOT NULL DEFAULT 0 COMMENT '排名位置',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_article_period` (`article_id`, `period_type`, `period_date`),
    KEY `idx_period_rank` (`period_type`, `period_date`, `rank_position`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='热门文章排行表';

-- ============================================
-- Done
-- ============================================
