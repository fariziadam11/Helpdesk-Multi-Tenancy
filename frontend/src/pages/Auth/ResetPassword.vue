<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { authApi } from '@/api/auth'
import { useToast } from '@/composables/useToast'
import PasswordStrength from '@/components/shared/PasswordStrength.vue'
import Password from 'primevue/password'
import Button from 'primevue/button'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const toast = useToast()

const token = ref('')
const password = ref('')
const confirmPassword = ref('')
const isLoading = ref(false)
const isSuccess = ref(false)

onMounted(() => {
  token.value = (route.query.token as string) || ''
  if (!token.value) {
    toast.error('Token reset password tidak valid')
    router.push('/login')
  }
})

const passwordMismatch = ref(false)

const handleSubmit = async () => {
  if (!password.value || !confirmPassword.value) {
    toast.error('Semua field harus diisi')
    return
  }

  // Validate password criteria
  const hasMinLength = password.value.length >= 6
  const hasUpperCase = /[A-Z]/.test(password.value)
  const hasSpecialChars = (password.value.match(/[^a-zA-Z0-9]/g) || []).length >= 2

  if (!hasMinLength || !hasUpperCase || !hasSpecialChars) {
    toast.error('Password tidak memenuhi kriteria yang ditentukan')
    return
  }

  if (password.value !== confirmPassword.value) {
    passwordMismatch.value = true
    toast.error('Password tidak cocok')
    return
  }

  passwordMismatch.value = false
  isLoading.value = true

  try {
    await authApi.resetPassword(token.value, password.value)
    isSuccess.value = true
    toast.success('Password berhasil direset')
    
    // Redirect to login after 2 seconds
    setTimeout(() => {
      router.push('/login')
    }, 2000)
  } catch (error: any) {
    console.error('Reset password error:', error)
    toast.error(error.message || 'Gagal reset password. Token mungkin sudah kadaluarsa.')
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-container">
      <div class="auth-card">
        <div class="auth-header">
          <h2 class="auth-title">{{ t('auth.resetPassword.title') }}</h2>
          <p class="auth-subtitle">{{ t('auth.resetPassword.subtitle') }}</p>
        </div>

        <div v-if="isSuccess" class="success-message">
          <i class="pi pi-check-circle"></i>
          <p>{{ t('auth.resetPassword.successMessage') }}</p>
          <p class="redirect-text">{{ t('auth.resetPassword.redirecting') }}</p>
        </div>

        <form v-else @submit.prevent="handleSubmit" class="auth-form">
          <div class="form-item">
            <label class="form-label">
              {{ t('auth.resetPassword.newPassword') }} <span class="required">*</span>
            </label>
            <Password
              id="password"
              v-model="password"
              :placeholder="t('auth.resetPassword.newPasswordPlaceholder')"
              :disabled="isLoading"
              :toggleMask="true"
              :feedback="false"
              :inputProps="{ autocomplete: 'new-password' }"
              inputClass="form-input"
              class="password-input"
            />
            <PasswordStrength :password="password" />
          </div>

          <div class="form-item">
            <label class="form-label">
              {{ t('auth.resetPassword.confirmPassword') }} <span class="required">*</span>
            </label>
            <Password
              id="confirmPassword"
              v-model="confirmPassword"
              :placeholder="t('auth.resetPassword.confirmPasswordPlaceholder')"
              :disabled="isLoading"
              :toggleMask="true"
              :feedback="false"
              :inputProps="{ autocomplete: 'new-password' }"
              inputClass="form-input"
              class="password-input"
            />
            <small v-if="passwordMismatch" class="error-text">
              {{ t('auth.resetPassword.passwordMismatch') }}
            </small>
          </div>

          <div class="form-actions">
            <Button
              type="submit"
              :label="t('auth.resetPassword.submit')"
              :loading="isLoading"
              :disabled="isLoading"
              class="submit-btn"
            />
          </div>

          <div class="auth-footer">
            <p>
              {{ t('auth.resetPassword.rememberPassword') }}
              <a href="#" @click.prevent="router.push('/login')" class="auth-link">
                {{ t('auth.resetPassword.signIn') }}
              </a>
            </p>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--background);
  padding: 2rem 1rem;
}

.auth-container {
  width: 100%;
  max-width: 420px;
}

.auth-card {
  padding: 2.5rem;
  background-color: #ffffff;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.auth-header {
  text-align: center;
  margin-bottom: 2rem;
}

.auth-title {
  margin: 0 0 0.5rem 0;
  font-size: 2rem;
  font-weight: 600;
  color: var(--text-primary);
  letter-spacing: -0.02em;
}

.auth-subtitle {
  margin: 0;
  color: var(--text-secondary);
  font-size: 0.9375rem;
}

.success-message {
  text-align: center;
  padding: 2rem 0;
}

.success-message i {
  font-size: 4rem;
  color: #24a148;
  margin-bottom: 1rem;
}

.success-message p {
  color: var(--text-primary);
  margin-bottom: 1rem;
  line-height: 1.6;
}

.redirect-text {
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-label {
  display: block;
  margin-bottom: 0.375rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
}

.required {
  color: var(--error-color);
}

.password-input {
  width: 100%;
}

.password-input :deep(.p-password-input) {
  width: 100%;
}

.error-text {
  color: var(--error-color);
  font-size: 0.8125rem;
  margin-top: 0.25rem;
}

.form-actions {
  margin-top: 0.5rem;
}

.submit-btn {
  width: 100%;
  min-height: 3rem;
}

.auth-footer {
  margin-top: 1.5rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border-color);
  text-align: center;
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.auth-link {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 500;
  margin-left: 0.25rem;
  transition: color 0.15s;
}

.auth-link:hover {
  color: var(--primary-hover);
  text-decoration: underline;
}

@media (max-width: 640px) {
  .auth-card {
    padding: 2rem 1.5rem;
  }

  .auth-title {
    font-size: 1.75rem;
  }
}
</style>
