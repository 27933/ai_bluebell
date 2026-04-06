package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"errors"
	"time"

	"go.uber.org/zap"
)

// CreateComment 创建评论
func CreateComment(p *models.ParamCreateComment) (*models.Comment, error) {
	// 检查文章是否存在
	article, err := mysql.GetArticleById(p.ArticleID)
	if err != nil {
		zap.L().Error("mysql.GetArticleById() failed", zap.Int64("articleID", p.ArticleID), zap.Error(err))
		return nil, mysql.ErrorArticleNotExist
	}

	// 检查文章是否允许评论
	if !article.AllowComment {
		return nil, errors.New("article does not allow comments")
	}

	// 检查父评论是否存在（如果提供了父评论ID）
	if p.ParentID != nil && *p.ParentID > 0 {
		parentComment, err := mysql.GetCommentById(*p.ParentID)
		if err != nil {
			return nil, errors.New("parent comment not found")
		}
		if parentComment.ArticleID != p.ArticleID {
			return nil, errors.New("parent comment does not belong to this article")
		}
		// 只允许一级回复
		if parentComment.ParentID != nil {
			return nil, errors.New("only one level of replies is allowed")
		}
	}

	// 创建评论
	comment := &models.Comment{
		ID:        snowflake.GenID(),
		ArticleID: p.ArticleID,
		UserID:    p.UserID,
		ParentID:  p.ParentID,
		Content:   p.Content,
		Status:    string(models.CommentStatusActive),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存评论
	if err := mysql.CreateComment(comment); err != nil {
		zap.L().Error("mysql.CreateComment() failed", zap.Error(err))
		return nil, err
	}

	// 更新文章评论计数
	if err := mysql.UpdateArticleCommentCount(p.ArticleID); err != nil {
		zap.L().Error("mysql.UpdateArticleCommentCount() failed", zap.Error(err))
	}

	return comment, nil
}

// GetCommentList 获取评论列表
func GetCommentList(p *models.ParamCommentList) ([]*models.ApiComment, int64, error) {
	// 参数验证
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Size < 1 || p.Size > 50 {
		p.Size = 20
	}

	// 获取评论列表
	comments, total, err := mysql.GetCommentList(p.ArticleID, p.Page, p.Size)
	if err != nil {
		zap.L().Error("mysql.GetCommentList() failed", zap.Error(err))
		return nil, 0, err
	}

	return comments, total, nil
}

// UpdateComment 更新评论
func UpdateComment(id int64, userID int64, content string) error {
	// 检查评论是否存在且属于该用户
	comment, err := mysql.GetCommentById(id)
	if err != nil {
		return mysql.ErrorCommentNotExist
	}

	if comment.UserID != userID {
		return errors.New("you can only update your own comments")
	}

	if comment.Status == string(models.CommentStatusDeleted) {
		return errors.New("comment has been deleted")
	}

	// 更新评论
	if err := mysql.UpdateComment(id, content); err != nil {
		zap.L().Error("mysql.UpdateComment() failed", zap.Error(err))
		return err
	}

	return nil
}

// DeleteComment 删除评论
func DeleteComment(id int64, userID int64) error {
	// 检查评论是否存在
	comment, err := mysql.GetCommentById(id)
	if err != nil {
		return mysql.ErrorCommentNotExist
	}

	// 允许删除的条件：评论作者 OR 文章作者 OR 管理员
	if comment.UserID != userID {
		// 检查是否是文章作者
		article, err := mysql.GetArticleById(comment.ArticleID)
		if err != nil || article.AuthorID != userID {
			// 检查是否是管理员
			user, err := mysql.GetUserById(userID)
			if err != nil || user.Role != "admin" {
				return errors.New("no permission to delete this comment")
			}
		}
	}

	// 软删除评论
	if err := mysql.DeleteComment(id); err != nil {
		zap.L().Error("mysql.DeleteComment() failed", zap.Error(err))
		return err
	}

	// 减少文章评论计数
	if comment.ArticleID > 0 {
		if err := mysql.UpdateArticleCommentCount(comment.ArticleID); err != nil {
			zap.L().Warn("mysql.UpdateArticleCommentCount() failed", zap.Error(err))
		}
	}

	return nil
}

// GetUserComments 获取用户评论列表
func GetUserComments(userID int64, page, size int) ([]*models.Comment, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}

	comments, total, err := mysql.GetUserComments(userID, page, size)
	if err != nil {
		zap.L().Error("mysql.GetUserComments() failed", zap.Error(err))
		return nil, 0, err
	}

	return comments, total, nil
}