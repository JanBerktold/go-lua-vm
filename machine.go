package lua

import (
	"fmt"
	"os"
	"strings"
)

type VM struct {
	globalEnv *Environment
}

type Process struct {
	localEnv *Environment
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
		false,
	}
}

func (vm *VM) executeInstructions(instructions *[]Statement) {
	proc := vm.newProcess()
	proc.running = true
	currentEnvironment := proc.localEnv

	for _, instruc := range *instructions {

		if !proc.running {
			return
		}

		switch v := instruc.(type) {
		case VariableAssignment:
			if v.local {
				currentEnvironment.SetValue(v.name, v.value)
			} else {
				vm.globalEnv.SetValue(v.name, v.value)
			}
		default:
			fmt.Println("INVALID INSTRUCTION")
			os.Exit(1)
		}
	}

}
