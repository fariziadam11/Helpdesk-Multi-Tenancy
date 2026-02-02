<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import ToggleSwitch from 'primevue/toggleswitch'
import Password from 'primevue/password'
import ColorPicker from 'primevue/colorpicker'
import Divider from 'primevue/divider'
import FileUpload from 'primevue/fileupload'
import { tenantApi, type CreateTenantRequest } from '@/api/tenant'
import { uploadApi } from '@/api/upload'
import { useToast } from '@/composables/useToast'
import { logger } from '@/utils/logger'

const props = defineProps<{
  id?: string
}>()

const router = useRouter()
const toast = useToast()

const isEditMode = computed(() => !!props.id)
const isSubmitting = ref(false)
const isLoading = ref(false)

const form = ref<Partial<CreateTenantRequest>>({
  name: '',
  slug: '',
  logo_url: '',
  primary_color: '#6929C4',
  is_active: true,
  // InvGate Defaults
  invgate_base_url: 'https://cdn.invgate.net',
  invgate_company_id: undefined,
  invgate_group_id: undefined,
  invgate_location_id: undefined,
  invgate_username: '',
  invgate_password: '', // Empty means unchanged on Edit
  // Email Defaults
  email_domain: '',
  email_sender: ''
})

// Helper for color picker binding (it uses hex string without #)
const colorPickerValue = computed({
  get: () => form.value.primary_color?.replace('#', '') || '6929C4',
  set: (val: any) => {
    form.value.primary_color = `#${val}`
  }
})

const loadTenant = async () => {
  if (!props.id) return
  
  isLoading.value = true
  try {
    const data = await tenantApi.getById(props.id)
    form.value = {
      name: data.name,
      slug: data.slug,
      logo_url: data.logo_url || '',
      primary_color: data.primary_color || '#6929C4',
      is_active: data.is_active,
      invgate_base_url: data.invgate_base_url,
      invgate_company_id: data.invgate_company_id,
      invgate_group_id: data.invgate_group_id,
      invgate_location_id: data.invgate_location_id,
      invgate_username: data.invgate_username,
      invgate_password: '', // Never returned by API
      email_domain: data.email_domain || '',
      email_sender: data.email_sender || ''
    }
  } catch (error) {
    logger.error('Failed to load tenant', error)
    toast.error('Failed to load tenant details')
    router.push('/admin/tenants')
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  if (isEditMode.value) {
    loadTenant()
  }
})

const onUploadLogo = async (event: any) => {
    const file = event.files[0]
    try {
        const response = await uploadApi.uploadFile(file)
        // Store relative path (e.g. /uploads/uuid.jpg)
        form.value.logo_url = response.url
        toast.success('Logo uploaded')
    } catch (e: any) {
        toast.error('Upload failed: ' + (e.message || 'Unknown error'))
    }
}

const getPreviewUrl = (url: string) => {
    if (!url) return ''
    if (url.startsWith('http') || url.startsWith('data:')) return url
    if (url.startsWith('/uploads/')) {
        const apiBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
        const origin = new URL(apiBase).origin
        return `${origin}${url}`
    }
    return url
}

const handleSubmit = async () => {
    // Validation
    if (!form.value.name || !form.value.slug) {
        toast.error('Name and Slug are required')
        return
    }
    
    // InvGate Validation
    if (!form.value.invgate_base_url || !form.value.invgate_username || 
        form.value.invgate_company_id === undefined || 
        form.value.invgate_group_id === undefined || 
        form.value.invgate_location_id === undefined) {
        toast.error('All InvGate Configuration fields are required')
        return
    }

    if (!isEditMode.value && !form.value.invgate_password) {
        toast.error('InvGate Password is required for new tenants')
        return
    }

    isSubmitting.value = true
    try {
        const payload = { ...form.value }
        
        // Remove password if empty in edit mode
        if (isEditMode.value && !payload.invgate_password) {
            delete payload.invgate_password
        }

        if (isEditMode.value && props.id) {
            await tenantApi.update(props.id, payload)
            toast.success('Tenant updated successfully')
        } else {
            // Cast to CreateTenantRequest because we validated required fields
            await tenantApi.create(payload as CreateTenantRequest)
            toast.success('Tenant created successfully')
        }
        router.push('/admin/tenants')
    } catch (error: any) {
        logger.error('Failed to save tenant', error)
        const msg = error.response?.data?.error || error.message || 'Failed to save tenant'
        toast.error(msg)
    } finally {
        isSubmitting.value = false
    }
}

