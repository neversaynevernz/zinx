package ziface

import "net"

// 链接模块的抽象层

type IConn interface {

	// 启动链接 让当前的链接准备开始工作
	Start()

	// 关闭链接 结束当前链接的工作
	Stop()

	// 获取当前链接的对象 套接字
	GetTCPConnection() *net.TCPConn

	// 得到当前链接模块的链接ID
	GetConnID() uint32

	// 得到客户端链接的地址和端口
	RemoteAddr() net.Addr

	// 发送数据的方法
	Send(data []byte) error
}

// 定义一个处理业务链接的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
