<script setup lang="ts">
import { useRouter } from "vue-router";
import { useI18n } from "vue-i18n";
import { useDashboard } from "@/composables/useDashboard";
import { useAuthStore } from "@/stores/auth";
import ProgressSpinner from "primevue/progressspinner";
import Button from "primevue/button";

const { t } = useI18n();
const router = useRouter();
const authStore = useAuthStore();
const { stats, isLoading } = useDashboard();

const formatDate = (timestamp: number): string => {
  const date = new Date(timestamp * 1000);
  return date.toLocaleDateString("id-ID", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const getStatusColor = (status: string): string => {
  const statusColors: Record<string, string> = {
    New: "#6929C4",
    Open: "#6929C4",
    Pending: "#8A3FFC",
    Waiting: "#8A3FFC",
    Resolved: "#24a148",
    Closed: "#24a148",
    Rejected: "#da1e28",
    Canceled: "#6f6f6f",
  };
  return statusColors[status] || "#6f6f6f";
};
</script>

<template>
  <div class="dashboard-page">
    <div class="page-header">
      <div>
        <h1>{{ t('dashboard.title') }}</h1>
        <p class="welcome-text">{{ t('dashboard.welcome') }}, {{ authStore.fullName }}!</p>
      </div>
    </div>

    <div v-if="isLoading" class="loading-container">
      <ProgressSpinner />
    </div>

    <div v-else class="dashboard-content">
      <!-- Statistics Cards -->
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon total">
            <i class="pi pi-ticket"></i>
          </div>
          <div class="stat-content">
            <h3 class="stat-value">{{ stats.totalTickets }}</h3>
            <p class="stat-label">{{ t('dashboard.stats.totalTickets') }}</p>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon open">
            <i class="pi pi-clock"></i>
          </div>
          <div class="stat-content">
            <h3 class="stat-value">{{ stats.openTickets }}</h3>
            <p class="stat-label">{{ t('dashboard.stats.openTickets') }}</p>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon resolved">
            <i class="pi pi-check-circle"></i>
          </div>
          <div class="stat-content">
            <h3 class="stat-value">{{ stats.resolvedTickets }}</h3>
            <p class="stat-label">{{ t('dashboard.stats.resolvedTickets') }}</p>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon pending">
            <i class="pi pi-hourglass"></i>
          </div>
          <div class="stat-content">
            <h3 class="stat-value">{{ stats.pendingTickets }}</h3>
            <p class="stat-label">{{ t('dashboard.stats.pendingTickets') }}</p>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon rejected">
            <i class="pi pi-times-circle"></i>
          </div>
          <div class="stat-content">
            <h3 class="stat-value">{{ stats.rejectedTickets }}</h3>
            <p class="stat-label">{{ t('dashboard.stats.rejectedTickets') }}</p>
          </div>
        </div>
      </div>

      <!-- Recent Tickets & Quick Actions -->
      <div class="bottom-grid">
        <div class="recent-tickets-card">
          <div class="card-header">
            <h2 class="card-title">{{ t('dashboard.recentTickets.title') }}</h2>
            <Button
              size="small"
              :label="t('dashboard.recentTickets.viewAll')"
              @click="router.push('/tickets')"
            />
          </div>
          <div v-if="stats.recentTickets.length === 0" class="empty-state">
            <p>{{ t('dashboard.recentTickets.noTickets') }}</p>
          </div>
          <div v-else class="tickets-list">
            <div
              v-for="ticket in stats.recentTickets"
              :key="ticket.id"
              class="ticket-item"
              @click="router.push(`/tickets/${ticket.id}`)"
            >
              <div class="ticket-info">
                <h4 class="ticket-title">{{ ticket.title }}</h4>
                <p class="ticket-meta">
                  <span class="ticket-id">{{
                    ticket.wrk_ticket_id || `#${ticket.id}`
                  }}</span>
                  <span class="ticket-date">{{
                    formatDate(ticket.created_at)
                  }}</span>
                </p>
              </div>
              <div class="ticket-status">
                <span
                  class="status-badge"
                  :style="{
                    backgroundColor: getStatusColor(ticket.status) + '20',
                    color: getStatusColor(ticket.status),
                  }"
                >
                  {{ ticket.status }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div class="quick-actions-card">
          <h2 class="card-title">{{ t('dashboard.quickActions.title') }}</h2>
          <div class="actions-list">
            <button
              class="action-button"
              @click="router.push('/tickets/create')"
            >
              <i class="pi pi-plus-circle"></i>
              <span>{{ t('dashboard.quickActions.createTicket') }}</span>
            </button>
            <button class="action-button" @click="router.push('/tickets')">
              <i class="pi pi-list"></i>
              <span>{{ t('dashboard.quickActions.viewAllTickets') }}</span>
            </button>
            <button class="action-button" @click="router.push('/articles')">
              <i class="pi pi-book"></i>
              <span>{{ t('dashboard.quickActions.browseArticles') }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard-page {
  padding: 2rem;
  padding-top: 0.5rem;
  max-width: 1600px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 2rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.5rem 0;
}

.welcome-text {
  color: var(--text-secondary);
  font-size: 1rem;
  margin: 0;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 4rem;
}

.dashboard-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

/* Statistics Cards */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
}

.stat-card {
  background-color: #ffffff;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1.5rem;
  transition: box-shadow 0.2s ease;
}

.stat-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  flex-shrink: 0;
}

.stat-icon.total {
  background-color: #6929c420;
  color: #6929c4;
}

.stat-icon.open {
  background-color: #6929c420;
  color: #6929c4;
}

.stat-icon.resolved {
  background-color: #24a14820;
  color: #24a148;
}

.stat-icon.pending {
  background-color: #8a3ffc20;
  color: #8a3ffc;
}

.stat-icon.rejected {
  background-color: #da1e2820;
  color: #da1e28;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 2rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.25rem 0;
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
}

/* Bottom Grid */
.bottom-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 1.5rem;
}

