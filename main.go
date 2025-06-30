package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type AnalysisRequestJ struct {
	Code string `json:"code"`
}

type AnalysisResponseJ struct {
	IsValid      bool     `json:"isValid"`
	Tokens       []TokenJ `json:"tokens"`
	SyntaxErrors []string `json:"syntaxErrors"`
	SemanticInfo []string `json:"semanticInfo"`
}

func main() {
	r := mux.NewRouter()

	// CORS headers
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	// Endpoint para análisis de código Java
	r.HandleFunc("/analyze-java", analyzeJavaHandler).Methods("POST")
	
	// Mantener endpoints anteriores (comentados por ahora)
	// r.HandleFunc("/analyze-c", analyzeCHandler).Methods("POST")
	// r.HandleFunc("/analyze", analyzeHandler).Methods("POST")

	fmt.Println("Servidor iniciado en puerto 8080")
	fmt.Println("Endpoint disponible: /analyze-java (para código Java)")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(r)))
}

func analyzeJavaHandler(w http.ResponseWriter, r *http.Request) {
	var req AnalysisRequestJ
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	// Análisis léxico
	lexer := NewLexerJ(req.Code)
	tokens := lexer.Tokenize()

	// Análisis sintáctico
	parser := NewParserJ(tokens)
	syntaxErrors := parser.Parse()

	// Análisis semántico
	semantic := NewSemanticJ(tokens)
	semanticInfo := semantic.Analyze()

	// Verificar si hay errores semánticos
	hasSemanticErrors := false
	for _, info := range semanticInfo {
		if strings.Contains(info, "❌ ERROR SEMÁNTICO") || strings.Contains(info, "❌ ERROR SINTÁCTICO") || strings.Contains(info, "❌ ERROR LÉXICO") {
			hasSemanticErrors = true
			break
		}
	}

	response := AnalysisResponseJ{
		IsValid:      len(syntaxErrors) == 0 && !hasSemanticErrors,
		Tokens:       tokens,
		SyntaxErrors: syntaxErrors,
		SemanticInfo: semanticInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}