package model

import (
	"errors"
)

var (
	//ErrorUserNotExists 用户不存在
	ErrorUserNotExists = errors.New("用户不存在")

	//ErrorUserExists 用户已存在
	ErrorUserExists = errors.New("用户已存在")

	//ErrorUserPwd 密码不正确
	ErrorUserPwd = errors.New("密码不正确")

	//ErrorOnlineUserNotExists 在线用户不存在
	ErrorOnlineUserNotExists = errors.New("在线用户不存在")

	//ErrorNothingUnreadMesToAll 没有未读的@全体成员的消息
	ErrorNothingUnreadMesToAll = errors.New("没有未读的@全体成员消息")
)
