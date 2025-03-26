NAME=terraform-provider-proxmox
TARGETS=darwin linux windows
TERRAFORM_PLUGIN_EXTENSION=
VERSION=0.73.2# x-release-please-version

GOLANGCI_LINT_VERSION=v1.64.8# renovate: depName=golangci/golangci-lint datasource=github-releases

# check if opentofu is installed and use it if it is,
# otherwise use terraform
ifeq ($(shell tofu -version 2>/dev/null),)
	TERRAFORM_EXECUTABLE=terraform
else
	TERRAFORM_EXECUTABLE=tofu
endif

ifeq ($(OS),Windows_NT)
	TERRAFORM_PLATFORM=windows_amd64
	TERRAFORM_PLUGIN_CACHE_DIRECTORY=$$(cygpath -u "$(shell pwd -P)")/cache/plugins
	TERRAFORM_PLUGIN_EXTENSION=.exe
else
	TERRAFORM_PLATFORM=$$($(TERRAFORM_EXECUTABLE) -version | awk 'FNR == 2 {print $$2}')
	TERRAFORM_PLUGIN_CACHE_DIRECTORY=$(shell pwd -P)/cache/plugins
endif

TERRAFORM_PLUGIN_OUTPUT_DIRECTORY=./build
TERRAFORM_PLUGIN_EXECUTABLE=$(TERRAFORM_PLUGIN_OUTPUT_DIRECTORY)/$(NAME)_v$(VERSION)$(TERRAFORM_PLUGIN_EXTENSION)
TERRAFORM_PLUGIN_EXECUTABLE_EXAMPLE=$(TERRAFORM_PLUGIN_OUTPUT_DIRECTORY)/$(NAME)$(TERRAFORM_PLUGIN_EXTENSION)

default: build

.PHONY: clean
clean:
	rm -rf ./dist
	rm -rf ./cache
	rm -rf ./build
	rm -f ./example/test-api-creds-auth/outs_cred-tester__expect_*

.PHONY: build
build:
	mkdir -p "$(TERRAFORM_PLUGIN_OUTPUT_DIRECTORY)" "$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)"
	rm -f "$(TERRAFORM_PLUGIN_EXECUTABLE)"
	go build -o "$(TERRAFORM_PLUGIN_EXECUTABLE)"

.PHONY: example
example: example-build example-init example-apply example-destroy

.PHONY: example-apply
example-apply:
	export TF_CLI_CONFIG_FILE="$(shell pwd -P)/example.tfrc" \
		&& export TF_DISABLE_CHECKPOINT="true" \
		&& export TF_PLUGIN_CACHE_DIR="$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)" \
		&& cd ./example \
		&& $(TERRAFORM_EXECUTABLE) apply -auto-approve

.PHONY: example-build
example-build:
	mkdir -p "$(TERRAFORM_PLUGIN_OUTPUT_DIRECTORY)"
	rm -rf "$(TERRAFORM_PLUGIN_EXECUTABLE_EXAMPLE)"
	go build -o "$(TERRAFORM_PLUGIN_EXECUTABLE_EXAMPLE)"

.PHONY: example-destroy
example-destroy:
	export TF_CLI_CONFIG_FILE="$(shell pwd -P)/example.tfrc" \
		&& export TF_DISABLE_CHECKPOINT="true" \
		&& export TF_PLUGIN_CACHE_DIR="$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)" \
		&& cd ./example \
		&& $(TERRAFORM_EXECUTABLE) destroy -auto-approve

.PHONY: example-init
example-init:
	export TF_CLI_CONFIG_FILE="$(shell pwd -P)/example.tfrc" \
		&& export TF_DISABLE_CHECKPOINT="true" \
		&& export TF_PLUGIN_CACHE_DIR="$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)" \
		&& cd ./example \
		&& rm -f .terraform.lock.hcl \
		&& $(TERRAFORM_EXECUTABLE) init

.PHONY: example-plan
example-plan:
	export TF_CLI_CONFIG_FILE="$(shell pwd -P)/example.tfrc" \
		&& export TF_DISABLE_CHECKPOINT="true" \
		&& export TF_PLUGIN_CACHE_DIR="$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)" \
		&& cd ./example \
		&& $(TERRAFORM_EXECUTABLE) plan

.PHONY: test-api-creds-auth
test-api-creds-auth:
	rm -f ./example/test-api-creds-auth/outs_cred-tester__expect_*
	export TF_CLI_CONFIG_FILE="$(shell pwd -P)/example/test-api-creds-auth/example.tfrc" \
		&& export TF_DISABLE_CHECKPOINT="true" \
		&& export TF_PLUGIN_CACHE_DIR="$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)" \
		&& cd ./example/test-api-creds-auth \
		&& ./api-creds-tester.sh $(TERRAFORM_EXECUTABLE)

.PHONY: fmt
fmt:
	gofmt -s -w $$(find . -name '*.go')

.PHONY: init
init:
	go get ./...

.PHONY: test
test:
	go test ./...

.PHONY: testacc
testacc:
	@# explicitly add TF_ACC=1 to trigger the acceptance tests, `testacc.env` might be missing or incomplete
	@TF_ACC=1 env $$(cat testacc.env | xargs) go test --timeout=30m --tags=acceptance -count=1 -v github.com/bpg/terraform-provider-proxmox/fwprovider/...

.PHONY: lint
lint:
	# NOTE: This target runs only locally, not in CI. See .github/workflows/golangci-lint.yml for CI linting.
	@docker run -t --rm -v $$(pwd):/app -v ~/.cache/golangci-lint/$(GOLANGCI_LINT_VERSION):/root/.cache -w /app golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) golangci-lint run --fix

.PHONY: release-build
release-build:
	goreleaser build --clean --skip=validate

.PHONY: docs
docs:
	@mkdir -p ./build/docs-gen
	@go generate main.go

.PHONY: targets
targets: $(TARGETS)

.PHONY: $(TARGETS)
$(TARGETS):
	GOOS=$@ GOARCH=amd64 CGO_ENABLED=0 go build \
		-o "dist/$@/$(NAME)_v$(VERSION)-custom" \
		-a -ldflags '-extldflags "-static"'
	zip \
		-j "dist/$(NAME)_v$(VERSION)-custom_$@_amd64.zip" \
		"dist/$@/$(NAME)_v$(VERSION)-custom"
