package ws

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

		fmt.Println("Read websocket: ", msg.Message)

		// cmd := exec.Command("/home/spa", " < ", " in.sphere")
		// cmd.Dir = "/home/sparta-13Apr2023/bench"

		// // 创建一个Buffer，并将msg.Message的内容写入这个Buffer
		// var stdin bytes.Buffer
		// stdin.Write([]byte(msg.Message))
		// cmd.Stdin = &stdin

		// // 获取输出对象，可以从该对象中读取输出结果
		// output, err := cmd.Output()
		// if err != nil {
		// 	fmt.Println("cmd err: ", err)
		// }

		cmd := exec.Command("/home/spa")
		cmd.Dir = "/home/sparta-13Apr2023/bench"
		// 打开文件
		file, err := os.Open("/home/sparta-13Apr2023/bench/in.sphere")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// 将命令的标准输入重定向到文件
		cmd.Stdin = file

		// 创建一个管道来捕获命令的输出
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}

		// 开始执行命令
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		// 读取命令的输出
		output, err := ioutil.ReadAll(stdout)
		if err != nil {
			log.Fatal(err)
		}

		// 等待命令执行完成
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("The output: %s\n", output)

		// 打印输出结果
		fmt.Printf("%s\n", output)

		// output 格式化输出
		result := fmt.Sprintf("%s", output)
		c.Write(1, []byte(result))
		c.Write(1, []byte(result))
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
