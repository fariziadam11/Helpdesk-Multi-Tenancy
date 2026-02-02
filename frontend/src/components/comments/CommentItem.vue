<script setup lang="ts">
import { useI18n } from "vue-i18n";
import type { Comment, Attachment } from "@/api/types";
import { formatUnixTimestamp } from "@/utils/date";
import AttachmentPreview from "../attachments/AttachmentPreview.vue";
import { useUserName } from "@/composables/useUserName";

const { t } = useI18n();

interface Props {
  comment: Comment;
}

const props = defineProps<Props>();

const { displayName, isAgent } = useUserName(props.comment.author_id);



const getAttachments = (comment: Comment): Array<number | Attachment> => {
  const attachments = comment.attached_files || comment.attachments;
  if (!attachments || !Array.isArray(attachments)) return [];

  return attachments.map((item) => {
    if (typeof item === "number") {
      return item;
    }
    return item as Attachment;
  });
};
</script>

<template>
  <div :id="`commentItemTile-${comment.id || 'unknown'}`" class="comment-card">
    <div class="comment-item">
      <div class="comment-header">
        <span class="comment-author">
          <span class="comment-author-icon" aria-hidden="true">
            <i v-if="isAgent" class="pi pi-id-card"></i>
            <i v-else class="pi pi-user"></i>
          </span>
          <span class="comment-author-name">
            {{
              comment.author ||
              displayName ||
              `User #${comment.author_id || "Unknown"}`
            }}
          </span>
        </span>
        <span class="comment-date">
          {{ formatUnixTimestamp(comment.created_at) }}
        </span>
      </div>
      <div class="comment-body" v-html="comment.message || comment.comment || t('tickets.detailPage.noComments')"></div>
      <div
        v-if="getAttachments(comment).length > 0"
        class="comment-attachments"
      >
        <AttachmentPreview
          v-for="(attachment, index) in getAttachments(comment)"
          :key="
            typeof attachment === 'number' ? attachment : attachment.id || index
          "
          :attachment-id="
            typeof attachment === 'number' ? attachment : attachment.id
          "
          :attachment="typeof attachment === 'number' ? undefined : attachment"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.comment-card {
  background-color: #ffffff;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1rem;
}

.comment-item {
  padding: 0.5rem;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
}

.comment-author {
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
}

.comment-author-icon {
  display: inline-flex;
  align-items: center;
}

.comment-author-icon i {
  font-size: 1rem;
}

.comment-date {
  color: var(--text-secondary);
}

.comment-body {
  margin-bottom: 0.5rem;
  line-height: 1.6;
}

.comment-body :deep(p) {
  margin: 0 0 0.75rem 0;
}

.comment-body :deep(p:last-child) {
  margin-bottom: 0;
}

.comment-body :deep(br) {
  display: block;
  content: "";
  margin-top: 0.5rem;
}

.comment-body :deep(a) {
  color: #6929C4;
  text-decoration: underline;
  cursor: pointer;
}

.comment-body :deep(a:hover) {
  color: #4F2196;
  text-decoration: none;
}

.comment-body :deep(strong),
.comment-body :deep(b) {
  font-weight: 600;
}

.comment-body :deep(em),
.comment-body :deep(i) {
  font-style: italic;
}

.comment-body :deep(ul),
.comment-body :deep(ol) {
  margin: 0.5rem 0;
  padding-left: 1.5rem;
}

.comment-body :deep(li) {
  margin: 0.25rem 0;
}

.comment-attachments {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

@media (max-width: 768px) {
  .comment-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.25rem;
  }

  .comment-author-name {
    word-break: break-word;
  }
}
</style>
