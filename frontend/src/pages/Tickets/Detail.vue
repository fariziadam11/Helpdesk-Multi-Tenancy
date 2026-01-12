<script setup lang="ts">
import { computed, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import TicketDetailCard from "@/components/tickets/TicketDetailCard.vue";
import TicketActions from "@/components/tickets/TicketActions.vue";
import CommentList from "@/components/comments/CommentList.vue";
import AttachmentPreview from "@/components/attachments/AttachmentPreview.vue";
import SolutionModal from "@/components/tickets/SolutionModal.vue";
import UpdateTicketModal from "@/components/tickets/UpdateTicketModal.vue";
import { useTicketDetail } from "@/composables/useTicketDetail";
import { useComments } from "@/composables/useComments";
import { useTicketSolution } from "@/composables/useTicketSolution";
import { useUpdateTicket } from "@/composables/useUpdateTicket";
import { useCategories } from "@/composables/useCategories";
import { useTicketMeta } from "@/composables/useTicketMeta";
import type { Attachment } from "@/api/types";
import Breadcrumb from "primevue/breadcrumb";
import ProgressSpinner from "primevue/progressspinner";
import Button from "primevue/button";

const { t } = useI18n();

interface Props {
  id: string;
}

const props = defineProps<Props>();
const route = useRoute();
const router = useRouter();

const ticketId = computed(() =>
  parseInt(props.id || String(route.params.id), 10)
);

const {
  data: ticket,
  isLoading: ticketLoading,
  error: ticketError,
} = useTicketDetail(ticketId.value);
const { data: comments, isLoading: commentsLoading } = useComments(
  ticketId.value
);

const { acceptMutation, rejectMutation } = useTicketSolution(ticketId.value);
const { mutate: updateTicket, isPending: isUpdatePending } = useUpdateTicket(
  ticketId.value
);

const { data: categories } = useCategories();
const { data: ticketMeta } = useTicketMeta();

const showSolutionModal = ref(false);
const showUpdateModal = ref(false);
const solutionMode = ref<"accept" | "reject">("accept");

const isSolutionLoading = computed(
  () => acceptMutation.isPending.value || rejectMutation.isPending.value
);

const openAcceptSolution = () => {
  solutionMode.value = "accept";
  showSolutionModal.value = true;
};

const openRejectSolution = () => {
  solutionMode.value = "reject";
  showSolutionModal.value = true;
};

const handleSolutionSubmit = async (payload: {
  rating?: number;
  comment?: string;
}) => {
  try {
    if (solutionMode.value === "accept") {
      await acceptMutation.mutateAsync({
        rating: payload.rating || 5,
        ...(payload.comment && { comment: payload.comment }),
      });
    } else {
      if (!payload.comment) {
        throw new Error("Comment is required for rejecting solution");
      }
      await rejectMutation.mutateAsync({
        comment: payload.comment,
      });
    }
    showSolutionModal.value = false;
  } catch (error) {
    // Error handled by mutation
  }
};

const openUpdateModal = () => {
  showUpdateModal.value = true;
};

const handleUpdateSubmit = (payload: Record<string, any>) => {
  updateTicket(payload, {
    onSuccess: () => {
      showUpdateModal.value = false;
    },
  });
};

const getAttachments = (): Array<number | Attachment> => {
  if (!ticket.value?.attachments) return [];
  if (!Array.isArray(ticket.value.attachments)) return [];

  return ticket.value.attachments.map((item) => {
    if (typeof item === "number") {
      return item;
    }
    return item as Attachment;
  });
};

const breadcrumbItems = computed(() => [
  { label: t('tickets.list'), command: () => router.push("/tickets") },
  { label: ticket.value?.pretty_id || `#${ticketId.value}` },
]);
</script>

<template>
  <div class="ticket-detail-page">
    <Breadcrumb
      id="ticketDetailBreadcrumb"
      :model="breadcrumbItems"
      class="breadcrumb"
    />

    <div v-if="ticketLoading" class="loading-container">
      <ProgressSpinner
        id="ticketDetailLoading"
        style="width: 50px; height: 50px"
        strokeWidth="4"
      />
    </div>

    <div v-else-if="ticketError" class="error-container">
      <p>Error loading ticket: {{ ticketError.message }}</p>
      <Button
        id="ticketDetailBackBtn"
        :label="t('tickets.detailPage.backToList')"
        @click="router.push('/tickets')"
      />
    </div>

    <div v-else-if="ticket" class="ticket-content">
      <TicketActions
        :ticket="ticket"
        :is-update-pending="isUpdatePending"
        @update="openUpdateModal"
        @accept-solution="openAcceptSolution"
        @reject-solution="openRejectSolution"
      />

      <div class="content-grid">
        <div class="main-content">
          <TicketDetailCard :ticket="ticket" />
        </div>
        <div v-if="getAttachments().length > 0" class="sidebar-content">
          <div class="attachments-section">
            <h3>{{ t('tickets.detailPage.attachments') }}</h3>
            <div class="attachments-list">
              <AttachmentPreview
                v-for="(attachment, index) in getAttachments()"
                :key="
                  typeof attachment === 'number'
                    ? attachment
                    : attachment.id || index
                "
                :attachment-id="
                  typeof attachment === 'number' ? attachment : attachment.id
                "
                :attachment="
                  typeof attachment === 'number' ? undefined : attachment
                "
              />
            </div>
          </div>
        </div>
      </div>
      <div class="comments-section">
        <CommentList
          :comments="comments || []"
          :ticket-id="ticketId"
          :loading="commentsLoading"
        />
      </div>
    </div>

    <SolutionModal
      :open="showSolutionModal"
      :mode="solutionMode"
      :loading="isSolutionLoading"
      @close="showSolutionModal = false"
      @submit="handleSolutionSubmit"
    />

    <UpdateTicketModal
      :open="showUpdateModal"
      :ticket="ticket || null"
      :categories="categories"
      :types="ticketMeta?.types"
      :priorities="ticketMeta?.priorities"
      :loading="isUpdatePending"
      @close="showUpdateModal = false"
      @submit="handleUpdateSubmit"
    />
  </div>
</template>

<style scoped>
.ticket-detail-page {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.breadcrumb {
  margin-bottom: 1rem;
}

.loading-container,
.error-container {
  padding: 2rem;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
}

.error-container {
  color: var(--error-color);
}

.ticket-content {
  margin-top: 2rem;
}

.content-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 2rem;
  margin-top: 1rem;
}

.main-content {
  flex: 1;
}

.sidebar-content {
  flex-shrink: 0;
}

@media (min-width: 1024px) {
  .content-grid {
    grid-template-columns: 2fr 1fr;
  }
}

.comments-section {
  margin-top: 2rem;
}

.attachments-section {
  margin-top: 1rem;
}

.attachments-section h3 {
  margin-bottom: 1rem;
}

.attachments-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

@media (max-width: 768px) {
  .ticket-detail-page {
    padding: 1rem;
  }
}
</style>