const handleCancel = () => {
    router.push('/admin/tenants')
}

// Auto-generate slug from name if creating new
const handleNameInput = () => {
    if (!isEditMode.value && form.value.name) {
        form.value.slug = form.value.name
            .toLowerCase()
            .replace(/[^a-z0-9]+/g, '-')
            .replace(/^-+|-+$/g, '')
    }
}
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">{{ isEditMode ? 'Edit Tenant' : 'Create New Tenant' }}</h1>
    </div>

    <div class="form-card" v-if="!isLoading">
      <form @submit.prevent="handleSubmit" class="tenant-form">
        
        <!-- Basic Info -->
        <div class="form-section">
            <h3 class="section-title">Basic Information</h3>
            <div class="form-grid">
                <div class="form-group">
                    <label for="name" class="form-label">Company Name <span class="required">*</span></label>
                    <InputText id="name" v-model="form.name" @input="handleNameInput" placeholder="e.g. Acme Corp" class="w-full" :disabled="isSubmitting" />
                </div>
                <div class="form-group">
                    <label for="slug" class="form-label">Slug (Subdomain) <span class="required">*</span></label>
                    <InputText id="slug" v-model="form.slug" placeholder="e.g. acme-corp" class="w-full" :disabled="isEditMode || isSubmitting" /> 
                    <small class="helper-text">Used for subdomain: {{ form.slug || 'slug' }}.yourapp.com</small>
                </div>
                <div class="form-group flex justify-end items-center mt-6">
                     <span class="mr-2 text-sm font-medium text-gray-700">Active Status</span>
                     <ToggleSwitch v-model="form.is_active" :disabled="isSubmitting" />
                </div>
            </div>
        </div>

        <Divider />

        <!-- InvGate Configuration -->
        <div class="form-section">
            <h3 class="section-title">InvGate Configuration</h3>
            <p class="section-subtitle">Credentials for the ITSM integration.</p>
            
            <div class="form-grid">
                 <div class="form-group">
                    <label class="form-label">InvGate Base URL <span class="required">*</span></label>
                    <InputText v-model="form.invgate_base_url" placeholder="https://cdn.invgate.net" class="w-full" :disabled="isSubmitting" />
                </div>
                <div class="form-group">
                    <label class="form-label">InvGate Username <span class="required">*</span></label>
                    <InputText v-model="form.invgate_username" placeholder="API User" class="w-full" :disabled="isSubmitting" />
                </div>
                <div class="form-group">
                    <label class="form-label">InvGate Password <span class="required" v-if="!isEditMode">*</span></label>
                     <Password v-model="form.invgate_password" :feedback="false" toggleMask class="w-full" inputClass="w-full" :placeholder="isEditMode ? 'Leave blank to keep unchanged' : ''" :disabled="isSubmitting" />
                </div>
                 <div class="form-group">
                    <label class="form-label">Company ID <span class="required">*</span></label>
                    <InputNumber v-model="form.invgate_company_id" class="w-full" :useGrouping="false" :disabled="isSubmitting" />
                </div>
                 <div class="form-group">
                    <label class="form-label">Group ID <span class="required">*</span></label>
                    <InputNumber v-model="form.invgate_group_id" class="w-full" :useGrouping="false" :disabled="isSubmitting" />
                </div>
                 <div class="form-group">
                    <label class="form-label">Location ID <span class="required">*</span></label>
                    <InputNumber v-model="form.invgate_location_id" class="w-full" :useGrouping="false" :disabled="isSubmitting" />
                </div>
            </div>
        </div>

        <Divider />

        <!-- Email Configuration -->
        <div class="form-section">
            <h3 class="section-title">Email Configuration</h3>
            <p class="section-subtitle">Optional overrides for email notifications.</p>
            
            <div class="form-grid">
               <div class="form-group">
                    <label class="form-label">Email Domain</label>
                    <InputText v-model="form.email_domain" placeholder="e.g. acme.com (for auto-assignment)" class="w-full" :disabled="isSubmitting" />
                </div>
                <div class="form-group">
                    <label class="form-label">Email Sender Name</label>
                    <InputText v-model="form.email_sender" placeholder="e.g. Acme Support" class="w-full" :disabled="isSubmitting" />
                </div>
            </div>
        </div>

        <Divider />

        <!-- Branding -->
        <div class="form-section">
            <h3 class="section-title">Branding</h3>

            <div class="form-grid">
                 <div class="form-group">
                    <label for="primary_color" class="form-label">Primary Color</label>
                    <div class="flex items-center gap-3">
                         <ColorPicker v-model="colorPickerValue" />
                         <span class="text-sm font-mono bg-gray-100 px-2 py-1 rounded">{{ form.primary_color }}</span>
                         <InputText v-model="form.primary_color" placeholder="#RRGGBB" class="w-32" :disabled="isSubmitting" />
                    </div>
                </div>
                
                 <div class="form-group">
                    <label class="form-label">Logo</label>
                    <div class="flex flex-col gap-2">
                        <div class="flex gap-2 items-center">
                            <FileUpload mode="basic" name="file" auto customUpload @uploader="onUploadLogo" accept="image/*" chooseLabel="Upload Logo" class="p-button-outlined p-button-secondary" :maxFileSize="1000000" />
                            <Button v-if="form.logo_url" label="Delete Logo" icon="pi pi-trash" severity="danger" outlined @click="form.logo_url = ''" size="small" />
                        </div>
                        <InputText v-if="form.logo_url" v-model="form.logo_url" class="w-full text-sm bg-gray-50 text-gray-500" readonly />
                    </div>
                    
                    <div class="mt-2 p-4 border border-dashed border-gray-300 rounded flex justify-center bg-gray-50 preview-box">
                        <img v-if="form.logo_url" :src="getPreviewUrl(form.logo_url)" alt="Logo Preview" class="h-12 object-contain" />
                        <span v-else class="text-gray-400 text-sm">Logo preview</span>
                    </div>
                </div>
            </div>
            
            <!-- Preview Header -->
            <div class="mt-6 border rounded overflow-hidden">
                <div class="text-xs bg-gray-100 p-2 border-b">Header Preview</div>
                <header class="p-4 flex items-center justify-between" :style="{ backgroundColor: form.primary_color }">
                     <img v-if="form.logo_url" :src="getPreviewUrl(form.logo_url)" alt="Logo" class="h-8 object-contain" />
                     <div v-else class="text-white font-bold text-lg">{{ form.name || 'Company Name' }}</div>
                     
                     <div class="flex gap-2">
                        <div class="w-8 h-8 rounded-full bg-white opacity-20"></div>
                        <div class="w-8 h-8 rounded-full bg-white opacity-20"></div>
                     </div>
                </header>
            </div>
        </div>

        <div class="form-actions">
           <Button type="button" label="Cancel" severity="secondary" @click="handleCancel" :disabled="isSubmitting" />
           <Button type="submit" label="Save Tenant" icon="pi pi-check" :loading="isSubmitting" />
        </div>

      </form>
    </div>
    
    <div v-else class="loading-state">
        <i class="pi pi-spin pi-spinner text-4xl text-primary"></i>
    </div>
  </div>
</template>

<style scoped>
.page-container {
  padding: 2rem;
  max-width: 800px;
  margin: 0 auto;
}

.page-title {
  font-size: 1.875rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 2rem;
}

.form-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  padding: 2rem;
  border: 1px solid var(--border-color);
}

.form-section {
    margin-bottom: 0;
}

.section-title {
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 0.5rem;
}

.section-subtitle {
    font-size: 0.875rem;
    color: var(--text-secondary);
    margin-bottom: 1.5rem;
    margin-top: 0;
}

.form-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: 1.5rem;
}

@media (min-width: 640px) {
    .form-grid {
        grid-template-columns: 1fr 1fr;
    }
}

.form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
}

.form-label {
    font-weight: 500;
    font-size: 0.875rem;
    color: var(--text-primary);
}

.required {
    color: var(--error-color);
}

.helper-text {
    font-size: 0.75rem;
    color: var(--text-secondary);
}

.form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 2rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--border-color);
}

.loading-state {
    display: flex;
    justify-content: center;
    padding: 4rem;
}

.preview-box {
    min-height: 80px;
    display: flex;
    align-items: center;
    justify-content: center;
}

:deep(.p-password-input) {
    width: 100%;
}
</style>
