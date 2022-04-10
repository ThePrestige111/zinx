package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 单元测试
func TestDataPack(t *testing.T) {
	// 模拟的服务器
	listener, ListenErr := net.Listen("tcp", "127.0.0.1:7777")
	if ListenErr != nil {
		fmt.Println("server listen err:", ListenErr)
	}

	go func() {
		conn, AcceptErr := listener.Accept()
		if AcceptErr != nil {
			fmt.Println("server accept error", AcceptErr)
		}

		go func(conn net.Conn) {
			// 拆包过程
			dp := NewDataPack()

			for {
				// 第一次从conn中读，读取head
				head := make([]byte, dp.GetHeadLen())
				_, readErr := io.ReadFull(conn, head)
				if readErr != nil {
					fmt.Println("read head error...")
					break
				}

				msgHead, unpackErr := dp.Unpack(head)
				if unpackErr != nil {
					fmt.Println("server unpack error...")
					return
				}

				// 第二次从conn中读，根据head中的dataLen，再读取data中的内容
				if msgHead.GetMsgLen() > 0 {
					msg := msgHead.(*Message) // 转换类型
					msg.Data = make([]byte, msg.GetMsgLen())

					//根据dataLen的长度再次从io流中读取
					_, rErr := io.ReadFull(conn, msg.Data)
					if rErr != nil {
						fmt.Println("server unpack data err: ", rErr)
						return
					}

					// 读取完毕
					fmt.Println("msgId: ", msg.GetMsgID(), ", dataLen:", msg.GetMsgLen(), ", data:", string(msg.GetMsgData()))
				}
			}
		}(conn)
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err")
	}

	// 创建一个封包对象
	dp := NewDataPack()

	// 模拟粘包过程，封装两个msg一起发
	msg1 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte("hello"),
	}
	send1, _ := dp.Pack(msg1)
	msg2 := &Message{
		Id:      2,
		DataLen: 4,
		Data:    []byte("zinx"),
	}
	send2, _ := dp.Pack(msg2)

	send1 = append(send1, send2...)
	conn.Write(send1)

	// 客户端阻塞
	select {}
}
