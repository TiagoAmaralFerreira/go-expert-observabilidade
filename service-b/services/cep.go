package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var viaCEPURL = "https://viacep.com.br/ws/%s/json/"

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro"`
}

func GetCityByCEP(ctx context.Context, cep string) (string, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetAttributes(attribute.String("cep", cep))
	span.SetAttributes(attribute.String("viacep.url", viaCEPURL))

	url := fmt.Sprintf(viaCEPURL, cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		span.RecordError(err)
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		return "", fmt.Errorf("failed to get CEP info: %w", err)
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))

	if resp.StatusCode != http.StatusOK {
		span.RecordError(fmt.Errorf("viacep returned status %d", resp.StatusCode))
		return "", fmt.Errorf("viacep returned status %d", resp.StatusCode)
	}

	var data ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		span.RecordError(err)
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if data.Erro {
		span.SetAttributes(attribute.Bool("cep.not_found", true))
		return "", errors.New("can not find zipcode")
	}

	span.SetAttributes(attribute.String("city", data.Localidade))
	return data.Localidade, nil
}
