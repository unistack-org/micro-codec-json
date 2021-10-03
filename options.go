package json

import (
	codec "go.unistack.org/micro/v3/codec"
)

type unmarshalOptionsKey struct{}

func UnmarshalOptions(o JsonUnmarshalOptions) codec.Option {
	return codec.SetOption(unmarshalOptionsKey{}, o)
}

type marshalOptionsKey struct{}

func MarshalOptions(o JsonMarshalOptions) codec.Option {
	return codec.SetOption(marshalOptionsKey{}, o)
}
