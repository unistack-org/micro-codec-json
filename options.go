package json

import (
	codec "github.com/unistack-org/micro/v3/codec"
)

type unmarshalOptionsKey struct{}

func UnmarshalOptions(o JsonUnmarshalOptions) codec.Option {
	return codec.SetOption(unmarshalOptionsKey{}, o)
}

type marshalOptionsKey struct{}

func MarshalOptions(o JsonUnmarshalOptions) codec.Option {
	return codec.SetOption(marshalOptionsKey{}, o)
}
