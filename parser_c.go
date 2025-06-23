package main

import (
	"fmt"
)

type ParserC struct {
	tokens   []TokenC
	position int
	errors   []string
}

func NewParserC(tokens []TokenC) *ParserC {
	// Filter out whitespace tokens for parsing
	var filteredTokens []TokenC
	for _, token := range tokens {
		if token.Type != WHITESPACE_C {
			filteredTokens = append(filteredTokens, token)
		}
	}

	return &ParserC{
		tokens:   filteredTokens,
		position: 0,
		errors:   []string{},
	}
}

func (p *ParserC) addError(message string) {
	p.errors = append(p.errors, message)
}

func (p *ParserC) Parse() []string {
	for p.currentToken() != nil {
		token := p.currentToken()
		if token.Type == FOR_C {
			p.parseForStatement()
		} else if token.Type == DO_C {
			p.parseDoWhileStatement()
		} else if token.Type == TYPE_C {
			// Parsear declaraciones de variables
			p.parseVariableDeclaration()
		} else {
			p.position++
		}
	}
	return p.errors
}

func (p *ParserC) currentToken() *TokenC {
	if p.position >= len(p.tokens) {
		return nil
	}
	return &p.tokens[p.position]
}

func (p *ParserC) consume(expectedType TokenTypeC) bool {
	token := p.currentToken()
	if token == nil {
		p.addError(fmt.Sprintf("Se esperaba %s pero se llegó al final del código", expectedType))
		return false
	}

	if token.Type != expectedType {
		p.addError(fmt.Sprintf("Se esperaba %s pero se encontró %s '%s' en línea %d, columna %d",
			expectedType, token.Type, token.Value, token.Line, token.Column))
		return false
	}

	p.position++
	return true
}

func (p *ParserC) parseVariableDeclaration() {
	// Consumir tipo (int, float, etc.)
	if p.currentToken() != nil && p.currentToken().Type == TYPE_C {
		p.position++
	}

	// Nombre de variable
	if !p.consume(IDENTIFIER_C) {
		return
	}

	// Operador de asignación
	if !p.consume(ASSIGNMENT_C) {
		return
	}

	// Valor - verificar que sea válido
	token := p.currentToken()
	if token == nil {
		p.addError("Se esperaba valor en declaración")
		return
	}

	if token.Type == NUMBER_C {
		p.position++
	} else if token.Type == IDENTIFIER_C {
		p.position++
	} else if token.Type == UNKNOWN_C {
		p.addError(fmt.Sprintf("Número mal formado '%s' en línea %d, columna %d", token.Value, token.Line, token.Column))
		p.position++
		return
	} else {
		p.addError(fmt.Sprintf("Se esperaba número o identificador, se encontró %s", token.Type))
		return
	}

	// Punto y coma obligatorio en C
	if !p.consume(SEMICOLON_C) {
		return
	}
}

func (p *ParserC) parseForStatement() {
	// Verificar que comience con 'for'
	if !p.consume(FOR_C) {
		return
	}

	// Paréntesis de apertura
	if !p.consume(LPAREN_C) {
		return
	}

	// Inicialización: tipo identifier = value
	p.parseInitialization()

	// Punto y coma después de inicialización
	if !p.consume(SEMICOLON_C) {
		return
	}

	// Condición: identifier comparison value
	p.parseCondition()

	// Punto y coma después de condición
	if !p.consume(SEMICOLON_C) {
		return
	}

	// Incremento: identifier++ o ++identifier o identifier += value
	p.parseIncrement()

	// Paréntesis de cierre
	if !p.consume(RPAREN_C) {
		return
	}

	// Llave de apertura
	if !p.consume(LBRACE_C) {
		return
	}

	// Cuerpo del bucle (statements)
	p.parseStatements()

	// Llave de cierre
	if !p.consume(RBRACE_C) {
		return
	}
}

func (p *ParserC) parseInitialization() {
	// Verificar declaración de variable
	token := p.currentToken()
	if token == nil {
		p.addError("Se esperaba declaración de variable en inicialización")
		return
	}

	// Debe ser un tipo como int, float, etc.
	if token.Type == TYPE_C {
		p.position++
	}

	// Nombre de variable
	if !p.consume(IDENTIFIER_C) {
		return
	}

	// Operador de asignación
	if !p.consume(ASSIGNMENT_C) {
		return
	}

	// Valor
	token = p.currentToken()
	if token == nil {
		p.addError("Se esperaba valor en inicialización")
		return
	}

	if token.Type == NUMBER_C || token.Type == IDENTIFIER_C {
		p.position++
	} else {
		p.addError(fmt.Sprintf("Se esperaba número o identificador en inicialización, se encontró %s", token.Type))
	}
}

func (p *ParserC) parseCondition() {
	// Identificador
	if !p.consume(IDENTIFIER_C) {
		return
	}

	// Operador de comparación
	if !p.consume(COMPARISON_C) {
		return
	}

	// Valor
	token := p.currentToken()
	if token == nil {
		p.addError("Se esperaba valor en condición")
		return
	}

	if token.Type == NUMBER_C || token.Type == IDENTIFIER_C {
		p.position++
	} else {
		p.addError(fmt.Sprintf("Se esperaba número o identificador en condición, se encontró %s", token.Type))
	}
}

