package response

import (
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"iproto/frame"
	"net"
	"os"
	"path"
	"server/log"
	"strings"
)

const ServerDone = "done"

const PackageSize = 1460 * 2 // 当MTU=1500时 TCP最大有效负载为1460，这里取整数倍，减少读的次数。
//const PackageSize = 2048

type Unit struct {
	f    *frame.Frame
	conn *net.Conn
	file *os.File
}

func (u *Unit) setFrame(t frame.FrameType) {
	u.f = &frame.Frame{Type: t, Size: -1, Body: make([]byte, 1), Reserved: make([]byte, 1)}
}

func (u *Unit) flush() {
	u.f.Reserved = make([]byte, 1)
}

func New(conn *net.Conn) *Unit {
	return &Unit{
		conn: conn,
	}
}

func (u *Unit) SetFile(file *os.File) *Unit {
	u.file = file
	return u
}

func (u *Unit) Send() error {
	if u.file == nil {
		return errors.New("no file to send")
	}

	if e := u.SendInfo(); e != nil {
		return e
	}
	if e := u.SendData(); e != nil {
		return e
	}
	if e := u.SendSignal(ServerDone); e != nil {
		return e
	}
	return nil
}

func (u *Unit) SendSignal(signal string) error {
	u.setFrame(frame.FrameType_Conn)
	u.f.Body = []byte(signal)
	if e := u.send(); e != nil {
		return fmt.Errorf("SendSignal error: %v", e)
	}
	return nil
}

func (u *Unit) SendInfo() error {
	u.setFrame(frame.FrameType_Info)
	_, fname := path.Split(u.file.Name())
	u.f.Body = []byte(fname)
	//u.f.Size = int32(u.file.Size)
	if e := u.send(); e != nil {
		return fmt.Errorf("SendInfo error: %v", e)
	}
	return nil
}

// 还是需要设置一个默认包大小，这样才能保证第一次顺利通讯
func generateCommunicationInfo(fname string, speed int) string {
	return strings.Join([]string{fname, fmt.Sprintf("%d", speed)}, "|")
}

func (u *Unit) getFileSlice(begin, offset int) ([]byte, int, error) {
	b := make([]byte, offset)
	n, err := u.file.ReadAt(b, int64(begin))
	return b, n, err
}

func (u *Unit) SendData() error {
	u.setFrame(frame.FrameType_Data)
	var sent int
	space := PackageSize - proto.Size(u.f)
	for {
		slice, i, err := u.getFileSlice(sent, space)
		if err == io.EOF && i == 0 {
			break
		}

		u.f.Size = int32(i)
		u.f.Body = slice[:i]
		sent += i
		if e := u.send(); e != nil {
			return fmt.Errorf("SendData error: %v", e)
		}
		u.flush()
	}
	log.L.Sugar().Infof("file:%s, sent:%d B \n", u.file.Name(), sent)

	return nil
}

func (u *Unit) send() error {
	if n := PackageSize - proto.Size(u.f); n > 0 {
		e := make([]byte, n)
		u.f.Reserved = e
	}

	marshal, err2 := proto.Marshal(u.f)
	if err2 != nil {
		return fmt.Errorf("proto marshaling error: %v", err2)
	}

	if _, err := (*u.conn).Write(marshal); err != nil {
		return fmt.Errorf("write error: %v", err)

	}
	return nil
}

func (u *Unit) Read(conn *net.Conn) (*Unit, error) {
	buffer := make([]byte, PackageSize)

	n, err := (*conn).Read(buffer)
	if err != nil {
		log.L.Sugar().Errorf("read from conn err: %s\n", err)
		return nil, err
	}
	if n != PackageSize {
		log.L.Sugar().Errorf("Read %d bytes, expected %d\n", n, PackageSize)
		return nil, err
	}
	f := frame.Frame{}
	err = proto.Unmarshal(buffer, &f)

	if err != nil {
		log.L.Sugar().Errorf("unmarshal err : %s\n", err)
		return nil, err
	}
	return &Unit{f: &f}, nil
}
