package lua

// File is responsible for taking a list of tokens and transforming it into a list of executable statements,
// which will be passed towards the VM at execution stage. Bytecode is not compatible with other implementations.

type Statement interface {}

type FunctionCall struct {
	name string
	args int
}

type VariableAssignment struct {
	name string
	value interface{}
	local bool
}

type PushStack struct {

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

func CreateBytecode(tok []Token) *[]Statement {
	result := make([]Statement, 100, 1000)

	currentStatement := 0
	currentToken := 0

	for currentToken < len(tok) {
		token := tok[currentToken]

		if token.typ == identifier_token && tok[currentToken + 1].value == "=" {

			// Special case, if we want to assign a function to the variable
			if isAssignableTyp(tok[currentToken + 2].typ) {
				statment := VariableAssignment{token.value, tok[currentToken + 2].value, false}

				if currentToken > 0 {
					statment.local = tok[currentToken - 1].typ == keyword_token && tok[currentToken - 1].value == "local"
				}

				result[currentStatement] = statment
				currentStatement++
				currentToken += 2
			} else if tok[currentToken + 2].typ == keyword_token &&  tok[currentToken + 2].value == "function" {
				// Assign a function
				funcEnd := getFuncEnd(&tok, currentToken + 2)
				statment := VariableAssignment{token.value, CreateBytecode(tok[currentToken+3:funcEnd-1]), false}

				if currentToken > 0 {
					statment.local = tok[currentToken - 1].typ == keyword_token && tok[currentToken - 1].value == "local"
				}

				result[currentStatement] = statment
				currentStatement++
				currentToken = funcEnd + 1

			}

		}

		

		currentToken++
	}

	realResult := result[0:currentStatement]
	return &realResult
}