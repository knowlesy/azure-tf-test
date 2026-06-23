# EXPECTED ERROR: The resource block expects mockazurerm_virtual_network, but received mockazurerm_virtual_netwrk.
terraform { required_providers { mockazurerm = { source = "peterknowles/mockazurerm" } } }
provider "mockazurerm" {}

resource "mockazurerm_virtual_netwrk" "vnet" {
  name                = "typo-vnet"
  resource_group_name = "example-rg"
  location            = "East US"
  address_space       = ["10.0.0.0/16"]
}
