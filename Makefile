# regular expressions for excluded file patterns
EXCLUDE_COVERAGE_FILES="(/examples/)|(/\.tools/)"


# go-install-tool will 'go install' any package $2 and install it locally to $1.
# This will prevent that they are installed in the $USER/go/bin folder and different
# projects ca have different versions of the tools
PROJECT_DIR := $(shell dirname $(abspath $(firstword $(MAKEFILE_LIST))))

TOOLS_DIR ?= $(PROJECT_DIR)/.tools

# $(1) command name
# $(2) repo URL
# $(3) version
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e ;\
TMP_DIR=$$(mktemp -d) ;\
cd $$TMP_DIR ;\
go mod init tmp ;\
echo "Removing any outdated version of $(1)";\
rm -f $(1)*;\
echo "Downloading $(2)@$(3)" ;\
GOBIN=$(TOOLS_DIR) GOFLAGS="-mod=mod" go install "$(2)@$(3)" ;\
touch "$(1)-$(3)";\
rm -rf $$TMP_DIR ;\
}
endef

# prereqs binary dependencies
GOLANGCI_LINT = $(TOOLS_DIR)/golangci-lint

.PHONY: prereqs
prereqs:
	@echo "### Check if prerequisites are met, and installing missing dependencies"
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/v2/cmd/golangci-lint,v2.3.1)

.PHONY: fmt
fmt: prereqs
	@echo "### Formatting code and fixing imports"
	$(GOLANGCI_LINT) fmt

.PHONY: lint
lint: prereqs
	@echo "### Linting code"
	$(GOLANGCI_LINT) run ./... --timeout=6m

.PHONY: test
test:
	@echo "### Testing code"
	go test -race -mod vendor -a $$(go list ./... | grep -v /examples) -coverpkg=./... -coverprofile coverage.out

.PHONY: coverage-report
coverage-report:
	@echo "### Generating coverage report"
	go tool cover --func=coverage.out

.PHONY: coverage-report-html
coverage-report-html:
	@echo "### Generating HTML coverage report"
	go tool cover --html=coverage.out
