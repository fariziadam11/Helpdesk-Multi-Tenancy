<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { authApi } from '@/api/auth'
import { useToast } from '@/composables/useToast'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'

const { t } = useI18n()
const router = useRouter()
const toast = useToast()

const email = ref('')
const isLoading = ref(false)
const isSuccess = ref(false)

const handleSubmit = async () => {
  if (!email.value.trim()) {
    toast.error('Email harus diisi!')
    return
  }

  isLoading.value = true

  try {
    await authApi.forgotPassword(email.value)
    isSuccess.value = true
    toast.success('Link reset password telah dikirim ke email Anda')
  } catch (error: any) {
    console.error('Forgot password error:', error)
    toast.error(error.message || 'Gagal mengirim email reset password')
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
          <h2 class="auth-title">{{ t('auth.forgotPassword.title') }}</h2>
          <p class="auth-subtitle">{{ t('auth.forgotPassword.subtitle') }}</p>
        </div>

        <div v-if="isSuccess" class="success-message">
          <i class="pi pi-check-circle"></i>
          <p>{{ t('auth.forgotPassword.successMessage') }}</p>
          <Button
            :label="t('auth.forgotPassword.backToLogin')"
            @click="router.push('/login')"
            class="back-btn"
          />
        </div>

        <form v-else @submit.prevent="handleSubmit" class="auth-form">
          <div class="form-item">
            <label for="email" class="form-label">
              {{ t('auth.forgotPassword.email') }} <span class="required">*</span>
            </label>
            <InputText
              id="email"
              v-model="email"
              type="email"
              :placeholder="t('auth.forgotPassword.emailPlaceholder')"
              :disabled="isLoading"
              class="form-input"
            />
          </div>

          <div class="form-actions">
            <Button
              type="submit"
              :label="t('auth.forgotPassword.submit')"
              :loading="isLoading"
              :disabled="isLoading"
              class="submit-btn"
            />
          </div>

          <div class="auth-footer">
            <p>
              {{ t('auth.forgotPassword.rememberPassword') }}
              <a href="#" @click.prevent="router.push('/login')" class="auth-link">
                {{ t('auth.forgotPassword.signIn') }}
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
  margin-bottom: 2rem;
  line-height: 1.6;
}

.back-btn {
  width: 100%;
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

.form-input {
  width: 100%;
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
