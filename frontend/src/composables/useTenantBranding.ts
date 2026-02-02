import { computed } from 'vue'
import { useTenantStore } from '@/stores/tenant'

/**
 * Composable for accessing tenant branding data
 * Use this in components that need to display tenant-specific branding
 */
export function useTenantBranding() {
  const tenantStore = useTenantStore()

  const logoUrl = computed(() => {
    const url = tenantStore.logoUrl
    if (!url) return '/logo_white.svg'
    
    // Handle relative uploads
    if (url.startsWith('/uploads/')) {
        const apiBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
        const origin = new URL(apiBase).origin
        return `${origin}${url}`
    }
    
    return url
  })
  const primaryColor = computed(() => tenantStore.primaryColor)
  const tenantName = computed(() => tenantStore.tenantName)
  const tenantSlug = computed(() => tenantStore.tenantSlug)
  const isLoading = computed(() => tenantStore.isLoading)
  const error = computed(() => tenantStore.error)

  /**
   * Generate CSS variables for tenant branding
   * Can be applied to :root or specific elements
   */
  const brandingStyles = computed(() => ({
    '--tenant-primary-color': tenantStore.primaryColor,
    '--tenant-primary-hover': adjustBrightness(tenantStore.primaryColor, -15),
    '--tenant-primary-light': adjustBrightness(tenantStore.primaryColor, 40),
  }))

  /**
   * Check if tenant branding is loaded
   */
  const isBrandingLoaded = computed(() => !!tenantStore.tenant)

  return {
    logoUrl,
    primaryColor,
    tenantName,
    tenantSlug,
    isLoading,
    error,
    brandingStyles,
    isBrandingLoaded,
  }
}

/**
 * Adjust brightness of a hex color
 * @param hex - Hex color string (e.g., "#6929C4")
 * @param percent - Positive to lighten, negative to darken
 */
function adjustBrightness(hex: string, percent: number): string {
  // Remove # if present
  hex = hex.replace(/^#/, '')
  
  // Parse hex to RGB
  let r = parseInt(hex.substring(0, 2), 16)
  let g = parseInt(hex.substring(2, 4), 16)
  let b = parseInt(hex.substring(4, 6), 16)
  
  // Adjust brightness
  r = Math.min(255, Math.max(0, r + (r * percent) / 100))
  g = Math.min(255, Math.max(0, g + (g * percent) / 100))
  b = Math.min(255, Math.max(0, b + (b * percent) / 100))
  
  // Convert back to hex
  const toHex = (n: number) => Math.round(n).toString(16).padStart(2, '0')
  return `#${toHex(r)}${toHex(g)}${toHex(b)}`
}
