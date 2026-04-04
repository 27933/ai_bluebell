import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface User {
  id: string
  username: string
  email?: string
  nickname?: string
  bio?: string
  role: 'visitor' | 'reader' | 'author' | 'admin'
  status: string
  total_words?: number
  total_likes?: number
  created_at?: string
}

export interface Token {
  access_token: string
  refresh_token: string
  expires_in: number
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<Token | null>(null)
  const isLoggedIn = computed(() => !!token.value?.access_token)

  function setUser(newUser: User) {
    user.value = newUser
  }

  function setToken(newToken: Token) {
    token.value = newToken
    // 保存到 localStorage
    localStorage.setItem('token', JSON.stringify(newToken))
  }

  function logout() {
    user.value = null
    token.value = null
    localStorage.removeItem('token')
  }

  function loadTokenFromStorage() {
    const stored = localStorage.getItem('token')
    if (stored) {
      try {
        token.value = JSON.parse(stored)
      } catch (e) {
        console.error('Failed to parse token from storage', e)
      }
    }
  }

  return {
    user,
    token,
    isLoggedIn,
    setUser,
    setToken,
    logout,
    loadTokenFromStorage,
  }
})
