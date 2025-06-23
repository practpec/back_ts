package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Semantic struct {
	tokens      []Token
	variables   map[string]VariableInfo
	information []string
}

type VariableInfo struct {
	Name         string
	Type         string
	InitialValue string
	Line         int
	Column       int
}

func NewSemantic(tokens []Token) *Semantic {
	return &Semantic{
		tokens:      tokens,
		variables:   make(map[string]VariableInfo),
		information: []string{},
	}
}

func (s *Semantic) addInfo(message string) {
	s.information = append(s.information, message)
}
func (s *Semantic) detectMalformedNumbers() {
	for _, token := range s.tokens {
		if token.Type == UNKNOWN {
			// Verificar si parece un número mal formado
			if len(token.Value) > 0 && unicode.IsDigit(rune(token.Value[0])) {
				s.addInfo(fmt.Sprintf("❌ ERROR LÉXICO: Número mal formado '%s' en línea %d, columna %d", 
					token.Value, token.Line, token.Column))
			}
		}
	}
}

func (s *Semantic) detectInvalidExpressions() {
	for i := 0; i < len(s.tokens)-1; i++ {
		currentToken := s.tokens[i]
		nextToken := s.tokens[i+1]
		
		// Verificar patrones inválidos de tokens consecutivos
		if currentToken.Type == NUMBER && nextToken.Type == IDENTIFIER {
			// Verificar que no haya operador entre ellos
			if currentToken.Line == nextToken.Line {
				s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Número '%s' seguido de identificador '%s' sin operador en línea %d", 
					currentToken.Value, nextToken.Value, currentToken.Line))
			}
		}
		
		if currentToken.Type == IDENTIFIER && nextToken.Type == NUMBER {
			if currentToken.Line == nextToken.Line && !s.isReservedWord(currentToken.Value) {
				s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Identificador '%s' seguido de número '%s' sin operador en línea %d", 
					currentToken.Value, nextToken.Value, currentToken.Line))
			}
		}
		
		if currentToken.Type == NUMBER && nextToken.Type == NUMBER {
			if currentToken.Line == nextToken.Line {
				s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Dos números consecutivos '%s' '%s' sin operador en línea %d", 
					currentToken.Value, nextToken.Value, currentToken.Line))
			}
		}
		
		if currentToken.Type == IDENTIFIER && nextToken.Type == IDENTIFIER {
			if currentToken.Line == nextToken.Line && 
				!s.isReservedWord(currentToken.Value) && 
				!s.isReservedWord(nextToken.Value) {
				s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Dos identificadores consecutivos '%s' '%s' sin operador en línea %d", 
					currentToken.Value, nextToken.Value, currentToken.Line))
			}
		}
	}
}
func (s *Semantic) analyzeDoWhileLoop() {
	doFound := false
	whileFound := false
	var conditionVar string
	
	for i := 0; i < len(s.tokens); i++ {
		token := s.tokens[i]
		
		if token.Type == DO {
			doFound = true
			s.addInfo("Bucle 'do-while' detectado - Analizando estructura")
		}
		
		if token.Type == WHILE && doFound {
			whileFound = true
			s.addInfo("Cláusula 'while' encontrada en bucle do-while")
			
			// Buscar la condición después de while
			if i+2 < len(s.tokens) && s.tokens[i+1].Type == LPAREN {
				condPos := i + 2
				for condPos < len(s.tokens) && s.tokens[condPos].Type != RPAREN {
					if s.tokens[condPos].Type == IDENTIFIER {
						conditionVar = s.tokens[condPos].Value
						break
					}
					condPos++
				}
				
				if conditionVar != "" {
					s.addInfo(fmt.Sprintf("Variable en condición do-while: '%s'", conditionVar))
					
					// Verificar si la variable está declarada
					if _, exists := s.variables[conditionVar]; !exists {
						s.addInfo(fmt.Sprintf("❌ ERROR SEMÁNTICO: Variable '%s' en condición do-while no está declarada", conditionVar))
					} else {
						s.addInfo(fmt.Sprintf("✓ Variable '%s' en condición do-while está correctamente declarada", conditionVar))
					}
				}
			}
		}
	}
	
	if doFound && whileFound {
		s.addInfo("✓ Estructura do-while completa detectada")
	} else if doFound && !whileFound {
		s.addInfo("❌ ERROR SEMÁNTICO: Bucle 'do' sin cláusula 'while' correspondiente")
	}
}
func (s *Semantic) detectUndeclaredVariables() {
	// Encontrar todas las variables usadas
	usedVariables := make(map[string][]int) // variable -> posiciones donde se usa
	
	for i, token := range s.tokens {
		if token.Type == IDENTIFIER {
			// Excluir palabras reservadas como 'console', 'log', 'out', 'println', etc.
			if !s.isReservedWord(token.Value) {
				usedVariables[token.Value] = append(usedVariables[token.Value], i)
			}
		}
	}
	
	// Verificar que todas las variables usadas estén declaradas
	for varName, positions := range usedVariables {
		if _, declared := s.variables[varName]; !declared {
			s.addInfo(fmt.Sprintf("❌ ERROR SEMÁNTICO: Variable '%s' usada sin declarar (línea %d)", 
				varName, s.tokens[positions[0]].Line))
		}
	}
}

