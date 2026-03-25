import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api'

export interface Customer {
  id: string
  name: string
  email: string
  phone: string
  document: string
  address: string
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
export function useCustomers(page = 1, pageSize = 10, search = '') {
  return useQuery({
    queryKey: ['customers', page, pageSize, search],
    queryFn: async () => {
      const params = new URLSearchParams()
      params.append('page', page.toString())
      params.append('pageSize', pageSize.toString())
      if (search) {
        params.append('name_like', search)
      }
      
      const { data } = await api.get<SearchResponse<Customer>>(`/customers/search?${params.toString()}`)
      return data
    },
  })
}

// Mutations
export function useCreateCustomer() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (payload: Partial<Customer>) => {
      const { data } = await api.post<Customer>('/customers/', payload)
      return data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['customers'] })
    },
  })
}
