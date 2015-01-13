package lua

// File is responsible for taking a list of tokens and transforming it into a list of executable statements,
// which will be passed towards the VM at execution stage. Bytecode is not compatible with other implementations.

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

type MulOperation struct {
}

type DivOperation struct {
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

func EvaluateStatement(tar *[]Statement, counter *int, tok []Token, count, lastType int) (int, bool, int) {
	switch tok[count].typ {
	case language_token:
		var stat Statement

		if tok[count].value == "+" {
			stat = AddOperation{}
		} else if tok[count].value == "-" {
			stat = SubOperation{}
		} else if tok[count].value == "*" {
			stat = MulOperation{}
		} else if tok[count].value == "/" {
			stat = DivOperation{}
		} else if tok[count].value == "#" {
			stat = TableLengthOperation{}
		}

		if stat != nil {
			if tok[count+1].typ >= identifier_token {
				EvaluateStatement(tar, counter, tok, count+1, tok[count].typ)
			}
			AddStatement(tar, counter, stat)
			return tok[count].typ, false, count + 2
		}
	case identifier_token:
		if lastType >= identifier_token {
			return tok[count].typ, true, count + 1
		}
		PushVariable(tar, counter, tok[count].value)
	case number_token, string_token:
		if lastType >= identifier_token {
			return tok[count].typ, true, count + 1
		}
		PushValue(tar, counter, tok[count].real)
	}

	return tok[count].typ, false, count + 1
}

func EvaluateStatements(tar *[]Statement, counter *int, tok []Token) (endAt int) {
	count := 0
	lastType := 0
	suceed := false
	for count < len(tok) {
		lastType, suceed, count = EvaluateStatement(tar, counter, tok, count, lastType)
		if suceed {
			return count
		}
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
			EvaluateStatements(&result, &currentStatement, tok[currentToken+1:len(tok)])
			VariableAssign(&result, &currentStatement, tok[currentToken-1].value, currentToken-2 >= 0 && tok[currentToken-2].value == "local")
		}

		currentToken++
	}

	realResult := result[0:currentStatement]
	return &realResult
}
