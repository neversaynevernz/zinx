package main

import (
	"fmt"
	"github.com/neversaynevernz/zinx/ziface"
	"github.com/neversaynevernz/zinx/znet"
)

// ping test 自定义累 继承 BaseRouter
type PingRouter struct {
	znet.BaseRouter
}

// Test PreHandle
func (p *PingRouter) Handle(request ziface.IRequest) {

	fmt.Println("Call Router Handle...")
	// 先读取客户端的消息 再回写 ping...ping...ping

	fmt.Println("recv from client: msgID= ", request.GetMsgID())
	fmt.Println("recv from client: data= ", string(request.GetData()))

	request.GetConnection().SendMsg(1, []byte("ping...ping...ping "))
}

func main() {
	// 1 创建一个 server 句柄 使用 Zinx 的 API
	s := znet.NewServer("[Zinx V1.5.0]")

	// 给当前 Zinx 框架添加一个自定义的 Router
	s.AddRouter(&PingRouter{})

	// 启动 server
	s.Serve()
}
