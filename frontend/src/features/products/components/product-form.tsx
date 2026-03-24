// @ts-nocheck
import { useState } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate } from '@tanstack/react-router'
import { toast } from 'sonner'
import { Save, ArrowLeft } from 'lucide-react'

import { useCreateProduct } from '../api/use-products'
import { useCategories } from '@/features/catalog/api/use-categories'
import { usePeriods } from '@/features/catalog/api/use-periods'
import { useStyles } from '@/features/catalog/api/use-styles'

import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'

const productSchema = z.object({
  sku: z.string().min(1, 'SKU é obrigatório'),
  title: z.string().min(1, 'Título é obrigatório'),
  description: z.string().optional(),
  condition: z.string(),
  status: z.string(),
  location: z.string().optional(),
  dimensions: z.string().optional(),
  weight: z.coerce.number().min(0).optional(),
  acquisitionCost: z.coerce.number().min(0),
  basePrice: z.coerce.number().min(0),
  currentPrice: z.coerce.number().min(0),
  categoryId: z.string().min(1, 'Categoria é obrigatória'),
  periodId: z.string().optional(),
  styleId: z.string().optional(),
  manufacturingYear: z.coerce.number().optional(),
  isConsigned: z.boolean().default(false),
})

type ProductFormValues = z.infer<typeof productSchema>

