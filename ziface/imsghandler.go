package ziface

/*
消息管理抽象层
*/
type IMsgHandler interface {

	// 调度/执行对应的 Router 消息处理方法
	DoMsgHandler(request IRequest)

	//为消息添加具体的处理逻辑
	AddRouter(uint32, IRouter)

	// 启动Worker工作池
	StartWorkerPool()

	// 将消息交给消息任务队列处理
	SendMsgToTaskQueue(request IRequest)
}
