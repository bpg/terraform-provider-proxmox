GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
NAME=terraform-provider-proxmox
TARGETS=darwin linux windows
TERRAFORM_PLUGIN_EXTENSION=
VERSION=0.14.0# x-release-please-version
VERSION_EXAMPLE=9999.0.0

ifeq ($(OS),Windows_NT)
	TERRAFORM_PLATFORM=windows_amd64
	TERRAFORM_PLUGIN_CACHE_DIRECTORY=$$(cygpath -u "$(shell pwd -P)")/cache/plugins
	TERRAFORM_PLUGIN_EXTENSION=.exe
else
	TERRAFORM_PLATFORM=$$(terraform -version | awk 'FNR == 2 {print $$2}')
	TERRAFORM_PLUGIN_CACHE_DIRECTORY=$(shell pwd -P)/cache/plugins
endif

TERRAFORM_PLUGIN_DIRECTORY=$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)/registry.terraform.io/bpg/proxmox/$(VERSION)/$(TERRAFORM_PLATFORM)
TERRAFORM_PLUGIN_DIRECTORY_EXAMPLE=$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)/registry.terraform.io/bpg/proxmox/$(VERSION_EXAMPLE)/$(TERRAFORM_PLATFORM)
TERRAFORM_PLUGIN_EXECUTABLE=$(TERRAFORM_PLUGIN_DIRECTORY)/$(NAME)_v$(VERSION)_x4$(TERRAFORM_PLUGIN_EXTENSION)
TERRAFORM_PLUGIN_EXECUTABLE_EXAMPLE=$(TERRAFORM_PLUGIN_DIRECTORY_EXAMPLE)/$(NAME)_v$(VERSION_EXAMPLE)_x4$(TERRAFORM_PLUGIN_EXTENSION)

default: build

build:
	mkdir -p "$(TERRAFORM_PLUGIN_DIRECTORY)"
	rm -f "$(TERRAFORM_PLUGIN_EXECUTABLE)"
	go build -o "$(TERRAFORM_PLUGIN_EXECUTABLE)"

example: example-build example-init example-apply example-destroy

example-apply:
	export TF_CLI_CONFIG_FILE="$(shell pwd -P)/example.tfrc" \
		&& export TF_DISABLE_CHECKPOINT="true" \
		&& export TF_PLUGIN_CACHE_DIR="$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)" \
		&& cd ./example \
		&& terraform apply -auto-approve

example-build:
	rm -rf "$(TERRAFORM_PLUGIN_DIRECTORY_EXAMPLE)"
	mkdir -p "$(TERRAFORM_PLUGIN_DIRECTORY_EXAMPLE)"
	go build -o "$(TERRAFORM_PLUGIN_EXECUTABLE_EXAMPLE)"

example-destroy:
	export TF_CLI_CONFIG_FILE="$(shell pwd -P)/example.tfrc" \
		&& export TF_DISABLE_CHECKPOINT="true" \
		&& export TF_PLUGIN_CACHE_DIR="$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)" \
		&& cd ./example \
		&& terraform destroy -auto-approve

example-init:
	export TF_CLI_CONFIG_FILE="$(shell pwd -P)/example.tfrc" \
		&& export TF_DISABLE_CHECKPOINT="true" \
		&& export TF_PLUGIN_CACHE_DIR="$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)" \
		&& cd ./example \
		&& rm -f .terraform.lock.hcl \
		&& terraform init

example-plan:
	export TF_CLI_CONFIG_FILE="$(shell pwd -P)/example.tfrc" \
		&& export TF_DISABLE_CHECKPOINT="true" \
		&& export TF_PLUGIN_CACHE_DIR="$(TERRAFORM_PLUGIN_CACHE_DIRECTORY)" \
		&& cd ./example \
		&& terraform plan

fmt:
	gofmt -s -w $(GOFMT_FILES)

init:
	go get ./...

targets: $(TARGETS)

test:
	go test -v ./...

$(TARGETS):
	GOOS=$@ GOARCH=amd64 CGO_ENABLED=0 go build \
		-o "dist/$@/$(NAME)_v$(VERSION)-custom_x4" \
		-a -ldflags '-extldflags "-static"'
	zip \
		-j "dist/$(NAME)_v$(VERSION)-custom_$@_amd64.zip" \
		"dist/$@/$(NAME)_v$(VERSION)-custom_x4"

.PHONY: build example example-apply example-destroy example-init example-plan fmt init targets test $(TARGETS)
