package ziface

/* 连接管理模块 */

type IConnectionManager interface {
	// Add 添加连接
	Add(conn IConnection)

	// Remove 删除连接
	Remove(conn IConnection)

	// GetConn 根据connID获取连接
	GetConn(connID uint32) (IConnection, error)

	// LenConn 得到当前连接数
	LenConn() int

	// ClearConn 清除并终止所有连接
	ClearConn()
}
