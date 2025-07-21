package main

// Versión NO optimizada del parser que usa concatenación de strings
// y múltiples operaciones ineficientes para demostrar el impacto en rendimiento

type ParserUnoptimized struct {
	tokens   []Token
	position int
	errors   []string
}

func NewParserUnoptimized(tokens []Token) *ParserUnoptimized {
	// Ineficiente: filtrar tokens usando concatenación de strings
	var filteredTokens []Token
	for _, token := range tokens {
		tokenTypeStr := string(token.Type)
		whitespaceStr := "W" + "H" + "I" + "T" + "E" + "S" + "P" + "A" + "C" + "E"
		if tokenTypeStr != whitespaceStr {
			filteredTokens = append(filteredTokens, token)
		}
	}
	
	return &ParserUnoptimized{
		tokens:   filteredTokens,
		position: 0,
		errors:   []string{},
	}
}

func (p *ParserUnoptimized) addErrorUnoptimized(message string) {
	// Ineficiente: usar concatenación para agregar timestamp o prefijo
	errorPrefix := "E" + "R" + "R" + "O" + "R" + ": "
	fullMessage := errorPrefix + message
	p.errors = append(p.errors, fullMessage)
}

func (p *ParserUnoptimized) ParseUnoptimized() []string {
	for p.currentTokenUnoptimized() != nil {
		token := p.currentTokenUnoptimized()
		
		// Ineficiente: crear strings para comparación
		forStr := "F" + "O" + "R"
		doStr := "D" + "O"
		keywordStr := "K" + "E" + "Y" + "W" + "O" + "R" + "D"
		typeStr := "T" + "Y" + "P" + "E"
		
		tokenTypeStr := string(token.Type)
		
		if tokenTypeStr == forStr {
			p.parseForStatementUnoptimized()
		} else if tokenTypeStr == doStr {
			p.parseDoWhileStatementUnoptimized()
		} else if tokenTypeStr == keywordStr || tokenTypeStr == typeStr {
			p.parseVariableDeclarationUnoptimized()
		} else {
			p.position++
		}
	}
	return p.errors
}

func (p *ParserUnoptimized) currentTokenUnoptimized() *Token {
	if p.position >= len(p.tokens) {
		return nil
	}
	return &p.tokens[p.position]
}

func (p *ParserUnoptimized) consumeUnoptimized(expectedType TokenType) bool {
	token := p.currentTokenUnoptimized()
	if token == nil {
		// Ineficiente: concatenar strings para mensaje de error
		expectedStr := string(expectedType)
		message := "Se esperaba " + expectedStr + " pero se llegó al final del código"
		p.addErrorUnoptimized(message)
		return false
	}
	
	// Ineficiente: convertir tipos a strings para comparar
	expectedStr := string(expectedType)
	actualStr := string(token.Type)
	
	if actualStr != expectedStr {
		// Ineficiente: múltiples concatenaciones para mensaje de error
		lineStr := p.intToStringInefficiently(token.Line)
		colStr := p.intToStringInefficiently(token.Column)
		message := "Se esperaba " + expectedStr + " pero se encontró " + actualStr + " '" + token.Value + "' en línea " + lineStr + ", columna " + colStr
		p.addErrorUnoptimized(message)
		return false
	}
	
	p.position++
	return true
}

// Función ineficiente para convertir int a string
func (p *ParserUnoptimized) intToStringInefficiently(num int) string {
	// Ineficiente: convertir dígito por dígito usando concatenación
	if num == 0 {
		return "0"
	}
	
	if num < 0 {
		return "-" + p.intToStringInefficiently(-num)
	}
	
	result := ""
	for num > 0 {
		digit := num % 10
		digitStr := p.digitToStringInefficiently(digit)
		result = digitStr + result
		num = num / 10
	}
	return result
}

