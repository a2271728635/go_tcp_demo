package message

//这里定义消息类型
const (
	LoginMesType = "LoginMes"
	//LoginMesType 客户端登录信息 常量

	LoginResMesType = "LoginResMes"
	//LoginResMesType 服务器端登录的返回信息 常量

	RegisterMesType = "RegisterMes"
	//RegisterMesType 注册信息 常量

	RegisterResMesType = "RegisterResMes"
	//RegisterResMesType 服务器端注册返回信息 常量

	NotifyUserStatusMesType = "NotifyUserStatusMes"
	//NotifyUserStatusMesType 服务器端向客户端,推送用户上线了 常量

	SmsMesType = "SmsMes"
	//SmsMesType 用于客户端发送群发消息 常量

	SmsMesResType = "SmsMesRes"
	//SmsMesResType 客户端发送群发消息后的,服务端接收后,返回给所有在线客户端消息 常量

	SmsPToPMesType = "SmsPToPMes"
	//SmsPToPMesType 用于客户端发送私聊信息 常量

	SmsPToPMesResType = "SmsPToPMesRes"
	//SmsPToPMesType 客户端发送私聊消息后的,服务端接收后,返回给指定客户端消息 常量

	SmsFindUserStatusByIDType = "SmsFindUserStatusByID"
	//SmsFindUserStatusType 客户端查找用户在线状态

	ChangeUserOnlineStatusMesType = "ChangeUserOnlineStatusMes"
	//ChangeUserOnlineStatusMesType 客户端发送,改变在线状态的消息 常量

	ChangeUserOnlineStatusResMesType = "ChangeUserOnlineStatusResMes"
	//ChangeUserOnlineStatusResMesType 服务器端发送,改变在线状态的返回消息 常量
	//使用Code struct 404为失败 200为成功

	NotifyUserOtherStatusMesType = "NotifyUserOtherStatusMes"
	//NotifyUserOtherStatusMesType 服务器端向客户端,推送非上线的用户状态改变 常量

	GetToAllMesType = "GetToAllMes"
	//GetToAllMesType 客户端请求redis上的@全体成员消息 常量

	GetToAllResMesType = "GetToAllResMes"
	//GetToAllResMesType 服务器端向客户端发送,所有@all的消息

)

//这里定义用户在线状态类型
const (
	UserOnline     = iota //在线 0
	UserOffline           //离线 1
	UserBusyStatus        //忙 2
)

//Message 是消息结构体
type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"data"` //消息内容
	//Message 是消息结构体
}

//Code 是服务器端给客户端的返回码
type Code struct {
	Code  int    //返回码 int
	Error string //错误信息
}

//LoginMes 用户登录信息
type LoginMes struct {
	UserID   int    `json:"userId"`   //用户id
	UserPwd  string `json:"userPwd"`  //用户密码
	UserName string `json:"userName"` //用户名
	//LoginMes 用户登录信息
}

//LoginResMes 登录后的返回信息
type LoginResMes struct {
	Code    int    `json:"code"`    //返回状态码 500未注册 200成功
	UsersID []int  `json:"usersId"` //一个保存在线用户的id切片
	UserID  int    `json:"userId"`  //登录用户的id
	Error   string `json:"error"`   //返回错误信息
	//LoginResMes 登录后的返回信息
}

//RegisterMes 注册信息
type RegisterMes struct {
	User User `json:"user"`
	//RegisterMes 注册信息
}

//RegisterResMes 注册后的返回信息
type RegisterResMes struct {
	Code  int    `json:"code"`  //返回状态码 400 该用户已被占用 200注册成功
	Error string `json:"error"` //返回错误信息
	//RegisterResMes 注册后的返回信息
}

//NotifyUserStatusMes 用于服务器端,向客户端推送,用户当前在线状态
type NotifyUserStatusMes struct {
	UserID     int `json:"userId"`     //推送的用户的id
	UserStatus int `json:"userStatus"` //推送的用户的状态
	//NotifyUserStatusMes 用于服务器端,向客户端推送,用户当前在线状态
}

//SmsMes 用于客户端发送全体消息
type SmsMes struct {
	Content string `json:"content"` //消息内容
	User    User   `json:"user"`
}

//SmsMesRes 客户端发送全体消息后的,服务端接收后,返回给指定客户端消息
type SmsMesRes struct {
	Content string `json:"content"` //消息内容
	User    User   `json:"user"`
}

//SmsPToPMes 用于客户端发送私聊信息
type SmsPToPMes struct {
	Content   string `json:"content"`   //消息内容
	FromUser  User   `json:"fromUser"`  //发送者
	ReachUser User   `json:"reachUser"` //接受者
}

//ChangeUserOnlineStatusMes 用户客户端发送,改变在线状态消息
type ChangeUserOnlineStatusMes struct {
	User   User   `json:"user"`   //改变的用户
	Status string `json:"status"` //改变的用户状态
}

//Contents 用于保存@all消息的内容
type Contents struct {
	FromUser User   `json:"fromUser"` //发送者
	Time     int64  `json:"time"`     //发送时间戳
	Content  string `json:"content"`  //发送的内容
}
