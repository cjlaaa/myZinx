package main

import (
	"zinx/znet"
)

// 基于Zinx框架来开发的 服务器端应用程序
func main(){
	// 1创建一个server句柄, 使用zinx的api
	s := znet.NewServer("[zinx V0.2]")
	// 2启动server
	s.Serve()
}