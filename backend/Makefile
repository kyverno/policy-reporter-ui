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
CONTROLLER_GEN_VERSION             ?= v0.18.0
CLIENT_GEN                         ?= $(TOOLS_DIR)/client-gen
LISTER_GEN                         ?= $(TOOLS_DIR)/lister-gen
INFORMER_GEN                       ?= $(TOOLS_DIR)/informer-gen
OPENAPI_GEN                        ?= $(TOOLS_DIR)/openapi-gen
REGISTER_GEN                       ?= $(TOOLS_DIR)/register-gen
DEEPCOPY_GEN                       ?= $(TOOLS_DIR)/deepcopy-gen
CODE_GEN_VERSION                   ?= v0.33.4

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
CLIENT_PACKAGE              := $(PACKAGE)/pkg/crd/client
CLIENTSET_PACKAGE           := $(CLIENT_PACKAGE)/clientset
LISTERS_PACKAGE             := $(CLIENT_PACKAGE)/listers
INFORMERS_PACKAGE           := $(CLIENT_PACKAGE)/informers
CRDS_PATH                   := ${PWD}/config/crds

.PHONY: codegen-client-clientset
codegen-client-clientset: $(CLIENT_GEN)
	@echo Generate clientset... >&2
	@rm -rf ./pkg/crd/client/clientset && mkdir -p ./pkg/crd/client/clientset
	@$(CLIENT_GEN) \
		--go-header-file ./scripts/boilerplate.go.txt \
		--clientset-name versioned \
		--output-dir ./pkg/crd/client/clientset \
		--output-pkg $(CLIENTSET_PACKAGE) \
		--input-base github.com/kyverno/policy-reporter-ui \
		--input ./pkg/crd/api/customboard/v1alpha1

.PHONY: codegen-client-listers
codegen-client-listers: ## Generate listers
codegen-client-listers: $(LISTER_GEN)
	@echo Generate listers... >&2
	@rm -rf ./pkg/crd/client/listers && mkdir -p ./pkg/crd/client/listers
	@$(LISTER_GEN) \
		--go-header-file ./scripts/boilerplate.go.txt \
		--output-dir ./pkg/crd/client/listers \
		--output-pkg $(LISTERS_PACKAGE) \
		./pkg/crd/api/customboard/v1alpha1

.PHONY: codegen-client-informers
codegen-client-informers: $(INFORMER_GEN)
	@echo Generate informers... >&2
	@rm -rf ./pkg/crd/client/informers && mkdir -p ./pkg/crd/client/informers
	@$(INFORMER_GEN) \
		--go-header-file ./scripts/boilerplate.go.txt \
		--output-dir ./pkg/crd/client/informers \
		--output-pkg $(INFORMERS_PACKAGE) \
		--versioned-clientset-package $(CLIENTSET_PACKAGE)/versioned \
		--listers-package $(LISTERS_PACKAGE) \
		./pkg/crd/api/customboard/v1alpha1

.PHONY: codegen-register
codegen-register: $(REGISTER_GEN) ## Generate types registrations
	@echo Generate registration... >&2
	@$(REGISTER_GEN) --go-header-file=./scripts/boilerplate.go.txt --output-file zz_generated.register.go ./pkg/crd/api/...

.PHONY: codegen-deepcopy
codegen-deepcopy: $(DEEPCOPY_GEN) ## Generate deep copy functions
	@echo Generate deep copy functions... >&2
	@$(DEEPCOPY_GEN) --go-header-file ./scripts/boilerplate.go.txt --output-file zz_generated.deepcopy.go ./pkg/crd/api/...

.PHONY: codegen-crds
codegen-crds: $(CONTROLLER_GEN)
	@echo Generate crds... >&2
	@rm -rf $(CRDS_PATH)/customboard && mkdir -p $(CRDS_PATH)/customboard
	@$(CONTROLLER_GEN) \
		paths=./pkg/crd/api/customboard/v1alpha1/... \
		crd:crdVersions=v1,ignoreUnexportedFields=true,generateEmbeddedObjectMeta=false \
		output:dir=$(CRDS_PATH)/customboard
