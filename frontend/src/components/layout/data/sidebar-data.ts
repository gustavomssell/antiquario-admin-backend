import {
  LayoutDashboard,
  Package,
  Layers,
  Clock,
  Palette,
  Droplets,
  Tags,
  Users,
  Truck,
  ShoppingCart,
  CreditCard,
  Wrench,
  SearchCheck,
  CalendarCheck,
  Gavel,
} from 'lucide-react'
import { type SidebarData } from '../types'

export const sidebarData: SidebarData = {
  user: {
    name: 'Administrador',
    email: 'admin@antiquaos.com',
    avatar: '',
  },
  teams: [
    {
      name: 'AntiquaOS',
      logo: Package,
      plan: 'Enterprise',
    },
  ],
  navGroups: [
    {
      title: 'Principal',
      items: [
        {
          title: 'Dashboard',
          url: '/',
          icon: LayoutDashboard,
        },
        {
          title: 'Inventário (Produtos)',
          url: '/products',
          icon: Package,
        },
      ],
    },
    {
      title: 'Catálogo Base',
      items: [
        {
          title: 'Categorias',
          url: '/catalog/categories',
          icon: Layers,
        },
        {
          title: 'Períodos Históricos',
          url: '/catalog/periods',
          icon: Clock,
        },
        {
          title: 'Estilos de Época',
          url: '/catalog/styles',
          icon: Palette,
        },
        {
          title: 'Materiais',
          url: '/catalog/materials',
          icon: Droplets,
        },
        {
          title: 'Tags e Identificadores',
          url: '/catalog/tags',
          icon: Tags,
        },
      ],
    },
    {
      title: 'Comercial e Transações',
      items: [
        {
          title: 'Clientes',
          url: '/commercial/customers',
          icon: Users,
        },
        {
          title: 'Fornecedores',
          url: '/commercial/suppliers',
          icon: Truck,
        },
        {
          title: 'Vendas PVD',
          url: '/commercial/sales',
          icon: ShoppingCart,
        },
        {
          title: 'Aquisições',
          url: '/commercial/acquisitions',
          icon: CreditCard,
        },
      ],
    },
    {
      title: 'Centro de Operações e Restauro',
      items: [
        {
          title: 'Ordens de Restauro',
          url: '/operations/restorations',
          icon: Wrench,
        },
        {
          title: 'Avaliações Técnicas',
          url: '/operations/appraisals',
          icon: SearchCheck,
        },
        {
          title: 'Reservas Flexíveis',
          url: '/operations/reservations',
          icon: CalendarCheck,
        },
        {
          title: 'Saguão de Leilões',
          url: '/operations/auctions',
          badge: 'Live',
          icon: Gavel,
        },
      ],
    },
  ],
}
