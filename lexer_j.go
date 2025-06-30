package main

import (
	"unicode"
)

type TokenTypeJ string

const (
	PUBLIC_J      TokenTypeJ = "PUBLIC"
	CLASS_J       TokenTypeJ = "CLASS"
	STATIC_J      TokenTypeJ = "STATIC"
	VOID_J        TokenTypeJ = "VOID"
	MAIN_J        TokenTypeJ = "MAIN"
	IF_J          TokenTypeJ = "IF"
	ELSE_J        TokenTypeJ = "ELSE"
	IDENTIFIER_J  TokenTypeJ = "IDENTIFIER"
	NUMBER_J      TokenTypeJ = "NUMBER"
	STRING_J      TokenTypeJ = "STRING"
	OPERATOR_J    TokenTypeJ = "OPERATOR"
	LPAREN_J      TokenTypeJ = "LPAREN"
	RPAREN_J      TokenTypeJ = "RPAREN"
	LBRACE_J      TokenTypeJ = "LBRACE"
	RBRACE_J      TokenTypeJ = "RBRACE"
	LBRACKET_J    TokenTypeJ = "LBRACKET"
	RBRACKET_J    TokenTypeJ = "RBRACKET"
	SEMICOLON_J   TokenTypeJ = "SEMICOLON"
	DOT_J         TokenTypeJ = "DOT"
	COMMA_J       TokenTypeJ = "COMMA"
	KEYWORD_J     TokenTypeJ = "KEYWORD"
	COMPARISON_J  TokenTypeJ = "COMPARISON"
	ASSIGNMENT_J  TokenTypeJ = "ASSIGNMENT"
	TYPE_J        TokenTypeJ = "TYPE"
	WHITESPACE_J  TokenTypeJ = "WHITESPACE"
	UNKNOWN_J     TokenTypeJ = "UNKNOWN"
)

type TokenJ struct {
	Type     TokenTypeJ `json:"type"`
	Value    string     `json:"value"`
	Position int        `json:"position"`
	Line     int        `json:"line"`
	Column   int        `json:"column"`
}

type LexerJ struct {
	input    string
	position int
	line     int
	column   int
}

func NewLexerJ(input string) *LexerJ {
	return &LexerJ{
		input:    input,
		position: 0,
		line:     1,
		column:   1,
	}
}

func (l *LexerJ) Tokenize() []TokenJ {
	var tokens []TokenJ
	keywords := map[string]TokenTypeJ{
		"public":    PUBLIC_J,
		"class":     CLASS_J,
		"static":    STATIC_J,
		"void":      VOID_J,
		"main":      MAIN_J,
		"if":        IF_J,
		"else":      ELSE_J,
		"return":    KEYWORD_J,
		"break":     KEYWORD_J,
		"continue":  KEYWORD_J,
		"while":     KEYWORD_J,
		"for":       KEYWORD_J,
		"int":       TYPE_J,
		"String":    TYPE_J,
		"boolean":   TYPE_J,
		"double":    TYPE_J,
		"float":     TYPE_J,
		"char":      TYPE_J,
		"long":      TYPE_J,
		"System":    KEYWORD_J,
		"out":       KEYWORD_J,
		"println":   KEYWORD_J,
		"print":     KEYWORD_J,
		"equals":    KEYWORD_J,
		"args":      KEYWORD_J,
		"true":      KEYWORD_J,
		"false":     KEYWORD_J,
		"null":      KEYWORD_J,
	}

	for l.position < len(l.input) {
		if unicode.IsSpace(rune(l.input[l.position])) {
			l.consumeWhitespace()
			continue
		}

		// Números
		if unicode.IsDigit(rune(l.input[l.position])) {
			token := l.consumeNumber()
			tokens = append(tokens, token)
			continue
		}

		// Identificadores y palabras clave
		if unicode.IsLetter(rune(l.input[l.position])) || l.input[l.position] == '_' {
			token := l.consumeIdentifier()
			if tokenType, exists := keywords[token.Value]; exists {
				token.Type = tokenType
			}
			tokens = append(tokens, token)
			continue
		}

		// Strings con comillas dobles
		if l.input[l.position] == '"' {
			token := l.consumeString()
			tokens = append(tokens, token)
			continue
		}

		// Comentarios de una línea //
		if l.position+1 < len(l.input) && l.input[l.position] == '/' && l.input[l.position+1] == '/' {
			l.consumeLineComment()
			continue
		}

		// Comentarios de bloque /* */
		if l.position+1 < len(l.input) && l.input[l.position] == '/' && l.input[l.position+1] == '*' {
			l.consumeBlockComment()
			continue
		}

		// Operadores y símbolos
		token := l.consumeOperatorOrSymbol()
		tokens = append(tokens, token)
	}

	return tokens
}

