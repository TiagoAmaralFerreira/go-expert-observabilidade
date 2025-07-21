# 🚀 Guia de Início Rápido

## Pré-requisitos

1. **Docker e Docker Compose** instalados
2. **Chave da WeatherAPI** (gratuita em https://www.weatherapi.com/)

## ⚡ Execução Rápida

### 1. Configure as variáveis de ambiente

```bash
cp env.example .env
# Edite o arquivo .env e adicione sua chave da WeatherAPI
```

### 2. Execute o sistema

```bash
docker-compose up --build
```

### 3. Teste a API

```bash
# Teste com CEP válido
curl -X POST http://localhost:8081/weather \
  -H "Content-Type: application/json" \
  -d '{"cep": "01310100"}'

# Ou use o script de teste
./test-api.sh
```

### 4. Visualize os traces

Acesse: http://localhost:9411

## 📊 O que você verá

### Resposta da API
```json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

### Traces no Zipkin
- Tempo de resposta de cada operação
- Dependências entre serviços
- Detalhes das chamadas para APIs externas

## 🛠️ Comandos Úteis

```bash
# Ver logs
docker-compose logs service-a
docker-compose logs service-b

# Parar serviços
docker-compose down

# Reconstruir
docker-compose up --build

# Executar em background
docker-compose up -d
```

## 🔧 Troubleshooting

### Problema: "weather API key not configured"
**Solução**: Configure a variável `WEATHER_API_KEY` no arquivo `.env`

### Problema: "can not find zipcode"
**Solução**: Use um CEP válido (ex: "01310100" para São Paulo)

### Problema: Serviços não iniciam
**Solução**: Verifique se as portas 8081, 8082 e 9411 estão livres 