func (s *Semantic) isReservedWord(word string) bool {
	reservedWords := map[string]bool{
		"console": true,
		"log":     true,
		"system":  true,
		"out":     true,
		"println": true,
		"print":   true,
		"length":  true,
	}
	return reservedWords[strings.ToLower(word)]
}

func (s *Semantic) Analyze() []string {
	s.analyzeVariableDeclarations()
	s.analyzeForLoop()
	s.checkVariableUsage()
	s.analyzeInfiniteLoop()
	s.detectUndeclaredVariables() // Nueva función
	return s.information
}

func (s *Semantic) analyzeVariableDeclarations() {
	for i := 0; i < len(s.tokens); i++ {
		token := s.tokens[i]
		
		// Buscar declaraciones de variables (let, const, var, int, etc.)
		if token.Type == KEYWORD || token.Type == TYPE {
			if i+1 < len(s.tokens) && s.tokens[i+1].Type == IDENTIFIER {
				varName := s.tokens[i+1].Value
				varType := s.inferType(token.Value)
				
				// Verificar si hay declaración de tipo TypeScript (: tipo)
				nextPos := i + 2
				if nextPos < len(s.tokens) && s.tokens[nextPos].Type == COLON {
					nextPos++ // saltar ':'
					if nextPos < len(s.tokens) && (s.tokens[nextPos].Type == TYPE || s.tokens[nextPos].Type == IDENTIFIER) {
						varType = s.tokens[nextPos].Value
						nextPos++ // saltar tipo
					}
				}
				
				var initialValue string
				if nextPos < len(s.tokens) && s.tokens[nextPos].Type == ASSIGNMENT {
					nextPos++ // saltar '='
					if nextPos < len(s.tokens) {
						initialValue = s.tokens[nextPos].Value
					}
				}
				
				s.variables[varName] = VariableInfo{
					Name:         varName,
					Type:         varType,
					InitialValue: initialValue,
					Line:         token.Line,
					Column:       token.Column,
				}
				
				s.addInfo(fmt.Sprintf("Variable '%s' declarada como tipo '%s' con valor inicial '%s' en línea %d", 
					varName, varType, initialValue, token.Line))
				
				// Saltar los tokens procesados
				i = nextPos
			}
		}
	}
}

func (s *Semantic) analyzeForLoop() {
	forFound := false
	var loopVar string
	var conditionVar string
	var incrementVar string
	var startValue, endValue int
	var hasCondition bool
	
	for i := 0; i < len(s.tokens); i++ {
		token := s.tokens[i]
		
		if token.Type == FOR {
			forFound = true
			s.addInfo("Bucle 'for' detectado - Analizando estructura")
		}
		
		// Analizar inicialización del bucle
		if forFound && token.Type == IDENTIFIER && i > 0 && 
			(s.tokens[i-1].Type == KEYWORD || s.tokens[i-1].Type == TYPE) {
			loopVar = token.Value
			
			// Buscar valor inicial
			if i+2 < len(s.tokens) && s.tokens[i+1].Type == ASSIGNMENT {
				if val, err := strconv.Atoi(s.tokens[i+2].Value); err == nil {
					startValue = val
					s.addInfo(fmt.Sprintf("Variable de control '%s' inicializada con valor %d", loopVar, startValue))
				}
			}
		}
		
		// Analizar condición del bucle
		if forFound && token.Type == COMPARISON && i > 0 && i+1 < len(s.tokens) {
			hasCondition = true
			operator := token.Value
			
			if s.tokens[i-1].Type == IDENTIFIER && s.tokens[i+1].Type == NUMBER {
				conditionVar = s.tokens[i-1].Value
				if val, err := strconv.Atoi(s.tokens[i+1].Value); err == nil {
					endValue = val
					s.addInfo(fmt.Sprintf("Condición: '%s %s %d' - Variable de control se compara con %d", 
						conditionVar, operator, endValue, endValue))
					
					// Verificar coherencia de la condición
					s.checkLoopConditionCoherence(operator, startValue, endValue)
				}
			}
		}
		
		// Analizar incremento
		if forFound && token.Type == INCREMENT {
			if i > 0 && s.tokens[i-1].Type == IDENTIFIER {
				incrementVar = s.tokens[i-1].Value
			} else if i+1 < len(s.tokens) && s.tokens[i+1].Type == IDENTIFIER {
				incrementVar = s.tokens[i+1].Value
			}
			
			s.addInfo(fmt.Sprintf("Incremento detectado para variable '%s' (%s)", incrementVar, token.Value))
		}
	}
	
	// Verificar consistencia de variables en el bucle
	if forFound && loopVar != "" {
		s.checkLoopVariableConsistency(loopVar, conditionVar, incrementVar)
	}
	
	// Análisis completo del bucle
	if forFound && hasCondition {
		iterations := s.calculateIterations(startValue, endValue)
		if iterations > 0 {
			s.addInfo(fmt.Sprintf("El bucle ejecutará aproximadamente %d iteraciones", iterations))
		}
	}
}

