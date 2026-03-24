import { Link } from '@tanstack/react-router'
import { Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { ProductsTable } from './components/products-table'

export function Products() {
  return (
    <div className='flex flex-col space-y-8 p-8 pt-4'>
      <div className='flex items-center justify-between space-y-2'>
        <div>
          <h2 className='text-2xl font-bold tracking-tight'>Inventário de Obras</h2>
          <p className='text-muted-foreground'>
            Gestão avançada do acervo. Visualize status, valores e crie fichas técnicas detalhadas.
          </p>
        </div>
        <Button asChild>
          <Link to="/products/new">
            <Plus className='mr-2 h-4 w-4' /> Adicionar Peça (Ficha Técnica)
          </Link>
        </Button>
      </div>
      <ProductsTable />
    </div>
  )
}
