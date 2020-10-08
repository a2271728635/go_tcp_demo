package process

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/common/message"
	"go_code/chatroom/server/model"
	"go_code/chatroom/server/utils"
	"net"
)

//MySmsProcess 用于全局
var (
	MySmsProcess *SmsProcess
)

//SmsProcess 处理和消息有关
type SmsProcess struct {
	Conn net.Conn //是谁在调用这个结构体
}

//SendGroupMes 处理用户群发消息
func (SmsProcessThis *SmsProcess) SendGroupMes(mes *message.Message) (err error) {

	//取出mes的内容
	var smsMes message.SmsMes
	err = json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("SendGroupMes() json.Unmarshal fail=", err.Error())
		return
	}

	//组装到@全体成员的离线消息中
	SmsProcessThis.SaveToAllMes(smsMes.Content, smsMes.User)
	//------------------------------------------------------------------

	//组装smsMesResMes
	var smsMesResMes message.Message
	smsMesResMes.Type = message.SmsMesResType

	//组装smsMesResMes中的内容
	var returnData message.SmsMesRes
	returnData.Content = smsMes.Content
	returnData.User = smsMes.User

	data, err := json.Marshal(returnData)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//将序列化后的值赋值给smsMesResmes.Data
	smsMesResMes.Data = string(data)

	//将mes序列化后准备发送
	data, err = json.Marshal(smsMesResMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//------------------------------------------------------------------

	//拿到服务器端的在线用户列表,将消息群发给在线用户
	for id, up := range MyUserMgr.onlineUsers {
		//因为在线用户包括自己,所有不发送给自己
		if id == smsMes.User.UserID {
			continue
		}
		SmsProcessThis.SendToAllOnlineUser(data, up.Conn)

	}
	return
}

//SendToAllOnlineUser 向所有在线用户发送消息
func (SmsProcessThis *SmsProcess) SendToAllOnlineUser(data []byte, conn net.Conn) (err error) {
	tf := &utils.Transfer{
		Conn: conn,
	}
	err = tf.WritePkg(data)
	//fmt.Println("服务器端已发送")

	if err != nil {
		fmt.Println("SendToAllOnlineUser fail=", err.Error())
		return
	}
	return
}

//SendSmsToUser 处理在线用户私聊
func (SmsProcessThis *SmsProcess) SendSmsToUser(mes *message.Message) (err error) {

	//取出mes的内容
	var smsPToPMes message.SmsPToPMes
	err = json.Unmarshal([]byte(mes.Data), &smsPToPMes)
	if err != nil {
		fmt.Println("SendSmsToUser() json.Unmarshal fail=", err.Error())
		return
	}

	//组装sms
	var sms message.Message
	sms.Type = message.SmsPToPMesResType

	//组装sms中的内容
	var returnData message.SmsPToPMes
	returnData.Content = smsPToPMes.Content
	returnData.FromUser = smsPToPMes.FromUser
	returnData.ReachUser = smsPToPMes.ReachUser

	data, err := json.Marshal(returnData)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//将序列化后的值赋值给sms.Data
	sms.Data = string(data)

	//将sms序列化后准备发送
	data, err = json.Marshal(sms)
	if err != nil {
		fmt.Println("json.Marshal err=", err.Error())
		return
	}
	//找到对方的Conn
	ReachProcess, ok := MyUserMgr.onlineUsers[smsPToPMes.ReachUser.UserID]
	if !ok {
		fmt.Println("对方已离线,转入离线留言逻辑")
		return
	}

	tf := &utils.Transfer{
		Conn: ReachProcess.Conn,
	}
	err = tf.WritePkg(data)
	//fmt.Println("服务器端已发送")

	if err != nil {
		fmt.Println("SendToAllOnlineUser fail=", err.Error())
		return
	}
	return
}

//GetToAllMes 获取@全体用户的云端保存的消息记录,并返回给客户端
func (SmsProcessThis *SmsProcess) GetToAllMes(mes *message.Message) (err error) {

	toClientMap := make(map[int]string)
	toClientMap, err = model.MySmsDao.GetToAllMes()
	if err != nil {
		fmt.Println("GetToAllMes() err=", err.Error())
		return
	}

	//组装sms
	var sms message.Message
	sms.Type = message.GetToAllResMesType

	//序列化
	data, err := json.Marshal(toClientMap)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//将序列化后的值赋值给sms.Data
	sms.Data = string(data)

	//将sms序列化后准备发送
	data, err = json.Marshal(sms)
	if err != nil {
		fmt.Println("json.Marshal err=", err.Error())
		return
	}

	tf := &utils.Transfer{
		Conn: SmsProcessThis.Conn,
	}
	err = tf.WritePkg(data)

	if err != nil {
		fmt.Println("GetToAllMes fail=", err.Error())
		return
	}
	return
}

//SaveToAllMes 将@全体成员的消息保存到云端
func (SmsProcessThis *SmsProcess) SaveToAllMes(content string, user message.User) (err error) {
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.User = user
	err = model.MySmsDao.SaveToAllMes(&smsMes)
	if err != nil {
		fmt.Println("SaveToAllMes() err=", err.Error())
		return
	}
	return
}
