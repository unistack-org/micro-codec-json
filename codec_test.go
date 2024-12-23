package json

import (
	"bytes"
	"testing"

	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/codec"
)

func TestRawMessage(t *testing.T) {
	b := &broker.Message{}
	buf, err := NewCodec().Marshal(b)
	if err != nil {
		panic(err)
	}
	_ = buf
}

func TestFrame(t *testing.T) {
	s := &codec.Frame{Data: []byte("test")}

	buf, err := NewCodec().Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf, []byte(`test`)) {
		t.Fatalf("bytes not equal %s != %s", buf, `test`)
	}
}

func TestFrameFlatten(t *testing.T) {
	s := &struct {
		One  string
		Name *codec.Frame `json:"name" codec:"flatten"`
	}{
		One:  "xx",
		Name: &codec.Frame{Data: []byte("test")},
	}

	buf, err := NewCodec(codec.Flatten(true)).Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf, []byte(`test`)) {
		t.Fatalf("bytes not equal %s != %s", buf, `test`)
	}
}

func TestStructByTag(t *testing.T) {
	type Str struct {
		Name []string `json:"name" codec:"flatten"`
	}

	val := &Str{Name: []string{"first", "second"}}

	c := NewCodec(codec.Flatten(true))
	buf, err := c.Marshal(val)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf, []byte(`["first","second"]`)) {
		t.Fatalf("invalid marshal: %s != %s", buf, []byte(`["first","second"]`))
	}

	err = c.Unmarshal([]byte(`["1","2"]`), val)
	if err != nil {
		t.Fatal(err)
	}

	if len(val.Name) != 2 {
		t.Fatalf("invalid unmarshal: %v", val)
	}
}
