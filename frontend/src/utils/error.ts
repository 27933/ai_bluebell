import { ElMessage } from 'element-plus'

export interface ApiError {
  code: number
  msg: string
  data?: any
}

// 错误码映射
const errorCodeMap: Record<number, string> = {
  1000: '成功',
  1001: '请求参数错误',
  1002: '用户名已存在',
  1003: '用户名不存在',
  1004: '用户名或密码错误',
  1005: '服务繁忙，请稍后重试',
  1006: '需要登录',
  1007: '无效的token',
  1013: '你没有权限执行此操作',
  1014: '文章不存在',
}

export function handleApiError(error: any) {
  let message = '发生错误'

  if (error.response?.data?.code) {
    const code = error.response.data.code
    message = error.response.data.msg || errorCodeMap[code] || `错误 ${code}`
  } else if (error.message) {
    message = error.message
  }

  ElMessage.error(message)
  return message
}

export function getErrorMessage(code: number): string {
  return errorCodeMap[code] || `未知错误 (${code})`
}
