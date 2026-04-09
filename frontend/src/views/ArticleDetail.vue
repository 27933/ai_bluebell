<template>
  <div class="article-container" :class="{ loading }">
    <div v-if="article" class="article-content">
      <!-- 文章头部 -->
      <div class="article-header">
        <div class="article-title-row">
          <h1 class="article-title">{{ article.title }}</h1>
          <span v-if="article.is_featured" class="featured-badge">
            <i class="bi bi-star-fill"></i> 精选
          </span>
        </div>

        <div class="author-info">
          <div class="author-avatar">{{ getInitial(article.author?.nickname || article.author?.username) }}</div>
          <div>
            <div class="author-name">{{ article.author?.nickname || article.author?.username || 'Unknown' }}</div>
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
        <div class="article-text markdown-body" v-html="renderMarkdown(article.content)" />
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
        <button v-if="canDelete" class="btn btn-outline btn-danger" @click="handleDeleteArticle">
          <i class="bi bi-trash"></i> 删除
        </button>
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
        <div v-if="topLevelComments.length === 0 && authStore.isLoggedIn" class="empty-state">
          <p>暂无评论，成为第一个评论者吧！</p>
        </div>

        <div v-if="topLevelComments.length === 0 && !authStore.isLoggedIn" class="empty-state">
          <p>请 <router-link to="/login">登录</router-link> 后发表评论</p>
        </div>

        <div v-else class="comments-list">
          <div v-for="comment in topLevelComments" :key="comment.id" class="comment">
            <!-- 墓碑：已删除但有回复的评论 -->
            <template v-if="comment.status === 'deleted'">
              <div class="comment-deleted">
                <i class="bi bi-slash-circle"></i> 该评论已被删除
              </div>
            </template>
            <template v-else>
              <div class="comment-author">
                <div class="author-avatar-sm">{{ getInitial(comment.author?.nickname || comment.author?.username) }}</div>
                <div class="comment-author-info">
                  <div class="author-name">{{ comment.author?.nickname || comment.author?.username || '匿名用户' }}</div>
                  <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
                </div>
              </div>
              <template v-if="editingCommentId === comment.id">
                <div class="edit-form">
                  <textarea
                    v-model="editingContent"
                    class="form-control"
                    rows="3"
                    autofocus
                  />
                  <div class="reply-form-actions">
                    <button
                      class="btn btn-primary btn-sm"
                      @click="handleUpdateComment(comment.id)"
                      :disabled="commentLoading"
                    >
                      {{ commentLoading ? '保存中...' : '保存' }}
                    </button>
                    <button class="btn btn-outline btn-sm" @click="cancelEdit">取消</button>
                  </div>
                </div>
              </template>
              <template v-else>
                <div class="comment-content">{{ comment.content }}</div>
                <div class="comment-actions">
                  <a
                    v-if="authStore.isLoggedIn"
                    href="javascript:void(0)"
                    class="action-link"
                    @click="toggleReply(comment.id)"
                  >
                    <i class="bi bi-reply"></i> 回复
                  </a>
                  <a
                    v-if="canEditComment(comment)"
                    href="javascript:void(0)"
                    class="action-link"
                    @click="startEdit(comment)"
                  >
                    <i class="bi bi-pencil"></i> 编辑
                  </a>
                  <a
                    v-if="canDeleteComment(comment)"
                    href="javascript:void(0)"
                    class="action-link danger"
                    @click="handleDeleteComment(comment.id)"
                  >
                    <i class="bi bi-trash"></i> 删除
                  </a>
                </div>
              </template>
            </template>

            <!-- 回复表单 -->
            <div v-if="replyingTo === comment.id" class="reply-form">
              <textarea
                v-model="replyContent"
                class="form-control"
                rows="2"
                :placeholder="`回复 @${comment.author?.nickname || comment.author?.username}...`"
                autofocus
              />
              <div class="reply-form-actions">
                <button
                  class="btn btn-primary btn-sm"
                  @click="handleSubmitReply(comment.id)"
                  :disabled="commentLoading"
                >
                  {{ commentLoading ? '发送中...' : '发送回复' }}
                </button>
                <button class="btn btn-outline btn-sm" @click="replyingTo = null">取消</button>
              </div>
            </div>

            <!-- 嵌套回复 -->
            <div v-if="getReplies(comment.id).length > 0" class="replies-list">
              <div v-for="reply in getReplies(comment.id)" :key="reply.id" class="comment reply-comment">
                <div class="comment-author">
                  <div class="author-avatar-sm">{{ getInitial(reply.author?.nickname || reply.author?.username) }}</div>
                  <div class="comment-author-info">
                    <div class="author-name">{{ reply.author?.nickname || reply.author?.username || '匿名用户' }}</div>
                    <span class="comment-time">{{ formatDate(reply.created_at) }}</span>
                  </div>
                </div>
                <template v-if="editingCommentId === reply.id">
                  <div class="edit-form">
                    <textarea
                      v-model="editingContent"
                      class="form-control"
                      rows="2"
                      autofocus
                    />
                    <div class="reply-form-actions">
                      <button
                        class="btn btn-primary btn-sm"
                        @click="handleUpdateComment(reply.id)"
                        :disabled="commentLoading"
                      >
                        {{ commentLoading ? '保存中...' : '保存' }}
                      </button>
                      <button class="btn btn-outline btn-sm" @click="cancelEdit">取消</button>
                    </div>
                  </div>
                </template>
                <template v-else>
                  <div class="comment-content">{{ reply.content }}</div>
                  <div class="comment-actions">
                    <a
                      v-if="canEditComment(reply)"
                      href="javascript:void(0)"
                      class="action-link"
                      @click="startEdit(reply)"
                    >
                      <i class="bi bi-pencil"></i> 编辑
                    </a>
                    <a
                      v-if="canDeleteComment(reply)"
                      href="javascript:void(0)"
                      class="action-link danger"
                      @click="handleDeleteComment(reply.id)"
                    >
                      <i class="bi bi-trash"></i> 删除
                    </a>
                  </div>
                </template>
              </div>
            </div>
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
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import apiClient from '../services/api'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)
const commentLoading = ref(false)
const isLiked = ref(false)
const replyingTo = ref<string | null>(null)
const replyContent = ref('')
const editingCommentId = ref<string | null>(null)
const editingContent = ref('')

