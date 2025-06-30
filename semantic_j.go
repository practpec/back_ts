package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type SemanticJ struct {
	tokens      []TokenJ
	variables   map[string]VariableInfoJ
	information []string
}

type VariableInfoJ struct {
	Name         string
	Type         string
	InitialValue string
	Line         int
	Column       int
}

func NewSemanticJ(tokens []TokenJ) *SemanticJ {
	return &SemanticJ{
		tokens:      tokens,
		variables:   make(map[string]VariableInfoJ),
		information: []string{},
	}
}

func (s *SemanticJ) addInfo(message string) {
	s.information = append(s.information, message)
}

func (s *SemanticJ) Analyze() []string {
	s.analyzeClassStructure()
	s.analyzeVariableDeclarations()
	s.analyzeIfStatements()
	s.analyzeMethodCalls()
	s.checkVariableUsage()
	s.validateDataTypes()
	s.detectMalformedTokens()
	return s.information
}

func (s *SemanticJ) analyzeClassStructure() {
	hasPublic := false
	hasClass := false
	hasMain := false
	className := ""

	for i, token := range s.tokens {
		switch token.Type {
		case PUBLIC_J:
			if !hasPublic {
				hasPublic = true
				s.addInfo("✓ Modificador 'public' encontrado - Clase accesible públicamente")
			}
		case CLASS_J:
			if hasPublic && !hasClass {
				hasClass = true
				// Buscar nombre de la clase
				if i+1 < len(s.tokens) && s.tokens[i+1].Type == IDENTIFIER_J {
					className = s.tokens[i+1].Value
					s.addInfo(fmt.Sprintf("✓ Clase '%s' declarada correctamente", className))
				}
			}
		case MAIN_J:
			if !hasMain {
				hasMain = true
				s.addInfo("✓ Método main encontrado - Punto de entrada del programa")
			}
		}
	}

	if !hasPublic {
		s.addInfo("❌ ERROR SEMÁNTICO: Falta modificador 'public' en la clase")
	}
	if !hasClass {
		s.addInfo("❌ ERROR SEMÁNTICO: Falta declaración de clase")
	}
	if !hasMain {
		s.addInfo("❌ ERROR SEMÁNTICO: Falta método main - El programa no tiene punto de entrada")
	}

	if hasPublic && hasClass && hasMain {
		s.addInfo("✓ Estructura de clase Java básica completa")
	}
}

func (s *SemanticJ) analyzeVariableDeclarations() {
	for i := 0; i < len(s.tokens); i++ {
		token := s.tokens[i]

		// Buscar declaraciones: tipo identificador = valor;
		if token.Type == TYPE_J {
			if i+1 < len(s.tokens) && s.tokens[i+1].Type == IDENTIFIER_J {
				varName := s.tokens[i+1].Value
				varType := token.Value

				var initialValue string
				nextPos := i + 2
				if nextPos < len(s.tokens) && s.tokens[nextPos].Type == ASSIGNMENT_J {
					nextPos++ // saltar '='
					if nextPos < len(s.tokens) {
						if s.tokens[nextPos].Type == STRING_J || s.tokens[nextPos].Type == NUMBER_J {
							initialValue = s.tokens[nextPos].Value
						} else if s.tokens[nextPos].Type == IDENTIFIER_J {
							initialValue = s.tokens[nextPos].Value
						}
					}
				}

				s.variables[varName] = VariableInfoJ{
					Name:         varName,
					Type:         varType,
					InitialValue: initialValue,
					Line:         token.Line,
					Column:       token.Column,
				}

				// Validar tipo vs valor
				s.validateTypeAssignment(varType, initialValue, token.Line)

				s.addInfo(fmt.Sprintf("Variable '%s' declarada como '%s' con valor '%s' en línea %d",
					varName, varType, initialValue, token.Line))

				// Saltar tokens procesados
				i = nextPos
			}
		}
	}
}

