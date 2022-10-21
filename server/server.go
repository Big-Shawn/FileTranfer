package main

import (
	"io"
	"net"
	"os"
	"strings"
)

func main() {

	host := "127.0.0.1:8088"
	server := buildServer(host)
	targetPath := "/home/shawn/server"

	for {
		accept, err := server.Accept()
		if err != nil {
			panic(err)
		}
		go resolveConnection(accept, targetPath)
	}

}

func receiveFileName(conn net.Conn, path string) *os.File {
	bytes := make([]byte, 100)
	n, err := conn.Read(bytes)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		panic(err)

	}
	fileName := string(bytes[:n])
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	open, err := os.Open(path + fileName)
	if err != nil {
		panic(err)
	}
	return open
}

func receiveFile(handler io.Writer, conn net.Conn) {
	for {
		readBuffer := make([]byte, 2048)
		n, err := conn.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				return
			}
			panic("Receive Error :" + err.Error())
		}
		_, err = handler.Write(readBuffer[:n])
		if err != nil {
			panic("Receive Write Error: " + err.Error())
		}
	}

}

func resolveConnection(conn net.Conn, path string) {

	for {
		handler := receiveFileName(conn, path)
		if handler == nil {
			continue
		}
		receiveFile(handler, conn)
	}

}

func buildServer(host string) net.Listener {
	listen, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
	}
	return listen

}
