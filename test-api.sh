#!/bin/bash

# Script de teste para a API de temperatura por CEP
# Uso: ./test-api.sh

echo "🧪 Testando API de Temperatura por CEP"
echo "======================================"

# Função para testar uma requisição
test_request() {
    local description="$1"
    local cep="$2"
    local expected_status="$3"
    
    echo ""
    echo "📋 Teste: $description"
    echo "CEP: $cep"
    echo "Status esperado: $expected_status"
    
    response=$(curl -s -w "\n%{http_code}" -X POST http://localhost:8081/weather \
        -H "Content-Type: application/json" \
        -d "{\"cep\": \"$cep\"}")
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    echo "Status recebido: $http_code"
    echo "Resposta: $body"
    
    if [ "$http_code" = "$expected_status" ]; then
        echo "✅ SUCESSO"
    else
        echo "❌ FALHA - Status esperado: $expected_status, recebido: $http_code"
    fi
}

# Aguarda os serviços estarem prontos
echo "⏳ Aguardando serviços estarem prontos..."
sleep 5

# Testes
test_request "CEP válido (São Paulo)" "01310100" "200"
test_request "CEP válido (Rio de Janeiro)" "20040020" "200"
test_request "CEP com menos de 8 dígitos" "1234567" "422"
test_request "CEP com mais de 8 dígitos" "123456789" "422"
test_request "CEP inexistente" "99999999" "404"
test_request "CEP com letras" "abc12345" "422"

echo ""
echo "🎉 Testes concluídos!"
echo ""
echo "💡 Para visualizar os traces, acesse: http://localhost:9411" 