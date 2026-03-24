import { Bar, BarChart, ResponsiveContainer, XAxis, YAxis } from 'recharts'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card' 
import { RecentActivities } from './recent-sales'

// Dados mais realistas para um antiquário
const data = [
  { name: 'Jan', total: 32000 },
  { name: 'Fev', total: 28000 },
  { name: 'Mar', total: 45000 },
  { name: 'Abr', total: 38000 },
  { name: 'Mai', total: 52000 },
  { name: 'Jun', total: 41000 },
  { name: 'Jul', total: 35000 },
  { name: 'Ago', total: 48000 },
  { name: 'Set', total: 43000 },
  { name: 'Out', total: 39000 },
  { name: 'Nov', total: 55000 },
  { name: 'Dez', total: 62000 },
]

export function Overview() {
  return (
    <>
      {/* KPIs Financeiros Principais */}
      <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-4'>
        <Card>
          <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
            <CardTitle className='text-sm font-medium'>Vendas do Mês</CardTitle>
            <svg
              xmlns='http://www.w3.org/2000/svg'
              viewBox='0 0 24 24'
              fill='none'
              stroke='currentColor'
              strokeLinecap='round'
              strokeLinejoin='round'
              strokeWidth='2'
              className='text-muted-foreground h-4 w-4'
            >
              <path d='M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6' />
            </svg>
          </CardHeader>
          <CardContent>
            <div className='text-2xl font-bold'>R$ 45.231,89</div>
            <p className='text-muted-foreground text-xs'>
              +20,1% em relação ao mês passado
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
            <CardTitle className='text-sm font-medium'>Margem Bruta</CardTitle>
            <svg
              xmlns='http://www.w3.org/2000/svg'
              viewBox='0 0 24 24'
              fill='none'
              stroke='currentColor'
              strokeLinecap='round'
              strokeLinejoin='round'
              strokeWidth='2'
              className='text-muted-foreground h-4 w-4'
            >
              <path d='M13 2L3 14h9l-1 8 10-12h-9l1-8z' />
            </svg>
          </CardHeader>
          <CardContent>
            <div className='text-2xl font-bold'>28,5%</div>
            <p className='text-muted-foreground text-xs'>
              +3,2% em relação ao mês passado
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
            <CardTitle className='text-sm font-medium'>Ticket Médio</CardTitle>
            <svg
              xmlns='http://www.w3.org/2000/svg'
              viewBox='0 0 24 24'
              fill='none'
              stroke='currentColor'
              strokeLinecap='round'
              strokeLinejoin='round'
              strokeWidth='2'
              className='text-muted-foreground h-4 w-4'
            >
              <path d='M22 12h-4l-3 9L9 3l-3 9H2' />
            </svg>
          </CardHeader>
          <CardContent>
            <div className='text-2xl font-bold'>R$ 1.923</div>
            <p className='text-muted-foreground text-xs'>
              +15% em relação ao mês passado
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
            <CardTitle className='text-sm font-medium'>
              Peças Vendidas
            </CardTitle>
            <svg
              xmlns='http://www.w3.org/2000/svg'
              viewBox='0 0 24 24'
              fill='none'
              stroke='currentColor'
              strokeLinecap='round'
              strokeLinejoin='round'
              strokeWidth='2'
              className='text-muted-foreground h-4 w-4'
            >
              <path d='M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2' />
              <circle cx='9' cy='7' r='4' />
              <path d='M22 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75' />
            </svg>
          </CardHeader>
          <CardContent>
            <div className='text-2xl font-bold'>24</div>
            <p className='text-muted-foreground text-xs'>
              +6 peças desde o mês passado
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Indicadores de Estoque */}
      <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-3'>
        <Card>
          <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
            <CardTitle className='text-sm font-medium'>
              Valor do Inventário
            </CardTitle>
            <svg
              xmlns='http://www.w3.org/2000/svg'
              viewBox='0 0 24 24'
              fill='none'
              stroke='currentColor'
              strokeLinecap='round'
              strokeLinejoin='round'
              strokeWidth='2'
              className='text-muted-foreground h-4 w-4'
            >
              <rect x='3' y='3' width='18' height='18' rx='2' ry='2' />
              <circle cx='9' cy='9' r='2' />
              <path d='M21 15l-3.086-3.086a2 2 0 0 0-2.828 0L6 21' />
            </svg>
          </CardHeader>
          <CardContent>
            <div className='text-2xl font-bold'>R$ 480.000</div>
            <p className='text-muted-foreground text-xs'>
              312 peças em estoque
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
            <CardTitle className='text-sm font-medium'>
              Em Consignação
            </CardTitle>
            <svg
              xmlns='http://www.w3.org/2000/svg'
              viewBox='0 0 24 24'
              fill='none'
              stroke='currentColor'
              strokeLinecap='round'
              strokeLinejoin='round'
              strokeWidth='2'
              className='text-muted-foreground h-4 w-4'
            >
              <path d='M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z' />
              <line x1='3' y1='6' x2='21' y2='6' />
              <path d='M16 10a4 4 0 0 1-8 0' />
            </svg>
          </CardHeader>
          <CardContent>
            <div className='text-2xl font-bold'>87 peças</div>
            <p className='text-muted-foreground text-xs'>
              R$ 120.000 em valor
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-2'>
            <CardTitle className='text-sm font-medium'>
              Estoque Parado
            </CardTitle>
            <svg
              xmlns='http://www.w3.org/2000/svg'
              viewBox='0 0 24 24'
              fill='none'
              stroke='currentColor'
              strokeLinecap='round'
              strokeLinejoin='round'
              strokeWidth='2'
              className='text-muted-foreground h-4 w-4'
            >
              <circle cx='12' cy='12' r='10' />
              <polyline points='12,6 12,12 16,14' />
            </svg>
          </CardHeader>
          <CardContent>
            <div className='text-2xl font-bold'>15 peças</div>
            <p className='text-muted-foreground text-xs'>
              +6 meses sem movimento
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Alertas Importantes
      <Card className='border-orange-200 bg-orange-50 dark:border-orange-800 dark:bg-orange-950'>
        <CardHeader className='pb-3'>
          <CardTitle className='flex items-center gap-2 text-sm font-medium text-orange-800 dark:text-orange-200'>
            <svg
              xmlns='http://www.w3.org/2000/svg'
              viewBox='0 0 24 24'
              fill='none'
              stroke='currentColor'
              strokeLinecap='round'
              strokeLinejoin='round'
              strokeWidth='2'
              className='h-4 w-4'
            >
              <path d='M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z' />
              <line x1='12' y1='9' x2='12' y2='13' />
              <line x1='12' y1='17' x2='12.01' y2='17' />
            </svg>
            Alertas Importantes
          </CardTitle>
        </CardHeader>
        <CardContent className='space-y-2'>
          <div className='text-sm text-orange-800 dark:text-orange-200'>
            <span className='font-medium'>3 consignações</span> vencem nos
            próximos 7 dias
          </div>
          <div className='text-sm text-orange-800 dark:text-orange-200'>
            <span className='font-medium'>5 peças</span> precisam de reavaliação
            de preço
          </div>
          <div className='text-sm text-orange-800 dark:text-orange-200'>
            <span className='font-medium'>8 peças</span> com informações
            incompletas
          </div>
        </CardContent>
      </Card> */}

      {/* Gráfico Principal + Atividades Recentes */}
      <div className='grid grid-cols-1 gap-4 lg:grid-cols-7'>
        <Card className='col-span-1 lg:col-span-4'>
          <CardHeader>
            <CardTitle>Vendas Mensais</CardTitle>
            <CardDescription>
              Evolução do faturamento ao longo do ano
            </CardDescription>
          </CardHeader>
          <CardContent className='ps-2'>
            <ResponsiveContainer width='100%' height={350}>
              <BarChart data={data}>
                <XAxis
                  dataKey='name'
                  stroke='#888888'
                  fontSize={12}
                  tickLine={false}
                  axisLine={false}
                />
                <YAxis
                  stroke='#888888'
                  fontSize={12}
                  tickLine={false}
                  axisLine={false}
                  tickFormatter={(value) => `R$${(value / 1000).toFixed(0)}k`}
                />
                <Bar
                  dataKey='total'
                  fill='currentColor'
                  radius={[4, 4, 0, 0]}
                  className='fill-primary'
                />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        <Card className='col-span-1 lg:col-span-3'>
          <CardHeader>
            <CardTitle>Atividades Recentes</CardTitle>
            <CardDescription>
              Últimas movimentações do seu antiquário
            </CardDescription>
          </CardHeader>
          <CardContent>
            <RecentActivities />
          </CardContent>
        </Card>
      </div>
    </>
  )
}
