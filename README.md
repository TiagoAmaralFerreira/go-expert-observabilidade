# Sistema de Temperatura por CEP com OTEL

Sistema distribuído em Go que recebe um CEP, identifica a cidade e retorna o clima atual com implementação de OpenTelemetry e Zipkin para observabilidade.

## 🏗️ Arquitetura

O sistema é composto por dois serviços:

- **Serviço A**: Responsável pelo input e validação do CEP
- **Serviço B**: Responsável pela orquestração e busca de dados de temperatura

## 🚀 Funcionalidades

- Validação de CEP (8 dígitos)
- Busca de localização por CEP via ViaCEP
- Consulta de temperatura via WeatherAPI
- Conversão automática de temperaturas (Celsius, Fahrenheit, Kelvin)
- Tracing distribuído com OpenTelemetry
- Visualização de traces no Zipkin

## 📋 Pré-requisitos

- Docker e Docker Compose
- Conta na [WeatherAPI](https://www.weatherapi.com/) (gratuita)

## ⚙️ Configuração

1. **Clone o repositório:**
   ```bash
   git clone <seu-repositorio>
   cd observabilidade
   ```

2. **Configure as variáveis de ambiente:**
   Copie o arquivo de exemplo e configure:
   ```bash
   cp env.example .env
   ```
   
   Edite o arquivo `.env`:
   ```env
   WEATHER_API_KEY=sua_chave_aqui
   SERVICE_A_PORT=8080
   SERVICE_B_PORT=8081
   ZIPKIN_URL=http://localhost:9411
   OTEL_COLLECTOR_URL=http://localhost:4317
   ```

## 🏃‍♂️ Executando

### Com Docker Compose (Recomendado)
```bash
docker-compose up --build
```

### Desenvolvimento Local
```bash
# Terminal 1 - Serviço A
cd service-a
go run main.go

# Terminal 2 - Serviço B  
cd service-b
go run main.go

# Terminal 3 - Zipkin
docker run -d -p 9411:9411 openzipkin/zipkin
```

## 📡 Endpoints

### Serviço A - POST /cep
Recebe um CEP e encaminha para o Serviço B.

**Request Body:**
```json
{
  "cep": "29902555"
}
```

**Resposta de sucesso (200):**
```json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

**Resposta de erro (422):**
```json
{
  "message": "invalid zipcode"
}
```

### Serviço B - GET /weather/{cep}
Consulta direta de temperatura por CEP.

**Exemplo:**
```bash
curl http://localhost:8081/weather/29902555
```

## 🧪 Testes

```bash
# Teste do Serviço A
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "29902555"}'

# Teste do Serviço B
curl http://localhost:8081/weather/29902555
```

## 📊 Observabilidade

### Zipkin UI
Acesse http://localhost:9411 para visualizar os traces distribuídos.

### Métricas Coletadas
- Tempo de resposta do ViaCEP
- Tempo de resposta do WeatherAPI
- Traces completos da requisição

## 📁 Estrutura do Projeto

```
├── service-a/           # Serviço A - Input e validação
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── service-b/           # Serviço B - Orquestração
│   ├── main.go
│   ├── services/
│   ├── go.mod
│   └── Dockerfile
├── docker-compose.yml   # Orquestração dos serviços
├── otel-collector-config.yaml
└── README.md
```

## 🔧 Tecnologias

- **Go**: Linguagem principal
- **OpenTelemetry**: Observabilidade
- **Zipkin**: Visualização de traces
- **ViaCEP**: API de CEPs brasileiros
- **WeatherAPI**: API de clima
- **Docker**: Containerização

## 🧾 Evidências

<img width="1914" height="1029" alt="image" src="https://github.com/user-attachments/assets/72b15e57-216a-4df6-8c8a-93ba10068359" />

<img width="1565" height="643" alt="image" src="https://github.com/user-attachments/assets/13d06617-53de-4093-b685-ce00a24817e0" />
