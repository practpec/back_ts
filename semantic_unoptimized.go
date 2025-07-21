package main

import (
	"strconv"
	"strings"
	"unicode"
)

// Versión NO optimizada del analizador semántico que usa concatenación de strings
// y múltiples operaciones ineficientes para demostrar el impacto en rendimiento

type SemanticUnoptimized struct {
	tokens      []Token
	variables   map[string]VariableInfo
	information []string
}

func NewSemanticUnoptimized(tokens []Token) *SemanticUnoptimized {
	return &SemanticUnoptimized{
		tokens:      tokens,
		variables:   make(map[string]VariableInfo),
		information: []string{},
	}
}

func (s *SemanticUnoptimized) addInfoUnoptimized(message string) {
	s.information = append(s.information, message)
}

func (s *SemanticUnoptimized) AnalyzeUnoptimized() []string {
	s.analyzeVariableDeclarationsUnoptimized()
	s.analyzeForLoopUnoptimized()
	s.checkVariableUsageUnoptimized()
	s.analyzeInfiniteLoopUnoptimized()
	s.detectUndeclaredVariablesUnoptimized()
	s.detectMalformedNumbersUnoptimized()
	s.detectInvalidExpressionsUnoptimized()
	s.analyzeDoWhileLoopUnoptimized()
	return s.information
}

func (s *SemanticUnoptimized) detectMalformedNumbersUnoptimized() {
	for _, token := range s.tokens {
		// Ineficiente: crear string para comparación de tipo
		unknownType := "U" + "N" + "K" + "N" + "O" + "W" + "N"
		tokenTypeStr := string(token.Type)
		
		if tokenTypeStr == unknownType {
			// Verificar si parece un número mal formado usando concatenación ineficiente
			if len(token.Value) > 0 && unicode.IsDigit(rune(token.Value[0])) {
				// Ineficiente: múltiples concatenaciones para mensaje de error
				errorMsg := "❌ ERROR LÉXICO: Número mal formado '"
				errorMsg = errorMsg + token.Value
				errorMsg = errorMsg + "' en línea "
				errorMsg = errorMsg + s.intToStringInefficiently(token.Line)
				errorMsg = errorMsg + ", columna "
				errorMsg = errorMsg + s.intToStringInefficiently(token.Column)
				s.addInfoUnoptimized(errorMsg)
			}
		}
	}
}

