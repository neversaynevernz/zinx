package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/neversaynevernz/zinx/utils"
	"github.com/neversaynevernz/zinx/ziface"
)

// 封包、拆包的具体模块
type DataPack struct{}

// 封装拆包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包的头的长度
func (d *DataPack) GetHeadLen() uint32 {
	// Datalen uint32 (4个字节) + ID uint32 (4个字节)
	return 8
}

// 封包方法
// |datalen|msgid|data|
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放 bytes 字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	// 将dataLen 写进databuff 中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// 将MsgId 写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 将data 写进databuff 中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

// 拆包方法
// 将包的Head信息都出来 之后再根据head信息里的data的长度 再进行一次读
func (d *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压head信息 得到 datalen 和 MsgID
	msg := &Message{}
	// 读 datalen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// 读 MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断datalen 是否超出了设置的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv!")
	}

	// 为什么不读 data ????
	return msg, nil
}
