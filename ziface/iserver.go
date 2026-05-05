package ziface

type IServer interface {

	// 启动服务器
	Start() error

	// 停止服务器
	Stop()

	// 运行服务器
	Serve() error

	// 路由功能: 给当前的服务注册一个路由方法 供客户端的链接处理调用
	AddRouter(msgid uint32, router IRouter)

	//获取当前server 的链接管理器
	GetConnManager() IConnManager
}
