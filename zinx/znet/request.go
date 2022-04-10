package znet

import "zinx/ziface"

type Request struct {
	// 已经个客户端建立的连接
	conn ziface.IConnection

	// 客户端请求的数据
	msg ziface.Imessage
}

// GetConnection 得到当前连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData 得到消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetMsgData()
}

// GetMsgID 得到消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
