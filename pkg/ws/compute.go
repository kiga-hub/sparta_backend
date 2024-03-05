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
	"strings"
	"sync"
	"time"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
	"github.com/kiga-hub/sparta_backend/pkg/models"
	"github.com/kiga-hub/sparta_backend/pkg/utils"
)

// ComputeReadLoop -
func (c *Conn) ComputeReadLoop() {
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
			circleName, err := sparta.ProcessSparta(GetConfig().DataDir, surfName)
			if err != nil {
				err := c.Write(1, []byte("Error in processing sparta"))
				if err != nil {
					fmt.Printf(utils.ErrorMsg, err)
					return
				}
				return
			}
			if _, err := c.ComputeSpartaResult(c.ctx, circleName, GetConfig().SpaExec); err != nil {
				err := c.Write(1, []byte("----------Forced interrupt!----------"))
				if err != nil {
					fmt.Printf(utils.ErrorMsg, err)
					return
				}
				return
			}

			// if sparta.IsGridToParaView {
			// 	c.Grid2Paraview(c.ctx, filepath.Dir(circleName), GetConfig().ScriptDir)
			// }
		}()
	}
}

// ComputeWriteLoop - send message to client
func (c *Conn) ComputeWriteLoop() {
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

// ComputeDefault -
func (c *Conn) ComputeDefault(circleName string, spaExe string) string {
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

// ComputeSpartaResult -
func (c *Conn) ComputeSpartaResult(ctx context.Context, circleName string, spaExe string) (string, error) {
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

	err = c.Write(1, []byte("----------Start sparta compute!----------"))
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
	}
	fmt.Println("----------Start sparta compute!----------")

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

	err = c.Write(1, []byte("----------Complete sparta compute!----------"))
	if err != nil {
		fmt.Printf(utils.ErrorMsg, err)
	}
	fmt.Println("----------Complete sparta compute!----------")
	return filepath.Dir(circleName), nil
}
