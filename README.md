# terraform-provider-mockazurerm

A custom offline mock Terraform provider for Azure, built using the `terraform-plugin-framework`.

This provider creates a highly realistic, offline-only validation engine that perfectly mimics the behavior of the official `hashicorp/azurerm` provider. It strictly operates locally and performs NO actual cloud deployments to Azure. Its primary purpose is parsing standard Azure HCL configurations, validating dependency hierarchies, and generating accurate mock state models locally for rapid pipeline CI/CD testing, compliance auditing, or module validation.

## Features

- **Zero Cloud Interaction:** Completely mocks the Azure Resource Manager (ARM) API.
- **Local State Generation:** All CRUD operations generate standard Azure resource ID strings and mock computed attributes directly into a `.mockazurerm_db.json` database.
- **Dependency Integrity:** Ensures parent/child relationships (like Subnets referencing Virtual Networks) are strictly validated in local memory.
- **Schema Parity:** Includes standard major resources like `azurerm_resource_group`, `azurerm_linux_virtual_machine`, `azurerm_kubernetes_cluster`, `azurerm_storage_account`, etc.

## Quick Start

### 1. Build the Provider

Ensure you have Go 1.21+ installed. Clone this repository and run:

```bash
go build -o terraform-provider-mockazurerm
```

### 2. Configure Terraform Dev Overrides

Configure Terraform to use your locally compiled provider by creating or updating a `.terraformrc` file:

```hcl
provider_installation {
  dev_overrides {
    "peterknowles/mockazurerm" = "/path/to/your/repository/azure-tf-test"
  }
  direct {}
}
```

Set the environment variable to point to this configuration:

```bash
export TF_CLI_CONFIG_FILE="/path/to/your/repository/azure-tf-test/test/.terraformrc"
```

### 3. Write HCL and Apply

Create a `main.tf` file using the mock provider:

```hcl
terraform {
  required_providers {
    mockazurerm = {
      source  = "peterknowles/mockazurerm"
    }
  }
}

provider "mockazurerm" {
  subscription_id = "00000000-0000-0000-0000-000000000000"
}

resource "mockazurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}
```

Run Terraform:

```bash
terraform plan
terraform apply
```

Check the generated `.mockazurerm_db.json` file to see your offline mock deployment state!