func (s *SemanticUnoptimized) detectInvalidExpressionsUnoptimized() {
	for i := 0; i < len(s.tokens)-1; i++ {
		currentToken := s.tokens[i]
		nextToken := s.tokens[i+1]
		
		// Ineficientes: crear strings para tipos
		numberType := "N" + "U" + "M" + "B" + "E" + "R"
		identifierType := "I" + "D" + "E" + "N" + "T" + "I" + "F" + "I" + "E" + "R"
		
		currentTypeStr := string(currentToken.Type)
		nextTypeStr := string(nextToken.Type)
		
		// Verificar patrones inválidos usando concatenaciones ineficientes
		if currentTypeStr == numberType && nextTypeStr == identifierType {
			if currentToken.Line == nextToken.Line {
				errorMsg := "❌ ERROR SINTÁCTICO: Número '"
				errorMsg = errorMsg + currentToken.Value
				errorMsg = errorMsg + "' seguido de identificador '"
				errorMsg = errorMsg + nextToken.Value
				errorMsg = errorMsg + "' sin operador en línea "
				errorMsg = errorMsg + s.intToStringInefficiently(currentToken.Line)
				s.addInfoUnoptimized(errorMsg)
			}
		}
		
		if currentTypeStr == identifierType && nextTypeStr == numberType {
			if currentToken.Line == nextToken.Line && !s.isReservedWordUnoptimized(currentToken.Value) {
				errorMsg := "❌ ERROR SINTÁCTICO: Identificador '"
				errorMsg = errorMsg + currentToken.Value
				errorMsg = errorMsg + "' seguido de número '"
				errorMsg = errorMsg + nextToken.Value
				errorMsg = errorMsg + "' sin operador en línea "
				errorMsg = errorMsg + s.intToStringInefficiently(currentToken.Line)
				s.addInfoUnoptimized(errorMsg)
			}
		}
		
		if currentTypeStr == numberType && nextTypeStr == numberType {
			if currentToken.Line == nextToken.Line {
				errorMsg := "❌ ERROR SINTÁCTICO: Dos números consecutivos '"
				errorMsg = errorMsg + currentToken.Value
				errorMsg = errorMsg + "' '"
				errorMsg = errorMsg + nextToken.Value
				errorMsg = errorMsg + "' sin operador en línea "
				errorMsg = errorMsg + s.intToStringInefficiently(currentToken.Line)
				s.addInfoUnoptimized(errorMsg)
			}
		}
		
		if currentTypeStr == identifierType && nextTypeStr == identifierType {
			if currentToken.Line == nextToken.Line && 
				!s.isReservedWordUnoptimized(currentToken.Value) && 
				!s.isReservedWordUnoptimized(nextToken.Value) {
				errorMsg := "❌ ERROR SINTÁCTICO: Dos identificadores consecutivos '"
				errorMsg = errorMsg + currentToken.Value
				errorMsg = errorMsg + "' '"
				errorMsg = errorMsg + nextToken.Value
				errorMsg = errorMsg + "' sin operador en línea "
				errorMsg = errorMsg + s.intToStringInefficiently(currentToken.Line)
				s.addInfoUnoptimized(errorMsg)
			}
		}
	}
}

func (s *SemanticUnoptimized) analyzeDoWhileLoopUnoptimized() {
	doFound := false
	whileFound := false
	var conditionVar string
	
	// Ineficientes: crear strings para tipos
	doType := "D" + "O"
	whileType := "W" + "H" + "I" + "L" + "E"
	lparenType := "L" + "P" + "A" + "R" + "E" + "N"
	rparenType := "R" + "P" + "A" + "R" + "E" + "N"
	identifierType := "I" + "D" + "E" + "N" + "T" + "I" + "F" + "I" + "E" + "R"
	
	for i := 0; i < len(s.tokens); i++ {
		token := s.tokens[i]
		tokenTypeStr := string(token.Type)
		
		if tokenTypeStr == doType {
			doFound = true
			s.addInfoUnoptimized("Bucle 'do-while' detectado - Analizando estructura")
		}
		
		if tokenTypeStr == whileType && doFound {
			whileFound = true
			s.addInfoUnoptimized("Cláusula 'while' encontrada en bucle do-while")
			
			// Buscar la condición después de while
			if i+2 < len(s.tokens) {
				nextTokenTypeStr := string(s.tokens[i+1].Type)
				if nextTokenTypeStr == lparenType {
					condPos := i + 2
					for condPos < len(s.tokens) {
						condTokenTypeStr := string(s.tokens[condPos].Type)
						if condTokenTypeStr == rparenType {
							break
						}
						if condTokenTypeStr == identifierType {
							conditionVar = s.tokens[condPos].Value
							break
						}
						condPos++
					}
					
					if conditionVar != "" {
						msg := "Variable en condición do-while: '"
						msg = msg + conditionVar
						msg = msg + "'"
						s.addInfoUnoptimized(msg)
						
						// Verificar si la variable está declarada
						if _, exists := s.variables[conditionVar]; !exists {
							errorMsg := "❌ ERROR SEMÁNTICO: Variable '"
							errorMsg = errorMsg + conditionVar
							errorMsg = errorMsg + "' en condición do-while no está declarada"
							s.addInfoUnoptimized(errorMsg)
						} else {
							successMsg := "✓ Variable '"
							successMsg = successMsg + conditionVar
							successMsg = successMsg + "' en condición do-while está correctamente declarada"
							s.addInfoUnoptimized(successMsg)
						}
					}
				}
			}
		}
	}
	
	if doFound && whileFound {
		s.addInfoUnoptimized("✓ Estructura do-while completa detectada")
	} else if doFound && !whileFound {
		s.addInfoUnoptimized("❌ ERROR SEMÁNTICO: Bucle 'do' sin cláusula 'while' correspondiente")
	}
}

