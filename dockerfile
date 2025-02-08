# Etapa de construcción
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copiar mod y sum para descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o /task-manager ./cmd/main.go

# Etapa final
FROM alpine:3.18

WORKDIR /app

# Copiar el binario compilado
COPY --from=builder /task-manager /app/task-manager
COPY web/templates /app/web/templates
COPY web/static /app/web/static

# Puerto expuesto
EXPOSE 8080

# Comando de ejecución
CMD ["/app/task-manager"]