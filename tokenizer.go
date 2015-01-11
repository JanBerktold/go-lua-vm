package lua

// TODO:
// - Fix language tokens

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type Token struct {
	typ   int
	value string
	line  int
}

type continueReading func(rune) bool

func isWhiteSpace(r rune) bool {
	return (r >= 9 && r <= 13) || (r == 32)
}

func isKeyword(s string) bool {
	return keywordMap[s]
}

func runeToAscii(r rune) string {
	if r < 128 {
		return string(r)
	} else {
		return "\\u" + strconv.FormatInt(int64(r), 16)
	}
}

func readUntil(reader *io.RuneReader, start string, fn continueReading) (string, rune) {
	buffer := bytes.NewBufferString(start)

	for {
		newRune, _, runeErr := (*reader).ReadRune()
		if fn(newRune) && runeErr == nil {
			buffer.WriteRune(newRune)
		} else {
			return buffer.String(), newRune
		}
	}
}

func isToken(s string) (bool, bool) {
	for _, token := range tokens {
		if len(s) <= len(token) && token[0:len(s)] == s {
			return true, len(s) == len(token)
		}
	}
	return false, false
}

func Tokenize(code io.Reader) []Token {
	runeReader := code.(io.RuneReader)

	result := make([]Token, 100, 1000)
	tokenNum := 0

	var savedRune bool
	var readRune rune
	var line int

	for {
		var str string

		if savedRune {
			savedRune = false
		} else {
			var err error
			readRune, _, err = runeReader.ReadRune()

			if err != nil {
				break
			}
		}

		// increase line counter
		if readRune == '\n' {
			line++
		}

		// skip whitespace
		if isWhiteSpace(readRune) || readRune == 0 {
			continue
		}

		// language tokens
		if isRuneToken, isCompleted := isToken(runeToAscii(readRune)); isRuneToken {

			if isCompleted {
				result[tokenNum] = Token{language_token, runeToAscii(readRune), line}
				tokenNum++
				continue
			}

			buffer := bytes.NewBufferString("")
			buffer.WriteRune(readRune)

			for {
				newRune, _, runeErr := runeReader.ReadRune()
				if !isCompleted && runeErr == nil {
					buffer.WriteRune(newRune)
					stillToken, newCompleted := isToken(buffer.String())
					if newCompleted {
						break
					} else if !stillToken {
						buffer.UnreadRune()
						break
					}

				} else {
					if runeErr == nil {
						savedRune = true
						readRune = newRune
					}
					break
				}
			}

			result[tokenNum] = Token{language_token, strings.TrimSpace(buffer.String()), line}
			tokenNum++
			continue
		}

		// string
		if readRune == 34 || readRune == 39 {
			str, readRune = readUntil(&runeReader, "", func(r rune) bool {
				return r != readRune
			})

			result[tokenNum] = Token{string_token, str, line}
			tokenNum++
			continue
		}

		// any non-number identifer
		if !unicode.IsDigit(readRune) {
			str, readRune = readUntil(&runeReader, runeToAscii(readRune), func(r rune) bool {
				return unicode.IsNumber(r) || unicode.IsLetter(r)
			})
			savedRune = true

			token := Token{0, str, line}

			if isKeyword(str) {
				token.typ = keyword_token
			} else {
				token.typ = identifier_token
			}

			result[tokenNum] = token
			tokenNum++
			continue
		}

		// numbers
		if unicode.IsDigit(readRune) {
			str, readRune = readUntil(&runeReader, runeToAscii(readRune), func(r rune) bool {
				return unicode.IsDigit(r) || r == 46
			})
			savedRune = true

			result[tokenNum] = Token{number_token, str, line}
			tokenNum++
			continue
		}

	}

	return result[0:tokenNum]
}
