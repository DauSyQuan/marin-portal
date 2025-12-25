import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

// Import Bootstrap & Custom CSS
import 'bootstrap/dist/css/bootstrap.min.css'
import 'bootstrap/dist/js/bootstrap.bundle.min.js'
import 'leaflet/dist/leaflet.css'
import './assets/main.css' // Copy phần CSS root, body vào đây

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')