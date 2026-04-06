import apiClient from '../services/api'

// ===== 类型定义 =====

export interface AdminUser {
  id: string
  username: string
  nickname?: string
  email?: string
  role: string
  status: string
  created_at: string
}

export interface AdminUserListResp {
  list: AdminUser[]
  total: number
  page: number
  size: number
  pages: number
}

export interface AdminArticle {
  id: string
  title: string
  summary: string
  word_count: number
  author_id: string
  author_username: string
  author_nickname?: string
  status: string
  is_featured: boolean
  like_count: number
  comment_count: number
  view_count: number
  slug: string
  created_at: string
  updated_at: string
}

export interface AdminArticleListResp {
  list: AdminArticle[]
  total: number
  page: number
  size: number
  pages: number
}

export interface SystemOverview {
  user_count: number
  article_count: number
  comment_count: number
  today_new_user_count: number
  today_new_article_count: number
}

export interface DailyStatItem {
  date: string
  new_user_count: number
  new_article_count: number
  new_comment_count: number
}

// ===== 用户管理 =====

export function getUserList(params: {
  role?: string
  status?: string
  page?: number
  size?: number
}) {
  return apiClient.get<any, { code: number; data: AdminUserListResp }>('/admin/users', { params })
}

export function getUserDetail(id: string) {
  return apiClient.get<any, { code: number; data: AdminUser }>(`/admin/users/${id}`)
}

export function updateUserRole(id: string, role: string) {
  return apiClient.patch<any, { code: number; data: null }>(`/admin/users/${id}/role`, { role })
}

export function updateUserStatus(id: string, status: string) {
  return apiClient.patch<any, { code: number; data: null }>(`/admin/users/${id}/status`, { status })
}

export function batchUpdateUserStatus(userIds: string[], status: string) {
  return apiClient.patch<any, { code: number; data: { updated_count: number } }>(
    '/admin/users/batch/status',
    { user_ids: userIds.map(Number), status }
  )
}

// ===== 文章管理 =====

export function getAdminArticleList(params: {
  status?: string
  keyword?: string
  author_name?: string
  page?: number
  size?: number
}) {
  return apiClient.get<any, { code: number; data: AdminArticleListResp }>('/admin/articles', { params })
}

export function setArticleFeatured(id: string, isFeatured: boolean) {
  return apiClient.patch<any, { code: number; data: null }>(`/admin/articles/${id}/featured`, {
    is_featured: isFeatured,
  })
}

// ===== 统计 =====

export function getSystemOverview() {
  return apiClient.get<any, { code: number; data: SystemOverview }>('/admin/stats/overview')
}

export function getSystemDailyStats(params?: { start_date?: string; end_date?: string }) {
  return apiClient.get<any, { code: number; data: DailyStatItem[] }>('/admin/stats/daily', { params })
}
