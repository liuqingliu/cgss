//IPC框架进行单元测试
package ipc

import (
	"testing"
)

type EchoServer struct { //定义一个空结构体EchoServer（输出服务端）
}

func (server *EchoServer) Handle(requst string) string { //函数Handle是结构体server接口的处理函数，传递参数request，响应值string类型
	return "ECHO:" + request
}

func (server *EchoServer) Name() string { //函数Name是结构体server的接口函数，返回值为string类型
	return "EchoServer"
}

func TestIpc(t *testing.T) { //函数testipc，传入参数是testing.T类型的
	server := NewIpcServer(&EchoServer{}) // 用NewIpcserver来创建一个server服务端

	client1 := NewIpcClient(server) // 用NewIpcClient来创建两个client客户端
	client2 := NewIpcClient(server)

	resp1 := client1.Call("From Client1") //客户端1/2分别输出不同信息，调用call函数
	resp2 := client2.Call("From Client2")

	if resp1 != "ECHO:From Client1" || resp2 != "ECHO:From Client2" { //判断返回值是否属于预期
		t.Error("IpcClient.Call failed. resp1:", resp1, "resp2:", resp2) //输出错误参数
	}

	client1.Close() //关闭2客户端
	client2.Close()

}
