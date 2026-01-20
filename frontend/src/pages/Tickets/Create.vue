<script setup lang="ts">
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import { useCreateTicket } from "@/composables/useCreateTicket";
import { useCategories } from "@/composables/useCategories";
import { useTicketMeta } from "@/composables/useTicketMeta";
import { useFormValidation } from "@/composables/useFormValidation";
import FileUpload from "@/components/shared/FileUpload.vue";
import { TICKET_DEFAULTS } from "@/utils/constants";
import Breadcrumb from "primevue/breadcrumb";
import InputText from "primevue/inputtext";
import Textarea from "primevue/textarea";
import Select from "primevue/select";
import Button from "primevue/button";
import ProgressSpinner from "primevue/progressspinner";
import { onMounted } from "vue";

const { t } = useI18n();
const router = useRouter();

const title = ref("");
const description = ref("");
const categoryId = ref<string>("");
const typeId = ref<string>("");
const priorityId = ref<string>("");
const dateOccurred = ref<string>("");
const files = ref<File[]>([]);

const now = new Date();
const yyyy = now.getFullYear();
const mm = String(now.getMonth() + 1).padStart(2, "0");
const dd = String(now.getDate()).padStart(2, "0");
const today = `${yyyy}-${mm}-${dd}`;

if (today) {
  dateOccurred.value = today;
}

const { data: categories, isLoading: categoriesLoading } = useCategories();
const { data: ticketMeta, isLoading: metaLoading } = useTicketMeta();
const { mutate: createTicket, isPending, error } = useCreateTicket();

const categoryOptions = computed(() => categories.value || []);
const typeOptions = computed(() => ticketMeta.value?.types || []);
const priorityOptions = computed(() => ticketMeta.value?.priorities || []);

// Auto-select "Problem" type on mount
onMounted(() => {
  if (ticketMeta.value?.types) {
    const problemType = ticketMeta.value.types.find(t => t.name.toLowerCase() === 'problem')
    if (problemType) {
      typeId.value = String(problemType.id)
    }
  }
})

const validation = useFormValidation({
  fields: {
    title: { value: title, required: true },
    description: { value: description, required: true },
    categoryId: { value: categoryId, required: true },
    priorityId: { value: priorityId, required: true },
    dateOccurred: { value: dateOccurred, required: true },
  },
});

const breadcrumbItems = computed(() => [
  { label: t('tickets.list'), command: () => router.push("/tickets") },
  { label: t('tickets.create') },
]);

const handleSubmit = () => {
  validation.markAllTouched();

  if (!validation.isFormValid.value) {
    return;
  }

  const dateTimestamp = dateOccurred.value
    ? Math.floor(new Date(dateOccurred.value).getTime() / 1000)
    : Math.floor(Date.now() / 1000);

  createTicket({
    source_id: TICKET_DEFAULTS.SOURCE_ID,
    category_id: parseInt(categoryId.value, 10),
    type_id: parseInt(typeId.value, 10),
    priority_id: parseInt(priorityId.value, 10),
    title: title.value,
    description: description.value,
    date_ocurred: dateTimestamp,
    attachments: files.value.length > 0 ? files.value : undefined,
  });
};
</script>

