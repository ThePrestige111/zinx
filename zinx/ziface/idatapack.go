package ziface

/* 封包、拆包模块，处理TCP粘包问题 */

type IDataPack interface {
	// GetHeadLen 获取包的头部长度
	GetHeadLen() uint32

	// Pack 封包方法
	Pack(msg Imessage) ([]byte, error)

	// Unpack 拆包方法
	Unpack([]byte) (Imessage, error)
}
