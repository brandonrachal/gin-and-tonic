.PHONY: all build clean test vet help

BIN_DIR := bin

COMMANDS := migration_client api_server

all: vet test clean build

build: $(addprefix $(BIN_DIR)/, $(COMMANDS))

$(BIN_DIR)/%: cmd/%/main.go
	@mkdir -p $(BIN_DIR)
	@echo "Building $@"
	go build -o $@ ./cmd/$(notdir $*)/main.go

clean:
	@echo "Cleaning up $(BIN_DIR)"
	@find $(BIN_DIR) -type f ! -name ".gitkeep" -delete

test:
	go test ./...

vet:
	go vet ./...

help:
	@echo "Available targets:"
	@echo "  all: Vets the code, tests it, cleans the bin, builds the programs"
	@echo "  build: Builds all programs"
	@echo "  clean: Removes all programs from bin folder"
	@echo "  vet: Vet all Go code in the project"
	@echo "  test: Runs all Go tests in the project"
	@echo "  help: Displays this help message"