func (s *SemanticUnoptimized) detectUndeclaredVariablesUnoptimized() {
	// Encontrar todas las variables usadas
	usedVariables := make(map[string][]int)
	
	identifierType := "I" + "D" + "E" + "N" + "T" + "I" + "F" + "I" + "E" + "R"
	
	for i, token := range s.tokens {
		tokenTypeStr := string(token.Type)
		if tokenTypeStr == identifierType {
			if !s.isReservedWordUnoptimized(token.Value) {
				usedVariables[token.Value] = append(usedVariables[token.Value], i)
			}
		}
	}
	
	// Verificar que todas las variables usadas estén declaradas
	for varName, positions := range usedVariables {
		if _, declared := s.variables[varName]; !declared {
			errorMsg := "❌ ERROR SEMÁNTICO: Variable '"
			errorMsg = errorMsg + varName
			errorMsg = errorMsg + "' usada sin declarar (línea "
			errorMsg = errorMsg + s.intToStringInefficiently(s.tokens[positions[0]].Line)
			errorMsg = errorMsg + ")"
			s.addInfoUnoptimized(errorMsg)
		}
	}
}

func (s *SemanticUnoptimized) isReservedWordUnoptimized(word string) bool {
	// Ineficiente: crear strings y hacer múltiples comparaciones
	wordLower := strings.ToLower(word)
	
	console := "c" + "o" + "n" + "s" + "o" + "l" + "e"
	log := "l" + "o" + "g"
	system := "s" + "y" + "s" + "t" + "e" + "m"
	out := "o" + "u" + "t"
	println := "p" + "r" + "i" + "n" + "t" + "l" + "n"
	print := "p" + "r" + "i" + "n" + "t"
	length := "l" + "e" + "n" + "g" + "t" + "h"
	
	return wordLower == console || wordLower == log || wordLower == system || 
		   wordLower == out || wordLower == println || wordLower == print || 
		   wordLower == length
}

func (s *SemanticUnoptimized) analyzeVariableDeclarationsUnoptimized() {
	// Ineficientes: crear strings para tipos
	keywordType := "K" + "E" + "Y" + "W" + "O" + "R" + "D"
	typeType := "T" + "Y" + "P" + "E"
	identifierType := "I" + "D" + "E" + "N" + "T" + "I" + "F" + "I" + "E" + "R"
	colonType := "C" + "O" + "L" + "O" + "N"
	assignmentType := "A" + "S" + "S" + "I" + "G" + "N" + "M" + "E" + "N" + "T"
	
	for i := 0; i < len(s.tokens); i++ {
		token := s.tokens[i]
		tokenTypeStr := string(token.Type)
		
		// Buscar declaraciones de variables
		if tokenTypeStr == keywordType || tokenTypeStr == typeType {
			if i+1 < len(s.tokens) {
				nextTokenTypeStr := string(s.tokens[i+1].Type)
				if nextTokenTypeStr == identifierType {
					varName := s.tokens[i+1].Value
					varType := s.inferTypeUnoptimized(token.Value)
					
					// Verificar si hay declaración de tipo TypeScript
					nextPos := i + 2
					if nextPos < len(s.tokens) {
						nextPosTokenTypeStr := string(s.tokens[nextPos].Type)
						if nextPosTokenTypeStr == colonType {
							nextPos++ // saltar ':'
							if nextPos < len(s.tokens) {
								typeTokenTypeStr := string(s.tokens[nextPos].Type)
								if typeTokenTypeStr == typeType || typeTokenTypeStr == identifierType {
									varType = s.tokens[nextPos].Value
									nextPos++
								}
							}
						}
					}
					
					var initialValue string
					if nextPos < len(s.tokens) {
						assignTokenTypeStr := string(s.tokens[nextPos].Type)
						if assignTokenTypeStr == assignmentType {
							nextPos++ // saltar '='
							if nextPos < len(s.tokens) {
								initialValue = s.tokens[nextPos].Value
							}
						}
					}
					
					s.variables[varName] = VariableInfo{
						Name:         varName,
						Type:         varType,
						InitialValue: initialValue,
						Line:         token.Line,
						Column:       token.Column,
					}
					
					// Ineficiente: múltiples concatenaciones para mensaje
					msg := "Variable '"
					msg = msg + varName
					msg = msg + "' declarada como tipo '"
					msg = msg + varType
					msg = msg + "' con valor inicial '"
					msg = msg + initialValue
					msg = msg + "' en línea "
					msg = msg + s.intToStringInefficiently(token.Line)
					s.addInfoUnoptimized(msg)
					
					i = nextPos
				}
			}
		}
	}
}

