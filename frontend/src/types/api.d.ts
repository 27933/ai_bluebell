declare module 'axios' {
  export interface AxiosInstance {
    get<T = any>(url: string, config?: any): Promise<T>
    post<T = any>(url: string, data?: any, config?: any): Promise<T>
    put<T = any>(url: string, data?: any, config?: any): Promise<T>
    delete<T = any>(url: string, config?: any): Promise<T>
    patch<T = any>(url: string, data?: any, config?: any): Promise<T>
  }
}

export interface ApiResponse<T = any> {
  code: number
  msg: string
  data: T
}
