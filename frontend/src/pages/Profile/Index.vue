<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { usersApi } from '@/api/users'
import { useToast } from '@/composables/useToast'
import PasswordStrength from '@/components/shared/PasswordStrength.vue'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'

const { t } = useI18n()
const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

const isSaving = ref(false)

// Form data
const formData = ref({
  name: '',
  lastname: '',
  email: '',
  password: '',
  confirmPassword: '',
})

// Form errors
const errors = ref({
  name: '',
  lastname: '',
  password: '',
  confirmPassword: '',
})

onMounted(() => {
  if (authStore.user) {
    formData.value.name = authStore.user.name || ''
    formData.value.lastname = authStore.user.lastname || ''
    formData.value.email = authStore.user.email || ''
  }
})

const validateForm = (): boolean => {
  let isValid = true
  errors.value = {
    name: '',
    lastname: '',
    password: '',
    confirmPassword: '',
  }

  if (!formData.value.name.trim()) {
    errors.value.name = t('profile.errors.nameRequired')
    isValid = false
  }

  if (!formData.value.lastname.trim()) {
    errors.value.lastname = t('profile.errors.lastnameRequired')
    isValid = false
  }

  // Password validation (only if password is provided)
  if (formData.value.password) {
    const hasMinLength = formData.value.password.length >= 6
    const hasUpperCase = /[A-Z]/.test(formData.value.password)
    const hasSpecialChars = (formData.value.password.match(/[^a-zA-Z0-9]/g) || []).length >= 2
    
    if (!hasMinLength) {
      errors.value.password = t('profile.errors.passwordTooShort')
      isValid = false
    } else if (!hasUpperCase) {
      errors.value.password = 'Kata sandi harus mengandung minimal 1 huruf kapital'
      isValid = false
    } else if (!hasSpecialChars) {
      errors.value.password = 'Kata sandi harus mengandung minimal 2 karakter spesial'
      isValid = false
    }

    if (formData.value.password !== formData.value.confirmPassword) {
      errors.value.confirmPassword = t('profile.errors.passwordMismatch')
      isValid = false
    }
  }

  return isValid
}

const handleSubmit = async () => {
  if (!validateForm()) {
    return
  }

  isSaving.value = true

  try {
    const updateData: { name: string; lastname: string; password?: string } = {
      name: formData.value.name.trim(),
      lastname: formData.value.lastname.trim(),
    }

    // Only include password if it's provided
    if (formData.value.password) {
      updateData.password = formData.value.password
    }

    const updatedUser = await usersApi.updateProfile(updateData)

    // Update auth store with new user data
    authStore.updateUser({
      name: updatedUser.name || formData.value.name,
      lastname: updatedUser.lastname || formData.value.lastname,
      email: updatedUser.email || formData.value.email,
    })

    toast.success(t('profile.updateSuccess'))

    // Clear password fields after successful update
    formData.value.password = ''
    formData.value.confirmPassword = ''
  } catch (error: any) {
    console.error('Failed to update profile:', error)
    toast.error(error.response?.data?.message || t('profile.updateError'))
  } finally {
    isSaving.value = false
  }
}

const handleCancel = () => {
  router.push('/dashboard')
}
</script>

