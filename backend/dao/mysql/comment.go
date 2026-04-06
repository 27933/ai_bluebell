package mysql

import (
	"bluebell/models"
	"database/sql"
	"time"
)

// CreateComment 创建评论
func CreateComment(comment *models.Comment) error {
	sqlStr := `INSERT INTO comments (
		article_id, user_id, parent_id, content, ip_address, user_agent
	) VALUES (?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(sqlStr,
		comment.ArticleID, comment.UserID, comment.ParentID, comment.Content,
		comment.IPAddress, comment.UserAgent,
	)
	if err != nil {
		return err
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	comment.ID = commentID

	return nil
}

// GetCommentById 根据ID获取评论
func GetCommentById(id int64) (*models.Comment, error) {
	comment := new(models.Comment)
	sqlStr := `SELECT id, article_id, user_id, parent_id, content, like_count, status,
		ip_address, user_agent, created_at, updated_at
		FROM comments WHERE id = ?`

	err := db.Get(comment, sqlStr, id)
	if err == sql.ErrNoRows {
		return nil, ErrorCommentNotExist
	}
	return comment, err
}

// GetCommentList 获取评论列表（含作者信息）
func GetCommentList(articleId int64, page, size int) ([]*models.ApiComment, int64, error) {
	var total int64

	// 计算总数：活跃评论 + 有活跃回复的已删除父评论
	countSql := `SELECT COUNT(*) FROM comments c
		WHERE c.article_id = ? AND (
			c.status = 'active'
			OR (c.status = 'deleted' AND c.parent_id IS NULL AND EXISTS (
				SELECT 1 FROM comments r WHERE r.parent_id = c.id AND r.status = 'active'
			))
		)`
	err := db.Get(&total, countSql, articleId)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据：活跃评论 + 有活跃回复的已删除父评论（墓碑模式）
	type commentRow struct {
		ID        int64     `db:"id"`
		Content   string    `db:"content"`
		LikeCount int       `db:"like_count"`
		Status    string    `db:"status"`
		ParentID  *int64    `db:"parent_id"`
		CreatedAt time.Time `db:"created_at"`
		AuthorID  int64     `db:"author_id"`
		Username  string    `db:"username"`
		Nickname  string    `db:"nickname"`
		Avatar    string    `db:"avatar"`
	}

	querySql := `SELECT c.id, c.content, c.like_count, c.status, c.parent_id, c.created_at,
		u.id AS author_id, u.username, COALESCE(u.nickname, '') AS nickname, COALESCE(u.avatar, '') AS avatar
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.article_id = ? AND (
			c.status = 'active'
			OR (c.status = 'deleted' AND c.parent_id IS NULL AND EXISTS (
				SELECT 1 FROM comments r WHERE r.parent_id = c.id AND r.status = 'active'
			))
		)
		ORDER BY c.created_at ASC
		LIMIT ? OFFSET ?`

	var rows []commentRow
	err = db.Select(&rows, querySql, articleId, size, (page-1)*size)
	if err != nil {
		return nil, 0, err
	}

	comments := make([]*models.ApiComment, 0, len(rows))
	for _, row := range rows {
		content := row.Content
		author := models.CommentAuthor{
			ID:       row.AuthorID,
			Username: row.Username,
			Nickname: row.Nickname,
			Avatar:   row.Avatar,
		}
		// 已删除的评论：隐藏内容和作者信息（墓碑占位）
		if row.Status == "deleted" {
			content = "该评论已被删除"
			author = models.CommentAuthor{}
		}
		comments = append(comments, &models.ApiComment{
			ID:        row.ID,
			Content:   content,
			LikeCount: row.LikeCount,
			Status:    row.Status,
			ParentID:  row.ParentID,
			CreatedAt: row.CreatedAt,
			Author:    author,
		})
	}

	return comments, total, nil
}

// GetCommentListWithReplies 获取评论列表（包含一级回复）
func GetCommentListWithReplies(articleId int64, page, size int) ([]*models.Comment, []*models.Comment, int64, error) {
	// 获取所有评论
	var allComments []*models.Comment
	sqlStr := `SELECT id, article_id, user_id, parent_id, content, like_count, status,
		ip_address, user_agent, created_at, updated_at
		FROM comments
		WHERE article_id = ? AND status = 'active'
		ORDER BY created_at DESC`

	err := db.Select(&allComments, sqlStr, articleId)
	if err != nil {
		return nil, nil, 0, err
	}

	// 分离主评论和回复
	var mainComments []*models.Comment
	var replies []*models.Comment

	for _, comment := range allComments {
		if comment.ParentID == nil {
			mainComments = append(mainComments, comment)
		} else {
			replies = append(replies, comment)
		}
	}

	// 分页处理主评论
	total := int64(len(mainComments))
	start := (page - 1) * size
	end := start + size
	if end > len(mainComments) {
		end = len(mainComments)
	}
	if start < len(mainComments) {
		mainComments = mainComments[start:end]
	} else {
		mainComments = []*models.Comment{}
	}

	return mainComments, replies, total, nil
}

// UpdateComment 更新评论
func UpdateComment(id int64, content string) error {
	sqlStr := "UPDATE comments SET content = ?, updated_at = NOW() WHERE id = ?"
	result, err := db.Exec(sqlStr, content, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrorCommentNotExist
	}

	return nil
}

// DeleteComment 删除评论（软删除）
func DeleteComment(id int64) error {
	sqlStr := "UPDATE comments SET status = 'deleted', updated_at = NOW() WHERE id = ?"
	result, err := db.Exec(sqlStr, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrorCommentNotExist
	}

	return nil
}

// UpdateCommentLikeCount 更新评论点赞数
func UpdateCommentLikeCount(commentId int64, delta int) error {
	sqlStr := "UPDATE comments SET like_count = like_count + ? WHERE id = ?"
	_, err := db.Exec(sqlStr, delta, commentId)
	return err
}

// GetUserComments 获取用户的评论列表
func GetUserComments(userId int64, page, size int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	// 计算总数
	countSql := `SELECT COUNT(*) FROM comments WHERE user_id = ? AND status = 'active'`
	err := db.Get(&total, countSql, userId)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	querySql := `SELECT c.id, c.article_id, c.user_id, c.parent_id, c.content, c.like_count, c.status,
		c.ip_address, c.user_agent, c.created_at, c.updated_at,
		a.title as article_title
		FROM comments c
		JOIN articles a ON c.article_id = a.id
		WHERE c.user_id = ? AND c.status = 'active'
		ORDER BY c.created_at DESC
		LIMIT ? OFFSET ?`

	type CommentWithArticle struct {
		models.Comment
		ArticleTitle string `db:"article_title"`
	}

	var results []CommentWithArticle
	err = db.Select(&results, querySql, userId, size, (page-1)*size)
	if err != nil {
		return nil, 0, err
	}

	// 转换结果
	for _, result := range results {
		comments = append(comments, &result.Comment)
	}

	return comments, total, nil
}