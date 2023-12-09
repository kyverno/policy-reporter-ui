############
# DEFAULTS #
############

USE_CONFIG           ?= standard,no-ingress,in-cluster,all-read-rbac
KUBECONFIG           ?= ""
PIP                  ?= "pip3"
GO 					 ?= go
BUILD 				 ?= build
IMAGE_TAG 			 ?= 2.0.0-preview

#############
# VARIABLES #
#############

GIT_SHA             := $(shell git rev-parse HEAD)
KOCACHE             ?= /tmp/ko-cache
GOOS                ?= $(shell go env GOOS)
GOARCH              ?= $(shell go env GOARCH)
REGISTRY            ?= ghcr.io
OWNER               ?= kyverno
KO_REGISTRY         := ko.local
LD_FLAGS            := "-s -w"
LOCAL_PLATFORM      := linux/$(GOARCH)
PLATFORMS           := linux/arm64,linux/amd64,linux/s390x
KO_TAGS             := $(GIT_SHA)
IMAGE   			:= policy-reporter-ui
REPO                := $(REGISTRY)/$(OWNER)/$(IMAGE)
COMMA               := ,

ifndef VERSION
APP_VERSION         := $(GIT_SHA)
else
APP_VERSION         := $(VERSION)
endif


#########
# TOOLS #
#########
TOOLS_DIR      					   := $(PWD)/.tools
KO             					   := $(TOOLS_DIR)/ko
KO_VERSION     					   := main
GCI                                := $(TOOLS_DIR)/gci
GCI_VERSION                        := v0.9.1
GOFUMPT                            := $(TOOLS_DIR)/gofumpt
GOFUMPT_VERSION                    := v0.4.0

$(KO):
	@echo Install ko... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/google/ko@$(KO_VERSION)

$(GCI):
	@echo Install gci... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/daixiang0/gci@$(GCI_VERSION)

$(GOFUMPT):
	@echo Install gofumpt... >&2
	@GOBIN=$(TOOLS_DIR) go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)


.PHONY: gci
gci: $(GCI)
	@echo "Running gci"
	@$(GCI) write -s standard -s default -s "prefix(github.com/kyverno/policy-reporter-ui)" ./backend

.PHONY: gofumpt
gofumpt: $(GOFUMPT)
	@echo "Running gofumpt"
	@$(GOFUMPT) -w ./backend

.PHONY: fmt
fmt: gci gofumpt

.PHONY: build-frontend
build-frontend:
	@echo Build frontend with bun... >&2
	@cd frontend && bun install && bun run generate

.PHONY: ko-build
ko-build: $(KO)
	@echo Build image with ko... >&2
	@rm -rf backend/kodata
	@cp -r frontend/dist backend/kodata
	@cd backend && KO_DATA_DATE_EPOCH=$(git log -1 --format='%ct') LDFLAGS=$(LD_FLAGS) KOCACHE=$(KOCACHE) KO_DOCKER_REPO=$(KO_REGISTRY) \
		$(KO) build . --tags=$(KO_TAGS) --platform=$(LOCAL_PLATFORM)

.PHONY: ko-publish
ko-publish: $(KO)
	@echo Publishing image "$(KO_TAGS)" with ko... >&2
	@rm -rf backend/kodata
	@cp -r frontend/dist backend/kodata
	@cp -r backend/templates backend/kodata/email-templates
	@cd backend && KO_DATA_DATE_EPOCH=$(git log -1 --format='%ct') LDFLAGS=$(LD_FLAGS) KOCACHE=$(KOCACHE) KO_DOCKER_REPO=$(REPO) \
		$(KO) build . --bare --tags=$(KO_TAGS) --push --platform=$(PLATFORMS)
