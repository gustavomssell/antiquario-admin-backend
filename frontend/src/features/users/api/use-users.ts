import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api'

export interface User {
  id: number
  userName: string
  email: string
  firstName: string
  lastName: string
  status: boolean
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
export function useUsers(page = 1, pageSize = 10, search = '') {
  return useQuery({
    queryKey: ['users', page, pageSize, search],
    queryFn: async () => {
      const params = new URLSearchParams()
      params.append('page', page.toString())
      params.append('pageSize', pageSize.toString())
      if (search) {
        params.append('userName_like', search)
      }
      
      const { data } = await api.get<SearchResponse<User>>(`/users/search?${params.toString()}`)
      return data
    },
  })
}

// Mutations
export function useCreateUser() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (payload: any) => {
      // The Go controller expects 'user' instead of 'userName' in Request
      const mappedPayload = {
        ...payload,
        user: payload.userName
      }
      const { data } = await api.post<User>('/user/', mappedPayload)
      return data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] })
    },
  })
}
