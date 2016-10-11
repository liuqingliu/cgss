//添加用户、删除用户、列出用户和广播
package cg

import (
	"encoding/json"
	"errors"
	"ipc"
	"sync" //引入”锁“的包
)

var _ ipc.Server = &CenterServer{} //确认实现了server接口

type Message struct { //message消息内容结构体
	From    string "from"
	To      string "to"
	Content string "content"
}

type Room struct { //房间结构体
}

type CenterServer struct { //中央处理服务结构体内容
	servers map[string]ipc.Server //类型为map[string]ipc.Server的map参数
	players []*Player             //类型为player的参数
	rooms   []*Room               //类型为room的参数
	mutex   sync.RWMutex          //类型为sync.RWMutex的读写锁
}

func NewCenterServer() *CenterServer { //创建返回值为centerserver类型指针的函数
	servers := make(map[string]ipc.Server)                   //创建实体servers，map类型
	players := make([]*Player, 0)                            //创建player类型的结构体
	return &CenterServer{servers: servers, players: players} //返回中央处理服务
}

func (server *CenterServer) addPlayer(params string) error { //添加player到中心服务器server的函数，传入参数params，返回值error
	player := NewPlayer() //创建一个player

	err := json.Unmarshal([]byte(params), &player) //反序列化json格式的player
	if err != nil {
		return err
	}

	server.mutex.Lock() //锁
	//解锁 ？为什么锁了后直接又解锁。。。是可以同时到锁这里，只有当解锁后才会添加下面的player！
	defer server.mutex.Unlock()
	//直接添加到登录到服务器上的players
	server.players = append(server.players, player)

	return nil
}

//在中心服务器server中删除player(下线)
func (server *CenterServer) removePlayer(params string) error {
	//锁住，只能单请求
	server.mutex.Lock()
	defer server.mutex.Unlock()
	//循环当前服务器上在线的用户
	for i, v := range server.players {
		//判断在线用户中和要删除的用户名是否相同
		if v.Name == params {
			if len(server.players) == 1 {
				//若当前在线服务器上用户只有一个，则直接清空
				server.players = make([]*Player, 0)
			} else if i == len(server.players)-1 {
				//若当前找到的是最后一个位置的player，则取所有前面的players
				server.players = server.players[:i-1]
			} else if i == 0 {
				//若当前找到的是第一个位置的player，则取所有后面的players
				server.players = server.players[1:]
			} else {
				//若当前找到的是中间位置的player，则取所有前面，和所有后面，除去当前i位置
				server.players = append(server.players[:i-1], server.players[:i+1]...)
			}
			return nil
		}
	}

	return errors.New("Player not found")
}

//在线服务器展示所有player玩家信息
func (server *CenterServer) listPlayer(params string) (players string, err error) {
	server.mutex.RLock()
	defer server.mutex.RUnlock()

	if len(server.players) > 0 {
		//json化players
		b, _ := json.Marshal(server.players)
		players = string(b)
	} else {
		err = errors.New("No play online.")
	}

	return
}

//在线服务器广播消息
func (server *CenterServer) broadCast(params string) error {
	var message Message
	//将传入的消息反序列化
	err := json.Unmarshal([]byte(params), &message)
	if err != nil {
		return err
	}

	server.mutex.Lock()
	defer server.mutex.Unlock()

	if len(server.players) > 0 {
		for _, player := range server.players {
			//将广播的消息传入给每个player
			player.mq <- &message
		}
	} else {
		err = errors.New("No player online.")
	}

	return err
}

//在线服务器指令处理方法
func (server *CenterServer) Handle(method, params string) *ipc.Response {
	switch method {
	//添加在线玩家
	case "addplayer":
		err := server.addPlayer(params)
		if err != nil {
			return &ipc.Response{Code: err.Error()}
		}
	//移除在线玩家
	case "removeplayer":
		err := server.removePlayer(params)
		if err != nil {
			return &ipc.Response{Code: err.Error()}
		}
	//显示所有玩家
	case "listplayer":
		players, err := server.listPlayer(params)
		if err != nil {
			return &ipc.Response{Code: err.Error()}
		}
		return &ipc.Response{"200", players}
	//广播消息
	case "broadcast":
		err := server.broadCast(params)
		if err != nil {
			return &ipc.Response{Code: err.Error()}
		}
		return &ipc.Response{Code: "200"}
	//错误指令处理
	default:
		return &ipc.Response{Code: "404", Body: method + ":" + params}
	}
	return &ipc.Response{Code: "200"}
}

//返回所需要的中心服务器名称
func (server *CenterServer) Name() string {
	return "CenterServer"
}