function getInitial(name?: string): string {
  if (!name) return '？'
  return name.charAt(0).toUpperCase()
}

interface Author {
  id: string
  username: string
  nickname?: string
  avatar?: string
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
  author?: Author
  is_featured: boolean
}

interface Comment {
  id: string
  content: string
  status?: string
  author?: Author
  parent_id?: string | null
  created_at: string
}

const article = ref<Article | null>(null)
const comments = ref<Comment[]>([])

const newComment = reactive({
  content: '',
})

const topLevelComments = computed(() =>
  comments.value.filter((c) => !c.parent_id)
)

function getReplies(parentId: string): Comment[] {
  return comments.value.filter((c) => c.parent_id === parentId)
}

function canDeleteComment(comment: Comment): boolean {
  if (!authStore.isLoggedIn || !authStore.user) return false
  const uid = authStore.user.id
  const isCommentAuthor = comment.author?.id === uid
  // API 返回顶层 author_id，不是嵌套 author.id
  const isArticleAuthor =
    article.value?.author_id === uid &&
    (authStore.user.role === 'author' || authStore.user.role === 'admin')
  const isAdmin = authStore.user.role === 'admin'
  return isCommentAuthor || isArticleAuthor || isAdmin
}

function canEditComment(comment: Comment): boolean {
  if (!authStore.isLoggedIn || !authStore.user) return false
  return comment.author?.id === authStore.user.id
}

function startEdit(comment: Comment) {
  editingCommentId.value = comment.id
  editingContent.value = comment.content
  replyingTo.value = null
}

