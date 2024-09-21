package manager

import (
	"github.com/joeyCheek888/go-std/httpx/websocket/handler"
)

// EventRegister 用户建立连接事件
func (manager *manager) EventRegister(client *Client) {
	manager.AddClients(client)
	handler.RegisterCallback(client)
}

// EventUnregister 用户断开连接
func (manager *manager) EventUnregister(client *Client) {
	manager.DelClients(client)
	handler.CloseConnCallback(client)
	// 关闭 chan
	client.close()
}
