package ziface

import "net"

/* 连接模块  */

type IConnection interface {
	// Start 启动连接 让当前的连接准备开始工作
	Start()

	// Stop 停止连接 结束当前连接的工作
	Stop()

	// GetTCPConnection 获取当前连接绑定的socket conn
	GetTCPConnection() *net.TCPConn

	// GetConnID 获取当前连接的模块的连接ID
	GetConnID() uint32

	// GetRemoteAddr 获取远程客户端的TCP状态 IP port
	GetRemoteAddr() net.Addr

	// SendMsg 发送数据
	SendMsg(msgID uint32, data []byte) error

	// SetProperty 设置链接属性
	SetProperty(key string, value interface{})

	// GetProperty 获取链接属性
	GetProperty(key string) (interface{}, error)

	// RemoveProperty 删除链接属性
	RemoveProperty(key string) error
}

// HandleFunc 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
