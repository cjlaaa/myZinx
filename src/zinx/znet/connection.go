package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

//连接模块
type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn

	//连接的ID
	ConnID uint32

	//当前的连接状态
	isClosed bool

	//告知当前连接已经退出/停止的channel(由Reader告知Writer退出)
	ExitChan chan bool

	//无缓冲的管道,用于读/写Goroutine之间的消息通信
	msgChan chan []byte

	//消息的管理MsgID和对应的处理业务API关系
	MsgHandle ziface.IMsgHandle
}

//初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		MsgHandle: msgHandler,
		isClosed:  false,
		msgChan:   make(chan []byte),
		ExitChan:  make(chan bool, 1),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println("[Reader is exit!],connID = ", c.ConnID, ", remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//创建一个拆包解包的对象
		dp := NewDataPack()

		//读取客户端的Msg Head 二进制流8个字节,
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}

		//拆包,得到MsgID和MsgDatalen放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		//根据dataLen再次读取Data,放在Msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		//从路由中,找到注册绑定的Conn对应的router调用
		//根据绑定好的MsgID找到对应处理api业务执行
		go c.MsgHandle.DoMsgHandler(&req)
	}
}

//写消息
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println("[conn Writer exit!]", c.RemoteAddr().String())

	//不断的阻塞的等待channel的消息, 进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error, ", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经退出,此时Writer也要退出
			return
		}

	}
}

//启动连接 让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start().. ConnID = ", c.ConnID)
	//启动从当前连接读数据的业务
	go c.StartReader()
	//启动从当前连接写数据的业务
	go c.StartWriter()
}

//停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)

	//如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//关闭socket连接
	c.Conn.Close()

	//告知Writer关闭
	c.ExitChan <- true

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

//获取当前连接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接模块的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的 TCP状态 IP Port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//提供一个SendMsg方法 将我们要发送给客户端的数据,先进行封包,在发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	//将data进行封包 MsgDataLen|MsgID|Data
	dp := NewDataPack()

	//MsgDataLen|MsgId|Data
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}

	//将数据发送给客户端
	c.msgChan <- binaryMsg

	return nil
}
