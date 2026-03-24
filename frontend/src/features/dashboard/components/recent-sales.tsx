export function RecentActivities() {
  return (
    <div className='space-y-8'>
      {/* Venda Recente */}
      <div className='flex items-center gap-4'>
        <div className='flex h-9 w-9 items-center justify-center rounded-full bg-green-100 dark:bg-green-900'>
          <svg
            xmlns='http://www.w3.org/2000/svg'
            viewBox='0 0 24 24'
            fill='none'
            stroke='currentColor'
            strokeLinecap='round'
            strokeLinejoin='round'
            strokeWidth='2'
            className='h-4 w-4 text-green-600 dark:text-green-400'
          >
            <path d='M12 2v20M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6' />
          </svg>
        </div>
        <div className='flex flex-1 flex-wrap items-center justify-between'>
          <div className='space-y-1'>
            <p className='text-sm leading-none font-medium'>
              Venda: Relógio de Bolso Omega
            </p>
            <p className='text-muted-foreground text-xs'>
              Cliente: João Silva • há 2 horas
            </p>
          </div>
          <div className='font-medium text-green-600'>+R\$ 3.200</div>
        </div>
      </div>

      {/* Nova Peça Adicionada */}
      <div className='flex items-center gap-4'>
        <div className='flex h-9 w-9 items-center justify-center rounded-full bg-blue-100 dark:bg-blue-900'>
          <svg
            xmlns='http://www.w3.org/2000/svg'
            viewBox='0 0 24 24'
            fill='none'
            stroke='currentColor'
            strokeLinecap='round'
            strokeLinejoin='round'
            strokeWidth='2'
            className='h-4 w-4 text-blue-600 dark:text-blue-400'
          >
            <circle cx='12' cy='12' r='10' />
            <line x1='12' y1='8' x2='12' y2='16' />
            <line x1='8' y1='12' x2='16' y2='12' />
          </svg>
        </div>
        <div className='flex flex-1 flex-wrap items-center justify-between'>
          <div className='space-y-1'>
            <p className='text-sm leading-none font-medium'>
              Nova peça: Vaso Ming Dynasty
            </p>
            <p className='text-muted-foreground text-xs'>
              Porcelana • Século XVIII • há 4 horas
            </p>
          </div>
          <div className='font-medium text-blue-600'>R\$ 8.500</div>
        </div>
      </div>

      {/* Atualização de Preço */}
      <div className='flex items-center gap-4'>
        <div className='flex h-9 w-9 items-center justify-center rounded-full bg-orange-100 dark:bg-orange-900'>
          <svg
            xmlns='http://www.w3.org/2000/svg'
            viewBox='0 0 24 24'
            fill='none'
            stroke='currentColor'
            strokeLinecap='round'
            strokeLinejoin='round'
            strokeWidth='2'
            className='h-4 w-4 text-orange-600 dark:text-orange-400'
          >
            <path d='M3 3v5h5M3 8l4-4 4 4 4-4 4 4v5' />
            <path d='M21 21v-5h-5M21 16l-4 4-4-4-4 4-4-4v-5' />
          </svg>
        </div>
        <div className='flex flex-1 flex-wrap items-center justify-between'>
          <div className='space-y-1'>
            <p className='text-sm leading-none font-medium'>
              Preço ajustado: Escultura Bronze
            </p>
            <p className='text-muted-foreground text-xs'>
              De R\$ 5.000 para R\$ 4.500 • há 6 horas
            </p>
          </div>
          <div className='font-medium text-orange-600'>-10%</div>
        </div>
      </div>

      {/* Novo Cliente */}
      <div className='flex items-center gap-4'>
        <div className='flex h-9 w-9 items-center justify-center rounded-full bg-purple-100 dark:bg-purple-900'>
          <svg
            xmlns='http://www.w3.org/2000/svg'
            viewBox='0 0 24 24'
            fill='none'
            stroke='currentColor'
            strokeLinecap='round'
            strokeLinejoin='round'
            strokeWidth='2'
            className='h-4 w-4 text-purple-600 dark:text-purple-400'
          >
            <path d='M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2' />
            <circle cx='9' cy='7' r='4' />
            <path d='M22 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75' />
          </svg>
        </div>
        <div className='flex flex-1 flex-wrap items-center justify-between'>
          <div className='space-y-1'>
            <p className='text-sm leading-none font-medium'>
              Novo cliente: Maria Antunes
            </p>
            <p className='text-muted-foreground text-xs'>
              Interesse em móveis coloniais • há 1 dia
            </p>
          </div>
          <div className='font-medium text-purple-600'>Novo</div>
        </div>
      </div>

      {/* Peça Finalizada Restauração */}
      <div className='flex items-center gap-4'>
        <div className='flex h-9 w-9 items-center justify-center rounded-full bg-emerald-100 dark:bg-emerald-900'>
          <svg
            xmlns='http://www.w3.org/2000/svg'
            viewBox='0 0 24 24'
            fill='none'
            stroke='currentColor'
            strokeLinecap='round'
            strokeLinejoin='round'
            strokeWidth='2'
            className='h-4 w-4 text-emerald-600 dark:text-emerald-400'
          >
            <polyline points='20,6 9,17 4,12' />
          </svg>
        </div>
        <div className='flex flex-1 flex-wrap items-center justify-between'>
          <div className='space-y-1'>
            <p className='text-sm leading-none font-medium'>
              Restauração concluída: Quadro a Óleo
            </p>
            <p className='text-muted-foreground text-xs'>
              Século XIX • Pronto para venda • há 1 dia
            </p>
          </div>
          <div className='font-medium text-emerald-600'>Concluído</div>
        </div>
      </div>

      {/* Consulta de Cliente */}
      <div className='flex items-center gap-4'>
        <div className='flex h-9 w-9 items-center justify-center rounded-full bg-gray-100 dark:bg-gray-800'>
          <svg
            xmlns='http://www.w3.org/2000/svg'
            viewBox='0 0 24 24'
            fill='none'
            stroke='currentColor'
            strokeLinecap='round'
            strokeLinejoin='round'
            strokeWidth='2'
            className='h-4 w-4 text-gray-600 dark:text-gray-400'
          >
            <circle cx='11' cy='11' r='8' />
            <path d='m21 21-4.35-4.35' />
          </svg>
        </div>
        <div className='flex flex-1 flex-wrap items-center justify-between'>
          <div className='space-y-1'>
            <p className='text-sm leading-none font-medium'>
              Consulta: Cristaleira Francesa
            </p>
            <p className='text-muted-foreground text-xs'>
              Cliente: Pedro Santos • há 2 dias
            </p>
          </div>
          <div className='font-medium text-gray-600'>Visualizado</div>
        </div>
      </div>
    </div>
  )
}
