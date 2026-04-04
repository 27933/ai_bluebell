<template>
  <div class="article-container" :class="{ loading }">
    <div v-if="article" class="article-content">
      <!-- 文章头部 -->
      <div class="article-header">
        <h1 class="article-title">{{ article.title }}</h1>

        <div class="author-info">
          <div class="author-avatar">{{ getInitial(article.author_id) }}</div>
          <div>
            <div class="author-name">{{ article.author || 'Unknown' }}</div>
            <small class="text-muted">发布于 {{ formatDate(article.created_at) }}</small>
          </div>
        </div>

        <div class="article-stats">
          <span><i class="bi bi-eye"></i> {{ article.view_count }} 次浏览</span>
          <span><i class="bi bi-heart"></i> {{ article.like_count }} 个赞</span>
          <span><i class="bi bi-chat"></i> {{ article.comment_count }} 条评论</span>
        </div>
      </div>

      <!-- 文章内容 -->
      <div class="article-body">
        <div class="article-text">{{ article.content }}</div>
      </div>

      <!-- 操作按钮 -->
      <div class="article-actions">
        <button
          class="btn btn-primary"
          :class="{ liked: isLiked }"
          @click="handleLike"
          :disabled="!authStore.isLoggedIn"
        >
          <i class="bi" :class="isLiked ? 'bi-heart-fill' : 'bi-heart'"></i>
          {{ isLiked ? '已赞' : '赞' }} ({{ article.like_count }})
        </button>
        <router-link v-if="canEdit" :to="`/write/${article.id}`">
          <button class="btn btn-outline">
            <i class="bi bi-pencil"></i> 编辑
          </button>
        </router-link>
      </div>

      <div class="divider"></div>

      <!-- 评论区域 -->
      <div class="comments-section">
        <h3 class="comments-title">评论 ({{ comments.length }})</h3>

        <!-- 评论表单 -->
        <div v-if="authStore.isLoggedIn" class="comment-form">
          <textarea
            v-model="newComment.content"
            class="form-control"
            rows="3"
            placeholder="说说你的看法..."
          />
          <button
            class="btn btn-primary mt-2"
            @click="handleSubmitComment"
            :disabled="commentLoading"
          >
            {{ commentLoading ? '发表中...' : '发表评论' }}
          </button>
        </div>

        <!-- 评论列表 -->
        <div v-if="comments.length === 0 && authStore.isLoggedIn" class="empty-state">
          <p>暂无评论，成为第一个评论者吧！</p>
        </div>

        <div v-if="comments.length === 0 && !authStore.isLoggedIn" class="empty-state">
          <p>请 <router-link to="/login">登录</router-link> 后发表评论</p>
        </div>

        <div v-else class="comments-list">
          <div v-for="comment in comments" :key="comment.id" class="comment">
            <div class="comment-author">
              <div class="author-avatar-sm">{{ getInitial(comment.user_id) }}</div>
              <div>
                <div class="author-name">{{ comment.user_name || '匿名用户' }}</div>
                <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
              </div>
            </div>
            <div class="comment-content">{{ comment.content }}</div>
            <small><a href="javascript:void(0)">回复</a></small>
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="!loading" class="not-found">
      <p>文章不存在</p>
    </div>

    <div v-else class="loading-state">
      <p>加载中...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import apiClient from '../services/api'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const authStore = useAuthStore()
const loading = ref(false)
const commentLoading = ref(false)
const isLiked = ref(false)

function getInitial(id?: string): string {
  if (!id) return '？'
  return String(id).charAt(0).toUpperCase()
}

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
      console.error('加载失败:', response.msg)
    }
  } catch (error: any) {
    console.error('加载失败:', error.message)
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
    alert('请先登录')
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
    console.error('操作失败:', error)
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
    alert('请输入评论内容')
    return
  }

  commentLoading.value = true

  try {
    const response = await apiClient.post('/comments', {
      article_id: article.value?.id,
      content: newComment.content,
    })

    if (response.code === 1000) {
      alert('评论发表成功')
      newComment.content = ''
      loadComments()
      if (article.value) {
        article.value.comment_count += 1
      }
    } else {
      alert(response.msg || '发表失败')
    }
  } catch (error: any) {
    alert('发表失败：' + (error.message || '未知错误'))
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
:root {
  --primary-color: #2563eb;
  --border-color: #e2e8f0;
  --bg-light: #f8fafc;
  --text-primary: #0f172a;
  --text-secondary: #64748b;
  --text-tertiary: #94a3b8;
}

.article-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem 1.25rem;
}

