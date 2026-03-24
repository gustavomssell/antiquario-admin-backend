import { useState, useRef, useEffect } from 'react'

import { cn } from '@/lib/utils'
import {
  Search as SearchIcon,
  Gavel,
  History,
  MoreVertical,
  CheckCircle2,
  CircleAlert,
} from 'lucide-react'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import { Badge } from '@/components/ui/badge'
import { useWebSocket } from '@/lib/use-websocket'

// Mock Data for Auction Lots (Equivalent to "Chats" list)
const auctionLots = [
  {
    id: 'LOTE-001',
    title: 'Vaso Ming Dinastia (Céramica Azul)',
    image: 'https://images.unsplash.com/photo-1613946069412-38f7f1ff0b65?auto=format&fit=crop&q=80&w=200',
    currentBid: 45000,
    status: 'LIVE',
    bids: [
      { sender: 'System', amount: 15000, time: '14:00' },
      { sender: 'Client #89', amount: 20000, time: '14:05' },
      { sender: 'Client #12', amount: 45000, time: '14:08' },
    ]
  },
  {
    id: 'LOTE-002',
    title: 'Relógio Patek Philippe 1940',
    image: 'https://images.unsplash.com/photo-1524805444758-089113d48a6d?auto=format&fit=crop&q=80&w=200',
    currentBid: 0,
    status: 'WAITING',
    bids: []
  },
  {
    id: 'LOTE-003',
    title: 'Poltrona Mole - Sergio Rodrigues',
    image: 'https://images.unsplash.com/photo-1567538096630-e0c55bd6374c?auto=format&fit=crop&q=80&w=200',
    currentBid: 120000,
    status: 'SOLD',
    bids: [
      { sender: 'System', amount: 80000, time: '12:00' },
      { sender: 'You', amount: 120000, time: '12:30' },
    ]
  }
]

