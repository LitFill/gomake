package templat

var MakeTempl = `COMPILER := go
BINNAME := {{.ProgName}}

BUILDCMD := $(COMPILER) build
OUTPUT := -o $(BINNAME)
FLAGS := -v
VERSION := 0.0.1

RUNCMD := $(COMPILER) run

.PHONY: all build win run clean package release doc changelog gh help 

all: build win ## Build the binary for Linux and Windows

build: main.go ## Build the binary for Linux
	@echo "Building $(BINNAME) for Linux"
	@$(BUILDCMD) $(OUTPUT) $(FLAGS)

win: main.go ## Build the binary for a niche gaming os (Windows)
	@echo "Building $(BINNAME) for Windows"
	@$(BUILDCMD) $(OUTPUT).exe $(FLAGS)

run: main.go ## Run the main.go
	@echo "Running $(BINNAME)"
	@$(RUNCMD) $(FLAGS) .

clean: ## Clean up
	@echo "Cleaning up"
	@rm -f $(BINNAME)*

package: all ## Package the binary for release
	@echo "Packaging $(BINNAME) for release"
	@tar -czf "$(BINNAME)-$(VERSION).tar.gz" "$(BINNAME)" "$(BINNAME).exe"

release: package doc ## Create a release on GitHub
	@echo "Creating release $(VERSION) on GitHub"
	@git tag -a v$(VERSION) -m "Version $(VERSION)"
	@git push origin v$(VERSION)
	@gh release create v$(VERSION) "$(BINNAME)-$(VERSION).tar.gz" --generate-notes --notes-from-tag --title "$(VERSION)" --notes "Release $(VERSION), view changelogs in CHANGELOG.md"

doc: ## Create doc/scc.html
	@echo "Creating scc documentation in html"
	@mkdir -p "doc"
	@touch "doc/scc.html"
	@scc --overhead 1.0 --no-gen -n "scc.html" -s "complexity" -f "html" > doc/scc.html

changelog: ## Generate CHANGELOG.md
  @echo "Generating CHANGELOG.md"
  @git-chglog > CHANGELOG.md

help: ## Prints help for targets with comments
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
`
