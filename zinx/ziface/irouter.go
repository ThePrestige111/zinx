package ziface

/* 路由模块 */

type IRouter interface {
	// PreHandle 处理Conn业务之前的钩子方法hook
	PreHandle(request IRequest)

	// Handle 处理Conn业务中的主方法hook
	Handle(request IRequest)

	// PostHandle 处理Conn业务之后的钩子方法Hook
	PostHandle(request IRequest)
}
