package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kiga-hub/websocket/pkg/utils"
	"github.com/labstack/echo/v4"
)

// WebsocketServer -
type WebsocketServer struct {
	clients  map[*Client]bool
	upGrader websocket.Upgrader
}

// NewServer -
func NewServer() *WebsocketServer {
	return &WebsocketServer{
		clients:  make(map[*Client]bool),
		upGrader: websocket.Upgrader{},
	}
}

// SocketMessage this is socket msg struct
type SocketMessage struct {
	Message string `json:"message"`
}

// HandleConnections -
func (server *WebsocketServer) HandleConnections(c echo.Context) error {
	ws, err := server.upGrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	client := NewBroadcastClient(ws, server)
	server.clients[client] = true

	go client.Read()

	return c.JSON(http.StatusOK, utils.SuccessJSONData(nil))
}

// RemoveClient -
func (server *WebsocketServer) RemoveClient(client *Client) {
	delete(server.clients, client)
}

// RemoveAllConn -
func (server *WebsocketServer) RemoveAllConn() {
	for client := range server.clients {
		if err := client.Close(); err != nil {
			return
		}
		delete(server.clients, client)
	}
}
