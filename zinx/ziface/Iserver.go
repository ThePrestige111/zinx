package ziface

/* 定义服务器接口模块 */

type IServer interface {
	// Start 开启服务器
	Start()

	// Stop 关闭服务器
	Stop()

	// Server 运行服务器
	Server()

	// AddRouter 路由功能:给当前的服务注册一个路由方法，供客户端的处理使用
	AddRouter(msgID uint32, router IRouter)

	// GetConnectionManager 获取当前服务器的ConnectionManager
	GetConnectionManager() IConnectionManager

	// SetOnConnStart 注册SetOnConnStart钩子函数方法
	SetOnConnStart(func(connection IConnection))

	// SetOnConnStop 注册OnConnStop钩子函数方法
	SetOnConnStop(func(connection IConnection))

	// CallOnConnStart 调用OnConnStart钩子函数的方法
	CallOnConnStart(connection IConnection)

	// CallOnConnStop 调用OnConnStop钩子函数方法
	CallOnConnStop(connection IConnection)
}
