<template>
  <div class="dashboard-container">
    <!-- 欢迎信息 -->
    <div class="welcome-card">
      <h2>欢迎回来，{{ authStore.user?.nickname || authStore.user?.username }}！</h2>
      <p>继续创作，分享你的想法和经验。</p>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-number">{{ stats.published }}</div>
        <div class="stat-label">已发布文章</div>
      </div>

      <div class="stat-card">
        <div class="stat-number">{{ stats.totalViews }}</div>
        <div class="stat-label">总浏览数</div>
      </div>

      <div class="stat-card">
        <div class="stat-number">{{ stats.totalLikes }}</div>
        <div class="stat-label">获赞总数</div>
      </div>

      <div class="stat-card">
        <div class="stat-number">{{ stats.totalComments }}</div>
        <div class="stat-label">总评论数</div>
      </div>
    </div>

    <!-- 阅读趋势图 -->
    <div class="trend-card">
      <div class="trend-header">
        <h3>阅读趋势</h3>
        <div class="btn-group">
          <button
            :class="{ active: timeRange === 'week' }"
            @click="timeRange = 'week'"
            class="btn-sm"
          >
            周统计
          </button>
          <button
            :class="{ active: timeRange === 'month' }"
            @click="timeRange = 'month'"
            class="btn-sm"
          >
            月统计
          </button>
        </div>
      </div>
      <div id="trend-chart" style="height: 300px; width: 100%;" />
    </div>

    <!-- 文章管理 -->
    <div class="articles-card">
      <div class="articles-header">
        <h3>我的文章</h3>
        <router-link to="/write" class="btn-primary">
          <i class="bi bi-plus"></i> 写新文章
        </router-link>
      </div>

      <!-- 筛选按钮 -->
      <div class="filter-buttons">
        <button
          :class="{ active: filterStatus === 'all' }"
          @click="filterStatus = 'all'"
          class="btn-filter"
        >
          全部
        </button>
        <button
          :class="{ active: filterStatus === 'published' }"
          @click="filterStatus = 'published'"
          class="btn-filter"
        >
          已发布
        </button>
        <button
          :class="{ active: filterStatus === 'draft' }"
          @click="filterStatus = 'draft'"
          class="btn-filter"
        >
          草稿
        </button>
      </div>

      <!-- 文章表格 -->
      <div class="articles-table-wrapper">
        <table class="articles-table" v-if="filteredArticles.length > 0">
          <thead>
            <tr>
              <th>标题</th>
              <th>状态</th>
              <th>浏览</th>
              <th>赞</th>
              <th>评论</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="article in filteredArticles" :key="article.id">
              <td><strong>{{ article.title }}</strong></td>
              <td>
                <span class="badge" :class="article.status === 'published' ? 'badge-success' : 'badge-warning'">
                  {{ article.status === 'published' ? '已发布' : '草稿' }}
                </span>
              </td>
              <td>{{ article.view_count || 0 }}</td>
              <td>{{ article.like_count || 0 }}</td>
              <td>{{ article.comment_count || 0 }}</td>
              <td>
                <div class="action-buttons">
                  <router-link :to="`/article/${article.id}`">
                    <button class="btn-sm btn-outline">查看</button>
                  </router-link>
                  <router-link :to="`/write/${article.id}`">
                    <button class="btn-sm btn-outline">编辑</button>
                  </router-link>
                  <button class="btn-sm btn-outline btn-danger" @click="handleDelete(article.id)">删除</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-else class="empty-state">
          <p>暂无文章</p>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination" v-if="total > pageSize">
        <button
          :disabled="currentPage === 1"
          @click="currentPage--; loadArticles()"
          class="btn-sm"
        >
          上一页
        </button>
        <span class="page-info">第 {{ currentPage }} 页，共 {{ Math.ceil(total / pageSize) }} 页</span>
        <button
          :disabled="currentPage >= Math.ceil(total / pageSize)"
          @click="currentPage++; loadArticles()"
          class="btn-sm"
        >
          下一页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import apiClient from '../services/api'
import { useAuthStore } from '../stores/auth'
import * as echarts from 'echarts'

const authStore = useAuthStore()

// 统计数据
const stats = reactive({
  published: 0,
  totalViews: 0,
  totalLikes: 0,
  totalComments: 0,
})

