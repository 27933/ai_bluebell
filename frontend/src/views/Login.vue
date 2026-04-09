<template>
  <div class="auth-page">
    <!-- 左侧品牌区 -->
    <div class="brand-panel">
      <div class="brand-content">
        <div class="brand-logo">
          <i class="bi bi-flower1"></i>
        </div>
        <h1 class="brand-name">Bluebell</h1>
        <p class="brand-tagline">知识博客平台</p>
        <div class="brand-features">
          <div class="feature-item">
            <i class="bi bi-pencil-square"></i>
            <span>Markdown 编辑器</span>
          </div>
          <div class="feature-item">
            <i class="bi bi-bar-chart-line"></i>
            <span>数据统计仪表板</span>
          </div>
          <div class="feature-item">
            <i class="bi bi-shield-check"></i>
            <span>JWT 双 Token 认证</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧表单区 -->
    <div class="form-panel">
      <div class="form-wrapper">
        <div class="form-header">
          <h2>欢迎回来</h2>
          <p>登录你的 Bluebell 账户</p>
        </div>

        <el-form ref="formRef" :model="form" :rules="rules" @submit.prevent="handleSubmit">
          <el-form-item prop="username">
            <el-input
              v-model="form.username"
              placeholder="用户名"
              size="large"
              clearable
              :prefix-icon="UserIcon"
            />
          </el-form-item>

          <el-form-item prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="密码"
              size="large"
              show-password
              :prefix-icon="LockIcon"
            />
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              native-type="submit"
              :loading="loading"
              size="large"
              style="width: 100%"
            >
              登录
            </el-button>
          </el-form-item>
        </el-form>

        <div class="form-footer">
          还没有账户？<router-link to="/register">立即注册</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, markRaw } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import apiClient from '../services/api'
import { useAuthStore } from '../stores/auth'

const UserIcon = markRaw(User)
const LockIcon = markRaw(Lock)

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  username: '',
  password: '',
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
  ],
}

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    loading.value = true
    try {
      const response = await apiClient.post('/auth/login', {
        username: form.username,
        password: form.password,
      })

      if (response.code === 1000) {
        const data = response.data
        authStore.setUser(data.user)
        authStore.setToken(data.token)

        ElMessage.success('登录成功')

        setTimeout(() => {
          const redirect = route.query.redirect as string
          if (redirect) {
            router.push(redirect)
          } else if (data.user.role === 'author' || data.user.role === 'admin') {
            router.push('/dashboard')
          } else {
            router.push('/')
          }
        }, 500)
      } else {
        ElMessage.error(response.msg || '登录失败')
      }
    } catch (error: any) {
      ElMessage.error(error.message || '登录失败，请稍后重试')
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.auth-page {
  display: flex;
  min-height: 100vh;
}

/* 左侧品牌区 */
.brand-panel {
  flex: 1;
  background: linear-gradient(135deg, #2563eb 0%, #7c3aed 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  position: relative;
  overflow: hidden;
}

.brand-panel::before {
  content: '';
  position: absolute;
  width: 400px;
  height: 400px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.05);
  top: -100px;
  right: -100px;
}

.brand-panel::after {
  content: '';
  position: absolute;
  width: 300px;
  height: 300px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.05);
  bottom: -80px;
  left: -80px;
}

.brand-content {
  position: relative;
  z-index: 1;
  text-align: center;
  color: white;
}

.brand-logo {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.95;
}

.brand-name {
  font-size: 2.5rem;
  font-weight: 800;
  margin: 0 0 0.5rem 0;
  letter-spacing: -0.5px;
}

.brand-tagline {
  font-size: 1.1rem;
  opacity: 0.85;
  margin: 0 0 2.5rem 0;
}

.brand-features {
  display: flex;
  flex-direction: column;
  gap: 0.9rem;
  text-align: left;
  background: rgba(255, 255, 255, 0.12);
  border-radius: 12px;
  padding: 1.25rem 1.5rem;
  backdrop-filter: blur(4px);
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 0.95rem;
  opacity: 0.9;
}

.feature-item i {
  font-size: 1.1rem;
}

/* 右侧表单区 */
.form-panel {
  width: 480px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 3rem 2.5rem;
  background: #fff;
}

.form-wrapper {
  width: 100%;
  max-width: 360px;
}

.form-header {
  margin-bottom: 2rem;
}

.form-header h2 {
  font-size: 1.75rem;
  font-weight: 700;
  color: #0f172a;
  margin: 0 0 0.5rem 0;
}

.form-header p {
  color: #64748b;
  margin: 0;
  font-size: 0.95rem;
}

.form-footer {
  text-align: center;
  color: #64748b;
  font-size: 0.9rem;
  margin-top: 1.5rem;
}

.form-footer a {
  color: #2563eb;
  font-weight: 500;
  margin-left: 4px;
  text-decoration: none;
}

.form-footer a:hover {
  text-decoration: underline;
}

/* 移动端：单列，渐变 header 在顶部 */
@media (max-width: 768px) {
  .auth-page {
    flex-direction: column;
  }

  .brand-panel {
    flex: none;
    padding: 2rem 1.5rem;
    min-height: 220px;
  }

  .brand-features {
    display: none;
  }

  .brand-name {
    font-size: 2rem;
  }

  .form-panel {
    width: 100%;
    flex: 1;
    padding: 2rem 1.5rem;
  }
}
</style>
