package tabtable

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/msw-x/moon/ustring"
)

type Writer struct {
	writer    *tabwriter.Writer
	buf       *bytes.Buffer
	bufwriter *bufio.Writer
}

func New() *Writer {
	w := &Writer{}
	w.Init(DefaultContext())
	return w
}

func (o *Writer) Init(ctx Context) {
	o.writer = new(tabwriter.Writer)
	o.buf = bytes.NewBuffer(nil)
	o.bufwriter = bufio.NewWriter(o.buf)
	o.writer.Init(o.bufwriter, ctx.MinWidth, ctx.TabWidth, ctx.Padding, ctx.PadChar, ctx.Flags)
}

func (o *Writer) Write(a ...any) {
	s := ""
	for _, v := range a {
		s += fmt.Sprint(v) + "\t"
	}
	s += "\n"
	fmt.Fprint(o.writer, s)
}

func (o *Writer) Writef(format string, a ...any) {
	format += "\n"
	fmt.Fprintf(o.writer, format, a...)
}

func (o *Writer) String() string {
	o.writer.Flush()
	o.bufwriter.Flush()
	s := fmt.Sprint(o.buf)
	s = ustring.TransformLines(s, ustring.TrimBackWhitespaces)
	s = strings.TrimSuffix(s, "\n")
	return s
}

func (o *Writer) Print() {
	fmt.Print(o.String())
}
