package ws

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
		c.CalculateSpartaResult2(circleName, GetConfig().SpaExec)

		if sparta.IsDumpGrid {
			c.Grid2Paraview(filepath.Dir(circleName), GetConfig().ScriptDir)
			// s.Grid2Paraview(filepath.Dir(circleName))
		}

		// var dataStrs []string
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable x index", msg.Data["variable x index"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable y index", msg.Data["variable y index"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable z index", msg.Data["variable z index"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable n equal", msg.Data["variable n equal"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "variable fnum equal", msg.Data["variable fnum equal"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "seed", msg.Data["seed"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "dimension", msg.Data["dimension"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "global nrho", msg.Data["global nrho"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "global fnum", msg.Data["global fnum"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "timestep", msg.Data["timestep"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "global gridcut", msg.Data["global gridcut"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "global surfmax", msg.Data["global surfmax"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "boundary", msg.Data["boundary"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "create_box", msg.Data["create_box"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "create_grid", msg.Data["create_grid"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "balance_grid", msg.Data["balance_grid"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "species ar.species", msg.Data["species ar.species"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "mixture air Ar frac", msg.Data["mixture air Ar frac"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "mixture air group", msg.Data["mixture air group"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "mixture air Ar vstream", msg.Data["mixture air Ar vstream"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "fix in emit/face air", msg.Data["fix in emit/face air"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "collide vss air", msg.Data["collide vss air"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "read_surf", msg.Data["read_surf"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "surf_collide 1 diffuse", msg.Data["surf_collide 1 diffuse"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "surf_modify", msg.Data["surf_modify"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "create_particles air n", msg.Data["create_particles air n"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "fix", msg.Data["fix"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "collide_modify", "vremax 100 yes"))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "compute g grid all all", msg.Data["compute g grid all all"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "compute max reduce max", msg.Data["compute max reduce max"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "stats_style", msg.Data["stats_style"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "stats", msg.Data["stats"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "run", msg.Data["run"]))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "collide_modify", "vremax 100 no"))
		// dataStrs = append(dataStrs, fmt.Sprintf("%s %s", "run", msg.Data["run"]))

		// dataStr := strings.Join(dataStrs, "\n")

		// Print the result
		// fmt.Println(dataStr)

		// cmd := exec.Command("/home/sparta/build/src/spa_")
		// cmd.Dir = "/home/sparta/examples/test"
		// // Open the file
		// file, err := os.Open("/home/sparta/examples/test/in.circle")
		// if err != nil {
		// 	fmt.Printf(utils.ErrorMsg, err)
		// 	return
		// }
		// defer file.Close()
		// // Redirect the command's stdin to the file
		// cmd.Stdin = file

		// // redirect the command's stdin to the string
		// // cmd.Stdin = strings.NewReader(dataStr)

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

		// // Read the command's output in a separate goroutine to prevent blocking
		// output, err := io.ReadAll(stdout)
		// if err != nil {
		// 	fmt.Printf(utils.ErrorMsg, err)
		// 	return
		// }
		// // Wait for the command to finish
		// if err := cmd.Wait(); err != nil {
		// 	fmt.Printf(utils.ErrorMsg, err)
		// 	return
		// }

		// // Print the output
		// fmt.Printf("The output: %s\n", output)
		// fmt.Printf("%s\n", output)

		// Format the output

		// for i := 0; i < 10; i++ {
		// 	bytes := sparta.ToBytes()

		// 	// // Write the result to the client
		// 	err = c.Write(1, bytes)
		// 	if err != nil {
		// 		fmt.Printf(utils.ErrorMsg, err)
		// 		return
		// 	}
		// 	time.Sleep(1 * time.Second)
		// }
		// ra := time.Now().Format("2006-01-02 15:04:05")

		// convert sparta to []byte
		// bytes := sparta.ToBytes()

		// // // Write the result to the client
		// err = c.Write(1, bytes)
		// if err != nil {
		// 	fmt.Printf(utils.ErrorMsg, err)
		// 	return
		// }

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

func (c *Conn) CalculateSpartaResult(circleName string, spaExe string) string {
	cmd := exec.Command(spaExe)
	cmd.Dir = filepath.Dir(circleName)

	file, err := os.Open(circleName)
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return ""
	}
	defer file.Close()

	cmd.Stdin = file

	reader, writer, _ := os.Pipe()
	multiWriter := io.MultiWriter(os.Stdout, writer)
	cmd.Stdout = multiWriter
	cmd.Stderr = multiWriter

	go func() {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			err = c.Write(1, []byte(line))
			if err != nil {
				fmt.Printf(utils.ErrorMsg, err)
				break
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf(utils.ErrorMsg, err)
		}
	}()

	// if err := cmd.Start(); err != nil {
	// 	fmt.Printf(utils.ErrorMsg, err)
	// 	return ""
	// }

	err = c.Write(1, []byte("begin to calculate"))
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
	}

	// if err := cmd.Wait(); err != nil {
	// 	fmt.Printf(utils.ErrorMsg, err)
	// 	return ""
	// }

	err = cmd.Run()
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
		return ""
	}

	writer.Close()

	return filepath.Dir(circleName)
}

type WebSocketWriter struct {
	ws *websocket.Conn
}

func (wsw WebSocketWriter) Write(p []byte) (n int, err error) {
	err = wsw.ws.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *Conn) CalculateSpartaResult22(circleName string, spaExe string) string {
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

	// 实时处理标准输出
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			fmt.Println("STDOUT:", scanner.Text())
		}
		wg.Done()
	}()

	// 实时处理错误输出
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

func (c *Conn) CalculateSpartaResult2(circleName string, spaExe string) string {
	cmd := exec.Command("bash", "-c", spaExe+" < "+circleName)
	cmd.Dir = filepath.Dir(circleName)
	// using pty to start command
	f, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	defer f.Close()

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

		// buf := make([]byte, 1024)
		// for {
		// 	n, err := f.Read(buf)
		// 	if err != nil {
		// 		if err == io.EOF {
		// 			// end of file
		// 			fmt.Println("Done!")
		// 		}
		// 		// else {
		// 		// occur error
		// 		// fmt.Fprintln(os.Stderr, "读取错误:", err)
		// 		// }
		// 		break // 无论是EOF还是其他错误，都结束循环
		// 	}
		// 	err = c.Write(1, buf[:n])
		// 	if err != nil {
		// 		fmt.Printf(utils.ErrorMsg, err)
		// 		return
		// 	}
		// 	// fmt.Print(string(buf[:n]))
		// }
	}()

	// 等待命令执行完成
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "command wait err:", err)
	}
	return filepath.Dir(circleName)
}

// Grid2Paraview -
func (c *Conn) Grid2Paraview(dir, scriptDir string) {
	go func() {
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

		// Wait for the command to finish
		if err := cmd.Wait(); err != nil {
			fmt.Printf(utils.ErrorMsg, err)
			return
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf(utils.ErrorMsg, err)
			return
		}
	}()
}