func (s *SemanticUnoptimized) analyzeForLoopUnoptimized() {
	forFound := false
	var loopVar string
	var conditionVar string
	var incrementVar string
	var startValue, endValue int
	var hasCondition bool
	
	// Ineficientes: crear strings para tipos
	forType := "F" + "O" + "R"
	identifierType := "I" + "D" + "E" + "N" + "T" + "I" + "F" + "I" + "E" + "R"
	keywordType := "K" + "E" + "Y" + "W" + "O" + "R" + "D"
	typeType := "T" + "Y" + "P" + "E"
	assignmentType := "A" + "S" + "S" + "I" + "G" + "N" + "M" + "E" + "N" + "T"
	comparisonType := "C" + "O" + "M" + "P" + "A" + "R" + "I" + "S" + "O" + "N"
	numberType := "N" + "U" + "M" + "B" + "E" + "R"
	incrementType := "I" + "N" + "C" + "R" + "E" + "M" + "E" + "N" + "T"
	
	for i := 0; i < len(s.tokens); i++ {
		token := s.tokens[i]
		tokenTypeStr := string(token.Type)
		
		if tokenTypeStr == forType {
			forFound = true
			s.addInfoUnoptimized("Bucle 'for' detectado - Analizando estructura")
		}
		
		// Analizar inicialización del bucle
		if forFound && tokenTypeStr == identifierType && i > 0 {
			prevTokenTypeStr := string(s.tokens[i-1].Type)
			if prevTokenTypeStr == keywordType || prevTokenTypeStr == typeType {
				loopVar = token.Value
				
				// Buscar valor inicial
				if i+2 < len(s.tokens) {
					assignTokenTypeStr := string(s.tokens[i+1].Type)
					if assignTokenTypeStr == assignmentType {
						if val, err := strconv.Atoi(s.tokens[i+2].Value); err == nil {
							startValue = val
							msg := "Variable de control '"
							msg = msg + loopVar
							msg = msg + "' inicializada con valor "
							msg = msg + s.intToStringInefficiently(startValue)
							s.addInfoUnoptimized(msg)
						}
					}
				}
			}
		}
		
		// Analizar condición del bucle
		if forFound && tokenTypeStr == comparisonType && i > 0 && i+1 < len(s.tokens) {
			hasCondition = true
			operator := token.Value
			
			prevTokenTypeStr := string(s.tokens[i-1].Type)
			nextTokenTypeStr := string(s.tokens[i+1].Type)
			
			if prevTokenTypeStr == identifierType && nextTokenTypeStr == numberType {
				conditionVar = s.tokens[i-1].Value
				if val, err := strconv.Atoi(s.tokens[i+1].Value); err == nil {
					endValue = val
					msg := "Condición: '"
					msg = msg + conditionVar
					msg = msg + " "
					msg = msg + operator
					msg = msg + " "
					msg = msg + s.intToStringInefficiently(endValue)
					msg = msg + "' - Variable de control se compara con "
					msg = msg + s.intToStringInefficiently(endValue)
					s.addInfoUnoptimized(msg)
					
					// Verificar coherencia de la condición
					s.checkLoopConditionCoherenceUnoptimized(operator, startValue, endValue)
				}
			}
		}
		
		// Analizar incremento
		if forFound && tokenTypeStr == incrementType {
			if i > 0 {
				prevTokenTypeStr := string(s.tokens[i-1].Type)
				if prevTokenTypeStr == identifierType {
					incrementVar = s.tokens[i-1].Value
				}
			}
			if i+1 < len(s.tokens) {
				nextTokenTypeStr := string(s.tokens[i+1].Type)
				if nextTokenTypeStr == identifierType {
					incrementVar = s.tokens[i+1].Value
				}
			}
			
			msg := "Incremento detectado para variable '"
			msg = msg + incrementVar
			msg = msg + "' ("
			msg = msg + token.Value
			msg = msg + ")"
			s.addInfoUnoptimized(msg)
		}
	}
	
	// Verificar consistencia de variables en el bucle
	if forFound && loopVar != "" {
		s.checkLoopVariableConsistencyUnoptimized(loopVar, conditionVar, incrementVar)
	}
	
	// Análisis completo del bucle
	if forFound && hasCondition {
		iterations := s.calculateIterationsUnoptimized(startValue, endValue)
		if iterations > 0 {
			msg := "El bucle ejecutará aproximadamente "
			msg = msg + s.intToStringInefficiently(iterations)
			msg = msg + " iteraciones"
			s.addInfoUnoptimized(msg)
		}
	}
}

