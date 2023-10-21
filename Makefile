REPO ?= ghcr.io/kyverno/policy-reporter-ui
IMAGE_TAG ?= 1.9.1
LD_FLAGS="-s -w"
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x

.PHONY: docker-build
docker-build:
	@docker buildx build --progress plain --platform $(PLATFORMS) --tag $(REPO):$(IMAGE_TAG) . --build-arg LD_FLAGS=$(LD_FLAGS)
	@docker buildx build --progress plain --platform $(PLATFORMS) --tag $(REPO):latest . --build-arg LD_FLAGS=$(LD_FLAGS)

.PHONY: docker-push
docker-push:
	@docker buildx build --progress plain --platform $(PLATFORMS) --tag $(REPO):$(IMAGE_TAG) --tag $(REPO):latest . --build-arg LD_FLAGS=$(LD_FLAGS) --push

.PHONY: docker-push-dev
docker-push-dev:
	@docker buildx build --progress plain --platform $(PLATFORMS) --tag $(REPO):dev . --build-arg LD_FLAGS=$(LD_FLAGS) --push
