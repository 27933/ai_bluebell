<template>
  <el-header class="navbar">
    <div class="navbar-container">
      <router-link to="/" class="logo">
        <h1>Bluebell Blog</h1>
      </router-link>

      <div class="nav-menu">
        <router-link to="/" class="nav-link">首页</router-link>
        <router-link to="/articles" class="nav-link">浏览</router-link>

        <!-- 登录前 -->
        <template v-if="!authStore.isLoggedIn">
          <router-link to="/login" class="nav-link">登录</router-link>
          <router-link to="/register" class="nav-link">注册</router-link>
        </template>

        <!-- 登录后 -->
        <template v-else>
          <!-- Reader 菜单 -->
          <template v-if="authStore.user?.role === 'reader'">
            <router-link to="/profile" class="nav-link">个人资料</router-link>
          </template>

          <!-- Author 和 Admin 菜单 -->
          <template v-if="authStore.user?.role === 'author' || authStore.user?.role === 'admin'">
            <router-link to="/write" class="nav-link">写文章</router-link>
            <router-link to="/dashboard" class="nav-link">仪表板</router-link>
            <router-link to="/profile" class="nav-link">个人资料</router-link>
          </template>

          <!-- Admin 菜单 -->
          <template v-if="authStore.user?.role === 'admin'">
            <router-link to="/admin" class="nav-link">管理后台</router-link>
          </template>

          <!-- 用户菜单 -->
          <el-dropdown @command="handleCommand" class="user-menu">
            <span class="user-name">{{ authStore.user?.username }}</span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">我的主页</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </div>
    </div>
  </el-header>
</template>

<script setup lang="ts">
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

function handleCommand(command: string) {
  if (command === 'logout') {
    authStore.logout()
    router.push('/login')
  } else if (command === 'profile') {
    router.push('/profile')
  }
}
</script>

<style scoped>
.navbar {
  background-color: #fff;
  border-bottom: 1px solid #e0e0e0;
  padding: 0 20px;
  height: 60px;
}

.navbar-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.logo {
  text-decoration: none;
}

.logo h1 {
  margin: 0;
  font-size: 24px;
  color: #409eff;
}

.nav-menu {
  display: flex;
  gap: 30px;
  align-items: center;
}

.nav-link {
  text-decoration: none;
  color: #333;
  transition: color 0.3s;
}

.nav-link:hover {
  color: #409eff;
}

.user-menu {
  cursor: pointer;
}

.user-name {
  color: #333;
}
</style>