@media (max-width: 900px) {
  .bottom-grid {
    grid-template-columns: 1fr;
  }
}

.recent-tickets-card,
.quick-actions-card {
  background-color: #ffffff;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 1.5rem;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.card-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 1.5rem 0;
}

/* Recent Tickets */
.tickets-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.ticket-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background-color: var(--surface);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.ticket-item:hover {
  border-color: var(--primary-color);
  background-color: #f4f4f4;
  transform: translateX(4px);
}

.ticket-info {
  flex: 1;
}

.ticket-title {
  font-size: 1rem;
  font-weight: 500;
  color: var(--text-primary);
  margin: 0 0 0.5rem 0;
}

.ticket-meta {
  display: flex;
  gap: 1rem;
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0;
}

.ticket-id {
  font-weight: 500;
}

.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
}

.empty-state {
  text-align: center;
  padding: 2rem;
  color: var(--text-secondary);
}

/* Quick Actions */
.quick-actions-card .card-title {
  margin-bottom: 1.5rem;
}

.actions-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.action-button {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  background-color: var(--surface);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
}

.action-button:hover {
  background-color: #f4f4f4;
  border-color: var(--primary-color);
  transform: translateX(4px);
}

.action-button i {
  font-size: 1.25rem;
  color: var(--primary-color);
}

@media (max-width: 1024px) {
  .charts-grid {
    grid-template-columns: 1fr;
  }

  .chart-card {
    min-height: 400px;
  }

  .chart-container {
    min-height: 300px;
  }
}

@media (max-width: 768px) {
  .dashboard-page {
    padding: 1rem;
  }

  .page-header h1 {
    font-size: 1.5rem;
  }

  .welcome-text {
    font-size: 0.875rem;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
  }

  .stat-card {
    padding: 1rem;
  }

  .stat-icon {
    width: 50px;
    height: 50px;
    font-size: 1.25rem;
  }

  .stat-value {
    font-size: 1.5rem;
  }

  .charts-grid,
  .bottom-grid {
    grid-template-columns: 1fr;
  }

  .chart-card {
    padding: 1rem;
    min-height: 400px;
  }

  .chart-container {
    min-height: 300px;
  }

  .recent-tickets-card,
  .quick-actions-card {
    padding: 1rem;
  }

  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .card-header ::deep .p-button {
    width: 100%;
  }

  .actions-list {
    gap: 0.5rem;
  }

  .action-button {
    padding: 0.875rem;
    font-size: 0.875rem;
    width: 100%;
    justify-content: flex-start;
  }

  .action-button i {
    font-size: 1.125rem;
  }
}

@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .stat-card {
    flex-direction: row;
    gap: 1rem;
  }

  .bottom-grid {
    grid-template-columns: 1fr;
  }

  .recent-tickets-card,
  .quick-actions-card {
    padding: 0.875rem;
  }

  .card-title {
    font-size: 1.125rem;
  }

  .action-button {
    padding: 0.75rem;
    font-size: 0.8125rem;
  }

  .action-button i {
    font-size: 1rem;
  }

  .ticket-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.75rem;
  }

  .ticket-status {
    width: 100%;
  }

  .status-badge {
    display: inline-block;
  }
}
</style>
