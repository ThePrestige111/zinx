package apis

import (
	"MMO_Game/core"
	__ "MMO_Game/pb"
	"fmt"
	"google.golang.org/protobuf/proto"
	"zinx/ziface"
	"zinx/znet"
)

type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	// 解析客户端传来的协议
	protoMsg := &__.Position{}
	unmarshallErr := proto.Unmarshal(request.GetData(), protoMsg)
	if unmarshallErr != nil {
		fmt.Println("Move Position Unmarshal err: ", unmarshallErr)
		return
	}

	// 得到当前发送消息的是哪个玩家
	pid, getPropertyErr := request.GetConnection().GetProperty("pid")
	if getPropertyErr != nil {
		fmt.Println("Get Player ID err: ", getPropertyErr)
		return
	}

	fmt.Printf("Player pid = %d, move(%f, %f, %f, %f)\n", pid, protoMsg.X, protoMsg.Y, protoMsg.Z, protoMsg.V)
	player := core.WorldManagerObj.GetPlayerByPid(pid.(int32))

	// 给其他玩家广播当前玩家的位置消息
	player.UpdatePos(protoMsg.X, protoMsg.Y, protoMsg.Z, protoMsg.V)
}
