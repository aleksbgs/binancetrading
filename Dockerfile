# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install OS dependencies
RUN apk add --no-cache git make protobuf openssl

# Install Go CLI tools
RUN go install github.com/swaggo/swag/cmd/swag@latest \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH="/go/bin:$PATH"

# Copy source code
COPY . .

# Generate Swagger docs
RUN swag init --generalInfo cmd/server/main.go --output docs

# Compile protobufs
RUN protoc \
  --proto_path=internal/grpc/proto \
  --go_out=paths=source_relative:internal/grpc/proto \
  --go-grpc_out=paths=source_relative:internal/grpc/proto \
  internal/grpc/proto/candlestick.proto

# Build binary
RUN go mod tidy && go build -o trading-service ./cmd/server

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /root/

# Copy built binary and resources
COPY --from=builder /app/trading-service .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/internal/grpc/proto ./internal/grpc/proto
COPY --from=builder /app/cert ./cert
COPY --from=builder /app/.env .
COPY --from=builder /app/config.yaml .

EXPOSE 8080 50051

CMD ["./trading-service"]