func (p *ParserC) parseIncrement() {
	token := p.currentToken()
	if token == nil {
		p.addError("Se esperaba expresión de incremento")
		return
	}

	// Puede ser ++i, i++, i+=1, etc.
	if token.Type == INCREMENT_C {
		p.position++
		if !p.consume(IDENTIFIER_C) {
			return
		}
	} else if token.Type == IDENTIFIER_C {
		p.position++
		nextToken := p.currentToken()
		if nextToken == nil {
			p.addError("Se esperaba operador de incremento")
			return
		}

		if nextToken.Type == INCREMENT_C {
			p.position++
		} else if nextToken.Type == ASSIGNMENT_C {
			p.position++
			// Consumir valor
			valueToken := p.currentToken()
			if valueToken == nil {
				p.addError("Se esperaba valor después del operador de asignación")
				return
			}
			if valueToken.Type == NUMBER_C || valueToken.Type == IDENTIFIER_C {
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

func (p *ParserC) parseStatements() {
	// Parsear el cuerpo del bucle
	for p.currentToken() != nil && p.currentToken().Type != RBRACE_C {
		token := p.currentToken()

		// Declaraciones de variables o asignaciones
		if token.Type == TYPE_C {
			p.parseVariableDeclaration()
		} else if token.Type == IDENTIFIER_C || token.Type == KEYWORD_C {
			p.parseStatement()
		} else {
			p.position++
		}
	}
}

func (p *ParserC) parseStatement() {
	token := p.currentToken()
	if token == nil {
		return
	}

	// Identificador (variable, función, etc.)
	if token.Type == IDENTIFIER_C || token.Type == KEYWORD_C {
		p.position++

		// Puede tener punto para acceso a propiedades (no común en C básico)
		if p.currentToken() != nil && p.currentToken().Value == "." {
			p.position++ // consume '.'
			if p.currentToken() != nil && p.currentToken().Type == IDENTIFIER_C {
				p.position++ // consume método
			}
		}

		// Paréntesis para llamada a función
		if p.currentToken() != nil && p.currentToken().Type == LPAREN_C {
			p.position++

			// Argumentos (simplificado)
			for p.currentToken() != nil && p.currentToken().Type != RPAREN_C {
				if p.currentToken().Type == UNKNOWN_C {
					p.addError(fmt.Sprintf("Token inválido '%s' en expresión en línea %d, columna %d",
						p.currentToken().Value, p.currentToken().Line, p.currentToken().Column))
				}
				p.position++
			}

			if p.currentToken() != nil && p.currentToken().Type == RPAREN_C {
				p.position++
			}
		} else if p.currentToken() != nil && p.currentToken().Type == ASSIGNMENT_C {
			// Es una asignación
			p.position++ // consume '='

			// Parsear expresión del lado derecho
			p.parseExpression()
		}

		// Punto y coma obligatorio en C
		if p.currentToken() != nil && p.currentToken().Type == SEMICOLON_C {
			p.position++
		}
	}
}

func (p *ParserC) parseExpression() {
	// Parsear expresión simple (número, variable, o operación)
	lastTokenType := ""

	for p.currentToken() != nil &&
		(p.currentToken().Type == NUMBER_C ||
			p.currentToken().Type == IDENTIFIER_C ||
			p.currentToken().Type == OPERATOR_C ||
			p.currentToken().Type == UNKNOWN_C) {

		currentToken := p.currentToken()

		if currentToken.Type == UNKNOWN_C {
			p.addError(fmt.Sprintf("Token inválido '%s' en expresión en línea %d, columna %d",
				currentToken.Value, currentToken.Line, currentToken.Column))
		}

		// Verificar secuencias inválidas
		if lastTokenType == "NUMBER_C" && currentToken.Type == IDENTIFIER_C {
			p.addError(fmt.Sprintf("Error de sintaxis: número '%s' seguido de identificador '%s' sin operador en línea %d, columna %d",
				p.tokens[p.position-1].Value, currentToken.Value, currentToken.Line, currentToken.Column))
		}

		lastTokenType = string(currentToken.Type)
		p.position++

		// Salir si llegamos a un delimitador
		if p.currentToken() != nil &&
			(p.currentToken().Type == SEMICOLON_C ||
				p.currentToken().Type == RBRACE_C ||
				p.currentToken().Type == RPAREN_C) {
			break
		}
	}
}

func (p *ParserC) parseDoWhileStatement() {
	// Verificar que comience con 'do'
	if !p.consume(DO_C) {
		return
	}

	// Llave de apertura
	if !p.consume(LBRACE_C) {
		return
	}

	// Cuerpo del bucle
	p.parseStatements()

	// Llave de cierre
	if !p.consume(RBRACE_C) {
		return
	}

	// Palabra clave 'while'
	if !p.consume(WHILE_C) {
		return
	}

	// Paréntesis de apertura
	if !p.consume(LPAREN_C) {
		return
	}

	// Condición
	p.parseCondition()

	// Paréntesis de cierre
	if !p.consume(RPAREN_C) {
		return
	}

	// Punto y coma final obligatorio en C
	if !p.consume(SEMICOLON_C) {
		return
	}
}