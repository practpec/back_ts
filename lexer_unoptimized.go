package main

import (
	"unicode"
)

// Versión NO optimizada del lexer que usa concatenación de strings
// y múltiples conversiones innecesarias para demostrar el impacto en rendimiento

type LexerUnoptimized struct {
	input    string
	position int
	line     int
	column   int
}

func NewLexerUnoptimized(input string) *LexerUnoptimized {
	return &LexerUnoptimized{
		input:    input,
		position: 0,
		line:     1,
		column:   1,
	}
}

func (l *LexerUnoptimized) TokenizeUnoptimized() []Token {
	var tokens []Token
	
	// Versión ineficiente: crear el map en cada llamada en lugar de usar una variable global
	keywords := l.createKeywordsMapInefficiently()
	
	for l.position < len(l.input) {
		if unicode.IsSpace(rune(l.input[l.position])) {
			l.consumeWhitespaceUnoptimized()
			continue
		}
		
		// Números - versión ineficiente
		if unicode.IsDigit(rune(l.input[l.position])) {
			token := l.consumeNumberUnoptimized()
			tokens = append(tokens, token)
			continue
		}
		
		// Identificadores y palabras clave - versión ineficiente
		if unicode.IsLetter(rune(l.input[l.position])) {
			token := l.consumeIdentifierUnoptimized()
			if tokenType, exists := keywords[token.Value]; exists {
				token.Type = tokenType
			}
			tokens = append(tokens, token)
			continue
		}
		
		// Strings - versión ineficiente
		if l.input[l.position] == '"' || l.input[l.position] == '\'' {
			token := l.consumeStringUnoptimized()
			tokens = append(tokens, token)
			continue
		}
		
		// Operadores y símbolos - versión ineficiente
		token := l.consumeOperatorOrSymbolUnoptimized()
		tokens = append(tokens, token)
	}
	
	return tokens
}

// Crear keywords map de forma ineficiente (en cada llamada)
func (l *LexerUnoptimized) createKeywordsMapInefficiently() map[string]TokenType {
	// Ineficiente: crear el map cada vez en lugar de usar una variable global
	keywords := make(map[string]TokenType)
	
	// Ineficiente: usar concatenación de strings en lugar de literales
	forWord := "f" + "o" + "r"
	doWord := "d" + "o"
	whileWord := "w" + "h" + "i" + "l" + "e"
	letWord := "l" + "e" + "t"
	constWord := "c" + "o" + "n" + "s" + "t"
	varWord := "v" + "a" + "r"
	intWord := "i" + "n" + "t"
	stringWord := "s" + "t" + "r" + "i" + "n" + "g"
	numberWord := "n" + "u" + "m" + "b" + "e" + "r"
	booleanWord := "b" + "o" + "o" + "l" + "e" + "a" + "n"
	consoleWord := "c" + "o" + "n" + "s" + "o" + "l" + "e"
	
	keywords[forWord] = FOR
	keywords[doWord] = DO
	keywords[whileWord] = WHILE
	keywords[letWord] = KEYWORD
	keywords[constWord] = KEYWORD
	keywords[varWord] = KEYWORD
	keywords[intWord] = TYPE
	keywords[stringWord] = TYPE
	keywords[numberWord] = TYPE
	keywords[booleanWord] = TYPE
	keywords[consoleWord] = KEYWORD
	
	return keywords
}

func (l *LexerUnoptimized) consumeWhitespaceUnoptimized() {
	for l.position < len(l.input) && unicode.IsSpace(rune(l.input[l.position])) {
		// Ineficiente: crear string temporal para cada caracter
		char := string(l.input[l.position])
		if char == "\n" {
			l.line++
			l.column = 1
		} else {
			l.column++
		}
		l.position++
	}
}

