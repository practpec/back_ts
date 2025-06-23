package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type SemanticC struct {
	tokens      []TokenC
	variables   map[string]VariableInfoC
	information []string
}

type VariableInfoC struct {
	Name         string
	Type         string
	InitialValue string
	Line         int
	Column       int
}

func NewSemanticC(tokens []TokenC) *SemanticC {
	return &SemanticC{
		tokens:      tokens,
		variables:   make(map[string]VariableInfoC),
		information: []string{},
	}
}

func (s *SemanticC) addInfo(message string) {
	s.information = append(s.information, message)
}

func (s *SemanticC) detectMalformedNumbers() {
	for _, token := range s.tokens {
		if token.Type == UNKNOWN_C {
			// Verificar si parece un número mal formado
			if len(token.Value) > 0 && unicode.IsDigit(rune(token.Value[0])) {
				s.addInfo(fmt.Sprintf("❌ ERROR LÉXICO: Número mal formado '%s' en línea %d, columna %d",
					token.Value, token.Line, token.Column))
			}
		}
	}
}

func (s *SemanticC) detectInvalidExpressions() {
	for i := 0; i < len(s.tokens)-1; i++ {
		currentToken := s.tokens[i]
		nextToken := s.tokens[i+1]

		// Verificar patrones inválidos de tokens consecutivos
		if currentToken.Type == NUMBER_C && nextToken.Type == IDENTIFIER_C {
			if currentToken.Line == nextToken.Line {
				s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Número '%s' seguido de identificador '%s' sin operador en línea %d",
					currentToken.Value, nextToken.Value, currentToken.Line))
			}
		}

		if currentToken.Type == IDENTIFIER_C && nextToken.Type == NUMBER_C {
			if currentToken.Line == nextToken.Line && !s.isReservedWord(currentToken.Value) {
				s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Identificador '%s' seguido de número '%s' sin operador en línea %d",
					currentToken.Value, nextToken.Value, currentToken.Line))
			}
		}

		if currentToken.Type == NUMBER_C && nextToken.Type == NUMBER_C {
			if currentToken.Line == nextToken.Line {
				s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Dos números consecutivos '%s' '%s' sin operador en línea %d",
					currentToken.Value, nextToken.Value, currentToken.Line))
			}
		}

		if currentToken.Type == IDENTIFIER_C && nextToken.Type == IDENTIFIER_C {
			if currentToken.Line == nextToken.Line &&
				!s.isReservedWord(currentToken.Value) &&
				!s.isReservedWord(nextToken.Value) {
				s.addInfo(fmt.Sprintf("❌ ERROR SINTÁCTICO: Dos identificadores consecutivos '%s' '%s' sin operador en línea %d",
					currentToken.Value, nextToken.Value, currentToken.Line))
			}
		}
	}
}

