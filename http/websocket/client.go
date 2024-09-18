package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go-std/snowflake"
	"runtime/debug"
)

// Client 用户连接
type Client struct {
	Addr          string          // 客户端地址
	Socket        *websocket.Conn // 用户连接
	Send          chan []byte     // 待发送的数据
	ClientID      int64
	Token         string
	FirstTime     int64 // 首次连接事件
	HeartbeatTime int64 // 用户上次心跳时间
}

// NewClient 初始化
func NewClient(addr string, token string, socket *websocket.Conn, firstTime int64) (client *Client) {
	client = &Client{
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan []byte, 100),
		ClientID:      snowflake.GenerateID(),
		Token:         token,
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
	}
	return
}

// 读取客户端数据
func (c *Client) Read() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("client socket read recover ", string(debug.Stack()), r)
		}
	}()
	defer func() {
		close(c.Send)
	}()
	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			return
		}

		// 处理程序
		// fmt.Println("读取客户端数据 处理:", string(message))
		defaultProcessData(c, message)
	}
}

// 向客户端写数据
func (c *Client) Write() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("client socket write recover", string(debug.Stack()), r)
		}
	}()
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 发送数据错误 关闭连接
				return
			}
			_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// SendMsg 发送数据
func (c *Client) SendMsg(msg []byte) {
	if c == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			// fmt.Println("SendMsg stop:", r, string(debug.Stack()))
		}
	}()
	c.Send <- msg
}

// close 关闭客户端连接
func (c *Client) close() {
	c.Socket.Close()
}

// Heartbeat 用户心跳
func (c *Client) Heartbeat(currentTime int64) {
	c.HeartbeatTime = currentTime

	return
}

// IsHeartbeatTimeout 心跳超时
func (c *Client) IsHeartbeatTimeout(currentTime int64, heartbeatTimeout int64) (timeout bool) {
	if c.HeartbeatTime+heartbeatTimeout <= currentTime {
		timeout = true
	}
	return
}
