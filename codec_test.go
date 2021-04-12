package json

import (
	"bytes"
	"testing"
)

func TestReadBody(t *testing.T) {
	s := &struct {
		Name string
	}{}
	c := NewCodec()
	b := bytes.NewReader(nil)
	err := c.ReadBody(b, s)
	if err != nil {
		t.Fatal(err)
	}
}
