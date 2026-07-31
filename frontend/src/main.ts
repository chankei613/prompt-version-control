import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHashHistory } from 'vue-router'
import 'virtual:uno.css'
import '@unocss/reset/tailwind.css'
import './assets/globals.css'
import App from './App.vue'
import PromptsView from './pages/PromptsView.vue'
import PromptDetailView from './pages/PromptDetailView.vue'
import HelpView from './pages/HelpView.vue'
import SettingsView from './pages/SettingsView.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', redirect: '/prompts' },
    { path: '/prompts', component: PromptsView },
    { path: '/prompts/:id', component: PromptDetailView },
    { path: '/help', component: HelpView },
    { path: '/settings', component: SettingsView },
  ],
})

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
