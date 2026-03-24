package websocket

import (
	"net/http"

	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	ws "github.com/gbrayhan/microservices-go/src/infrastructure/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Configurable CORS validation
	},
}

type IWebSocketController interface {
	ServeWS(ctx *gin.Context)
}

type WebSocketController struct {
	Hub    *ws.Hub
	Logger *logger.Logger
}

func NewWebSocketController(hub *ws.Hub, loggerInstance *logger.Logger) IWebSocketController {
	return &WebSocketController{
		Hub:    hub,
		Logger: loggerInstance,
	}
}

func (c *WebSocketController) ServeWS(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		c.Logger.Error("Failed to upgrade connection to WebSocket", zap.Error(err))
		return
	}
	client := &ws.Client{Hub: c.Hub, Conn: conn, Send: make(chan ws.Message, 256)}
	client.Hub.Register <- client

	// Start reading and writing asynchronously inside the goroutines.
	go client.WritePump()
	go client.ReadPump()
}
