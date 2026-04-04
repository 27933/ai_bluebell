<template>
  <div class="article-container" v-loading="loading">
    <div v-if="article" class="article-content">
      <div class="article-header">
        <h1>{{ article.title }}</h1>
        <div class="article-meta">
          <span>浏览: {{ article.view_count }}</span>
          <span>点赞: {{ article.like_count }}</span>
          <span>评论: {{ article.comment_count }}</span>
          <span>{{ formatDate(article.created_at) }}</span>
        </div>
      </div>

      <div class="article-body">
        <div class="article-text">{{ article.content }}</div>
      </div>

      <div class="article-actions">
        <el-button type="primary" @click="handleLike">
          <i :class="isLiked ? 'el-icon-star-on' : 'el-icon-star-off'" />
          {{ isLiked ? '已赞' : '点赞' }}
        </el-button>
        <router-link v-if="canEdit" :to="`/write/${article.id}`">
          <el-button type="info">
            <i class="el-icon-edit"></i>
            编辑
          </el-button>
        </router-link>
      </div>

      <el-divider></el-divider>

      <!-- 评论列表 -->
      <div class="comments-section">
        <h3>评论 ({{ comments.length }})</h3>

        <div v-if="comments.length === 0" class="empty-comments">
          <el-empty description="暂无评论" />
        </div>

        <div v-else class="comments-list">
          <div v-for="comment in comments" :key="comment.id" class="comment-item">
            <div class="comment-header">
              <strong>{{ comment.user_name || '匿名用户' }}</strong>
              <span class="comment-date">{{ formatDate(comment.created_at) }}</span>
            </div>
            <div class="comment-body">{{ comment.content }}</div>
          </div>
        </div>

        <!-- 评论表单 -->
        <div v-if="authStore.isLoggedIn" class="comment-form">
          <el-form :model="newComment">
            <el-form-item>
              <el-input
                v-model="newComment.content"
                type="textarea"
                rows="3"
                placeholder="发表评论..."
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleSubmitComment" :loading="commentLoading">
                发表评论
              </el-button>
            </el-form-item>
          </el-form>
        </div>

        <div v-else class="login-prompt">
          <p>请 <router-link to="/login">登录</router-link> 后发表评论</p>
        </div>
      </div>
    </div>

    <div v-else-if="!loading" class="not-found">
      <el-empty description="文章不存在" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import apiClient from '../services/api'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const authStore = useAuthStore()
const loading = ref(false)
const commentLoading = ref(false)
const isLiked = ref(false)

interface Article {
  id: string
  title: string
  content: string
  view_count: number
  like_count: number
  comment_count: number
  created_at: string
  author_id?: string
}

interface Comment {
  id: string
  content: string
  user_name?: string
  created_at: string
}

const article = ref<Article | null>(null)
const comments = ref<Comment[]>([])

const newComment = reactive({
  content: '',
})

const canEdit = computed(() => {
  return (
    article.value &&
    authStore.isLoggedIn &&
    (authStore.user?.role === 'author' || authStore.user?.role === 'admin')
  )
})

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('zh-CN')
}

async function loadArticle() {
  const articleId = route.params.id as string
  loading.value = true

  try {
    const response = await apiClient.get(`/articles/${articleId}`)

    if (response.code === 1000) {
      article.value = response.data
    } else {
      ElMessage.error(response.msg || '加载失败')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '加载失败')
  } finally {
    loading.value = false
  }
}

async function loadComments() {
  const articleId = route.params.id as string

  try {
    const response = await apiClient.get('/comments', {
      params: {
        article_id: articleId,
      },
    })

    if (response.code === 1000) {
      comments.value = response.data.list || []
    }
  } catch (error) {
    console.error('加载评论失败:', error)
  }
}

async function handleLike() {
  if (!authStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }

  try {
    if (isLiked.value) {
      await apiClient.delete('/likes', {
        params: {
          target_id: article.value?.id,
          target_type: 'article',
        },
      })
    } else {
      await apiClient.post('/likes', {
        target_id: article.value?.id,
        target_type: 'article',
      })
    }

    isLiked.value = !isLiked.value
    if (article.value) {
      article.value.like_count += isLiked.value ? 1 : -1
    }
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  }
}

async function loadLikeStatus() {
  const articleId = route.params.id as string
  if (!authStore.isLoggedIn) return // 未登录不查询

  try {
    const response = await apiClient.get('/likes/status', {
      params: {
        target_type: 'article',
        target_id: articleId,
      },
    })
    if (response.code === 1000) {
      isLiked.value = response.data?.is_liked || false
    }
  } catch (error) {
    console.error('Failed to load like status:', error)
  }
}

async function handleSubmitComment() {
  if (!newComment.content.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }

  commentLoading.value = true

  try {
    const response = await apiClient.post('/comments', {
      article_id: article.value?.id,
      content: newComment.content,
    })

    if (response.code === 1000) {
      ElMessage.success('评论发表成功')
      newComment.content = ''
      loadComments()
      if (article.value) {
        article.value.comment_count += 1
      }
    } else {
      ElMessage.error(response.msg || '发表失败')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '发表失败')
  } finally {
    commentLoading.value = false
  }
}

onMounted(() => {
  loadArticle()
  loadComments()
  loadLikeStatus()
})
</script>

<style scoped>
.article-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 40px 20px;
}

.article-header {
  margin-bottom: 30px;
  border-bottom: 1px solid #eee;
  padding-bottom: 20px;
}

.article-header h1 {
  font-size: 32px;
  color: #333;
  margin: 0 0 10px 0;
  line-height: 1.4;
}

.article-meta {
  display: flex;
  gap: 20px;
  color: #999;
  font-size: 14px;
}

.article-body {
  margin: 30px 0;
}

.article-text {
  color: #333;
  line-height: 1.8;
  font-size: 16px;
  white-space: pre-wrap;
  word-break: break-word;
}

.article-actions {
  margin: 30px 0;
  display: flex;
  gap: 10px;
}

.comments-section {
  margin-top: 40px;
}

.comments-section h3 {
  color: #333;
  margin-bottom: 20px;
}

.comments-list {
  margin-bottom: 30px;
}

.comment-item {
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
  margin-bottom: 15px;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 14px;
}

.comment-date {
  color: #999;
}

.comment-body {
  color: #333;
  line-height: 1.6;
}

.comment-form {
  background-color: #f5f7fa;
  padding: 20px;
  border-radius: 4px;
}

.login-prompt {
  text-align: center;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 4px;
  color: #666;
}

.login-prompt a {
  color: #409eff;
  text-decoration: none;
}

.login-prompt a:hover {
  text-decoration: underline;
}

.empty-comments {
  padding: 40px 0;
}

.not-found {
  padding: 100px 0;
}
</style>
