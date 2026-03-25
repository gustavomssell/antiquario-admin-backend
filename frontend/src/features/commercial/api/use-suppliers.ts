import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api'

export interface Supplier {
  id: string
  name: string
  email: string
  phone: string
  document: string
  address: string
  specialty: string
  notes: string
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
export function useSuppliers(page = 1, pageSize = 10, search = '') {
  return useQuery({
    queryKey: ['suppliers', page, pageSize, search],
    queryFn: async () => {
      const params = new URLSearchParams()
      params.append('page', page.toString())
      params.append('pageSize', pageSize.toString())
      if (search) {
        params.append('name_like', search)
      }
      
      const { data } = await api.get<SearchResponse<Supplier>>(`/suppliers/search?${params.toString()}`)
      return data
    },
  })
}

// Mutations
export function useCreateSupplier() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (payload: Partial<Supplier>) => {
      const { data } = await api.post<Supplier>('/suppliers/', payload)
      return data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['suppliers'] })
    },
  })
}
