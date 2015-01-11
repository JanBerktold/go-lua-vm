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
	value interface{}
	local bool
}

type PushStack struct {
	value interface{}
}

type ReturnValue struct {
	amount int
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
	(*tar)[*counter] = PushStack{
		value,
	}
	(*counter)++
}

func CreateBytecode(tok []Token) *[]Statement {
	result := make([]Statement, 100, 1000)

	currentStatement := 0
	currentToken := 0

	for currentToken < len(tok) {
		token := tok[currentToken]

		if token.typ == identifier_token && tok[currentToken+1].value == "=" {

			// Special case, if we want to assign a function to the variable
			if isAssignableTyp(tok[currentToken+2].typ) {
				statment := VariableAssignment{token.value, tok[currentToken+2].real, false}

				if currentToken > 0 {
					statment.local = tok[currentToken-1].typ == keyword_token && tok[currentToken-1].value == "local"
				}

				result[currentStatement] = statment
				currentStatement++
				currentToken += 2
			} else if tok[currentToken+2].typ == keyword_token && tok[currentToken+2].value == "function" {
				// Assign a function
				funcEnd := getFuncEnd(&tok, currentToken+2)
				statment := VariableAssignment{token.value, CreateBytecode(tok[currentToken+3 : funcEnd-1]), false}

				if currentToken > 0 {
					statment.local = tok[currentToken-1].typ == keyword_token && tok[currentToken-1].value == "local"
				}

				result[currentStatement] = statment
				currentStatement++
				currentToken = funcEnd + 1

			}

			continue
		}

		// Function
		if token.typ == identifier_token && tok[currentToken+1].value == "(" {
			paramEnd := getFuncParamEnd(&tok, currentToken+1)

			if currentToken-1 > 0 && tok[currentToken-1].value == "function" {
				// Assignment
				funcEnd := getFuncEnd(&tok, currentToken+1)
				statment := VariableAssignment{token.value, CreateBytecode(tok[paramEnd+1 : funcEnd-1]), false}

				if currentToken > 0 {
					statment.local = tok[currentToken-1].typ == keyword_token && tok[currentToken-1].value == "local"
				}

				result[currentStatement] = statment
				currentStatement++
				currentToken = funcEnd + 1
			} else {
				// Call

				//statment := FunctionCall{token.value, 2}

			}

		}

		// RETURN STATEMENT
		if token.typ == keyword_token && token.value == "return" {
			amount := 0

			for currentToken < len(tok) {
				currentToken++
				if isAssignableTyp(tok[currentToken].typ) {
					amount++
					PushValue(&result, &currentStatement, tok[currentToken].real)
					if currentToken+1 >= len(tok) || tok[currentToken+1].typ != language_token || tok[currentToken+1].value != "," {
						break
					} else {
						currentToken++
					}
				} else {
					break
				}
			}

			result[currentStatement] = ReturnValue{
				amount,
			}
			currentStatement++
		}

		currentToken++
	}

	realResult := result[0:currentStatement]
	return &realResult
}
