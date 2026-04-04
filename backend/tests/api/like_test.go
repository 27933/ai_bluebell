package api

import (
	"fmt"
	"testing"
)

// TestLike_Article_NotLogin 测试未登录点赞文章
func TestLike_Article_NotLogin(t *testing.T) {
	body := map[string]interface{}{
		"target_type": "article",
		"target_id":   1,
	}

	resp, statusCode, err := DoPost("/likes", body, "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回需要登录错误
	if !AssertError(t, resp, statusCode, 1006) {
		t.Logf("响应: %+v", resp)
	}
}

// TestLike_Article_Success 测试成功点赞文章
func TestLike_Article_Success(t *testing.T) {
	// 创建作者和文章
	authorToken := createTestUserAndLogin(t)
	if authorToken == "" {
		t.Skip("无法获取作者Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, authorToken, "published")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	// 另一用户点赞
	likerToken := createTestUserAndLogin(t)
	if likerToken == "" {
		t.Skip("无法获取点赞者Token，跳过测试")
		return
	}

	body := map[string]interface{}{
		"target_type": "article",
		"target_id":   articleID,
	}

	resp, statusCode, err := DoPost("/likes", body, likerToken)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestLike_Article_NotExist 测试点赞不存在的文章
func TestLike_Article_NotExist(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	body := map[string]interface{}{
		"target_type": "article",
		"target_id":   999999999,
	}

	resp, statusCode, err := DoPost("/likes", body, token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回参数错误
	if !AssertError(t, resp, statusCode, 1001) {
		t.Logf("响应: %+v", resp)
	}
}

// TestLike_Comment_Success 测试成功点赞评论
func TestLike_Comment_Success(t *testing.T) {
	// 创建文章
	authorToken := createTestUserAndLogin(t)
	if authorToken == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, authorToken, "published")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	// 创建评论
	commentID := createTestComment(t, authorToken, articleID, "测试评论")
	if commentID == 0 {
		t.Skip("无法创建测试评论，跳过测试")
		return
	}

	// 点赞评论
	likerToken := createTestUserAndLogin(t)
	if likerToken == "" {
		t.Skip("无法获取点赞者Token，跳过测试")
		return
	}

	body := map[string]interface{}{
		"target_type": "comment",
		"target_id":   commentID,
	}

	resp, statusCode, err := DoPost("/likes", body, likerToken)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestLike_Duplicate 测试重复点赞
func TestLike_Duplicate(t *testing.T) {
	// 创建文章
	authorToken := createTestUserAndLogin(t)
	if authorToken == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, authorToken, "published")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	// 用户点赞
	likerToken := createTestUserAndLogin(t)
	if likerToken == "" {
		t.Skip("无法获取点赞者Token，跳过测试")
		return
	}

	body := map[string]interface{}{
		"target_type": "article",
		"target_id":   articleID,
	}

	// 第一次点赞
	resp1, _, err := DoPost("/likes", body, likerToken)
	if err != nil {
		t.Fatalf("第一次点赞失败: %v", err)
	}
	if resp1.Code != 1000 {
		t.Fatalf("第一次点赞应成功: %v", resp1.Msg)
	}

	// 第二次点赞（重复）
	resp2, statusCode2, err := DoPost("/likes", body, likerToken)
	if err != nil {
		t.Fatalf("第二次点赞请求失败: %v", err)
	}

	// 重复点赞应该不报错（幂等）
	AssertSuccess(t, resp2, statusCode2)
}

// TestUnlike_NotLogin 测试未登录取消点赞
func TestUnlike_NotLogin(t *testing.T) {
	resp, statusCode, err := DoDeleteWithQuery("/likes?target_type=article&target_id=1", "")
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 期望返回需要登录错误
	if !AssertError(t, resp, statusCode, 1006) {
		t.Logf("响应: %+v", resp)
	}
}

// TestUnlike_Success 测试成功取消点赞
func TestUnlike_Success(t *testing.T) {
	// 创建文章
	authorToken := createTestUserAndLogin(t)
	if authorToken == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, authorToken, "published")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	// 用户点赞
	likerToken := createTestUserAndLogin(t)
	if likerToken == "" {
		t.Skip("无法获取点赞者Token，跳过测试")
		return
	}

	likeBody := map[string]interface{}{
		"target_type": "article",
		"target_id":   articleID,
	}

	_, _, _ = DoPost("/likes", likeBody, likerToken)

	// 取消点赞
	resp, statusCode, err := DoDeleteWithQuery(fmt.Sprintf("/likes?target_type=article&target_id=%d", articleID), likerToken)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// TestUnlike_NotLiked 测试取消未点赞的内容
func TestUnlike_NotLiked(t *testing.T) {
	// 创建文章
	authorToken := createTestUserAndLogin(t)
	if authorToken == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, authorToken, "published")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	// 用户未点赞直接取消
	userToken := createTestUserAndLogin(t)
	if userToken == "" {
		t.Skip("无法获取用户Token，跳过测试")
		return
	}

	resp, statusCode, err := DoDeleteWithQuery(fmt.Sprintf("/likes?target_type=article&target_id=%d", articleID), userToken)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	// 取消未点赞的内容应该不报错
	AssertSuccess(t, resp, statusCode)
}

// TestLikeStatus_Get 测试获取点赞状态
func TestLikeStatus_Get(t *testing.T) {
	// 创建文章
	authorToken := createTestUserAndLogin(t)
	if authorToken == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	articleID := createTestArticle(t, authorToken, "published")
	if articleID == 0 {
		t.Skip("无法创建测试文章，跳过测试")
		return
	}

	// 用户点赞
	likerToken := createTestUserAndLogin(t)
	if likerToken == "" {
		t.Skip("无法获取点赞者Token，跳过测试")
		return
	}

	likeBody := map[string]interface{}{
		"target_type": "article",
		"target_id":   articleID,
	}

	_, _, _ = DoPost("/likes", likeBody, likerToken)

	// 查询点赞状态
	resp, statusCode, err := DoGetWithQuery(fmt.Sprintf("/likes/status?target_type=article&target_id=%d", articleID), likerToken)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
		return
	}

	data, err := ExtractData(resp)
	if err != nil {
		t.Fatalf("提取数据失败: %v", err)
	}

	if data["is_liked"] != true {
		t.Errorf("期望 is_liked=true，实际: %v", data["is_liked"])
	}
}

// TestUserLikes_List 测试获取用户点赞列表
func TestUserLikes_List(t *testing.T) {
	token := createTestUserAndLogin(t)
	if token == "" {
		t.Skip("无法获取Token，跳过测试")
		return
	}

	resp, statusCode, err := DoGetWithQuery("/user/likes?target_type=article", token)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if !AssertSuccess(t, resp, statusCode) {
		t.Logf("响应: %+v", resp)
	}
}

// DoDeleteWithQuery 发送带查询参数的 DELETE 请求
func DoDeleteWithQuery(path string, token string) (*ResponseData, int, error) {
	return DoRequest("DELETE", path, nil, token)
}
