import type { ColumnDef } from '@tanstack/react-table'
import type { User } from '../api/use-users'
import { DataTableColumnHeader } from '@/components/data-table/column-header'
import { MoreHorizontal, Edit, Trash, ShieldAlert } from 'lucide-react'
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

export const columns: ColumnDef<User>[] = [
  {
    accessorKey: 'userName',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='Username / Credencial' />
    ),
    cell: ({ row }) => <span className="font-semibold">{row.original.userName}</span>
  },
  {
    accessorKey: 'firstName',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='Nome Completo' />
    ),
    cell: ({ row }) => <span>{row.original.firstName} {row.original.lastName}</span>
  },
  {
    accessorKey: 'email',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='E-mail Contato' />
    ),
    cell: ({ row }) => <span className="text-muted-foreground">{row.original.email}</span>
  },
  {
    accessorKey: 'status',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='Status Conta' />
    ),
    cell: ({ row }) => {
      const isActive = row.original.status
      return isActive ? (
        <Badge variant="outline" className="bg-green-100 text-green-700">Ativo</Badge>
      ) : (
        <Badge variant="destructive">Desativado</Badge>
      )
    }
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
            <DropdownMenuLabel>Gestão</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
              <Edit className="mr-2 h-4 w-4" /> Editar Perfil
            </DropdownMenuItem>
            <DropdownMenuItem>
              <ShieldAlert className="mr-2 h-4 w-4" /> Redefinir Senha
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="text-red-600">
              <Trash className="mr-2 h-4 w-4" /> Revogar Acesso
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      )
    },
  },
]
