package manager

import (
	"sync"
	"time"
)

const (
	heartbeatTimeout = 3 * 60 // 用户心跳超时时间
)

var (
	Manager *manager
)

// Event
// @Description: manager
type Event struct {
	Register   chan *Member // 连接连接处理
	Unregister chan *Member // 断开连接处理程序
	Broadcast  chan []byte  // 广播 向全部成员发送数据
}

// 连接管理
type manager struct {
	Members     map[*Member]bool  // 全部的连接
	MembersLock sync.RWMutex      // 读写锁
	Users       map[int64]*Member // 登录的用户
	UserLock    sync.RWMutex      // 读写锁
	Event
}

// init 引入该包后自动创建连接管理
func init() {
	Manager = &manager{
		Members:     make(map[*Member]bool),
		MembersLock: sync.RWMutex{},
		Users:       make(map[int64]*Member),
		UserLock:    sync.RWMutex{},
		Event: Event{
			Register:   make(chan *Member, 1000),
			Unregister: make(chan *Member, 1000),
			Broadcast:  make(chan []byte, 1000),
		},
	}
	return
}

// InClient 判断连接是否存在
func (manager *manager) InClient(client *Member) (ok bool) {
	manager.MembersLock.RLock()
	defer manager.MembersLock.RUnlock()

	// 连接存在，在添加
	_, ok = manager.Members[client]
	return
}

// GetClients 获取所有客户端
func (manager *manager) GetClients() (clients map[*Member]bool) {
	clients = make(map[*Member]bool)
	manager.ClientsRange(func(client *Member, value bool) (result bool) {
		clients[client] = value
		return true
	})
	return
}

// ClientsRange 遍历
func (manager *manager) ClientsRange(f func(client *Member, value bool) (result bool)) {
	manager.MembersLock.RLock()
	defer manager.MembersLock.RUnlock()
	for key, value := range manager.Members {
		result := f(key, value)
		if result == false {
			return
		}
	}
	return
}

// GetClientsLen GetClientsLen
func (manager *manager) GetClientsLen() (clientsLen int) {
	clientsLen = len(manager.Members)
	return
}

// AddClients 添加客户端
func (manager *manager) AddClients(client *Member) {
	manager.MembersLock.Lock()
	defer manager.MembersLock.Unlock()
	manager.Members[client] = true
	manager.Users[client.ID] = client
}

// DelClients 删除客户端
func (manager *manager) DelClients(client *Member) {
	manager.MembersLock.Lock()
	defer manager.MembersLock.Unlock()
	if _, ok := manager.Members[client]; ok {
		delete(manager.Members, client)
	}

	if _, ok := manager.Users[client.ID]; ok {
		delete(manager.Users, client.ID)
	}
}

// GetUserClient 获取用户的连接
func (manager *manager) GetUserClient(ID int64) (client *Member) {
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
func (manager *manager) GetUserClients() (clients []*Member) {
	clients = make([]*Member, 0)
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
		Manager.ClientsRange(func(client *Member, value bool) (result bool) {
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
func GetUserClient(ID int64) (client *Member) {
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
