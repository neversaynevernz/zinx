package znet

import (
	"fmt"
	"net"
	"neversaynevernz/zinx/ziface"
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
}

// 初始化链接
func NewConnection(conn *net.TCPConn, connID uint32, handleAPI ziface.HandleFunc) *Connection {
	return &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: handleAPI,
		isClosed:  false,
		ExitChan:  make(chan bool),
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Start reader ...")
	defer fmt.Printf("ConnID[%d], Reader exits, remote addr: %s", c.ConnID, c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户的数据到缓存中
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err:", err)
			continue
		}

		// 调用当前链接绑定的 handleAPI
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Printf("Connection[%d], handle err: %s\n", c.ConnID, err)
			break
		}
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
func (c *Connection) Send(data []byte) error {
	//return c.Conn.Write(data)
	return nil
}
