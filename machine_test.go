package lua

import (
	"testing"
)

func TestNumberVariableRetrievement(t *testing.T) {
	vm := New()
	vm.ExecuteString(`x = 10`)

	if i, suceed := vm.GetGlobalVariable("x").(float64); !suceed || i != 10 {
		t.Fatal("Number variable could not be retrieved")
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
