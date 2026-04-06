## 1. API 层

- [x] 1.1 创建 `src/api/admin.ts`，封装所有 admin 接口函数（用户列表、角色修改、状态修改、批量状态、文章列表、精选设置、统计概览、每日统计）

## 2. 路由与布局

- [x] 2.1 创建 `src/views/admin/AdminLayout.vue`（侧边栏 + 顶部栏 + 内容区 `<router-view>`）
- [x] 2.2 在 `src/router/index.ts` 新增 `/admin` 路由组

## 3. 系统统计页

- [x] 3.1 创建 `src/views/admin/StatsView.vue`，调用 `GET /admin/stats/overview` 展示 4 个数据卡片（总用户、总文章、总评论、总点赞）
- [x] 3.2 在统计页调用 `GET /admin/stats/daily`，用 ECharts 折线图展示最近 30 天用户/文章/评论每日新增趋势（三条折线）

## 4. 用户管理页

- [x] 4.1 创建 `src/views/admin/UsersView.vue`，展示用户列表表格（用户名、角色、状态、注册时间），支持角色/状态筛选和分页
- [x] 4.2 实现修改角色功能：行内下拉选择器 → 确认弹框 → 调用 `PATCH /admin/users/:id/role`，自身不可修改
- [x] 4.3 实现封禁/解封功能：操作列按钮 → 确认弹框 → 调用 `PATCH /admin/users/:id/status`，自身不可封禁
- [x] 4.4 实现批量封禁：表格多选 → 批量操作栏 → 调用 `PATCH /admin/users/batch/status`

## 5. 文章管理页

- [x] 5.1 创建 `src/views/admin/ArticlesView.vue`，展示文章列表表格（标题、作者用户名、状态标签、创建时间、是否精选），支持分页
- [x] 5.2 实现状态筛选、关键词搜索、作者名搜索（变更查询参数重新请求）
- [x] 5.3 实现设置/取消精选：操作列按钮 → 调用 `PATCH /admin/articles/:id/featured` → 更新行状态

## 6. 收尾

- [x] 6.1 在普通用户导航栏（admin 角色可见）添加"管理后台"入口链接，指向 `/admin`
- [x] 6.2 端到端验证：用 admin 账号登录，逐一测试三个页面的核心操作；用 author 账号访问 `/admin` 验证 403 跳转
