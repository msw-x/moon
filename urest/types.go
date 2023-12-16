package urest

type Void struct{}

type Text string

func (o *Text) Set(s string) {
	*o = Text(s)
}

type Image struct {
	Type string
	Data []byte
}
