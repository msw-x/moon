package ujson

import (
	"encoding/json"
)

type Optional[T any] struct {
	Has   bool
	Value *T
}

func (o *Optional[T]) UnmarshalJSON(b []byte) error {
	o.Has = true
	if string(b) == "null" {
		return nil
	}
	o.Value = new(T)
	return json.Unmarshal(b, o.Value)
}
