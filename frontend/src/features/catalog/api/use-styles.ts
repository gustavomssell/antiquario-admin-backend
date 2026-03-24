import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api'

export interface Style {
  id: string
  name: string
  description: string
  createdAt: string
  updatedAt: string
}

interface SearchResponse<T> {
  data: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

// Queries
export function useStyles(page = 1, pageSize = 10, search = '') {
  return useQuery({
    queryKey: ['styles', page, pageSize, search],
    queryFn: async () => {
      const params = new URLSearchParams()
      params.append('page', page.toString())
      params.append('pageSize', pageSize.toString())
      if (search) {
        params.append('name_like', search)
      }
      
      const { data } = await api.get<SearchResponse<Style>>(`/styles/search?${params.toString()}`)
      return data
    },
  })
}

// Mutations
export function useCreateStyle() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (payload: { name: string; description?: string }) => {
      const { data } = await api.post<Style>('/styles', payload)
      return data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['styles'] })
    },
  })
}
