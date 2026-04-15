import { createApp } from 'vue'
import App from './app/App.vue'
import router from './app/router.js'
import './app/styles/main.css'

createApp(App).use(router).mount('#app')
