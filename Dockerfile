# Build stage
FROM golang:1.22.5-alpine AS builder

# Instala dependências necessárias
RUN apk add --no-cache git

WORKDIR /app

# Copia os arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest

# Instala ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copia o binário compilado
COPY --from=builder /app/main .

# Expõe a porta
EXPOSE 8080

# Define a variável de ambiente padrão
ENV PORT=8080

# Executa a aplicação
CMD ["./main"]