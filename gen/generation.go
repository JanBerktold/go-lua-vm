package gen

import (
	"lua"
	"io"
	"bytes"
)

func GetCodeFromTokens(tok []Tokens) string {
	var buffers bytes.Buffer
	WriteCodeFromTokens(tok, Buffer)
	return buffers.String()
}

func WriteCodeFromTokens(tok []Tokens, target io.Writer) {

}