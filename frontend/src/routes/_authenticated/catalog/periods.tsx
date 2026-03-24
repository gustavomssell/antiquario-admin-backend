import { createFileRoute } from '@tanstack/react-router'
import { Periods } from '@/features/catalog/periods'

export const Route = createFileRoute('/_authenticated/catalog/periods')({
  component: Periods,
})
