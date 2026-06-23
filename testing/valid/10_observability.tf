terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

resource "mockazurerm_resource_group" "rg" {
  name     = "obs-rg"
  location = "East US"
}

resource "mockazurerm_log_analytics_workspace" "law" {
  name                = "core-law"
  location            = mockazurerm_resource_group.rg.location
  resource_group_name = mockazurerm_resource_group.rg.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}
