<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import Select from 'primevue/select'
import { setCookie } from '@/utils/cookies'
import { tenantApi } from '@/api/tenant'
import { useAuthStore } from '@/stores/auth'
import { useTenantStore } from '@/stores/tenant'
import { useToast } from '@/composables/useToast'
import type { TenantPublicInfo } from '@/api/tenant'

const authStore = useAuthStore()
const tenantStore = useTenantStore()
const toast = useToast()

const tenants = ref<TenantPublicInfo[]>([])
const selectedTenant = ref<TenantPublicInfo | null>(null)
const isLoading = ref(false)

// Only show switcher for Super Admin
const isSuperAdmin = computed(() => {
  if (!authStore.user) return false
  return authStore.user.email.includes('admin') || tenantStore.tenantSlug === 'default'
})

const loadTenants = async () => {
  if (!isSuperAdmin.value) return
  
  try {
    isLoading.value = true
    const data = await tenantApi.list()
    tenants.value = data
  } catch (error) {
    console.error('Failed to load tenants', error)
  } finally {
    isLoading.value = false
  }
}

// Watch for auth state changes to load tenants once admin is confirmed
watch(isSuperAdmin, (newValue) => {
    if (newValue) {
        if (tenantStore.tenant) {
            selectedTenant.value = tenantStore.tenant
        }
        loadTenants()
    }
}, { immediate: true })

const onTenantChange = (event: { value: TenantPublicInfo }) => {
  const newTenant = event.value
  if (!newTenant) return

  // 1. Update Cookie
  setCookie('tenant_id', newTenant.id, { expires: 7, path: '/' })
  
  // 2. Update Store
  tenantStore.setTenant(newTenant) // Updates branding immediately
  
  // 3. Reload page to ensure all components/API clients refresh with new context
  toast.info(`Switched to ${newTenant.name}`)
  
  setTimeout(() => {
     window.location.reload()
  }, 500)
}
</script>

<template>
  <div v-if="isSuperAdmin" class="tenant-switcher">
    <Select
      v-model="selectedTenant"
      :options="tenants"
      optionLabel="name"
      placeholder="Select Tenant"
      class="tenant-dropdown"
      :loading="isLoading"
      @change="onTenantChange"
    >
      <template #value="slotProps">
        <div v-if="slotProps.value" class="flex items-center gap-2">
           <i class="pi pi-building text-sm"></i>
           <span>{{ slotProps.value.name }}</span>
        </div>
        <span v-else>
            {{ slotProps.placeholder }}
        </span>
      </template>
      <template #option="slotProps">
        <div class="flex items-center gap-2">
          <span>{{ slotProps.option.name }}</span>
          <span v-if="slotProps.option.slug === 'default'" class="text-xs text-gray-400">(Default)</span>
        </div>
      </template>
    </Select>
  </div>
</template>

<style scoped>
.tenant-switcher {
  margin-right: 1rem;
  min-width: 200px;
}

.tenant-dropdown {
  width: 100%;
  border: none;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 6px;
}

:deep(.p-select-label) {
  color: white;
  font-size: 0.875rem;
  padding: 0.5rem 0.75rem;
}

:deep(.p-select-trigger) {
  color: rgba(255, 255, 255, 0.7);
  width: 2rem;
}

/* Select panel styles */
:deep(.p-select-panel) {
  border: 1px solid var(--border-color);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
</style>
