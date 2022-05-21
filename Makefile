.PHONY: setup-pre-commit-ci
setup-pre-commit-ci:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: setup
setup: setup-pre-commit-ci
	pre-commit install --hook-type commit-msg --hook-type pre-commit

.PHONY: test
test:
	go test ./... -covermode=count -coverprofile=coverage.out -count=1

.PHONY: coverage
coverage: test
	go tool cover -html=coverage.out

.PHONY: build
build:
	go build