func (s *SemanticUnoptimized) checkLoopConditionCoherenceUnoptimized(operator string, start, end int) {
	// Ineficientes: crear strings para operadores
	lessEqual := "<" + "="
	less := "<"
	greaterEqual := ">" + "="
	greater := ">"
	
	if operator == lessEqual || operator == less {
		if start > end {
			s.addInfoUnoptimized("⚠️ ADVERTENCIA: La condición del bucle podría nunca ser verdadera (valor inicial mayor que final)")
		}
	} else if operator == greaterEqual || operator == greater {
		if start < end {
			s.addInfoUnoptimized("⚠️ ADVERTENCIA: La condición del bucle podría nunca ser verdadera (valor inicial menor que final)")
		}
	}
}

func (s *SemanticUnoptimized) calculateIterationsUnoptimized(start, end int) int {
	if start <= end {
		return end - start + 1
	}
	return 0
}

func (s *SemanticUnoptimized) checkVariableUsageUnoptimized() {
	usedVars := make(map[string]bool)
	
	identifierType := "I" + "D" + "E" + "N" + "T" + "I" + "F" + "I" + "E" + "R"
	
	for _, token := range s.tokens {
		tokenTypeStr := string(token.Type)
		if tokenTypeStr == identifierType {
			if _, exists := s.variables[token.Value]; exists {
				usedVars[token.Value] = true
			}
		}
	}
	
	for varName := range s.variables {
		if !usedVars[varName] {
			msg := "⚠️ Variable '"
			msg = msg + varName
			msg = msg + "' declarada pero no utilizada"
			s.addInfoUnoptimized(msg)
		} else {
			msg := "✓ Variable '"
			msg = msg + varName
			msg = msg + "' declarada y utilizada correctamente"
			s.addInfoUnoptimized(msg)
		}
	}
}

