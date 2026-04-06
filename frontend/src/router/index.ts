import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue'),
  },
  {
    path: '/search',
    name: 'Search',
    component: () => import('../views/Search.vue'),
  },
  {
    path: '/tags',
    name: 'Tags',
    component: () => import('../views/Tags.vue'),
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue'),
  },
  {
    path: '/article/:id',
    name: 'ArticleDetail',
    component: () => import('../views/ArticleDetail.vue'),
  },
  {
    path: '/write',
    name: 'WriteArticle',
    component: () => import('../views/WriteArticle.vue'),
    meta: { requiresAuth: true, requiredRole: 'author' },
  },
  {
    path: '/write/:id',
    name: 'EditArticle',
    component: () => import('../views/WriteArticle.vue'),
    meta: { requiresAuth: true, requiredRole: 'author' },
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../views/Dashboard.vue'),
    meta: { requiresAuth: true, requiredRole: 'author' },
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('../views/Profile.vue'),
    meta: { requiresAuth: true },
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()
  const { requiresAuth, requiredRole } = to.meta as {
    requiresAuth?: boolean
    requiredRole?: string
  }

  // 已登录用户访问登录/注册页 → 跳回首页
  if ((to.name === 'Login' || to.name === 'Register') && authStore.isLoggedIn) {
    return next({ name: 'Home' })
  }

  // 需要登录
  if (requiresAuth && !authStore.isLoggedIn) {
    return next({ name: 'Login', query: { redirect: to.fullPath } })
  }

  // 需要特定角色（author / admin 均可）
  if (requiredRole === 'author') {
    const role = authStore.user?.role
    if (role !== 'author' && role !== 'admin') {
      return next({ name: 'Home' })
    }
  }

  next()
})

export default router
