<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  password: string
}

const props = defineProps<Props>()

interface PasswordRequirement {
  label: string
  met: boolean
}

const requirements = computed<PasswordRequirement[]>(() => {
  const pwd = props.password
  
  return [
    {
      label: 'Minimum 6 karakter',
      met: pwd.length >= 6
    },
    {
      label: 'Minimum 1 huruf kapital',
      met: /[A-Z]/.test(pwd)
    },
    {
      label: 'Minimum 2 karakter spesial (.,-_();: dll)',
      met: (pwd.match(/[^a-zA-Z0-9]/g) || []).length >= 2
    }
  ]
})

const metCount = computed(() => requirements.value.filter(r => r.met).length)
const totalCount = computed(() => requirements.value.length)
const strengthPercentage = computed(() => (metCount.value / totalCount.value) * 100)

const strengthColor = computed(() => {
  if (strengthPercentage.value === 100) return '#24a148' // Green - All met
  return '#da1e28' // Red - Not all met
})

const strengthLabel = computed(() => {
  if (strengthPercentage.value === 100) return 'Semua Kriteria Terpenuhi'
  return 'Kriteria Belum Terpenuhi'
})
</script>

<template>
  <div v-if="password" class="password-strength">
    <div class="strength-header">
      <span class="strength-label">Kriteria Kata Sandi: <strong>{{ strengthLabel }}</strong></span>
      <span class="strength-count">{{ metCount }}/{{ totalCount }}</span>
    </div>
    
    <div class="strength-bar-container">
      <div 
        class="strength-bar" 
        :style="{ 
          width: `${strengthPercentage}%`,
          backgroundColor: strengthColor
        }"
      ></div>
    </div>

    <ul class="requirements-list">
      <li 
        v-for="(req, index) in requirements" 
        :key="index"
        :class="{ 'met': req.met }"
      >
        <i :class="req.met ? 'pi pi-check-circle' : 'pi pi-circle'"></i>
        <span>{{ req.label }}</span>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.password-strength {
  margin-top: 0.75rem;
  padding: 1rem;
  background-color: #f4f4f4;
  border-radius: 4px;
  border: 1px solid var(--border-color);
}

.strength-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.strength-label {
  font-size: 0.875rem;
  color: var(--text-primary);
}

.strength-label strong {
  font-weight: 600;
}

.strength-count {
  font-size: 0.8125rem;
  color: var(--text-secondary);
  font-weight: 500;
}

.strength-bar-container {
  width: 100%;
  height: 8px;
  background-color: #e0e0e0;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 0.75rem;
}

.strength-bar {
  height: 100%;
  transition: width 0.3s ease, background-color 0.3s ease;
  border-radius: 4px;
}

.requirements-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.requirements-list li {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.8125rem;
  color: var(--text-secondary);
  transition: color 0.2s ease;
}

.requirements-list li.met {
  color: #24a148;
}

.requirements-list li i {
  font-size: 1rem;
  flex-shrink: 0;
}

.requirements-list li.met i {
  color: #24a148;
}

.requirements-list li:not(.met) i {
  color: #8d8d8d;
}
</style>
