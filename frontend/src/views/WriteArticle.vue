<template>
  <div class="write-container">
    <div class="write-header">
      <h1>{{ isEdit ? '编辑文章' : '写文章' }}</h1>
      <div class="header-actions">
        <el-button @click="handleSaveDraft" :loading="draftLoading">
          <i class="el-icon-document-copy"></i>
          保存草稿
        </el-button>
        <el-button type="primary" @click="handlePublish" :loading="publishLoading">
          <i class="el-icon-upload2"></i>
          {{ isEdit ? '更新' : '发布' }}
        </el-button>
        <el-button v-if="isEdit" type="danger" @click="handleDelete" :loading="deleteLoading">
          <i class="el-icon-delete"></i>
          删除文章
        </el-button>
      </div>
    </div>

    <div class="editor-section">
      <!-- 元数据表单 -->
      <div class="article-meta">
        <el-form :model="form" label-width="80px">
          <el-form-item label="标题">
            <el-input v-model="form.title" placeholder="输入文章标题..." maxlength="200" />
          </el-form-item>

          <el-form-item label="摘要">
            <el-input
              v-model="form.summary"
              type="textarea"
              rows="2"
              placeholder="输入文章摘要..."
              maxlength="300"
            />
          </el-form-item>

          <el-form-item label="精选">
            <el-switch v-model="form.is_featured" />
          </el-form-item>
        </el-form>
      </div>

      <!-- Markdown 编辑器 -->
      <div class="editor-wrapper">
        <div class="editor-tabs">
          <div :class="['tab', { active: editorMode === 'edit' }]" @click="editorMode = 'edit'">
            编辑
          </div>
          <div :class="['tab', { active: editorMode === 'preview' }]" @click="editorMode = 'preview'">
            预览
          </div>
          <div :class="['tab', { active: editorMode === 'split' }]" @click="editorMode = 'split'">
            分屏
          </div>
          <div class="tab-spacer" />
          <button class="upload-img-btn" @click="triggerImageUpload" :disabled="imageUploading">
            <i class="bi bi-image"></i>
            {{ imageUploading ? '上传中...' : '插入图片' }}
          </button>
          <input
            ref="imageInputRef"
            type="file"
            accept=".jpg,.jpeg,.png,.gif,.webp"
            style="display:none"
            @change="handleImageUpload"
          />
        </div>

        <div class="editor-content">
          <!-- 编辑模式 -->
          <textarea
            v-if="editorMode === 'edit'"
            v-model="form.content"
            class="markdown-editor"
            placeholder="用 Markdown 格式编写文章内容..."
          />

          <!-- 分屏模式 -->
          <div v-if="editorMode === 'split'" class="split-view">
            <textarea v-model="form.content" class="markdown-editor split" placeholder="Markdown..." />
            <div class="preview-pane">
              <div class="preview-content" v-html="renderMarkdown(form.content)" />
            </div>
          </div>

          <!-- 预览模式 -->
          <div v-if="editorMode === 'preview'" class="preview-pane">
            <div class="preview-content" v-html="renderMarkdown(form.content)" />
          </div>
        </div>
      </div>

      <!-- 草稿状态 -->
      <div v-if="draftStatus" class="draft-status">
        <i class="el-icon-success"></i>
        {{ draftStatus }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import apiClient from '../services/api'
import { useAuthStore } from '../stores/auth'
import { marked } from 'marked'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// 编辑模式
const editorMode = ref<'edit' | 'preview' | 'split'>('split')
const isEdit = computed(() => !!route.params.id)

// 表单数据
const form = reactive({
  title: '',
  content: '',
  summary: '',
  status: 'draft' as 'draft' | 'published',
  is_featured: false,
})

// 加载状态
const draftLoading = ref(false)
const publishLoading = ref(false)
const deleteLoading = ref(false)
const draftStatus = ref('')
const imageUploading = ref(false)
const imageInputRef = ref<HTMLInputElement | null>(null)

// 草稿自动保存
let autoSaveTimer: ReturnType<typeof setInterval> | null = null

function renderMarkdown(content: string): string {
  if (!content) return ''
  try {
    return marked(content) as string
  } catch {
    return `<p>${content.replace(/\n/g, '<br>')}</p>`
  }
}

function triggerImageUpload() {
  imageInputRef.value?.click()
}

async function handleImageUpload(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return

  imageUploading.value = true
  try {
    const formData = new FormData()
    formData.append('file', file)
    const response = await apiClient.post('/upload/image', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    if (response.code === 1000) {
      const url = response.data.url
      const imageMarkdown = `![${file.name}](${url})`
      // 插入到内容末尾（或光标位置）
      form.content = form.content
        ? form.content + '\n' + imageMarkdown
        : imageMarkdown
      ElMessage.success('图片上传成功')
    } else {
      ElMessage.error(response.msg || '上传失败')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '上传失败')
  } finally {
    imageUploading.value = false
    // 清空 input，允许重复上传同一文件
    if (imageInputRef.value) imageInputRef.value.value = ''
  }
}

function getDraftKey(): string {
  const articleId = route.params.id as string
  return articleId ? `article_draft_${articleId}` : 'article_draft_new'
}

function saveDraftToLocalStorage() {
  const draft = {
    title: form.title,
    content: form.content,
    summary: form.summary,
    status: form.status,
    is_featured: form.is_featured,
    timestamp: new Date().toISOString(),
  }
  localStorage.setItem(getDraftKey(), JSON.stringify(draft))
  draftStatus.value = `自动保存于 ${new Date().toLocaleTimeString('zh-CN')}`
}

function loadDraftFromLocalStorage() {
  const draft = localStorage.getItem(getDraftKey())
  if (draft) {
    try {
      const data = JSON.parse(draft)
      form.title = data.title || ''
      form.content = data.content || ''
      form.summary = data.summary || ''
      form.status = data.status || 'draft'
      form.is_featured = data.is_featured || false
      draftStatus.value = `从 ${new Date(data.timestamp).toLocaleTimeString('zh-CN')} 恢复草稿`
    } catch (error) {
      console.error('Failed to load draft:', error)
    }
  }
}

async function loadArticle() {
  if (!isEdit.value) return

  const articleId = route.params.id as string
  try {
    const response = await apiClient.get(`/articles/${articleId}`)
    if (response.code === 1000) {
      const article = response.data
      form.title = article.title || ''
      form.content = article.content || ''
      form.summary = article.summary || ''
      form.status = article.status || 'published'
      form.is_featured = article.is_featured || false
      draftStatus.value = ''
    } else {
      ElMessage.error(response.msg || '加载文章失败')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '加载文章失败')
  }
}

async function handleSaveDraft() {
  if (!form.title.trim()) {
    ElMessage.warning('请输入文章标题')
    return
  }

  draftLoading.value = true
  try {
    let response
    if (isEdit.value) {
      const articleId = route.params.id as string
      response = await apiClient.put(`/author/articles/${articleId}`, {
        title: form.title,
        content: form.content,
        summary: form.summary || form.content.substring(0, 100),
        status: 'draft',
        is_featured: form.is_featured,
        allow_comment: true,
      })
    } else {
      response = await apiClient.post('/articles', {
        title: form.title,
        content: form.content,
        summary: form.summary || form.content.substring(0, 100),
        status: 'draft',
        is_featured: form.is_featured,
        allow_comment: true,
      })
    }

    if (response.code === 1000) {
      saveDraftToLocalStorage()
      ElMessage.success('草稿已保存')
    } else {
      ElMessage.error(response.msg || '保存失败')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '保存失败')
  } finally {
    draftLoading.value = false
  }
}

async function handlePublish() {
  if (!form.title.trim()) {
    ElMessage.warning('请输入文章标题')
    return
  }
  if (!form.content.trim()) {
    ElMessage.warning('请输入文章内容')
    return
  }

  publishLoading.value = true
  try {
    let response

    if (isEdit.value) {
      // 编辑现有文章
      const articleId = route.params.id as string
      response = await apiClient.put(`/author/articles/${articleId}`, {
        title: form.title,
        content: form.content,
        summary: form.summary || form.content.substring(0, 100),
        status: 'published',
        is_featured: form.is_featured,
        allow_comment: true,
      })
    } else {
      // 创建新文章，点"发布"按钮强制 published
      response = await apiClient.post('/articles', {
        title: form.title,
        content: form.content,
        summary: form.summary || form.content.substring(0, 100),
        status: 'published',
        is_featured: form.is_featured,
        allow_comment: true,
      })
    }

    if (response.code === 1000) {
      ElMessage.success(isEdit.value ? '文章已更新' : '文章已发布')
      localStorage.removeItem(getDraftKey())
      // 返回仪表板或首页
      setTimeout(() => {
        if (authStore.user?.role === 'author' || authStore.user?.role === 'admin') {
          router.push('/dashboard')
        } else {
          router.push('/')
        }
      }, 500)
    } else {
      ElMessage.error(response.msg || '发布失败')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '发布失败')
  } finally {
    publishLoading.value = false
  }
}

async function handleDelete() {
  if (!isEdit.value) return

  try {
    await ElMessageBox.confirm('确定要删除这篇文章吗？此操作不可撤销。', '警告', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
    })

    deleteLoading.value = true
    const articleId = route.params.id as string
    const response = await apiClient.delete(`/author/articles/${articleId}`)

    if (response.code === 1000) {
      ElMessage.success('文章已删除')
      localStorage.removeItem(getDraftKey())
      setTimeout(() => {
        router.push('/dashboard')
      }, 500)
    } else {
      ElMessage.error(response.msg || '删除失败')
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  } finally {
    deleteLoading.value = false
  }
}

onMounted(() => {
  if (isEdit.value) {
    loadArticle()
  } else {
    loadDraftFromLocalStorage()
  }

  // 启动自动保存（每30秒）
  autoSaveTimer = setInterval(() => {
    if (form.title.trim() || form.content.trim()) {
      saveDraftToLocalStorage()
    }
  }, 30000)
})

onBeforeUnmount(() => {
  if (autoSaveTimer) {
    clearInterval(autoSaveTimer)
  }
})
</script>

<style scoped>
.write-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
}

.write-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
}

