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
func (p *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	//
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// Test PreHandle
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping...ping...ping... error")
	}
}

// Test PreHandle
func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping\n"))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

func main() {
	// 1 创建一个 server 句柄 使用 Zinx 的 API
	s := znet.NewServer("[Zinx V1.3.0]")

	// 给当前 Zinx 框架添加一个自定义的 Router
	s.AddRouter(&PingRouter{})

	// 启动 server
	s.Serve()
}
