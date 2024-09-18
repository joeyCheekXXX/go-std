package websocket

// EventRegister 用户建立连接事件
func (manager *manager) EventRegister(client *Client) {
	manager.AddClients(client)
	registerCallback(client)
}

// EventUnregister 用户断开连接
func (manager *manager) EventUnregister(client *Client) {
	manager.DelClients(client)
	closeConnCallback(client)
	// 关闭 chan
	client.close()
}
