package lua

import (
	"fmt"
	"strings"
	"testing"
)

func compareTokens(first, second []Token) bool {

	if len(first) != len(second) {
		return false
	}

	return true
}

func printTokens(first []Token) {
	for _, token := range first {
		fmt.Printf("[%10v, %4v, %s]\n", tokenMapping[token.typ], token.line, token.value)
	}
}

func TestBasic(t *testing.T) {

	code := `
	function fact (n)
    	if n == 0 then
    		return 1
    	else
     		return n * fact(n-1)
    	end
    end
    
    print("enter a number:")`

	tokens := Tokenize(strings.NewReader(code))

	if !compareTokens(tokens, tokens) {
		t.Fail()
	}

}

func TestFunctionCallTokens(t *testing.T) {

	code := `print("Hello world")`

	tokens := Tokenize(strings.NewReader(code))

	if len(tokens) != 4 {
		t.Fail()
	}
}

func TestMultiParameterFunctionCallTokens(t *testing.T) {

	code := `
	function testy123(this, hopefully, works)
		print(hopefully)
	end

	testy123("1gadg", 126, 'adgiodagu')
	`

	tokens := Tokenize(strings.NewReader(code))

	if len(tokens) != 22 {
		t.Fail()
	}
}

func TestVariableAssignmentToken(t *testing.T) {

	code := "local x = 25432125"
	tokens := Tokenize(strings.NewReader(code))

	if len(tokens) != 4 {
		t.Fail()
	}
}
