//IPC框架的客户端
package ipc

import (
	"encoding/json"
)

type IpcClient struct { //定义ipc客户端内容结构体
	conn chan string // conn 为 string 类型的通道
}

//函数NewIpcClient新建Ipc客户端，传入参数为ipcserver（服务端）结构体，返回为客户端结构体
func NewIpcClient(server *IpcServer) *IpcClient {
	c := server.Connect() //连接服务端，返回给c，并给客户端结构体赋值
	return &IpcClient{c}
}

//函数call，是客户端ipcclient定义的client的方法，传入参数为string类型的method和params，返回参数为 *Response类型的resp和error
func (client *IpcClient) Call(method, params string) (resp *Response, err error) {
	req := &Request{method, params} //获取request请求结构体，参数为method，params

	var b []byte               //定义type类型数组b
	b, err = json.Marshal(req) //将req序列化为json格式串
	if err != nil {
		return
	}
	client.conn <- string(b) //将json串b传入通道client.conn
	str := <-client.conn     //获取client.conn通道字符串

	var resp1 Response                        //定义response类型结构体resp1
	err = json.Unmarshal([]byte(str), &resp1) //反序列化str为结构体resp1
	resp = &resp1                             //将resp1赋值给resp

	return
}

func (client *IpcClient) Close() { //ipc客户端client定义一个close函数
	client.conn <- "CLOSE" //将close字符串扔进client.conn通道
}
