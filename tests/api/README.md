# Bluebell API 集成测试

## 概述

本目录包含 Bluebell 博客系统的完整 API 集成测试套件，覆盖所有公开接口、认证接口和管理员接口。

## 测试结构

```
tests/api/
├── main_test.go       # 测试基础设施（HTTP客户端、公共函数）
├── auth_test.go       # 认证相关接口测试
├── article_test.go    # 文章相关接口测试
├── comment_test.go    # 评论相关接口测试
├── like_test.go       # 点赞相关接口测试
├── category_test.go   # 栏目相关接口测试
├── tag_test.go        # 标签相关接口测试
├── admin_test.go      # 管理员接口测试
└── README.md          # 本文档
```

## 测试重点

测试主要验证以下业务逻辑：

1. **权限控制**
   - 未登录用户无法访问需要认证的接口
   - 普通用户无法访问管理员接口
   - 非作者无法更新/删除他人的文章/评论

2. **状态管理**
   - 文章状态的转换（草稿 -> 发布 -> 下线）
   - 评论的增删改查
   - 点赞的状态管理

3. **数据验证**
   - 必填字段验证
   - 唯一性验证（用户名、栏目名）
   - 关联数据存在性验证

4. **错误处理**
   - 不存在的资源返回 404 或适当的错误码
   - 无效的参数返回 400 错误
   - 权限不足返回 401/403 错误

## 环境要求

- Docker 容器 `bluebell-ai` 正在运行
- API 服务监听在 `localhost:8084`
- MySQL 和 Redis 服务可访问

## 运行测试

### 在主机上运行

```bash
# 进入项目目录
cd /root/code/bluebell

# 运行所有测试
go test -v ./tests/api/... -count=1

# 只运行认证相关测试
go test -v ./tests/api -run TestAuth -count=1

# 只运行文章相关测试
go test -v ./tests/api -run TestArticle -count=1

# 运行特定测试用例
go test -v ./tests/api -run TestAuth_Register_Success -count=1
```

### 在 Docker 容器中运行

```bash
# 使用脚本运行所有测试
./scripts/run_api_tests.sh

# 运行详细输出
./scripts/run_api_tests.sh -v

# 只运行特定测试
./scripts/run_api_tests.sh -run TestAuth
```

## 测试数据

测试使用随机生成的数据进行隔离，确保测试之间不会相互影响：

- 用户名：`testuser_<timestamp>`
- 邮箱：`test_<timestamp>@test.com`
- 文章标题：`测试文章_<timestamp>`
- 栏目名称：`测试栏目_<timestamp>`
- 标签名称：`测试标签_<timestamp>`

## 辅助函数

### 认证相关

```go
// 创建测试用户并登录，返回 token
token := createTestUserAndLogin(t)

// 使用指定凭据登录
token := loginWithCredentials(t, username, password)

// 从响应中提取 token
token := extractTokenFromResponse(t, resp)
```

### 数据创建

```go
// 创建测试文章
articleID := createTestArticle(t, token, "published")

// 创建带标签的文章
articleID := createTestArticleWithTags(t, token, []string{"Go", "测试"})

// 创建测试评论
commentID := createTestComment(t, token, articleID, "评论内容")

// 创建测试栏目
categoryID := createTestCategory(t, token)
```

### 请求和断言

```go
// 发送请求
resp, statusCode, err := DoGet("/articles", "")
resp, statusCode, err := DoPost("/auth/login", body, "")
resp, statusCode, err := DoPut("/author/articles/1", body, token)
resp, statusCode, err := DoDelete("/author/articles/1", token)
resp, statusCode, err := DoPatch("/author/articles/1/status", body, token)

// 断言成功
AssertSuccess(t, resp, statusCode)

// 断言错误
AssertError(t, resp, statusCode, 1001)  // 1001 是 CodeInvalidParam

// 提取响应数据
data, err := ExtractData(resp)
```

## 状态码说明

| 状态码 | 含义 | 说明 |
|--------|------|------|
| 1000 | 成功 | 请求处理成功 |
| 1001 | 参数错误 | 请求参数无效或缺失 |
| 1002 | 用户已存在 | 注册时用户名已存在 |
| 1003 | 用户不存在 | 登录时用户不存在 |
| 1004 | 密码错误 | 登录密码错误 |
| 1005 | 服务繁忙 | 服务器内部错误 |
| 1006 | 需要登录 | 未提供有效的认证信息 |
| 1007 | Token无效 | Token 过期或无效 |
| 1008 | 栏目不存在 | 栏目 ID 不存在 |
| 1009 | 栏目已存在 | 栏目名称已存在 |
| 1014 | 文章不存在 | 文章 ID 不存在 |

## 注意事项

1. **测试独立性**：每个测试用例应该独立运行，不依赖其他测试的结果
2. **数据清理**：测试会自动清理创建的测试数据（用户、文章等）
3. **管理员权限**：部分管理员接口需要管理员账号，测试中使用了权限验证而非实际管理员操作
4. **跳过测试**：当前置条件不满足时（如无法获取 Token），测试会被跳过而非失败

## 添加新测试

```go
func TestNewFeature_Xxx(t *testing.T) {
    // 准备数据
    token := createTestUserAndLogin(t)
    if token == "" {
        t.Skip("无法获取Token，跳过测试")
        return
    }

    // 发送请求
    body := map[string]interface{}{
        "key": "value",
    }
    resp, statusCode, err := DoPost("/api/endpoint", body, token)
    if err != nil {
        t.Fatalf("请求失败: %v", err)
    }

    // 验证结果
    if !AssertSuccess(t, resp, statusCode) {
        t.Logf("响应: %+v", resp)
        return
    }

    // 验证响应数据
    data, err := ExtractData(resp)
    if err != nil {
        t.Fatalf("提取数据失败: %v", err)
    }

    if data["expected_field"] == nil {
        t.Error("响应缺少 expected_field 字段")
    }
}
```

## 故障排除

### 测试失败：服务无法访问

```bash
# 检查容器状态
docker ps | grep bluebell-ai

# 启动容器
docker start bluebell-ai

# 进入容器检查服务
docker exec -it bluebell-ai bash
curl http://localhost:8084/ping
```

### 测试失败：数据库连接错误

检查 `conf/dev.yaml` 中的数据库配置是否正确，确保：
- MySQL 服务可访问
- 数据库 `bluebell` 存在
- 用户名密码正确

### 测试失败：数据冲突

这是正常现象，测试会自动跳过。如果出现大量冲突，可能是：
- 测试频率过高，时间戳重复
- 数据库中残留测试数据

解决方案：
```bash
# 清理测试数据（慎用！）
# 在 MySQL 中执行：
# DELETE FROM user WHERE username LIKE 'testuser_%';
# DELETE FROM article WHERE title LIKE '测试文章_%';
```
