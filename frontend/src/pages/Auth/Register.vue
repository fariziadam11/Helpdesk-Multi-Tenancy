<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useRegister } from '@/composables/useAuth'
import { useToast } from '@/composables/useToast'
import GuestHeader from '@/components/GuestHeader.vue'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'

const { t } = useI18n()

const router = useRouter()
const toast = useToast()

const name = ref('')
const lastname = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')

const { mutate: register, isPending } = useRegister()

const hasUpperCase = computed(() => {
  return /[A-Z]/.test(password.value)
})

const isFormValid = computed(() => {
  return (
    name.value.trim() !== '' &&
    lastname.value.trim() !== '' &&
    email.value.trim() !== '' &&
    password.value.length >= 6 &&
    hasUpperCase.value &&
    password.value === confirmPassword.value
  )
})

const passwordMismatch = computed(() => {
  return (
    confirmPassword.value !== '' && password.value !== confirmPassword.value
  )
})

const handleSubmit = (e?: Event) => {
  if (e) {
    e.preventDefault()
  }

  if (!isFormValid.value || passwordMismatch.value) {
    return
  }

  // Validasi duplicate firstname, lastname, dan email
  const trimmedName = name.value.trim()
  const trimmedLastname = lastname.value.trim()
  const trimmedEmail = email.value.trim()

  // Check if firstname and lastname are the same
  if (trimmedName.toLowerCase() === trimmedLastname.toLowerCase()) {
    toast.error(t('auth.register.validation.nameLastnameSame'))
    return
  }

  // Check if firstname and email are the same
  if (trimmedName.toLowerCase() === trimmedEmail.toLowerCase()) {
    toast.error(t('auth.register.validation.nameEmailSame'))
    return
  }

  // Check if lastname and email are the same
  if (trimmedLastname.toLowerCase() === trimmedEmail.toLowerCase()) {
    toast.error(t('auth.register.validation.lastnameEmailSame'))
    return
  }

  register({
    name: trimmedName,
    lastname: trimmedLastname,
    email: trimmedEmail,
    password: password.value,
  })
}
</script>

<template>
  <GuestHeader />
  <div class="auth-page">
    <div class="auth-container">
      <div id="registerTile" class="auth-card">
        <div class="auth-header">
          <h2 class="auth-title">{{ t('auth.register.title') }}</h2>
          <p class="auth-subtitle">{{ t('auth.register.subtitle') }}</p>
        </div>

        <form @submit.prevent="handleSubmit" class="auth-form">
          <div class="form-row">
            <div class="form-item form-item-half">
              <label for="registerName" class="form-label">
                {{ t('auth.register.firstName') }} <span class="required">*</span>
              </label>
              <InputText
                id="registerName"
                v-model="name"
                :placeholder="t('auth.register.firstNamePlaceholder')"
                :disabled="isPending"
                autocomplete="given-name"
                class="form-input"
              />
            </div>

            <div class="form-item form-item-half">
              <label for="registerLastname" class="form-label">
                {{ t('auth.register.lastName') }} <span class="required">*</span>
              </label>
              <InputText
                id="registerLastname"
                v-model="lastname"
                :placeholder="t('auth.register.lastNamePlaceholder')"
                :disabled="isPending"
                autocomplete="family-name"
                class="form-input"
              />
            </div>
          </div>

          <div class="form-item">
            <label for="registerEmail" class="form-label">
              {{ t('auth.register.email') }} <span class="required">*</span>
            </label>
            <InputText
              id="registerEmail"
              v-model="email"
              type="email"
              :placeholder="t('auth.register.emailPlaceholder')"
              :disabled="isPending"
              autocomplete="email"
              class="form-input"
            />
          </div>

          <div class="form-item">
            <label class="form-label">
              {{ t('auth.register.password') }} <span class="required">*</span>
            </label>
            <Password
              id="registerPassword"
              v-model="password"
              :placeholder="t('auth.register.passwordPlaceholder')"
              :disabled="isPending"
              :toggleMask="true"
              :feedback="false"
              :inputProps="{ autocomplete: 'new-password' }"
              inputClass="form-input"
              class="password-input"
            />
            <div class="input-feedback">
              <small v-if="password.length > 0 && password.length < 6" class="helper-text">
                <span class="icon"><i class="pi pi-exclamation-triangle"></i></span> {{ t('auth.register.passwordMinLength') }}
              </small>
              <small v-if="password.length >= 6 && !hasUpperCase" class="helper-text">
                <span class="icon"><i class="pi pi-exclamation-triangle"></i></span> {{ t('auth.register.passwordNeedsUppercase') }}
              </small>
            </div>
          </div>

          <div class="form-item">
            <label class="form-label">
              {{ t('auth.register.confirmPassword') }} <span class="required">*</span>
            </label>
            <Password
              id="registerConfirmPassword"
              v-model="confirmPassword"
              :placeholder="t('auth.register.confirmPasswordPlaceholder')"
              :disabled="isPending"
              :toggleMask="true"
              :feedback="false"
              :inputProps="{ autocomplete: 'new-password' }"
              inputClass="form-input"
              class="password-input"
            />
            <div class="input-feedback">
              <small v-if="passwordMismatch" class="helper-text">
                <span class="icon"><i class="pi pi-exclamation-triangle"></i></span> {{ t('auth.register.passwordMismatch') }}
              </small>
            </div>
          </div>

          <div class="form-actions">
            <Button
              type="submit"
              :label="t('auth.register.createAccount')"
              :loading="isPending"
              :disabled="!isFormValid || isPending || passwordMismatch"
              class="submit-btn"
            />
          </div>

          <div class="auth-footer">
            <p>
              {{ t('auth.register.haveAccount') }}
              <a href="#" @click.prevent="router.push('/login')" class="auth-link">
                {{ t('auth.register.signIn') }}
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

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.form-item-half {
  flex: 1;
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

  .form-row {
    grid-template-columns: 1fr;
  }

  .auth-title {
    font-size: 1.75rem;
  }
}
</style>
