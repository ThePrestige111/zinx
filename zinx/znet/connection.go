package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

// Connection  连接模块
type Connection struct {
	TCPServer    ziface.IServer         //当前conn隶属于哪个Server
	Conn         *net.TCPConn           // 当前连接的TCP socket
	ConnID       uint32                 // 连接ID
	IsClosed     bool                   // 当前连接状态
	ExitChan     chan bool              // 告知当前连接已经退出的channel
	msgChan      chan []byte            // 无缓冲通道，用于读写goroutine之间的通信
	MsgHandler   ziface.IMsgHandler     // 消息的管理msgID和对应的处理业务API关系
	property     map[string]interface{} // 链接属性集合
	propertyLock sync.RWMutex           // 保护链接属性的锁
}

// NewConnection 初始化连接模块
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TCPServer:  server,
		Conn:       conn,
		ConnID:     connID,
		IsClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		MsgHandler: msgHandler,
		property:   make(map[string]interface{}),
	}

	// 将conn加入到ConnectionManager
	c.TCPServer.GetConnectionManager().Add(c)

	return c
}

// StartWriter 专门给客户端发送消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("Writer Goroutine is running...")
	defer fmt.Println("[Writer is exit!] ", c.GetRemoteAddr().String())

	// 不断地阻塞等待channel的消息，写给客户端
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data err ", err)
				return
			}
		case <-c.ExitChan:
			// 代表reader已经退出，此时writer也要退出
			return

		}
	}
}

// Start 启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)

	// 启动当前读数据的业务
	go c.StartReader()
	// 启动当前写数据的业务
	go c.StartWriter()

	// 按照开发者传递进来的 创建链接之后需要调用的处理业务，执行对应的Hook函数
	c.TCPServer.CallOnConnStart(c)
}

// StartReader 读取业务
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("[Reader is exit!] connID = ", c.ConnID, " remote addr is", c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		//// 读取客户端的数据到Buff中，最大512字节
		//buf := make([]byte, 512)
		//_, err := c.Conn.Read(buf)
		//
		//if err != nil {
		//	fmt.Println("receive buf err ", err)
		//	continue
		//}

		// 创建一个拆包解包对象
		dp := NewDataPack()

		// 读取客户端的Msg Head的二进制流的前8个字节
		headData := make([]byte, dp.GetHeadLen())
		_, ReadHeadErr := io.ReadFull(c.GetTCPConnection(), headData)
		if ReadHeadErr != nil {
			fmt.Println("read headData err")
			break
		}

		// 拆包，得到msgID和msgData放于msg中
		msg, unPackErr := dp.Unpack(headData)
		if unPackErr != nil {
			fmt.Println("unpack error", unPackErr)
			break
		}

		// 根据dataLen 再次读取Data
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			_, readDataErr := io.ReadFull(c.GetTCPConnection(), data)
			if readDataErr != nil {
				fmt.Println("read data err ", readDataErr)
				break
			}
		}

		msg.SetMsgData(data)
		// 得到当前Conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中找到绑定的conn对应的router调用
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// Stop 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)

	// 如果当前连接已经关闭
	if c.IsClosed == true {
		return
	}

	c.IsClosed = true

	// 调用开发者注册的 销毁链接之前的Hook函数
	c.TCPServer.CallOnConnStop(c)

	// 关闭socket连接
	err := c.Conn.Close()
	if err != nil {
		fmt.Printf("stop failed")
		return
	}

	// 告知writer关闭
	c.ExitChan <- true

	// 将当前连接从ConnectionManager中移除
	c.TCPServer.GetConnectionManager().Remove(c)

	// 关闭管道
	close(c.ExitChan)
	close(c.msgChan)
}

// GetTCPConnection 获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接的模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// GetRemoteAddr 获取远程客户端的TCP状态 IP port
func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg  发送数据
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.IsClosed == true {
		return errors.New("connection closed when send msg")
	}

	// 将data进行封包
	dp := NewDataPack()

	binaryMsg, packErr := dp.Pack(NewMsgPackage(msgID, data))
	if packErr != nil {
		fmt.Println("pack error msg id = ", msgID)
		return errors.New("pack error msg")
	}

	// 将数据发送给客户端
	c.msgChan <- binaryMsg
	return nil
}

// SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}

	return nil, errors.New("no property found")
}

// RemoveProperty 删除链接属性
func (c *Connection) RemoveProperty(key string) error {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if _, ok := c.property[key]; ok {
		delete(c.property, key)
		return nil
	}
	return errors.New("remove property failed")
}
