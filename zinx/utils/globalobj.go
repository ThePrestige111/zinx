package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

/*	存储一切有关Zinx框架的全局参数，供其他用户配置 */

type GlobalObj struct {
	TcpServer ziface.IServer // 当前Zinx全局的Server对象
	Host      string         // 当前服务器主机监听的IP
	TcpPort   int            // 当前服务器主机监听的端口号
	Name      string         // 当前服务器的名称

	Version          string // 当前zinx的版本号
	MaxConn          int    // 当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 // 当前zinx框架数据包的最大值
	WorkerPoolSize   uint32 // 当前worker工作池的中goroutine的数量
	MaxWorkerTaskLen uint32 // 每个worker对应消息队列的任务数量的最大值
}

// GlobalObject 定义一个全局变量
var GlobalObject *GlobalObj

// Reload 从zinx.json去加载用户自定义参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// json文件数据解析到struct
	UnmarshalErr := json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(UnmarshalErr)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		// 如果配置文件没有加载，默认的值
		Name:             "ZinxServerApp",
		Version:          "V0.4",
		TcpPort:          9090,
		Host:             "0.0.0.1",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1000,
	}

	//从config/zinx.json加载数据
	GlobalObject.Reload()
}
