<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useTenantStore } from '@/stores/tenant'
import AppHeader from '@/components/AppHeader.vue'
import GuestHeader from '@/components/GuestHeader.vue'
import Toast from '@/components/Toast.vue'
import { RouterView } from 'vue-router'

const route = useRoute()
const authStore = useAuthStore()
const tenantStore = useTenantStore()

// Initialize tenant on app mount
onMounted(async () => {
  // First try to load from storage (for faster initial render)
  tenantStore.initializeFromStorage()
  
  // Then detect and fetch fresh tenant info
  await tenantStore.detectTenant()
})

// Apply tenant branding colors as CSS variables
watch(
  () => tenantStore.primaryColor,
  (color) => {
    if (color) {
      document.documentElement.style.setProperty('--primary-color', color)
      document.documentElement.style.setProperty('--primary-hover', adjustBrightness(color, 25))
      document.documentElement.style.setProperty('--primary-active', adjustBrightness(color, -15))
    }
  },
  { immediate: true }
)

// Helper to adjust color brightness
function adjustBrightness(hex: string, percent: number): string {
  hex = hex.replace(/^#/, '')
  let r = parseInt(hex.substring(0, 2), 16)
  let g = parseInt(hex.substring(2, 4), 16)
  let b = parseInt(hex.substring(4, 6), 16)
  r = Math.min(255, Math.max(0, r + (r * percent) / 100))
  g = Math.min(255, Math.max(0, g + (g * percent) / 100))
  b = Math.min(255, Math.max(0, b + (b * percent) / 100))
  const toHex = (n: number) => Math.round(n).toString(16).padStart(2, '0')
  return `#${toHex(r)}${toHex(g)}${toHex(b)}`
}

const showHeader = computed(() => {
  // Show header for all authenticated users, except on login/register pages
  return authStore.isAuthenticated && route.name !== 'login' && route.name !== 'register'
})

const showGuestHeader = computed(() => {
  return !authStore.isAuthenticated && route.name !== 'login' && route.name !== 'register'
})
</script>

<template>
  <div id="app">
    <AppHeader v-if="showHeader" />
    <GuestHeader v-if="showGuestHeader" />
    <RouterView />
    <Toast />
  </div>
</template>

<style>
:root {
  --primary-color: #6929C4;
  --primary-hover: #8A3FFC;
  --primary-active: #4F2196;
  --text-primary: #161616;
  --text-secondary: #525252;
  --background: #ffffff;
  --surface: #f4f4f4;
  --border-color: #e0e0e0;
  --error-color: #da1e28;
  --success-color: #24a148;
}

#app {
  min-height: 100vh;
  background-color: var(--background);
  color: var(--text-primary);
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: 'IBM Plex Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI',
    'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans',
    'Helvetica Neue', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* PrimeVue Theme Customization */
.p-component {
  font-family: 'IBM Plex Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI',
    'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans',
    'Helvetica Neue', sans-serif;
}
</style>

