package stack

//API interface to create custom operations
type ExternalOperation interface {
	//token value
	Command() Command
	//custom operation
	Action() func(*Context) error
}
type Command string
