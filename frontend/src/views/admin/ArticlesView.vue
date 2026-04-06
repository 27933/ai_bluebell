<template>
  <div class="articles-page">
    <!-- 筛选/搜索栏 -->
    <div class="filter-bar">
      <el-select v-model="filterStatus" placeholder="状态" style="width: 120px" @change="handleFilter">
        <el-option label="全部状态" value="all" />
        <el-option label="已发布" value="published" />
        <el-option label="草稿" value="draft" />
        <el-option label="已下线" value="offline" />
      </el-select>
      <el-input
        v-model="keyword"
        placeholder="搜索标题/内容"
        style="width: 200px"
        clearable
        @keyup.enter="handleFilter"
        @clear="handleFilter"
      >
        <template #prefix><i class="bi bi-search"></i></template>
      </el-input>
      <el-input
        v-model="authorName"
        placeholder="按作者名搜索"
        style="width: 160px"
        clearable
        @keyup.enter="handleFilter"
        @clear="handleFilter"
      >
        <template #prefix><i class="bi bi-person"></i></template>
      </el-input>
      <el-button type="primary" @click="handleFilter">搜索</el-button>
    </div>

    <!-- 表格 -->
    <div class="table-card">
      <el-table :data="articles" v-loading="loading">
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
        <el-table-column prop="author_username" label="作者" width="120" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" size="small">
              {{ statusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="精选" width="70">
          <template #default="{ row }">
            <i
              v-if="row.is_featured"
              class="bi bi-star-fill"
              style="color: #f59e0b; font-size: 1rem;"
            ></i>
            <i v-else class="bi bi-star" style="color: #cbd5e1; font-size: 1rem;"></i>
          </template>
        </el-table-column>
        <el-table-column label="浏览/点赞" width="110">
          <template #default="{ row }">
            <span style="color:#64748b; font-size:0.85rem;">
              {{ row.view_count }} / {{ row.like_count }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="130">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="!row.is_featured"
              size="small"
              type="warning"
              text
              @click="toggleFeatured(row, true)"
            >设为精选</el-button>
            <el-button
              v-else
              size="small"
              type="info"
              text
              @click="toggleFeatured(row, false)"
            >取消精选</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="loadArticles"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getAdminArticleList, setArticleFeatured } from '../../api/admin'
import type { AdminArticle } from '../../api/admin'

const articles = ref<AdminArticle[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(false)
const filterStatus = ref('all')
const keyword = ref('')
const authorName = ref('')

onMounted(() => loadArticles())

async function loadArticles() {
  loading.value = true
  try {
    const res = await getAdminArticleList({
      status: filterStatus.value,
      keyword: keyword.value || undefined,
      author_name: authorName.value || undefined,
      page: page.value,
      size: pageSize,
    })
    if (res.code === 1000) {
      articles.value = res.data.list || []
      total.value = res.data.total
    }
  } catch {
    ElMessage.error('加载文章列表失败')
  } finally {
    loading.value = false
  }
}

function handleFilter() {
  page.value = 1
  loadArticles()
}

async function toggleFeatured(row: AdminArticle, featured: boolean) {
  try {
    const res = await setArticleFeatured(row.id, featured)
    if (res.code === 1000) {
      row.is_featured = featured
      ElMessage.success(featured ? '已设为精选' : '已取消精选')
    }
  } catch {
    ElMessage.error('操作失败')
  }
}

function statusLabel(status: string): string {
  const map: Record<string, string> = {
    published: '已发布',
    draft: '草稿',
    offline: '已下线',
  }
  return map[status] || status
}

function statusTagType(status: string): 'success' | 'info' | 'danger' | 'warning' {
  const map: Record<string, 'success' | 'info' | 'danger' | 'warning'> = {
    published: 'success',
    draft: 'info',
    offline: 'danger',
  }
  return map[status] || 'info'
}

function formatDate(str: string): string {
  return str ? new Date(str).toLocaleDateString('zh-CN') : '-'
}
</script>

<style scoped>
.articles-page {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.table-card {
  background: white;
  border-radius: 10px;
  padding: 1rem;
  box-shadow: 0 1px 3px rgba(0,0,0,0.07);
}

.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 1rem;
}
</style>
