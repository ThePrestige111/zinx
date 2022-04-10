package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

/* 连接管理模块 */

type ConnectionManager struct {
	connections map[uint32]ziface.IConnection // 管理的连接信息集合
	conLock     sync.RWMutex                  //保护连接集合的读写锁
}

// NewConnManager 初始化
func NewConnManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// Add 添加连接
func (cm *ConnectionManager) Add(conn ziface.IConnection) {
	cm.conLock.Lock()
	defer cm.conLock.Unlock()

	// 将conn加入到ConnectionManager
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnectionManager successfully! connID =", conn.GetConnID())
}

// Remove 删除连接
func (cm *ConnectionManager) Remove(conn ziface.IConnection) {
	cm.conLock.Lock()
	defer cm.conLock.Unlock()

	// 将conn从ConnectionManager中移除
	delete(cm.connections, conn.GetConnID())
	fmt.Println("connID =", conn.GetConnID(), " remove from ConnectionManager")
}

// GetConn 根据connID获取连接
func (cm *ConnectionManager) GetConn(connID uint32) (ziface.IConnection, error) {
	cm.conLock.RLock()
	defer cm.conLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection is not found")
	}
}

// LenConn 得到当前连接数
func (cm *ConnectionManager) LenConn() int {
	return len(cm.connections)
}

// ClearConn 清除并终止所有连接
func (cm *ConnectionManager) ClearConn() {
	cm.conLock.Lock()
	defer cm.conLock.Unlock()

	for connID, conn := range cm.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(cm.connections, connID)
	}
	fmt.Println("Clear all connections success!")
}
