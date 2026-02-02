<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useLogin } from '@/composables/useAuth'
import GuestHeader from '@/components/GuestHeader.vue'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'

const { t } = useI18n()

const router = useRouter()

const email = ref('')
const password = ref('')

const { mutate: login, isPending } = useLogin()

const isFormValid = computed(() => {
  return email.value.trim() !== '' && password.value.length >= 6
})

const handleSubmit = (e?: Event) => {
  if (e) {
    e.preventDefault()
  }

  if (!isFormValid.value) {
    return
  }

  login({
    email: email.value.trim(),
    password: password.value,
  })
}
</script>

<template>
  <GuestHeader />
  <div class="auth-page">
    <div class="auth-container">
      <div id="loginTile" class="auth-card">
        <div class="auth-header">
          <h2 class="auth-title">{{ t('auth.login.title') }}</h2>
          <p class="auth-subtitle">{{ t('auth.login.subtitle') }}</p>
        </div>

        <form @submit.prevent="handleSubmit" class="auth-form">
          <div class="form-item">
            <label for="email" class="form-label">
              {{ t('auth.login.email') }} <span class="required">*</span>
            </label>
            <InputText
              id="email"
              v-model="email"
              type="email"
              :placeholder="t('auth.login.emailPlaceholder')"
              :disabled="isPending"
              autocomplete="email"
              class="form-input"
            />
          </div>

          <div class="form-item">
            <label class="form-label">
              {{ t('auth.login.password') }} <span class="required">*</span>
            </label>
            <Password
              id="password"
              v-model="password"
              :placeholder="t('auth.login.passwordPlaceholder')"
              :disabled="isPending"
              :toggleMask="true"
              :feedback="false"
              :inputProps="{ autocomplete: 'current-password' }"
              inputClass="form-input"
              class="password-input"
            />
            <div class="forgot-password-link">
              <a href="#" @click.prevent="router.push('/forgot-password')" class="auth-link">
                {{ t('auth.login.forgotPassword') }}
              </a>
            </div>
          </div>

          <div class="form-actions">
            <Button
              type="submit"
              :label="t('auth.login.signIn')"
              :loading="isPending"
              :disabled="!isFormValid || isPending"
              class="submit-btn"
            />
          </div>

          <div class="auth-footer">
            <p>
              {{ t('auth.login.noAccount') }}
              <a href="#" @click.prevent="router.push('/register')" class="auth-link">
                {{ t('auth.login.createAccount') }}
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
  display: block;
  width: 100%;
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

.password-input {
  width: 100%;
}

.password-input :deep(.p-password-input) {
  width: 100%;
}

.forgot-password-link {
  text-align: right;
  margin-top: 0.5rem;
}

.forgot-password-link .auth-link {
  font-size: 0.875rem;
}

.input-feedback {
  min-height: 1.25rem;
  margin-top: 0.25rem;
}

.helper-text {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  font-size: 0.8125rem;
  color: var(--text-secondary);
  margin: 0;
}

.icon {
  display: inline-flex;
  align-items: center;
}

.icon i {
  font-size: 0.875rem;
}

.form-actions {
  margin-top: 0.5rem;
}

.submit-btn {
  width: 100%;
  min-height: 3rem;
}

.submit-btn :deep(.p-button) {
  width: 100%;
  background-color: #6929C4;
  border-color: #6929C4;
}

.submit-btn :deep(.p-button:hover:not(:disabled)) {
  background-color: #8A3FFC;
  border-color: #8A3FFC;
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
