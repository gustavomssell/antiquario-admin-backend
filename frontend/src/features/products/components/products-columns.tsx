import { ColumnDef } from '@tanstack/react-table'
import { Product } from '../api/use-products'
import { DataTableColumnHeader } from '@/components/data-table/column-header'
import { MoreHorizontal, Edit, Trash, ExternalLink } from 'lucide-react'
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
import { Link } from '@tanstack/react-router'

export const columns: ColumnDef<Product>[] = [
  {
    accessorKey: 'sku',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='SKU' />
    ),
    cell: ({ row }) => <span className="font-mono text-xs">{row.original.sku}</span>
  },
  {
    accessorKey: 'title',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='Título da Peça' />
    ),
    cell: ({ row }) => <span className="font-medium">{row.original.title}</span>
  },
  {
    accessorKey: 'status',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='Status' />
    ),
    cell: ({ row }) => {
      const status = row.original.status
      if (status === 'AVAILABLE') return <Badge variant="default" className="bg-green-600">Disponível</Badge>
      if (status === 'RESERVED') return <Badge variant="secondary" className="bg-yellow-600 text-white">Reservado</Badge>
      if (status === 'SOLD') return <Badge variant="outline">Vendido</Badge>
      if (status === 'IN_RESTORATION') return <Badge variant="destructive">Em Restauro</Badge>
      return <Badge>{status}</Badge>
    }
  },
  {
    accessorKey: 'currentPrice',
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title='Preço Atual' />
    ),
    cell: ({ row }) => {
      const price = parseFloat(row.original.currentPrice.toString())
      const formatted = new Intl.NumberFormat('pt-BR', {
        style: 'currency',
        currency: 'BRL',
      }).format(price)
      return <span>{formatted}</span>
    }
  },
  {
    id: 'actions',
    cell: ({ row }) => {
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
            <DropdownMenuItem asChild>
              {/* @ts-expect-error Tanstack router lacks strict route mappings here */}
              <Link to={`/products/${row.original.id}/edit`}>
                <Edit className="mr-2 h-4 w-4" /> Editar Ficha Técnica
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem>
              <ExternalLink className="mr-2 h-4 w-4" /> Visualizar na Loja
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="text-red-600">
              <Trash className="mr-2 h-4 w-4" /> Excluir
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      )
    },
  },
]