function cancelEdit() {
  editingCommentId.value = null
  editingContent.value = ''
}

async function handleUpdateComment(commentId: string) {
  if (!editingContent.value.trim()) return
  commentLoading.value = true
  try {
    const response = await apiClient.put(`/comments/${commentId}`, {
      content: editingContent.value.trim(),
    })
    if (response.code === 1000) {
      const c = comments.value.find((c) => c.id === commentId)
      if (c) c.content = editingContent.value.trim()
      cancelEdit()
    } else {
      ElMessage.error(response.msg || '编辑失败')
    }
  } catch (error: any) {
    ElMessage.error('编辑失败')
  } finally {
    commentLoading.value = false
  }
}

function toggleReply(commentId: string) {
  if (replyingTo.value === commentId) {
    replyingTo.value = null
  } else {
    replyingTo.value = commentId
    replyContent.value = ''
  }
}

const canEdit = computed(() => {
  if (!article.value || !authStore.isLoggedIn || !authStore.user) return false
  const uid = authStore.user.id
  return article.value.author_id === uid || authStore.user.role === 'admin'
})

const canDelete = computed(() => canEdit.value)

function renderMarkdown(content: string): string {
  if (!content) return ''
  try {
    const raw = marked(content) as string
    return DOMPurify.sanitize(raw)
  } catch {
    return DOMPurify.sanitize(`<p>${content.replace(/\n/g, '<br>')}</p>`)
  }
}

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
    ElMessage.warning('请先登录')
    return
  }

  try {
    if (isLiked.value) {
      await apiClient.delete('/likes', {
        params: {
          target_id: Number(article.value?.id),
          target_type: 'article',
        },
      })
    } else {
      await apiClient.post('/likes', {
        target_id: Number(article.value?.id),
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

async function handleDeleteComment(commentId: string) {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '删除确认', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
    })
  } catch {
    return // 用户取消
  }

  try {
    const response = await apiClient.delete(`/comments/${commentId}`)
    if (response.code === 1000) {
      ElMessage.success('评论已删除')
      // 重新拉取评论列表，让后端决定删除评论是否需要墓碑展示
      await loadComments()
      if (article.value) {
        article.value.comment_count = Math.max(0, article.value.comment_count - 1)
      }
    } else {
      ElMessage.error(response.msg || '删除失败')
    }
  } catch (error: any) {
    ElMessage.error('删除失败：' + (error.message || '未知错误'))
  }
}

async function handleDeleteArticle() {
  try {
    await ElMessageBox.confirm('确定要删除这篇文章吗？此操作不可撤销。', '删除确认', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
    })
  } catch {
    return
  }
  try {
    const response = await apiClient.delete(`/author/articles/${article.value?.id}`)
    if (response.code === 1000) {
      ElMessage.success('文章已删除')
      router.push('/')
    } else {
      ElMessage.error(response.msg || '删除失败')
    }
  } catch (error: any) {
    ElMessage.error('删除失败：' + (error.message || '未知错误'))
  }
}

async function handleSubmitReply(parentId: string) {
  if (!replyContent.value.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }

  commentLoading.value = true
  try {
    const response = await apiClient.post('/comments', {
      article_id: Number(article.value?.id),
      parent_id: Number(parentId),
      content: replyContent.value,
    })

    if (response.code === 1000) {
      ElMessage.success('回复成功')
      replyContent.value = ''
      replyingTo.value = null
      await loadComments()
      if (article.value) {
        article.value.comment_count += 1
      }
    } else {
      ElMessage.error(response.msg || '回复失败')
    }
  } catch (error: any) {
    ElMessage.error('回复失败：' + (error.message || '未知错误'))
  } finally {
    commentLoading.value = false
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
      article_id: Number(article.value?.id),
      content: newComment.content,
    })

    if (response.code === 1000) {
      ElMessage.success('评论发表成功')
      newComment.content = ''
      await loadComments()
      if (article.value) {
        article.value.comment_count += 1
      }
    } else {
      ElMessage.error(response.msg || '发表失败')
    }
  } catch (error: any) {
    ElMessage.error('发表失败：' + (error.message || '未知错误'))
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

.article-title-row {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.article-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
  line-height: 1.4;
  flex: 1;
}

.featured-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  background-color: #fef3c7;
  color: #b45309;
  border: 1px solid #fde68a;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.25rem 0.6rem;
  border-radius: 20px;
  white-space: nowrap;
  margin-top: 0.5rem;
  flex-shrink: 0;
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
  word-break: break-word;
}

