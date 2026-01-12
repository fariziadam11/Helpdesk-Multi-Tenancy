import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import vueQueryPlugin from './plugins/vueQuery'
import piniaPlugin from './plugins/pinia'
import { useAuthStore } from './stores/auth'
import i18n from './i18n'

// PrimeVue
import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'
import 'primeicons/primeicons.css'

const app = createApp(App)

// PrimeVue Configuration
app.use(PrimeVue, {
  theme: {
    preset: Aura,
    options: {
      darkModeSelector: false,
      cssLayer: false,
    },
  },
})

app.use(i18n)

app.use(piniaPlugin)
app.use(vueQueryPlugin)

const authStore = useAuthStore()
authStore.initializeAuth()

app.use(router)

app.mount('#app')
