package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// SignUpHandler 处理注册请求的函数
// @Summary 用户注册
// @Description 新用户注册账号
// @Tags 认证
// @Accept json
// @Produce json
// @Param signup body models.ParamSignUp true "注册参数"
// @Success 200 {object} controller._ResponseSuccess "注册成功"
// @Failure 400 {object} controller._ResponseError "参数错误"
// @Failure 400 {object} controller._ResponseError "用户已存在"
// @Router /api/v1/auth/signup [post]
func SignUpHandler(c *gin.Context) {
	// 1. 获取参数和参数校验
	p := new(models.ParamSignUp)
	// ShouldBindJSON() - gin框架中的函数: 将请求中json格式的数据绑定到我们定义的结构体中
	// 			即对于每一个携带请求的路由, 我们都需要定义绑定请求参数的结构体
	//			-> param.go 文件专门定义绑定请求参数的结构体
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		// errs.Translate(trans)) 翻译validator类型错误 即英文错误翻译成中文错误
		// removeTopStruct(errs.Translate(trans)) 用于去掉返回的结构体信息 因为前端没必要知道后端的结构体名字
		// 两个函数合起来就达到返回一个  tag(json):错误信息 的效果
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2. 业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// LoginHandler 登录
func LoginHandlerOld(c *gin.Context) {
	// 转发到auth.go中的LoginHandler
	authLoginHandler(c)
}

func authLoginHandler(c *gin.Context) {
	// 1. 获取参数
	p := new(models.ParamUserLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 业务逻辑处理
	user, err := logic.Login(&models.ParamLogin{
		Username: p.Username,
		Password: p.Password,
	})
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		ResponseError(c, CodeInvalidPassword)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, user)
}
