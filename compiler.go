package compiler

import (
	"errors"
	"fmt"
	"github.com/gdtrp/brainfuck/stack"
	"io"
)

type Compiler struct {
	commands map[stack.Command]stack.ExternalOperation
}

func (c *Compiler) registerOperation(operation stack.ExternalOperation) error {
	if _, ok := c.commands[operation.Command()]; ok {
		return errors.New(fmt.Sprintf("operation %v already present in the supported commands list", operation.Command()))
	}
	c.commands[operation.Command()] = operation
	return nil
}

/*
compile provided script. read byte data from reader and write outgoing bytes to writer. all unsupported tokens will be ignored
*/
func (c Compiler) Compile(script io.Reader, reader io.Reader, writer io.Writer) error {
	context := stack.NewContext(reader, writer)
	token := make([]byte, 1)
	for {
		if _, err := script.Read(token); err == nil {
			if operation, found := c.commands[stack.Command(token)]; found {
				if err := context.Execute(operation); err != nil {
					return err
				}
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		} else {
			return err
		}

	}
	return context.ValidateExecution()
}

/*
create new compiler. additional operations can also be provided. command name overlapping is not allowed
*/
func New(ops ...stack.ExternalOperation) (Compiler, error) {
	result := Compiler{
		commands: make(map[stack.Command]stack.ExternalOperation),
	}
	for _, o := range stack.GetDefaultOperations() {
		if err := result.registerOperation(o); err != nil {
			return result, err
		}
	}
	for _, o := range ops {
		if err := result.registerOperation(o); err != nil {
			return result, err
		}
	}
	return result, nil
}
