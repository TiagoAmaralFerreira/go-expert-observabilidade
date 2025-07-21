package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var viaCEPURL = "https://viacep.com.br/ws/%s/json/"

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro"`
}

func GetCityByCEP(cep string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(viaCEPURL, cep))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if data.Erro {
		return "", errors.New("can not find zipcode")
	}

	return data.Localidade, nil
}