.write-header h1 {
  margin: 0;
  font-size: 28px;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.editor-section {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.article-meta {
  padding: 20px;
  border-bottom: 1px solid #eee;
  background-color: #f5f7fa;
}

.editor-wrapper {
  display: flex;
  flex-direction: column;
  height: 700px;
}

.editor-tabs {
  display: flex;
  align-items: center;
  background-color: #f5f7fa;
  border-bottom: 1px solid #eee;
}

.tab-spacer {
  flex: 1;
}

.upload-img-btn {
  margin-right: 10px;
  padding: 6px 14px;
  font-size: 13px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background: white;
  color: #606266;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s;
}

.upload-img-btn:hover:not(:disabled) {
  color: #409eff;
  border-color: #409eff;
}

.upload-img-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.tab {
  padding: 15px 20px;
  text-align: center;
  cursor: pointer;
  color: #666;
  transition: all 0.3s;
}

.tab.active {
  background-color: white;
  color: #409eff;
  border-bottom: 2px solid #409eff;
  margin-bottom: -1px;
}

.tab:hover {
  background-color: white;
}

.editor-content {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.markdown-editor {
  flex: 1;
  padding: 20px;
  border: none;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  resize: none;
  outline: none;
  background-color: white;
  color: #333;
  line-height: 1.6;
}

.markdown-editor.split {
  border-right: 1px solid #eee;
  background-color: #f9f9f9;
}

.split-view {
  display: flex;
  width: 100%;
  overflow: hidden;
}

.preview-pane {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  background-color: white;
}

.preview-content {
  line-height: 1.8;
  color: #333;
  word-break: break-word;
}

.preview-content h1,
.preview-content h2,
.preview-content h3 {
  margin-top: 20px;
  margin-bottom: 10px;
  color: #333;
  font-weight: bold;
}

.preview-content h1 {
  font-size: 24px;
}

.preview-content h2 {
  font-size: 20px;
}

.preview-content h3 {
  font-size: 18px;
}

.preview-content p {
  margin-bottom: 12px;
}

.preview-content ul,
.preview-content ol {
  margin-left: 20px;
  margin-bottom: 12px;
}

.preview-content li {
  margin-bottom: 6px;
}

.preview-content code {
  background-color: #f5f5f5;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 14px;
}

.preview-content pre {
  background-color: #f5f5f5;
  padding: 12px;
  border-radius: 4px;
  overflow-x: auto;
  margin-bottom: 12px;
}

.preview-content blockquote {
  border-left: 4px solid #409eff;
  padding-left: 12px;
  color: #666;
  margin-left: 0;
  margin-bottom: 12px;
}

.draft-status {
  padding: 10px 20px;
  background-color: #f0f9ff;
  color: #0ea5e9;
  font-size: 12px;
  border-top: 1px solid #bae6fd;
  text-align: right;
}

.draft-status i {
  margin-right: 4px;
}
</style>
