package model

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/common/message"
	"time"

	"github.com/garyburd/redigo/redis"
)

//当服务端启动后,初始化一个SmsDao实例
//将其设置为全局变量,在需要去redis操作时,就直接使用
var (
	MySmsDao *SmsDao
)

//SmsDao 结构体
type SmsDao struct {
	pool *redis.Pool
}

//NewSmsDao 使用工厂模式,创建一个SmsDao实例
func NewSmsDao(pool *redis.Pool) (smsDao *SmsDao) {
	smsDao = &SmsDao{
		pool: pool,
	}
	return
}

//SaveToAllMes 云端保存@全体成员的消息
func (smsDaoThis *SmsDao) SaveToAllMes(content *message.SmsMes) (err error) {

	conn := smsDaoThis.pool.Get()
	defer conn.Close()

	var contents message.Contents
	contents.FromUser.UserID = content.User.UserID
	contents.Time = time.Now().Unix()
	contents.Content = content.Content

	data, err := json.Marshal(contents)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//序列化完成之后,在redis进行Hset
	_, err = conn.Do("RPush", "toallmes", string(data))
	if err != nil {
		fmt.Println("写入离线@all消息,写入list失败,err=", err)
	}
	return
}

//GetToAllMes 读取云端保存@全体成员 的消息
func (smsDaoThis *SmsDao) GetToAllMes() (ToClientMesMap map[int]string, err error) {
	conn := smsDaoThis.pool.Get()
	defer conn.Close()

	res, err := redis.Values(conn.Do("lrange", "toallmes", 0, -1))
	if err != nil {
		if err == redis.ErrNil {
			//在该list中，并没有@全体成员消息
			err = ErrorNothingUnreadMesToAll
		}
		fmt.Println("查询@全体消息是否存在时,err=", err)
		return
	}

	ToClientMesMap = make(map[int]string)

	for key, value := range res {
		ToClientMesMap[key] = string(value.([]byte))
	}

	return ToClientMesMap, nil

}
