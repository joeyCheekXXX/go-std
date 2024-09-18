package websocket

import (
	"fmt"
	"sync"
	"time"
)

const (
	heartbeatTimeout = 3 * 60 // 用户心跳超时时间
)

var (
	Manager *manager
)

type Event struct {
	Register   chan *Client // 连接连接处理
	Unregister chan *Client // 断开连接处理程序
	Broadcast  chan []byte  // 广播 向全部成员发送数据
}

// 连接管理
type manager struct {
	Clients     map[*Client]bool  // 全部的连接
	ClientsLock sync.RWMutex      // 读写锁
	Users       map[int64]*Client // 登录的用户 // appID+uuid
	UserLock    sync.RWMutex      // 读写锁
	Event
}

// init 引入该包后自动创建连接管理
func init() {
	Manager = &manager{
		Clients:     make(map[*Client]bool),
		ClientsLock: sync.RWMutex{},
		Users:       make(map[int64]*Client),
		UserLock:    sync.RWMutex{},
		Event: Event{
			Register:   make(chan *Client, 1000),
			Unregister: make(chan *Client, 1000),
			Broadcast:  make(chan []byte, 1000),
		},
	}
	return
}

// InClient 判断连接是否存在
func (manager *manager) InClient(client *Client) (ok bool) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	// 连接存在，在添加
	_, ok = manager.Clients[client]
	return
}

// GetClients 获取所有客户端
func (manager *manager) GetClients() (clients map[*Client]bool) {
	clients = make(map[*Client]bool)
	manager.ClientsRange(func(client *Client, value bool) (result bool) {
		clients[client] = value
		return true
	})
	return
}

// ClientsRange 遍历
func (manager *manager) ClientsRange(f func(client *Client, value bool) (result bool)) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()
	for key, value := range manager.Clients {
		result := f(key, value)
		if result == false {
			return
		}
	}
	return
}

// GetClientsLen GetClientsLen
func (manager *manager) GetClientsLen() (clientsLen int) {
	clientsLen = len(manager.Clients)
	return
}

// AddClients 添加客户端
func (manager *manager) AddClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()
	manager.Clients[client] = true
	manager.Users[client.ClientID] = client
}

// DelClients 删除客户端
func (manager *manager) DelClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()
	if _, ok := manager.Clients[client]; ok {
		delete(manager.Clients, client)
	}

	if _, ok := manager.Users[client.ClientID]; ok {
		delete(manager.Users, client.ClientID)
	}
}

// GetUserClient 获取用户的连接
func (manager *manager) GetUserClient(ID int64) (client *Client) {
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	if value, ok := manager.Users[ID]; ok {
		client = value
	}
	return
}

// GetUsersLen GetClientsLen
func (manager *manager) GetUsersLen() (userLen int) {
	userLen = len(manager.Users)
	return
}

// GetUserClients 获取用户的key
func (manager *manager) GetUserClients() (clients []*Client) {
	clients = make([]*Client, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	for _, v := range manager.Users {
		clients = append(clients, v)
	}
	return
}

// sendAll 向全部成员发送数据
func (manager *manager) sendAll(message []byte) {
	clients := manager.GetUserClients()
	for _, conn := range clients {
		conn.SendMsg(message)
	}
}

// Start 事件管道处理程序
func (manager *manager) Start() {
	for {
		select {
		case conn := <-manager.Register:
			// 建立连接事件
			manager.EventRegister(conn)
		case conn := <-manager.Unregister:
			// 断开连接事件
			manager.EventUnregister(conn)
		case message := <-manager.Broadcast:
			// 广播事件
			clients := manager.GetClients()
			for conn := range clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
				}
			}
		}
	}
}

// GetManagerInfo 获取管理者信息
func GetManagerInfo(isDebug string) (managerInfo map[string]interface{}) {
	managerInfo = make(map[string]interface{})
	managerInfo["clientsLen"] = Manager.GetClientsLen()        // 客户端连接数
	managerInfo["usersLen"] = Manager.GetUsersLen()            // 登录用户数
	managerInfo["chanRegisterLen"] = len(Manager.Register)     // 未处理连接事件数
	managerInfo["chanUnregisterLen"] = len(Manager.Unregister) // 未处理退出登录事件数
	managerInfo["chanBroadcastLen"] = len(Manager.Broadcast)   // 未处理广播事件数
	if isDebug == "true" {
		addrList := make([]string, 0)
		Manager.ClientsRange(func(client *Client, value bool) (result bool) {
			addrList = append(addrList, client.Addr)
			return true
		})
		// users := Manager.GetUserKeys()
		managerInfo["clients"] = addrList // 客户端列表
		// managerInfo["users"] = users      // 登录用户列表
	}
	return
}

// GetUserClient 获取用户所在的连接
func GetUserClient(ID int64) (client *Client) {
	client = Manager.GetUserClient(ID)
	return
}

// ClearTimeoutConnections 定时清理超时连接
func ClearTimeoutConnections() {
	if Manager == nil {
		return
	}
	currentTime := time.Now().Unix()

	clients := Manager.GetClients()
	for client := range clients {
		if client.IsHeartbeatTimeout(currentTime, heartbeatTimeout) {
			_ = client.close
		}
	}
}

// AllSendMessages 全员广播
func AllSendMessages(Symbol string, ID int64, data string) {
	fmt.Println("全员广播", Symbol, ID, data)
	Manager.sendAll([]byte(data))
}
