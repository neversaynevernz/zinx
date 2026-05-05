package ziface

/*
	链接管理抽象层
*/

type IConnManager interface {
	Add(conn IConnection)                   // 添加链接
	Remove(conn IConnection)                // 删除链接
	Get(connID uint32) (IConnection, error) // 获取链家
	Len() int                               // 链接的个数
	ClearConn()                             // 删除并链接
}
