package main

import (
	"unicode"
)

type TokenTypeC string

const (
	FOR_C         TokenTypeC = "FOR"
	DO_C          TokenTypeC = "DO"
	WHILE_C       TokenTypeC = "WHILE"
	IDENTIFIER_C  TokenTypeC = "IDENTIFIER"
	NUMBER_C      TokenTypeC = "NUMBER"
	OPERATOR_C    TokenTypeC = "OPERATOR"
	LPAREN_C      TokenTypeC = "LPAREN"
	RPAREN_C      TokenTypeC = "RPAREN"
	LBRACE_C      TokenTypeC = "LBRACE"
	RBRACE_C      TokenTypeC = "RBRACE"
	SEMICOLON_C   TokenTypeC = "SEMICOLON"
	STRING_C      TokenTypeC = "STRING"
	KEYWORD_C     TokenTypeC = "KEYWORD"
	COMPARISON_C  TokenTypeC = "COMPARISON"
	INCREMENT_C   TokenTypeC = "INCREMENT"
	ASSIGNMENT_C  TokenTypeC = "ASSIGNMENT"
	TYPE_C        TokenTypeC = "TYPE"
	WHITESPACE_C  TokenTypeC = "WHITESPACE"
	UNKNOWN_C     TokenTypeC = "UNKNOWN"
)

type TokenC struct {
	Type     TokenTypeC `json:"type"`
	Value    string     `json:"value"`
	Position int        `json:"position"`
	Line     int        `json:"line"`
	Column   int        `json:"column"`
}

type LexerC struct {
	input    string
	position int
	line     int
	column   int
}

func NewLexerC(input string) *LexerC {
	return &LexerC{
		input:    input,
		position: 0,
		line:     1,
		column:   1,
	}
}

func (l *LexerC) Tokenize() []TokenC {
	var tokens []TokenC
	keywords := map[string]TokenTypeC{
		"for":      FOR_C,
		"do":       DO_C,
		"while":    WHILE_C,
		"if":       KEYWORD_C,
		"else":     KEYWORD_C,
		"return":   KEYWORD_C,
		"break":    KEYWORD_C,
		"continue": KEYWORD_C,
		"int":      TYPE_C,
		"float":    TYPE_C,
		"double":   TYPE_C,
		"char":     TYPE_C,
		"void":     TYPE_C,
		"long":     TYPE_C,
		"short":    TYPE_C,
		"printf":   KEYWORD_C,
		"scanf":    KEYWORD_C,
		"include":  KEYWORD_C,
		"main":     KEYWORD_C,
		"stdio":    KEYWORD_C,
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

		// Strings
		if l.input[l.position] == '"' || l.input[l.position] == '\'' {
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

func (l *LexerC) consumeWhitespace() {
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

func (l *LexerC) consumeLineComment() {
	for l.position < len(l.input) && l.input[l.position] != '\n' {
		l.position++
		l.column++
	}
}

func (l *LexerC) consumeBlockComment() {
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

func (l *LexerC) consumeNumber() TokenC {
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

		return TokenC{
			Type:     UNKNOWN_C,
			Value:    l.input[start:l.position],
			Position: start,
			Line:     l.line,
			Column:   startCol,
		}
	}

	return TokenC{
		Type:     NUMBER_C,
		Value:    l.input[start:l.position],
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *LexerC) consumeIdentifier() TokenC {
	start := l.position
	startCol := l.column

	for l.position < len(l.input) && (unicode.IsLetter(rune(l.input[l.position])) ||
		unicode.IsDigit(rune(l.input[l.position])) || l.input[l.position] == '_') {
		l.position++
		l.column++
	}

	return TokenC{
		Type:     IDENTIFIER_C,
		Value:    l.input[start:l.position],
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *LexerC) consumeString() TokenC {
	start := l.position
	startCol := l.column
	quote := l.input[l.position]
	l.position++ // consume opening quote
	l.column++

	for l.position < len(l.input) && l.input[l.position] != quote {
		if l.input[l.position] == '\\' {
			l.position++ // skip escape character
			l.column++
		}
		l.position++
		l.column++
	}

	if l.position < len(l.input) {
		l.position++ // consume closing quote
		l.column++
	}

	return TokenC{
		Type:     STRING_C,
		Value:    l.input[start:l.position],
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *LexerC) consumeOperatorOrSymbol() TokenC {
	start := l.position
	startCol := l.column
	char := l.input[l.position]

	l.position++
	l.column++

	// Check for triple-character operators
	if l.position+1 < len(l.input) {
		threeChar := string(char) + string(l.input[l.position]) + string(l.input[l.position+1])
		if threeChar == "===" || threeChar == "!==" {
			l.position += 2
			l.column += 2
			return TokenC{
				Type:     COMPARISON_C,
				Value:    threeChar,
				Position: start,
				Line:     l.line,
				Column:   startCol,
			}
		}
	}

	// Check for multi-character operators
	if l.position < len(l.input) {
		twoChar := string(char) + string(l.input[l.position])
		switch twoChar {
		case "<=", ">=", "==", "!=", "++", "--", "+=", "-=", "*=", "/=":
			l.position++
			l.column++
			return TokenC{
				Type:     l.getOperatorType(twoChar),
				Value:    twoChar,
				Position: start,
				Line:     l.line,
				Column:   startCol,
			}
		}
	}

	// Single character operators/symbols
	tokenType := l.getSingleCharType(char)
	return TokenC{
		Type:     tokenType,
		Value:    string(char),
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *LexerC) getOperatorType(op string) TokenTypeC {
	switch op {
	case "<=", ">=", "==", "!=", "<", ">":
		return COMPARISON_C
	case "++", "--":
		return INCREMENT_C
	case "+=", "-=", "*=", "/=", "=":
		return ASSIGNMENT_C
	default:
		return OPERATOR_C
	}
}

func (l *LexerC) getSingleCharType(char byte) TokenTypeC {
	switch char {
	case '(':
		return LPAREN_C
	case ')':
		return RPAREN_C
	case '{':
		return LBRACE_C
	case '}':
		return RBRACE_C
	case ';':
		return SEMICOLON_C
	case '=':
		return ASSIGNMENT_C
	case '<', '>':
		return COMPARISON_C
	case '+', '-', '*', '/':
		return OPERATOR_C
	default:
		return UNKNOWN_C
	}
}