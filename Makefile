# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@go build -o ganaudiencia .
	@GOOS=windows go build -o ganaudiencia.exe .


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


# CLI
dist:
	@go build -o build/cli/linux/main cmd/cli/main.go
	@cp -rf driver/chromedriver-linux64 build/cli/linux
	@cp -rf driver/geckodriver-linux64 build/cli/linux
	
	@GOOS=windows go build -o build/cli/win/main.exe cmd/cli/main.go
	@mkdir -p ~/rdpshare/gondwana
	@cp build/cli/win/main.exe ~/rdpshare/gondwana
	@cp -rf driver/chromedriver-win64 build/cli/win
	@cp -rf driver/chromedriver-win64 ~/rdpshare/gondwana


cli:
	@go run cmd/cli/main.go


# Docker
docker-build:
	@docker build -t ganaudiencia:latest -f Containerfile .

docker-up:
	@docker compose up -d

docker-down:
	@docker compose down


.PHONY: all build run test clean watch
