COMPILER := go
BINNAME := gomake

BUILDCMD := $(COMPILER) build
OUTPUT := -o $(BINNAME)
FLAGS := -v
VERSION := 1.0.6

RUNCMD := $(COMPILER) run

.PHONY: all build run clean win help release gh doc

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

release: build win ## Package the binary for release
	@if [ -f "$(BINNAME)" ] && [ -f "$(BINNAME).exe" ]; then \
		echo "Packaging $(BINNAME) for release"; \
		tar -czf "$(BINNAME)-$(VERSION).tar.gz" "$(BINNAME)" "$(BINNAME).exe"; \
	else \
	        echo "Error: Binary $(BINNAME) is missing for Linux or Windows."; \
	        echo "Try running this command:"; \
	        echo "\tmake"; \
	fi

gh: release ## Create a release on GitHub
	@echo "Creating release $(VERSION) on GitHub"
	@git tag -a v$(VERSION) -m "Version $(VERSION)"
	@git push origin v$(VERSION)
	@gh release create v$(VERSION) "$(BINNAME)-$(VERSION).tar.gz" --title "$(VERSION)" --notes "Release $(VERSION)"

doc: ## Create doc/scc.html
	@echo "Creating scc documentation in html"
	@mkdir -p "doc"
	@touch "doc/scc.html"
	@scc --overhead 1.0 --no-gen -n "scc.html" -s "complexity" -f "html" > doc/scc.html

help: ## Prints help for targets with comments
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
