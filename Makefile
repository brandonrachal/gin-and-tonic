.PHONY: all build clean test help vet

# Define the output directory for binaries
BIN_DIR := bin

# List of command directories to build
COMMANDS := migration_client api_server

all: build

build: $(addprefix $(BIN_DIR)/, $(COMMANDS))

$(BIN_DIR)/%: cmd/%/main.go
	@mkdir -p $(BIN_DIR)
	@echo "Building $@"
	go build -o $@ ./cmd/$(notdir $*)/main.go

clean:
	@echo "Cleaning up $(BIN_DIR)"
	@rm -rf $(BIN_DIR)

test:
	go test ./...

vet:
	go vet ./...

help:
	@echo "Available targets:"
	@echo "  test: Runs all Go tests in the project"
	@echo "  help: Displays this help message"