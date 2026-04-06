## 上下文

后端已有完整的 admin API（`/api/v1/admin/*`），均通过 `AdminAuthMiddleware` 做角色校验。
前端目前只有面向普通用户的布局（`DefaultLayout`），没有管理后台入口。
本次变更在前端新增独立的 `/admin` 路由树和布局，不改动现有任何页面。

## 目标 / 非目标

**目标：**
- 独立布局：管理后台与普通用户界面完全隔离，路由、样式互不影响
- 覆盖三个核心管理功能：用户管理、文章管理、系统统计
- 路由级权限守卫：非 admin 角色访问 `/admin/*` 立即跳转
- 复用现有的 axios 实例和 auth store，不引入新的状态管理方案

**非目标：**
- 不实现性能监控图表（`/admin/metrics/*`），留给后续版本
- 不实现管理员对文章内容的直接编辑（已有作者编辑功能）
- 不实现操作日志审计

## 决策

### 决策 1：独立布局而非嵌套在 DefaultLayout 里

**选择**：新建 `AdminLayout.vue`，在路由配置中使用 `component: AdminLayout` 作为父路由。

**理由**：管理后台通常有不同的导航结构（侧边栏 vs 顶部导航）、不同的视觉风格（更密集的表格布局）。独立布局避免样式污染，也方便未来独立部署。

**替代方案**：在 DefaultLayout 内增加条件渲染 → 代码耦合，不利于维护。

### 决策 2：路由守卫用 beforeEnter 而非全局 beforeEach

**选择**：在 `/admin` 路由组上加 `beforeEnter` 守卫，检查 token 解析出的 role 是否为 `admin`。

**理由**：只影响 admin 路由树，不干扰现有全局守卫逻辑，职责清晰。

### 决策 3：admin API 单独封装为 `src/api/admin.ts`

**选择**：新建 `admin.ts`，复用现有 axios 实例（request.ts），不新建实例。

**理由**：admin 接口同样需要 Authorization header，复用 request.ts 拦截器即可，无需重复实现 token 注入。

### 决策 4：用户/文章列表用前端分页（每次请求一页）

**选择**：每次翻页发起新请求，不做全量加载。

**理由**：管理员数据量大，全量加载内存压力高；后端已支持 page/size 参数。

## 风险 / 权衡

- **[风险] token 中 role 字段被篡改** → 后端 AdminAuthMiddleware 是最终防线，前端守卫只做 UX 层防护，可接受
- **[风险] admin 文章列表无删除功能** → 当前 API 无 `DELETE /admin/articles/:id`，本次不实现删除；精选设置和状态筛选已足够满足内容管理需求
- **[权衡] 不引入新 UI 组件库** → 复用已有的 Element Plus，避免引入额外依赖
