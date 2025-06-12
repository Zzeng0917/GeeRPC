package codec

import "io"

type Header struct {
	ServiceMethod string //服务名和方法名
	Seq           uint64 //请求ID
	Error         string //错误信息
}

type Codec interface {
	io.Closer                   //释放编解码器资源
	ReadHeader(*Header) error   //从输入流读取请求头
	ReadBody(interface{}) error //从输入流读取消息体
	Write(*Header, interface{}) error
	//将消息头和消息体序列化为字节流，并写入输出流
}

type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "applicaiton/json"
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
