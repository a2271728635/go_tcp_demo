package main

import (
	"fmt"
	"go_code/chatroom/common/message"
	process2 "go_code/chatroom/server/process"
	"go_code/chatroom/server/utils"
	"io"
	"net"
)

//Processor 是处理客户端发过来的消息类别，决定调用函数的控制台
type Processor struct {
	Conn net.Conn
}

//根据客户端发过来的消息类别，决定调用函数
//和main同包 小写
func (thisProcessor *Processor) serverProcessMes(mes *message.Message) (err error) {

	switch mes.Type {
	case message.LoginMesType:
		//处理客户端登录请求
		//创建一个UserProcess实例
		up := &process2.UserProcess{
			Conn: thisProcessor.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		//处理客户端注册请求
		up := &process2.UserProcess{
			Conn: thisProcessor.Conn,
		}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:
		//处理客户端群发请求
		smsp := &process2.SmsProcess{}
		err = smsp.SendGroupMes(mes)
	case message.SmsPToPMesType:
		//处理客户端私聊请求
		smsp := &process2.SmsProcess{}
		err = smsp.SendSmsToUser(mes)
	case message.ChangeUserOnlineStatusMesType:
		//处理客户端改变用户状态请求
		up := &process2.UserProcess{
			Conn: thisProcessor.Conn,
		}
		err = up.ChangeUserOnlineStatus(mes)
	case message.GetToAllMesType:
		//处理用户想获取@all消息请求
		smsp := &process2.SmsProcess{
			Conn: thisProcessor.Conn,
		}
		err = smsp.GetToAllMes(mes)
	default:
		fmt.Println("错误类别")
		return
	}
	return
}

func (thisProcessor *Processor) process2() (err error) {
	//读取客户端发过来的东西
	for {

		//创建一个tf实例 去工具类里面找read
		tf := &utils.Transfer{
			Conn: thisProcessor.Conn,
		}

		mes, err := tf.ReadPkg()

		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端已关闭，服务器端也关闭")
				return err
			}
			fmt.Println("readPkg err=", err)
			return err
		}

		err = thisProcessor.serverProcessMes(&mes)
		//fmt.Println("serverProcessMes err=", err)
		if err != nil {
			return err
		}
	}
}
