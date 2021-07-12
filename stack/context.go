package stack

import (
	"errors"
	"io"
)

//Context struct contains all execution data
type Context struct {
	//current memory
	Memory []byte
	//current memory cell index
	CurrentIdx int
	//output writer
	Writer io.Writer
	//input reader
	Reader io.Reader
	//stack struct responsible for command execution order
	Stack *Stack
}

const defaultMemorySize = 65536

func NewContext(reader io.Reader, writer io.Writer) *Context {
	return NewContextWithMemorySize(reader, writer, defaultMemorySize)
}
func NewContextWithMemorySize(reader io.Reader, writer io.Writer, size int) *Context {
	return &Context{
		Memory:     make([]byte, size),
		CurrentIdx: 0,
		Writer:     writer,
		Reader:     reader,
		Stack:      &Stack{},
	}
}

//returns current memory cell index
func (c *Context) GetIndex() int {
	return c.CurrentIdx
}

//sets current memory cell index
func (c *Context) SetIndex(index int) error {
	if err := validate(index, c.Memory); err != nil {
		return err
	}
	c.CurrentIdx = index
	return nil
}

//returns byte value of memory cell index
func (c *Context) GetCurrentByte() byte {
	b, _ := c.GetByte(c.CurrentIdx)
	return b
}

//returns byte value of provided cell index
func (c *Context) GetByte(index int) (byte, error) {
	if err := validate(index, c.Memory); err != nil {
		return 0, err
	}
	return c.Memory[index], nil
}
func validate(index int, memory []byte) error {
	if index >= len(memory) || index < 0 {
		return errors.New("index is out of range")
	}
	return nil
}

//sets current cell index value
func (c *Context) SetCurrentByte(b byte) error{
	return c.SetByte(c.CurrentIdx, b)
}

//set byte value of provided cell index
func (c *Context) SetByte(index int, b byte) error {
	if err := validate(index, c.Memory); err != nil {
		return err
	}
	c.Memory[index] = b
	return nil
}

//execute next operation from stack
func (c *Context) Execute(operation ExternalOperation) error {

	internal, ok := operation.(internalOperation)
	if ok && internal.OnAdd() != nil {
		if err := internal.OnAdd()(c); err != nil {
			return err
		}
	}
	c.Stack.push(operation)

	if ok && internal.AfterAdd() != nil {
		if err := internal.AfterAdd()(c); err != nil {
			return err
		}
	}
	if !c.Stack.isSkipExecution() {
		for c.Stack.hasNext() {
			op := c.Stack.pop()
			if err := op.Action()(c); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Context) ValidateExecution() error {
	return c.Stack.validateExecution()
}
