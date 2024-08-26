package main

import (
	"google.golang.org/protobuf/proto"
	"iproto/frame"
	"log"
	"net"
	"os"
)

const PackageSize = 2048
const ClientReady = "ready"
const ServerDone = "done"
const FilePath = "file.txt"

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

		f, e := readFrame(&conn)
		if e != nil {
			log.Println("readFrame error:", err)

			break
		}
		if f.Type == frame.FrameType_Info {
			newfile, e := os.OpenFile(string(f.Body), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if e != nil {
				log.Println("os.Open error:", e)
				continue
			}

			for ; e == nil; f, e = readFrame(&conn) {
				if e != nil {
					log.Println("readFrame error:", e)
					continue
				}
				if f.Type == frame.FrameType_Data {
					n, e := newfile.Write(f.Body)
					if e != nil {
						log.Println("newfile write error:", e)
					}
					log.Println("newfile write size:", n)
				} else if f.Type == frame.FrameType_Conn && string(f.Body) == ServerDone {
					newfile.Sync()
					newfile.Close()
					break
				}
			}

		}

	}

}

func readFrame(conn *net.Conn) (*frame.Frame, error) {
	buffer := make([]byte, PackageSize)

	n, err := (*conn).Read(buffer)
	if err != nil {
		log.Printf("read from conn err: %s\n", err)
		return nil, err
	}
	if n != PackageSize {
		log.Printf("Read %d bytes, expected %d\n", n, PackageSize)
		return nil, err
	}
	f := frame.Frame{}
	err = proto.Unmarshal(buffer, &f)

	if err != nil {
		log.Printf("unmarshal err : %s\n", err)
		return nil, err
	}
	return &f, nil
}