func (l *LexerJ) consumeWhitespace() {
	for l.position < len(l.input) && unicode.IsSpace(rune(l.input[l.position])) {
		if l.input[l.position] == '\n' {
			l.line++
			l.column = 1
		} else {
			l.column++
		}
		l.position++
	}
}

func (l *LexerJ) consumeLineComment() {
	for l.position < len(l.input) && l.input[l.position] != '\n' {
		l.position++
		l.column++
	}
}

func (l *LexerJ) consumeBlockComment() {
	l.position += 2 // skip /*
	l.column += 2

	for l.position+1 < len(l.input) {
		if l.input[l.position] == '*' && l.input[l.position+1] == '/' {
			l.position += 2
			l.column += 2
			break
		}
		if l.input[l.position] == '\n' {
			l.line++
			l.column = 1
		} else {
			l.column++
		}
		l.position++
	}
}

func (l *LexerJ) consumeNumber() TokenJ {
	start := l.position
	startCol := l.column

	for l.position < len(l.input) && (unicode.IsDigit(rune(l.input[l.position])) || l.input[l.position] == '.') {
		l.position++
		l.column++
	}

	// Verificar si hay letras inmediatamente después del número (error)
	if l.position < len(l.input) && unicode.IsLetter(rune(l.input[l.position])) {
		for l.position < len(l.input) && (unicode.IsLetter(rune(l.input[l.position])) || unicode.IsDigit(rune(l.input[l.position]))) {
			l.position++
			l.column++
		}

		return TokenJ{
			Type:     UNKNOWN_J,
			Value:    l.input[start:l.position],
			Position: start,
			Line:     l.line,
			Column:   startCol,
		}
	}

	return TokenJ{
		Type:     NUMBER_J,
		Value:    l.input[start:l.position],
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *LexerJ) consumeIdentifier() TokenJ {
	start := l.position
	startCol := l.column

	for l.position < len(l.input) && (unicode.IsLetter(rune(l.input[l.position])) ||
		unicode.IsDigit(rune(l.input[l.position])) || l.input[l.position] == '_') {
		l.position++
		l.column++
	}

	return TokenJ{
		Type:     IDENTIFIER_J,
		Value:    l.input[start:l.position],
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *LexerJ) consumeString() TokenJ {
	start := l.position
	startCol := l.column
	l.position++ // consume opening quote
	l.column++

	for l.position < len(l.input) && l.input[l.position] != '"' {
		if l.input[l.position] == '\\' {
			l.position++ // skip escape character
			l.column++
		}
		if l.position < len(l.input) {
			l.position++
			l.column++
		}
	}

	if l.position < len(l.input) {
		l.position++ // consume closing quote
		l.column++
	}

	return TokenJ{
		Type:     STRING_J,
		Value:    l.input[start:l.position],
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *LexerJ) consumeOperatorOrSymbol() TokenJ {
	start := l.position
	startCol := l.column
	char := l.input[l.position]

	l.position++
	l.column++

	// Check for multi-character operators
	if l.position < len(l.input) {
		twoChar := string(char) + string(l.input[l.position])
		switch twoChar {
		case "<=", ">=", "==", "!=", "&&", "||":
			l.position++
			l.column++
			return TokenJ{
				Type:     COMPARISON_J,
				Value:    twoChar,
				Position: start,
				Line:     l.line,
				Column:   startCol,
			}
		}
	}

	// Single character operators/symbols
	tokenType := l.getSingleCharType(char)
	return TokenJ{
		Type:     tokenType,
		Value:    string(char),
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *LexerJ) getSingleCharType(char byte) TokenTypeJ {
	switch char {
	case '(':
		return LPAREN_J
	case ')':
		return RPAREN_J
	case '{':
		return LBRACE_J
	case '}':
		return RBRACE_J
	case '[':
		return LBRACKET_J
	case ']':
		return RBRACKET_J
	case ';':
		return SEMICOLON_J
	case '.':
		return DOT_J
	case ',':
		return COMMA_J
	case '=':
		return ASSIGNMENT_J
	case '<', '>':
		return COMPARISON_J
	case '+', '-', '*', '/':
		return OPERATOR_J
	default:
		return UNKNOWN_J
	}
}