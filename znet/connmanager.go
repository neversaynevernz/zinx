package znet

import (
	"errors"
	"fmt"
	"github.com/neversaynevernz/zinx/ziface"
	"sync"
)

/*
	链接管理模块
*/

type ConnManager struct {
	Connections map[uint32]ziface.IConnection // 管理的链接集合
	ConnLock    sync.RWMutex                  // 保护链接集合的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		Connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加链接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源 map, 加写锁
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()

	// 将 conn 加入 ConnManager 中
	cm.Connections[conn.GetConnID()] = conn

	fmt.Println("connID = ", conn.GetConnID(), "add to ConnManager successfully: conn num = ", cm.Len())
}

// 删除链接
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源 map, 加写锁
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()

	delete(cm.Connections, conn.GetConnID())
	fmt.Println("connID = ", conn.GetConnID(), "remove from ConnManager successfully: conn num = ", cm.Len())
}

// 获取链家
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源 map, 加读锁
	cm.ConnLock.RLock()
	defer cm.ConnLock.RUnlock()
	if c, ok := cm.Connections[connID]; ok {
		return c, nil
	}
	return nil, errors.New("connection not found")
}

// 链接的个数
func (cm *ConnManager) Len() int {
	return len(cm.Connections)
}

// 清除并终止所有链接
func (cm *ConnManager) ClearConn() {
	// 保护共享资源 map, 加写锁
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()

	// 删除 conn 并停止conn工作
	for connID, conn := range cm.Connections {
		// 停止
		conn.Stop()
		// 删除
		delete(cm.Connections, connID)
	}
	fmt.Println("Clear All connections successfully: conn num = ", cm.Len())

	// cm.Connections = make(map[uint32]ziface.IConnection)
}
