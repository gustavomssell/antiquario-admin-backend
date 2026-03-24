import { CategoriesTable } from './components/categories-table'
import { CategoryDialog } from './components/category-dialog'

export function Categories() {
  return (
    <div className='flex flex-col space-y-8 p-8 pt-4'>
      <div className='flex items-center justify-between space-y-2'>
        <div>
          <h2 className='text-2xl font-bold tracking-tight'>Categorias Bases</h2>
          <p className='text-muted-foreground'>
            Gerencie as categorias estruturais exclusivas do Antiquário.
          </p>
        </div>
        <CategoryDialog />
      </div>
      <CategoriesTable />
    </div>
  )
}