func (s *SemanticC) analyzeDoWhileLoop() {
	doFound := false
	whileFound := false
	var conditionVar string

	for i := 0; i < len(s.tokens); i++ {
		token := s.tokens[i]

		if token.Type == DO_C {
			doFound = true
			s.addInfo("Bucle 'do-while' detectado - Analizando estructura")
		}

		if token.Type == WHILE_C && doFound {
			whileFound = true
			s.addInfo("Cláusula 'while' encontrada en bucle do-while")

			// Buscar la condición después de while
			if i+2 < len(s.tokens) && s.tokens[i+1].Type == LPAREN_C {
				condPos := i + 2
				for condPos < len(s.tokens) && s.tokens[condPos].Type != RPAREN_C {
					if s.tokens[condPos].Type == IDENTIFIER_C {
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

func (s *SemanticC) detectUndeclaredVariables() {
	// Encontrar todas las variables usadas
	usedVariables := make(map[string][]int) // variable -> posiciones donde se usa

	for i, token := range s.tokens {
		if token.Type == IDENTIFIER_C {
			// Excluir palabras reservadas como 'printf', 'scanf', etc.
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

func (s *SemanticC) isReservedWord(word string) bool {
	reservedWords := map[string]bool{
		"printf":  true,
		"scanf":   true,
		"main":    true,
		"stdio":   true,
		"include": true,
		"return":  true,
		"if":      true,
		"else":    true,
		"for":     true,
		"while":   true,
		"do":      true,
		"break":   true,
		"continue": true,
	}
	return reservedWords[strings.ToLower(word)]
}



func (s *SemanticC) analyzeVariableDeclarations() {
	for i := 0; i < len(s.tokens); i++ {
		token := s.tokens[i]

		// Buscar declaraciones de variables (int, float, double, etc.)
		if token.Type == TYPE_C {
			if i+1 < len(s.tokens) && s.tokens[i+1].Type == IDENTIFIER_C {
				varName := s.tokens[i+1].Value
				varType := token.Value

				var initialValue string
				nextPos := i + 2
				if nextPos < len(s.tokens) && s.tokens[nextPos].Type == ASSIGNMENT_C {
					nextPos++ // saltar '='
					if nextPos < len(s.tokens) {
						initialValue = s.tokens[nextPos].Value
					}
				}

				s.variables[varName] = VariableInfoC{
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

func (s *SemanticC) analyzeForLoop() {
	forFound := false
	var loopVar string
	var conditionVar string
	var incrementVar string
	var startValue, endValue int
	var hasCondition bool

	for i := 0; i < len(s.tokens); i++ {
		token := s.tokens[i]

		if token.Type == FOR_C {
			forFound = true
			s.addInfo("Bucle 'for' detectado - Analizando estructura")
		}

		// Analizar inicialización del bucle
		if forFound && token.Type == IDENTIFIER_C && i > 0 &&
			s.tokens[i-1].Type == TYPE_C {
			loopVar = token.Value

			// Buscar valor inicial
			if i+2 < len(s.tokens) && s.tokens[i+1].Type == ASSIGNMENT_C {
				if val, err := strconv.Atoi(s.tokens[i+2].Value); err == nil {
					startValue = val
					s.addInfo(fmt.Sprintf("Variable de control '%s' inicializada con valor %d", loopVar, startValue))
				}
			}
		}

		// Analizar condición del bucle
		if forFound && token.Type == COMPARISON_C && i > 0 && i+1 < len(s.tokens) {
			hasCondition = true
			operator := token.Value

			if s.tokens[i-1].Type == IDENTIFIER_C && s.tokens[i+1].Type == NUMBER_C {
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
		if forFound && token.Type == INCREMENT_C {
			if i > 0 && s.tokens[i-1].Type == IDENTIFIER_C {
				incrementVar = s.tokens[i-1].Value
			} else if i+1 < len(s.tokens) && s.tokens[i+1].Type == IDENTIFIER_C {
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

func (s *SemanticC) checkLoopConditionCoherence(operator string, start, end int) {
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

func (s *SemanticC) calculateIterations(start, end int) int {
	if start <= end {
		return end - start + 1
	}
	return 0
}

func (s *SemanticC) checkVariableUsage() {
	usedVars := make(map[string]bool)

	for _, token := range s.tokens {
		if token.Type == IDENTIFIER_C {
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
			_ = varInfo // evitar warning de variable no usada
		}
	}
}

func (s *SemanticC) analyzeInfiniteLoop() {
	// Detectar posibles bucles infinitos
	hasIncrement := false
	hasValidCondition := false

	for _, token := range s.tokens {
		if token.Type == INCREMENT_C {
			hasIncrement = true
		}
		if token.Type == COMPARISON_C {
			hasValidCondition = true
		}
	}

	if !hasIncrement && hasValidCondition {
		s.addInfo("⚠️ POSIBLE BUCLE INFINITO: No se detectó incremento en la variable de control")
	} else if hasIncrement && hasValidCondition {
		s.addInfo("✓ Estructura de bucle válida: tiene condición e incremento")
	}
}

func (s *SemanticC) checkLoopVariableConsistency(loopVar, conditionVar, incrementVar string) {
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

// Agregar esta función en semantic_c.go

func (s *SemanticC) detectSelfReferenceInInitialization() {
	for i := 0; i < len(s.tokens)-3; i++ {
		// Buscar patrón: TYPE IDENTIFIER = IDENTIFIER
		if s.tokens[i].Type == TYPE_C &&
			i+1 < len(s.tokens) && s.tokens[i+1].Type == IDENTIFIER_C &&
			i+2 < len(s.tokens) && s.tokens[i+2].Type == ASSIGNMENT_C &&
			i+3 < len(s.tokens) && s.tokens[i+3].Type == IDENTIFIER_C {
			
			varName := s.tokens[i+1].Value
			initValue := s.tokens[i+3].Value
			
			// Si la variable se inicializa con ella misma
			if varName == initValue {
				s.addInfo(fmt.Sprintf("❌ ERROR SEMÁNTICO: Variable '%s' se inicializa con ella misma en línea %d - Esto causa comportamiento indefinido", 
					varName, s.tokens[i].Line))
			}
		}
	}
}

// También agregar esta función para detectar uso de variables antes de declaración
func (s *SemanticC) detectUseBeforeDeclaration() {
	declaredVars := make(map[string]int) // variable -> línea donde se declara
	
	// Primera pasada: recopilar todas las declaraciones
	for i := 0; i < len(s.tokens)-1; i++ {
		if s.tokens[i].Type == TYPE_C && 
			i+1 < len(s.tokens) && s.tokens[i+1].Type == IDENTIFIER_C {
			varName := s.tokens[i+1].Value
			declaredVars[varName] = s.tokens[i].Line
		}
	}
	
	// Segunda pasada: verificar uso antes de declaración
	for i := 0; i < len(s.tokens); i++ {
		token := s.tokens[i]
		if token.Type == IDENTIFIER_C && !s.isReservedWord(token.Value) {
			// Verificar si se usa antes de ser declarada
			if declaredLine, exists := declaredVars[token.Value]; exists {
				if token.Line < declaredLine {
					s.addInfo(fmt.Sprintf("❌ ERROR SEMÁNTICO: Variable '%s' usada en línea %d antes de ser declarada en línea %d", 
						token.Value, token.Line, declaredLine))
				}
			}
		}
	}
}

// Modificar la función Analyze() para incluir estas verificaciones:
func (s *SemanticC) Analyze() []string {
	s.analyzeVariableDeclarations()
	s.detectSelfReferenceInInitialization()  // Nueva función
	s.detectUseBeforeDeclaration()           // Nueva función
	s.analyzeForLoop()
	s.analyzeDoWhileLoop()
	s.checkVariableUsage()
	s.analyzeInfiniteLoop()
	s.detectUndeclaredVariables()
	s.detectMalformedNumbers()
	s.detectInvalidExpressions()
	return s.information
}