# 🚀 Guia Rápido - Sistema de Temperatura por CEP

Este guia mostra como executar rapidamente o sistema de temperatura por CEP com observabilidade.

## 📋 Pré-requisitos

- Docker e Docker Compose instalados
- Conta na [WeatherAPI](https://www.weatherapi.com/) (gratuita)

## ⚙️ Configuração Rápida

1. **Clone e configure:**
   ```bash
   git clone <seu-repositorio>
   cd observabilidade
   cp env.example .env
   ```

2. **Configure sua chave da WeatherAPI:**
   ```bash
   # Edite o arquivo .env e adicione sua chave
   nano .env
   ```

3. **Execute o sistema:**
   ```bash
   docker-compose up --build
   ```

## 🧪 Testando

Após os serviços estarem rodando, execute:

```bash
./test-api.sh
```

Ou teste manualmente:

```bash
# Teste do Serviço A
curl -X POST http://localhost:8080/cep \
  -H "Content-Type: application/json" \
  -d '{"cep": "29902555"}'

# Teste do Serviço B
curl http://localhost:8081/weather/29902555
```

## 📊 Observabilidade

- **Zipkin UI**: http://localhost:9411
- **Logs dos serviços**: `docker-compose logs -f`

## 🏗️ Arquitetura

```
┌─────────────┐    HTTP    ┌─────────────┐    HTTP    ┌─────────────┐
│   Cliente   │ ────────── │ Serviço A   │ ────────── │ Serviço B   │
│             │            │ (Porta 8080)│            │ (Porta 8081)│
└─────────────┘            └─────────────┘            └─────────────┘
                                   │                         │
                                   │ OTEL                    │ OTEL
                                   ▼                         ▼
                            ┌─────────────┐            ┌─────────────┐
                            │   Zipkin    │            │ ViaCEP API  │
                            │ (Porta 9411)│            │             │
                            └─────────────┘            └─────────────┘
                                                              │
                                                              │ HTTP
                                                              ▼
                                                       ┌─────────────┐
                                                       │WeatherAPI   │
                                                       │             │
                                                       └─────────────┘
```

## 🔧 Serviços

### Serviço A (Porta 8080)
- **Endpoint**: `POST /cep`
- **Função**: Validação de CEP e encaminhamento para Serviço B
- **Validações**: 8 dígitos, formato string

### Serviço B (Porta 8081)
- **Endpoints**: 
  - `GET /weather/{cep}`
  - `POST /weather`
- **Função**: Busca cidade por CEP e temperatura
- **APIs**: ViaCEP + WeatherAPI

## 📝 Respostas

### Sucesso (200)
```json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

### Erro - CEP Inválido (422)
```json
{
  "message": "invalid zipcode"
}
```

### Erro - CEP Não Encontrado (404)
```json
{
  "message": "can not find zipcode"
}
```

## 🛠️ Desenvolvimento Local

Para desenvolvimento sem Docker:

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

## 🐛 Troubleshooting

1. **Erro de conexão**: Verifique se as portas estão livres
2. **Erro de API**: Configure corretamente a chave da WeatherAPI
3. **Traces não aparecem**: Aguarde alguns segundos para o Zipkin processar

## 📚 Documentação Completa

Veja o [README.md](README.md) para documentação detalhada. 