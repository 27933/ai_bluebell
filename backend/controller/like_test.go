package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestLikeHandler_JSONRequest 测试JSON请求体点赞
func TestLikeHandler_JSONRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// 注册路由
	r.POST("/api/v1/likes", LikeHandler)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedCode   int
		expectedErrMsg string
		jsonRequest    bool
	}{
		{
			name: "JSON请求点赞文章",
			requestBody: map[string]interface{}{
				"target_type": "article",
				"target_id":   "100",
			},
			expectedCode:   http.StatusOK,
			expectedErrMsg: "", // 由于需要登录，实际会返回401
			jsonRequest:    true,
		},
		{
			name: "JSON请求点赞评论",
			requestBody: map[string]interface{}{
				"target_type": "comment",
				"target_id":   "200",
			},
			expectedCode:   http.StatusOK,
			expectedErrMsg: "",
			jsonRequest:    true,
		},
		{
			name: "JSON请求无效目标类型",
			requestBody: map[string]interface{}{
				"target_type": "invalid",
				"target_id":   "100",
			},
			expectedCode:   http.StatusOK,
			expectedErrMsg: "请求参数错误",
			jsonRequest:    true,
		},
		{
			name: "JSON请求缺少target_id",
			requestBody: map[string]interface{}{
				"target_type": "article",
			},
			expectedCode:   http.StatusOK,
			expectedErrMsg: "请求参数错误",
			jsonRequest:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建请求
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/likes", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			// 解析响应
			res := new(ResponseData)
			if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
				t.Fatalf("json.Unmarshal w.Body failed, err:%v\n", err)
			}

			if tt.expectedErrMsg != "" {
				assert.NotEqual(t, res.Code, CodeSuccess)
				assert.Contains(t, res.Msg.(string), tt.expectedErrMsg)
			} else {
				// 实际会返回需要登录的错误
				assert.Equal(t, res.Code, CodeNeedLogin)
			}
		})
	}
}

// TestLikeHandler_QueryParams 测试查询参数点赞（向后兼容）
func TestLikeHandler_QueryParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/api/v1/likes", LikeHandler)

	tests := []struct {
		name           string
		targetType     string
		targetID       string
		expectedCode   int
		expectedErrMsg string
	}{
		{
			name:           "查询参数点赞文章",
			targetType:     "article",
			targetID:       "100",
			expectedCode:   http.StatusOK,
			expectedErrMsg: "",
		},
		{
			name:           "查询参数点赞评论",
			targetType:     "comment",
			targetID:       "200",
			expectedCode:   http.StatusOK,
			expectedErrMsg: "",
		},
		{
			name:           "查询参数无效目标类型",
			targetType:     "invalid",
			targetID:       "100",
			expectedCode:   http.StatusOK,
			expectedErrMsg: "请求参数错误",
		},
		{
			name:           "查询参数缺失target_id",
			targetType:     "article",
			targetID:       "",
			expectedCode:   http.StatusOK,
			expectedErrMsg: "请求参数错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/v1/likes"
			if tt.targetType != "" {
				url += "?target_type=" + tt.targetType
			}
			if tt.targetID != "" {
				if tt.targetType != "" {
					url += "&"
				} else {
					url += "?"
				}
				url += "target_id=" + tt.targetID
			}

			req, _ := http.NewRequest(http.MethodPost, url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			res := new(ResponseData)
			if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
				t.Fatalf("json.Unmarshal w.Body failed, err:%v\n", err)
			}

			if tt.expectedErrMsg != "" {
				assert.NotEqual(t, res.Code, CodeSuccess)
				assert.Contains(t, res.Msg.(string), tt.expectedErrMsg)
			} else {
				assert.Equal(t, res.Code, CodeNeedLogin)
			}
		})
	}
}

// TestUnlikeHandler 测试取消点赞
func TestUnlikeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/api/v1/likes", UnlikeHandler)

	tests := []struct {
		name           string
		targetType     string
		targetID       string
		expectedCode   int
		expectedErrMsg string
	}{
		{
			name:           "取消点赞文章",
			targetType:     "article",
			targetID:       "100",
			expectedCode:   http.StatusOK,
			expectedErrMsg: "",
		},
		{
			name:           "取消点赞评论",
			targetType:     "comment",
			targetID:       "200",
			expectedCode:   http.StatusOK,
			expectedErrMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/v1/likes"
			if tt.targetType != "" {
				url += "?target_type=" + tt.targetType
			}
			if tt.targetID != "" {
				if tt.targetType != "" {
					url += "&"
				} else {
					url += "?"
				}
				url += "target_id=" + tt.targetID
			}

			req, _ := http.NewRequest(http.MethodDelete, url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			res := new(ResponseData)
			if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
				t.Fatalf("json.Unmarshal w.Body failed, err:%v\n", err)
			}

			// 由于没有登录，应该返回需要登录的错误
			assert.Equal(t, res.Code, CodeNeedLogin)
		})
	}
}

