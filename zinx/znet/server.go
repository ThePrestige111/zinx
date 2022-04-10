package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

// Server IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的ip版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int
	// 当前server消息管理模块，用来绑定msgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandler
	// 该Server的连接管理器
	ConnManager ziface.IConnectionManager
	// 该Server创建连接后自动调用Hook函数-OnConnStart
	OnConnStart func(conn ziface.IConnection)
	// 该Server销毁连接后自动调用Hook函数-OnConnStop
	OnConnStop func(conn ziface.IConnection)
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Success")
}

// NewServer 初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
	}
	return s
}

// SetOnConnStart 注册SetOnConnStart钩子函数方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 注册OnConnStop钩子函数方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用OnConnStart钩子函数的方法
func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println(">>>>>Call OnConnStart<<<<<")
		s.OnConnStart(connection)
	}
}

// CallOnConnStop 调用OnConnStop钩子函数方法
func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println(">>>>>Call OnConnStop<<<<<")
		s.OnConnStop(connection)
	}
}

// Start 开启服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, listenning at %s, Port is %d\n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)

	// 将Start从同步改为异步
	go func() {
		// 0. 开启工作池
		s.MsgHandler.StartWorkerPool()

		// 1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err:", err)
			return
		}
		// 2. 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err ", err)
			return
		}

		var cid uint32 = 0
		fmt.Println("start Zinx server successfully, ", s.Name, "is listening")

		// 3. 阻塞的等待客户端连接，处理客户端连接业务
		for {
			// 如果有客户端连接进来，阻塞返回
			conn, e := listener.AcceptTCP()
			if e != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 判断是否超过最大连接
			if s.ConnManager.LenConn() > utils.GlobalObject.MaxConn {
				fmt.Println("====== connections out of maxConnection =======")
				conn.Close()
				continue
			}

			// 将处理新连接的方法和 conn进行绑定 得到连接拓展模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动当前连接的业务
			go dealConn.Start()
		}
	}()
}

// Stop 关闭服务器
func (s *Server) Stop() {
	s.ConnManager.ClearConn()
	fmt.Println("[STOP] Zinx Server", s.Name)
}

// Server 运行服务器
func (s *Server) Server() {
	s.Start()

	// 阻塞状态
	select {}
}

// GetConnectionManager 获取当前服务器的ConnectionManager
func (s *Server) GetConnectionManager() ziface.IConnectionManager {
	return s.ConnManager
}
