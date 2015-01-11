package lua

func sliceToMap(s []string) map[string]bool {
	result := make(map[string]bool)
	for _, str := range s {
		result[str] = true
	}
	return result
}

const (

	// TOKEN TYPES
	keyword_token = iota
	language_token
	identifier_token
	string_token
	number_token
)

var (

	// KEYWORDS
	keywords   = []string{"and", "break", "do", "else", "elseif", "end", "false", "for", "function", "if", "in", "local", "nil", "not", "or", "repeat", "return", "then", "true", "until", "while"}
	keywordMap = sliceToMap(keywords)

	scopeStart = []string{"do", "elseif", "for", "function", "repeat", "while"}
	scopeMap = sliceToMap(scopeStart)
	
	// LANGUAGE TOKENS
	// NOTE: All tokens have to be sorted after their length, starting with the longest one
	tokens = []string{"...", "..", "==", "~=", "<=", ">=", "+", "-", "*", "/", "%", "^", "#", "<", ">", "=", "(", ")", "{", "}", "[", "]", ";", ":", ",", "."}

	// DEBUGGING
	tokenMapping = map[int]string{
		keyword_token:    "KEYWORD",
		identifier_token: "IDENTIFIER",
		language_token:   "LANGUAGE",
		string_token:     "STRING",
		number_token:     "NUMBER",
	}
)
