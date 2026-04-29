package znet

import (
	"fmt"
	"github.com/neversaynevernz/zinx/utils"
	"github.com/neversaynevernz/zinx/ziface"
	"net"
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

	// 当前的 server添加 router,
	// server注册的链接对应的业务
	Router ziface.IRouter
}

// 启动服务器
func (s *Server) Start() error {
	fmt.Printf("[Zinx] Server Name : %s, listenner at IP: %s, Port: %d is starting\n",
		utils.GlobalObject.Name,
		utils.GlobalObject.Host,
		utils.GlobalObject.TcpPort,
	)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize,
	)

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

		fmt.Println("start zinx server success", s.Name, "Listening...")

		// 阻塞的等待客户端链接 处理客户端链接业务(读写)
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 初始化ID
			var cid uint32
			cid = 0

			// 新链接业务方法 和 conn 绑定
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			go dealConn.Start()
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

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("[Zinx] Add Router Success")
}

/*
初始化Server模块
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
	return s
}
