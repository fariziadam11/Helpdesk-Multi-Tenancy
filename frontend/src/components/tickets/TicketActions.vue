<script setup lang="ts">
import { useI18n } from "vue-i18n";
import type { Ticket } from "@/api/types";
import Button from "primevue/button";

const { t } = useI18n();

interface Props {
  ticket: Ticket;
  isUpdatePending?: boolean;
}

interface Emits {
  (e: "update"): void;
  (e: "accept-solution"): void;
  (e: "reject-solution"): void;
}

defineProps<Props>();
const emit = defineEmits<Emits>();
</script>

<template>
  <div class="ticket-header-actions">
    <div class="ticket-header-actions-right">
      <Button
        id="ticketDetailUpdateBtn"
        :label="t('tickets.detailPage.actions.updateTicket')"
        severity="secondary"
        size="small"
        :disabled="isUpdatePending"
        @click="emit('update')"
      />
      <Button
        v-if="ticket.status === 'Resolved'"
        id="ticketDetailAcceptSolutionBtn"
        :label="t('tickets.detailPage.actions.acceptSolution')"
        size="small"
        @click="emit('accept-solution')"
      />
      <Button
        v-if="ticket.status === 'Resolved'"
        id="ticketDetailRejectSolutionBtn"
        :label="t('tickets.detailPage.actions.rejectSolution')"
        severity="secondary"
        outlined
        size="small"
        @click="emit('reject-solution')"
      />
    </div>
  </div>
</template>

<style scoped>
.ticket-header-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 1.5rem;
  margin-bottom: 0.5rem;
}

.ticket-header-actions-right {
  display: flex;
  gap: 0.5rem;
}
</style>
