package process

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/client/utils"
	"go_code/chatroom/common/message"
	"net"
	"os"
)

//UserProcess 是客户端处理用户相关的struct
type UserProcess struct {
}

//Islogin 登录判断
func (thisUserProcess *UserProcess) Islogin(userID int, userPwd string) (err error) {
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	//准备延时关闭
	defer conn.Close()

	//准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType

	//用户登录信息
	var loginMes message.LoginMes
	loginMes.UserID = userID
	loginMes.UserPwd = userPwd

	//序列化用户信息
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//将序列化后的值给mes.Data
	mes.Data = string(data)

	//序列化mes本身
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	/*
	   //这里已经被优化到utils.go文件的WritePkg()工具函数
	   	var pkgLen uint32
	   	pkgLen = uint32(len(data))
	   	var buf [4]byte
	   	binary.BigEndian.PutUint32(buf[:], pkgLen)
	   	//发送信息长度
	   	n, err := conn.Write(buf[:4])
	   	if n != 4 || err != nil {
	   		fmt.Println("conn.Write(buf) err=", err)
	   		return
	   	}

	   	fmt.Printf("客户端发送消息长度成功,长度为=%d 内容是=%v", len(data), string(data))

	   	//发送消息本身
	   	_, err = conn.Write(data)
	   	if err != nil {
	   		fmt.Println("conn.Write(data) err=", err)
	   		return
	   	}
	*/

	tf := &utils.Transfer{
		Conn: conn,
	}

	//准备向服务器端发送消息
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("WritePkg(data) err=", err)
	}

	//处理服务器端发过来的消息
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn)  err=", err)
		return
	}

	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)

	if loginResMes.Code == 200 {

		//初始化CurUser
		CurUser.Conn = conn
		CurUser.UserID = userID
		CurUser.UserStatus = message.UserOnline

		fmt.Println("客户端:登录成功")
		//显示当前聊天室在线用户
		fmt.Println("当前在线人数:", len(loginResMes.UsersID))
		fmt.Println("当前聊天室在线用户:")
		for _, v := range loginResMes.UsersID {
			//注意,这里的UsersID是登录后,服务器端返回的,在线用户的 []int

			//这里设置聊天室在线人数不显示自己
			if v == userID {
				continue
			}

			fmt.Println("用户id为:", v)

			//对客户端的 onlineUser 进行初始化
			//将从服务器端获取到的 UsersID []int 装入客户端 onlineUser map[int]*User 中
			user := &message.User{
				UserID:     v,
				UserStatus: message.UserOnline,
			}
			onlineUser[v] = user

		}
		fmt.Printf("---恭喜id为 %d 用户登录成功--- \n", loginResMes.UserID)

		//开启协程,保持客户端和服务器端的通讯不断开
		go serverProcessMes(conn)

		for {
			ShowMenu() //登录成功后显示二级菜单
		}

	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}

//Register 注册判断
func (thisUserProcess *UserProcess) Register(userID int, userPwd string, userName string) (err error) {

	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	//准备延时关闭
	defer conn.Close()

	//准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.RegisterMesType

	//用户注册信息
	var registerMes message.RegisterMes
	registerMes.User.UserID = userID
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	//序列化用户信息
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//将序列化后的值给mes.Data
	mes.Data = string(data)

	//序列化mes本身
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	tf := &utils.Transfer{
		Conn: conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册时,客户端发送到服务器端出错:", err)
	}

	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn)  err=", err)
		return
	}

	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)

	if registerResMes.Code == 200 {
		fmt.Println("注册成功,请重新登录")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}

	return
}

//ChangeUserOnlineStatus 客户端改变在线状态
func (thisUserProcess *UserProcess) ChangeUserOnlineStatus(status int) (err error) {
	var mes message.Message
	mes.Type = message.ChangeUserOnlineStatusMesType

	var change message.ChangeUserOnlineStatusMes
	change.User.UserID = CurUser.UserID
	//change.User.UserPwd = CurUser.UserPwd
	change.User.UserStatus = status

	//开始序列化
	data, err := json.Marshal(change)
	if err != nil {
		fmt.Println("ChangeUserStatus() json.Marshal fail=", err.Error())
		return
	}

	mes.Data = string(data)

	//序列化mes后,准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("ChangeUserStatus() json.Marshal fail=", err.Error())
		return
	}

	//开始发送
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMes() tf.WritePkg(data) fail=", err.Error())
		return
	}
	return
}
