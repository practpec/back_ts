package main

import (
	"strconv"
)

type Parser struct {
	tokens   []Token
	position int
	errors   []string
}

// Pool de strings para reutilizar mensajes de error comunes
var (
	errorEOF = "Se llegó al final del código inesperadamente"
	errorSemicolon = "Se esperaba punto y coma"
	errorIdentifier = "Se esperaba identificador"
	errorAssignment = "Se esperaba operador de asignación"
	errorValue = "Se esperaba valor"
)

func NewParser(tokens []Token) *Parser {
	// Pre-filtrar tokens de whitespace una sola vez
	filteredTokens := make([]Token, 0, len(tokens))
	for i := range tokens {
		if tokens[i].Type != WHITESPACE {
			filteredTokens = append(filteredTokens, tokens[i])
		}
	}
	
	return &Parser{
		tokens:   filteredTokens,
		position: 0,
		errors:   make([]string, 0, 4), // Pre-allocar con capacidad estimada
	}
}

func (p *Parser) Parse() []string {
	for p.position < len(p.tokens) {
		token := &p.tokens[p.position]
		
		switch token.Type {
		case FOR:
			p.parseForStatement()
		case DO:
			p.parseDoWhileStatement()
		case KEYWORD, TYPE:
			p.parseVariableDeclaration()
		default:
			p.position++
		}
	}
	return p.errors
}

// Función inline optimizada para agregar errores
func (p *Parser) addError(message string) {
	p.errors = append(p.errors, message)
}

// Función optimizada para obtener token actual
func (p *Parser) currentToken() *Token {
	if p.position >= len(p.tokens) {
		return nil
	}
	return &p.tokens[p.position]
}

// Función optimizada para consumir tokens
func (p *Parser) consume(expectedType TokenType) bool {
	token := p.currentToken()
	if token == nil {
		p.addError(errorEOF)
		return false
	}
	
	if token.Type != expectedType {
		// Usar strconv optimizado para convertir números
		p.addError("Se esperaba " + string(expectedType) + " pero se encontró " + string(token.Type) + 
			" '" + token.Value + "' en línea " + strconv.Itoa(token.Line) + 
			", columna " + strconv.Itoa(token.Column))
		return false
	}
	
	p.position++
	return true
}

func (p *Parser) parseVariableDeclaration() {
	// Verificar palabra clave/tipo
	if p.currentToken() != nil && (p.currentToken().Type == KEYWORD || p.currentToken().Type == TYPE) {
		p.position++
	}
	
	// Nombre de variable
	if !p.consume(IDENTIFIER) {
		return
	}
	
	// Verificar declaración de tipo TypeScript opcional
	if p.currentToken() != nil && p.currentToken().Type == COLON {
		p.position++ // consume ':'
		if p.currentToken() != nil && (p.currentToken().Type == TYPE || p.currentToken().Type == IDENTIFIER) {
			p.position++
		} else {
			p.addError("Se esperaba tipo después de ':'")
			return
		}
	}
	
	// Operador de asignación
	if !p.consume(ASSIGNMENT) {
		return
	}
	
	// Valor
	token := p.currentToken()
	if token == nil {
		p.addError(errorValue)
		return
	}
	
	switch token.Type {
	case NUMBER, IDENTIFIER:
		p.position++
	case UNKNOWN:
		p.addError("Número mal formado '" + token.Value + "' en línea " + 
			strconv.Itoa(token.Line) + ", columna " + strconv.Itoa(token.Column))
		p.position++
		return
	default:
		p.addError("Se esperaba número o identificador, se encontró " + string(token.Type))
		return
	}
	
	// Punto y coma opcional
	if p.currentToken() != nil && p.currentToken().Type == SEMICOLON {
		p.position++
	} else {
		nextToken := p.currentToken()
		if nextToken != nil && nextToken.Line == token.Line {
			p.addError("Se esperaba punto y coma o salto de línea después de la declaración en línea " + 
				strconv.Itoa(token.Line))
		}
	}
}

func (p *Parser) parseForStatement() {
	if !p.consume(FOR) { return }
	if !p.consume(LPAREN) { return }
	
	p.parseInitialization()
	if !p.consume(SEMICOLON) { return }
	
	p.parseCondition()
	if !p.consume(SEMICOLON) { return }
	
	p.parseIncrement()
	if !p.consume(RPAREN) { return }
	if !p.consume(LBRACE) { return }
	
	p.parseStatements()
	if !p.consume(RBRACE) { return }
}

func (p *Parser) parseInitialization() {
	token := p.currentToken()
	if token == nil {
		p.addError("Se esperaba declaración de variable en inicialización")
		return
	}
	
	if token.Type == KEYWORD || token.Type == TYPE {
		p.position++
	}
	
	if !p.consume(IDENTIFIER) { return }
	
	// Tipo TypeScript opcional
	if p.currentToken() != nil && p.currentToken().Type == COLON {
		p.position++
		if p.currentToken() != nil && (p.currentToken().Type == TYPE || p.currentToken().Type == IDENTIFIER) {
			p.position++
		}
	}
	
	if !p.consume(ASSIGNMENT) { return }
	
	token = p.currentToken()
	if token == nil {
		p.addError("Se esperaba valor en inicialización")
		return
	}
	
	if token.Type == NUMBER || token.Type == IDENTIFIER {
		p.position++
	} else {
		p.addError("Se esperaba número o identificador en inicialización, se encontró " + string(token.Type))
	}
}

