package templat

var LibMake = `VERSION := 0.0.1

.PHONY: release help doc

release: doc ## Create a release on GitHub
	@echo "Creating release $(VERSION) on GitHub"
	@git tag -a v$(VERSION) -m "Version $(VERSION)"
	@git push origin v$(VERSION)
	@gh release create v$(VERSION) --generate-notes --notes-from-tag --notes "Release $(VERSION), view changelogs in CHANGELOG.md"

docs: ## Create docs/scc.html
	@if ! [ -x "$(shell which scc)" ]; then \
		echo "scc is not installed"; \
		echo "installing scc..."; \
		go install github.com/boyter/scc/v3@latest; \
	fi
	@echo "Creating scc documentation in html"
	@mkdir -p "docs"
	@touch "docs/scc.html"
	@scc --overhead 1.0 --no-gen -n "scc.html" -s "complexity" -f "html" > docs/scc.html

help: ## Prints help for targets with comments
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
`
