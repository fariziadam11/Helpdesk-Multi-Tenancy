<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Ticket, Category } from '@/api/types'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Select from 'primevue/select'

const { t } = useI18n()

interface Props {
  open: boolean
  ticket: Ticket | null
  categories?: Category[]
  types?: Array<{ id: number; name: string }>
  priorities?: Array<{ id: number; name: string }>
  loading?: boolean
}

interface Emits {
  (e: 'close'): void
  (e: 'submit', payload: Record<string, any>): void
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
})

const emit = defineEmits<Emits>()

const updateTitle = ref('')
const updateDescription = ref('')
const updateCategoryId = ref<number | null>(null)
const updateTypeId = ref<number | null>(null)
const updatePriorityId = ref<number | null>(null)

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen && props.ticket) {
      updateTitle.value = props.ticket.title || ''
      updateDescription.value = props.ticket.description || ''
      updateCategoryId.value = props.ticket.category_id || null
      updateTypeId.value = props.ticket.type_id || null
      updatePriorityId.value = props.ticket.priority_id || null
    }
  },
  { immediate: true },
)

const handleClose = () => {
  emit('close')
}

const handleSubmit = () => {
  if (!props.ticket) return

  const payload: Record<string, any> = {}

  if (updateTitle.value.trim() !== props.ticket.title) {
    payload.title = updateTitle.value.trim()
  }
  if (updateDescription.value.trim() !== props.ticket.description) {
    payload.description = updateDescription.value.trim()
  }
  if (
    updateCategoryId.value &&
    updateCategoryId.value !== props.ticket.category_id
  ) {
    payload.category_id = updateCategoryId.value
  }
  if (
    updateTypeId.value &&
    updateTypeId.value !== props.ticket.type_id
  ) {
    payload.type_id = updateTypeId.value
  }
  if (
    updatePriorityId.value &&
    updatePriorityId.value !== props.ticket.priority_id
  ) {
    payload.priority_id = updatePriorityId.value
  }

  if (Object.keys(payload).length === 0) {
    return
  }

  emit('submit', payload)
}
</script>

<template>
  <Dialog
    id="ticketUpdateModal"
    :visible="open"
    header="Update Ticket"
    :modal="true"
    :closable="true"
    :style="{ width: '600px' }"
    @update:visible="handleClose"
  >
    <template #header>
      <div class="modal-header-content">
      <h3>{{ t('tickets.detailPage.updateTicket') }}</h3>
        <p class="ticket-update-description">
          Update ticket information. Only changed fields will be updated.
        </p>
      </div>
    </template>

    <div class="update-form">
      <div class="form-item">
        <label for="updateTicketTitle" class="form-label">
          {{ t('tickets.form.title') }} <span class="required">*</span>
        </label>
        <InputText
          id="updateTicketTitle"
          v-model="updateTitle"
          :placeholder="t('tickets.form.titlePlaceholder')"
          :disabled="loading"
          class="form-input"
        />
      </div>

      <div class="form-item">
        <label for="updateTicketDescription" class="form-label">
          {{ t('tickets.form.description') }} <span class="required">*</span>
        </label>
        <Textarea
          id="updateTicketDescription"
          v-model="updateDescription"
          :placeholder="t('tickets.form.descriptionPlaceholder')"
          :rows="6"
          :disabled="loading"
          class="form-input"
        />
      </div>

      <div class="form-item">
        <label for="updateTicketCategory" class="form-label">
          {{ t('tickets.form.category') }} <span class="required">*</span>
        </label>
        <Select
          id="updateTicketCategory"
          v-model="updateCategoryId"
          :options="categories"
          optionLabel="name"
          optionValue="id"
          :placeholder="t('tickets.form.selectCategory')"
          :disabled="loading"
          class="form-input"
        />
      </div>

      <div class="form-item">
        <label for="updateTicketType" class="form-label">
          Type <span class="required">*</span>
        </label>
        <Select
          id="updateTicketType"
          v-model="updateTypeId"
          :options="types"
          optionLabel="name"
          optionValue="id"
          placeholder="Select type"
          :disabled="loading"
          class="form-input"
        />
      </div>

      <div class="form-item">
        <label for="updateTicketPriority" class="form-label">
          Priority <span class="required">*</span>
        </label>
        <Select
          id="updateTicketPriority"
          v-model="updatePriorityId"
          :options="priorities"
          optionLabel="name"
          optionValue="id"
          placeholder="Select priority"
          :disabled="loading"
          class="form-input"
        />
      </div>
    </div>

    <template #footer>
      <div class="modal-footer">
        <Button
          id="ticketUpdateModalCancelBtn"
          :label="t('tickets.detailPage.cancel')"
          severity="secondary"
          :disabled="loading"
          @click="handleClose"
        />
        <Button
          id="ticketUpdateModalSubmitBtn"
          :label="t('tickets.detailPage.update')"
          :disabled="loading"
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

.ticket-update-description {
  margin: 0;
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.update-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-top: 1rem;
}

.form-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-label {
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

.modal-footer {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
}
</style>
