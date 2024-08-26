package main

import (
	"google.golang.org/protobuf/proto"
	"iproto/frame"
	"log"
	"net"
	"server/cache"
)

const PackageSize = 2048
const ClientReady = "ready"
const ServerDone = "done"

// const FilePath = "file.txt"
const FilePath = "/Users/big-shawn/Downloads/Telegram Lite/example_video.mp4"

func main() {
	host := "127.0.0.1:8879"
	listen, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
	}

	// 缓存将要发送的文件信息
	if err = cache.LoadFiles(FilePath); err != nil {
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

func resolveConnection(conn *net.Conn) {
	defer (*conn).Close()

	for _, f := range cache.Files {
		SendFileInfo(conn, f.Name, f.Size)
		SendFileData(conn, f.Content)
		SendComplete(conn)
	}

}

func SendFileInfo(conn *net.Conn, name string, size int) {
	f := GetAFrame(frame.FrameType_Info)
	f.Body = []byte(name)

	writePacket(conn, f)
}

func GetAFrame(t frame.FrameType) *frame.Frame {
	return &frame.Frame{Type: t, Size: -1, Body: make([]byte, 1), Reserved: make([]byte, 1)}
}

func SendFileData(conn *net.Conn, content *[]byte) {
	sent, total := 0, len(*content)
	space := PackageSize - proto.Size(GetAFrame(frame.FrameType_Data))

	for sent < total {
		fn := GetAFrame(frame.FrameType_Data)
		if space >= total {
			fn.Size = int32(total)
			fn.Body = (*content)[sent:total]
			sent += total
			writePacket(conn, fn)
		} else {
			fn.Size = int32(space)
			fn.Body = (*content)[sent : sent+space]
			sent += space
			writePacket(conn, fn)
		}
	}

}

func SendComplete(conn *net.Conn) {
	f := GetAFrame(frame.FrameType_Conn)
	f.Body = []byte(ServerDone)
	writePacket(conn, f)
}

func writePacket(conn *net.Conn, f *frame.Frame) {
	log.Printf("full size : %d", proto.Size(f))
	if n := PackageSize - proto.Size(f); n > 0 {
		e := make([]byte, n)
		f.Reserved = e
	}

	marshal, err2 := proto.Marshal(f)
	if err2 != nil {
		log.Printf("marshal err : %s\n", err2)
	}
	n, err := (*conn).Write(marshal)
	if err != nil {
		log.Printf("write to conn err: %s\n", err)
	}
	log.Printf("wrote %d bytes\n", n)
}
