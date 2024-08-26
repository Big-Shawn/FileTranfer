package main

import (
	"google.golang.org/protobuf/proto"
	"io"
	"iproto/frame"
	"log"
	"net"
	"os"
)

const PackageSize = 2048
const ClientReady = "ready"
const ServerDone = "done"
const FilePath = ""

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

	open, err := os.Open(FilePath)
	if err != nil {
		log.Printf("open file err: %s\n", err)
		return
	}
	all, err := io.ReadAll(open)
	defer open.Close()

	SendFileInfo(conn, open.Name(), len(all))
	SendFileData(conn, &all)
	SendComplete(conn)

}

func SendFileInfo(conn *net.Conn, name string, size int) {
	f := frame.Frame{
		Type: frame.FrameType_Info,
		Size: int32(size),
		Body: []byte(name),
	}
	writePacket(conn, &f)
}

func GetAFrame(t frame.FrameType) *frame.Frame {
	return &frame.Frame{Type: t, Size: -1}
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
	f := frame.Frame{
		Type: frame.FrameType_Conn,
		Size: -1,
		Body: []byte(ServerDone),
	}
	writePacket(conn, &f)
}

func writePacket(conn *net.Conn, f *frame.Frame) {
	e := make([]byte, PackageSize-proto.Size(f))
	f.Reserved = e
	n, err := (*conn).Write(e)
	if err != nil {
		log.Printf("write to conn err: %s\n", err)
	}
	log.Printf("wrote %d bytes\n", n)
}
