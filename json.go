// Package json provides a json codec
package json

import (
	"encoding/json"
	"io"

	"github.com/unistack-org/micro/v3/codec"
)

type jsonCodec struct {
}

func (c *jsonCodec) ReadHeader(conn io.ReadWriter, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *jsonCodec) ReadBody(conn io.ReadWriter, b interface{}) error {
	if b == nil {
		return nil
	}
	return json.NewDecoder(conn).Decode(b)
}

func (c *jsonCodec) Write(conn io.ReadWriter, m *codec.Message, b interface{}) error {
	if b == nil {
		return nil
	}
	return json.NewEncoder(conn).Encode(b)
}

func (c *jsonCodec) Marshal(b interface{}) ([]byte, error) {
	if b == nil {
		return nil, nil
	}
	return json.Marshal(b)
}

func (c *jsonCodec) Unmarshal(b []byte, v interface{}) error {
	if b == nil {
		return nil
	}
	return json.Unmarshal(b, v)
}

func (c *jsonCodec) String() string {
	return "json"
}

func NewCodec() codec.Codec {
	return &jsonCodec{}
}
