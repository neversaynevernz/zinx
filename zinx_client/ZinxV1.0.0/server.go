package main

// local debug
//import "neversaynevernz/zinx/znet"

// 远程仓库的包
import "github.com/neversaynevernz/zinx/znet"

func main() {
	// 1创建服务端
	s := znet.NewServer("[Zinx V1.1.0]")

	// 启动 server
	s.Serve()
}
