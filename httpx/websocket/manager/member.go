package manager

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joeyCheek888/go-std/snowflake"
	"runtime/debug"
	"time"
)

// Member 用户连接
type Member struct {
	ID            int64           // 客户端ID
	Addr          string          // 客户端地址
	Socket        *websocket.Conn // 用户连接
	Send          chan []byte     // 待发送的数据
	Token         string          // 客户端Token
	FirstTime     int64           // 首次连接事件
	HeartbeatTime int64           // 用户上次心跳时间
}

// NewMember 初始化
func NewMember(socket *websocket.Conn, token string) *Member {
	firstTime := time.Now().Unix()

	return &Member{
		Addr:          socket.RemoteAddr().String(),
		Socket:        socket,
		Send:          make(chan []byte, 100),
		ID:            snowflake.GenerateID(),
		Token:         token,
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
	}
}

// 读取客户端数据
func (c *Member) Read(event *Event) {
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

		event.processMessageFunc(c, message)
	}
}

// 向客户端写数据
func (c *Member) Write(event *Event) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("client socket write recover", string(debug.Stack()), r)
		}
	}()
	defer func() {
		event.UnregisterChan <- c
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
func (c *Member) SendMsg(msg []byte) {
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
func (c *Member) close() {
	c.Socket.Close()
}

// Heartbeat 用户心跳
func (c *Member) Heartbeat(currentTime int64) {
	c.HeartbeatTime = currentTime
	return
}

// IsHeartbeatTimeout 心跳超时
func (c *Member) IsHeartbeatTimeout(currentTime int64, heartbeatTimeout int64) (timeout bool) {
	if c.HeartbeatTime+heartbeatTimeout <= currentTime {
		timeout = true
	}
	return
}
