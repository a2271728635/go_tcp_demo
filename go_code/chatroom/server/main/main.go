package main

import (
	"fmt"
	"go_code/chatroom/server/model"
	"net"
	"time"
)

//处理客户端发过来的
func pro(conn net.Conn) {

	//记得关闭
	defer conn.Close()

	//这里调用总控制台
	processor := &Processor{
		Conn: conn,
	}
	err := processor.process2()

	if err != nil {
		fmt.Println("服务器端调用总控制台时出错,err=", err)
		return
	}
}

//这里写魔法函数,完成对UserDao的初始化任务
func initUserDao() {
	//注意 这里的pool是一个全局变量
	//应先初始化main.initPool,再初始化main.initUserDao
	model.MyUserDao = model.NewUserDao(pool)
}

//这里写魔法函数,完成对SmsDao的初始化任务
func initSmsDao() {
	//注意 这里的pool是一个全局变量
	//应先初始化main.initPool,再初始化main.initSmsDao
	model.MySmsDao = model.NewSmsDao(pool)
}

func main() {

	//当服务器启动时,初始化redis的连接池
	initPool("localhost:6379", 16, 0, 300*time.Second)
	//初始化UserDao
	initUserDao()

	//初始化SmsDao
	initSmsDao()

	fmt.Println("新结构服务器端的8889端口监听")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	//关闭
	defer listen.Close()
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}

	//监听端口有成功后，等待客户端来链接服务器
	for {
		fmt.Println("等待客户端来链接")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
			continue
		}
		//一旦链接成功，则启动一个协程和客户端保持通讯
		go pro(conn)
	}
}
