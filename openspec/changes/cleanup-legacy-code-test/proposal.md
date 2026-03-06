## 为什么

Bluebell 项目从旧版论坛系统（基于 Post/Community/Vote）迁移到新版 CMS 文章系统（基于 Article/Category/Tag/Like）。在项目演进过程中，遗留了大量未使用的代码文件，包括 Controller、Logic、DAO 和 Models 层。这些代码：

1. 未被 router 注册，无法通过任何 API 访问
2. 增加了代码维护复杂度，新开发者容易产生困惑
3. 拖慢了编译速度
4. 占用了代码审查和搜索的注意力

现在项目架构已经稳定，是时候清理这些遗留代码，保持代码库的整洁。

## 变更内容

### 删除的代码文件

**Controller 层（4个文件）:**
- `controller/post.go` - 旧版帖子功能 Handler
- `controller/community.go` - 旧版社区功能 Handler
- `controller/vote.go` - 投票功能 Handler
- `controller/post_test.go` - 帖子功能测试文件

**Logic 层（3个文件）:**
- `logic/post.go` - 帖子业务逻辑 (~6,157行，最大遗留文件)
- `logic/community.go` - 社区业务逻辑
- `logic/vote.go` - 投票业务逻辑

**DAO 层 - MySQL（3个文件）:**
- `dao/mysql/post.go` - 帖子数据访问
- `dao/mysql/community.go` - 社区数据访问
- `dao/mysql/post_test.go` - 测试文件

**DAO 层 - Redis（2个文件）:**
- `dao/redis/post.go` - 帖子 Redis 操作
- `dao/redis/vote.go` - 投票 Redis 操作 (ZSet)

**Models 层（2个文件）:**
- `models/post.go` - Post, ParamPostList, ParamVote 等模型
- `models/community.go` - Community 模型

**BREAKING:** 这些删除是破坏性变更，但这些代码未被任何活跃代码引用，实际影响为0。

### 验证步骤

1. 确认 router 中无引用
2. 确认活跃代码中无引用
3. 编译验证通过
4. 运行测试验证通过

## 功能 (Capabilities)

### 新增功能

本项目为代码清理类变更，不涉及新增功能或修改现有功能的行为。无需创建新的 capability 规范文件。

### 修改功能

无。本变更只涉及代码删除，不修改现有功能的行为规范。

## 影响

### 代码影响
- 删除约 14 个文件，总计 ~7,000+ 行代码
- 减少项目编译时间
- 降低新开发者的学习成本

### API 影响
- 无 API 删除（这些接口本就未注册到 router）
- 无 API 变更

### 数据库影响
- 无数据库变更建议（如需清理对应表，需单独评估）

### 依赖影响
- 无外部依赖变更

### 回滚方案
- 所有删除的文件可通过 Git 历史恢复
- 建议先在一个独立 commit 中删除，便于回滚
