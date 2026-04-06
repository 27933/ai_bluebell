<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <template #header>
        <div class="card-header">
          <span class="title">个人资料</span>
        </div>
      </template>

      <div v-if="profileLoading" class="profile-content" style="text-align:center;padding:20px;color:#999;">
        加载中...
      </div>

      <div v-else-if="!authStore.user" class="profile-content" style="text-align:center;padding:20px;color:#999;">
        无法加载用户信息，请 <router-link to="/login">重新登录</router-link>
      </div>

      <div v-else class="profile-content">
        <el-form :model="form" label-width="100px">
          <el-form-item label="用户名">
            <span>{{ authStore.user.username }}</span>
          </el-form-item>

          <el-form-item label="邮箱">
            <el-input v-model="form.email" placeholder="请输入邮箱" />
          </el-form-item>

          <el-form-item label="昵称">
            <el-input v-model="form.nickname" placeholder="请输入昵称" />
          </el-form-item>

          <el-form-item label="个人介绍">
            <el-input
              v-model="form.bio"
              type="textarea"
              rows="4"
              placeholder="请输入个人介绍"
            />
          </el-form-item>

          <el-form-item label="角色">
            <span>{{ roleText }}</span>
          </el-form-item>

          <el-form-item label="加入时间">
            <span>{{ formatDate(authStore.user.created_at) }}</span>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="handleUpdate" :loading="loading">
              更新资料
            </el-button>
            <el-button @click="handleLogout">退出登录</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import apiClient from '../services/api'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)
const profileLoading = ref(true)

const form = reactive({
  email: '',
  nickname: '',
  bio: '',
})

const roleText = computed(() => {
  const roles: Record<string, string> = {
    reader: '读者',
    author: '作者',
    admin: '管理员',
    visitor: '访客',
  }
  return roles[authStore.user?.role || 'visitor']
})

function formatDate(date?: string) {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

async function handleUpdate() {
  // 前端邮箱格式校验
  if (form.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    ElMessage.error('邮箱格式不正确')
    return
  }

  loading.value = true
  try {
    // 只发送非空字段，避免空字符串触发后端验证
    const data: Record<string, string> = {}
    if (form.nickname) data.nickname = form.nickname
    if (form.email) data.email = form.email
    if (form.bio) data.bio = form.bio

    const response = await apiClient.put('/auth/profile', data)

    if (response.code === 1000) {
      ElMessage.success('资料更新成功')
      if (authStore.user) {
        if (form.email) authStore.user.email = form.email
        if (form.nickname) authStore.user.nickname = form.nickname
        if (form.bio) authStore.user.bio = form.bio
      }
    } else {
      ElMessage.error(response.msg || '更新失败')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '更新失败')
  } finally {
    loading.value = false
  }
}

function handleLogout() {
  authStore.logout()
  ElMessage.success('已退出登录')
  setTimeout(() => {
    router.push('/login')
  }, 500)
}

async function loadProfile() {
  try {
    const response = await apiClient.get('/auth/profile')
    if (response.code === 1000) {
      const data = response.data
      form.email = data.email || ''
      form.nickname = data.nickname || ''
      form.bio = data.bio || ''
      authStore.setUser(data)
    }
  } catch (error: any) {
    console.error('加载用户信息失败:', error)
    if (authStore.user) {
      form.email = authStore.user.email || ''
      form.nickname = authStore.user.nickname || ''
      form.bio = authStore.user.bio || ''
    }
  } finally {
    profileLoading.value = false
  }
}

onMounted(() => {
  loadProfile()
})
</script>

<style scoped>
.profile-container {
  max-width: 600px;
  margin: 20px auto;
  padding: 20px;
}

.profile-card {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  text-align: center;
}

.title {
  font-size: 20px;
  font-weight: bold;
  color: #333;
}

.profile-content {
  padding: 20px 0;
}
</style>
