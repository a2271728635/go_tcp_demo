package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go_code/chatroom/common/message"
	"net"
)

//Transfer 将读写方法关联到结构体中
type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte //传输时使用的缓冲
}

//ReadPkg 是读取信息的工具函数
func (thisUtils *Transfer) ReadPkg() (mes message.Message, err error) {

	//buf := make([]byte, 8096)

	fmt.Println("读取客户端发送的")
	_, err = thisUtils.Conn.Read(thisUtils.Buf[:4])
	if err != nil {
		//err = errors.New("readPkg head err")
		return
	}
	//fmt.Println("读到的buf=", buf[:4])
	var pakLen uint32
	pakLen = binary.BigEndian.Uint32(thisUtils.Buf[0:4])

	//根据 pakLen 读取内容
	n, err := thisUtils.Conn.Read(thisUtils.Buf[:pakLen])
	if n != int(pakLen) || err != nil {
		//fmt.Println("conn.Read err=", err)
		return
	}
	//反序列化
	//!!&mes
	err = json.Unmarshal(thisUtils.Buf[:pakLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

//WritePkg 是发送信息的工具函数
func (thisUtils *Transfer) WritePkg(data []byte) (err error) {
	//先发送长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	//var buf [4]byte
	binary.BigEndian.PutUint32(thisUtils.Buf[0:4], pkgLen)

	//发送信息长度
	n, err := thisUtils.Conn.Write(thisUtils.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(buf) err=", err)
		return
	}

	//发送消息本身
	n, err = thisUtils.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(data) err=", err)
		return
	}

	return
}
