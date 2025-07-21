# Weather API

API para consultar temperatura baseada em CEP brasileiro.

## 🚀 Funcionalidades

- Consulta de temperatura por CEP
- Suporte a múltiplas unidades de temperatura (Celsius, Fahrenheit, Kelvin)
- Validação de CEP
- Tratamento de erros robusto
- CORS habilitado

## 📋 Pré-requisitos

- Go 1.22+
- Conta na [WeatherAPI](https://www.weatherapi.com/) (gratuita)

## ⚙️ Configuração

1. **Clone o repositório:**
   ```bash
   git clone <seu-repositorio>
   cd cloud-run
   ```

2. **Configure as variáveis de ambiente:**
   Crie um arquivo `.env` na raiz do projeto:
   ```env
   WEATHER_API_KEY=sua_chave_aqui
   PORT=8082
   ```

3. **Instale as dependências:**
   ```bash
   go mod tidy
   ```

## 🏃‍♂️ Executando

### Desenvolvimento
```bash
go run main.go
```

### Produção
```bash
go build -o main .
./main
```

### Docker
```bash
docker build -t weather-api .
docker run -p 8080:8080 weather-api
```

## 📡 Endpoints

### GET /weather/{CEP}

Consulta a temperatura para um CEP específico.

**Parâmetros:**
- `CEP`: CEP brasileiro (8 dígitos, sem hífen)

**Exemplo de requisição:**
```bash
curl http://localhost:8082/weather/08141140
```

**Resposta de sucesso:**
```json
{
  "temp_C": 25.5,
  "temp_F": 77.9,
  "temp_K": 298.65
}
```

**Respostas de erro:**
```json
{
  "message": "invalid zipcode"
}
```
```json
{
  "message": "can not find zipcode"
}
```

## 🧪 Testes

```bash
go test ./services -v
```

## 🚀 Deploy no Cloud Run

```bash
gcloud run deploy weather-api --source . --region us-central1 --allow-unauthenticated
```

## 📁 Estrutura do Projeto

```
├── handlers/          # Handlers HTTP
│   └── weather.go
├── services/          # Lógica de negócio
│   ├── cep.go
│   ├── weather.go
│   └── cep_test.go
├── main.go           # Ponto de entrada
├── Dockerfile        # Configuração Docker
├── go.mod           # Dependências Go
└── .env             # Variáveis de ambiente
```

## 🔧 Tecnologias

- **Go**: Linguagem principal
- **WeatherAPI**: API de clima
- **ViaCEP**: API de CEPs brasileiros
- **Docker**: Containerização
- **Cloud Run**: Plataforma de deploy

## 📝 Licença

Este projeto está sob a licença MIT. 