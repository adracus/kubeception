.PHONY: all
all: check test install

.PHONY: install-requirements
install-requirements:
	@cd $$(mktemp -d); GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.16.0; cd -

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
	@golangci-lint run

.PHONY: check-format
check-format:
	@unformatted="$$(go fmt -n ./pkg/...)"
	@if [ "$$unformatted" != "" ]; then echo "Unformatted files detected: $unformatted"; exit 1; fi

.PHONY: check
check: check-format lint

.PHONY: start
start:
	@go run cmd/manager/main.go

.PHONY: install
install:
	@go install cmd/manager/main.go

