package lua

import (
	"reflect"
	"strings"
	"testing"
)

type StatmentTest struct {
	code string
	result int
	print bool
}

func printStatements(t *testing.T, st *[]Statement) {
	for _, stat := range *st {
		t.Logf("%v, %v\n", reflect.TypeOf(stat), stat)
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
		StatmentTest{"return 123, 'agfda', 02", 4, false },
	}
	
	for _, test := range tests {
		tokens := Tokenize(strings.NewReader(test.code))
		statements := CreateBytecode(tokens)
		if len(*statements) != test.result {
			t.Errorf("Test %q failed. Got %v statment(s). Expected %v.", test.code, test.result, len(*statements))
		}
		if len(*statements) != test.result || test.print {
			t.Logf("Generated statements from code %q:", test.code)
			printStatements(t, statements)
		}
	}
}
