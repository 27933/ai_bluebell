package controller

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UploadImageHandler 上传图片的处理函数（预留）
// @Summary 上传图片
// @Description 上传图片文件（预留接口）
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "图片文件"
// @Success 200 {object} controller._ResponseUpload "上传成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/upload/image [post]
func UploadImageHandler(c *gin.Context) {
	// 1. 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		zap.L().Error("get upload file failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 校验文件类型和大小
	// TODO: 实现文件类型和大小校验

	// 3. 保存文件
	// 生成文件名：时间戳 + 原始扩展名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().Unix(), ext)

	// TODO: 实现文件保存逻辑
	// 实际项目中应该保存到云存储或指定目录

	// 4. 返回响应
	ResponseSuccess(c, gin.H{
		"url": fmt.Sprintf("/uploads/images/%s", filename),
		"filename": filename,
		"size": file.Size,
	})
}

// UploadAttachmentHandler 上传附件的处理函数（预留）
// @Summary 上传附件
// @Description 上传附件文件（预留接口）
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "附件文件"
// @Success 200 {object} controller._ResponseUpload "上传成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/upload/attachment [post]
func UploadAttachmentHandler(c *gin.Context) {
	// 1. 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		zap.L().Error("get upload file failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 校验文件类型和大小
	// TODO: 实现文件类型和大小校验

	// 3. 保存文件
	// TODO: 实现文件保存逻辑

	// 4. 返回响应
	ResponseSuccess(c, gin.H{
		"url": fmt.Sprintf("/uploads/attachments/%s", file.Filename),
		"filename": file.Filename,
		"size": file.Size,
	})
}