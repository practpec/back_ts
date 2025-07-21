package main

import (
	"strconv"
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

// Mapa estático global para palabras reservadas (máxima eficiencia)
var reservedWords = map[string]bool{
	"console": true,
	"log":     true,
	"system":  true,
	"out":     true,
	"println": true,
	"print":   true,
	"length":  true,
}

// Strings constantes para tipos de inferencia (evitar creaciones repetidas)
const (
	typeNumber   = "number"
	typeString   = "string"
	typeVariable = "variable"
	typeConstant = "constant"
	typeUnknown  = "unknown"
)

func NewSemantic(tokens []Token) *Semantic {
	return &Semantic{
		tokens:      tokens,
		variables:   make(map[string]VariableInfo, 16), // Pre-allocar con capacidad
		information: make([]string, 0, 32),             // Pre-allocar
	}
}

func (s *Semantic) addInfo(message string) {
	s.information = append(s.information, message)
}

func (s *Semantic) Analyze() []string {
	s.analyzeVariableDeclarations()
	s.analyzeForLoop()
	s.checkVariableUsage()
	s.analyzeInfiniteLoop()
	s.detectUndeclaredVariables()
	s.detectMalformedNumbers()
	s.detectInvalidExpressions()
	s.analyzeDoWhileLoop()
	return s.information
}

func (s *Semantic) detectMalformedNumbers() {
	for i := range s.tokens {
		if s.tokens[i].Type == UNKNOWN && len(s.tokens[i].Value) > 0 && 
		   unicode.IsDigit(rune(s.tokens[i].Value[0])) {
			s.addInfo("❌ ERROR LÉXICO: Número mal formado '" + s.tokens[i].Value + 
				"' en línea " + strconv.Itoa(s.tokens[i].Line) + 
				", columna " + strconv.Itoa(s.tokens[i].Column))
		}
	}
}

func (s *Semantic) detectInvalidExpressions() {
	for i := 0; i < len(s.tokens)-1; i++ {
		current := &s.tokens[i]
		next := &s.tokens[i+1]
		
		if current.Line != next.Line {
			continue
		}
		
		switch {
		case current.Type == NUMBER && next.Type == IDENTIFIER:
			s.addInfo("❌ ERROR SINTÁCTICO: Número '" + current.Value + 
				"' seguido de identificador '" + next.Value + 
				"' sin operador en línea " + strconv.Itoa(current.Line))
		case current.Type == IDENTIFIER && next.Type == NUMBER && !s.isReservedWord(current.Value):
			s.addInfo("❌ ERROR SINTÁCTICO: Identificador '" + current.Value + 
				"' seguido de número '" + next.Value + 
				"' sin operador en línea " + strconv.Itoa(current.Line))
		case current.Type == NUMBER && next.Type == NUMBER:
			s.addInfo("❌ ERROR SINTÁCTICO: Dos números consecutivos '" + current.Value + 
				"' '" + next.Value + "' sin operador en línea " + strconv.Itoa(current.Line))
		case current.Type == IDENTIFIER && next.Type == IDENTIFIER && 
			!s.isReservedWord(current.Value) && !s.isReservedWord(next.Value):
			s.addInfo("❌ ERROR SINTÁCTICO: Dos identificadores consecutivos '" + current.Value + 
				"' '" + next.Value + "' sin operador en línea " + strconv.Itoa(current.Line))
		}
	}
}

func (s *Semantic) analyzeDoWhileLoop() {
	doFound := false
	whileFound := false
	var conditionVar string
	
	for i := 0; i < len(s.tokens); i++ {
		token := &s.tokens[i]
		
		if token.Type == DO {
			doFound = true
			s.addInfo("Bucle 'do-while' detectado - Analizando estructura")
		}
		
		if token.Type == WHILE && doFound {
			whileFound = true
			s.addInfo("Cláusula 'while' encontrada en bucle do-while")
			
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
					s.addInfo("Variable en condición do-while: '" + conditionVar + "'")
					
					if _, exists := s.variables[conditionVar]; !exists {
						s.addInfo("❌ ERROR SEMÁNTICO: Variable '" + conditionVar + 
							"' en condición do-while no está declarada")
					} else {
						s.addInfo("✓ Variable '" + conditionVar + 
							"' en condición do-while está correctamente declarada")
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
	usedVariables := make(map[string]int, 16) // Pre-allocar
	
	for i := range s.tokens {
		if s.tokens[i].Type == IDENTIFIER && !s.isReservedWord(s.tokens[i].Value) {
			if _, exists := usedVariables[s.tokens[i].Value]; !exists {
				usedVariables[s.tokens[i].Value] = i
			}
		}
	}
	
	for varName, tokenIndex := range usedVariables {
		if _, declared := s.variables[varName]; !declared {
			s.addInfo("❌ ERROR SEMÁNTICO: Variable '" + varName + 
				"' usada sin declarar (línea " + strconv.Itoa(s.tokens[tokenIndex].Line) + ")")
		}
	}
}

// Función optimizada con lookup directo
func (s *Semantic) isReservedWord(word string) bool {
	return reservedWords[word]
}

func (s *Semantic) analyzeVariableDeclarations() {
	for i := 0; i < len(s.tokens); i++ {
		token := &s.tokens[i]
		
		if (token.Type == KEYWORD || token.Type == TYPE) && 
		   i+1 < len(s.tokens) && s.tokens[i+1].Type == IDENTIFIER {
			
			varName := s.tokens[i+1].Value
			varType := s.inferType(token.Value)
			
			nextPos := i + 2
			if nextPos < len(s.tokens) && s.tokens[nextPos].Type == COLON {
				nextPos++
				if nextPos < len(s.tokens) && 
				   (s.tokens[nextPos].Type == TYPE || s.tokens[nextPos].Type == IDENTIFIER) {
					varType = s.tokens[nextPos].Value
					nextPos++
				}
			}
			
			var initialValue string
			if nextPos < len(s.tokens) && s.tokens[nextPos].Type == ASSIGNMENT {
				nextPos++
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
			
			s.addInfo("Variable '" + varName + "' declarada como tipo '" + varType + 
				"' con valor inicial '" + initialValue + "' en línea " + strconv.Itoa(token.Line))
			
			i = nextPos
		}
	}
}

func (s *Semantic) analyzeForLoop() {
	forFound := false
	var loopVar, conditionVar, incrementVar string
	var startValue, endValue int
	var hasCondition bool
	
	for i := 0; i < len(s.tokens); i++ {
		token := &s.tokens[i]
		
		if token.Type == FOR {
			forFound = true
			s.addInfo("Bucle 'for' detectado - Analizando estructura")
		}
		
		if forFound && token.Type == IDENTIFIER && i > 0 && 
		   (s.tokens[i-1].Type == KEYWORD || s.tokens[i-1].Type == TYPE) {
			loopVar = token.Value
			
			if i+2 < len(s.tokens) && s.tokens[i+1].Type == ASSIGNMENT {
				if val, err := strconv.Atoi(s.tokens[i+2].Value); err == nil {
					startValue = val
					s.addInfo("Variable de control '" + loopVar + "' inicializada con valor " + 
						strconv.Itoa(startValue))
				}
			}
		}
		
		if forFound && token.Type == COMPARISON && i > 0 && i+1 < len(s.tokens) {
			hasCondition = true
			operator := token.Value
			
			if s.tokens[i-1].Type == IDENTIFIER && s.tokens[i+1].Type == NUMBER {
				conditionVar = s.tokens[i-1].Value
				if val, err := strconv.Atoi(s.tokens[i+1].Value); err == nil {
					endValue = val
					s.addInfo("Condición: '" + conditionVar + " " + operator + " " + 
						strconv.Itoa(endValue) + "' - Variable de control se compara con " + 
						strconv.Itoa(endValue))
					
					s.checkLoopConditionCoherence(operator, startValue, endValue)
				}
			}
		}
		
		if forFound && token.Type == INCREMENT {
			if i > 0 && s.tokens[i-1].Type == IDENTIFIER {
				incrementVar = s.tokens[i-1].Value
			} else if i+1 < len(s.tokens) && s.tokens[i+1].Type == IDENTIFIER {
				incrementVar = s.tokens[i+1].Value
			}
			
			s.addInfo("Incremento detectado para variable '" + incrementVar + "' (" + token.Value + ")")
		}
	}
	
	if forFound && loopVar != "" {
		s.checkLoopVariableConsistency(loopVar, conditionVar, incrementVar)
	}
	
	if forFound && hasCondition {
		iterations := s.calculateIterations(startValue, endValue)
		if iterations > 0 {
			s.addInfo("El bucle ejecutará aproximadamente " + strconv.Itoa(iterations) + " iteraciones")
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
	usedVars := make(map[string]bool, len(s.variables))
	
	for i := range s.tokens {
		if s.tokens[i].Type == IDENTIFIER {
			if _, exists := s.variables[s.tokens[i].Value]; exists {
				usedVars[s.tokens[i].Value] = true
			}
		}
	}
	
	for varName := range s.variables {
		if !usedVars[varName] {
			s.addInfo("⚠️ Variable '" + varName + "' declarada pero no utilizada")
		} else {
			s.addInfo("✓ Variable '" + varName + "' declarada y utilizada correctamente")
		}
	}
}

func (s *Semantic) analyzeInfiniteLoop() {
	hasIncrement := false
	hasValidCondition := false
	
	for i := range s.tokens {
		switch s.tokens[i].Type {
		case INCREMENT:
			hasIncrement = true
		case COMPARISON:
			hasValidCondition = true
		}
	}
	
	if !hasIncrement && hasValidCondition {
		s.addInfo("⚠️ POSIBLE BUCLE INFINITO: No se detectó incremento en la variable de control")
	} else if hasIncrement && hasValidCondition {
		s.addInfo("✓ Estructura de bucle válida: tiene condición e incremento")
	}
}

// Función optimizada con switch
func (s *Semantic) inferType(declaration string) string {
	switch declaration {
	case "int":
		return typeNumber
	case "string":
		return typeString
	case "let":
		return typeVariable
	case "const":
		return typeConstant
	case "var":
		return typeVariable
	default:
		return typeUnknown
	}
}

func (s *Semantic) checkLoopVariableConsistency(loopVar, conditionVar, incrementVar string) {
	if conditionVar != "" && conditionVar != loopVar {
		s.addInfo("❌ ERROR SEMÁNTICO: Variable en condición '" + conditionVar + 
			"' no coincide con variable de control '" + loopVar + "'")
	} else if conditionVar == loopVar {
		s.addInfo("✓ Variable de condición '" + conditionVar + 
			"' coincide correctamente con variable de control")
	}
	
	if incrementVar != "" && incrementVar != loopVar {
		s.addInfo("❌ ERROR SEMÁNTICO: Variable en incremento '" + incrementVar + 
			"' no coincide con variable de control '" + loopVar + "'")
	} else if incrementVar == loopVar {
		s.addInfo("✓ Variable de incremento '" + incrementVar + 
			"' coincide correctamente con variable de control")
	}
	
	if conditionVar == "" {
		s.addInfo("⚠️ ADVERTENCIA: No se detectó variable en la condición del bucle")
	}
	
	if incrementVar == "" {
		s.addInfo("⚠️ ADVERTENCIA: No se detectó variable en el incremento del bucle")
	}
}