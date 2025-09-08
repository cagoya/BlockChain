import axios from 'axios'

// 创建axios实例
const instance = axios.create({
  baseURL: 'http://localhost:8888',
  timeout: 10000,
})

// 请求拦截器 - 自动添加token
instance.interceptors.request.use(
  (config) => {
    // 检查是否是登录或注册请求，这些请求不需要token
    const isAuthRequest = config.url?.includes('/api/account/login') || 
                         config.url?.includes('/api/account/register')
    
    if (!isAuthRequest) {
      // 从localStorage获取token
      const token = localStorage.getItem('userToken')
      
      // 如果token存在，添加到请求头
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
    }
    
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理401错误
instance.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    // 如果返回401未授权错误，清除token并跳转到登录页
    if (error.response?.status === 401) {
      localStorage.removeItem('userToken')
      localStorage.removeItem('userInfo')
      // 清除axios默认请求头
      delete instance.defaults.headers.common['Authorization']
      
      // 跳转到登录页
      window.location.href = '/login'
    }
    
    return Promise.reject(error)
  }
)

export default instance
