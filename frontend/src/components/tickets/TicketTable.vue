<script setup lang="ts">
import { computed } from "vue";
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import type { Ticket, PaginationMeta } from "@/api/types";
import { formatUnixTimestamp } from "@/utils/date";
import Tag from "primevue/tag";

const { t } = useI18n();

interface Props {
  tickets: Ticket[];
  loading?: boolean;
  searchQuery?: string;
  pagination?: PaginationMeta;
  currentPage?: number;
  pageLimit?: number;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  searchQuery: "",
  currentPage: 1,
  pageLimit: 5, // Items per page
});

interface Emits {
  (e: "page-change", page: number): void;
  (e: "limit-change", limit: number): void;
}

const emit = defineEmits<Emits>();

const router = useRouter();

// Client-side filtering for search (since backend doesn't support search yet)
const filteredTickets = computed(() => {
  if (!props.searchQuery) return props.tickets;
  const query = props.searchQuery.toLowerCase();
  return props.tickets.filter(
    (ticket) =>
      ticket.title?.toLowerCase().includes(query) ||
      ticket.wrk_ticket_id?.toLowerCase().includes(query) ||
      ticket.inv_gate_id?.toString().toLowerCase().includes(query) ||
      ticket.status?.toLowerCase().includes(query)
  );
});

// Use server-side pagination metadata if available
const totalItems = computed(
  () => props.pagination?.total ?? filteredTickets.value.length
);

const totalPages = computed(
  () =>
    props.pagination?.total_pages ??
    Math.ceil(filteredTickets.value.length / props.pageLimit)
);

const handlePageChange = (newPage: number) => {
  if (
    newPage >= 1 &&
    newPage <= totalPages.value &&
    newPage !== props.currentPage
  ) {
    emit("page-change", newPage);
  }
};

const handleLimitChange = (event: Event) => {
  const target = event.target as HTMLSelectElement;
  const newLimit = Number(target.value);
  if (newLimit > 0 && newLimit !== props.pageLimit) {
    emit("limit-change", newLimit);
  }
};

// Generate page numbers to display
const pageNumbers = computed(() => {
  const pages: (number | string)[] = [];
  const current = props.currentPage;
  const total = totalPages.value;

  if (total <= 7) {
    // Show all pages if 7 or less
    for (let i = 1; i <= total; i++) {
      pages.push(i);
    }
  } else {
    // Always show first page
    pages.push(1);

    if (current > 3) {
      pages.push("...");
    }

    // Show pages around current
    const start = Math.max(2, current - 1);
    const end = Math.min(total - 1, current + 1);

    for (let i = start; i <= end; i++) {
      pages.push(i);
    }

    if (current < total - 2) {
      pages.push("...");
    }

    // Always show last page
    pages.push(total);
  }

  return pages;
});

const handleRowClick = (ticket: Ticket) => {
  router.push(`/tickets/${ticket.id}`);
};

const getStatusSeverity = (status: string) => {
  const statusLower = status?.toLowerCase() || "";
  if (statusLower.includes("pending")) return "info";
  if (statusLower.includes("solved") || statusLower.includes("closed"))
    return "success";
  if (statusLower.includes("in progress")) return "warn";
  return "secondary";
};

const getTicketDisplayId = (ticket: Ticket) => {
  if (ticket.wrk_ticket_id) return ticket.wrk_ticket_id;
  const baseId = ticket.inv_gate_id ?? ticket.id;
  return `WRK-#${baseId}`;
};
</script>

