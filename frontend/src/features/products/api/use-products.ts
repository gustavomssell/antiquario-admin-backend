import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api'

export interface Product {
  id: string
  sku: string
  title: string
  description: string
  dimensions: string
  weight: number
  condition: string
  status: string
  location: string
  acquisitionCost: number
  basePrice: number
  currentPrice: number
  categoryId: string
  periodId?: string
  styleId?: string
  manufacturingYear?: number
  isConsigned: boolean
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

export function useProducts(page = 1, pageSize = 10, search = '') {
  return useQuery({
    queryKey: ['products', page, pageSize, search],
    queryFn: async () => {
      const params = new URLSearchParams()
      params.append('page', page.toString())
      params.append('pageSize', pageSize.toString())
      if (search) {
        params.append('title_like', search)
      }
      
      const { data } = await api.get<SearchResponse<Product>>(`/products/search?${params.toString()}`)
      return data
    },
  })
}

export function useCreateProduct() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (payload: any) => {
      const { data } = await api.post<Product>('/products/', payload)
      return data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] })
    },
  })
}
