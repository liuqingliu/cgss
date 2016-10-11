//在线玩家的管理
//我们为每个玩家都起了一个独立的goroutine，监听所有发送给他们
//的聊天信息，一旦收到就即时打印到控制台上
package cg

import (
	"fmt"
)

type Player struct { //定义player结构体属性内容
	Name  string        "name"  //姓名
	Level int           "level" //等级
	Exp   int           "exp"   //
	Room  int           "room"  //房间号
	mq    chan *Message //等待收取的消息
}

func NewPlayer() *Player { //函数NewPlayer创建player
	m := make(chan *Message, 1024)    //创建一个缓冲区为1024的channel，类型为message
	player := &Player{"", 0, 0, 0, m} //创建一个player，赋初值

	go func(p *Player) { // gorountine 匿名函数，处理实体player，参数p为player类
		for {
			msg := <-p.mq                                         //获取player类channel里的p实体mq值
			fmt.Println(p.Name, "received message:", msg.Content) //打印收到message
		}
	}(player)

	return player
}
