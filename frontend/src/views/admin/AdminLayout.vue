<template>
  <div class="admin-layout">
    <!-- 侧边栏 -->
    <aside class="admin-sidebar">
      <div class="sidebar-header">
        <router-link to="/" class="sidebar-brand">
          <i class="bi bi-journal-text"></i>
          <span>Bluebell</span>
        </router-link>
        <div class="sidebar-badge">管理后台</div>
      </div>

      <nav class="sidebar-nav">
        <router-link to="/admin/stats" class="sidebar-link" active-class="active">
          <i class="bi bi-bar-chart-line"></i>
          <span>系统统计</span>
        </router-link>
        <router-link to="/admin/users" class="sidebar-link" active-class="active">
          <i class="bi bi-people"></i>
          <span>用户管理</span>
        </router-link>
        <router-link to="/admin/articles" class="sidebar-link" active-class="active">
          <i class="bi bi-file-earmark-text"></i>
          <span>文章管理</span>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <router-link to="/" class="sidebar-link">
          <i class="bi bi-arrow-left-circle"></i>
          <span>返回前台</span>
        </router-link>
      </div>
    </aside>

    <!-- 主区域 -->
    <div class="admin-main">
      <!-- 顶部栏 -->
      <header class="admin-topbar">
        <div class="topbar-title">{{ pageTitle }}</div>
        <div class="topbar-user">
          <div class="user-avatar-sm">{{ getInitial(authStore.user?.username) }}</div>
          <span>{{ authStore.user?.username }}</span>
        </div>
      </header>

      <!-- 内容区 -->
      <main class="admin-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()
const route = useRoute()

const pageTitle = computed(() => {
  const map: Record<string, string> = {
    '/admin/stats': '系统统计',
    '/admin/users': '用户管理',
    '/admin/articles': '文章管理',
  }
  return map[route.path] || '管理后台'
})

function getInitial(name?: string): string {
  if (!name) return '?'
  return name.charAt(0).toUpperCase()
}
</script>

<style scoped>
.admin-layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
  background: #f1f5f9;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

/* ===== 侧边栏 ===== */
.admin-sidebar {
  width: 220px;
  flex-shrink: 0;
  background: #1e293b;
  color: #cbd5e1;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.sidebar-header {
  padding: 1.25rem 1rem 1rem;
  border-bottom: 1px solid #334155;
}

.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: white;
  text-decoration: none;
  font-size: 1.1rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.sidebar-brand i {
  color: #818cf8;
}

.sidebar-badge {
  display: inline-block;
  background: #4f46e5;
  color: white;
  font-size: 0.7rem;
  padding: 0.15rem 0.5rem;
  border-radius: 4px;
  font-weight: 600;
  letter-spacing: 0.05em;
}

.sidebar-nav {
  flex: 1;
  padding: 0.75rem 0;
}

.sidebar-link {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1.25rem;
  color: #94a3b8;
  text-decoration: none;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.2s;
  border-left: 3px solid transparent;
}

.sidebar-link:hover {
  background: #334155;
  color: #e2e8f0;
}

.sidebar-link.active {
  background: #334155;
  color: white;
  border-left-color: #818cf8;
}

.sidebar-link i {
  font-size: 1rem;
  width: 18px;
  text-align: center;
}

.sidebar-footer {
  padding: 0.75rem 0;
  border-top: 1px solid #334155;
}

/* ===== 主区域 ===== */
.admin-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.admin-topbar {
  background: white;
  border-bottom: 1px solid #e2e8f0;
  padding: 0 1.5rem;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
}

.topbar-title {
  font-size: 1rem;
  font-weight: 600;
  color: #1e293b;
}

.topbar-user {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #475569;
  font-size: 0.9rem;
}

.user-avatar-sm {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: linear-gradient(135deg, #4f46e5, #7c3aed);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 0.8rem;
}

.admin-content {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem;
}
</style>
