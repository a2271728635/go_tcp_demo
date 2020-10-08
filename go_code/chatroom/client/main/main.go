package main

import (
	"fmt"
	"go_code/chatroom/client/process"
)

var userID int
var userPwd string
var userName string

func main() {
	//用于接收用户的选择
	var key int
	//用于判断是否还要继续显示菜单
	//var loop = true

	for true {
		fmt.Println("---欢迎登陆 By LiaoYi---")
		fmt.Println("1.登陆")
		fmt.Println("2.注册")
		fmt.Println("3.退出")
		fmt.Println("请选择1-3")

		fmt.Scanf("%d \n", &key)
		switch key {
		case 1:
			fmt.Println("登陆LY群")
			fmt.Println("输入id(仅数字)")
			fmt.Scanf("%d\n", &userID)

			fmt.Println("输入密码")
			fmt.Scanf("%s\n", &userPwd)

			fmt.Println("userId=", userID)
			fmt.Println("userPwd=", userPwd)

			up := &process.UserProcess{}
			up.Islogin(userID, userPwd)
			//loop = false
		case 2:
			fmt.Println("注册用户")

			fmt.Println("请输入注册用户id(仅数字)")
			fmt.Scanf("%d\n", &userID)

			fmt.Println("请输入注册用户密码")
			fmt.Scanf("%s\n", &userPwd)

			fmt.Println("请输入注册用户昵称")
			fmt.Scanf("%s\n", &userName)

			fmt.Println("userId=", userID)
			fmt.Println("userPwd=", userPwd)
			fmt.Println("userPwd=", userName)
			up := &process.UserProcess{}
			up.Register(userID, userPwd, userName)
			//loop = false
		case 3:
			fmt.Println("退出")
			break
			//loop = false
		default:
			fmt.Println("无效输入，请输入正确的值")
			//loop = true
		}
	}
}
