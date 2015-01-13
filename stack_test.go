package lua

import (
	"testing"
)

func TestStackPop(t *testing.T) {
	stack := NewStack()

	for i := 0; i < 10; i++ {
		stack.Push(i)
	}

	for i := 9; i > -1; i-- {
		if stack.Len() != i+1 {
			t.Fatalf("Stack object's length is wrong: %v. Expected: %v.", stack.Len(), i+1)
		}
		ret := stack.Pop()
		if ret != i {
			t.Fatalf("Stack object returned wrong number: %v. Expected: %v.", ret, i)
		}
	}

}
