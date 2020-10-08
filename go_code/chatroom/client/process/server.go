package process

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/client/utils"
	"go_code/chatroom/common/message"
	"net"
	"os"
)

//ShowOnlieMenu 展示在线状态选择菜单
func ShowOnlieMenu(key int) {

}

//ShowMenu 展示登录成功后的菜单
func ShowMenu() {

	var key int //获取用户输入的菜单内容

	loop := true        //是否显示菜单
	var content string  //发送消息的内容
	var reachUserID int //reachUserID 要私聊的用户ID

	//创建SmsProcess实例,方便反复使用
	SmsProcess := &SmsProcess{}
	Up := &UserProcess{}
	for loop {

		fmt.Println("---1.显示在线用户列表---")
		fmt.Println("---2.发送全体消息---")
		fmt.Println("---3.发送私聊消息---")
		fmt.Println("---4.消息列表---")
		fmt.Println("---5.变更登录状态---")
		fmt.Println("---请输入(1-5)---")

		fmt.Scanf("%d \n", &key)

		switch key {
		case 1:
			fmt.Println("显示在线用户列表")

			outputOnlineUser()
		case 2:
			fmt.Println("发送全体消息")

			fmt.Println("请输入发送消息:")
			fmt.Scanf("%s \n", &content)
			SmsProcess.SendGroupMes(content)
		case 3:
			fmt.Println("请输入要私聊的用户id(仅数字)")
			fmt.Scanf("%d\n", &reachUserID)
			fmt.Println("请输入私聊内容")
			fmt.Scanf("%s \n", &content)
			fmt.Printf("你对%d说:%v \n", reachUserID, content)
			SmsProcess.SendSmsToUser(content, reachUserID)
		case 4:
			fmt.Println("消息列表")
		MesList:
			for {
				fmt.Println("---1.@全体成员---")
				fmt.Println("---2.私聊---")
				fmt.Println("---3.返回上一级菜单---")
				fmt.Scanf("%d \n", &key)
				switch key {
				case 1:
					SmsProcess.GetToAllMesForServer()
					break MesList
				case 2:
					fmt.Println("服务器端还没做")
					break MesList
				case 3:
					break MesList
				}
			}

		case 5:
		Online:
			for {
				fmt.Println("---1.在线---")
				fmt.Println("---2.离线---")
				fmt.Println("---3.忙碌---")
				fmt.Println("---4.返回上一级菜单---")
				fmt.Scanf("%d \n", &key)
				switch key {
				case 1:
					Up.ChangeUserOnlineStatus(message.UserOnline)
					fmt.Println("当前在线状态已变更为:【在线】")
					break Online
				case 2:
					Up.ChangeUserOnlineStatus(message.UserOffline)
					fmt.Println("当前在线状态已变更为:【离线】")
					os.Exit(0)
					break Online
				case 3:
					Up.ChangeUserOnlineStatus(message.UserBusyStatus)
					fmt.Println("当前在线状态已变更为:【忙碌】")
					break Online
				case 4:
					break Online
				default:
					fmt.Println("输入信息错误,请重新输入")
				}
			}
		default:
			fmt.Println("输入信息错误,请重新输入")

		}
	}
}

//serverProcessMes 保持和服务器端的通讯
func serverProcessMes(conn net.Conn) {

	//创建Transfer实例，不停的读取服务器端发过来的消息
	tf := &utils.Transfer{
		Conn: conn,
	}

	for {
		fmt.Println("客户端正在等待读取服务器端发送的消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg err=", err)
			return
		}

		//fmt.Printf("mes=%v", mes)

		//等待服务器端发送消息
		switch mes.Type {
		case message.NotifyUserStatusMesType: //服务器端发送到客户端的消息,有人上线

			//拿到从服务器端发过来的,用户在线状态更新的消息
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)

			//将服务器端推送过来的NotifyUserStatusMes,在客户端的map进行保存和维护
			updateUserStatus(&notifyUserStatusMes)

			//打印上线消息
			outputUpdateUserStatus(&notifyUserStatusMes)
		case message.SmsMesResType: //服务器端发送过来了,某客户端的群发消息
			outputGroupMes(&mes)

		case message.SmsPToPMesResType:
			//处理服务器端发送过来的私聊信息
			outputPToPMes(&mes)
		case message.ChangeUserOnlineStatusResMesType:
			//处理服务器端发送过来的,关于本客户端的在线状态改变的返回值
			changeUserStatus(&mes)

		case message.NotifyUserOtherStatusMesType:
			//处理服务器端发送过来的,关于其他客户端的在线状态改变的返回值

			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)

			//将服务器端推送过来的NotifyUserStatusMes,在客户端的map进行保存和维护
			updateUserStatus(&notifyUserStatusMes)

		case message.GetToAllResMesType:
			//处理服务器发送过来的,@all的消息

			allClientMap := make(map[int]string)

			json.Unmarshal([]byte(mes.Data), &allClientMap)

			getToAllMes(&allClientMap)
		default:
			fmt.Println("mes.Type=", mes.Type)
			fmt.Println("mes.Data=", mes.Data)
			fmt.Printf("无法识别的消息类型,请更新客户端")
		}
	}
}
