package ziface

type IServer interface {

	// 启动服务器
	Start() error

	// 停止服务器
	Stop() error	

	// 运行服务器
	Serve() error
}	

