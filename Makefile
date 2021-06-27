GO ?= go
BUILD ?= build
REPO ?= fjogeleit/policy-reporter-ui
IMAGE_TAG ?= 0.11.0
LD_FLAGS="-s -w"

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
	CGO_ENABLED=0 $(GO) build -v -ldflags="$(LD_FLAGS)" $(GOFLAGS) -o $(BUILD)/policyreporter-ui .

.PHONY: docker-build
docker-build:
	@docker buildx build --progress plane --platform linux/arm64,linux/amd64 --tag $(REPO):$(IMAGE_TAG) . --build-arg LD_FLAGS=$(LD_FLAGS)
	@docker buildx build --progress plane --platform linux/arm64,linux/amd64 --tag $(REPO):latest . --build-arg LD_FLAGS=$(LD_FLAGS)

.PHONY: docker-push
docker-push:
	@docker buildx build --progress plane --platform linux/arm64,linux/amd64 --tag $(REPO):$(IMAGE_TAG) . --build-arg LD_FLAGS=$(LD_FLAGS) --push