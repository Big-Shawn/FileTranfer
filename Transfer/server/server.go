package main

import (
	"net"
	"server/cache"
	"server/config"
	"server/log"
	"server/response"
)

func main() {

	if e := config.Load(); e != nil {
		panic(e)
	}
	if e := log.Init(); e != nil {
		panic(e)
	}

	// 缓存将要发送的文件信息
	cache.LoadFiles(config.C.Path)

	listen, err := net.Listen("tcp", config.C.Port)
	if err != nil {
		panic(err)
	}

	// 接受数据客户端的连接
	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}

		go resolveConnection(&conn)
	}

}

func resolveConnection(conn *net.Conn) {
	defer (*conn).Close()
	r := response.New(conn)

	for _, f := range cache.Files {
		if err := r.SetFile(f.Handler).Send(); err != nil {
			log.L.Sugar().Errorf("file sent error: %s", err)
		}
	}

}
