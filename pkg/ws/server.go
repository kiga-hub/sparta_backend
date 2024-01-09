package ws

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kiga-hub/websocket/pkg/utils"
	"github.com/labstack/echo/v4"
)

type WebsocketServer struct {
	clients   map[*Client]bool
	broadcast chan []byte
	upgrader  websocket.Upgrader
}

func NewServer() *WebsocketServer {
	return &WebsocketServer{
		clients:   make(map[*Client]bool),
		broadcast: make(chan []byte),
		upgrader:  websocket.Upgrader{},
	}
}

// SocketMessage Define our message object
type SocketMessage struct {
	Message string `json:"message"`
}

func (server *WebsocketServer) HandleConnections(c echo.Context) error {

	ws, err := server.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	client := NewBroadcastClient(ws, server)
	server.clients[client] = true
	fmt.Println("Connect websocket")
	go client.Read()

	return c.JSON(http.StatusOK, utils.SuccessJSONData(nil))
}

// func (server *WebsocketServer) HandleMessages() {
// 	for {
// 		msg := <-server.broadcast

// 		for client := range server.clients {
// 			err := client.Write(websocket.TextMessage, msg)
// 			if err != nil {
// 				log.Printf("error: %v", err)
// 				client.Close()
// 				delete(server.clients, client)
// 			}
// 		}
// 	}
// }

func (server *WebsocketServer) RemoveClient(client *Client) {
	delete(server.clients, client)
}

func (server *WebsocketServer) RemoveAllConn() {
	for client := range server.clients {
		client.Close()
		delete(server.clients, client)
	}
}
