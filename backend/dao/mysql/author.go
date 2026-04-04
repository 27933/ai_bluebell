package mysql

import (
	"bluebell/models"
)

// GetAuthorInfoByUsername 根据用户名获取作者信息
func GetAuthorInfoByUsername(username string) (*models.AuthorInfoResponse, error) {
	var authorInfo models.AuthorInfoResponse

	// 使用JOIN查询聚合作者信息
	sql := `SELECT
			u.username,
			u.nickname,
			u.bio,
			DATE_FORMAT(u.created_at, '%Y-%m-%d') as join_date,
			COUNT(a.id) as article_count,
			COALESCE(SUM(a.view_count), 0) as total_views,
			COALESCE(SUM(a.like_count), 0) as total_likes
		FROM users u
		LEFT JOIN articles a ON u.id = a.author_id AND a.status = 'published'
		WHERE u.username = ?
		GROUP BY u.id, u.username, u.nickname, u.bio, u.created_at`

	err := db.Get(&authorInfo, sql, username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil // 用户不存在
		}
		return nil, err
	}

	return &authorInfo, nil
}