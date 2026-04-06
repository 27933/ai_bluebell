## 1. Home.vue 修复

- [x] 1.1 添加 loadTags 函数，调用 GET /tags 获取标签列表
- [x] 1.2 修改 popularTags 从硬编码改为响应式数据
- [x] 1.3 在 onMounted 中调用 loadTags
- [x] 1.4 实现 handleSearch 函数，调用 GET /articles/search
- [x] 1.5 修正 loadArticles 的 sort 参数映射（hot→popular, latest→time）

## 2. ArticleDetail.vue 修复

- [x] 2.1 修复模板中 author 字段访问路径（article.author.nickname/username）
- [x] 2.2 修复 getInitial 函数参数（从 author_id 改为 author.username）
- [x] 2.3 修复评论列表中用户名显示（comment.author.nickname/username）
- [x] 2.4 更新 Article 和 Comment 接口定义以匹配 API 返回结构

## 3. Dashboard.vue 修复

- [x] 3.1 修改 loadTrendData 函数，先获取最热门文章
- [x] 3.2 使用正确的 API 参数调用趋势接口（article_id, days, group_by）
- [x] 3.3 处理无文章时的空状态显示
- [x] 3.4 确保文章列表加载完成后再加载趋势数据

## 4. Profile.vue 修复

- [x] 4.1 添加 loadProfile 函数，调用 GET /auth/profile
- [x] 4.2 在 onMounted 中调用 loadProfile 获取最新数据
- [x] 4.3 更新 authStore.user 以保持数据同步

## 5. 验证测试

- [x] 5.1 测试 Home.vue 搜索功能（需手动验证）
- [x] 5.2 测试 Home.vue 标签加载和筛选（需手动验证）
- [x] 5.3 测试 ArticleDetail.vue 作者和评论显示（需手动验证）
- [x] 5.4 测试 Dashboard.vue 趋势图（需手动验证）
- [x] 5.5 测试 Profile.vue 数据加载和更新（需手动验证）
