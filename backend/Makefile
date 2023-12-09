GO ?= go
BUILD ?= build
LD_FLAGS="-s -w"
GOFLAGS ?=

all: build

.PHONY: clean
clean:
	rm -rf $(BUILD)

.PHONY: prepare
prepare:
	mkdir -p $(BUILD)

.PHONY: test
test:
	go test -v ./... -timeout=120s

.PHONY: coverage
coverage:
	go test -v ./... -covermode=count -coverprofile=coverage.out -timeout=120s

.PHONY: build
build: prepare
	CGO_ENABLED=0 $(GO) build -v -ldflags=$(LD_FLAGS) $(GOFLAGS) -o $(BUILD)/policyreporter-ui .

.PHONY: fmt
fmt:
	$(call print-target)
	@echo "Running gci"
	@go run github.com/daixiang0/gci@v0.9.1 write -s standard -s default -s "prefix(github.com/kyverno/policy-reporter-ui)" .
	@echo "Running gofumpt"
	@go run mvdan.cc/gofumpt@v0.4.0 -w .
