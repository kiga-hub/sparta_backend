package ws

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
	"github.com/kiga-hub/sparta_backend/pkg/models"
	"github.com/kiga-hub/sparta_backend/pkg/utils"
)

// ConvertToVisualFormatReadLoop -
func (c *Conn) ConvertToVisualFormatReadLoop() {
	defer c.Close()
	for {
		msgType, data, err := c.ws.ReadMessage()
		if err != nil {
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
			}

			circleName := filepath.Join(GetConfig().DataDir, "in.circle")

			if sparta.IsGridToParaView {
				// Convert to paraview file
				c.Grid2Paraview(c.ctx, filepath.Dir(circleName), GetConfig().ScriptDir)
			}
		}()
	}
}

// ConvertToVisualFormatWriteLoop -
func (c *Conn) ConvertToVisualFormatWriteLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
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

// Grid2Paraview -
func (c *Conn) Grid2Paraview(ctx context.Context, dir, scriptDir string) {

	err := c.Write(1, []byte("----------Start convert grid to paraview!----------"))
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
	}
	fmt.Println("----------Start convert grid to paraview!----------")

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
	fmt.Println("----------Complete the grid to paraview!----------")
}
