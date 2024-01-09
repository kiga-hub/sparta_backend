package ws

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"

	"github.com/gorilla/websocket"
)

const errorMsg = "error: %v"

// Client -
type Client struct {
	socket *websocket.Conn
	server *WebsocketServer
}

// NewBroadcastClient - Create a new client
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

// Read Read message from client
func (c *Client) Read() {
	defer c.disconnect()
	for {
		if !c.processMessage() {
			break
		}
	}
}

func (c *Client) disconnect() {
	//断开连接
	c.server.RemoveClient(c)
	err := c.Close()
	if err != nil {
		return
	}
}

// processMessage -
func (c *Client) processMessage() bool {
	var msg SocketMessage
	err := c.socket.ReadJSON(&msg)
	if err != nil {
		fmt.Printf(errorMsg, err)
		delete(c.server.clients, c)
		return false
	}

	cmd := exec.Command("/home/spa")
	cmd.Dir = "/home/sparta-13Apr2023/bench"
	// Open the file
	file, err := os.Open("/home/sparta-13Apr2023/bench/in.sphere")
	if err != nil {
		fmt.Printf(errorMsg, err)
		return false
	}

	// Redirect the command's stdin to the file
	cmd.Stdin = file

	// Create a pipe to capture the command's output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf(errorMsg, err)
		return false
	}

	// Start executing the command
	if err := cmd.Start(); err != nil {
		fmt.Printf(errorMsg, err)
		return false
	}

	// Read the command's output in a separate goroutine to prevent blocking
	go func() {
		output, err := ioutil.ReadAll(stdout)
		if err != nil {
			fmt.Printf(errorMsg, err)
			return
		}

		// Close the file after reading
		err = file.Close()
		if err != nil {
			fmt.Printf(errorMsg, err)
			return
		}

		// Print the output
		fmt.Printf("The output: %s\n", output)
		fmt.Printf("%s\n", output)

		// Format the output
		result := fmt.Sprintf("%s", output)

		// Write the result to the client
		err = c.Write(1, []byte(result))
		if err != nil {
			fmt.Printf(errorMsg, err)
			return
		}
	}()

	return true
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