func (p *ParserUnoptimized) digitToStringInefficiently(digit int) string {
	// Ineficiente: usar concatenación en lugar de lookup table
	switch digit {
	case 0:
		return "0"
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	case 5:
		return "5"
	case 6:
		return "6"
	case 7:
		return "7"
	case 8:
		return "8"
	case 9:
		return "9"
	default:
		return "?"
	}
}

func (p *ParserUnoptimized) parseVariableDeclarationUnoptimized() {
	// Ineficiente: crear strings para tipos
	keywordType := "K" + "E" + "Y" + "W" + "O" + "R" + "D"
	typeType := "T" + "Y" + "P" + "E"
	
	// Consumir palabra clave (let, const, var) o tipo (int, etc.)
	if p.currentTokenUnoptimized() != nil {
		tokenTypeStr := string(p.currentTokenUnoptimized().Type)
		if tokenTypeStr == keywordType || tokenTypeStr == typeType {
			p.position++
		}
	}
	
	// Nombre de variable
	identifierType := "I" + "D" + "E" + "N" + "T" + "I" + "F" + "I" + "E" + "R"
	if !p.consumeUnoptimized(TokenType(identifierType)) {
		return
	}
	
	// Verificar si hay declaración de tipo TypeScript (: tipo)
	colonType := "C" + "O" + "L" + "O" + "N"
	if p.currentTokenUnoptimized() != nil {
		tokenTypeStr := string(p.currentTokenUnoptimized().Type)
		if tokenTypeStr == colonType {
			p.position++ // consume ':'
			// Consumir tipo
			if p.currentTokenUnoptimized() != nil {
				currentTypeStr := string(p.currentTokenUnoptimized().Type)
				if currentTypeStr == typeType || currentTypeStr == identifierType {
					p.position++
				} else {
					p.addErrorUnoptimized("Se esperaba tipo después de ':'")
					return
				}
			}
		}
	}
	
	// Operador de asignación
	assignmentType := "A" + "S" + "S" + "I" + "G" + "N" + "M" + "E" + "N" + "T"
	if !p.consumeUnoptimized(TokenType(assignmentType)) {
		return
	}
	
	// Valor - verificar que sea válido
	token := p.currentTokenUnoptimized()
	if token == nil {
		p.addErrorUnoptimized("Se esperaba valor en declaración")
		return
	}
	
	numberType := "N" + "U" + "M" + "B" + "E" + "R"
	unknownType := "U" + "N" + "K" + "N" + "O" + "W" + "N"
	tokenTypeStr := string(token.Type)
	
	if tokenTypeStr == numberType {
		p.position++
	} else if tokenTypeStr == identifierType {
		p.position++
	} else if tokenTypeStr == unknownType {
		// Ineficiente: múltiples concatenaciones para mensaje de error
		lineStr := p.intToStringInefficiently(token.Line)
		colStr := p.intToStringInefficiently(token.Column)
		message := "Número mal formado '" + token.Value + "' en línea " + lineStr + ", columna " + colStr
		p.addErrorUnoptimized(message)
		p.position++
		return
	} else {
		message := "Se esperaba número o identificador, se encontró " + tokenTypeStr
		p.addErrorUnoptimized(message)
		return
	}
	
	// Punto y coma (opcional)
	semicolonType := "S" + "E" + "M" + "I" + "C" + "O" + "L" + "O" + "N"
	if p.currentTokenUnoptimized() != nil {
		tokenTypeStr := string(p.currentTokenUnoptimized().Type)
		if tokenTypeStr == semicolonType {
			p.position++
		} else {
			// Verificar si la siguiente línea tiene contenido
			nextToken := p.currentTokenUnoptimized()
			if nextToken != nil && nextToken.Line == token.Line {
				lineStr := p.intToStringInefficiently(token.Line)
				message := "Se esperaba punto y coma o salto de línea después de la declaración en línea " + lineStr
				p.addErrorUnoptimized(message)
			}
		}
	}
}

