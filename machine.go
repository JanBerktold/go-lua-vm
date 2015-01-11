package lua

import (
	"os"
	"strings"
)

type VM struct {

}

type Process struct {
	globalEnv *Environment
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

func (v *VM) executeTokens(tokens []Token) {

}
