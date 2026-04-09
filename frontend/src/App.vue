<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import Navbar from './components/Navbar.vue'
import apiClient from './services/api'
import { useAuthStore } from './stores/auth'

const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()

// token 存在但 user 缺失时（旧会话或首次使用新代码），自动拉取用户信息
// 在组件 onMounted 中执行，确保 Vue 响应式系统完全就绪
onMounted(async () => {
  if (authStore.token?.access_token) {
    try {
      const response = await apiClient.get('/auth/profile') as any
      if (response.code === 1000) {
        authStore.setUser(response.data)
      } else if (response.code === 1015) {
        // 账号已被封禁
        authStore.logout()
        ElMessage.error('账号已被封禁，请联系管理员')
        router.push('/login')
      }
    } catch {
      // token 已失效，保持未登录状态
    }
  }
})
</script>

<template>
  <div class="app">
    <Navbar v-if="!route.path.startsWith('/admin')" />
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body {
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
    sans-serif;
  color: #334155;
  background-color: #f8fafc;
}

#app {
  height: 100%;
}

.app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.main-content {
  flex: 1;
}

a {
  color: inherit;
  text-decoration: none;
}

button {
  font-family: inherit;
}

input,
select,
textarea {
  font-family: inherit;
}

/* 保持Bootstrap container默认，让页面自适应 */

/* ===== 全局色系变量（原型一致） ===== */
:root {
  --primary-color: #2563eb;
  --secondary-color: #64748b;
  --danger-color: #ef4444;

  /* 覆盖 Element Plus 默认主色，与项目色系统一 */
  --el-color-primary: #2563eb;
  --el-color-primary-light-3: #5b8af0;
  --el-color-primary-light-5: #93b4f5;
  --el-color-primary-light-7: #c9d9fa;
  --el-color-primary-light-8: #dce8fc;
  --el-color-primary-light-9: #eef3fe;
  --el-color-primary-dark-2: #1d4ed8;
}
</style>