func (p *ParserUnoptimized) parseForStatementUnoptimized() {
	// Ineficiente: crear string para tipo FOR
	forType := "F" + "O" + "R"
	if !p.consumeUnoptimized(TokenType(forType)) {
		return
	}
	
	// Paréntesis de apertura
	lparenType := "L" + "P" + "A" + "R" + "E" + "N"
	if !p.consumeUnoptimized(TokenType(lparenType)) {
		return
	}
	
	// Inicialización
	p.parseInitializationUnoptimized()
	
	// Punto y coma después de inicialización
	semicolonType := "S" + "E" + "M" + "I" + "C" + "O" + "L" + "O" + "N"
	if !p.consumeUnoptimized(TokenType(semicolonType)) {
		return
	}
	
	// Condición
	p.parseConditionUnoptimized()
	
	// Punto y coma después de condición
	if !p.consumeUnoptimized(TokenType(semicolonType)) {
		return
	}
	
	// Incremento
	p.parseIncrementUnoptimized()
	
	// Paréntesis de cierre
	rparenType := "R" + "P" + "A" + "R" + "E" + "N"
	if !p.consumeUnoptimized(TokenType(rparenType)) {
		return
	}
	
	// Llave de apertura
	lbraceType := "L" + "B" + "R" + "A" + "C" + "E"
	if !p.consumeUnoptimized(TokenType(lbraceType)) {
		return
	}
	
	// Cuerpo del bucle
	p.parseStatementsUnoptimized()
	
	// Llave de cierre
	rbraceType := "R" + "B" + "R" + "A" + "C" + "E"
	if !p.consumeUnoptimized(TokenType(rbraceType)) {
		return
	}
}

func (p *ParserUnoptimized) parseInitializationUnoptimized() {
	token := p.currentTokenUnoptimized()
	if token == nil {
		p.addErrorUnoptimized("Se esperaba declaración de variable en inicialización")
		return
	}
	
	// Ineficiente: crear strings para tipos
	keywordType := "K" + "E" + "Y" + "W" + "O" + "R" + "D"
	typeType := "T" + "Y" + "P" + "E"
	tokenTypeStr := string(token.Type)
	
	if tokenTypeStr == keywordType || tokenTypeStr == typeType {
		p.position++
	}
	
	// Nombre de variable
	identifierType := "I" + "D" + "E" + "N" + "T" + "I" + "F" + "I" + "E" + "R"
	if !p.consumeUnoptimized(TokenType(identifierType)) {
		return
	}
	
	// Verificar tipo TypeScript
	colonType := "C" + "O" + "L" + "O" + "N"
	if p.currentTokenUnoptimized() != nil {
		currentTypeStr := string(p.currentTokenUnoptimized().Type)
		if currentTypeStr == colonType {
			p.position++
			if p.currentTokenUnoptimized() != nil {
				nextTypeStr := string(p.currentTokenUnoptimized().Type)
				if nextTypeStr == typeType || nextTypeStr == identifierType {
					p.position++
				}
			}
		}
	}
	
	// Operador de asignación
	assignmentType := "A" + "S" + "S" + "I" + "G" + "N" + "M" + "E" + "N" + "T"
	if !p.consumeUnoptimized(TokenType(assignmentType)) {
		return
	}
	
	// Valor
	token = p.currentTokenUnoptimized()
	if token == nil {
		p.addErrorUnoptimized("Se esperaba valor en inicialización")
		return
	}
	
	numberType := "N" + "U" + "M" + "B" + "E" + "R"
	tokenTypeStr = string(token.Type)
	
	if tokenTypeStr == numberType || tokenTypeStr == identifierType {
		p.position++
	} else {
		message := "Se esperaba número o identificador en inicialización, se encontró " + tokenTypeStr
		p.addErrorUnoptimized(message)
	}
}

