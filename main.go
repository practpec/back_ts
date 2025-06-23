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

type AnalysisRequest struct {
	Code string `json:"code"`
}

type AnalysisResponse struct {
	IsValid      bool     `json:"isValid"`
	Tokens       []Token  `json:"tokens"`
	SyntaxErrors []string `json:"syntaxErrors"`
	SemanticInfo []string `json:"semanticInfo"`
}

func main() {
	r := mux.NewRouter()
	
	// CORS headers
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	
	r.HandleFunc("/analyze", analyzeHandler).Methods("POST")
	
	fmt.Println("Servidor iniciado en puerto 8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(r)))
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	var req AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}
	
	// Análisis léxico
	lexer := NewLexer(req.Code)
	tokens := lexer.Tokenize()
	
	// Análisis sintáctico
	parser := NewParser(tokens)
	syntaxErrors := parser.Parse()
	
	// Análisis semántico
	semantic := NewSemantic(tokens)
	semanticInfo := semantic.Analyze()
	
	// Verificar si hay errores semánticos
	hasSemanticErrors := false
	for _, info := range semanticInfo {
		if strings.Contains(info, "❌ ERROR SEMÁNTICO") {
			hasSemanticErrors = true
			break
		}
	}
	
	response := AnalysisResponse{
		IsValid:      len(syntaxErrors) == 0 && !hasSemanticErrors,
		Tokens:       tokens,
		SyntaxErrors: syntaxErrors,
		SemanticInfo: semanticInfo,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}