func (s *SemanticJ) validateTypeAssignment(varType, value string, line int) {
	if value == "" {
		return
	}

	switch varType {
	case "int":
		if _, err := strconv.Atoi(value); err != nil && !s.isVariableOfType(value, "int") {
			s.addInfo(fmt.Sprintf("❌ ERROR SEMÁNTICO: Valor '%s' no es compatible con tipo 'int' en línea %d", value, line))
		} else {
			s.addInfo(fmt.Sprintf("✓ Asignación válida: valor '%s' compatible con tipo 'int'", value))
		}
	case "String":
		if !strings.HasPrefix(value, "\"") || !strings.HasSuffix(value, "\"") {
			if !s.isVariableOfType(value, "String") {
				s.addInfo(fmt.Sprintf("❌ ERROR SEMÁNTICO: Valor '%s' no es compatible con tipo 'String' en línea %d", value, line))
			}
		} else {
			s.addInfo(fmt.Sprintf("✓ Asignación válida: String '%s' correctamente declarado", value))
		}
	case "boolean":
		if value != "true" && value != "false" && !s.isVariableOfType(value, "boolean") {
			s.addInfo(fmt.Sprintf("❌ ERROR SEMÁNTICO: Valor '%s' no es compatible con tipo 'boolean' en línea %d", value, line))
		}
	}
}

func (s *SemanticJ) isVariableOfType(varName, expectedType string) bool {
	if varInfo, exists := s.variables[varName]; exists {
		return varInfo.Type == expectedType
	}
	return false
}

func (s *SemanticJ) analyzeIfStatements() {
	for i := 0; i < len(s.tokens); i++ {
		if s.tokens[i].Type == IF_J {
			s.addInfo("Estructura 'if' detectada - Analizando condición")
			
			// Buscar condición entre paréntesis
			parenStart := i + 1
			if parenStart < len(s.tokens) && s.tokens[parenStart].Type == LPAREN_J {
				// Analizar condición
				condStart := parenStart + 1
				condEnd := s.findMatchingParen(parenStart)
				
				if condEnd > condStart {
					s.analyzeCondition(condStart, condEnd)
				}
			}
		}
	}
}

func (s *SemanticJ) findMatchingParen(start int) int {
	parenCount := 1
	for i := start + 1; i < len(s.tokens); i++ {
		if s.tokens[i].Type == LPAREN_J {
			parenCount++
		} else if s.tokens[i].Type == RPAREN_J {
			parenCount--
			if parenCount == 0 {
				return i
			}
		}
	}
	return -1
}

func (s *SemanticJ) analyzeCondition(start, end int) {
	hasComparison := false
	hasMethodCall := false
	var leftOperand, operator, rightOperand string

	for i := start; i < end; i++ {
		token := s.tokens[i]
		
		if token.Type == IDENTIFIER_J {
			// Verificar si la variable está declarada
			if _, exists := s.variables[token.Value]; !exists && !s.isReservedWord(token.Value) {
				s.addInfo(fmt.Sprintf("❌ ERROR SEMÁNTICO: Variable '%s' usada en condición pero no está declarada", token.Value))
			}
			
			if leftOperand == "" {
				leftOperand = token.Value
			}
			
			// Verificar si es llamada a método (.equals)
			if i+2 < end && s.tokens[i+1].Type == DOT_J && s.tokens[i+2].Value == "equals" {
				hasMethodCall = true
				s.analyzeEqualsMethod(i, end)
			}
		} else if token.Type == COMPARISON_J {
			hasComparison = true
			operator = token.Value
		} else if token.Type == NUMBER_J && rightOperand == "" {
			rightOperand = token.Value
		}
	}

	if hasComparison && leftOperand != "" && rightOperand != "" {
		s.validateComparison(leftOperand, operator, rightOperand)
	}

	if hasMethodCall {
		s.addInfo("✓ Llamada a método .equals() detectada para comparación de Strings")
	}
}

