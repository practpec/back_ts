
package main

import (
	"fmt"
)

type Parser struct {
	tokens   []Token
	position int
	errors   []string
}

func (p *Parser) addError(message string) {
	p.errors = append(p.errors, message)
}
func (p *Parser) parseVariableDeclaration() {
	// Consumir palabra clave (let, const, var) o tipo (int, etc.)
	if p.currentToken() != nil && (p.currentToken().Type == KEYWORD || p.currentToken().Type == TYPE) {
		p.position++
	}
	
	// Nombre de variable
	if !p.consume(IDENTIFIER) {
		return
	}
	
	// Verificar si hay declaración de tipo TypeScript (: tipo)
	if p.currentToken() != nil && p.currentToken().Type == COLON {
		p.position++ // consume ':'
		// Consumir tipo
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
	
	// Valor - verificar que sea válido
	token := p.currentToken()
	if token == nil {
		p.addError("Se esperaba valor en declaración")
		return
	}
	
	if token.Type == NUMBER {
		p.position++
	} else if token.Type == IDENTIFIER {
		p.position++
	} else if token.Type == UNKNOWN {
		p.addError(fmt.Sprintf("Número mal formado '%s' en línea %d, columna %d", token.Value, token.Line, token.Column))
		p.position++
		return
	} else {
		p.addError(fmt.Sprintf("Se esperaba número o identificador, se encontró %s", token.Type))
		return
	}
	
	// Punto y coma (opcional - solo advertencia si falta)
	if p.currentToken() != nil && p.currentToken().Type == SEMICOLON {
		p.position++
	} else {
		// Verificar si la siguiente línea tiene contenido
		nextToken := p.currentToken()
		if nextToken != nil && nextToken.Line == token.Line {
			p.addError(fmt.Sprintf("Se esperaba punto y coma o salto de línea después de la declaración en línea %d", token.Line))
		}
	}
}

func NewParser(tokens []Token) *Parser {
	// Filter out whitespace tokens for parsing
	var filteredTokens []Token
	for _, token := range tokens {
		if token.Type != WHITESPACE {
			filteredTokens = append(filteredTokens, token)
		}
	}
	
	return &Parser{
		tokens:   filteredTokens,
		position: 0,
		errors:   []string{},
	}
}

func (p *Parser) Parse() []string {
	for p.currentToken() != nil {
		token := p.currentToken()
		if token.Type == FOR {
			p.parseForStatement()
		} else if token.Type == DO {
			p.parseDoWhileStatement()
		} else if token.Type == KEYWORD || token.Type == TYPE {
			// Parsear declaraciones de variables
			p.parseVariableDeclaration()
		} else {
			p.position++
		}
	}
	return p.errors
}

func (p *Parser) currentToken() *Token {
	if p.position >= len(p.tokens) {
		return nil
	}
	return &p.tokens[p.position]
}

func (p *Parser) consume(expectedType TokenType) bool {
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

func (p *Parser) parseForStatement() {
	// Verificar que comience con 'for'
	if !p.consume(FOR) {
		return
	}
	
	// Paréntesis de apertura
	if !p.consume(LPAREN) {
		return
	}
	
	// Inicialización: let/const/var identifier = value
	p.parseInitialization()
	
	// Punto y coma después de inicialización
	if !p.consume(SEMICOLON) {
		return
	}
	
	// Condición: identifier comparison value
	p.parseCondition()
	
	// Punto y coma después de condición
	if !p.consume(SEMICOLON) {
		return
	}
	
	// Incremento: identifier++ o ++identifier o identifier += value
	p.parseIncrement()
	
	// Paréntesis de cierre
	if !p.consume(RPAREN) {
		return
	}
	
	// Llave de apertura
	if !p.consume(LBRACE) {
		return
	}
	
	// Cuerpo del bucle (statements)
	p.parseStatements()
	
	// Llave de cierre
	if !p.consume(RBRACE) {
		return
	}
}

func (p *Parser) parseInitialization() {
	// Verificar declaración de variable
	token := p.currentToken()
	if token == nil {
		p.addError("Se esperaba declaración de variable en inicialización")
		return
	}
	
	// Puede ser let, const, var, o tipo como int
	if token.Type == KEYWORD || token.Type == TYPE {
		p.position++
	}
	
	// Nombre de variable
	if !p.consume(IDENTIFIER) {
		return
	}
	
	// Verificar si hay declaración de tipo TypeScript (: tipo)
	if p.currentToken() != nil && p.currentToken().Type == COLON {
		p.position++ // consume ':'
		// Consumir tipo
		if p.currentToken() != nil && (p.currentToken().Type == TYPE || p.currentToken().Type == IDENTIFIER) {
			p.position++
		}
	}
	
	// Operador de asignación
	if !p.consume(ASSIGNMENT) {
		return
	}
	
	// Valor
	token = p.currentToken()
	if token == nil {
		p.addError("Se esperaba valor en inicialización")
		return
	}
	
	if token.Type == NUMBER || token.Type == IDENTIFIER {
		p.position++
	} else {
		p.addError(fmt.Sprintf("Se esperaba número o identificador en inicialización, se encontró %s", token.Type))
	}
}

func (p *Parser) parseCondition() {
	// Identificador
	if !p.consume(IDENTIFIER) {
		return
	}
	
	// Operador de comparación
	if !p.consume(COMPARISON) {
		return
	}
	
	// Valor
	token := p.currentToken()
	if token == nil {
		p.addError("Se esperaba valor en condición")
		return
	}
	
	if token.Type == NUMBER || token.Type == IDENTIFIER {
		p.position++
	} else {
		p.addError(fmt.Sprintf("Se esperaba número o identificador en condición, se encontró %s", token.Type))
	}
}

func (p *Parser) parseIncrement() {
	token := p.currentToken()
	if token == nil {
		p.addError("Se esperaba expresión de incremento")
		return
	}
	
	// Puede ser ++i, i++, i+=1, etc.
	if token.Type == INCREMENT {
		p.position++
		if !p.consume(IDENTIFIER) {
			return
		}
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
			// Consumir valor
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
	// Parsear el cuerpo del bucle - simplificado para el ejemplo
	for p.currentToken() != nil && p.currentToken().Type != RBRACE {
		token := p.currentToken()
		
		// Ejemplo básico: console.log o system.out.println
		if token.Type == IDENTIFIER || token.Type == KEYWORD {
			p.parseStatement()
		} else {
			p.position++
		}
	}
}

func (p *Parser) parseStatement() {
	// Simplificado - puede ser una llamada a función como console.log() o asignación
	token := p.currentToken()
	if token == nil {
		return
	}
	
	// Identificador (console, system, variable, etc.)
	if token.Type == IDENTIFIER || token.Type == KEYWORD {
		p.position++
		
		// Puede tener punto para acceso a propiedades
		if p.currentToken() != nil && p.currentToken().Value == "." {
			p.position++ // consume '.'
			if p.currentToken() != nil && p.currentToken().Type == IDENTIFIER {
				p.position++ // consume método
			}
		}
		
		// Paréntesis para llamada a función
		if p.currentToken() != nil && p.currentToken().Type == LPAREN {
			p.position++
			
			// Argumentos (simplificado)
			for p.currentToken() != nil && p.currentToken().Type != RPAREN {
				if p.currentToken().Type == UNKNOWN {
					p.addError(fmt.Sprintf("Token inválido '%s' en expresión en línea %d, columna %d", 
						p.currentToken().Value, p.currentToken().Line, p.currentToken().Column))
				}
				p.position++
			}
			
			if p.currentToken() != nil && p.currentToken().Type == RPAREN {
				p.position++
			}
		} else if p.currentToken() != nil && p.currentToken().Type == ASSIGNMENT {
			// Es una asignación
			p.position++ // consume '='
			
			// Parsear expresión del lado derecho
			p.parseExpression()
		}
		
		// Punto y coma (opcional)
		if p.currentToken() != nil && p.currentToken().Type == SEMICOLON {
			p.position++
		}
	}
}

func (p *Parser) parseExpression() {
	// Parsear expresión simple (número, variable, o operación)
	lastTokenType := ""
	
	for p.currentToken() != nil && 
		(p.currentToken().Type == NUMBER || 
		 p.currentToken().Type == IDENTIFIER || 
		 p.currentToken().Type == OPERATOR ||
		 p.currentToken().Type == UNKNOWN) {
		
		currentToken := p.currentToken()
		
		if currentToken.Type == UNKNOWN {
			p.addError(fmt.Sprintf("Token inválido '%s' en expresión en línea %d, columna %d", 
				currentToken.Value, currentToken.Line, currentToken.Column))
		}
		
		// Verificar secuencias inválidas
		if lastTokenType == "NUMBER" && currentToken.Type == IDENTIFIER {
			p.addError(fmt.Sprintf("Error de sintaxis: número '%s' seguido de identificador '%s' sin operador en línea %d, columna %d", 
				p.tokens[p.position-1].Value, currentToken.Value, currentToken.Line, currentToken.Column))
		}
		
		if lastTokenType == "IDENTIFIER" && currentToken.Type == NUMBER {
			p.addError(fmt.Sprintf("Error de sintaxis: identificador '%s' seguido de número '%s' sin operador en línea %d, columna %d", 
				p.tokens[p.position-1].Value, currentToken.Value, currentToken.Line, currentToken.Column))
		}
		
		if lastTokenType == "NUMBER" && currentToken.Type == NUMBER {
			p.addError(fmt.Sprintf("Error de sintaxis: dos números consecutivos '%s' '%s' sin operador en línea %d, columna %d", 
				p.tokens[p.position-1].Value, currentToken.Value, currentToken.Line, currentToken.Column))
		}
		
		if lastTokenType == "IDENTIFIER" && currentToken.Type == IDENTIFIER {
			p.addError(fmt.Sprintf("Error de sintaxis: dos identificadores consecutivos '%s' '%s' sin operador en línea %d, columna %d", 
				p.tokens[p.position-1].Value, currentToken.Value, currentToken.Line, currentToken.Column))
		}
		
		lastTokenType = string(currentToken.Type)
		p.position++
		
		// Salir si llegamos a un delimitador
		if p.currentToken() != nil && 
			(p.currentToken().Type == SEMICOLON || 
			 p.currentToken().Type == RBRACE ||
			 p.currentToken().Type == RPAREN) {
			break
		}
	}
}

func (p *Parser) parseDoWhileStatement() {
	// Verificar que comience con 'do'
	if !p.consume(DO) {
		return
	}
	
	// Llave de apertura
	if !p.consume(LBRACE) {
		return
	}
	
	// Cuerpo del bucle
	p.parseStatements()
	
	// Llave de cierre
	if !p.consume(RBRACE) {
		return
	}
	
	// Palabra clave 'while'
	if !p.consume(WHILE) {
		return
	}
	
	// Paréntesis de apertura
	if !p.consume(LPAREN) {
		return
	}
	
	// Condición
	p.parseCondition()
	
	// Paréntesis de cierre
	if !p.consume(RPAREN) {
		return
	}
	
	// Punto y coma final
	if !p.consume(SEMICOLON) {
		return
	}
}