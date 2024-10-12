package manager

import (
	"net/http"
	"sync"
	"time"
)

const (
	heartbeatTimeout = 3 * 60 // 用户心跳超时时间
)

// Manager 连接管理
type Manager struct {
	Lock        sync.RWMutex      // 读写锁
	Members     map[*Member]bool  // 全部的连接
	MemberIdRef map[int64]*Member // 连接的ID索引
	Event       *Event            // 事件管理
}

func NewManager() *Manager {
	return &Manager{
		Members:     make(map[*Member]bool),
		MemberIdRef: make(map[int64]*Member),
		Lock:        sync.RWMutex{},
		Event:       NewEvent(),
	}
}

// WithRegisterConnFunc
//
//	@Description: 设置注册连接回调函数
//	@receiver manager
//	@param registerConnFunc
func (manager *Manager) WithRegisterConnFunc(registerConnFunc func(member *Member)) {
	manager.Event.registerConnFunc = registerConnFunc
}

// WithCloseConnFunc
//
//	@Description: 设置关闭连接回调函数
//	@receiver manager
//	@param registerConnFunc
func (manager *Manager) WithCloseConnFunc(closeConnFunc func(member *Member)) {
	manager.Event.closeConnFunc = closeConnFunc
}

// WithProcessMessageFunc
//
//	@Description: 设置消息处理回调函数
//	@receiver manager
//	@param registerConnFunc
func (manager *Manager) WithProcessMessageFunc(processMessageFunc func(member *Member, message []byte)) {
	manager.Event.processMessageFunc = processMessageFunc
}

// WithCheckTokenFunc
//
//	@Description: 设置检查token回调函数
//	@receiver manager
//	@param checkTokenFunc
func (manager *Manager) WithCheckTokenFunc(checkTokenFunc func(token string) bool) {
	manager.Event.CheckTokenFunc = checkTokenFunc
}

// WithCheckOriginFunc
//
//	@Description: 设置检查origin回调函数 用于检查请求的合法性，是否跨域伪造等
//	@receiver manager
//	@param checkOriginFunc
func (manager *Manager) WithCheckOriginFunc(checkOriginFunc func(r *http.Request) bool) {
	manager.Event.CheckOriginFunc = checkOriginFunc
}

// GetAllMember
//
//	@Description: 获取所有客户端
//	@receiver manager
//	@return clients
func (manager *Manager) GetAllMember() (clients map[*Member]bool) {
	clients = make(map[*Member]bool)
	manager.MembersRange(func(client *Member, value bool) (result bool) {
		clients[client] = value
		return true
	})
	return
}

// MembersRange
//
//	@Description: 遍历客户端
//	@receiver manager
//	@param f
func (manager *Manager) MembersRange(f func(client *Member, value bool) (result bool)) {
	manager.Lock.RLock()
	defer manager.Lock.RUnlock()
	for key, value := range manager.Members {
		result := f(key, value)
		if result == false {
			return
		}
	}
	return
}

// GetAllMembersLen
//
//	@Description: 获取客户端数量
//	@receiver manager
//	@return clientsLen
func (manager *Manager) GetAllMembersLen() (clientsLen int) {
	clientsLen = len(manager.Members)
	return
}

// AddMember
//
//	@Description: 添加客户端
//	@receiver manager
//	@param member
func (manager *Manager) AddMember(member *Member) {
	manager.Lock.Lock()
	defer manager.Lock.Unlock()

	manager.Members[member] = true
	manager.MemberIdRef[member.ID] = member
}

// DelMember
//
//	@Description:  删除客户端
//	@receiver manager
//	@param member
func (manager *Manager) DelMember(member *Member) {
	manager.Lock.Lock()
	defer manager.Lock.Unlock()

	if _, ok := manager.Members[member]; ok {
		delete(manager.Members, member)
	}

	if _, ok := manager.MemberIdRef[member.ID]; ok {
		delete(manager.MemberIdRef, member.ID)
	}
}

// GetMemberByID
//
//	@Description: 获取连接 根据ID
//	@receiver manager
//	@param ID
//	@return client
func (manager *Manager) GetMemberByID(ID int64) (client *Member) {
	manager.Lock.RLock()
	defer manager.Lock.RUnlock()
	if value, ok := manager.MemberIdRef[ID]; ok {
		client = value
	}
	return
}

// GetMembersLen
//
//	@Description: 获取所有连接数量
//	@receiver manager
//	@return userLen
func (manager *Manager) GetMembersLen() (userLen int) {
	userLen = len(manager.MemberIdRef)
	return
}

// sendAll
//
//	@Description: 向全部成员发送数据
//	@receiver manager
//	@param message
func (manager *Manager) sendAll(message []byte) {
	for conn, _ := range manager.GetAllMember() {
		conn.SendMsg(message)
	}
}

// Start
//
//	@Description: 事件管道处理程序
//	@receiver manager
func (manager *Manager) Start() {
	for {
		select {
		case member := <-manager.Event.RegisterChan:
			// 建立连接事件
			manager.Event.registerConnFunc(member)
			manager.AddMember(member)

			go func() {
				member.Read(manager.Event)
				member.Write(manager.Event)
			}()

		case member := <-manager.Event.UnregisterChan:
			// 断开连接事件
			manager.DelMember(member)
			manager.Event.closeConnFunc(member)

			// 关闭 member chan
			member.close()

		case message := <-manager.Event.BroadcastChan:
			// 广播事件
			manager.sendAll(message)
		}
	}
}

// GetManagerInfo
//
//	@Description: 获取管理者信息
//	@receiver manager
//	@param isDebug
//	@return managerInfo
func (manager *Manager) GetManagerInfo(isDebug string) (managerInfo map[string]interface{}) {
	managerInfo = make(map[string]interface{})
	managerInfo["clientsLen"] = manager.GetAllMembersLen()               // 客户端连接数
	managerInfo["usersLen"] = manager.GetMembersLen()                    // 登录用户数
	managerInfo["chanRegisterLen"] = len(manager.Event.RegisterChan)     // 未处理连接事件数
	managerInfo["chanUnregisterLen"] = len(manager.Event.UnregisterChan) // 未处理退出登录事件数
	managerInfo["chanBroadcastLen"] = len(manager.Event.BroadcastChan)   // 未处理广播事件数
	if isDebug == "true" {
		addrList := make([]string, 0)
		manager.MembersRange(func(client *Member, value bool) (result bool) {
			addrList = append(addrList, client.Addr)
			return true
		})
		managerInfo["clients"] = addrList // 客户端列表
	}
	return
}

// ClearTimeoutConnections
//
//	@Description: 定时清理超时连接
//	@receiver manager
func (manager *Manager) ClearTimeoutConnections() {
	if manager == nil {
		return
	}
	currentTime := time.Now().Unix()

	clients := manager.GetAllMember()
	for client := range clients {
		if client.IsHeartbeatTimeout(currentTime, heartbeatTimeout) {
			_ = client.close
		}
	}
}
