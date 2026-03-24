import { createFileRoute } from '@tanstack/react-router'
import { LiveAuctionRoom } from '@/features/operations/auctions'

export const Route = createFileRoute('/_authenticated/operations/auctions')({
  component: LiveAuctionRoom,
})
