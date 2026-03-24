import { createFileRoute } from '@tanstack/react-router'
import { Tags } from '@/features/catalog/tags'

export const Route = createFileRoute('/_authenticated/catalog/tags')({
  component: Tags,
})
