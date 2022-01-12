
package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"myDemo/protobufDemo/pb"
)

func main() {
	//定义一个Person结构对象
	person := &pb.Person{
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

	//编码
	//将person对象,就是将protobuf的message进行序列化,得到一个二进制文件
	data, err := proto.Marshal(person)
	//data就是我们要进行网络传输的数据,对端需要按照Message Person格式进行解析
	if err != nil {
		fmt.Println("marshal err: ", err)
	}

	//解码
	newPerson := &pb.Person{}
	err = proto.Unmarshal(data, newPerson)
	if err != nil {
		fmt.Println("unmarshal err: ", err)
	}
	fmt.Println("源数据: ", person)
	fmt.Println("解码之后的数据: ", newPerson)
}
