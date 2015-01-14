package lua

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type VM struct {
	globalEnv *Environment
}

type Process struct {
	localEnv *Environment
	stack    *Stack
	running  bool
}

func (v *VM) ExecuteString(code string) {
	tokens := Tokenize(strings.NewReader(code))
	v.executeTokens(tokens)
}

func (v *VM) ExecuteFile(file *os.File) {
	defer file.Close()
	tokens := Tokenize(file)
	v.executeTokens(tokens)
}

func (v *VM) GetGlobalVariable(key string) interface{} {
	return v.globalEnv.SearchValue(key)
}

func (v *VM) SetGlobalVariable(key string, value interface{}) {
	v.globalEnv.SetValue(key, value)
}

func New() *VM {
	return &VM{
		NewEnv(nil),
	}
}

func (v *VM) executeTokens(tokens []Token) {
	instructions := CreateBytecode(tokens)
	v.executeInstructions(instructions)
}

func (v *VM) newProcess() *Process {
	return &Process{
		NewEnv(v.globalEnv),
		NewStack(),
		false,
	}
}

func (vm *VM) executeInstructions(instructions *[]Statement) interface{} {
	proc := vm.newProcess()
	proc.running = true
	currentEnvironment := proc.localEnv

	for _, instruc := range *instructions {

		if !proc.running {
			return nil
		}

		switch v := instruc.(type) {
		case VariableAssignment:
			fmt.Printf("ASSIGN VARIABLE %v to %v\n", v.name, proc.stack.Peek())
			if v.local {
				currentEnvironment.SetValue(v.name, proc.stack.Pop())
			} else {
				vm.globalEnv.SetValue(v.name, proc.stack.Pop())
			}
		case PushValueStack:
			proc.stack.Push(v.value)
		case PushVariableStack:
			fmt.Printf("PUSH VARAIBLE %v VALUE %v\n", v.name, currentEnvironment.SearchValue(v.name))

			proc.stack.Push(currentEnvironment.SearchValue(v.name))
		case ReturnValue:
			if currentEnvironment == proc.localEnv {
				return proc.stack.Pop()
			}
		case AddOperation:
			proc.stack.Push(proc.stack.PopFloat64() + proc.stack.PopFloat64())
		case SubOperation:
			secondNum := proc.stack.PopFloat64()
			firstNum := proc.stack.PopFloat64()
			proc.stack.Push(firstNum - secondNum)
		case MulOperation:
			proc.stack.Push(proc.stack.PopFloat64() * proc.stack.PopFloat64())
		case DivOperation:
			secondNum := proc.stack.PopFloat64()
			firstNum := proc.stack.PopFloat64()
			proc.stack.Push(firstNum / secondNum)
		default:
			fmt.Printf("INVALID INSTRUCTION %v\n", reflect.TypeOf(v))
		}
	}

	return nil
}
