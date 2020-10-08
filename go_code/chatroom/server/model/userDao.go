package model

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/common/message"

	"github.com/garyburd/redigo/redis"
)

//当服务端启动后,初始化一个userDao实例
//将其设置为全局变量,在需要去redis操作时,就直接使用
var (
	MyUserDao *UserDao
)

//UserDao 结构体完成对User 结构体的各种操作
type UserDao struct {
	pool *redis.Pool
}

//NewUserDao 使用工厂模式,创建一个UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

//getUserById 通过传入的id判断用户是否存在,存在err == nil,不存在err = ErrorUserNotExists
func (userDaoThis *UserDao) getUserByID(conn redis.Conn, id int) (user *User, err error) {

	//通过传入的id去redis查询该用户是否存在
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		if err == redis.ErrNil {
			//在users hash 中，并没有该id
			err = ErrorUserNotExists
		}
		fmt.Println("查询id是否存在时,err=", err)
		return
	}

	user = &User{}

	//反序列化
	json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return

}

//Login 用于登录校验
func (userDaoThis *UserDao) Login(userID int, userPwd string) (user *User, err error) {

	//先从UserDao 的链接池中取一个连接
	conn := userDaoThis.pool.Get()
	defer conn.Close()
	user, err = userDaoThis.getUserByID(conn, userID)
	if err != nil {
		return
	}

	if user.UserPwd != userPwd {
		err = ErrorUserPwd
		return
	}
	return
}

//Register 注册注册
func (userDaoThis *UserDao) Register(user *message.User) (err error) {

	//先从UserDao 的链接池中取一个连接
	conn := userDaoThis.pool.Get()
	defer conn.Close()
	_, err = userDaoThis.getUserByID(conn, user.UserID)
	if err == nil {
		err = ErrorUserExists
		return
	}

	//说明注册的用户id不存在
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//序列化完成之后,在redis进行Hset
	_, err = conn.Do("HSet", "users", user.UserID, string(data))
	if err != nil {
		fmt.Println("注册用户时,写入Hset失败,err=", err)
	}
	return
}