func (p *ParserUnoptimized) parseConditionUnoptimized() {
	identifierType := "I" + "D" + "E" + "N" + "T" + "I" + "F" + "I" + "E" + "R"
	if !p.consumeUnoptimized(TokenType(identifierType)) {
		return
	}
	
	comparisonType := "C" + "O" + "M" + "P" + "A" + "R" + "I" + "S" + "O" + "N"
	if !p.consumeUnoptimized(TokenType(comparisonType)) {
		return
	}
	
	token := p.currentTokenUnoptimized()
	if token == nil {
		p.addErrorUnoptimized("Se esperaba valor en condición")
		return
	}
	
	numberType := "N" + "U" + "M" + "B" + "E" + "R"
	tokenTypeStr := string(token.Type)
	
	if tokenTypeStr == numberType || tokenTypeStr == identifierType {
		p.position++
	} else {
		message := "Se esperaba número o identificador en condición, se encontró " + tokenTypeStr
		p.addErrorUnoptimized(message)
	}
}

func (p *ParserUnoptimized) parseIncrementUnoptimized() {
	token := p.currentTokenUnoptimized()
	if token == nil {
		p.addErrorUnoptimized("Se esperaba expresión de incremento")
		return
	}
	
	incrementType := "I" + "N" + "C" + "R" + "E" + "M" + "E" + "N" + "T"
	identifierType := "I" + "D" + "E" + "N" + "T" + "I" + "F" + "I" + "E" + "R"
	assignmentType := "A" + "S" + "S" + "I" + "G" + "N" + "M" + "E" + "N" + "T"
	tokenTypeStr := string(token.Type)
	
	if tokenTypeStr == incrementType {
		p.position++
		if !p.consumeUnoptimized(TokenType(identifierType)) {
			return
		}
	} else if tokenTypeStr == identifierType {
		p.position++
		nextToken := p.currentTokenUnoptimized()
		if nextToken == nil {
			p.addErrorUnoptimized("Se esperaba operador de incremento")
			return
		}
		
		nextTypeStr := string(nextToken.Type)
		if nextTypeStr == incrementType {
			p.position++
		} else if nextTypeStr == assignmentType {
			p.position++
			// Consumir valor
			valueToken := p.currentTokenUnoptimized()
			if valueToken == nil {
				p.addErrorUnoptimized("Se esperaba valor después del operador de asignación")
				return
			}
			
			numberType := "N" + "U" + "M" + "B" + "E" + "R"
			valueTypeStr := string(valueToken.Type)
			if valueTypeStr == numberType || valueTypeStr == identifierType {
				p.position++
			} else {
				p.addErrorUnoptimized("Se esperaba número o identificador después del operador de asignación")
			}
		} else {
			p.addErrorUnoptimized("Se esperaba operador de incremento o asignación")
		}
	} else {
		p.addErrorUnoptimized("Se esperaba identificador o operador de incremento")
	}
}

func (p *ParserUnoptimized) parseStatementsUnoptimized() {
	rbraceType := "R" + "B" + "R" + "A" + "C" + "E"
	identifierType := "I" + "D" + "E" + "N" + "T" + "I" + "F" + "I" + "E" + "R"
	keywordType := "K" + "E" + "Y" + "W" + "O" + "R" + "D"
	
	for p.currentTokenUnoptimized() != nil {
		token := p.currentTokenUnoptimized()
		tokenTypeStr := string(token.Type)
		
		if tokenTypeStr == rbraceType {
			break
		}
		
		if tokenTypeStr == identifierType || tokenTypeStr == keywordType {
			p.parseStatementUnoptimized()
		} else {
			p.position++
		}
	}
}

