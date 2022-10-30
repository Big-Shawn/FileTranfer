package frame

type Frame []byte

const FrameSize = 1024

// getFrames 根据数据帧的定义返回本次从数据流中获取的数据切片
func (f *Frame) getFrames(receive []byte) [][]byte {
	//
	return nil

}
func SendHeader(len int) {

}
