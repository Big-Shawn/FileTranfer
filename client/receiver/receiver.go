package receiver

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"iproto/frame"
	"log"
	"net"
	"os"
)

const PackageSize = 2048
const ServerDone = "done"

type Unit struct {
	f    *frame.Frame
	conn *net.Conn
	file *os.File
	size int64
	path string
}

func New(p string) *Unit {
	return &Unit{path: p}
}

func (u *Unit) Handle(m []byte) error {
	f := &frame.Frame{}
	err := proto.Unmarshal(m, f)
	if err != nil {
		return err
	}
	u.f = f
	if err := u.Landing(); err != nil {
		return err
	}
	return nil
}

func (u *Unit) Read(conn *net.Conn) error {
	msg := make(chan []byte, 20)
	stop := make(chan struct{})
	defer close(msg)

	go func() {
		pkg := make([]byte, 0, PackageSize)
		for phase := range msg {
			available := PackageSize - len(pkg)

			if len(phase) < available {
				pkg = append(pkg, phase...)
			} else {
				pkg = append(pkg, phase[:available]...)
			}

			if len(pkg) == PackageSize {
				err := u.Handle(pkg)
				if err != nil {
					stop <- struct{}{}
					log.Printf("\n protobuf handle error: %s\n", err)
					break
				}
				pkg = make([]byte, 0, PackageSize)
			}

			if len(phase) > available {
				pkg = append(pkg, phase[available:]...)
			}
		}
	}()

	go func() {
		<-stop

		(*conn).Close()
	}()

	for {

		b := make([]byte, PackageSize)
		n, err := (*conn).Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		msg <- b[:n]
	}

	return nil
}

func (u *Unit) Landing() error {
	if u.f.Type == frame.FrameType_Info {
		file, err := os.OpenFile(u.path+"/"+string(u.f.Body), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
		u.file = file
		fmt.Printf("\nfile transfer start : %s", string(u.f.Body))
	}

	if u.f.Type == frame.FrameType_Data {
		n, err := u.file.Write(u.f.Body)
		if err != nil {
			return err
		}
		u.size += int64(n)
		fmt.Printf("\nfile transferring received: %d KB \u001B[A ", u.size/1024)
	}

	if u.f.Type == frame.FrameType_Conn && string(u.f.Body) == ServerDone {
		if e := u.file.Sync(); e != nil {
			return e

		}
		if e := u.file.Close(); e != nil {
			return e
		}
		fmt.Printf("\nfile transfer completed")
	}
	return nil

}
