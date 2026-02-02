<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { locale } = useI18n()

const languages = [
  { code: 'id', name: 'ID', flag: '🇮🇩' },
  { code: 'en', name: 'EN', flag: '🇬🇧' },
]

const changeLanguage = (lang: string) => {
  locale.value = lang
  localStorage.setItem('locale', lang)
}
</script>

<template>
  <div class="language-switcher">
    <button
      v-for="lang in languages"
      :key="lang.code"
      class="lang-btn"
      :class="{ active: locale === lang.code }"
      @click="changeLanguage(lang.code)"
      :title="`Switch to ${lang.name}`"
    >
      <span class="flag">{{ lang.flag }}</span>
      <span class="code">{{ lang.name }}</span>
    </button>
  </div>
</template>

<style scoped>
.language-switcher {
  display: flex;
  gap: 0.5rem;
}

.lang-btn {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.5rem 0.75rem;
  border: 1px solid #e0e0e0;
  background: #ffffff;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.875rem;
}

.lang-btn:hover {
  border-color: #6929C4;
  background: #f8f4ff;
}

.lang-btn.active {
  background: #6929C4;
  color: #ffffff;
  border-color: #6929C4;
}

.flag {
  font-size: 1rem;
  line-height: 1;
}

.code {
  font-weight: 500;
}
</style>
