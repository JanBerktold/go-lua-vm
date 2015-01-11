package lua

// TODO:
// - Fix bug where identifier is dropping the next rune
// - Fix language tokens

import (
	"io"
	"strconv"
	"bytes"
	"unicode"
	"fmt"
)

type Token struct {
	typ int
	value string
}

type continueReading func(rune) bool

func isWhiteSpace(r rune) bool {
	return (r >= 9 && r <= 13) || (r == 32)
}

func isKeyword(s string) bool {
	return keywordMap[s]
}

func RuneToAscii(r rune) string {
    if r < 128 {
        return string(r)
    } else {
        return "\\u" + strconv.FormatInt(int64(r), 16)
    }
}

func ReadUntil(reader *io.RuneReader, start string, fn continueReading) (string, rune) {
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

func IsToken(s string) (bool, bool) {
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

	for {
		var readRune rune
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

		// skip whitespace
		if isWhiteSpace(readRune) {
			continue
		}

		// language tokens
		if isToken, isCompleted := IsToken(RuneToAscii(readRune)); isToken {

			if isCompleted {
				result[tokenNum] = Token{language_token, RuneToAscii(readRune)}
				tokenNum++
				continue
			}

			buffer := bytes.NewBufferString("")
			buffer.WriteRune(readRune)

			for {
				newRune, _, runeErr := runeReader.ReadRune()
				if !isCompleted && runeErr == nil {
					buffer.WriteRune(newRune)
					if _, newCompleted := IsToken(buffer.String()); newCompleted {
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

			result[tokenNum] = Token{language_token, buffer.String()}
			tokenNum++
			continue
		}

		// string
		if readRune == 34 || readRune == 39 {
			str, readRune = ReadUntil(&runeReader, "", func(r rune) bool {
				return r != readRune
			})

			result[tokenNum] = Token{string_token, str}
			tokenNum++
			continue
		}

		// any non-number identifer
		if !unicode.IsDigit(readRune) {
			str, readRune = ReadUntil(&runeReader, RuneToAscii(readRune), func(r rune) bool {
				return unicode.IsNumber(r) || unicode.IsLetter(r)
			})

			fmt.Println("MISSED RUNE: " + RuneToAscii(readRune))


			token := Token{0, str}

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
			str, readRune = ReadUntil(&runeReader, RuneToAscii(readRune), func(r rune) bool {
				return unicode.IsDigit(r) || r == 46
			})
			savedRune = true
		
			result[tokenNum] = Token{number_token, str}
			tokenNum++
			continue
		}

	}
	
	return result[0:tokenNum]
}