func (s *Semantic) checkLoopConditionCoherence(operator string, start, end int) {
	switch operator {
	case "<=", "<":
		if start > end {
			s.addInfo("⚠️ ADVERTENCIA: La condición del bucle podría nunca ser verdadera (valor inicial mayor que final)")
		}
	case ">=", ">":
		if start < end {
			s.addInfo("⚠️ ADVERTENCIA: La condición del bucle podría nunca ser verdadera (valor inicial menor que final)")
		}
	}
}

func (s *Semantic) calculateIterations(start, end int) int {
	if start <= end {
		return end - start + 1
	}
	return 0
}

func (s *Semantic) checkVariableUsage() {
	usedVars := make(map[string]bool)
	
	for _, token := range s.tokens {
		if token.Type == IDENTIFIER {
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
			varInfo=varInfo
		}
	}
}

func (s *Semantic) analyzeInfiniteLoop() {
	// Detectar posibles bucles infinitos
	hasIncrement := false
	hasValidCondition := false
	
	for _, token := range s.tokens {
		if token.Type == INCREMENT {
			hasIncrement = true
		}
		if token.Type == COMPARISON {
			hasValidCondition = true
		}
	}
	
	if !hasIncrement && hasValidCondition {
		s.addInfo("⚠️ POSIBLE BUCLE INFINITO: No se detectó incremento en la variable de control")
	} else if hasIncrement && hasValidCondition {
		s.addInfo("✓ Estructura de bucle válida: tiene condición e incremento")
	}
}

func (s *Semantic) inferType(declaration string) string {
	switch strings.ToLower(declaration) {
	case "int":
		return "number"
	case "string":
		return "string"
	case "let":
		return "variable"
	case "const":
		return "constant"
	case "var":
		return "variable"
	default:
		return "unknown"
	}
}

func (s *Semantic) checkLoopVariableConsistency(loopVar, conditionVar, incrementVar string) {
	// Verificar que la variable de condición sea la misma que la declarada
	if conditionVar != "" && conditionVar != loopVar {
		s.addInfo(fmt.Sprintf("❌ ERROR SEMÁNTICO: Variable en condición '%s' no coincide con variable de control '%s'", 
			conditionVar, loopVar))
	} else if conditionVar == loopVar {
		s.addInfo(fmt.Sprintf("✓ Variable de condición '%s' coincide correctamente con variable de control", conditionVar))
	}
	
	// Verificar que la variable de incremento sea la misma que la declarada
	if incrementVar != "" && incrementVar != loopVar {
		s.addInfo(fmt.Sprintf("❌ ERROR SEMÁNTICO: Variable en incremento '%s' no coincide con variable de control '%s'", 
			incrementVar, loopVar))
	} else if incrementVar == loopVar {
		s.addInfo(fmt.Sprintf("✓ Variable de incremento '%s' coincide correctamente con variable de control", incrementVar))
	}
	
	// Verificar que la variable de control esté siendo utilizada en todas las partes
	if conditionVar == "" {
		s.addInfo("⚠️ ADVERTENCIA: No se detectó variable en la condición del bucle")
	}
	
	if incrementVar == "" {
		s.addInfo("⚠️ ADVERTENCIA: No se detectó variable en el incremento del bucle")
	}
}