package controller

// ResCode 为每个状态码定义一个类型
type ResCode int64

// 定义状态码 方便在处理的时候直接调用
const (
	CodeSuccess         ResCode = 1000
	CodeInvalidParam    ResCode = 1001
	CodeUserExist       ResCode = 1002
	CodeUserNotExist    ResCode = 1003
	CodeInvalidPassword ResCode = 1004
	CodeServerBusy      ResCode = 1005

	CodeNeedLogin    ResCode = 1006
	CodeInvalidToken ResCode = 1007

	// 栏目相关错误码
	CodeCategoryNotExist         ResCode = 1008
	CodeCategoryExist            ResCode = 1009
	CodeCategoryNameInvalid      ResCode = 1010
	CodeNoPermissionCategory     ResCode = 1011
	CodeArticleAlreadyInCategory ResCode = 1012

	// 导出相关错误码
	CodeNoPermission    ResCode = 1013
	CodeArticleNotExist ResCode = 1014
)

// codeMsgMap 定义一个Map结构用来保存与状态码对应的状态信息
var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNotExist:    "用户名不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeNeedLogin:    "需要登录",
	CodeInvalidToken: "无效的token",

	// 栏目相关错误码
	CodeCategoryNotExist:      "栏目不存在",
	CodeCategoryExist:         "栏目已存在",
	CodeCategoryNameInvalid:   "栏目名称无效",
	CodeNoPermissionCategory:  "无权限操作该栏目",
	CodeArticleAlreadyInCategory: "文章已在该栏目中",

	// 导出相关错误码
	CodeNoPermission:    "没有权限",
	CodeArticleNotExist: "文章不存在",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
