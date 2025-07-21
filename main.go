package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

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

type PerformanceMetrics struct {
	ExecutionTime   string  `json:"executionTime"`
	MemoryUsage     string  `json:"memoryUsage"`
	AllocatedBytes  uint64  `json:"allocatedBytes"`
	TotalAllocs     uint64  `json:"totalAllocs"`
	GCCycles        uint32  `json:"gcCycles"`
	CPUUsage        float64 `json:"cpuUsage"`
}

type AnalysisWithMetrics struct {
	AnalysisResponse
	Metrics PerformanceMetrics `json:"metrics"`
}

func main() {
	r := mux.NewRouter()
	
	// CORS headers
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	
	// Endpoints existentes
	r.HandleFunc("/analyze", analyzeHandler).Methods("POST")
	
	// Nuevos endpoints para comparación de rendimiento
	r.HandleFunc("/analyze-optimized", analyzeOptimizedHandler).Methods("POST")
	r.HandleFunc("/analyze-unoptimized", analyzeUnoptimizedHandler).Methods("POST")
	
	fmt.Println("Servidor iniciado en puerto 8080")
	fmt.Println("Endpoints disponibles:")
	fmt.Println("  POST /analyze - Análisis existente")
	fmt.Println("  POST /analyze-optimized - Análisis optimizado con métricas")
	fmt.Println("  POST /analyze-unoptimized - Análisis NO optimizado con métricas")
	
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(r)))
}

// Handler existente sin cambios
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
		if strings.Contains(info, "❌ ERROR SEMÁNTICO") || strings.Contains(info, "❌ ERROR SINTÁCTICO") || strings.Contains(info, "❌ ERROR LÉXICO") {
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

// Nuevo handler para análisis optimizado con métricas
func analyzeOptimizedHandler(w http.ResponseWriter, r *http.Request) {
	var req AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}
	
	// Medir métricas antes del análisis
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	
	startTime := time.Now()
	startCPU := time.Now()
	
	// Análisis léxico optimizado (código existente)
	lexer := NewLexer(req.Code)
	tokens := lexer.Tokenize()
	
	// Análisis sintáctico optimizado (código existente)
	parser := NewParser(tokens)
	syntaxErrors := parser.Parse()
	
	// Análisis semántico optimizado (código existente)
	semantic := NewSemantic(tokens)
	semanticInfo := semantic.Analyze()
	
	// Medir métricas después del análisis
	executionTime := time.Since(startTime)
	cpuTime := time.Since(startCPU)
	
	runtime.GC()
	runtime.ReadMemStats(&m2)
	
	// Verificar si hay errores semánticos
	hasSemanticErrors := false
	for _, info := range semanticInfo {
		if strings.Contains(info, "❌ ERROR SEMÁNTICO") || strings.Contains(info, "❌ ERROR SINTÁCTICO") || strings.Contains(info, "❌ ERROR LÉXICO") {
			hasSemanticErrors = true
			break
		}
	}
	
	metrics := PerformanceMetrics{
		ExecutionTime:   executionTime.String(),
		MemoryUsage:     fmt.Sprintf("%.2f KB", float64(m2.Alloc-m1.Alloc)/1024),
		AllocatedBytes:  m2.Alloc - m1.Alloc,
		TotalAllocs:     m2.TotalAlloc - m1.TotalAlloc,
		GCCycles:        m2.NumGC - m1.NumGC,
		CPUUsage:        cpuTime.Seconds(),
	}
	
	response := AnalysisWithMetrics{
		AnalysisResponse: AnalysisResponse{
			IsValid:      len(syntaxErrors) == 0 && !hasSemanticErrors,
			Tokens:       tokens,
			SyntaxErrors: syntaxErrors,
			SemanticInfo: semanticInfo,
		},
		Metrics: metrics,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Nuevo handler para análisis NO optimizado con métricas
func analyzeUnoptimizedHandler(w http.ResponseWriter, r *http.Request) {
	var req AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}
	
	// Medir métricas antes del análisis
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	
	startTime := time.Now()
	startCPU := time.Now()
	
	// Análisis léxico NO optimizado
	lexerUnoptimized := NewLexerUnoptimized(req.Code)
	tokens := lexerUnoptimized.TokenizeUnoptimized()
	
	// Análisis sintáctico NO optimizado
	parserUnoptimized := NewParserUnoptimized(tokens)
	syntaxErrors := parserUnoptimized.ParseUnoptimized()
	
	// Análisis semántico NO optimizado
	semanticUnoptimized := NewSemanticUnoptimized(tokens)
	semanticInfo := semanticUnoptimized.AnalyzeUnoptimized()
	
	// Medir métricas después del análisis
	executionTime := time.Since(startTime)
	cpuTime := time.Since(startCPU)
	
	runtime.GC()
	runtime.ReadMemStats(&m2)
	
	// Verificar si hay errores semánticos
	hasSemanticErrors := false
	for _, info := range semanticInfo {
		if strings.Contains(info, "❌ ERROR SEMÁNTICO") || strings.Contains(info, "❌ ERROR SINTÁCTICO") || strings.Contains(info, "❌ ERROR LÉXICO") {
			hasSemanticErrors = true
			break
		}
	}
	
	metrics := PerformanceMetrics{
		ExecutionTime:   executionTime.String(),
		MemoryUsage:     fmt.Sprintf("%.2f KB", float64(m2.Alloc-m1.Alloc)/1024),
		AllocatedBytes:  m2.Alloc - m1.Alloc,
		TotalAllocs:     m2.TotalAlloc - m1.TotalAlloc,
		GCCycles:        m2.NumGC - m1.NumGC,
		CPUUsage:        cpuTime.Seconds(),
	}
	
	response := AnalysisWithMetrics{
		AnalysisResponse: AnalysisResponse{
			IsValid:      len(syntaxErrors) == 0 && !hasSemanticErrors,
			Tokens:       tokens,
			SyntaxErrors: syntaxErrors,
			SemanticInfo: semanticInfo,
		},
		Metrics: metrics,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}