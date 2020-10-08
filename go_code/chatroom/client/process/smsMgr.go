package process

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/common/message"
	"time"
)

func outputGroupMes(mes *message.Message) {

	//打印服务器发送到客户端的群发消息
	var smsMesRes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMesRes)
	if err != nil {
		fmt.Println("outputGroupMes() json.Unmarshal fail=", err.Error())
		return
	}
	//显示
	sms := fmt.Sprintf("用户\t%d @全体成员:\t%s", smsMesRes.User.UserID, smsMesRes.Content)
	fmt.Println(sms)
	fmt.Println()
}

func outputPToPMes(mes *message.Message) {
	//打印服务器发送到客户端的私聊消息
	var smsPToPMes message.SmsPToPMes
	err := json.Unmarshal([]byte(mes.Data), &smsPToPMes)
	if err != nil {
		fmt.Println("outputPToPMes() json.Unmarshal fail=", err.Error())
		return
	}
	//显示
	sms := fmt.Sprintf("用户\t%d 对你说:\t%s", smsPToPMes.FromUser.UserID, smsPToPMes.Content)
	fmt.Println(sms)
	fmt.Println()
}

func getToAllMes(mesMap *map[int]string) {
	//该函数处理返回的 @all 的消息map
	
	var value message.Contents
	for _, v := range *mesMap {
		err := json.Unmarshal([]byte(v), &value)
		if err != nil {
			fmt.Println("json.Marshal err=", err)
			return
		}
		formatTimeStr:=time.Unix(value.Time,0).Format("2006-01-02 15:04:05")
		fmt.Printf("%v 在 %v @全体成员 : %v",value.FromUser.UserID,formatTimeStr,value.Content)
		fmt.Println()
	}
}
