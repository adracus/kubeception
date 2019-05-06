.PHONY: ci
ci: install-requirements verify

.PHONY: all
all: ci install

.PHONY: install-requirements
install-requirements:
	@cd tools && GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@692dacb773b703162c091c2d8c59f9cd2d6801db; cd -
	@cd tools && GO111MODULE=on go get github.com/golang/mock/mockgen@v1.3.0; cd -

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

.PHONY: verify
verify: check test

.PHONY: start
start:
	@go run cmd/manager/main.go

.PHONY: install
install:
	@go install cmd/manager/main.go