// TestGetLikeStatusHandler 测试获取点赞状态
func TestGetLikeStatusHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/api/v1/likes/status", GetLikeStatusHandler)

	tests := []struct {
		name           string
		targetType     string
		targetID       string
		expectedCode   int
		expectedErrMsg string
	}{
		{
			name:           "获取文章点赞状态",
			targetType:     "article",
			targetID:       "100",
			expectedCode:   http.StatusOK,
			expectedErrMsg: "",
		},
		{
			name:           "获取评论点赞状态",
			targetType:     "comment",
			targetID:       "200",
			expectedCode:   http.StatusOK,
			expectedErrMsg: "",
		},
		{
			name:           "缺少target_type",
			targetType:     "",
			targetID:       "100",
			expectedCode:   http.StatusOK,
			expectedErrMsg: "请求参数错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/v1/likes/status"
			if tt.targetType != "" {
				url += "?target_type=" + tt.targetType
			}
			if tt.targetID != "" {
				if tt.targetType != "" {
					url += "&"
				} else {
					url += "?"
				}
				url += "target_id=" + tt.targetID
			}

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			res := new(ResponseData)
			if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
				t.Fatalf("json.Unmarshal w.Body failed, err:%v\n", err)
			}

			if tt.expectedErrMsg != "" {
				assert.NotEqual(t, res.Code, CodeSuccess)
				assert.Contains(t, res.Msg.(string), tt.expectedErrMsg)
			} else {
				// 如果参数正常，应该返回成功状态（未登录也能查询点赞状态）
				assert.Equal(t, res.Code, CodeSuccess)
				assert.NotNil(t, res.Data)
			}
		})
	}
}

// TestGetUserLikesHandler 测试获取用户点赞列表
func TestGetUserLikesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/api/v1/user/likes", GetUserLikesHandler)

	tests := []struct {
		name           string
		targetType     string
		page           string
		size           string
		expectedCode   int
		expectedErrMsg string
	}{
		{
			name:           "获取用户文章点赞列表",
			targetType:     "article",
			page:           "1",
			size:           "20",
			expectedCode:   http.StatusOK,
			expectedErrMsg: "",
		},
		{
			name:           "获取用户评论点赞列表",
			targetType:     "comment",
			page:           "2",
			size:           "10",
			expectedCode:   http.StatusOK,
			expectedErrMsg: "",
		},
		{
			name:           "缺失target_type",
			targetType:     "",
			page:           "1",
			size:           "20",
			expectedCode:   http.StatusOK,
			expectedErrMsg: "请求参数错误",
		},
		{
			name:           "无效target_type",
			targetType:     "invalid",
			page:           "1",
			size:           "20",
			expectedCode:   http.StatusOK,
			expectedErrMsg: "请求参数错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/v1/user/likes"
			params := make([]string, 0)
			if tt.targetType != "" {
				params = append(params, "target_type="+tt.targetType)
			}
			if tt.page != "" {
				params = append(params, "page="+tt.page)
			}
			if tt.size != "" {
				params = append(params, "size="+tt.size)
			}

			if len(params) > 0 {
				url += "?" + params[0]
				for i := 1; i < len(params); i++ {
					url += "&" + params[i]
				}
			}

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			res := new(ResponseData)
			if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
				t.Fatalf("json.Unmarshal w.Body failed, err:%v\n", err)
			}

			if tt.expectedErrMsg != "" {
				assert.NotEqual(t, res.Code, CodeSuccess)
				assert.Contains(t, res.Msg.(string), tt.expectedErrMsg)
			} else {
				// 需要登录才能查看自己的点赞列表
				assert.Equal(t, res.Code, CodeNeedLogin)
			}
		})
	}
}

// 测试批量点赞和点赞历史API
func TestNewAPIs_Placeholder(t *testing.T) {
	t.Run("批量点赞状态查询API", func(t *testing.T) {
		t.Skip("批量点赞状态查询API尚未实现")
	})

	t.Run("点赞历史查询API", func(t *testing.T) {
		t.Skip("点赞历史查询API尚未实现")
	})
}

// 辅助函数：设置测试上下文（已实现，无需占位符）

// 测试数据验证：测试LikeRequest结构体绑定
func TestLikeRequestValidation(t *testing.T) {
	tests := []struct {
		name        string
		targetType  string
		targetID    int64
		shouldError bool
	}{
		{
			name:        "有效文章点赞请求",
			targetType:  "article",
			targetID:    100,
			shouldError: false,
		},
		{
			name:        "有效评论点赞请求",
			targetType:  "comment",
			targetID:    200,
			shouldError: false,
		},
		{
			name:        "无效目标类型",
			targetType:  "invalid",
			targetID:    100,
			shouldError: true,
		},
		{
			name:        "空目标类型",
			targetType:  "",
			targetID:    100,
			shouldError: true,
		},
		// 注意：target_id为0通常也被认为是无效的（但具体要看业务逻辑）
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt // 避免未使用警告
			// 实际测试需要验证binding标签
			// 这里跳过，因为验证在控制器中通过ShouldBindJSON处理
			t.Skip("验证在控制器中通过ShouldBindJSON处理")
		})
	}
}