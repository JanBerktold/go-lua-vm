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
	identifier_token
	language_token
	string_token
	number_token
	
)

var (

	// KEYWORDS
	keywords = []string{"and", "break", "do", "else", "elseif", "end", "false", "for", "function", "if", "in", "local", "nil", "not", "or", "repeat", "return", "then", "true", "until", "while"}
	keywordMap = sliceToMap(keywords)

	// LANGUAGE TOKENS
	tokens = []string{"+", "-", "*", "/", "%", "^", "#", "==", "~=", "<=", ">=", "<", ">", "=", "(", ")", "{", "}", "[", "]", ";", ":", ",", ".", "..", "..."}

	// DEBUGGING
	tokenMapping = map[int]string {
		keyword_token: "KEYWORD",
		identifier_token: "IDENTIFIER",
		language_token: "LANGUAGE",
		string_token: "STRING",
		number_token: "NUMBER",
	}
)