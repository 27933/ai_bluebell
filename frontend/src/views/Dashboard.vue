<template>
  <div class="dashboard-container">
    <div class="dashboard-header">
      <h1>仪表板</h1>
      <router-link to="/write" class="write-btn">
        <el-button type="primary">
          <i class="el-icon-plus"></i>
          写新文章
        </el-button>
      </router-link>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <el-card class="stat-card">
        <div class="stat-number">{{ stats.published }}</div>
        <div class="stat-label">已发布文章</div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-number">{{ stats.totalViews }}</div>
        <div class="stat-label">总浏览数</div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-number">{{ stats.totalLikes }}</div>
        <div class="stat-label">总获赞数</div>
      </el-card>

      <el-card class="stat-card">
        <div class="stat-number">{{ stats.totalComments }}</div>
        <div class="stat-label">总评论数</div>
      </el-card>
    </div>

    <!-- 阅读趋势图 -->
    <el-card class="chart-card">
      <template #header>
        <div class="chart-header">
          <span>阅读趋势</span>
          <el-button-group size="small">
            <el-button
              :type="timeRange === 'week' ? 'primary' : 'default'"
              @click="timeRange = 'week'"
            >
              周
            </el-button>
            <el-button
              :type="timeRange === 'month' ? 'primary' : 'default'"
              @click="timeRange = 'month'"
            >
              月
            </el-button>
          </el-button-group>
        </div>
      </template>
      <div id="trend-chart" style="height: 400px" />
    </el-card>

    <!-- 文章列表 -->
    <el-card class="articles-card">
      <template #header>
        <div class="card-header">
          <span>我的文章</span>
          <el-button-group>
            <el-button
              :type="filterStatus === 'all' ? 'primary' : 'default'"
              @click="filterStatus = 'all'"
              size="small"
            >
              全部
            </el-button>
            <el-button
              :type="filterStatus === 'published' ? 'primary' : 'default'"
              @click="filterStatus = 'published'"
              size="small"
            >
              已发布
            </el-button>
            <el-button
              :type="filterStatus === 'draft' ? 'primary' : 'default'"
              @click="filterStatus = 'draft'"
              size="small"
            >
              草稿
            </el-button>
          </el-button-group>
        </div>
      </template>

      <el-table
        :data="filteredArticles"
        stripe
        v-loading="loading"
        style="width: 100%"
      >
        <el-table-column prop="title" label="标题" width="300" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'published' ? 'success' : 'info'">
              {{ row.status === 'published' ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="view_count" label="浏览" width="80" align="center" />
        <el-table-column prop="like_count" label="点赞" width="80" align="center" />
        <el-table-column prop="comment_count" label="评论" width="80" align="center" />
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" align="center">
          <template #default="{ row }">
            <router-link :to="`/article/${row.id}`">
              <el-button link type="primary" size="small">查看</el-button>
            </router-link>
            <router-link :to="`/write/${row.id}`">
              <el-button link type="primary" size="small">编辑</el-button>
            </router-link>
            <el-button link type="danger" size="small" @click="handleDelete(row.id)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="prev, pager, next, total"
        style="margin-top: 20px; text-align: right"
        @current-change="loadArticles"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import apiClient from '../services/api'
import { useAuthStore } from '../stores/auth'
import * as echarts from 'echarts'

const authStore = useAuthStore()
const loading = ref(false)
const trendLoading = ref(false)

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

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('zh-CN')
}

async function loadArticles() {
  loading.value = true
  try {
    // 获取当前用户的文章
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
      ElMessage.error(response.msg || '加载失败')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function handleDelete(articleId: string) {
  try {
    await ElMessageBox.confirm('确定要删除这篇文章吗？', '警告', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
    })

    const response = await apiClient.delete(`/author/articles/${articleId}`)
    if (response.code === 1000) {
      ElMessage.success('文章已删除')
      loadArticles()
    } else {
      ElMessage.error(response.msg || '删除失败')
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

async function loadTrendData() {
  trendLoading.value = true
  try {
    const response = await apiClient.get('/article-stats/trend', {
      params: {
        time_range: timeRange.value, // 'week' 或 'month'
        group_by: timeRange.value === 'week' ? 'hour' : 'day',
      },
    })
    if (response.code === 1000) {
      renderChart(response.data)
    } else {
      ElMessage.error(response.msg || '加载趋势数据失败')
    }
  } catch (error: any) {
    console.error('Failed to load trend data:', error)
    // 不显示错误提示，因为可能数据不足
  } finally {
    trendLoading.value = false
  }
}

function renderChart(data: any) {
  const chartDOM = document.getElementById('trend-chart') as HTMLElement
  if (!chartDOM) return

  const myChart = echarts.init(chartDOM)

  // 模拟数据处理（如果后端返回格式不同，需调整）
  const labels = data.labels || data.dates || []
  const values = data.views || data.counts || []

  const option = {
    title: {
      text: '',
      textStyle: {
        color: '#333',
      },
    },
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
          color: '#ddd',
        },
      },
      axisLabel: {
        color: '#666',
        fontSize: 12,
      },
    },
    yAxis: {
      type: 'value',
      axisLine: {
        lineStyle: {
          color: '#ddd',
        },
      },
      axisLabel: {
        color: '#666',
        fontSize: 12,
      },
      splitLine: {
        lineStyle: {
          color: '#eee',
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
          color: '#409eff',
          width: 2,
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(64, 158, 255, 0.3)' },
            { offset: 1, color: 'rgba(64, 158, 255, 0)' },
          ]),
        },
        itemStyle: {
          color: '#409eff',
          borderWidth: 2,
          borderColor: '#fff',
        },
      },
    ],
  }
  myChart.setOption(option)

  // 响应式重绘
  window.addEventListener('resize', () => {
    myChart.resize()
  })
}

onMounted(() => {
  loadArticles()
  loadTrendData()
})

// 监听时间范围变化
watch(() => timeRange.value, () => {
  loadTrendData()
})
</script>

<style scoped>
.dashboard-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
}

.dashboard-header h1 {
  margin: 0;
  font-size: 28px;
  color: #333;
}

.write-btn {
  text-decoration: none;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  text-align: center;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.stat-number {
  font-size: 32px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 10px;
}

.stat-label {
  color: #666;
  font-size: 14px;
}

.articles-card {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.chart-card {
  margin-bottom: 30px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}
</style>
