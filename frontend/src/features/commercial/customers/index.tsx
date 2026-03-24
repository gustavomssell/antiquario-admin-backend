import { CustomersTable } from './components/customers-table'
import { CustomerDialog } from './components/customer-dialog'

export function Customers() {
  return (
    <div className='flex flex-col space-y-8 p-8 pt-4'>
      <div className='flex items-center justify-between space-y-2'>
        <div>
          <h2 className='text-2xl font-bold tracking-tight'>Gestão de Clientes (CRM)</h2>
          <p className='text-muted-foreground'>
            Carteira de colecionadores, galerias e compradores ativos do Antiquário.
          </p>
        </div>
        <CustomerDialog />
      </div>
      <CustomersTable />
    </div>
  )
}
