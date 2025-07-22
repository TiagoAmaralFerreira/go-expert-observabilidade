#!/bin/bash

# Script para testar a API de temperatura por CEP
# Certifique-se de que os serviços estão rodando antes de executar este script

echo "🧪 Testando Sistema de Temperatura por CEP"
echo "=========================================="

# Configurações
SERVICE_A_URL="http://localhost:8080"
SERVICE_B_URL="http://localhost:8081"
TEST_CEP="29902555"

echo ""
echo "📡 Testando Serviço A (POST /cep)"
echo "----------------------------------"

# Teste 1: CEP válido
echo "✅ Teste 1: CEP válido ($TEST_CEP)"
response=$(curl -s -X POST "$SERVICE_A_URL/cep" \
  -H "Content-Type: application/json" \
  -d "{\"cep\": \"$TEST_CEP\"}")
echo "Resposta: $response"
echo ""

# Teste 2: CEP inválido (formato)
echo "❌ Teste 2: CEP inválido (formato)"
response=$(curl -s -X POST "$SERVICE_A_URL/cep" \
  -H "Content-Type: application/json" \
  -d '{"cep": "12345"}')
echo "Resposta: $response"
echo ""

# Teste 3: CEP inválido (tipo)
echo "❌ Teste 3: CEP inválido (tipo)"
response=$(curl -s -X POST "$SERVICE_A_URL/cep" \
  -H "Content-Type: application/json" \
  -d '{"cep": 12345678}')
echo "Resposta: $response"
echo ""

echo "📡 Testando Serviço B (GET /weather/{cep})"
echo "-------------------------------------------"

# Teste 4: CEP válido via GET
echo "✅ Teste 4: CEP válido via GET ($TEST_CEP)"
response=$(curl -s "$SERVICE_B_URL/weather/$TEST_CEP")
echo "Resposta: $response"
echo ""

# Teste 5: CEP inválido via GET
echo "❌ Teste 5: CEP inválido via GET"
response=$(curl -s "$SERVICE_B_URL/weather/12345")
echo "Resposta: $response"
echo ""

# Teste 6: CEP inexistente
echo "❌ Teste 6: CEP inexistente"
response=$(curl -s "$SERVICE_B_URL/weather/99999999")
echo "Resposta: $response"
echo ""

echo "📡 Testando Serviço B (POST /weather)"
echo "-------------------------------------"

# Teste 7: CEP válido via POST
echo "✅ Teste 7: CEP válido via POST ($TEST_CEP)"
response=$(curl -s -X POST "$SERVICE_B_URL/weather" \
  -H "Content-Type: application/json" \
  -d "{\"cep\": \"$TEST_CEP\"}")
echo "Resposta: $response"
echo ""

echo "🎯 Testes concluídos!"
echo ""
echo "📊 Para visualizar os traces, acesse: http://localhost:9411"
echo "🔍 Para verificar os logs dos serviços, use: docker-compose logs -f" 