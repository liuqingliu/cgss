// IPC框架的服务器端

package ipc

import (
	"encoding/json"
	"fmt"
)

type Request struct { //请求内容结构体
	Method string "method" //json字段
	Params string "params"
}

type Response struct { //响应内容结构体
	Code string "code"
	Body string "body"
}

type Server interface { //服务端接口
	Name() string                           //名称string
	Handle(method, params string) *Response //处理(方法，参数) *Response
}

type IpcServer struct { //IpcServer服务内容结构体
	Server //服务端接口，可以定义多个服务层
}

func NewIpcServer(server Server) *IpcServer { //创建IpcServer方法，返回值为*ipcserver
	return &IpcServer{server} //相当于给IpServer结构体赋值,Server:server
}

/*
结构体 server 有一个 方法 connect  ， 这个方法的返回值 有两个， 一个 是chan 类型， 一个是string 类型
*/
func (server *IpcServer) Connect() chan string {
	session := make(chan string, 0) //创建一个session是string类型的通道，长度为0的缓冲区
	go func(c chan string) {        // 另起轻量级线程实现goroutine，匿名函数，传入参数c，c为string类型的chan通道
		for {
			request := <-c          //从c中取数据到request
			if request == "CLOSE" { //判断从通道取到的是否为close
				break
			}
			var req Request //定义req为Request类型
			//Unmarshal用于反序列化json的函数 根据data将数据反序列化到传入的对象中
			err := json.Unmarshal([]byte(request), &req) //将json反序列化成struct对象 req.method,req.param
			if err != nil {
				fmt.Println("Invalid request format:", request)
			}
			resp := server.Handle(req.Method, req.Params) //resp是server方法handle处理后的返回值,类型为*response，req.method,req.params为传入参数
			//Marshal 用于将struct对象序列化到json对象中
			b, err := json.Marshal(resp) //将resp内容变为b的json串中
			c <- string(b)               //写入通道c，内容为b的json串
		}
		fmt.Println("Session closed.")
	}(session)

	fmt.Println("A new session has been created successfully.")
	return session
}
