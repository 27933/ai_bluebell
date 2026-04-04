package controller

import (
	"bluebell/logic"
	"bluebell/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetSystemOverviewHandler 获取系统概览的处理函数
// @Summary 获取系统概览
// @Description 获取系统概览统计数据（用户、文章、评论总数等）
// @Tags 管理员
// @Accept json
// @Produce json
// @Success 200 {object} models.SystemOverview "系统概览数据"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/admin/stats/overview [get]
func GetSystemOverviewHandler(c *gin.Context) {
	// 获取系统概览数据
	overview, err := logic.GetSystemOverview()
	if err != nil {
		zap.L().Error("logic.GetSystemOverview() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(c, overview)
}

// GetSystemDailyStatsHandler 获取系统日统计的处理函数
// @Summary 获取系统日统计
// @Description 获取系统每日统计数据（用户、文章、评论新增数量）
// @Tags 管理员
// @Accept json
// @Produce json
// @Param start_date query string false "开始日期（YYYY-MM-DD）"
// @Param end_date query string false "结束日期（YYYY-MM-DD）"
// @Success 200 {object} models.DailyStats "日统计数据"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/admin/stats/daily [get]
func GetSystemDailyStatsHandler(c *gin.Context) {
	// 获取查询参数
	p := &models.ParamStatsDaily{
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
	}

	// 获取日统计数据
	stats, err := logic.GetSystemDailyStats(p)
	if err != nil {
		zap.L().Error("logic.GetSystemDailyStats() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(c, stats)
}