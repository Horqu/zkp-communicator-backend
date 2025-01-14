.PHONY: run build test
.ONESHELL:

run:
	@echo "Starting services..."
	# Starting services in background
	go run cmd/api-gateway/main.go & pid1=$$!; \
	go run cmd/zkp-service/main.go & pid2=$$!; \
	go run cmd/auth-service/main.go & pid3=$$!; \

	# Kill all services when SIGINT is received
	trap "kill $$pid1; kill $$pid2; kill $$pid3 && exit 0" SIGINT; \

	# Wait for all services to finish
	wait $$pid1; wait $$pid2; wait $$pid3; \

build:
	go mod tidy
	go get github.com/gin-gonic/gin
	go get go.mau.fi/libsignal/ecc
	go get golang.org/x/crypto/sha3@v0.25.0
	go get github.com/go-playground/validator/v10@v10.20.0
	go get github.com/gin-gonic/gin/binding@v1.10.0
	go get github.com/mattn/go-isatty@v0.0.20
	go get golang.org/x/net/idna@v0.25.0
	go get github.com/stretchr/testify
	go build -o bin/api-gateway cmd/api-gateway/main.go
	go build -o bin/zkp-service cmd/zkp-service/main.go
	go build -o bin/auth-service cmd/auth-service/main.go
	go build -o bin/encryption internal/encryption/encryption.go

test:
	cd internal/encryption && go test