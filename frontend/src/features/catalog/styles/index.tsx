import { StylesTable } from './components/styles-table'
import { StyleDialog } from './components/style-dialog'

export function Styles() {
  return (
    <div className='flex flex-col space-y-8 p-8 pt-4'>
      <div className='flex items-center justify-between space-y-2'>
        <div>
          <h2 className='text-2xl font-bold tracking-tight'>Estilos de Época</h2>
          <p className='text-muted-foreground'>
            Defina os fluxos artísticos e características estilísticas proeminentes que guiam seu inventário.
          </p>
        </div>
        <StyleDialog />
      </div>
      <StylesTable />
    </div>
  )
}
