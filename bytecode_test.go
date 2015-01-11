package lua

import (
	"strings"
	"testing"
	"fmt"
)

func TestNumberAssignment(t *testing.T) {

	code := "local x = 25432125"
	tokens := Tokenize(strings.NewReader(code))
	statements := CreateBytecode(tokens)

	if len(*statements) != 1 {
		t.Fail()
	}

	if stat, suceed := (*statements)[0].(VariableAssignment); suceed {
		if stat.name != "x" || !stat.local {
			t.Fail()
		}
		if n, succeedNumber := stat.value.(float64); succeedNumber {
			if n != 25432125 {
				t.Fail()
			}
		} else {
			t.Fail()
		}
	} else {
		t.Fail()
	}

}

func TestStringAssignment(t *testing.T) {

	code := "x = 'hello there123'"
	tokens := Tokenize(strings.NewReader(code))
	statements := CreateBytecode(tokens)

	if len(*statements) != 1 {
		t.Fail()
	}

	if stat, suceed := (*statements)[0].(VariableAssignment); suceed {
		if stat.name != "x" || stat.local {
			t.Fail()
		}
		if n, succeedNumber := stat.value.(string); succeedNumber {
			if n != "hello there123" {
				t.Fail()
			}
		} else {
			t.Fail()
		}
	} else {
		t.Fail()
	}

}

func TestFunctionAssignment(t *testing.T) {

	code := "x = function() print('hello') end"
	tokens := Tokenize(strings.NewReader(code))
	statements := CreateBytecode(tokens)

	if len(*statements) != 1 {
		t.Fail()
	}
}

func TestFunctionCall(t *testing.T) {

	code := "print('hello')"
	tokens := Tokenize(strings.NewReader(code))
	statements := CreateBytecode(tokens)

	if len(*statements) != 1 {
		t.Fail()
	}
}

func TestFunctionMultipliParamsCall(t *testing.T) {

	code := "print('hello', 123, 'hi', 25.02)"
	tokens := Tokenize(strings.NewReader(code))
	statements := CreateBytecode(tokens)

	if len(*statements) != 1 {
		t.Fail()
	}
}

func TestSingleReturn(t *testing.T) {
	code := "return 123, 'agfda', 02"
	tokens := Tokenize(strings.NewReader(code))
	statements := CreateBytecode(tokens)

	if len(*statements) != 4 {
		t.Fail()
	}
}