func (s *SemanticJ) validateComparison(left, operator, right string) {
	// Obtener tipos de los operandos
	leftType := s.getOperandType(left)
	rightType := s.getOperandType(right)

	s.addInfo(fmt.Sprintf("Comparación detectada: '%s %s %s'", left, operator, right))

	// Validar que los tipos sean compatibles
	if leftType == "int" && rightType == "int" {
		s.addInfo("✓ Comparación numérica válida: ambos operandos son de tipo int")
	} else if leftType == "String" && rightType == "String" {
		if operator == "==" {
			s.addInfo("⚠️ ADVERTENCIA: Comparación de Strings con '==' - Se recomienda usar .equals()")
		} else {
			s.addInfo("❌ ERROR SEMÁNTICO: Los Strings no se pueden comparar con operadores como >, <, etc.")
		}
	} else if leftType != rightType {
		s.addInfo(fmt.Sprintf("❌ ERROR SEMÁNTICO: Comparación entre tipos incompatibles: %s vs %s", leftType, rightType))
	} else {
		s.addInfo("✓ Comparación válida")
	}
}

func (s *SemanticJ) getOperandType(operand string) string {
	// Si es un número literal
	if _, err := strconv.Atoi(operand); err == nil {
		return "int"
	}
	
	// Si es una variable declarada
	if varInfo, exists := s.variables[operand]; exists {
		return varInfo.Type
	}
	
	// Si es un string literal
	if strings.HasPrefix(operand, "\"") && strings.HasSuffix(operand, "\"") {
		return "String"
	}
	
	return "unknown"
}

func (s *SemanticJ) analyzeEqualsMethod(start, end int) {
	// Analizar llamada a .equals("string")
	varName := s.tokens[start].Value
	
	// Verificar que la variable sea de tipo String
	if varInfo, exists := s.variables[varName]; exists {
		if varInfo.Type == "String" {
			s.addInfo(fmt.Sprintf("✓ Método .equals() usado correctamente en variable String '%s'", varName))
		} else {
			s.addInfo(fmt.Sprintf("❌ ERROR SEMÁNTICO: Método .equals() usado en variable '%s' de tipo '%s' - Solo válido para Strings", varName, varInfo.Type))
		}
	} else {
		s.addInfo(fmt.Sprintf("❌ ERROR SEMÁNTICO: Variable '%s' no declarada en llamada a .equals()", varName))
	}
	
	// Verificar argumento del equals
	for i := start; i < end; i++ {
		if s.tokens[i].Value == "equals" && i+2 < end {
			if s.tokens[i+1].Type == LPAREN_J && s.tokens[i+2].Type == STRING_J {
				argValue := s.tokens[i+2].Value
				s.addInfo(fmt.Sprintf("✓ Argumento de .equals() es válido: %s", argValue))
			}
		}
	}
}

func (s *SemanticJ) analyzeMethodCalls() {
	for i := 0; i < len(s.tokens); i++ {
		token := s.tokens[i]
		
		// Detectar cualquier intento de usar System (correcto o incorrecto)
		if token.Type == IDENTIFIER_J && (strings.Contains(strings.ToLower(token.Value), "system") || 
			token.Value == "System" || token.Value == "Syste" || token.Value == "sistem" ||
			strings.Contains(strings.ToLower(token.Value), "syste")) {
			
			s.analyzeSystemOutCall(i)
		}
		
		// Detectar tokens que contengan "out" o "println" fuera de contexto
		if token.Type == IDENTIFIER_J && (strings.Contains(token.Value, "out") || 
			strings.Contains(token.Value, "println")) {
			s.analyzeOrphanedTokens(i)
		}
	}
}

