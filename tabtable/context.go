package tabtable

type Context struct {
	MinWidth int
	TabWidth int
	Padding  int
	PadChar  byte
	Flags    uint
}

func DefaultContext() Context {
	return Context{
		Padding: 2,
		PadChar: ' ',
	}
}
