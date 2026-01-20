<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { formatFileSize, removeFileFromArray } from '@/utils/file'

interface Props {
  files: File[]
  disabled?: boolean
  multiple?: boolean
  accept?: string
  id?: string
  maxFiles?: number
  maxFileSize?: number // in bytes
  allowedFormats?: string[] // file extensions without dot
}

interface Emits {
  (e: 'update:files', files: File[]): void
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  multiple: true,
  accept: undefined,
  id: undefined,
  maxFiles: undefined,
  maxFileSize: undefined,
  allowedFormats: undefined,
})

const emit = defineEmits<Emits>()

const inputId = computed(() => props.id || `fileInput-${Math.random().toString(36).substr(2, 9)}`)

const localFiles = ref<File[]>([...props.files])
const errorMessage = ref<string>('')

watch(
  () => props.files,
  (newFiles) => {
    localFiles.value = [...newFiles]
  },
  { deep: true },
)

const validateFile = (file: File): string | null => {
  // Check file size
  if (props.maxFileSize && file.size > props.maxFileSize) {
    const maxSizeMB = (props.maxFileSize / (1024 * 1024)).toFixed(0)
    return `File "${file.name}" melebihi ukuran maksimum ${maxSizeMB}MB`
  }

  // Check file format
  if (props.allowedFormats && props.allowedFormats.length > 0) {
    const fileExt = file.name.split('.').pop()?.toLowerCase()
    if (!fileExt || !props.allowedFormats.includes(fileExt)) {
      return `File "${file.name}" harus berformat ${props.allowedFormats.join(', ')}`
    }
  }

  return null
}

const handleFileChange = (event: Event) => {
  errorMessage.value = ''
  const input = event.target as HTMLInputElement
  const selectedFiles = Array.from(input.files || [])

  // Check max files limit
  if (props.maxFiles && localFiles.value.length + selectedFiles.length > props.maxFiles) {
    errorMessage.value = `Maksimum ${props.maxFiles} file dapat diupload`
    input.value = ''
    return
  }

  // Validate each file
  for (const file of selectedFiles) {
    const error = validateFile(file)
    if (error) {
      errorMessage.value = error
      input.value = ''
      return
    }
  }

  const newFiles = [...localFiles.value, ...selectedFiles]
  localFiles.value = newFiles
  emit('update:files', newFiles)
  input.value = ''
}

const removeFile = (index: number) => {
  const newFiles = removeFileFromArray(localFiles.value, index)
  localFiles.value = newFiles
  emit('update:files', newFiles)
}
</script>

<template>
  <div class="file-upload">
    <div class="file-upload-wrapper">
      <input
        :id="inputId"
        type="file"
        :multiple="multiple"
        :accept="accept"
        class="bx--file-input"
        :disabled="disabled"
        @change="handleFileChange"
      />
      <label :for="inputId" class="bx--file-label">
        <span class="bx--file-label-text">Choose files</span>
        <span class="bx--file-label-button">Browse</span>
      </label>
    </div>
    <div v-if="localFiles.length > 0" class="file-list">
      <div v-for="(file, index) in localFiles" :key="index" class="file-item">
        <span class="file-name">{{ file.name }}</span>
        <span class="file-size">({{ formatFileSize(file.size) }})</span>
        <button
          type="button"
          class="file-remove-btn"
          :disabled="disabled"
          @click="removeFile(index)"
        >
          <i class="pi pi-times"></i>
        </button>
      </div>
    </div>
    <small v-if="errorMessage" class="error-text">{{ errorMessage }}</small>
    <small v-else-if="localFiles.length > 0" class="helper-text">
      {{ localFiles.length }} file(s) selected
    </small>
  </div>
</template>

<style scoped>
.file-upload {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.file-upload-wrapper {
  position: relative;
}

.bx--file-input {
  position: absolute;
  width: 0.1px;
  height: 0.1px;
  opacity: 0;
  overflow: hidden;
  z-index: -1;
}

.bx--file-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  padding: 0.875rem 1rem;
  background-color: var(--surface);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  transition: all 0.15s ease;
}

.bx--file-label:hover {
  background-color: #f4f4f4;
  border-color: var(--primary-color);
}

.bx--file-label:active {
  background-color: #e8e8e8;
}

.bx--file-label-text {
  flex: 1;
  color: var(--text-primary);
  font-size: 0.9375rem;
}

.bx--file-label-button {
  padding: 0.5rem 1rem;
  background-color: var(--primary-color);
  color: #ffffff;
  border: 1px solid var(--primary-color);
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 400;
  transition: all 0.11s cubic-bezier(0.2, 0, 0.38, 0.9);
}

.bx--file-label:hover .bx--file-label-button {
  background-color: var(--primary-hover);
  border-color: var(--primary-hover);
}

.bx--file-input:disabled + .bx--file-label {
  opacity: 0.5;
  cursor: not-allowed;
  background-color: var(--surface);
}

.bx--file-input:disabled + .bx--file-label:hover {
  background-color: var(--surface);
  border-color: var(--border-color);
}

.file-list {
  margin-top: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  background-color: #ffffff;
  border: 1px solid var(--border-color);
  border-radius: 4px;
}

.file-name {
  flex: 1;
  color: var(--text-primary);
  font-size: 0.875rem;
}

.file-size {
  color: var(--text-secondary);
  font-size: 0.8125rem;
}

.file-remove-btn {
  background: none;
  border: none;
  color: var(--error-color);
  cursor: pointer;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  transition: background-color 0.15s;
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.file-remove-btn i {
  font-size: 1rem;
}

.file-remove-btn:hover:not(:disabled) {
  background-color: #f4f4f4;
}

.file-remove-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.helper-text {
  font-size: 0.8125rem;
  color: var(--text-secondary);
  margin-top: 0.25rem;
}

.error-text {
  font-size: 0.8125rem;
  color: #da1e28;
  margin-top: 0.25rem;
}
</style>

