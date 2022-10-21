package main

import (
	"io"
	"io/fs"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func main() {
	host := "127.0.0.1:8088"
	connection := buildConnection(host)

	path := "/home/shawn/clienFile"
	info, _ := getDirInfo(path)

	// 考虑多个文件同时传输时的问题，即传输的文件名可能和文件内容不一致时的情况
	//for _, file := range info {
	//	filePath, name := readFile(file, path)
	//	send(connection, []byte("start"))
	//	send(connection, []byte(name))
	//	open, err := os.Open(filePath)
	//	if err != nil {
	//		panic("File Open err : " + err.Error())
	//	}
	//	for {
	//		readBuffer := make([]byte, 2048)
	//		n, err := open.Read(readBuffer)
	//		if err != nil {
	//			if err == io.EOF {
	//				break
	//			}
	//			panic("File Read err : " + err.Error())
	//		}
	//		send(connection, readBuffer[:n])
	//		// defer open.Close() ?
	//	}
	//	open.Close()
	//	send(connection, []byte("ok"))
	//}

	for _, file := range info {
		filePath, name := readFile(file, path)
		send(connection, []byte(name))
		sendFile(filePath, connection)
	}

}

func buildConnection(host string) net.Conn {
	dial, err := net.Dial("tcp", host)
	if err != nil {
		panic(err)
	}
	return dial
}

func getDirInfo(path string) ([]fs.FileInfo, int) {
	//获取目录下所有文件信息
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	count := len(files)

	return files, count
}

func readFile(file fs.FileInfo, path string) (filePath, fileName string) {
	if file.IsDir() {
		panic(file.Name() + " is a dir")
	}
	if suffix := strings.HasSuffix(path, "/"); !suffix {
		path += "/"
	}

	filePath = path + file.Name()

	return filePath, fileName
}

func send(target io.Writer, content []byte) {
	// 1. 发送文件片段
	// 2. 发送文件开始发送信号，文件发送结束信号
	_, err := target.Write(content)
	if err != nil {
		panic(err)
	}
}

func sendFile(filePath string, connection net.Conn) {
	open, err := os.Open(filePath)
	if err != nil {
		panic("File Open err : " + err.Error())
	}
	for {
		readBuffer := make([]byte, 2048)
		n, err := open.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic("File Read err : " + err.Error())
		}
		send(connection, readBuffer[:n])
		// defer open.Close() ?
	}
	defer open.Close()
}