<template>
  <div class="profile-page">
    <div class="page-header">
      <button class="back-button" @click="handleCancel">
        <i class="pi pi-arrow-left"></i>
        {{ t('profile.backToDashboard') }}
      </button>
    </div>

    <div class="profile-container">
      <div class="profile-card">
        <div class="card-header">
          <h1 class="card-title">
            <i class="pi pi-user-edit"></i>
            {{ t('profile.title') }}
          </h1>
          <p class="card-description">{{ t('profile.description') }}</p>
        </div>

        <form @submit.prevent="handleSubmit" class="profile-form">
          <!-- Name Field -->
          <div class="form-group">
            <label for="name" class="form-label">
              {{ t('profile.fields.name') }}
              <span class="required">*</span>
            </label>
            <InputText
              id="name"
              v-model="formData.name"
              :class="{ 'p-invalid': errors.name }"
              :disabled="isSaving"
              class="form-input"
            />
            <small v-if="errors.name" class="error-message">{{ errors.name }}</small>
          </div>

          <!-- Last Name Field -->
          <div class="form-group">
            <label for="lastname" class="form-label">
              {{ t('profile.fields.lastname') }}
              <span class="required">*</span>
            </label>
            <InputText
              id="lastname"
              v-model="formData.lastname"
              :class="{ 'p-invalid': errors.lastname }"
              :disabled="isSaving"
              class="form-input"
            />
            <small v-if="errors.lastname" class="error-message">{{ errors.lastname }}</small>
          </div>

          <!-- Email Field (Disabled) -->
          <div class="form-group">
            <label for="email" class="form-label">
              {{ t('profile.fields.email') }}
            </label>
            <InputText
              id="email"
              v-model="formData.email"
              disabled
              class="form-input"
            />
            <small class="field-hint">{{ t('profile.hints.emailDisabled') }}</small>
          </div>

          <!-- Password Field -->
          <div class="form-group">
            <label for="password" class="form-label">
              {{ t('profile.fields.password') }}
            </label>
            <Password
              id="password"
              v-model="formData.password"
              :class="{ 'p-invalid': errors.password }"
              :disabled="isSaving"
              :feedback="false"
              toggleMask
              class="form-input"
            />
            <small v-if="errors.password" class="error-message">{{ errors.password }}</small>
            <small v-else class="field-hint">{{ t('profile.hints.passwordOptional') }}</small>
            <PasswordStrength :password="formData.password" />
          </div>

          <!-- Confirm Password Field -->
          <div class="form-group">
            <label for="confirmPassword" class="form-label">
              {{ t('profile.fields.confirmPassword') }}
            </label>
            <Password
              id="confirmPassword"
              v-model="formData.confirmPassword"
              :class="{ 'p-invalid': errors.confirmPassword }"
              :disabled="isSaving"
              :feedback="false"
              toggleMask
              class="form-input"
            />
            <small v-if="errors.confirmPassword" class="error-message">{{ errors.confirmPassword }}</small>
          </div>

          <!-- Form Actions -->
          <div class="form-actions">
            <Button
              type="button"
              :label="t('profile.actions.cancel')"
              severity="secondary"
              :disabled="isSaving"
              @click="handleCancel"
            />
            <Button
              type="submit"
              :label="t('profile.actions.save')"
              :loading="isSaving"
            />
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.profile-page {
  padding: 2rem;
  max-width: 800px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 2rem;
}

.back-button {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: none;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 0.5rem 1rem;
  color: var(--text-primary);
  cursor: pointer;
  font-size: 0.875rem;
  transition: all 0.2s ease;
}

.back-button:hover {
  background-color: #E8DAFF;
  border-color: #6929C4;
  color: #6929C4;
}

.back-button i {
  font-size: 0.875rem;
}

.profile-container {
  display: flex;
  justify-content: center;
}

.profile-card {
  background: #ffffff;
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 2rem;
  width: 100%;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.card-header {
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 2px solid var(--border-color);
}

.card-title {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 1.75rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.5rem 0;
}

.card-title i {
  color: #6929C4;
  font-size: 1.5rem;
}

.card-description {
  color: var(--text-secondary);
  font-size: 0.9375rem;
  margin: 0;
}

.profile-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-label {
  font-size: 0.9375rem;
  font-weight: 500;
  color: var(--text-primary);
}

.required {
  color: #da1e28;
  margin-left: 0.25rem;
}

.form-input {
  width: 100%;
}

.form-input :deep(.p-inputtext) {
  width: 100%;
}

.error-message {
  color: #da1e28;
  font-size: 0.8125rem;
  margin-top: 0.25rem;
}

.field-hint {
  color: var(--text-secondary);
  font-size: 0.8125rem;
  margin-top: 0.25rem;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border-color);
}

@media (max-width: 768px) {
  .profile-page {
    padding: 1rem;
  }

  .profile-card {
    padding: 1.5rem;
  }

  .card-title {
    font-size: 1.5rem;
  }

  .form-actions {
    flex-direction: column-reverse;
  }

  .form-actions button {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .card-title {
    font-size: 1.25rem;
  }

  .back-button {
    width: 100%;
    justify-content: center;
  }
}
</style>
