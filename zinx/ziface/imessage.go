package ziface

/* 消息模块 将请求的消息封装到一个Message中  */

type Imessage interface {
	GetMsgID() uint32   // 获取消息的ID
	GetMsgLen() uint32  // 获取消息的长度
	GetMsgData() []byte // 获取消息的数据
	SetMsgId(uint32)    // 设置消息的ID
	SetMsgLen(uint32)   // 设置消息的长度
	SetMsgData([]byte)  // 设置消息的数据
}
