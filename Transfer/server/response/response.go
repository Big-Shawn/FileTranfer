package response

import "net"

type Data struct {
	FileName string `json:"file_name,omitempty"`
	FileSize string `json:"file_size,omitempty"`
	ClientId string `json:"client_id,omitempty"`
	FileData []byte `json:"file_data,omitempty"`
}

type Body struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data Data   `json:"data"`
}

// 发送响应
func (r *Body) Send(conn net.Conn) {
	// json 编码
	// 加上

}
