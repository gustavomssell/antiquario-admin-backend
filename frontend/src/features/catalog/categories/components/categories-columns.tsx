import type { ColumnDef } from '@tanstack/react-table'
import type { Category } from '../../api/use-categories'
import { DataTableColumnHeader } from '@/components/data-table/column-header'
import { CategoryActions } from './categories-actions'

export const columns: ColumnDef<Category>[] = [

  {
    accessorKey: 'name',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='Nome da Categoria' />
    ),
  },
  {
    accessorKey: 'description',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='Descrição Restrita' />
    ),
  },
  {
    id: 'actions',
    cell: ({ row }) => {
      return <CategoryActions category={row.original} />
    },
  },
]
