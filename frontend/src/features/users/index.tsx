import { UsersTable } from './components/users-table'
import { UserDialog } from './components/user-dialog'

export function Users() {
  return (
    <div className='flex flex-col space-y-8 p-8 pt-4'>
      <div className='flex items-center justify-between space-y-2'>
        <div>
          <h2 className='text-2xl font-bold tracking-tight'>Gestão de Acessos</h2>
          <p className='text-muted-foreground'>
            Controle os colaboradores, leiloeiros e curadores autorizados no sistema.
          </p>
        </div>
        <UserDialog />
      </div>
      <UsersTable />
    </div>
  )
}
