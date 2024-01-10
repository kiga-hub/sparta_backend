package ws

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
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
	// Convert the Data field to a string with newline-separated key-value pairs
	var dataStrs []string
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable x index", msg.Data["variable x index"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable y index", msg.Data["variable y index"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable z index", msg.Data["variable z index"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable n equal", msg.Data["variable n equal"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable fnum equal", msg.Data["variable fnum equal"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "seed", msg.Data["seed"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "dimension", msg.Data["dimension"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "global nrho", msg.Data["global nrho"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "global fnum", msg.Data["global fnum"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "timestep", msg.Data["timestep"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "global gridcut", msg.Data["global gridcut"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "global surfmax", msg.Data["global surfmax"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "boundary", msg.Data["boundary"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "create_box", msg.Data["create_box"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "create_grid", msg.Data["create_grid"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "balance_grid", msg.Data["balance_grid"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "species ar.species", msg.Data["species ar.species"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "mixture air Ar frac", msg.Data["mixture air Ar frac"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "mixture air group", msg.Data["mixture air group"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "mixture air Ar vstream", msg.Data["mixture air Ar vstream"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "fix in emit/face air", msg.Data["fix in emit/face air"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "collide vss air", msg.Data["collide vss air"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "read_surf", msg.Data["read_surf"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "surf_collide 1 diffuse", msg.Data["surf_collide 1 diffuse"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "surf_modify", msg.Data["surf_modify"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "create_particles air n", msg.Data["create_particles air n"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "fix", msg.Data["fix"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "collide_modify", "vremax 100 yes"))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "compute g grid all all", msg.Data["compute g grid all all"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "compute max reduce max", msg.Data["compute max reduce max"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "stats_style", msg.Data["stats_style"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "stats", msg.Data["stats"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "run", msg.Data["run"]))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "collide_modify", "vremax 100 no"))
	dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "run", msg.Data["run"]))

	dataStr := strings.Join(dataStrs, "\n")

	// Print the result
	fmt.Println(dataStr)

	cmd := exec.Command("/home/spa_")
	cmd.Dir = "/home/sparta-13Apr2023/bench"
	// Open the file
	// file, err := os.Open("/home/sparta-13Apr2023/bench/in.sphere")
	// if err != nil {
	// 	fmt.Printf(errorMsg, err)
	// 	return false
	// }

	// Redirect the command's stdin to the file
	// cmd.Stdin = file

	// redirect the command's stdin to the string
	cmd.Stdin = strings.NewReader(dataStr)

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
		output, err := io.ReadAll(stdout)
		if err != nil {
			fmt.Printf(errorMsg, err)
			return
		}

		// // Close the file after reading
		// err = file.Close()
		// if err != nil {
		// 	fmt.Printf(errorMsg, err)
		// 	return
		// }

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
