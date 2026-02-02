<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import { tenantApi, type Tenant } from '@/api/tenant'
import { useToast } from '@/composables/useToast'
import { logger } from '@/utils/logger'

const router = useRouter()
const toast = useToast()

const tenants = ref<Tenant[]>([])
const isLoading = ref(false)

const loadTenants = async () => {
    isLoading.value = true
    try {
        const data = await tenantApi.list()
        tenants.value = data
    } catch (error) {
        logger.error('Failed to load tenants', error)
        toast.error('Failed to load tenants')
    } finally {
        isLoading.value = false
    }
}

onMounted(() => {
    loadTenants()
})

const getLogoUrl = (url: string) => {
    if (!url) return ''
    if (url.startsWith('http') || url.startsWith('data:')) return url
    if (url.startsWith('/uploads/')) {
        const apiBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
        const origin = new URL(apiBase).origin
        return `${origin}${url}`
    }
    return url
}

const handleEdit = (tenant: Tenant) => {
    router.push(`/admin/tenants/${tenant.id}/edit`)
}

const handleToggleStatus = async (tenant: Tenant) => {
    try {
        const newStatus = !tenant.is_active
        await tenantApi.updateStatus(tenant.id, newStatus)
        tenant.is_active = newStatus
        toast.success(`Tenant ${newStatus ? 'activated' : 'deactivated'} successfully`)
    } catch (error) {
        logger.error('Failed to update tenant status', error)
        toast.error('Failed to update tenant status')
    }
}

const handleCreate = () => {
    router.push('/admin/tenants/create')
}
</script>

<template>
    <div class="page-container">
        <div class="page-header">
            <div>
                <h1 class="page-title">Tenant Management</h1>
                <p class="page-subtitle">Manage system tenants and companies</p>
            </div>
            <Button label="Add Tenant" icon="pi pi-plus" @click="handleCreate" />
        </div>

        <div class="card">
            <DataTable :value="tenants" :loading="isLoading" stripedRows responsiveLayout="scroll">
                <template #empty>No tenants found.</template>
                
                <Column field="name" header="Company Name" sortable>
                    <template #body="{ data }">
                        <div class="flex items-center gap-2">
                             <img v-if="data.logo_url" :src="getLogoUrl(data.logo_url)" alt="Logo" class="w-8 h-8 object-contain rounded bg-gray-100 p-1" />
                             <div v-else class="w-8 h-8 rounded bg-gray-200 flex items-center justify-center text-xs font-bold text-gray-500">
                                {{ data.name.substring(0, 2).toUpperCase() }}
                             </div>
                             <span class="font-medium">{{ data.name }}</span>
                        </div>
                    </template>
                </Column>
                
                <Column field="slug" header="Slug" sortable>
                    <template #body="{ data }">
                        <code class="bg-gray-100 px-2 py-1 rounded text-sm text-gray-700">{{ data.slug }}</code>
                    </template>
                </Column>

                <Column field="primary_color" header="Brand Color">
                    <template #body="{ data }">
                        <div class="flex items-center gap-2">
                            <div class="w-6 h-6 rounded border border-gray-200" :style="{ backgroundColor: data.primary_color }"></div>
                            <span class="text-sm text-gray-600">{{ data.primary_color }}</span>
                        </div>
                    </template>
                </Column>
                
                <Column field="is_active" header="Status" sortable>
                     <template #body="{ data }">
                        <Tag :severity="data.is_active ? 'success' : 'danger'" :value="data.is_active ? 'Active' : 'Inactive'" />
                    </template>
                </Column>
                
                <Column header="Actions" alignFrozen="right" :exportable="false" style="min-width: 8rem">
                    <template #body="{ data }">
                        <div class="flex gap-2">
                            <Button icon="pi pi-pencil" severity="secondary" text rounded aria-label="Edit" @click="handleEdit(data)" />
                            <Button 
                                :icon="data.is_active ? 'pi pi-ban' : 'pi pi-check-circle'" 
                                :severity="data.is_active ? 'danger' : 'success'" 
                                text rounded 
                                :aria-label="data.is_active ? 'Deactivate' : 'Activate'"
                                @click="handleToggleStatus(data)" 
                                v-tooltip.top="data.is_active ? 'Deactivate' : 'Activate'"
                            />
                        </div>
                    </template>
                </Column>
            </DataTable>
        </div>
    </div>
</template>

<style scoped>
.page-container {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-title {
  font-size: 1.875rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.5rem 0;
}

.page-subtitle {
  color: var(--text-secondary);
  margin: 0;
}

.card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  padding: 1.5rem;
  border: 1px solid var(--border-color);
}
</style>
