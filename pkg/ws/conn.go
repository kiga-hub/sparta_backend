package ws

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
	"github.com/kiga-hub/arc/logging"
	"github.com/kiga-hub/sparta_backend/pkg/models"
	"github.com/kiga-hub/sparta_backend/pkg/utils"
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

// SendMessage - client read message
type SendMessage struct {
	// websocket.TextMessage
	Type int
	Data []byte
}

// Conn - websocket connect struct
type Conn struct {
	ws        *websocket.Conn   // websocket
	outChan   chan *SendMessage // write queue
	logger    logging.ILogger
	closeChan chan struct{}
	isClosed  bool
	closeLock sync.RWMutex
}

// NewConn - new one
func NewConn(conn *websocket.Conn, logger logging.ILogger) *Conn {
	srv := &Conn{
		ws:        conn,
		outChan:   make(chan *SendMessage, 1024*4),
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
		//read
		// var msg SocketMessage
		// err := c.ws.ReadJSON(&msg)
		// if err != nil {
		// 	websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure)
		// 	return
		// }

		msgType, data, err := c.ws.ReadMessage()
		if err != nil {
			websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure)
			return
		}
		if msgType == websocket.CloseMessage {
			c.logger.Infof("websocket param error: %d:%v", msgType, data)
			return
		}

		sparta := &models.Sparta{}
		if err = json.Unmarshal(data, sparta); err != nil {
			c.logger.Infof("json.Unmarshal")

			err := c.Write(1, []byte("Error in passing parameters"))
			if err != nil {
				fmt.Printf(utils.ErrorMsg, err)
				return
			}
			return
		}

		surfName := strings.Replace(sparta.UploadStlName, "stl", "surf", -1)
		circleName := sparta.ProcessSparta(GetConfig().DataDir, surfName)
		c.CalculateSpartaResult(circleName, GetConfig().SpaExec)

		if sparta.IsDumpGrid {
			err := c.Write(1, []byte("Start convert grid to paraview!"))
			if err != nil {
				fmt.Printf(utils.ErrorMsg, err)
			}
			c.Grid2Paraview(filepath.Dir(circleName), GetConfig().ScriptDir)
			err = c.Write(1, []byte("Done!"))
			if err != nil {
				fmt.Printf(utils.ErrorMsg, err)
			}
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
			fmt.Println("close websocket remote address is:", c.ws.RemoteAddr().String())
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
	case c.outChan <- &SendMessage{Type, data}:
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

func (c *Conn) ReadFromStdin() {
	for {
		var data []byte
		_, err := os.Stdin.Read(data)
		if err != nil {
			fmt.Printf(utils.ErrorMsg, err)
			return
		}

		err = c.Write(1, data)
		if err != nil {
			fmt.Printf(utils.ErrorMsg, err)
			return
		}
	}
}

func (c *Conn) CalculateDefault(circleName string, spaExe string) string {
	cmd := exec.Command("bash", "-c", spaExe+" < "+circleName)
	cmd.Dir = filepath.Dir(circleName)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf(utils.ErrorMsg, err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			fmt.Println("STDOUT:", scanner.Text())
		}
		wg.Done()
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			fmt.Println("STDERR:", scanner.Text())
		}
		wg.Done()
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command:", err)
		return ""
	}

	return filepath.Dir(circleName)
}

func (c *Conn) CalculateSpartaResult(circleName string, spaExe string) string {
	cmd := exec.Command("bash", "-c", spaExe+" < "+circleName)
	cmd.Dir = filepath.Dir(circleName)
	// using pty to start command
	f, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = c.Write(1, []byte("Begin!"))
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
	}

	go func() {

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			// Send each line to the websocket
			err = c.Write(1, []byte(scanner.Text()))
			if err != nil {
				fmt.Printf(utils.ErrorMsg, err)
				return
			}
		}

	}()

	// wait for command to finish
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "command wait err:", err)
	}
	err = c.Write(1, []byte("Done!"))
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
	}
	return filepath.Dir(circleName)
}

// Grid2Paraview -
func (c *Conn) Grid2Paraview(dir, scriptDir string) {
	// do grid2paraview. pvpython grid2paraview.py circle.txt output -r tmp.grid.1000
	txtFile := filepath.Join(dir, "in.txt")
	outputDir := dir + "/output/"
	tmpGridDir := filepath.Join(dir, "tmp.grid.*")

	// Delete the outputDir directory, TODO need to keep historical files
	if err := utils.ClearDir(outputDir); err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return
	}

	cmd := exec.Command("pvpython", "grid2paraview.py", txtFile, outputDir, "-r", tmpGridDir)
	cmd.Dir = filepath.Join(scriptDir, "paraview")

	// Create a pipe to capture the command's output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return
	}

	// Start executing the command
	if err := cmd.Start(); err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return
	}

	// Create a new Scanner that will read from stdout
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		// Send each line to the websocket
		err = c.Write(1, []byte(scanner.Text()))
		if err != nil {
			fmt.Printf(utils.ErrorMsg, err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return
	}
	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return
	}

}