.article-container.loading {
  opacity: 0.6;
  pointer-events: none;
}

/* 文章头部 */
.article-header {
  margin-bottom: 2rem;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 2rem;
}

.article-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 1rem 0;
  line-height: 1.4;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
}

.author-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary-color), #7c3aed);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 1.2rem;
}

.author-name {
  font-weight: 600;
  color: var(--text-primary);
}

.text-muted {
  color: var(--text-tertiary);
  font-size: 0.9rem;
}

.article-stats {
  display: flex;
  gap: 1.5rem;
  color: var(--text-tertiary);
  font-size: 0.9rem;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 1rem;
  margin-bottom: 1rem;
}

.article-stats span {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

/* 文章内容 */
.article-body {
  margin: 2rem 0;
}

.article-text {
  color: var(--text-primary);
  line-height: 1.8;
  font-size: 1.05rem;
  white-space: pre-wrap;
  word-break: break-word;
}

/* 操作按钮 */
.article-actions {
  margin: 2rem 0;
  display: flex;
  gap: 1rem;
  align-items: center;
}

.btn {
  padding: 0.6rem 1rem;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  text-decoration: none;
  font-size: 0.95rem;
}

.btn-primary {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.btn-primary:hover {
  background: #1d4ed8;
  border-color: #1d4ed8;
}

.btn-primary.liked {
  background: #dc2626;
  border-color: #dc2626;
}

.btn-outline {
  background: white;
  color: var(--primary-color);
  border-color: var(--primary-color);
}

.btn-outline:hover {
  background: var(--bg-light);
}

/* 分割线 */
.divider {
  height: 1px;
  background-color: var(--border-color);
  margin: 2rem 0;
}

/* 评论部分 */
.comments-section {
  margin-top: 2rem;
}

.comments-title {
  font-weight: 700;
  margin-bottom: 1.5rem;
  color: var(--text-primary);
}

.comment-form {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  border: 1px solid var(--border-color);
  margin-bottom: 1.5rem;
}

.form-control {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-family: inherit;
  font-size: 1rem;
  color: var(--text-primary);
}

.form-control:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.mt-2 {
  margin-top: 0.5rem;
}

.comments-list {
  margin-bottom: 1.5rem;
}

.comment {
  padding: 1rem 0;
  border-bottom: 1px solid var(--border-color);
}

.comment:last-child {
  border-bottom: none;
}

.comment-author {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 0.5rem;
}

.author-avatar-sm {
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
  flex-shrink: 0;
}

.comment-time {
  font-size: 0.85rem;
  color: var(--text-tertiary);
}

.comment-content {
  color: var(--text-primary);
  line-height: 1.6;
  margin-bottom: 0.5rem;
}

.comment small a {
  color: var(--primary-color);
  text-decoration: none;
}

.comment small a:hover {
  text-decoration: underline;
}

.empty-state {
  text-align: center;
  padding: 1.5rem;
  background: var(--bg-light);
  border-radius: 8px;
  color: var(--text-secondary);
}

.empty-state a {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 500;
}

.not-found,
.loading-state {
  text-align: center;
  padding: 2rem;
  color: var(--text-secondary);
}

@media (max-width: 768px) {
  .article-container {
    padding: 1rem 0.75rem;
  }

  .article-title {
    font-size: 1.5rem;
  }

  .article-stats {
    flex-direction: column;
    gap: 0.5rem;
  }

  .article-actions {
    flex-wrap: wrap;
  }

  .btn {
    font-size: 0.9rem;
    padding: 0.5rem 0.75rem;
  }
}
</style>
