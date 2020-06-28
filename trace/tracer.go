package trace

import (
	"fmt"
	"io"
)

type Tracer interface {
	// The ...interface{} argument type means that Trace accepts zero or more arguments of any type
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}
func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

func New(w io.Writer) Tracer{
	return &tracer{out: w}
}