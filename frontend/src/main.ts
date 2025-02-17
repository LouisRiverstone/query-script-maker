import {createApp} from 'vue'
import App from './App.vue'
import './style.css';
import routes from './routes'

createApp(App).use(routes).mount('#app')