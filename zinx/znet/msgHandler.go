package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandler struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// worker工作池中worker的数量
	WorkPoolSize uint32
}

// NewMsgHandler 创建MsgHandle
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:         make(map[uint32]ziface.IRouter),
		TaskQueue:    make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}

// DoMsgHandler 调度/执行对应的Router消息处理方法
func (m *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	// 从Request中找到msgID
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is not found! Need Register!")
		return
	}
	// 根据msgID调度相应的方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (m *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	// 判断当前m绑定的API处理方法是否存在
	if _, ok := m.Apis[msgID]; ok {
		// 已经被注册
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}
	m.Apis[msgID] = router
	fmt.Println("Add api msgID = ", msgID, " success!")
}

// StartWorkerPool 启动一个Worker工作池（只发生一次）
func (m *MsgHandler) StartWorkerPool() {
	// 根据WorkPoolSize分别开启Worker，每个Worker代表一个goroutine
	for i := 0; i < int(m.WorkPoolSize); i++ {
		// 给当前worker开辟一个对应的channel：0号worker用0号channel
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go m.StartWorker(i, m.TaskQueue[i])
	}
}

// StartWorker 启动一个Worker
func (m *MsgHandler) StartWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started...")

	// 不断的阻塞等待对应消息队列的消息
	for {
		select {
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue，由Worker处理
func (m *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 将消息平均分配给worker
	// 根据客户端建立ConnID来进行分配
	workID := request.GetConnection().GetConnID() % m.WorkPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		" request msgID = ", request.GetMsgID(), " to Worker = ", workID)

	// 将消息发送给对应Worker的TaskQueue
	m.TaskQueue[workID] <- request
}
