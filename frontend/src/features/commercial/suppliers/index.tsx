import { SuppliersTable } from './components/suppliers-table'
import { SupplierDialog } from './components/supplier-dialog'

export function Suppliers() {
  return (
    <div className='flex flex-col space-y-8 p-8 pt-4'>
      <div className='flex items-center justify-between space-y-2'>
        <div>
          <h2 className='text-2xl font-bold tracking-tight'>Gestão de Fornecedores</h2>
          <p className='text-muted-foreground'>
            Gestão de caçadores de antiguidades, galerias parceiras e conselheiros.
          </p>
        </div>
        <SupplierDialog />
      </div>
      <SuppliersTable />
    </div>
  )
}
