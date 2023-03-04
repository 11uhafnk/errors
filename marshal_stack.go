package errors

var (
	stackSourceFileName     = "source"
	stackSourceLineName     = "line"
	stackSourceFunctionName = "func"
)

type state struct {
	b []byte
}

// Write implement fmt.Formatter interface.
func (s *state) Write(b []byte) (n int, err error) {
	s.b = b
	return len(b), nil
}

// Width implement fmt.Formatter interface.
func (s *state) Width() (wid int, ok bool) {
	return 0, false
}

// Precision implement fmt.Formatter interface.
func (s *state) Precision() (prec int, ok bool) {
	return 0, false
}

// Flag implement fmt.Formatter interface.
func (s *state) Flag(c int) bool {
	return false
}

func frameField(f Frame, s *state, c rune) string {
	f.Format(s, c)
	return string(s.b)
}

// ZerologMarshalStack implements pkg/errors stack trace marshaling.
//
//	zerolog.ErrorStackMarshaler = errors.ZerologMarshalStack
func ZerologMarshalStack(err error) interface{} {
	var st stackTracer
	if !As(err, &st) {
		return nil
	}

	sterr, ok := err.(stackTracer)
	for !ok {
		err = Unwrap(err)
		sterr, ok = err.(stackTracer)
	}
	trace := sterr.StackTrace()
	s := &state{}
	out := make([]map[string]string, 0, len(trace))
	for _, frame := range trace {
		out = append(out, map[string]string{
			stackSourceFileName:     frameField(frame, s, 's'),
			stackSourceLineName:     frameField(frame, s, 'd'),
			stackSourceFunctionName: frameField(frame, s, 'n'),
		})
	}
	return out
}
