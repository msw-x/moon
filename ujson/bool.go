package ujson

import "github.com/msw-x/moon/parse"

type Bool bool

func (o *Bool) UnmarshalJSON(b []byte) error {
	v, err := parse.Bool(unquote(b))
	*o = Bool(v)
	return err
}

func (o Bool) Value() bool {
	return bool(o)
}
