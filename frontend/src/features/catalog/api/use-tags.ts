import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api'

export interface Tag {
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
export function useTags(page = 1, pageSize = 10, search = '') {
  return useQuery({
    queryKey: ['tags', page, pageSize, search],
    queryFn: async () => {
      const params = new URLSearchParams()
      params.append('page', page.toString())
      params.append('pageSize', pageSize.toString())
      if (search) {
        params.append('name_like', search)
      }
      
      const { data } = await api.get<SearchResponse<Tag>>(`/tags/search?${params.toString()}`)
      return data
    },
  })
}

// Mutations
export function useCreateTag() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (payload: { name: string; description?: string }) => {
      const { data } = await api.post<Tag>('/tags/', payload)
      return data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tags'] })
    },
  })
}
