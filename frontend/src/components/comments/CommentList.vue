<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import CommentItem from './CommentItem.vue'
import FileUpload from '@/components/shared/FileUpload.vue'
import type { Comment } from '@/api/types'
import { useCreateComment } from '@/composables/useComments'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import Textarea from 'primevue/textarea'

const { t } = useI18n()

interface Props {
  comments: Comment[]
  ticketId: number
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
})

const visibleComments = computed(() => {
  return props.comments.filter((comment) => {
    return comment.customer_visible !== false
  })
})

const solutionComments = computed(() => {
  return visibleComments.value.filter((comment) => comment.is_solution === true)
})

const regularComments = computed(() => {
  return visibleComments.value.filter((comment) => comment.is_solution !== true)
})

const showModal = ref(false)
const message = ref('')
const files = ref<File[]>([])

const { mutate: createComment, isPending } = useCreateComment(props.ticketId)

const handleSubmit = () => {
  if (!message.value.trim()) return

  createComment(
    {
      message: message.value,
      attachments: files.value.length > 0 ? files.value : undefined,
    },
    {
      onSuccess: () => {
        message.value = ''
        files.value = []
        showModal.value = false
      },
    }
  )
}

const handleModalClose = () => {
  showModal.value = false
}
</script>

<template>
  <div class="comment-list">
    <div class="comment-list-header">
      <h3>
        Comments ({{ visibleComments.length }})
        <span v-if="solutionComments.length" class="solution-count">
          • Solution {{ solutionComments.length }}
        </span>
      </h3>
      <Button id="commentListAddBtn" :label="t('tickets.detailPage.addComment')" @click="showModal = true" />
    </div>

    <div v-if="loading" class="loading">{{ t('common.loading') }}</div>
    <div v-else-if="visibleComments.length === 0" class="empty">
      {{ t('tickets.detailPage.noComments') }}
    </div>
    <div v-else>
      <div
        v-if="solutionComments.length > 0"
        class="comments solution-comments"
      >
        <h4 class="section-title">Solution</h4>
        <CommentItem
          v-for="comment in solutionComments"
          :key="comment.id"
          :comment="comment"
        />
      </div>

      <div
        v-if="regularComments.length > 0"
        class="comments regular-comments"
      >
        <h4 class="section-title">
          {{ solutionComments.length ? 'Other comments' : 'Comments' }}
        </h4>
        <CommentItem
          v-for="comment in regularComments"
          :key="comment.id"
          :comment="comment"
        />
      </div>
    </div>

    <Dialog
      id="commentListModal"
      :visible="showModal"
      :header="t('tickets.detailPage.addComment')"
      :modal="true"
      :closable="true"
      :style="{ width: '600px' }"
      @update:visible="handleModalClose"
    >
      <div class="comment-form">
        <div class="form-item">
          <label for="commentListModalTextarea" class="form-label">
            Message
          </label>
          <Textarea
            id="commentListModalTextarea"
            v-model="message"
            :placeholder="t('tickets.detailPage.writeComment')"
            :rows="5"
            class="form-input"
          />
        </div>
        <div class="form-item">
          <label for="commentListModalFileInput" class="form-label">
            Attachments <span class="optional-text">(optional)</span>
          </label>
          <FileUpload
            id="commentListModalFileInput"
            v-model:files="files"
            :disabled="isPending"
          />
        </div>
      </div>

      <template #footer>
        <div class="modal-footer">
          <Button
            id="commentListModalCancelBtn"
            :label="t('tickets.detailPage.cancel')"
            severity="secondary"
            @click="handleModalClose"
          />
          <Button
            id="commentListModalSubmitBtn"
            :label="isPending ? t('common.loading') : t('tickets.detailPage.postComment')"
            :disabled="!message.trim() || isPending"
            :loading="isPending"
            @click="handleSubmit"
          />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<style scoped>
.comment-list {
  margin-top: 2rem;
}

.comment-list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.comments {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.solution-comments {
  margin-bottom: 1.5rem;
}

.section-title {
  margin: 0 0 0.5rem;
  font-size: 0.9rem;
  font-weight: 600;
}

.solution-count {
  font-size: 0.8rem;
  font-weight: 400;
  color: var(--text-secondary);
  margin-left: 0.25rem;
}

.loading,
.empty {
  padding: 2rem;
  text-align: center;
  color: var(--text-secondary);
}

.comment-form {
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

.optional-text {
  color: var(--text-secondary);
  font-weight: 400;
  font-size: 0.8125rem;
}

.form-input {
  width: 100%;
}

.modal-footer {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
}

@media (max-width: 768px) {
  .comment-list-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .comment-list-header :deep(.p-button) {
    width: 100%;
  }

  .modal-footer {
    flex-direction: column-reverse;
  }

  .modal-footer :deep(.p-button) {
    width: 100%;
  }
}
</style>