<template>
  <div class="ticket-table">
    <div v-if="loading" class="loading-skeleton">
      <div v-for="i in 5" :key="i" class="skeleton-row">
        <div class="skeleton-cell"></div>
        <div class="skeleton-cell"></div>
        <div class="skeleton-cell"></div>
        <div class="skeleton-cell"></div>
      </div>
    </div>
    <div v-else class="table-container">
      <table class="ticket-data-table">
        <thead>
          <tr>
            <th class="col-ticket-id">{{ t('tickets.table.id') }}</th>
            <th class="col-title">{{ t('tickets.table.title') }}</th>
            <th class="col-status">{{ t('tickets.table.status') }}</th>
            <th class="col-created">{{ t('tickets.table.createdAt') }}</th>
            <th class="col-actions">{{ t('tickets.table.actions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="ticket in filteredTickets"
            :key="ticket.id"
            class="table-row"
            @click="handleRowClick(ticket)"
          >
            <td data-label="Ticket ID" class="col-ticket-id">
              {{ getTicketDisplayId(ticket) }}
            </td>
            <td data-label="Title" class="col-title">{{ ticket.title }}</td>
            <td data-label="Status" class="col-status">
              <Tag
                :id="`ticketTableStatusTag-${ticket.id}`"
                :severity="getStatusSeverity(ticket.status)"
              >
                {{ ticket.status }}
              </Tag>
            </td>
            <td data-label="Created At" class="col-created">
              {{ formatUnixTimestamp(ticket.created_at) }}
            </td>
            <td data-label="Actions" class="col-actions">
              <button
                :id="`ticketTableDetailBtn-${ticket.id}`"
                class="btn-detail"
                @click.stop="router.push(`/tickets/${ticket.id}`)"
              >
                {{ t('tickets.table.view') }}
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="!loading && totalItems > 0" class="pagination-container">
      <div class="pagination-info">
        Showing
        {{ pagination ? (pagination.page - 1) * pagination.limit + 1 : 1 }}
        to
        {{
          pagination
            ? Math.min(pagination.page * pagination.limit, pagination.total)
            : filteredTickets.length
        }}
        of
        {{ totalItems }}
        tickets
      </div>

      <div v-if="totalPages > 1" class="carbon-pagination">
        <!-- Page size selector -->
        <div class="pagination-left">
          <label class="pagination-label">{{ t('tickets.itemsPerPage') }}</label>
          <select
            :value="pageLimit"
            @change="handleLimitChange"
            class="bx--select-input"
          >
            <option value="5">5</option>
            <option value="10">10</option>
            <option value="20">20</option>
            <option value="50">50</option>
          </select>
        </div>

        <!-- Page navigation -->
        <div class="pagination-right">
          <!-- Left side: Previous then First (single arrow dekat double arrow) -->
          <button
            class="pagination-button"
            :disabled="currentPage === 1"
            @click="handlePageChange(currentPage - 1)"
            aria-label="Previous page"
            type="button"
          >
            <i class="pi pi-angle-left"></i>
          </button>

          <button
            class="pagination-button"
            :disabled="currentPage === 1"
            @click="handlePageChange(1)"
            aria-label="First page"
            type="button"
          >
            <i class="pi pi-angle-double-left"></i>
          </button>

          <!-- Page numbers -->
          <div class="pagination-pages">
            <button
              v-for="(page, index) in pageNumbers"
              :key="index"
              class="pagination-page-button"
              :class="{
                'pagination-page-active': page === currentPage,
                'pagination-ellipsis': page === '...',
              }"
              :disabled="page === '...'"
              @click="page !== '...' && handlePageChange(page as number)"
              type="button"
            >
              {{ page }}
            </button>
          </div>

          <!-- Right side: Last then Next (double arrow dekat single arrow) -->
          <button
            class="pagination-button"
            :disabled="currentPage === totalPages"
            @click="handlePageChange(totalPages)"
            aria-label="Last page"
            type="button"
          >
            <i class="pi pi-angle-double-right"></i>
          </button>

          <button
            class="pagination-button"
            :disabled="currentPage === totalPages"
            @click="handlePageChange(currentPage + 1)"
            aria-label="Next page"
            type="button"
          >
            <i class="pi pi-angle-right"></i>
          </button>
        </div>
      </div>
    </div>
    <div v-else-if="!loading && totalItems === 0" class="empty-state">
      <p>
        No tickets found{{ props.searchQuery ? " matching your search" : "" }}.
      </p>
    </div>
  </div>
</template>

<style scoped>
.ticket-table {
  margin-top: 1rem;
}

.table-container {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.ticket-data-table {
  width: 100%;
  border-collapse: collapse;
  background-color: #ffffff;
  border: 1px solid var(--border-color);
}

.ticket-data-table thead {
  background-color: #ffffff;
}

.ticket-data-table thead th {
  background-color: #ffffff;
  padding: 0.75rem 1rem;
  text-align: left;
  font-weight: 600;
  font-size: 0.875rem;
  line-height: 1.29;
  letter-spacing: 0.16px;
  color: var(--text-primary);
  border-bottom: 1px solid var(--border-color);
  white-space: nowrap;
}

.ticket-data-table thead th.col-actions {
  text-align: right;
}

.ticket-data-table tbody tr {
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
  transition: background-color 0.15s ease;
  background-color: #ffffff;
}

.ticket-data-table tbody tr:hover {
  background-color: var(--surface);
}

.ticket-data-table tbody tr:hover .btn-detail {
  background-color: var(--primary-hover);
  border-color: var(--primary-hover);
}

.ticket-data-table tbody tr:last-child {
  border-bottom: none;
}

.ticket-data-table tbody td {
  padding: 0.75rem 1rem;
  font-size: 0.875rem;
  line-height: 1.29;
  letter-spacing: 0.16px;
  color: var(--text-primary);
  vertical-align: middle;
}

.ticket-data-table tbody td.col-actions {
  text-align: right;
}

.btn-detail {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  font-weight: 400;
  line-height: 1.29;
  letter-spacing: 0.16px;
  background-color: var(--primary-color);
  color: #ffffff;
  border: 1px solid var(--primary-color);
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.11s cubic-bezier(0.2, 0, 0.38, 0.9);
  text-align: center;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  position: relative;
  outline: none;
  font-family: inherit;
  min-height: 2.5rem;
}

.btn-detail:hover {
  background-color: var(--primary-hover);
  border-color: var(--primary-hover);
}

.btn-detail:active {
  background-color: var(--primary-active);
  border-color: var(--primary-active);
}

.btn-detail:focus-visible {
  outline: 2px solid var(--primary-color);
  outline-offset: -2px;
}

.table-row {
  cursor: pointer;
}

.loading-skeleton {
  padding: 1rem;
}

.skeleton-row {
  display: flex;
  gap: 1rem;
  padding: 0.75rem 0;
  border-bottom: 1px solid var(--border-color);
}

.skeleton-cell {
  flex: 1;
  height: 1rem;
  background-color: #e0e0e0;
  border-radius: 2px;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.pagination-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 1.5rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border-color);
  flex-wrap: wrap;
  gap: 1rem;
}

.pagination-info {
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.empty-state {
  padding: 3rem 1rem;
  text-align: center;
  color: var(--text-secondary);
}

.empty-state p {
  margin: 0;
  font-size: 0.9375rem;
}

/* Carbon Design System style pagination */
.carbon-pagination {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.pagination-left {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.pagination-label {
  color: var(--text-secondary);
  font-size: 0.875rem;
  font-weight: 400;
  white-space: nowrap;
}

.bx--select-input {
  padding: 0.5rem 2rem 0.5rem 0.75rem;
  background-color: var(--surface);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  color: var(--text-primary);
  font-size: 0.875rem;
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 16 16'%3E%3Cpath fill='%23161616' d='M8 11L3 6h10z'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 0.5rem center;
  padding-right: 2rem;
  min-width: 5rem;
  height: 2.5rem;
  transition: background-color 0.11s cubic-bezier(0.2, 0, 0.38, 0.9);
}

.bx--select-input:hover {
  background-color: #e8e8e8;
}

.bx--select-input:focus {
  outline: 2px solid var(--primary-color);
  outline-offset: -2px;
  border-color: var(--primary-color);
}

.bx--select-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background-color: var(--surface);
}

.pagination-right {
  display: flex;
  align-items: center;
  gap: 0.125rem;
}

.pagination-button {
  min-width: 2.5rem;
  height: 2.5rem;
  padding: 0;
  background-color: #ffffff;
  border: 1px solid var(--border-color);
  color: var(--text-primary);
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background-color 0.11s cubic-bezier(0.2, 0, 0.38, 0.9);
  border-left: none;
}

.pagination-button:first-child {
  border-left: 1px solid var(--border-color);
}

.pagination-button:hover:not(:disabled) {
  background-color: var(--surface);
}

.pagination-button:focus:not(:disabled) {
  outline: 2px solid var(--primary-color);
  outline-offset: -2px;
  border-color: var(--primary-color);
}

.pagination-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background-color: var(--surface);
}

.pagination-button i {
  font-size: 1rem;
  display: block;
}

.pagination-pages {
  display: flex;
  gap: 0.125rem;
  margin: 0 0.5rem;
  align-items: center;
}

.pagination-page-button {
  min-width: 2.5rem;
  height: 2.5rem;
  padding: 0;
  background-color: #ffffff;
  border: 1px solid var(--border-color);
  color: var(--text-primary);
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: background-color 0.11s cubic-bezier(0.2, 0, 0.38, 0.9);
  font-weight: 400;
  font-size: 0.875rem;
  border-left: none;
}

.pagination-page-button:first-child {
  border-left: 1px solid var(--border-color);
}

.pagination-page-button:hover:not(:disabled):not(.pagination-page-active) {
  background-color: var(--surface);
}

.pagination-page-button:focus:not(:disabled) {
  outline: 2px solid var(--primary-color);
  outline-offset: -2px;
  border-color: var(--primary-color);
}

.pagination-page-button.pagination-page-active {
  background-color: var(--primary-color);
  color: #ffffff;
  border-color: var(--primary-color);
}

.pagination-page-button.pagination-page-active:hover {
  background-color: var(--primary-hover);
  border-color: var(--primary-hover);
}

.pagination-page-button.pagination-ellipsis {
  cursor: default;
  background-color: transparent;
  border: none;
  color: var(--text-secondary);
  min-width: auto;
  padding: 0 0.5rem;
}

.pagination-page-button.pagination-ellipsis:hover {
  background-color: transparent;
}

@media (max-width: 1024px) {
  .ticket-data-table {
    font-size: 0.875rem;
  }

  .ticket-data-table thead th,
  .ticket-data-table tbody td {
    padding: 0.625rem 0.75rem;
  }

  .btn-detail {
    padding: 0.375rem 0.75rem;
    font-size: 0.8125rem;
    min-height: 2rem;
  }
}

@media (max-width: 768px) {
  .table-container {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    margin: 0 -1rem;
    padding: 0 1rem;
  }

  .ticket-data-table {
    min-width: 600px;
  }

  .ticket-data-table thead th,
  .ticket-data-table tbody td {
    padding: 0.75rem 0.5rem;
    font-size: 0.8125rem;
  }

  .ticket-data-table thead th.col-created,
  .ticket-data-table tbody td.col-created {
    display: none;
  }

  .btn-detail {
    padding: 0.5rem 0.875rem;
    font-size: 0.8125rem;
    min-height: 2.25rem;
  }

  .pagination-container {
    flex-direction: column;
    gap: 1rem;
  }

  .pagination-info {
    text-align: center;
  }

  .carbon-pagination {
    flex-direction: column;
    gap: 1rem;
  }

  .pagination-left,
  .pagination-right {
    width: 100%;
    justify-content: center;
  }

  .pagination-pages {
    flex-wrap: wrap;
    justify-content: center;
  }
}

@media (max-width: 480px) {
  .table-container {
    margin: 0;
    padding: 0;
  }

  .ticket-data-table {
    border: 0;
    min-width: 100%;
  }

  .ticket-data-table thead {
    display: none;
  }

  .ticket-data-table tbody tr {
    display: block;
    margin-bottom: 1rem;
    border: 1px solid var(--border-color);
    border-radius: 8px;
    padding: 1rem;
    background-color: #ffffff;
    cursor: pointer;
  }

  .ticket-data-table tbody tr:hover {
    background-color: var(--surface);
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
  }

  .ticket-data-table tbody td {
    display: block;
    padding: 0.75rem 0;
    text-align: left;
    border-bottom: 1px solid var(--border-color);
  }

  .ticket-data-table tbody td:last-child {
    border-bottom: none;
    padding-top: 1rem;
  }

  .ticket-data-table tbody td::before {
    content: attr(data-label);
    font-weight: 600;
    color: var(--text-secondary);
    display: block;
    margin-bottom: 0.5rem;
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .ticket-data-table tbody td.col-actions {
    text-align: left;
  }

  .ticket-data-table tbody td.col-actions::before {
    margin-bottom: 0.75rem;
  }

  .btn-detail {
    width: 100%;
    padding: 0.75rem 1rem;
    font-size: 0.875rem;
    min-height: 2.75rem;
  }

  .pagination-pages {
    flex-wrap: wrap;
  }

  .pagination-page-button {
    min-width: 2rem;
    padding: 0.5rem;
    font-size: 0.8125rem;
  }

  .pagination-left {
    flex-direction: column;
    gap: 0.5rem;
  }

  .pagination-label {
    font-size: 0.8125rem;
  }

  .pagination-info {
    font-size: 0.8125rem;
  }
}
</style>
