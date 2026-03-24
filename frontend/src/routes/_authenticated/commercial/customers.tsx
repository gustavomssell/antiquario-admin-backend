import { createFileRoute } from '@tanstack/react-router'
import { Customers } from '@/features/commercial/customers'

export const Route = createFileRoute('/_authenticated/commercial/customers')({
  component: Customers,
})
