package process

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/common/message"
	"go_code/chatroom/server/model"
	"go_code/chatroom/server/utils"
	"net"
)

//UserProcess 处理user相关业务
type UserProcess struct {
	Conn             net.Conn
	Userid           int //该Conn属于哪个用户id
	UserOnlineStatus int //该用户的在线状态
}

//NotifyOtherOnlineUser 通知其他全部在线用户,某个用户在线状态的改变
//注意,这里是传入的userID外的其他人
func (UserProcessThis *UserProcess) NotifyOtherOnlineUser(userID int, userStatus int) {

	if userStatus == message.UserOnline { //说明是通知上线的
		for id, up := range MyUserMgr.onlineUsers {
			//因为是上线通知,所以不会通知自己
			if id == userID {
				continue
			}
			up.NotifyMeChangeUserStatus(userID, message.UserOnline)
		}
	} else { //是通知其他状态的
		for id, up := range MyUserMgr.onlineUsers {
			//这里也不打算让客户端存自己的在线状态
			if id == userID {
				continue
			}
			up.NotifyMeChangeUserStatus(userID, userStatus)
		}
	}
}

//NotifyMeChangeUserStatus 根据传入的id,通知我(我来自于*UserProcess),该用户的在线状态变化
func (UserProcessThis *UserProcess) NotifyMeChangeUserStatus(userID int, userStatus int) {
	//组装NotifyUserStatusMes
	var mes message.Message
	var notifyUserStatusMes message.NotifyUserStatusMes

	switch userStatus {
	case message.UserOnline:
		mes.Type = message.NotifyUserStatusMesType
		notifyUserStatusMes.UserStatus = message.UserOnline
	default:
		mes.Type = message.NotifyUserOtherStatusMesType
		notifyUserStatusMes.UserStatus = userStatus
	}

	notifyUserStatusMes.UserID = userID

	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//将序列化后的notifyUserStatusMes赋值给 mes.Data
	mes.Data = string(data)

	//将mes序列化后准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//发送信息
	tf := &utils.Transfer{
		Conn: UserProcessThis.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline WritePkg err=", err)
		return
	}
}

//ServerProcessRegister 处理注册请求
func (UserProcessThis *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	//申明消息类型为注册消息类型
	var resMes message.Message
	resMes.Type = message.RegisterMesType

	//申明一个RegisterResMes
	var registerResMes message.RegisterResMes

	err = model.MyUserDao.Register(&registerMes.User)

	if err != nil {
		if err == model.ErrorUserExists {
			registerResMes.Code = 505
			registerResMes.Error = model.ErrorUserExists.Error()
		} else {
			registerResMes.Code = 505
			registerResMes.Error = "注册时产生了未知错误"
		}
	} else {
		registerResMes.Code = 200
		registerResMes.Error = ""
	}

	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	resMes.Data = string(data)

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	tf := &utils.Transfer{
		Conn: UserProcessThis.Conn,
	}
	err = tf.WritePkg(data)
	return

}

//ServerProcessLogin 处理登录请求
func (UserProcessThis *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	//申明消息类型为登录消息类型
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	//申明一个LoginResMes
	var loginResMes message.LoginResMes

	/*
		//暂时写死，后面改写成查数据库
		if loginMes.UserId == 1 && loginMes.UserPwd == "123" {
			loginResMes.Code = 200
			loginResMes.Error = ""
		} else {
			loginResMes.Code = 500
			loginResMes.Error = "用户不存在"
		}
	*/

	//去redis数据库完成账户验证
	_, err = model.MyUserDao.Login(loginMes.UserID, loginMes.UserPwd)
	if err != nil {

		if err == model.ErrorUserNotExists {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ErrorUserPwd {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器端未知错误"
		}

	} else {

		loginResMes.Code = 200
		loginResMes.Error = ""

		//登录成功后,将客户端发送过来的UserID,放入到UserProcess中,用于后面的AddOnlineUser()
		UserProcessThis.Userid = loginMes.UserID

		//登录成功后,设置用户在线状态为在线
		UserProcessThis.UserOnlineStatus = 0

		//登录成功后将该用户放入UserMgr中,记做在线用户
		MyUserMgr.AddOnlineUser(UserProcessThis)

		//在登录返回消息中,开始放入在线用户的id切片,用于客户端显示在线列表
		for id := range MyUserMgr.onlineUsers {
			loginResMes.UsersID = append(loginResMes.UsersID, id)
		}

		//登录后的返回信息,放入登录成功的该用户的id,用于客户端显示
		loginResMes.UserID = loginMes.UserID

		//通知其他所有人,我上线了
		UserProcessThis.NotifyOtherOnlineUser(loginMes.UserID, message.UserOnline)

		fmt.Println("成功登录")
	}

	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	resMes.Data = string(data)

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	tf := &utils.Transfer{
		Conn: UserProcessThis.Conn,
	}
	err = tf.WritePkg(data)
	return
}

//ChangeUserOnlineStatus 服务器端改变用户在线状态
func (UserProcessThis *UserProcess) ChangeUserOnlineStatus(mes *message.Message) (err error) {

	var changeMes message.ChangeUserOnlineStatusMes
	err = json.Unmarshal([]byte(mes.Data), &changeMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	userID := changeMes.User.UserID

	userStatus := changeMes.User.UserStatus

	//服务器端的返回值
	var changeResMes message.Code

	OnlienUser, ok := MyUserMgr.onlineUsers[userID]

	if !ok {
		//这个在线用户不存在,无法改变用户的在线状态
		changeResMes.Code = 404
		changeResMes.Error = model.ErrorOnlineUserNotExists.Error()
		return
	}

	OnlienUser.UserOnlineStatus = userStatus

	if userStatus == 1 {
		//说明用户想退出
		changeResMes.Code = 200
		changeResMes.Error = ""
		//delete(MyUserMgr.onlineUsers, userID)
		MyUserMgr.DelOnlineUser(userID)
	}

	//通知其他客户端用户,该客户端用户的在线状态改变
	UserProcessThis.NotifyOtherOnlineUser(userID, userStatus)

	//已经改变用户的在线状态
	var resMes message.Message
	resMes.Type = message.ChangeUserOnlineStatusResMesType

	//通知客户端
	data, err := json.Marshal(changeResMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err.Error())
		return
	}
	resMes.Data = string(data)

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err.Error())
		return
	}

	tf := &utils.Transfer{
		Conn: UserProcessThis.Conn,
	}
	err = tf.WritePkg(data)
	return
}
