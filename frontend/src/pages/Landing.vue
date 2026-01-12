<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAllArticles } from '@/composables/useAllArticles'
import { useAuthStore } from '@/stores/auth'
import type { ArticleCategory } from '@/api/articles'
import ProgressSpinner from 'primevue/progressspinner'
import Button from 'primevue/button'

const { t } = useI18n()

const router = useRouter()
const authStore = useAuthStore()

// Categories
const categories: ArticleCategory[] = [
  { id: 20, name: 'CRM' },
  { id: 16, name: 'ESS' },
  { id: 22, name: 'Intranet' },
  { id: 23, name: 'Job Portal' },
  { id: 17, name: 'Kehadiran' },
  { id: 21, name: 'LMS' },
  { id: 24, name: 'Pengaturan Perusahaan' },
  { id: 19, name: 'Penggajian' },
  { id: 18, name: 'Personalia' },
]

const { articles, isLoading, error } = useAllArticles()

const searchQuery = ref('')
const selectedCategories = ref<number[]>([])
const sortBy = ref<'newest' | 'popular' | 'alphabetical'>('newest')
const currentPage = ref(1)
const itemsPerPage = 9

// Count articles per category
const categoryArticleCounts = computed(() => {
  const counts: Record<number, number> = {}
  categories.forEach(cat => {
    counts[cat.id] = (articles.value || []).filter(
      article => article.category_id === cat.id
    ).length
  })
  return counts
})

// Filter articles
const filteredArticles = computed(() => {
  let result = articles.value || []

  if (!authStore.isAuthenticated) {
    result = result.filter(article => !article.is_private)
  }

  if (selectedCategories.value.length > 0) {
    result = result.filter(article => selectedCategories.value.includes(article.category_id))
  }

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(article =>
      article.title.toLowerCase().includes(query) ||
      (article.content && article.content.toLowerCase().includes(query))
    )
  }

  return result
})

// Sort filtered articles
const sortedArticles = computed(() => {
  const result = [...filteredArticles.value]
  
  switch (sortBy.value) {
    case 'newest':
      return result.sort((a, b) => b.creation_date - a.creation_date)
    case 'popular':
      return result.sort((a, b) => b.views - a.views)
    case 'alphabetical':
      return result.sort((a, b) => a.title.localeCompare(b.title))
    default:
      return result
  }
})

// Paginate sorted articles
const paginatedArticles = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage
  const end = start + itemsPerPage
  return sortedArticles.value.slice(start, end)
})

const totalPages = computed(() => {
  return Math.ceil(sortedArticles.value.length / itemsPerPage)
})

const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    // Scroll to articles section
    const articlesSection = document.querySelector('.articles-section')
    if (articlesSection) {
      articlesSection.scrollIntoView({ behavior: 'smooth', block: 'start' })
    }
  }
}

const getCategoryName = (categoryId: number): string => {
  const category = categories.find(cat => cat.id === categoryId)
  return category?.name || 'Unknown'
}