func (s *SemanticJ) analyzeSystemOutCall(startPos int) {
	if startPos >= len(s.tokens) {
		return
	}
	
	systemToken := s.tokens[startPos]
	
	// 1. Verificar que "System" esté correctamente escrito
	if systemToken.Value != "System" {
		s.addInfo(fmt.Sprintf("❌ ERROR LÉXICO: '%s' está mal escrito - debe ser exactamente 'System' en línea %d", 
			systemToken.Value, systemToken.Line))
		return
	}
	
	// 2. Verificar que haya exactamente un punto después de System
	pos := startPos + 1
	if pos >= len(s.tokens) || s.tokens[pos].Type != DOT_J {
		s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Falta punto (.) después de 'System' en línea %d", systemToken.Line))
		return
	}
	pos++
	
	// 3. Verificar que sea exactamente "out"
	if pos >= len(s.tokens) {
		s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Llamada incompleta - falta 'out' después de 'System.' en línea %d", systemToken.Line))
		return
	}
	
	outToken := s.tokens[pos]
	if outToken.Value != "out" {
		s.addInfo(fmt.Sprintf("❌ ERROR LÉXICO: '%s' está mal escrito - debe ser exactamente 'out' en línea %d", 
			outToken.Value, outToken.Line))
		return
	}
	pos++
	
	// 4. Verificar que haya exactamente un punto después de out
	if pos >= len(s.tokens) || s.tokens[pos].Type != DOT_J {
		s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Falta punto (.) después de 'out' en línea %d", outToken.Line))
		return
	}
	pos++
	
	// 5. Verificar que sea exactamente "println"
	if pos >= len(s.tokens) {
		s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Llamada incompleta - falta 'println' después de 'System.out.' en línea %d", systemToken.Line))
		return
	}
	
	printlnToken := s.tokens[pos]
	if printlnToken.Value != "println" {
		s.addInfo(fmt.Sprintf("❌ ERROR LÉXICO: '%s' está mal escrito - debe ser exactamente 'println' en línea %d", 
			printlnToken.Value, printlnToken.Line))
		return
	}
	pos++
	
	// 6. Verificar paréntesis de apertura
	if pos >= len(s.tokens) || s.tokens[pos].Type != LPAREN_J {
		s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Falta paréntesis de apertura '(' después de 'println' en línea %d", printlnToken.Line))
		return
	}
	pos++
	
	// 7. Analizar argumentos dentro de los paréntesis
	argumentsValid := s.analyzeSystemOutArguments(pos, systemToken.Line)
	
	if argumentsValid {
		s.addInfo("✓ System.out.println() es completamente válido")
	}
}

