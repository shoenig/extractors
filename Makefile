SHELL = sh

default: test

.PHONY: test
test:
	@echo "--> Running Tests ..."
	@go test -count=1 -v -race ./...

.PHONY: vet
vet:
	@echo "--> Vet Go Sources ..."
	@go vet ./...

.PHONY: lint
lint: vet
	@echo "--> Lint ..."
	@golangci-lint run --config .github/workflows/scripts/golangci.yaml
