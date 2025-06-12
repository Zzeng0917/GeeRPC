package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

var _ Codec = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush()
		//将缓冲区的数据强制写入连接
		if err != nil {
			_ = c.Close()
		}
		// 关闭连接
	}() //无论程序以何种方式退出都执行

	//序列化消息头为二进制数据，并写入缓冲区，若失败记录日志并返回错误
	if err := c.enc.Encode(h); err != nil {
		log.Println("rpc codec:gob error encoding header:", err)
		return err
	}
	//类似地，序列化消息体
	if err := c.enc.Encode(body); err != nil {
		log.Println("rpc codec:gob error encoding body:", err)
		return err
	}
	return nil //所有步骤成功，返回nil表示没有错误
}

func (c *GobCodec) Close() error {
	return c.conn.Close()
}
