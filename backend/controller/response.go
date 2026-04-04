package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
{
	"code": 10000, // 程序中的错误码
	"msg": xx,     // 提示信息
	"data": {},    // 数据
}

*/

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// ResponseError 封装返回指定结构的错误信息, 便于前端处理
func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		//code.Msg() 根据code 状态码得到响应的信息
		Msg:  code.Msg(),
		Data: nil,
	})
}

// ResponseErrorWithMsg 自定义的错误信息
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}

// _ResponseError 错误响应结构体（用于Swagger文档）
type _ResponseError struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
}

// _ResponseSuccess 成功响应结构体（用于Swagger文档）
type _ResponseSuccess struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// _ResponseLogin 登录响应结构体（用于Swagger文档）
type _ResponseLogin struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		User  interface{} `json:"user"`
		Token interface{} `json:"token"`
	} `json:"data"`
}

// _ResponseUserProfile 用户资料响应结构体（用于Swagger文档）
type _ResponseUserProfile struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

// _ResponseRefreshToken 刷新Token响应结构体（用于Swagger文档）
type _ResponseRefreshToken struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

// _ResponseArticleList 文章列表响应结构体（用于Swagger文档）
type _ResponseArticleList struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		List  interface{} `json:"list"`
		Total int64       `json:"total"`
		Page  int         `json:"page"`
		Size  int         `json:"size"`
		Pages int64       `json:"pages"`
	} `json:"data"`
}

// _ResponseArticleDetail 文章详情响应结构体（用于Swagger文档）
type _ResponseArticleDetail struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

// _ResponseCreateArticle 创建文章响应结构体（用于Swagger文档）
type _ResponseCreateArticle struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

// _ResponseFeaturedArticles 精选文章响应结构体（用于Swagger文档）
type _ResponseFeaturedArticles struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

// _ResponseTrendingArticles 热门文章响应结构体（用于Swagger文档）
type _ResponseTrendingArticles struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		List interface{} `json:"list"`
		Page int         `json:"page"`
		Size int         `json:"size"`
	} `json:"data"`
}

// _ResponseArticleStats 文章统计响应结构体（用于Swagger文档）
type _ResponseArticleStats struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

// _ResponseArticleTrend 文章趋势响应结构体（用于Swagger文档）
type _ResponseArticleTrend struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		ArticleID int64       `json:"article_id"`
		Days        int         `json:"days"`
		GroupBy     string      `json:"group_by"`
		Trend       interface{} `json:"trend"`
	} `json:"data"`
}

// _ResponseArticleDailyStats 文章日统计响应结构体（用于Swagger文档）
type _ResponseArticleDailyStats struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

// _ResponseBatchArticleStats 批量文章统计响应结构体（用于Swagger文档）
type _ResponseBatchArticleStats struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		Stats interface{} `json:"stats"`
	} `json:"data"`
}

// _ResponseAuthorInfo 作者信息响应结构体（用于Swagger文档）
type _ResponseAuthorInfo struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

// _ResponseAuthorArticles 作者文章响应结构体（用于Swagger文档）
type _ResponseAuthorArticles struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		List  interface{} `json:"list"`
		Total int64       `json:"total"`
		Page  int         `json:"page"`
		Size  int         `json:"size"`
		Pages int64       `json:"pages"`
	} `json:"data"`
}

// _ResponseCategoryList 栏目列表响应结构体（用于Swagger文档）
type _ResponseCategoryList struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		List interface{} `json:"list"`
	} `json:"data"`
}

// _ResponseCategory 栏目详情响应结构体（用于Swagger文档）
type _ResponseCategory struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

// _ResponseTagList 标签列表响应结构体（用于Swagger文档）
type _ResponseTagList struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		List interface{} `json:"list"`
	} `json:"data"`
}

// _ResponseTag 标签响应结构体（用于Swagger文档）
type _ResponseTag struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

// _ResponseCommentList 评论列表响应结构体（用于Swagger文档）
type _ResponseCommentList struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		List  interface{} `json:"list"`
		Total int64       `json:"total"`
		Page  int         `json:"page"`
		Size  int         `json:"size"`
		Pages int64       `json:"pages"`
	} `json:"data"`
}

// _ResponseComment 评论响应结构体（用于Swagger文档）
type _ResponseComment struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

// _ResponseLikeStatus 点赞状态响应结构体（用于Swagger文档）
type _ResponseLikeStatus struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		IsLiked bool `json:"is_liked"`
	} `json:"data"`
}

// _ResponseLikeList 点赞列表响应结构体（用于Swagger文档）
type _ResponseLikeList struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		List  interface{} `json:"list"`
		Total int64       `json:"total"`
		Page  int         `json:"page"`
		Size  int         `json:"size"`
		Pages int64       `json:"pages"`
	} `json:"data"`
}

// _ResponseExportArticle 导出文章响应结构体（用于Swagger文档）
type _ResponseExportArticle struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		Filename string `json:"filename"`
		Content  string `json:"content"`
		Size     int64  `json:"size"`
	} `json:"data"`
}

// _ResponseExportBatch 批量导出响应结构体（用于Swagger文档）
type _ResponseExportBatch struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		BatchID string      `json:"batch_id"`
		Files   interface{} `json:"files"`
		Count   int         `json:"file_count"`
		Size    int64       `json:"total_size"`
	} `json:"data"`
}

// _ResponseUpload 上传响应结构体（用于Swagger文档）
type _ResponseUpload struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		URL      string `json:"url"`
		Filename string `json:"filename"`
		Size     int64  `json:"size"`
	} `json:"data"`
}

// _ResponseBatchUpdate 批量更新响应结构体（用于Swagger文档）
type _ResponseBatchUpdate struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		UpdatedCount int `json:"updated_count"`
	} `json:"data"`
}

// _ResponseMetricsHistory 指标历史响应结构体（用于Swagger文档）
type _ResponseMetricsHistory struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		MetricType string      `json:"metric_type"`
		StartTime  int64       `json:"start_time"`
		EndTime    int64       `json:"end_time"`
		Data       interface{} `json:"data"`
	} `json:"data"`
}

// _ResponsePostList 帖子列表响应结构体（用于Swagger文档）
type _ResponsePostList struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		List  interface{} `json:"list"`
		Total int64       `json:"total"`
		Page  int         `json:"page"`
		Size  int         `json:"size"`
		Pages int64       `json:"pages"`
	} `json:"data"`
}

// _ResponseUserList 用户列表响应结构体（用于Swagger文档）
type _ResponseUserList struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		List  interface{} `json:"list"`
		Total int64       `json:"total"`
		Page  int         `json:"page"`
		Size  int         `json:"size"`
		Pages int64       `json:"pages"`
	} `json:"data"`
}

// _ResponseUserDetail 用户详情响应结构体（用于Swagger文档）
type _ResponseUserDetail struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}
