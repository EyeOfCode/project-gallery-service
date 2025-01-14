FROM golang:1.23-alpine

WORKDIR /app

# Install tools
RUN apk add --no-cache git \
    && go install github.com/air-verse/air@latest \
    && go install github.com/swaggo/swag/cmd/swag@latest \
    && apk del git

# Copy dependency files first
COPY go.mod go.sum ./
RUN go mod download

RUN go mod tidy

# Copy source code
COPY . .

# Copy air config
COPY .air.toml* ./

# Generate swagger
RUN swag init -g cmd/api/main.go

EXPOSE 8080

CMD ["air"]