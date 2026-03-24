import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { ConfigDrawer } from '@/components/config-drawer'
import { Header } from '@/components/layout/header'
import { Main } from '@/components/layout/main'
import { TopNav } from '@/components/layout/top-nav'
import { ProfileDropdown } from '@/components/profile-dropdown'
import { Search } from '@/components/search'
import { ThemeSwitch } from '@/components/theme-switch'
import { Analytics } from './components/analytics'
import { Overview } from './components/overview'

const topNav = [
  {
    title: 'Visão Geral',
    href: 'dashboard/overview',
    isActive: true,
    disabled: false,
  },
  {
    title: 'Clientes',
    href: 'dashboard/customers',
    isActive: false,
    disabled: true,
  },
  {
    title: 'Produtos',
    href: 'dashboard/products',
    isActive: false,
    disabled: true,
  },
  {
    title: 'Configurações',
    href: 'dashboard/settings',
    isActive: false,
    disabled: true,
  },
]

export function Dashboard() {
  return (
    <>
      {/* ===== Top Heading ===== */}
      <Header fixed>
        <TopNav links={topNav} />
        <div className='ms-auto flex items-center space-x-4'>
          <Search />
          <ThemeSwitch />
          <ConfigDrawer />
          <ProfileDropdown />
        </div>
      </Header>

      {/* ===== Main ===== */}
      <Main>
        <div className='mb-2 flex items-center justify-between space-y-2'>
          <h1 className='text-2xl font-bold tracking-tight'>Dashboard</h1>
          <div className='flex items-center space-x-2'>
            <Button>Relatório</Button>
          </div>
        </div>
        <Tabs
          orientation='vertical'
          defaultValue='overview'
          className='space-y-4'
        >
          <div className='w-full overflow-x-auto pb-2'>
            <TabsList>
              <TabsTrigger value='overview'>Visão Geral</TabsTrigger>
              <TabsTrigger value='analytics'>Análises</TabsTrigger>
              <TabsTrigger value='reports' disabled>
                Relatórios
              </TabsTrigger>
              <TabsTrigger value='notifications' disabled>
                Notificações
              </TabsTrigger>
            </TabsList>
          </div>
          <TabsContent value='overview' className='space-y-4'>
            <Overview />
          </TabsContent>
          <TabsContent value='analytics' className='space-y-4'>
            <Analytics />
          </TabsContent>
        </Tabs>
      </Main>
    </>
  )
}
