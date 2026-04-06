<template>
  <div class="tags-container">
    <div class="container py-4">
      <h2 class="page-title">热门标签</h2>

      <div v-if="loading" class="loading-state">
        <p>加载中...</p>
      </div>

      <div v-else-if="allTags.length === 0" class="empty-state">
        <p>暂无标签</p>
      </div>

      <template v-else>
        <!-- 标签网格 -->
        <div class="tags-grid">
          <router-link
            v-for="tag in paginatedTags"
            :key="tag.id"
            :to="{ path: '/', query: { tag: tag.name } }"
            class="tag-card"
          >
            <div class="tag-name">{{ tag.name }}</div>
            <div class="tag-count">{{ tag.article_count || 0 }} 篇文章</div>
          </router-link>
        </div>

        <!-- 分页 -->
        <nav v-if="totalPages > 1" class="mt-4">
          <ul class="pagination justify-content-center">
            <li :class="{ disabled: currentPage === 1 }" class="page-item">
              <a href="javascript:void(0)" class="page-link" @click="currentPage > 1 && currentPage--">
                上一页
              </a>
            </li>
            <li v-for="page in displayPages" :key="page" :class="{ active: currentPage === page, disabled: page === '...' }" class="page-item">
              <a href="javascript:void(0)" class="page-link" @click="typeof page === 'number' && (currentPage = page)">
                {{ page }}
              </a>
            </li>
            <li :class="{ disabled: currentPage >= totalPages }" class="page-item">
              <a href="javascript:void(0)" class="page-link" @click="currentPage < totalPages && currentPage++">
                下一页
              </a>
            </li>
          </ul>
          <div class="page-info">
            共 {{ allTags.length }} 个标签，第 {{ currentPage }}/{{ totalPages }} 页
          </div>
        </nav>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import apiClient from '../services/api'

interface Tag {
  id: string
  name: string
  description?: string
  slug?: string
  article_count?: number
}

const allTags = ref<Tag[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = 36

// 计算总页数
const totalPages = computed(() => Math.ceil(allTags.value.length / pageSize) || 1)

// 计算当前页的标签
const paginatedTags = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return allTags.value.slice(start, start + pageSize)
})

// 计算显示的页码（带省略号）
const displayPages = computed(() => {
  const total = totalPages.value
  const current = currentPage.value
  const pages: (number | string)[] = []

  if (total <= 7) {
    for (let i = 1; i <= total; i++) pages.push(i)
  } else {
    pages.push(1)
    if (current > 3) pages.push('...')

    const start = Math.max(2, current - 1)
    const end = Math.min(total - 1, current + 1)
    for (let i = start; i <= end; i++) pages.push(i)

    if (current < total - 2) pages.push('...')
    pages.push(total)
  }

  return pages
})

async function loadTags() {
  loading.value = true
  try {
    const response = await apiClient.get('/tags')
    if (response.code === 1000) {
      // API 返回格式是 data: [...] 或 data: { list: [...] }
      allTags.value = Array.isArray(response.data) ? response.data : (response.data.list || [])
    }
  } catch (error) {
    console.error('加载标签失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadTags()
})
</script>

<style scoped>
.tags-container {
  background-color: #f8fafc;
  min-height: calc(100vh - 60px);
}

.page-title {
  font-weight: 700;
  margin-bottom: 2rem;
  color: #0f172a;
}

.tags-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 1rem;
}

.tag-card {
  text-align: center;
  text-decoration: none;
  padding: 1.25rem 1rem;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  transition: all 0.3s;
  background: white;
}

.tag-card:hover {
  border-color: var(--primary-color);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.tag-name {
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--primary-color);
  margin-bottom: 0.25rem;
  word-break: break-word;
}

.tag-count {
  color: #94a3b8;
  font-size: 0.85rem;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: 2rem;
  color: #64748b;
}

/* 分页 */
.pagination {
  display: flex;
  list-style: none;
  padding: 0;
  margin: 0;
  gap: 0.25rem;
}

.page-item .page-link {
  display: block;
  padding: 0.5rem 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 4px;
  color: var(--primary-color);
  text-decoration: none;
  transition: all 0.3s;
}

.page-item .page-link:hover {
  background-color: #f8fafc;
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

.page-info {
  text-align: center;
  margin-top: 1rem;
  color: #64748b;
  font-size: 0.9rem;
}

.justify-content-center {
  justify-content: center;
}

.mt-4 {
  margin-top: 1.5rem;
}
</style>
