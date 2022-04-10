package main

import (
	"MMO_Game/apis"
	"MMO_Game/core"
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

// 当前客户端建立连接后的hook函数
func onConnectionAdd(conn ziface.IConnection) {
	// 创建一个player对象
	player := core.NewPlayer(conn)

	// 发送客户端MsgID：1的消息: 同步当前player的id给客户端
	player.SyncPid()

	// 发送客户端MsgID：200的消息：同步当前player的初始位置给客户端
	player.BroadCastStartPosition()

	// 将新上线的玩家添加到world manager中
	core.WorldManagerObj.AddPlayer(player)

	// 将该链接绑定一个玩家ID,记录当前连接是属于哪个玩家的
	conn.SetProperty("pid", player.Pid)

	// 同步周边玩家，告知他们当前玩家已上线，广播当前玩家的位置信息
	player.SyncSurrounding()

	fmt.Println(">>>>> Player ", player.Pid, " is online <<<<<")
}

// 给当前链接断开之前出发的hook函数
func onConnectionLost(conn ziface.IConnection) {
	// 获取该链接对应的玩家
	pid, err := conn.GetProperty("pid")
	if err != nil {
		fmt.Println("Get Player err: ", err)
		return
	}
	player := core.WorldManagerObj.GetPlayerByPid(pid.(int32))

	// 玩家下线
	player.OffLine()

	fmt.Println(">>>>> Player ", player.Pid, " is offline <<<<<")
}

func main() {
	// 创建zinx server句柄
	s := znet.NewServer("MMO Game Zinx")

	// 链接创建和销毁的钩子函数
	s.SetOnConnStart(onConnectionAdd)
	s.SetOnConnStop(onConnectionLost)

	// 注册一些路由业务
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})

	//启动服务
	s.Server()
}
