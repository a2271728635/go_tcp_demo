package model

import (
	"go_code/chatroom/common/message"
	"net"
)

//CurUser 在客户端,因为很多地方需要使用到curUser
//将其作为一个全局变量,为了方便管理
//全局变量放到client->process->userMgr.go->var CurUser
type CurUser struct {
	Conn net.Conn
	message.User
}
