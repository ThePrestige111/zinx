package core

import (
	__ "MMO_Game/pb"
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"sync"
	"zinx/ziface"
)

type Player struct {
	Pid  int32              // 玩家id
	Conn ziface.IConnection // 当前玩家的连接
	X    float32            // 平面X坐标
	Y    float32            // 高度
	Z    float32            // 平面Y坐标
	V    float32            // 旋转的0-360角度
}

var PidGen int32 = 1  // 用来生成玩家id的计数器
var IDLock sync.Mutex // 保护pid的锁

// NewPlayer 创建一个玩家
func NewPlayer(conn ziface.IConnection) *Player {
	// 生成一个玩家ID
	IDLock.Lock()
	id := PidGen
	PidGen++
	IDLock.Unlock()

	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(20)), // 随机在160坐标点，基于X轴若干偏移
		Y:    0,
		Z:    float32(140 + rand.Intn(10)), // 随机在140坐标点，基于Y轴若干偏移
		V:    0,
	}
	return p
}

// SendMessage 发送给客户端消息的方法，主要是将protobuf数据序列化后，再调用zinx的sendMsg方法
func (p *Player) SendMessage(msgID uint32, data proto.Message) {
	// 将Proto Message结构体序列化
	msg, marshalErr := proto.Marshal(data)
	if marshalErr != nil {
		fmt.Println("marshal msg err: ", marshalErr)
		return
	}

	// 将二进制文件 通过zinx框架的SendMsg将数据发送给客户端
	if p.Conn == nil {
		fmt.Println("connection is off")
		return
	}

	sendMsgErr := p.Conn.SendMsg(msgID, msg)
	if sendMsgErr != nil {
		fmt.Println("send message err: ", sendMsgErr)
	}
	return
}

// SyncPid 告知客户端玩家pid，同步已经生成的玩家id给客户端
func (p *Player) SyncPid() {
	protoMsg := &__.SyncPid{
		Pid: p.Pid,
	}

	p.SendMessage(1, protoMsg)
}

// BroadCastStartPosition 广播玩家自己的出生地点
func (p *Player) BroadCastStartPosition() {
	protoMsg := &__.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &__.BroadCast_P{
			P: &__.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	p.SendMessage(200, protoMsg)
}

// Talk 玩家广播世界聊天消息
func (p *Player) Talk(content string) {
	//1 组建MsgID：200 proto数据
	protoMsg := &__.BroadCast{
		Pid: p.Pid,
		Tp:  1, //代表聊天广播
		Data: &__.BroadCast_Content{
			Content: content,
		},
	}

	// 得到当前世界所有的在线玩家
	players := WorldManagerObj.GetAllPlayers()

	// 向所有的玩家（包括自己）发送MsgID:200的消息
	for _, player := range players {
		// player分别给对应的客户端发送消息
		player.SendMessage(200, protoMsg)
	}
}

// SyncSurrounding 同步玩家上线的信息
func (p *Player) SyncSurrounding() {
	// 获取当前玩家周围有哪些
	pids := WorldManagerObj.AoiManager.GetPlayerIdsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldManagerObj.GetPlayerByPid(int32(pid)))
	}

	// 将当前玩家的位置通过MsgID:200 发送给周围的玩家（让其他人看到自己）
	// 组建当前玩家的位置消息
	protoMsg200 := &__.BroadCast{
		Pid: p.Pid,
		Tp:  2, // 广播坐标
		Data: &__.BroadCast_P{
			P: &__.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 告知客户端玩家pid，同步已经生成的玩家id给客户端
	for _, player := range players {
		player.SendMessage(200, protoMsg200)
	}

	// 将周围的全部玩家位置信息通过MsgID：202 发送给房前玩家的客户端（让自己看到其他人）
	// 组建MsgID:202的消息
	playersProtoMsg := make([]*__.Player, 0, len(players))
	for _, player := range players {
		onePlayer := &__.Player{
			Pid: player.Pid,
			P: &__.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		playersProtoMsg = append(playersProtoMsg, onePlayer)
	}

	// 封装SyncPlayer protoBuff数据
	protoMsg202 := &__.SyncPlayers{
		Ps: playersProtoMsg[:], // 拷贝数据
	}

	p.SendMessage(202, protoMsg202)
}

// UpdatePos 更新当前玩家的坐标
func (p *Player) UpdatePos(X, Y, Z, V float32) {
	// 更新当前玩家的坐标
	p.X = X
	p.Y = Y
	p.Z = Z
	p.V = V

	// 组建广播协议
	protoMsg := &__.BroadCast{
		Pid: p.Pid,
		Tp:  4, // 4-移动位置更新
		Data: &__.BroadCast_P{
			P: &__.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 获取当前玩家的九宫格之内的玩家
	players := p.GetSurroundingPlayers()

	// 一次给每个玩家对应的客户端发送当前玩家位置更新的消息
	for _, player := range players {
		player.SendMessage(200, protoMsg)
	}
}

// GetSurroundingPlayers 获取当前玩家周边的玩家
func (p *Player) GetSurroundingPlayers() []*Player {
	// 得到当前九宫格内所有玩家的pID
	pids := WorldManagerObj.AoiManager.GetPlayerIdsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))

	// 将所有pid对应的player放到players中
	for _, pid := range pids {
		player := WorldManagerObj.GetPlayerByPid(int32(pid))
		players = append(players, player)
	}
	return players
}

// OffLine 玩家下线
func (p *Player) OffLine() {
	// 得到当前玩家周边的九宫格内的都有哪些玩家
	players := p.GetSurroundingPlayers()

	// 给周围玩家广播MsgID:201的信息
	protoMsg := &__.SyncPid{
		Pid: p.Pid,
	}

	for _, player := range players {
		player.SendMessage(201, protoMsg)
	}

	// 从地图和服务器中删除玩家
	WorldManagerObj.AoiManager.RemovePidByPos(int(p.Pid), p.X, p.Z)
	WorldManagerObj.RemovePlayerByPid(p.Pid)
}
