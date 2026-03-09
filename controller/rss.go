package controller

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"bluebell/logic"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RSSItem RSS条目
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Author      string `xml:"author"`
}

// RSSChannel RSS频道
type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	PubDate     string    `xml:"pubDate"`
	Items       []RSSItem `xml:"item"`
}

// RSS RSS根元素
type RSS struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel RSSChannel `xml:"channel"`
}

// GetRSSHandler 获取RSS订阅的处理函数
// @Summary 获取RSS订阅
// @Description 获取最新文章的RSS订阅
// @Tags RSS
// @Accept json
// @Produce xml
// @Success 200 {string} string "RSS XML"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/rss [get]
func GetRSSHandler(c *gin.Context) {
	// 1. 获取文章列表（最新的20篇已发布文章）
	articles, err := logic.GetRSSArticles(20)
	if err != nil {
		zap.L().Error("logic.GetRSSArticles() failed", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error generating RSS feed")
		return
	}

	// 2. 构建RSS内容
	items := make([]RSSItem, 0, len(articles))
	for _, article := range articles {
		item := RSSItem{
			Title:       article.Title,
			Link:        fmt.Sprintf("https://example.com/articles/%d", article.ID),
			Description: article.Summary,
			PubDate:     article.CreatedAt.Format(time.RFC1123Z),
			Author:      fmt.Sprintf("Author%d", article.AuthorID), // 简化处理，使用作者ID
		}
		items = append(items, item)
	}

	rss := RSS{
		Version: "2.0",
		Channel: RSSChannel{
			Title:       "知识博客",
			Link:        "https://example.com",
			Description: "知识博客 - 分享技术，传播知识",
			Language:    "zh-CN",
			PubDate:     time.Now().Format(time.RFC1123Z),
			Items:       items,
		},
	}

	// 3. 返回JSON格式的RSS数据
	ResponseSuccess(c, gin.H{
		"title":       rss.Channel.Title,
		"link":        rss.Channel.Link,
		"description": rss.Channel.Description,
		"language":    rss.Channel.Language,
		"pub_date":    rss.Channel.PubDate,
		"items":       rss.Channel.Items,
	})
}