func (s *SemanticJ) analyzeSystemOutArguments(startPos int, systemLine int) bool {
	pos := startPos
	argumentCount := 0
	hasErrors := false
	hasStringConcatenation := false
	
	// Buscar hasta el paréntesis de cierre
	for pos < len(s.tokens) && s.tokens[pos].Type != RPAREN_J {
		token := s.tokens[pos]
		
		if token.Type == STRING_J {
			argumentCount++
			s.addInfo(fmt.Sprintf("✓ Argumento String válido: %s", token.Value))
			pos++
			
			// Verificar si hay concatenación después del string
			if pos < len(s.tokens) && s.tokens[pos].Type == OPERATOR_J && s.tokens[pos].Value == "+" {
				hasStringConcatenation = true
				s.addInfo("✓ Concatenación detectada - Operador '+' válido para unir strings")
				pos++ // consumir el +
				continue // continuar analizando el siguiente elemento
			}
			
			// Si hay otro token que NO sea ) o + después del string, es error
			if pos < len(s.tokens) && s.tokens[pos].Type != RPAREN_J && 
			   !(s.tokens[pos].Type == OPERATOR_J && s.tokens[pos].Value == "+") {
				nextToken := s.tokens[pos]
				if nextToken.Type == OPERATOR_J && nextToken.Value != "+" {
					s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Operador '%s' NO permitido dentro de println() - Solo '+' para concatenación en línea %d", 
						nextToken.Value, nextToken.Line))
					hasErrors = true
				} else if nextToken.Type != COMMA_J {
					s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Token '%s' NO permitido después de String en println() en línea %d", 
						nextToken.Value, nextToken.Line))
					hasErrors = true
				}
				pos++
			}
			
		} else if token.Type == IDENTIFIER_J {
			argumentCount++
			// Verificar que la variable esté declarada
			if _, exists := s.variables[token.Value]; exists {
				s.addInfo(fmt.Sprintf("✓ Variable '%s' válida como argumento", token.Value))
			} else {
				s.addInfo(fmt.Sprintf("❌ ERROR SEMÁNTICO: Variable '%s' NO declarada en println() en línea %d", 
					token.Value, token.Line))
				hasErrors = true
			}
			pos++
			
			// Verificar si hay concatenación después de la variable
			if pos < len(s.tokens) && s.tokens[pos].Type == OPERATOR_J && s.tokens[pos].Value == "+" {
				hasStringConcatenation = true
				s.addInfo("✓ Concatenación detectada - Operador '+' válido para unir elementos")
				pos++ // consumir el +
				continue // continuar analizando el siguiente elemento
			}
			
		} else if token.Type == NUMBER_J {
			argumentCount++
			s.addInfo(fmt.Sprintf("✓ Número válido como argumento: %s", token.Value))
			pos++
			
			// Verificar concatenación después del número
			if pos < len(s.tokens) && s.tokens[pos].Type == OPERATOR_J && s.tokens[pos].Value == "+" {
				hasStringConcatenation = true
				s.addInfo("✓ Concatenación detectada - Operador '+' válido para unir elementos")
				pos++ // consumir el +
				continue
			}
			
		} else if token.Type == OPERATOR_J && token.Value == "+" {
			// Si llegamos aquí, significa que hay un + sin elemento previo válido
			s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Operador '+' sin elemento previo válido en línea %d", 
				token.Line))
			hasErrors = true
			pos++
			
		} else if token.Type == OPERATOR_J && token.Value != "+" {
			s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Operador '%s' NO permitido dentro de println() - Solo '+' para concatenación en línea %d", 
				token.Value, token.Line))
			hasErrors = true
			pos++
			
		} else if token.Type == COMMA_J {
			s.addInfo(fmt.Sprintf("❌ ERROR SEMÁNTICO: System.out.println() SOLO acepta UN argumento - múltiples argumentos NO permitidos en línea %d", 
				token.Line))
			hasErrors = true
			pos++
			
		} else if token.Type == UNKNOWN_J {
			s.addInfo(fmt.Sprintf("❌ ERROR LÉXICO: Token inválido '%s' dentro de println() en línea %d", 
				token.Value, token.Line))
			hasErrors = true
			pos++
			
		} else {
			s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Token '%s' NO permitido dentro de println() en línea %d", 
				token.Value, token.Line))
			hasErrors = true
			pos++
		}
	}
	
	// Verificar paréntesis de cierre
	if pos >= len(s.tokens) || s.tokens[pos].Type != RPAREN_J {
		s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Falta paréntesis de cierre ')' en println() en línea %d", systemLine))
		hasErrors = true
	}
	
	// Verificar cantidad de argumentos considerando concatenación
	if argumentCount == 0 {
		s.addInfo("ℹ️ println() sin argumentos - imprimirá línea vacía")
	} else if hasStringConcatenation {
		s.addInfo("✓ Concatenación de strings válida - Java permite unir múltiples elementos con '+'")
	}
	
	return !hasErrors
}

func (s *SemanticJ) analyzeOrphanedTokens(pos int) {
	token := s.tokens[pos]
	
	// Detectar tokens como "outprintln" o similares
	if strings.Contains(token.Value, "outprint") || strings.Contains(token.Value, "println") {
		if !s.isPrecededBySystemDot(pos) {
			s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: '%s' debe ser parte de 'System.out.println()' en línea %d", 
				token.Value, token.Line))
		}
	}
	
	// Detectar "out" suelto
	if token.Value == "out" && !s.isPrecededBySystemDot(pos) {
		s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: 'out' debe ser parte de 'System.out.println()' en línea %d", 
			token.Line))
	}
}

func (s *SemanticJ) isPrecededBySystemDot(pos int) bool {
	if pos < 2 {
		return false
	}
	return s.tokens[pos-2].Value == "System" && s.tokens[pos-1].Type == DOT_J
}

