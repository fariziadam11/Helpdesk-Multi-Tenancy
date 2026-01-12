<script setup lang="ts">
import { ref, watch } from "vue";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import TicketTable from "@/components/tickets/TicketTable.vue";
import { useTickets } from "@/composables/useTickets";
import Button from "primevue/button";
import InputText from "primevue/inputtext";
import ProgressSpinner from "primevue/progressspinner";

const { t } = useI18n()

const router = useRouter();
const searchQuery = ref("");
const currentPage = ref(1);
const pageLimit = ref(5); // Items per page

const { tickets, pagination, isLoading, error, goToPage, setLimit } =
  useTickets(currentPage, pageLimit);

// Reset to page 1 when search query changes
watch(searchQuery, () => {
  if (currentPage.value !== 1) {
    currentPage.value = 1;
  }
});
</script>

<template>
  <div class="tickets-list-page">
    <div class="page-header">
      <h1>{{ t('tickets.title') }}</h1>
      <Button
        id="ticketsListCreateBtn"
        :label="t('tickets.create')"
        @click="router.push('/tickets/create')"
      />
    </div>

    <div class="page-content">
      <div class="search-section">
        <InputText
          id="ticketsListSearch"
          v-model="searchQuery"
          :placeholder="t('tickets.searchPlaceholder')"
          class="search-input"
        />
      </div>

      <div v-if="isLoading" class="loading-container">
        <ProgressSpinner
          id="ticketsListLoading"
          style="width: 50px; height: 50px"
          strokeWidth="4"
        />
      </div>

      <div v-else-if="error" class="error-container">
        <p>Error loading tickets: {{ error.message }}</p>
      </div>

      <TicketTable
        v-else-if="tickets"
        :tickets="tickets"
        :loading="isLoading"
        :search-query="searchQuery"
        :pagination="pagination"
        :current-page="currentPage"
        :page-limit="pageLimit"
        @page-change="goToPage"
        @limit-change="setLimit"
      />
    </div>
  </div>
</template>

<style scoped>
.tickets-list-page {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-header h1 {
  margin: 0;
  font-size: 2rem;
  font-weight: 600;
}

.page-content {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.search-section {
  margin-bottom: 1rem;
}

.search-input {
  width: 100%;
}

.loading-container,
.error-container {
  padding: 2rem;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
}

.error-container {
  color: var(--error-color);
}

@media (max-width: 768px) {
  .tickets-list-page {
    padding: 1rem;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .page-header h1 {
    font-size: 1.5rem;
  }
}
</style>
