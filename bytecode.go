package lua

// File is responsible for taking a list of tokens and transforming it into a list of executable statements,
// which will be passed towards the VM at execution stage. Bytecode is not compatible with other implementations.

import (
	"fmt"
	"reflect"
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

func isAssignableTyp(t int) bool {
	return t >= identifier_token
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

func PushValue(tar *[]Statement, counter *int, value interface{}) {
	(*tar)[*counter] = PushValueStack{
		value,
	}
	(*counter)++
}

func PushVariable(tar *[]Statement, counter *int, name string) {
	(*tar)[*counter] = PushVariableStack{
		name,
	}
	(*counter)++
}

func VariableAssign(tar *[]Statement, counter *int, name string, local bool) {
	(*tar)[*counter] = VariableAssignment{
		name,
		local,
	}
	(*counter)++
}

func AddStatement(tar *[]Statement, counter *int, value interface{}) {
	typBefore1 := reflect.TypeOf((*tar)[*counter - 1])
	fmt.Println(typBefore1)

	(*tar)[*counter] = value
	(*counter)++
}

func EvaluateStatement(tar *[]Statement, counter *int, tok []Token) (endAt int) {
	count := 0
	lastType := -1
	for count < len(tok) {

		switch tok[count].typ {
		case language_token:
			var stat Statement

			if tok[count].value == "+" {
				stat = AddOperation{}
			} else if tok[count].value == "-" {
				stat = SubOperation{}
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
			EvaluateStatement(&result, &currentStatement, tok[currentToken + 1:len(tok)])
			VariableAssign(&result, &currentStatement, tok[currentToken - 1].value, false)
		}

		currentToken++
	}

	realResult := result[0:currentStatement]
	return &realResult
}
