GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
NAME=$$(grep TerraformProviderName proxmoxtf/version.go | grep -o -e 'terraform-provider-[a-z]*')
TARGETS=darwin linux windows
VERSION=$$(grep TerraformProviderVersion proxmoxtf/version.go | grep -o -e '[0-9]\.[0-9]\.[0-9]')

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
	cd ./example && terraform init

example-plan:
	cd ./example && terraform plan

fmt:
	gofmt -w $(GOFMT_FILES)

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
