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

// TryRecordVisit 同步尝试记录访问，返回是否为新访问
// 供 GetArticleDetailHandler 调用，确保响应中的 view_count 与实际一致
func TryRecordVisit(articleID int64, userID *int64, ipAddress string) (inserted bool) {
	now := time.Now()
	visitDate := now.Format("2006-01-02")

	// Redis 限流（仅对未登录用户）
	if userID == nil {
		redisKey := RedisKeyPrefixIPVisit + strconv.FormatInt(articleID, 10) + ":" + ipAddress + ":" + visitDate
		currentCount, err := redis.GetVisitCount(redisKey)
		if err != nil {
			zap.L().Error("get visit count from redis failed", zap.Error(err))
		} else if currentCount >= IPVisitLimit {
			zap.L().Info("IP visit limit reached",
				zap.Int64("article_id", articleID),
				zap.String("ip", ipAddress),
				zap.Int("limit", IPVisitLimit))
			return false
		}
		if err := redis.IncrVisitCount(redisKey, RedisKeyExpireIPVisit); err != nil {
			zap.L().Error("incr visit count failed", zap.Error(err))
		}
	}

	// DB 限流（仅对未登录用户）
	if userID == nil {
		canVisit, count, err := mysql.CheckIPVisitLimit(articleID, ipAddress, now, IPVisitLimit)
		if err != nil {
			zap.L().Error("check ip visit limit failed", zap.Error(err))
			return false
		}
		if !canVisit {
			zap.L().Info("IP visit limit in database",
				zap.Int64("article_id", articleID),
				zap.String("ip", ipAddress),
				zap.Int("count", count))
			return false
		}
	}

	// 原子插入（INSERT IGNORE），重复访问返回 false
	ok, err := mysql.RecordArticleVisit(articleID, userID, ipAddress, now)
	if err != nil {
		zap.L().Error("record article visit failed", zap.Error(err))
		return false
	}
	if ok {
		if err := mysql.UpdateArticleView(articleID, now); err != nil {
			zap.L().Error("update article view count failed", zap.Error(err))
		}
	}
	return ok
}

// RecordArticleViewWithAntiCheat 异步版（供其他场景使用）
func RecordArticleViewWithAntiCheat(articleID int64, userID *int64, ipAddress string) error {
	go TryRecordVisit(articleID, userID, ipAddress)
	return nil
}