<template>
  <div class="create-ticket-page">
    <Breadcrumb
      id="createTicketBreadcrumb"
      :model="breadcrumbItems"
      class="breadcrumb"
    />

    <div class="page-header">
      <h1>{{ t('tickets.create') }}</h1>
    </div>

    <form @submit.prevent="handleSubmit" class="ticket-form">
      <div class="form-item">
        <label for="createTicketTitleInput" class="form-label">
          {{ t('tickets.form.title') }} <span class="required">*</span>
        </label>
        <InputText
          id="createTicketTitleInput"
          v-model="title"
          :placeholder="t('tickets.form.titlePlaceholder')"
          :disabled="isPending"
          :invalid="validation.isFieldInvalid('title')"
          @blur="validation.markFieldTouched('title')"
          class="form-input"
        />
        <small v-if="validation.isFieldInvalid('title')" class="error-text">
          {{ validation.getFieldError("title") }}
        </small>
      </div>

      <div class="form-item">
        <label for="createTicketDescriptionTextarea" class="form-label">
          {{ t('tickets.form.description') }} <span class="required">*</span>
        </label>
        <Textarea
          id="createTicketDescriptionTextarea"
          v-model="description"
          :placeholder="t('tickets.form.descriptionPlaceholder')"
          :rows="6"
          :disabled="isPending"
          :invalid="validation.isFieldInvalid('description')"
          @blur="validation.markFieldTouched('description')"
          class="form-input"
        />
        <small
          v-if="validation.isFieldInvalid('description')"
          class="error-text"
        >
          {{ validation.getFieldError("description") }}
        </small>
      </div>

      <div class="form-item">
        <label for="createTicketCategorySelect" class="form-label">
          {{ t('tickets.form.category') }} <span class="required">*</span>
        </label>
        <Select
          id="createTicketCategorySelect"
          v-model="categoryId"
          :options="categoryOptions"
          optionLabel="name"
          optionValue="id"
          :placeholder="t('tickets.form.selectCategory')"
          :disabled="categoriesLoading || isPending"
          :invalid="validation.isFieldInvalid('categoryId')"
          @blur="validation.markFieldTouched('categoryId')"
          class="form-input"
        />
        <ProgressSpinner
          v-if="categoriesLoading"
          style="width: 20px; height: 20px; margin-top: 0.5rem"
          strokeWidth="4"
        />
        <small
          v-if="validation.isFieldInvalid('categoryId')"
          class="error-text"
        >
          {{ validation.getFieldError("categoryId") }}
        </small>
      </div>

      <div class="form-item">
        <label for="createTicketPrioritySelect" class="form-label">
          {{ t('tickets.form.priority') }} <span class="required">*</span>
        </label>
        <Select
          id="createTicketPrioritySelect"
          v-model="priorityId"
          :options="priorityOptions"
          optionLabel="name"
          optionValue="id"
          :placeholder="t('tickets.form.selectPriority')"
          :disabled="metaLoading || isPending"
          :invalid="validation.isFieldInvalid('priorityId')"
          @blur="validation.markFieldTouched('priorityId')"
          class="form-input"
        />
        <ProgressSpinner
          v-if="metaLoading"
          style="width: 20px; height: 20px; margin-top: 0.5rem"
          strokeWidth="4"
        />
        <small
          v-if="validation.isFieldInvalid('priorityId')"
          class="error-text"
        >
          {{ validation.getFieldError("priorityId") }}
        </small>
      </div>

      <div class="form-item">
        <label for="createTicketDateOccurred" class="form-label">
          {{ t('tickets.form.dateOccurred') }} <span class="required">*</span>
        </label>
        <div class="date-input-wrapper">
          <input
            id="createTicketDateOccurred"
            v-model="dateOccurred"
            type="date"
            class="date-input"
            :class="{
              'date-input-invalid': validation.isFieldInvalid('dateOccurred'),
            }"
            :required="true"
            :disabled="isPending"
            :max="today"
            @blur="validation.markFieldTouched('dateOccurred')"
          />
        </div>
        <small
          v-if="validation.isFieldInvalid('dateOccurred')"
          class="error-text"
        >
          {{ validation.getFieldError("dateOccurred") }}
        </small>
        <small v-else class="helper-text">
          {{ t('tickets.form.dateOccurredHelper') }}
        </small>
      </div>

      <div class="form-item">
        <label for="createTicketFileInput" class="form-label">
          {{ t('tickets.form.screenshot') }} <span class="optional-text">({{ t('tickets.form.optional') }})</span>
        </label>
        <FileUpload 
          v-model:files="files" 
          :disabled="isPending"
          accept="image/jpeg,image/jpg,image/png"
          :allowed-formats="['jpg', 'jpeg', 'png']"
          :max-file-size="2 * 1024 * 1024"
          :max-files="4"
        />
        <small class="helper-text">Format: JPG, JPEG, PNG • Maksimum 2MB per file • Maksimum 4 file</small>
      </div>

      <div v-if="error" class="error-message">{{ t('tickets.form.error') }}: {{ error.message }}</div>

      <div class="form-actions">
        <Button
          id="createTicketCancelBtn"
          :label="t('tickets.form.cancel')"
          severity="secondary"
          @click="router.push('/tickets')"
        />
        <Button
          id="createTicketSubmitBtn"
          type="submit"
          :label="isPending ? t('tickets.form.creating') : t('tickets.form.submit')"
          :disabled="!validation.isFormValid || isPending"
          :loading="isPending"
        />
      </div>
    </form>
  </div>
