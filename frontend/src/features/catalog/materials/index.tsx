import { MaterialsTable } from './components/materials-table'
import { MaterialDialog } from './components/material-dialog'

export function Materials() {
  return (
    <div className='flex flex-col space-y-8 p-8 pt-4'>
      <div className='flex items-center justify-between space-y-2'>
        <div>
          <h2 className='text-2xl font-bold tracking-tight'>Composições e Materiais</h2>
          <p className='text-muted-foreground'>
            Tabela de elementos brutos e ligas usados para construção das antiguidades.
          </p>
        </div>
        <MaterialDialog />
      </div>
      <MaterialsTable />
    </div>
  )
}
