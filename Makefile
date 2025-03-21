.PHONY: clean run build test db-build db-run debug-env
.ONESHELL:

clean:
	@echo "Cleaning up..."
	rm -rf bin
	rm -rf vendor

debug-env:
	@echo "PostgreSQL User: zkp_user"
	@echo "PostgreSQL Password: zkp_password"
	@echo "PostgreSQL Database: zkp_db"

run:
	@echo "Starting services..."
	# Starting services in background
	go run -mod=vendor cmd/api-gateway/main.go & pid1=$$!; \

	# Kill all services when SIGINT is received
	trap "kill $$pid1 && exit 0" SIGINT; \

	# Wait for all services to finish
	wait $$pid1; \

build:
	go mod tidy
	go mod vendor
	go build -mod=vendor -o bin/api-gateway cmd/api-gateway/main.go
	go build -mod=vendor -o bin/zkp-service cmd/zkp-service/main.go
	go build -mod=vendor -o bin/auth-service cmd/auth-service/main.go
	go build -mod=vendor -o bin/messaging cmd/messaging-service/main.go
	go build -mod=vendor -o bin/contacts cmd/contacts-service/main.go
	go build -mod=vendor -o bin/encryption internal/encryption/encryption.go

test:
	cd internal/encryption && go test -mod=vendor
	cd ../..
	cd internal/zkp && go test -mod=vendor

db-build:
	@echo "Building Docker image for PostgreSQL..."
	docker build -t zkp-postgres .

db-run:
	@echo "Running PostgreSQL container with --env-file..."
	docker run -d \
		--name zkp-postgres \
		-p 5432:5432 \
		zkp-postgres