// 趋势图数据
const timeRange = ref<'week' | 'month'>('week')

// 文章列表
const articles = ref<any[]>([])
const filterStatus = ref<'all' | 'published' | 'draft'>('all')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const filteredArticles = computed(() => {
  if (filterStatus.value === 'all') {
    return articles.value
  }
  return articles.value.filter(article => article.status === filterStatus.value)
})

async function loadArticles() {
  try {
    const response = await apiClient.get('/author/articles', {
      params: {
        page: currentPage.value,
        size: pageSize.value,
      },
    })

    if (response.code === 1000) {
      articles.value = response.data.list || []
      total.value = response.data.total || 0

      // 计算统计数据
      stats.published = articles.value.filter(a => a.status === 'published').length
      stats.totalViews = articles.value.reduce((sum, a) => sum + (a.view_count || 0), 0)
      stats.totalLikes = articles.value.reduce((sum, a) => sum + (a.like_count || 0), 0)
      stats.totalComments = articles.value.reduce((sum, a) => sum + (a.comment_count || 0), 0)
    } else {
      console.error('加载文章失败:', response.msg)
    }
  } catch (error: any) {
    console.error('加载文章出错:', error)
  }
}

async function handleDelete(articleId: string) {
  if (!confirm('确定要删除这篇文章吗？')) {
    return
  }

  try {
    const response = await apiClient.delete(`/author/articles/${articleId}`)
    if (response.code === 1000) {
      alert('文章已删除')
      loadArticles()
    } else {
      alert(response.msg || '删除失败')
    }
  } catch (error: any) {
    alert('删除失败：' + (error.message || '未知错误'))
  }
}

async function loadTrendData() {
  try {
    const response = await apiClient.get('/article-stats/trend', {
      params: {
        time_range: timeRange.value,
        group_by: timeRange.value === 'week' ? 'hour' : 'day',
      },
    })
    if (response.code === 1000) {
      renderChart(response.data)
    } else {
      console.error('加载趋势数据失败:', response.msg)
    }
  } catch (error: any) {
    console.error('加载趋势数据出错:', error)
  }
}

function renderChart(data: any) {
  const chartDOM = document.getElementById('trend-chart') as HTMLElement
  if (!chartDOM) return

  const myChart = echarts.init(chartDOM)

  const labels = data.labels || data.dates || []
  const values = data.views || data.counts || []

  const option = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(0, 0, 0, 0.7)',
      borderColor: '#ccc',
      textStyle: {
        color: '#fff',
      },
    },
    grid: {
      left: '3%',
      right: '3%',
      bottom: '3%',
      top: '5%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      data: labels,
      boundaryGap: false,
      axisLine: {
        lineStyle: {
          color: '#e2e8f0',
        },
      },
      axisLabel: {
        color: '#94a3b8',
        fontSize: 12,
      },
    },
    yAxis: {
      type: 'value',
      axisLine: {
        lineStyle: {
          color: '#e2e8f0',
        },
      },
      axisLabel: {
        color: '#94a3b8',
        fontSize: 12,
      },
      splitLine: {
        lineStyle: {
          color: '#f1f5f9',
        },
      },
    },
    series: [
      {
        name: '浏览数',
        type: 'line',
        data: values,
        smooth: true,
        lineStyle: {
          color: '#2563eb',
          width: 2,
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(37, 99, 235, 0.3)' },
            { offset: 1, color: 'rgba(37, 99, 235, 0)' },
          ]),
        },
        itemStyle: {
          color: '#2563eb',
          borderWidth: 2,
          borderColor: '#fff',
        },
      },
    ],
  }
  myChart.setOption(option)

  window.addEventListener('resize', () => {
    myChart.resize()
  })
}

onMounted(() => {
  loadArticles()
  loadTrendData()
})

watch(() => timeRange.value, () => {
  loadTrendData()
})
</script>

<style scoped>
:root {
  --primary-color: #2563eb;
  --secondary-color: #64748b;
  --danger-color: #ef4444;
  --bg-color: #f8fafc;
  --card-bg: #ffffff;
  --border-color: #e2e8f0;
  --text-primary: #0f172a;
  --text-secondary: #64748b;
  --text-tertiary: #94a3b8;
}

.dashboard-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1.25rem;
  background-color: var(--bg-color);
}

