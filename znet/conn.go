package znet

import (
	"errors"
	"fmt"
	"github.com/neversaynevernz/zinx/utils"
	"io"
	"net"

	"github.com/neversaynevernz/zinx/ziface"
)

// 实现链接模块
type Connection struct {

	// 当前Conn隶属于哪个Server
	TcpServer ziface.IServer

	//当前链接的套接字
	Conn *net.TCPConn

	// 链接的ID
	ConnID uint32

	// 链接状态
	isClosed bool

	// 当前链接的处理业务方法的API
	handleAPI ziface.HandleFunc

	// 告知当前链接已经退出/停止的 channel
	// 由 reader 告诉 writer 退出
	ExitChan chan bool

	// 无缓冲通道 用于读、写 goroutine之间的消息通信
	msgChan chan []byte

	// 消息的管理MsgID和对应业务API处理关系
	MsgHandler ziface.IMsgHandler
}

// 初始化链接
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}

	//将 conn 加入到 ConnManager 中
	c.TcpServer.GetConnManager().Add(c)

	return c
}

func (c *Connection) StartReader() {

	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Printf("[Reader Goroutine is exit!], connID[%d], remote addr: %s\n", c.ConnID, c.RemoteAddr().String())
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
			break
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

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制， 将消息发送给Worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中 找到注册绑定的 conn 对应的 router调用
			// 根据绑定好的MsgID 找到对应的api业务处理
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// 主要用来写消息的 goroutine
// 专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println("Writer Goroutine is exit!]", c.RemoteAddr().String())
	defer c.Stop()

	// 不断的阻塞的等待channel 的消息 进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error:", err)
				return
			}
		case <-c.ExitChan:
			// 代表 Reader 已经退出,此时 Writer 需退出
			return
		}
	}

}

// 启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Printf("Connection Start:[%d]\n", c.ConnID)

	// 启动从当前链接的读数据的业务
	go c.StartReader()

	// 启动从当前链接的写数据的业务
	go c.StartWriter()

	// 按照开发者传递进来的 创建链接之后需要调用的处理业务, 执行对应的Hook函数
	c.TcpServer.CallOnConnStart(c)
}

// 关闭链接 结束当前链接的工作
func (c *Connection) Stop() {

	fmt.Printf("Connection stop:[%d]\n", c.ConnID)

	if c.isClosed == true {
		return
	}

	c.isClosed = true

	// 按照开发者传递进来的 销毁链接之前需要调用的处理业务, 执行对应的Hook函数
	c.TcpServer.CallOnConnStop(c)

	// 关闭 socket链接
	c.Conn.Close()

	// 告知 Writer 关闭
	c.ExitChan <- true

	// 将当前链从 ConnMgr 中摘除掉
	c.TcpServer.GetConnManager().Remove(c)

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
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
	c.msgChan <- binaryMsg

	return nil
}
