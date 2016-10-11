//CenterClient匿名组合了IpcClient
package cg

import (
	"encoding/json"
	"errors"

	"ipc"
)

type CenterClient struct { //中央服务器客户端
	*ipc.IpcClient
}

//中央服务器添加新用户
func (client *CenterClient) AddPlayer(player *Player) error {
	b, err := json.Marshal(*player)
	if err != nil {
		return err
	}
	//调用ipc客户端的call方法
	resp, err := client.Call("addplayer", string(b))
	if err == nil && resp.Code == "200" {
		return nil
	}

	return err
}

//中央服务器删除用户
func (client *CenterClient) RemovePlayer(name string) error {
	ret, _ := client.Call("removeplayer", name)
	if ret.Code == "200" {
		return nil
	}

	return errors.New(ret.Code)
}

//中央服务器展示用户列表
func (client *CenterClient) ListPlayer(params string) (ps []*Player, err error) {
	resp, _ := client.Call("listplayer", params)
	if resp.Code != "200" {
		err = errors.New(resp.Code)
		return
	}

	err = json.Unmarshal([]byte(resp.Body), &ps)
	return
}

//中央服务器广播消息
func (client *CenterClient) Broadcast(message string) error {
	m := &Message{Content: message}

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	resp, _ := client.Call("broadcast", string(b))
	if resp.Code == "200" {
		return nil
	}

	return errors.New(resp.Code)
}
