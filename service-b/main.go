package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"service-b/services"

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

type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
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
		semconv.ServiceName("service-b"),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("service-b")
	return tp, nil
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func validateCEP(cep string) bool {
	return len(cep) == 8
}

func processCEP(ctx context.Context, cep string) (*WeatherResponse, int, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	// Valida o CEP
	if !validateCEP(cep) {
		span.SetAttributes(semconv.HTTPStatusCode(422))
		return nil, 422, fmt.Errorf("invalid zipcode")
	}

	// Busca a cidade pelo CEP
	city, err := services.GetCityByCEP(ctx, cep)
	if err != nil {
		if err.Error() == "can not find zipcode" {
			span.SetAttributes(semconv.HTTPStatusCode(404))
			return nil, 404, fmt.Errorf("can not find zipcode")
		}
		span.RecordError(err)
		return nil, 500, fmt.Errorf("error getting city information")
	}

	// Busca a temperatura
	tempC, err := services.GetTemperature(ctx, city)
	if err != nil {
		span.RecordError(err)
		return nil, 500, fmt.Errorf("error getting weather information")
	}

	response := &WeatherResponse{
		City:  city,
		TempC: tempC,
		TempF: tempC*1.8 + 32,
		TempK: tempC + 273,
	}

	span.SetAttributes(semconv.HTTPStatusCode(200))
	return response, 200, nil
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		// Extrai o CEP da URL /weather/{cep}
		path := strings.TrimPrefix(r.URL.Path, "/weather/")
		if path == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "CEP parameter required"})
			return
		}

		response, statusCode, err := processCEP(ctx, path)
		if err != nil {
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
			return
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)

	case "POST":
		// Decodifica o request
		var cepReq CEPRequest
		if err := json.NewDecoder(r.Body).Decode(&cepReq); err != nil {
			span.RecordError(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid request format"})
			return
		}

		response, statusCode, err := processCEP(ctx, cepReq.CEP)
		if err != nil {
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
			return
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "method not allowed"})
	}
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

	http.HandleFunc("/weather", corsMiddleware(weatherHandler))
	http.HandleFunc("/weather/", corsMiddleware(weatherHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Service B starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
