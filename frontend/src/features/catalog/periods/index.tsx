import { PeriodsTable } from './components/periods-table'
import { PeriodDialog } from './components/period-dialog'

export function Periods() {
  return (
    <div className='flex flex-col space-y-8 p-8 pt-4'>
      <div className='flex items-center justify-between space-y-2'>
        <div>
          <h2 className='text-2xl font-bold tracking-tight'>Períodos Históricos</h2>
          <p className='text-muted-foreground'>
            Gerencie os períodos e séculos que definem as peças de catálogo.
          </p>
        </div>
        <PeriodDialog />
      </div>
      <PeriodsTable />
    </div>
  )
}
