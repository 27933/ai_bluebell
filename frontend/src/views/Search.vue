<template>
  <div class="search-container">
    <div class="container py-4">
      <h2 class="page-title">搜索文章</h2>

      <!-- 搜索框 -->
      <div class="search-box mb-4">
        <div class="input-group">
          <input
            v-model="keyword"
            type="text"
            class="form-control"
            placeholder="输入关键词搜索..."
            @keyup.enter="handleSearch"
            autofocus
          />
          <button class="btn btn-primary" type="button" @click="handleSearch">
            <i class="bi bi-search"></i> 搜索
          </button>
        </div>
      </div>

      <!-- 搜索结果统计 -->
      <div v-if="hasSearched" class="search-stats mb-3">
        <span v-if="articles.length > 0">
          找到 <strong>{{ total }}</strong> 条结果
        </span>
        <span v-else>
          未找到相关文章
        </span>
      </div>

      <!-- 搜索结果列表 -->
      <div v-if="loading" class="loading-state">
        <p>搜索中...</p>
      </div>

      <div v-else-if="!hasSearched" class="hint-state">
        <p>输入关键词开始搜索</p>
      </div>

      <div v-else-if="articles.length === 0" class="empty-state">
        <p>没有找到相关文章，试试其他关键词？</p>
      </div>

      <div v-else class="search-results">
        <router-link
          v-for="article in articles"
          :key="article.id"
          :to="`/article/${article.id}`"
          class="article-card"
        >
          <div class="card-body">
            <h5 class="article-title" v-html="highlightKeyword(article.title)"></h5>
            <p class="article-summary" v-html="highlightKeyword(article.summary)"></p>
            <div class="article-meta">
              <div class="author-info">
                <div class="author-avatar">{{ getInitial(article.author?.username) }}</div>
                <div>
                  <div class="author-name">{{ article.author?.nickname || article.author?.username || 'Unknown' }}</div>
                  <small class="text-muted">{{ timeAgo(article.created_at) }}</small>
                </div>
              </div>
              <div>
                <i class="bi bi-eye"></i> {{ article.view_count }}
              </div>
            </div>
          </div>
        </router-link>
      </div>

      <!-- 分页 -->
      <nav v-if="!loading && articles.length > 0 && totalPages > 1" class="mt-4">
        <ul class="pagination justify-content-center">
          <li :class="{ disabled: currentPage === 1 }" class="page-item">
            <a href="javascript:void(0)" class="page-link" @click="currentPage > 1 && (currentPage--, handleSearch())">
              上一页
            </a>
          </li>
          <li v-for="page in totalPages" :key="page" :class="{ active: currentPage === page }" class="page-item">
            <a href="javascript:void(0)" class="page-link" @click="currentPage = page; handleSearch()">
              {{ page }}
            </a>
          </li>
          <li :class="{ disabled: currentPage >= totalPages }" class="page-item">
            <a href="javascript:void(0)" class="page-link" @click="currentPage < totalPages && (currentPage++, handleSearch())">
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
import { useRoute } from 'vue-router'
import apiClient from '../services/api'

interface Author {
  id: string
  username: string
  nickname?: string
}

interface Article {
  id: string
  title: string
  summary: string
  author?: Author
  view_count: number
  created_at: string
}

const route = useRoute()
const keyword = ref('')
const articles = ref<Article[]>([])
const loading = ref(false)
const hasSearched = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const totalPages = computed(() => Math.ceil(total.value / pageSize.value) || 1)

function getInitial(name?: string): string {
  if (!name) return '?'
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

function highlightKeyword(text: string): string {
  if (!keyword.value || !text) return text
  const regex = new RegExp(`(${keyword.value})`, 'gi')
  return text.replace(regex, '<mark>$1</mark>')
}

async function handleSearch() {
  if (!keyword.value.trim()) return

  loading.value = true
  hasSearched.value = true

  try {
    const response = await apiClient.get('/articles/search', {
      params: {
        keyword: keyword.value,
        page: currentPage.value,
        size: pageSize.value,
      },
    })

    if (response.code === 1000) {
      articles.value = response.data.list || []
      total.value = response.data.total || 0
    }
  } catch (error) {
    console.error('搜索失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  // 如果 URL 带有 keyword 参数，自动搜索
  const queryKeyword = route.query.keyword as string
  if (queryKeyword) {
    keyword.value = queryKeyword
    handleSearch()
  }
})
</script>

<style scoped>
.search-container {
  background-color: #f8fafc;
  min-height: calc(100vh - 60px);
}

.page-title {
  font-weight: 700;
  margin-bottom: 1.5rem;
  color: #0f172a;
}

.search-box .input-group {
  max-width: 600px;
}

.form-control {
  border: 1px solid #e2e8f0;
  border-radius: 6px 0 0 6px;
  padding: 0.75rem 1rem;
  font-size: 1rem;
}

.form-control:focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.btn-primary {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
  border-radius: 0 6px 6px 0;
  padding: 0.75rem 1.25rem;
}

.search-stats {
  color: #64748b;
  font-size: 0.95rem;
}

.search-results {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.article-card {
  background: white;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  text-decoration: none;
  color: inherit;
  transition: all 0.3s;
}

.article-card:hover {
  border-color: var(--primary-color);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.card-body {
  padding: 1.5rem;
}

.article-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #0f172a;
  margin-bottom: 0.5rem;
}

.article-summary {
  color: #64748b;
  font-size: 0.95rem;
  line-height: 1.6;
  margin-bottom: 1rem;
}

.article-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.85rem;
  color: #94a3b8;
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

.loading-state,
.empty-state,
.hint-state {
  text-align: center;
  padding: 3rem;
  color: #64748b;
}

.pagination {
  display: flex;
  list-style: none;
  padding: 0;
  gap: 0.25rem;
}

.page-item .page-link {
  padding: 0.5rem 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  color: var(--primary-color);
  text-decoration: none;
}

.page-item.active .page-link {
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.page-item.disabled .page-link {
  color: #cbd5e1;
  pointer-events: none;
}

:deep(mark) {
  background-color: #fef08a;
  padding: 0 2px;
}

@media (max-width: 768px) {
  .search-box .input-group {
    max-width: 100%;
    flex-direction: row;
  }

  .search-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }

  .pagination {
    flex-wrap: wrap;
    gap: 0.25rem;
  }

  .result-item {
    padding: 1rem;
  }
}
</style>
