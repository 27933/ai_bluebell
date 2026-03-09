package mysql

import (
	"bluebell/models"
	"database/sql"
	"fmt"
	"strings"
)

// CreateTag 创建标签
func CreateTag(tag *models.Tag) error {
	sqlStr := `INSERT INTO tags (name, description, slug) VALUES (?, ?, ?)`
	result, err := db.Exec(sqlStr, tag.Name, tag.Description, tag.Slug)
	if err != nil {
		return err
	}

	tagID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	tag.ID = tagID
	return nil
}

// GetTagById 根据ID获取标签
func GetTagById(id int64) (*models.Tag, error) {
	tag := new(models.Tag)
	sqlStr := `SELECT id, name, description, slug, created_at, updated_at FROM tags WHERE id = ?`
	err := db.Get(tag, sqlStr, id)
	if err == sql.ErrNoRows {
		return nil, ErrorTagNotExist
	}
	return tag, err
}

// GetTagByName 根据名称获取标签
func GetTagByName(name string) (*models.Tag, error) {
	tag := new(models.Tag)
	sqlStr := `SELECT id, name, description, slug, created_at, updated_at FROM tags WHERE name = ?`
	err := db.Get(tag, sqlStr, name)
	if err == sql.ErrNoRows {
		return nil, ErrorTagNotExist
	}
	return tag, err
}

// GetAllTags 获取所有标签
func GetAllTags() ([]*models.Tag, error) {
	var tags []*models.Tag
	sqlStr := `SELECT id, name, description, slug, created_at, updated_at FROM tags ORDER BY name`
	err := db.Select(&tags, sqlStr)
	return tags, err
}

// GetTagsByArticleId 获取文章的所有标签
func GetTagsByArticleId(articleId int64) ([]*models.Tag, error) {
	var tags []*models.Tag
	sqlStr := `SELECT t.id, t.name, t.description, t.slug, t.created_at, t.updated_at
		FROM tags t
		JOIN article_tags at ON t.id = at.tag_id
		WHERE at.article_id = ?
		ORDER BY t.name`

	err := db.Select(&tags, sqlStr, articleId)
	return tags, err
}

// GetTagsByAuthorId 获取作者使用的所有标签
func GetTagsByAuthorId(authorId int64) ([]*models.Tag, error) {
	var tags []*models.Tag
	sqlStr := `SELECT DISTINCT t.id, t.name, t.description, t.slug, t.created_at, t.updated_at
		FROM tags t
		JOIN article_tags at ON t.id = at.tag_id
		JOIN articles a ON at.article_id = a.id
		WHERE a.author_id = ?
		ORDER BY t.name`

	err := db.Select(&tags, sqlStr, authorId)
	return tags, err
}

// UpdateTag 更新标签
func UpdateTag(id int64, updates map[string]interface{}) error {
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
	sqlStr := fmt.Sprintf("UPDATE tags SET %s, updated_at = NOW() WHERE id = ?",
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
		return ErrorTagNotExist
	}

	return nil
}

// DeleteTag 删除标签
func DeleteTag(id int64) error {
	// 检查标签是否被使用
	var count int64
	checkSql := `SELECT COUNT(*) FROM article_tags WHERE tag_id = ?`
	err := db.Get(&count, checkSql, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrorTagInUse
	}

	result, err := db.Exec("DELETE FROM tags WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrorTagNotExist
	}

	return nil
}

// GetTagsByArticleID 获取文章的标签列表
func GetTagsByArticleID(articleID int64) ([]models.Tag, error) {
	var tags []models.Tag
	sqlStr := `SELECT t.id, t.name, t.description, t.slug, t.created_at, t.updated_at
				FROM tags t
				INNER JOIN article_tags at ON t.id = at.tag_id
				WHERE at.article_id = ?
				ORDER BY t.name`
	err := db.Select(&tags, sqlStr, articleID)
	return tags, err
}

// GetTagStats 获取标签统计信息
func GetTagStats() ([]*models.ApiTag, error) {
	var tags []*models.ApiTag
	sqlStr := `SELECT
		t.id, t.name, t.description, t.slug, t.created_at,
		COUNT(DISTINCT at.article_id) as article_count
		FROM tags t
		LEFT JOIN article_tags at ON t.id = at.tag_id
		LEFT JOIN articles a ON at.article_id = a.id AND a.status = 'published'
		GROUP BY t.id
		ORDER BY article_count DESC, t.name`

	err := db.Select(&tags, sqlStr)
	return tags, err
}

// AddArticleTags 为文章添加标签
func AddArticleTags(articleId int64, tagIds []int64) error {
	if len(tagIds) == 0 {
		return nil
	}

	// 构建批量插入SQL
	valueStrings := make([]string, 0, len(tagIds))
	valueArgs := make([]interface{}, 0, len(tagIds)*2)

	for _, tagId := range tagIds {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, articleId, tagId)
	}

	sqlStr := fmt.Sprintf("INSERT INTO article_tags (article_id, tag_id) VALUES %s ON DUPLICATE KEY UPDATE article_id = article_id",
		strings.Join(valueStrings, ", "))

	_, err := db.Exec(sqlStr, valueArgs...)
	return err
}

// RemoveArticleTags 移除文章的所有标签
func RemoveArticleTags(articleId int64) error {
	_, err := db.Exec("DELETE FROM article_tags WHERE article_id = ?", articleId)
	return err
}

// GetArticleTagIds 获取文章的所有标签ID
func GetArticleTagIds(articleId int64) ([]int64, error) {
	var tagIds []int64
	sqlStr := `SELECT tag_id FROM article_tags WHERE article_id = ?`
	err := db.Select(&tagIds, sqlStr, articleId)
	return tagIds, err
}

// GetTagBySlug 根据slug获取标签
func GetTagBySlug(slug string) (*models.Tag, error) {
	tag := new(models.Tag)
	sqlStr := `SELECT id, name, description, slug, created_at, updated_at FROM tags WHERE slug = ?`
	err := db.Get(tag, sqlStr, slug)
	if err == sql.ErrNoRows {
		return nil, ErrorTagNotExist
	}
	return tag, err
}

// GetAuthorTags 获取作者使用的所有标签
func GetAuthorTags(authorId int64) ([]*models.Tag, error) {
	var tags []*models.Tag
	sqlStr := `SELECT DISTINCT t.id, t.name, t.description, t.slug, t.created_at, t.updated_at
		FROM tags t
		JOIN article_tags at ON t.id = at.tag_id
		JOIN articles a ON at.article_id = a.id
		WHERE a.author_id = ? AND a.status = 'published'
		ORDER BY t.name`

	err := db.Select(&tags, sqlStr, authorId)
	return tags, err
}

// CountArticlesByTag 统计标签下的文章数量
func CountArticlesByTag(tagID int64) (int64, error) {
	var count int64
	sqlStr := `SELECT COUNT(DISTINCT a.id)
		FROM articles a
		JOIN article_tags at ON a.id = at.article_id
		WHERE at.tag_id = ? AND a.status = 'published'`

	err := db.Get(&count, sqlStr, tagID)
	return count, err
}