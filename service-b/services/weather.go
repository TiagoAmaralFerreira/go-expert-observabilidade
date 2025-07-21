package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func normalizeCityName(city string) string {
	replacements := map[string]string{
		"ã": "a", "á": "a", "â": "a", "à": "a",
		"é": "e", "ê": "e", "è": "e",
		"í": "i", "î": "i", "ì": "i",
		"ó": "o", "ô": "o", "ò": "o",
		"ú": "u", "û": "u", "ù": "u",
		"ç": "c",
		"ñ": "n",
	}

	normalized := city
	for accented, plain := range replacements {
		normalized = strings.ReplaceAll(normalized, accented, plain)
		normalized = strings.ReplaceAll(normalized, strings.ToUpper(accented), strings.ToUpper(plain))
	}

	return normalized
}

func GetTemperature(ctx context.Context, city string) (float64, error) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		span.RecordError(errors.New("weather API key not configured"))
		return 0, errors.New("weather API key not configured")
	}

	normalizedCity := normalizeCityName(city)
	span.SetAttributes(attribute.String("city.original", city))
	span.SetAttributes(attribute.String("city.normalized", normalizedCity))

	baseURL := "http://api.weatherapi.com/v1/current.json"
	params := url.Values{}
	params.Add("key", apiKey)
	params.Add("q", normalizedCity)
	params.Add("aqi", "no")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	span.SetAttributes(attribute.String("weather.api.url", baseURL))

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		span.RecordError(err)
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		return 0, fmt.Errorf("failed to get weather info: %w", err)
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))

	if resp.StatusCode != http.StatusOK {
		span.RecordError(fmt.Errorf("weather API returned status %d", resp.StatusCode))
		return 0, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var data WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		span.RecordError(err)
		return 0, fmt.Errorf("failed to decode weather response: %w", err)
	}

	span.SetAttributes(attribute.Float64("temperature.celsius", data.Current.TempC))
	return data.Current.TempC, nil
}
