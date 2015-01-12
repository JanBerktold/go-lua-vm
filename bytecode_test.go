package lua

import (
	"strings"
	"testing"
	"fmt"
)

func TestNumberAssignment(t *testing.T) {

	code := "local x = 25432125 + 21"
	tokens := Tokenize(strings.NewReader(code))
	statements := CreateBytecode(tokens)

	fmt.Println(statements)

	if len(*statements) != 1 {
		t.FailNow()
	}



}

func TestStringAssignment(t *testing.T) {

	code := "x = 'hello there123'"
	tokens := Tokenize(strings.NewReader(code))
	statements := CreateBytecode(tokens)



	if len(*statements) != 1 {
		t.FailNow()
	}


}

func TestFunctionAssignment(t *testing.T) {

	code := "x = function() print('hello') end"
	tokens := Tokenize(strings.NewReader(code))
	statements := CreateBytecode(tokens)

	if len(*statements) != 1 {
		t.FailNow()
	}
}

func TestFunctionCall(t *testing.T) {

	code := "print('hello')"
	tokens := Tokenize(strings.NewReader(code))
	statements := CreateBytecode(tokens)

	if len(*statements) != 1 {
		t.FailNow()
	}
}

func TestFunctionMultipliParamsCall(t *testing.T) {

	code := "print('hello', 123, 'hi', 25.02)"
	tokens := Tokenize(strings.NewReader(code))
	statements := CreateBytecode(tokens)

	if len(*statements) != 1 {
		t.FailNow()
	}
}

func TestSingleReturn(t *testing.T) {
	code := "return 123, 'agfda', 02"
	tokens := Tokenize(strings.NewReader(code))
	statements := CreateBytecode(tokens)

	if len(*statements) != 4 {
		t.FailNow()
	}
}
