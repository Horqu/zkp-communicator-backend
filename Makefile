.PHONY: run build test

run:
	go run cmd/api-gateway/main.go

build:
	go get github.com/gin-gonic/gin
	go get go.mau.fi/libsignal/ecc
	go get golang.org/x/crypto/sha3@v0.25.0
	go get github.com/go-playground/validator/v10@v10.20.0
	go get github.com/gin-gonic/gin/binding@v1.10.0
	go get github.com/mattn/go-isatty@v0.0.20
	go get golang.org/x/net/idna@v0.25.0
	go get github.com/stretchr/testify
	go build -o bin/api-gateway cmd/api-gateway/main.go
	go build -o bin/encryption internal/encryption/encryption.go

test:
	cd internal/encryption && go test