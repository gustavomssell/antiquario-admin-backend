import { ColumnDef } from '@tanstack/react-table'
import { Supplier } from '../../api/use-suppliers'
import { DataTableColumnHeader } from '@/components/data-table/column-header'
import { MoreHorizontal, Edit, Trash, Truck } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

export const columns: ColumnDef<Supplier>[] = [
  {
    accessorKey: 'name',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='Nome do Fornecedor' />
    ),
    cell: ({ row }) => <span className="font-semibold">{row.original.name}</span>
  },
  {
    accessorKey: 'specialty',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='Especialidade' />
    ),
    cell: ({ row }) => {
      return row.original.specialty ? (
        <Badge variant="outline">{row.original.specialty}</Badge>
      ) : (
        <span className="text-muted-foreground">-</span>
      )
    }
  },
  {
    accessorKey: 'phone',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='Contato (Telefone)' />
    ),
  },
  {
    accessorKey: 'document',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='CPF/CNPJ' />
    ),
    cell: ({ row }) => <span className="font-mono text-xs">{row.original.document}</span>
  },
  {
    id: 'actions',
    cell: () => {
      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant='ghost' className='h-8 w-8 p-0'>
              <span className='sr-only'>Menu</span>
              <MoreHorizontal className='h-4 w-4' />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align='end'>
            <DropdownMenuLabel>Ações</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
              <Edit className="mr-2 h-4 w-4" /> Editar Informações
            </DropdownMenuItem>
            <DropdownMenuItem>
              <Truck className="mr-2 h-4 w-4" /> Ver Aquisições
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="text-red-600">
              <Trash className="mr-2 h-4 w-4" /> Excluir Fornecedor
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      )
    },
  },
]
