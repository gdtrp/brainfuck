package compiler

import (
	"github.com/gdtrp/brainfuck/stack"
	"bytes"
	"testing"
)

var scripts = []struct {
	name   string
	script string
	result []byte
	input  []byte
}{

	{"test empty loop execution", `[---[->-[[>++]]]]++.`, []byte{2}, nil},
	{"simple script with increment and output", "+++++++++.", []byte{9}, nil},

	{"simple script with newlines and tabs", "++	\n++\n++\n+++.", []byte{9}, nil},

	{"simple script with input", "+++.>,++.", []byte{3, 9}, []byte{7}},

	{"jabh", `+++[>+++++<-]>[>+>+++>+>++>+++++>++<[++<]>---]>->-.[>++>+<<--]>--.--.+.>
   >>++.<<.<------.+.+++++.>>-.<++++.<--.>>>.<<---.<.-->-.>+.[+++++.---<]>>
   [.--->]<<.<+.++.++>+++[.<]`, []byte("Just another brainfuck hacker"), nil},

	{"count lines words and chars", `>>>+>>>>>+>>+>>+[<<],[
      	    -[-[-[-[-[-[-[-[<+>-[>+<-[>-<-[-[-[<++[<++++++>-]<
      	        [>>[-<]<[>]<-]>>[<+>-[<->[-]]]]]]]]]]]]]]]]
      	    <[-<<[-]+>]<<[>>>>>>+<<<<<<-]>[>]>>>>>>>+>[
      	        <+[
      	            >+++++++++<-[>-<-]++>[<+++++++>-[<->-]+[+>>>>>>]]
      	            <[>+<-]>[>>>>>++>[-]]+<
      	        ]>[-<<<<<<]>>>>
      	    ],
      	]+<++>>>[[+++++>>>>>>]<+>+[[<++++++++>-]<.<<<<<]>>>>>>>>]`, []byte("\t1\t4\t22\n"), []byte("test1 test2\ntest3 tttt")},
	{"simple nested loop", "++++++++++[+[>+++<-]]>.", []byte("!"), nil},
	{"nested loop", `+++[>+++++<-]>[>+>+++>+>++>+++++>++<[++<]>---]>.>.>.`, []byte("-K-"), nil},
}

func TestCompiler_Scripts(t *testing.T) {
	compiler, error := New()
	if error != nil {
		t.Fatalf("error should be nil")
	}
	for _, test := range scripts {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			if error := compiler.Compile(bytes.NewBuffer([]byte(test.script)), bytes.NewBuffer(test.input), &buf); error != nil {
				t.Fatalf("unexpected error %v", error)
			}
			result := buf.Bytes()
			if len(result) != len(test.result) {
				t.Fatalf("wrong result value expected %v but was %v", test.result, result)
			}
			for i, v := range test.result {
				if result[i] != v {
					t.Fatalf("wrong result value expected %v but was %v", test.result, result)
				}
			}
		})
	}

}

type CustomOperation struct {
	command stack.Command
	action  func(ctx *stack.Context) error
}

func (c CustomOperation) Command() stack.Command {
	return c.command
}
func (c CustomOperation) Action() func(ctx *stack.Context) error {
	return c.action
}
func TestCompilerWithCustomCommand(t *testing.T) {
	compiler, error := New(CustomOperation{command: "*", action: func(ctx *stack.Context) error {
		ctx.SetCurrentByte(ctx.GetCurrentByte() * 2)
		return nil
	}})
	var buf bytes.Buffer
	if error != nil {
		t.Fatalf("error should be nil")
	}
	if error = compiler.Compile(bytes.NewBuffer([]byte("++*.>+++*.")), nil, &buf); error != nil {
		t.Fatalf("error should be nil")
	}
	result := buf.Bytes()
	if result[0] != 4 || result[1] != 6 {
		t.Fatalf("wrong value, expected %v but was %v", []byte{4, 6}, result)
	}
}

func TestCompilerWithSameCustomCommand(t *testing.T) {
	_, error := New(CustomOperation{command: "*", action: func(ctx *stack.Context) error {
		ctx.SetCurrentByte(ctx.GetCurrentByte() * 2)
		return nil
	}}, CustomOperation{command: "*", action: func(ctx *stack.Context) error {
		ctx.SetCurrentByte(ctx.GetCurrentByte() * 2)
		return nil
	}})
	if error == nil {
		t.Fatalf("error shouldn't be nil")
	}
}

func TestCompilerWithExistingCustomCommand(t *testing.T) {
	_, error := New(CustomOperation{command: "+", action: func(ctx *stack.Context) error {
		ctx.SetCurrentByte(ctx.GetCurrentByte() * 2)
		return nil
	}})
	if error == nil {
		t.Fatalf("error shouldn't be nil")
	}
}

var errorScripts = []struct {
	name   string
	script string
	input  []byte
}{

	{"mismatched closing bracket", `[[[[[-]-]`, nil},
	{"mismatched closing bracket", `[--`, nil},
	{"mismatched opening bracket", `[--]]`, nil},
	{"overflow", `+[>+]`, nil},
}

func TestCompiler_ErrorScript(t *testing.T) {
	compiler, error := New()
	if error != nil {
		t.Fatalf("error should be nil")
	}
	for _, test := range errorScripts {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			if error := compiler.Compile(bytes.NewBuffer([]byte(test.script)), bytes.NewBuffer(test.input), &buf); error == nil {
				t.Fatalf("error must be present")
			}
		})
	}

}
