# EXPECTED ERROR: Storage account names must be between 3 and 24 characters in length and may contain numbers and lowercase letters only.
terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

resource "mockazurerm_storage_account" "sa" {
  name                     = "Invalid-Name-With-Caps!" # Fails schema validation (simulated logic)
  resource_group_name      = "example-rg"
  location                 = "East US"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