const formatDate = (timestamp: number): string => {
  const date = new Date(timestamp * 1000)
  return date.toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

const truncateContent = (content: string, maxLength: number = 120): string => {
  if (!content) return ''
  const text = content.replace(/<[^>]*>/g, '')
  if (text.length <= maxLength) return text
  return text.substring(0, maxLength) + '...'
}

const goToArticle = (articleId: number, categoryId: number) => {
  router.push({
    path: `/articles/${articleId}`,
    query: { category_id: categoryId.toString() },
  })
}

const toggleCategory = (categoryId: number) => {
  const index = selectedCategories.value.indexOf(categoryId)
  if (index > -1) {
    selectedCategories.value.splice(index, 1)
  } else {
    selectedCategories.value.push(categoryId)
  }
}

const isCategorySelected = (categoryId: number): boolean => {
  return selectedCategories.value.includes(categoryId)
}
</script>

<template>
  <div class="landing-page">
    <!-- Hero Section -->
    <section class="hero">
      <div class="hero-container">
        <div class="hero-badge">
          <i class="pi pi-book"></i>
          <span>{{ t('landing.hero.badge') }}</span>
        </div>
        <h1 class="hero-title">{{ t('landing.hero.title') }}</h1>
        <p class="hero-description">
          {{ t('landing.hero.description') }}
        </p>

        <!-- Quick Stats -->
        <div class="stats-row">
          <div class="stat-item">
            <i class="pi pi-file"></i>
            <span>{{ t('landing.hero.stats.articles', { count: articles?.length || 0 }) }}</span>
          </div>
          <div class="stat-item">
            <i class="pi pi-th-large"></i>
            <span>{{ t('landing.hero.stats.categories', { count: categories.length }) }}</span>
          </div>
          <div class="stat-item">
            <i class="pi pi-users"></i>
            <span>{{ t('landing.hero.stats.support') }}</span>
          </div>
        </div>
      </div>
    </section>

    <!-- Category Filter -->
    <section class="categories-section">
      <div class="container">
        <h2 class="section-title">{{ t('landing.categories.title') }}</h2>
        <div class="categories-grid">
          <button
            v-for="category in categories"
            :key="category.id"
            class="category-chip"
            :class="{ active: isCategorySelected(category.id) }"
            @click="toggleCategory(category.id)"
          >
            <i class="pi pi-folder"></i>
            <span>{{ category.name }}</span>
            <span class="category-count">{{ categoryArticleCounts[category.id] || 0 }}</span>
            <i v-if="isCategorySelected(category.id)" class="pi pi-check check-icon"></i>
          </button>
        </div>
      </div>
    </section>

    <!-- Articles Section -->
    <section class="articles-section">
      <div class="container">
        <!-- Loading -->
        <div v-if="isLoading" class="state-container">
          <ProgressSpinner style="width: 50px; height: 50px" />
          <p class="state-text">Memuat artikel...</p>
        </div>

        <!-- Error -->
        <div v-else-if="error" class="state-container">
          <i class="pi pi-exclamation-circle state-icon error"></i>
          <p class="state-text">{{ error instanceof Error ? error.message : 'Gagal memuat artikel' }}</p>
          <Button label="Coba Lagi" @click="router.go(0)" />
        </div>

        <!-- Empty -->
        <div v-else-if="sortedArticles.length === 0" class="state-container">
          <i class="pi pi-inbox state-icon"></i>
          <p class="state-text">
            {{ searchQuery || selectedCategories.length > 0 ? t('landing.articles.noResults') : t('landing.articles.empty') }}
          </p>
          <Button
            v-if="searchQuery || selectedCategories.length > 0"
            :label="t('landing.articles.resetFilter')"
            severity="secondary"
            @click="searchQuery = ''; selectedCategories = []"
          />
        </div>

        <!-- Articles Grid -->
        <div v-else>
          <div class="articles-controls">
            <div class="articles-header">
              <h2 class="section-title">
                {{ selectedCategories.length > 0 ? t('landing.articles.selected') : t('landing.articles.all') }}
              </h2>
              <span class="articles-count">{{ t('landing.articles.count', { count: sortedArticles.length }) }}</span>
            </div>
            
            <div class="controls-row">
              <!-- Search Box -->
              <div class="search-box-articles">
                <i class="pi pi-search search-icon-articles"></i>
                <input
                  v-model="searchQuery"
                  type="text"
                  class="search-input-articles"
                  :placeholder="t('landing.articles.searchPlaceholder')"
                />
              </div>
              
              <!-- Sort Controls -->
              <div class="sort-controls">
                <label>{{ t('landing.articles.sort.label') }}</label>
                <select v-model="sortBy" class="sort-select">
                  <option value="newest">{{ t('landing.articles.sort.newest') }}</option>
                  <option value="popular">{{ t('landing.articles.sort.popular') }}</option>
                  <option value="alphabetical">{{ t('landing.articles.sort.alphabetical') }}</option>
                </select>
              </div>
            </div>
          </div>

          <div class="articles-grid">
            <article
              v-for="article in paginatedArticles"
              :key="article.id"
              class="article-card"
              @click="goToArticle(article.id, article.category_id)"
            >
              <div class="article-header">
                <span class="article-category">{{ getCategoryName(article.category_id) }}</span>
                <div class="article-stats">
                  <span v-if="article.rating > 0" class="stat">
                    <i class="pi pi-star-fill"></i>
                    {{ article.rating.toFixed(1) }}
                  </span>
                  <span class="stat">
                    <i class="pi pi-eye"></i>
                    {{ article.views }}
                  </span>
                </div>
              </div>

              <h3 class="article-title">{{ article.title }}</h3>
              <p v-if="article.content" class="article-excerpt">
                {{ truncateContent(article.content) }}
              </p>

              <div class="article-footer">
                <span class="article-date">
                  <i class="pi pi-calendar"></i>
                  {{ formatDate(article.creation_date) }}
                </span>
                <span class="read-more">
                  {{ t('landing.articles.readMore') }} <i class="pi pi-arrow-right"></i>
                </span>
              </div>
            </article>
          </div>
          
          <!-- Pagination -->
          <div v-if="totalPages > 1" class="pagination">
            <button
              class="pagination-btn"
              :disabled="currentPage === 1"
              @click="goToPage(currentPage - 1)"
            >
              <i class="pi pi-chevron-left"></i>
            </button>
            
            <button
              v-for="page in totalPages"
              :key="page"
              class="pagination-btn"
              :class="{ active: currentPage === page }"
              @click="goToPage(page)"
            >
              {{ page }}
            </button>
            
            <button
              class="pagination-btn"
              :disabled="currentPage === totalPages"
              @click="goToPage(currentPage + 1)"
            >
              <i class="pi pi-chevron-right"></i>
            </button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.landing-page {
  min-height: 100vh;
  background: #fafafa;
}

/* Hero Section */
.hero {
  background: linear-gradient(135deg, #6929C4 0%, #8A3FFC 50%, #A56EFF 100%);
  padding: 4rem 2rem 4rem;
  position: relative;
  overflow: hidden;
  min-height: 100vh;
  display: flex;
  align-items: center;
}

.hero::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.05'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
  opacity: 0.3;
}

.hero-container {
  max-width: 900px;
  margin: 0 auto;
  text-align: center;
  position: relative;
  z-index: 1;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  padding: 0.5rem 1.25rem;
  border-radius: 50px;
  color: #ffffff;
  font-size: 0.875rem;
  font-weight: 500;
  margin-bottom: 1.5rem;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.hero-title {
  font-size: 3.5rem;
  font-weight: 700;
  color: #ffffff;
  margin: 0 0 1.5rem 0;
  line-height: 1.2;
  letter-spacing: -0.02em;
}

.hero-description {
  font-size: 1.25rem;
  color: rgba(255, 255, 255, 0.9);
  margin: 0 0 3rem 0;
  line-height: 1.6;
  max-width: 700px;
  margin-left: auto;
  margin-right: auto;
}

.search-box {
  position: relative;
  max-width: 600px;
  margin: 0 auto 3rem;
}

.search-icon {
  position: absolute;
  left: 1.5rem;
  top: 50%;
  transform: translateY(-50%);
  color: #6929C4;
  font-size: 1.25rem;
  z-index: 1;
}

.search-input {
  width: 100%;
  padding: 1.25rem 1.5rem 1.25rem 4rem;
  font-size: 1rem;
  border: none;
  border-radius: 50px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
  transition: all 0.3s ease;
}

.search-input:focus {
  box-shadow: 0 15px 50px rgba(0, 0, 0, 0.3);
  transform: translateY(-2px);
}

.stats-row {
  display: flex;
  justify-content: center;
  gap: 3rem;
  flex-wrap: wrap;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  color: #ffffff;
  font-size: 0.9375rem;
  font-weight: 500;
}

.stat-item i {
  font-size: 1.25rem;
  opacity: 0.9;
}

/* Categories Section */
.categories-section {
  padding: 4rem 2rem;
  background: #ffffff;
}

.container {
  max-width: 1400px;
  margin: 0 auto;
}

.section-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 2rem 0;
  text-align: center;
}

.categories-grid {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 0.75rem;
}

.category-chip {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: #f4f4f4;
  border: 2px solid transparent;
  border-radius: 50px;
  color: var(--text-primary);
  font-size: 0.9375rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.category-chip:hover {
  background: #e8daff;
  border-color: #6929C4;
  transform: translateY(-2px);
}

.category-chip.active {
  background: #6929C4;
  color: #ffffff;
  border-color: #6929C4;
}

.check-icon {
  font-size: 0.875rem;
}

.category-count {
  background: rgba(105, 41, 196, 0.15);
  color: #6929C4;
  padding: 0.25rem 0.625rem;
  border-radius: 50px;
  font-size: 0.75rem;
  font-weight: 600;
  margin-left: auto;
}

.category-chip.active .category-count {
  background: rgba(255, 255, 255, 0.25);
  color: #ffffff;
}

/* Articles Section */
.articles-section {
  padding: 4rem 2rem;
  background: #fafafa;
}

.state-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  gap: 1.5rem;
}

.state-icon {
  font-size: 4rem;
  color: #8d8d8d;
}

.state-icon.error {
  color: var(--error-color);
}

.state-text {
  font-size: 1.125rem;
  color: var(--text-secondary);
  text-align: center;
  margin: 0;
}

.articles-controls {
  margin-bottom: 2rem;
}

.articles-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.controls-row {
  display: flex;
  gap: 1rem;
  align-items: center;
  flex-wrap: wrap;
  width: 100%;
}

.search-box-articles {
  position: relative;
  flex: 1;
  min-width: 250px;
}

.search-icon-articles {
  position: absolute;
  left: 1rem;
  top: 50%;
  transform: translateY(-50%);
  color: #8d8d8d;
  font-size: 1rem;
}

.search-input-articles {
  width: 100%;
  padding: 0.75rem 1rem 0.75rem 3rem;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-size: 0.9375rem;
  transition: all 0.2s ease;
}

.search-input-articles:focus {
  outline: none;
  border-color: #6929C4;
  box-shadow: 0 0 0 3px rgba(105, 41, 196, 0.1);
}

.sort-controls {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.sort-controls label {
  font-size: 0.9375rem;
  color: var(--text-secondary);
  font-weight: 500;
}

.sort-select {
  padding: 0.625rem 1rem;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  font-size: 0.9375rem;
  background: #ffffff;
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 150px;
}

.sort-select:hover {
  border-color: #6929C4;
}

.sort-select:focus {
  outline: none;
  border-color: #6929C4;
  box-shadow: 0 0 0 3px rgba(105, 41, 196, 0.1);
}

.articles-count {
  color: var(--text-secondary);
  font-size: 0.9375rem;
}

/* Pagination */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 0.5rem;
  margin-top: 3rem;
  flex-wrap: wrap;
}

.pagination-btn {
  min-width: 40px;
  height: 40px;
  padding: 0.5rem 1rem;
  border: 1px solid #e0e0e0;
  background: #ffffff;
  color: var(--text-primary);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.9375rem;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: center;
}

.pagination-btn:hover:not(:disabled) {
  border-color: #6929C4;
  color: #6929C4;
}

.pagination-btn.active {
  background: #6929C4;
  color: #ffffff;
  border-color: #6929C4;
}

.pagination-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.articles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1.5rem;
}

