## 上下文

### 当前状态

Bluebell 项目包含两个版本的代码：

1. **旧版论坛系统**（遗留代码）
   - 核心实体：`Post`（帖子）、`Community`（社区）、`Vote`（投票）
   - 存储：MySQL + Redis ZSet
   - 排序：基于投票分数的热榜算法
   - 状态：代码存在，但 router 未注册

2. **新版 CMS 文章系统**（活跃代码）
   - 核心实体：`Article`（文章）、`Category`（栏目）、`Tag`（标签）、`Like`（点赞）
   - 存储：MySQL 为主，Redis 仅少量缓存
   - 排序：基于时间或访问量
   - 状态：所有接口已注册，正常服务

### 遗留代码分布

```
14个文件，约7,000+行代码：

controller/     4个文件    ~250行
├── post.go              帖子Handler (CreatePost, GetPostDetail, GetPostList)
├── community.go         社区Handler
├── vote.go              投票Handler
└── post_test.go         测试文件

logic/          3个文件    ~6,300行
├── post.go              帖子业务逻辑（最大文件）
├── community.go         社区业务
└── vote.go              投票业务

dao/mysql/      3个文件    ~150行
├── post.go              帖子数据访问
├── community.go         社区数据访问
└── post_test.go         测试文件

dao/redis/      2个文件    ~120行
├── post.go              帖子Redis操作
└── vote.go              投票ZSet操作

models/         2个文件    ~60行
├── post.go              Post, ParamPostList等
└── community.go         Community模型
```

## 目标 / 非目标

**目标：**
- 删除所有未使用的遗留代码文件
- 确保项目编译通过
- 确保现有测试通过
- 建立验证流程，确认删除的安全性

**非目标：**
- 清理数据库表（本次只处理代码，数据保留）
- 重构活跃代码
- 添加新功能
- 修改现有 API 行为

## 决策

### 决策 1：分阶段删除策略

**选择：** 分两个阶段删除

**理由：**
- 阶段1：删除 Controller/Logic/DAO 层（低风险，这些只被自己引用）
- 阶段2：删除 Models 层（稍高风险，可能被其他模块引用到模型定义）

**替代方案考虑：**
- 一次性全部删除：风险略高，如果 models 被引用会导致编译失败
- 保留所有代码：不符合项目演进方向

### 决策 2：使用 Git 管理删除操作

**选择：** 所有删除操作放在一个独立的 Git commit 中

**理由：**
- 便于回滚（如果发现误删）
- 代码审查时可以清晰看到删除的文件列表
- Git 历史保留完整文件内容

**回滚命令：**
```bash
git revert <commit-hash>
```

### 决策 3：验证顺序

**选择：** 删除前按顺序验证

1. 静态检查：确认 router 中无引用
2. 编译检查：go build 通过
3. 测试检查：go test 通过

**理由：**
- 在删除前就发现问题，避免反复操作
- 编译是最可靠的验证方式（Go 是编译型语言）

## 风险 / 权衡

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|----------|
| 误删活跃代码 | 低 | 高 | 1. 双重验证引用检查<br>2. 编译验证<br>3. 独立 commit 便于回滚 |
| models 被其他模块引用 | 中 | 中 | 阶段2单独验证 models 的引用情况 |
| 删除后发现需要恢复功能 | 低 | 低 | 1. Git 历史完整保留<br>2. 功能需求可通过新架构实现 |
| 编译失败但难以定位 | 低 | 低 | 分阶段删除，每次验证 |

**权衡：**
- **短期成本**（清理时间） vs **长期收益**（维护成本降低）
- 建议接受短期成本，获得长期代码库健康

## 迁移计划

### 阶段 1：删除 Controller/Logic/DAO

```bash
# 1. 验证无引用
grep -r "CreatePostHandler\|GetPostDetailHandler\|CommunityHandler\|VotePostHandler" router/

# 2. 删除文件
rm controller/post.go controller/community.go controller/vote.go controller/post_test.go
rm logic/post.go logic/community.go logic/vote.go
rm dao/mysql/post.go dao/mysql/community.go dao/mysql/post_test.go
rm dao/redis/post.go dao/redis/vote.go

# 3. 编译验证
go build ./...

# 4. 测试验证
go test ./...

# 5. 提交
git add -A
git commit -m "cleanup: remove legacy post/community/vote controller, logic and dao files"
```

### 阶段 2：删除 Models（待阶段1验证后）

```bash
# 1. 验证 models 无引用
grep -r "type Post\|type Community\|type ParamPostList\|type ParamVote" --include="*.go" . | grep -v "models/post.go\|models/community.go"

# 2. 删除文件
rm models/post.go models/community.go

# 3. 编译验证
go build ./...

# 4. 测试验证
go test ./...

# 5. 提交
git add -A
git commit -m "cleanup: remove legacy post/community models"
```

## 验证清单

- [ ] Router 中无遗留 Handler 引用
- [ ] 活跃 Controller 无遗留 Logic 调用
- [ ] 活跃 Logic 无遗留 DAO 调用
- [ ] 无遗留 Models 类型引用
- [ ] `go build ./...` 编译通过
- [ ] `go test ./...` 测试通过

## Open Questions

无待定决策。本变更范围清晰，实施方案明确。
