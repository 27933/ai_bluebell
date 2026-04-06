<template>
  <div class="stats-page">
    <!-- 概览卡片 -->
    <div class="stat-cards">
      <div class="stat-card">
        <div class="stat-icon users"><i class="bi bi-people-fill"></i></div>
        <div class="stat-info">
          <div class="stat-value">{{ overview?.user_count ?? '-' }}</div>
          <div class="stat-label">总用户数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon articles"><i class="bi bi-file-earmark-text-fill"></i></div>
        <div class="stat-info">
          <div class="stat-value">{{ overview?.article_count ?? '-' }}</div>
          <div class="stat-label">总文章数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon comments"><i class="bi bi-chat-fill"></i></div>
        <div class="stat-info">
          <div class="stat-value">{{ overview?.comment_count ?? '-' }}</div>
          <div class="stat-label">总评论数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon today"><i class="bi bi-calendar-check-fill"></i></div>
        <div class="stat-info">
          <div class="stat-value">{{ overview?.today_new_user_count ?? '-' }}</div>
          <div class="stat-label">今日新增用户</div>
        </div>
      </div>
    </div>

    <!-- 趋势图 -->
    <div class="chart-card">
      <div class="chart-header">
        <h3>最近 30 天新增趋势</h3>
      </div>
      <div ref="chartRef" class="chart-container"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import { getSystemOverview, getSystemDailyStats } from '../../api/admin'
import type { SystemOverview, DailyStatItem } from '../../api/admin'
import { ElMessage } from 'element-plus'

const overview = ref<SystemOverview | null>(null)
const chartRef = ref<HTMLElement | null>(null)
let chart: echarts.ECharts | null = null

onMounted(async () => {
  await Promise.all([loadOverview(), loadDailyStats()])
})

onUnmounted(() => {
  chart?.dispose()
})

async function loadOverview() {
  try {
    const res = await getSystemOverview()
    if (res.code === 1000) {
      overview.value = res.data
    }
  } catch {
    ElMessage.error('加载概览数据失败')
  }
}

async function loadDailyStats() {
  // 最近 30 天
  const end = new Date()
  const start = new Date()
  start.setDate(start.getDate() - 29)
  const fmt = (d: Date) => d.toISOString().slice(0, 10)

  try {
    const res = await getSystemDailyStats({ start_date: fmt(start), end_date: fmt(end) })
    if (res.code === 1000 && chartRef.value) {
      renderChart(res.data || [])
    }
  } catch {
    ElMessage.error('加载趋势数据失败')
  }
}

function renderChart(data: DailyStatItem[]) {
  if (!chartRef.value) return
  chart = echarts.init(chartRef.value)
  chart.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['新增用户', '新增文章', '新增评论'], bottom: 0 },
    grid: { left: 40, right: 20, top: 20, bottom: 40 },
    xAxis: {
      type: 'category',
      data: data.map(d => d.date.slice(5, 10)),  // "2026-04-06T00:00:00Z" → "04-06"
      axisLabel: { rotate: 30, fontSize: 11 },
    },
    yAxis: { type: 'value', minInterval: 1 },
    series: [
      {
        name: '新增用户',
        type: 'line',
        smooth: true,
        data: data.map(d => d.new_user_count),
        itemStyle: { color: '#4f46e5' },
      },
      {
        name: '新增文章',
        type: 'line',
        smooth: true,
        data: data.map(d => d.new_article_count),
        itemStyle: { color: '#059669' },
      },
      {
        name: '新增评论',
        type: 'line',
        smooth: true,
        data: data.map(d => d.new_comment_count),
        itemStyle: { color: '#d97706' },
      },
    ],
  })
}
</script>

<style scoped>
.stats-page {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

/* 概览卡片 */
.stat-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
}

.stat-card {
  background: white;
  border-radius: 10px;
  padding: 1.25rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: 0 1px 3px rgba(0,0,0,0.07);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.4rem;
  color: white;
  flex-shrink: 0;
}

.stat-icon.users    { background: #4f46e5; }
.stat-icon.articles { background: #059669; }
.stat-icon.comments { background: #d97706; }
.stat-icon.today    { background: #0891b2; }

.stat-value {
  font-size: 1.6rem;
  font-weight: 700;
  color: #1e293b;
  line-height: 1;
}

.stat-label {
  font-size: 0.85rem;
  color: #64748b;
  margin-top: 0.25rem;
}

/* 趋势图 */
.chart-card {
  background: white;
  border-radius: 10px;
  padding: 1.25rem;
  box-shadow: 0 1px 3px rgba(0,0,0,0.07);
}

.chart-header h3 {
  font-size: 1rem;
  font-weight: 600;
  color: #1e293b;
  margin: 0 0 1rem;
}

.chart-container {
  height: 300px;
}

@media (max-width: 900px) {
  .stat-cards {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
