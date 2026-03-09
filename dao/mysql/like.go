package mysql

import (
	"bluebell/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// CreateLike 创建点赞记录（原始版本）
func CreateLike(userID int64, targetType models.TargetType, targetID int64) error {
	sqlStr := `INSERT INTO likes (user_id, target_type, target_id) VALUES (?, ?, ?)`
	_, err := db.Exec(sqlStr, userID, targetType, targetID)
	return err
}

// CreateLikeOptimistic 乐观锁方式创建点赞记录
// 使用INSERT IGNORE防止重复点赞，返回是否成功插入
func CreateLikeOptimistic(userID int64, targetType models.TargetType, targetID int64) (bool, error) {
	sqlStr := `INSERT IGNORE INTO likes (user_id, target_type, target_id) VALUES (?, ?, ?)`
	result, err := db.Exec(sqlStr, userID, targetType, targetID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

// OptimisticLike 乐观锁点赞（原子操作）
// 如果目标类型为article，同时更新文章点赞计数
func OptimisticLike(userID int64, targetType models.TargetType, targetID int64) error {
	// 使用事务确保一致性（如果需要同时更新计数的话）
	return WithTransaction(func(tx *sqlx.Tx) error {
		// 尝试插入点赞记录（IGNORE方式）
		sqlStr := `INSERT IGNORE INTO likes (user_id, target_type, target_id) VALUES (?, ?, ?)`
		result, err := tx.Exec(sqlStr, userID, targetType, targetID)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		// 如果没有插入成功（重复点赞），直接返回成功（或可以返回已存在错误）
		if rowsAffected == 0 {
			// 可以选择返回一个特定错误，或者直接返回nil
			return nil // 或 return errors.New("already liked")
		}

		// 成功插入，更新对应目标的点赞计数
		switch targetType {
		case models.TargetTypeArticle:
			// 原子更新文章点赞计数
			updateSql := `UPDATE articles SET like_count = like_count + 1 WHERE id = ?`
			if _, err := tx.Exec(updateSql, targetID); err != nil {
				return err
			}
		case models.TargetTypeComment:
			// 原子更新评论点赞计数
			updateSql := `UPDATE comments SET like_count = like_count + 1 WHERE id = ?`
			if _, err := tx.Exec(updateSql, targetID); err != nil {
				return err
			}
		default:
			return nil // 其他类型不更新计数
		}

		return nil
	})
}

// OptimisticUnlike 乐观锁取消点赞（原子操作）
func OptimisticUnlike(userID int64, targetType models.TargetType, targetID int64) error {
	// 使用事务确保一致性
	return WithTransaction(func(tx *sqlx.Tx) error {
		// 删除点赞记录
		sqlStr := `DELETE FROM likes WHERE user_id = ? AND target_type = ? AND target_id = ?`
		result, err := tx.Exec(sqlStr, userID, targetType, targetID)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		// 如果没有删除成功（点赞记录不存在）
		if rowsAffected == 0 {
			// 可以选择返回一个特定错误
			return nil // 或 return errors.New("not liked")
		}

		// 成功删除，更新对应目标的点赞计数
		switch targetType {
		case models.TargetTypeArticle:
			// 原子更新文章点赞计数
			updateSql := `UPDATE articles SET like_count = like_count - 1 WHERE id = ? AND like_count > 0`
			if _, err := tx.Exec(updateSql, targetID); err != nil {
				return err
			}
		case models.TargetTypeComment:
			// 原子更新评论点赞计数
			updateSql := `UPDATE comments SET like_count = like_count - 1 WHERE id = ? AND like_count > 0`
			if _, err := tx.Exec(updateSql, targetID); err != nil {
				return err
			}
		default:
			return nil // 其他类型不更新计数
		}

		return nil
	})
}

// DeleteLike 删除点赞记录
func DeleteLike(userID int64, targetType models.TargetType, targetID int64) error {
	sqlStr := `DELETE FROM likes WHERE user_id = ? AND target_type = ? AND target_id = ?`
	result, err := db.Exec(sqlStr, userID, targetType, targetID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrorLikeNotExist
	}

	return nil
}

// CheckLikeExists 检查点赞记录是否存在
func CheckLikeExists(userID int64, targetType models.TargetType, targetID int64) (bool, error) {
	var count int64
	sqlStr := `SELECT COUNT(*) FROM likes WHERE user_id = ? AND target_type = ? AND target_id = ?`
	err := db.Get(&count, sqlStr, userID, targetType, targetID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetLikeCount 获取点赞数量
func GetLikeCount(targetType models.TargetType, targetID int64) (int64, error) {
	var count int64
	sqlStr := `SELECT COUNT(*) FROM likes WHERE target_type = ? AND target_id = ?`
	err := db.Get(&count, sqlStr, targetType, targetID)
	return count, err
}

// GetUserLikes 获取用户的点赞列表（返回targetID列表）
func GetUserLikes(userID int64, targetType models.TargetType, page, size int) ([]int64, int64, error) {
	var targetIDs []int64
	var total int64

	// 计算总数
	countSql := `SELECT COUNT(*) FROM likes WHERE user_id = ? AND target_type = ?`
	err := db.Get(&total, countSql, userID, targetType)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	querySql := `SELECT target_id FROM likes
		WHERE user_id = ? AND target_type = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	err = db.Select(&targetIDs, querySql, userID, targetType, size, (page-1)*size)
	if err != nil {
		return nil, 0, err
	}

	return targetIDs, total, nil
}

// GetUserLikeDetails 获取用户的点赞详情列表（包含完整信息）
func GetUserLikeDetails(userID int64, targetType models.TargetType, page, size int) ([]*models.Like, int64, error) {
	var likes []*models.Like
	var total int64

	// 计算总数
	countSql := `SELECT COUNT(*) FROM likes WHERE user_id = ? AND target_type = ?`
	err := db.Get(&total, countSql, userID, targetType)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	querySql := `SELECT id, user_id, target_type, target_id, created_at FROM likes
		WHERE user_id = ? AND target_type = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	err = db.Select(&likes, querySql, userID, targetType, size, (page-1)*size)
	if err != nil {
		return nil, 0, err
	}

	return likes, total, nil
}

// GetLikeStatus 获取点赞状态
func GetLikeStatus(userID int64, targetType models.TargetType, targetID int64) (*models.ApiLikeStatus, error) {
	var like models.Like
	sqlStr := `SELECT id, created_at FROM likes
		WHERE user_id = ? AND target_type = ? AND target_id = ?`

	err := db.Get(&like, sqlStr, userID, targetType, targetID)
	if err == sql.ErrNoRows {
		// 未点赞
		count, err := GetLikeCount(targetType, targetID)
		if err != nil {
			return nil, err
		}
		return &models.ApiLikeStatus{
			IsLiked:   false,
			LikeCount: int(count),
		}, nil
	}
	if err != nil {
		return nil, err
	}

	// 已点赞
	count, err := GetLikeCount(targetType, targetID)
	if err != nil {
		return nil, err
	}

	return &models.ApiLikeStatus{
		IsLiked:   true,
		LikeCount: int(count),
		CreatedAt: like.CreatedAt,
	}, nil
}

// GetBatchLikeStatus 批量获取点赞状态
func GetBatchLikeStatus(userID int64, targets []models.TargetInfo) ([]*models.BatchLikeStatusResponse, error) {
	if len(targets) == 0 {
		return nil, nil
	}

	results := make([]*models.BatchLikeStatusResponse, 0, len(targets))

	// 由于需要查询每个目标，这里使用批量查询会更好
	// 但为了简化实现，我们循环查询每个目标
	// 未来可以优化为批量查询

	for _, target := range targets {
		// 获取点赞状态
		likeStatus, err := GetLikeStatus(userID, models.TargetType(target.TargetType), target.TargetID)
		if err != nil {
			// 如果查询出错，返回部分结果
			// 可以选择跳过这个目标或返回错误
			results = append(results, &models.BatchLikeStatusResponse{
				TargetType: target.TargetType,
				TargetID:   target.TargetID,
				IsLiked:    nil, // nil表示查询失败
				LikeCount:  0,
			})
			continue
		}

		// 构建响应
		status := &models.BatchLikeStatusResponse{
			TargetType: target.TargetType,
			TargetID:   target.TargetID,
			LikeCount:  likeStatus.LikeCount,
		}

		// 注意：如果用户未登录，likeStatus可能不包含用户是否点赞的信息
		// 在实际实现中，需要检查是否为有效的用户点赞状态
		// 这里简化处理
		if userID > 0 {
			isLiked := likeStatus.IsLiked
			status.IsLiked = &isLiked
		} else {
			status.IsLiked = nil // 未登录用户无法获取点赞状态
		}

		results = append(results, status)
	}

	return results, nil
}