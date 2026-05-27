.PHONY: build run test lint integration load clean

# 1. Build the runtime container
build:
	docker compose build goboxd

# 2. Run the server
run:
	docker compose up goboxd

# 3. Run unit tests inside the isolated builder container
test:
	docker build --target builder -t goboxd-builder .
	docker run --rm goboxd-builder go test -v ./...

# 4. Run the Go linter to check for code quality
lint:
	docker build --target builder -t goboxd-builder .
	docker run --rm goboxd-builder golangci-lint run ./...

# 5. Integration tests (To be written)
integration:
	@echo "Integration tests will be run here"

# 6. Load tests (For Phase 3)
load:
	@echo "Load testing script will run here"

# Cleanup
clean:
	docker compose down -v