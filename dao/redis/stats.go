package redis

import (
	"strconv"
	"time"
)

// GetVisitCount 从Redis获取访问次数
func GetVisitCount(key string) (int, error) {
	val, err := client.Get(key).Result()
	if err == Nil {
		return 0, nil // key不存在，返回0
	} else if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// IncrVisitCount 增加访问次数
func IncrVisitCount(key string, expireSeconds int) error {
	// 使用管道确保原子性
	pipe := client.Pipeline()
	incrCmd := pipe.Incr(key)
	pipe.Expire(key, time.Duration(expireSeconds)*time.Second)

	_, err := pipe.Exec()
	if err != nil {
		return err
	}

	// 返回增加后的值（可选）
	_ = incrCmd.Val()

	return nil
}

// SetVisitCount 设置访问次数（带过期时间）
func SetVisitCount(key string, count int, expireSeconds int) error {
	return client.Set(key, count, time.Duration(expireSeconds)*time.Second).Err()
}