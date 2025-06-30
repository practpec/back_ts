package main

import (
	"fmt"
)

type ParserJ struct {
	tokens   []TokenJ
	position int
	errors   []string
}

func NewParserJ(tokens []TokenJ) *ParserJ {
	// Filter out whitespace tokens for parsing
	var filteredTokens []TokenJ
	for _, token := range tokens {
		if token.Type != WHITESPACE_J {
			filteredTokens = append(filteredTokens, token)
		}
	}

	return &ParserJ{
		tokens:   filteredTokens,
		position: 0,
		errors:   []string{},
	}
}

func (p *ParserJ) addError(message string) {
	p.errors = append(p.errors, message)
}

func (p *ParserJ) Parse() []string {
	// Analizar estructura de clase Java
	p.parseClass()
	return p.errors
}

func (p *ParserJ) currentToken() *TokenJ {
	if p.position >= len(p.tokens) {
		return nil
	}
	return &p.tokens[p.position]
}

func (p *ParserJ) consume(expectedType TokenTypeJ) bool {
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

func (p *ParserJ) parseClass() {
	// public class NombreClase {
	if !p.consume(PUBLIC_J) {
		return
	}

	if !p.consume(CLASS_J) {
		return
	}

	// Nombre de la clase
	if !p.consume(IDENTIFIER_J) {
		return
	}

	// Llave de apertura de clase
	if !p.consume(LBRACE_J) {
		return
	}

	// Método main
	p.parseMainMethod()

	// Llave de cierre de clase
	if !p.consume(RBRACE_J) {
		return
	}
}

func (p *ParserJ) parseMainMethod() {
	// public static void main(String[] args) {
	if !p.consume(PUBLIC_J) {
		return
	}

	if !p.consume(STATIC_J) {
		return
	}

	if !p.consume(VOID_J) {
		return
	}

	if !p.consume(MAIN_J) {
		return
	}

	if !p.consume(LPAREN_J) {
		return
	}

	// String[] args - Manejo correcto de arrays
	if !p.consume(TYPE_J) { // String
		return
	}

	if !p.consume(LBRACKET_J) { // [
		return
	}

	if !p.consume(RBRACKET_J) { // ]
		return
	}

	// args debe ser IDENTIFIER, no KEYWORD
	token := p.currentToken()
	if token == nil {
		p.addError("Se esperaba nombre de parámetro después de String[]")
		return
	}
	
	if token.Type == IDENTIFIER_J || (token.Type == KEYWORD_J && token.Value == "args") {
		p.position++ // Aceptar tanto IDENTIFIER como "args" por compatibilidad
	} else {
		p.addError(fmt.Sprintf("Se esperaba identificador para parámetro, se encontró %s '%s'", token.Type, token.Value))
		return
	}

	if !p.consume(RPAREN_J) {
		return
	}

	// Llave de apertura del método
	if !p.consume(LBRACE_J) {
		return
	}

	// Cuerpo del método
	p.parseMethodBody()

	// Llave de cierre del método
	if !p.consume(RBRACE_J) {
		return
	}
}

func (p *ParserJ) parseMethodBody() {
	for p.currentToken() != nil && p.currentToken().Type != RBRACE_J {
		token := p.currentToken()

		if token.Type == TYPE_J {
			// Declaración de variable
			p.parseVariableDeclaration()
		} else if token.Type == IF_J {
			// Estructura if
			p.parseIfStatement()
		} else if token.Type == IDENTIFIER_J {
			// Posible llamada a método o asignación
			p.parseMethodCallOrAssignment()
		} else {
			p.position++
		}
	}
}

func (p *ParserJ) parseVariableDeclaration() {
	// Tipo variable = valor;
	if !p.consume(TYPE_J) {
		return
	}

	if !p.consume(IDENTIFIER_J) {
		return
	}

	if !p.consume(ASSIGNMENT_J) {
		return
	}

	// Valor (puede ser número, string, o expresión)
	token := p.currentToken()
	if token == nil {
		p.addError("Se esperaba valor en declaración de variable")
		return
	}

	if token.Type == NUMBER_J || token.Type == STRING_J || token.Type == IDENTIFIER_J {
		p.position++
	} else {
		p.addError(fmt.Sprintf("Se esperaba valor válido en declaración, se encontró %s", token.Type))
		return
	}

	if !p.consume(SEMICOLON_J) {
		return
	}
}

func (p *ParserJ) parseIfStatement() {
	// if (condición) { bloque }
	if !p.consume(IF_J) {
		return
	}

	if !p.consume(LPAREN_J) {
		return
	}

	// Analizar condición
	p.parseCondition()

	if !p.consume(RPAREN_J) {
		return
	}

	if !p.consume(LBRACE_J) {
		return
	}

	// Cuerpo del if
	p.parseIfBody()

	if !p.consume(RBRACE_J) {
		return
	}
}

func (p *ParserJ) parseCondition() {
	// Puede ser: variable > numero, variable.equals("string"), etc.
	token := p.currentToken()
	if token == nil {
		p.addError("Se esperaba condición en if")
		return
	}

	if token.Type == IDENTIFIER_J {
		p.position++

		// Verificar si hay .equals() o operador de comparación
		nextToken := p.currentToken()
		if nextToken != nil {
			if nextToken.Type == DOT_J {
				// Llamada a método como .equals()
				p.position++ // consume '.'
				if p.currentToken() != nil && p.currentToken().Value == "equals" {
					p.position++ // consume 'equals'
					if !p.consume(LPAREN_J) {
						return
					}
					// Argumento del equals
					if p.currentToken() != nil && p.currentToken().Type == STRING_J {
						p.position++
					}
					if !p.consume(RPAREN_J) {
						return
					}
				}
			} else if nextToken.Type == COMPARISON_J {
				// Operador de comparación
				p.position++ // consume operador
				// Valor a comparar
				if p.currentToken() != nil && (p.currentToken().Type == NUMBER_J || p.currentToken().Type == IDENTIFIER_J) {
					p.position++
				}
			}
		}
	}
}

func (p *ParserJ) parseIfBody() {
	// Puede contener llamadas a System.out.println() u otras declaraciones
	for p.currentToken() != nil && p.currentToken().Type != RBRACE_J {
		if p.currentToken().Type == IDENTIFIER_J && p.currentToken().Value == "System" {
			p.parseSystemOut()
		} else {
			p.position++
		}
	}
}

func (p *ParserJ) parseMethodCallOrAssignment() {
	// System.out.println() o asignación
	if p.currentToken() != nil && p.currentToken().Value == "System" {
		p.parseSystemOut()
	} else {
		// Asignación simple
		p.position++ // identifier
		if p.currentToken() != nil && p.currentToken().Type == ASSIGNMENT_J {
			p.position++ // =
			if p.currentToken() != nil {
				p.position++ // valor
			}
			if p.currentToken() != nil && p.currentToken().Type == SEMICOLON_J {
				p.position++ // ;
			}
		}
	}
}

func (p *ParserJ) parseSystemOut() {
	// System.out.println("mensaje");
	if !p.consume(IDENTIFIER_J) { // System
		return
	}

	if !p.consume(DOT_J) {
		return
	}

	if p.currentToken() != nil && p.currentToken().Value == "out" {
		p.position++
	} else {
		p.addError("Se esperaba 'out' después de 'System.'")
		return
	}

	if !p.consume(DOT_J) {
		return
	}

	if p.currentToken() != nil && p.currentToken().Value == "println" {
		p.position++
	} else {
		p.addError("Se esperaba 'println' después de 'System.out.'")
		return
	}

	if !p.consume(LPAREN_J) {
		return
	}

	// Argumento (string o variable)
	token := p.currentToken()
	if token != nil && (token.Type == STRING_J || token.Type == IDENTIFIER_J) {
		p.position++
	} else {
		p.addError("Se esperaba argumento en System.out.println()")
		return
	}

	if !p.consume(RPAREN_J) {
		return
	}

	if !p.consume(SEMICOLON_J) {
		return
	}
}