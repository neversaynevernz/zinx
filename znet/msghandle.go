package znet

import (
	"fmt"
	"github.com/neversaynevernz/zinx/ziface"
)

/*
消息处理模块的实现
*/
type MsgHandle struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

// 初始化MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// 调度/执行对应的 Router 消息处理方法
func (mh *MsgHandle) DoMsgHandler(r ziface.IRequest) {

	// 从 request中找到msgID
	handler, ok := mh.Apis[r.GetMsgID()]
	if !ok {
		fmt.Println("no handler， please register first! msgid=", r.GetMsgID())
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
		panic("Repeat Register Api, msg id : " + fmt.Sprint(msgId))
	}

	// 注册路由
	//添加 msg 与API 绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add Router msgId :", msgId, "success")
}
