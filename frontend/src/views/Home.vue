<template>
  <div class="page-wrapper">
    <div class="container py-4">
      <!-- 搜索和筛选 -->
      <div class="row mb-4">
        <div class="col-md-8">
          <div class="input-group">
            <input
              v-model="searchQuery"
              type="text"
              class="form-control"
              placeholder="搜索文章..."
            />
            <button class="btn btn-primary" type="button" @click="handleSearch">
              <i class="bi bi-search"></i> 搜索
            </button>
          </div>
        </div>
        <div class="col-md-4">
          <select v-model="sortBy" class="form-select">
            <option value="hot">按热度排序</option>
            <option value="latest">按时间排序</option>
            <option value="comments">按评论排序</option>
          </select>
        </div>
      </div>

      <!-- 标签筛选 -->
      <div class="mb-4">
        <div class="d-flex gap-2 flex-wrap">
          <span
            v-for="tag in popularTags"
            :key="tag"
            class="tag"
            @click="searchQuery = tag"
          >
            {{ tag }}
          </span>
        </div>
      </div>

      <!-- 文章列表 -->
      <div v-if="loading" class="empty-state">
        <p>加载中...</p>
      </div>

      <div v-else-if="articles.length === 0" class="empty-state">
        <p>暂无文章</p>
      </div>

      <div v-else class="row g-4 mb-4">
        <!-- 文章卡片 -->
        <div
          v-for="article in articles"
          :key="article.id"
          class="col-md-6"
        >
          <router-link :to="`/article/${article.id}`" class="article-card">
            <div class="card-body">
              <h5 class="article-title">{{ article.title }}</h5>
              <p class="article-summary">{{ article.summary }}</p>
              <div class="mb-2">
                <span
                  v-for="tag in (article.tags || []).slice(0, 2)"
                  :key="tag"
                  class="tag"
                >
                  {{ tag }}
                </span>
              </div>
              <div class="article-meta">
                <div class="author-info">
                  <div class="author-avatar">{{ getInitial(article.author?.username) }}</div>
                  <div>
                    <div class="author-name">{{ article.author?.username || 'Unknown' }}</div>
                    <small class="text-muted">{{ timeAgo(article.created_at) }}</small>
                  </div>
                </div>
                <div>
                  <i class="bi bi-eye"></i> {{ article.view_count }} &nbsp;
                  <i class="bi bi-heart"></i> {{ article.like_count }} &nbsp;
                  <i class="bi bi-chat"></i> {{ article.comment_count }}
                </div>
              </div>
            </div>
          </router-link>
        </div>
      </div>

      <!-- 分页 -->
      <nav v-if="!loading && articles.length > 0" aria-label="Page navigation">
        <ul class="pagination justify-content-center">
          <li :class="{ disabled: currentPage === 1 }" class="page-item">
            <a
              href="javascript:void(0)"
              class="page-link"
              @click="currentPage > 1 && (currentPage--, loadArticles())"
            >
              上一页
            </a>
          </li>
          <li
            v-for="page in totalPages"
            :key="page"
            :class="{ active: currentPage === page }"
            class="page-item"
          >
            <a
              href="javascript:void(0)"
              class="page-link"
              @click="currentPage = page; loadArticles()"
            >
              {{ page }}
            </a>
          </li>
          <li :class="{ disabled: currentPage >= totalPages }" class="page-item">
            <a
              href="javascript:void(0)"
              class="page-link"
              @click="currentPage < totalPages && (currentPage++, loadArticles())"
            >
              下一页
            </a>
          </li>
        </ul>
      </nav>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import apiClient from '../services/api'

interface Author {
  id: string
  username: string
}

interface Article {
  id: string
  title: string
  summary: string
  author: Author
  tags?: string[]
  view_count: number
  like_count: number
  comment_count: number
  created_at: string
}

const articles = ref<Article[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(6)
const total = ref(0)
const searchQuery = ref('')
const sortBy = ref('hot')

const popularTags = [
  'Vue.js',
  'TypeScript',
  '后端开发',
  '数据库',
  '性能优化',
]

const totalPages = computed(() => Math.ceil(total.value / pageSize.value) || 1)

function getInitial(name?: string): string {
  if (!name) return '？'
  return name.charAt(0).toUpperCase()
}

function timeAgo(dateString: string): string {
  const date = new Date(dateString)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) return '今天'
  if (days === 1) return '1 天前'
  if (days < 7) return `${days} 天前`
  return `${Math.floor(days / 7)} 周前`
}

