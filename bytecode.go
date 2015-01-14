package lua

// File is responsible for taking a list of tokens and transforming it into a list of executable statements,
// which will be passed towards the VM at execution stage. Bytecode is not compatible with other implementations.

import "fmt"
type Statement interface{}

type FunctionCall struct {
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

type IntegerPair struct {
	first, second int
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

		if (*tok)[n].typ == language_token {

			if (*tok)[n].real == ")" {
				level--
			} else if (*tok)[n].real == "(" {
				level++
			}
		}

	}
	return n
}

func SearchMatchingSymbol(tok []Token, params map[string]int) int {
	level := 1
	count := 0
	for count = 0; count < len(tok) && level > 0; count++ {
		if value, exist := params[tok[count].value]; exist {
			level += value
		}
	}
	return count

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

func SplitParameters(tok []Token) []IntegerPair {
	if len(tok) > 0 {
		result := make([]IntegerPair, 10, 100)

		mapSymbols := map[string]int{
			"(": 1,
			")": -1,
			"{": 1,
			"}": -1,
		}

		level := 0
		count := 0
		lastHit := 0
		for n, t := range tok {
			if value, exist := mapSymbols[t.value]; exist {
				level += value
			}
			if t.value == "," && level == 0 {
				result[count] = IntegerPair{lastHit, n - 1}
				count++
				lastHit = n + 1
			}
		}

		if lastHit != 0 {
			lastHit++
		}
		result[count] = IntegerPair{lastHit-1, len(tok)}
		count++

		return result[0:count]
	}

	return []IntegerPair{}
}

func HandleFunctionCall(tar *[]Statement, tok []Token, count int, counter *int) {
	seperateEnd := SearchMatchingSymbol(tok[count+1:len(tok)], map[string]int{
		"(": 1,
		")": -1,
	}) + count - 1

	// Parameter list goes from count to seperateEnd
	params := SplitParameters(tok[count+1:seperateEnd+1])
				
	for _, pair := range params {
		EvaluateStatements(tar, counter, tok[pair.first+count+1:pair.second+count+2])
	}

	AddStatement(tar, counter, FunctionCall{
		len(params),
	})
}

func EvaluateStatement(tar *[]Statement, counter *int, tok []Token, count, lastType int) (int, bool, int) {
	fmt.Printf("LENGTH: %v\n", len(tok))
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
		} else if tok[count].value == "(" {
			if lastType == language_token {
				seperateEnd := SearchMatchingSymbol(tok[count+1:len(tok)], map[string]int{
					"(": 1,
					")": -1,
				}) + count

				(*counter) = (*counter) - 1
				savedOperation := (*tar)[*counter]
				(*tar)[*counter] = nil

				EvaluateStatements(tar, counter, tok[count+1:seperateEnd])
				AddStatement(tar, counter, savedOperation)
				return tok[count].typ, false, seperateEnd
			} else if lastType == identifier_token  {
				HandleFunctionCall(tar, tok, count, counter)
			}
		}

		if stat != nil {
			returnAt := count + 1
			fmt.Printf("COUNT %v", count+1)
			if len(tok) >= count+1 && tok[count+1].typ >= identifier_token {
				EvaluateStatement(tar, counter, tok, count+1, tok[count].typ)
				returnAt++
			}
			AddStatement(tar, counter, stat)
			return tok[count].typ, false, returnAt
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

func GetNextSign(tok []Token) int {
	for n := 0; n < len(tok); n++ {
		if tok[n].value == "=" {
				fmt.Printf("RETURN %v", n)
			return n
		}
	}
	return len(tok)
}

func CreateBytecode(tok []Token) *[]Statement {
	result := make([]Statement, 100, 1000)

	currentStatement := 0
	currentToken := 0

	for currentToken < len(tok) {
		token := tok[currentToken]

		if token.typ == language_token && token.value == "=" {
			EvaluateStatements(&result, &currentStatement, tok[currentToken+1:currentToken + GetNextSign(tok[currentToken+1:len(tok)])])
			VariableAssign(&result, &currentStatement, tok[currentToken-1].value, currentToken-2 >= 0 && tok[currentToken-2].value == "local")
		} else if token.typ == language_token && token.value == "(" && tok[currentToken-1].typ == identifier_token {
			PushVariable(&result, &currentStatement, tok[currentToken-1].value)
			HandleFunctionCall(&result, tok, currentToken, &currentStatement)
		}

		currentToken++
	}

	realResult := result[0:currentStatement]
	return &realResult
}
