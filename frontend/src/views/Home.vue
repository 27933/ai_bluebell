<template>
  <div class="home-container">
    <div class="home-header">
      <h1>Bluebell 博客</h1>
      <p>发现优质内容，分享知识故事</p>
    </div>

    <div class="articles-section">
      <div class="section-header">
        <h2>最新文章</h2>
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="handlePageChange"
        />
      </div>

      <el-empty v-if="loading" description="加载中..." />

      <div v-else-if="articles.length === 0" class="empty-state">
        <el-empty description="暂无文章" />
      </div>

      <div v-else class="articles-grid">
        <el-card v-for="article in articles" :key="article.id" class="article-card">
          <div class="article-header">
            <router-link :to="`/article/${article.id}`" class="article-title">
              {{ article.title }}
            </router-link>
            <span v-if="article.is_featured" class="featured-badge">精选</span>
          </div>

          <p class="article-summary">{{ article.summary }}</p>

          <div class="article-meta">
            <span class="meta-item">
              <i class="el-icon-view" />
              {{ article.view_count }}
            </span>
            <span class="meta-item">
              <i class="el-icon-thumb" />
              {{ article.like_count }}
            </span>
            <span class="meta-item">
              <i class="el-icon-chat-dot-round" />
              {{ article.comment_count }}
            </span>
            <span class="meta-date">{{ formatDate(article.created_at) }}</span>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import apiClient from '../services/api'
import { ElMessage } from 'element-plus'

interface Article {
  id: string
  title: string
  summary: string
  view_count: number
  like_count: number
  comment_count: number
  is_featured: boolean
  created_at: string
}

const articles = ref<Article[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('zh-CN')
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
      ElMessage.error(response.msg || '加载失败')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '加载失败')
  } finally {
    loading.value = false
  }
}

function handlePageChange() {
  loadArticles()
}

onMounted(() => {
  loadArticles()
})
</script>

<style scoped>
.home-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 40px 20px;
}

.home-header {
  text-align: center;
  margin-bottom: 40px;
  padding: 40px 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 8px;
}

.home-header h1 {
  font-size: 40px;
  margin: 0 0 10px 0;
}

.home-header p {
  font-size: 16px;
  margin: 0;
  opacity: 0.9;
}

.articles-section {
  margin-top: 40px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h2 {
  margin: 0;
  color: #333;
  font-size: 24px;
}

.articles-grid {
  display: grid;
  gap: 20px;
}

.article-card {
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
}

.article-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.article-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.article-title {
  font-size: 18px;
  font-weight: bold;
  color: #333;
  text-decoration: none;
  flex: 1;
}

.article-title:hover {
  color: #409eff;
}

.featured-badge {
  background-color: #f56c6c;
  color: white;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  margin-left: 8px;
}

.article-summary {
  color: #666;
  line-height: 1.6;
  margin: 10px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-meta {
  display: flex;
  gap: 20px;
  font-size: 12px;
  color: #999;
  margin-top: 12px;
}

.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.meta-date {
  margin-left: auto;
}

.empty-state {
  padding: 60px 20px;
  text-align: center;
}
</style>
