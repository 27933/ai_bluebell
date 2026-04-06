<template>
  <!-- 导航栏 -->
  <nav class="navbar navbar-expand-lg">
    <div class="container">
      <router-link to="/" class="navbar-brand">
        <i class="bi bi-journal-text"></i> Bluebell
      </router-link>

      <div class="navbar-nav ms-auto">
        <router-link to="/" class="nav-item">
          <a class="nav-link">首页</a>
        </router-link>
        <router-link to="/tags" class="nav-item">
          <a class="nav-link">标签</a>
        </router-link>

        <!-- 未登录状态 -->
        <template v-if="!authStore.isLoggedIn">
          <router-link to="/login" class="nav-item auth-link">
            <a class="nav-link" style="color: var(--primary-color); font-weight: 600;">登录</a>
          </router-link>
          <router-link to="/register" class="nav-item">
            <a class="nav-link" style="color: var(--primary-color);">注册</a>
          </router-link>
        </template>

        <!-- 已登录状态 - 用户菜单 -->
        <div v-else class="nav-item user-menu">
          <button class="user-menu-toggle" @click="toggleUserMenu">
            <div class="user-avatar">{{ getInitial(authStore.user?.username) }}</div>
            <span>{{ authStore.user?.username }}</span>
            <i class="bi bi-chevron-down" style="font-size: 0.8rem;"></i>
          </button>

          <div v-show="userMenuOpen" class="user-dropdown" @mouseenter="cancelClose" @mouseleave="delayClose">
            <!-- 作者菜单项 -->
            <router-link
              v-if="authStore.user?.role === 'author' || authStore.user?.role === 'admin'"
              to="/write"
              class="dropdown-link"
              @click="userMenuOpen = false"
            >
              <i class="bi bi-pencil"></i> 写文章
            </router-link>
            <router-link
              v-if="authStore.user?.role === 'author' || authStore.user?.role === 'admin'"
              to="/dashboard"
              class="dropdown-link"
              @click="userMenuOpen = false"
            >
              <i class="bi bi-bar-chart"></i> 仪表板
            </router-link>

            <!-- 普通菜单项 -->
            <router-link to="/profile" class="dropdown-link" @click="userMenuOpen = false">
              <i class="bi bi-user-circle"></i> 个人资料
            </router-link>

            <!-- 管理员菜单项 -->
            <router-link
              v-if="authStore.user?.role === 'admin'"
              to="/admin"
              class="dropdown-link"
              @click="userMenuOpen = false"
            >
              <i class="bi bi-gear"></i> 管理后台
            </router-link>

            <div class="divider"></div>

            <button class="dropdown-link logout" @click="handleLogout">
              <i class="bi bi-box-arrow-right"></i> 退出登录
            </button>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const router = useRouter()
const userMenuOpen = ref(false)
let closeTimer: ReturnType<typeof setTimeout> | null = null

function getInitial(name?: string): string {
  if (!name) return '？'
  return name.charAt(0).toUpperCase()
}

function toggleUserMenu() {
  cancelClose()
  userMenuOpen.value = !userMenuOpen.value
  // 点击打开后，300ms 后开始监听外部点击
  if (userMenuOpen.value) {
    setTimeout(() => {
      document.addEventListener('click', handleOutsideClick)
    }, 50)
  }
}

function delayClose() {
  closeTimer = setTimeout(() => {
    userMenuOpen.value = false
  }, 300)
}

function cancelClose() {
  if (closeTimer) {
    clearTimeout(closeTimer)
    closeTimer = null
  }
}

function handleOutsideClick(e: MouseEvent) {
  const menu = document.querySelector('.user-menu')
  if (menu && !menu.contains(e.target as Node)) {
    userMenuOpen.value = false
    document.removeEventListener('click', handleOutsideClick)
  }
}

async function handleLogout() {
  authStore.logout()
  userMenuOpen.value = false
  document.removeEventListener('click', handleOutsideClick)
  router.push('/login')
}
</script>

<style scoped>
/* ===== 导航栏 ===== */
.navbar {
  background-color: white;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  padding: 1rem 0;
  position: sticky;
  top: 0;
  z-index: 100;
}

/* 使用 Bootstrap 默认 .container 响应式宽度 */
.container {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.navbar-brand {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--primary-color) !important;
  text-decoration: none;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.navbar-nav {
  display: flex;
  align-items: center;
  gap: 0;
  list-style: none;
}

.nav-item {
  position: relative;
  display: flex;
  align-items: center;
}

.nav-link {
  color: var(--secondary-color) !important;
  font-weight: 500;
  margin: 0 0.5rem;
  transition: all 0.3s;
  text-decoration: none;
  display: block;
  padding: 0.5rem 0;
}

.nav-link:hover {
  color: var(--primary-color) !important;
}

.nav-link.active {
  color: var(--primary-color) !important;
}

.ms-auto {
  margin-left: auto;
}

/* ===== 用户菜单 ===== */
.user-menu {
  position: relative;
  display: inline-block;
  margin-left: 1rem;
}

.user-menu-toggle {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.5rem 0.75rem;
  background: none;
  border: none;
  cursor: pointer;
  color: var(--secondary-color);
  font-weight: 500;
  transition: all 0.3s;
  border-radius: 6px;
}

.user-menu-toggle:hover {
  background-color: #f1f5f9;
  color: var(--primary-color);
}

.user-avatar {
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
}

.user-dropdown {
  display: block;
  position: absolute;
  right: 0;
  top: 100%;
  margin-top: 0.5rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  min-width: 200px;
  z-index: 1000;
  animation: slideDown 0.2s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.dropdown-link {
  display: block;
  width: 100%;
  padding: 0.75rem 1rem;
  border: none;
  background: none;
  text-align: left;
  cursor: pointer;
  color: #334155;
  font-size: 0.95rem;
  transition: all 0.3s;
  text-decoration: none;
  font-weight: normal;
}

.dropdown-link:hover {
  background-color: #f8fafc;
  color: var(--primary-color);
}

.dropdown-link:first-child {
  border-radius: 8px 8px 0 0;
}

.divider {
  height: 1px;
  background-color: #e2e8f0;
  margin: 0.5rem 0;
}

.logout {
  color: var(--danger-color);
}

.logout:hover {
  background-color: #fef2f2;
}

.logout:last-child {
  border-radius: 0 0 8px 8px;
}
</style>
