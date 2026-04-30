package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 负责测试 datapack 拆包 封包的单元测试
func TestDataPack(t *testing.T) {
	// 模拟服务器
	//1 创建 sockerTCP
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	// 创建一个 go 承载负责从客户端处理业务
	go func() {
		for {
			//2 从客户端读取数据 拆包处理
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept err:", err)
				return
			}
			go func(conn net.Conn) {
				// 处理客户端请求
				// ------> 拆包的过程 <------
				// 定义一个拆包的对象
				dp := NewDataPack()
				for {
					// 第一次从 conn 读 把包的head 读出来
					headData := make([]byte, dp.GetHeadLen())
					if _, err := io.ReadFull(conn, headData); err != nil {
						fmt.Println("read head err:", err)
						break
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err:", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						// msg 只有数据的 需要第二次读取
						// 第二次从 conn 读 把head的datalen 再读取data
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						if _, err := io.ReadFull(conn, msg.Data); err != nil {
							fmt.Println("server unpack data err:", err)
							return
						}
						// 完整的一个消息已经读取完毕了
						fmt.Println("------>Recv MsgID: ", msg.Id,
							",datalen: ", msg.DataLen,
							",data: ", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	// 模拟客户端
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		fmt.Println("server connect err:", err)
		return
	}

	// 创建一个封包对象 dp
	dp := NewDataPack()
	// 模拟粘包过程， 封装俩个 msg 一同发送
	// 封装第一个msg1
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}
	// 封装第二个msg2
	msg2 := &Message{
		Id:      2,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 err:", err)
		return
	}

	// 将俩个包粘一起
	sendData1 = append(sendData1, sendData2...)

	// 一次性发送给服务端
	conn.Write(sendData1)

	// 客户端阻塞
	select {}
}