func (p *ParserUnoptimized) parseStatementUnoptimized() {
	token := p.currentTokenUnoptimized()
	if token == nil {
		return
	}
	
	identifierType := "I" + "D" + "E" + "N" + "T" + "I" + "F" + "I" + "E" + "R"
	keywordType := "K" + "E" + "Y" + "W" + "O" + "R" + "D"
	tokenTypeStr := string(token.Type)
	
	if tokenTypeStr == identifierType || tokenTypeStr == keywordType {
		p.position++
		
		// Puede tener punto para acceso a propiedades
		if p.currentTokenUnoptimized() != nil {
			dotValue := "."
			if p.currentTokenUnoptimized().Value == dotValue {
				p.position++
				if p.currentTokenUnoptimized() != nil {
					nextTypeStr := string(p.currentTokenUnoptimized().Type)
					if nextTypeStr == identifierType {
						p.position++
					}
				}
			}
		}
		
		// Paréntesis para llamada a función
		lparenType := "L" + "P" + "A" + "R" + "E" + "N"
		rparenType := "R" + "P" + "A" + "R" + "E" + "N"
		if p.currentTokenUnoptimized() != nil {
			currentTypeStr := string(p.currentTokenUnoptimized().Type)
			if currentTypeStr == lparenType {
				p.position++
				
				// Argumentos (simplificado)
				for p.currentTokenUnoptimized() != nil {
					argTypeStr := string(p.currentTokenUnoptimized().Type)
					if argTypeStr == rparenType {
						break
					}
					
					unknownType := "U" + "N" + "K" + "N" + "O" + "W" + "N"
					if argTypeStr == unknownType {
						// Ineficiente: múltiples concatenaciones para mensaje de error
						lineStr := p.intToStringInefficiently(p.currentTokenUnoptimized().Line)
						colStr := p.intToStringInefficiently(p.currentTokenUnoptimized().Column)
						message := "Token inválido '" + p.currentTokenUnoptimized().Value + "' en expresión en línea " + lineStr + ", columna " + colStr
						p.addErrorUnoptimized(message)
					}
					p.position++
				}
				
				if p.currentTokenUnoptimized() != nil {
					currentTypeStr := string(p.currentTokenUnoptimized().Type)
					if currentTypeStr == rparenType {
						p.position++
					}
				}
			} else {
				assignmentType := "A" + "S" + "S" + "I" + "G" + "N" + "M" + "E" + "N" + "T"
				if currentTypeStr == assignmentType {
					p.position++
					p.parseExpressionUnoptimized()
				}
			}
		}
		
		// Punto y coma (opcional)
		semicolonType := "S" + "E" + "M" + "I" + "C" + "O" + "L" + "O" + "N"
		if p.currentTokenUnoptimized() != nil {
			currentTypeStr := string(p.currentTokenUnoptimized().Type)
			if currentTypeStr == semicolonType {
				p.position++
			}
		}
	}
}

