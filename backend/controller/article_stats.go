package controller

import (
	"fmt"
	"time"

	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetTrendingArticlesHandler 获取热门文章排行的处理函数
// @Summary 获取热门文章排行
// @Description 获取热门文章排行（日榜/周榜/月榜）
// @Tags 统计
// @Accept json
// @Produce json
// @Param period query string false "排行周期" default(daily) enums(daily,weekly,monthly)
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} controller._ResponseTrendingArticles "热门文章列表"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/articles/trending [get]
func GetTrendingArticlesHandler(c *gin.Context) {
	// 1. 获取参数
	p := &models.TrendingArticleQuery{
		PeriodType: "daily",
		Page:       1,
		Size:       20,
	}

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetTrendingArticlesHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取热门文章排行
	articles, err := logic.GetTrendingArticles(p.PeriodType, p.Page, p.Size)
	if err != nil {
		zap.L().Error("logic.GetTrendingArticles() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, gin.H{
		"list": articles,
		"page": p.Page,
		"size": p.Size,
	})
}

// GetAuthorTrendHandler 获取作者所有文章的访问趋势汇总
// @Summary 获取作者访问趋势
// @Description 获取当前登录作者所有文章的每日访问量汇总（用于仪表板折线图）
// @Tags 统计
// @Security ApiKeyAuth
// @Param days query int false "统计天数" default(7) minimum(1) maximum(90)
// @Success 200 {object} controller._ResponseArticleTrend "访问趋势数据"
// @Router /api/v1/author/stats/trend [get]
func GetAuthorTrendHandler(c *gin.Context) {
	authorID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	days := 7
	if d := c.Query("days"); d != "" {
		if _, err := fmt.Sscan(d, &days); err != nil || days < 1 || days > 90 {
			days = 7
		}
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	trendData, err := mysql.GetAuthorDailyTrend(authorID, startDate, endDate)
	if err != nil {
		zap.L().Error("mysql.GetAuthorDailyTrend failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, gin.H{
		"days":  days,
		"trend": trendData,
	})
}

// GetArticleTrendHandler 获取文章访问趋势的处理函数
// @Summary 获取文章访问趋势
// @Description 获取指定文章的访问趋势数据
// @Tags 统计
// @Accept json
// @Produce json
// @Param article_id query int true "文章ID"
// @Param days query int false "统计天数" default(30) minimum(1) maximum(90)
// @Param group_by query string false "分组方式" default(day) enums(hour,day,week,month)
// @Success 200 {object} controller._ResponseArticleTrend "访问趋势数据"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/article-stats/trend [get]
func GetArticleTrendHandler(c *gin.Context) {
	// 1. 获取参数
	p := &models.ArticleStatsQuery{}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetArticleTrendHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取文章访问趋势
	trendData, err := logic.GetArticleTrendData(p.ArticleID, p.Days, p.GroupBy)
	if err != nil {
		zap.L().Error("logic.GetArticleTrendData() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, gin.H{
		"article_id": p.ArticleID,
		"days":       p.Days,
		"group_by":   p.GroupBy,
		"trend":      trendData,
	})
}

// GetArticleStatsHandler 获取文章统计数据的处理函数
// @Summary 获取文章统计数据
// @Description 获取指定文章的统计数据（浏览量、独立访客）
// @Tags 统计
// @Accept json
// @Produce json
// @Param article_id query int true "文章ID"
// @Param days query int false "统计天数" default(30) minimum(1) maximum(90)
// @Param fields query string false "返回字段" default(views,uv) enums(views,uv)
// @Success 200 {object} controller._ResponseArticleStats "文章统计数据"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/article-stats [get]
func GetArticleStatsHandler(c *gin.Context) {
	// 1. 获取参数
	articleIDStr := c.Query("article_id")
	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil {
		zap.L().Error("GetArticleStatsHandler with invalid article_id", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days < 1 || days > 90 {
		zap.L().Error("GetArticleStatsHandler with invalid days", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取文章统计数据
	stats, err := logic.GetArticleStatsWithUV(articleID, days)
	if err != nil {
		zap.L().Error("logic.GetArticleStatsWithUV() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, stats)
}

// BatchGetArticleStatsHandler 批量获取文章统计数据的处理函数
// @Summary 批量获取文章统计数据
// @Description 批量获取多篇文章的统计数据
// @Tags 统计
// @Accept json
// @Produce json
// @Param ids query string true "文章ID列表，逗号分隔"
// @Success 200 {object} controller._ResponseBatchArticleStats "批量统计数据"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/article-stats/batch [get]
func BatchGetArticleStatsHandler(c *gin.Context) {
	// 1. 获取参数
	idsStr := c.Query("ids")
	if idsStr == "" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 解析文章ID列表
	idStrs := strings.Split(idsStr, ",")
	var articleIDs []int64
	for _, idStr := range idStrs {
		id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
		if err != nil {
			zap.L().Error("BatchGetArticleStatsHandler with invalid id", zap.String("id", idStr))
			ResponseError(c, CodeInvalidParam)
			return
		}
		articleIDs = append(articleIDs, id)
	}

	if len(articleIDs) == 0 || len(articleIDs) > 100 {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 批量获取文章统计数据
	stats, err := logic.BatchGetArticleStats(articleIDs)
	if err != nil {
		zap.L().Error("logic.BatchGetArticleStats() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, gin.H{
		"stats": stats,
	})
}

// RecordArticleViewWithAntiCheatHandler 使用防刷机制记录文章访问的处理函数
// @Summary 记录文章访问（防刷）
// @Description 记录文章访问，带防刷机制（同一IP每天最多10次）
// @Tags 统计
// @Accept json
// @Produce json
// @Param article_id query int true "文章ID"
// @Success 200 {object} controller._ResponseSuccess "记录成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Router /api/v1/articles/view [post]
func RecordArticleViewWithAntiCheatHandler(c *gin.Context) {
	// 1. 获取参数
	articleIDStr := c.Query("article_id")
	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil {
		zap.L().Error("RecordArticleViewWithAntiCheatHandler with invalid article_id", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取用户ID（如果有）
	var userID *int64
	if currentUserID, err := getCurrentUserID(c); err == nil {
		userID = &currentUserID
	}

	// 获取客户端IP
	ip := c.ClientIP()

	// 2. 记录文章访问
	if err := logic.RecordArticleViewWithAntiCheat(articleID, userID, ip); err != nil {
		zap.L().Error("logic.RecordArticleViewWithAntiCheat() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, gin.H{
		"message": "访问记录成功",
	})
}