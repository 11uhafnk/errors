package errors

import (
	"reflect"
	"testing"
	// "github.com/stretchr/testify/assert"
)

func simple() error {
	return New("simple")
}
func test_withWrap() error {
	return Wrap(simple(), "wrap")
}
func test_withMessage() error {
	return WithMessage(simple(), "msg")
}

func TestMarshalStack(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want interface{}
	}{
		{name: "nil", err: nil, want: nil},
		{name: "simple", err: simple(), want: []map[string]string{
			{"func": "simple", "line": "10", "source": "marshal_stack_test.go"},
			{"func": "TestMarshalStack", "line": "26", "source": "marshal_stack_test.go"},
			{"func": "tRunner", "line": "1576", "source": "testing.go"},
			{"func": "goexit", "line": "1598", "source": "asm_amd64.s"},
		}},
		{name: "wrap", err: test_withWrap(), want: []map[string]string{
			{"func": "simple", "line": "10", "source": "marshal_stack_test.go"},
			{"func": "test_withWrap", "line": "13", "source": "marshal_stack_test.go"},
			{"func": "TestMarshalStack", "line": "32", "source": "marshal_stack_test.go"},
			{"func": "tRunner", "line": "1576", "source": "testing.go"},
			{"func": "goexit", "line": "1598", "source": "asm_amd64.s"},
		}},
		{name: "with message", err: test_withMessage(), want: []map[string]string{
			{"func": "simple", "line": "10", "source": "marshal_stack_test.go"},
			{"func": "test_withMessage", "line": "16", "source": "marshal_stack_test.go"},
			{"func": "TestMarshalStack", "line": "39", "source": "marshal_stack_test.go"},
			{"func": "tRunner", "line": "1576", "source": "testing.go"},
			{"func": "goexit", "line": "1598", "source": "asm_amd64.s"},
		}},
		{name: "double wrap", err: Wrap(test_withWrap(), "wrap msg"), want: []map[string]string{
			{"func": "simple", "line": "10", "source": "marshal_stack_test.go"},
			{"func": "test_withWrap", "line": "13", "source": "marshal_stack_test.go"},
			{"func": "TestMarshalStack", "line": "46", "source": "marshal_stack_test.go"},
			{"func": "tRunner", "line": "1576", "source": "testing.go"},
			{"func": "goexit", "line": "1598", "source": "asm_amd64.s"},
		}},
		{name: "double message", err: WithMessage(test_withMessage(), "msg2"), want: []map[string]string{
			{"func": "simple", "line": "10", "source": "marshal_stack_test.go"},
			{"func": "test_withMessage", "line": "16", "source": "marshal_stack_test.go"},
			{"func": "TestMarshalStack", "line": "53", "source": "marshal_stack_test.go"},
			{"func": "tRunner", "line": "1576", "source": "testing.go"},
			{"func": "goexit", "line": "1598", "source": "asm_amd64.s"},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ZerologMarshalStack(tt.err)
			// assert.Equal(t, tt.want, got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZerologMarshalStack() = %v, want %v", got, tt.want)
			}
		})
	}
}
