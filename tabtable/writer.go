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

func (this *Writer) Init(ctx Context) {
	this.writer = new(tabwriter.Writer)
	this.buf = bytes.NewBuffer(nil)
	this.bufwriter = bufio.NewWriter(this.buf)
	this.writer.Init(this.bufwriter, ctx.MinWidth, ctx.TabWidth, ctx.Padding, ctx.PadChar, ctx.Flags)
}

func (this *Writer) Write(a ...any) {
	s := ""
	for _, v := range a {
		s += fmt.Sprint(v) + "\t"
	}
	s += "\n"
	fmt.Fprint(this.writer, s)
}

func (this *Writer) Writef(format string, a ...any) {
	format += "\n"
	fmt.Fprintf(this.writer, format, a...)
}

func (this *Writer) String() string {
	this.writer.Flush()
	this.bufwriter.Flush()
	s := fmt.Sprint(this.buf)
	lines := strings.Split(s, "\n")
	for n, line := range lines {
		lines[n] = ustring.TrimBackWhitespaces(line)
	}
	s = strings.Join(lines, "\n")
	return s
}

func (this *Writer) Print() {
	fmt.Print(this.String())
}