.article-card {
  background: #ffffff;
  border: 1px solid #e0e0e0;
  border-radius: 12px;
  padding: 1.75rem;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.article-card:hover {
  border-color: #6929C4;
  box-shadow: 0 8px 24px rgba(105, 41, 196, 0.15);
  transform: translateY(-4px);
}

.article-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.article-category {
  background: linear-gradient(135deg, #e8daff 0%, #f0e6ff 100%);
  color: #6929C4;
  padding: 0.375rem 1rem;
  border-radius: 50px;
  font-size: 0.8125rem;
  font-weight: 600;
  white-space: nowrap;
}

.article-stats {
  display: flex;
  gap: 1rem;
  font-size: 0.8125rem;
  color: var(--text-secondary);
}

.stat {
  display: flex;
  align-items: center;
  gap: 0.375rem;
}

.stat i {
  font-size: 0.875rem;
}

.article-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-excerpt {
  font-size: 0.9375rem;
  color: var(--text-secondary);
  line-height: 1.6;
  margin: 0;
  flex: 1;
}

.article-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 1rem;
  border-top: 1px solid #e0e0e0;
  font-size: 0.875rem;
}

.article-date {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--text-secondary);
}

.read-more {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #6929C4;
  font-weight: 600;
  transition: gap 0.2s ease;
}

