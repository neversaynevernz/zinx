package znet

import (
	"fmt"
	"github.com/neversaynevernz/zinx/utils"
	"github.com/neversaynevernz/zinx/ziface"
)

/*
消息处理模块的实现
*/
type MsgHandle struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter

	// 负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest

	// 业务工作Worker池的 worker 数量
	WorkerPoolSize uint32
}

// 初始化MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, // 从全局配置中获取
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// 调度/执行对应的 Router 消息处理方法
func (mh *MsgHandle) DoMsgHandler(r ziface.IRequest) {

	// 从 request中找到msgID
	handler, ok := mh.Apis[r.GetMsgID()]
	if !ok {
		fmt.Println("no handler, please register first! msgid=", r.GetMsgID())
		return
	}

	// 根据msgID调度对应的router业务
	handler.PreHandle(r)
	handler.Handle(r)
	handler.PostHandle(r)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {

	// 判断 当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		// 已经注册
		panic("Repeat Register Api, msg id: " + fmt.Sprint(msgId))
	}

	// 注册路由
	//添加 msg 与API 绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add Router msgId :", msgId, "success")
}

// 启动一个Worker 工作池
// （开启工作池的动作只能发生一次， 一个zinx框架只能有一个worker工作池）
func (mh *MsgHandle) StartWorkerPool() {
	// 根据 workerPoolSize 分别开启Worker 每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个 worker 被启动
		// 1 当前的 worker 对应的channel消息队列 开辟空间 第0个worker 就用第0个channel...
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2 启动当前的 Worker  阻塞等待消息从channel传递过来
		go mh.SendOneWorker(i, mh.TaskQueue[i])
	}
}

// 去启动一个Worker工作流程
func (mh *MsgHandle) SendOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID=", workerID, "is starting...")
	// 不断的阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息过来，出列的就是一个客户端的 Request, 执行当前Request绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 将消息交给TaskQueue，由Worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1 将消息平均分配给不同的 worker
	// 根据客户端建立的ConnID 来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID =", request.GetConnection().GetConnID(),
		"request MsgID = ", request.GetMsgID(),
		"to WorkerID = ", workerID)
	// 2 将消息发送给对应的 worker的 TaskQueue 即可
	mh.TaskQueue[workerID] <- request
}
