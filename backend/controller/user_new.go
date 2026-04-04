package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetUserListHandler 获取用户列表（管理员用）
// @Summary 获取用户列表
// @Description 管理员获取用户列表，支持角色和状态筛选
// @Tags 管理员
// @Accept json
// @Produce json
// @Param role query string false "用户角色" default(all) enums(all,visitor,reader,author,admin)
// @Param status query string false "用户状态" default(all) enums(all,active,inactive)
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} controller._ResponseUserList "用户列表"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/admin/users [get]
func GetUserListHandler(c *gin.Context) {
	// 1. 获取参数
	role := c.DefaultQuery("role", "all")
	status := c.DefaultQuery("status", "all")
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "20")

	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}

	// 2. 获取用户列表
	users, total, err := logic.GetUserList(role, status, page, size)
	if err != nil {
		zap.L().Error("logic.GetUserList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, gin.H{
		"list": users,
		"total": total,
		"page": page,
		"size": size,
		"pages": (total + int64(size) - 1) / int64(size),
	})
}

// UpdateUserStatusHandler 更新用户状态（管理员用）
// @Summary 更新用户状态
// @Description 管理员更新用户状态（激活/禁用）
// @Tags 管理员
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param status body object true "用户状态"
// @Success 200 {object} controller._ResponseSuccess "更新成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/admin/users/{id}/status [patch]
func UpdateUserStatusHandler(c *gin.Context) {
	// 1. 获取参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	p := new(models.User)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("UpdateUserStatusHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 更新用户状态
	if err := logic.UpdateUserStatus(id, p.Status); err != nil {
		zap.L().Error("logic.UpdateUserStatus() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// GetUserDetailHandler 获取用户详情（管理员用）
// @Summary 获取用户详情
// @Description 管理员获取用户详细信息
// @Tags 管理员
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} controller._ResponseUserDetail "用户详情"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/admin/users/{id} [get]
func GetUserDetailHandler(c *gin.Context) {
	// 1. 获取参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取用户详情
	user, err := logic.GetUserDetail(id)
	if err != nil {
		zap.L().Error("logic.GetUserDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, user)
}

// UpdateUserRoleHandler 更新用户角色（管理员用）
// @Summary 更新用户角色
// @Description 管理员更新用户角色（admin/author/reader）
// @Tags 管理员
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param role body object true "角色信息"
// @Success 200 {object} controller._ResponseSuccess "更新成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/admin/users/{id}/role [patch]
func UpdateUserRoleHandler(c *gin.Context) {
	// 1. 获取参数
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 解析请求参数
	var p struct {
		Role string `json:"role" binding:"required,oneof=admin author reader"`
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("UpdateUserRoleHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 3. 更新用户角色
	if err := logic.UpdateUserRole(id, p.Role); err != nil {
		zap.L().Error("logic.UpdateUserRole() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 4. 返回响应
	ResponseSuccess(c, nil)
}

// BatchUpdateUserStatusHandler 批量更新用户状态（管理员用）
// @Summary 批量更新用户状态
// @Description 管理员批量更新用户状态（激活/禁用）
// @Tags 管理员
// @Accept json
// @Produce json
// @Param request body object true "批量更新请求"
// @Success 200 {object} controller._ResponseBatchUpdate "更新成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 403 {object} controller._ResponseError "无权限"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/admin/users/batch/status [patch]
func BatchUpdateUserStatusHandler(c *gin.Context) {
	// 1. 解析请求参数
	var p struct {
		UserIDs []int64 `json:"user_ids" binding:"required,min=1,max=100"`
		Status  string  `json:"status" binding:"required,oneof=active inactive"`
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("BatchUpdateUserStatusHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 批量更新用户状态
	if err := logic.BatchUpdateUserStatus(p.UserIDs, p.Status); err != nil {
		zap.L().Error("logic.BatchUpdateUserStatus() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, gin.H{
		"updated_count": len(p.UserIDs),
	})
}