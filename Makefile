.PHONY: all
all: verify generate install

.PHONY: generate
generate:
	@go generate ./pkg/...

.PHONY: format
format:
	@go fmt ./pkg/...

.PHONY: test
test:
	@go test ./pkg/...

.PHONY: lint
lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run

.PHONY: check-format
check-format:
	@unformatted="$$(go fmt -n ./pkg/...)"
	@if [ "$$unformatted" != "" ]; then echo "Unformatted files detected: $unformatted"; exit 1; fi

.PHONY: check
check: check-format lint

.PHONY: verify
verify: check test

.PHONY: start
start:
	@go run cmd/manager/main.go

.PHONY: install
install:
	@go install cmd/manager/main.go

