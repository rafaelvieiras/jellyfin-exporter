# Etapa de construção
FROM golang:1.18 AS build

# Defina a pasta de trabalho dentro do contêiner
WORKDIR /build

# Copie os arquivos go.mod e go.sum para baixar as dependências
COPY go.mod .
COPY go.sum .

# Baixe as dependências do projeto
RUN go mod download

# Copie o código fonte para o contêiner
COPY . .

# Construa a aplicação
RUN CGO_ENABLED=0 go build -o jellyfin-exporter ./main.go

# Etapa final para uma imagem mais leve
FROM alpine:latest

# Copie o binário da etapa de construção para a imagem final
COPY --from=build /build/jellyfin-exporter /app/jellyfin-exporter

# Defina a pasta de trabalho dentro do contêiner
WORKDIR /app

# Defina a porta que a aplicação vai expor
EXPOSE 8080

# Comando para iniciar a aplicação
ENTRYPOINT ["./jellyfin-exporter"]
