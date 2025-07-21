# Sistema de Temperatura por CEP com OTEL e Zipkin

Sistema distribuído em Go que implementa OpenTelemetry (OTEL) e Zipkin para observabilidade, permitindo consultar a temperatura atual de uma localização através de um CEP brasileiro.

## 🏗️ Arquitetura

O sistema é composto por dois serviços:

### Serviço A (Porta 8081)
- **Responsabilidade**: Validação de entrada e orquestração
- **Endpoint**: `POST /weather`
- **Input**: `{"cep": "29902555"}`
- **Validações**:
  - CEP deve conter exatamente 8 dígitos
  - CEP deve ser uma string
- **Respostas**:
  - `422`: CEP inválido
  - `200`: Sucesso (encaminha para Serviço B)

### Serviço B (Porta 8082)
- **Responsabilidade**: Orquestração de APIs externas
- **Endpoint**: `POST /weather`
- **Funcionalidades**:
  - Busca cidade pelo CEP (ViaCEP API)
  - Busca temperatura atual (WeatherAPI)
  - Conversão para Celsius, Fahrenheit e Kelvin
- **Respostas**:
  - `200`: `{"city": "São Paulo", "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.65}`
  - `422`: CEP inválido
  - `404`: CEP não encontrado

## 🚀 Tecnologias

- **Go 1.22+**: Linguagem principal
- **OpenTelemetry**: Instrumentação e observabilidade
- **Zipkin**: Backend de tracing distribuído
- **Docker & Docker Compose**: Containerização
- **ViaCEP API**: Consulta de CEPs brasileiros
- **WeatherAPI**: Dados meteorológicos

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
   ```bash
   cp env.example .env
   # Edite o arquivo .env e adicione sua chave da WeatherAPI
   ```

3. **Execute o sistema:**
   ```bash
   docker-compose up --build
   ```

## 🏃‍♂️ Executando

### Desenvolvimento Local

Para executar os serviços individualmente:

```bash
# Serviço A
cd service-a
go mod tidy
go run main.go

# Serviço B
cd service-b
go mod tidy
go run main.go
```

### Produção com Docker

```bash
# Construir e executar todos os serviços
docker-compose up --build

# Executar em background
docker-compose up -d --build

# Parar os serviços
docker-compose down
```

## 📡 Testando a API

### Exemplo de Requisição

```bash
curl -X POST http://localhost:8081/weather \
  -H "Content-Type: application/json" \
  -d '{"cep": "29902555"}'
```

### Resposta de Sucesso

```json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

### Respostas de Erro

```json
// CEP inválido (422)
{
  "message": "invalid zipcode"
}

// CEP não encontrado (404)
{
  "message": "can not find zipcode"
}
```

## 🔍 Observabilidade com Zipkin

### Acessando o Zipkin

1. Abra o navegador em: `http://localhost:9411`
2. Clique em "Find traces"
3. Selecione um serviço e clique em "Find traces"

### Spans Implementados

- **Serviço A**:
  - Validação de CEP
  - Encaminhamento para Serviço B

- **Serviço B**:
  - Busca de cidade por CEP (ViaCEP)
  - Busca de temperatura (WeatherAPI)
  - Conversão de unidades

### Tracing Distribuído

O sistema implementa tracing distribuído entre os serviços, permitindo visualizar:
- Tempo de resposta de cada operação
- Dependências entre serviços
- Erros e exceções
- Atributos customizados (CEP, cidade, temperatura)

## 📁 Estrutura do Projeto

```
├── service-a/                 # Serviço A - Validação
│   ├── main.go               # Ponto de entrada
│   ├── Dockerfile            # Containerização
│   └── go.mod               # Dependências
├── service-b/                 # Serviço B - Orquestração
│   ├── main.go               # Ponto de entrada
│   ├── services/             # Lógica de negócio
│   │   ├── cep.go           # Serviço ViaCEP
│   │   └── weather.go       # Serviço WeatherAPI
│   ├── Dockerfile            # Containerização
│   └── go.mod               # Dependências
├── docker-compose.yml         # Orquestração
├── env.example               # Variáveis de ambiente
└── README-OTEL.md           # Documentação
```

## 🧪 Testes

### Testando CEPs Válidos

```bash
# CEP válido
curl -X POST http://localhost:8081/weather \
  -H "Content-Type: application/json" \
  -d '{"cep": "29902555"}'

# CEP válido (São Paulo)
curl -X POST http://localhost:8081/weather \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310100"}'
```

### Testando CEPs Inválidos

```bash
# CEP com menos de 8 dígitos
curl -X POST http://localhost:8081/weather \
  -H "Content-Type: application/json" \
  -d '{"cep": "1234567"}'

# CEP com mais de 8 dígitos
curl -X POST http://localhost:8081/weather \
  -H "Content-Type: application/json" \
  -d '{"cep": "123456789"}'

# CEP inexistente
curl -X POST http://localhost:8081/weather \
  -H "Content-Type: application/json" \
  -d '{"cep": "99999999"}'
```

## 🔧 Monitoramento

### Logs dos Serviços

```bash
# Ver logs do Serviço A
docker-compose logs service-a

# Ver logs do Serviço B
docker-compose logs service-b

# Ver logs do Zipkin
docker-compose logs zipkin
```

### Métricas de Performance

No Zipkin, você pode visualizar:
- Tempo de resposta por operação
- Taxa de erro
- Dependências entre serviços
- Latência de rede

## 🚀 Deploy

### Ambiente de Produção

Para deploy em produção, considere:

1. **Configurar um backend persistente para o Zipkin** (MySQL, Elasticsearch)
2. **Usar um collector OTEL** para processar traces
3. **Implementar health checks** nos serviços
4. **Configurar métricas** (Prometheus + Grafana)
5. **Implementar rate limiting** e autenticação

### Exemplo de Deploy no Kubernetes

```yaml
# deployment-service-a.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: service-a
spec:
  replicas: 3
  selector:
    matchLabels:
      app: service-a
  template:
    metadata:
      labels:
        app: service-a
    spec:
      containers:
      - name: service-a
        image: service-a:latest
        ports:
        - containerPort: 8081
        env:
        - name: SERVICE_B_URL
          value: "http://service-b:8082"
```

## 📝 Licença

Este projeto está sob a licença MIT. 