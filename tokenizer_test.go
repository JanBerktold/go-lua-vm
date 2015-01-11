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
		fmt.Printf("[%10v, %s]\n", tokenMapping[token.typ], token.value)
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
	printTokens(tokens)

	if !compareTokens(tokens, tokens) {
		t.Fail()
	}

}

func TestFunctionCall(t *testing.T) {

	code := `print("Hello world")`

	tokens := Tokenize(strings.NewReader(code))
	printTokens(tokens)

	if len(tokens) != 4 {
		t.Fail()
	}
}
