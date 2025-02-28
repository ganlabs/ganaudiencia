# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@if command -v rsrc > /dev/null; then \
            rsrc -ico static/ganlabs.ico; \
        else \
            read -p "Go's 'rsrc' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/akavel/rsrc@latest; \
                rsrc -ico static/ganlabs.ico; \
            else \
                echo "You chose not to install rsrc. Exiting..."; \
                exit 1; \
            fi; \
        fi
	@go build -o dist/ganaudiencia .
	@GOOS=windows go build -o dist/ganaudiencia.exe .
	@echo "Build complete"

# @GOOS=windows go build -ldflags="-H windowsgui" -o dist/ganaudiencia.exe .

# Run the application
run:
	@go run .

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v



# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload

watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

# Live Reload

dev: watch
	@echo "Watching..."

# CLI
dist: build
	@mkdir -p ~/rdpshare/gan
	@cp dist/ganaudiencia.exe ~/rdpshare/gan
	@cp dist/ganaudiencia.exe ~/winshare/gan



cli:
	@go run cmd/cli/main.go


# Docker
docker-build:
	@docker build -t ganaudiencia:latest -f Containerfile .

docker-up:
	@docker compose up -d

docker-up-dev:
	@docker compose up

docker-down:
	@docker compose down
	@docker image rm ganaudiencia-app:latest 


.PHONY: all build run test clean watch
