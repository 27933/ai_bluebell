package logic

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"bluebell/dao/mysql"
	"bluebell/models"
)

// ExportArticle 导出单篇文章
func ExportArticle(userID int64, userRole string, articleID int64) (*models.ExportArticleResponse, error) {
	// 获取文章详情
	article, err := mysql.GetArticleById(articleID)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, models.ErrorArticleNotExist
	}

	// 权限验证：作者只能导出自己的文章，管理员可以导出所有文章
	if userRole != "admin" && article.AuthorID != userID {
		return nil, models.ErrorNoPermission
	}

	// 获取作者信息
	author, err := mysql.GetUserById(article.AuthorID)
	if err != nil {
		return nil, err
	}

	// 获取文章标签
	tags, err := mysql.GetTagsByArticleID(articleID)
	if err != nil {
		return nil, err
	}

	// 生成Markdown内容
	content := generateMarkdownContent(article, author, tags)

	// 生成文件名
	filename := generateMarkdownFilename(article.Title, article.CreatedAt)

	return &models.ExportArticleResponse{
		Filename: filename,
		Content:  content,
		Size:     int64(len(content)),
	}, nil
}

// ExportArticlesBatch 批量导出文章
func ExportArticlesBatch(userID int64, userRole string, articleIDs []int64) (*models.ExportBatchResponse, error) {
	// 验证文章数量
	if len(articleIDs) == 0 || len(articleIDs) > 50 {
		return nil, mysql.ErrorInvalidParam
	}

	// 获取所有文章
	articles := make([]*models.Article, 0, len(articleIDs))
	for _, articleID := range articleIDs {
		article, err := mysql.GetArticleById(articleID)
		if err != nil {
			return nil, err
		}
		if article == nil {
			return nil, models.ErrorArticleNotExist
		}

		// 权限验证
		if userRole != "admin" && article.AuthorID != userID {
			return nil, models.ErrorNoPermission
		}

		articles = append(articles, article)
	}

	// 生成批量导出信息
	exportFiles := make([]*models.ExportFileInfo, 0, len(articles))
	totalSize := int64(0)

	for _, article := range articles {
		// 获取作者信息
		author, err := mysql.GetUserById(article.AuthorID)
		if err != nil {
			return nil, err
		}

		// 获取文章标签
		tags, err := mysql.GetTagsByArticleID(article.ID)
		if err != nil {
			return nil, err
		}

		// 生成Markdown内容
		content := generateMarkdownContent(article, author, tags)

		// 生成文件名
		filename := generateMarkdownFilename(article.Title, article.CreatedAt)

		fileInfo := &models.ExportFileInfo{
			Filename: filename,
			Content:  content,
			Size:     int64(len(content)),
		}

		exportFiles = append(exportFiles, fileInfo)
		totalSize += fileInfo.Size
	}

	// 生成批次ID（用于后续下载）
	batchID := fmt.Sprintf("export-%d-%d", time.Now().Unix(), userID)

	// TODO: 将文件信息存储到临时缓存中，供下载使用
	// 这里简化处理，直接返回文件列表

	return &models.ExportBatchResponse{
		BatchID:   batchID,
		Files:     exportFiles,
		FileCount: int64(len(exportFiles)),
		TotalSize: totalSize,
	}, nil
}

// generateMarkdownContent 生成Markdown内容
func generateMarkdownContent(article *models.Article, author *models.User, tags []models.Tag) string {
	var buf bytes.Buffer

	// YAML front matter
	buf.WriteString("---\n")
	buf.WriteString(fmt.Sprintf("title: \"%s\"\n", escapeYAML(article.Title)))
	buf.WriteString(fmt.Sprintf("author: \"%s\"\n", author.Username))
	buf.WriteString(fmt.Sprintf("created_at: \"%s\"\n", article.CreatedAt.Format("2006-01-02 15:04:05")))
	buf.WriteString(fmt.Sprintf("updated_at: \"%s\"\n", article.UpdatedAt.Format("2006-01-02 15:04:05")))

	// 标签
	if len(tags) > 0 {
		tagNames := make([]string, len(tags))
		for i, tag := range tags {
			tagNames[i] = tag.Name
		}
		buf.WriteString(fmt.Sprintf("tags: [%s]\n", strings.Join(tagNames, ", ")))
	}

	buf.WriteString(fmt.Sprintf("word_count: %d\n", article.WordCount))
	buf.WriteString(fmt.Sprintf("view_count: %d\n", article.ViewCount))
	buf.WriteString(fmt.Sprintf("like_count: %d\n", article.LikeCount))
	buf.WriteString(fmt.Sprintf("comment_count: %d\n", article.CommentCount))
	buf.WriteString("---\n\n")

	// 文章标题
	buf.WriteString(fmt.Sprintf("# %s\n\n", article.Title))

	// 文章摘要（如果有）
	if article.Summary != "" {
		buf.WriteString("> ")
		buf.WriteString(strings.ReplaceAll(article.Summary, "\n", "\n> "))
		buf.WriteString("\n\n")
	}

	// 文章正文
	buf.WriteString(article.Content)
	buf.WriteString("\n")

	return buf.String()
}

// generateMarkdownFilename 生成Markdown文件名
func generateMarkdownFilename(title string, createdAt time.Time) string {
	// 清理标题中的特殊字符
	cleanTitle := strings.TrimSpace(title)
	cleanTitle = strings.ReplaceAll(cleanTitle, "/", "-")
	cleanTitle = strings.ReplaceAll(cleanTitle, "\\", "-")
	cleanTitle = strings.ReplaceAll(cleanTitle, ":", "-")
	cleanTitle = strings.ReplaceAll(cleanTitle, "*", "-")
	cleanTitle = strings.ReplaceAll(cleanTitle, "?", "-")
	cleanTitle = strings.ReplaceAll(cleanTitle, "\"", "-")
	cleanTitle = strings.ReplaceAll(cleanTitle, "<", "-")
	cleanTitle = strings.ReplaceAll(cleanTitle, ">", "-")
	cleanTitle = strings.ReplaceAll(cleanTitle, "|", "-")

	// 限制标题长度
	if len(cleanTitle) > 50 {
		cleanTitle = cleanTitle[:50]
	}

	// 生成日期后缀
	dateStr := createdAt.Format("20060102")

	return fmt.Sprintf("%s-%s.md", cleanTitle, dateStr)
}

// escapeYAML 转义YAML中的特殊字符
func escapeYAML(str string) string {
	str = strings.ReplaceAll(str, "\\", "\\\\")
	str = strings.ReplaceAll(str, "\"", "\\\"")
	return str
}