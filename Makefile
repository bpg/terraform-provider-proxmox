GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
NAME=$$(grep TerraformProviderName proxmoxtf/version.go | grep -o -e 'terraform-provider-[a-z]*')
TARGETS=darwin linux windows
TERRAFORM_PLUGIN_EXTENSION=
VERSION=$$(grep TerraformProviderVersion proxmoxtf/version.go | grep -o -e '[0-9]\.[0-9]\.[0-9]')

ifeq ($(OS),Windows_NT)
	TERRAFORM_CACHE_DIRECTORY=$$(cygpath -u "$(APPDATA)")/terraform.d/plugins
	TERRAFORM_PLATFORM=windows_amd64
	TERRAFORM_PLUGIN_EXTENSION=.exe
else
	TERRAFORM_CACHE_DIRECTORY=$(HOME)/terraform.d/plugins
	UNAME_S=$$(shell uname -s)

	ifeq ($(UNAME_S),Darwin)
		TERRAFORM_PLATFORM=darwin_amd64
	else
		TERRAFORM_PLATFORM=linux_amd64
	endif
endif

TERRAFORM_PLUGIN_DIRECTORY=$(TERRAFORM_CACHE_DIRECTORY)/terraform.danitso.com/provider/proxmox/$(VERSION)/$(TERRAFORM_PLATFORM)
TERRAFORM_PLUGIN_EXECUTABLE=$(TERRAFORM_PLUGIN_DIRECTORY)/$(NAME)_v$(VERSION)_x4$(TERRAFORM_PLUGIN_EXTENSION)

default: build

build:
	go build -o "bin/$(NAME)_v$(VERSION)-custom_x4"

example: example-init example-apply example-apply example-destroy

example-apply:
	cd ./example && terraform apply -auto-approve

example-destroy:
	cd ./example && terraform destroy -auto-approve

example-init:
	rm -f "example/$(NAME)_v"*
	go build -o "example/$(NAME)_v$(VERSION)-custom_x4"

	mkdir -p "$(TERRAFORM_PLUGIN_DIRECTORY)"
	rm -f "$(TERRAFORM_PLUGIN_EXECUTABLE)"
	cp "example/$(NAME)_v$(VERSION)-custom_x4" "$(TERRAFORM_PLUGIN_EXECUTABLE)"

	cd ./example && terraform init

example-plan:
	cd ./example && terraform plan

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
