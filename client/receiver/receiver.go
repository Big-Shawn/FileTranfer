package receiver

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"iproto/frame"
	"log"
	"net"
	"os"
	"path/filepath"
)

const PackageSize = 1460 * 2 // 当MTU=1500典型值时 TCP最大有效负载为1460，这里取整数倍，减少读的次数。
// const PackageSize = 2048
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
	msg := make(chan []byte, 100)
	stop := make(chan error, 1)
	defer (*conn).Close()

	go func() {
		pkg := make([]byte, 0, PackageSize)
		var err error
		for phase := range msg {
			if err != nil {
				continue
				// when producer is more than faster than consumer,
				// consume buffer so that the producer can detect the stop signal
			}
			available := PackageSize - len(pkg)

			if len(phase) < available {
				pkg = append(pkg, pkg...)
			} else {
				pkg = append(pkg, phase[:available]...)
			}

			if len(pkg) == PackageSize {
				err = u.Handle(pkg)
				if err != nil {
					stop <- fmt.Errorf("error handling package file: %+v: %v, lenth: %d, pkg: %v", u.file, err, len(pkg), pkg)
					continue
				}
				pkg = make([]byte, 0, PackageSize)
			}

			if len(phase) > available {
				pkg = append(pkg, phase[available:]...)
			}
		}
		if err == nil {
			stop <- nil
		}
	}()

	for {

		b := make([]byte, PackageSize)
		n, err := (*conn).Read(b)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			return err
		}
		select {
		case e := <-stop:
			close(msg)
			return e
		case msg <- b[:n]:
			{
			}
		}
	}
	close(msg)
	return <-stop
}

// Read2 will be slower than read
func (u *Unit) Read2(conn *net.Conn) error {
	msg := make(chan []byte, 100)
	stop := make(chan error)
	defer (*conn).Close()
	go func() {
		for {
			b := make([]byte, PackageSize)
			n, err := (*conn).Read(b)
			if err == io.EOF && n == 0 {
				break
			}
			if err != nil {
				stop <- err
				break
			}
			msg <- b[:n]
		}
		close(msg)
		stop <- nil
	}()

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
				return err
			}
			pkg = make([]byte, 0, PackageSize)
		}

		if len(phase) > available {
			pkg = append(pkg, phase[available:]...)
		}
	}

	return <-stop
}

func (u *Unit) Landing() error {
	if u.f.Type == frame.FrameType_Info {
		file, err := os.OpenFile(u.path+"/"+string(u.f.Body), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
		u.file = file
		log.Printf("file transfer start : %s \n", string(u.f.Body))
	}

	if u.f.Type == frame.FrameType_Data {
		n, err := u.file.Write(u.f.Body)
		if err != nil {
			return err
		}
		u.size += int64(n)
		fmt.Println()
		fmt.Printf("\u001B[A file transferring received: %d KB", u.size/1024)
	}

	if u.f.Type == frame.FrameType_Conn && string(u.f.Body) == ServerDone {
		if e := u.file.Sync(); e != nil {
			return e
		}
		if e := u.file.Close(); e != nil {
			return e
		}
		fmt.Println()
		log.Printf("file received completed: %s, size: %d KB, %d B \n",
			filepath.Base(u.file.Name()),
			u.size/1024,
			u.size)
	}
	return nil

}
