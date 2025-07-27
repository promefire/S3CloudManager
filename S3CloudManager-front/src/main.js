import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
/* eslint-disable */
import './assets/css/materialize.min.css'
/* eslint-disable */
import './assets/js/materialize.min.js'

const app = createApp(App)

app.use(router)

app.mount('#app')
