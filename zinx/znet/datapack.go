package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

type DataPack struct{}

// NewDataPack 拆包封包实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包的头部长度
func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen uint32 4 字节 + DataID uint32 4 字节
	return 8
}

// Pack 封包方法
func (dp *DataPack) Pack(msg ziface.Imessage) ([]byte, error) {
	// 创建一个存放Bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将dataLen写进dataBuff中
	lenErr := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen())
	if lenErr != nil {
		return nil, lenErr
	}

	// 将MsgID写进dataBuff中
	IDErr := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID())
	if IDErr != nil {
		return nil, IDErr
	}

	// 将MsgData写进dataBuff中
	DataErr := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgData())
	if DataErr != nil {
		return nil, DataErr
	}
	return dataBuff.Bytes(), nil
}

// Unpack 拆包方法
func (dp *DataPack) Unpack(binaryData []byte) (ziface.Imessage, error) {
	// 创建一个存放二进制数据的IOReader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压head信息，得到dataLen和MsgID
	msg := &Message{}

	// 读dataLen
	LenErr := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen)
	if LenErr != nil {
		return nil, LenErr
	}

	// 读MsgID
	IDErr := binary.Read(dataBuff, binary.LittleEndian, &msg.Id)
	if IDErr != nil {
		return nil, IDErr
	}

	// 判断DataLen是否超过我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize >= 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("package too large")
	}

	return msg, nil
}
