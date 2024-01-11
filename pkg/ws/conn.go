package ws

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/kiga-hub/arc/logging"
	"github.com/gorilla/websocket"
)

const (
	// wrtie wait time
	writeWait = 60 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 20 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	// maxMessageSize = 512
)

// Message - client read message
type Message struct {
	// websocket.TextMessage 
	Type int
	Data []byte
}

// Conn - websocket connect struct
type Conn struct {
	ws        *websocket.Conn // websocket
	outChan   chan *Message   // write queue
	logger    logging.ILogger
	closeChan chan struct{}
	isClosed  bool
	closeLock sync.RWMutex
}

// NewConn - new one
func NewConn(conn *websocket.Conn, logger logging.ILogger) *Conn {
	srv := &Conn{
		ws:        conn,
		outChan:   make(chan *Message, 1024*4),
		logger:    logger,
		closeChan: make(chan struct{}),
		closeLock: sync.RWMutex{},
	}
	return srv
}

// ReadLoop -
func (c *Conn) ReadLoop() {
	defer c.Close()

	for {
		msgType, data, err := c.ws.ReadMessage()
		if err != nil {
			websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure)
			return
		}
		if msgType == websocket.CloseMessage {
			c.logger.Infof("websocket param error: %d:%v", msgType, data)
			return
		}
	}
}

// WriteLoop - send message to client
func (c *Conn) WriteLoop() {
	defer c.Close()

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-c.closeChan:
			// get close notice
			return
		case msg := <-c.outChan:
			// get a response
			if err := c.ws.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
				c.logger.Error(err)
				return
			}
			// write to websocket
			if err := c.ws.WriteMessage(msg.Type, msg.Data); err != nil {
				c.logger.Errorw(fmt.Sprintf("send to client err:%s", err.Error()))
				// shut donw
				return
			}
		case <-ticker.C:
			// timeout
			if err := c.ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.logger.Errorw(fmt.Sprintf("ping:%s", err.Error()))
				return
			}
		}
	}
}

// Write - write message to queue
func (c *Conn) Write(Type int, data []byte) error {
	select {
	case c.outChan <- &Message{Type, data}:
	case <-c.closeChan:
		return errors.New("ws write - already closed")
	}
	return nil
}

// Close -
func (c *Conn) Close() {
	c.closeLock.Lock()
	defer c.closeLock.Unlock()

	if !c.isClosed {
		c.isClosed = true
		c.ws.Close()
		close(c.closeChan)
	}
}

// IsClosed -
func (c *Conn) IsClosed() bool {
	c.closeLock.RLock()
	defer c.closeLock.RUnlock()

	return c.isClosed
}
