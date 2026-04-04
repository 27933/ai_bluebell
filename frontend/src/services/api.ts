import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig } from 'axios'
import { AxiosError } from 'axios'
import { useAuthStore } from '../stores/auth'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8084/api/v1'

// 创建 axios 实例
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    if (authStore.token?.access_token) {
      config.headers.Authorization = `Bearer ${authStore.token.access_token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
apiClient.interceptors.response.use(
  (response) => {
    return response.data
  },
  async (error: AxiosError) => {
    const authStore = useAuthStore()
    const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean }

    // 处理 401 错误（Token 过期）
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      try {
        const refreshToken = authStore.token?.refresh_token
        if (refreshToken) {
          const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
            refresh_token: refreshToken,
          })
          const newToken = response.data.data.token
          authStore.setToken(newToken)
          // 使用新 token 重试原请求
          if (originalRequest.headers) {
            originalRequest.headers.Authorization = `Bearer ${newToken.access_token}`
          }
          return apiClient(originalRequest)
        }
      } catch (refreshError) {
        authStore.logout()
        window.location.href = '/login'
        return Promise.reject(refreshError)
      }
    }

    // 处理其他错误
    const errorMessage = (error.response?.data as any)?.msg || error.message || 'API 请求失败'
    const apiError = new Error(errorMessage)
    return Promise.reject(apiError)
  }
)

export default apiClient
