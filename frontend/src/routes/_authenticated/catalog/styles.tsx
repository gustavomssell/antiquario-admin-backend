import { createFileRoute } from '@tanstack/react-router'
import { Styles } from '@/features/catalog/styles'

export const Route = createFileRoute('/_authenticated/catalog/styles')({
  component: Styles,
})