.article-card:hover .read-more {
  gap: 0.75rem;
}

/* CTA Section */
.cta-section {
  background: linear-gradient(135deg, #4F2196 0%, #6929C4 100%);
  padding: 5rem 2rem;
}

.cta-container {
  max-width: 700px;
  margin: 0 auto;
}

.cta-content {
  text-align: center;
  color: #ffffff;
}

.cta-icon {
  font-size: 4rem;
  margin-bottom: 1.5rem;
  opacity: 0.9;
}

.cta-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin: 0 0 1rem 0;
}

.cta-description {
  font-size: 1.125rem;
  margin: 0 0 2.5rem 0;
  opacity: 0.9;
  line-height: 1.6;
}

/* Responsive */
@media (max-width: 768px) {
  .hero {
    padding: 3rem 1.5rem 3rem;
    min-height: 100vh;
  }
  
  .controls-row {
    flex-direction: column;
  }
  
  .search-box-articles {
    width: 100%;
  }

  .hero-title {
    font-size: 2.25rem;
  }

  .hero-description {
    font-size: 1rem;
  }

  .stats-row {
    gap: 1.5rem;
  }

  .section-title {
    font-size: 1.5rem;
  }

  .articles-grid {
    grid-template-columns: 1fr;
  }

  .cta-title {
    font-size: 1.75rem;
  }

  .cta-description {
    font-size: 1rem;
  }
}
</style>