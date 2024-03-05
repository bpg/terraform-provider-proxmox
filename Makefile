NAME=terraform-provider-proxmox
TARGETS=darwin linux windows
TERRAFORM_PLUGIN_EXTENSION=
VERSION=0.48.1# x-release-please-version

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

.PHONY: build
build:
	mkdir -p "$(TERRAFORM_PLUGIN_OUTPUT_DIRECTORY)"
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
	@TF_ACC=1 env $$(cat testacc.env | xargs) go test -count=1 -v github.com/bpg/terraform-provider-proxmox/fwprovider/tests/...

.PHONY: lint
lint:
	go run -modfile=tools/go.mod github.com/golangci/golangci-lint/cmd/golangci-lint run --fix

.PHONY: release-build
release-build:
	go run -modfile=tools/go.mod github.com/goreleaser/goreleaser build --clean --skip=validate

.PHONY: docs
docs:
	@mkdir -p ./build/docs-gen
	@cd ./tools && go generate tools.go

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
