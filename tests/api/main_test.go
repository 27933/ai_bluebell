package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

// 配置信息
const (
	BaseURL = "http://localhost:8084"
	APIVersion = "/api/v1"
)

// HTTP 客户端
var client = &http.Client{
	Timeout: 30 * time.Second,
}

// ResponseData API 响应数据结构
type ResponseData struct {
	Code int64       `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

// TokenResponse Token 响应结构
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// LoginResponseData 登录响应数据
type LoginResponseData struct {
	User  map[string]interface{} `json:"user"`
	Token TokenResponse          `json:"token"`
}

// Article 文章结构
type Article struct {
	ID            int64    `json:"id,string"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	Summary       string   `json:"summary"`
	Status        string   `json:"status"`
	IsFeatured    bool     `json:"is_featured"`
	ViewCount     int      `json:"view_count"`
	LikeCount     int      `json:"like_count"`
	CommentCount  int      `json:"comment_count"`
	AllowComment  bool     `json:"allow_comment"`
	Tags          []Tag    `json:"tags"`
	AuthorID      int64    `json:"author_id,string"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

// Tag 标签结构
type Tag struct {
	ID          int64  `json:"id,string"`
	Name        string `json:"name"`
	ArticleCount int   `json:"article_count"`
}

// Category 栏目结构
type Category struct {
	ID          int64  `json:"id,string"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Comment 评论结构
type Comment struct {
	ID        int64  `json:"id,string"`
	Content   string `json:"content"`
	ArticleID int64  `json:"article_id,string"`
	UserID    int64  `json:"user_id,string"`
	CreatedAt string `json:"created_at"`
}

// User 用户结构
type User struct {
	ID       int64  `json:"id,string"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

// TestMain 测试入口
func TestMain(m *testing.M) {
	// 检查服务是否可用
	if err := checkService(); err != nil {
		fmt.Printf("服务检查失败: %v\n", err)
		fmt.Println("请确保服务已启动: docker exec -it bluebell-ai bash 后运行应用")
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

// checkService 检查服务是否可用
func checkService() error {
	resp, err := client.Get(BaseURL + "/ping")
	if err != nil {
		return fmt.Errorf("无法连接到服务: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("服务返回非200状态码: %d", resp.StatusCode)
	}
	return nil
}

// DoRequest 发送 HTTP 请求
func DoRequest(method, path string, body interface{}, token string) (*ResponseData, int, error) {
	url := BaseURL + APIVersion + path

	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		bodyReader = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	var result ResponseData
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("解析响应失败: %w, body: %s", err, string(respBody))
	}

	return &result, resp.StatusCode, nil
}

// DoGet 发送 GET 请求
func DoGet(path string, token string) (*ResponseData, int, error) {
	return DoRequest(http.MethodGet, path, nil, token)
}

// DoPost 发送 POST 请求
func DoPost(path string, body interface{}, token string) (*ResponseData, int, error) {
	return DoRequest(http.MethodPost, path, body, token)
}

// DoPut 发送 PUT 请求
func DoPut(path string, body interface{}, token string) (*ResponseData, int, error) {
	return DoRequest(http.MethodPut, path, body, token)
}

// DoDelete 发送 DELETE 请求
func DoDelete(path string, token string) (*ResponseData, int, error) {
	return DoRequest(http.MethodDelete, path, nil, token)
}

// DoPatch 发送 PATCH 请求
func DoPatch(path string, body interface{}, token string) (*ResponseData, int, error) {
	return DoRequest(http.MethodPatch, path, body, token)
}

// DoGetWithQuery 发送带查询参数的 GET 请求
func DoGetWithQuery(path string, token string) (*ResponseData, int, error) {
	return DoRequest(http.MethodGet, path, nil, token)
}

// ExtractData 从响应中提取 data 字段
func ExtractData(resp *ResponseData) (map[string]interface{}, error) {
	if resp.Data == nil {
		return nil, nil
	}

	dataBytes, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// AssertSuccess 断言响应成功
func AssertSuccess(t *testing.T, resp *ResponseData, statusCode int) bool {
	t.Helper()

	if statusCode != http.StatusOK {
		t.Errorf("期望状态码200，实际得到: %d", statusCode)
		return false
	}

	if resp.Code != 1000 {
		t.Errorf("期望业务码1000，实际得到: %d, msg: %v", resp.Code, resp.Msg)
		return false
	}

	return true
}

// AssertError 断言响应错误
func AssertError(t *testing.T, resp *ResponseData, statusCode int, expectedCode int64) bool {
	t.Helper()

	if statusCode != http.StatusOK {
		t.Errorf("期望HTTP状态码200，实际得到: %d", statusCode)
		return false
	}

	if resp.Code != expectedCode {
		t.Errorf("期望业务码%d，实际得到: %d, msg: %v", expectedCode, resp.Code, resp.Msg)
		return false
	}

	return true
}

// AssertCode 断言状态码
func AssertCode(t *testing.T, resp *ResponseData, expectedCode int64) bool {
	t.Helper()

	if resp.Code != expectedCode {
		t.Errorf("期望业务码%d，实际得到: %d, msg: %v", expectedCode, resp.Code, resp.Msg)
		return false
	}

	return true
}

// GenerateRandomUsername 生成随机用户名
func GenerateRandomUsername() string {
	return fmt.Sprintf("testuser_%d", time.Now().UnixNano())
}

// GenerateRandomEmail 生成随机邮箱
func GenerateRandomEmail() string {
	return fmt.Sprintf("test_%d@test.com", time.Now().UnixNano())
}
