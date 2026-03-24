import { useEffect, useRef, useState } from 'react'
import { useAuthStore } from '@/stores/auth-store'
import { toast } from 'sonner'


// Using env or localhost standard
const WS_URL = 'ws://localhost:8080/v1/ws'

export function useWebSocket() {
  const [isConnected, setIsConnected] = useState(false)
  const [messages, setMessages] = useState<any[]>([])
  const wsRef = useRef<WebSocket | null>(null)
  const { auth } = useAuthStore()

  useEffect(() => {
    // We can append token to URL or just connect
    const ws = new WebSocket(`${WS_URL}?token=${auth.accessToken}`)
    wsRef.current = ws

    ws.onopen = () => {
      setIsConnected(true)
      toast.success('Sincronizado ao Hub Principal (Ao Vivo)')
    }

    ws.onmessage = (event) => {
      try {
        const payload = JSON.parse(event.data)
        setMessages((prev) => [...prev, payload])
      } catch (e) {
        setMessages((prev) => [...prev, { type: 'TEXT', content: event.data, timestamp: new Date().toISOString() }])
      }
    }

    ws.onclose = () => {
      setIsConnected(false)
      toast.error('Conexão em Tempo Real perdida.')
    }

    return () => {
      ws.close()
    }
  }, [auth.accessToken])

  const sendMessage = (message: any) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message))
    } else {
      toast.error('Desconectado! Impossível enviar comando.')
    }
  }

  return { isConnected, messages, sendMessage }
}
