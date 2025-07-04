BINARY=terraform-provider-graylog

default: build

build:
	go install
	@echo "Binary has been compiled to $(shell go env GOPATH)/bin/${BINARY}"
	@echo "In order to have Terraform pick this up you will need to add the following to your $$HOME/.terraformrc file:"
	@echo "  provider_installation {"
	@echo "    dev_overrides {"
	@echo "      \"SanchosPancho/graylog\" = \"$(shell go env GOPATH)/bin\""
	@echo "    }"
	@echo "    direct {}"
	@echo "  }"
	@echo ""
	@echo "This should only be used during development. See https://www.terraform.io/docs/commands/cli-config.html#development-overrides-for-provider-developers for details."

lint:
	golangci-lint run

test:
	go test -v -cover -failfast ./...

testacc:
	TF_ACC=1 go test -v -cover -failfast ./...

clean:
	rm -f "$(shell go env GOPATH)/bin/$(BINARY)"

updatedeps:
	go get -u && go mod tidy
