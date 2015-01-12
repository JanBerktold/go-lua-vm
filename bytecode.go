package lua

// File is responsible for taking a list of tokens and transforming it into a list of executable statements,
// which will be passed towards the VM at execution stage. Bytecode is not compatible with other implementations.

import (
	"fmt"
)

type Statement interface{}

type FunctionCall struct {
	name string
	args int
}

type VariableAssignment struct {
	name  string
	local bool
}

type PushVariableStack struct {
	name string
}

type PushValueStack struct {
	value interface{}
}

type ReturnValue struct {
	amount int
}

type AddOperation struct {
}

type SubOperation struct {
}

type TableLengthOperation struct {
}

func getFuncEnd(tok *[]Token, n int) int {
	level := 1
	for level > 0 && n < len(*tok) {
		n++

		if (*tok)[n].typ == keyword_token {

			if (*tok)[n].value == "end" {
				level--
			} else if scopeMap[(*tok)[n].value] {
				level++
			}
		}

	}
	return n
}

func getFuncParamEnd(tok *[]Token, n int) int {
	level := 1
	for level > 0 && n+1 < len(*tok) {
		n++

		if (*tok)[n].typ == keyword_token {

			if (*tok)[n].value == ")" {
				level--
			} else if (*tok)[n].value == "(" {
				level++
			}
		}

	}
	return n
}

func AddStatement(tar *[]Statement, counter *int, stat Statement) {
	for (*tar)[*counter] != nil {
		(*counter)++
	}
	(*tar)[*counter] = stat
	(*counter)++
}

func PushValue(tar *[]Statement, counter *int, value interface{}) {
	AddStatement(tar, counter, PushValueStack{
		value,
	})
}

func PushVariable(tar *[]Statement, counter *int, name string) {
	AddStatement(tar, counter, PushVariableStack{
		name,
	})
}

func VariableAssign(tar *[]Statement, counter *int, name string, local bool) {
	AddStatement(tar, counter, VariableAssignment{
		name,
		local,
	})
}

func EvaluateStatement(tar *[]Statement, counter *int, tok []Token) (endAt int) {
	count := 0
	lastType := -1
	firstArith := true
	for count < len(tok) {

		switch tok[count].typ {
		case language_token:
			var stat Statement

			if tok[count].value == "+" {
				if firstArith {
					fmt.Println("RESET")
					(*counter)++
				}
				AddStatement(tar, counter, AddOperation{})
				if firstArith {
					firstArith = false
					(*counter)--
					(*counter)--
				}
			} else if tok[count].value == "-" {
				stat = SubOperation{}
			} else if tok[count].value == "#" {
				stat = TableLengthOperation{}
			}

			if stat != nil {
				AddStatement(tar, counter, stat)
			}
		case identifier_token:
			if lastType >= identifier_token {
				return count - 1
			}
			PushVariable(tar, counter, tok[count].value)
		case number_token, string_token:
			if lastType >= identifier_token {
				return count - 1
			}
			PushValue(tar, counter, tok[count].real)
		}

		lastType = tok[count].typ
		count++
	}
	return count
}

func CreateBytecode(tok []Token) *[]Statement {
	result := make([]Statement, 100, 1000)

	currentStatement := 0
	currentToken := 0

	for currentToken < len(tok) {
		token := tok[currentToken]

		if token.typ == language_token && token.value == "=" {
			EvaluateStatement(&result, &currentStatement, tok[currentToken+1:len(tok)])
			VariableAssign(&result, &currentStatement, tok[currentToken-1].value, currentToken - 2 >= 0 && tok[currentToken - 2].value == "local")
		}

		currentToken++
	}

	realResult := result[0:currentStatement]
	return &realResult
}
