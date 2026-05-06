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

	fmt.Println("Call PingRouter Handle...")
	// 先读取客户端的消息 再回写 ping...ping...ping

	fmt.Println("recv from client: msgID= ", request.GetMsgID())
	fmt.Println("recv from client: data= ", string(request.GetData()))

	request.GetConnection().SendMsg(0, []byte("ping...ping...ping "))
}

// ping test 自定义累 继承 BaseRouter
type HelloZinxRouter struct {
	znet.BaseRouter
}

// Test PreHandle
func (p *HelloZinxRouter) Handle(request ziface.IRequest) {

	fmt.Println("Call HelloZinxRouter Handle...")

	fmt.Println("recv from client: msgID= ", request.GetMsgID())
	fmt.Println("recv from client: data= ", string(request.GetData()))

	request.GetConnection().SendMsg(1, []byte("Hello, Welcome to Zinx!"))
}

// 创建链接之后的执行钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("conn DoConnectionBegin")
	if err := conn.SendMsg(202, []byte("DoConnectionBegin")); err != nil {
		fmt.Println(err)
	}
}

// 销毁链接之前的执行钩子函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("conn DoConnectionLost, coonid ", conn.GetConnID())
	// 广播下线 类似的

}

func main() {
	// 1 创建一个 server 句柄 使用 Zinx 的 API
	s := znet.NewServer("[Zinx V1.8.0]")

	// 2 注册链接的Hook
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 3 给当前 Zinx 框架添加一个自定义的 Router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	// 4 启动 server
	s.Serve()
}
