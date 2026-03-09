-- 访问统计功能增强 - 数据库迁移脚本

-- 1. 创建访问记录表（用于防刷和详细分析）
CREATE TABLE IF NOT EXISTS article_visits (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    article_id BIGINT NOT NULL,
    user_id BIGINT DEFAULT NULL,
    ip_address VARCHAR(45) NOT NULL,
    visit_date DATE NOT NULL,
    visit_hour TINYINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_article_date (article_id, visit_date),
    INDEX idx_ip_date (ip_address, visit_date),
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. 创建热门文章排行表（每日更新）
CREATE TABLE IF NOT EXISTS trending_articles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    article_id BIGINT NOT NULL,
    period_type ENUM('daily', 'weekly', 'monthly') NOT NULL,
    period_date DATE NOT NULL,
    view_count INT DEFAULT 0,
    unique_view_count INT DEFAULT 0,
    rank_position INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_article_period (article_id, period_type, period_date),
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. 为article_stats表添加复合索引
-- 注意：MySQL 5.7不支持IF NOT EXISTS，需要先检查索引是否存在
-- 这里使用存储过程来安全添加索引
DELIMITER //
CREATE PROCEDURE AddIndexIfNotExists(IN tableName VARCHAR(64), IN indexName VARCHAR(64), IN indexColumns VARCHAR(255))
BEGIN
    DECLARE indexExists INT DEFAULT 0;

    SELECT COUNT(*) INTO indexExists
    FROM information_schema.statistics
    WHERE table_schema = DATABASE()
    AND table_name = tableName
    AND index_name = indexName;

    IF indexExists = 0 THEN
        SET @sql = CONCAT('ALTER TABLE ', tableName, ' ADD INDEX ', indexName, ' (', indexColumns, ')');
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
    END IF;
END //
DELIMITER ;

-- 安全添加索引
CALL AddIndexIfNotExists('article_stats', 'idx_date_views', 'date, views');
CALL AddIndexIfNotExists('article_stats', 'idx_article_date', 'article_id, date');

-- 4. 为trending_articles表添加索引
CALL AddIndexIfNotExists('trending_articles', 'idx_period_rank', 'period_type, period_date, rank_position');
CALL AddIndexIfNotExists('trending_articles', 'idx_unique_views', 'period_type, period_date, unique_view_count DESC');

-- 5. 初始化今日热门文章数据（存储过程）
DELIMITER //
CREATE PROCEDURE InitTrendingArticles(IN p_period_type VARCHAR(10), IN p_target_date DATE)
BEGIN
    DECLARE v_start_date DATE;

    -- 计算统计开始日期
    IF p_period_type = 'daily' THEN
        SET v_start_date = p_target_date;
    ELSEIF p_period_type = 'weekly' THEN
        SET v_start_date = DATE_SUB(p_target_date, INTERVAL 6 DAY);
    ELSEIF p_period_type = 'monthly' THEN
        SET v_start_date = DATE_SUB(p_target_date, INTERVAL 29 DAY);
    END IF;

    -- 清空指定周期的旧数据
    DELETE FROM trending_articles WHERE period_type = p_period_type AND period_date = p_target_date;

    -- 插入新的排行数据
    INSERT INTO trending_articles (article_id, period_type, period_date, view_count, unique_view_count, rank_position)
    SELECT
        a.id as article_id,
        p_period_type as period_type,
        p_target_date as period_date,
        COALESCE(SUM(s.views), 0) as view_count,
        COALESCE(COUNT(DISTINCT v.ip_address), 0) as unique_view_count,
        0 as rank_position
    FROM articles a
    LEFT JOIN article_stats s ON a.id = s.article_id
        AND s.date >= v_start_date AND s.date <= p_target_date
    LEFT JOIN article_visits v ON a.id = v.article_id
        AND v.visit_date >= v_start_date AND v.visit_date <= p_target_date
    WHERE a.status = 'published'
    GROUP BY a.id
    HAVING unique_view_count > 0
    ORDER BY unique_view_count DESC, view_count DESC
    LIMIT 100;

    -- 更新排行位置
    SET @rank := 0;
    UPDATE trending_articles
    SET rank_position = (@rank := @rank + 1)
    WHERE period_type = p_period_type AND period_date = p_target_date
    ORDER BY unique_view_count DESC, view_count DESC;
END //
DELIMITER ;

-- 6. 创建定时任务事件（需要MySQL事件调度器支持）
-- 查看事件调度器状态：SHOW VARIABLES LIKE 'event_scheduler';
-- 开启事件调度器：SET GLOBAL event_scheduler = ON;

-- 每日凌晨2点更新日榜
CREATE EVENT IF NOT EXISTS update_daily_trending
ON SCHEDULE EVERY 1 DAY STARTS DATE_ADD(CURDATE(), INTERVAL 2 HOUR)
DO CALL InitTrendingArticles('daily', CURDATE() - INTERVAL 1 DAY);

-- 每周一凌晨3点更新周榜
CREATE EVENT IF NOT EXISTS update_weekly_trending
ON SCHEDULE EVERY 1 WEEK STARTS DATE_ADD(DATE_ADD(CURDATE(), INTERVAL 1 WEEK), INTERVAL 3 HOUR)
DO CALL InitTrendingArticles('weekly', CURDATE() - INTERVAL 1 WEEK);

-- 每月1号凌晨4点更新月榜
CREATE EVENT IF NOT EXISTS update_monthly_trending
ON SCHEDULE EVERY 1 MONTH STARTS DATE_ADD(DATE_ADD(CURDATE(), INTERVAL 1 MONTH), INTERVAL 4 HOUR)
DO CALL InitTrendingArticles('monthly', CURDATE() - INTERVAL 1 MONTH);

-- 7. 插入测试数据（可选）
-- 可以手动执行以下语句测试排行功能：
-- CALL InitTrendingArticles('daily', CURDATE());

-- 8. 清理过期访问记录（保留90天）
CREATE EVENT IF NOT EXISTS cleanup_old_visits
ON SCHEDULE EVERY 1 DAY STARTS DATE_ADD(CURDATE(), INTERVAL 3 HOUR)
DO DELETE FROM article_visits WHERE visit_date < DATE_SUB(CURDATE(), INTERVAL 90 DAY);