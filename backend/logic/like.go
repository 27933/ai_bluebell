package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

// Like 点赞（优化版，使用乐观锁和事务）
func Like(userID int64, targetType models.TargetType, targetID int64) error {
	// 检查点赞目标是否存在和状态
	switch targetType {
	case models.TargetTypeArticle:
		// 检查文章是否存在
		article, err := mysql.GetArticleById(targetID)
		if err != nil {
			return mysql.ErrorArticleNotExist
		}
		if article.Status != string(models.ArticleStatusPublished) {
			return errors.New("article is not published")
		}
		if !article.AllowComment { // 检查是否允许点赞（假设与评论设置相同）
			return errors.New("article does not allow likes")
		}
	case models.TargetTypeComment:
		// 检查评论是否存在
		comment, err := mysql.GetCommentById(targetID)
		if err != nil {
			return mysql.ErrorCommentNotExist
		}
		if comment.Status != string(models.CommentStatusActive) {
			return errors.New("comment is not active")
		}
	default:
		return fmt.Errorf("unsupported target type: %s", targetType)
	}

	// 使用乐观锁方式点赞（包含事务和计数更新）
	if err := mysql.OptimisticLike(userID, targetType, targetID); err != nil {
		zap.L().Error("mysql.OptimisticLike() failed", zap.Error(err))
		return err
	}

	return nil
}

// Unlike 取消点赞（优化版，使用乐观锁和事务）
func Unlike(userID int64, targetType models.TargetType, targetID int64) error {
	// 检查点赞目标是否存在（可选，但可以提前验证）
	switch targetType {
	case models.TargetTypeArticle:
		// 检查文章是否存在（可选）
		if _, err := mysql.GetArticleById(targetID); err != nil {
			return mysql.ErrorArticleNotExist
		}
	case models.TargetTypeComment:
		// 检查评论是否存在（可选）
		if _, err := mysql.GetCommentById(targetID); err != nil {
			return mysql.ErrorCommentNotExist
		}
	default:
		return fmt.Errorf("unsupported target type: %s", targetType)
	}

	// 使用乐观锁方式取消点赞（包含事务和计数更新）
	if err := mysql.OptimisticUnlike(userID, targetType, targetID); err != nil {
		zap.L().Error("mysql.OptimisticUnlike() failed", zap.Error(err))
		return err
	}

	// 注意：OptimisticUnlike如果没有找到记录会返回nil
	// 如果需要严格的错误检查，可以在这里添加额外的验证
	// 但通常不检查也可以，因为"取消一个不存在的点赞"可以视为成功

	return nil
}

// GetLikeStatus 获取点赞状态
func GetLikeStatus(userID int64, targetType models.TargetType, targetID int64) (*models.ApiLikeStatus, error) {
	// 验证targetType
	if targetType != models.TargetTypeArticle && targetType != models.TargetTypeComment {
		return nil, fmt.Errorf("invalid target type: %s", targetType)
	}

	// 调用DAO层获取点赞状态
	status, err := mysql.GetLikeStatus(userID, targetType, targetID)
	if err != nil {
		zap.L().Error("mysql.GetLikeStatus() failed", zap.Error(err))
		return nil, err
	}

	return status, nil
}

// GetUserLikes 获取用户点赞列表
func GetUserLikes(userID int64, targetType models.TargetType, page, size int) ([]*models.Like, int64, error) {
	// 参数验证
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}

	// 验证targetType
	if targetType != models.TargetTypeArticle && targetType != models.TargetTypeComment {
		return nil, 0, fmt.Errorf("invalid target type: %s", targetType)
	}

	// 获取点赞详细信息列表（包含创建时间和完整信息）
	likes, total, err := mysql.GetUserLikeDetails(userID, targetType, page, size)
	if err != nil {
		zap.L().Error("mysql.GetUserLikeDetails() failed", zap.Error(err))
		return nil, 0, err
	}

	return likes, total, nil
}

