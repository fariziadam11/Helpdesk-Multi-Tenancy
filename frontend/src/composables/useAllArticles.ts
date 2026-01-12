import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { articlesApi } from '@/api/articles'
import type { Article } from '@/api/articles'

export const useAllArticles = () => {
  const query = useQuery({
    queryKey: ['articles', 'all'],
    queryFn: () => articlesApi.getAll(),
    staleTime: 5 * 60 * 1000, // 5 minutes
  })

  const articles = computed<Article[]>(() => {
    return query.data.value || []
  })
  
  return {
    ...query,
    articles,
  }
}

