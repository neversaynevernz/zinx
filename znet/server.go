package znet

import (
	"fmt"
	"net"
	"neversaynevernz/zinx/ziface"
)

type Server struct {

	// 服务器的名称
	Name string

	// 服务器的IP版本
	IPVersion string

	// 服务器监听的IP
	IP string

	// 服务器监听的端口
	Port int
}

// 启动服务器
func (s *Server) Start() error {

	fmt.Printf("[Zinx] Server Name: %s, listener at IP: %s, Port: %d is starting...\n", s.Name, s.IP, s.Port)

	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}

		// 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen tcp error: ", err)
			return
		}
		defer listener.Close()

		// 阻塞的等待客户端链接 处理客户端链接业务(读写)
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 已经建立链接 做一个回显业务
			// 这里我们又开了一个协程 处理客户端的链接业务
			// 这样做的好处是 可以同时处理多个客户端链接
			go s.handleConnection(conn)
		}
	}()

	return nil
}

// 停止服务器
func (s *Server) Stop() error {
	// TODO 将一些服务器的资源或者一些已经开辟的链接进行停止或者回收
	return nil
}

// 运行服务器
func (s *Server) Serve() error {

	s.Start()

	// TODO 做一些服务器启动之后的额外业务

	select {}

	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return
			continue
		}

		fmt.Printf("recv client buf %s, cnt: %d\n", buf[:n], n) // 回显业务
		_, err = conn.Write(buf[:n])
		if err != nil {
			return
		}
	}
}

/*
初始化Server模块
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
