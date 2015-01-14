package lua

import (
	"testing"
)

type ArithmeticTest struct {
	equation string
	result   float64
}

func TestArithmeticRetrievement(t *testing.T) {
	vm := New()

	tests := []ArithmeticTest{
		ArithmeticTest{`x = 10`, 10},
		ArithmeticTest{`x = 50 * 2 + 21`, 121},
		ArithmeticTest{`x = 2 + 21 + 50 + 46 + 8`, 127},
		ArithmeticTest{`x = 50 - 54 + 98 - 2`, 92},
		ArithmeticTest{`x = 50 / 10`, 5},
		ArithmeticTest{`x = 50 / 10 + 2.21`, 7.21},
		ArithmeticTest{`x = 50 / 10 * 10`, 50},
		ArithmeticTest{`x = 2 * (20 + 10)`, 60},
		ArithmeticTest{`x = 2 * 20 + 10`, 50},
		ArithmeticTest{`y = 20 / 10 x = y * 2`, 4},
	}

	for _, test := range tests {
		vm.ExecuteString(test.equation)
		if i, suceed := vm.GetGlobalVariable("x").(float64); !suceed || i != test.result {
			t.Fatalf("Number variable of equation %q could not be retrieved or wrong value. \nExpected: %v. Recieved: %v\n", test.equation, test.result, i)
		}
	}

}

func TestStringVariableRetrievement(t *testing.T) {
	vm := New()
	vm.ExecuteString(`x = "abcadoi 35235 0123/39#'+"`)

	if i, suceed := vm.GetGlobalVariable("x").(string); !suceed || i != "abcadoi 35235 0123/39#'+" {
		t.Fatal("String variable could not be retrieved")
	}
}

func TestLocalVariableRetrievement(t *testing.T) {
	vm := New()
	vm.ExecuteString(`local x = 134315654`)

	if value := vm.GetGlobalVariable("x"); value != nil {
		t.Fatal("Variable could be retrieved even though declared as local")
	}
}
