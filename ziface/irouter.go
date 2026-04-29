package ziface

// 路由抽象接口
// 路由的数据都是 IRequest
type IRouter interface {
	// 在处理conn业务之前的钩子方法 Hook
	PreHandle(request IRequest)
	// 在处理conn业务子方法 Hook
	Handle(request IRequest)
	// 在处理conn业务之后的钩子方法 Hook
	PostHandle(request IRequest)
}
