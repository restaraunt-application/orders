# ==============================================================================
# Running linters
fmt: ## Format code
	gofumpt -l -w $$(go list -f {{.Dir}} ./... | grep -v /vendor/)

lint: ## Run lint go code
	golangci-lint run
# ==============================================================================
# Modules support
tidy:
	go mod tidy
	go mod vendor
# ==============================================================================
