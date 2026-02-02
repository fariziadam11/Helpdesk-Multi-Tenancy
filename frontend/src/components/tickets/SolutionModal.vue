<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'

const { t } = useI18n()

interface Props {
  open: boolean
  mode: 'accept' | 'reject'
  loading?: boolean
}

interface Emits {
  (e: 'close'): void
  (e: 'submit', payload: { rating?: number; comment?: string }): void
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
})

const emit = defineEmits<Emits>()

const solutionRating = ref(5)
const solutionComment = ref('')
const solutionCommentError = ref('')

const ratingOptions = [
  { label: '1', value: 1 },
  { label: '2', value: 2 },
  { label: '3', value: 3 },
  { label: '4', value: 4 },
  { label: '5', value: 5 },
]

const isCommentRequired = computed(() => {
  return props.mode === 'accept' && solutionRating.value < 4
})

const isCommentValid = computed(() => {
  if (props.mode === 'reject') {
    return solutionComment.value.trim().length > 0
  }
  // For accept mode: comment required only if rating < 4
  if (isCommentRequired.value) {
    return solutionComment.value.trim().length > 0
  }
  return true
})

const handleClose = () => {
  solutionComment.value = ''
  solutionCommentError.value = ''
  emit('close')
}

const handleSubmit = () => {
  solutionCommentError.value = ''

  // For reject mode, comment is always required
  if (props.mode === 'reject' && !solutionComment.value.trim()) {
    solutionCommentError.value = t('tickets.detailPage.solution.commentRequired')
    return
  }

  // For accept mode, comment required only if rating < 4
  if (props.mode === 'accept' && isCommentRequired.value && !solutionComment.value.trim()) {
    solutionCommentError.value = 'Comment is required when rating is less than 4'
    return
  }

  const payload: { rating?: number; comment?: string } = {}

  if (props.mode === 'accept') {
    payload.rating = solutionRating.value
    // Only include comment if provided (required when rating < 4)
    if (solutionComment.value.trim()) {
      payload.comment = solutionComment.value.trim()
    }
  } else {
    // Reject mode always requires comment
    payload.comment = solutionComment.value.trim()
  }

  emit('submit', payload)
  solutionComment.value = ''
  solutionCommentError.value = ''
}

const handleCommentInput = () => {
  if (solutionCommentError.value) {
    solutionCommentError.value = ''
  }
}
</script>

<template>
  <Dialog
    id="ticketSolutionModal"
    :visible="open"
    :header="mode === 'accept' ? t('tickets.detailPage.solution.accept') : t('tickets.detailPage.solution.reject')"
    :modal="true"
    :closable="true"
    :style="{ width: '500px' }"
    @update:visible="handleClose"
  >
    <template #header>
      <div class="modal-header-content">
        <h3>{{ mode === 'accept' ? t('tickets.detailPage.solution.accept') : t('tickets.detailPage.solution.reject') }}</h3>
        <p class="ticket-solution-description">
          {{
            mode === 'accept'
              ? 'Please rate the solution. Comment is required if rating is less than 4.'
              : 'Please provide a comment for rejecting this solution.'
          }}
        </p>
      </div>
    </template>

    <div class="solution-form">
      <div
        v-if="mode === 'accept'"
        class="solution-field solution-rating-field"
      >
        <label for="solutionRating" class="solution-label">
          {{ t('tickets.detailPage.solution.rating') }}
        </label>
        <Select
          id="solutionRating"
          v-model="solutionRating"
          :options="ratingOptions"
          optionLabel="label"
          optionValue="value"
          class="form-input"
        />
      </div>

      <div class="solution-field">
        <label for="solutionComment" class="solution-label">
          Comment
          <span v-if="mode === 'reject' || isCommentRequired" class="required-indicator">*</span>
        </label>
        <Textarea
          id="solutionComment"
          v-model="solutionComment"
          :placeholder="t('tickets.detailPage.solution.commentPlaceholder')"
          :rows="4"
          :invalid="!!solutionCommentError"
          @input="handleCommentInput"
          class="form-input"
        />
        <small v-if="solutionCommentError" class="error-text">
          {{ solutionCommentError }}
        </small>
      </div>
    </div>

    <template #footer>
      <div class="modal-footer">
        <Button
          id="ticketSolutionModalCancelBtn"
          :label="t('tickets.detailPage.cancel')"
          severity="secondary"
          :disabled="loading"
          @click="handleClose"
        />
        <Button
          id="ticketSolutionModalSubmitBtn"
          :label="mode === 'accept' ? t('tickets.detailPage.solution.submit') : t('tickets.detailPage.solution.submit')"
          :disabled="loading || !isCommentValid"
          :loading="loading"
          @click="handleSubmit"
        />
      </div>
    </template>
  </Dialog>
</template>

<style scoped>
.modal-header-content h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.25rem;
  font-weight: 600;
}

.solution-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-top: 1rem;
}

.solution-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.solution-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 0.25rem;
  display: block;
}

.required-indicator {
  color: var(--error-color);
}

.ticket-solution-description {
  margin: 0;
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.form-input {
  width: 100%;
}

.error-text {
  color: var(--error-color);
  font-size: 0.8125rem;
  margin-top: 0.25rem;
}

.modal-footer {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
}
</style>
