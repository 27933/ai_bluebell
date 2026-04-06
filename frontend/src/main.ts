import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'bootstrap/dist/css/bootstrap.min.css'
import 'bootstrap-icons/font/bootstrap-icons.css'
import './style.css'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import { useAuthStore } from './stores/auth'

const app = createApp(App)

app.use(ElementPlus)
app.use(createPinia())
app.use(router)

// 初始化认证状态（只恢复 token 和已缓存的 user）
const authStore = useAuthStore()
authStore.loadTokenFromStorage()

app.mount('#app')
