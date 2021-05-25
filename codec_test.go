package json

import (
	"bytes"
	"testing"
)

func TestStructByTag(t *testing.T) {
	type Str struct {
		Name []string `json:"name" codec:"flatten"`
	}

	val := &Str{Name: []string{"first", "second"}}

	c := NewCodec()
	buf, err := c.Marshal(val)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf, []byte(`["first","second"]`)) {
		t.Fatalf("invalid marshal: %s", buf)
	}

	err = c.Unmarshal([]byte(`["1","2"]`), val)
	if err != nil {
		t.Fatal(err)
	}

	if len(val.Name) != 2 {
		t.Fatalf("invalid unmarshal: %v", val)
	}

}

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
