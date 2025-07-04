# ATTENTION
This repository was forked from https://github.com/zahiar/terraform-provider-graylog
First of all - thank's to @zahiar for his hard work on this project.

This fork will be adapted for version 5+ when time permits.
By using this provider, you accept full responsibility for any consequences that may occur.

# Terraform Provider: Graylog
This is a Terraform provider for managing resources within [Graylog](https://docs.graylog.org/).

## Getting Started
As this provider is published to the public [Terraform Registry](https://registry.terraform.io/providers/SanchosPancho/graylog),
you can install it like so (for Terraform 0.14+):
```hcl
provider "graylog" {
  web_endpoint_uri = "http://example.com/api"
  api_version      = "v4"
}

terraform {
  required_providers {
    graylog = {
      source  = "SanchosPancho/graylog"
    }
  }
}
```

For more detailed instructions and documentation on the resources and data sources supported, please go to
[Terraform Registry](https://registry.terraform.io/providers/SanchosPancho/graylog/latest/docs).

## Maintenance
This provider is maintained during free time, so if you are interested in helping to develop this further, you
are more than welcome to submit a pull request or raise a ticket if you'd prefer.

## Development

### Requirements
If you do wish to help develop this, you will need the following installed:
* [Go](http://www.golang.org) (see `go.mod` file for the correct version to install)
* [Go Linter](https://formulae.brew.sh/formula/golangci-lint)
* [GOPATH](http://golang.org/doc/code.html#GOPATH) (is correctly setup)
* [Terraform](https://www.terraform.io/downloads.html) (0.14+)

### Building
Simply run `make build`, and it will compile and create a binary, as well as print-out instructions
on how to configure Terraform to use this locally built provider.
```shell
$ make build
```

### Testing

#### Unit Tests
```shell
$ make test
```

### Acceptance Tests
```shell
$ make testacc
```

### Documentation
Every data source or resource added must have an accompanying docs page (see `docs` directory for examples).

Docs are written using Markdown, and you can use [this page](https://registry.terraform.io/tools/doc-preview) to preview what your docs will look like when rendered.
