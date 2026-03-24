import { TagsTable } from './components/tags-table'
import { TagDialog } from './components/tag-dialog'

export function Tags() {
  return (
    <div className='flex flex-col space-y-8 p-8 pt-4'>
      <div className='flex items-center justify-between space-y-2'>
        <div>
          <h2 className='text-2xl font-bold tracking-tight'>Rótulos Classificatórios (Tags)</h2>
          <p className='text-muted-foreground'>
            Gestão de rótulos avulsos para metadados secundários de obras e produtos.
          </p>
        </div>
        <TagDialog />
      </div>
      <TagsTable />
    </div>
  )
}