/* Markdown 渲染样式（:deep 穿透 scoped） */
:deep(.markdown-body h1),
:deep(.markdown-body h2),
:deep(.markdown-body h3),
:deep(.markdown-body h4),
:deep(.markdown-body h5),
:deep(.markdown-body h6) {
  margin: 1.2em 0 0.6em;
  font-weight: 600;
  line-height: 1.4;
}
:deep(.markdown-body h1) { font-size: 1.8rem; }
:deep(.markdown-body h2) { font-size: 1.5rem; }
:deep(.markdown-body h3) { font-size: 1.25rem; }
:deep(.markdown-body p)  { margin: 0.8em 0; }
:deep(.markdown-body img) {
  max-width: 100%;
  border-radius: 6px;
  margin: 0.5em 0;
}
:deep(.markdown-body code) {
  background: #f3f4f6;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: monospace;
  font-size: 0.9em;
}
:deep(.markdown-body pre) {
  background: #f3f4f6;
  padding: 1em;
  border-radius: 6px;
  overflow-x: auto;
}
:deep(.markdown-body pre code) {
  background: none;
  padding: 0;
}
:deep(.markdown-body blockquote) {
  border-left: 4px solid #ddd;
  padding-left: 1em;
  color: #666;
  margin: 0.8em 0;
}
:deep(.markdown-body ul),
:deep(.markdown-body ol) {
  padding-left: 1.5em;
  margin: 0.8em 0;
}
:deep(.markdown-body hr) {
  border: none;
  border-top: 2px solid #d1d5db;
  margin: 1.5em 0;
}
:deep(.markdown-body a) {
  color: #409eff;
  text-decoration: none;
}
:deep(.markdown-body a:hover) { text-decoration: underline; }

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

.btn-danger {
  color: #dc3545;
  border-color: #dc3545;
}

.btn-danger:hover {
  background: #dc3545;
  color: white;
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

.comment-deleted {
  color: var(--text-secondary, #999);
  font-style: italic;
  padding: 0.4rem 0;
  font-size: 0.9rem;
}

.comment-author-info {
  flex: 1;
}

.comment-actions {
  display: flex;
  gap: 1rem;
  margin-top: 0.4rem;
}

.action-link {
  color: var(--text-tertiary);
  font-size: 0.85rem;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  transition: color 0.2s;
}

.action-link:hover {
  color: var(--primary-color);
}

.action-link.danger:hover {
  color: #dc2626;
}

.reply-form,
.edit-form {
  margin-top: 0.75rem;
  padding: 0.75rem;
  background: #f8fafc;
  border-radius: 6px;
  border: 1px solid var(--border-color);
}

.reply-form-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.btn-sm {
  padding: 0.35rem 0.75rem;
  font-size: 0.875rem;
}

.replies-list {
  margin-top: 0.75rem;
  margin-left: 2.5rem;
  border-left: 2px solid var(--border-color);
  padding-left: 1rem;
}

.reply-comment {
  padding: 0.75rem 0;
  border-bottom: 1px solid #f1f5f9;
}

.reply-comment:last-child {
  border-bottom: none;
  padding-bottom: 0;
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

  /* 嵌套回复缩进在小屏减小 */
  .replies-list {
    margin-left: 1rem;
  }
}
</style>
