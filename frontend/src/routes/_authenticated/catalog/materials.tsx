import { createFileRoute } from '@tanstack/react-router'
import { Materials } from '@/features/catalog/materials'

export const Route = createFileRoute('/_authenticated/catalog/materials')({
  component: Materials,
})