/* 欢迎卡片 */
.welcome-card {
  background: linear-gradient(135deg, var(--primary-color), #7c3aed);
  border-radius: 12px;
  padding: 2rem;
  color: white;
  margin-bottom: 2rem;
}

.welcome-card h2 {
  margin: 0 0 0.5rem 0;
  font-size: 1.5rem;
  font-weight: 700;
}

.welcome-card p {
  margin: 0;
  font-size: 1rem;
  opacity: 0.95;
}

/* 统计卡片 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: var(--card-bg);
  border-radius: 8px;
  padding: 1.5rem;
  border: 1px solid var(--border-color);
  text-align: center;
}

.stat-number {
  font-size: 2rem;
  font-weight: 700;
  color: var(--primary-color);
  margin-bottom: 0.5rem;
}

.stat-label {
  font-size: 0.9rem;
  color: var(--text-secondary);
}

/* 趋势图卡片 */
.trend-card {
  background: var(--card-bg);
  border-radius: 12px;
  padding: 2rem;
  border: 1px solid var(--border-color);
  margin-bottom: 2rem;
}

.trend-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.trend-header h3 {
  margin: 0;
  font-weight: 700;
  color: var(--text-primary);
  font-size: 1.1rem;
}

.btn-group {
  display: flex;
  gap: 0.5rem;
}

.btn-group .btn-sm {
  padding: 0.4rem 0.8rem;
  border: 1px solid var(--border-color);
  background: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
  color: var(--text-secondary);
  transition: all 0.3s;
}

.btn-group .btn-sm:hover {
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.btn-group .btn-sm.active {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

/* 文章管理卡片 */
.articles-card {
  background: var(--card-bg);
  border-radius: 12px;
  padding: 2rem;
  border: 1px solid var(--border-color);
}

.articles-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.articles-header h3 {
  margin: 0;
  font-weight: 700;
  color: var(--text-primary);
  font-size: 1.1rem;
}

.btn-primary {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.6rem 1rem;
  background: var(--primary-color);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.3s;
  font-size: 0.95rem;
}

.btn-primary:hover {
  background: #1d4ed8;
  transform: translateY(-1px);
}

/* 筛选按钮 */
.filter-buttons {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}

.btn-filter {
  padding: 0.4rem 0.8rem;
  border: 1px solid var(--border-color);
  background: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
  color: var(--text-secondary);
  transition: all 0.3s;
  font-weight: 500;
}

.btn-filter:hover {
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.btn-filter.active {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

/* 表格 */
.articles-table-wrapper {
  overflow-x: auto;
  margin-bottom: 1.5rem;
}

.articles-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.95rem;
}

.articles-table thead {
  background-color: var(--bg-color);
}

.articles-table th {
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: var(--text-secondary);
  border-bottom: 1px solid var(--border-color);
}

.articles-table td {
  padding: 1rem;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-primary);
}

.articles-table tbody tr:hover {
  background-color: var(--bg-color);
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 4px;
  font-size: 0.85rem;
  font-weight: 500;
}

.badge-success {
  background-color: #d1fae5;
  color: #047857;
}

.badge-warning {
  background-color: #fef3c7;
  color: #92400e;
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.btn-sm {
  padding: 0.4rem 0.8rem;
  border: 1px solid var(--border-color);
  background: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
  color: var(--primary-color);
  transition: all 0.3s;
  text-decoration: none;
  display: inline-block;
}

.btn-sm:hover {
  background-color: #f1f5f9;
}

.btn-sm.btn-outline {
  border-color: var(--primary-color);
}

.btn-sm.btn-danger {
  color: var(--danger-color);
  border-color: var(--danger-color);
}

.btn-sm.btn-danger:hover {
  background-color: #fef2f2;
}

.empty-state {
  text-align: center;
  padding: 2rem 1rem;
  color: var(--text-tertiary);
}

/* 分页 */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 1.5rem;
}

.page-info {
  color: var(--text-secondary);
  font-size: 0.95rem;
}

@media (max-width: 768px) {
  .dashboard-container {
    padding: 1rem;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
  }

  .articles-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .articles-table th,
  .articles-table td {
    padding: 0.75rem;
    font-size: 0.85rem;
  }

  .action-buttons {
    flex-direction: column;
  }
}
</style>
