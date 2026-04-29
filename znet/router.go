package znet

import "neversaynevernz/zinx/ziface"

// 实现 router 时先嵌入 BaseRouter基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct{}

// 在处理conn业务之前的钩子方法 Hook
func (r *BaseRouter) PreHandle(request ziface.IRequest) {

}

// 在处理conn业务子方法 Hook
func (r *BaseRouter) Handle(request ziface.IRequest) {

}

// 在处理conn业务之后的钩子方法 Hook
func (r *BaseRouter) PostHandle(request ziface.IRequest) {

}
