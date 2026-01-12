import { createI18n } from 'vue-i18n'
import id from './locales/id'
import en from './locales/en'

const i18n = createI18n({
  legacy: false, // Use Composition API
  locale: localStorage.getItem('locale') || 'id',
  fallbackLocale: 'id',
  messages: {
    id,
    en,
  },
})

export default i18n
