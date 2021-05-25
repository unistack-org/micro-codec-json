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
		return json.Marshal(nv)
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
		return json.Unmarshal(b, nv)
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

	var err error
	if nv, nerr := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); nerr == nil {
		err = json.NewDecoder(conn).Decode(nv)
	} else {
		err = json.NewDecoder(conn).Decode(v)
	}

	if err == io.EOF {
		return nil
	}

	return err
}

func (c *jsonCodec) Write(conn io.Writer, m *codec.Message, v interface{}) error {
	switch m := v.(type) {
	case nil:
		return nil
	case *codec.Frame:
		_, err := conn.Write(m.Data)
		return err
	}

	var err error
	if nv, nerr := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); nerr == nil {
		err = json.NewEncoder(conn).Encode(nv)
	} else {
		err = json.NewEncoder(conn).Encode(v)
	}

	return err
}

func (c *jsonCodec) String() string {
	return "json"
}

func NewCodec() codec.Codec {
	return &jsonCodec{}
}