func (s *SemanticJ) checkVariableUsage() {
	usedVars := make(map[string]bool)
	
	for _, token := range s.tokens {
		if token.Type == IDENTIFIER_J {
			if _, exists := s.variables[token.Value]; exists {
				usedVars[token.Value] = true
			}
		}
	}
	
	for varName, varInfo := range s.variables {
		if !usedVars[varName] {
			s.addInfo(fmt.Sprintf("⚠️ Variable '%s' declarada pero no utilizada", varName))
		} else {
			s.addInfo(fmt.Sprintf("✓ Variable '%s' declarada y utilizada correctamente", varName))
			_ = varInfo // evitar warning
		}
	}
}

func (s *SemanticJ) validateDataTypes() {
	intCount := 0
	stringCount := 0
	
	for _, varInfo := range s.variables {
		switch varInfo.Type {
		case "int":
			intCount++
		case "String":
			stringCount++
		}
	}
	
	s.addInfo(fmt.Sprintf("📊 Resumen de tipos: %d variables int, %d variables String", intCount, stringCount))
	
	if intCount > 0 && stringCount > 0 {
		s.addInfo("✓ Uso diverso de tipos de datos - Buena práctica")
	}
}

func (s *SemanticJ) detectMalformedTokens() {
	for i, token := range s.tokens {
		if token.Type == UNKNOWN_J {
			// Detectar errores específicos de println
			if strings.Contains(token.Value, "print") {
				if token.Value == "prntln" {
					s.addInfo(fmt.Sprintf("❌ ERROR LÉXICO: '%s' está mal escrito - debe ser 'println' (falta 'i') en línea %d", 
						token.Value, token.Line))
				} else if strings.Contains(token.Value, "printl") {
					s.addInfo(fmt.Sprintf("❌ ERROR LÉXICO: '%s' está mal escrito - debe ser 'println' en línea %d", 
						token.Value, token.Line))
				} else {
					s.addInfo(fmt.Sprintf("❌ ERROR LÉXICO: '%s' está mal escrito - debe ser exactamente 'println' en línea %d", 
						token.Value, token.Line))
				}
			} else if len(token.Value) > 0 && unicode.IsDigit(rune(token.Value[0])) {
				s.addInfo(fmt.Sprintf("❌ ERROR LÉXICO: Número mal formado '%s' en línea %d", 
					token.Value, token.Line))
			} else {
				s.addInfo(fmt.Sprintf("❌ ERROR LÉXICO: Token inválido '%s' en línea %d", 
					token.Value, token.Line))
			}
		}
		
		// Detectar identificadores que podrían ser errores de println
		if token.Type == IDENTIFIER_J && strings.Contains(token.Value, "print") {
			if token.Value != "print" {  // "print" solo es válido
				s.addInfo(fmt.Sprintf("❌ ERROR LÉXICO: '%s' parece un error de escritura - debe ser exactamente 'println' en línea %d", 
					token.Value, token.Line))
			}
		}
		
		// NO marcar error para concatenación válida de strings
		if token.Type == STRING_J && i+1 < len(s.tokens) {
			nextToken := s.tokens[i+1]
			// Solo marcar error si NO es el operador + para concatenación
			if nextToken.Type == OPERATOR_J && nextToken.Value != "+" {
				s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Operador '%s' después de String NO permitido en línea %d - Use '+' para concatenación", 
					nextToken.Value, nextToken.Line))
			}
		}
	}
}

func (s *SemanticJ) isReservedWord(word string) bool {
	reservedWords := map[string]bool{
		"System":   true,
		"out":      true,
		"println":  true,
		"print":    true,
		"equals":   true,
		"length":   true,
		"args":     true,
		"public":   true,
		"class":    true,
		"static":   true,
		"void":     true,
		"main":     true,
		"if":       true,
		"else":     true,
		"true":     true,
		"false":    true,
		"null":     true,
	}
	return reservedWords[strings.ToLower(word)]
}