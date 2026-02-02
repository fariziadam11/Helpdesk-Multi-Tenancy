import { http } from './http'

export interface ArticleAttachment {
  id: number
  url: string
  name: string
}

export interface Article {
  id: number
  title: string
  author_id: number
  responsible_id: number
  content: string
  creation_date: number
  last_update_date: number
  solved_requests: string | number
  views: number
  category_id: number
  is_private: boolean
  attachments: ArticleAttachment[]
  rating: number
}

export interface ArticlesResponse {
  data: Article[]
}

export interface ArticleCategory {
  id: number
  name: string
}

export const articlesApi = {
  getByCategory: async (categoryId: number): Promise<ArticlesResponse> => {
    const response = await http.get<ArticlesResponse>(`/articles?category_id=${categoryId}`)
    return response.data
  },
  getAll: async (): Promise<Article[]> => {
    try {
      // Use hardcoded categories (same as Articles/Index.vue and Landing.vue)
      const categories = [
        { id: 20, name: 'CRM' },
        { id: 16, name: 'ESS' },
        { id: 22, name: 'Intranet' },
        { id: 17, name: 'Kehadiran' },
        { id: 21, name: 'LMS' },
        { id: 24, name: 'Pengaturan Perusahaan' },
        { id: 19, name: 'Penggajian' },
        { id: 18, name: 'Personalia' },
        { id: 32, name: 'Applicant' },
        { id: 31, name: 'Employer' },
      ]
      
      if (!categories || categories.length === 0) {
        return []
      }
      
      // Fetch articles for all categories in parallel
      const articlePromises = categories.map(category =>
        http.get<ArticlesResponse>(`/articles?category_id=${category.id}`)
          .then(res => res.data.data || [])
          .catch(() => []) // Return empty array if category fails
      )
      
      const articleArrays = await Promise.all(articlePromises)
      // Flatten and deduplicate by id
      const allArticles = articleArrays.flat()
      const uniqueArticles = Array.from(
        new Map(allArticles.map(article => [article.id, article])).values()
      )
      
      return uniqueArticles
    } catch (error) {
      console.error('Failed to fetch all articles:', error)
      return []
    }
  },
}

