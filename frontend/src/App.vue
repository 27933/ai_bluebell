<script setup lang="ts">
import { onMounted } from 'vue'
import Navbar from './components/Navbar.vue'
import apiClient from './services/api'
import { useAuthStore } from './stores/auth'

const authStore = useAuthStore()

// token 存在但 user 缺失时（旧会话或首次使用新代码），自动拉取用户信息
// 在组件 onMounted 中执行，确保 Vue 响应式系统完全就绪
onMounted(async () => {
  if (authStore.token?.access_token && !authStore.user) {
    try {
      const response = await apiClient.get('/auth/profile') as any
      if (response.code === 1000) {
        authStore.setUser(response.data)
      }
    } catch {
      // token 已失效，保持未登录状态
    }
  }
})
</script>

<template>
  <div class="app">
    <Navbar />
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
}
</style>

