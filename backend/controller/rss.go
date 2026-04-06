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

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Author      string `xml:"author"`
}

type rssChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []rssItem `xml:"item"`
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

// GetRSSHandler 返回标准 RSS 2.0 XML，可直接被 RSS 阅读器订阅
func GetRSSHandler(c *gin.Context) {
	articles, err := logic.GetRSSArticles(20)
	if err != nil {
		zap.L().Error("logic.GetRSSArticles() failed", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error generating RSS feed")
		return
	}

	// 用请求 Host 构建链接，兼容本地开发和生产环境
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)

	items := make([]rssItem, 0, len(articles))
	for _, a := range articles {
		author := a.AuthorUsername
		if a.AuthorNickname != "" {
			author = a.AuthorNickname
		}
		items = append(items, rssItem{
			Title:       a.Title,
			Link:        fmt.Sprintf("%s/article/%d", baseURL, a.ID),
			Description: a.Summary,
			PubDate:     a.CreatedAt.Format(time.RFC1123Z),
			Author:      author,
		})
	}

	feed := rssFeed{
		Version: "2.0",
		Channel: rssChannel{
			Title:         "Bluebell",
			Link:          baseURL,
			Description:   "Bluebell - 分享技术，传播知识",
			Language:      "zh-CN",
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			Items:         items,
		},
	}

	output, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error generating RSS feed")
		return
	}

	c.Data(http.StatusOK, "application/rss+xml; charset=utf-8",
		append([]byte(xml.Header), output...))
}