func (l *LexerUnoptimized) consumeNumberUnoptimized() Token {
	start := l.position
	startCol := l.column
	
	// Ineficiente: usar concatenación de strings en lugar de slicing
	number := ""
	
	for l.position < len(l.input) && (unicode.IsDigit(rune(l.input[l.position])) || l.input[l.position] == '.') {
		// Ineficiente: concatenar caracter por caracter
		number = number + string(l.input[l.position])
		l.position++
		l.column++
	}
	
	// Verificar si hay letras inmediatamente después del número (error)
	if l.position < len(l.input) && unicode.IsLetter(rune(l.input[l.position])) {
		// Ineficiente: continuar concatenando
		for l.position < len(l.input) && (unicode.IsLetter(rune(l.input[l.position])) || unicode.IsDigit(rune(l.input[l.position]))) {
			number = number + string(l.input[l.position])
			l.position++
			l.column++
		}
		
		return Token{
			Type:     UNKNOWN,
			Value:    number, // Usar la concatenación ineficiente
			Position: start,
			Line:     l.line,
			Column:   startCol,
		}
	}
	
	return Token{
		Type:     NUMBER,
		Value:    number, // Usar la concatenación ineficiente
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *LexerUnoptimized) consumeIdentifierUnoptimized() Token {
	start := l.position
	startCol := l.column
	
	// Ineficiente: usar concatenación de strings
	identifier := ""
	
	for l.position < len(l.input) && (unicode.IsLetter(rune(l.input[l.position])) || 
		unicode.IsDigit(rune(l.input[l.position])) || l.input[l.position] == '_') {
		// Ineficiente: concatenar caracter por caracter
		identifier = identifier + string(l.input[l.position])
		l.position++
		l.column++
	}
	
	return Token{
		Type:     IDENTIFIER,
		Value:    identifier, // Usar la concatenación ineficiente
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *LexerUnoptimized) consumeStringUnoptimized() Token {
	start := l.position
	startCol := l.column
	quote := l.input[l.position]
	
	// Ineficiente: usar concatenación
	str := string(quote) // Empezar con la comilla
	l.position++ // consume opening quote
	l.column++
	
	for l.position < len(l.input) && l.input[l.position] != quote {
		if l.input[l.position] == '\\' {
			// Ineficiente: concatenar escape character
			str = str + string(l.input[l.position])
			l.position++ // skip escape character
			l.column++
		}
		// Ineficiente: concatenar caracter por caracter
		str = str + string(l.input[l.position])
		l.position++
		l.column++
	}
	
	if l.position < len(l.input) {
		// Ineficiente: concatenar comilla final
		str = str + string(l.input[l.position])
		l.position++ // consume closing quote
		l.column++
	}
	
	return Token{
		Type:     STRING,
		Value:    str, // Usar la concatenación ineficiente
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *LexerUnoptimized) consumeOperatorOrSymbolUnoptimized() Token {
	start := l.position
	startCol := l.column
	char := l.input[l.position]
	
	l.position++
	l.column++
	
	// Ineficiente: crear strings temporales para comparaciones
	charStr := string(char)
	
	// Check for triple-character operators
	if l.position+1 < len(l.input) {
		// Ineficiente: concatenar strings para crear tokens de 3 caracteres
		secondChar := string(l.input[l.position])
		thirdChar := string(l.input[l.position+1])
		threeChar := charStr + secondChar + thirdChar
		
		if threeChar == ("=" + "=" + "=") || threeChar == ("!" + "=" + "=") {
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
		// Ineficiente: concatenar strings para crear tokens de 2 caracteres
		secondChar := string(l.input[l.position])
		twoChar := charStr + secondChar
		
		// Ineficiente: múltiples concatenaciones para comparar
		if twoChar == ("<" + "=") || twoChar == (">" + "=") || 
		   twoChar == ("=" + "=") || twoChar == ("!" + "=") || 
		   twoChar == ("+" + "+") || twoChar == ("-" + "-") || 
		   twoChar == ("+" + "=") || twoChar == ("-" + "=") {
			l.position++
			l.column++
			return Token{
				Type:     l.getOperatorTypeUnoptimized(twoChar),
				Value:    twoChar,
				Position: start,
				Line:     l.line,
				Column:   startCol,
			}
		}
	}
	
	// Single character operators/symbols
	tokenType := l.getSingleCharTypeUnoptimized(char)
	return Token{
		Type:     tokenType,
		Value:    charStr, // Usar string creado ineficientemente
		Position: start,
		Line:     l.line,
		Column:   startCol,
	}
}

func (l *LexerUnoptimized) getOperatorTypeUnoptimized(op string) TokenType {
	// Ineficiente: crear strings temporales para cada comparación
	lessThanEqual := "<" + "="
	greaterThanEqual := ">" + "="
	equalEqual := "=" + "="
	notEqual := "!" + "="
	lessThan := "<"
	greaterThan := ">"
	tripleEqual := "=" + "=" + "="
	tripleNotEqual := "!" + "=" + "="
	increment := "+" + "+"
	decrement := "-" + "-"
	plusEqual := "+" + "="
	minusEqual := "-" + "="
	equal := "="
	
	if op == lessThanEqual || op == greaterThanEqual || op == equalEqual || 
	   op == notEqual || op == lessThan || op == greaterThan || 
	   op == tripleEqual || op == tripleNotEqual {
		return COMPARISON
	} else if op == increment || op == decrement {
		return INCREMENT
	} else if op == plusEqual || op == minusEqual || op == equal {
		return ASSIGNMENT
	}
	return OPERATOR
}

func (l *LexerUnoptimized) getSingleCharTypeUnoptimized(char byte) TokenType {
	// Ineficiente: crear strings para cada comparación
	charStr := string(char)
	leftParen := "("
	rightParen := ")"
	leftBrace := "{"
	rightBrace := "}"
	semicolon := ";"
	colon := ":"
	equal := "="
	lessThan := "<"
	greaterThan := ">"
	plus := "+"
	minus := "-"
	multiply := "*"
	divide := "/"
	
	if charStr == leftParen {
		return LPAREN
	} else if charStr == rightParen {
		return RPAREN
	} else if charStr == leftBrace {
		return LBRACE
	} else if charStr == rightBrace {
		return RBRACE
	} else if charStr == semicolon {
		return SEMICOLON
	} else if charStr == colon {
		return COLON
	} else if charStr == equal {
		return ASSIGNMENT
	} else if charStr == lessThan || charStr == greaterThan {
		return COMPARISON
	} else if charStr == plus || charStr == minus || charStr == multiply || charStr == divide {
		return OPERATOR
	}
	return UNKNOWN
}