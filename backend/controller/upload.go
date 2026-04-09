package controller

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// allowedMagicBytes 各图片格式的文件头魔数
var allowedMagicBytes = []struct {
	mime   string
	header []byte
}{
	{"image/jpeg", []byte{0xFF, 0xD8, 0xFF}},
	{"image/png", []byte{0x89, 0x50, 0x4E, 0x47}},
	{"image/gif", []byte{0x47, 0x49, 0x46, 0x38}},
	{"image/webp", []byte{0x52, 0x49, 0x46, 0x46}}, // RIFF....WEBP
}

func checkImageMIME(f io.Reader) bool {
	buf := make([]byte, 8)
	if _, err := io.ReadFull(f, buf); err != nil {
		return false
	}
	for _, magic := range allowedMagicBytes {
		if len(buf) >= len(magic.header) {
			match := true
			for i, b := range magic.header {
				if buf[i] != b {
					match = false
					break
				}
			}
			if match {
				return true
			}
		}
	}
	return false
}

const (
	uploadDir    = "./uploads/images"
	maxImageSize = 5 << 20 // 5MB
)

var allowedImageExts = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true,
	".gif": true, ".webp": true,
}

// UploadImageHandler 上传图片
// @Summary 上传图片
// @Description 上传图片文件，保存到服务器本地
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "图片文件（jpg/jpeg/png/gif/webp，最大5MB）"
// @Success 200 {object} controller._ResponseUpload "上传成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 401 {object} controller._ResponseError "未登录"
// @Failure 500 {object} controller._ResponseError "服务器错误"
// @Security ApiKeyAuth
// @Router /api/v1/upload/image [post]
func UploadImageHandler(c *gin.Context) {
	// 1. 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		zap.L().Error("get upload file failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 校验大小
	if file.Size > maxImageSize {
		ResponseErrorWithMsg(c, CodeInvalidParam, "图片大小不能超过 5MB")
		return
	}

	// 3. 校验扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExts[ext] {
		ResponseErrorWithMsg(c, CodeInvalidParam, "仅支持 jpg/jpeg/png/gif/webp 格式")
		return
	}

	// 4. 校验文件头魔数（防止改扩展名绕过）
	src, err := file.Open()
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	defer src.Close()
	if !checkImageMIME(src) {
		ResponseErrorWithMsg(c, CodeInvalidParam, "文件内容与格式不符，请上传真实图片")
		return
	}

	// 6. 确保目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		zap.L().Error("create upload dir failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 7. 生成唯一文件名，避免覆盖
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join(uploadDir, filename)

	// 8. 保存文件
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		zap.L().Error("save upload file failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 9. 返回可访问的 URL
	ResponseSuccess(c, gin.H{
		"url":      fmt.Sprintf("/uploads/images/%s", filename),
		"filename": filename,
		"size":     file.Size,
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