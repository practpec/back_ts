package main

import (
	"unicode"
)

type TokenType string

const (
	FOR         TokenType = "FOR"
	DO          TokenType = "DO"
	WHILE       TokenType = "WHILE"
	IDENTIFIER  TokenType = "IDENTIFIER"
	NUMBER      TokenType = "NUMBER"
	OPERATOR    TokenType = "OPERATOR"
	LPAREN      TokenType = "LPAREN"
	RPAREN      TokenType = "RPAREN"
	LBRACE      TokenType = "LBRACE"
	RBRACE      TokenType = "RBRACE"
	SEMICOLON   TokenType = "SEMICOLON"
	COLON       TokenType = "COLON"
	STRING      TokenType = "STRING"
	KEYWORD     TokenType = "KEYWORD"
	COMPARISON  TokenType = "COMPARISON"
	INCREMENT   TokenType = "INCREMENT"
	ASSIGNMENT  TokenType = "ASSIGNMENT"
	TYPE        TokenType = "TYPE"
	WHITESPACE  TokenType = "WHITESPACE"
	UNKNOWN     TokenType = "UNKNOWN"
)

type Token struct {
	Type     TokenType `json:"type"`
	Value    string    `json:"value"`
	Position int       `json:"position"`
	Line     int       `json:"line"`
	Column   int       `json:"column"`
}

type Lexer struct {
	input    string
	position int
	line     int
	column   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		position: 0,
		line:   1,
		column: 1,
	}
}

func (l *Lexer) Tokenize() []Token {
	var tokens []Token
	keywords := map[string]TokenType{
		"for":    FOR,
		"do":     DO,
		"while":  WHILE,
		"let":    KEYWORD,
		"const":  KEYWORD,
		"var":    KEYWORD,
		"int":    TYPE,
		"string": TYPE,
		"number": TYPE,
		"boolean": TYPE,
		"console": KEYWORD,
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
		if unicode.IsLetter(rune(l.input[l.position])) {
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
		
		// Operadores y símbolos
		token := l.consumeOperatorOrSymbol()
		tokens = append(tokens, token)
	}
	
	return tokens
}

func (l *Lexer) consumeWhitespace() {
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

func (l *Lexer) consumeNumber() Token {
	start := l.position
	startCol := l.column
	
	for l.position < len(l.input) && (unicode.IsDigit(rune(l.input[l.position])) || l.input[l.position] == '.') {
		l.position++
		l.column++
	}
	
	// Verificar si hay letras inmediatamente después del número (error)
	if l.position < len(l.input) && unicode.IsLetter(rune(l.input[l.position])) {
		// Consumir las letras para incluirlas en el token de error
		for l.position < len(l.input) && (unicode.IsLetter(rune(l.input[l.position])) || unicode.IsDigit(rune(l.input[l.position]))) {
			l.position++
			l.column++
		}
		
		return Token{
			Type:     UNKNOWN,
			Value:    l.input[start:l.position],
			Position: start,
			Line:     l.line,
			Column:   startCol,
		}
	}
	
	return Token{
		Type:     NUMBER,
		Value:    l.input[start:l.position],
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *Lexer) consumeIdentifier() Token {
	start := l.position
	startCol := l.column
	
	for l.position < len(l.input) && (unicode.IsLetter(rune(l.input[l.position])) || 
		unicode.IsDigit(rune(l.input[l.position])) || l.input[l.position] == '_') {
		l.position++
		l.column++
	}
	
	return Token{
		Type:     IDENTIFIER,
		Value:    l.input[start:l.position],
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *Lexer) consumeString() Token {
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
	
	return Token{
		Type:     STRING,
		Value:    l.input[start:l.position],
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *Lexer) consumeOperatorOrSymbol() Token {
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
			return Token{
				Type:     COMPARISON,
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
		case "<=", ">=", "==", "!=", "++", "--", "+=", "-=":
			l.position++
			l.column++
			return Token{
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
	return Token{
		Type:     tokenType,
		Value:    string(char),
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *Lexer) getOperatorType(op string) TokenType {
	switch op {
	case "<=", ">=", "==", "!=", "<", ">", "===", "!==":
		return COMPARISON
	case "++", "--":
		return INCREMENT
	case "+=", "-=", "=":
		return ASSIGNMENT
	default:
		return OPERATOR
	}
}

func (l *Lexer) getSingleCharType(char byte) TokenType {
	switch char {
	case '(':
		return LPAREN
	case ')':
		return RPAREN
	case '{':
		return LBRACE
	case '}':
		return RBRACE
	case ';':
		return SEMICOLON
	case ':':
		return COLON
	case '=':
		return ASSIGNMENT
	case '<', '>':
		return COMPARISON
	case '+', '-', '*', '/':
		return OPERATOR
	default:
		return UNKNOWN
	}
}