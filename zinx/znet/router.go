package znet

import "zinx/ziface"

// BaseRouter 实现Router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类的方法进行重写
// Router继承BaseRouter的好处就是，不需要用户将所有方法都实现
type BaseRouter struct{}

// PreHandle 处理Conn业务之前的钩子方法hook
func (be *BaseRouter) PreHandle(request ziface.IRequest) {}

// Handle 处理Conn业务中的主方法hook
func (be *BaseRouter) Handle(request ziface.IRequest) {}

// PostHandle 处理Conn业务之后的钩子方法Hook
func (be *BaseRouter) PostHandle(request ziface.IRequest) {}