func (p *Parser) parseCondition() {
	if !p.consume(IDENTIFIER) { return }
	if !p.consume(COMPARISON) { return }
	
	token := p.currentToken()
	if token == nil {
		p.addError("Se esperaba valor en condición")
		return
	}
	
	if token.Type == NUMBER || token.Type == IDENTIFIER {
		p.position++
	} else {
		p.addError("Se esperaba número o identificador en condición, se encontró " + string(token.Type))
	}
}

func (p *Parser) parseIncrement() {
	token := p.currentToken()
	if token == nil {
		p.addError("Se esperaba expresión de incremento")
		return
	}
	
	if token.Type == INCREMENT {
		p.position++
		if !p.consume(IDENTIFIER) { return }
	} else if token.Type == IDENTIFIER {
		p.position++
		nextToken := p.currentToken()
		if nextToken == nil {
			p.addError("Se esperaba operador de incremento")
			return
		}
		
		if nextToken.Type == INCREMENT {
			p.position++
		} else if nextToken.Type == ASSIGNMENT {
			p.position++
			valueToken := p.currentToken()
			if valueToken == nil {
				p.addError("Se esperaba valor después del operador de asignación")
				return
			}
			if valueToken.Type == NUMBER || valueToken.Type == IDENTIFIER {
				p.position++
			} else {
				p.addError("Se esperaba número o identificador después del operador de asignación")
			}
		} else {
			p.addError("Se esperaba operador de incremento o asignación")
		}
	} else {
		p.addError("Se esperaba identificador o operador de incremento")
	}
}

func (p *Parser) parseStatements() {
	for p.currentToken() != nil && p.currentToken().Type != RBRACE {
		token := p.currentToken()
		
		if token.Type == IDENTIFIER || token.Type == KEYWORD {
			p.parseStatement()
		} else {
			p.position++
		}
	}
}

func (p *Parser) parseStatement() {
	token := p.currentToken()
	if token == nil { return }
	
	if token.Type == IDENTIFIER || token.Type == KEYWORD {
		p.position++
		
		// Acceso a propiedades con punto
		if p.currentToken() != nil && p.currentToken().Value == "." {
			p.position++
			if p.currentToken() != nil && p.currentToken().Type == IDENTIFIER {
				p.position++
			}
		}
		
		// Llamada a función
		if p.currentToken() != nil && p.currentToken().Type == LPAREN {
			p.position++
			
			// Argumentos simplificados
			for p.currentToken() != nil && p.currentToken().Type != RPAREN {
				if p.currentToken().Type == UNKNOWN {
					p.addError("Token inválido '" + p.currentToken().Value + 
						"' en expresión en línea " + strconv.Itoa(p.currentToken().Line) + 
						", columna " + strconv.Itoa(p.currentToken().Column))
				}
				p.position++
			}
			
			if p.currentToken() != nil && p.currentToken().Type == RPAREN {
				p.position++
			}
		} else if p.currentToken() != nil && p.currentToken().Type == ASSIGNMENT {
			p.position++
			p.parseExpression()
		}
		
		// Punto y coma opcional
		if p.currentToken() != nil && p.currentToken().Type == SEMICOLON {
			p.position++
		}
	}
}

func (p *Parser) parseExpression() {
	lastTokenType := UNKNOWN
	
	for p.currentToken() != nil {
		currentToken := p.currentToken()
		
		// Verificar si es un token válido en expresión
		if currentToken.Type != NUMBER && currentToken.Type != IDENTIFIER && 
		   currentToken.Type != OPERATOR && currentToken.Type != UNKNOWN {
			break
		}
		
		if currentToken.Type == UNKNOWN {
			p.addError("Token inválido '" + currentToken.Value + 
				"' en expresión en línea " + strconv.Itoa(currentToken.Line) + 
				", columna " + strconv.Itoa(currentToken.Column))
		}
		
		// Verificar secuencias inválidas usando switch optimizado
		switch {
		case lastTokenType == NUMBER && currentToken.Type == IDENTIFIER:
			p.addError("Error de sintaxis: número seguido de identificador sin operador en línea " + 
				strconv.Itoa(currentToken.Line))
		case lastTokenType == IDENTIFIER && currentToken.Type == NUMBER:
			p.addError("Error de sintaxis: identificador seguido de número sin operador en línea " + 
				strconv.Itoa(currentToken.Line))
		case lastTokenType == NUMBER && currentToken.Type == NUMBER:
			p.addError("Error de sintaxis: dos números consecutivos sin operador en línea " + 
				strconv.Itoa(currentToken.Line))
		case lastTokenType == IDENTIFIER && currentToken.Type == IDENTIFIER:
			p.addError("Error de sintaxis: dos identificadores consecutivos sin operador en línea " + 
				strconv.Itoa(currentToken.Line))
		}
		
		lastTokenType = currentToken.Type
		p.position++
		
		// Verificar delimitadores
		if p.currentToken() != nil {
			tokenType := p.currentToken().Type
			if tokenType == SEMICOLON || tokenType == RBRACE || tokenType == RPAREN {
				break
			}
		}
	}
}

func (p *Parser) parseDoWhileStatement() {
	if !p.consume(DO) { return }
	if !p.consume(LBRACE) { return }
	
	p.parseStatements()
	
	if !p.consume(RBRACE) { return }
	if !p.consume(WHILE) { return }
	if !p.consume(LPAREN) { return }
	
	p.parseCondition()
	
	if !p.consume(RPAREN) { return }
	if !p.consume(SEMICOLON) { return }
}