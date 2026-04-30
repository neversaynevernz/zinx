package utils

import (
	"encoding/json"
	"errors"
	"github.com/neversaynevernz/zinx/ziface"
	"os"
)

// 存储框架全局配置
// 一些参数可以通过 zinx.Json 由用户配置
type GlobalObj struct {

	// Server
	TcpServer ziface.IServer //当前Zinx全局的Server 对象
	Host      string         //当前服务器监听的IP
	TcpPort   int            //当前服务器监听的端口号
	Name      string         //当前服务器的名称

	//Zinx
	Version        string // 当前Zinx 版本号
	MaxConn        int    //当前服务器主机允许的最大链接数
	MaxPackageSize uint32 //当前Zinx框架数据包的最大值
}

/*
定义一个全局的对外
*/

var confpath string = "conf/zinx.json"

var GlobalObject *GlobalObj

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, os.ErrNotExist)
}

func (g *GlobalObj) Reload() {

	if ok := fileExists(confpath); !ok {
		return
	}

	data, err := os.ReadFile(confpath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
	提供一个 init 对象初始化当前的 GlobalObject
*/

func init() {
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "v1.6.0",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1024,
		MaxPackageSize: 4096,
	}
	//从用户配置加载 用户自定义参数
	GlobalObject.Reload()
}