func (s *SemanticUnoptimized) analyzeInfiniteLoopUnoptimized() {
	hasIncrement := false
	hasValidCondition := false
	
	incrementType := "I" + "N" + "C" + "R" + "E" + "M" + "E" + "N" + "T"
	comparisonType := "C" + "O" + "M" + "P" + "A" + "R" + "I" + "S" + "O" + "N"
	
	for _, token := range s.tokens {
		tokenTypeStr := string(token.Type)
		if tokenTypeStr == incrementType {
			hasIncrement = true
		}
		if tokenTypeStr == comparisonType {
			hasValidCondition = true
		}
	}
	
	if !hasIncrement && hasValidCondition {
		s.addInfoUnoptimized("⚠️ POSIBLE BUCLE INFINITO: No se detectó incremento en la variable de control")
	} else if hasIncrement && hasValidCondition {
		s.addInfoUnoptimized("✓ Estructura de bucle válida: tiene condición e incremento")
	}
}

func (s *SemanticUnoptimized) inferTypeUnoptimized(declaration string) string {
	// Ineficiente: crear strings para cada comparación
	declarationLower := strings.ToLower(declaration)
	
	intType := "i" + "n" + "t"
	stringType := "s" + "t" + "r" + "i" + "n" + "g"
	letType := "l" + "e" + "t"
	constType := "c" + "o" + "n" + "s" + "t"
	varType := "v" + "a" + "r"
	
	if declarationLower == intType {
		return "number"
	} else if declarationLower == stringType {
		return "string"
	} else if declarationLower == letType {
		return "variable"
	} else if declarationLower == constType {
		return "constant"
	} else if declarationLower == varType {
		return "variable"
	}
	return "unknown"
}

func (s *SemanticUnoptimized) checkLoopVariableConsistencyUnoptimized(loopVar, conditionVar, incrementVar string) {
	// Verificar que la variable de condición sea la misma que la declarada
	if conditionVar != "" && conditionVar != loopVar {
		errorMsg := "❌ ERROR SEMÁNTICO: Variable en condición '"
		errorMsg = errorMsg + conditionVar
		errorMsg = errorMsg + "' no coincide con variable de control '"
		errorMsg = errorMsg + loopVar
		errorMsg = errorMsg + "'"
		s.addInfoUnoptimized(errorMsg)
	} else if conditionVar == loopVar {
		successMsg := "✓ Variable de condición '"
		successMsg = successMsg + conditionVar
		successMsg = successMsg + "' coincide correctamente con variable de control"
		s.addInfoUnoptimized(successMsg)
	}
	
	// Verificar que la variable de incremento sea la misma que la declarada
	if incrementVar != "" && incrementVar != loopVar {
		errorMsg := "❌ ERROR SEMÁNTICO: Variable en incremento '"
		errorMsg = errorMsg + incrementVar
		errorMsg = errorMsg + "' no coincide con variable de control '"
		errorMsg = errorMsg + loopVar
		errorMsg = errorMsg + "'"
		s.addInfoUnoptimized(errorMsg)
	} else if incrementVar == loopVar {
		successMsg := "✓ Variable de incremento '"
		successMsg = successMsg + incrementVar
		successMsg = successMsg + "' coincide correctamente con variable de control"
		s.addInfoUnoptimized(successMsg)
	}
	
	// Verificar que la variable de control esté siendo utilizada en todas las partes
	if conditionVar == "" {
		s.addInfoUnoptimized("⚠️ ADVERTENCIA: No se detectó variable en la condición del bucle")
	}
	
	if incrementVar == "" {
		s.addInfoUnoptimized("⚠️ ADVERTENCIA: No se detectó variable en el incremento del bucle")
	}
}

// Función ineficiente para convertir int a string usando strconv pero con concatenaciones
func (s *SemanticUnoptimized) intToStringInefficiently(num int) string {
	// Usar strconv pero de forma ineficiente con concatenaciones
	baseStr := strconv.Itoa(num)
	result := ""
	for _, char := range baseStr {
		result = result + string(char)
	}
	return result
}