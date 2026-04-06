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
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  function setToken(newToken: Token) {
    token.value = newToken
    localStorage.setItem('token', JSON.stringify(newToken))
  }

  function logout() {
    user.value = null
    token.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  function loadTokenFromStorage() {
    const storedToken = localStorage.getItem('token')
    if (storedToken) {
      try {
        token.value = JSON.parse(storedToken)
      } catch (e) {
        console.error('Failed to parse token from storage', e)
      }
    }
    const storedUser = localStorage.getItem('user')
    if (storedUser) {
      try {
        user.value = JSON.parse(storedUser)
      } catch (e) {
        console.error('Failed to parse user from storage', e)
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
