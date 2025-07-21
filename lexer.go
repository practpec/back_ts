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

// Mapa global estático para máximo rendimiento
var keywords = map[string]TokenType{
	"for":     FOR,
	"do":      DO,
	"while":   WHILE,
	"let":     KEYWORD,
	"const":   KEYWORD,
	"var":     KEYWORD,
	"int":     TYPE,
	"string":  TYPE,
	"number":  TYPE,
	"boolean": TYPE,
	"console": KEYWORD,
}

// Arrays estáticos para operadores de múltiples caracteres (máxima eficiencia)
var twoCharOps = [...]string{"<=", ">=", "==", "!=", "++", "--", "+=", "-="}
var threeCharOps = [...]string{"===", "!=="}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:    input,
		position: 0,
		line:     1,
		column:   1,
	}
}

func (l *Lexer) Tokenize() []Token {
	// Pre-allocar slice con capacidad estimada para evitar re-allocaciones
	tokens := make([]Token, 0, len(l.input)/4)
	
	for l.position < len(l.input) {
		char := l.input[l.position]
		
		// Optimización: switch en lugar de múltiples if para caracteres comunes
		switch {
		case unicode.IsSpace(rune(char)):
			l.consumeWhitespace()
		case unicode.IsDigit(rune(char)):
			tokens = append(tokens, l.consumeNumber())
		case unicode.IsLetter(rune(char)):
			token := l.consumeIdentifier()
			// Lookup directo en mapa estático
			if tokenType, exists := keywords[token.Value]; exists {
				token.Type = tokenType
			}
			tokens = append(tokens, token)
		case char == '"' || char == '\'':
			tokens = append(tokens, l.consumeString())
		default:
			tokens = append(tokens, l.consumeOperatorOrSymbol())
		}
	}
	
	return tokens
}

func (l *Lexer) consumeWhitespace() {
	for l.position < len(l.input) {
		char := l.input[l.position]
		if !unicode.IsSpace(rune(char)) {
			break
		}
		if char == '\n' {
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
	
	// Optimización: avanzar directamente sin conversiones innecesarias
	for l.position < len(l.input) {
		char := l.input[l.position]
		if !unicode.IsDigit(rune(char)) && char != '.' {
			break
		}
		l.position++
		l.column++
	}
	
	// Verificar si hay letras después (error)
	if l.position < len(l.input) && unicode.IsLetter(rune(l.input[l.position])) {
		for l.position < len(l.input) && (unicode.IsLetter(rune(l.input[l.position])) || unicode.IsDigit(rune(l.input[l.position]))) {
			l.position++
			l.column++
		}
		
		return Token{
			Type:     UNKNOWN,
			Value:    l.input[start:l.position], // Slicing directo - máxima eficiencia
			Position: start,
			Line:     l.line,
			Column:   startCol,
		}
	}
	
	return Token{
		Type:     NUMBER,
		Value:    l.input[start:l.position], // Slicing directo
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *Lexer) consumeIdentifier() Token {
	start := l.position
	startCol := l.column
	
	// Loop optimizado con condiciones simplificadas
	for l.position < len(l.input) {
		char := l.input[l.position]
		if !unicode.IsLetter(rune(char)) && !unicode.IsDigit(rune(char)) && char != '_' {
			break
		}
		l.position++
		l.column++
	}
	
	return Token{
		Type:     IDENTIFIER,
		Value:    l.input[start:l.position], // Slicing directo
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *Lexer) consumeString() Token {
	start := l.position
	startCol := l.column
	quote := l.input[l.position]
	l.position++
	l.column++
	
	// Loop optimizado para strings
	for l.position < len(l.input) && l.input[l.position] != quote {
		if l.input[l.position] == '\\' && l.position+1 < len(l.input) {
			l.position += 2 // Skip escape + next char
			l.column += 2
		} else {
			l.position++
			l.column++
		}
	}
	
	if l.position < len(l.input) {
		l.position++
		l.column++
	}
	
	return Token{
		Type:     STRING,
		Value:    l.input[start:l.position], // Slicing directo
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *Lexer) consumeOperatorOrSymbol() Token {
	start := l.position
	startCol := l.column
	
	// Optimización: verificar operadores de 3 caracteres primero con array estático
	if l.position+2 < len(l.input) {
		threeChar := l.input[l.position : l.position+3]
		for _, op := range threeCharOps {
			if threeChar == op {
				l.position += 3
				l.column += 3
				return Token{
					Type:     COMPARISON,
					Value:    threeChar,
					Position: start,
					Line:     l.line,
					Column:   startCol,
				}
			}
		}
	}
	
	// Verificar operadores de 2 caracteres con array estático
	if l.position+1 < len(l.input) {
		twoChar := l.input[l.position : l.position+2]
		for _, op := range twoCharOps {
			if twoChar == op {
				l.position += 2
				l.column += 2
				return Token{
					Type:     getOperatorType(twoChar),
					Value:    twoChar,
					Position: start,
					Line:     l.line,
					Column:   startCol,
				}
			}
		}
	}
	
	// Operador de un caracter con switch optimizado
	char := l.input[l.position]
	l.position++
	l.column++
	
	tokenType := getSingleCharType(char)
	return Token{
		Type:     tokenType,
		Value:    string(char),
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

// Funciones optimizadas con switch statements
func getOperatorType(op string) TokenType {
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

func getSingleCharType(char byte) TokenType {
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