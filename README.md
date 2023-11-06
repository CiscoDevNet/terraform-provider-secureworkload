# Cisco Secure Workload (SecureWorkload) Terraform Provider

> **Note:** this Terraform provider is now publically available on the [Terraform Registry](https://registry.terraform.io/providers/CiscoDevNet/secureworkload/latest).
 
Terraform Provider for managing Cisco Secure Workload (SecureWorkload) resources.

## Usage

### Using the Terraform Registry

Create a `main.tf` file with the following content, save, and run `terraform init` from a terminal window in the same directory as `main.tf`:

```hcl
terraform {
  required_providers {
    secureworkload = {
      source = "CiscoDevNet/secureworkload"
      version = "0.1.0"
    }
  }
}

provider "secureworkload" {
  api_key                  = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_secret               = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_url                  = "https://tenant.secureworkloadpreview.com"
  disable_tls_verification = false
}

resource "secureworkload_filter" "filter" {
  name            = "Terraform created filter"
  query_type      = "eq"
  query_field     = "ip"
  query_value     = "10.0.0.1"
  app_scope_id = "5ceea87b497d4f753baf85bc"
}
```

### Building and Consuming

1. Build the plugin

```bash
make build
```

2. Copy the plugin to your terraform plugin directory, e.g.

```
mkdir ~/.terraform.d/plugins/darwin_amd64
cp terraform-provider-secureworkload ~/.terraform.d/plugins/darwin_amd64
```

3.Add plugin to terraform for the current module you are working on

```bash
cd /path/to/terraform/module
terraform init -plugin-dir ~/.terraform.d/plugins/darwin_amd64
```

4.Write terraform code using this provider.

```hcl
provider "secureworkload" {
  api_key                  = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_secret               = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
  api_url                  = "https://tenant.secureworkloadpreview.com"
  disable_tls_verification = false
}

resource "secureworkload_filter" "filter" {
  name            = "Terraform created filter"
  query_type      = "eq"
  query_field     = "ip"
  query_value     = "10.0.0.1"
  app_scope_id = "5ceea87b497d4f753baf85bc"
}
```

More [example terraform modules for managing secureworkload resources with this provider.](./examples)

## Development

### Testing

Tests can be executed via

```bash
make test
```

When the test process is running any variable set in a top level `.env` file in this project will be available to the tests as an environment variable.

Example `.env` file

```bash
VARIABLE=value
```

This file is gitignored to prevent any sensitive material such as api keys from being published.

## Publishing

To build binaries for mac, linux(amd64), windows(x86), run

```bash
make cross-compile
```

The built binaries will be placed in the [bin directory](./bin).
