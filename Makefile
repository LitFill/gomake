COMPILER := go
BINNAME := gomake

BUILDCMD := $(COMPILER) build
OUTPUT := -o $(BINNAME)
FLAGS := -v

RUNCMD := $(COMPILER) run

.PHONY: all build run clean win help

all: build ## Build the binary for Linux

build: main.go ## Actually build the binary
	@echo "Building $(BINNAME) for Linux"
	@$(BUILDCMD) $(OUTPUT) $(FLAGS)

win: main.go ## Build the binary for Windows
	@echo "Building $(BINNAME) for Windows"
	@$(BUILDCMD) $(OUTPUT).exe $(FLAGS)

run: main.go ## Run the main.go
	@echo "Running $(BINNAME)"
	@$(RUNCMD) $(FLAGS) $^

clean: ## Clean up
	@echo "Cleaning up"
	@rm -f $(BINNAME)*

help: ## Prints help for targets with comments
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
