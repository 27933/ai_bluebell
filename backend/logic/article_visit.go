package logic

import (
	"strconv"
	"time"

	"bluebell/dao/mysql"
	"bluebell/dao/redis"

	"go.uber.org/zap"
)

const (
	// IPVisitLimit 同一IP每天对同一文章的最大访问次数
	IPVisitLimit = 10
	// RedisKeyPrefixIPVisit Redis中IP访问记录的key前缀
	RedisKeyPrefixIPVisit = "article:ipvisit:"
	// RedisKeyExpireIPVisit IP访问记录的过期时间（24小时）
	RedisKeyExpireIPVisit = 24 * 3600 // 24小时
)

// RecordArticleViewWithAntiCheat 记录文章访问（带防刷机制）
func RecordArticleViewWithAntiCheat(articleID int64, userID *int64, ipAddress string) error {
	// 获取当前时间
	now := time.Now()
	visitDate := now.Format("2006-01-02")

	// 检查IP访问限制（仅对未登录用户）
	if userID == nil {
		// 构建Redis key
		redisKey := RedisKeyPrefixIPVisit + strconv.FormatInt(articleID, 10) + ":" + ipAddress + ":" + visitDate

		// 获取当前访问次数
		currentCount, err := redis.GetVisitCount(redisKey)
		if err != nil {
			zap.L().Error("get visit count from redis failed", zap.Error(err))
			// 不中断主流程，继续记录访问
		} else if currentCount >= IPVisitLimit {
			// 达到访问限制，不记录此次访问
			zap.L().Info("IP visit limit reached",
				zap.Int64("article_id", articleID),
				zap.String("ip", ipAddress),
				zap.Int("limit", IPVisitLimit))
			return nil
		}

		// 增加访问计数
		if err := redis.IncrVisitCount(redisKey, RedisKeyExpireIPVisit); err != nil {
			zap.L().Error("incr visit count failed", zap.Error(err))
		}
	}

	// 记录访问到数据库（异步执行）
	go func() {
		// 检查数据库中的访问限制
		if userID == nil {
			canVisit, count, err := mysql.CheckIPVisitLimit(articleID, ipAddress, now, IPVisitLimit)
			if err != nil {
				zap.L().Error("check ip visit limit failed", zap.Error(err))
				return
			}
			if !canVisit {
				zap.L().Info("IP visit limit in database",
					zap.Int64("article_id", articleID),
					zap.String("ip", ipAddress),
					zap.Int("count", count))
				return
			}
		}

		// 记录访问
		if err := mysql.RecordArticleVisit(articleID, userID, ipAddress, now); err != nil {
			zap.L().Error("record article visit failed", zap.Error(err))
			return
		}

		// 更新文章总浏览量（原有的统计逻辑）
		if err := mysql.UpdateArticleView(articleID, now); err != nil {
			zap.L().Error("update article view count failed", zap.Error(err))
		}
	}()

	return nil
}