export function LiveAuctionRoom() {
  const [search, setSearch] = useState('')
  const [selectedLot, setSelectedLot] = useState<typeof auctionLots[0] | null>(null)
  const [bidInput, setBidInput] = useState('')
  
  // Real-time hook integration
  const { isConnected, messages, sendMessage } = useWebSocket()
  const scrollRef = useRef<HTMLDivElement>(null)

  // Auto-scroll logic for bids
  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight
    }
  }, [messages, selectedLot?.bids])

  const handlePlaceBid = (e: React.FormEvent) => {
    e.preventDefault()
    if (!bidInput || !selectedLot) return

    const amount = parseFloat(bidInput)
    // Send via WebSocket (Simulated action merging with local state for now due to mock)
    sendMessage({
      type: 'PLACE_BID',
      action: 'BID',
      lotId: selectedLot.id,
      amount: amount,
      timestamp: new Date().toISOString(),
    })
    
    setSelectedLot(prev => {
      if (!prev) return prev
      return {
        ...prev,
        currentBid: amount,
        bids: [...prev.bids, { sender: 'You', amount, time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) }]
      }
    })
    setBidInput('')
  }

  const filteredLots = auctionLots.filter((lot) =>
    lot.title.toLowerCase().includes(search.toLowerCase()) || lot.id.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <div className='flex h-[calc(100vh-8rem)] w-full overflow-hidden rounded-xl border bg-background shadow-xl mt-4 max-w-7xl mx-auto'>
      
      {/* LEFT PANEL: List of Lots (Equivalent to Chats List) */}
      <div className='flex w-full flex-col sm:w-80 border-r bg-muted/10'>
        <div className='p-4 border-b bg-card'>
          <div className='flex items-center justify-between mb-4'>
            <div className='flex items-center gap-2'>
              <Gavel className="w-5 h-5 text-primary" />
              <h1 className='text-lg font-bold'>Lotes do Leilão</h1>
            </div>
            {isConnected ? (
              <Badge className="bg-green-600 px-2 py-0 text-[10px]"><CheckCircle2 className="w-3 h-3 mr-1" /> ON</Badge>
            ) : (
              <Badge variant="destructive" className="px-2 py-0 text-[10px]"><CircleAlert className="w-3 h-3 mr-1" /> OFF</Badge>
            )}
          </div>
          <div className='relative'>
            <SearchIcon className='absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground' />
            <input
              type='text'
              className='w-full rounded-md border border-input bg-background pl-9 pr-4 py-2 text-sm focus:outline-hidden focus:ring-1 focus:ring-primary'
              placeholder='Buscar Lote ou Peça...'
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </div>
        </div>

        <ScrollArea className='flex-1'>
          <div className="p-2 space-y-1">
            {filteredLots.map((lot) => (
              <button
                key={lot.id}
                onClick={() => setSelectedLot(lot)}
                className={cn(
                  'flex w-full items-center gap-3 rounded-lg px-3 py-3 text-left transition-colors hover:bg-muted',
                  selectedLot?.id === lot.id && 'bg-primary/10 hover:bg-primary/15'
                )}
              >
                <Avatar className="h-10 w-10 border">
                  <AvatarImage src={lot.image} alt={lot.id} />
                  <AvatarFallback>{lot.id.substring(5, 8)}</AvatarFallback>
                </Avatar>
                <div className="flex-1 overflow-hidden">
                  <div className="flex justify-between items-center mb-1">
                    <span className="text-sm font-semibold truncate">{lot.id}</span>
                    <span className={cn(
                      "text-[10px] font-bold px-1.5 py-0.5 rounded-sm",
                      lot.status === 'LIVE' && "bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300 animate-pulse",
                      lot.status === 'WAITING' && "bg-slate-100 text-slate-700",
                      lot.status === 'SOLD' && "bg-green-100 text-green-700"
                    )}>
                      {lot.status}
                    </span>
                  </div>
                  <p className="text-xs text-muted-foreground truncate">{lot.title}</p>
                  <p className="text-xs font-medium text-primary mt-1">
                    Lance: {new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(lot.currentBid)}
                  </p>
                </div>
              </button>
            ))}
          </div>
        </ScrollArea>
      </div>

      {/* RIGHT PANEL: Live Bid Feed (Equivalent to Chat Room) */}
      {selectedLot ? (
        <div className='flex flex-1 flex-col bg-card'>
          
          {/* Top Header of the Room */}
          <div className='flex items-center justify-between border-b px-6 py-4 shadow-xs'>
            <div className='flex items-center gap-4'>
              <Avatar className="h-12 w-12 border-2 border-primary/20">
                <AvatarImage src={selectedLot.image} />
                <AvatarFallback>{selectedLot.id}</AvatarFallback>
              </Avatar>
              <div>
                <h2 className='text-lg font-bold leading-tight'>{selectedLot.title}</h2>
                <div className="flex items-center text-sm text-muted-foreground mt-0.5">
                  <span className="font-mono">{selectedLot.id}</span>
                  <Separator orientation="vertical" className="h-3 mx-2" />
                  <span>Lance Vigente: <strong className="text-foreground">{new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(selectedLot.currentBid)}</strong></span>
                </div>
              </div>
            </div>
            <div>
              <Button variant='ghost' size='icon'>
                <MoreVertical className='h-5 w-5 text-muted-foreground' />
              </Button>
            </div>
          </div>

          {/* Conversation / Bid History */}
          <div className='flex-1 overflow-y-auto p-4 bg-muted/5' ref={scrollRef}>
             <div className="flex flex-col space-y-4">
                <div className="text-center my-4">
                   <Badge variant="outline" className="text-xs font-normal text-muted-foreground bg-background">
                     Leilão Aberto Pelo Leiloeiro
                   </Badge>
                </div>

                {selectedLot.bids.map((bid, index) => {
                  const isYou = bid.sender === 'You'
                  const isSystem = bid.sender === 'System'

                  if (isSystem) {
                    return (
                      <div key={index} className="flex justify-center my-2">
                        <div className="bg-primary/5 border border-primary/10 rounded-full px-4 py-1 text-xs text-primary flex items-center">
                          <History className="w-3 h-3 mr-2" />
                          Lance Base Registrado: <strong>{new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(bid.amount)}</strong>
                        </div>
                      </div>
                    )
                  }

                  return (
                    <div key={index} className={cn("flex w-full", isYou ? "justify-end" : "justify-start")}>
                      <div className={cn(
                        "max-w-[70%] rounded-2xl px-4 py-2 shadow-sm relative",
                        isYou ? "bg-primary text-primary-foreground rounded-tr-sm" : "bg-background border rounded-tl-sm"
                      )}>
                        <div className="flex justify-between items-baseline mb-1 gap-4">
                          <span className={cn("text-[10px] font-bold uppercase tracking-wider", isYou ? "text-primary-foreground/80" : "text-muted-foreground")}>
                            {bid.sender}
                          </span>
                          <span className={cn("text-[10px]", isYou ? "text-primary-foreground/60" : "text-muted-foreground/60")}>
                            {bid.time}
                          </span>
                        </div>
                        <div className="flex items-center gap-2">
                           <Gavel className={cn("h-4 w-4", isYou ? "text-primary-foreground/80" : "text-green-600")} />
                           <p className="text-xl font-bold tracking-tight">
                             {new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(bid.amount)}
                           </p>
                        </div>
                      </div>
                    </div>
                  )
                })}

                {/* Incoming WebSocket Messages visually appended here for effect */}
                {messages.map((msg, idx) => (
                   <div key={`ws-${idx}`} className="flex justify-start">
                      <div className="bg-background border border-blue-200 rounded-2xl rounded-tl-sm px-4 py-2 shadow-sm max-w-[70%]">
                        <div className="flex justify-between items-baseline mb-1 gap-4">
                          <span className="text-[10px] font-bold uppercase tracking-wider text-blue-600">
                            {msg.sender || 'Participant'}
                          </span>
                          <span className="text-[10px] text-muted-foreground/60">Ao Vivo</span>
                        </div>
                        <p className="text-lg font-bold tracking-tight text-foreground">
                          {msg.content || JSON.stringify(msg)}
                        </p>
                      </div>
                   </div>
                ))}
             </div>
          </div>

          {/* Action Bar (Send Bid) */}
          <div className='p-4 bg-background border-t'>
            <form onSubmit={handlePlaceBid} className='flex items-center gap-2'>
              <div className='flex items-center pl-3 pr-1 py-1 flex-1 bg-muted/30 border rounded-full focus-within:ring-1 focus-within:ring-primary focus-within:bg-background transition-all'>
                <span className="text-muted-foreground font-semibold">R$</span>
                <input
                  type='number'
                  placeholder={`Mínimo: ${selectedLot.currentBid + 1000}`}
                  className='w-full bg-transparent px-2 py-2 focus:outline-hidden font-medium'
                  value={bidInput}
                  onChange={(e) => setBidInput(e.target.value)}
                  disabled={selectedLot.status !== 'LIVE' || !isConnected}
                />
              </div>
              <Button 
                type='submit' 
                size='icon' 
                className='h-12 w-12 rounded-full shadow-lg bg-blue-600 hover:bg-blue-700'
                disabled={selectedLot.status !== 'LIVE' || !bidInput || !isConnected}
              >
                <Gavel className='h-5 w-5' />
              </Button>
            </form>
            {selectedLot.status !== 'LIVE' && (
              <p className="text-center text-xs text-muted-foreground mt-2">
                O leilão para este lote está {selectedLot.status === 'SOLD' ? 'Encerrado' : 'Aguardando Início'}.
              </p>
            )}
          </div>

        </div>
      ) : (
        <div className='flex flex-1 flex-col items-center justify-center bg-card text-center p-8'>
          <div className='flex h-20 w-20 items-center justify-center rounded-full bg-muted mb-6 shadow-inner'>
            <Gavel className='h-10 w-10 text-muted-foreground/50' />
          </div>
          <h2 className='text-2xl font-semibold'>Saguão de Leilões</h2>
          <p className='text-muted-foreground mt-2 max-w-sm'>
            Selecione um lote no painel à esquerda para acompanhar os lances ao vivo ou efetuar uma oferta.
          </p>
        </div>
      )}
    </div>
  )
}