async function loadArticles() {
  loading.value = true
  try {
    const response = await apiClient.get('/articles', {
      params: {
        page: currentPage.value,
        size: pageSize.value,
        status: 'published',
      },
    })

    if (response.code === 1000) {
      articles.value = response.data.list || []
      total.value = response.data.total || 0
    } else {
      console.error('加载失败:', response.msg)
    }
  } catch (error: any) {
    console.error('加载失败:', error.message)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  if (searchQuery.value.trim()) {
    // 暂时在当前列表中过滤搜索结果
    // TODO: 调用后端搜索 API
    console.log('搜索关键词:', searchQuery.value)
  }
}

onMounted(() => {
  loadArticles()
})
</script>

<style scoped>
:root,
:host {
  --primary-color: #2563eb;
  --secondary-color: #64748b;
  --danger-color: #ef4444;
}

.page-wrapper {
  background-color: #f8fafc;
  min-height: 100%;
  padding: 0;
}

.container {
  max-width: 1140px;
  margin: 0 auto;
  padding: 0 1rem;
  width: 100%;
}

.py-4 {
  padding-top: 2rem;
  padding-bottom: 2rem;
}

.mb-4 {
  margin-bottom: 1.5rem;
}

.mb-2 {
  margin-bottom: 0.5rem;
}

/* Bootstrap utility classes (override if needed) */
/* Note: Most of these come from Bootstrap CDN now */

/* 表单 */
.input-group {
  display: flex;
  gap: 0.5rem;
  align-items: stretch;
}

.form-control,
.form-select {
  border: 1px solid #e2e8f0 !important;
  border-radius: 6px !important;
  padding: 0.75rem 1rem !important;
  font-size: 1rem !important;
  transition: all 0.3s !important;
}

.form-control {
  flex: 1 !important;
  color: #334155 !important;
  background-color: white !important;
  min-width: 0;
}

.form-control:focus,
.form-select:focus {
  outline: none !important;
  border-color: var(--primary-color) !important;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1) !important;
}

.form-select {
  width: 100% !important;
  color: #334155 !important;
  background-color: white !important;
  cursor: pointer !important;
}

/* 按钮 */
.btn {
  padding: 0.75rem 1.25rem !important;
  border: none !important;
  border-radius: 6px !important;
  cursor: pointer !important;
  font-weight: 500 !important;
  transition: all 0.3s !important;
  display: inline-flex !important;
  align-items: center !important;
  gap: 0.5rem !important;
  font-size: 1rem !important;
  white-space: nowrap !important;
  flex-shrink: 0 !important;
  min-width: max-content !important;
}

.btn-primary {
  background-color: #2563eb !important;
  color: white !important;
  border: none !important;
}

.btn-primary:hover {
  background-color: #1d4ed8 !important;
  transform: translateY(-1px) !important;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.3) !important;
}

/* 标签 */
.tag {
  display: inline-block;
  background-color: #eef2ff;
  color: var(--primary-color);
  padding: 0.25rem 0.75rem;
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: 500;
  margin-right: 0.5rem;
  margin-bottom: 0.5rem;
  cursor: pointer;
  transition: all 0.3s;
}

.tag:hover {
  background-color: var(--primary-color);
  color: white;
}

.d-flex {
  display: flex;
}

.gap-2 {
  gap: 0.5rem;
}

.flex-wrap {
  flex-wrap: wrap;
}

/* 文章卡片 */
.article-card {
  border: none;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
  display: block;
  text-decoration: none;
  color: inherit;
  height: 100%;
}

.article-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
}

.card-body {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.article-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #0f172a;
  margin-bottom: 0.5rem;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.2;
}

.article-summary {
  color: #64748b;
  font-size: 0.95rem;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 1rem;
  flex-grow: 1;
}

.article-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.85rem;
  color: #94a3b8;
  border-top: 1px solid #e2e8f0;
  padding-top: 1rem;
  margin-top: auto;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.author-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary-color), #7c3aed);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 0.9rem;
}

.author-name {
  font-weight: 600;
  color: #0f172a;
}

.text-muted {
  color: #94a3b8;
  font-size: 0.9rem;
}

/* 分页 */
.pagination {
  display: flex;
  list-style: none;
  padding: 0;
  margin: 0;
  gap: 0.25rem;
}

.page-item {
  list-style: none;
}

.page-link {
  display: block;
  padding: 0.5rem 0.75rem;
  color: var(--primary-color);
  text-decoration: none;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  transition: all 0.3s;
}

.page-link:hover {
  background-color: #f8fafc;
}

.page-item.active .page-link {
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.page-item.disabled .page-link {
  color: #cbd5e1;
  cursor: not-allowed;
  pointer-events: none;
}

.justify-content-center {
  justify-content: center;
}

.empty-state {
  text-align: center;
  padding: 2rem 1rem;
  color: #64748b;
}

@media (max-width: 768px) {
  .col-md-8,
  .col-md-6,
  .col-md-4 {
    flex: 0 0 100%;
  }

  .input-group {
    flex-direction: column;
  }

  .pagination {
    flex-wrap: wrap;
  }

  .article-meta {
    flex-direction: column;
    align-items: flex-start;
  }

  .col-md-6 {
    flex: 0 0 100%;
  }
}
</style>
