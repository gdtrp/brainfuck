package stack

import "errors"

type LinkedElement interface {
	Next() LinkedElement
	Previous() LinkedElement
	HasMoreElements() bool
	setNext(LinkedElement)
	setPrevious(LinkedElement)
	//configure links of current element with stack
	ConfigureLink(*Stack)
	//rewind to first element in stack or first element of current loop
	RewindToStart() LinkedElement
	//return closest operation
	CurrentOperation() OperationalElement
	//return link to current executing loop
	CurrentLoop() LoopElement
	//return link to upper level nested loop
	GetPreviousLoop() LoopElement
}
type OperationalElement interface {
	LinkedElement
	Operation() ExternalOperation
}

type LoopElement interface {
	LinkedElement
	firstElement() OperationalElement
	setFirstElement(OperationalElement)
}

//Stack contains all execution state
type Stack struct {
	//link to next element to execute
	nextElement LinkedElement
	//link to loop needs to be skipped.
	//if set all subsequent operations will not be executed until loop will not be closed
	skip LoopElement
	//link to current executed element
	current LinkedElement
	//link to current executing loop
	currentLoop LoopElement
	//link to last added element to stack.
	lastAdded LinkedElement
}

//High level struct contains links to previous and next elements in execution order
type Link struct {
	next     LinkedElement
	previous LinkedElement
}
type OperationContainer struct {
	Link
	//Operation to execute
	operation ExternalOperation
	//Current loop link
	loop LoopElement
}

//LoopContainer struct
type LoopContainer struct {
	Link
	//Link to first loop element
	firstLoopElement OperationalElement
}

func (c *OperationContainer) Loop() LoopElement {
	return c.loop
}
func (c *OperationContainer) CurrentLoop() LoopElement {
	return c.loop
}
func (c *LoopContainer) CurrentLoop() LoopElement {
	return c
}
func (c *OperationContainer) Operation() ExternalOperation {
	return c.operation
}

func (c *OperationContainer) CurrentOperation() OperationalElement {
	return c
}

func (c *OperationContainer) ConfigureLink(stack *Stack) {
	c.linkWithLoop(stack.currentLoop)
	linkPrevious(stack.lastAdded, c)
	stack.nextElement = c
	stack.lastAdded = c
}

func (c *OperationContainer) linkWithLoop(element LoopElement) {
	if element != nil {
		c.loop = element
		if element.firstElement() == nil {
			element.setFirstElement(c)
		}
	}
}
func linkPrevious(element LinkedElement, c LinkedElement) {
	if element != nil {
		c.setPrevious(element)
		element.setNext(c)
	}
}
func (c *LoopContainer) ConfigureLink(stack *Stack) {
	linkPrevious(stack.lastAdded, c)
	stack.currentLoop = c
	stack.lastAdded = nil
}

func (c *LoopContainer) firstElement() OperationalElement {
	return c.firstLoopElement
}

func (c *LoopContainer) setFirstElement(element OperationalElement) {
	c.firstLoopElement = element
}
func (c *LoopContainer) CurrentOperation() OperationalElement {
	return c.firstElement()
}
func (c *Link) GetPreviousLoop() LoopElement {
	prev := c.Previous()
	if prev != nil {
		return prev.CurrentLoop()
	}
	return nil
}

func (c *Link) Next() LinkedElement {
	return c.next
}
func (c *Link) setNext(element LinkedElement) {
	c.next = element
}
func (c *Link) Previous() LinkedElement {
	return c.previous
}
func (c *Link) RewindToStart() LinkedElement {
	loop := c.Previous()
	for loop.Previous() != nil {
		loop = loop.Previous()
	}
	return loop
}
func (c *Link) setPrevious(element LinkedElement) {
	c.previous = element
}

func (c *Link) HasMoreElements() bool {
	return c.Next() != nil
}

func newStack() *Stack {
	return &Stack{}
}

func (s *Stack) addToStack(operation ExternalOperation) {
	newOp := &OperationContainer{operation: operation}
	newOp.ConfigureLink(s)
}
func (s *Stack) addLoopContainer() {
	newOp := &LoopContainer{}
	newOp.ConfigureLink(s)
}

//check if has next element in stack
func (s *Stack) hasNext() bool {
	return s.nextElement != nil
}

//retrieve element from stack and set next
func (s *Stack) next() ExternalOperation {
	current := s.nextElement.CurrentOperation()
	result := current.Operation()
	s.current = current
	s.nextElement = current.Next()
	return result
}

//add new operation to stack
func (s *Stack) prepareStack(operation ExternalOperation, ctx *Context) error {
	internal, ok := operation.(internalOperation)
	if ok && internal.OnAdd() != nil {
		if error := internal.OnAdd()(ctx); error != nil {
			return error
		}
	}
	s.addToStack(operation)
	if ok && internal.AfterAdd() != nil {
		if error := internal.AfterAdd()(ctx); error != nil {
			return error
		}
	}
	return nil
}

//add new operation to stack and execute
func (s *Stack) execute(operation ExternalOperation, ctx *Context) error {
	if error := s.prepareStack(operation, ctx); error != nil {
		return error
	}
	//specific case for loops which needs to be added but without execution (covers excludes look-ahead requirement)
	if !s.isSkipExecution() {
		for s.hasNext() {
			op := s.next()
			if err := op.Action()(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

//add loop element to stack and set it current
func (s *Stack) initLoop() {
	s.addLoopContainer()
}

//mark current loop as finished. returns error if initLoop method wasn't called
func (s *Stack) terminateLoop() error {
	if s.currentLoop == nil {
		return errors.New("missing start loop")
	}
	s.lastAdded = s.currentLoop
	if s.skip == s.currentLoop {
		s.skip = nil
	}
	s.currentLoop = s.currentLoop.GetPreviousLoop()
	return nil
}

//rewind loop to beginning
func (s *Stack) endLoop() {
	s.nextElement = s.current.RewindToStart()
}
func (s *Stack) isSkipExecution() bool {
	return s.skip != nil
}

//ends loop execution and jump to next element.
//if next element is empty that means that the loop execution needs to be skipped,
//because current loop is not read fully (covers excludes look-ahead requirement)
func (s *Stack) breakLoop() {
	s.nextElement = s.current.CurrentLoop().Next()
	if !s.current.HasMoreElements() {
		s.skip = s.current.CurrentLoop()
	}
}

func (s *Stack) validateExecution(c *Context) error {
	if s.currentLoop != nil {
		return errors.New("missing closing brackets")
	}
	return nil
}
