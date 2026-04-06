import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig } from 'axios'
import { AxiosError } from 'axios'
import { ElMessage } from 'element-plus'
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
  async (response) => {
    const data = response.data

    // 后端所有错误均返回 HTTP 200，code 1007 表示 token 无效，需在此处刷新
    if (data.code === 1007) {
      const config = response.config as AxiosRequestConfig & { _retry?: boolean }
      if (!config._retry) {
        config._retry = true
        const authStore = useAuthStore()
        const refreshToken = authStore.token?.refresh_token

        if (refreshToken) {
          try {
            const refreshResp = await axios.post(`${API_BASE_URL}/auth/refresh`, {
              refresh_token: refreshToken,
            })
            // 后端 refresh 响应：{ code: 1000, data: { access_token, refresh_token, expires_in } }
            const newToken = refreshResp.data.data
            authStore.setToken(newToken)
            if (config.headers) {
              config.headers['Authorization'] = `Bearer ${newToken.access_token}`
            }
            return apiClient(config)
          } catch {
            // 10.2 Token 失效提示
            authStore.logout()
            ElMessage.warning('登录已过期，请重新登录')
            setTimeout(() => { window.location.href = '/login' }, 1500)
            return Promise.reject(new Error('登录已过期，请重新登录'))
          }
        } else {
          // 10.2 无 refresh token
          authStore.logout()
          ElMessage.warning('请先登录')
          setTimeout(() => { window.location.href = '/login' }, 1500)
          return Promise.reject(new Error('请先登录'))
        }
      }
    }

    return data
  },
  async (error: AxiosError) => {
    // 10.1 全局错误处理：网络层异常（后端目前正常情况均返回 HTTP 200）
    const authStore = useAuthStore()
    const originalRequest = error.config as AxiosRequestConfig & { _retry?: boolean }

    // HTTP 401（备用，以防后端将来返回标准 HTTP 状态码）
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      try {
        const refreshToken = authStore.token?.refresh_token
        if (refreshToken) {
          const resp = await axios.post(`${API_BASE_URL}/auth/refresh`, {
            refresh_token: refreshToken,
          })
          const newToken = resp.data.data
          authStore.setToken(newToken)
          if (originalRequest.headers) {
            originalRequest.headers['Authorization'] = `Bearer ${newToken.access_token}`
          }
          return apiClient(originalRequest)
        }
      } catch (refreshError) {
        authStore.logout()
        ElMessage.warning('登录已过期，请重新登录')
        setTimeout(() => { window.location.href = '/login' }, 1500)
        return Promise.reject(refreshError)
      }
    }

    // 请求超时
    if (error.code === 'ECONNABORTED') {
      ElMessage.error('请求超时，请稍后重试')
      return Promise.reject(error)
    }

    // 无网络 / 无法连接服务器
    if (!error.response) {
      ElMessage.error('网络连接失败，请检查网络')
      return Promise.reject(error)
    }

    // 服务器 5xx 错误
    if (error.response.status >= 500) {
      ElMessage.error('服务器异常，请稍后重试')
      return Promise.reject(error)
    }

    const errorMessage = (error.response?.data as any)?.msg || error.message || '请求失败'
    return Promise.reject(new Error(errorMessage))
  }
)

export default apiClient
