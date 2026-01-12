<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import LanguageSwitcher from './LanguageSwitcher.vue'

const { t } = useI18n()
const router = useRouter()

const goToLanding = () => {
  router.push('/')
}

const goToLogin = () => {
  router.push('/login')
}

const goToRegister = () => {
  router.push('/register')
}
</script>

<template>
  <header class="guest-header">
    <div class="header-content">
      <div class="header-left">
        <button class="brand" type="button" @click="goToLanding">
          <img src="/logo_white.svg" alt="Werk logo" class="brand-logo" />
        </button>
      </div>
      <div class="header-right">
        <LanguageSwitcher />
        <button class="register-btn" @click="goToRegister">
          {{ t('nav.register') }}
        </button>
        <button class="login-btn" @click="goToLogin">
          {{ t('nav.login') }}
        </button>
      </div>
    </div>
  </header>
</template>

<style scoped>
.guest-header {
  background-color: #6929C4;
  border-bottom: 1px solid #262626;
  padding: 0 2rem;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 3rem;
}

.header-left {
  display: flex;
  align-items: center;
  flex: 0 0 auto;
}

.brand {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: none;
  border: none;
  padding: 0;
  margin: 0;
  cursor: pointer;
  color: #ffffff;
}

.brand:focus-visible {
  outline: 2px solid #ffffff;
  outline-offset: 2px;
}

.brand-logo {
  height: 1.4rem;
  width: auto;
  display: block;
}

.app-title {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 500;
  color: #ffffff;
  white-space: nowrap;
}

.app-title:hover {
  color: #f4f4f4;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex: 0 0 auto;
  margin-left: auto;
}

.register-btn {
  background: #ffffff;
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: #6929C4;
  padding: 0.5rem 1.25rem;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.register-btn:hover {
  background: #f4f4f4;
  border-color: #ffffff;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.register-btn:focus-visible {
  outline: 2px solid #ffffff;
  outline-offset: 2px;
}

.login-btn {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: #ffffff;
  padding: 0.5rem 1.25rem;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.login-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  border-color: rgba(255, 255, 255, 0.5);
}

.login-btn:focus-visible {
  outline: 2px solid #ffffff;
  outline-offset: 2px;
}

@media (max-width: 768px) {
  .guest-header {
    padding: 0 1rem;
  }

  .app-title {
    font-size: 1rem;
  }

  .brand-logo {
    height: 1.2rem;
  }

  .header-right {
    gap: 0.5rem;
  }

  .register-btn,
  .login-btn {
    padding: 0.5rem 0.875rem;
    font-size: 0.8125rem;
  }
}

@media (max-width: 480px) {
  .guest-header {
    padding: 0 0.75rem;
  }

  .app-title {
    font-size: 0.875rem;
  }

  .brand-logo {
    height: 1rem;
  }

  .header-right {
    gap: 0.375rem;
  }

  .register-btn,
  .login-btn {
    padding: 0.375rem 0.625rem;
    font-size: 0.75rem;
  }
}

@media (max-width: 360px) {
  .guest-header {
    padding: 0 0.5rem;
  }

  .brand {
    gap: 0.25rem;
  }

  .app-title {
    font-size: 0.75rem;
  }

  .brand-logo {
    height: 0.875rem;
  }

  .header-right {
    gap: 0.25rem;
  }

  .register-btn,
  .login-btn {
    padding: 0.25rem 0.5rem;
    font-size: 0.6875rem;
  }
}
</style>

