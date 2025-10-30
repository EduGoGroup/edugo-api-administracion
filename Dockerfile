# Dockerfile para API Administración
# Puerto: 8081

FROM golang:1.21-alpine AS builder

# Instalar dependencias del sistema
RUN apk add --no-cache git

# Establecer directorio de trabajo
WORKDIR /app

# Copiar go.mod y go.sum
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Etapa final - imagen ligera
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar binario compilado desde builder
COPY --from=builder /app/main .

# Exponer puerto
EXPOSE 8081

# Comando de inicio
CMD ["./main"]
