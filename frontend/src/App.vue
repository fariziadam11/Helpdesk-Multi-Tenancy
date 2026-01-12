<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import AppHeader from '@/components/AppHeader.vue'
import GuestHeader from '@/components/GuestHeader.vue'
import Toast from '@/components/Toast.vue'
import { RouterView } from 'vue-router'

const route = useRoute()
const authStore = useAuthStore()

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
