package process

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/client/model"
	"go_code/chatroom/common/message"
)

//onlineUser 客户端需要维护的所有用户在线状态的map
var onlineUser map[int]*message.User = make(map[int]*message.User, 10)

//CurUser 此客户端维护的和服务器端通讯的Conn,在登录成功后进行初始化
var CurUser model.CurUser

//outputOnlineUser 在客户端显示当前在线的用户
func outputOnlineUser() {
	fmt.Println("当前所有在线用户:")
	for id := range onlineUser {
		fmt.Println("用户id:", id)
	}
}

//outputUpdateUserStatus 按照聊天室的逻辑,登录才提醒,改变在线状态不提醒
//输出用户在线在线状态的变化
func outputUpdateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	fmt.Println("有新用户登录,id:", notifyUserStatusMes.UserID)
}

//该函数处理返回的NotifyUserStatusMes
//更新客户端的用户map
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	//为了防止重复赋值,即上线之后下线的情况
	user, ok := onlineUser[notifyUserStatusMes.UserID]
	if ok { //如果存在,说明本地map有该用户的在线状态信息
		user.UserStatus = notifyUserStatusMes.UserStatus
	} else { //如果不存在,说明没有,进行创建赋值
		user = &message.User{
			UserID:     notifyUserStatusMes.UserID,
			UserStatus: notifyUserStatusMes.UserStatus,
		}
	}
	onlineUser[notifyUserStatusMes.UserID] = user
}

//该函数处理返回的ChangeUserOnlineStatusResMes
//打印改变自身在线状态的结果
func changeUserStatus(mes *message.Message) {

	var changeUserStatusMes message.Code
	err := json.Unmarshal([]byte(mes.Data), &changeUserStatusMes)
	if err != nil {
		fmt.Println("changeUserStatus() json.Unmarshal fail=", err.Error())
		return
	}
	switch changeUserStatusMes.Code {
	case 200:
		fmt.Println("在线状态修改成功")
	case 404:
		fmt.Println("在线状态修改失败")
		fmt.Println("原因:", changeUserStatusMes.Error)
	}
}