export function ProductForm() {
  const navigate = useNavigate()
  const createMutation = useCreateProduct()
  
  // Load Catalog data for Selectors
  const { data: categoriesData } = useCategories(1, 100)
  const { data: periodsData } = usePeriods(1, 100)
  const { data: stylesData } = useStyles(1, 100)

  const form = useForm<ProductFormValues>({
    resolver: zodResolver(productSchema),
    defaultValues: {
      sku: '',
      title: '',
      description: '',
      condition: 'EXCELLENT',
      status: 'AVAILABLE',
      location: 'Showroom Principal',
      dimensions: '',
      weight: 0,
      acquisitionCost: 0,
      basePrice: 0,
      currentPrice: 0,
      categoryId: '',
      periodId: '',
      styleId: '',
      manufacturingYear: undefined,
      isConsigned: false,
    },
  })

  function onSubmit(values: ProductFormValues) {
    createMutation.mutate(values, {
      onSuccess: () => {
        toast.success('Peça catalogada e adicionada ao inventário!')
        navigate({ to: '/products' })
      },
      onError: () => {
        toast.error('Erro ao salvar o produto.')
      },
    })
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className='space-y-8'>
        <div className='flex items-center justify-between'>
          <div>
            <h2 className='text-2xl font-bold tracking-tight'>Nova Ficha Técnica</h2>
            <p className='text-muted-foreground'>
              Cadastre todos os detalhes físicos, históricos e comerciais da peça.
            </p>
          </div>
          <div className='flex space-x-3'>
            <Button variant="outline" type="button" onClick={() => navigate({ to: '/products' })}>
              <ArrowLeft className='mr-2 h-4 w-4' /> Cancelar
            </Button>
            <Button type='submit' disabled={createMutation.isPending}>
              {createMutation.isPending ? 'Salvando...' : (
                <><Save className='mr-2 h-4 w-4' /> Salvar no Acervo</>
              )}
            </Button>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          {/* Informações Básicas */}
          <Card>
            <CardHeader>
              <CardTitle>Identificação Básica</CardTitle>
              <CardDescription>Dados principais da antiguidade.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <FormField
                  control={form.control}
                  name='sku'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Código / SKU</FormLabel>
                      <FormControl><Input placeholder='EX: ANTIQ-001' {...field} /></FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name='status'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Status no Acervo</FormLabel>
                      <Select onValueChange={field.onChange} defaultValue={field.value}>
                        <FormControl><SelectTrigger><SelectValue placeholder="Status" /></SelectTrigger></FormControl>
                        <SelectContent>
                          <SelectItem value="AVAILABLE">Disponível</SelectItem>
                          <SelectItem value="RESERVED">Reservado (Sinal)</SelectItem>
                          <SelectItem value="IN_RESTORATION">Em Restauro</SelectItem>
                          <SelectItem value="SOLD">Vendido</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <FormField
                control={form.control}
                name='title'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Título/Nome da Peça</FormLabel>
                    <FormControl><Input placeholder='Vaso Ming Dinastia...' {...field} /></FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name='description'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Laudo de Descrição Curatorial</FormLabel>
                    <FormControl><Textarea className="h-32" placeholder='Detalhes históricos, marcas de uso...' {...field} /></FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>

          {/* Classificação e História */}
          <Card>
            <CardHeader>
              <CardTitle>Taxonomia e História</CardTitle>
              <CardDescription>Classificação oficial do catálogo.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <FormField
                control={form.control}
                name='categoryId'
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Categoria Base</FormLabel>
                    <Select onValueChange={field.onChange} defaultValue={field.value}>
                      <FormControl><SelectTrigger><SelectValue placeholder="Selecione Categoria" /></SelectTrigger></FormControl>
                      <SelectContent>
                        {categoriesData?.data.map((cat) => (
                          <SelectItem key={cat.id} value={cat.id}>{cat.name}</SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className="grid grid-cols-2 gap-4">
                <FormField
                  control={form.control}
                  name='periodId'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Período Histórico</FormLabel>
                      <Select onValueChange={field.onChange} defaultValue={field.value}>
                        <FormControl><SelectTrigger><SelectValue placeholder="Opcional" /></SelectTrigger></FormControl>
                        <SelectContent>
                          {periodsData?.data.map((p) => (
                            <SelectItem key={p.id} value={p.id}>{p.name}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name='styleId'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Estilo Artístico</FormLabel>
                      <Select onValueChange={field.onChange} defaultValue={field.value}>
                        <FormControl><SelectTrigger><SelectValue placeholder="Opcional" /></SelectTrigger></FormControl>
                        <SelectContent>
                          {stylesData?.data.map((s) => (
                            <SelectItem key={s.id} value={s.id}>{s.name}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <FormField
                  control={form.control}
                  name='condition'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Estado de Conservação</FormLabel>
                      <Select onValueChange={field.onChange} defaultValue={field.value}>
                        <FormControl><SelectTrigger><SelectValue placeholder="Classificação" /></SelectTrigger></FormControl>
                        <SelectContent>
                          <SelectItem value="MINT">Perfeito (Mint)</SelectItem>
                          <SelectItem value="EXCELLENT">Excelente</SelectItem>
                          <SelectItem value="GOOD">Bom (Marcas de Uso)</SelectItem>
                          <SelectItem value="FAIR">Razoável</SelectItem>
                          <SelectItem value="POOR">Necessita Restauro</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name='manufacturingYear'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Ano Aproximado</FormLabel>
                      <FormControl><Input type="number" placeholder='1890' {...field} /></FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            </CardContent>
          </Card>

          {/* Comercial e Valores */}
          <Card className="md:col-span-2">
            <CardHeader>
              <CardTitle>Financeiro e Estoque</CardTitle>
              <CardDescription>Valores, margens e logística.</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                <FormField
                  control={form.control}
                  name='acquisitionCost'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Custo de Aquisição (R$)</FormLabel>
                      <FormControl><Input type="number" step="0.01" {...field} /></FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name='basePrice'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Preço de Catálogo (R$)</FormLabel>
                      <FormControl><Input type="number" step="0.01" {...field} /></FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name='currentPrice'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Preço de Venda Praticado (R$)</FormLabel>
                      <FormControl><Input type="number" step="0.01" className="font-bold text-green-700" {...field} /></FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name='location'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Prateleira / Depósito</FormLabel>
                      <FormControl><Input placeholder='Galeria A, Prateleira 2' {...field} /></FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
              
              <Separator className="my-6" />

              <div className="flex items-center space-x-8">
                <FormField
                  control={form.control}
                  name='isConsigned'
                  render={({ field }) => (
                    <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4 shadow-sm w-80">
                      <div className="space-y-0.5">
                        <FormLabel className="text-base">Peça Consignada?</FormLabel>
                        <FormDescription>Ative se a peça pertence a terceiros.</FormDescription>
                      </div>
                      <FormControl>
                        <Switch checked={field.value} onCheckedChange={field.onChange} />
                      </FormControl>
                    </FormItem>
                  )}
                />
                
                <div className="flex-1 grid grid-cols-2 gap-4">
                   <FormField
                    control={form.control}
                    name='dimensions'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Dimensões Físicas</FormLabel>
                        <FormControl><Input placeholder='AxLxP (ex: 40x20x15 cm)' {...field} /></FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <FormField
                    control={form.control}
                    name='weight'
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Peso Bruto (Kg)</FormLabel>
                        <FormControl><Input type="number" step="0.01" {...field} /></FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </form>
    </Form>
  )
}
