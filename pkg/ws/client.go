package ws

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/gorilla/websocket"
)

// Client -
type Client struct {
	socket *websocket.Conn
	server *WebsocketServer
}

// NewBroadcastClient - 创建websocket客户端
func NewBroadcastClient(conn *websocket.Conn, server *WebsocketServer) *Client {
	return &Client{
		socket: conn,
		server: server,
	}
}

// Close -
func (c *Client) Close() error {
	return c.socket.Close()
}

// Read - 监控断线事件
func (c *Client) Read() {
	defer func() {
		//断开连接
		c.server.RemoveClient(c)
		c.Close()
	}()
	for {
		var msg SocketMessage
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("error: %v", err)
			delete(c.server.clients, c)
			break
		}

		fmt.Println("Read websocket: ", msg)
		c.Write(1, []byte(msg.Message+" it is response"))
	}
}

// Write -
func (c *Client) Write(msgType int, data []byte) error {
	if err := c.socket.SetWriteDeadline(time.Now().Add(time.Second * 10)); err != nil {
		return errors.WithStack(err)
	}
	if err := c.socket.WriteMessage(msgType, data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
