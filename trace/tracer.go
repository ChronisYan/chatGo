package trace

type Tracer interface {
	// The ...interface{} argument type means that Trace accepts zero or more arguments of any type
	Trace(...interface{})
}
