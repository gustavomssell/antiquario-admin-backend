import { createFileRoute } from '@tanstack/react-router'
import { ProductForm } from '@/features/products/components/product-form'

export const Route = createFileRoute('/_authenticated/products/new')({
  component: () => (
    <div className='p-8 pt-6'>
      <ProductForm />
    </div>
  ),
})
