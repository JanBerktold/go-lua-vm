package lua

import (
	"testing"
)

func TestNumberVariableRetrievement(t *testing.T) {
	vm := New()
	vm.ExecuteString(`x = 10`)

	if i, suceed := vm.GetGlobalVariable("x").(float64); !suceed || i != 10 {
		t.Fail()
	}
}

func TestStringVariableRetrievement(t *testing.T) {
	vm := New()
	vm.ExecuteString(`x = "abcadoi 35235 0123/39#'+"`)

	if i, suceed := vm.GetGlobalVariable("x").(string); !suceed || i != "abcadoi 35235 0123/39#'+" {
		t.Fail()
	}
}

func TestLocalVariableRetrievement(t *testing.T) {
	vm := New()
	vm.ExecuteString(`local x = 134315654`)

	if value := vm.GetGlobalVariable("x"); value != nil {
		t.Fail()
	}
}
