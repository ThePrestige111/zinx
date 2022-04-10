package apis

import (
	"MMO_Game/core"
	__ "MMO_Game/pb"
	"fmt"
	"google.golang.org/protobuf/proto"
	"zinx/ziface"
	"zinx/znet"
)

// WorldChatApi 世界聊天
type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	// 解析客户端传递进来的proto协议
	protoMsg := &__.Talk{}
	unmarshallErr := proto.Unmarshal(request.GetData(), protoMsg)
	if unmarshallErr != nil {
		fmt.Println("unmarshall talk err: ", unmarshallErr)
	}

	// 当前聊天数据是哪个玩家发送的
	pid, getPropertyErr := request.GetConnection().GetProperty("pid")
	if getPropertyErr != nil {
		fmt.Println("getProperty err:", getPropertyErr)
	}

	// 根据pid得到对应的player对象
	player := core.WorldManagerObj.GetPlayerByPid(pid.(int32))

	// 将这个消息传递给其他全部在线的玩家
	player.Talk(protoMsg.Content)
}
