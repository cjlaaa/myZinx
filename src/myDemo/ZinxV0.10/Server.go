package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"myDemo/protobufDemo/pb"
	"zinx/ziface"
	"zinx/znet"
)

// 基于Zinx框架来开发的 服务器端应用程序
//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")

	person := &pb.Person{}
	err := proto.Unmarshal(request.GetData(), person)
	if err != nil {
		fmt.Println("unmarshal err: ", err)
	}
	fmt.Println("源数据: ", person)

	//编码
	//将person对象,就是将protobuf的message进行序列化,得到一个二进制文件
	data, err := proto.Marshal(person)
	//data就是我们要进行网络传输的数据,对端需要按照Message Person格式进行解析
	if err != nil {
		fmt.Println("marshal err: ", err)
	}
	err = request.GetConnection().SendMsg(101, data)
	if err != nil {
		fmt.Println(err)
	}

	//先读取客户端的数据,再回写ping..ping..ping
	//fmt.Println("recv from client: msgID = ", request.GetMsgID(),
	//	",data = ", string(request.GetData()))
	//
	//err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	//if err != nil {
	//	fmt.Println(err)
	//}
}

//hello Zinx test 自定义路由
type HelloZinxRouter struct {
	znet.BaseRouter
}

//Test Handle
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle...")

	//先读取客户端的数据,再回写ping..ping..ping

	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
		",data = ", string(request.GetData()))

	person := &pb.Person{}
	err := proto.Unmarshal(request.GetData(), person)
	if err != nil {
		fmt.Println("unmarshal err: ", err)
	}
	fmt.Println("源数据: ", person)

	clientData := &pb.Person{
		Name:   "hENRYcHANG",
		Age:    32,
		Emails: []string{"cjlaaa@gmail.com", "cjlaaa@126.com"},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{
				Number: "13888888888",
				Type:   pb.PhoneType_MOBILE,
			},
			&pb.PhoneNumber{
				Number: "88886666",
				Type:   pb.PhoneType_HOME,
			},
			&pb.PhoneNumber{
				Number: "18666666666",
				Type:   pb.PhoneType_WORK,
			},
		},
	}

	clientByte, err := proto.Marshal(clientData)

	err = request.GetConnection().SendMsg(201, clientByte)
	if err != nil {
		fmt.Println(err)
	}
}

//创建链接之后执行钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("====> DoConnectionBegin os Called... ")
	if err := conn.SendMsg(202, []byte("DoConnection BEGIN")); err != nil {
		fmt.Println(err)
	}

	//给当前的连接设置一些属性
	fmt.Println("Set conn Name, Home ...")
	conn.SetProperty("Name", "机智的常总")
	conn.SetProperty("Github", "https://github.com/cjlaaa")
	conn.SetProperty("Home", "https://cjlaaa.github.io")
}

//连接断开之前的需要执行的函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("====> DoConnectionList is Called... ")
	fmt.Println("conn ID = ", conn.GetConnID(), " is Lost ...")

	//获取连接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}
	if github, err := conn.GetProperty("Github"); err == nil {
		fmt.Println("Github = ", github)
	}
	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Home = ", home)
	}
}

func main() {
	// 1创建一个server句柄, 使用zinx的api
	s := znet.NewServer("[zinx V0.10]")
	// 2注册连接Hook钩子函数
	//s.SetOnConnStart(DoConnectionBegin)
	//s.SetOnConnStop(DoConnectionLost)
	// 3给当前zinx框架添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	// 4启动server
	s.Serve()
}
