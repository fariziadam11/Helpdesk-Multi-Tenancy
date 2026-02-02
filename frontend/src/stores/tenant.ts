import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { tenantApi, type TenantPublicInfo } from '@/api/tenant'
import { getCookie, setCookie, removeCookie } from '@/utils/cookies'
import { logger } from '@/utils/logger'

const TENANT_COOKIE_NAME = 'tenant'
const TENANT_ID_COOKIE_NAME = 'tenant_id'

export const useTenantStore = defineStore('tenant', () => {
  const tenant = ref<TenantPublicInfo | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Computed properties for easy access to tenant branding
  const tenantId = computed(() => tenant.value?.id || '')
  const tenantName = computed(() => tenant.value?.name || 'Helpdesk')
  const tenantSlug = computed(() => tenant.value?.slug || '')
  const logoUrl = computed(() => tenant.value?.logo_url || '/logo_white.svg')
  const primaryColor = computed(() => tenant.value?.primary_color || '#6929C4')

  /**
   * Extract tenant slug from subdomain
   * e.g., tenant1.localhost -> tenant1
   * e.g., tenant1.app.example.com -> tenant1
   */
  function extractSlugFromSubdomain(): string | null {
    const hostname = window.location.hostname
    
    // Skip localhost without subdomain
    if (hostname === 'localhost') {
      return null
    }

    const parts = hostname.split('.')

    // Handle localhost with subdomain (e.g., tenant1.localhost)
    if (parts.length === 2 && parts[1] === 'localhost') {
      return parts[0] ?? null
    }

    // Handle regular domain (e.g., tenant1.app.example.com)
    if (parts.length >= 3) {
      return parts[0] ?? null
    }

    return null
  }

  /**
   * Get tenant slug from various sources
   * Priority: subdomain > cookie > localStorage
   */
  function getStoredSlug(): string | null {
    // 1. Try subdomain first
    const subdomainSlug = extractSlugFromSubdomain()
    if (subdomainSlug) {
      return subdomainSlug
    }

    // 2. Try cookie
    const cookieData = getCookie(TENANT_COOKIE_NAME)
    if (cookieData) {
      try {
        const parsed = JSON.parse(cookieData)
        return parsed.slug || null
      } catch {
        return null
      }
    }

    // 3. Try localStorage (for development)
    if (typeof window !== 'undefined') {
      return localStorage.getItem('tenant_slug')
    }

    return null
  }

  /**
   * Get stored tenant ID (for API requests before tenant is loaded)
   */
  function getStoredTenantId(): string | null {
    // 1. Try cookie
    const cookieId = getCookie(TENANT_ID_COOKIE_NAME)
    if (cookieId) {
      return cookieId
    }

    // 2. Try localStorage (for development)
    if (typeof window !== 'undefined') {
      return localStorage.getItem('tenant_id')
    }

    return null
  }

  /**
   * Fetch tenant info from API and store it
   */
  async function fetchTenantInfo(slug: string): Promise<boolean> {
    isLoading.value = true
    error.value = null

    try {
      const data = await tenantApi.getPublicInfo(slug)
      tenant.value = data

      // Store tenant branding in cookie for persistence
      // Note: tenant_id is managed by auth store after login
      setCookie(TENANT_COOKIE_NAME, JSON.stringify(data), {
        expires: 7,
        path: '/',
        secure: import.meta.env.PROD,
        sameSite: 'lax',
      })

      logger.info('Tenant loaded', { slug, name: data.name })
      return true
    } catch (err) {
      error.value = 'Failed to load tenant'
      logger.error('Failed to fetch tenant info', err)
      return false
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Detect and load tenant on app initialization
   */
  async function detectTenant(): Promise<boolean> {
    const slug = getStoredSlug()
    
    if (!slug) {
      logger.warn('No tenant slug detected')
      error.value = 'Tenant not identified'
      return false
    }

    return fetchTenantInfo(slug)
  }

  /**
   * Manually set tenant (for development/testing)
   */
  function setTenant(tenantData: TenantPublicInfo) {
    tenant.value = tenantData
    
    // Store tenant branding in cookie
    // Note: tenant_id is managed by auth store after login
    setCookie(TENANT_COOKIE_NAME, JSON.stringify(tenantData), {
      expires: 7,
      path: '/',
      secure: import.meta.env.PROD,
      sameSite: 'lax',
    })
  }

  /**
   * Clear tenant data
   */
  function clearTenant() {
    tenant.value = null
    removeCookie(TENANT_COOKIE_NAME, { path: '/' })
    removeCookie(TENANT_ID_COOKIE_NAME, { path: '/' })
    
    if (typeof window !== 'undefined') {
      localStorage.removeItem('tenant_slug')
      localStorage.removeItem('tenant_id')
    }
  }

  /**
   * Initialize from stored data (for API requests before full tenant load)
   */
  function initializeFromStorage() {
    const cookieData = getCookie(TENANT_COOKIE_NAME)
    if (cookieData) {
      try {
        tenant.value = JSON.parse(cookieData)
      } catch {
        // Ignore parse errors
      }
    }
  }

  return {
    // State
    tenant,
    isLoading,
    error,
    
    // Computed
    tenantId,
    tenantName,
    tenantSlug,
    logoUrl,
    primaryColor,
    
    // Actions
    detectTenant,
    fetchTenantInfo,
    setTenant,
    clearTenant,
    getStoredTenantId,
    initializeFromStorage,
  }
})
