package znet

import (
	"fmt"
	"github.com/neversaynevernz/zinx/utils"
	"github.com/neversaynevernz/zinx/ziface"
	"net"
	"runtime"
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

	// 当前的 server 的消息管理模块
	// 用来绑定MsgID和对应的处理业务API 关系
	MsgHandler ziface.IMsgHandler

	// 该 server 的链接管理器
	ConnMgr ziface.IConnManager

	// 该server创建链接之后自动调用Hook函数 --OnConnStart
	OnConnStart func(conn ziface.IConnection)

	// 该server销毁链接之后自动调用Hook函数 --OnConnStop
	OnConnStop func(conn ziface.IConnection)
}

// 启动服务器
func (s *Server) Start() error {

	runtime.GOMAXPROCS(1)

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

		// 0 开启消息队列及Worker工作池
		s.MsgHandler.StartWorkerPool()

		// 1 获取一个TCP的Addr
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

		// 初始化ID
		var cid uint32
		cid = 0

		// 阻塞的等待客户端链接 处理客户端链接业务(读写)
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 设置最大链接个数的判断. 如果超过最大连接数则关闭此新的链接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//TODO 给客户端相应一个超出最大链接的错误包
				fmt.Println("[Zinx] MaxConn exceeded = ", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			// 新链接业务方法 和 conn 绑定
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			go dealConn.Start()
		}
	}()

	return nil
}

// 停止服务器
func (s *Server) Stop() {
	// TODO 将一些服务器的资源或者一些已经开辟的链接进行停止或者回收
	fmt.Println("[Stop] Zinx Server Name:", utils.GlobalObject.Name)
	s.ConnMgr.ClearConn()
}

// 运行服务器
func (s *Server) Serve() error {

	s.Start()

	// TODO 做一些服务器启动之后的额外业务

	select {}

	return nil
}

// 路由功能: 给当前的服务注册一个路由方法 供客户端的链接处理调用
func (s *Server) AddRouter(msgid uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgid, router)
	fmt.Println("[Zinx] Add Router Success")
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnMgr
}

/*
初始化Server模块
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

// 注册 OnConnStart钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 注册 OnConnStart钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用 OnConnStop钩子函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("[Zinx] Call OnConnStart success")
		s.OnConnStart(conn)
	}
}

// 调用 OnConnStop钩子函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("[Zinx] Call OnConnStop success")
		s.OnConnStop(conn)
	}
}
