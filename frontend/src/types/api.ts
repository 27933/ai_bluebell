export interface ApiResponse<T = any> {
  code: number
  msg: string
  data: T
}

export interface LoginResponse {
  user: {
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
  token: {
    access_token: string
    refresh_token: string
    expires_in: number
  }
}

export interface Article {
  id: string
  title: string
  content: string
  summary: string
  view_count: number
  like_count: number
  comment_count: number
  is_featured: boolean
  created_at: string
  word_count?: number
  slug?: string
  status?: string
  allow_comment?: boolean
}

export interface Comment {
  id: string
  content: string
  user_name?: string
  user_id?: string
  article_id: string
  created_at: string
}

export interface ListResponse<T> {
  list: T[]
  total: number
  page: number
  size: number
  pages: number
}
