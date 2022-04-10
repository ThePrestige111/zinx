package ziface

/* 客户端请求消息模块 将客户端请求的连接信息和请求的数据包装到一个Request中 */

type IRequest interface {
	// GetConnection 得到当前连接
	GetConnection() IConnection

	// GetData 得到消息数据
	GetData() []byte

	// GetMsgID 得到消息ID
	GetMsgID() uint32
}
