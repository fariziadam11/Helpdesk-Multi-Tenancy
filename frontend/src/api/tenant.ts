import { http } from './http'

export interface TenantPublicInfo {
  id: string
  name: string
  slug: string
  logo_url?: string
  primary_color?: string
}

export interface Tenant extends TenantPublicInfo {
  is_active: boolean
  invgate_company_id: number
  invgate_group_id: number
  invgate_location_id: number
  invgate_base_url: string
  invgate_username: string
  email_domain?: string
  email_sender?: string
  created_at?: string
  updated_at?: string
}

export interface CreateTenantRequest {
  name: string
  slug: string
  invgate_company_id: number
  invgate_group_id: number
  invgate_location_id: number
  invgate_base_url: string
  invgate_username: string
  invgate_password: string
  email_domain?: string
  email_sender?: string
  logo_url?: string
  primary_color?: string
  is_active?: boolean
}

export interface UpdateTenantRequest extends Partial<CreateTenantRequest> {}

interface ApiResponse<T> {
  success: boolean
  data: T
  error?: string
}

export const tenantApi = {
  /**
   * Get public tenant info by slug (for branding)
   * This endpoint doesn't require authentication
   */
  async getPublicInfo(slug: string): Promise<TenantPublicInfo> {
    const response = await http.get<ApiResponse<TenantPublicInfo>>(`/tenants/${slug}/info`)
    return response.data.data
  },

  /**
   * List all tenants (Admin only)
   */
  async list(): Promise<Tenant[]> {
    const response = await http.get<ApiResponse<Tenant[]>>('/admin/tenants')
    return response.data.data
  },

  /**
   * Get tenant by ID (Admin only)
   */
  async getById(id: string): Promise<Tenant> {
    const response = await http.get<ApiResponse<Tenant>>(`/admin/tenants/${id}`)
    return response.data.data
  },

  /**
   * Create new tenant (Admin only)
   */
  async create(data: CreateTenantRequest): Promise<Tenant> {
    const response = await http.post<ApiResponse<Tenant>>('/admin/tenants', data)
    return response.data.data
  },

  /**
   * Update tenant (Admin only)
   */
  async update(id: string, data: UpdateTenantRequest): Promise<Tenant> {
    const response = await http.put<ApiResponse<Tenant>>(`/admin/tenants/${id}`, data)
    return response.data.data
  },

  /**
   * Delete tenant (Admin only)
   */
  async delete(id: string): Promise<void> {
    await http.delete(`/admin/tenants/${id}`)
  },

  /**
   * Toggle tenant status (Admin only)
   */
  async updateStatus(id: string, isActive: boolean): Promise<Tenant> {
    const response = await http.patch<ApiResponse<Tenant>>(`/admin/tenants/${id}/status`, { is_active: isActive })
    return response.data.data
  },
}
