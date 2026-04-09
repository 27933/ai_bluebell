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
            <span>Markdown 编辑器，实时预览</span>
          </div>
          <div class="feature-item">
            <i class="bi bi-people"></i>
            <span>评论互动，分享想法</span>
          </div>
          <div class="feature-item">
            <i class="bi bi-graph-up"></i>
            <span>数据统计，了解影响力</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧表单区 -->
    <div class="form-panel">
      <div class="form-wrapper">
        <div class="form-header">
          <h2>创建账户</h2>
          <p>加入 Bluebell，开始你的创作之旅</p>
        </div>

        <el-form ref="formRef" :model="form" :rules="rules" @submit.prevent="handleSubmit">
          <el-form-item prop="username">
            <el-input
              v-model="form.username"
              placeholder="用户名（3-20 个字符）"
              size="large"
              clearable
              :prefix-icon="UserIcon"
            />
          </el-form-item>

          <el-form-item prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="密码（至少 8 位）"
              size="large"
              show-password
              :prefix-icon="LockIcon"
            />
          </el-form-item>

          <el-form-item prop="confirm_password">
            <el-input
              v-model="form.confirm_password"
              type="password"
              placeholder="确认密码"
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
              注册
            </el-button>
          </el-form-item>
        </el-form>

        <div class="form-footer">
          已有账户？<router-link to="/login">立即登录</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, markRaw } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import apiClient from '../services/api'

const UserIcon = markRaw(User)
const LockIcon = markRaw(Lock)

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  username: '',
  password: '',
  confirm_password: '',
})

const validateConfirmPassword = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度应为 3-20 个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 8, max: 100, message: '密码至少 8 位', trigger: 'blur' },
  ],
  confirm_password: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' },
  ],
}

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (!valid) return

    loading.value = true
    try {
      const response = await apiClient.post('/auth/signup', {
        username: form.username,
        password: form.password,
        re_password: form.confirm_password,
      })

      if (response.code === 1000) {
        ElMessage.success('注册成功，请登录')
        setTimeout(() => {
          router.push('/login')
        }, 1000)
      } else {
        ElMessage.error(response.msg || '注册失败')
      }
    } catch (error: any) {
      ElMessage.error(error.message || '注册失败，请稍后重试')
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
