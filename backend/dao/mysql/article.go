package mysql

import (
	"bluebell/models"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// CreateArticle 创建文章
func CreateArticle(article *models.Article) error {
	sqlStr := `INSERT INTO articles (
		title, content, summary, word_count, author_id, status,
		is_featured, allow_comment, slug, meta_keywords, meta_description,
		ip_address, user_agent, extra
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(sqlStr,
		article.Title, article.Content, article.Summary, article.WordCount,
		article.AuthorID, article.Status, article.IsFeatured, article.AllowComment,
		article.Slug, article.MetaKeywords, article.MetaDescription,
		article.IPAddress, article.UserAgent, article.Extra,
	)
	if err != nil {
		return err
	}

	articleID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	article.ID = articleID

	return nil
}

// GetArticleById 根据ID获取文章
func GetArticleById(id int64) (*models.Article, error) {
	article := new(models.Article)
	sqlStr := `SELECT
		id, title, content, summary, word_count, author_id, status,
		is_featured, featured_at, allow_comment, like_count, comment_count,
		slug, meta_keywords, meta_description, view_count,
		ip_address, user_agent, extra, created_at, updated_at
		FROM articles WHERE id = ?`

	err := db.Get(article, sqlStr, id)
	if err == sql.ErrNoRows {
		return nil, ErrorArticleNotExist
	}
	return article, err
}

// GetArticleList 获取文章列表
func GetArticleList(param *models.ParamArticleList) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	// 构建查询条件
	conditions := []string{"1=1"}
	args := []interface{}{}

	if param.Status != "" && param.Status != "all" {
		conditions = append(conditions, "status = ?")
		args = append(args, param.Status)
	}

	if param.AuthorID > 0 {
		conditions = append(conditions, "author_id = ?")
		args = append(args, param.AuthorID)
	}

	if param.AuthorName != "" {
		// 按作者用户名搜索，需要联表查询
		conditions = append(conditions, `author_id IN (
			SELECT id FROM users
			WHERE username LIKE ?
		)`)
		args = append(args, "%"+param.AuthorName+"%")
	}

	if param.Keyword != "" {
		conditions = append(conditions, "(title LIKE ? OR content LIKE ?)")
		args = append(args, "%"+param.Keyword+"%", "%"+param.Keyword+"%")
	}

	if param.Tag != "" {
		// 标签查询需要联表
		conditions = append(conditions, `id IN (
			SELECT DISTINCT article_id FROM article_tags
			WHERE tag_id IN (
				SELECT id FROM tags WHERE name = ?
			)
		)`)
		args = append(args, param.Tag)
	}

	// 计算总数
	countSql := `SELECT COUNT(*) FROM articles WHERE ` + strings.Join(conditions, " AND ")
	err := db.Get(&total, countSql, args...)
	if err != nil {
		return nil, 0, err
	}

	// 构建排序
	orderBy := "created_at DESC"
	if param.Sort == "popular" {
		if param.Days > 0 {
			// 热门查询需要联表统计
			startDate := time.Now().AddDate(0, 0, -param.Days).Format("2006-01-02")
			conditions = append(conditions, `id IN (
				SELECT article_id FROM article_stats
				WHERE date >= ?
				GROUP BY article_id
				ORDER BY SUM(views) DESC
			)`)
			args = append(args, startDate)
		}
		orderBy = "view_count DESC, created_at DESC"
	}

	// 查询数据
	querySql := fmt.Sprintf(`SELECT
		id, title, content, summary, word_count, author_id, status,
		is_featured, featured_at, allow_comment, like_count, comment_count,
		slug, meta_keywords, meta_description, view_count,
		ip_address, user_agent, extra, created_at, updated_at
		FROM articles
		WHERE %s
		ORDER BY %s
		LIMIT ? OFFSET ?`,
		strings.Join(conditions, " AND "), orderBy)

	args = append(args, param.Size, (param.Page-1)*param.Size)

	err = db.Select(&articles, querySql, args...)
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetFeaturedArticles 获取精选文章
func GetFeaturedArticles(limit int) ([]*models.Article, error) {
	var articles []*models.Article
	sqlStr := `SELECT
		id, title, content, summary, word_count, author_id, status,
		is_featured, featured_at, allow_comment, like_count, comment_count,
		slug, meta_keywords, meta_description, view_count,
		ip_address, user_agent, extra, created_at, updated_at
		FROM articles
		WHERE status = 'published' AND is_featured = TRUE
		ORDER BY featured_at DESC
		LIMIT ?`

	err := db.Select(&articles, sqlStr, limit)
	return articles, err
}

// GetArticlesByTag 根据标签获取文章
func GetArticlesByTag(tagID int64, param *models.ParamArticleList) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	// 计算总数
	countSql := `SELECT COUNT(DISTINCT a.id)
		FROM articles a
		JOIN article_tags at ON a.id = at.article_id
		WHERE at.tag_id = ? AND a.status = 'published'`
	err := db.Get(&total, countSql, tagID)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	querySql := `SELECT DISTINCT
		a.id, a.title, a.content, a.summary, a.word_count, a.author_id, a.status,
		a.is_featured, a.featured_at, a.allow_comment, a.like_count, a.comment_count,
		a.slug, a.meta_keywords, a.meta_description, a.view_count,
		a.ip_address, a.user_agent, a.extra, a.created_at, a.updated_at
		FROM articles a
		JOIN article_tags at ON a.id = at.article_id
		WHERE at.tag_id = ? AND a.status = 'published'
		ORDER BY a.created_at DESC
		LIMIT ? OFFSET ?`

	err = db.Select(&articles, querySql, tagID, param.Size, (param.Page-1)*param.Size)
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// UpdateArticleTags 更新文章标签
func UpdateArticleTags(articleId int64, tagIds []int64) error {
	// 先删除现有标签
	_, err := db.Exec("DELETE FROM article_tags WHERE article_id = ?", articleId)
	if err != nil {
		return err
	}

	// 如果没有新标签，直接返回
	if len(tagIds) == 0 {
		return nil
	}

	// 批量插入新标签
	valueStrings := make([]string, 0, len(tagIds))
	valueArgs := make([]interface{}, 0, len(tagIds)*2)

	for _, tagId := range tagIds {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, articleId, tagId)
	}

	sqlStr := fmt.Sprintf("INSERT INTO article_tags (article_id, tag_id) VALUES %s",
		strings.Join(valueStrings, ", "))

	_, err = db.Exec(sqlStr, valueArgs...)
	return err
}

// UpdateArticleCommentCount 更新文章评论数
func UpdateArticleCommentCount(articleId int64) error {
	sqlStr := `UPDATE articles SET
		comment_count = (SELECT COUNT(*) FROM comments WHERE article_id = ? AND status = 'approved'),
		updated_at = NOW()
		WHERE id = ?`

	_, err := db.Exec(sqlStr, articleId, articleId)
	return err
}

// UpdateArticle 更新文章
func UpdateArticle(id int64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	setParts := make([]string, 0, len(updates))
	args := make([]interface{}, 0, len(updates)+1)

	for field, value := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = ?", field))
		args = append(args, value)
	}

	args = append(args, id)
	sqlStr := fmt.Sprintf("UPDATE articles SET %s, updated_at = NOW() WHERE id = ?",
		strings.Join(setParts, ", "))

	result, err := db.Exec(sqlStr, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrorArticleNotExist
	}

	return nil
}

// DeleteArticle 删除文章（硬删除）
func DeleteArticle(id int64) error {
	result, err := db.Exec("DELETE FROM articles WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrorArticleNotExist
	}

	return nil
}

// UpdateArticleStatus 更新文章状态
func UpdateArticleStatus(id int64, status string) error {
	var sqlStr string
	if status == "published" && status != "draft" {
		// 发布文章时设置发布时间
		sqlStr = "UPDATE articles SET status = ?, updated_at = NOW() WHERE id = ?"
	} else {
		sqlStr = "UPDATE articles SET status = ?, updated_at = NOW() WHERE id = ?"
	}

	result, err := db.Exec(sqlStr, status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrorArticleNotExist
	}

	return nil
}

// UpdateArticleFeatured 更新文章精选状态
func UpdateArticleFeatured(id int64, isFeatured bool) error {
	var featuredAt interface{}
	if isFeatured {
		featuredAt = time.Now()
	} else {
		featuredAt = nil
	}

	sqlStr := "UPDATE articles SET is_featured = ?, featured_at = ?, updated_at = NOW() WHERE id = ?"
	result, err := db.Exec(sqlStr, isFeatured, featuredAt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrorArticleNotExist
	}

	return nil
}

// GetArticlesByAuthor 获取作者的文章列表
func GetArticlesByAuthor(authorID int64, status string, page, size int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	conditions := []string{"author_id = ?"}
	args := []interface{}{authorID}

	if status != "" && status != "all" {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}

	// 计算总数
	q, params, err := sqlx.In("SELECT COUNT(*) FROM articles WHERE "+strings.Join(conditions, " AND "), args...)
	if err != nil {
		return nil, 0, err
	}
	countSql := db.Rebind(q)
	err = db.Get(&total, countSql, params...)
	if err != nil {
		return nil, 0, err
	}

	// 查询数据
	q, params2, err := sqlx.In(`SELECT
		id, title, content, summary, word_count, author_id, status,
		is_featured, featured_at, allow_comment, like_count, comment_count,
		slug, meta_keywords, meta_description, view_count,
		ip_address, user_agent, extra, created_at, updated_at
		FROM articles
		WHERE `+strings.Join(conditions, " AND ")+`
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`, args...)
	if err != nil {
		return nil, 0, err
	}
	querySql := db.Rebind(q)
	queryArgs := append(params2, size, (page-1)*size)
	err = db.Select(&articles, querySql, queryArgs...)
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}