</template>

<style scoped>
.create-ticket-page {
  padding: 2rem;
  max-width: 800px;
  margin: 0 auto;
}

.breadcrumb {
  margin-bottom: 1rem;
}

.page-header {
  margin: 2rem 0;
}

.page-header h1 {
  margin: 0;
  font-size: 2rem;
  font-weight: 600;
}

.ticket-form {
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
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 0.375rem;
}

.optional-text {
  color: var(--text-secondary);
  font-weight: 400;
  font-size: 0.8125rem;
}

.required {
  color: var(--error-color);
}

.helper-text {
  font-size: 0.8125rem;
  color: var(--text-secondary);
  margin-top: 0.25rem;
}

.error-text {
  color: var(--error-color);
  font-size: 0.8125rem;
  margin-top: 0.25rem;
}

.form-input {
  width: 100%;
}

.date-input-wrapper {
  position: relative;
  width: 100%;
}

.date-input {
  width: 100%;
  padding: 0.875rem 1rem;
  font-size: 0.9375rem;
  background-color: var(--surface);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  color: var(--text-primary);
  transition: all 0.15s ease;
  font-family: inherit;
  line-height: 1.5;
  min-height: 2.5rem;
}

.date-input:focus {
  outline: 2px solid var(--primary-color);
  outline-offset: -2px;
  border-color: var(--primary-color);
  background-color: var(--surface);
}

.date-input:disabled {
  background-color: var(--surface);
  color: #8d8d8d;
  cursor: not-allowed;
  opacity: 0.6;
  border-color: var(--border-color);
}

.date-input-invalid {
  border-color: var(--error-color);
  box-shadow: 0 0 0 1px var(--error-color);
}

.date-input-invalid:focus {
  outline-color: var(--error-color);
  border-color: var(--error-color);
  box-shadow: 0 0 0 2px var(--error-color);
}

.date-input::-webkit-calendar-picker-indicator {
  cursor: pointer;
  opacity: 1;
  filter: invert(0.5);
}

.date-input::-webkit-calendar-picker-indicator:hover {
  opacity: 0.8;
}

.date-input:disabled::-webkit-calendar-picker-indicator {
  cursor: not-allowed;
  opacity: 0.3;
}

.error-message {
  padding: 1rem;
  background-color: #fef0f0;
  color: var(--error-color);
  border-radius: 4px;
  border: 1px solid var(--error-color);
  margin-bottom: 1rem;
}

.form-actions {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border-color);
}

@media (max-width: 768px) {
  .create-ticket-page {
    padding: 1rem;
  }

  .page-header {
    margin: 1rem 0;
  }

  .page-header h1 {
    font-size: 1.5rem;
  }

  .form-actions {
    flex-direction: column-reverse;
    gap: 0.5rem;
    width: 100%;
  }

  .form-actions :deep(.p-button) {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .form-actions {
    gap: 0.5rem;
    width: 100%;
  }

  .form-actions :deep(.p-button) {
    width: 100%;
    min-height: 2.5rem;
    padding: 0.75rem 1rem;
    font-size: 0.875rem;
  }
}
</style>
