## 上下文

前端项目使用 Vue 3 + TypeScript，通过 `apiClient`（基于 axios）与后端 API 通信。API 文档定义在 `backend/docs/API_INTEGRATION.md`，基础 URL 为 `/api/v1`。

当前问题：
1. 部分页面的 API 调用参数与后端接口不匹配
2. 部分功能使用硬编码数据或 TODO 占位符
3. 响应数据的字段映射与实际 API 返回结构不一致

## 目标 / 非目标

**目标：**
- 修复所有前端页面的 API 调用，使其与后端接口完全匹配
- 确保数据结构正确映射，页面能正常显示数据
- 保持代码风格与现有代码一致

**非目标：**
- 不添加新功能或新页面
- 不修改 API 层的 `apiClient` 配置
- 不调整 UI 样式或布局
- 不处理后端 API 的任何修改

## 决策

### 1. Home.vue 搜索功能

**决策**：直接调用 `GET /articles/search` 接口

**实现方式**：
```typescript
async function handleSearch() {
  if (!searchQuery.value.trim()) return
  
  const response = await apiClient.get('/articles/search', {
    params: {
      keyword: searchQuery.value,
      page: currentPage.value,
      size: pageSize.value,
    },
  })
  // 复用现有的文章列表渲染逻辑
}
```

**替代方案考虑**：
- 前端过滤：不可行，无法搜索未加载的数据
- 使用通用 `/articles` 接口的 `keyword` 参数：可行但语义不如专用搜索接口清晰

### 2. Home.vue 标签获取

**决策**：在 `onMounted` 时调用 `GET /tags` 获取标签列表

**实现方式**：
```typescript
const popularTags = ref<{id: string, name: string}[]>([])

async function loadTags() {
  const response = await apiClient.get('/tags')
  if (response.code === 1000) {
    popularTags.value = response.data.list || []
  }
}
```

### 3. Home.vue 排序参数映射

**决策**：前端保持用户友好的选项名称，发送请求时转换为 API 参数

**映射关系**：
| 前端值 | API sort 参数 |
|--------|---------------|
| `hot` | `popular` |
| `latest` | `time` |
| `comments` | `time`（API 不支持按评论排序，降级为时间排序）|

### 4. ArticleDetail.vue 数据映射

**决策**：修改模板和脚本中的字段访问路径

**字段映射**：
- 作者名：`article.author` → `article.author?.nickname || article.author?.username`
- 作者头像初始：`article.author_id` → `article.author?.username`
- 评论用户名：`comment.user_name` → `comment.author?.nickname || comment.author?.username`
- 评论用户 ID：`comment.user_id` → `comment.author?.id`

### 5. Dashboard.vue 趋势图

**决策**：由于 API 需要 `article_id` 参数，改为展示用户最热门文章的趋势

**实现方式**：
1. 先获取用户文章列表，找出浏览量最高的文章
2. 使用该文章的 ID 调用趋势 API
3. 如果没有文章，显示空状态

```typescript
async function loadTrendData() {
  // 找出浏览量最高的文章
  const topArticle = articles.value
    .filter(a => a.status === 'published')
    .sort((a, b) => (b.view_count || 0) - (a.view_count || 0))[0]
  
  if (!topArticle) return // 显示空状态
  
  const response = await apiClient.get('/article-stats/trend', {
    params: {
      article_id: topArticle.id,
      days: timeRange.value === 'week' ? 7 : 30,
      group_by: 'day',
    },
  })
}
```

**替代方案考虑**：
- 显示所有文章的汇总趋势：API 不支持
- 显示多篇文章的趋势：可行但 UI 复杂度增加

### 6. Profile.vue 用户数据

**决策**：在 `onMounted` 时调用 `GET /auth/profile` 获取最新数据

**实现方式**：
```typescript
async function loadProfile() {
  const response = await apiClient.get('/auth/profile')
  if (response.code === 1000) {
    const data = response.data
    form.email = data.email || ''
    form.nickname = data.nickname || ''
    form.bio = data.bio || ''
    // 同时更新 authStore
    authStore.setUser(data)
  }
}
```

## 风险 / 权衡

| 风险 | 缓解措施 |
|------|----------|
| Dashboard 趋势图只显示单篇文章数据 | 在 UI 上明确标注"热门文章趋势"，用户可理解 |
| 按评论排序功能降级 | 可在 UI 上移除此选项，或保留但提示用户 |
| API 调用失败时的用户体验 | 复用现有的错误处理逻辑，保持一致 |
| 并发请求过多影响性能 | Home.vue 的文章和标签可并行加载；Dashboard 需串行（先文章后趋势）|
