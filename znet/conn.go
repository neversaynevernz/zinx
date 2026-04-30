package znet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/neversaynevernz/zinx/ziface"
)

// 实现链接模块
type Connection struct {

	//当前链接的套接字
	Conn *net.TCPConn

	// 链接的ID
	ConnID uint32

	// 链接状态
	isClosed bool

	// 当前链接的处理业务方法的API
	handleAPI ziface.HandleFunc

	// 告知当前链接停止的 chan
	ExitChan chan bool

	// 当前链接处理的 方法 Router
	Router ziface.IRouter
}

// 初始化链接
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool),
	}
}

func (c *Connection) StartReader() {

	fmt.Println("Start reader ...")
	defer fmt.Printf("ConnID[%d], Reader exits, remote addr: %s", c.ConnID, c.RemoteAddr().String())
	defer c.Stop()

	for {
		
		//创建一个拆包解包的对象
		dp := NewDataPack()

		// 读取客户端的Msg Head 二进制流 8 个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head err:", err)
			break
		}

		// 拆包 得到 msgID 和 datalen 放在msg 消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err:", err)
			return
		}

		// 根据 dataLen 再次读取 Data 放在msg.Data 中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err:", err)
				break
			}
		}

		msg.SetData(data)

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 执行注册的路由方法
		go func(req ziface.IRequest) {
			// 从路由中 找到注册绑定的 conn 对应的 router调用
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.PostHandle(req)
		}(&req)
	}
}

// 启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Printf("Connection Start:[%d]\n", c.ConnID)

	// 启动从当前链接的读数据的业务
	go c.StartReader()

	// TODO当前链接的写数据的业务
}

// 关闭链接 结束当前链接的工作
func (c *Connection) Stop() {

	fmt.Printf("Connection stop:[%d]\n", c.ConnID)

	if c.isClosed == true {
		return
	}

	c.isClosed = true

	// 关闭 socket链接
	c.Conn.Close()

	close(c.ExitChan)
}

// 获取当前链接的对象 套接字
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 得到当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 得到客户端链接的地址和端口
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据的方法
// 提供一个 SendMsg 方法 将我们要发送给客户端的数据 先进行封包
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	// 将 data 进行封包 |MsgDataLen|MsgID|Data|
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return err
	}
	// 将数据发送给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write msg id: ", msgId, "error:", err)
		return errors.New("conn write error")
	}
	return nil
}
