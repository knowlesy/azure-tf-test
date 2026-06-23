terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

resource "mockazurerm_resource_group" "rg" {
  name     = "foundation-rg"
  location = "East US"
}

resource "mockazurerm_storage_account" "sa" {
  name                     = "foundationsa123"
  resource_group_name      = mockazurerm_resource_group.rg.name
  location                 = mockazurerm_resource_group.rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