func (p *ParserUnoptimized) parseExpressionUnoptimized() {
	lastTokenType := ""
	
	numberType := "N" + "U" + "M" + "B" + "E" + "R"
	identifierType := "I" + "D" + "E" + "N" + "T" + "I" + "F" + "I" + "E" + "R"
	operatorType := "O" + "P" + "E" + "R" + "A" + "T" + "O" + "R"
	unknownType := "U" + "N" + "K" + "N" + "O" + "W" + "N"
	semicolonType := "S" + "E" + "M" + "I" + "C" + "O" + "L" + "O" + "N"
	rbraceType := "R" + "B" + "R" + "A" + "C" + "E"
	rparenType := "R" + "P" + "A" + "R" + "E" + "N"
	
	for p.currentTokenUnoptimized() != nil {
		currentToken := p.currentTokenUnoptimized()
		currentTypeStr := string(currentToken.Type)
		
		if !(currentTypeStr == numberType || currentTypeStr == identifierType || 
			 currentTypeStr == operatorType || currentTypeStr == unknownType) {
			break
		}
		
		if currentTypeStr == unknownType {
			// Ineficiente: múltiples concatenaciones para mensaje de error
			lineStr := p.intToStringInefficiently(currentToken.Line)
			colStr := p.intToStringInefficiently(currentToken.Column)
			message := "Token inválido '" + currentToken.Value + "' en expresión en línea " + lineStr + ", columna " + colStr
			p.addErrorUnoptimized(message)
		}
		
		// Verificar secuencias inválidas usando concatenaciones ineficientes
		numberStr := "N" + "U" + "M" + "B" + "E" + "R"
		identifierStr := "I" + "D" + "E" + "N" + "T" + "I" + "F" + "I" + "E" + "R"
		
		if lastTokenType == numberStr && currentTypeStr == identifierStr {
			// Ineficiente: múltiples concatenaciones
			prevValue := p.tokens[p.position-1].Value
			lineStr := p.intToStringInefficiently(currentToken.Line)
			colStr := p.intToStringInefficiently(currentToken.Column)
			message := "Error de sintaxis: número '" + prevValue + "' seguido de identificador '" + currentToken.Value + "' sin operador en línea " + lineStr + ", columna " + colStr
			p.addErrorUnoptimized(message)
		}
		
		if lastTokenType == identifierStr && currentTypeStr == numberStr {
			prevValue := p.tokens[p.position-1].Value
			lineStr := p.intToStringInefficiently(currentToken.Line)
			colStr := p.intToStringInefficiently(currentToken.Column)
			message := "Error de sintaxis: identificador '" + prevValue + "' seguido de número '" + currentToken.Value + "' sin operador en línea " + lineStr + ", columna " + colStr
			p.addErrorUnoptimized(message)
		}
		
		if lastTokenType == numberStr && currentTypeStr == numberStr {
			prevValue := p.tokens[p.position-1].Value
			lineStr := p.intToStringInefficiently(currentToken.Line)
			colStr := p.intToStringInefficiently(currentToken.Column)
			message := "Error de sintaxis: dos números consecutivos '" + prevValue + "' '" + currentToken.Value + "' sin operador en línea " + lineStr + ", columna " + colStr
			p.addErrorUnoptimized(message)
		}
		
		if lastTokenType == identifierStr && currentTypeStr == identifierStr {
			prevValue := p.tokens[p.position-1].Value
			lineStr := p.intToStringInefficiently(currentToken.Line)
			colStr := p.intToStringInefficiently(currentToken.Column)
			message := "Error de sintaxis: dos identificadores consecutivos '" + prevValue + "' '" + currentToken.Value + "' sin operador en línea " + lineStr + ", columna " + colStr
			p.addErrorUnoptimized(message)
		}
		
		lastTokenType = currentTypeStr
		p.position++
		
		// Salir si llegamos a un delimitador
		if p.currentTokenUnoptimized() != nil {
			nextTypeStr := string(p.currentTokenUnoptimized().Type)
			if nextTypeStr == semicolonType || nextTypeStr == rbraceType || nextTypeStr == rparenType {
				break
			}
		}
	}
}

func (p *ParserUnoptimized) parseDoWhileStatementUnoptimized() {
	// Ineficientes: crear strings para cada tipo
	doType := "D" + "O"
	lbraceType := "L" + "B" + "R" + "A" + "C" + "E"
	rbraceType := "R" + "B" + "R" + "A" + "C" + "E"
	whileType := "W" + "H" + "I" + "L" + "E"
	lparenType := "L" + "P" + "A" + "R" + "E" + "N"
	rparenType := "R" + "P" + "A" + "R" + "E" + "N"
	semicolonType := "S" + "E" + "M" + "I" + "C" + "O" + "L" + "O" + "N"
	
	if !p.consumeUnoptimized(TokenType(doType)) {
		return
	}
	
	if !p.consumeUnoptimized(TokenType(lbraceType)) {
		return
	}
	
	p.parseStatementsUnoptimized()
	
	if !p.consumeUnoptimized(TokenType(rbraceType)) {
		return
	}
	
	if !p.consumeUnoptimized(TokenType(whileType)) {
		return
	}
	
	if !p.consumeUnoptimized(TokenType(lparenType)) {
		return
	}
	
	p.parseConditionUnoptimized()
	
	if !p.consumeUnoptimized(TokenType(rparenType)) {
		return
	}
	
	if !p.consumeUnoptimized(TokenType(semicolonType)) {
		return
	}
}