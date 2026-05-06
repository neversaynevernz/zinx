package main

import (
	"fmt"
	"github.com/neversaynevernz/zinx/znet"
	"io"
	"net"
	"time"
)

// 客户端模拟
func main() {

	fmt.Println("client1 start...")

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	for {

		// 发送封包的 msg 消息
		// msg 0 的消息
		dp := znet.NewDataPack()

		binary, err := dp.Pack(znet.NewMessage(1, []byte("Zinx v1.5.0 client0 test message")))
		if err != nil {
			fmt.Println("pack error:", err)
			return
		}

		if _, err := conn.Write(binary); err != nil {
			fmt.Println("write error:", err)
			return
		}

		// 服务器应该回复一个 message数据 例如 msgid:1 ping ping ping

		// 读取流中的 head 部分，得到ID 和datalen

		//再根据Datalen在进行二次读取

		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := conn.Read(binaryHead); err != nil {
			fmt.Println("read head error:", err)
			return
		}

		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msghead error:", err)
			return
		}

		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error:", err)
			}

			fmt.Println("------>Recv Server MsgID: ", msg.Id,
				",datalen: ", msg.DataLen,
				",data: ", string(msg.Data))

		}

		// cpu 阻塞
		time.Sleep(1 * time.Second)
	}
}
