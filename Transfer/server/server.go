package main

import (
	"crypto/md5"
	"net"
	"server/response"
)

func main() {
	host := "127.0.0.1:8879"
	listen, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
	}

	// 接受数据客户端的连接

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}

		go resolveConnection(conn)
	}

}

func resolveConnection(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	uuid := md5.Sum([]byte(addr))
	uuids := make([]byte, len(uuid))
	// 数组转切片
	uuids = append(uuids, uuid[:]...)
	frame.Header()
	response := response.Body{
		Code: 0,
		Msg:  string(uuids),
		Data: response.Data{},
	}

}
