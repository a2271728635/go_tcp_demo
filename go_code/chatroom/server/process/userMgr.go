package process

import (
	"fmt"
	"go_code/chatroom/common/message"
)

//服务器端的在线列表有且只有一个
//所以UserMgr 到处都会用到,设置为全局变量
var (
	MyUserMgr *UserMgr
)

//UserMgr 用于维护用户在线列表
type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

//对 onlineUsers这个map 完成初始化
func init() {
	MyUserMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

//AddOnlineUser 登录之后,将在线人数添加进onlineUsers
func (UserMgrThis *UserMgr) AddOnlineUser(up *UserProcess) {
	UserMgrThis.onlineUsers[up.Userid] = up
}

//DelOnlineUser 退出登录之后,将退出的删除
func (UserMgrThis *UserMgr) DelOnlineUser(userID int) {
	delete(UserMgrThis.onlineUsers, userID)
}

//GetAllOnlineUser 返回当前所有的在线用户
func (UserMgrThis *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return UserMgrThis.onlineUsers
}

//GetOnlineUserByID 根据传入id返回对应的值
func (UserMgrThis *UserMgr) GetOnlineUserByID(userID int) (up *UserProcess, err error) {

	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserID = userID
	notifyUserStatusMes.UserStatus = message.UserOnline

	up, isOk := UserMgrThis.onlineUsers[userID]
	if isOk {

	} else {
		//查不到和服务器端的连接,说明不在线
		err = fmt.Errorf("用户%d 不在线", userID)
		return
	}
	return
}
