package services

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCityByCEP(t *testing.T) {
	// Cria um servidor de teste que simula a API ViaCEP
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ws/01001000/json/" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"localidade":"São Paulo","erro":false}`))
		} else if r.URL.Path == "/ws/00000000/json/" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"erro":true}`))
		} else if r.URL.Path == "/ws/99999999/json/" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"erro":true}`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer testServer.Close()

	// Substitui a URL base temporariamente para os testes
	originalURL := viaCEPURL
	viaCEPURL = testServer.URL + "/ws/%s/json/"
	defer func() { viaCEPURL = originalURL }() // Restaura após os testes

	tests := []struct {
		name    string
		cep     string
		want    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "CEP válido",
			cep:     "01001000",
			want:    "São Paulo",
			wantErr: false,
		},
		{
			name:    "CEP inválido",
			cep:     "00000000",
			want:    "",
			wantErr: true,
			errMsg:  "can not find zipcode",
		},
		{
			name:    "CEP não encontrado",
			cep:     "99999999",
			want:    "",
			wantErr: true,
			errMsg:  "can not find zipcode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCityByCEP(tt.cep)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetCityByCEP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("GetCityByCEP() error message = %v, want %v", err.Error(), tt.errMsg)
			}

			if got != tt.want {
				t.Errorf("GetCityByCEP() = %v, want %v", got, tt.want)
			}
		})
	}
}
