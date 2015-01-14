package lua

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type StatmentTest struct {
	code string
	result int
	print bool
}

func printStatements(st *[]Statement) {
	for _, stat := range *st {
		fmt.Printf("%v, %v\n", reflect.TypeOf(stat), stat)
	}
}

func TestBasicStatments(t *testing.T) {
	
	tests := []StatmentTest{
		StatmentTest{"local x = 25432125 + 21", 4, false },
		StatmentTest{"local x = 25432125 + 21 + 50 + 46 + 8", 10, false },
		StatmentTest{"x = 'hello there123'", 2, false },
		StatmentTest{"x = #testaobd + 21", 5, false },
		StatmentTest{"x = function() print('hello') end", 1, false },
		StatmentTest{"print('hello')", 3, false },
		StatmentTest{"print('hello', 123, 'hi', 25.02)", 6, false },
		StatmentTest{"print('hello', 123, 'hi', 25.02)", 6, false },
	}
	
	for _, test := range tests {
		tokens := Tokenize(strings.NewReader(test.code))
		statements := CreateBytecode(tokens)
		if len(*statements) != test.result {
			t.FailNow()
		}
	}
}

func TestSingleReturn(t *testing.T) {
	code := "return 123, 'agfda', 02"
	tokens := Tokenize(strings.NewReader(code))
	statements := CreateBytecode(tokens)

	printStatements(statements)
	if len(*statements) != 4 {
		t.FailNow()
	}
}
