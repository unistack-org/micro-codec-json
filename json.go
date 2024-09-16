// Package json provides a json codec
package json

import (
	"bytes"
	"encoding/json"
	"reflect"

	pb "go.unistack.org/micro-proto/v3/codec"
	"go.unistack.org/micro/v3/codec"
	rutil "go.unistack.org/micro/v3/util/reflect"
)

var _ codec.Codec = (*jsonCodec)(nil)

var (
	DefaultMarshalOptions = JsonMarshalOptions{
		EscapeHTML: true,
	}

	DefaultUnmarshalOptions = JsonUnmarshalOptions{
		DisallowUnknownFields: false,
		UseNumber:             false,
	}
)

type JsonMarshalOptions struct {
	EscapeHTML bool
}

type JsonUnmarshalOptions struct {
	DisallowUnknownFields bool
	UseNumber             bool
}

type jsonCodec struct {
	opts codec.Options
}

func (c *jsonCodec) Marshal(v interface{}, opts ...codec.Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, err := rutil.StructFieldByTag(v, options.TagName, "flatten"); err == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		return m.Data, nil
	case *pb.Frame:
		return m.Data, nil
	}

	marshalOptions := DefaultMarshalOptions
	if options.Context != nil {
		if f, ok := options.Context.Value(marshalOptionsKey{}).(JsonMarshalOptions); ok {
			marshalOptions = f
		}
	}

	if !marshalOptions.EscapeHTML {
		w := bytes.NewBuffer(nil)
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(marshalOptions.EscapeHTML)
		err := enc.Encode(v)
		buf := w.Bytes()
		return buf[:len(buf)-1], err
	}

	return json.Marshal(v)
}

func (c *jsonCodec) Unmarshal(b []byte, v interface{}, opts ...codec.Option) error {
	if len(b) == 0 || v == nil {
		return nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, err := rutil.StructFieldByTag(v, options.TagName, "flatten"); err == nil {
			v = nv
			rv := reflect.ValueOf(v)
			if rv.Kind() != reflect.Pointer &&
				rv.Kind() != reflect.Map {
				v = reflect.New(rv.Type()).Interface()
			}
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		m.Data = b
		return nil
	case *pb.Frame:
		m.Data = b
		return nil
	}

	unmarshalOptions := DefaultUnmarshalOptions
	if options.Context != nil {
		if f, ok := options.Context.Value(unmarshalOptionsKey{}).(JsonUnmarshalOptions); ok {
			unmarshalOptions = f
		}
	}

	if unmarshalOptions.DisallowUnknownFields || unmarshalOptions.UseNumber {
		dec := json.NewDecoder(bytes.NewBuffer(b))
		if unmarshalOptions.DisallowUnknownFields {
			dec.DisallowUnknownFields()
		}
		if unmarshalOptions.UseNumber {
			dec.UseNumber()
		}

		return dec.Decode(v)
	}

	return json.Unmarshal(b, v)
}

func (c *jsonCodec) String() string {
	return "json"
}

func NewCodec(opts ...codec.Option) codec.Codec {
	return &jsonCodec{opts: codec.NewOptions(opts...)}
}
