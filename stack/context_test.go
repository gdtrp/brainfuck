package stack

import (
	"bytes"
	"errors"
	"testing"
)

func TestSetCurrentByte(t *testing.T) {
	var input []byte
	var output []byte

	ctx := NewContextWithMemorySize(bytes.NewBuffer(input), bytes.NewBuffer(output), 5)
	ctx.SetCurrentByte(6)
	if len(ctx.Memory) != 5 {
		t.Errorf("wrong memory size should be 1 but was %v", ctx.Memory[0])
	}
	if ctx.Memory[0] != 6 {
		t.Errorf("wrong byte value should be 6 but was %v", ctx.Memory[0])
	}

	ctx.SetIndex(4)
	ctx.SetCurrentByte(6)

	if ctx.Memory[4] != 6 {
		t.Errorf("wrong byte value should be 6 but was %v", ctx.Memory[0])
	}
}
func TestSetByte(t *testing.T) {

	var input []byte
	var output []byte

	ctx := NewContextWithMemorySize(bytes.NewBuffer(input), bytes.NewBuffer(output), 5)
	if error := ctx.SetByte(3, 8); error != nil {
		t.Fatalf("unexpected error")
	}
	if ctx.Memory[3] != 8 {
		t.Errorf("wrong byte value should be 8 but was %v", ctx.Memory[0])
	}
	if error := ctx.SetByte(6, 9); error == nil {
		t.Errorf("error expected")
	}
}
func TestSetIndex(t *testing.T) {

	var input []byte
	var output []byte

	ctx := NewContextWithMemorySize(bytes.NewBuffer(input), bytes.NewBuffer(output), 5)
	if error := ctx.SetIndex(4); error != nil {
		t.Fatalf("unexpected error")
	}
	if ctx.CurrentIdx != 4 {
		t.Errorf("wrong index value should be 4 but was %v", ctx.CurrentIdx)
	}
	if error := ctx.SetIndex(6); error == nil {
		t.Errorf("error expected")
	}
}
func TestGetIndex(t *testing.T) {

	var input []byte
	var output []byte

	ctx := NewContextWithMemorySize(bytes.NewBuffer(input), bytes.NewBuffer(output), 5)
	if error := ctx.SetIndex(4); error != nil {
		t.Fatalf("unexpected error")
	}
	if ctx.GetIndex() != 4 {
		t.Errorf("wrong index value should be 4 but was %v", ctx.GetIndex())
	}

}
func TestGetByte(t *testing.T) {

	var input []byte
	var output []byte

	ctx := NewContextWithMemorySize(bytes.NewBuffer(input), bytes.NewBuffer(output), 5)
	ctx.SetByte(4, 5)
	if val, err := ctx.GetByte(4); val != 5 || err != nil {
		t.Errorf("wrong byte value should be 5 but was %v", val)
	}
	if _, err := ctx.GetByte(6); err == nil {
		t.Errorf("error expected")
	}
}
func TestGetCurrentByte(t *testing.T) {

	var input []byte
	var output []byte

	ctx := NewContextWithMemorySize(bytes.NewBuffer(input), bytes.NewBuffer(output), 5)

	if ctx.GetCurrentByte() != 0 {
		t.Errorf("wrong byte value should be 0 but was %v", ctx.GetCurrentByte())
	}
	ctx.SetByte(4, 5)
	ctx.SetIndex(4)
	if ctx.GetCurrentByte() != 5 {
		t.Errorf("wrong byte value should be 5 but was %v", ctx.GetCurrentByte())
	}
}
func TestExecute(t *testing.T) {

	var input []byte
	var output []byte

	ctx := NewContextWithMemorySize(bytes.NewBuffer(input), bytes.NewBuffer(output), 5)
	executed := false
	ctx.Execute(operation{
		token: "test",
		action: func(context *Context) error {
			executed = true
			return nil
		},
	})
	if !executed {
		t.Errorf("action wasn't executed")
	}
	error := ctx.Execute(operation{
		token: "test",
		action: func(context *Context) error {
			return errors.New("test error")
		},
	})
	if error == nil {
		t.Errorf("error expected")
	}
}
