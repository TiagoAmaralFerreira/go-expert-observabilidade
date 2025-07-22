package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type CEPRequest struct {
	CEP string `json:"cep"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

var tracer trace.Tracer

func initTracer() (*sdktrace.TracerProvider, error) {
	zipkinURL := os.Getenv("ZIPKIN_URL")
	if zipkinURL == "" {
		zipkinURL = "http://localhost:9411"
	}

	exporter, err := zipkin.New(zipkinURL + "/api/v2/spans")
	if err != nil {
		return nil, fmt.Errorf("failed to create zipkin exporter: %w", err)
	}

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("service-a"),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("service-a")
	return tp, nil
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func validateCEP(cep string) bool {
	// Verifica se contém exatamente 8 dígitos
	matched, _ := regexp.MatchString(`^\d{8}$`, cep)
	return matched
}

func forwardToServiceB(ctx context.Context, cep string) ([]byte, int, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	serviceBURL := os.Getenv("SERVICE_B_URL")
	if serviceBURL == "" {
		serviceBURL = "http://localhost:8081"
	}

	// Cria o payload para o Serviço B
	payload := CEPRequest{CEP: cep}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		span.RecordError(err)
		return nil, 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Cria a requisição HTTP
	req, err := http.NewRequestWithContext(ctx, "POST", serviceBURL+"/weather", bytes.NewBuffer(jsonData))
	if err != nil {
		span.RecordError(err)
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Faz a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		return nil, 0, fmt.Errorf("failed to send request to service B: %w", err)
	}
	defer resp.Body.Close()

	// Lê a resposta do Serviço B
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.RecordError(err)
		return nil, 0, fmt.Errorf("failed to read response body: %w", err)
	}

	// Copia a resposta do Serviço B para o cliente
	span.SetAttributes(semconv.HTTPStatusCode(resp.StatusCode))

	return body, resp.StatusCode, nil
}

func cepHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	// Verifica se é POST
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "method not allowed"})
		return
	}

	// Decodifica o request
	var cepReq CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&cepReq); err != nil {
		span.RecordError(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid request format"})
		return
	}

	span.SetAttributes(semconv.HTTPMethod("POST"))
	span.SetAttributes(semconv.HTTPURL(r.URL.String()))

	// Valida o CEP
	if !validateCEP(cepReq.CEP) {
		span.SetAttributes(semconv.HTTPStatusCode(422))
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid zipcode"})
		return
	}

	// Encaminha para o Serviço B
	body, statusCode, err := forwardToServiceB(ctx, cepReq.CEP)
	if err != nil {
		span.RecordError(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "error forwarding to service B"})
		return
	}

	w.WriteHeader(statusCode)
	w.Write(body)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Inicializa o tracer
	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	http.HandleFunc("/cep", corsMiddleware(cepHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Service A starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
