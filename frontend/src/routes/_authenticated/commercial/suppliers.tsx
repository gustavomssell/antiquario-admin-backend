import { createFileRoute } from '@tanstack/react-router'
import { Suppliers } from '@/features/commercial/suppliers'

export const Route = createFileRoute('/_authenticated/commercial/suppliers')({
  component: Suppliers,
})
