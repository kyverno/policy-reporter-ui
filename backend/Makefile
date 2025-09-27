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
PLATFORMS           := all
IMAGE   			:= policy-reporter-ui
REPO                := $(REGISTRY)/$(OWNER)/$(IMAGE)
KO_TAGS             := $(shell git rev-parse --short HEAD)
COMMA               := ,
PACKAGE             ?= github.com/kyverno/policy-reporter-ui

ifndef VERSION
KO_TAGS         := $(shell git rev-parse --short HEAD)
else
KO_TAGS         := $(VERSION)
endif


#########
# TOOLS #
#########

TOOLS_DIR                          := $(PWD)/.tools
CONTROLLER_GEN                     := $(TOOLS_DIR)/controller-gen
CONTROLLER_GEN_VERSION             ?= v0.16.1
CLIENT_GEN                         ?= $(TOOLS_DIR)/client-gen
LISTER_GEN                         ?= $(TOOLS_DIR)/lister-gen
INFORMER_GEN                       ?= $(TOOLS_DIR)/informer-gen
OPENAPI_GEN                        ?= $(TOOLS_DIR)/openapi-gen
REGISTER_GEN                       ?= $(TOOLS_DIR)/register-gen
DEEPCOPY_GEN                       ?= $(TOOLS_DIR)/deepcopy-gen
CODE_GEN_VERSION                   ?= v0.28.0

$(CONTROLLER_GEN):
	@echo Install controller-gen... >&2
	@cd ./hack/controller-gen && GOBIN=$(TOOLS_DIR) go install

$(CLIENT_GEN):
	@echo Install client-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/client-gen@$(CODE_GEN_VERSION)

$(LISTER_GEN):
	@echo Install lister-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/lister-gen@$(CODE_GEN_VERSION)

$(INFORMER_GEN):
	@echo Install informer-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/informer-gen@$(CODE_GEN_VERSION)

$(OPENAPI_GEN):
	@echo Install openapi-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/openapi-gen@$(CODE_GEN_VERSION)

$(REGISTER_GEN):
	@echo Install register-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/register-gen@$(CODE_GEN_VERSION)

$(DEEPCOPY_GEN):
	@echo Install deepcopy-gen... >&2
	@GOBIN=$(TOOLS_DIR) go install k8s.io/code-generator/cmd/deepcopy-gen@$(CODE_GEN_VERSION)

$(GEN_CRD_API_REFERENCE_DOCS):
	@echo Install gen-crd-api-reference-docs... >&2
	@GOBIN=$(TOOLS_DIR) go install github.com/ahmetb/gen-crd-api-reference-docs@$(GEN_CRD_API_REFERENCE_DOCS_VERSION)

###########
# CODEGEN #
###########

OUT_PACKAGE                 := $(PACKAGE)/pkg/crd/client/customboard
INPUT_DIRS                  := $(PACKAGE)/pkg/crd/api/customboard/v1alpha1
CLIENT_INPUT_DIRS           := $(PACKAGE)/pkg/crd/api/customboard/v1alpha1
CLIENTSET_PACKAGE           := $(OUT_PACKAGE)/clientset
LISTERS_PACKAGE             := $(OUT_PACKAGE)/listers
INFORMERS_PACKAGE           := $(OUT_PACKAGE)/informers
CRDS_PATH                   := ${PWD}/config/crds

.PHONY: codegen-client-clientset
codegen-client-clientset: $(CLIENT_GEN) ## Generate clientset
	@echo Generate clientset... >&2
	@rm -rf $(CLIENTSET_PACKAGE) && mkdir -p $(CLIENTSET_PACKAGE)
	@$(CLIENT_GEN) \
		--go-header-file ./scripts/boilerplate.go.txt \
		--clientset-name versioned \
		--output-package $(CLIENTSET_PACKAGE) \
		--input-base "" \
		--input $(CLIENT_INPUT_DIRS)

.PHONY: codegen-client-listers
codegen-client-listers: $(LISTER_GEN) ## Generate listers
	@echo Generate listers... >&2
	@rm -rf $(LISTERS_PACKAGE) && mkdir -p $(LISTERS_PACKAGE)
	@$(LISTER_GEN) \
		--go-header-file ./scripts/boilerplate.go.txt \
		--output-package $(LISTERS_PACKAGE) \
		--input-dirs $(CLIENT_INPUT_DIRS)

.PHONY: codegen-client-informers
codegen-client-informers: $(INFORMER_GEN) ## Generate informers
	@$(INFORMER_GEN) \
		--go-header-file ./scripts/boilerplate.go.txt \
		--output-package $(INFORMERS_PACKAGE) \
		--input-dirs $(CLIENT_INPUT_DIRS) \
		--versioned-clientset-package $(CLIENTSET_PACKAGE)/versioned \
		--listers-package $(LISTERS_PACKAGE)

.PHONY: codegen-client-wrappers
codegen-client-wrappers: codegen-client-clientset $(GOIMPORTS) ## Generate client wrappers
	@echo Generate client wrappers... >&2
	@go run ./hack/main.go
	@$(GOIMPORTS) -w ./pkg/clients
	@go fmt ./pkg/clients/...

.PHONY: codegen-register
codegen-register: $(REGISTER_GEN) ## Generate types registrations
	@echo Generate registration... >&2
	@$(REGISTER_GEN) \
		--go-header-file=./scripts/boilerplate.go.txt \
		--input-dirs=$(INPUT_DIRS)

.PHONY: codegen-deepcopy
codegen-deepcopy: $(DEEPCOPY_GEN) ## Generate deep copy functions
	@echo Generate deep copy functions... >&2
	@$(DEEPCOPY_GEN) \
		--go-header-file=./scripts/boilerplate.go.txt \
		--input-dirs=$(INPUT_DIRS) \
		--output-file-base=zz_generated.deepcopy

.PHONY: codegen-crds
codegen-crds: $(CONTROLLER_GEN)
	@echo Generate policy reporter crds... >&2
	@rm -rf $(CRDS_PATH) && mkdir -p $(CRDS_PATH)
	@$(CONTROLLER_GEN) paths=./pkg/crd/api/customboard/... crd:crdVersions=v1,ignoreUnexportedFields=true,generateEmbeddedObjectMeta=false output:dir=$(CRDS_PATH)
