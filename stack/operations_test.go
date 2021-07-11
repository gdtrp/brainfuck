package stack

import (
	"bytes"
	"testing"
)

func TestIncrOperation(t *testing.T) {
	context := Context{
		Memory:     []byte{0, 1, 2, 3, 4},
		CurrentIdx: 2,
	}
	error := incr.action(&context)

	if error != nil {
		t.Error("error")
	}
	if context.CurrentIdx != 2 {
		t.Error("current idx shouldn't been changed")
	}
	if bytes.Compare(context.Memory, []byte{0, 1, 3, 3, 4}) != 0 {
		t.Error("arrays aren't equal")
	}

	context = Context{
		Memory:     []byte{0, 1, 2, 3, 4},
		CurrentIdx: 6,
	}
}

func TestDecrOperation(t *testing.T) {
	context := Context{
		Memory:     []byte{0, 1, 2, 3, 4},
		CurrentIdx: 2,
	}
	error := decr.action(&context)
	if error != nil {
		t.Error("error")
	}

	if context.CurrentIdx != 2 {
		t.Error("current idx shouldn't been changed")
	}
	if bytes.Compare(context.Memory, []byte{0, 1, 1, 3, 4}) != 0 {
		t.Error("arrays aren't equal")
	}
	context = Context{
		Memory:     []byte{0, 1, 2, 3, 4},
		CurrentIdx: 0,
	}
	error = decr.action(&context)
	if error != nil {
		t.Error("error")
	}

	if bytes.Compare(context.Memory, []byte{255, 1, 2, 3, 4}) != 0 {
		t.Error("arrays aren't equal")
	}
}

func TestIncrPointerOperation(t *testing.T) {
	context := Context{
		Memory:     []byte{0, 1, 2, 3, 4},
		CurrentIdx: 2,
	}
	error := ip.action(&context)

	if context.CurrentIdx != 3 {
		t.Error("index didn't change")
	}

	if error != nil || bytes.Compare(context.Memory, []byte{0, 1, 2, 3, 4}) != 0 {
		t.Error("arrays aren't equal")
	}

	context = Context{
		Memory:     []byte{0, 1, 2, 3, 4},
		CurrentIdx: 5,
	}
	error = ip.action(&context)
	if error == nil {
		t.Error("error should be returned")
	}
}

func TestDecrPointerOperation(t *testing.T) {
	context := Context{
		Memory:     []byte{0, 1, 2, 3, 4},
		CurrentIdx: 2,
	}
	error := dp.action(&context)

	if context.CurrentIdx != 1 {
		t.Error("index didn't change")
	}

	if error != nil || bytes.Compare(context.Memory, []byte{0, 1, 2, 3, 4}) != 0 {
		t.Error("arrays aren't equal")
	}

	context = Context{
		Memory:     []byte{0, 1, 2, 3, 4},
		CurrentIdx: 0,
	}
	error = dp.action(&context)
	if error == nil {
		t.Error("error should be returned")
	}
}

func TestPrintPointerOperation(t *testing.T) {
	writer := bytes.NewBuffer([]byte{})
	context := Context{
		Memory:     []byte{0, 1, 2, 3, 4},
		CurrentIdx: 2,
		Writer:     writer,
	}
	error := output.action(&context)
	result := writer.Bytes()
	if len(result) != 1 || error != nil || result[0] != 2 {
		t.Error("wrong result")
	}
	if context.CurrentIdx != 2 {
		t.Error("index changed")
	}

	if bytes.Compare(context.Memory, []byte{0, 1, 2, 3, 4}) != 0 {
		t.Error("arrays aren't equal")
	}

}

func TestSetPointerOperation(t *testing.T) {
	reader := bytes.NewBuffer([]byte{6})
	context := Context{
		Memory:     []byte{0, 1, 2, 3, 4},
		CurrentIdx: 2,
		Reader:     reader,
	}
	error := input.action(&context)

	if error != nil || context.Memory[2] != 6 {
		t.Error("wrong result")
	}
	if context.CurrentIdx != 2 {
		t.Error("index changed")
	}

	if bytes.Compare(context.Memory, []byte{0, 1, 6, 3, 4}) != 0 {
		t.Error("arrays aren't equal")
	}
}
