<script setup lang="ts">
import { useI18n } from "vue-i18n";
import type { Ticket } from "@/api/types";
import { formatUnixTimestamp } from "@/utils/date";
import Tag from "primevue/tag";

const { t } = useI18n();

interface Props {
  ticket: Ticket;
}

defineProps<Props>();

const getStatusSeverity = (status: string) => {
  const statusLower = status?.toLowerCase() || "";
  if (statusLower.includes("pending")) return "info";
  if (statusLower.includes("solved") || statusLower.includes("closed"))
    return "success";
  if (statusLower.includes("in progress")) return "warn";
  return "secondary";
};
</script>

<template>
  <div id="ticketDetailCardTile" class="ticket-detail-card">
    <h3>{{ t('tickets.detailPage.ticketDetails') }}</h3>
    <div class="detail-list">
      <div class="detail-row">
        <div class="detail-label">{{ t('tickets.detailPage.ticketId') }}</div>
        <div class="detail-value">
          {{ ticket.pretty_id || `#${ticket.id}` }}
        </div>
      </div>
      <div class="detail-row">
        <div class="detail-label">{{ t('tickets.detailPage.title') }}</div>
        <div class="detail-value">{{ ticket.title }}</div>
      </div>
      <div class="detail-row">
        <div class="detail-label">{{ t('tickets.detailPage.status') }}</div>
        <div class="detail-value">
          <Tag
            id="ticketDetailCardStatusTag"
            :severity="getStatusSeverity(ticket.status)"
          >
            {{ ticket.status }}
          </Tag>
        </div>
      </div>
      <div class="detail-row">
        <div class="detail-label">{{ t('tickets.detailPage.description') }}</div>
        <div class="detail-value">
          {{ ticket.description || "N/A" }}
        </div>
      </div>
      <div class="detail-row">
        <div class="detail-label">{{ t('tickets.detailPage.createdAt') }}</div>
        <div class="detail-value">
          {{ formatUnixTimestamp(ticket.created_at) }}
        </div>
      </div>
      <div v-if="ticket.last_update" class="detail-row">
        <div class="detail-label">{{ t('tickets.detailPage.lastUpdate') }}</div>
        <div class="detail-value">
          {{ formatUnixTimestamp(ticket.last_update) }}
        </div>
      </div>
      <div v-if="ticket.closed_at" class="detail-row">
        <div class="detail-label">Closed At</div>
        <div class="detail-value">
          {{ formatUnixTimestamp(ticket.closed_at) }}
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.ticket-detail-card {
  background-color: #ffffff;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 1.5rem;
}

.ticket-detail-card h3 {
  margin: 0 0 1.5rem 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
}

.detail-list {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.detail-row {
  display: grid;
  grid-template-columns: 150px 1fr;
  gap: 1rem;
  padding: 1rem 0;
  border-bottom: 1px solid var(--border-color);
}

.detail-row:last-child {
  border-bottom: none;
}

.detail-label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-secondary);
}

.detail-value {
  font-size: 0.875rem;
  color: var(--text-primary);
  word-break: break-word;
}

@media (max-width: 768px) {
  .detail-row {
    grid-template-columns: 1fr;
    gap: 0.5rem;
    padding: 0.875rem 0;
  }

  .detail-label {
    font-size: 0.8125rem;
  }

  .detail-value {
    font-size: 0.8125rem;
  }
}
</style>
