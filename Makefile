CERT_DIR=cert
SERVER_CRT=$(CERT_DIR)/server.crt
SERVER_KEY=$(CERT_DIR)/server.key
CA_CRT=$(CERT_DIR)/ca.crt

.PHONY: build run test proto docker docker-up docker-down certs clean

# Go
build:
	go build -o trading-service ./cmd/server

run: build
	./trading-service

test:
	go test ./...

# Protobuf
proto:
	protoc \
	  --proto_path=internal/grpc/proto \
	  --go_out=paths=source_relative:internal/grpc/proto \
	  --go-grpc_out=paths=source_relative:internal/grpc/proto \
	  internal/grpc/proto/candlestick.proto


# Docker
docker:
	docker build -t trading-service .

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

# TLS Certificates (self-signed)
certs:
	mkdir -p $(CERT_DIR)
	openssl req -x509 -newkey rsa:4096 -sha256 -days 365 -nodes \
		-keyout $(SERVER_KEY) -out $(SERVER_CRT) \
		-subj "/C=US/ST=State/L=Local/O=Trading/OU=Service/CN=localhost" \
		-addext "subjectAltName = DNS:localhost"
	cp $(SERVER_CRT) $(CA_CRT)


clean:
	rm -rf trading-service cert