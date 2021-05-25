// Package json provides a json codec
package json

import (
	"encoding/json"
	"io"

	"github.com/unistack-org/micro/v3/codec"
	rutil "github.com/unistack-org/micro/v3/util/reflect"
)

type jsonCodec struct{}

const (
	flattenTag = "flatten"
)

func (c *jsonCodec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case nil:
		return nil, nil
	case *codec.Frame:
		return m.Data, nil
	}

	if nv, err := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); err == nil {
		v = nv
	}

	return json.Marshal(v)
}

func (c *jsonCodec) Unmarshal(b []byte, v interface{}) error {
	if len(b) == 0 {
		return nil
	}
	switch m := v.(type) {
	case nil:
		return nil
	case *codec.Frame:
		m.Data = b
		return nil
	}

	if nv, err := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); err == nil {
		v = nv
	}

	return json.Unmarshal(b, v)
}

func (c *jsonCodec) ReadHeader(conn io.Reader, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *jsonCodec) ReadBody(conn io.Reader, v interface{}) error {
	switch m := v.(type) {
	case nil:
		return nil
	case *codec.Frame:
		buf, err := io.ReadAll(conn)
		if err != nil {
			return err
		} else if len(buf) == 0 {
			return nil
		}
		m.Data = buf
		return nil
	}

	buf, err := io.ReadAll(conn)
	if err != nil {
		return err
	} else if len(buf) == 0 {
		return nil
	}

	if nv, nerr := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); nerr == nil {
		v = nv
	}

	return c.Unmarshal(buf, v)
}

func (c *jsonCodec) Write(conn io.Writer, m *codec.Message, v interface{}) error {
	switch m := v.(type) {
	case nil:
		return nil
	case *codec.Frame:
		_, err := conn.Write(m.Data)
		return err
	}

	if nv, nerr := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); nerr == nil {
		v = nv
	}

	buf, err := c.Marshal(v)
	if err != nil {
		return err
	}
	_, err = conn.Write(buf)

	return err
}

func (c *jsonCodec) String() string {
	return "json"
}

func NewCodec() codec.Codec {
	return &jsonCodec{}
}
