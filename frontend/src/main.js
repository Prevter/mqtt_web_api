import { createApp } from 'vue'
import { getRouter } from './router'
import './style.css'
import App from './App.vue'

const app = createApp(App)
app.use(getRouter())
app.mount('#app')