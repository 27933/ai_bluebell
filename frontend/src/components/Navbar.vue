<template>
  <nav class="navbar">
    <div class="container">
      <router-link to="/" class="navbar-brand">
        <i class="bi bi-journal-text"></i> Bluebell
      </router-link>

      <!-- 桌面端导航 -->
      <div class="navbar-nav ms-auto desktop-nav">
        <router-link to="/" class="nav-item">
          <a class="nav-link">首页</a>
        </router-link>
        <router-link to="/tags" class="nav-item">
          <a class="nav-link">标签</a>
        </router-link>

        <template v-if="!authStore.isLoggedIn">
          <router-link to="/login" class="nav-item auth-link">
            <a class="nav-link" style="color: var(--primary-color); font-weight: 600;">登录</a>
          </router-link>
          <router-link to="/register" class="nav-item">
            <a class="nav-link" style="color: var(--primary-color);">注册</a>
          </router-link>
        </template>

        <div v-else class="nav-item user-menu">
          <button class="user-menu-toggle" @click="toggleUserMenu">
            <div class="user-avatar">{{ getInitial(authStore.user?.username) }}</div>
            <span>{{ authStore.user?.username }}</span>
            <i class="bi bi-chevron-down" style="font-size: 0.8rem;"></i>
          </button>
          <div v-show="userMenuOpen" class="user-dropdown" @mouseenter="cancelClose" @mouseleave="delayClose">
            <router-link v-if="authStore.user?.role === 'author' || authStore.user?.role === 'admin'" to="/write" class="dropdown-link" @click="userMenuOpen = false">
              <i class="bi bi-pencil"></i> 写文章
            </router-link>
            <router-link v-if="authStore.user?.role === 'author' || authStore.user?.role === 'admin'" to="/dashboard" class="dropdown-link" @click="userMenuOpen = false">
              <i class="bi bi-bar-chart"></i> 仪表板
            </router-link>
            <router-link to="/profile" class="dropdown-link" @click="userMenuOpen = false">
              <i class="bi bi-user-circle"></i> 个人资料
            </router-link>
            <router-link v-if="authStore.user?.role === 'admin'" to="/admin" class="dropdown-link" @click="userMenuOpen = false">
              <i class="bi bi-gear"></i> 管理后台
            </router-link>
            <div class="divider"></div>
            <button class="dropdown-link logout" @click="handleLogout">
              <i class="bi bi-box-arrow-right"></i> 退出登录
            </button>
          </div>
        </div>
      </div>

      <!-- 移动端汉堡按钮 -->
      <button class="hamburger" @click="toggleMobileMenu" :class="{ open: mobileMenuOpen }">
        <span></span>
        <span></span>
        <span></span>
      </button>
    </div>

    <!-- 移动端展开菜单 -->
    <div class="mobile-menu" :class="{ open: mobileMenuOpen }">
      <div class="mobile-nav-links">
        <router-link to="/" class="mobile-nav-link" @click="closeMobileMenu">
          <i class="bi bi-house"></i> 首页
        </router-link>
        <router-link to="/tags" class="mobile-nav-link" @click="closeMobileMenu">
          <i class="bi bi-tags"></i> 标签
        </router-link>

        <template v-if="!authStore.isLoggedIn">
          <router-link to="/login" class="mobile-nav-link primary" @click="closeMobileMenu">
            <i class="bi bi-box-arrow-in-right"></i> 登录
          </router-link>
          <router-link to="/register" class="mobile-nav-link" @click="closeMobileMenu">
            <i class="bi bi-person-plus"></i> 注册
          </router-link>
        </template>

        <template v-else>
          <div class="mobile-user-info">
            <div class="user-avatar">{{ getInitial(authStore.user?.username) }}</div>
            <span>{{ authStore.user?.username }}</span>
          </div>
          <router-link v-if="authStore.user?.role === 'author' || authStore.user?.role === 'admin'" to="/write" class="mobile-nav-link" @click="closeMobileMenu">
            <i class="bi bi-pencil"></i> 写文章
          </router-link>
          <router-link v-if="authStore.user?.role === 'author' || authStore.user?.role === 'admin'" to="/dashboard" class="mobile-nav-link" @click="closeMobileMenu">
            <i class="bi bi-bar-chart"></i> 仪表板
          </router-link>
          <router-link to="/profile" class="mobile-nav-link" @click="closeMobileMenu">
            <i class="bi bi-person-circle"></i> 个人资料
          </router-link>
          <router-link v-if="authStore.user?.role === 'admin'" to="/admin" class="mobile-nav-link" @click="closeMobileMenu">
            <i class="bi bi-gear"></i> 管理后台
          </router-link>
          <div class="mobile-divider"></div>
          <button class="mobile-nav-link logout" @click="handleLogout">
            <i class="bi bi-box-arrow-right"></i> 退出登录
          </button>
        </template>
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
const mobileMenuOpen = ref(false)
let closeTimer: ReturnType<typeof setTimeout> | null = null

function toggleMobileMenu() {
  mobileMenuOpen.value = !mobileMenuOpen.value
  document.body.style.overflow = mobileMenuOpen.value ? 'hidden' : ''
}

function closeMobileMenu() {
  mobileMenuOpen.value = false
  document.body.style.overflow = ''
}

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
  mobileMenuOpen.value = false
  document.body.style.overflow = ''
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

/* ===== 汉堡按钮（移动端） ===== */
.hamburger {
  display: none;
  flex-direction: column;
  justify-content: space-between;
  width: 24px;
  height: 18px;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
}

.hamburger span {
  display: block;
  height: 2px;
  background: var(--secondary-color);
  border-radius: 2px;
  transition: all 0.3s;
  transform-origin: center;
}

.hamburger.open span:nth-child(1) {
  transform: translateY(8px) rotate(45deg);
}
.hamburger.open span:nth-child(2) {
  opacity: 0;
}
.hamburger.open span:nth-child(3) {
  transform: translateY(-8px) rotate(-45deg);
}

/* ===== 移动端菜单 ===== */
.mobile-menu {
  display: none;
  background: white;
  border-top: 1px solid #e2e8f0;
  max-height: 0;
  overflow: hidden;
  transition: max-height 0.3s ease;
}

.mobile-menu.open {
  max-height: 100vh;
}

.mobile-nav-links {
  padding: 0.5rem 0 1rem;
}

.mobile-nav-link {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  width: 100%;
  padding: 0.875rem 1.5rem;
  color: #334155;
  font-size: 1rem;
  font-weight: 500;
  text-decoration: none;
  background: none;
  border: none;
  cursor: pointer;
  text-align: left;
  transition: background 0.2s;
}

.mobile-nav-link:hover,
.mobile-nav-link:active {
  background: #f8fafc;
  color: var(--primary-color);
}

.mobile-nav-link.primary {
  color: var(--primary-color);
  font-weight: 600;
}

.mobile-nav-link.logout {
  color: var(--danger-color);
}

.mobile-user-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1.5rem;
  border-bottom: 1px solid #e2e8f0;
  margin-bottom: 0.25rem;
  font-weight: 600;
  color: #1e293b;
}

.mobile-divider {
  height: 1px;
  background: #e2e8f0;
  margin: 0.5rem 0;
}

/* ===== 响应式断点 ===== */
@media (max-width: 768px) {
  .desktop-nav {
    display: none !important;
  }

  .hamburger {
    display: flex;
  }

  .mobile-menu {
    display: block;
  }
}
</style>
