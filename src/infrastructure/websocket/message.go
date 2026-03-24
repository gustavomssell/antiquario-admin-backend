package websocket

type Message struct {
	Type    string      `json:"type"`    // e.g., "auction_bid", "notification"
	Payload interface{} `json:"payload"` // the actual data payload
}
