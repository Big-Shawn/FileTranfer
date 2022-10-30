package main

import (
	"client/frame"
	"net"
)

// 建立tcp 数据连接
// 获取服务器分配的客户端标识
// 获取文件夹下的文件信息
// 传输文件
// 		1. 数据分片，向主机发送
// 		2. 等待主机回送接收确认
//  	3. 重复 1 2 直至文件发送完成
// 		4. 发送文件结束信号，结束该次文件传输

var serverResponse chan frame.Frame

func main() {
	host := "127.0.0.1:8879"
	conn, err := net.Dial("tcp", host)
	if err != nil {
		panic(err)
	}
	// receive from server
	for {
		buf := make(frame.Frame, frame.Size)
		conn.Read(buf)
		buf.Get()
		buf.Send()
	}

}
