package stack

import (
	"io"
)

//internal interface for system commands
type internalOperation interface {
	ExternalOperation
	//action will be executed before adding operation to stack. Can be nil
	OnAdd() func(*Context) error
	//action will be executed after adding operation to stack. Can be nil
	AfterAdd() func(*Context) error
}

type operation struct {
	token    Command
	action   func(*Context) error
	onAdd    func(*Context) error
	afterAdd func(*Context) error
}

func (o operation) Command() Command {
	return o.token
}
func (o operation) Action() func(*Context) error {
	return o.action
}
func (o operation) OnAdd() func(*Context) error {
	return o.onAdd
}
func (o operation) AfterAdd() func(*Context) error {
	return o.afterAdd
}

func GetDefaultOperations() []operation {
	return []operation{
		incr, decr, ip, dp, output, input, startLoop, endLoop,
	}
}

//increment operation
var incr = operation{
	token: "+",
	action: func(ctx *Context) error {
		return ctx.SetCurrentByte(ctx.GetCurrentByte() + 1)
	},
}

//decrement operation
var decr = operation{
	token: "-",
	action: func(ctx *Context) error {
		return ctx.SetCurrentByte(ctx.GetCurrentByte() - 1)
	},
}

//increase current index operation
var ip = operation{
	token: ">",
	action: func(ctx *Context) error {
		return ctx.SetIndex(ctx.GetIndex() + 1)
	},
}

//decrease current index operation
var dp = operation{
	token: "<",
	action: func(ctx *Context) error {
		return ctx.SetIndex(ctx.GetIndex() - 1)
	},
}

//print current index value operation
var output = operation{
	token: ".",

	action: func(ctx *Context) error {
		_, error := ctx.Writer.Write([]byte{ctx.GetCurrentByte()})
		return error
	},
}

//read current output byte to the current index
var input = operation{
	token: ",",
	action: func(ctx *Context) error {
		b := make([]byte, 1)
		if _, error := ctx.Reader.Read(b); error == nil {
			return ctx.SetCurrentByte(b[0])
		} else {
			if error == io.EOF {
				return nil
			}
			return error
		}
	},
}

//start loop operation
var startLoop = operation{
	token: "[",

	onAdd: func(ctx *Context) error {
		ctx.Stack.initLoop()
		return nil
	},
	action: func(ctx *Context) error {
		if ctx.GetCurrentByte() == 0 {
			ctx.Stack.breakLoop()
		}
		return nil
	},
}

//end loop operation
var endLoop = operation{
	token: "]",
	afterAdd: func(ctx *Context) error {
		return ctx.Stack.terminateLoop()
	},
	action: func(ctx *Context) error {
		ctx.Stack.endLoop()
		return nil
	},
}
