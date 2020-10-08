package process

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/client/utils"
	"go_code/chatroom/common/message"
)

//SmsProcess 用于发送消息
type SmsProcess struct {
}

//SendGroupMes 发送消息到全部人
func (SmsProcessThis *SmsProcess) SendGroupMes(content string) (err error) {

	//开始组装Mes
	var mes message.Message
	mes.Type = message.SmsMesType

	//开始组装SmsMes
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.User.UserID = CurUser.UserID
	smsMes.User.UserStatus = CurUser.UserStatus

	//开始序列化
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes() json.Marshal fail=", err.Error())
		return
	}

	mes.Data = string(data)

	//序列化mes后,准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes() json.Marshal fail=", err.Error())
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

//SendSmsToUser 发送私聊信息
func (SmsProcessThis *SmsProcess) SendSmsToUser(content string, ReachID int) (err error) {
	//开始组装Mes
	var mes message.Message
	mes.Type = message.SmsPToPMesType

	//开始组装SmsPToPMes
	var SmsPToPMes message.SmsPToPMes
	SmsPToPMes.Content = content
	SmsPToPMes.FromUser.UserID = CurUser.UserID
	SmsPToPMes.ReachUser.UserID = ReachID
	//开始序列化
	data, err := json.Marshal(SmsPToPMes)
	if err != nil {
		fmt.Println("SendSmsToUser() json.Marshal fail=", err.Error())
		return
	}

	mes.Data = string(data)

	//序列化mes后,准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendSmsToUser() json.Marshal fail=", err.Error())
		return
	}

	//开始发送
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendSmsToUser() tf.WritePkg(data) fail=", err.Error())
		return
	}
	return
}

//GetToAllMesForServer 发送,想从服务器端获取@all的消息
func (SmsProcessThis *SmsProcess) GetToAllMesForServer() (err error) {
	//开始组装Mes
	var mes message.Message
	mes.Type = message.GetToAllMesType

	//序列化mes后,准备发送
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("GetToAllMesForServer() json.Marshal fail=", err.Error())
		return
	}

	//开始发送
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("GetToAllMesForServer() tf.WritePkg(data) fail=", err.Error())
		return
	}
	return

}
