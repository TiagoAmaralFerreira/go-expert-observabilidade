package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
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

func GetTemperature(city string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return 0, errors.New("weather API key not configured")
	}

	normalizedCity := normalizeCityName(city)
	baseURL := "http://api.weatherapi.com/v1/current.json"
	params := url.Values{}
	params.Add("key", apiKey)
	params.Add("q", normalizedCity)
	params.Add("aqi", "no")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	var data WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.Current.TempC, nil
}
