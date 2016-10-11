//主程序具备模拟用户游戏过程和管理员功能（发通告）
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"cg"
	"ipc"
)

var centerClient *cg.CenterClient //定义中央服务器

//启动中央服务器，通过ipc去创建服务器和客户端
func startCenterService() error {
	server := ipc.NewIpcServer(&cg.CenterServer{}) //服务器
	client := ipc.NewIpcClient(server)             //客户端
	centerClient = &cg.CenterClient{client}        //服务器指令

	return nil
}

//帮助指令
func Help(args []string) int {
	fmt.Println(`
    Commands:
        login <userbane><level><exp>
        logout <username>
        send <message>
        listplayer
        quit(q)
        help(h)
    `)
	return 0
}

//退出程序
func Quit(args []string) int {
	return 1
}

//退出“中心服务器”
func Logout(args []string) int {
	//判断是否为2个参数
	if len(args) != 2 {
		fmt.Println("USAGE: logout <username>")
		return 0
	}
	//中心服务器处理退出玩家
	centerClient.RemovePlayer(args[1])

	return 0
}

//登录”中心服务器“
func Login(args []string) int {
	//4个参数
	if len(args) != 4 {
		fmt.Println("USAGE: login <username><level><exp>")
		return 0
	}
	//第3个参数验证
	level, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Invaild Parameter : <level> should be an integer.")
		return 0
	}
	//第4个参数验证
	exp, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Println("Invaild Parameter : <exp> should be an integer.")
		return 0
	}
	//创建新玩家，并将传入的参数赋值
	player := cg.NewPlayer()
	player.Name = args[1]
	player.Level = level
	player.Exp = exp
	//中心服务器添加玩家
	err = centerClient.AddPlayer(player)
	if err != nil {
		fmt.Println("Faild adding player", err)
	}

	return 0
}

//展示所有玩家
func ListPlayer(args []string) int {
	ps, err := centerClient.ListPlayer("")
	if err != nil {
		fmt.Println("Faild. ", err)
	} else {
		for i, v := range ps {
			fmt.Println(i+1, ":", v)
		}
	}
	return 0
}

//发送消息给在线玩家
func Send(args []string) int {
	message := strings.Join(args[1:], " ")
	//服务器广播消息
	err := centerClient.Broadcast(message)
	if err != nil {
		fmt.Println("Faild. ", err)
	}

	return 0
}

//获取操作命令map
func GetCommandHandlers() map[string]func(args []string) int {
	return map[string]func([]string) int{
		"help":       Help,
		"h":          Help,
		"quit":       Quit,
		"q":          Quit,
		"login":      Login,
		"logout":     Logout,
		"listplayer": ListPlayer,
		"send":       Send,
	}
}

func main() {
	//提示文字
	fmt.Println("Casual Game Server Soluion")

	//开始启动服务器
	startCenterService()

	//列出帮助指令
	Help(nil)

	//定义一个buffer
	r := bufio.NewReader(os.Stdin)

	//获取操作指令map
	handlers := GetCommandHandlers()

	for {
		fmt.Print("command> ")
		//循环读取每行指令
		b, _, _ := r.ReadLine()
		line := string(b)
		//分离参数为数组
		tokens := strings.Split(line, " ")
		//判断是否存在符合要求的handler指令
		if handler, ok := handlers[tokens[0]]; ok {
			ret := handler(tokens)
			if ret != 0 {
				break
			}
		} else {
			fmt.Println("Unknown command : ", tokens[0])
		}
	}
}
