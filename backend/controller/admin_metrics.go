package controller

import (
	"bluebell/logic"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetSystemMetricsHandler 获取系统性能指标
// @Summary 获取系统性能指标
// @Description 获取系统实时性能指标（CPU、内存、磁盘、活跃用户数）
// @Tags 管理员
// @Accept json
// @Produce json
// @Success 200 {object} models.SystemMetrics "系统性能指标"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/admin/metrics/realtime [get]
func GetSystemMetricsHandler(c *gin.Context) {
	// 获取实时系统指标
	metrics, err := logic.GetSystemMetrics()
	if err != nil {
		zap.L().Error("logic.GetSystemMetrics() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, metrics)
}

// GetSystemMetricsHistoryHandler 获取系统性能历史数据
// @Summary 获取系统性能历史数据
// @Description 获取指定时间段内的系统性能历史数据
// @Tags 管理员
// @Accept json
// @Produce json
// @Param start_time query int false "开始时间戳（秒）"
// @Param end_time query int false "结束时间戳（秒）"
// @Param metric_type query string false "指标类型" default(cpu) enums(cpu,memory,disk)
// @Success 200 {object} controller._ResponseMetricsHistory "历史数据"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/admin/metrics/history [get]
func GetSystemMetricsHistoryHandler(c *gin.Context) {
	// 获取查询参数
	startTimeStr := c.DefaultQuery("start_time", "")
	endTimeStr := c.DefaultQuery("end_time", "")
	metricType := c.DefaultQuery("metric_type", "cpu") // 默认查询CPU

	// 解析时间参数
	var startTime, endTime int64
	var err error

	if startTimeStr != "" {
		startTime, err = strconv.ParseInt(startTimeStr, 10, 64)
		if err != nil {
			ResponseError(c, CodeInvalidParam)
			return
		}
	} else {
		// 默认查询最近1小时
		startTime = time.Now().Add(-1 * time.Hour).Unix()
	}

	if endTimeStr != "" {
		endTime, err = strconv.ParseInt(endTimeStr, 10, 64)
		if err != nil {
			ResponseError(c, CodeInvalidParam)
			return
		}
	} else {
		endTime = time.Now().Unix()
	}

	// 验证时间范围
	if startTime >= endTime {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 获取历史数据（这里简化处理，实际应该从数据库或缓存获取）
	historyData, err := logic.GetSystemMetricsHistory(startTime, endTime, metricType)
	if err != nil {
		zap.L().Error("logic.GetSystemMetricsHistory() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, gin.H{
		"metric_type": metricType,
		"start_time":  startTime,
		"end_time":    endTime,
		"data":        historyData,
	})
}