package ws

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
	ctx       context.Context
	cancel    context.CancelFunc
	firstRun  bool
	count     int
}

// NewConn - new one
func NewConn(conn *websocket.Conn, logger logging.ILogger) *Conn {
	ctx, cancel := context.WithCancel(context.Background())
	srv := &Conn{
		ws:        conn,
		outChan:   make(chan *SendMessage, 1024*4),
		logger:    logger,
		closeChan: make(chan struct{}),
		closeLock: sync.RWMutex{},
		ctx:       ctx,
		cancel:    cancel,
		firstRun:  true,
		count:     0,
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
		c.count++
		go func() {
			fmt.Println("This is the "+strconv.Itoa(c.count)+" round of computation", "Time: ", time.Now().Format("2006-01-02 15:04:05"))
			outputInfo := fmt.Sprintf("This is the %d round of computation Time: %s", c.count, time.Now().Format("2006-01-02 15:04:05"))
			if err := c.Write(1, []byte(outputInfo)); err != nil {
				fmt.Printf(utils.ErrorMsg, err)
				return
			}

			if !c.firstRun {
				c.cancel()
				c.ctx, c.cancel = context.WithCancel(context.Background())
			} else {
				c.firstRun = false
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
			if _, err := c.CalculateSpartaResult(c.ctx, circleName, GetConfig().SpaExec); err != nil {
				err := c.Write(1, []byte("----------Forced interrupt!----------"))
				if err != nil {
					fmt.Printf(utils.ErrorMsg, err)
					return
				}
				return
			}

			if sparta.IsGridToParaView {
				c.Grid2Paraview(c.ctx, filepath.Dir(circleName), GetConfig().ScriptDir)
			}
		}()
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
		c.cancel()
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

func (c *Conn) CalculateSpartaResult(ctx context.Context, circleName string, spaExe string) (string, error) {
	cmd := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("%s < %s", spaExe, circleName))
	cmd.Dir = filepath.Dir(circleName)

	// Delete all files ending in *.ppm under the /data directory
	err := filepath.Walk(GetConfig().DataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == ".ppm" {
			err := os.Remove(path)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return "", err
	}

	err = c.Write(1, []byte("----------Start sparta calculate!----------"))
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
	}

	f, err := pty.Start(cmd)
	if err != nil {
		return "", err
	}
	defer f.Close()

	done := make(chan error, 1)
	go func() {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				done <- ctx.Err()
				return
			default:
				err = c.Write(1, []byte(scanner.Text()))
				if err != nil {
					fmt.Printf(utils.ErrorMsg, err)
					done <- err
					return
				}
			}
		}
		done <- nil //scanner.Err()
	}()

	select {
	case <-ctx.Done():
		// Context was cancelled, kill the process
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		<-done // Wait for goroutine to finish
		return "", ctx.Err()
	case err := <-done:
		// Command finished
		if err != nil {
			fmt.Fprintln(os.Stderr, "command finished with error:", err)
		}
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "command wait err:", err)
		return "", err
	}

	err = c.Write(1, []byte("----------Complete sparta calculate!----------"))
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
	}
	return filepath.Dir(circleName), nil
}

// Grid2Paraview -
func (c *Conn) Grid2Paraview(ctx context.Context, dir, scriptDir string) {

	err := c.Write(1, []byte("----------Start convert grid to paraview!----------"))
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
	}

	// do grid2paraview. pvpython grid2paraview.py circle.txt output -r tmp.grid.1000
	txtFile := filepath.Join(dir, "in.txt")
	outputDir := dir + "/output/"
	tmpGridDir := filepath.Join(dir, "tmp.grid.*")

	// Delete the outputDir directory, TODO If need to keep historical files
	if err := utils.ClearDir(outputDir); err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return
	}

	cmd := exec.CommandContext(c.ctx, "pvpython", "grid2paraview.py", txtFile, outputDir, "-r", tmpGridDir)
	cmd.Dir = filepath.Join(scriptDir, "paraview")

	// // Create a pipe to capture the command's output
	// stdout, err := cmd.StdoutPipe()
	// if err != nil {
	// 	fmt.Printf(utils.ErrorMsg, err)
	// 	return
	// }

	// // Start executing the command
	// if err := cmd.Start(); err != nil {
	// 	fmt.Printf(utils.ErrorMsg, err)
	// 	return
	// }

	f, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	done := make(chan error, 1)
	go func() {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				// Context was cancelled, stop reading
				err = c.Write(1, []byte("Forced interrupt!"))
				if err != nil {
					fmt.Printf(utils.ErrorMsg, err)
				}

				done <- ctx.Err()
				return
			default:
				err = c.Write(1, []byte(scanner.Text()))
				if err != nil {
					fmt.Printf(utils.ErrorMsg, err)
					done <- err
					return
				}
			}
		}
		done <- nil //scanner.Err()
	}()

	select {
	case <-ctx.Done():
		// Context was cancelled, kill the process
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		<-done // Wait for goroutine to finish
		return
	case err := <-done:
		// Command finished
		if err != nil {
			fmt.Fprintln(os.Stderr, "command finished with error:", err)
		}
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return
	}

	err = c.Write(1, []byte("----------Complete the grid to paraview!----------"))
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
	}

}
